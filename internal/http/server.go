package http

import (
	"context"
	"github.com/caddyserver/certmagic"
	"github.com/rs/zerolog"
	"github.com/symaster1995/ms-starter/cmd/rest/flags"
	"net"
	"net/http"
	"os"
	"time"
)

type Server struct {
	log        *zerolog.Logger
	Addr       string
	Listener   net.Listener
	Domain     string
	HttpServer *http.Server
}

func NewServer(opts *flags.ApiConfig, logger *zerolog.Logger, api *ApiBackend) *Server {

	handler := NewHandler(logger)
	handler.configureRouter()
	handler.ItemService = api.ItemService

	httpServer := &http.Server{
		Addr:              opts.HttpBindAddress,
		ReadHeaderTimeout: opts.HttpReadHeaderTimeout,
		ReadTimeout:       opts.HttpReadTimeout,
		WriteTimeout:      opts.HttpWriteTimeout,
		Handler:           handler,
	}

	return &Server{
		Addr:       opts.HttpBindAddress,
		Domain:     opts.Domain,
		HttpServer: httpServer,
		log:        logger,
	}
}

func (s *Server) Open() (err error) {

	if s.Domain != "" {
		s.Listener, err = certmagic.Listen([]string{s.Domain})
	}

	s.Listener, err = net.Listen("tcp", s.Addr)

	if err != nil {
		s.log.Error().Err(err).Msg("Failed to set up TCP listener")
		return err
	}

	go func(log *zerolog.Logger) {
		log.Debug().Str("address", s.Addr).Msg("Server Listening")
		if err := s.HttpServer.Serve(s.Listener); err != http.ErrServerClosed {
			log.Error().Err(err).Msg("Failed to serve HTTP")
			os.Exit(1)
		}
		log.Info().Msg("Stopping")
	}(s.log)
	return nil
}

func (s *Server) Close() error {
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	return s.HttpServer.Shutdown(shutdownCtx)
}
