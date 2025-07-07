package http

// RegisterRequest represents the expected payload for user registration.
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	OTPRef   string `json:"otp_ref"`
	OTPCode  string `json:"otp_code"`
}

// LoginRequest represents the expected payload for user login.
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// RefreshRequest represents the expected payload for refreshing an access token.
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// OAuthLoginRequest represents OAuth login payload.
type OAuthLoginRequest struct {
	Provider    string `json:"provider"`
	ProviderUID string `json:"provider_uid"`
	Username    string `json:"username"`
}

// LogoutRequest represents the expected payload for logging out.
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// CheckEmailRequest represents the payload to verify if an email is already registered.
type CheckEmailRequest struct {
	Username string `json:"username"`
}

// SendOTPRequest represents the payload to request an OTP to be sent.
type SendOTPRequest struct {
	Email   string `json:"email"`
	Purpose string `json:"purpose"`
}

// VerifyOTPRequest represents the payload for verifying an OTP code.
type VerifyOTPRequest struct {
	Email       string `json:"email"`
	Ref         string `json:"ref"`
	Code        string `json:"code"`
	Purpose     string `json:"purpose"`
	NewPassword string `json:"new_password"`
}

// ResetPasswordRequest represents the payload for resetting a password
// using an OTP reference and code.
type ResetPasswordRequest struct {
	Email       string `json:"email"`
	Ref         string `json:"ref"`
	Code        string `json:"code"`
	NewPassword string `json:"new_password"`
}
