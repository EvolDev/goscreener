package index

import (
	"bytes"
	"net/http"
)

const version = "goscreener v0.0.3"

type Handler struct {
}

func NewIndexHandler() *Handler {
	return &Handler{}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var buffer bytes.Buffer
	buffer.WriteString(version)
	_, err := w.Write(buffer.Bytes())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
