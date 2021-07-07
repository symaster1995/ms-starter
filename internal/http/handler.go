package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"github.com/symaster1995/ms-starter/internal/models"
	"net/http"
	"time"
)

type Handler struct {
	name        string
	router      chi.Router
	log         *zerolog.Logger
	ItemService models.ItemService
}

func NewHandler(logger *zerolog.Logger) *Handler {

	return &Handler{
		name: "Handler",
		log:  logger,
	}

	//h.ItemService = postgres.NewItemService()
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	h.router.ServeHTTP(w, r)
}

func (h *Handler) configureRouter() {

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	//r.Use(middleware.Logger)

	r.Use(zeroLogger(h.log))

	r.Mount("/items", h.mountItemsRouter())
	h.router = r
}

func zeroLogger(l *zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()

			defer func(){
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
//
//func zapLogger(l *zap.Logger) func(next http.Handler) http.Handler {
//	return func(next http.Handler) http.Handler {
//		fn := func(w http.ResponseWriter, r *http.Request) {
//			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
//
//			t1 := time.Now()
//			defer func() {
//				l.Info("Request",
//					zap.String("method", r.Method),
//					zap.String("path", r.URL.Path),
//					zap.String("proto", r.Proto),
//					zap.Int("status", ww.Status()),
//					zap.Duration("lat", time.Since(t1)),
//				)
//			}()
//
//			next.ServeHTTP(ww, r)
//		}
//		return http.HandlerFunc(fn)
//	}
//}
