package otp

import (
	"testing"
	"time"
)

func TestGenerateCode(t *testing.T) {
	code, err := generateCode()
	if err != nil {
		t.Fatalf("generateCode returned error: %v", err)
	}
	if len(code) != 6 {
		t.Errorf("expected code length 6, got %d", len(code))
	}
}

func TestVerifyOTP(t *testing.T) {
	svc := &GmailOTPService{otps: make(map[string]otpEntry)}
	svc.otps["user@example.com"] = otpEntry{Code: "123456", ExpiresAt: time.Now().Add(1 * time.Minute)}

	if !svc.VerifyOTP("user@example.com", "123456") {
		t.Errorf("expected OTP verification success")
	}
	if svc.VerifyOTP("user@example.com", "123456") {
		t.Errorf("OTP should not verify twice")
	}
	svc.otps["user2@example.com"] = otpEntry{Code: "999999", ExpiresAt: time.Now().Add(-1 * time.Minute)}
	if svc.VerifyOTP("user2@example.com", "999999") {
		t.Errorf("expired OTP should fail")
	}
}
