package usecase

import (
	"context"
	"invoice_project/internal/invoice/domain"
	"invoice_project/internal/invoice/repository"
	"invoice_project/pkg/apperror"

	"github.com/gofiber/fiber/v2"
)

type InvoiceDocumentUsecase interface {
	CreateDocument(ctx context.Context, doc *domain.InvoiceDocument, items []domain.InvoiceItem) error
	GetDocument(ctx context.Context, id uint) (*domain.InvoiceDocument, error)
}

type documentUC struct {
	repo repository.InvoiceDocumentRepository
}

func NewInvoiceDocumentUsecase(repo repository.InvoiceDocumentRepository) InvoiceDocumentUsecase {
	return &documentUC{repo: repo}
}

func (u *documentUC) CreateDocument(ctx context.Context, doc *domain.InvoiceDocument, items []domain.InvoiceItem) error {
	if doc == nil {
		return apperror.New(fiber.StatusBadRequest)
	}
	doc.ID = 0
	return u.repo.CreateDocument(ctx, doc, items)
}

func (u *documentUC) GetDocument(ctx context.Context, id uint) (*domain.InvoiceDocument, error) {
	if id == 0 {
		return nil, apperror.New(fiber.StatusBadRequest)
	}
	return u.repo.GetDocument(ctx, id)
}
