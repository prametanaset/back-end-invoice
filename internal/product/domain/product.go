package domain

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	StoreID   uuid.UUID `gorm:"type:uuid;not null" json:"store_id"`
	Sku       string    `gorm:"type:text;uniqueIndex" json:"sku"`
	Name      string    `gorm:"type:text" json:"name"`
	Price     float64   `gorm:"not null" json:"price"`      // float64 ถ้ามีทศนิยม
	VatType   string    `gorm:"not null;type:text" json:"vat_type"`
	VatRate   int       `gorm:"not null" json:"vat_rate"`   
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	ProductImage  *ProductImage    `gorm:"foreignKey:ProductID" json:"product_image,omitempty"`

}

type ProductImage struct {
	ID        uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	ProductID uint    `gorm:"not null" json:"product_id"`
	Product   Product `gorm:"foreignKey:ProductID;references:ID" json:"product"`
	Url       string  `gorm:"type:text" json:"url"`
}
