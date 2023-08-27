package models

import "time"

type (
	Visit struct {
		ShortyID  string     `json:"shorty_id" db:"shorty_id"`
		Referer   string     `json:"referer" db:"referer"`
		UserIP    string     `json:"user_ip" db:"user_ip"`
		UserAgent string     `json:"user_agent" db:"user_agent"`
		CreatedAt *time.Time `json:"created_at" db:"created_at"`
		UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	}

	VisitInput struct {
		ShortyID  string `json:"shorty_id"`
		Referer   string `json:"referer"`
		UserIP    string `json:"user_ip"`
		UserAgent string `json:"user_agent"`
	}
)

func (v *Visit) BeforeCreate() {
	now := time.Now()

	if v.CreatedAt == nil {
		v.CreatedAt = &now
	}

	if v.UpdatedAt == nil {
		v.UpdatedAt = &now
	}
}
