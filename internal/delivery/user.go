package delivery

import (
	"errors"
	"forum/internal/service"
	"forum/models"
	"net/http"
	"strings"
)

func (h *Handler) userProfilePage(w http.ResponseWriter, r *http.Request) {
	user := h.userIdentity(w, r)
	username := strings.TrimPrefix(r.URL.Path, "/profile/")
	userPage, err := h.Services.GetUserByUsername(username)
	if err != nil {
		h.errorPage(w, http.StatusNotFound, err.Error())
		return
	}
	if r.Method != http.MethodGet {
		h.errorPage(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	posts, err := h.Services.GetPostByUsername(userPage.Username, r.URL.Query())
	if err != nil {
		if errors.Is(err, service.ErrInvalidQuery) {
			h.errorPage(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			return
		}
		h.errorPage(w, http.StatusInternalServerError, err.Error())
		return
	}
	notifications, err := h.Services.GetAllNotificationForUser(user.Username)
	if err != nil {
		h.errorPage(w, http.StatusInternalServerError, err.Error())
		return
	}
	info := models.Info{
		User:          user,
		ProfileUser:   userPage,
		Posts:         posts,
		Notifications: notifications,
	}
	if err := h.Tmpl.ExecuteTemplate(w, "user.html", info); err != nil {
		h.errorPage(w, http.StatusInternalServerError, err.Error())
	}
}
