package delivery

import (
	"forum/internal/service"
	"html/template"
	"net/http"
)

type Handler struct {
	tmpl     *template.Template
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		tmpl:     template.Must(template.ParseGlob("./templates/*.html")),
		services: services,
	}
}

func (h *Handler) InitRoutes() *http.ServeMux {
	m := new(middleware)
	m.addMiddleware(h.authMiddleware)
	m.addMiddleware(h.loggingMiddleware)

	mux := http.NewServeMux()
	mux.HandleFunc("/", m.use(h.homePage))

	mux.HandleFunc("/auth/signin", m.use(h.signIn))
	mux.HandleFunc("/auth/signup", m.use(h.signUp))
	mux.HandleFunc("/auth/logout", m.use(h.logOut))

	mux.HandleFunc("/post/", m.use(h.post))
	mux.HandleFunc("/post/create", m.use(h.createPost))
	// mux.HandleFunc("/post/delete/", m.use(h.deletePost))
	// mux.HandleFunc("/post/change/", m.use(h.changePost))

	// mux.HandleFunc("/comment/delete/", m.use(h.deleteComment))
	// mux.HandleFunc("/comment/change/", m.use(h.changeComment))

	mux.HandleFunc("/like", m.use(h.likeReaction))
	mux.HandleFunc("/dislike", m.use(h.dislikeReaction))

	// mux.HandleFunc("/profile/", m.use(h.userProfilePage))

	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	return mux
}
