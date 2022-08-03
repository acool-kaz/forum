package delivery

import (
	"forum/models"
	"net/http"
	"strings"
)

func (h *Handler) userPage(w http.ResponseWriter, r *http.Request) {
	user := h.userIdentity(w, r)
	if user == (models.User{}) {
		h.errorPage(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}
	if r.Method != http.MethodGet {
		h.errorPage(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	posts, err := h.Services.GetPostByUsername(user.Username)
	if err != nil {
		h.errorPage(w, http.StatusInternalServerError, err.Error())
		return
	}
	info := models.Info{
		User:  user,
		Posts: posts,
	}
	if err := h.Tmpl.ExecuteTemplate(w, "user.html", info); err != nil {
		h.errorPage(w, http.StatusInternalServerError, err.Error())
	}
}

func (h *Handler) anotherUserPage(w http.ResponseWriter, r *http.Request) {
	userPage := strings.TrimPrefix(r.URL.Path, "/profile/")
	user := h.userIdentity(w, r)
	if r.Method != http.MethodGet {
		h.errorPage(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	posts, err := h.Services.GetPostByUsername(userPage)
	if err != nil {
		h.errorPage(w, http.StatusInternalServerError, err.Error())
		return
	}
	info := models.Info{
		User:  user,
		Posts: posts,
	}
	if err := h.Tmpl.ExecuteTemplate(w, "user.html", info); err != nil {
		h.errorPage(w, http.StatusInternalServerError, err.Error())
	}
}
