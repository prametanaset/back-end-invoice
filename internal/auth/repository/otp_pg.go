package repository

import (
	"errors"
	"time"

	"invoice_project/internal/auth/domain"

	"gorm.io/gorm"
)

// OTPRepository manages OTP persistence.
type OTPRepository interface {
	CreateOTP(o *domain.OTP) error
	GetActiveOTP(dest, purpose string) (*domain.OTP, error)
	MarkUsed(id uint64) error
	IncrementAttempts(id uint64) error
	RevokeOTP(id uint64) error
}

type otpPG struct{ db *gorm.DB }

// NewOTPRepository returns a PostgreSQL implementation of OTPRepository.
func NewOTPRepository(db *gorm.DB) OTPRepository {
	return &otpPG{db: db}
}

func (r *otpPG) CreateOTP(o *domain.OTP) error {
	return r.db.Create(o).Error
}

func (r *otpPG) GetActiveOTP(dest, purpose string) (*domain.OTP, error) {
	var otp domain.OTP
	err := r.db.Where("destination = ? AND purpose = ? AND used_at IS NULL AND revoked_at IS NULL AND expires_at > ?",
		dest, purpose, time.Now()).Order("created_at desc").First(&otp).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &otp, nil
}

func (r *otpPG) MarkUsed(id uint64) error {
	return r.db.Model(&domain.OTP{}).Where("id = ?", id).Update("used_at", time.Now()).Error
}

func (r *otpPG) IncrementAttempts(id uint64) error {
	return r.db.Model(&domain.OTP{}).Where("id = ?", id).UpdateColumn("attempts", gorm.Expr("attempts + 1")).Error
}

func (r *otpPG) RevokeOTP(id uint64) error {
	return r.db.Model(&domain.OTP{}).Where("id = ?", id).Update("revoked_at", time.Now()).Error
}
