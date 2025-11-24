package user

import (
	"context"
	"errors"
	"strings"

	"github.com/max-kriv0s/go-microservices-edu/iam/internal/model"
)

func (s *service) Register(ctx context.Context, userInfo model.UserInfo, password string) (string, error) {
	userInfo.Email = strings.ToLower(userInfo.Email)

	// Ищем пользователя по email
	_, err := s.userRepository.GetByEmail(ctx, userInfo.Email)
	if err == nil {
		return "", model.ErrLoginOrEmailExists
	}
	if !errors.Is(err, model.ErrUserNotFound) {
		return "", err
	}

	// Ищем пользователя по login
	_, err = s.userRepository.GetByLogin(ctx, userInfo.Login)
	if err == nil {
		return "", model.ErrLoginOrEmailExists
	}
	if !errors.Is(err, model.ErrUserNotFound) {
		return "", err
	}

	// Создаем нового пользователя
	passwordHash, err := s.hashPassword(password)
	if err != nil {
		return "", model.ErrFailedToHashPassword
	}

	userUUID, err := s.userRepository.Create(ctx, userInfo, passwordHash)
	if err != nil {
		return "", err
	}

	return userUUID, nil
}
