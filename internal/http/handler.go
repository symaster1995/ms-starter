package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"github.com/symaster1995/ms-starter/internal/products"
	"net/http"
	"time"
)

type RootHandler struct {
	name   string
	router chi.Router
	log    *zerolog.Logger
}

func NewRootHandler(logger *zerolog.Logger, api *ApiBackend) *RootHandler {

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	r.Use(zeroLogger(logger))

	itemsHandler := products.NewItemHandler(logger, api.ItemService)

	r.Mount("/items", itemsHandler)

	return &RootHandler{
		router: r,
		name:   "RootHandler",
		log:    logger,
	}
}

func (rH *RootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rH.router.ServeHTTP(w, r)
}

func zeroLogger(l *zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()

			defer func() {
				l.Info().
					Str("method", r.Method).
					Str("url", r.URL.Path).
					Int("status", ww.Status()).
					Int("size", ww.BytesWritten()).
					Dur("duration", time.Since(t1)).
					Msg("")
			}()
			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}
