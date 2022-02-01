package products

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"github.com/symaster1995/ms-starter/internal/products/models"
	"github.com/symaster1995/ms-starter/pkg/errors"
	response "github.com/symaster1995/ms-starter/pkg/http"
	"net/http"
	"strconv"
)

type ItemHandler struct {
	name string
	chi.Router
	log         *zerolog.Logger
	ItemService models.ItemService
}

func NewItemHandler(logger *zerolog.Logger, itemService models.ItemService) *ItemHandler {

	ih := &ItemHandler{
		name:        "ItemHandler",
		log:         logger,
		ItemService: itemService,
	}

	r := chi.NewRouter()
	r.Use(
		middleware.RequestID,
		middleware.RealIP,
		middleware.Recoverer,
	)

	r.Route("/", func(r chi.Router) {
		r.Get("/", ih.handleGetItems)
		r.Post("/", ih.handleCreateItem)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", ih.handleGetItem)
			r.Patch("/", ih.handleUpdateItem)
			r.Delete("/", ih.handleDeleteItem)
		})
	})

	ih.Router = r
	return ih
}

func (i *ItemHandler) handleGetItem(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		i.log.Error().Err(err).Msg("parsing id failed")
		response.ErrorJSON(w, errors.Errorf(errors.ErrInvalid, "invalid id format"))
		return
	}

	item, err := i.ItemService.FindItemByID(r.Context(), id)
	if err != nil {
		i.log.Error().Err(err).Msg("find item failed")
		response.ErrorJSON(w, err)
		return
	}

	response.RenderJSON(w, http.StatusOK, item)
	return
}

func (i *ItemHandler) handleGetItems(w http.ResponseWriter, r *http.Request) {
	var filter models.ItemFilter

	if err := decodeJson(r, &filter); err != nil {
		i.log.Error().Err(err).Msg("listing item failed")
		response.ErrorJSON(w, errors.Errorf(errors.ErrInvalid, "invalid JSON body"))
		return
	}

	items, n, err := i.ItemService.FindItems(r.Context(), filter)
	if err != nil {
		i.log.Error().Err(err).Msg("listing item failed")
		response.ErrorJSON(w, err)
		return
	}

	response.RenderJSON(w, http.StatusOK, itemListResponse{Item: items, Count: n})
	return
}

func (i *ItemHandler) handleCreateItem(w http.ResponseWriter, r *http.Request) {
	var item models.Item

	if err := decodeJson(r, &item); err != nil {
		i.log.Error().Err(err).Msg("creating item failed")
		response.ErrorJSON(w, errors.Errorf(errors.ErrInvalid, "invalid JSON body"))
		return
	}

	err := i.ItemService.CreateItem(r.Context(), &item)
	if err != nil {
		i.log.Error().Err(err).Msg("creating item failed")
		response.ErrorJSON(w, err)
		return
	}

	response.RenderJSON(w, http.StatusCreated, item)
	return
}

func (i *ItemHandler) handleUpdateItem(w http.ResponseWriter, r *http.Request) {
	var itemUpd models.ItemUpdate

	id, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		i.log.Error().Err(err).Msg("parsing id failed")
		response.ErrorJSON(w, errors.Errorf(errors.ErrInvalid, "invalid id format"))
		return
	}

	if err := decodeJson(r, &itemUpd); err != nil {
		i.log.Error().Err(err).Msg("updating item failed")
		response.ErrorJSON(w, errors.Errorf(errors.ErrInvalid, "invalid JSON body"))
		return
	}

	item, err := i.ItemService.UpdateItem(r.Context(), id, itemUpd)
	if err != nil {
		i.log.Error().Err(err).Msg("updating item failed")
		response.ErrorJSON(w, err)
		return
	}

	response.RenderJSON(w, http.StatusOK, item)
	return
}

func (i *ItemHandler) handleDeleteItem(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		i.log.Error().Err(err).Msg("parsing id failed")
		response.ErrorJSON(w, errors.Errorf(errors.ErrInvalid, "invalid id format"))
		return
	}

	if err := i.ItemService.DeleteItem(r.Context(), id); err != nil {
		i.log.Error().Err(err).Msg("deleting item failed")
		response.ErrorJSON(w, err)
		return
	}

	response.RenderJSON(w, http.StatusOK, `{"message": "item deleted"}`)
	return
}

func decodeJson(r *http.Request, data interface{}) error {
	return json.NewDecoder(r.Body).Decode(data)
}

type itemListResponse struct {
	Item  []*models.Item `json:"items"`
	Count int            `json:"count"`
}
