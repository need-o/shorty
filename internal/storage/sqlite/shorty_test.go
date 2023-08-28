package sqlite

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/need-o/shorty/internal/migrate"
	"github.com/need-o/shorty/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestSqliteGetShorty(t *testing.T) {
	ctx := context.Background()

	t.Run("valid get shorty", func(t *testing.T) {
		deleteDB()
		s := NewShortyStorage(createDB())

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
		deleteDB()
		s := NewShortyStorage(createDB())

		_, err := s.GetShorty(ctx, "test")
		assert.ErrorIs(t, err, models.ErrShortyNotFound)
	})
}

func TestSqliteCreateShorty(t *testing.T) {
	ctx := context.Background()

	t.Run("create valid shorty", func(t *testing.T) {
		deleteDB()
		s := NewShortyStorage(createDB())

		input := models.Shorty{
			URL: "https://example.com",
		}

		err := s.CreateShorty(ctx, &input)
		assert.NoError(t, err)
		assert.NotNil(t, input.CreatedAt)
		assert.NotNil(t, input.UpdatedAt)
	})

	t.Run("create valid shorty with ID", func(t *testing.T) {
		deleteDB()
		s := NewShortyStorage(createDB())

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
		deleteDB()
		s := NewShortyStorage(createDB())

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

func TestSqliteCreateVisit(t *testing.T) {
	ctx := context.Background()

	t.Run("create valid visit", func(t *testing.T) {
		deleteDB()
		s := NewShortyStorage(createDB())

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

func TestSqliteGetVisits(t *testing.T) {
	ctx := context.Background()

	t.Run("get existing visits", func(t *testing.T) {
		deleteDB()
		s := NewShortyStorage(createDB())
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
		deleteDB()
		s := NewShortyStorage(createDB())

		visits, err := s.GetVisits(ctx, "not_existing")
		assert.NoError(t, err)

		assert.True(t, len(visits) == 0)
	})
}

func createDB() *sqlx.DB {
	db, err := sqlx.Open("sqlite3", "shorty_test.db")
	if err != nil {
		log.Fatal(err)
	}

	err = migrate.RunForSqlite3(db.DB, "file://../../../migrations")
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func deleteDB() {
	os.Remove("shorty_test.db")
}
