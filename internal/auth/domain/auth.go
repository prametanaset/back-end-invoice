package domain

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Username  string    `gorm:"unique;not null" json:"username"`
	Password  string    `gorm:"not null" json:"-"`
	Role      string    `gorm:"type:varchar(20);default:'user'" json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// RefreshToken เก็บ refresh token ของแต่ละ user
type RefreshToken struct {
	ID        uint           `gorm:"primaryKey;autoIncrement"`
	Token     string         `gorm:"unique;not null"`
	UserID    uint           `gorm:"not null;index"`
	User      User           `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	ExpiredAt time.Time      `gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	DeleteAt  gorm.DeletedAt `gorm:"index" json:"delete_at"`
}
