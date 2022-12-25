package delivery

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

func (h *Handler) errorPage(w http.ResponseWriter, r *http.Request, status int, msg string) {
	w.WriteHeader(status)
	fmt.Printf("\r%s %s [%s]\t%s%s - %d - %s\n", time.Now().Format("2006/01/02 15:04:05"), r.Proto, r.Method, r.Host, r.RequestURI, status, http.StatusText(status))
	fmt.Print(msg)
	data := struct {
		Status  int
		Message string
	}{
		Status:  status,
		Message: http.StatusText(status),
	}
	if data.Status != 500 {
		temp := strings.Split(msg, ":")
		data.Message = temp[len(temp)-1]
	}
	if err := h.tmpl.ExecuteTemplate(w, "error.html", data); err != nil {
		fmt.Fprintf(w, "%d - %s\n", data.Status, data.Message)
	}
}
