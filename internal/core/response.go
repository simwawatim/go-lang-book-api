package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

const maxBodyBytes = 1 << 20 // 1 MB

func JSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if v != nil {
		_ = json.NewEncoder(w).Encode(v)
	}
}

func Error(w http.ResponseWriter, status int, msg string) {
	JSON(w, status, map[string]string{"error": msg})
}

// Decode parses a JSON body into dst and rejects unknown fields.
func Decode(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return fmt.Errorf("invalid JSON body: %w", err)
	}
	return nil
}

// WriteError maps service errors to HTTP status codes.
func WriteError(w http.ResponseWriter, err error) {
	var ve *ValidationError
	switch {
	case errors.As(err, &ve):
		JSON(w, http.StatusUnprocessableEntity, map[string]any{
			"error":  ve.Error(),
			"fields": ve.Fields,
		})
	case errors.Is(err, ErrNotFound):
		Error(w, http.StatusNotFound, "resource not found")
	case errors.Is(err, ErrConflict):
		Error(w, http.StatusConflict, err.Error())
	case errors.Is(err, ErrUnauthorized):
		Error(w, http.StatusUnauthorized, err.Error())
	case errors.Is(err, ErrForbidden):
		Error(w, http.StatusForbidden, "you do not have permission to do this")
	default:
		log.Printf("internal error: %v", err)
		Error(w, http.StatusInternalServerError, "internal server error")
	}
}
