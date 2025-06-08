package usecase

import (
	"errors"

	"invoice_project/internal/invoice/domain"
	"invoice_project/internal/invoice/repository"
)

type InvoiceUsecase interface {
	CreateInvoice(customer string, amount float64, createdBy uint) (*domain.Invoice, error)
	GetInvoice(id uint, userID uint) (*domain.Invoice, error)
	ListInvoices(userID uint) ([]domain.Invoice, error)
}

type invoiceUC struct {
	repo repository.InvoiceRepository
}

func NewInvoiceUsecase(repo repository.InvoiceRepository) InvoiceUsecase {
	return &invoiceUC{repo: repo}
}

func (u *invoiceUC) CreateInvoice(customer string, amount float64, createdBy uint) (*domain.Invoice, error) {
	inv := &domain.Invoice{
		Customer:    customer,
		Amount:      amount,
		CreatedByID: createdBy,
	}
	err := u.repo.CreateInvoice(inv)
	if err != nil {
		return nil, err
	}
	return inv, nil
}

func (u *invoiceUC) GetInvoice(id uint, userID uint) (*domain.Invoice, error) {
	inv, err := u.repo.GetInvoiceByID(id)
	if err != nil {
		return nil, err
	}
	if inv == nil || inv.CreatedByID != userID {
		return nil, errors.New("invoice not found or unauthorized")
	}
	return inv, nil
}

func (u *invoiceUC) ListInvoices(userID uint) ([]domain.Invoice, error) {
	return u.repo.ListInvoicesByUser(userID)
}
