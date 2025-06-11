package domain

import (
	"time"

	"github.com/google/uuid"
)

type Customer struct {
	ID           uint         `gorm:"primaryKey;autoIncrement" json:"id"`
	StoreID      uuid.UUID    `gorm:"type:uuid;not null" json:"store_id"`
	CustomerType string       `gorm:"type:text" json:"customer_type"`
	Status       string       `gorm:"type:text" json:"status"`
	CreatedAt    time.Time    `json:"created_at"`
	CreatedBy    string       `gorm:"type:text" json:"created_by"`
	UpdatedAt    time.Time    `json:"updated_at"`
	UpdatedBy    string       `gorm:"type:text" json:"updated_by"`

	// ✅ Add these relationships
	CompanyCustomer  *CompanyCustomer    `gorm:"foreignKey:CustomerID" json:"company_customer,omitempty"`
	PersonCustomer   *PersonCustomer     `gorm:"foreignKey:CustomerID" json:"person_customer,omitempty"`
	CustomerAddress  *CustomerAddress    `gorm:"foreignKey:CustomerID" json:"customer_address,omitempty"`
	CustomerContact  []CustomerContact   `gorm:"foreignKey:CustomerID" json:"customer_contacts,omitempty"`
}

type CompanyCustomer struct {
	ID           uint     `gorm:"primaryKey;autoIncrement" json:"id"`
	CustomerID   uint     `gorm:"not null" json:"customer_id"`
	Customer     Customer `gorm:"foreignKey:CustomerID;references:ID" json:"customer"`
	CompanyName  string   `gorm:"type:text" json:"company_name"`
	VatNo        string   `gorm:"type:text" json:"vat_no"`
	BranchNo     int      `json:"branch_no"`
}

type PersonCustomer struct {
	ID         uint     `gorm:"primaryKey;autoIncrement" json:"id"`
	CustomerID   uint     `gorm:"not null" json:"customer_id"`
	Customer   Customer `gorm:"foreignKey:CustomerID;references:ID" json:"customer"`
	FirstName  string   `gorm:"type:text" json:"first_name"`
	LastName   string   `gorm:"type:text" json:"last_name"`
	VatNo      string   `gorm:"type:text" json:"vat_no"`
}

type CustomerAddress struct {
	ID             uint     `gorm:"primaryKey;autoIncrement" json:"id"`
	CustomerID   uint     `gorm:"not null" json:"customer_id"`
	Customer       Customer `gorm:"foreignKey:CustomerID;references:ID" json:"customer"`
	AddressLine1   string   `gorm:"type:text" json:"address_line1"`
	AddressLine2   string   `gorm:"type:text" json:"address_line2"`
	ProvinceID     int      `gorm:"not null" json:"province_id"`
	DistrictsID    int      `gorm:"not null" json:"districts_id"`
	SubdistrictsID int      `gorm:"not null" json:"subdistricts_id"`
	PostalCode     string   `gorm:"size:10;not null" json:"postal_code"`
}

type CustomerContact struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CustomerID   uint     `gorm:"not null" json:"customer_id"`
	Customer     Customer  `gorm:"foreignKey:CustomerID;references:ID" json:"customer"`
	ContactType  string    `gorm:"type:text" json:"contact_type"`
	ContactValue string    `gorm:"type:text" json:"contact_value"`
	CreatedAt    time.Time `json:"created_at"`
}
