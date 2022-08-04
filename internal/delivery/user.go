package delivery

import (
	"forum/models"
	"net/http"
	"strings"
)

func (h *Handler) anotherUserPage(w http.ResponseWriter, r *http.Request) {
	user := h.userIdentity(w, r)
	username := strings.TrimPrefix(r.URL.Path, "/profile/")
	currentUser, err := h.Services.GetUserByUsername(username)
	if err != nil {
		h.errorPage(w, http.StatusNotFound, err.Error())
		return
	}
	if r.Method != http.MethodGet {
		h.errorPage(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	posts, err := h.Services.GetPostByUsername(currentUser.Username)
	if err != nil {
		h.errorPage(w, http.StatusInternalServerError, err.Error())
		return
	}
	info := models.Info{
		User:        user,
		ProfileUser: currentUser,
		Posts:       posts,
	}
	if err := h.Tmpl.ExecuteTemplate(w, "user.html", info); err != nil {
		h.errorPage(w, http.StatusInternalServerError, err.Error())
	}
}