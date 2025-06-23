package http

import (
	"time"

	"invoice_project/internal/auth/usecase"
	merchUC "invoice_project/internal/merchant/usecase"
	"invoice_project/pkg/apperror"
	"invoice_project/pkg/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AuthHandler struct {
	authUC     usecase.AuthUsecase
	merchantUC merchUC.MerchantUsecase
	jwtSecret  string
}

func NewAuthHandler(authUC usecase.AuthUsecase, merchantUC merchUC.MerchantUsecase, jwtSecret string) *AuthHandler {
	return &AuthHandler{
		authUC:     authUC,
		merchantUC: merchantUC,
		jwtSecret:  jwtSecret,
	}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var body RegisterRequest
	if err := c.BodyParser(&body); err != nil {
		return apperror.New(fiber.StatusBadRequest)
	}
	if err := h.authUC.Register(body.Username, body.Password); err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "user registered"})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var body LoginRequest
	if err := c.BodyParser(&body); err != nil {
		return apperror.New(fiber.StatusBadRequest)
	}
	accessToken, refreshToken, err := h.authUC.Login(body.Username, body.Password)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"expires_at":    time.Now().Add(15 * time.Minute),
	})
}

func (h *AuthHandler) OAuthLogin(c *fiber.Ctx) error {
	var body OAuthLoginRequest
	if err := c.BodyParser(&body); err != nil {
		return apperror.New(fiber.StatusBadRequest)
	}
	if body.Provider == "" || body.ProviderUID == "" {
		return apperror.New(fiber.StatusBadRequest)
	}
	username := body.Username
	if username == "" {
		username = body.Provider + "_" + body.ProviderUID
	}
	accessToken, refreshToken, err := h.authUC.OAuthLogin(body.Provider, body.ProviderUID, username)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"expires_at":    time.Now().Add(15 * time.Minute),
	})
}

func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	var body RefreshRequest
	if err := c.BodyParser(&body); err != nil {
		return apperror.New(fiber.StatusBadRequest)
	}
	newAccess, newRefresh, err := h.authUC.RefreshAccessToken(body.RefreshToken)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"access_token":  newAccess,
		"refresh_token": newRefresh,
		"expires_at":    time.Now().Add(15 * time.Minute),
	})
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	var body LogoutRequest
	if err := c.BodyParser(&body); err != nil {
		return apperror.New(fiber.StatusBadRequest)
	}
	if body.RefreshToken == "" {
		return apperror.New(fiber.StatusBadRequest)
	}

	if err := h.authUC.Logout(body.RefreshToken); err != nil {
		return err
	}
	return c.JSON(fiber.Map{"message": "logged out"})
}

func (h *AuthHandler) CheckEmail(c *fiber.Ctx) error {
	var body CheckEmailRequest
	if err := c.BodyParser(&body); err != nil {
		return apperror.New(fiber.StatusBadRequest)
	}
	if body.Username == "" {
		return apperror.New(fiber.StatusBadRequest)
	}
	taken, err := h.authUC.IsUsernameTaken(body.Username)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"taken": taken})
}

func (h *AuthHandler) Me(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return apperror.New(fiber.StatusUnauthorized)
	}
	user, err := h.authUC.GetProfile(userID)
	if err != nil {
		return err
	}

	resp := fiber.Map{"user": user}

	m, err := h.merchantUC.GetMyMerchant(userID)
	if err != nil {
		return err
	}
	if m != nil {
		merchResp := fiber.Map{"merchant": m}
		switch m.MerchantType {
		case "person":
			p, err := h.merchantUC.GetPerson(m.ID)
			if err != nil {
				return err
			}
			if p != nil {
				merchResp["person"] = p
			}
		case "company":
			comp, err := h.merchantUC.GetCompany(m.ID)
			if err != nil {
				return err
			}
			if comp != nil {
				merchResp["company"] = comp
			}
		}
		resp["merchant_info"] = merchResp
	}

	return c.JSON(resp)
}

// RegisterRoutes สำหรับ auth
func (h *AuthHandler) RegisterRoutes(app *fiber.App) {
	apiAuth := app.Group("/auth")
	apiAuth.Post("/register", h.Register)
	apiAuth.Post("/check-email", h.CheckEmail)
	apiAuth.Post("/login", h.Login)
	apiAuth.Post("/oauth-login", h.OAuthLogin)
	apiAuth.Post("/refresh", h.Refresh)
	apiAuth.Post("/logout", h.Logout)
	app.Get("/me", middleware.RequireRoles("user", "admin"), h.Me)
}
