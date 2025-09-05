package clients

import (
	"fmt"
	
	pb "github.com/Wembie/ecommerce.users.manager.lib.protos/v1/libgo"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)


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
