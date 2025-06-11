package http

import (
	"invoice_project/internal/customer/domain"
	"time"
)

// CreateCustomerRequest represents the full payload to create a customer.
type CreateCustomerRequest struct {
	Customer       domain.Customer           `json:"customer"`
	Person         *domain.PersonCustomer    `json:"person,omitempty"`
	Company        *domain.CompanyCustomer   `json:"company,omitempty"`
	Address        *domain.CustomerAddress   `json:"address,omitempty"`
	Contacts       []domain.CustomerContact  `json:"contacts,omitempty"`
}

type UpdateCustomerRequest struct {
    Customer domain.Customer         `json:"customer"`
    Person   *domain.PersonCustomer  `json:"person,omitempty"`
    Company  *domain.CompanyCustomer `json:"company,omitempty"`
    Address  *domain.CustomerAddress `json:"address,omitempty"`
    Contacts []domain.CustomerContact `json:"contacts,omitempty"`
}

type PersonCustomerRequest struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	VatNo     string `json:"vat_no"`
}

type CompanyCustomerRequest struct {
	CompanyName string `json:"company_name" validate:"required"`
	VatNo       string `json:"vat_no"`
	BranchNo    int    `json:"branch_no"`
}

type CustomerAddressRequest struct {
	AddressLine1   string `json:"address_line1" validate:"required"`
	AddressLine2   string `json:"address_line2"`
	ProvinceID     int    `json:"province_id" validate:"required"`
	DistrictsID    int    `json:"districts_id" validate:"required"`
	SubdistrictsID int    `json:"subdistricts_id" validate:"required"`
	PostalCode     string `json:"postal_code" validate:"required"`
}

type CustomerContactRequest struct {
	ContactType  string    `json:"contact_type" validate:"required"`
	ContactValue string    `json:"contact_value" validate:"required"`
	CreatedAt    time.Time `json:"created_at"`
}
