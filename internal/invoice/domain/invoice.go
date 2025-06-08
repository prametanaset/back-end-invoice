package domain

import "time"

type Invoice struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Customer    string    `gorm:"not null" json:"customer"`
	Amount      float64   `gorm:"not null" json:"amount"`
	CreatedAt   time.Time `json:"created_at"`
	CreatedByID uint      `json:"created_by"`
}
