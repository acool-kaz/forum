package delivery

import (
	"errors"
	"forum/internal/models"
	"net/http"
)

func (h *Handler) likeReaction(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/like" {
		h.errorPage(w, r, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}

	if r.Method != http.MethodPost {
		h.errorPage(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	userId := r.Context().Value(userId).(uint)
	if userId == 0 {
		h.errorPage(w, r, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	postId := r.URL.Query().Get("post")
	commentId := r.URL.Query().Get("comment")

	if err := h.services.Reaction.Set(r.Context(), postId, commentId, models.LikeReaction, userId); err != nil {
		if errors.Is(err, models.ErrInvalidReaction) {
			h.errorPage(w, r, http.StatusNotFound, err.Error())
			return
		}
		h.errorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	http.Redirect(w, r, "/post/"+postId, http.StatusSeeOther)
}

func (h *Handler) dislikeReaction(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/dislike" {
		h.errorPage(w, r, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}

	if r.Method != http.MethodPost {
		h.errorPage(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	userId := r.Context().Value(userId).(uint)
	if userId == 0 {
		h.errorPage(w, r, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	postId := r.URL.Query().Get("post")
	commentId := r.URL.Query().Get("comment")

	if err := h.services.Reaction.Set(r.Context(), postId, commentId, models.DislikeReaction, userId); err != nil {
		if errors.Is(err, models.ErrInvalidReaction) {
			h.errorPage(w, r, http.StatusNotFound, err.Error())
			return
		}
		h.errorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	http.Redirect(w, r, "/post/"+postId, http.StatusSeeOther)
}
