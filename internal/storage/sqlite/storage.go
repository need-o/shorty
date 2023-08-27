package sqlite

import "github.com/jmoiron/sqlx"

type Storage struct {
	Shorty *ShortyStorage
}

func NewStorage(db *sqlx.DB) *Storage {
	return &Storage{
		Shorty: NewShortyStorage(db),
	}
}
