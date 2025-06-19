package usersvc

import (
	"context"
	"time"

	"server-events/internal/services/usersvc/models"

	"github.com/go-kit/log"
)

type (
	UserSvc interface {
		CreateUser(ctx context.Context, user *models.User) (bool, error)
		GetUsers(userID string, stream chan<- string) error
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
func (s *userSvc) GetUsers(userID string, stream chan<- string) error {
	for i := 0; i < 5; i++ {
		stream <- "Message " + userID
		time.Sleep(1 * time.Second)
	}

	close(stream)
	return nil
}
