package usecase

import (
	"context"
	"time"

	"invoice_project/internal/auth/domain"
	"invoice_project/internal/auth/repository"
	"invoice_project/pkg/apperror"
	"invoice_project/pkg/otp"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// OTPUsecase defines sending and verifying OTP codes.
type OTPUsecase interface {
	SendOTP(ctx context.Context, email string) error
	VerifyOTP(ctx context.Context, email, code string) error
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

const otpPurposeVerifyEmail = "verify_email"

func (u *otpUC) SendOTP(ctx context.Context, email string) error {
	user, err := u.authRepo.GetUserByUsername(email)
	if err != nil {
		return err
	}
	if user == nil {
		return apperror.New(fiber.StatusNotFound)
	}
	code, err := u.svc.SendOTP(ctx, email)
	if err != nil {
		return err
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(code), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	otp := &domain.OTP{
		Purpose:     otpPurposeVerifyEmail,
		Destination: email,
		CodeHash:    string(hashed),
		ExpiresAt:   time.Now().Add(5 * time.Minute),
	}
	return u.otpRepo.CreateOTP(otp)
}

func (u *otpUC) VerifyOTP(ctx context.Context, email, code string) error {
	user, err := u.authRepo.GetUserByUsername(email)
	if err != nil {
		return err
	}
	if user == nil {
		return apperror.New(fiber.StatusNotFound)
	}
	otpRec, err := u.otpRepo.GetActiveOTP(email, otpPurposeVerifyEmail)
	if err != nil {
		return err
	}
	if otpRec == nil {
		return apperror.New(fiber.StatusBadRequest)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(otpRec.CodeHash), []byte(code)); err != nil {
		_ = u.otpRepo.IncrementAttempts(otpRec.ID)
		return apperror.New(fiber.StatusBadRequest)
	}
	if err := u.otpRepo.MarkUsed(otpRec.ID); err != nil {
		return err
	}
	if err := u.authRepo.SetUserVerified(user.ID); err != nil {
		return err
	}
	return nil
}
