package http

import (
	"strconv"

	"invoice_project/internal/location/usecase"
	"invoice_project/pkg/apperror"
	"invoice_project/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

// LocationHandler handles HTTP requests for thai location data

type LocationHandler struct {
	uc usecase.LocationUsecase
}

func NewLocationHandler(uc usecase.LocationUsecase) *LocationHandler {
	return &LocationHandler{uc: uc}
}

func (h *LocationHandler) listGeographies(c *fiber.Ctx) error {
	geos, err := h.uc.ListGeographies(c.Context())
	if err != nil {
		return apperror.New(fiber.StatusInternalServerError)
	}
	return c.JSON(geos)
}

func (h *LocationHandler) listProvinces(c *fiber.Ctx) error {
	gidStr := c.Query("geo_id")
	gid, _ := strconv.Atoi(gidStr)
	provinces, err := h.uc.ListProvinces(c.Context(), gid)
	if err != nil {
		return apperror.New(fiber.StatusInternalServerError)
	}
	return c.JSON(provinces)
}

func (h *LocationHandler) listAmphures(c *fiber.Ctx) error {
	pidStr := c.Query("province_id")
	pid, _ := strconv.Atoi(pidStr)
	amphures, err := h.uc.ListAmphures(c.Context(), pid)
	if err != nil {
		return apperror.New(fiber.StatusInternalServerError)
	}
	return c.JSON(amphures)
}

func (h *LocationHandler) listTambons(c *fiber.Ctx) error {
	aidStr := c.Query("amphure_id")
	aid, _ := strconv.Atoi(aidStr)
	tambons, err := h.uc.ListTambons(c.Context(), aid)
	if err != nil {
		return apperror.New(fiber.StatusInternalServerError)
	}
	return c.JSON(tambons)
}

// RegisterRoutes registers location endpoints
func (h *LocationHandler) RegisterRoutes(app *fiber.App) {
	api := app.Group("/locations", middleware.RequireRoles("user", "admin"))
	api.Get("/geographies", h.listGeographies)
	api.Get("/provinces", h.listProvinces)
	api.Get("/amphures", h.listAmphures)
	api.Get("/tambons", h.listTambons)
}
