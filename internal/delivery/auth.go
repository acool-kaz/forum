package delivery

import (
	"errors"
	"forum/internal/models"
	"net/http"
	"time"
)

func (h *Handler) signUp(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(userId).(uint)
	if user != 0 {
		h.errorPage(w, r, http.StatusBadRequest, "you already in")
		return
	}
	if r.URL.Path != "/auth/signup" {
		h.errorPage(w, r, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	switch r.Method {
	case http.MethodGet:
		if err := h.tmpl.ExecuteTemplate(w, "register.html", nil); err != nil {
			h.errorPage(w, r, http.StatusInternalServerError, err.Error())
			return
		}
	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			h.errorPage(w, r, http.StatusInternalServerError, err.Error())
			return
		}
		email, ok := r.Form["email"]
		if !ok {
			h.errorPage(w, r, http.StatusBadRequest, "email field not found")
			return
		}
		username, ok := r.Form["username"]
		if !ok {
			h.errorPage(w, r, http.StatusBadRequest, "username field not found")
			return
		}
		password, ok := r.Form["password"]
		if !ok {
			h.errorPage(w, r, http.StatusBadRequest, "password field not found")
			return
		}
		verifyPassword, ok := r.Form["verifyPassword"]
		if !ok {
			h.errorPage(w, r, http.StatusBadRequest, "verifyPassword field not found")
			return
		}
		user := models.User{
			Email:          email[0],
			Username:       username[0],
			Password:       password[0],
			VerifyPassword: verifyPassword[0],
		}
		_, err = h.services.User.Create(r.Context(), user)
		if err != nil {
			if errors.Is(err, models.ErrInvalidUser) {
				h.errorPage(w, r, http.StatusBadRequest, err.Error())
				return
			}
			h.errorPage(w, r, http.StatusInternalServerError, err.Error())
			return
		}
		var session models.Session
		session, err = h.services.Session.GenerateSessionToken(r.Context(), user.Username, user.Password)
		if err != nil {
			if errors.Is(err, models.ErrUserNotFound) {
				h.errorPage(w, r, http.StatusBadRequest, err.Error())
				return
			}
			h.errorPage(w, r, http.StatusInternalServerError, err.Error())
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   session.Token,
			Expires: session.ExpiresAt,
			Path:    "/",
		})
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (h *Handler) signIn(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(userId).(uint)
	if user != 0 {
		h.errorPage(w, r, http.StatusBadRequest, "you already in")
		return
	}
	if r.URL.Path != "/auth/signin" {
		h.errorPage(w, r, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	switch r.Method {
	case http.MethodGet:
		if err := h.tmpl.ExecuteTemplate(w, "login.html", nil); err != nil {
			h.errorPage(w, r, http.StatusInternalServerError, err.Error())
			return
		}
	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			h.errorPage(w, r, http.StatusInternalServerError, err.Error())
			return
		}
		username, ok := r.Form["username"]
		if !ok {
			h.errorPage(w, r, http.StatusBadRequest, "username field not found")
			return
		}
		password, ok := r.Form["password"]
		if !ok {
			h.errorPage(w, r, http.StatusBadRequest, "password field not found")
			return
		}
		var session models.Session
		session, err = h.services.Session.GenerateSessionToken(r.Context(), username[0], password[0])
		if err != nil {
			if errors.Is(err, models.ErrUserNotFound) {
				h.errorPage(w, r, http.StatusBadRequest, err.Error())
				return
			}
			h.errorPage(w, r, http.StatusInternalServerError, err.Error())
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   session.Token,
			Expires: session.ExpiresAt,
			Path:    "/",
		})
		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		h.errorPage(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
	}
}

func (h *Handler) logOut(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(userId).(uint)
	if user == 0 {
		h.errorPage(w, r, http.StatusBadRequest, "cant log-out, without log-in")
		return
	}
	if r.Method != http.MethodGet {
		h.errorPage(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	if r.URL.Path != "/auth/logout" {
		h.errorPage(w, r, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}

	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			h.errorPage(w, r, http.StatusUnauthorized, err.Error())
			return
		}
		h.errorPage(w, r, http.StatusBadRequest, err.Error())
		return
	}

	if err = h.services.Session.DeleteSessionToken(r.Context(), c.Value); err != nil {
		h.errorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now(),
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
