package delivery

import (
	"forum/models"
	"net/http"
	"strings"
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
	var posts []models.Post
	var err error
	if len(r.URL.Query()) == 0 {
		posts, err = h.Services.GetAllPost()
		if err != nil {
			h.errorPage(w, http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		for key, val := range r.URL.Query() {
			switch key {
			case "category":
				posts, err = h.Services.GetPostsByCategory(strings.Join(val, ""))
				if err != nil {
					h.errorPage(w, http.StatusInternalServerError, err.Error())
					return
				}
			case "time":
				if user == (models.User{}) {
					h.errorPage(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
					return
				}
				posts, err = h.Services.GetPostsByTime(strings.Join(val, ""))
				if err != nil {
					h.errorPage(w, http.StatusInternalServerError, err.Error())
					return
				}
			case "likes":
				if user == (models.User{}) {
					h.errorPage(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
					return
				}
				posts, err = h.Services.GetPostsByLike(strings.Join(val, ""))
				if err != nil {
					h.errorPage(w, http.StatusInternalServerError, err.Error())
					return
				}
			default:
				h.errorPage(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
				return
			}
		}
	}
	d := models.Info{
		Posts: posts,
		User:  user,
	}
	if err := h.Tmpl.ExecuteTemplate(w, "index.html", d); err != nil {
		h.errorPage(w, http.StatusInternalServerError, err.Error())
	}
}
