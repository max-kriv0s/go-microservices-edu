package auth

import (
	"time"

	"github.com/max-kriv0s/go-microservices-edu/iam/internal/repository"
	def "github.com/max-kriv0s/go-microservices-edu/iam/internal/service"
)

var _ def.AuthService = (*service)(nil)

type service struct {
	userService       def.UserService
	sessionRepository repository.SessonRepository
	sessionTTL        time.Duration
}

func NewService(userService def.UserService, sessionRepository repository.SessonRepository, sessionTTL time.Duration) *service {
	return &service{
		userService:       userService,
		sessionRepository: sessionRepository,
		sessionTTL:        sessionTTL,
	}
}
