package main

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"github.com/symaster1995/ms-starter/internal/config"
	"github.com/symaster1995/ms-starter/internal/http"
	productsDB "github.com/symaster1995/ms-starter/internal/products/database"
	postgres "github.com/symaster1995/ms-starter/pkg/database"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"time"
)

type Launcher struct {
	log        *zerolog.Logger
	store      postgres.DB
	httpServer *http.Server
	apiBackend *http.ApiBackend
}

func NewLauncher() *Launcher {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	log := zerolog.New(output).With().Timestamp().Logger()

	return &Launcher{
		log: &log,
	}
}

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	v := viper.New()
	o := config.NewConfig(v)

	c := make(chan os.Signal, 1) //make a channel to listen for errors
	signal.Notify(c, os.Interrupt)

	go func() { <-c; cancel() }() //goroutine to call cancel() signal if data is received from channel <-c

	m := NewLauncher()

	if err := m.run(o); err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}

	<-ctx.Done() //wait for cancel() signal

	if err := m.Shutdown(); err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
}

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

	//Create http Server
	if err := m.RunServer(cfg.ApiConfig); err != nil {
		return err
	}
	return nil
}

func (m *Launcher) RunServer(apiConfig *config.ApiConfig) error {
	m.httpServer = http.NewServer(apiConfig, m.log, m.apiBackend)
	return m.httpServer.Open()
}

func (m *Launcher) Shutdown() (err error) {

	if err := m.httpServer.Close(); err != nil { //Close http Server
		m.log.Error().Err(err).Msg("Failed to close http server")
		return err
	}
	return nil
}
