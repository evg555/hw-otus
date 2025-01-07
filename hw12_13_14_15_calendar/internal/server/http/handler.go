package internalhttp

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
	router *mux.Router
	app    Application
	Logger Logger
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

func (h *Handler) HomeHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello World"))
}
