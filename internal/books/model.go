// Package books implements CRUD for books. Anyone can read;
// only logged-in users can create, and only the owner can edit or delete.
package books

import (
	"strings"
	"time"
	"unicode/utf8"

	"github.com/simwawatim/book-api/internal/core"
)

type Book struct {
	ID            int64     `json:"id"`
	Title         string    `json:"title"`
	Author        string    `json:"author"`
	ISBN          string    `json:"isbn"`
	PublishedYear *int      `json:"published_year"`
	Description   string    `json:"description"`
	OwnerID       int64     `json:"owner_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// BookInput is the request body for create (POST) and update (PUT).
type BookInput struct {
	Title         string `json:"title"`
	Author        string `json:"author"`
	ISBN          string `json:"isbn"`
	PublishedYear *int   `json:"published_year"`
	Description   string `json:"description"`
}

func (in *BookInput) Normalize() {
	in.Title = strings.TrimSpace(in.Title)
	in.Author = strings.TrimSpace(in.Author)
	in.ISBN = strings.ReplaceAll(strings.TrimSpace(in.ISBN), "-", "")
	in.Description = strings.TrimSpace(in.Description)
}

func (in *BookInput) Validate() error {
	var v core.Validator
	v.Check(in.Title != "", "title", "is required")
	v.Check(utf8.RuneCountInString(in.Title) <= 255, "title", "must be at most 255 characters")
	v.Check(in.Author != "", "author", "is required")
	v.Check(utf8.RuneCountInString(in.Author) <= 255, "author", "must be at most 255 characters")
	v.Check(in.ISBN == "" || len(in.ISBN) == 10 || len(in.ISBN) == 13, "isbn", "must be 10 or 13 characters")
	if in.PublishedYear != nil {
		y := *in.PublishedYear
		v.Check(y >= 0 && y <= time.Now().Year()+1, "published_year", "is not a valid year")
	}
	v.Check(utf8.RuneCountInString(in.Description) <= 5000, "description", "must be at most 5000 characters")
	return v.Err()
}

type ListParams struct {
	Search string
	Limit  int
	Offset int
}

type ListResult struct {
	Data   []Book `json:"data"`
	Total  int    `json:"total"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}
