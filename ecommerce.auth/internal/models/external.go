package models

import "github.com/google/uuid"

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Success      bool    `json:"success"`
	Token        *string `json:"token,omitempty"`
	ErrorMessage *string `json:"error_message,omitempty"`
}

type ValidateUserResponse struct {
	Valid    bool      	`json:"valid"`
	UserID   *uuid.UUID `json:"user_id,omitempty"`
	Username *string    `json:"username,omitempty"`
	Email    *string    `json:"email,omitempty"`
}