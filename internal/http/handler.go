package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/symaster1995/ms-starter/internal/models"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Handler struct {
	name        string
	router      chi.Router
	log         *zap.Logger
	slog *zap.SugaredLogger
	ItemService models.ItemService
}

func NewHandler() *Handler {

	return &Handler{
		name:        "Handler",
		log:         zap.NewNop(),
	}

	//h.ItemService = postgres.NewItemService()
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	h.router.ServeHTTP(w, r)
}

func (h *Handler) configureRouter(logger *zap.Logger) {

	h.log = logger.With(zap.String("service", "handler"))

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	//r.Use(middleware.Logger)
	r.Use(zapLogger(logger))

	r.Mount("/items", h.mountItemsRouter())

	h.router = r
}

// Wrap Go-Chi Middleware with Zap Logger

func zapLogger(l *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()
			defer func() {
				l.Info("Request",
					zap.String("method", r.Method),
					zap.String("path", r.URL.Path),
					zap.String("proto", r.Proto),
					zap.Int("status", ww.Status()),
					zap.Duration("lat", time.Since(t1)),
				)
			}()

			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}
