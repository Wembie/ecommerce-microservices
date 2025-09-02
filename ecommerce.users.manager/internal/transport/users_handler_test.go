package transport_test

import (
	"errors"
	"testing"
	"time"

	"ecommerce.users.manager/internal/mocks"
	"ecommerce.users.manager/internal/models"
	"ecommerce.users.manager/internal/transport"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	pb "github.com/Wembie/ecommerce.users.manager.lib.protos/v1/libgo"
)

func TestPingHandler_CreateUser(t *testing.T) {
	setup := mocks.NewTestSetup()
	
	tests := []struct {
		name        string
		request     *pb.CreateUserRequest
		setupMocks  func(*mocks.MockService)
		expectError bool
	}{
		{
			name: "successful create",
			request: &pb.CreateUserRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMocks: func(svc *mocks.MockService) {
				user := &models.User{
					ID:           uuid.New(),
					Username:     "testuser",
					Email:        "test@example.com",
					PasswordHash: "hashedpassword",
					CreatedAt:    time.Now(),
				}
				svc.On("CreateUser", setup.Ctx, mock.Anything, mock.Anything).Return(user, nil)
			},
			expectError: false,
		},
		{
			name: "conversion error - empty username",
			request: &pb.CreateUserRequest{
				Username: "",
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMocks: func(svc *mocks.MockService) {},
			expectError: true,
		},
		{
			name: "conversion error - empty email",
			request: &pb.CreateUserRequest{
				Username: "testuser",
				Email:    "",
				Password: "password123",
			},
			setupMocks: func(svc *mocks.MockService) {},
			expectError: true,
		},
		{
			name: "conversion error - empty password",
			request: &pb.CreateUserRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "",
			},
			setupMocks: func(svc *mocks.MockService) {},
			expectError: true,
		},
		{
			name: "service error",
			request: &pb.CreateUserRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMocks: func(svc *mocks.MockService) {
				svc.On("CreateUser", setup.Ctx, mock.Anything, mock.Anything).Return(nil, errors.New("service error"))
			},
			expectError: true,
		},
		{
			name:        "nil request should return error",
			request:     nil,
			setupMocks:  func(svc *mocks.MockService) {},
			expectError: true,
		},
		{
			name: "response conversion success",
			request: &pb.CreateUserRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMocks: func(svc *mocks.MockService) {
				user := &models.User{
					ID:           uuid.New(),
					Username:     "testuser",
					Email:        "test@example.com",
					PasswordHash: "hashedpassword",
					CreatedAt:    time.Now(),
				}
				svc.On("CreateUser", setup.Ctx, mock.Anything, mock.Anything).Return(user, nil)
			},
			expectError: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := setup.NewMockService()
			h := transport.NewHandler(setup.Logger, mockSvc, setup.Config)
			tt.setupMocks(mockSvc)
			
			resp, err := h.CreateUser(setup.Ctx, tt.request)
			
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
			}
			
			mockSvc.AssertExpectations(t)
		})
	}
}

func TestPingHandler_GetUser(t *testing.T) {
	setup := mocks.NewTestSetup()
	validID := uuid.New()
	
	tests := []struct {
		name        string
		id          string
		setupMocks  func(*mocks.MockService, uuid.UUID)
		expectError bool
	}{
		{
			name: "successful get",
			id:   validID.String(),
			setupMocks: func(svc *mocks.MockService, id uuid.UUID) {
				user := &models.User{
					ID:           id,
					Username:     "testuser",
					Email:        "test@example.com",
					PasswordHash: "hashedpassword",
					CreatedAt:    time.Now(),
				}
				svc.On("GetUser", setup.Ctx, mock.Anything, mock.Anything).Return(user, nil)
			},
			expectError: false,
		},
		{
			name:        "invalid uuid",
			id:          "invalid-uuid",
			setupMocks:  func(svc *mocks.MockService, id uuid.UUID) {},
			expectError: true,
		},
		{
			name: "service error",
			id:   validID.String(),
			setupMocks: func(svc *mocks.MockService, id uuid.UUID) {
				svc.On("GetUser", setup.Ctx, mock.Anything, mock.Anything).Return(nil, errors.New("service error"))
			},
			expectError: true,
		},
		{
			name: "response conversion success",
			id:   validID.String(),
			setupMocks: func(svc *mocks.MockService, id uuid.UUID) {
				user := &models.User{
					ID:           id,
					Username:     "testuser",
					Email:        "test@example.com",
					PasswordHash: "hashedpassword",
					CreatedAt:    time.Now(),
				}
				svc.On("GetUser", setup.Ctx, mock.Anything, mock.Anything).Return(user, nil)
			},
			expectError: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := setup.NewMockService()
			h := transport.NewHandler(setup.Logger, mockSvc, setup.Config)
			id, _ := uuid.Parse(tt.id)
			req := &pb.GetUserRequest{Id: tt.id}
			tt.setupMocks(mockSvc, id)
			
			resp, err := h.GetUser(setup.Ctx, req)
			
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
			}
			
			mockSvc.AssertExpectations(t)
		})
	}
}

func TestPingHandler_UpdateUser(t *testing.T) {
	setup := mocks.NewTestSetup()
	validID := uuid.New()
	newUsername := "newusername"
	newEmail := "newemail@example.com"
	newPassword := "newpassword123"
	
	tests := []struct {
		name        string
		request     *pb.UpdateUserRequest
		setupMocks  func(*mocks.MockService)
		expectError bool
	}{
		{
			name: "successful update",
			request: &pb.UpdateUserRequest{
				Id: validID.String(),
			},
			setupMocks: func(svc *mocks.MockService) {
				user := &models.User{
					ID:           validID,
					Username:     "testuser",
					Email:        "test@example.com",
					PasswordHash: "hashedpassword",
					CreatedAt:    time.Now(),
				}
				svc.On("UpdateUser", setup.Ctx, mock.Anything, mock.Anything).Return(user, nil)
			},
			expectError: false,
		},
		{
			name: "invalid uuid",
			request: &pb.UpdateUserRequest{
				Id: "invalid-uuid",
			},
			setupMocks:  func(svc *mocks.MockService) {},
			expectError: true,
		},
		{
			name: "service error",
			request: &pb.UpdateUserRequest{
				Id: validID.String(),
			},
			setupMocks: func(svc *mocks.MockService) {
				svc.On("UpdateUser", setup.Ctx, mock.Anything, mock.Anything).Return(nil, errors.New("service error"))
			},
			expectError: true,
		},
		{
			name: "conversion error for request",
			request: &pb.UpdateUserRequest{
				Id: "invalid-uuid",
			},
			setupMocks:  func(svc *mocks.MockService) {},
			expectError: true,
		},
		{
			name: "successful update with optional fields",
			request: &pb.UpdateUserRequest{
				Id:       validID.String(),
				Username: &newUsername,
				Email:    &newEmail,
				Password: &newPassword,
			},
			setupMocks: func(svc *mocks.MockService) {
				user := &models.User{
					ID:           validID,
					Username:     newUsername,
					Email:        newEmail,
					PasswordHash: "hashedpassword",
					CreatedAt:    time.Now(),
				}
				svc.On("UpdateUser", setup.Ctx, mock.Anything, mock.Anything).Return(user, nil)
			},
			expectError: false,
		},
		{
			name: "service returns nil (no rows affected)",
			request: &pb.UpdateUserRequest{
				Id: validID.String(),
			},
			setupMocks: func(svc *mocks.MockService) {
				svc.On("UpdateUser", setup.Ctx, mock.Anything, mock.Anything).Return(nil, nil)
			},
			expectError: true, // Should return error when service returns nil
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := setup.NewMockService()
			h := transport.NewHandler(setup.Logger, mockSvc, setup.Config)
			tt.setupMocks(mockSvc)
			
			resp, err := h.UpdateUser(setup.Ctx, tt.request)
			
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
			}
			
			mockSvc.AssertExpectations(t)
		})
	}
}

func TestPingHandler_DeleteUser(t *testing.T) {
	setup := mocks.NewTestSetup()
	validID := uuid.New()
	
	tests := []struct {
		name        string
		request     *pb.DeleteUserRequest
		setupMocks  func(*mocks.MockService)
		expectError bool
	}{
		{
			name: "successful delete",
			request: &pb.DeleteUserRequest{
				Id: validID.String(),
			},
			setupMocks: func(svc *mocks.MockService) {
				svc.On("DeleteUser", setup.Ctx, mock.Anything, mock.Anything).Return(true, nil)
			},
			expectError: false,
		},
		{
			name: "invalid uuid",
			request: &pb.DeleteUserRequest{
				Id: "invalid-uuid",
			},
			setupMocks:  func(svc *mocks.MockService) {},
			expectError: true,
		},
		{
			name: "service error",
			request: &pb.DeleteUserRequest{
				Id: validID.String(),
			},
			setupMocks: func(svc *mocks.MockService) {
				svc.On("DeleteUser", setup.Ctx, mock.Anything, mock.Anything).Return(false, errors.New("service error"))
			},
			expectError: true,
		},
		{
			name: "delete failed (returns false)",
			request: &pb.DeleteUserRequest{
				Id: validID.String(),
			},
			setupMocks: func(svc *mocks.MockService) {
				svc.On("DeleteUser", setup.Ctx, mock.Anything, mock.Anything).Return(false, nil)
			},
			expectError: false, // Service handles this case, not an error
		},
		{
			name:        "nil request should return error",
			request:     nil,
			setupMocks:  func(svc *mocks.MockService) {},
			expectError: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := setup.NewMockService()
			h := transport.NewHandler(setup.Logger, mockSvc, setup.Config)
			tt.setupMocks(mockSvc)
			
			resp, err := h.DeleteUser(setup.Ctx, tt.request)
			
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
			}
			
			mockSvc.AssertExpectations(t)
		})
	}
}

func TestPingHandler_AuthenticateUser(t *testing.T) {
	setup := mocks.NewTestSetup()
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.test"
	
	tests := []struct {
		name        string
		request     *pb.AuthRequest
		setupMocks  func(*mocks.MockService)
		expectError bool
	}{
		{
			name: "successful authentication",
			request: &pb.AuthRequest{
				Username: "testuser",
				Password: "password123",
			},
			setupMocks: func(svc *mocks.MockService) {
				authResp := &models.AuthResponse{
					Success: true,
					Token:   &token,
				}
				svc.On("AuthenticateUser", setup.Ctx, mock.Anything, mock.Anything).Return(authResp, nil)
			},
			expectError: false,
		},
		{
			name: "failed authentication",
			request: &pb.AuthRequest{
				Username: "testuser",
				Password: "wrongpassword",
			},
			setupMocks: func(svc *mocks.MockService) {
				errorMsg := "invalid credentials"
				authResp := &models.AuthResponse{
					Success:      false,
					ErrorMessage: &errorMsg,
				}
				svc.On("AuthenticateUser", setup.Ctx, mock.Anything, mock.Anything).Return(authResp, nil)
			},
			expectError: false, // Not an error, just failed auth
		},
		{
			name: "conversion error - empty username",
			request: &pb.AuthRequest{
				Username: "",
				Password: "password123",
			},
			setupMocks:  func(svc *mocks.MockService) {},
			expectError: true,
		},
		{
			name: "conversion error - empty password",
			request: &pb.AuthRequest{
				Username: "testuser",
				Password: "",
			},
			setupMocks:  func(svc *mocks.MockService) {},
			expectError: true,
		},
		{
			name: "service error",
			request: &pb.AuthRequest{
				Username: "testuser",
				Password: "password123",
			},
			setupMocks: func(svc *mocks.MockService) {
				svc.On("AuthenticateUser", setup.Ctx, mock.Anything, mock.Anything).Return(nil, errors.New("service error"))
			},
			expectError: true,
		},
		{
			name:        "nil request should return error",
			request:     nil,
			setupMocks:  func(svc *mocks.MockService) {},
			expectError: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := setup.NewMockService()
			h := transport.NewHandler(setup.Logger, mockSvc, setup.Config)
			tt.setupMocks(mockSvc)
			
			resp, err := h.AuthenticateUser(setup.Ctx, tt.request)
			
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
			}
			
			mockSvc.AssertExpectations(t)
		})
	}
}

func TestPingHandler_ValidateUser(t *testing.T) {
	setup := mocks.NewTestSetup()
	userID := uuid.New()
	username := "testuser"
	email := "test@example.com"
	
	tests := []struct {
		name        string
		request     *pb.ValidateUserRequest
		setupMocks  func(*mocks.MockService)
		expectError bool
	}{
		{
			name: "successful validation",
			request: &pb.ValidateUserRequest{
				Token: "valid-jwt-token",
			},
			setupMocks: func(svc *mocks.MockService) {
				validateResp := &models.ValidateUserResponse{
					Valid:    true,
					UserID:   &userID,
					Username: &username,
					Email:    &email,
				}
				svc.On("ValidateUser", setup.Ctx, mock.Anything, mock.Anything).Return(validateResp, nil)
			},
			expectError: false,
		},
		{
			name: "invalid token",
			request: &pb.ValidateUserRequest{
				Token: "invalid-token",
			},
			setupMocks: func(svc *mocks.MockService) {
				validateResp := &models.ValidateUserResponse{
					Valid: false,
				}
				svc.On("ValidateUser", setup.Ctx, mock.Anything, mock.Anything).Return(validateResp, nil)
			},
			expectError: false, // Not an error, just invalid token
		},
		{
			name: "conversion error - empty token",
			request: &pb.ValidateUserRequest{
				Token: "",
			},
			setupMocks:  func(svc *mocks.MockService) {},
			expectError: true,
		},
		{
			name: "service error",
			request: &pb.ValidateUserRequest{
				Token: "valid-jwt-token",
			},
			setupMocks: func(svc *mocks.MockService) {
				svc.On("ValidateUser", setup.Ctx, mock.Anything, mock.Anything).Return(nil, errors.New("service error"))
			},
			expectError: true,
		},
		{
			name:        "nil request should return error",
			request:     nil,
			setupMocks:  func(svc *mocks.MockService) {},
			expectError: true,
		},
		{
			name: "response conversion success",
			request: &pb.ValidateUserRequest{
				Token: "valid-jwt-token",
			},
			setupMocks: func(svc *mocks.MockService) {
				validateResp := &models.ValidateUserResponse{
					Valid:    true,
					UserID:   &userID,
					Username: &username,
					Email:    &email,
				}
				svc.On("ValidateUser", setup.Ctx, mock.Anything, mock.Anything).Return(validateResp, nil)
			},
			expectError: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := setup.NewMockService()
			h := transport.NewHandler(setup.Logger, mockSvc, setup.Config)
			tt.setupMocks(mockSvc)
			
			resp, err := h.ValidateUser(setup.Ctx, tt.request)
			
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
			}
			
			mockSvc.AssertExpectations(t)
		})
	}
}