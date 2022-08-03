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
	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/comment/like/"))
	if err != nil {
		h.errorPage(w, http.StatusNotFound, err.Error())
		return
	}
	if r.Method != http.MethodPost {
		h.errorPage(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	if user == (models.User{}) {
		h.errorPage(w, http.StatusUnauthorized, "can not like comment")
		return
	}
	if err := h.Services.LikeComment(id, user.Username); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.errorPage(w, http.StatusNotFound, err.Error())
			return
		}
		h.errorPage(w, http.StatusInternalServerError, err.Error())
		return
	}
	postId, err := h.Services.GetPostIdByCommentId(id)
	if err != nil {
		h.errorPage(w, http.StatusInternalServerError, err.Error())
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/post/%d", postId), http.StatusSeeOther)
}

func (h *Handler) dislikeComment(w http.ResponseWriter, r *http.Request) {
	user := h.userIdentity(w, r)
	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/comment/dislike/"))
	if err != nil {
		h.errorPage(w, http.StatusNotFound, err.Error())
		return
	}
	if r.Method != http.MethodPost {
		h.errorPage(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	if user == (models.User{}) {
		h.errorPage(w, http.StatusUnauthorized, "can not dislike comment")
		return
	}
	if err := h.Services.DislikeComment(id, user.Username); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.errorPage(w, http.StatusNotFound, err.Error())
			return
		}
		h.errorPage(w, http.StatusInternalServerError, err.Error())
		return
	}
	postId, err := h.Services.GetPostIdByCommentId(id)
	if err != nil {
		h.errorPage(w, http.StatusInternalServerError, err.Error())
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/post/%d", postId), http.StatusSeeOther)
}
