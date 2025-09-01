package utils_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	pb "github.com/Wembie/ecommerce.users.manager.lib.protos/v1/libgo"
	"ecommerce.users.manager/internal/models"
	"ecommerce.users.manager/internal/utils"
)

func TestConvertToUserProto(t *testing.T) {
	now := time.Now()
	updated := now.Add(time.Hour)

	validID := uuid.New()

	tests := []struct {
		name        string
		input       *models.User
		expected    *time.Time
		expectNil   bool
		expectError bool
	}{
		{
			name:        "nil input should return error",
			input:       nil,
			expectError: true,
		},
		{
			name: "valid model without UpdatedAt",
			input: &models.User{
				ID:           validID,
				Username:     "Wembie",
				Email:        "wembie@example.com",
				PasswordHash: "hashedpass",
				CreatedAt:    now,
				UpdatedAt:    nil,
			},
			expectNil: true,
		},
		{
			name: "valid model with UpdatedAt",
			input: &models.User{
				ID:           validID,
				Username:     "Juan Acosta",
				Email:        "juanes@example.com",
				PasswordHash: "hashedpass2",
				CreatedAt:    now,
				UpdatedAt:    &updated,
			},
			expected:  &updated,
			expectNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := utils.ConvertToUserProto(tt.input)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, resp)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.input.ID.String(), resp.Id)
			assert.Equal(t, tt.input.Username, resp.Username)
			assert.Equal(t, tt.input.Email, resp.Email)
			assert.Equal(t, tt.input.PasswordHash, resp.PasswordHash)
			assert.WithinDuration(t, tt.input.CreatedAt, resp.CreatedAt.AsTime(), time.Second)

			if tt.expectNil {
				assert.Nil(t, resp.UpdatedAt)
			} else {
				assert.WithinDuration(t, tt.expected.UTC(), resp.UpdatedAt.AsTime().UTC(), time.Second)
			}
		})
	}
}

func TestConvertToCreateUserModel(t *testing.T) {
	tests := []struct {
		name        string
		input       *pb.CreateUserRequest
		expectError bool
	}{
		{"nil input", nil, true},
		{"missing username", &pb.CreateUserRequest{Email: "a@b.com", Password: "123"}, true},
		{"missing email", &pb.CreateUserRequest{Username: "juan", Password: "123"}, true},
		{"invalid email", &pb.CreateUserRequest{Username: "juan", Email: "not-an-email", Password: "123"}, true},
		{"missing password", &pb.CreateUserRequest{Username: "juan", Email: "a@b.com"}, true},
		{"valid", &pb.CreateUserRequest{Username: "juan", Email: "a@b.com", Password: "123"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := utils.ConvertToCreateUserModel(tt.input)
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.input.Username, resp.Username)
				assert.Equal(t, tt.input.Email, resp.Email)
				assert.Equal(t, tt.input.Password, resp.Password)
			}
		})
	}
}

func TestConvertToGetUserModel(t *testing.T) {
	validID := uuid.New()
	tests := []struct {
		name        string
		input       *pb.GetUserRequest
		expectError bool
	}{
		{"nil input", nil, true},
		{"invalid id", &pb.GetUserRequest{Id: "bad-id"}, true},
		{"valid id", &pb.GetUserRequest{Id: validID.String()}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := utils.ConvertToGetUserModel(tt.input)
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, validID, resp.ID)
			}
		})
	}
}

func TestConvertToUpdateUserModel(t *testing.T) {
	validID := uuid.New()
	username := "juan"
	email := "juan@example.com"
	password := "123"

	tests := []struct {
		name        string
		input       *pb.UpdateUserRequest
		expectError bool
	}{
		{"nil input", nil, true},
		{"invalid id", &pb.UpdateUserRequest{Id: "bad-id"}, true},
		{"empty username", &pb.UpdateUserRequest{Id: validID.String(), Username: &[]string{""}[0]}, true},
		{"invalid email", &pb.UpdateUserRequest{Id: validID.String(), Email: &[]string{"bad"}[0]}, true},
		{"empty password", &pb.UpdateUserRequest{Id: validID.String(), Password: &[]string{""}[0]}, true},
		{"valid update", &pb.UpdateUserRequest{Id: validID.String(), Username: &username, Email: &email, Password: &password}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := utils.ConvertToUpdateUserModel(tt.input)
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, validID, resp.ID)
				assert.Equal(t, username, *resp.Username)
				assert.Equal(t, email, *resp.Email)
				assert.Equal(t, password, *resp.Password)
			}
		})
	}
}

func TestConvertToDeleteUserModel(t *testing.T) {
	validID := uuid.New()
	tests := []struct {
		name        string
		input       *pb.DeleteUserRequest
		expectError bool
	}{
		{"nil input", nil, true},
		{"invalid id", &pb.DeleteUserRequest{Id: "bad-id"}, true},
		{"valid id", &pb.DeleteUserRequest{Id: validID.String()}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := utils.ConvertToDeleteUserModel(tt.input)
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, validID, resp.ID)
			}
		})
	}
}

func TestConvertToDeleteUserProto(t *testing.T) {
	tests := []struct {
		name        string
		input       *models.DeleteUserResponse
		expectError bool
	}{
		{"nil input", nil, true},
		{"valid success true", &models.DeleteUserResponse{Success: true}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := utils.ConvertToDeleteUserProto(tt.input)
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.input.Success, resp.Success)
			}
		})
	}
}

func TestConvertToAuthUserModel(t *testing.T) {
	tests := []struct {
		name        string
		input       *pb.AuthRequest
		expectError bool
	}{
		{"nil input", nil, true},
		{"missing username", &pb.AuthRequest{Password: "123"}, true},
		{"missing password", &pb.AuthRequest{Username: "juan"}, true},
		{"valid", &pb.AuthRequest{Username: "juan", Password: "123"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := utils.ConvertToAuthUserModel(tt.input)
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.input.Username, resp.Username)
				assert.Equal(t, tt.input.Password, resp.Password)
			}
		})
	}
}

func TestConvertToAuthUserProto(t *testing.T) {
	token := "abc"
	errMsg := "wrong pass"

	tests := []struct {
		name        string
		input       *models.AuthResponse
		expectError bool
	}{
		{"nil input", nil, true},
		{"valid success", &models.AuthResponse{Success: true, Token: &token}, false},
		{"valid failure", &models.AuthResponse{Success: false, ErrorMessage: &errMsg}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := utils.ConvertToAuthUserProto(tt.input)
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.input.Success, resp.Success)
				assert.Equal(t, tt.input.Token, resp.Token)
				assert.Equal(t, tt.input.ErrorMessage, resp.ErrorMessage)
			}
		})
	}
}

func TestConvertToValidateUserModel(t *testing.T) {
	tests := []struct {
		name        string
		input       *pb.ValidateUserRequest
		expectError bool
	}{
		{"nil input", nil, true},
		{"missing token", &pb.ValidateUserRequest{}, true},
		{"valid", &pb.ValidateUserRequest{Token: "valid-token"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := utils.ConvertToValidateUserModel(tt.input)
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.input.Token, resp.Token)
			}
		})
	}
}

func TestConvertToValidateUserProto(t *testing.T) {
	validID := uuid.New()
	username := "juan"
	email := "juan@example.com"

	tests := []struct {
		name        string
		input       *models.ValidateUserResponse
		expectError bool
	}{
		{"nil input", nil, true},
		{"valid with userID", &models.ValidateUserResponse{Valid: true, UserID: &validID, Username: &username, Email: &email}, false},
		{"valid without userID", &models.ValidateUserResponse{Valid: false}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := utils.ConvertToValidateUserProto(tt.input)
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.input.Valid, resp.Valid)
				if tt.input.UserID != nil {
					assert.Equal(t, tt.input.UserID.String(), *resp.UserId)
				} else {
					assert.Nil(t, resp.UserId)
				}
			}
		})
	}
}
