package domain

import "time"

// InvoiceDocument represents an issued invoice with summary information.
type InvoiceDocument struct {
	ID             uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	DocumentType   string    `gorm:"size:50" json:"document_type"`
	DocumentNumber int       `json:"document_number"`
	ReferenceID    *uint     `json:"reference_id"`
	StoreID        *string   `gorm:"type:uuid" json:"store_id"`
	CustomerID     *uint     `json:"customer_id"`
	IssueDate      time.Time `gorm:"type:date" json:"issue_date"`
	Status         string    `gorm:"size:50" json:"status"`
	BuyerName      string    `gorm:"size:255" json:"buyer_name"`
	BuyerTaxID     string    `gorm:"size:100" json:"buyer_tax_id"`
	BuyerAddress   string    `gorm:"type:text" json:"buyer_address"`
	SellerName     string    `gorm:"size:255" json:"seller_name"`
	SellerTaxID    string    `gorm:"size:100" json:"seller_tax_id"`
	SellerAddress  string    `gorm:"type:text" json:"seller_address"`
	Subtotal       int       `json:"subtotal"`
	DiscountType   int       `json:"discount_type"`
	DiscountValue  int       `json:"discount_value"`
	DiscountAmount int       `json:"discount_amount"`
	VatAmount      int       `json:"vat_amount"`
	GrandTotal     int       `json:"grand_total"`
	Remarks        string    `gorm:"type:text" json:"remarks"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`

	Items     []InvoiceItem      `gorm:"foreignKey:DocumentID" json:"items,omitempty"`
	Timelines []DocumentTimeline `gorm:"foreignKey:DocumentID" json:"timelines,omitempty"`
}

// InvoiceItem is a line item within an invoice document.
type InvoiceItem struct {
	ID          uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	DocumentID  uint    `gorm:"not null" json:"document_id"`
	ProductID   *uint   `json:"product_id"`
	ProductName string  `gorm:"size:255" json:"product_name"`
	Sku         string  `gorm:"size:100" json:"sku"`
	Qty         int     `json:"qty"`
	UnitPrice   int     `json:"unit_price"`
	Discount    int     `json:"discount"`
	VatType     string  `gorm:"size:50" json:"vat_type"`
	VatRate     float64 `gorm:"type:numeric(5,2)" json:"vat_rate"`
	LineTotal   float64 `gorm:"type:numeric(12,2)" json:"line_total"`
}

// DocumentTimeline records status changes for a document.
type DocumentTimeline struct {
	ID                uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	DocumentID        uint      `gorm:"not null" json:"document_id"`
	RelatedDocumentID *uint     `json:"related_document_id"`
	EventType         string    `gorm:"size:100" json:"event_type"`
	OldStatus         string    `gorm:"size:50" json:"old_status"`
	NewStatus         string    `gorm:"size:50" json:"new_status"`
	ChangedBy         string    `gorm:"size:100" json:"changed_by"`
	ChangedAt         time.Time `gorm:"autoCreateTime" json:"changed_at"`
	Note              string    `gorm:"type:text" json:"note"`
}
