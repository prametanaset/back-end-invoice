package http

import "invoice_project/internal/invoice/domain"

// CreateInvoiceRequest represents the expected payload for creating an invoice.
type CreateInvoiceRequest struct {
	Customer string  `json:"customer"`
	Amount   float64 `json:"amount"`
}

// CreateInvoiceDocumentRequest payload for creating invoice document
// along with its items.
type CreateInvoiceDocumentRequest struct {
	Document domain.InvoiceDocument `json:"document"`
	Items    []domain.InvoiceItem   `json:"items"`
}
