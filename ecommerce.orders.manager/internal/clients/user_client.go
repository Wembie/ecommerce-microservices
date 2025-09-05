package clients

import (
	"context"
	"fmt"

	"ecommerce.orders.manager/internal/models"

	pb "github.com/Wembie/ecommerce.users.manager.lib.protos/v1/libgo"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

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
	Valid    bool     `json:"valid"`
	UserId   *string  `json:"user_id,omitempty"`
	Username *string  `json:"username,omitempty"`
	Email    *string  `json:"email,omitempty"`
}

type UserClient interface {
	AuthenticateUser(ctx context.Context, req *AuthRequest) (*AuthResponse, error)
	ValidateUser(ctx context.Context, token string) (*ValidateUserResponse, error)
	ValidateUserInfo(ctx context.Context, token string) (*models.UserInfo, error)
	GetUser(ctx context.Context, userID uuid.UUID) (*models.UserInfo, error)
	Close() error
}

type userClient struct {
	conn   *grpc.ClientConn
	client pb.UserServiceClient
	logger *zap.Logger
}

func NewUserClient(address string, logger *zap.Logger) (UserClient, error) {
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC client: %w", err)
	}

	conn.Connect()

	client := pb.NewUserServiceClient(conn)

	return &userClient{
		conn:   conn,
		client: client,
		logger: logger,
	}, nil
}


func (c *userClient) AuthenticateUser(ctx context.Context, req *AuthRequest) (*AuthResponse, error) {
	c.logger.Info("Authenticating user", zap.String("username", req.Username))

	pbReq := &pb.AuthRequest{
		Username: req.Username,
		Password: req.Password,
	}

	resp, err := c.client.AuthenticateUser(ctx, pbReq)
	if err != nil {
		c.logger.Error("Failed to authenticate user", zap.Error(err))
		return nil, err
	}

	var token *string
	if resp.Token != nil && *resp.Token != "" {
		token = resp.Token
	}

	var errorMessage *string
	if resp.ErrorMessage != nil && *resp.ErrorMessage != "" {
		errorMessage = resp.ErrorMessage
	}

	return &AuthResponse{
		Success:      resp.Success,
		Token:        token,
		ErrorMessage: errorMessage,
	}, nil
}

func (c *userClient) ValidateUser(ctx context.Context, token string) (*ValidateUserResponse, error) {
	c.logger.Info("Validating user token")

	req := &pb.ValidateUserRequest{
		Token: token,
	}

	resp, err := c.client.ValidateUser(ctx, req)
	if err != nil {
		c.logger.Error("Failed to validate user", zap.Error(err))
		return nil, err
	}

	validateResponse := &ValidateUserResponse{
		Valid: resp.Valid,
	}

	if resp.UserId != nil {
		validateResponse.UserId = resp.UserId
	}
	if resp.Username != nil {
		validateResponse.Username = resp.Username
	}
	if resp.Email != nil {
		validateResponse.Email = resp.Email
	}

	return validateResponse, nil
}

func (c *userClient) ValidateUserInfo(ctx context.Context, token string) (*models.UserInfo, error) {
	c.logger.Info("Validating user token for user info")

	req := &pb.ValidateUserRequest{
		Token: token,
	}

	resp, err := c.client.ValidateUser(ctx, req)
	if err != nil {
		c.logger.Error("Failed to validate user", zap.Error(err))
		return nil, err
	}

	if !resp.Valid {
		c.logger.Warn("Invalid token provided")
		return nil, fmt.Errorf("invalid token")
	}

	userID, err := uuid.Parse(resp.GetUserId())
	if err != nil {
		c.logger.Error("Failed to parse user ID", zap.Error(err))
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	userInfo := &models.UserInfo{
		ID:       userID,
		Username: resp.GetUsername(),
		Email:    resp.GetEmail(),
	}

	return userInfo, nil
}

func (c *userClient) GetUser(ctx context.Context, userID uuid.UUID) (*models.UserInfo, error) {
	c.logger.Info("Getting user info", zap.String("user_id", userID.String()))

	req := &pb.GetUserRequest{
		Id: userID.String(),
	}

	resp, err := c.client.GetUser(ctx, req)
	if err != nil {
		c.logger.Error("Failed to get user", zap.Error(err))
		return nil, err
	}

	parsedID, err := uuid.Parse(resp.Id)
	if err != nil {
		c.logger.Error("Failed to parse user ID", zap.Error(err))
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	userInfo := &models.UserInfo{
		ID:       parsedID,
		Username: resp.Username,
		Email:    resp.Email,
	}

	return userInfo, nil
}

func (c *userClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}