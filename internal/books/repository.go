package books

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/simwawatim/book-api/internal/core"
)

const bookColumns = `id, title, author, isbn, published_year, description, owner_id, created_at, updated_at`

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func scanBook(row pgx.Row, b *Book) error {
	return row.Scan(&b.ID, &b.Title, &b.Author, &b.ISBN, &b.PublishedYear,
		&b.Description, &b.OwnerID, &b.CreatedAt, &b.UpdatedAt)
}

func (r *Repository) List(ctx context.Context, p ListParams) ([]Book, int, error) {
	var total int
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM books
		 WHERE $1 = '' OR title ILIKE '%' || $1 || '%' OR author ILIKE '%' || $1 || '%'`,
		p.Search,
	).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(ctx,
		`SELECT `+bookColumns+` FROM books
		 WHERE $1 = '' OR title ILIKE '%' || $1 || '%' OR author ILIKE '%' || $1 || '%'
		 ORDER BY id DESC
		 LIMIT $2 OFFSET $3`,
		p.Search, p.Limit, p.Offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	list := make([]Book, 0, p.Limit)
	for rows.Next() {
		var b Book
		if err := scanBook(rows, &b); err != nil {
			return nil, 0, err
		}
		list = append(list, b)
	}
	return list, total, rows.Err()
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*Book, error) {
	var b Book
	err := scanBook(r.db.QueryRow(ctx, `SELECT `+bookColumns+` FROM books WHERE id = $1`, id), &b)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, core.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *Repository) Create(ctx context.Context, in BookInput, ownerID int64) (*Book, error) {
	var b Book
	err := scanBook(r.db.QueryRow(ctx,
		`INSERT INTO books (title, author, isbn, published_year, description, owner_id)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING `+bookColumns,
		in.Title, in.Author, in.ISBN, in.PublishedYear, in.Description, ownerID,
	), &b)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *Repository) Update(ctx context.Context, id int64, in BookInput) (*Book, error) {
	var b Book
	err := scanBook(r.db.QueryRow(ctx,
		`UPDATE books
		 SET title = $2, author = $3, isbn = $4, published_year = $5,
		     description = $6, updated_at = NOW()
		 WHERE id = $1
		 RETURNING `+bookColumns,
		id, in.Title, in.Author, in.ISBN, in.PublishedYear, in.Description,
	), &b)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, core.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM books WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return core.ErrNotFound
	}
	return nil
}
