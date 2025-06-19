package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents application user credentials.
type User struct {
	ID           uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Username     string         `gorm:"unique;not null" json:"username"`
	PasswordHash string         `gorm:"not null" json:"-"`
	Role         string         `gorm:"type:varchar(30);not null" json:"role"`
	IsVerified   bool           `gorm:"default:false" json:"is_verified"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// RefreshToken เก็บ refresh token ของแต่ละ user
// UserSession stores refresh tokens and session info for each user.
type UserSession struct {
	ID           uint      `gorm:"primaryKey;autoIncrement"`
	UserID       uuid.UUID `gorm:"type:uuid;not null;index"`
	UserAgent    string    `json:"user_agent"`
	IPAddress    string    `gorm:"type:inet" json:"ip_address"`
	RefreshToken string    `gorm:"unique;not null"`
	ExpiresAt    time.Time `gorm:"not null"`
	Revoked      bool      `gorm:"default:false"`
	CreatedAt    time.Time `json:"created_at"`
	User         User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}
