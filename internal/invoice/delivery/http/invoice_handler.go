package http

import (
	"invoice_project/internal/invoice/usecase"
	"invoice_project/pkg/apperror"
	"invoice_project/pkg/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type InvoiceHandler struct {
	invUC usecase.InvoiceUsecase
}

func NewInvoiceHandler(invUC usecase.InvoiceUsecase) *InvoiceHandler {
	return &InvoiceHandler{
		invUC: invUC,
	}
}

func (h *InvoiceHandler) Create(c *fiber.Ctx) error {
	var body CreateInvoiceRequest
	if err := c.BodyParser(&body); err != nil {
		return apperror.New(fiber.StatusBadRequest)
	}
	userID := c.Locals("user_id").(uuid.UUID)
	inv, err := h.invUC.CreateInvoice(body.Customer, body.Amount, userID)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(inv)
}

func (h *InvoiceHandler) GetByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	uuidID, err := uuid.Parse(idParam)
	if err != nil {
		return apperror.New(fiber.StatusBadRequest)
	}
	userID := c.Locals("user_id").(uuid.UUID)
	inv, err := h.invUC.GetInvoice(uuidID, userID)
	if err != nil {
		return err
	}
	return c.JSON(inv)
}

func (h *InvoiceHandler) List(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)
	invoices, err := h.invUC.ListInvoices(userID)
	if err != nil {
		return err
	}
	return c.JSON(invoices)
}

func (h *InvoiceHandler) RegisterRoutes(app *fiber.App) {
	api := app.Group("/invoices", middleware.RequireRoles("user", "admin"))
	api.Post("/", h.Create)
	api.Get("/", h.List)
	api.Get("/:id", h.GetByID)
}
