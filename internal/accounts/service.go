package accounts

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"

	"github.com/simwawatim/book-api/internal/core"
)

var errInvalidCredentials = fmt.Errorf("%w: invalid email or password", core.ErrUnauthorized)

type Service struct {
	repo   *Repository
	tokens *core.TokenManager
}

func NewService(repo *Repository, tokens *core.TokenManager) *Service {
	return &Service{repo: repo, tokens: tokens}
}

func (s *Service) Register(ctx context.Context, in RegisterInput) (*AuthResponse, error) {
	in.Name = strings.TrimSpace(in.Name)
	in.Email = strings.ToLower(strings.TrimSpace(in.Email))

	var v core.Validator
	v.Check(in.Name != "", "name", "is required")
	v.Check(utf8.RuneCountInString(in.Name) <= 100, "name", "must be at most 100 characters")
	v.Check(in.Email != "", "email", "is required")
	v.Check(isEmail(in.Email), "email", "must be a valid email address")
	v.Check(utf8.RuneCountInString(in.Password) >= 8, "password", "must be at least 8 characters")
	v.Check(len(in.Password) <= 72, "password", "must be at most 72 bytes") // bcrypt limit
	if err := v.Err(); err != nil {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &User{Name: in.Name, Email: in.Email, PasswordHash: string(hash)}
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}
	return s.issueToken(user)
}

func (s *Service) Login(ctx context.Context, in LoginInput) (*AuthResponse, error) {
	email := strings.ToLower(strings.TrimSpace(in.Email))

	user, err := s.repo.GetByEmail(ctx, email)
	if errors.Is(err, core.ErrNotFound) {
		return nil, errInvalidCredentials
	}
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(in.Password)); err != nil {
		return nil, errInvalidCredentials
	}
	return s.issueToken(user)
}

func (s *Service) Me(ctx context.Context, userID int64) (*User, error) {
	return s.repo.GetByID(ctx, userID)
}

func (s *Service) issueToken(u *User) (*AuthResponse, error) {
	token, exp, err := s.tokens.Generate(u.ID)
	if err != nil {
		return nil, err
	}
	return &AuthResponse{Token: token, ExpiresAt: exp, User: u}, nil
}

func isEmail(s string) bool {
	addr, err := mail.ParseAddress(s)
	return err == nil && addr.Address == s
}
