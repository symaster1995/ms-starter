package models

import (
	"context"
	"time"
)

type ItemService interface {

	//FindItemByID returns a single Item by ID
	FindItemByID(ctx context.Context, id int) (*Item, error)

	// FindItem returns a first Item by filter
	FindItem(ctx context.Context, filter ItemFilter) (*Item, error)

	// FindItems returns a list of Items matched by filter
	FindItems(ctx context.Context, filter ItemFilter) ([]*Item, int, error)

	// CreateItem creates new item
	CreateItem(ctx context.Context, i *Item) error

	//UpdateItem updates new item
	//returns the new item state after changes are applied
	UpdateItem(ctx context.Context, id int, upd ItemUpdate) (*Item, error)

	//DeleteItem removes an item by id
	DeleteItem(ctx context.Context, id int) error
}

type Item struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ItemFilter struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`

	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}

type ItemUpdate struct {
	Name string `json:"name"`
}
