package otp

import (
	"context"
	"fmt"
	"net/smtp"
	"time"
)

// SMTPOTPService implements the Service interface using a generic SMTP server.
type SMTPOTPService struct {
	host      string
	port      int
	username  string
	password  string
	fromEmail string
	otps      map[string]otpEntry
	sendMail  func(addr string, a smtp.Auth, from string, to []string, msg []byte) error
}

// NewSMTPOTPService creates a new SMTPOTPService instance.
func NewSMTPOTPService(host string, port int, username, password, from string) *SMTPOTPService {
	return &SMTPOTPService{
		host:      host,
		port:      port,
		username:  username,
		password:  password,
		fromEmail: from,
		otps:      make(map[string]otpEntry),
		sendMail:  smtp.SendMail,
	}
}

// SendOTP generates an OTP code and delivers it using SMTP.
func (s *SMTPOTPService) SendOTP(ctx context.Context, to, ref string) (string, error) {
	code, err := generateCode()
	if err != nil {
		return "", err
	}
	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	auth := smtp.PlainAuth("", s.username, s.password, s.host)
	msg := []byte(buildOTPEmail(to, code, ref))
	if err := s.sendMail(addr, auth, s.fromEmail, []string{to}, msg); err != nil {
		return "", err
	}
	s.otps[to] = otpEntry{Code: code, ExpiresAt: time.Now().Add(5 * time.Minute)}
	return code, nil
}

// VerifyOTP verifies the code for the given recipient.
func (s *SMTPOTPService) VerifyOTP(to, code string) bool {
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
