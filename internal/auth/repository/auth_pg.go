package repository

import (
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"invoice_project/internal/auth/domain"
)

type AuthRepository interface {
	CreateUser(user *domain.User) error
	GetUserByUsername(username string) (*domain.User, error)
	GetUserByID(id uuid.UUID) (*domain.User, error)
	CreateSession(session *domain.UserSession) error
	GetSessionByToken(token string) (*domain.UserSession, error)
	RevokeSession(token string) error
	DeleteAllSessionsForUser(userID uuid.UUID) error
}

type authPG struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authPG{db: db}
}

func (r *authPG) CreateUser(user *domain.User) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashed)
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

func (r *authPG) GetUserByID(id uuid.UUID) (*domain.User, error) {
	var user domain.User
	err := r.db.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *authPG) CreateSession(s *domain.UserSession) error {
	return r.db.Create(s).Error
}

func (r *authPG) GetSessionByToken(token string) (*domain.UserSession, error) {
	var sess domain.UserSession
	err := r.db.
		Preload("User").
		Where("refresh_token = ?", token).
		First(&sess).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &sess, nil
}

func (r *authPG) RevokeSession(token string) error {
	return r.db.Model(&domain.UserSession{}).
		Where("refresh_token = ?", token).
		Update("revoked", true).Error
}

func (r *authPG) DeleteAllSessionsForUser(userID uuid.UUID) error {
	return r.db.Where("user_id = ?", userID).Delete(&domain.UserSession{}).Error
}
