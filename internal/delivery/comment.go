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

func (h *Handler) likeComment(w http.ResponseWriter, r *http.Request) {
	user := h.userIdentity(w, r)
	if user == (models.User{}) {
		h.errorPage(w, http.StatusUnauthorized, "can not like comment")
		return
	}
	if r.Method != http.MethodPost {
		h.errorPage(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/comment/like/"))
	if err != nil {
		h.errorPage(w, http.StatusNotFound, err.Error())
		return
	}
	comment, err := h.Services.GetCommentById(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.errorPage(w, http.StatusNotFound, err.Error())
			return
		}
		h.errorPage(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := h.Services.LikeComment(comment.Id, user.Username); err != nil {
		h.errorPage(w, http.StatusInternalServerError, err.Error())
		return
	}
	// post, err := h.Services.GetPostByCommentId(id)
	// if err != nil {
	// 	h.errorPage(w, http.StatusInternalServerError, err.Error())
	// 	return
	// }
	if comment.Creater != user.Username {
		newNotify := models.Notify{
			From:        user.Username,
			To:          comment.Creater,
			Description: "post comment like",
			PostId:      comment.PostId,
		}
		if err := h.Services.AddNewNotify(newNotify); err != nil {
			h.errorPage(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
	http.Redirect(w, r, fmt.Sprintf("/post/%d", comment.PostId), http.StatusSeeOther)
}

func (h *Handler) dislikeComment(w http.ResponseWriter, r *http.Request) {
	user := h.userIdentity(w, r)
	if user == (models.User{}) {
		h.errorPage(w, http.StatusUnauthorized, "can not dislike comment")
		return
	}
	if r.Method != http.MethodPost {
		h.errorPage(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/comment/dislike/"))
	if err != nil {
		h.errorPage(w, http.StatusNotFound, err.Error())
		return
	}
	comment, err := h.Services.GetCommentById(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.errorPage(w, http.StatusNotFound, err.Error())
			return
		}
		h.errorPage(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := h.Services.DislikeComment(id, user.Username); err != nil {
		h.errorPage(w, http.StatusInternalServerError, err.Error())
		return
	}
	// post, err := h.Services.GetPostByCommentId(id)
	// if err != nil {
	// 	h.errorPage(w, http.StatusInternalServerError, err.Error())
	// 	return
	// }
	if comment.Creater != user.Username {
		newNotify := models.Notify{
			From:        user.Username,
			To:          comment.Creater,
			Description: "post comment dislike",
			PostId:      comment.PostId,
		}
		if err := h.Services.AddNewNotify(newNotify); err != nil {
			h.errorPage(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
	http.Redirect(w, r, fmt.Sprintf("/post/%d", comment.PostId), http.StatusSeeOther)
}
