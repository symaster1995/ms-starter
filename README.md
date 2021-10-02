MS-STARTER ![test](https://github.com/symaster1995/ms-starter/actions/workflows/test.yml/badge.svg)
========

A REST API boilerplate with logging and configuration.

## Application Domain
This project demonstrates a real-world inventory application and mainly uses the `item` as a model.

It also uses interface for basic CRUD (Create/ Read/ Update/ Delete) operations and abides the inward dependency rule.

This way we can easily test or change other database/caching layers without as they share a common interface.

## Project Structure
This project uses but not limited to the standard project layout structure.
The code is organized with the following approach:

1. `api` - This folder contains API Details and Specs.
2. `internal` - This folder contains all packages that are private and exclusive to your application.
3. `pkg` - This folder contains public packages - `database` (database configurations), `errors` (custom errors), etc.
4. `cmd` - This folder contains main executable subpackages that ties the project together - `cmd/rest`.

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