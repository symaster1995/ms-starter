package http

import products "github.com/symaster1995/ms-starter/internal/products/models"

type ApiBackend struct {
	ItemService products.ItemService
}
