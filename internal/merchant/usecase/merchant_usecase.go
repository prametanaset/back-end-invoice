package usecase

import (
	"invoice_project/internal/merchant/domain"
	"invoice_project/internal/merchant/repository"
	"invoice_project/pkg/apperror"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type MerchantUsecase interface {
	GetMyMerchant(userID uuid.UUID) (*domain.Merchant, error)
	CreateMerchant(userID uuid.UUID, merchantType string) (*domain.Merchant, error)
	CreateStore(merchantID uuid.UUID, name string, branch string, addr domain.StoreAddress) (*domain.Store, error)
	ListStores(merchantID uuid.UUID) ([]domain.Store, error)
	AddPersonInfo(merchantID uuid.UUID, firstName, lastName string, vatNo *string) (*domain.PersonMerchant, error)
	AddCompanyInfo(merchantID uuid.UUID, companyName, vatNo string) (*domain.CompanyMerchant, error)
	AddContact(merchantID uuid.UUID, contactType, contactValue string) (*domain.MerchantContact, error)
	ListContacts(merchantID uuid.UUID) ([]domain.MerchantContact, error)
	GetPerson(merchantID uuid.UUID) (*domain.PersonMerchant, error)
	GetCompany(merchantID uuid.UUID) (*domain.CompanyMerchant, error)
	HasStore(userID uuid.UUID) (bool, error)
}

// StoreAddressInput holds address fields when creating a store.
type StoreAddressInput struct {
	AddressLine1  string
	SubdistrictID int
	DistrictID    int
	ProvinceID    int
	PostalCode    string
}

// ToDomain converts the input to a domain.StoreAddress.
func (in StoreAddressInput) ToDomain() domain.StoreAddress {
	return domain.StoreAddress{
		AddressLine1:  in.AddressLine1,
		SubdistrictID: in.SubdistrictID,
		DistrictID:    in.DistrictID,
		ProvinceID:    in.ProvinceID,
		PostalCode:    in.PostalCode,
	}
}

type merchantUC struct{ repo repository.MerchantRepository }

func NewMerchantUsecase(repo repository.MerchantRepository) MerchantUsecase {
	return &merchantUC{repo: repo}
}

func (u *merchantUC) GetMyMerchant(userID uuid.UUID) (*domain.Merchant, error) {
	return u.repo.GetMerchantByUser(userID)
}

func (u *merchantUC) CreateMerchant(userID uuid.UUID, merchantType string) (*domain.Merchant, error) {
	if merchantType != "person" && merchantType != "company" {
		return nil, apperror.New(fiber.StatusBadRequest)
	}

	// Prevent creating a merchant when a store already exists for the user.
	hasStore, err := u.HasStore(userID)
	if err != nil {
		return nil, err
	}
	if hasStore {
		return nil, apperror.New(fiber.StatusConflict)
	}

	mt, err := u.repo.GetMerchantTypeByName(merchantType)
	if err != nil {
		return nil, err
	}
	if mt == nil {
		return nil, apperror.New(fiber.StatusBadRequest)
	}

	exist, err := u.repo.GetMerchantByUserAndType(userID, mt.ID)
	if err != nil {
		return nil, err
	}
	if exist != nil {
		return nil, apperror.New(fiber.StatusConflict)
	}

	m := &domain.Merchant{UserID: userID, MerchantTypeID: mt.ID}
	if err := u.repo.CreateMerchant(m); err != nil {
		return nil, err
	}
	m.MerchantType = *mt
	return m, nil
}

func (u *merchantUC) CreateStore(merchantID uuid.UUID, name string, branch string, addr domain.StoreAddress) (*domain.Store, error) {
	s := &domain.Store{MerchantID: merchantID, StoreName: name, BranchNo: branch}
	if err := u.repo.CreateStore(s, &addr); err != nil {
		return nil, err
	}
	return s, nil
}

func (u *merchantUC) ListStores(merchantID uuid.UUID) ([]domain.Store, error) {
	return u.repo.ListStores(merchantID)
}

func (u *merchantUC) AddPersonInfo(merchantID uuid.UUID, firstName, lastName string, vatNo *string) (*domain.PersonMerchant, error) {
	m, err := u.repo.GetMerchant(merchantID)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, apperror.New(fiber.StatusNotFound)
	}
	if m.MerchantType.Name != "person" {
		return nil, apperror.New(fiber.StatusBadRequest)
	}
	p := &domain.PersonMerchant{
		MerchantID: merchantID,
		FirstName:  firstName,
		LastName:   lastName,
		VatNo:      vatNo,
	}
	if err := u.repo.CreatePerson(p); err != nil {
		return nil, err
	}
	return p, nil
}

func (u *merchantUC) AddCompanyInfo(merchantID uuid.UUID, companyName, vatNo string) (*domain.CompanyMerchant, error) {
	m, err := u.repo.GetMerchant(merchantID)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, apperror.New(fiber.StatusNotFound)
	}
	if m.MerchantType.Name != "company" {
		return nil, apperror.New(fiber.StatusBadRequest)
	}
	c := &domain.CompanyMerchant{
		MerchantID:  merchantID,
		CompanyName: companyName,
		VatNo:       vatNo,
	}
	if err := u.repo.CreateCompany(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (u *merchantUC) AddContact(merchantID uuid.UUID, contactType, contactValue string) (*domain.MerchantContact, error) {
	contact := &domain.MerchantContact{
		MerchantID:   merchantID,
		ContactType:  contactType,
		ContactValue: contactValue,
	}
	if err := u.repo.CreateContact(contact); err != nil {
		return nil, err
	}
	return contact, nil
}

func (u *merchantUC) ListContacts(merchantID uuid.UUID) ([]domain.MerchantContact, error) {
	return u.repo.ListContacts(merchantID)
}

func (u *merchantUC) GetPerson(merchantID uuid.UUID) (*domain.PersonMerchant, error) {
	return u.repo.GetPerson(merchantID)
}

func (u *merchantUC) GetCompany(merchantID uuid.UUID) (*domain.CompanyMerchant, error) {
	return u.repo.GetCompany(merchantID)
}

func (u *merchantUC) HasStore(userID uuid.UUID) (bool, error) {
	m, err := u.repo.GetMerchantByUser(userID)
	if err != nil {
		return false, err
	}
	if m == nil {
		return false, nil
	}
	stores, err := u.repo.ListStores(m.ID)
	if err != nil {
		return false, err
	}
	return len(stores) > 0, nil
}
