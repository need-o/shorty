package shorty

import (
	"context"
	"shorty/internal/models"
	"strings"

	"github.com/google/uuid"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type (
	storage interface {
		Get(ctx context.Context, id string) (*models.Shortening, error)
		Create(ctx context.Context, sh models.Shortening) error
		Update(ctx context.Context, sh models.Shortening) error
	}

	Shorty struct {
		storage storage
	}
)

func New(storage storage) *Shorty {
	return &Shorty{
		storage: storage,
	}
}

func (s *Shorty) Get(ctx context.Context, id string) (*models.Shortening, error) {
	return s.storage.Get(ctx, id)
}

func (s *Shorty) Create(ctx context.Context, in models.InputShortening) (string, error) {
	sh := models.Shortening{
		ID:  in.ID,
		URL: in.URL,
	}

	if sh.ID == "" {
		sh.ID = NewID(uuid.New().ID())
	}

	err := s.storage.Create(ctx, sh)
	if err != nil {
		return "", err
	}

	return sh.ID, nil
}

func (s *Shorty) Redirect(ctx context.Context, id string) (*models.RedirectShortening, error) {
	sh, err := s.storage.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	sh.Visits++

	return &models.RedirectShortening{URL: sh.URL}, s.storage.Update(ctx, *sh)
}

func NewID(number uint32) string {
	length := len(alphabet)
	var b strings.Builder

	b.Grow(10)
	for ; number > 0; number = number / uint32(length) {
		b.WriteByte(alphabet[(number % uint32(length))])
	}

	return b.String()
}