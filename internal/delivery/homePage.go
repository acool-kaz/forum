package delivery

import (
	"errors"
	"forum/internal/service"
	"forum/models"
	"net/http"
)

func (h *Handler) homePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		h.errorPage(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	if r.Method != http.MethodGet {
		h.errorPage(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	user := h.userIdentity(w, r)
	var (
		posts         []models.Post
		notifications []models.Notification
		err           error
	)
	if user != (models.User{}) {
		notifications, err = h.Services.GetAllNotificationForUser(user.Username)
		if err != nil {
			h.errorPage(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
	if len(r.URL.Query()) == 0 {
		posts, err = h.Services.GetAllPost()
		if err != nil {
			h.errorPage(w, http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		posts, err = h.Services.GetAllPostBy(user, r.URL.Query())
		if err != nil {
			if errors.Is(err, service.ErrUserNotFound) {
				h.errorPage(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
				return
			}
			h.errorPage(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
	info := models.Info{
		Posts:         posts,
		User:          user,
		Notifications: notifications,
	}
	if err := h.Tmpl.ExecuteTemplate(w, "index.html", info); err != nil {
		h.errorPage(w, http.StatusInternalServerError, err.Error())
	}
}
