package http

import (
	"net/http/httptest"
	"testing"
	"time"

	"invoice_project/internal/auth/domain"
	"invoice_project/internal/auth/usecase"
	"invoice_project/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

type stubRepo struct{}

func (s *stubRepo) CreateUser(user *domain.User) error { return nil }
func (s *stubRepo) GetUserByUsername(username string) (*domain.User, error) {
	return nil, nil
}
func (s *stubRepo) GetUserByID(id uint) (*domain.User, error) {
	if id == 1 {
		return &domain.User{ID: 1, Username: "tester", Role: "user"}, nil
	}
	return nil, nil
}
func (s *stubRepo) SaveRefreshToken(token *domain.RefreshToken) error        { return nil }
func (s *stubRepo) GetRefreshToken(raw string) (*domain.RefreshToken, error) { return nil, nil }
func (s *stubRepo) RevokeRefreshToken(raw string) error                      { return nil }
func (s *stubRepo) DeleteAllRefreshTokensForUser(uid uint) error             { return nil }

func TestAuthHandler_Me(t *testing.T) {
	uc := usecase.NewAuthUsecase(&stubRepo{}, "secret", 15, 24)
	h := NewAuthHandler(uc, "secret")

	app := fiber.New()
	app.Use(middleware.JWTMiddlewareExcept("secret"))
	h.RegisterRoutes(app)

	token, err := middleware.GenerateJWTWithExpiry("secret", 1, "user", time.Minute, "access")
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	req := httptest.NewRequest("GET", "/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test error: %v", err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}
