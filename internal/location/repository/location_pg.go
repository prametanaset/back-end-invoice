package repository

import (
	"invoice_project/internal/location/domain"

	"gorm.io/gorm"
)

// LocationRepository provides access to thai location tables
// Minimal read-only functions

type LocationRepository interface {
	ListGeographies() ([]domain.Geography, error)
	ListProvinces(geoID int) ([]domain.Province, error)
	ListAmphures(provinceID int) ([]domain.Amphure, error)
	ListTambons(amphureID int) ([]domain.Tambon, error)
}

type locationPG struct {
	db *gorm.DB
}

func NewLocationRepository(db *gorm.DB) LocationRepository {
	return &locationPG{db: db}
}

func (r *locationPG) ListGeographies() ([]domain.Geography, error) {
	var geos []domain.Geography
	if err := r.db.Order("id").Find(&geos).Error; err != nil {
		return nil, err
	}
	return geos, nil
}

func (r *locationPG) ListProvinces(geoID int) ([]domain.Province, error) {
	var provinces []domain.Province
	q := r.db.Order("id")
	if geoID != 0 {
		q = q.Where("geography_id = ?", geoID)
	}
	if err := q.Find(&provinces).Error; err != nil {
		return nil, err
	}
	return provinces, nil
}

func (r *locationPG) ListAmphures(provinceID int) ([]domain.Amphure, error) {
	var amphures []domain.Amphure
	q := r.db.Order("id")
	if provinceID != 0 {
		q = q.Where("province_id = ?", provinceID)
	}
	if err := q.Find(&amphures).Error; err != nil {
		return nil, err
	}
	return amphures, nil
}

func (r *locationPG) ListTambons(amphureID int) ([]domain.Tambon, error) {
	var tambons []domain.Tambon
	q := r.db.Order("id")
	if amphureID != 0 {
		q = q.Where("amphure_id = ?", amphureID)
	}
	if err := q.Find(&tambons).Error; err != nil {
		return nil, err
	}
	return tambons, nil
}
