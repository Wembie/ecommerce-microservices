package service

import (
	"context"
	"time"

	"ecommerce.users.manager/internal/models"
	"ecommerce.users.manager/internal/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func (s *UsersService) CreateUser(ctx context.Context, logger *zap.Logger, request *models.CreateUserRequest) (*models.User, error) {
	logger.Info("Creating User", zap.Any("request", request))

	hashed, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("Failed to hash password", zap.Error(err))
		return nil, err
	}

	model := &models.User{
		Username:     request.Username,
		Email:        request.Email,
		PasswordHash: string(hashed),
		CreatedAt:    time.Now(),
	}

	if err := s.userRepo.CreateGeneric(ctx, logger, model); err != nil {
		logger.Error("Failed to create user", zap.Error(err))
		return nil, err
	}

	logger.Info("User created successfully", zap.String("id", model.ID.String()))
	return model, nil
}

func (s *UsersService) GetUser(ctx context.Context, logger *zap.Logger, request *models.GetUserRequest) (*models.User, error) {
	logger.Info("Getting User", zap.String("id", request.ID.String()))

	model := &models.User{ID: request.ID}
	if err := s.userRepo.GetGeneric(ctx, logger, model); err != nil {
		logger.Error("Failed to get user", zap.Error(err))
		return nil, err
	}

	logger.Info("User retrieved successfully", zap.Any("user", model))
	return model, nil
}

func (s *UsersService) UpdateUser(ctx context.Context, logger *zap.Logger, request *models.UpdateUserRequest) (*models.User, error) {
	logger.Info("Updating User", zap.Any("request", request))

	if request.Password != nil {
		hashed, err := bcrypt.GenerateFromPassword([]byte(*request.Password), bcrypt.DefaultCost)
		if err != nil {
			logger.Error("Failed to hash new password", zap.Error(err))
			return nil, err
		}
		passwordHash := string(hashed)
		request.Password = &passwordHash
	}

	rowsAffected, err := s.userRepo.UpdateGeneric(ctx, logger, request)
	if err != nil {
		logger.Error("Failed to update user", zap.Error(err))
		return nil, err
	}

	if rowsAffected == 0 {
		logger.Warn("No rows were updated, user may not exist", zap.Any("request", request))
		return nil, nil
	}

	model := &models.User{ID: request.ID}
	if err := s.userRepo.GetGeneric(ctx, logger, model); err != nil {
		logger.Error("Failed to retrieve updated user", zap.Error(err))
		return nil, err
	}

	logger.Info("User updated successfully", zap.Any("user", model))
	return model, nil
}

func (s *UsersService) DeleteUser(ctx context.Context, logger *zap.Logger, request *models.DeleteUserRequest) (bool, error) {
	logger.Info("Deleting User", zap.String("id", request.ID.String()))

	model := &models.User{ID: request.ID}

	err := s.userRepo.DeleteGeneric(ctx, logger, model)
	if err != nil {
		logger.Error("Failed to delete user", zap.Error(err))
		return false, err
	}

	logger.Info("User deleted successfully", zap.String("id", request.ID.String()))
	return true, nil
}


func (s *UsersService) AuthenticateUser(ctx context.Context, logger *zap.Logger, request *models.AuthRequest) (*models.AuthResponse, error) {
	logger.Info("Authenticating User", zap.String("username", request.Username))

	user, err := s.userRepo.GetByUsernameOrEmail(ctx, logger, request.Username)
	if err != nil {
		logger.Error("User not found", zap.Error(err))
		return &models.AuthResponse{
			Success:      false,
			ErrorMessage: utils.StringPtr("invalid credentials"),
		}, nil
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.Password)); err != nil {
		logger.Error("Invalid password", zap.Error(err))
		return &models.AuthResponse{
			Success:      false,
			ErrorMessage: utils.StringPtr("invalid credentials"),
		}, nil
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"email":    user.Email,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(s.jwtKey)
	if err != nil {
		logger.Error("Failed to sign token", zap.Error(err))
		return &models.AuthResponse{
			Success:      false,
			ErrorMessage: utils.StringPtr("failed to generate token"),
		}, nil
	}

	resp := &models.AuthResponse{
		Token:   &tokenString,
		Success: true,
	}
	logger.Info("User authenticated successfully", zap.String("username", user.Username))
	return resp, nil
}

func (s *UsersService) ValidateUser(ctx context.Context, logger *zap.Logger, request *models.ValidateUserRequest) (*models.ValidateUserResponse, error) {
	logger.Info("Validating User", zap.String("token", request.Token))

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(request.Token, claims, func(t *jwt.Token) (interface{}, error) {
		return s.jwtKey, nil
	})
	if err != nil || !token.Valid {
		logger.Error("Invalid token", zap.Error(err))
		return &models.ValidateUserResponse{Valid: false}, nil
	}

	userIDStr, okID := claims["user_id"].(string)
	username, okUsername := claims["username"].(string)
	email, okEmail := claims["email"].(string)

	if !okID || !okUsername || !okEmail {
		logger.Error("Token missing required claims")
		return &models.ValidateUserResponse{Valid: false}, nil
	}

	userID, err := uuid.Parse(userIDStr)
		if err != nil {
		logger.Error("Invalid user_id format in token", zap.Error(err))
		return &models.ValidateUserResponse{Valid: false}, nil
	}

	user := &models.User{ID: userID}
	err = s.userRepo.GetGeneric(ctx, logger, user)
	if err != nil {
		logger.Error("User not found from token", zap.Error(err))
		return &models.ValidateUserResponse{Valid: false}, nil
	}

	resp := &models.ValidateUserResponse{
		Valid:    true,
		UserID:   &userID,
		Username: &username,
		Email:    &email,
	}
	logger.Info("Token validated successfully", zap.String("username", user.Username))
	return resp, nil
}