package converter

import (
	"encoding/json"

	"github.com/max-kriv0s/go-microservices-edu/iam/internal/model"
	repoModel "github.com/max-kriv0s/go-microservices-edu/iam/internal/repository/model"
)

func UserInfoToRepoModel(userInfo model.UserInfo, passwordHash string) (repoModel.User, error) {
	var notificationMethod = make(map[repoModel.ProviderName]repoModel.Target, len(userInfo.NotificationMethod))

	for providerName, target := range userInfo.NotificationMethod {
		notificationMethod[repoModel.ProviderName(providerName)] = repoModel.Target(target)
	}

	notificationMethodBytes, err := json.Marshal(notificationMethod)
	if err != nil {
		return repoModel.User{}, err
	}

	return repoModel.User{
		Login:              userInfo.Login,
		Email:              userInfo.Email,
		NotificationMethod: notificationMethodBytes,
		PasswordHash:       passwordHash,
	}, nil
}

func UserToModel(repoUser repoModel.User) (model.User, error) {
	var notificationMethod = make(map[model.ProviderName]model.Target)
	if len(repoUser.NotificationMethod) > 0 {
		var tmp map[string]string

		err := json.Unmarshal(repoUser.NotificationMethod, &tmp)
		if err != nil {
			return model.User{}, err
		}

		for k, v := range tmp {
			notificationMethod[model.ProviderName(k)] = model.Target(v)
		}
	}

	user := model.User{
		UUID: *repoUser.UUID,
		UserInfo: model.UserInfo{
			Login:              repoUser.Login,
			Email:              repoUser.Email,
			NotificationMethod: notificationMethod,
		},
		PasswordHash: repoUser.PasswordHash,
		CreatedAt:    repoUser.CreatedAt,
		UpdatedAt:    repoUser.UpdatedAt,
	}

	return user, nil
}
