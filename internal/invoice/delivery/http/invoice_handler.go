package http

import (
	"strconv"

	"invoice_project/internal/invoice/usecase"
	"invoice_project/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

type InvoiceHandler struct {
	invUC      usecase.InvoiceUsecase
	authSecret string
}

func NewInvoiceHandler(invUC usecase.InvoiceUsecase, authSecret string) *InvoiceHandler {
	return &InvoiceHandler{
		invUC:      invUC,
		authSecret: authSecret,
	}
}

func (h *InvoiceHandler) Create(c *fiber.Ctx) error {
	var body struct {
		Customer string  `json:"customer"`
		Amount   float64 `json:"amount"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
	}
	userID := c.Locals("user_id").(uint)
	inv, err := h.invUC.CreateInvoice(body.Customer, body.Amount, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(inv)
}

func (h *InvoiceHandler) GetByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid invoice ID"})
	}
	userID := c.Locals("user_id").(uint)
	inv, err := h.invUC.GetInvoice(uint(id), userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(inv)
}

func (h *InvoiceHandler) List(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	invoices, err := h.invUC.ListInvoices(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(invoices)
}

func (h *InvoiceHandler) RegisterRoutes(app *fiber.App) {
	api := app.Group("/invoices", middleware.JWTMiddleware(h.authSecret))
	api.Post("/", h.Create)
	api.Get("/", h.List)
	api.Get("/:id", h.GetByID)
}
