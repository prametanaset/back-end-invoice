package http

import (
	"invoice_project/internal/invoice/usecase"
	"invoice_project/pkg/apperror"
	"invoice_project/pkg/middleware"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type DocumentHandler struct {
	uc usecase.InvoiceDocumentUsecase
}

func NewDocumentHandler(uc usecase.InvoiceDocumentUsecase) *DocumentHandler {
	return &DocumentHandler{uc: uc}
}

func (h *DocumentHandler) Create(c *fiber.Ctx) error {
	var req CreateInvoiceDocumentRequest
	if err := c.BodyParser(&req); err != nil {
		return apperror.New(fiber.StatusBadRequest)
	}
	if err := h.uc.CreateDocument(c.Context(), &req.Document, req.Items); err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "created"})
}

func (h *DocumentHandler) Get(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return apperror.New(fiber.StatusBadRequest)
	}
	doc, err := h.uc.GetDocument(c.Context(), uint(id))
	if err != nil {
		return err
	}
	if doc == nil {
		return apperror.New(fiber.StatusNotFound)
	}
	return c.JSON(doc)
}

func (h *DocumentHandler) RegisterRoutes(app *fiber.App) {
	api := app.Group("/invoice-documents", middleware.RequireRoles("user", "admin"))
	api.Post("/", h.Create)
	api.Get("/:id", h.Get)
}
