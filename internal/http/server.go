package http

import (
	"context"
	"github.com/caddyserver/certmagic"
	"github.com/symaster1995/ms-starter/cmd/rest/flags"
	"go.uber.org/zap"
	"net"
	"net/http"
	"os"
	"time"
)

type Server struct {
	logger     *zap.Logger
	Addr       string
	Listener   net.Listener
	Domain     string
	httpServer *http.Server
}

func NewServer(opts *flags.ApiOpts, logger *zap.Logger) *Server {

	httpLogger := logger.With(zap.String("service", "http"))

	httpServer := &http.Server{
		Addr:              opts.HttpBindAddress,
		ReadHeaderTimeout: opts.HttpReadHeaderTimeout,
		ReadTimeout:       opts.HttpReadTimeout,
		WriteTimeout:      opts.HttpWriteTimeout,
		ErrorLog:          zap.NewStdLog(httpLogger),
	}

	return &Server{
		Addr:       opts.HttpBindAddress,
		Domain:     opts.Domain,
		httpServer: httpServer,
	}
}

func (s *Server) Open(logger *zap.Logger) (err error) {
	log := logger.With(zap.String("service", "tcp-listener"))

	if s.Domain != "" {
		s.Listener, err = certmagic.Listen([]string{s.Domain})
	}
	s.Listener, err = net.Listen("tcp", s.Addr)

	if err != nil {
		log.Error("Failed to set up TCP listener", zap.String("addr", s.Addr), zap.Error(err))
		return err
	}

	go func(log *zap.Logger) {
		log.Info("Listening", zap.String("transport", "http"), zap.String("addr", s.Addr))
		if err := s.httpServer.Serve(s.Listener); err != http.ErrServerClosed {
			log.Error("Failed to serve HTTP", zap.Error(err))
			os.Exit(1)
		}
		log.Info("Stopping")
	}(log)
	return nil
}

func (s *Server) Close() error {
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	return s.httpServer.Shutdown(shutdownCtx)
}
