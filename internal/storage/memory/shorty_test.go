package memory

import (
	"context"
	"testing"

	"github.com/need-o/shorty/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestGetShorty(t *testing.T) {
	t.Run("valid get shorty", func(t *testing.T) {
		ctx := context.Background()
		s := NewShortyStorage()

		input := models.Shorty{
			ID:  "test",
			URL: "https://example.com",
		}

		err := s.CreateShorty(ctx, &input)
		assert.NoError(t, err)

		sh, err := s.GetShorty(ctx, input.ID)
		assert.NoError(t, err)

		assert.Equal(t, sh.ID, input.ID)
		assert.Equal(t, sh.URL, input.URL)
		assert.NotNil(t, sh.CreatedAt)
		assert.NotNil(t, sh.UpdatedAt)
	})

	t.Run("not found get shorty", func(t *testing.T) {
		ctx := context.Background()
		s := NewShortyStorage()

		_, err := s.GetShorty(ctx, "test")
		assert.ErrorIs(t, err, models.ErrShortyNotFound)
	})
}

func TestCreateShorty(t *testing.T) {
	ctx := context.Background()

	t.Run("create valid shorty", func(t *testing.T) {
		s := NewShortyStorage()

		input := models.Shorty{
			URL: "https://example.com",
		}

		err := s.CreateShorty(ctx, &input)
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

		err := s.CreateShorty(ctx, &input)
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

		err := s.CreateShorty(ctx, &input)
		assert.NoError(t, err)

		err = s.CreateShorty(ctx, &input)
		assert.ErrorIs(t, err, models.ErrShortyExists)
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
		count := 10

		input := models.Visit{
			ShortyID:  "test",
			Referer:   "https://example.com",
			UserIP:    "127.0.0.1",
			UserAgent: "",
		}

		for i := 0; i < count; i++ {
			err := s.CreateVisit(ctx, &input)
			assert.NoError(t, err)
		}

		visits, err := s.GetVisits(ctx, input.ShortyID)
		assert.NoError(t, err)

		assert.True(t, len(visits) == count)
	})

	t.Run("get not existing visits", func(t *testing.T) {
		s := NewShortyStorage()

		visits, err := s.GetVisits(ctx, "not_existing")
		assert.NoError(t, err)

		assert.True(t, len(visits) == 0)
	})
}
