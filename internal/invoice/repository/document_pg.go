package repository

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
	"invoice_project/internal/invoice/domain"
)

type InvoiceDocumentRepository interface {
	CreateDocument(ctx context.Context, doc *domain.InvoiceDocument, items []domain.InvoiceItem) error
	GetDocument(ctx context.Context, id uint) (*domain.InvoiceDocument, error)
}

type documentPG struct {
	db *gorm.DB
}

func NewInvoiceDocumentRepository(db *gorm.DB) InvoiceDocumentRepository {
	return &documentPG{db: db}
}

func (r *documentPG) CreateDocument(ctx context.Context, doc *domain.InvoiceDocument, items []domain.InvoiceItem) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(doc).Error; err != nil {
			return err
		}
		for i := range items {
			items[i].DocumentID = doc.ID
		}
		if len(items) > 0 {
			if err := tx.Create(&items).Error; err != nil {
				return err
			}
		}
		tl := domain.DocumentTimeline{
			DocumentID: doc.ID,
			EventType:  "created",
			NewStatus:  doc.Status,
			ChangedAt:  time.Now(),
		}
		if err := tx.Create(&tl).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *documentPG) GetDocument(ctx context.Context, id uint) (*domain.InvoiceDocument, error) {
	var doc domain.InvoiceDocument
	err := r.db.WithContext(ctx).
		Preload("Items").
		Preload("Timelines").
		First(&doc, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &doc, nil
}
