package clients

import (
	"context"

	"ecommerce.auth/internal/models"

	pb "github.com/Wembie/ecommerce.users.manager.lib.protos/v1/libgo"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (c *userClient) AuthenticateUser(ctx context.Context, req *models.AuthRequest) (*models.AuthResponse, error) {
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

	return &models.AuthResponse{
		Success:      resp.Success,
		Token:        token,
		ErrorMessage: errorMessage,
	}, nil
}

func (c *userClient) ValidateUser(ctx context.Context, token string) (*models.ValidateUserResponse, error) {
	c.logger.Info("Validating user token")

	req := &pb.ValidateUserRequest{
		Token: token,
	}

	resp, err := c.client.ValidateUser(ctx, req)
	if err != nil {
		c.logger.Error("Failed to validate user", zap.Error(err))
		return nil, err
	}

	validateResponse := &models.ValidateUserResponse{
		Valid: resp.Valid,
	}

	if resp.UserId != nil {
		uid, err := uuid.Parse(*resp.UserId)
		if err != nil {
			c.logger.Error("Failed to parse user ID as UUID", zap.Error(err))
		} else {
			validateResponse.UserID = &uid
		}
	}
	if resp.Username != nil {
		validateResponse.Username = resp.Username
	}
	if resp.Email != nil {
		validateResponse.Email = resp.Email
	}

	return validateResponse, nil
}

func (c *userClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}