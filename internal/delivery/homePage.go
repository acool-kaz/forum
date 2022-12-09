package delivery

import (
	"forum/models"
	"net/http"
)

func (h *Handler) homePage(w http.ResponseWriter, r *http.Request) {
	var (
		posts         []models.Post
		notifications []models.Notification
		err           error
	)
	if r.URL.Path != "/" {
		h.errorPage(w, r, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	if r.Method != http.MethodGet {
		h.errorPage(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	user := h.ctx.Value(uCtx).(models.User)
	if len(r.URL.Query()) == 0 {
		posts, err = h.services.GetAllPost()
		if err != nil {
			h.errorPage(w, r, http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		posts, err = h.services.GetAllPostBy(user, r.URL.Query())
		if err != nil {
			h.errorPage(w, r, http.StatusInternalServerError, err.Error())
			return
		}
	}
	notifications, err = h.services.GetAllNotificationForUser(user)
	if err != nil {
		h.errorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	info := models.Info{
		Posts:         posts,
		User:          user,
		Notifications: notifications,
	}
	if err := h.tmpl.ExecuteTemplate(w, "index.html", info); err != nil {
		h.errorPage(w, r, http.StatusInternalServerError, err.Error())
	}
}
