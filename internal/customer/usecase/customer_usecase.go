package usecase

import (
	"context"
	"invoice_project/internal/customer/domain"
	"invoice_project/internal/customer/repository"
	"invoice_project/pkg/apperror"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// CustomerUseCase - interface สำหรับ business logic
type CustomerUseCase interface {
	CreateCustomer(
		ctx context.Context,
		customer *domain.Customer,
		person *domain.PersonCustomer,
		company *domain.CompanyCustomer,
		address *domain.CustomerAddress,
		contacts []domain.CustomerContact,
	) error
	GetCustomerByID(ctx context.Context, id uint) (*domain.Customer, error)
	ListCustomer(ctx context.Context, storeID uuid.UUID) ([]domain.Customer, error)
	UpdateCustomer(
		ctx context.Context,
		customer *domain.Customer,
		person *domain.PersonCustomer,
		company *domain.CompanyCustomer,
		address *domain.CustomerAddress,
		contacts []domain.CustomerContact,
	) error
	DeleteCustomer(ctx context.Context, id uint) error
}

// customerUseCase - implements CustomerUseCase
type customerUseCase struct {
	repo repository.CustomerRepository
}

func NewCustomerUseCase(repo repository.CustomerRepository) CustomerUseCase {
	return &customerUseCase{
		repo: repo,
	}
}

func (uc *customerUseCase) CreateCustomer(
	ctx context.Context,
	customer *domain.Customer,
	person *domain.PersonCustomer,
	company *domain.CompanyCustomer,
	address *domain.CustomerAddress,
	contacts []domain.CustomerContact,
) error {
	if customer == nil {
		return apperror.New(fiber.StatusBadRequest)
	}

	customer.ID = 0 // ป้องกัน override

	return uc.repo.CreateCustomer(ctx, customer, person, company, address, contacts)
}


// GetCustomerByID ดึงลูกค้าตาม ID
func (uc *customerUseCase) GetCustomerByID(ctx context.Context, id uint) (*domain.Customer, error) {
	return uc.repo.GetCustomer(ctx, id)
}


// ListCustomer ดึงลูกค้าทั้งหมดของร้าน
func (uc *customerUseCase) ListCustomer(ctx context.Context, storeID uuid.UUID) ([]domain.Customer, error) {
	return uc.repo.ListCustomer(storeID)
}

// UpdateCustomer อัปเดตข้อมูลลูกค้า
func (uc *customerUseCase) UpdateCustomer(
	ctx context.Context,
	customer *domain.Customer,
	person *domain.PersonCustomer,
	company *domain.CompanyCustomer,
	address *domain.CustomerAddress,
	contacts []domain.CustomerContact,
) error {
	if customer == nil {
		return apperror.New(fiber.StatusBadRequest)
	}
	return uc.repo.UpdateCustomer(ctx, customer, person, company, address, contacts)
}

// DeleteCustomer ลบลูกค้า
func (uc *customerUseCase) DeleteCustomer(ctx context.Context, id uint) error {
	return uc.repo.DeleteCustomer(id)
}
