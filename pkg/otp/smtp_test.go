package otp

import (
	"context"
	"net/smtp"
	"testing"
	"time"
)

func TestSMTPSendOTP(t *testing.T) {
	var sent bool
	svc := NewSMTPOTPService("smtp.example.com", 587, "user", "pass", "from@example.com")
	svc.sendMail = func(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
		sent = true
		return nil
	}
	code, err := svc.SendOTP(context.Background(), "to@example.com")
	if err != nil {
		t.Fatalf("SendOTP returned error: %v", err)
	}
	if len(code) != 6 {
		t.Errorf("expected code length 6, got %d", len(code))
	}
	if !sent {
		t.Errorf("sendMail not called")
	}
	if _, ok := svc.otps["to@example.com"]; !ok {
		t.Errorf("OTP not stored")
	}
}

func TestSMTPVerifyOTP(t *testing.T) {
	svc := &SMTPOTPService{otps: make(map[string]otpEntry)}
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
