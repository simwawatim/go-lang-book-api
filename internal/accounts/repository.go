package accounts

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/simwawatim/book-api/internal/core"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, u *User) error {
	err := r.db.QueryRow(ctx,
		`INSERT INTO users (name, email, password_hash)
		 VALUES ($1, $2, $3)
		 RETURNING id, created_at`,
		u.Name, u.Email, u.PasswordHash,
	).Scan(&u.ID, &u.CreatedAt)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return fmt.Errorf("email %w", core.ErrConflict)
	}
	return err
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (*User, error) {
	return r.getOne(ctx, `WHERE email = $1`, email)
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*User, error) {
	return r.getOne(ctx, `WHERE id = $1`, id)
}

func (r *Repository) getOne(ctx context.Context, where string, arg any) (*User, error) {
	var u User
	err := r.db.QueryRow(ctx,
		`SELECT id, name, email, password_hash, created_at FROM users `+where, arg,
	).Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, core.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}
