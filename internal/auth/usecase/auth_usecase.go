package usecase

import (
	"errors"
	"time"

	"invoice_project/internal/auth/domain"
	"invoice_project/internal/auth/repository"
	"invoice_project/pkg/middleware"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase interface {
	Register(username, password string) error
	Login(username, password string) (accessToken string, refreshToken string, err error)
	RefreshAccessToken(oldRefreshToken string) (newAccessToken string, newRefreshToken string, err error)
	Logout(refreshToken string) error
}

type authUC struct {
	repo          repository.AuthRepository
	jwtSecret     string
	jwtIssuer     string
	jwtAudience   string
	accessExpiry  time.Duration
	refreshExpiry time.Duration
	validate      *validator.Validate
}

func NewAuthUsecase(repo repository.AuthRepository, jwtSecret, jwtIssuer, jwtAudience string, accessMin int, refreshHours int) AuthUsecase {
	return &authUC{
		repo:          repo,
		jwtSecret:     jwtSecret,
		jwtIssuer:     jwtIssuer,
		jwtAudience:   jwtAudience,
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
		return err
	}
	existing, err := u.repo.GetUserByUsername(username)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New("username already taken")
	}
	user := &domain.User{
		Username: username,
		Password: password,
		Role:     "user",
	}
	return u.repo.CreateUser(user)
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
		return "", "", err
	}
	user, err := u.repo.GetUserByUsername(username)
	if err != nil {
		return "", "", err
	}
	if user == nil {
		return "", "", errors.New("invalid credentials")
	}
	// เปรียบเทียบรหัสผ่าน
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", "", errors.New("invalid credentials")
	}
	// สร้าง access token
	accessToken, err := middleware.GenerateJWTWithExpiry(u.jwtSecret, user.ID, user.Username, user.Role, u.jwtIssuer, []string{u.jwtAudience}, u.accessExpiry)
	if err != nil {
		return "", "", err
	}
	// สร้าง refresh token (random string หรือ JWT ก็ได้ ในที่นี้ใช้ JWT ง่าย ๆ)
	refreshToken, err := middleware.GenerateJWTWithExpiry(u.jwtSecret, user.ID, user.Username, user.Role, u.jwtIssuer, []string{u.jwtAudience}, u.refreshExpiry)
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
		return "", "", errors.New("refresh token not found or revoked")
	}
	// เช็คว่า expired หรือยัง
	if existing.ExpiredAt.Before(time.Now()) {
		// ถ้า expired ให้ลบทิ้งและคืน error
		if err := u.repo.DeleteRefreshToken(oldRefreshToken); err != nil {
			return "", "", err
		}
		return "", "", errors.New("refresh token expired")
	}
	// ถ้า valid ก็สร้าง access + refresh ใหม่
	user := &existing.User
	newAccess, err := middleware.GenerateJWTWithExpiry(u.jwtSecret, user.ID, user.Username, user.Role, u.jwtIssuer, []string{u.jwtAudience}, u.accessExpiry)
	if err != nil {
		return "", "", err
	}
	newRefresh, err := middleware.GenerateJWTWithExpiry(u.jwtSecret, user.ID, user.Username, user.Role, u.jwtIssuer, []string{u.jwtAudience}, u.refreshExpiry)
	if err != nil {
		return "", "", err
	}
	// อัปเดตใน DB: ลบ old แล้วเพิ่ม new
	if err := u.repo.DeleteRefreshToken(oldRefreshToken); err != nil {
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
	return u.repo.DeleteRefreshToken(refreshToken)
}
