package models_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"ecommerce.users.manager/internal/models"
)

func TestUser_JSONSerialization(t *testing.T) {
	now := time.Now()
	id := uuid.New()

	user := models.User{
		ID:           id,
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashedpwd",
		CreatedAt:    now,
		UpdatedAt:    &now,
	}

	data, err := json.Marshal(user)
	assert.NoError(t, err)

	var decoded models.User
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)

	assert.Equal(t, user.ID, decoded.ID)
	assert.Equal(t, user.Username, decoded.Username)
	assert.Equal(t, user.Email, decoded.Email)
	assert.Equal(t, user.PasswordHash, decoded.PasswordHash)
}

func TestCreateUserRequest_JSON(t *testing.T) {
	req := models.CreateUserRequest{
		Username: "newuser",
		Email:    "new@example.com",
		Password: "plainpwd",
	}

	data, err := json.Marshal(req)
	assert.NoError(t, err)

	var decoded models.CreateUserRequest
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)

	assert.Equal(t, req.Username, decoded.Username)
	assert.Equal(t, req.Email, decoded.Email)
	assert.Equal(t, req.Password, decoded.Password)
}

func TestUpdateUserRequest_OptionalFields(t *testing.T) {
	id := uuid.New()
	username := "updated"
	email := "updated@example.com"
	password := "newpwd"

	req := models.UpdateUserRequest{
		ID:        id,
		Username:  &username,
		Email:     &email,
		Password:  &password,
		UpdatedAt: time.Now(),
	}

	assert.Equal(t, id, req.ID)
	assert.NotNil(t, req.Username)
	assert.Equal(t, "updated", *req.Username)
}

func TestAuthResponse_JSONWithOptional(t *testing.T) {
	token := "jwt-token"
	resp := models.AuthResponse{
		Success: true,
		Token:   &token,
	}

	data, err := json.Marshal(resp)
	assert.NoError(t, err)

	jsonStr := string(data)
	assert.Contains(t, jsonStr, "jwt-token")
	assert.Contains(t, jsonStr, "success")
}

func TestValidateUserResponse_JSON(t *testing.T) {
	userID := uuid.New()
	username := "testuser"
	email := "test@example.com"

	resp := models.ValidateUserResponse{
		Valid:    true,
		UserID:   &userID,
		Username: &username,
		Email:    &email,
	}

	data, err := json.Marshal(resp)
	assert.NoError(t, err)

	var decoded models.ValidateUserResponse
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)

	assert.True(t, decoded.Valid)
	assert.Equal(t, userID, *decoded.UserID)
	assert.Equal(t, username, *decoded.Username)
	assert.Equal(t, email, *decoded.Email)
}
