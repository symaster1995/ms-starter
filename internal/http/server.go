package http

import (
	"context"
	"fmt"
	"github.com/caddyserver/certmagic"
	"github.com/rs/zerolog"
	"github.com/symaster1995/ms-starter/internal/config"
	"net"
	"net/http"
	"os"
	"runtime"
	"time"
)

type Server struct {
	log        *zerolog.Logger
	Addr       string
	Listener   net.Listener
	Domain     string
	HttpServer *http.Server
}

func NewServer(apiConfig *config.ApiConfig, logger *zerolog.Logger, api *ApiBackend) *Server {

	handler := NewRootHandler(logger, api.ItemService)

	httpServer := &http.Server{
		ReadHeaderTimeout: apiConfig.HttpReadHeaderTimeout,
		ReadTimeout:       apiConfig.HttpReadTimeout,
		WriteTimeout:      apiConfig.HttpWriteTimeout,
		Handler:           handler,
	}

	return &Server{
		Addr:       apiConfig.HttpBindAddress,
		Domain:     apiConfig.Domain,
		HttpServer: httpServer,
		log:        logger,
	}
}

func (s *Server) Open() (err error) {

	if s.UseTLS() {
		s.Listener, err = certmagic.Listen([]string{s.Domain})
	}

	s.Listener, err = net.Listen("tcp", s.Addr)

	if err != nil {
		s.log.Error().Err(err).Msg("Failed to set up TCP listener")
		return err
	}

	go func(log *zerolog.Logger) {
		log.Debug().Str("address", s.Addr).Str("OS", runtime.GOOS).Str("Arch", runtime.GOARCH).Msg("Server Listening")
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

func (s *Server) UseTLS() bool {
	return s.Domain != ""
}

func (s *Server) Scheme() string {
	if s.UseTLS() {
		return "https"
	}
	return "http"
}

func (s *Server) Port() int {
	if s.Listener == nil {
		return 0
	}
	return s.Listener.Addr().(*net.TCPAddr).Port
}

func (s *Server) URL() string {
	scheme, port := s.Scheme(), s.Port()

	// Use localhost unless a domain is specified.
	domain := "localhost"
	if s.Domain != "" {
		domain = s.Domain
	}

	// Return without port if using standard ports.
	if (scheme == "http" && port == 80) || (scheme == "https" && port == 443) {
		return fmt.Sprintf("%s://%s", s.Scheme(), domain)
	}
	return fmt.Sprintf("%s://%s:%d", s.Scheme(), domain, s.Port())
}
