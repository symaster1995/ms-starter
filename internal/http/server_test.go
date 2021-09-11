package http_test

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/symaster1995/ms-starter/cmd/rest/flags"
	msHttp "github.com/symaster1995/ms-starter/internal/http"
	"github.com/symaster1995/ms-starter/internal/mock"
	"io"
	"net/http"
	"os"
	"testing"
	"time"
)

type Server struct {
	*msHttp.Server
	ItemService mock.ItemService
}

func MustOpenServer(tb testing.TB) *Server {
	tb.Helper()

	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	log := zerolog.New(output).With().Timestamp().Logger()

	apiOpts := &flags.ApiConfig{
		HttpBindAddress:       ":69",
		HttpReadHeaderTimeout: 10 * time.Second,
		HttpReadTimeout:       1 * time.Second,
		HttpWriteTimeout:      1 * time.Second,
		Domain:                "",
	}

	apiBackend := &msHttp.ApiBackend{
		ItemService: mock.NewItemService(),
	}

	s := &Server{Server: msHttp.NewServer(apiOpts, &log, apiBackend)}

	// Begin running test server.
	if err := s.Open(); err != nil {
		tb.Fatal(err)
	}

	return s
}

func MustCloseServer(tb testing.TB, s *Server) {
	tb.Helper()
	if err := s.Close(); err != nil {
		tb.Fatal(err)
	}
}

func (s *Server) MustNewRequest(tb testing.TB, ctx context.Context, method, url string, body io.Reader) *http.Request {
	tb.Helper()
	// Create new net/http request with server's base URL.
	r, err := http.NewRequest(method, s.URL()+url, body)
	if err != nil {
		tb.Fatal(err)
	}

	return r
}
