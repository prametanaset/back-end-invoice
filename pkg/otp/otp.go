package otp

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

// Service defines OTP sending and verification behaviour.
type Service interface {
	SendOTP(ctx context.Context, to string) (string, error)
	VerifyOTP(to, code string) bool
}

// GmailOTPService implements the Service interface using Gmail API.
type GmailOTPService struct {
	srv       *gmail.Service
	fromEmail string
	otps      map[string]otpEntry
}

type otpEntry struct {
	Code      string
	ExpiresAt time.Time
}

func NewGmailOTPService(ctx context.Context, credentialsJSON, tokenJSON []byte, from string) (*GmailOTPService, error) {
	config, err := google.ConfigFromJSON(credentialsJSON, gmail.GmailSendScope)
	if err != nil {
		return nil, err
	}
	var token oauth2.Token
	if err := json.Unmarshal(tokenJSON, &token); err != nil {
		return nil, err
	}
	srv, err := gmail.NewService(ctx, option.WithTokenSource(config.TokenSource(ctx, &token)))
	if err != nil {
		return nil, err
	}
	return &GmailOTPService{srv: srv, fromEmail: from, otps: make(map[string]otpEntry)}, nil
}

func generateCode() (string, error) {
	b := make([]byte, 3)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	// Limit to 6 digits
	n := int(b[0])<<16 | int(b[1])<<8 | int(b[2])
	return fmt.Sprintf("%06d", n%1000000), nil
}

func (g *GmailOTPService) SendOTP(ctx context.Context, to string) (string, error) {
	code, err := generateCode()
	if err != nil {
		return "", err
	}
	msgStr := fmt.Sprintf("To: %s\r\nSubject: Your OTP Code\r\n\r\nYour OTP is %s", to, code)
	msg := &gmail.Message{Raw: base64.URLEncoding.EncodeToString([]byte(msgStr))}
	if _, err := g.srv.Users.Messages.Send("me", msg).Do(); err != nil {
		return "", err
	}
	g.otps[to] = otpEntry{Code: code, ExpiresAt: time.Now().Add(5 * time.Minute)}
	return code, nil
}

func (g *GmailOTPService) VerifyOTP(to, code string) bool {
	entry, ok := g.otps[to]
	if !ok || time.Now().After(entry.ExpiresAt) {
		return false
	}
	if entry.Code != code {
		return false
	}
	delete(g.otps, to)
	return true
}
