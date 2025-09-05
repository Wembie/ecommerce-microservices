package models

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAuthRequestJSON(t *testing.T) {
	req := AuthRequest{
		Username: "juan",
		Password: "secret",
	}

	data, err := json.Marshal(req)
	assert.NoError(t, err)

	var decoded AuthRequest
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)

	assert.Equal(t, req.Username, decoded.Username)
	assert.Equal(t, req.Password, decoded.Password)
}

func TestAuthResponseSuccess(t *testing.T) {
	token := "fake-jwt-token"
	resp := AuthResponse{
		Success: true,
		Token:   &token,
	}

	data, err := json.Marshal(resp)
	assert.NoError(t, err)


	jsonStr := string(data)
	assert.Contains(t, jsonStr, "fake-jwt-token")
	assert.NotContains(t, jsonStr, "error_message")

	var decoded AuthResponse
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)

	assert.True(t, decoded.Success)
	assert.Equal(t, token, *decoded.Token)
	assert.Nil(t, decoded.ErrorMessage)
}

func TestAuthResponseError(t *testing.T) {
	errMsg := "invalid credentials"
	resp := AuthResponse{
		Success:      false,
		ErrorMessage: &errMsg,
	}

	data, err := json.Marshal(resp)
	assert.NoError(t, err)

	// Debe contener el error pero no token
	jsonStr := string(data)
	assert.Contains(t, jsonStr, "invalid credentials")
	assert.NotContains(t, jsonStr, "token")

	var decoded AuthResponse
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)

	assert.False(t, decoded.Success)
	assert.Nil(t, decoded.Token)
	assert.Equal(t, errMsg, *decoded.ErrorMessage)
}

func TestValidateUserResponse(t *testing.T) {
	uid := uuid.New()
	username := "juan"
	email := "juan@example.com"

	resp := ValidateUserResponse{
		Valid:    true,
		UserID:   &uid,
		Username: &username,
		Email:    &email,
	}

	data, err := json.Marshal(resp)
	assert.NoError(t, err)

	var decoded ValidateUserResponse
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)

	assert.True(t, decoded.Valid)
	assert.Equal(t, uid, *decoded.UserID)
	assert.Equal(t, username, *decoded.Username)
	assert.Equal(t, email, *decoded.Email)
}
