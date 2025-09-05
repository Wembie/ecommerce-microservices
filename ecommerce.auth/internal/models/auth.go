package models

type TokenRequest struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

type TokenValidationRequest struct {
	Token string `json:"token" binding:"required"`
}

type TokenValidationResponse struct {
	Valid  bool   `json:"valid"`
	UserID string `json:"user_id,omitempty"`
	Email  string `json:"email,omitempty"`
}
