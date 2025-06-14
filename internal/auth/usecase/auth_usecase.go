package usecase

import (
	"time"

	"invoice_project/internal/auth/domain"
	"invoice_project/internal/auth/repository"
	"invoice_project/pkg/apperror"
	"invoice_project/pkg/middleware"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase interface {
	Register(username, password string) error
	Login(username, password string) (accessToken string, refreshToken string, err error)
	RefreshAccessToken(oldRefreshToken string) (newAccessToken string, newRefreshToken string, err error)
	Logout(refreshToken string) error
	GetProfile(userID uint) (*domain.User, error)
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
		Username: username,
		Password: password,
		Role:     "user",
	}
	if err := u.repo.CreateUser(user); err != nil {
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
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", "", apperror.New(fiber.StatusUnauthorized)
	}
	// สร้าง access token
	accessToken, err := middleware.GenerateJWTWithExpiry(u.jwtSecret, user.ID, user.Role, u.accessExpiry, "access")
	if err != nil {
		return "", "", err
	}
	// สร้าง refresh token (random string หรือ JWT ก็ได้ ในที่นี้ใช้ JWT ง่าย ๆ)
	refreshToken, err := middleware.GenerateJWTWithExpiry(u.jwtSecret, user.ID, user.Role, u.refreshExpiry, "refresh")
	if err != nil {
		return "", "", err
	}
	// บันทึก refresh token ลง DB
	rt := &domain.RefreshToken{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiredAt: time.Now().Add(u.refreshExpiry),
	}
	if err := u.repo.SaveRefreshToken(rt); err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func (u *authUC) RefreshAccessToken(oldRefreshToken string) (string, string, error) {
	// หา record จาก DB
	existing, err := u.repo.GetRefreshToken(oldRefreshToken)
	if err != nil {
		return "", "", err
	}
	if existing == nil {
		return "", "", apperror.New(fiber.StatusUnauthorized)
	}
	// เช็คว่า expired หรือยัง
	if existing.ExpiredAt.Before(time.Now()) {
		// ถ้า expired ให้ลบทิ้งและคืน error
		if err := u.repo.RevokeRefreshToken(oldRefreshToken); err != nil {
			return "", "", err
		}
		return "", "", apperror.New(fiber.StatusUnauthorized)
	}
	// ถ้า valid ก็สร้าง access + refresh ใหม่
	user := &existing.User
	newAccess, err := middleware.GenerateJWTWithExpiry(u.jwtSecret, user.ID, user.Role, u.accessExpiry, "access")
	if err != nil {
		return "", "", err
	}
	newRefresh, err := middleware.GenerateJWTWithExpiry(u.jwtSecret, user.ID, user.Role, u.refreshExpiry, "refresh")
	if err != nil {
		return "", "", err
	}
	// อัปเดตใน DB: ลบ old แล้วเพิ่ม new
	if err := u.repo.RevokeRefreshToken(oldRefreshToken); err != nil {
		return "", "", err
	}
	newRT := &domain.RefreshToken{
		Token:     newRefresh,
		UserID:    user.ID,
		ExpiredAt: time.Now().Add(u.refreshExpiry),
	}
	if err := u.repo.SaveRefreshToken(newRT); err != nil {
		return "", "", err
	}
	return newAccess, newRefresh, nil
}

func (u *authUC) Logout(refreshToken string) error {
	// ลบ refresh token ออกจาก DB เลย
	return u.repo.RevokeRefreshToken(refreshToken)
}

func (u *authUC) GetProfile(userID uint) (*domain.User, error) {
	user, err := u.repo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, apperror.New(fiber.StatusNotFound)
	}
	return user, nil
}
