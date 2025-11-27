package client

import (
	"context"

	"google.golang.org/grpc"

	authV1 "github.com/max-kriv0s/go-microservices-edu/shared/pkg/proto/auth/v1"
)

type IamServiceClient interface {
	Login(ctx context.Context, req *authV1.LoginRequest, opts ...grpc.CallOption) (*authV1.LoginResponse, error)
	Whoami(ctx context.Context, req *authV1.WhoamiRequest, opts ...grpc.CallOption) (*authV1.WhoamiResponse, error)
}
