package usecase

import (
	"context"
	"invoice_project/internal/product/domain"
	"invoice_project/internal/product/repository"
	"invoice_project/pkg/apperror"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ProductUseCase interface {
	CreateProduct(ctx context.Context, product *domain.Product, images []domain.ProductImage) error
	GetProduct(ctx context.Context, id uint) (*domain.Product, []domain.ProductImage, error)
	ListProducts(ctx context.Context, storeID uuid.UUID) ([]domain.Product, error)
	UpdateProduct(ctx context.Context, product *domain.Product, images []domain.ProductImage) error
	DeleteProduct(ctx context.Context, id uint) error
}


type productUseCase struct {
	repo repository.ProductRepository
}

func NewProductUseCase(repo repository.ProductRepository) ProductUseCase {
	return &productUseCase{
		repo: repo,
	}
}

func (uc *productUseCase) CreateProduct(ctx context.Context, product *domain.Product, productimg []domain.ProductImage) error {
	if product == nil {
		return apperror.New(fiber.StatusBadRequest)
	}
	product.ID = 0 // ป้องกัน override
	return uc.repo.CreateProduct(ctx, product, productimg)
}

// ดึงสินค้าโดย ID พร้อมรูป
func (uc *productUseCase) GetProduct(ctx context.Context, id uint) (*domain.Product, []domain.ProductImage, error) {
	return uc.repo.GetProduct(ctx, id)
}

// ดึงสินค้าทั้งหมดของร้าน
func (uc *productUseCase) ListProducts(ctx context.Context, storeID uuid.UUID) ([]domain.Product, error) {
	return uc.repo.ListProducts(ctx, storeID)
}

// อัปเดตสินค้า
func (uc *productUseCase) UpdateProduct(ctx context.Context, product *domain.Product, images []domain.ProductImage) error {
	if product == nil || product.ID == 0 {
		return apperror.New(fiber.StatusBadRequest)
	}
	return uc.repo.UpdateProduct(ctx, product, images)
}

// ลบสินค้า
func (uc *productUseCase) DeleteProduct(ctx context.Context, id uint) error {
	if id == 0 {
		return apperror.New(fiber.StatusBadRequest)
	}
	return uc.repo.DeleteProduct(ctx, id)
}
