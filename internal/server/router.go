// Package server connects URLs to handlers (like Django's urls.py).
package server

import (
	"context"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/simwawatim/book-api/internal/accounts"
	"github.com/simwawatim/book-api/internal/books"
	"github.com/simwawatim/book-api/internal/core"
	"github.com/simwawatim/book-api/internal/docs"
)

func NewRouter(db *pgxpool.Pool, tokens *core.TokenManager, corsOrigin string) http.Handler {
	accountsH := accounts.NewHandler(accounts.NewService(accounts.NewRepository(db), tokens))
	booksH := books.NewHandler(books.NewService(books.NewRepository(db)))
	auth := tokens.RequireAuth

	mux := http.NewServeMux()

	// health checks
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		core.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	mux.HandleFunc("GET /readyz", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		if err := db.Ping(ctx); err != nil {
			core.Error(w, http.StatusServiceUnavailable, "database unavailable")
			return
		}
		core.JSON(w, http.StatusOK, map[string]string{"status": "ready"})
	})

	mux.HandleFunc("POST /api/auth/register", accountsH.Register)
	mux.HandleFunc("POST /api/auth/login", accountsH.Login)
	mux.HandleFunc("GET /api/auth/me", auth(accountsH.Me))
	// books
	mux.HandleFunc("GET /api/books", booksH.List)
	mux.HandleFunc("GET /api/books/{id}", booksH.Get)
	mux.HandleFunc("POST /api/books", auth(booksH.Create))
	mux.HandleFunc("PUT /api/books/{id}", auth(booksH.Update))
	mux.HandleFunc("DELETE /api/books/{id}", auth(booksH.Delete))
	// API docs: open http://localhost:8743/docs
	docs.Register(mux)
	return core.Recoverer(core.Logger(core.CORS(corsOrigin)(mux)))
}
