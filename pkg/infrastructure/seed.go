package infrastructure

import (
	"errors"
	"log"

	"gorm.io/gorm"
	authModel "invoice_project/internal/auth/domain"
	merchModel "invoice_project/internal/merchant/domain"
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

// SeedMerchantTypes inserts default merchant types into the database if they do not already exist.
func SeedMerchantTypes(db *gorm.DB) {
	defaultTypes := []merchModel.MerchantType{
		{Name: "person"},
		{Name: "company"},
	}

	for _, t := range defaultTypes {
		var existing merchModel.MerchantType
		err := db.Where("name = ?", t.Name).First(&existing).Error
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				log.Printf("failed checking merchant type %s: %v", t.Name, err)
				continue
			}
		} else {
			// type already exists
			continue
		}

		if err := db.Create(&t).Error; err != nil {
			log.Printf("failed seeding merchant type %s: %v", t.Name, err)
		}
	}
}
