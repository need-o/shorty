package memory

import (
	"context"
	"shorty/internal/models"
	"sync"

	"github.com/jmoiron/sqlx"
)

type ShorteningStorage struct {
	m sync.Map
}

func NewShorteningStorage(db *sqlx.DB) *ShorteningStorage {
	return &ShorteningStorage{}
}

func (s *ShorteningStorage) Get(ctx context.Context, id string) (*models.Shortening, error) {
	sh, ok := s.m.Load(id)
	if !ok {
		return nil, models.ErrShorteningNotFound
	}

	result := sh.(models.Shortening)

	return &result, nil
}

func (s *ShorteningStorage) Create(ctx context.Context, sh models.Shortening) error {
	sh.BeforeCreate()

	return nil
}

func (s *ShorteningStorage) Update(ctx context.Context, sh models.Shortening) error {

	return nil
}
