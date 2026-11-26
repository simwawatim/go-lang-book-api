// Package core holds code shared by every domain package:
// errors, JSON helpers, JWT and middleware.
package core

import "errors"

var (
	ErrNotFound     = errors.New("not found")
	ErrConflict     = errors.New("already exists")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
)

// ValidationError carries per-field messages and is returned as HTTP 422.
type ValidationError struct {
	Fields map[string]string
}

func (e *ValidationError) Error() string { return "validation failed" }

// Validator collects field errors; returns nil from Err() when there are none.
type Validator struct {
	fields map[string]string
}

func (v *Validator) Check(ok bool, field, msg string) {
	if ok {
		return
	}
	if v.fields == nil {
		v.fields = map[string]string{}
	}
	if _, exists := v.fields[field]; !exists {
		v.fields[field] = msg
	}
}

func (v *Validator) Err() error {
	if len(v.fields) == 0 {
		return nil
	}
	return &ValidationError{Fields: v.fields}
}
