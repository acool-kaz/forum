package delivery

import (
	"fmt"
	"log"
	"net/http"
)

func (h *Handler) errorPage(w http.ResponseWriter, status int, msg string) {
	w.WriteHeader(status)
	log.Printf("%d - %s", status, msg)
	data := struct {
		Status  int
		Message string
	}{
		Status:  status,
		Message: http.StatusText(status),
	}
	if err := h.tmpl.ExecuteTemplate(w, "error.html", data); err != nil {
		fmt.Fprintf(w, "%d - %s\n", data.Status, data.Message)
	}
}
