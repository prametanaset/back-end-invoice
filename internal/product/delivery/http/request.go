package http

import (
	"invoice_project/internal/product/domain"
)

type CreateProductRequest struct {
	Product 		domain.Product  	   `json:"product"`
	ProductImage    []domain.ProductImage  `json:"product_image,omitempty"` // ส่งเป็น URL
}
