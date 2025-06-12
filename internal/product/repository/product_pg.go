package repository

import (
	"context"
	"invoice_project/internal/product/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)
type ProductRepository interface {
	CreateProduct(ctx context.Context, product *domain.Product, images []domain.ProductImage) error
	GetProduct(ctx context.Context, id uint) (*domain.Product, []domain.ProductImage, error)
	ListProducts(ctx context.Context, storeID uuid.UUID) ([]domain.Product, error)
	UpdateProduct(ctx context.Context, product *domain.Product, images []domain.ProductImage) error
	DeleteProduct(ctx context.Context, id uint) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *productRepository {
	return &productRepository{db: db}
}

// CreateProduct สร้างสินค้าใหม่พร้อมแนบรูปภาพ
func (r *productRepository) CreateProduct(ctx context.Context, product *domain.Product, productimg []domain.ProductImage) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(product).Error; err != nil {
			return err
		}

		for i := range productimg {
			productimg[i].ProductID = product.ID
		}

		if len(productimg) > 0 {
			if err := tx.Create(&productimg).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// GetProduct ดึงสินค้าและรูปภาพทั้งหมดจาก ID
func (r *productRepository) GetProduct(ctx context.Context, id uint) (*domain.Product, []domain.ProductImage, error) {
	var product domain.Product
	if err := r.db.WithContext(ctx).First(&product, id).Error; err != nil {
		return nil, nil, err
	}

	var images []domain.ProductImage
	if err := r.db.WithContext(ctx).
		Where("product_id = ?", id).
		Find(&images).Error; err != nil {
		return &product, nil, err
	}

	return &product, images, nil
}

// ListProducts ดึงสินค้าทั้งหมดในร้าน
func (r *productRepository) ListProducts(ctx context.Context, storeID uuid.UUID) ([]domain.Product, error) {
	var products []domain.Product
	err := r.db.WithContext(ctx).
		Where("store_id = ?", storeID).
		Find(&products).Error
	return products, err
}

// UpdateProduct อัปเดตสินค้าและรูปภาพใหม่ทั้งหมด
func (r *productRepository) UpdateProduct(ctx context.Context, product *domain.Product, images []domain.ProductImage) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// อัปเดตสินค้า
		if err := tx.Save(product).Error; err != nil {
			return err
		}

		// ลบรูปเดิมทั้งหมด
		if err := tx.Where("product_id = ?", product.ID).Delete(&domain.ProductImage{}).Error; err != nil {
			return err
		}

		// เพิ่มรูปใหม่
		for i := range images {
			images[i].ProductID = product.ID
		}
		if len(images) > 0 {
			if err := tx.Create(&images).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// DeleteProduct ลบสินค้าและรูปทั้งหมด
func (r *productRepository) DeleteProduct(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("product_id = ?", id).Delete(&domain.ProductImage{}).Error; err != nil {
			return err
		}
		if err := tx.Delete(&domain.Product{}, id).Error; err != nil {
			return err
		}
		return nil
	})
}