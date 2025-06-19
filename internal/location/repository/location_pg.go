package repository

import (
	"context"
	"invoice_project/internal/location/domain"

	"gorm.io/gorm"
)

type LocationRepository interface {
	GetProvinceAll(ctx context.Context) ([]domain.Province, error)
	GetProvinceByID(ctx context.Context, id uint) (*domain.Province, error)
	GetDistrictById(ctx context.Context, id uint) (*domain.District, error)
	GetDistricts(ctx context.Context, id uint) ([]domain.District, error)
	GetSubDistrictsById(ctx context.Context, id uint) (*domain.SubDistrict, error)
	GetSubDistricts(ctx context.Context, id uint) ([]domain.SubDistrict, error)
}

type locationRepository struct {
	db *gorm.DB
}


func NewLocationRepository(db *gorm.DB) LocationRepository {
	return &locationRepository{db}
}


// province
func (r *locationRepository) GetProvinceAll(ctx context.Context) ([]domain.Province, error) {
	var provinces []domain.Province
	if err := r.db.WithContext(ctx).Find(&provinces).Error; err != nil {
		return nil, err
	}
	return provinces, nil
}

func (r *locationRepository) GetProvinceByID(ctx context.Context, id uint) (*domain.Province, error) {
	var province domain.Province
	if err := r.db.WithContext(ctx).First(&province, id).Error; err != nil {
		return nil, err
	}
	return &province, nil
}


// district
func (r *locationRepository) GetDistrictById(ctx context.Context, id uint) (*domain.District, error) {
	var District domain.District
	if err := r.db.WithContext(ctx).First(&District, id).Error; err != nil {
		return nil, err
	}
	return &District, nil
}

func (r *locationRepository) GetDistricts(ctx context.Context, id uint) ([]domain.District, error) {
	var districts []domain.District
	if err := r.db.WithContext(ctx).Where("province_id = ?", id).Find(&districts).Error; err != nil {
		return nil, err
	}
	return districts, nil
}

//  sub-district
func (r *locationRepository) GetSubDistricts(ctx context.Context, id uint) ([]domain.SubDistrict, error) {
	var sdistricts []domain.SubDistrict
	if err := r.db.WithContext(ctx).Where("district_id = ?", id).Find(&sdistricts).Error; err != nil {
		return nil, err
	}
	return sdistricts, nil
}

func (r *locationRepository) GetSubDistrictsById(ctx context.Context, id uint) (*domain.SubDistrict, error) {
	var sdistricts domain.SubDistrict
	if err := r.db.WithContext(ctx).First(&sdistricts, id).Error; err != nil {
		return nil, err
	}
	return &sdistricts, nil
}