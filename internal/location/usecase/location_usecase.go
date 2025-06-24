package usecase

import (
	"context"
	"invoice_project/internal/location/domain"
	"invoice_project/internal/location/repository"
)

type LocationUsecase interface {
<<<<<<< HEAD
	GetProvinceAll(ctx context.Context) ([]domain.Province, error)
	GetProvinceByID(ctx context.Context, id uint) (*domain.Province, error)
	// District
	GetDistricts(ctx context.Context, id uint) ([]domain.District, error)
	GetDistrictById(ctx context.Context, id uint) (*domain.District, error)
	// Sub-District
	GetSubDistrictsById(ctx context.Context, id uint) (*domain.SubDistrict, error)
	GetSubDistricts(ctx context.Context, id uint) ([]domain.SubDistrict, error)
=======
	ListGeographies(ctx context.Context) ([]domain.Geography, error)
	ListProvinces(ctx context.Context, geoID int) ([]domain.Province, error)
	ListAmphures(ctx context.Context, provinceID int) ([]domain.Amphure, error)
	ListTambons(ctx context.Context, amphureID int) ([]domain.Tambon, error)
>>>>>>> fe689551f1395d734858e14c96b935403361a507
}

type locationUsecase struct {
	repo repository.LocationRepository
}

<<<<<<< HEAD
func NewLocationUseCase(repo repository.LocationRepository) LocationUsecase {
	return &locationUsecase{repo}
}

func (u *locationUsecase) GetProvinceAll(ctx context.Context) ([]domain.Province, error) {
	return u.repo.GetProvinceAll(ctx)
}

func (u *locationUsecase) GetProvinceByID(ctx context.Context, id uint) (*domain.Province, error) {
	return u.repo.GetProvinceByID(ctx, id)
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
=======
func NewLocationUsecase(repo repository.LocationRepository) LocationUsecase {
	return &locationUsecase{repo: repo}
}

func (uc *locationUsecase) ListGeographies(ctx context.Context) ([]domain.Geography, error) {
	return uc.repo.ListGeographies()
}

func (uc *locationUsecase) ListProvinces(ctx context.Context, geoID int) ([]domain.Province, error) {
	return uc.repo.ListProvinces(geoID)
}

func (uc *locationUsecase) ListAmphures(ctx context.Context, provinceID int) ([]domain.Amphure, error) {
	return uc.repo.ListAmphures(provinceID)
}

func (uc *locationUsecase) ListTambons(ctx context.Context, amphureID int) ([]domain.Tambon, error) {
	return uc.repo.ListTambons(amphureID)
}
>>>>>>> fe689551f1395d734858e14c96b935403361a507
