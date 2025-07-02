package domain

import "time"

// OTP stores a hashed OTP code and metadata for verification attempts.
type OTP struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement"`
	Purpose     string    `gorm:"type:text;not null"`
	Ref         string    `gorm:"type:text;not null"`
	Destination string    `gorm:"type:text;not null"`
	CodeHash    string    `gorm:"type:text;not null"`
	CreatedAt   time.Time `gorm:"not null;default:now()"`
	ExpiresAt   time.Time `gorm:"not null"`
	UsedAt      *time.Time
	RevokedAt   *time.Time
	Attempts    uint16 `gorm:"not null;default:0"`
}

func (OTP) TableName() string { return "otps" }
