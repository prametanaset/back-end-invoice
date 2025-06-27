package domain

// MerchantType defines a type of merchant such as "person" or "company".
type MerchantType struct {
	ID   uint   `gorm:"primaryKey;autoIncrement"`
	Name string `gorm:"size:30;unique;not null"`
}
