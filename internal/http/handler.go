package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"net/http"
)

type Handler struct {
	name   string
	router chi.Router
	log    *zap.Logger
}

func NewHandler() *Handler {
	return &Handler{
		name: "Handler",
		log:  zap.NewNop(),
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}


func (h *Handler) configureRouter() {

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	h.router = r
}
