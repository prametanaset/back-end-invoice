package usecase

import (
	"time"

	"invoice_project/internal/auth/domain"
	"invoice_project/internal/auth/repository"
	"invoice_project/pkg/apperror"
	"invoice_project/pkg/middleware"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase interface {
	Register(username, password string) error
	Login(username, password string) (accessToken string, refreshToken string, err error)
	OAuthLogin(provider, providerUID, username string) (accessToken string, refreshToken string, err error)
	RefreshAccessToken(oldRefreshToken string) (newAccessToken string, newRefreshToken string, err error)
	Logout(refreshToken string) error
	GetProfile(userID uuid.UUID) (*domain.User, error)
}

type authUC struct {
	repo          repository.AuthRepository
	jwtSecret     string
	accessExpiry  time.Duration
	refreshExpiry time.Duration
	validate      *validator.Validate
}

func NewAuthUsecase(repo repository.AuthRepository, jwtSecret string, accessMin int, refreshHours int) AuthUsecase {
	return &authUC{
		repo:          repo,
		jwtSecret:     jwtSecret,
		accessExpiry:  time.Duration(accessMin) * time.Minute,
		refreshExpiry: time.Duration(refreshHours) * time.Hour,
		validate:      validator.New(),
	}
}

func (u *authUC) Register(username, password string) error {
	// เช็คให้ username กับ password ไม่ว่าง และยาวไม่น้อย
	input := struct {
		Username string `validate:"required,min=3"`
		Password string `validate:"required,min=6"`
	}{
		Username: username,
		Password: password,
	}
	if err := u.validate.Struct(input); err != nil {
		return apperror.New(fiber.StatusBadRequest)
	}
	existing, err := u.repo.GetUserByUsername(username)
	if err != nil {
		return err
	}
	if existing != nil {
		return apperror.New(fiber.StatusBadRequest)
	}
	user := &domain.User{
		Username:     username,
		PasswordHash: password,
	}
	if err := u.repo.CreateUser(user); err != nil {
		return err
	}
	if err := u.repo.AssignRoleToUser(user.ID, "user"); err != nil {
		return err
	}
	return nil
}

func (u *authUC) Login(username, password string) (string, string, error) {
	// validate input
	input := struct {
		Username string `validate:"required,min=3"`
		Password string `validate:"required,min=6"`
	}{
		Username: username,
		Password: password,
	}
	if err := u.validate.Struct(input); err != nil {
		return "", "", apperror.New(fiber.StatusBadRequest)
	}
	user, err := u.repo.GetUserByUsername(username)
	if err != nil {
		return "", "", err
	}
	if user == nil {
		return "", "", apperror.New(fiber.StatusUnauthorized)
	}
	// เปรียบเทียบรหัสผ่าน
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", "", apperror.New(fiber.StatusUnauthorized)
	}
	role, err := u.repo.GetPrimaryRole(user.ID)
	if err != nil {
		return "", "", err
	}
	// สร้าง access token
	accessToken, err := middleware.GenerateJWTWithExpiry(u.jwtSecret, user.ID, role, u.accessExpiry, "access")
	if err != nil {
		return "", "", err
	}
	// สร้าง refresh token (random string หรือ JWT ก็ได้ ในที่นี้ใช้ JWT ง่าย ๆ)
	refreshToken, err := middleware.GenerateJWTWithExpiry(u.jwtSecret, user.ID, role, u.refreshExpiry, "refresh")
	if err != nil {
		return "", "", err
	}
	// บันทึก refresh token ลง DB
	sess := &domain.UserSession{
		UserID:       user.ID,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(u.refreshExpiry),
	}
	if err := u.repo.CreateSession(sess); err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

// OAuthLogin authenticates or registers a user using an external provider ID.
func (u *authUC) OAuthLogin(provider, providerUID, username string) (string, string, error) {
	// Check if login method already exists
	user, err := u.repo.GetUserByLoginMethod(provider, providerUID)
	if err != nil {
		return "", "", err
	}
	if user == nil {
		// Create user if not exists
		user = &domain.User{
			Username:     username,
			PasswordHash: uuid.NewString(), // random password
			IsVerified:   true,
		}
		if err := u.repo.CreateUser(user); err != nil {
			return "", "", err
		}
		if err := u.repo.AssignRoleToUser(user.ID, "user"); err != nil {
			return "", "", err
		}
		lm := &domain.UserLoginMethod{
			UserID:      user.ID,
			Provider:    provider,
			ProviderUID: providerUID,
		}
		if err := u.repo.CreateLoginMethod(lm); err != nil {
			return "", "", err
		}
	}
	role, err := u.repo.GetPrimaryRole(user.ID)
	if err != nil {
		return "", "", err
	}
	accessToken, err := middleware.GenerateJWTWithExpiry(u.jwtSecret, user.ID, role, u.accessExpiry, "access")
	if err != nil {
		return "", "", err
	}
	refreshToken, err := middleware.GenerateJWTWithExpiry(u.jwtSecret, user.ID, role, u.refreshExpiry, "refresh")
	if err != nil {
		return "", "", err
	}
	sess := &domain.UserSession{
		UserID:       user.ID,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(u.refreshExpiry),
	}
	if err := u.repo.CreateSession(sess); err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func (u *authUC) RefreshAccessToken(oldRefreshToken string) (string, string, error) {
	// หา record จาก DB
	existing, err := u.repo.GetSessionByToken(oldRefreshToken)
	if err != nil {
		return "", "", err
	}
	if existing == nil {
		return "", "", apperror.New(fiber.StatusUnauthorized)
	}
	if existing.RevokedAt != nil {
		return "", "", apperror.New(fiber.StatusUnauthorized)
	}
	// เช็คว่า expired หรือยัง
	if existing.ExpiresAt.Before(time.Now()) {
		// ถ้า expired ให้ลบทิ้งและคืน error
		if err := u.repo.RevokeSession(oldRefreshToken); err != nil {
			return "", "", err
		}
		return "", "", apperror.New(fiber.StatusUnauthorized)
	}
	// ถ้า valid ก็สร้าง access + refresh ใหม่
	user := &existing.User
	role, err := u.repo.GetPrimaryRole(user.ID)
	if err != nil {
		return "", "", err
	}
	newAccess, err := middleware.GenerateJWTWithExpiry(u.jwtSecret, user.ID, role, u.accessExpiry, "access")
	if err != nil {
		return "", "", err
	}
	newRefresh, err := middleware.GenerateJWTWithExpiry(u.jwtSecret, user.ID, role, u.refreshExpiry, "refresh")
	if err != nil {
		return "", "", err
	}
	// อัปเดตใน DB: ลบ old แล้วเพิ่ม new
	if err := u.repo.RevokeSession(oldRefreshToken); err != nil {
		return "", "", err
	}
	newSess := &domain.UserSession{
		UserID:       user.ID,
		RefreshToken: newRefresh,
		ExpiresAt:    time.Now().Add(u.refreshExpiry),
	}
	if err := u.repo.CreateSession(newSess); err != nil {
		return "", "", err
	}
	return newAccess, newRefresh, nil
}

func (u *authUC) Logout(refreshToken string) error {
	// Mark session revoked
	return u.repo.RevokeSession(refreshToken)
}

func (u *authUC) GetProfile(userID uuid.UUID) (*domain.User, error) {
	user, err := u.repo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, apperror.New(fiber.StatusNotFound)
	}
	return user, nil
}
