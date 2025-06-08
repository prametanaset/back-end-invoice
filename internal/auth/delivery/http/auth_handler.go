package http

import (
	"time"

	"invoice_project/internal/auth/usecase"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authUC    usecase.AuthUsecase
	jwtSecret string
}

func NewAuthHandler(authUC usecase.AuthUsecase, jwtSecret string) *AuthHandler {
	return &AuthHandler{
		authUC:    authUC,
		jwtSecret: jwtSecret,
	}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var body RegisterRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
	}
	if err := h.authUC.Register(body.Username, body.Password); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "user registered"})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var body LoginRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
	}
	accessToken, refreshToken, err := h.authUC.Login(body.Username, body.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}
	// ตอบกลับทั้ง 2 token
	// เราอาจจะเซ็ต cookie ให้ refresh token ด้วยก็ได้ แต่ในที่นี้เก็บเป็น response JSON ธรรมดา
	return c.JSON(fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"expires_at":    time.Now().Add(15 * time.Minute), // บอกวันหมดอายุ access token ประมาณนี้
	})
}

func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	var body RefreshRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
	}
	newAccess, newRefresh, err := h.authUC.RefreshAccessToken(body.RefreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{
		"access_token":  newAccess,
		"refresh_token": newRefresh,
		"expires_at":    time.Now().Add(15 * time.Minute),
	})
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	var body RefreshRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
	}
	if err := h.authUC.Logout(body.RefreshToken); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "logged out"})
}

// RegisterRoutes สำหรับ auth
func (h *AuthHandler) RegisterRoutes(app *fiber.App) {
	apiAuth := app.Group("/auth")
	apiAuth.Post("/register", h.Register)
	apiAuth.Post("/login", h.Login)
	apiAuth.Post("/refresh", h.Refresh)
	apiAuth.Post("/logout", h.Logout)
}
