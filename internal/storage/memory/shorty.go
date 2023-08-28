package memory

import (
	"context"
	"github.com/need-o/shorty/internal/models"
	"sync"
)

type ShortyStorage struct {
	shorty sync.Map
	visits []models.Visit
	m      sync.Mutex
}

func NewShortyStorage() *ShortyStorage {
	return &ShortyStorage{}
}

func (s *ShortyStorage) GetShorty(ctx context.Context, id string) (*models.Shorty, error) {
	sh, ok := s.shorty.Load(id)
	if !ok {
		return nil, models.ErrShortyNotFound
	}

	visits := []models.Visit{}
	for _, v := range s.visits {
		if v.ShortyID == id {
			visits = append(visits, v)
		}
	}

	result := sh.(*models.Shorty)
	result.Visits = visits

	return result, nil
}

func (s *ShortyStorage) CreateShorty(ctx context.Context, sh *models.Shorty) error {
	if _, ok := s.shorty.Load(sh.ID); ok {
		return models.ErrShortyExists
	}

	sh.BeforeCreate()
	s.shorty.Store(sh.ID, sh)

	return nil
}

func (s *ShortyStorage) CreateVisit(ctx context.Context, visit *models.Visit) error {
	s.m.Lock()
	defer s.m.Unlock()

	visit.BeforeCreate()
	s.visits = append(s.visits, *visit)

	return nil
}

func (s *ShortyStorage) GetVisits(ctx context.Context, shortyID string) ([]models.Visit, error) {
	s.m.Lock()
	defer s.m.Unlock()

	visits := []models.Visit{}

	for _, v := range s.visits {
		if v.ShortyID == shortyID {
			visits = append(visits, v)
		}
	}

	return visits, nil
}
