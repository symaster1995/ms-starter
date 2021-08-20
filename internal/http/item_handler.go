package http

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/symaster1995/ms-starter/internal/models"
	"github.com/symaster1995/ms-starter/pkg/errors"
	"net/http"
	"strconv"
)

func (h *Handler) mountItemsRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/", h.handleGetItems)
	r.Post("/", h.handleCreateItem)
	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", h.handleGetItem)
		r.Patch("/", h.handleUpdateItem)
		r.Delete("/", h.handleDeleteItem)
	})
	return r
}

func (h *Handler) handleGetItem(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		h.log.Error().Err(err).Msg("parsing id failed")
		ErrorJSON(w, errors.Errorf(errors.ErrInvalid, "invalid id format"))
		return
	}

	item, err := h.ItemService.FindItemByID(r.Context(), id)
	if err != nil {
		h.log.Error().Err(err).Msg("find item failed")
		ErrorJSON(w, err)
		return
	}

	RenderJSON(w, http.StatusOK, item)
	return
}

func (h *Handler) handleGetItems(w http.ResponseWriter, r *http.Request) {
	var filter models.ItemFilter

	if err := json.NewDecoder(r.Body).Decode(&filter); err != nil {
		h.log.Error().Err(err).Msg("listing item failed")
		ErrorJSON(w, errors.Errorf(errors.ErrInvalid, "invalid JSON body"))
		return
	}

	items, n, err := h.ItemService.FindItems(r.Context(), filter)
	if err != nil {
		h.log.Error().Err(err).Msg("listing item failed")
		ErrorJSON(w, err)
		return
	}

	RenderJSON(w, http.StatusOK, itemResponse{Item: items, Count: n})
	return
}

func (h *Handler) handleCreateItem(w http.ResponseWriter, r *http.Request) {
	var item models.Item

	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		h.log.Error().Err(err).Msg("creating item failed")
		ErrorJSON(w, errors.Errorf(errors.ErrInvalid, "invalid JSON body"))
		return
	}

	err := h.ItemService.CreateItem(r.Context(), &item)
	if err != nil {
		h.log.Error().Err(err).Msg("creating item failed")
		ErrorJSON(w, err)
		return
	}

	RenderJSON(w, http.StatusCreated, item)
	return
}

func (h *Handler) handleUpdateItem(w http.ResponseWriter, r *http.Request) {
	var itemUpd models.ItemUpdate

	id, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		h.log.Error().Err(err).Msg("parsing id failed")
		ErrorJSON(w, errors.Errorf(errors.ErrInvalid, "invalid id format"))
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&itemUpd); err != nil {
		h.log.Error().Err(err).Msg("updating item failed")
		ErrorJSON(w, errors.Errorf(errors.ErrInvalid, "invalid JSON body"))
		return
	}

	item, err := h.ItemService.UpdateItem(r.Context(), id, itemUpd)
	if err != nil {
		h.log.Error().Err(err).Msg("updating item failed")
		ErrorJSON(w, err)
		return
	}

	RenderJSON(w, http.StatusOK, item)
	return
}

func (h *Handler) handleDeleteItem(w http.ResponseWriter, r *http.Request) {}

type itemResponse struct {
	Item  []*models.Item `json:"items"`
	Count int            `json:"count"`
}