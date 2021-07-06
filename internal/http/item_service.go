package http

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/symaster1995/ms-starter/internal/models"
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

func (h *Handler) handleGetItem(w http.ResponseWriter, r *http.Request) {
	var filter models.ItemFilter

	if err := json.NewDecoder(r.Body).Decode(&filter); err != nil {
		RenderJSON(w, http.StatusBadRequest, errors.New("invalid JSON"))
		return
	}
}

func (h *Handler) handleGetItems(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "application/json")
	w.Write([]byte(`{}`))

}
func (h *Handler) handleCreateItem(w http.ResponseWriter, r *http.Request) {}
func (h *Handler) handleUpdateItem(w http.ResponseWriter, r *http.Request) {}
func (h *Handler) handleDeleteItem(w http.ResponseWriter, r *http.Request) {}
