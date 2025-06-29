package repository

import (
	"invoice_project/internal/merchant/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MerchantRepository interface {
	CreateMerchant(m *domain.Merchant) error
	GetMerchantByUser(userID uuid.UUID) (*domain.Merchant, error)
	GetMerchantByUserAndType(userID uuid.UUID, merchantTypeID uint) (*domain.Merchant, error)
	GetMerchant(id uuid.UUID) (*domain.Merchant, error)
	GetMerchantTypeByName(name string) (*domain.MerchantType, error)
	GetMerchantType(id uint) (*domain.MerchantType, error)
	CreateStore(store *domain.Store, addr *domain.StoreAddress) error
	ListStores(merchantID uuid.UUID) ([]domain.Store, error)
	CreatePerson(p *domain.PersonMerchant) error
	CreateCompany(c *domain.CompanyMerchant) error
	CreateContact(contact *domain.MerchantContact) error
	ListContacts(merchantID uuid.UUID) ([]domain.MerchantContact, error)
	GetPerson(merchantID uuid.UUID) (*domain.PersonMerchant, error)
	GetCompany(merchantID uuid.UUID) (*domain.CompanyMerchant, error)
}

type merchantPG struct{ db *gorm.DB }

func NewMerchantRepository(db *gorm.DB) MerchantRepository {
	return &merchantPG{db: db}
}

func (r *merchantPG) CreateMerchant(m *domain.Merchant) error {
	return r.db.Create(m).Error
}

func (r *merchantPG) GetMerchantByUser(userID uuid.UUID) (*domain.Merchant, error) {
	var m domain.Merchant
	err := r.db.Preload("MerchantType").Where("user_id = ?", userID).First(&m).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (r *merchantPG) GetMerchantByUserAndType(userID uuid.UUID, merchantTypeID uint) (*domain.Merchant, error) {
	var m domain.Merchant
	err := r.db.Where("user_id = ? AND merchant_type_id = ?", userID, merchantTypeID).First(&m).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (r *merchantPG) GetMerchant(id uuid.UUID) (*domain.Merchant, error) {
	var m domain.Merchant
	err := r.db.Preload("MerchantType").First(&m, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (r *merchantPG) GetMerchantTypeByName(name string) (*domain.MerchantType, error) {
	var mt domain.MerchantType
	err := r.db.Where("name = ?", name).First(&mt).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &mt, nil
}

func (r *merchantPG) GetMerchantType(id uint) (*domain.MerchantType, error) {
	var mt domain.MerchantType
	err := r.db.First(&mt, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &mt, nil
}

func (r *merchantPG) CreateStore(store *domain.Store, addr *domain.StoreAddress) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(store).Error; err != nil {
			return err
		}
		addr.StoreID = store.ID
		return tx.Create(addr).Error
	})
}

func (r *merchantPG) ListStores(merchantID uuid.UUID) ([]domain.Store, error) {
	var stores []domain.Store
	err := r.db.Where("merchant_id = ?", merchantID).Order("created_at desc").Find(&stores).Error
	if err != nil {
		return nil, err
	}
	return stores, nil
}

func (r *merchantPG) CreatePerson(p *domain.PersonMerchant) error {
	return r.db.Create(p).Error
}

func (r *merchantPG) CreateCompany(c *domain.CompanyMerchant) error {
	return r.db.Create(c).Error
}

func (r *merchantPG) CreateContact(contact *domain.MerchantContact) error {
	return r.db.Create(contact).Error
}

func (r *merchantPG) ListContacts(merchantID uuid.UUID) ([]domain.MerchantContact, error) {
	var contacts []domain.MerchantContact
	err := r.db.Where("merchant_id = ?", merchantID).Order("created_at desc").Find(&contacts).Error
	if err != nil {
		return nil, err
	}
	return contacts, nil
}

func (r *merchantPG) GetPerson(merchantID uuid.UUID) (*domain.PersonMerchant, error) {
	var p domain.PersonMerchant
	err := r.db.Where("merchant_id = ?", merchantID).First(&p).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

func (r *merchantPG) GetCompany(merchantID uuid.UUID) (*domain.CompanyMerchant, error) {
	var c domain.CompanyMerchant
	err := r.db.Where("merchant_id = ?", merchantID).First(&c).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}
