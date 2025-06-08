package http

// CreateInvoiceRequest represents the expected payload for creating an invoice.
type CreateInvoiceRequest struct {
	Customer string  `json:"customer"`
	Amount   float64 `json:"amount"`
}
