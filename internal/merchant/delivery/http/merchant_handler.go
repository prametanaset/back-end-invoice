package http

import (
	"invoice_project/internal/merchant/usecase"
	"invoice_project/pkg/apperror"
	"invoice_project/pkg/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type MerchantHandler struct {
	uc usecase.MerchantUsecase
}

func NewMerchantHandler(uc usecase.MerchantUsecase) *MerchantHandler {
	return &MerchantHandler{uc: uc}
}

func (h *MerchantHandler) CreateMerchant(c *fiber.Ctx) error {
	var body CreateMerchantRequest
	if err := c.BodyParser(&body); err != nil {
		return apperror.New(fiber.StatusBadRequest)
	}
	userID := c.Locals("user_id").(uuid.UUID)
	m, err := h.uc.CreateMerchant(userID, body.MerchantType)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(m)
}

func (h *MerchantHandler) CreateStore(c *fiber.Ctx) error {
	var body CreateStoreRequest
	if err := c.BodyParser(&body); err != nil {
		return apperror.New(fiber.StatusBadRequest)
	}
	merchantID, err := uuid.Parse(body.MerchantID)
	if err != nil {
		return apperror.New(fiber.StatusBadRequest)
	}
	addr := usecase.StoreAddressInput{
		AddressLine1:  body.AddressLine1,
		SubdistrictID: body.SubdistrictID,
		DistrictID:    body.DistrictID,
		ProvinceID:    body.ProvinceID,
		PostalCode:    body.PostalCode,
	}
	store, err := h.uc.CreateStore(merchantID, body.StoreName, body.BranchNo, addr.ToDomain())
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(store)
}

func (h *MerchantHandler) AddPerson(c *fiber.Ctx) error {
	var body AddPersonRequest
	if err := c.BodyParser(&body); err != nil {
		return apperror.New(fiber.StatusBadRequest)
	}
	merchantID, err := uuid.Parse(body.MerchantID)
	if err != nil {
		return apperror.New(fiber.StatusBadRequest)
	}
	var vat *string
	if body.VatNo != "" {
		vat = &body.VatNo
	}
	person, err := h.uc.AddPersonInfo(merchantID, body.FirstName, body.LastName, vat)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(person)
}

func (h *MerchantHandler) AddCompany(c *fiber.Ctx) error {
	var body AddCompanyRequest
	if err := c.BodyParser(&body); err != nil {
		return apperror.New(fiber.StatusBadRequest)
	}
	merchantID, err := uuid.Parse(body.MerchantID)
	if err != nil {
		return apperror.New(fiber.StatusBadRequest)
	}
	comp, err := h.uc.AddCompanyInfo(merchantID, body.CompanyName, body.VatNo)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(comp)
}

func (h *MerchantHandler) AddContact(c *fiber.Ctx) error {
	var body AddContactRequest
	if err := c.BodyParser(&body); err != nil {
		return apperror.New(fiber.StatusBadRequest)
	}
	merchantID, err := uuid.Parse(body.MerchantID)
	if err != nil {
		return apperror.New(fiber.StatusBadRequest)
	}
	contact, err := h.uc.AddContact(merchantID, body.ContactType, body.ContactValue)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(contact)
}

// RegisterMerchant handles registering a merchant, its store and optional contacts in one call.
func (h *MerchantHandler) RegisterMerchant(c *fiber.Ctx) error {
	var body RegisterMerchantRequest
	if err := c.BodyParser(&body); err != nil {
		return apperror.New(fiber.StatusBadRequest)
	}

	userID := c.Locals("user_id").(uuid.UUID)
	merchant, err := h.uc.CreateMerchant(userID, body.MerchantType)
	if err != nil {
		return err
	}

	switch merchant.MerchantType.Name {
	case "person":
		if body.Person == nil {
			return apperror.New(fiber.StatusBadRequest)
		}
		var vat *string
		if body.Person.VatNo != "" {
			vat = &body.Person.VatNo
		}
		if _, err := h.uc.AddPersonInfo(merchant.ID, body.Person.FirstName, body.Person.LastName, vat); err != nil {
			return err
		}
	case "company":
		if body.Company == nil {
			return apperror.New(fiber.StatusBadRequest)
		}
		if _, err := h.uc.AddCompanyInfo(merchant.ID, body.Company.CompanyName, body.Company.VatNo); err != nil {
			return err
		}
	default:
		return apperror.New(fiber.StatusBadRequest)
	}

	addr := usecase.StoreAddressInput{
		AddressLine1:  body.Store.AddressLine1,
		SubdistrictID: body.Store.SubdistrictID,
		DistrictID:    body.Store.DistrictID,
		ProvinceID:    body.Store.ProvinceID,
		PostalCode:    body.Store.PostalCode,
	}
	store, err := h.uc.CreateStore(merchant.ID, body.Store.StoreName, body.Store.BranchNo, addr.ToDomain())
	if err != nil {
		return err
	}

	var contacts []interface{}
	for _, ctt := range body.Contacts {
		contact, err := h.uc.AddContact(merchant.ID, ctt.ContactType, ctt.ContactValue)
		if err != nil {
			return err
		}
		contacts = append(contacts, contact)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"merchant": merchant,
		"store":    store,
		"contacts": contacts,
	})
}

func (h *MerchantHandler) ListContacts(c *fiber.Ctx) error {
	merchantIDStr := c.Query("merchant_id")
	id, err := uuid.Parse(merchantIDStr)
	if err != nil {
		return apperror.New(fiber.StatusBadRequest)
	}
	contacts, err := h.uc.ListContacts(id)
	if err != nil {
		return err
	}
	return c.JSON(contacts)
}

func (h *MerchantHandler) ListStores(c *fiber.Ctx) error {
	merchantIDStr := c.Query("merchant_id")
	id, err := uuid.Parse(merchantIDStr)
	if err != nil {
		return apperror.New(fiber.StatusBadRequest)
	}
	stores, err := h.uc.ListStores(id)
	if err != nil {
		return err
	}
	return c.JSON(stores)
}

func (h *MerchantHandler) RegisterRoutes(app *fiber.App) {
	api := app.Group("/merchants", middleware.RequireRoles("user", "admin"))
	api.Post("/register", h.RegisterMerchant)
	api.Post("/", h.CreateMerchant)
	api.Post("/stores", h.CreateStore)
	api.Get("/stores", h.ListStores)
	api.Post("/person", h.AddPerson)
	api.Post("/company", h.AddCompany)
	api.Post("/contacts", h.AddContact)
	api.Get("/contacts", h.ListContacts)
}
