package usecase

import (
	"invoice_project/internal/invoice/domain"
	"invoice_project/internal/invoice/repository"
	"invoice_project/pkg/apperror"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type InvoiceUsecase interface {
	CreateInvoice(customer string, amount float64, createdBy uuid.UUID) (*domain.Invoice, error)
	GetInvoice(id uuid.UUID, userID uuid.UUID) (*domain.Invoice, error)
	ListInvoices(userID uuid.UUID) ([]domain.Invoice, error)
}

type invoiceUC struct {
	repo repository.InvoiceRepository
}

func NewInvoiceUsecase(repo repository.InvoiceRepository) InvoiceUsecase {
	return &invoiceUC{repo: repo}
}

func (u *invoiceUC) CreateInvoice(customer string, amount float64, createdBy uuid.UUID) (*domain.Invoice, error) {
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

func (u *invoiceUC) GetInvoice(id uuid.UUID, userID uuid.UUID) (*domain.Invoice, error) {
	inv, err := u.repo.GetInvoiceByID(id)
	if err != nil {
		return nil, err
	}
	if inv == nil || inv.CreatedByID != userID {
		return nil, apperror.New(fiber.StatusNotFound)
	}
	return inv, nil
}

func (u *invoiceUC) ListInvoices(userID uuid.UUID) ([]domain.Invoice, error) {
	return u.repo.ListInvoicesByUser(userID)
}
