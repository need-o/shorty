package memory

import "github.com/jmoiron/sqlx"

type Storage struct {
	Shortenings *ShorteningStorage
}

func NewStorage(db *sqlx.DB) *Storage {
	return &Storage{
		Shortenings: NewShorteningStorage(db),
	}
}
