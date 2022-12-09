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
	user := h.ctx.Value(uCtx).(models.User)
	if user == (models.User{}) {
		h.errorPage(w, r, http.StatusUnauthorized, "can not like comment")
		return
	}
	if r.Method != http.MethodPost {
		h.errorPage(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/comment/like/"))
	if err != nil {
		h.errorPage(w, r, http.StatusNotFound, err.Error())
		return
	}
	comment, err := h.services.GetCommentById(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.errorPage(w, r, http.StatusNotFound, err.Error())
			return
		}
		h.errorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	if err := h.services.LikeComment(comment.Id, user.Username); err != nil {
		h.errorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	if comment.Creater != user.Username {
		newNotification := models.Notification{
			From:        user.Username,
			To:          comment.Creater,
			Description: "liked comment under the post",
			PostId:      comment.PostId,
		}
		if err := h.services.AddNewNotification(newNotification); err != nil {
			h.errorPage(w, r, http.StatusInternalServerError, err.Error())
			return
		}
	}
	http.Redirect(w, r, fmt.Sprintf("/post/%d", comment.PostId), http.StatusSeeOther)
}

func (h *Handler) dislikeComment(w http.ResponseWriter, r *http.Request) {
	user := h.ctx.Value(uCtx).(models.User)
	if user == (models.User{}) {
		h.errorPage(w, r, http.StatusUnauthorized, "can not dislike comment")
		return
	}
	if r.Method != http.MethodPost {
		h.errorPage(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/comment/dislike/"))
	if err != nil {
		h.errorPage(w, r, http.StatusNotFound, err.Error())
		return
	}
	comment, err := h.services.GetCommentById(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.errorPage(w, r, http.StatusNotFound, err.Error())
			return
		}
		h.errorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	if err := h.services.DislikeComment(id, user.Username); err != nil {
		h.errorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	if comment.Creater != user.Username {
		newNotification := models.Notification{
			From:        user.Username,
			To:          comment.Creater,
			Description: "disliked comment under the post",
			PostId:      comment.PostId,
		}
		if err := h.services.AddNewNotification(newNotification); err != nil {
			h.errorPage(w, r, http.StatusInternalServerError, err.Error())
			return
		}
	}
	http.Redirect(w, r, fmt.Sprintf("/post/%d", comment.PostId), http.StatusSeeOther)
}

func (h *Handler) deleteComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.errorPage(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	user := h.ctx.Value(uCtx).(models.User)
	if user == (models.User{}) {
		h.errorPage(w, r, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}
	commentId, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/comment/delete/"))
	if err != nil {
		h.errorPage(w, r, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	comment, err := h.services.GetCommentById(commentId)
	if err != nil {
		h.errorPage(w, r, http.StatusNotFound, err.Error())
		return
	}
	if user.Username != comment.Creater {
		h.errorPage(w, r, http.StatusBadRequest, "you cant delete this comment")
		return
	}
	if err := h.services.DeleteComment(comment); err != nil {
		h.errorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/post/%d", comment.PostId), http.StatusSeeOther)
}

func (h *Handler) changeComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.errorPage(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	user := h.ctx.Value(uCtx).(models.User)
	if user == (models.User{}) {
		h.errorPage(w, r, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}
	commentId, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/comment/change/"))
	if err != nil {
		h.errorPage(w, r, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	comment, err := h.services.GetCommentById(commentId)
	if err != nil {
		h.errorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	if user.Username != comment.Creater {
		h.errorPage(w, r, http.StatusBadRequest, "you can`t change this comment")
		return
	}
	if err := r.ParseForm(); err != nil {
		h.errorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	text, ok := r.Form["text"]
	if !ok {
		h.errorPage(w, r, http.StatusInternalServerError, "text field not found")
		return
	}
	comment.Text = strings.Join(text, "")
	if err := h.services.ChangeComment(comment); err != nil {
		h.errorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/post/%d", comment.PostId), http.StatusSeeOther)
}
