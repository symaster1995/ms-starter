package http

import (
	"encoding/json"
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

	if code > 300 {
		if typ, ok := data.(error); ok {
			data = &ErrorResponse{Error: typ.Error()}
		}
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
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
