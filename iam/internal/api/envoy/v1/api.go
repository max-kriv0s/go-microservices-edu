package v1

import (
	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"

	authV1 "github.com/max-kriv0s/go-microservices-edu/shared/pkg/proto/auth/v1"
)

type api struct {
	authv3.UnimplementedAuthorizationServer

	authV1API authV1.AuthServiceServer
}

func NewApi(authV1API authV1.AuthServiceServer) *api {
	return &api{
		authV1API: authV1API,
	}
}
