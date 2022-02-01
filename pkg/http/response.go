package http

import (
	"encoding/json"
	"github.com/symaster1995/ms-starter/pkg/errors"
	"log"
	"net/http"
)

//ErrInvalid:            400,
//ErrUnauthorized:       401,
//ErrNotFound:           404,
//ErrConflict:           409,
//ErrInternal:           500,
//ErrNotImplemented:     501,
//ErrServiceUnavailable: 503,

func RenderJSON(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		LogError(err)
		return
	}
	return
}

func ErrorJSON(w http.ResponseWriter, err error) {
	code, message := errors.ErrorCode(err), errors.ErrorMessage(err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(ErrorStatusCode(code))

	if err := json.NewEncoder(w).Encode(&ErrorResponse{Error: message}); err != nil {
		LogError(err)
		return
	}
	return
}

func LogError(err error) {
	log.Printf("[http] error: %s", err)
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func ErrorStatusCode(code string) int {
	if v, ok := codes[code]; ok {
		return v
	}
	return http.StatusInternalServerError
}

var codes = map[string]int{
	errors.ErrInvalid:            http.StatusBadRequest,
	errors.ErrUnauthorized:       http.StatusUnauthorized,
	errors.ErrNotFound:           http.StatusNotFound,
	errors.ErrConflict:           http.StatusConflict,
	errors.ErrInternal:           http.StatusInternalServerError,
	errors.ErrNotImplemented:     http.StatusNotImplemented,
	errors.ErrServiceUnavailable: http.StatusInternalServerError,
}
