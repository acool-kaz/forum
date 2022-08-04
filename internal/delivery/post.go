package delivery

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/models"
	"net/http"
	"strconv"
	"strings"
)

func (h *Handler) post(w http.ResponseWriter, r *http.Request) {
	user := h.userIdentity(w, r)
	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/post/"))
	if err != nil {
		h.errorPage(w, http.StatusBadRequest, "post not found")
	}
	post, err := h.Services.GetPostById(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.errorPage(w, http.StatusNotFound, err.Error())
			return
		}
		h.errorPage(w, http.StatusInternalServerError, err.Error())
		return
	}
	switch r.Method {
	case http.MethodGet:
		comments, err := h.Services.GetComments(post.Id)
		if err != nil {
			h.errorPage(w, http.StatusInternalServerError, err.Error())
			return
		}
		similarPosts, err := h.Services.GetSimilarPosts(post.Id)
		if err != nil {
			h.errorPage(w, http.StatusInternalServerError, err.Error())
			return
		}
		postLikes, err := h.Services.GetPostLikes(post.Id)
		if err != nil {
			h.errorPage(w, http.StatusInternalServerError, err.Error())
			return
		}
		postDisikes, err := h.Services.GetPostDislikes(post.Id)
		if err != nil {
			h.errorPage(w, http.StatusInternalServerError, err.Error())
			return
		}
		commentsLikes, err := h.Services.GetCommentLikes(post.Id)
		if err != nil {
			h.errorPage(w, http.StatusInternalServerError, err.Error())
			return
		}
		commentsDislikes, err := h.Services.GetCommentDislikes(post.Id)
		if err != nil {
			h.errorPage(w, http.StatusInternalServerError, err.Error())
			return
		}
		info := models.Info{
			SimilarPosts:     similarPosts,
			Post:             post,
			PostLikes:        postLikes,
			PostDislikes:     postDisikes,
			User:             user,
			Comments:         comments,
			CommentsLikes:    commentsLikes,
			CommentsDislikes: commentsDislikes,
		}
		if err := h.Tmpl.ExecuteTemplate(w, "post.html", info); err != nil {
			h.errorPage(w, http.StatusInternalServerError, err.Error())
		}
	case http.MethodPost:
		if user == (models.User{}) {
			h.errorPage(w, http.StatusUnauthorized, "cant post comment")
			return
		}
		if err := r.ParseForm(); err != nil {
			h.errorPage(w, http.StatusInternalServerError, err.Error())
			return
		}
		comment, ok := r.Form["comment"]
		if !ok {
			h.errorPage(w, http.StatusBadRequest, "comment field not foud")
			return
		}
		newComment := models.Comment{
			PostId:  post.Id,
			Creater: user.Username,
			Text:    comment[0],
		}
		if err := h.Services.Comment.CreateComment(newComment); err != nil {
			h.errorPage(w, http.StatusInternalServerError, err.Error())
			return
		}
		http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
	default:
		h.errorPage(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
	}
}

func (h *Handler) createPost(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post/create" {
		h.errorPage(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	user := h.userIdentity(w, r)
	if user == (models.User{}) {
		h.errorPage(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}
	switch r.Method {
	case http.MethodGet:
		info := models.Info{
			User: user,
		}
		if err := h.Tmpl.ExecuteTemplate(w, "createPost.html", info); err != nil {
			h.errorPage(w, http.StatusInternalServerError, err.Error())
		}
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			h.errorPage(w, http.StatusInternalServerError, err.Error())
			return
		}
		category, ok := r.Form["category"]
		if !ok {
			h.errorPage(w, http.StatusBadRequest, "category field not found")
			return
		}
		title, ok := r.Form["title"]
		if !ok {
			h.errorPage(w, http.StatusBadRequest, "title field not found")
			return
		}
		description, ok := r.Form["description"]
		if !ok {
			h.errorPage(w, http.StatusBadRequest, "description field not found")
			return
		}
		post := models.Post{
			Creater:     user.Username,
			Category:    category,
			Title:       title[0],
			Description: description[0],
		}
		if err := h.Services.CreatePost(post); err != nil {
			h.errorPage(w, http.StatusInternalServerError, err.Error())
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		h.errorPage(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
	}
}

func (h *Handler) likePost(w http.ResponseWriter, r *http.Request) {
	user := h.userIdentity(w, r)
	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/post/like/"))
	if err != nil {
		h.errorPage(w, http.StatusNotFound, err.Error())
		return
	}
	if r.Method != http.MethodPost {
		h.errorPage(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	if user == (models.User{}) {
		h.errorPage(w, http.StatusUnauthorized, "can not like post")
		return
	}
	if _, err := h.Services.GetPostById(id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.errorPage(w, http.StatusNotFound, err.Error())
			return
		}
		h.errorPage(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := h.Services.LikePost(id, user.Username); err != nil {
		h.errorPage(w, http.StatusInternalServerError, err.Error())
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/post/%d", id), http.StatusSeeOther)
}

func (h *Handler) dislikePost(w http.ResponseWriter, r *http.Request) {
	user := h.userIdentity(w, r)
	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/post/dislike/"))
	if err != nil {
		h.errorPage(w, http.StatusNotFound, err.Error())
		return
	}
	if r.Method != http.MethodPost {
		h.errorPage(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	if user == (models.User{}) {
		h.errorPage(w, http.StatusUnauthorized, "can not dislike post")
		return
	}
	if _, err := h.Services.GetPostById(id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.errorPage(w, http.StatusNotFound, err.Error())
			return
		}
		h.errorPage(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := h.Services.DislikePost(id, user.Username); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.errorPage(w, http.StatusNotFound, err.Error())
			return
		}
		h.errorPage(w, http.StatusInternalServerError, err.Error())
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/post/%d", id), http.StatusSeeOther)
}