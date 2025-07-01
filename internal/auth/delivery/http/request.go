package http

// RegisterRequest represents the expected payload for user registration.
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
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
	Email string `json:"email"`
}
