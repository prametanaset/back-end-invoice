package repository

import (
	"errors"

	"invoice_project/internal/invoice/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InvoiceRepository interface {
	CreateInvoice(inv *domain.Invoice) error
	GetInvoiceByID(id uuid.UUID) (*domain.Invoice, error)
	ListInvoicesByUser(userID uuid.UUID) ([]domain.Invoice, error)
}

type invoicePG struct {
	db *gorm.DB
}

func NewInvoiceRepository(db *gorm.DB) InvoiceRepository {
	return &invoicePG{db: db}
}

func (r *invoicePG) CreateInvoice(inv *domain.Invoice) error {
	return r.db.Create(inv).Error
}

func (r *invoicePG) GetInvoiceByID(id uuid.UUID) (*domain.Invoice, error) {
	var inv domain.Invoice
	err := r.db.First(&inv, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &inv, nil
}

func (r *invoicePG) ListInvoicesByUser(userID uuid.UUID) ([]domain.Invoice, error) {
	var invoices []domain.Invoice
	err := r.db.Where("created_by_id = ?", userID).Order("created_at desc").Find(&invoices).Error
	if err != nil {
		return nil, err
	}
	return invoices, nil
}
