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

type ShortyStorage struct {
	db *sqlx.DB
}

func NewShortyStorage(db *sqlx.DB) *ShortyStorage {
	return &ShortyStorage{
		db: db,
	}
}

func (s *ShortyStorage) Get(ctx context.Context, id string) (*models.Shorty, error) {
	sh := models.Shorty{}

	err := s.db.GetContext(ctx, &sh,
		`SELECT * FROM shorty WHERE id=$1`, id,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, models.ErrShortyNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("get shorty error: %w", err)
	}

	return &sh, nil
}

func (s *ShortyStorage) Create(ctx context.Context, sh *models.Shorty) error {
	sh.BeforeCreate()

	_, err := s.db.NamedExecContext(ctx,
		`INSERT INTO shorty (id, url, visits, created_at, updated_at) 
		 VALUES (:id, :url, :visits, :created_at, :updated_at);`, &sh,
	)

	var sqliteErr sqlite3.Error

	if errors.As(err, &sqliteErr) {
		if errors.Is(sqliteErr.Code, sqlite3.ErrConstraint) {
			return models.ErrShortyExists
		}
	}

	if err != nil {
		return fmt.Errorf("create shorty error: %w", err)
	}

	return nil
}

func (s *ShortyStorage) Update(ctx context.Context, sh *models.Shorty) error {
	sh.BeforeUpdate()

	_, err := s.db.NamedExecContext(ctx,
		`UPDATE shorty SET url=:url, visits=:visits, created_at=:created_at, updated_at=:updated_at
		 WHERE id=:id;`, &sh,
	)
	if err != nil {
		return fmt.Errorf("update shorty error: %w", err)
	}

	return nil
}