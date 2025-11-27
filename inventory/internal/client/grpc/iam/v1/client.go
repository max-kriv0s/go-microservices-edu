package v1

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	def "github.com/max-kriv0s/go-microservices-edu/inventory/internal/client/grpc"
	authV1 "github.com/max-kriv0s/go-microservices-edu/shared/pkg/proto/auth/v1"
)

var _ def.IamServiceClient = (*iamServiceClient)(nil)

type iamServiceClient struct {
	client authV1.AuthServiceClient
}

func NewIamServiceClient(client authV1.AuthServiceClient) *iamServiceClient {
	return &iamServiceClient{
		client: client,
	}
}

func (c *iamServiceClient) Whoami(ctx context.Context, req *authV1.WhoamiRequest, opts ...grpc.CallOption) (*authV1.WhoamiResponse, error) {
	return c.client.Whoami(ctx, req)
}

func (c *iamServiceClient) Login(ctx context.Context, req *authV1.LoginRequest, opts ...grpc.CallOption) (*authV1.LoginResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method Login is not implemented")
}
