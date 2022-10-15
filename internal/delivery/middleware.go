package delivery

import (
	"forum/models"
	"net/http"
	"time"
)

func (h *Handler) userIdentity(w http.ResponseWriter, r *http.Request) models.User {
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			return models.User{}
		}
		h.errorPage(w, http.StatusBadRequest, err.Error())
		return models.User{}
	}
	user, err := h.services.ParseSessionToken(c.Value)
	if err != nil {
		return models.User{}
	}
	if user.ExpiresAt.Before(time.Now()) {
		if err := h.services.DeleteSessionToken(c.Value); err != nil {
			h.errorPage(w, http.StatusInternalServerError, err.Error())
			return models.User{}
		}
		return models.User{}
	}
	return user
}

// const userCtx = "user"

// func (h *Handler) userIdentity(next http.HandlerFunc) http.HandlerFunc {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		c, err := r.Cookie("session_token")
// 		if err != nil {
// 			h.errorPage(w, http.StatusUnauthorized, err.Error())
// 			return
// 		}
// 		user, err := h.services.ParseSessionToken(c.Value)
// 		if err != nil {
// 			h.errorPage(w, http.StatusUnauthorized, err.Error())
// 			return
// 		}
// 		if user.ExpiresAt.Before(time.Now()) {
// 			if err := h.services.DeleteSessionToken(c.Value); err != nil {
// 				h.errorPage(w, http.StatusInternalServerError, err.Error())
// 				return
// 			}
// 			return
// 		}
// 		h.ctx = context.WithValue(h.ctx, userCtx, user)
// 		next.ServeHTTP(w, r)
// 	})
// }
