package http

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (h *Handler) mountItemsRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/", h.handleGetItems)
	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", h.handleGetItem)
	})

	return r
}

func (h *Handler) handleGetItem(w http.ResponseWriter, r *http.Request) {}
func (h *Handler) handleGetItems(w http.ResponseWriter, r *http.Request) {}
func (h *Handler) handleCreateItem(w http.ResponseWriter, r *http.Request) {}
func (h *Handler) handleUpdateItem(w http.ResponseWriter, r *http.Request) {}
func (h *Handler) handleDeleteItem(w http.ResponseWriter, r *http.Request) {}
