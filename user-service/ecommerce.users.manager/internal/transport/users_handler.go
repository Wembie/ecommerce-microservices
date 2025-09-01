package transport

import (
	"context"
	"ecommerce.users.manager/internal/config"
	"ecommerce.users.manager/internal/errors"
	"ecommerce.users.manager/internal/utils"

	pb "github.com/Wembie/ecommerce.users.manager.lib.protos/v1/libgo"
	"go.uber.org/zap"
)

func (h PingHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.User, error) {
	logger := config.CreateLoggerWithTraceID(h.log, ctx)
	logger.Info("Received CreateUser request", zap.Any("request", req))

	createReq, err := utils.ConvertToCreateUserModel(req)
	if err != nil {
		logger.Error("Error converting CreateUser request", zap.Error(err))
		return nil, errors.ToGRPCError(err)
	}

	user, err := h.pingSvc.CreateUser(ctx, logger, createReq)
	if err != nil {
		logger.Error("Error creating user", zap.Error(err))
		return nil, errors.ToGRPCError(err)
	}

	resp, err := utils.ConvertToUserProto(user)
	if err != nil {
		logger.Error("Error converting User response", zap.Error(err))
		return nil, errors.ToGRPCError(err)
	}

	logger.Info("CreateUser response", zap.Any("response", resp))
	return resp, nil
}

func (h PingHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
	logger := config.CreateLoggerWithTraceID(h.log, ctx)
	logger.Info("Received GetUser request", zap.Any("request", req))

	modelReq, err := utils.ConvertToGetUserModel(req)
	if err != nil {
		logger.Error("Invalid GetUser request", zap.Error(err))
		return nil, errors.ToGRPCError(err)
	}

	user, err := h.pingSvc.GetUser(ctx, logger, modelReq)
	if err != nil {
		logger.Error("Error getting user", zap.Error(err))
		return nil, errors.ToGRPCError(err)
	}

	resp, err := utils.ConvertToUserProto(user)
	if err != nil {
		logger.Error("Error converting User response", zap.Error(err))
		return nil, errors.ToGRPCError(err)
	}

	logger.Info("GetUser response", zap.Any("response", resp))
	return resp, nil
}

func (h PingHandler) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.User, error) {
	logger := config.CreateLoggerWithTraceID(h.log, ctx)
	logger.Info("Received UpdateUser request", zap.Any("request", req))

	updateReq, err := utils.ConvertToUpdateUserModel(req)
	if err != nil {
		logger.Error("Error converting UpdateUser request", zap.Error(err))
		return nil, errors.ToGRPCError(err)
	}

	user, err := h.pingSvc.UpdateUser(ctx, logger, updateReq)
	if err != nil {
		logger.Error("Error updating user", zap.Error(err))
		return nil, errors.ToGRPCError(err)
	}

	resp, err := utils.ConvertToUserProto(user)
	if err != nil {
		logger.Error("Error converting User response", zap.Error(err))
		return nil, errors.ToGRPCError(err)
	}

	logger.Info("UpdateUser response", zap.Any("response", resp))
	return resp, nil
}

func (h PingHandler) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	logger := config.CreateLoggerWithTraceID(h.log, ctx)
	logger.Info("Received DeleteUser request", zap.Any("request", req))

	modelReq, err := utils.ConvertToDeleteUserModel(req)
	if err != nil {
		logger.Error("Invalid DeleteUser request", zap.Error(err))
		return nil, errors.ToGRPCError(err)
	}

	success, err := h.pingSvc.DeleteUser(ctx, logger, modelReq)
	if err != nil {
		logger.Error("Error deleting user", zap.Error(err))
		return nil, errors.ToGRPCError(err)
	}

	resp := &pb.DeleteUserResponse{Success: success}
	logger.Info("DeleteUser response", zap.Any("response", resp))
	return resp, nil
}

func (h PingHandler) AuthenticateUser(ctx context.Context, req *pb.AuthRequest) (*pb.AuthResponse, error) {
	logger := config.CreateLoggerWithTraceID(h.log, ctx)
	logger.Info("Received AuthenticateUser request", zap.Any("request", req))

	modelReq, err := utils.ConvertToAuthUserModel(req)
	if err != nil {
		logger.Error("Invalid AuthUser request", zap.Error(err))
		return nil, errors.ToGRPCError(err)
	}

	authResp, err := h.pingSvc.AuthenticateUser(ctx, logger, modelReq)
	if err != nil {
		logger.Error("Error authenticating user", zap.Error(err))
		return nil, errors.ToGRPCError(err)
	}

	resp, err := utils.ConvertToAuthUserProto(authResp)
	if err != nil {
		logger.Error("Error converting AuthUser response", zap.Error(err))
		return nil, errors.ToGRPCError(err)
	}

	logger.Info("AuthenticateUser response", zap.Any("response", resp))
	return resp, nil
}

func (h PingHandler) ValidateUser(ctx context.Context, req *pb.ValidateUserRequest) (*pb.ValidateUserResponse, error) {
	logger := config.CreateLoggerWithTraceID(h.log, ctx)
	logger.Info("Received ValidateUser request", zap.Any("request", req))

	modelReq, err := utils.ConvertToValidateUserModel(req)
	if err != nil {
		logger.Error("Invalid ValidateUser request", zap.Error(err))
		return nil, errors.ToGRPCError(err)
	}

	validateResp, err := h.pingSvc.ValidateUser(ctx, logger, modelReq)
	if err != nil {
		logger.Error("Error validating user", zap.Error(err))
		return nil, errors.ToGRPCError(err)
	}

	resp, err := utils.ConvertToValidateUserProto(validateResp)
	if err != nil {
		logger.Error("Error converting ValidateUser response", zap.Error(err))
		return nil, errors.ToGRPCError(err)
	}

	logger.Info("ValidateUser response", zap.Any("response", resp))
	return resp, nil
}
