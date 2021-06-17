package main

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"github.com/symaster1995/ms-starter/cmd/rest/flags"
	"github.com/symaster1995/ms-starter/internal/http"
	"github.com/symaster1995/ms-starter/internal/postgres"
	"go.uber.org/zap"
	_ "net/http/pprof"
	"os"
	"os/signal"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	v := viper.New()
	o := flags.NewOpts(v)

	c := make(chan os.Signal, 1) //make a channel to listen for errors
	signal.Notify(c, os.Interrupt)

	go func() { <-c; cancel() }() //goroutine to call cancel() signal if data is received from channel <-c

	m := NewLauncher()

	if err := m.run(o); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	<-ctx.Done() //wait for cancel() signal

	if err := m.Shutdown(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type Launcher struct {
	log        *zap.Logger
	store      postgres.DB
	httpServer *http.Server
}

func NewLauncher() *Launcher {
	return &Launcher{
		log: zap.NewNop(),
	}
}

func (m *Launcher) run(opts *flags.ApiOpts) (err error) {
	m.log, _ = zap.NewDevelopment() //Initialize logger from zap
	defer m.log.Sync()
	m.httpServer = http.NewServer(opts, m.log) //Create http Server

	if err := m.httpServer.Open(m.log); err != nil { //Start http Server
		return err
	}
	return nil
}

func (m *Launcher) Shutdown() (err error) {

	if err := m.httpServer.Close(); err != nil { //Close http Server
		m.log.Error("Failed to close http server", zap.Error(err))
		return err
	}
	return nil
}
