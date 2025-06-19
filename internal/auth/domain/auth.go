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

// UserLoginMethod links a user with an external auth provider.
type UserLoginMethod struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	UserID      uuid.UUID `gorm:"type:uuid;not null;index"`
	Provider    string    `gorm:"size:30;not null"`
	ProviderUID string    `gorm:"size:191;not null"`
	CreatedAt   time.Time
	User        User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

// Role defines a permission role.
type Role struct {
	ID          int     `gorm:"primaryKey;autoIncrement"`
	Name        string  `gorm:"size:30;unique;not null"`
	Description *string `gorm:"type:text"`
}

// UserRole assigns a role to a user.
type UserRole struct {
	ID     uint      `gorm:"primaryKey;autoIncrement"`
	UserID uuid.UUID `gorm:"type:uuid;not null;index"`
	RoleID int       `gorm:"not null;index"`
	User   User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Role   Role      `gorm:"foreignKey:RoleID;constraint:OnDelete:CASCADE"`
}
