package domain

import (
	"time"

	"github.com/google/uuid"
)

type Invoice struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Customer    string    `gorm:"not null" json:"customer"`
	Amount      float64   `gorm:"not null" json:"amount"`
	CreatedAt   time.Time `json:"created_at"`
	CreatedByID uuid.UUID `json:"created_by"`
}
