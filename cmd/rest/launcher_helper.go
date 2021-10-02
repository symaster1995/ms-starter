package main

import (
	"context"
	"github.com/symaster1995/ms-starter/cmd/rest/flags"
	"github.com/symaster1995/ms-starter/internal/http"
	"github.com/symaster1995/ms-starter/internal/mock"
	"testing"
	"time"
)

type TestLauncher struct {
	*Launcher
}

func NewTestLauncher() *TestLauncher {
	return &TestLauncher{
		NewLauncher(),
	}
}

func (tl *TestLauncher) RunOrFail(tb testing.TB, ctx context.Context) {
	if err := tl.Run(tb, ctx); err != nil {
		tb.Fatal(err)
	}
}

func (tl *TestLauncher) Run(tb testing.TB, ctx context.Context) error {

	apiOpts := &flags.ApiConfig{
		HttpBindAddress:       ":6969",
		HttpReadHeaderTimeout: 10 * time.Second,
		HttpReadTimeout:       1 * time.Second,
		HttpWriteTimeout:      1 * time.Second,
		Domain:                "",
	}

	tl.apiBackend = &http.ApiBackend{
		ItemService: mock.NewItemService(),
	}

	tl.Launcher.httpServer = http.NewServer(apiOpts, tl.log, tl.apiBackend)

	return tl.Launcher.httpServer.Open()
}

func (tl *TestLauncher) ShutdownOrFail(tb testing.TB) {
	tb.Helper()
	if err := tl.Launcher.Shutdown(); err != nil {
		tb.Fatal(err)
	}
}
