package models

import "github.com/google/uuid"

type ProductInfo struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
}

type UserInfo struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
}

type TokenRequest struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Success      bool    `json:"success"`
	Token        *string `json:"token,omitempty"`
	ErrorMessage *string `json:"error_message,omitempty"`
}