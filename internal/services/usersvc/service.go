package usersvc

import (
	"context"

	"server-events/internal/services/usersvc/models"
	pbapiv1 "server-events/pkg/genproto/pb"

	"github.com/go-kit/log"
	"google.golang.org/protobuf/types/known/emptypb"
)

type (
	UserSvc interface {
		CreateUser(ctx context.Context, user *models.User) (bool, error)
		StreamUpdates(req *emptypb.Empty, stream pbapiv1.UserSvc_StreamUpdatesServer) (*models.Updates, error)
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

func (s *userSvc) StreamUpdates(req *emptypb.Empty, stream pbapiv1.UserSvc_StreamUpdatesServer) (*models.Updates, error) {
	return &models.Updates{
		Message: "message",
	}, nil
}
