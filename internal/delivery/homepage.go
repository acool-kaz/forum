package delivery

import (
	"forum/models"
	"net/http"
)

func (h *Handler) homePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		h.errorPage(w, r, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	if r.Method != http.MethodGet {
		h.errorPage(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	userId := r.Context().Value(userId).(uint)
	user, err := h.services.User.GetById(r.Context(), userId)
	if err != nil {
		h.errorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	posts, err := h.services.Post.GetAll(r.Context())
	if err != nil {
		h.errorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	info := models.Info{
		User:  user,
		Posts: posts,
	}
	if err := h.tmpl.ExecuteTemplate(w, "index.html", info); err != nil {
		h.errorPage(w, r, http.StatusInternalServerError, err.Error())
	}
}
