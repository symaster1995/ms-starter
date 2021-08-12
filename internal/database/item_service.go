package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/symaster1995/ms-starter/internal/models"
	"github.com/symaster1995/ms-starter/pkg/database"
	"strings"
	"time"
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

func (i *ItemService) FindItems(ctx context.Context, filter models.ItemFilter) ([]*models.Item, int, error) {

	var args []interface{}
	where := []string{"1 = 1"}

	if v := filter.Name; v != nil && *v != "" {
		where, args = append(where, "name = $1"), append(args, *v)
	}

	var items []*models.Item

	if err := i.db.InitTx(ctx, pgx.ReadCommitted, func(tx pgx.Tx) error {
		rows, err := tx.Query(ctx, `
			SELECT * FROM Items
			WHERE `+strings.Join(where, " AND ")+
			` ORDER BY id ASC `+database.FormatLimitOffset(filter.Limit, filter.Offset), args...,
		)

		if err != nil {
			return err
		}

		defer rows.Close()

		for rows.Next() {
			var item models.Item
			if err := rows.Err(); err != nil {
				return fmt.Errorf("failed to iterate: %w", err)
			}

			if err := rows.Scan(&item.ID, &item.Name, &item.CreatedAt, &item.UpdatedAt); err != nil {
				return err
			}

			items = append(items, &item)
		}

		return nil

	}); err != nil {
		return nil, 0, err
	}

	return items, len(items), nil
}

func (i *ItemService) CreateItem(ctx context.Context, item *models.Item) error {
	item.CreatedAt = time.Now()
	item.UpdatedAt = item.CreatedAt

	if err := i.db.InitTx(ctx, pgx.ReadCommitted, func(tx pgx.Tx) error {
		err := tx.QueryRow(ctx, `INSERT INTO Items (name, created_at, updated_at) VALUES ($1, $2, $3) RETURNING id`,
			item.Name, item.CreatedAt, item.UpdatedAt,
		).Scan(&item.ID)

		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (i *ItemService) UpdateItem(ctx context.Context, id int, upd models.ItemUpdate) (*models.Item, error) {

	item, err := i.FindItemByID(ctx, id)

	if err != nil {
		return nil, err
	}

	item.UpdatedAt = time.Now()
	item.Name = upd.Name

	if err := i.db.InitTx(ctx, pgx.ReadCommitted, func(tx pgx.Tx) error {
		rows , err := tx.Exec(ctx, `UPDATE items SET name = $1, updated_at = $2 WHERE id = $3`, upd.Name, item.UpdatedAt, id)

		if err != nil {
			return err
		}

		if rows.RowsAffected() == 0 {
			return fmt.Errorf("no rows updated")
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return item, nil
}

func (i *ItemService) DeleteItem(ctx context.Context, id int) error {
	return nil
}
