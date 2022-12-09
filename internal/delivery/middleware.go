package delivery

import (
	"context"
	"forum/models"
	"log"
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

const uCtx userCtx = "userCtx"

func (h *Handler) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("session_token")
		h.ctx = context.WithValue(h.ctx, uCtx, models.User{})
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		user, err := h.services.ParseSessionToken(c.Value)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		if user.ExpiresAt.Before(time.Now()) {
			if err := h.services.DeleteSessionToken(c.Value); err != nil {
				h.errorPage(w, r, http.StatusInternalServerError, err.Error())
				return
			}
			next.ServeHTTP(w, r)
			return
		}
		h.ctx = context.WithValue(h.ctx, uCtx, user)
		next.ServeHTTP(w, r)
	}
}

func (h *Handler) loggingMiddleware(router http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s\t[%s]\t%s%s", r.Proto, r.Method, r.Host, r.URL.Path)
		router.ServeHTTP(w, r)
	}
}
