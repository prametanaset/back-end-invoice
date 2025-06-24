package http

import (
	"invoice_project/internal/location/usecase"
	"invoice_project/pkg/middleware"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type LocationHandler struct {
	usecase usecase.LocationUsecase
}

type ReturnType struct {
	Value int
	Label string
}

func NewLocationHandler(u usecase.LocationUsecase) *LocationHandler {
	return &LocationHandler{u}
}

func (h *LocationHandler) GetProvinceList(c *fiber.Ctx) error {
	provinces, err := h.usecase.GetProvinceAll(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "ไม่สามารถดึงข้อมูลจังหวัดได้",
		})
	}

	// แปลงเป็น ReturnType
	var result []ReturnType
	for _, p := range provinces {
		result = append(result, ReturnType{
			Value: int(p.ID),
			Label: p.NameTh,
		})
	}

	return c.JSON(result)
}

func (h *LocationHandler) GetProvinceByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID ไม่ถูกต้อง",
		})
	}

	province, err := h.usecase.GetProvinceByID(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "ไม่พบจังหวัดที่ต้องการ",
		})
	}

	return c.JSON(ReturnType{
		Value: int(province.ID),
		Label: province.NameTh,
	})
}

func (h *LocationHandler) GetProvincesByGeoID(c *fiber.Ctx) error {
	idParam := c.Params("geo_id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "รหัสภูมิภาคไม่ถูกต้อง",
		})
	}

	provinces, err := h.usecase.GetProvinceByGeoID(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "ไม่สามารถดึงข้อมูลจังหวัดได้",
		})
	}

	var result []ReturnType
	for _, p := range provinces {
		result = append(result, ReturnType{
			Value: int(p.ID),
			Label: p.NameTh,
		})
	}

	return c.JSON(result)
}

func (h *LocationHandler) GetDistrictById(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID ไม่ถูกต้อง",
		})
	}

	district, err := h.usecase.GetDistrictById(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "ไม่พบเขตที่ต้องการ",
		})
	}

	return c.JSON(ReturnType{
		Value: int(district.ID),
		Label: district.NameTh,
	})
}

func (h *LocationHandler) GetDistricts(c *fiber.Ctx) error {
	idParam := c.Params("province_id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "รหัสจังหวัดไม่ถูกต้อง",
		})
	}

	districts, err := h.usecase.GetDistricts(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "ไม่สามารถดึงข้อมูลอำเภอได้",
		})
	}

	// แปลงเป็น ReturnType
	var result []ReturnType
	for _, d := range districts {
		result = append(result, ReturnType{
			Value: int(d.ID),
			Label: d.NameTh,
		})
	}

	return c.JSON(result)
}

func (h *LocationHandler) GetSubDistrictById(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID ไม่ถูกต้อง",
		})
	}

	sdistrict, err := h.usecase.GetSubDistrictsById(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "ไม่พบเขตที่ต้องการ",
		})
	}

	return c.JSON(sdistrict)
}

func (h *LocationHandler) GetZipCodeBySubDistrictID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID ไม่ถูกต้อง",
		})
	}

	zip, err := h.usecase.GetZipCodeBySubDistrictID(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "ไม่พบข้อมูลรหัสไปรษณีย์",
		})
	}

	return c.JSON(fiber.Map{"zip_code": zip})
}

func (h *LocationHandler) GetSubDistricts(c *fiber.Ctx) error {
	idParam := c.Params("district_id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "รหัสอำเภอไม่ถูกต้อง",
		})
	}

	districts, err := h.usecase.GetSubDistricts(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "ไม่สามารถดึงข้อมูลตำบลได้",
		})
	}

	// แปลงเป็น ReturnType
	var result []ReturnType
	for _, d := range districts {
		result = append(result, ReturnType{
			Value: int(d.ID),
			Label: d.NameTh,
		})
	}

	return c.JSON(result)
}

func (h *LocationHandler) RegisterRoutes(app *fiber.App) {
	api := app.Group("/locations", middleware.RequireRoles("user", "admin"))
	api.Get("/province/", h.GetProvinceList)
	api.Get("/province/:id", h.GetProvinceByID)
	api.Get("/geography/:geo_id/provinces", h.GetProvincesByGeoID)
	api.Get("/district/:id", h.GetDistrictById)
	api.Get("/province/:province_id/districts", h.GetDistricts)
	api.Get("/subdistrict/:id", h.GetSubDistrictById)
	api.Get("/subdistrict/:id/zip_code", h.GetZipCodeBySubDistrictID)
	api.Get("/districts/:district_id/subdistricts", h.GetSubDistricts)

}
