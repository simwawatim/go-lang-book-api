package books

import (
	"context"

	"github.com/simwawatim/book-api/internal/core"
)

const (
	defaultLimit = 20
	maxLimit     = 100
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) List(ctx context.Context, p ListParams) (*ListResult, error) {
	if p.Limit <= 0 {
		p.Limit = defaultLimit
	}
	if p.Limit > maxLimit {
		p.Limit = maxLimit
	}
	if p.Offset < 0 {
		p.Offset = 0
	}

	books, total, err := s.repo.List(ctx, p)
	if err != nil {
		return nil, err
	}
	return &ListResult{Data: books, Total: total, Limit: p.Limit, Offset: p.Offset}, nil
}

func (s *Service) Get(ctx context.Context, id int64) (*Book, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) Create(ctx context.Context, userID int64, in BookInput) (*Book, error) {
	in.Normalize()
	if err := in.Validate(); err != nil {
		return nil, err
	}
	return s.repo.Create(ctx, in, userID)
}

func (s *Service) Update(ctx context.Context, userID, id int64, in BookInput) (*Book, error) {
	in.Normalize()
	if err := in.Validate(); err != nil {
		return nil, err
	}
	if err := s.checkOwner(ctx, userID, id); err != nil {
		return nil, err
	}
	return s.repo.Update(ctx, id, in)
}

func (s *Service) Delete(ctx context.Context, userID, id int64) error {
	if err := s.checkOwner(ctx, userID, id); err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
}

func (s *Service) checkOwner(ctx context.Context, userID, bookID int64) error {
	book, err := s.repo.GetByID(ctx, bookID)
	if err != nil {
		return err
	}
	if book.OwnerID != userID {
		return core.ErrForbidden
	}
	return nil
}
