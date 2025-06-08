package repository

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"invoice_project/internal/auth/domain"
)

type AuthRepository interface {
	CreateUser(user *domain.User) error
	GetUserByUsername(username string) (*domain.User, error)
	SaveRefreshToken(token *domain.RefreshToken) error
	GetRefreshToken(rawToken string) (*domain.RefreshToken, error)
	DeleteRefreshToken(rawToken string) error
	DeleteAllRefreshTokensForUser(userID uint) error
}

type authPG struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authPG{db: db}
}

func (r *authPG) CreateUser(user *domain.User) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashed)
	if user.Role == "" {
		user.Role = "user"
	}
	return r.db.Create(user).Error
}

func (r *authPG) GetUserByUsername(username string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *authPG) SaveRefreshToken(token *domain.RefreshToken) error {
	return r.db.Create(token).Error
}

func (r *authPG) GetRefreshToken(rawToken string) (*domain.RefreshToken, error) {
	var rt domain.RefreshToken
	err := r.db.
		Preload("User").
		Where("token = ?", rawToken).
		First(&rt).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &rt, nil
}

func (r *authPG) DeleteRefreshToken(rawToken string) error {
	return r.db.Where("token = ?", rawToken).Delete(&domain.RefreshToken{}).Error
}

func (r *authPG) DeleteAllRefreshTokensForUser(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&domain.RefreshToken{}).Error
}
