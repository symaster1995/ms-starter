package models

import (
	"context"
	"time"
)

type ItemService interface {
	FindItemByID(ctx context.Context, id int) (*Item, error)

	FindItem(ctx context.Context, filter ItemFilter) (*Item, error)

	FindItems(ctx context.Context, filter ItemFilter) ([]*Item, int, error)

	CreateItem(ctx context.Context, i *Item) error

	UpdateItem(ctx context.Context, id int, upd ItemUpdate) (*Item, error)

	DeleteItem(ctx context.Context, id int) error
}

type Item struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ItemFilter struct {
	ID   int    `json:"id"`
	Name string `json:"name"`

	Limit      int `json:"limit"`
	Offset     int `json:"offset"`

}

type ItemUpdate struct {
	Name string `json:"name"`
}