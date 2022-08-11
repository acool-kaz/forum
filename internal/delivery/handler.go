package delivery

import (
	"forum/internal/service"
	"html/template"
	"net/http"
)

type Handler struct {
	Mux      *http.ServeMux
	Tmpl     *template.Template
	Services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		Mux:      http.NewServeMux(),
		Tmpl:     template.Must(template.ParseGlob("./ui/templates/*.html")),
		Services: services,
	}
}

func (h *Handler) InitRoutes() {
	h.Mux.HandleFunc("/", h.homePage)

	h.Mux.HandleFunc("/auth/signin", h.signIn)
	h.Mux.HandleFunc("/auth/signup", h.signUp)
	h.Mux.HandleFunc("/auth/logout", h.logOut)

	h.Mux.HandleFunc("/post/", h.post)
	h.Mux.HandleFunc("/post/create", h.createPost)

	h.Mux.HandleFunc("/post/like/", h.likePost)
	h.Mux.HandleFunc("/post/dislike/", h.dislikePost)
	h.Mux.HandleFunc("/comment/like/", h.likeComment)
	h.Mux.HandleFunc("/comment/dislike/", h.dislikeComment)

	h.Mux.HandleFunc("/profile/", h.userProfilePage)

	fs := http.FileServer(http.Dir("./ui/static"))
	h.Mux.Handle("/static/", http.StripPrefix("/static/", fs))
}
