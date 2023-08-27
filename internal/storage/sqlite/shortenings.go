package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"shorty/internal/models"

	"github.com/jmoiron/sqlx"
	"github.com/mattn/go-sqlite3"
)

type ShorteningStorage struct {
	db *sqlx.DB
}

func NewShorteningStorage(db *sqlx.DB) *ShorteningStorage {
	return &ShorteningStorage{
		db: db,
	}
}

func (s *ShorteningStorage) Get(ctx context.Context, id string) (*models.Shortening, error) {
	sh := models.Shortening{}

	err := s.db.GetContext(ctx, &sh,
		`SELECT * FROM shortenings WHERE id=$1`, id,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, models.ErrShorteningNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("get shortenings error: %w", err)
	}

	return &sh, nil
}

func (s *ShorteningStorage) Create(ctx context.Context, sh *models.Shortening) error {
	sh.BeforeCreate()

	_, err := s.db.NamedExecContext(ctx,
		`INSERT INTO shortenings (id, url, visits, created_at, updated_at) 
		 VALUES (:id, :url, :visits, :created_at, :updated_at);`, &sh,
	)

	var sqliteErr sqlite3.Error

	if errors.As(err, &sqliteErr) {
		if errors.Is(sqliteErr.Code, sqlite3.ErrConstraint) {
			return models.ErrShorteningExists
		}
	}

	if err != nil {
		return fmt.Errorf("create shortenings error: %w", err)
	}

	return nil
}

func (s *ShorteningStorage) Update(ctx context.Context, sh *models.Shortening) error {
	sh.BeforeUpdate()

	_, err := s.db.NamedExecContext(ctx,
		`UPDATE shortenings SET url=:url, visits=:visits, created_at=:created_at, updated_at=:updated_at
		 WHERE id=:id;`, &sh,
	)
	if err != nil {
		return fmt.Errorf("update shortenings error: %w", err)
	}

	return nil
}
