package mock

import (
	"context"
	"github.com/symaster1995/ms-starter/internal/models"
)

type ItemService struct {
	FindItemByIDFn func(context.Context, int) (*models.Item, error)
	FindItemsFn    func(context.Context, models.ItemFilter) ([]*models.Item, int, error)
	CreateItemFn   func(context.Context, *models.Item) error
	UpdateItemFn   func(context.Context, int, models.ItemUpdate) (*models.Item, error)
	DeleteItemFn   func(context.Context, int) error
}

func NewItemService() *ItemService {
	return &ItemService{
		FindItemByIDFn: func(context.Context, int) (*models.Item, error) { return nil, nil },
		FindItemsFn:    func(context.Context, models.ItemFilter) ([]*models.Item, int, error) { return nil, 0, nil },
		CreateItemFn:   func(context.Context, *models.Item) error { return nil },
		UpdateItemFn:   func(context.Context, int, models.ItemUpdate) (*models.Item, error) { return nil, nil },
		DeleteItemFn:   func(context.Context, int) error { return nil },
	}
}

func (i *ItemService) FindItemByID(ctx context.Context, id int) (*models.Item, error) {
	return i.FindItemByIDFn(ctx,id)
}

func (i *ItemService) FindItems(ctx context.Context, filter models.ItemFilter) ([]*models.Item, int, error) {
	return i.FindItemsFn(ctx, filter)
}

func (i *ItemService) CreateItem(ctx context.Context, item *models.Item) error {
	return i.CreateItemFn(ctx, item)
}

func (i *ItemService) UpdateItem(ctx context.Context, id int, upd models.ItemUpdate) (*models.Item, error) {
	return i.UpdateItemFn(ctx, id, upd)
}

func (i *ItemService) DeleteItem(ctx context.Context, id int) error {
	return i.DeleteItemFn(ctx, id)
}
