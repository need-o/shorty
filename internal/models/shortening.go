package models

import "time"

type (
	Shortening struct {
		ID        string     `json:"id" db:"id"`
		URL       string     `json:"url" db:"url"`
		Visits    int64      `json:"visits" db:"visits"`
		CreatedAt *time.Time `json:"created_at" db:"created_at"`
		UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	}

	ShortyInput struct {
		ID  string `json:"id"`
		URL string `json:"url"`
	}
)

func (s *Shortening) BeforeCreate() {
	now := time.Now()

	if s.CreatedAt == nil {
		s.CreatedAt = &now
	}

	if s.UpdatedAt == nil {
		s.UpdatedAt = &now
	}
}

func (s *Shortening) BeforeUpdate() {
	now := time.Now()

	if s.UpdatedAt == nil {
		s.UpdatedAt = &now
	}
}
