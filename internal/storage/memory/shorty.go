package memory

import (
	"context"
	"shorty/internal/models"
	"sync"
)

type ShortyStorage struct {
	m sync.Map
}

func NewShortyStorage() *ShortyStorage {
	return &ShortyStorage{}
}

func (s *ShortyStorage) Get(ctx context.Context, id string) (*models.Shorty, error) {
	sh, ok := s.m.Load(id)
	if !ok {
		return nil, models.ErrShortyNotFound
	}

	result := sh.(*models.Shorty)

	return result, nil
}

func (s *ShortyStorage) Create(ctx context.Context, sh *models.Shorty) error {
	if _, ok := s.m.Load(sh.ID); ok {
		return models.ErrShortyExists
	}

	sh.BeforeCreate()
	s.m.Store(sh.ID, sh)

	return nil
}

func (s *ShortyStorage) Update(ctx context.Context, sh *models.Shorty) error {
	if _, ok := s.m.Load(sh.ID); !ok {
		return models.ErrShortyNotFound
	}

	sh.BeforeUpdate()
	s.m.Store(sh.ID, sh)

	return nil
}
