package http

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v4"
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

	id, err := decodeId(r)
	if err != nil {
		h.log.Error().Err(err).Msg("Item Handler: Invalid id format")
		RenderJSON(w, http.StatusBadRequest, err)
		return
	}

	item, err := h.ItemService.FindItemByID(r.Context(), id)
	if err != nil {
		if err == pgx.ErrNoRows {
			h.log.Error().Err(err).Msg("Item Handler: Item not found")
			RenderJSON(w, http.StatusNotFound, errors.New(http.StatusText(http.StatusNotFound)))
			return
		}

		h.log.Error().Err(err).Msg("Item Handler: Internal Error")
		RenderJSON(w, http.StatusInternalServerError, err)
		return
	}

	RenderJSON(w, http.StatusOK, item)
	return
}

func (h *Handler) handleGetItems(w http.ResponseWriter, r *http.Request) {

	var filter models.ItemFilter

	if err := json.NewDecoder(r.Body).Decode(&filter); err != nil {
		h.log.Error().Err(err).Msg("Item Handler: Invalid Json")
		RenderJSON(w, http.StatusBadRequest, errors.New("invalid JSON"))
		return
	}

	items, n, err := h.ItemService.FindItems(r.Context(), filter)
	if err != nil {
		if err == pgx.ErrNoRows {
			h.log.Error().Err(err).Msg("Item Handler: Item not found")
			RenderJSON(w, http.StatusNotFound, errors.New("item not found"))
			return
		}
		h.log.Error().Err(err).Msg("Item Handler: Internal Error")
		RenderJSON(w, http.StatusInternalServerError, errors.New(http.StatusText(http.StatusInternalServerError)))
		return
	}

	RenderJSON(w, http.StatusOK, itemResponse{Item: items, Count: n})
	return

}
func (h *Handler) handleCreateItem(w http.ResponseWriter, r *http.Request) {}
func (h *Handler) handleUpdateItem(w http.ResponseWriter, r *http.Request) {}
func (h *Handler) handleDeleteItem(w http.ResponseWriter, r *http.Request) {}

func decodeId(r *http.Request) (int, error) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		return 0, errors.New("invalid id format")
	}
	return id, nil
}

type itemResponse struct {
	Item  []*models.Item `json:"items"`
	Count int            `json:"count"`
}
