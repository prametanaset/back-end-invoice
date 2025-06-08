package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"

	"invoice_project/internal/log/domain"
)

func Logger(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		// ดึงรหัสสถานะหลังจาก Handler ทำงานจบ
		status := c.Response().StatusCode()

		ip := c.IP()
		userAgent := c.Get("User-Agent")
		url := c.OriginalURL()
		var userIDPtr *uint
		if uid, ok := c.Locals("user_id").(uint); ok {
			userIDPtr = &uid
		}
		log.Warn(userIDPtr)
		var username string
		if name, ok := c.Locals("username").(string); ok {
			username = name
		}

		log := domain.UserLog{
			UserID:     userIDPtr,
			Username:   username,
			IPAddress:  ip,
			Action:     c.Method(),
			Resource:   url,
			DeviceInfo: userAgent,
			StartedAt:  start,
			Status:     status,
		}

		go func() {
			_ = db.Create(&log).Error
		}()

		return err
	}

}
