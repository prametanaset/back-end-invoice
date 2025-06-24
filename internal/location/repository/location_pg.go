package repository

import (
<<<<<<< HEAD
	"context"
=======
>>>>>>> fe689551f1395d734858e14c96b935403361a507
	"invoice_project/internal/location/domain"

	"gorm.io/gorm"
)

<<<<<<< HEAD
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
=======
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
>>>>>>> fe689551f1395d734858e14c96b935403361a507
		return nil, err
	}
	return provinces, nil
}

<<<<<<< HEAD
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
=======
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
>>>>>>> fe689551f1395d734858e14c96b935403361a507
