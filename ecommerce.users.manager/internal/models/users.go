package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:user_service.users"`

	ID           uuid.UUID  `bun:"id,pk,type:uuid,default:gen_random_uuid()" json:"id"`
	Username     string     `bun:",unique,notnull" json:"username"`
	Email        string     `bun:",unique,notnull" json:"email"`
	PasswordHash string     `bun:"password_hash,notnull" json:"password_hash"`
	CreatedAt    time.Time  `bun:",default:now()" json:"created_at"`
	UpdatedAt    *time.Time `bun:"updated_at,nullzero" json:"updated_at,omitempty"`
}

type CreateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
 
type GetUserRequest struct {
	ID uuid.UUID `json:"id"`
}

type UpdateUserRequest struct {
	bun.BaseModel `bun:"table:user_service.users"`

	ID           uuid.UUID  `bun:"id,pk,type:uuid" json:"id"`
	Username     *string    `json:"username,omitempty"`
	Email        *string    `json:"email,omitempty"`
	Password  	 *string    `bun:"password_hash" json:"password,omitempty"`
	UpdatedAt    time.Time  `bun:"updated_at,type:timestamptz" json:"-"`
}

type DeleteUserRequest struct {
	ID uuid.UUID `json:"id"`
}

type DeleteUserResponse struct {
	Success bool `json:"success"`
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

type ValidateUserRequest struct {
	Token string `json:"token"`
}

type ValidateUserResponse struct {
	Valid    bool       `json:"valid"`
	UserID   *uuid.UUID `json:"user_id,omitempty"`
	Username *string    `json:"username,omitempty"`
	Email    *string    `json:"email,omitempty"`
}
