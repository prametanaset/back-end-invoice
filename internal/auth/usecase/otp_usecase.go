package usecase

import (
	"context"
	"time"

	"invoice_project/internal/auth/domain"
	"invoice_project/internal/auth/repository"
	"invoice_project/pkg/apperror"
	"invoice_project/pkg/otp"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// OTPUsecase defines sending and verifying OTP codes.
type OTPUsecase interface {
	SendOTP(ctx context.Context, email, purpose string) (string, error)
	VerifyOTP(ctx context.Context, email, ref, code, purpose, newPassword string) (uuid.UUID, error)
	ResetPassword(userID uuid.UUID, newPassword string) error
}

type otpUC struct {
	authRepo repository.AuthRepository
	otpRepo  repository.OTPRepository
	svc      otp.Service
}

// NewOTPUsecase creates a new OTP usecase implementation.
func NewOTPUsecase(authRepo repository.AuthRepository, otpRepo repository.OTPRepository, svc otp.Service) OTPUsecase {
	return &otpUC{authRepo: authRepo, otpRepo: otpRepo, svc: svc}
}

const maxOTPAttempts = 5

func (u *otpUC) SendOTP(ctx context.Context, email, purpose string) (string, error) {
	if !domain.IsValidOTPPurpose(purpose) {
		return "", apperror.New(fiber.StatusBadRequest)
	}
	if purpose != string(domain.OTPPurposeVerifyEmail) {
		user, err := u.authRepo.GetUserByUsername(email)
		if err != nil {
			return "", err
		}
		if user == nil {
			return "", apperror.New(fiber.StatusNotFound)
		}
	}
	ref := uuid.NewString()
	code, err := u.svc.SendOTP(ctx, email, ref)
	if err != nil {
		return "", err
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(code), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	otp := &domain.OTP{
		Purpose:     purpose,
		Ref:         ref,
		Destination: email,
		CodeHash:    string(hashed),
		ExpiresAt:   time.Now().Add(5 * time.Minute),
	}
	return ref, u.otpRepo.CreateOTP(otp)
}

func (u *otpUC) VerifyOTP(ctx context.Context, email, ref, code, purpose, newPassword string) (uuid.UUID, error) {
	if !domain.IsValidOTPPurpose(purpose) {
		return uuid.Nil, apperror.New(fiber.StatusBadRequest)
	}
	var user *domain.User
	var err error
	if purpose != string(domain.OTPPurposeVerifyEmail) {
		user, err = u.authRepo.GetUserByUsername(email)
		if err != nil {
			return uuid.Nil, err
		}
		if user == nil {
			return uuid.Nil, apperror.New(fiber.StatusNotFound)
		}
	}
	otpRec, err := u.otpRepo.GetActiveOTP(email, purpose, ref)
	if err != nil {
		return uuid.Nil, err
	}
	if otpRec == nil {
		return uuid.Nil, apperror.New(fiber.StatusBadRequest)
	}
	if otpRec.Attempts >= maxOTPAttempts {
		_ = u.otpRepo.RevokeOTP(otpRec.ID)
		return uuid.Nil, apperror.New(fiber.StatusBadRequest)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(otpRec.CodeHash), []byte(code)); err != nil {
		_ = u.otpRepo.IncrementAttempts(otpRec.ID)
		if otpRec.Attempts+1 >= maxOTPAttempts {
			_ = u.otpRepo.RevokeOTP(otpRec.ID)
		}
		return uuid.Nil, apperror.New(fiber.StatusBadRequest)
	}
	if err := u.otpRepo.MarkUsed(otpRec.ID); err != nil {
		return uuid.Nil, err
	}
	if purpose == string(domain.OTPPurposeResetPassword) && newPassword != "" {
		if err := u.authRepo.UpdatePassword(user.ID, newPassword); err != nil {
			return uuid.Nil, err
		}
	}
	if purpose == string(domain.OTPPurposeResetPassword) {
		return user.ID, nil
	}
	return uuid.Nil, nil
}

func (u *otpUC) ResetPassword(userID uuid.UUID, newPassword string) error {
	if newPassword == "" {
		return apperror.New(fiber.StatusBadRequest)
	}
	return u.authRepo.UpdatePassword(userID, newPassword)
}
