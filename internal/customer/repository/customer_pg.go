package repository

import (
	"context"
	"invoice_project/internal/customer/domain"
	"invoice_project/pkg/apperror"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CustomerRepository interface {
	CreateCustomer(
		ctx context.Context,
		m *domain.Customer,
		person *domain.PersonCustomer,
		company *domain.CompanyCustomer,
		address *domain.CustomerAddress,
		contacts []domain.CustomerContact,
	) error
	GetCustomer(ctx context.Context, id uint) (*domain.Customer, error)
	ListCustomer(storeID uuid.UUID) ([]domain.Customer, error)
	UpdateCustomer(ctx context.Context,
		m *domain.Customer,
		person *domain.PersonCustomer,
		company *domain.CompanyCustomer,
		address *domain.CustomerAddress,
		contacts []domain.CustomerContact,) error
	DeleteCustomer(id uint) error

	WithTx(ctx context.Context, fn func(tx *gorm.DB) error) error
}


type customerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) *customerRepository {
	return &customerRepository{db: db}
}

// CreateCustomer เพิ่มลูกค้าใหม่
func (r *customerRepository) CreateCustomer(ctx context.Context,m *domain.Customer, person *domain.PersonCustomer, company *domain.CompanyCustomer, address *domain.CustomerAddress, contacts []domain.CustomerContact) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Insert Customer
		if err := tx.Create(m).Error; err != nil {
			return err
		}

		// Insert Person or Company based on customer type
		switch m.CustomerType {
		case "person":
			if person != nil {
				person.CustomerID = m.ID
				if err := tx.Create(person).Error; err != nil {
					return err
				}
			}
		case "company":
			if company != nil {
				company.CustomerID = m.ID
				if err := tx.Create(company).Error; err != nil {
					return err
				}
			}
		default:
			return apperror.New(fiber.StatusBadRequest)
		}

		// Insert Address
		if address != nil {
			address.CustomerID = m.ID
			if err := tx.Create(address).Error; err != nil {
				return err
			}
		}

		// Insert Contact(s)
		for i := range contacts {
			contacts[i].CustomerID = m.ID
			if err := tx.Create(&contacts[i]).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// GetCustomer ดึงข้อมูลลูกค้าตาม id พร้อม preload ความสัมพันธ์
func (r *customerRepository) GetCustomer(ctx context.Context, id uint) (*domain.Customer, error) {
	var customer domain.Customer
	err := r.db.WithContext(ctx).
		Preload("CompanyCustomer").
		Preload("PersonCustomer").
		Preload("CustomerAddress").
		Preload("CustomerContact").
		Where("id = ?", id).
		First(&customer).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

// ListCustomer ดึงรายชื่อลูกค้าทั้งหมดของร้าน
func (r *customerRepository) ListCustomer(storeID uuid.UUID) ([]domain.Customer, error) {
	var customers []domain.Customer
	err := r.db.Where("store_id = ?", storeID).Find(&customers).Error
	if err != nil {
		return nil, err
	}
	return customers, nil
}

// UpdateCustomer อัปเดตข้อมูลลูกค้า
func (r *customerRepository) UpdateCustomer(
	ctx context.Context,
	customer *domain.Customer,
	person *domain.PersonCustomer,
	company *domain.CompanyCustomer,
	address *domain.CustomerAddress,
	contacts []domain.CustomerContact,
) error {
	return r.WithTx(ctx, func(tx *gorm.DB) error {
		// อัปเดต customer
		if err := tx.Model(&domain.Customer{}).
			Where("id = ?", customer.ID).
			Updates(customer).Error; err != nil {
			return err
		}

		// อัปเดต person
		if person != nil {
			if err := tx.Model(&domain.PersonCustomer{}).
				Where("customer_id = ?", customer.ID).
				Updates(person).Error; err != nil {
				return err
			}
		}

		// อัปเดต company
		if company != nil {
			if err := tx.Model(&domain.CompanyCustomer{}).
				Where("customer_id = ?", customer.ID).
				Updates(company).Error; err != nil {
				return err
			}
		}

		// อัปเดต address
		if address != nil {
			if err := tx.Model(&domain.CustomerAddress{}).
				Where("customer_id = ?", customer.ID).
				Updates(address).Error; err != nil {
				return err
			}
		}

		// ลบ contacts เดิมทั้งหมดก่อน
		if err := tx.Where("customer_id = ?", customer.ID).
			Delete(&domain.CustomerContact{}).Error; err != nil {
			return err
		}
		// เพิ่ม contacts ใหม่
		if len(contacts) > 0 {
			for i := range contacts {
				contacts[i].CustomerID = customer.ID
			}
			if err := tx.Create(&contacts).Error; err != nil {
				return err
			}
		}

		return nil
	})
}


func (r *customerRepository) WithTx(ctx context.Context, fn func(tx *gorm.DB) error) error {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}


// DeleteCustomer ลบลูกค้า (ใช้ soft delete ผ่าน GORM)
func (r *customerRepository) DeleteCustomer(id uint) error {
	tx := r.db.Begin()

	if err := tx.Where("customer_id = ?", id).Delete(&domain.CustomerContact{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Where("customer_id = ?", id).Delete(&domain.CustomerAddress{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Where("customer_id = ?", id).Delete(&domain.CompanyCustomer{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Where("customer_id = ?", id).Delete(&domain.PersonCustomer{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Delete(&domain.Customer{}, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
