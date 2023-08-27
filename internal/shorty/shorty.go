package shorty

import (
	"context"
	"net/url"
	"shorty/internal/models"
	"strings"

	"github.com/google/uuid"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type (
	storage interface {
		Get(ctx context.Context, id string) (*models.Shorty, error)
		Create(ctx context.Context, sh *models.Shorty) error
		Update(ctx context.Context, sh *models.Shorty) error
		CreateVisit(ctx context.Context, visit *models.Visit) error
		GetVisits(ctx context.Context, shortyID string) ([]models.Visit, error)
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

func (s *Shorty) Get(ctx context.Context, id string) (*models.Shorty, error) {
	sh, err := s.storage.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	visits, err := s.storage.GetVisits(ctx, id)
	if err != nil {
		return nil, err
	}

	sh.Visits = visits

	return sh, nil
}

func (s *Shorty) Create(ctx context.Context, in models.ShortyInput) (*models.Shorty, error) {
	sh := models.Shorty{
		ID:  in.ID,
		URL: in.URL,
	}

	if sh.ID == "" {
		sh.ID = NewID(uuid.New().ID())
	}

	err := s.storage.Create(ctx, &sh)
	if err != nil {
		return &sh, err
	}

	return &sh, nil
}

func (s *Shorty) Redirect(ctx context.Context, v models.VisitInput) (*url.URL, error) {
	sh, err := s.storage.Get(ctx, v.ShortyID)
	if err != nil {
		return nil, err
	}

	url, err := url.Parse(sh.URL)
	if err != nil {
		return nil, err
	}

	visit := models.Visit{
		ShortyID:  v.ShortyID,
		Referer:   v.Referer,
		UserIP:    v.UserIP,
		UserAgent: v.UserAgent,
	}

	return url, s.storage.CreateVisit(ctx, &visit)
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
