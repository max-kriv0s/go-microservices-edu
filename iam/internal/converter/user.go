package converter

import (
	"github.com/max-kriv0s/go-microservices-edu/iam/internal/model"
	commonV1 "github.com/max-kriv0s/go-microservices-edu/shared/pkg/proto/common/v1"
	userV1 "github.com/max-kriv0s/go-microservices-edu/shared/pkg/proto/user/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func RegisterUserRequestToModel(req *userV1.UserRegistrationInfo) model.UserInfo {

	var notificationMethod = make(map[model.ProviderName]model.Target)
	for _, method := range req.Info.GetNotificationMethods() {
		notificationMethod[model.ProviderName(method.GetProviderName())] = model.Target(method.GetTarget())
	}

	return model.UserInfo{
		Login:              req.GetInfo().GetLogin(),
		Email:              req.GetInfo().GetEmail(),
		NotificationMethod: notificationMethod,
	}
}

func UserToProto(user model.User) *commonV1.User {
	methods := make([]*commonV1.NotificationMethod, 0, len(user.UserInfo.NotificationMethod))
	for provider, target := range user.UserInfo.NotificationMethod {
		methods = append(methods, &commonV1.NotificationMethod{
			ProviderName: string(provider),
			Target:       string(target),
		})
	}

	return &commonV1.User{
		Uuid: user.UUID,
		Info: &commonV1.UserInfo{
			Login:               user.UserInfo.Login,
			Email:               user.UserInfo.Email,
			NotificationMethods: methods,
		},
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
}
