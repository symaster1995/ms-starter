package database

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/symaster1995/ms-starter/internal/models"
	"github.com/symaster1995/ms-starter/pkg/database"
)

type ItemService struct {
	db *database.DB
}

func NewItemService(db *database.DB) *ItemService {
	return &ItemService{db: db}
}

func (i *ItemService) FindItemByID(ctx context.Context, id int) (*models.Item, error) {

	var item models.Item

	if err := i.db.InitTx(ctx, pgx.ReadCommitted, func(tx pgx.Tx) error {

		err := tx.QueryRow(ctx, `SELECT * FROM Items WHERE ID = $1`, id).Scan(&item.ID, &item.Name, &item.CreatedAt, &item.UpdatedAt)

		if err != nil {
			return err
		}
		return nil

	}); err != nil {
		return nil, err
	}

	return &item, nil
}

func (i *ItemService) FindItem(ctx context.Context, filter models.ItemFilter) (*models.Item, error) {
	return nil, nil
}

func (i *ItemService) FindItems(ctx context.Context, filter models.ItemFilter) ([]*models.Item, int, error) {
	return nil, 0, nil
}

func (i *ItemService) CreateItem(ctx context.Context, item *models.Item) error {
	return nil
}

func (i *ItemService) UpdateItem(ctx context.Context, id int, upd models.ItemUpdate) (*models.Item, error) {
	return nil, nil
}

func (i *ItemService) DeleteItem(ctx context.Context, id int) error {
	return nil
}
