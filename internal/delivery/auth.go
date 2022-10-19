package delivery

import (
	"errors"
	"forum/internal/service"
	"forum/models"
	"net/http"
	"time"
)

func (h *Handler) signUp(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/auth/signup" {
		h.errorPage(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	switch r.Method {
	case http.MethodGet:
		if err := h.tmpl.ExecuteTemplate(w, "register.html", nil); err != nil {
			h.errorPage(w, http.StatusInternalServerError, err.Error())
			return
		}
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			h.errorPage(w, http.StatusInternalServerError, err.Error())
			return
		}
		email, ok := r.Form["email"]
		if !ok {
			h.errorPage(w, http.StatusBadRequest, "email field not found")
			return
		}
		username, ok := r.Form["username"]
		if !ok {
			h.errorPage(w, http.StatusBadRequest, "username field not found")
			return
		}
		password, ok := r.Form["password"]
		if !ok {
			h.errorPage(w, http.StatusBadRequest, "password field not found")
			return
		}
		verifyPassword, ok := r.Form["verifyPassword"]
		if !ok {
			h.errorPage(w, http.StatusBadRequest, "verifyPassword field not found")
			return
		}
		user := models.User{
			Email:          email[0],
			Username:       username[0],
			Password:       password[0],
			VerifyPassword: verifyPassword[0],
		}
		if err := h.services.Auth.CreateUser(user); err != nil {
			if errors.Is(err, service.ErrInvalidUserName) ||
				errors.Is(err, service.ErrPasswordDontMatch) ||
				errors.Is(err, service.ErrInvalidEmail) ||
				errors.Is(err, service.ErrInvalidPassword) ||
				errors.Is(err, service.ErrUsernameTaken) {
				h.errorPage(w, http.StatusBadRequest, err.Error())
				return
			}
			h.errorPage(w, http.StatusInternalServerError, err.Error())
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (h *Handler) signIn(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/auth/signin" {
		h.errorPage(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	switch r.Method {
	case http.MethodGet:
		if err := h.tmpl.ExecuteTemplate(w, "login.html", nil); err != nil {
			h.errorPage(w, http.StatusInternalServerError, err.Error())
			return
		}
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			h.errorPage(w, http.StatusInternalServerError, err.Error())
			return
		}
		username, ok := r.Form["username"]
		if !ok {
			h.errorPage(w, http.StatusBadRequest, "username field not found")
			return
		}
		password, ok := r.Form["password"]
		if !ok {
			h.errorPage(w, http.StatusBadRequest, "password field not found")
			return
		}
		sessionToken, expiresAt, err := h.services.Auth.GenerateSessionToken(username[0], password[0])
		if err != nil {
			if errors.Is(err, service.ErrUserNotFound) {
				h.errorPage(w, http.StatusBadRequest, err.Error())
				return
			}
			h.errorPage(w, http.StatusInternalServerError, err.Error())
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   sessionToken,
			Expires: expiresAt,
			Path:    "/",
		})
		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		h.errorPage(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
	}
}

func (h *Handler) logOut(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/auth/logout" {
		h.errorPage(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	if r.Method != http.MethodGet {
		h.errorPage(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			h.errorPage(w, http.StatusUnauthorized, err.Error())
			return
		}
		h.errorPage(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.DeleteSessionToken(c.Value); err != nil {
		h.errorPage(w, http.StatusInternalServerError, err.Error())
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now(),
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
