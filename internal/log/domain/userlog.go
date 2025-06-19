package domain

import (
	"time"

	"github.com/google/uuid"
)

type UserLog struct {
	ID         uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID     *uuid.UUID `gorm:"type:uuid;index"`
	Username   string
	IPAddress  string `gorm:"size:45"`
	Action     string
	Resource   string
	DeviceInfo string
	StartedAt  time.Time
	Status     int
}
