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
	userID := c.Locals("user_id").(uint)
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
	api.Post("/", h.CreateMerchant)
	api.Post("/stores", h.CreateStore)
	api.Get("/stores", h.ListStores)
	api.Post("/person", h.AddPerson)
	api.Post("/company", h.AddCompany)
	api.Post("/contacts", h.AddContact)
	api.Get("/contacts", h.ListContacts)
}
