package request

type RegisterRequest struct {
	Username        string `json:"username" validate:"required,min=3,max=255,alphanum"`
	Password        string `json:"password" validate:"required,min=8,max=255,ascii"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required,min=3,max=255,alphanum"`
	Password string `json:"password" validate:"required,min=8,max=255,ascii"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
