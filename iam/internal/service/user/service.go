package user

import (
	"github.com/alexedwards/argon2id"

	"github.com/max-kriv0s/go-microservices-edu/iam/internal/repository"
	def "github.com/max-kriv0s/go-microservices-edu/iam/internal/service"
)

var _ def.UserService = (*service)(nil)

type service struct {
	userRepository repository.UserRepository
}

func NewService(userRepository repository.UserRepository) *service {
	return &service{
		userRepository: userRepository,
	}
}

func (s *service) hashPassword(password string) (string, error) {
	params := argon2id.DefaultParams
	return argon2id.CreateHash(password, params)
}

func (s *service) VerifyPassword(password, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(password, hash)
}
