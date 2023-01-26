package delivery

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"forum/internal/models"
	"net/http"
	"strings"
)

func (h *Handler) userProfilePage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.errorPage(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	username := strings.TrimPrefix(r.URL.Path, "/profile/")

	profileUser, err := h.services.User.GetOneBy(context.WithValue(r.Context(), models.Username, username))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.errorPage(w, r, http.StatusNotFound, err.Error())
			return
		}
		h.errorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	userId := r.Context().Value(userId).(uint)

	curUser, err := h.services.User.GetOneBy(context.WithValue(r.Context(), models.UserId, userId))
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			h.errorPage(w, r, http.StatusInternalServerError, err.Error())
			return
		}
	}

	ctx := r.Context()

	postFilter := r.URL.Query().Get("posts")
	if postFilter == "" {
		http.Redirect(w, r, fmt.Sprintf("/profile/%s?posts=created", username), http.StatusSeeOther)
		return
	}

	ctx = context.WithValue(ctx, models.ProfilePostFilter, postFilter)
	ctx = context.WithValue(ctx, models.UserId, profileUser.Id)

	posts, err := h.services.Post.GetAll(ctx)
	if err != nil {
		h.errorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	info := models.Info{
		User:        curUser,
		ProfileUser: profileUser,
		Posts:       posts,
	}

	if err := h.tmpl.ExecuteTemplate(w, "user.html", info); err != nil {
		h.errorPage(w, r, http.StatusInternalServerError, err.Error())
	}
}
