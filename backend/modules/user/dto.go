package user

type CreateUserRequest struct {
	Name     string  `json:"name" validate:"required"`
	Email    string  `json:"email" validate:"required,email"`
	Password string  `json:"password" validate:"required,min=8"`
	Role     string  `json:"role" validate:"required,oneof=admin lecturer student staff"`
	NIM      *string `json:"nim,omitempty"`
	NIP      *string `json:"nip,omitempty"`
}

type UpdateProfileRequest struct {
	Name      string  `json:"name" validate:"required"`
	AvatarURL *string `json:"avatar_url,omitempty"`
}
