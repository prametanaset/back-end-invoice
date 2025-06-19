package domain

import (
	"time"

	"github.com/google/uuid"
)

type Merchant struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID       uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	MerchantType string    `gorm:"type:varchar(20);not null" json:"merchant_type"`
}

type Store struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	MerchantID uuid.UUID `gorm:"type:uuid;not null" json:"merchant_id"`
	StoreName  string    `gorm:"size:255;not null" json:"store_name"`
	BranchNo   string    `gorm:"size:10" json:"branch_no"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type StoreAddress struct {
	ID            uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	StoreID       uuid.UUID `gorm:"type:uuid;not null" json:"store_id"`
	AddressLine1  string    `gorm:"not null" json:"address_line1"`
	SubdistrictID int       `gorm:"not null" json:"subdistrict_id"`
	DistrictID    int       `gorm:"not null" json:"district_id"`
	ProvinceID    int       `gorm:"not null" json:"province_id"`
	PostalCode    string    `gorm:"size:10;not null" json:"postal_code"`
}

// MerchantContact stores contact info such as phone or email for a merchant.
type MerchantContact struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	MerchantID   uuid.UUID `gorm:"type:uuid;not null" json:"merchant_id"`
	ContactType  string    `gorm:"size:50;not null" json:"contact_type"`
	ContactValue string    `gorm:"not null" json:"contact_value"`
	CreatedAt    time.Time `json:"created_at"`
}

// PersonMerchant stores personal merchant details for merchants of type "person".
type PersonMerchant struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	MerchantID uuid.UUID `gorm:"type:uuid;not null" json:"merchant_id"`
	FirstName  string    `gorm:"size:100;not null" json:"first_name"`
	LastName   string    `gorm:"size:100;not null" json:"last_name"`
	VatNo      *string   `gorm:"size:20" json:"vat_no"`
}

// CompanyMerchant stores company merchant details for merchants of type "company".
type CompanyMerchant struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	MerchantID  uuid.UUID `gorm:"type:uuid;not null" json:"merchant_id"`
	CompanyName string    `gorm:"size:255;not null" json:"company_name"`
	VatNo       string    `gorm:"size:20;not null" json:"vat_no"`
}
