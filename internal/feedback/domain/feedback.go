package domain

import (
	"time"

	"github.com/google/uuid"
)

type Feedback struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Score    	int64    `gorm:"not null" json:"score"`
	Comment      string   `gorm:"not null" json:"comment"`
	CreatedAt   time.Time `json:"created_at"`
	CreatedByID uuid.UUID `json:"created_by"`
}