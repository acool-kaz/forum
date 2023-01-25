package delivery

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/internal/models"
	"net/http"
	"strconv"
	"strings"
)

func (h *Handler) createPost(w http.ResponseWriter, r *http.Request) {
	var err error
	userId := r.Context().Value(userId).(uint)
	if userId == 0 {
		h.errorPage(w, r, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}
	if r.URL.Path != "/post/create" {
		h.errorPage(w, r, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	switch r.Method {
	case http.MethodGet:
		if err = h.tmpl.ExecuteTemplate(w, "createPost.html", nil); err != nil {
			h.errorPage(w, r, http.StatusInternalServerError, err.Error())
		}
	case http.MethodPost:
		if err = r.ParseMultipartForm(10 << 20); err != nil {
			h.errorPage(w, r, http.StatusInternalServerError, err.Error())
			return
		}
		tags, ok := r.Form["tags"]
		if !ok {
			h.errorPage(w, r, http.StatusBadRequest, "category field not found")
			return
		}
		title, ok := r.Form["title"]
		if !ok {
			h.errorPage(w, r, http.StatusBadRequest, "title field not found")
			return
		}
		description, ok := r.Form["description"]
		if !ok {
			h.errorPage(w, r, http.StatusBadRequest, "description field not found")
			return
		}
		file := r.MultipartForm.File["image"]
		post := models.Post{
			UserId:      userId,
			Title:       title[0],
			Tags:        tags[0],
			Description: description[0],
		}
		id, err := h.services.Post.Create(r.Context(), post, file)
		if err != nil {
			h.errorPage(w, r, http.StatusInternalServerError, err.Error())
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/post/%d", id), http.StatusSeeOther)
	default:
		h.errorPage(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
	}
}

func (h *Handler) post(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(userId).(uint)
	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/post/"))
	if err != nil {
		h.errorPage(w, r, http.StatusNotFound, "post not found")
		return
	}
	post, err := h.services.Post.GetById(r.Context(), uint(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.errorPage(w, r, http.StatusNotFound, http.StatusText(http.StatusNotFound))
			return
		}
		h.errorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	switch r.Method {
	case http.MethodGet:
		user, err := h.services.User.GetById(r.Context(), userId)
		if err != nil {
			h.errorPage(w, r, http.StatusInternalServerError, err.Error())
			return
		}
		info := models.Info{
			User:        user,
			ProfileUser: user,
			Post:        post,
		}
		if err := h.tmpl.ExecuteTemplate(w, "post.html", info); err != nil {
			h.errorPage(w, r, http.StatusInternalServerError, err.Error())
		}
	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			h.errorPage(w, r, http.StatusInternalServerError, err.Error())
			return
		}
		text, ok := r.Form["text"]
		if !ok {
			h.errorPage(w, r, http.StatusBadRequest, "comment field not found")
			return
		}
		comment := models.Comment{
			PostId: post.Id,
			UserId: userId,
			Text:   text[0],
		}
		if err = h.services.Comment.Create(r.Context(), comment); err != nil {
			h.errorPage(w, r, http.StatusInternalServerError, err.Error())
			return
		}
		http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
	default:
		h.errorPage(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
	}
}
