package service_test

import (
	"errors"
	"testing"
	"time"

	"ecommerce.users.manager/internal/mocks"
	"ecommerce.users.manager/internal/models"
	"ecommerce.users.manager/internal/service"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestUsersService_CreateUser(t *testing.T) {
	setup := mocks.NewTestSetup()
	
	tests := []struct {
		name        string
		request     *models.CreateUserRequest
		setupMocks  func(*mocks.MockRepository)
		expectError bool
	}{
		{
			name: "successful user creation",
			request: &models.CreateUserRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMocks: func(repo *mocks.MockRepository) {
				repo.On("CreateGeneric", setup.Ctx, mock.AnythingOfType("*zap.Logger"), mock.AnythingOfType("*models.User")).Return(nil).Run(func(args mock.Arguments) {
					user := args.Get(2).(*models.User)
					user.ID = uuid.New()
				})
			},
			expectError: false,
		},
		{
			name: "repository error should return error",
			request: &models.CreateUserRequest{
				Username: "testuser",
				Email:    "test@example.com", 
				Password: "password123",
			},
			setupMocks: func(repo *mocks.MockRepository) {
				repo.On("CreateGeneric", setup.Ctx, mock.AnythingOfType("*zap.Logger"), mock.AnythingOfType("*models.User")).Return(errors.New("creation failed"))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := setup.NewMockRepository()
			tt.setupMocks(mockRepo)

			svc := service.NewService(setup.Logger, mockRepo, setup.Config)

			result, err := svc.CreateUser(setup.Ctx, setup.Logger, tt.request)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.request.Username, result.Username)
				assert.Equal(t, tt.request.Email, result.Email)
				assert.NotEmpty(t, result.PasswordHash)
				assert.NotEqual(t, tt.request.Password, result.PasswordHash) // Password should be hashed
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUsersService_GetUser(t *testing.T) {
	setup := mocks.NewTestSetup()
	testID := uuid.New()

	tests := []struct {
		name        string
		request     *models.GetUserRequest
		setupMocks  func(*mocks.MockRepository)
		expectError bool
	}{
		{
			name: "successful get user",
			request: &models.GetUserRequest{
				ID: testID,
			},
			setupMocks: func(repo *mocks.MockRepository) {
				repo.On("GetGeneric", setup.Ctx, mock.AnythingOfType("*zap.Logger"), mock.AnythingOfType("*models.User")).Return(nil).Run(func(args mock.Arguments) {
					user := args.Get(2).(*models.User)
					user.ID = testID
					user.Username = "testuser"
					user.Email = "test@example.com"
					user.PasswordHash = "hashedpassword"
					user.CreatedAt = time.Now()
				})
			},
			expectError: false,
		},
		{
			name: "user not found should return error",
			request: &models.GetUserRequest{
				ID: testID,
			},
			setupMocks: func(repo *mocks.MockRepository) {
				repo.On("GetGeneric", setup.Ctx, mock.AnythingOfType("*zap.Logger"), mock.AnythingOfType("*models.User")).Return(errors.New("not found"))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := setup.NewMockRepository()
			tt.setupMocks(mockRepo)

			svc := service.NewService(setup.Logger, mockRepo, setup.Config)

			result, err := svc.GetUser(setup.Ctx, setup.Logger, tt.request)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, testID, result.ID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUsersService_UpdateUser(t *testing.T) {
	setup := mocks.NewTestSetup()
	testID := uuid.New()
	newUsername := "newusername"
	newEmail := "newemail@example.com"
	newPassword := "newpassword123"

	tests := []struct {
		name        string
		request     *models.UpdateUserRequest
		setupMocks  func(*mocks.MockRepository)
		expectError bool
		expectNil   bool
	}{
		{
			name: "successful update with password",
			request: &models.UpdateUserRequest{
				ID:       testID,
				Username: &newUsername,
				Email:    &newEmail,
				Password: &newPassword,
			},
			setupMocks: func(repo *mocks.MockRepository) {
				repo.On("UpdateGeneric", setup.Ctx, mock.AnythingOfType("*zap.Logger"), mock.AnythingOfType("*models.UpdateUserRequest")).Return(int64(1), nil)
				repo.On("GetGeneric", setup.Ctx, mock.AnythingOfType("*zap.Logger"), mock.AnythingOfType("*models.User")).Return(nil).Run(func(args mock.Arguments) {
					user := args.Get(2).(*models.User)
					user.ID = testID
					user.Username = newUsername
					user.Email = newEmail
					user.PasswordHash = "hashedpassword"
					user.CreatedAt = time.Now()
				})
			},
			expectError: false,
			expectNil:   false,
		},
		{
			name: "successful update without password",
			request: &models.UpdateUserRequest{
				ID:       testID,
				Username: &newUsername,
				Email:    &newEmail,
			},
			setupMocks: func(repo *mocks.MockRepository) {
				repo.On("UpdateGeneric", setup.Ctx, mock.AnythingOfType("*zap.Logger"), mock.AnythingOfType("*models.UpdateUserRequest")).Return(int64(1), nil)
				repo.On("GetGeneric", setup.Ctx, mock.AnythingOfType("*zap.Logger"), mock.AnythingOfType("*models.User")).Return(nil).Run(func(args mock.Arguments) {
					user := args.Get(2).(*models.User)
					user.ID = testID
					user.Username = newUsername
					user.Email = newEmail
					user.PasswordHash = "hashedpassword"
					user.CreatedAt = time.Now()
				})
			},
			expectError: false,
			expectNil:   false,
		},
		{
			name: "no rows affected should return nil",
			request: &models.UpdateUserRequest{
				ID:       testID,
				Username: &newUsername,
			},
			setupMocks: func(repo *mocks.MockRepository) {
				repo.On("UpdateGeneric", setup.Ctx, mock.AnythingOfType("*zap.Logger"), mock.AnythingOfType("*models.UpdateUserRequest")).Return(int64(0), nil)
			},
			expectError: false,
			expectNil:   true,
		},
		{
			name: "update error should return error",
			request: &models.UpdateUserRequest{
				ID:       testID,
				Username: &newUsername,
			},
			setupMocks: func(repo *mocks.MockRepository) {
				repo.On("UpdateGeneric", setup.Ctx, mock.AnythingOfType("*zap.Logger"), mock.AnythingOfType("*models.UpdateUserRequest")).Return(int64(0), errors.New("update failed"))
			},
			expectError: true,
			expectNil:   false,
		},
		{
			name: "get after update error should return error",
			request: &models.UpdateUserRequest{
				ID:       testID,
				Username: &newUsername,
			},
			setupMocks: func(repo *mocks.MockRepository) {
				repo.On("UpdateGeneric", setup.Ctx, mock.AnythingOfType("*zap.Logger"), mock.AnythingOfType("*models.UpdateUserRequest")).Return(int64(1), nil)
				repo.On("GetGeneric", setup.Ctx, mock.AnythingOfType("*zap.Logger"), mock.AnythingOfType("*models.User")).Return(errors.New("get failed"))
			},
			expectError: true,
			expectNil:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := setup.NewMockRepository()
			tt.setupMocks(mockRepo)

			svc := service.NewService(setup.Logger, mockRepo, setup.Config)

			result, err := svc.UpdateUser(setup.Ctx, setup.Logger, tt.request)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.expectNil {
					assert.Nil(t, result)
				} else {
					assert.NotNil(t, result)
					assert.Equal(t, testID, result.ID)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUsersService_DeleteUser(t *testing.T) {
	setup := mocks.NewTestSetup()
	testID := uuid.New()

	tests := []struct {
		name        string
		request     *models.DeleteUserRequest
		setupMocks  func(*mocks.MockRepository)
		expectError bool
	}{
		{
			name: "successful delete",
			request: &models.DeleteUserRequest{
				ID: testID,
			},
			setupMocks: func(repo *mocks.MockRepository) {
				repo.On("DeleteGeneric", setup.Ctx, mock.AnythingOfType("*zap.Logger"), mock.AnythingOfType("*models.User")).Return(nil)
			},
			expectError: false,
		},
		{
			name: "repository error should return error",
			request: &models.DeleteUserRequest{
				ID: testID,
			},
			setupMocks: func(repo *mocks.MockRepository) {
				repo.On("DeleteGeneric", setup.Ctx, mock.AnythingOfType("*zap.Logger"), mock.AnythingOfType("*models.User")).Return(errors.New("delete failed"))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := setup.NewMockRepository()
			tt.setupMocks(mockRepo)

			svc := service.NewService(setup.Logger, mockRepo, setup.Config)

			result, err := svc.DeleteUser(setup.Ctx, setup.Logger, tt.request)

			if tt.expectError {
				assert.Error(t, err)
				assert.False(t, result)
			} else {
				assert.NoError(t, err)
				assert.True(t, result)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUsersService_AuthenticateUser(t *testing.T) {
	setup := mocks.NewTestSetup()
	testID := uuid.New()
	username := "testuser"
	email := "test@example.com"
	password := "password123"
	
	// Generate a hash for testing
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	tests := []struct {
		name           string
		request        *models.AuthRequest
		setupMocks     func(*mocks.MockRepository)
		expectError    bool
		expectSuccess  bool
		expectToken    bool
	}{
		{
			name: "successful authentication",
			request: &models.AuthRequest{
				Username: username,
				Password: password,
			},
			setupMocks: func(repo *mocks.MockRepository) {
				user := &models.User{
					ID:           testID,
					Username:     username,
					Email:        email,
					PasswordHash: string(hashedPassword),
					CreatedAt:    time.Now(),
				}
				repo.On("GetByUsernameOrEmail", setup.Ctx, mock.AnythingOfType("*zap.Logger"), username).Return(user, nil)
			},
			expectError:   false,
			expectSuccess: true,
			expectToken:   true,
		},
		{
			name: "user not found",
			request: &models.AuthRequest{
				Username: "nonexistent",
				Password: password,
			},
			setupMocks: func(repo *mocks.MockRepository) {
				repo.On("GetByUsernameOrEmail", setup.Ctx, mock.AnythingOfType("*zap.Logger"), "nonexistent").Return(nil, errors.New("not found"))
			},
			expectError:   false,
			expectSuccess: false,
			expectToken:   false,
		},
		{
			name: "invalid password",
			request: &models.AuthRequest{
				Username: username,
				Password: "wrongpassword",
			},
			setupMocks: func(repo *mocks.MockRepository) {
				user := &models.User{
					ID:           testID,
					Username:     username,
					Email:        email,
					PasswordHash: string(hashedPassword),
					CreatedAt:    time.Now(),
				}
				repo.On("GetByUsernameOrEmail", setup.Ctx, mock.AnythingOfType("*zap.Logger"), username).Return(user, nil)
			},
			expectError:   false,
			expectSuccess: false,
			expectToken:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := setup.NewMockRepository()
			tt.setupMocks(mockRepo)

			svc := service.NewService(setup.Logger, mockRepo, setup.Config)

			result, err := svc.AuthenticateUser(setup.Ctx, setup.Logger, tt.request)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectSuccess, result.Success)
				
				if tt.expectToken {
					assert.NotNil(t, result.Token)
					assert.NotEmpty(t, *result.Token)
					assert.Nil(t, result.ErrorMessage)
				} else {
					assert.Nil(t, result.Token)
					assert.NotNil(t, result.ErrorMessage)
					assert.Equal(t, "invalid credentials", *result.ErrorMessage)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUsersService_ValidateUser(t *testing.T) {
	setup := mocks.NewTestSetup()
	testID := uuid.New()
	username := "testuser"
	email := "test@example.com"

	// Create a valid token for testing
	_ = service.NewService(setup.Logger, &mocks.MockRepository{}, setup.Config)
	validToken := generateTestToken(testID, username, email, setup.Config.JWTSecret)

	tests := []struct {
		name        string
		request     *models.ValidateUserRequest
		setupMocks  func(*mocks.MockRepository)
		expectError bool
		expectValid bool
	}{
		{
			name: "successful validation",
			request: &models.ValidateUserRequest{
				Token: validToken,
			},
			setupMocks: func(repo *mocks.MockRepository) {
				repo.On("GetGeneric", setup.Ctx, mock.AnythingOfType("*zap.Logger"), mock.AnythingOfType("*models.User")).Return(nil).Run(func(args mock.Arguments) {
					user := args.Get(2).(*models.User)
					user.ID = testID
					user.Username = username
					user.Email = email
					user.CreatedAt = time.Now()
				})
			},
			expectError: false,
			expectValid: true,
		},
		{
			name: "invalid token format",
			request: &models.ValidateUserRequest{
				Token: "invalid-token",
			},
			setupMocks: func(repo *mocks.MockRepository) {
				// No mock needed for invalid token
			},
			expectError: false,
			expectValid: false,
		},
		{
			name: "user not found from token",
			request: &models.ValidateUserRequest{
				Token: validToken,
			},
			setupMocks: func(repo *mocks.MockRepository) {
				repo.On("GetGeneric", setup.Ctx, mock.AnythingOfType("*zap.Logger"), mock.AnythingOfType("*models.User")).Return(errors.New("not found"))
			},
			expectError: false,
			expectValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := setup.NewMockRepository()
			tt.setupMocks(mockRepo)

			svc := service.NewService(setup.Logger, mockRepo, setup.Config)

			result, err := svc.ValidateUser(setup.Ctx, setup.Logger, tt.request)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectValid, result.Valid)
				
				if tt.expectValid {
					assert.NotNil(t, result.UserID)
					assert.Equal(t, testID, *result.UserID)
					assert.NotNil(t, result.Username)
					assert.Equal(t, username, *result.Username)
					assert.NotNil(t, result.Email)
					assert.Equal(t, email, *result.Email)
				} else {
					assert.Nil(t, result.UserID)
					assert.Nil(t, result.Username)
					assert.Nil(t, result.Email)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// Helper function to generate a test JWT token
func generateTestToken(userID uuid.UUID, username, email, jwtSecret string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  userID.String(),
		"username": username,
		"email":    email,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})
	
	tokenString, _ := token.SignedString([]byte(jwtSecret))
	return tokenString
}