package delivery

import (
	"context"
	"forum/internal/service"
	"html/template"
	"net/http"
)

type Handler struct {
	tmpl     *template.Template
	ctx      context.Context
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		tmpl:     template.Must(template.ParseGlob("./ui/templates/*.html")),
		ctx:      context.Background(),
		services: services,
	}
}

func (h *Handler) InitRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", h.homePage)

	mux.HandleFunc("/auth/signin", h.signIn)
	mux.HandleFunc("/auth/signup", h.signUp)
	mux.HandleFunc("/auth/logout", h.logOut)

	mux.HandleFunc("/post/", h.post)
	mux.HandleFunc("/post/create", h.createPost)
	mux.HandleFunc("/post/delete/", h.deletePost)
	mux.HandleFunc("/post/change/", h.changePost)
	mux.HandleFunc("/post/like/", h.likePost)
	mux.HandleFunc("/post/dislike/", h.dislikePost)

	mux.HandleFunc("/comment/delete/", h.deleteComment)
	mux.HandleFunc("/comment/change/", h.changeComment)
	mux.HandleFunc("/comment/like/", h.likeComment)
	mux.HandleFunc("/comment/dislike/", h.dislikeComment)

	mux.HandleFunc("/profile/", h.userProfilePage)

	fs := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	return mux
}
