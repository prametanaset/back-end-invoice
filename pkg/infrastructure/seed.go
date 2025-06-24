package infrastructure

import (
	"errors"
	"log"

	"gorm.io/gorm"
	authModel "invoice_project/internal/auth/domain"
)

// SeedRoles inserts default roles into the database if they do not already exist.
func SeedRoles(db *gorm.DB) {
	defaultRoles := []authModel.Role{
		{Name: "user"},
		{Name: "admin"},
	}

	for _, r := range defaultRoles {
		var existing authModel.Role
		err := db.Where("name = ?", r.Name).First(&existing).Error
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				log.Printf("failed checking role %s: %v", r.Name, err)
				continue
			}
		} else {
			// role already exists
			continue
		}

		if err := db.Create(&r).Error; err != nil {
			log.Printf("failed seeding role %s: %v", r.Name, err)
		}
	}
}
