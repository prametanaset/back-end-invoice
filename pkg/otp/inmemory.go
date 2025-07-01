package otp

import (
	"context"
	"fmt"
	"time"
)

// InMemoryOTPService stores OTPs in memory and prints them to stdout.
type InMemoryOTPService struct {
	otps map[string]otpEntry
}

// NewInMemoryOTPService creates a new InMemoryOTPService instance.
func NewInMemoryOTPService() *InMemoryOTPService {
	return &InMemoryOTPService{otps: make(map[string]otpEntry)}
}

// SendOTP generates an OTP and logs it. It returns the generated code.
func (s *InMemoryOTPService) SendOTP(ctx context.Context, to string) (string, error) {
	code, err := generateCode()
	if err != nil {
		return "", err
	}
	fmt.Printf("sending OTP %s to %s\n", code, to)
	s.otps[to] = otpEntry{Code: code, ExpiresAt: time.Now().Add(5 * time.Minute)}
	return code, nil
}

// VerifyOTP checks the OTP code for the given receiver.
func (s *InMemoryOTPService) VerifyOTP(to, code string) bool {
	entry, ok := s.otps[to]
	if !ok || time.Now().After(entry.ExpiresAt) {
		return false
	}
	if entry.Code != code {
		return false
	}
	delete(s.otps, to)
	return true
}
