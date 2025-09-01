package utils

import (
	"ecommerce.users.manager/internal/errors"
	"ecommerce.users.manager/internal/models"
	"github.com/asaskevich/govalidator"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/Wembie/ecommerce.users.manager.lib.protos/v1/libgo"
)

func ConvertToUserProto(m *models.User) (*pb.User, error) {
	if m == nil {
		return nil, errors.New(errors.ErrCodeInvalidArgument, "user model is nil")
	}

	return &pb.User{
		Id:          m.ID.String(),
		Username:    m.Username,
		Email:       m.Email,
		PasswordHash: m.PasswordHash,
		CreatedAt:   timestamppb.New(m.CreatedAt),
		UpdatedAt: func() *timestamppb.Timestamp {
			if m.UpdatedAt != nil {
				return timestamppb.New(*m.UpdatedAt)
			}
			return nil
		}(),
	}, nil
}

func ConvertToCreateUserModel(p *pb.CreateUserRequest) (*models.CreateUserRequest, error) {
	if p == nil {
		return nil, errors.New(errors.ErrCodeInvalidArgument, "CreateUserRequest is nil")
	}

	if p.GetUsername() == "" {
		return nil, errors.New(errors.ErrCodeInvalidArgument, "username is required")
	}

	if p.GetEmail() == "" {
		return nil, errors.New(errors.ErrCodeInvalidArgument, "email is required")
	}

	if !govalidator.IsEmail(p.GetEmail()) {
    	return nil, errors.New(errors.ErrCodeInvalidArgument, "invalid email format")
	}

	if p.GetPassword() == "" {
		return nil, errors.New(errors.ErrCodeInvalidArgument, "password is required")
	}

	return &models.CreateUserRequest{
		Username: p.GetUsername(),
		Email:    p.GetEmail(),
		Password: p.GetPassword(),
	}, nil
}

func ConvertToGetUserModel(p *pb.GetUserRequest) (*models.GetUserRequest, error) {
	if p == nil {
		return nil, errors.New(errors.ErrCodeInvalidArgument, "GetUserRequest is nil")
	}

	id, err := uuid.Parse(p.GetId())
	if err != nil {
		return nil, errors.New(errors.ErrCodeInvalidArgument, "invalid id: "+err.Error())
	}

	return &models.GetUserRequest{ID: id}, nil
}

func ConvertToUpdateUserModel(p *pb.UpdateUserRequest) (*models.UpdateUserRequest, error) {
	if p == nil {
		return nil, errors.New(errors.ErrCodeInvalidArgument, "UpdateUserRequest is nil")
	}

	id, err := uuid.Parse(p.GetId())
	if err != nil {
		return nil, errors.New(errors.ErrCodeInvalidArgument, "invalid id: "+err.Error())
	}

	req := &models.UpdateUserRequest{
		ID:        id,
		UpdatedAt: time.Now(),
	}

	if p.Username != nil {
		if *p.Username == "" {
			return nil, errors.New(errors.ErrCodeInvalidArgument, "username cannot be empty")
		}
		req.Username = p.Username
	}

	if p.Email != nil {
		if *p.Email == "" {
			return nil, errors.New(errors.ErrCodeInvalidArgument, "email cannot be empty")
		}
		if !govalidator.IsEmail(p.GetEmail()) {
			return nil, errors.New(errors.ErrCodeInvalidArgument, "invalid email format")
		}
		req.Email = p.Email
	}

	if p.Password != nil {
		if *p.Password == "" {
			return nil, errors.New(errors.ErrCodeInvalidArgument, "password cannot be empty")
		}
		req.Password = p.Password
	}

	return req, nil
}

func ConvertToDeleteUserModel(p *pb.DeleteUserRequest) (*models.DeleteUserRequest, error) {
	if p == nil {
		return nil, errors.New(errors.ErrCodeInvalidArgument, "DeleteUserRequest is nil")
	}

	id, err := uuid.Parse(p.GetId())
	if err != nil {
		return nil, errors.New(errors.ErrCodeInvalidArgument, "invalid id: "+err.Error())
	}

	return &models.DeleteUserRequest{ID: id}, nil
}

func ConvertToDeleteUserProto(m *models.DeleteUserResponse) (*pb.DeleteUserResponse, error) {
	if m == nil {
		return nil, errors.New(errors.ErrCodeInvalidArgument, "DeleteUserResponse is nil")
	}

	return &pb.DeleteUserResponse{
		Success: m.Success,
	}, nil
}

func ConvertToAuthUserModel(p *pb.AuthRequest) (*models.AuthRequest, error) {
	if p == nil {
		return nil, errors.New(errors.ErrCodeInvalidArgument, "AuthRequest is nil")
	}

	if p.GetUsername() == "" || p.GetPassword() == "" {
		return nil, errors.New(errors.ErrCodeInvalidArgument, "username and password are required")
	}

	return &models.AuthRequest{
		Username: p.GetUsername(),
		Password: p.GetPassword(),
	}, nil
}

func ConvertToAuthUserProto(m *models.AuthResponse) (*pb.AuthResponse, error) {
	if m == nil {
		return nil, errors.New(errors.ErrCodeInvalidArgument, "AuthResponse is nil")
	}

	return &pb.AuthResponse{
		Success:      m.Success,
		Token:        m.Token,
		ErrorMessage: m.ErrorMessage,
	}, nil
}

func ConvertToValidateUserModel(p *pb.ValidateUserRequest) (*models.ValidateUserRequest, error) {
	if p == nil {
		return nil, errors.New(errors.ErrCodeInvalidArgument, "ValidateUserRequest is nil")
	}

	if p.GetToken() == "" {
		return nil, errors.New(errors.ErrCodeInvalidArgument, "token is required")
	}

	return &models.ValidateUserRequest{Token: p.GetToken()}, nil
}

func ConvertToValidateUserProto(m *models.ValidateUserResponse) (*pb.ValidateUserResponse, error) {
	if m == nil {
		return nil, errors.New(errors.ErrCodeInvalidArgument, "ValidateUserResponse is nil")
	}

	var userID *string
	if m.UserID != nil {
		idStr := m.UserID.String()
		userID = &idStr
	}

	return &pb.ValidateUserResponse{
		Valid:    m.Valid,
		UserId:   userID,
		Username: m.Username,
		Email:    m.Email,
	}, nil
}