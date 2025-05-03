package usersvc

import (
	"context"

	"server-events/internal/services/usersvc/models"

	"github.com/go-kit/log"
)

type (
	UserSvc interface {
		CreateUser(ctx context.Context, user *models.User) (bool, error)
	}

	userSvc struct {
		logger log.Logger
	}
)

func NewService() UserSvc {
	var svc UserSvc
	svc = &userSvc{
		logger: log.NewNopLogger(),
	}

	return svc
}

func (s *userSvc) CreateUser(_ context.Context, user *models.User) (bool, error) {

	return true, nil

}
