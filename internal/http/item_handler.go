package http

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/symaster1995/ms-starter/internal/models"
	"net/http"
	"strconv"
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
	itemId := chi.URLParam(r, "id")
	if itemId == "" {
		h.log.Error().Msg("Invalid url id")
		RenderJSON(w, http.StatusBadRequest, errors.New("invalid url id"))
		return
	}

	id, err := strconv.Atoi(itemId)
	if err != nil {
		h.log.Error().Err(err).Msg("Invalid id format")
		RenderJSON(w, http.StatusBadRequest, errors.New("invalid id format"))
		return
	}

	item, err := h.ItemService.FindItemByID(r.Context(), id)
	if err != nil {
		if err.Error() == "no rows in result set" {
			h.log.Error().Err(err).Msg("Item not found")
			RenderJSON(w, http.StatusNotFound, err)
			return
		}
		h.log.Error().Err(err).Msg("Internal Error")
		RenderJSON(w, http.StatusInternalServerError, err)
		return
	}

	RenderJSON(w, http.StatusOK, item)
	return
}

func (h *Handler) handleGetItems(w http.ResponseWriter, r *http.Request) {

	var filter models.ItemFilter

	if err := json.NewDecoder(r.Body).Decode(&filter); err != nil {
		h.log.Error().Err(err).Msg("Invalid Json")
		RenderJSON(w, http.StatusBadRequest, errors.New("invalid JSON"))
		return
	}

	RenderJSON(w, http.StatusOK, `{}`)
	return

}
func (h *Handler) handleCreateItem(w http.ResponseWriter, r *http.Request) {}
func (h *Handler) handleUpdateItem(w http.ResponseWriter, r *http.Request) {}
func (h *Handler) handleDeleteItem(w http.ResponseWriter, r *http.Request) {}
