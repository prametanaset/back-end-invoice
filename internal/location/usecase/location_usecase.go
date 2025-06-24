package usecase

import (
	"context"
	"invoice_project/internal/location/domain"
	"invoice_project/internal/location/repository"
)

type LocationUsecase interface {
	GetProvinceAll(ctx context.Context) ([]domain.Province, error)
	GetProvinceByID(ctx context.Context, id uint) (*domain.Province, error)
	GetProvinceByGeoID(ctx context.Context, geoID uint) ([]domain.Province, error)
	// District
	GetDistricts(ctx context.Context, id uint) ([]domain.District, error)
	GetDistrictById(ctx context.Context, id uint) (*domain.District, error)
	// Sub-District
	GetSubDistrictsById(ctx context.Context, id uint) (*domain.SubDistrict, error)
	GetSubDistricts(ctx context.Context, id uint) ([]domain.SubDistrict, error)
	GetZipCodeBySubDistrictID(ctx context.Context, id uint) (int, error)
}

type locationUsecase struct {
	repo repository.LocationRepository
}

func NewLocationUseCase(repo repository.LocationRepository) LocationUsecase {
	return &locationUsecase{repo}
}

func (u *locationUsecase) GetProvinceAll(ctx context.Context) ([]domain.Province, error) {
	return u.repo.GetProvinceAll(ctx)
}

func (u *locationUsecase) GetProvinceByID(ctx context.Context, id uint) (*domain.Province, error) {
	return u.repo.GetProvinceByID(ctx, id)
}

func (u *locationUsecase) GetProvinceByGeoID(ctx context.Context, geoID uint) ([]domain.Province, error) {
	return u.repo.GetProvinceByGeoID(ctx, geoID)
}

// District

func (u *locationUsecase) GetDistricts(ctx context.Context, id uint) ([]domain.District, error) {
	return u.repo.GetDistricts(ctx, id)
}

func (u *locationUsecase) GetDistrictById(ctx context.Context, id uint) (*domain.District, error) {
	return u.repo.GetDistrictById(ctx, id)
}

// Sub-District
func (u *locationUsecase) GetSubDistricts(ctx context.Context, id uint) ([]domain.SubDistrict, error) {
	return u.repo.GetSubDistricts(ctx, id)
}

func (u *locationUsecase) GetSubDistrictsById(ctx context.Context, id uint) (*domain.SubDistrict, error) {
	return u.repo.GetSubDistrictsById(ctx, id)
}

func (u *locationUsecase) GetZipCodeBySubDistrictID(ctx context.Context, id uint) (int, error) {
	return u.repo.GetZipCodeBySubDistrictID(ctx, id)
}
