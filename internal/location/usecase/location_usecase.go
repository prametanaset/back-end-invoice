package usecase

import (
	"context"
	"invoice_project/internal/location/domain"
	"invoice_project/internal/location/repository"
)

type LocationUsecase interface {
	ListGeographies(ctx context.Context) ([]domain.Geography, error)
	ListProvinces(ctx context.Context, geoID int) ([]domain.Province, error)
	ListAmphures(ctx context.Context, provinceID int) ([]domain.Amphure, error)
	ListTambons(ctx context.Context, amphureID int) ([]domain.Tambon, error)
}

type locationUsecase struct {
	repo repository.LocationRepository
}

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
