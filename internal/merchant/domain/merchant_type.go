package domain

// MerchantType defines a type of merchant such as "person" or "company".
type MerchantType struct {
	ID   uint   `gorm:"primaryKey;autoIncrement"`
	Name string `gorm:"size:30;unique;not null"`
}

const (
	// MerchantTypePerson represents an individual merchant which is limited
	// to a single branch.
	MerchantTypePerson = "person"
	// MerchantTypeCompany represents a company merchant which may own
	// multiple branches.
	MerchantTypeCompany = "company"
)
