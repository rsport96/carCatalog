package models

import "github.com/gofrs/uuid"

type Car struct {
	ID     uuid.UUID `json:"id" db:"id"`
	RegNum string    `json:"regNum" db:"regNum"`
	Mark   string    `json:"mark" db:"mark"`
	Model  string    `json:"model" db:"model"`
	Year   int       `json:"year,omitempty" db:"year,omitempty"`
	Owner  People    `json:"owner" db:"owner"`
}
