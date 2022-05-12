MS-STARTER ![test](https://github.com/symaster1995/ms-starter/actions/workflows/test.yml/badge.svg)
========

A REST API boilerplate with logging and configuration.

## Application Domain
This project demonstrates a real-world inventory application that follows Domain Driven Design concepts and mainly uses the `product item` as a model.

It also uses interface for basic CRUD (Create/ Read/ Update/ Delete) operations and abides the inward dependency rule.

This way we can easily test or change other database/caching layers without as they share a common interface.

## Project Structure
This project uses but not limited to the standard project layout structure.
The code is organized with the following approach:

1. `api` - This folder contains API Details and Specs.
2. `internal` - This folder contains all domains and packages that are private and exclusive to your application.
3. `pkg` - This folder contains public packages - `database` (database configurations), `errors` (custom errors), etc.
4. `cmd` - This folder contains main executable subpackages that ties the project together - `cmd/rest`.

## How to use
### Domain Structure

The domains containing your exclusive business applications are located inside the internal folder.
The structure for domain packages are as follows:

1. `database` - This folder contains the database services for the domain package.
2. `models` - This folder contains the model entity for the package.
3. `domain_handler` - This file contains the logic for handling all requests with the associated domain.

### Tieing all the code
After creating a domain and fullfiling it's required files, we need to add some additional code to make it work.

In here we are adding domain's model to the collection of services `ApiBackend` struct in the internal http package:
```go
package http

import products "github.com/symaster1995/ms-starter/internal/products/models"

type ApiBackend struct {
	ItemService products.ItemService
}
```

Next is creating the domain's handler and adding it to the root handler in the internal http package:
```go
package http

import "github.com/symaster1995/ms-starter/internal/products"

func NewRootHandler(logger *zerolog.Logger, api *ApiBackend) *RootHandler {
	r := chi.NewRouter()
	itemsHandler := products.NewItemHandler(logger, api.ItemService)
	r.Mount("/items", itemsHandler)
}
```

Lastly we need to create a database service for the domain and add it to an instance of `ApiBackend` struct in the main package:
```go
package main

func (m *Launcher) run(cfg *config.Config) error {
	//Create new db instance
	db, err := postgres.NewDB(cfg.DBConfig.URL, m.log)
	if err != nil {
		m.log.Error().Err(err).Msg("Failed to create connection pool")
		return err
	}

	//Create item service
	itemService := productsDB.NewItemService(db)

	//Collection of services for easier integration
	m.apiBackend = &http.ApiBackend{
		ItemService: itemService,
	}
}
```

### Used Packages

--------------------------------------------------------------------------------------------------------------------
| Package                                               | Usage                                                    |
| :-----------------------------------------------------|:----------------------------------------------------------
| [chi](https://github.com/go-chi/chi)                  | Routing                                                  |
| [pgx](https://github.com/jackc/pgx/v4)                | Go driver and toolkit for PostgreSQL                     |
| [pgconn](https://github.com/jackc/pgconn)             | Checks PostgreSQL specific errors                        |
| [pgerrcode](https://github.com/jackc/pgerrcode)       | Constants for PostgreSQL error codes                     |
| [zerolog](https://github.com/rs/zerolog)              | Allows syntax logging                                    |
| [viper](https://github.com/spf13/viper)               | Configuration solution for Go applications               |
| [certmagic](https://github.com/caddyserver/certmagic) | Automates the obtaining and renewal of TLS certificates  |
| [go-cmp](https://github.com/google/go-cmp)            | Compares two values semantically for easier testing      |
--------------------------------------------------------------------------------------------------------------------
