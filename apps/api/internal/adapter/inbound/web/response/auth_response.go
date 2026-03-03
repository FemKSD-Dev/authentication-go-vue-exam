package response

type RegisterResponse struct {
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type MeResponse struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
}
