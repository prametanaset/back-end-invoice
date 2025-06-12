package http

import (
	"invoice_project/internal/product/usecase"
	"invoice_project/pkg/apperror"
	"invoice_project/pkg/middleware"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ProductHandler struct {
	uc usecase.ProductUseCase
}

func NewProductHandler(uc usecase.ProductUseCase) *ProductHandler {
	return &ProductHandler{uc: uc}
}

func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	var req CreateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return apperror.New(fiber.StatusBadRequest)
	}


	if req.Product.StoreID == uuid.Nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Store ID is required",
		})
	}

	err := h.uc.CreateProduct(c.Context(), &req.Product, req.ProductImage);
	if  err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create product",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "product created successfully",
	})
}


func (h *ProductHandler) GetProduct(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return apperror.New(fiber.StatusBadRequest)
	}

	product, images, err := h.uc.GetProduct(c.Context(), uint(id))
	if err != nil {
		return apperror.New(fiber.StatusNotFound)
	}

	return c.JSON(fiber.Map{
		"product": product,
		"images":  images,
	})
}

func (h *ProductHandler) ListProducts(c *fiber.Ctx) error {
	storeIDParam := c.Query("store_id")
	storeID, err := uuid.Parse(storeIDParam)
	if err != nil {
		return apperror.New(fiber.StatusBadRequest)
	}

	products, err := h.uc.ListProducts(c.Context(), storeID)
	if err != nil {
		return apperror.New(fiber.StatusBadRequest)
	}

	return c.JSON(products)
}

func (h *ProductHandler) UpdateProduct(c *fiber.Ctx) error {
	var req CreateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return apperror.New(fiber.StatusBadRequest)
	}

	err := h.uc.UpdateProduct(c.Context(), &req.Product, req.ProductImage)
	if err != nil {
		return apperror.New(fiber.StatusBadRequest)
	}

	return c.JSON(fiber.Map{"message": "Product updated successfully"})
}

func (h *ProductHandler) DeleteProduct(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return apperror.New(fiber.StatusBadRequest)
	}

	err = h.uc.DeleteProduct(c.Context(), uint(id))
	if err != nil {
		return apperror.New(fiber.StatusBadRequest)
	}

	return c.JSON(fiber.Map{"message": "Product deleted successfully"})
}

// ------------------------ Routes -------------------------

func (h *ProductHandler) RegisterRoutes(app *fiber.App) {
	api := app.Group("/products", middleware.RequireRoles("user", "admin"))
	api.Post("/", h.CreateProduct)
	api.Get("/", h.ListProducts)      // ?store_id=<uuid>
	api.Get("/:id", h.GetProduct)
	api.Put("/", h.UpdateProduct)     // Body: Product + Images
	api.Delete("/:id", h.DeleteProduct)
}