package http

type CreateMerchantRequest struct {
	MerchantType string `json:"merchant_type"`
}

type CreateStoreRequest struct {
	MerchantID    string `json:"merchant_id"`
	StoreName     string `json:"store_name"`
	BranchNo      string `json:"branch_no"`
	AddressLine1  string `json:"address_line1"`
	SubdistrictID int    `json:"subdistrict_id"`
	DistrictID    int    `json:"district_id"`
	ProvinceID    int    `json:"province_id"`
	PostalCode    string `json:"postal_code"`
}

type AddPersonRequest struct {
	MerchantID string `json:"merchant_id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	VatNo      string `json:"vat_no"`
}

type AddCompanyRequest struct {
	MerchantID  string `json:"merchant_id"`
	CompanyName string `json:"company_name"`
	VatNo       string `json:"vat_no"`
}

type AddContactRequest struct {
	MerchantID   string `json:"merchant_id"`
	ContactType  string `json:"contact_type"`
	ContactValue string `json:"contact_value"`
}

// RegisterMerchantRequest represents payload to register a merchant with store information in a single call.
type RegisterMerchantRequest struct {
	MerchantType string `json:"merchant_type"`
	Person       *struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		VatNo     string `json:"vat_no"`
	} `json:"person,omitempty"`
	Company *struct {
		CompanyName string `json:"company_name"`
		VatNo       string `json:"vat_no"`
	} `json:"company,omitempty"`
	Store struct {
		StoreName     string `json:"store_name"`
		BranchNo      string `json:"branch_no"`
		AddressLine1  string `json:"address_line1"`
		SubdistrictID int    `json:"subdistrict_id"`
		DistrictID    int    `json:"district_id"`
		ProvinceID    int    `json:"province_id"`
		PostalCode    string `json:"postal_code"`
	} `json:"store"`
	Contacts []struct {
		ContactType  string `json:"contact_type"`
		ContactValue string `json:"contact_value"`
	} `json:"contacts,omitempty"`
}
