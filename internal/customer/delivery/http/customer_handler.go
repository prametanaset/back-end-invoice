package http

import (
	"invoice_project/internal/customer/usecase"
	"invoice_project/pkg/middleware"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CustomerHandler struct {
	usecase usecase.CustomerUseCase
}

func NewCustomerHandler(uc usecase.CustomerUseCase) *CustomerHandler {
	return &CustomerHandler{usecase: uc}
}

func (h *CustomerHandler) CreateCustomer(c *fiber.Ctx) error {
	var req CreateCustomerRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// ตรวจสอบว่ามี Customer จริง
	if req.Customer.StoreID == uuid.Nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Store ID is required",
		})
	}

	// เรียกใช้ usecase
	err := h.usecase.CreateCustomer(c.Context(), &req.Customer, req.Person, req.Company, req.Address, req.Contacts)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create customer",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":  "Customer created successfully",
		"customer": req.Customer,
	})
}


// GET /customers/:id
func (h *CustomerHandler) GetCustomerByID(c *fiber.Ctx) error {
    idStr := c.Params("id")
    idUint, err := strconv.ParseUint(idStr, 10, 64)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
    }

    customer, err := h.usecase.GetCustomerByID(c.Context(), uint(idUint))
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Customer not found"})
    }

    return c.JSON(customer)
}


// GET /customers/store/:store_id
func (h *CustomerHandler) ListCustomer(c *fiber.Ctx) error {
	storeIDStr := c.Params("store_id")
	storeID, err := uuid.Parse(storeIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid store ID"})
	}

	customers, err := h.usecase.ListCustomer(c.Context(), storeID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to list customers"})
	}

	return c.JSON(customers)
}

// PUT /customers/:id
func (h *CustomerHandler) UpdateCustomer(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var req UpdateCustomerRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// ผูก customer_id ให้ทุก entity
	req.Customer.ID = uint(id)
	if req.Person != nil {
		req.Person.CustomerID = req.Customer.ID
	}
	if req.Company != nil {
		req.Company.CustomerID = req.Customer.ID
	}
	if req.Address != nil {
		req.Address.CustomerID = req.Customer.ID
	}
	for i := range req.Contacts {
		req.Contacts[i].CustomerID = req.Customer.ID
	}

	if err := h.usecase.UpdateCustomer(c.Context(), &req.Customer, req.Person, req.Company, req.Address, req.Contacts); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update customer"})
	}

	return c.JSON(fiber.Map{"message": "Customer updated successfully"})
}


// DELETE /customers/:id
func (h *CustomerHandler) DeleteCustomer(c *fiber.Ctx) error {
	idStr := c.Params("id")
	idUint, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	err = h.usecase.DeleteCustomer(c.Context(), uint(idUint))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete customer"})
	}

	return c.SendStatus(fiber.StatusNoContent)
}


func (h *CustomerHandler) RegisterRoutes(app *fiber.App) {
	api := app.Group("/customers", middleware.RequireRoles("user", "admin"))
	api.Post("/", h.CreateCustomer)
	api.Get("/:id", h.GetCustomerByID)
	api.Get("/store/:store_id", h.ListCustomer)
	api.Put("/:id", h.UpdateCustomer)
	api.Delete("/:id", h.DeleteCustomer)
}