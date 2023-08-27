package memory

import (
	"context"
	"shorty/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	t.Run("valid get shorty", func(t *testing.T) {
		ctx := context.Background()
		s := NewShortyStorage()

		input := models.Shorty{
			ID:  "test",
			URL: "https://example.com",
		}

		err := s.Create(ctx, &input)
		assert.NoError(t, err)

		sh, err := s.Get(ctx, input.ID)
		assert.NoError(t, err)

		assert.Equal(t, sh.ID, input.ID)
		assert.Equal(t, sh.URL, input.URL)
		assert.NotNil(t, sh.CreatedAt)
		assert.NotNil(t, sh.UpdatedAt)
	})

	t.Run("not found get shorty", func(t *testing.T) {
		ctx := context.Background()
		s := NewShortyStorage()

		_, err := s.Get(ctx, "test")
		assert.ErrorIs(t, err, models.ErrShortyNotFound)
	})
}

func TestCreate(t *testing.T) {
	ctx := context.Background()

	t.Run("create valid shorty", func(t *testing.T) {
		s := NewShortyStorage()

		input := models.Shorty{
			URL: "https://example.com",
		}

		err := s.Create(ctx, &input)
		assert.NoError(t, err)
		assert.NotNil(t, input.CreatedAt)
		assert.NotNil(t, input.UpdatedAt)
	})

	t.Run("create valid shorty with ID", func(t *testing.T) {
		s := NewShortyStorage()

		input := models.Shorty{
			ID:  "test",
			URL: "https://example.com",
		}

		err := s.Create(ctx, &input)
		assert.NoError(t, err)
		assert.NotNil(t, input.CreatedAt)
		assert.NotNil(t, input.UpdatedAt)
	})

	t.Run("create invalid shorty with existing ID", func(t *testing.T) {
		s := NewShortyStorage()

		input := models.Shorty{
			ID:  "test",
			URL: "https://example.com",
		}

		err := s.Create(ctx, &input)
		assert.NoError(t, err)

		err = s.Create(ctx, &input)
		assert.ErrorIs(t, err, models.ErrShortyExists)
	})
}

func TestUpdate(t *testing.T) {
	ctx := context.Background()

	t.Run("update valid shorty", func(t *testing.T) {
		s := NewShortyStorage()

		input := models.Shorty{
			ID:  "test",
			URL: "https://example.com",
		}

		err := s.Create(ctx, &input)
		assert.NoError(t, err)

		input.URL = "https://ya.ru"
		err = s.Update(ctx, &input)
		assert.NoError(t, err)

		changed, err := s.Get(ctx, input.ID)
		assert.NoError(t, err)

		assert.Equal(t, changed.URL, input.URL)
	})

	t.Run("update not found shorty", func(t *testing.T) {
		s := NewShortyStorage()

		input := models.Shorty{
			ID:  "test",
			URL: "https://example.com",
		}

		err := s.Update(ctx, &input)
		assert.ErrorIs(t, err, models.ErrShortyNotFound)
	})
}

func TestCreateVisit(t *testing.T) {
	ctx := context.Background()

	t.Run("create valid visit", func(t *testing.T) {
		s := NewShortyStorage()

		input := models.Visit{
			ShortyID:  "test",
			Referer:   "https://example.com",
			UserIP:    "127.0.0.1",
			UserAgent: "",
		}

		err := s.CreateVisit(ctx, &input)
		assert.NoError(t, err)
		assert.NotNil(t, input.CreatedAt)
		assert.NotNil(t, input.UpdatedAt)
	})
}

func TestGetVisits(t *testing.T) {
	ctx := context.Background()

	t.Run("get existing visits", func(t *testing.T) {
		s := NewShortyStorage()

		input := models.Visit{
			ShortyID:  "test",
			Referer:   "https://example.com",
			UserIP:    "127.0.0.1",
			UserAgent: "",
		}

		err := s.CreateVisit(ctx, &input)
		assert.NoError(t, err)

		err = s.CreateVisit(ctx, &input)
		assert.NoError(t, err)

		err = s.CreateVisit(ctx, &input)
		assert.NoError(t, err)

		visits, err := s.GetVisits(ctx, input.ShortyID)
		assert.NoError(t, err)

		assert.True(t, len(visits) == 3)
	})

	t.Run("get not existing visits", func(t *testing.T) {
		s := NewShortyStorage()

		visits, err := s.GetVisits(ctx, "not_existing")
		assert.NoError(t, err)

		assert.True(t, len(visits) == 0)
	})
}
