package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/max-kriv0s/go-microservices-edu/iam/internal/model"
	authV1 "github.com/max-kriv0s/go-microservices-edu/shared/pkg/proto/auth/v1"
	commonV1 "github.com/max-kriv0s/go-microservices-edu/shared/pkg/proto/common/v1"
)

func SessionInfoToWhoamiResponse(sessionInfo model.SessionInfo) *authV1.WhoamiResponse {
	return &authV1.WhoamiResponse{
		Session: &commonV1.Session{
			Uuid:      sessionInfo.Session.Uuid,
			CreatedAt: timestamppb.New(sessionInfo.Session.CreatedAt),
			UpdatedAt: timestamppb.New(sessionInfo.Session.UpdatedAt),
			ExpiresAt: timestamppb.New(sessionInfo.Session.ExpiresAt),
		},
		User: UserToProto(sessionInfo.User),
	}
}
