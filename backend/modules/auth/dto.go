package auth

// LoginRequest is the request body for login
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// LoginResponse is the response body for login
type LoginResponse struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	User         UserInfo `json:"user"`
}

// UserInfo is the public user information
type UserInfo struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Email     string  `json:"email"`
	Role      string  `json:"role"`
	NIM       *string `json:"nim,omitempty"`
	NIP       *string `json:"nip,omitempty"`
	AvatarURL *string `json:"avatar_url,omitempty"`
}

// RefreshRequest is the request body for token refresh
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// ForgotPasswordRequest is the request body for forgot password
type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// ResetPasswordRequest is the request body for password reset
type ResetPasswordRequest struct {
	Token    string `json:"token" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
}
