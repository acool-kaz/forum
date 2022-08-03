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
	user, err := h.Services.ParseSessionToken(c.Value)
	if err != nil {
		return models.User{}
	}
	if user.ExpiresAt.Before(time.Now()) {
		if err := h.Services.DeleteSessionToken(c.Value); err != nil {
			h.errorPage(w, http.StatusInternalServerError, err.Error())
			return models.User{}
		}
		return models.User{}
	}
	return user
}
