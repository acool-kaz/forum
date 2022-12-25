package delivery

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type middleware struct {
	allMiddlewares []func(http.HandlerFunc) http.HandlerFunc
}

func (m *middleware) addMiddleware(middle func(http.HandlerFunc) http.HandlerFunc) {
	m.allMiddlewares = append(m.allMiddlewares, middle)
}

func (m *middleware) use(router http.HandlerFunc) http.HandlerFunc {
	for _, middleware := range m.allMiddlewares {
		router = middleware(router)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		router.ServeHTTP(w, r)
	}
}

type userCtx string

const userId userCtx = "user_id"

func (h *Handler) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("session_token")
		if err != nil {
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), userId, uint(0))))
			return
		}
		session, err := h.services.ParseSessionToken(r.Context(), c.Value)
		if err != nil {
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), userId, uint(0))))
			return
		}
		if session.ExpiresAt.Before(time.Now()) {
			if err := h.services.DeleteSessionToken(r.Context(), c.Value); err != nil {
				h.errorPage(w, r, http.StatusInternalServerError, err.Error())
				return
			}
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), userId, uint(0))))
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), userId, session.UserId)))
	}
}

func (h *Handler) loggingMiddleware(router http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("\n%s %s [%s]\t%s%s - 200 - OK", time.Now().Format("2006/01/02 15:04:05"), r.Proto, r.Method, r.Host, r.RequestURI)
		router.ServeHTTP(w, r)
	}
}
