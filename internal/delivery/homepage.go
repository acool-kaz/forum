package delivery

import (
	"context"
	"forum/internal/models"
	"net/http"
)

func (h *Handler) homePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		h.errorPage(w, r, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}

	if r.Method != http.MethodGet {
		h.errorPage(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	userId := r.Context().Value(userId).(uint)

	user, err := h.services.User.GetOneBy(context.WithValue(r.Context(), models.UserId, userId))
	if err != nil {
		h.errorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	ctx := r.Context()

	likeFilter, ok := r.URL.Query()["likes"]
	if ok {
		ctx = context.WithValue(ctx, models.Filter, "likes-"+likeFilter[0])
	}

	timeFilter, ok := r.URL.Query()["time"]
	if ok {
		ctx = context.WithValue(ctx, models.Filter, "time-"+timeFilter[0])
	}

	category, ok := r.URL.Query()["category"]
	if ok {
		ctx = context.WithValue(ctx, models.Tags, category[0])
	}

	posts, err := h.services.Post.GetAll(ctx)
	if err != nil {
		h.errorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	info := models.Info{
		User:  user,
		Posts: posts,
	}

	if err := h.tmpl.ExecuteTemplate(w, "index.html", info); err != nil {
		h.errorPage(w, r, http.StatusInternalServerError, err.Error())
	}
}
