package apitransports

import (
	"context"
	"fmt"
	"server-events/internal/apiendpoints"
	"server-events/internal/services/usersvc/models"
	pbapiv1 "server-events/pkg/genproto/pb"

	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/go-kit/log"
)

type UserAPIServer struct {
	createUser grpctransport.Handler
	getUsers   grpctransport.Handler
}

func NewUserAPIServer(ep *apiendpoints.Endpoints, logger log.Logger) pbapiv1.UserSvcServer {
	return &UserAPIServer{
		createUser: grpctransport.NewServer(
			ep.CreateUser,
			decodeCreateUserRequest,
			encodeCreateUserResponse,
		),
		getUsers: grpctransport.NewServer(
			ep.GetUsers,
			decodeGetUsersRequest,
			encodeGetUsersResponse,
		),
	}
}

func (s *UserAPIServer) GetUsers(ctx context.Context, r *pbapiv1.GetUsersRequest) (*pbapiv1.GetUsersResponse, error) {
	_, resp, err := s.getUsers.ServeGRPC(ctx, r)

	if err != nil {
		return nil, err
	}

	return resp.(*pbapiv1.GetUsersResponse), nil
}
func decodeGetUsersRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pbapiv1.GetUsersRequest)
	if req == nil {
		fmt.Println("request is nil")
		return nil, fmt.Errorf("request is nil")
	}
	return apiendpoints.GetUserRequest{
		UserID: req.UserId,
	}, nil
}

func encodeGetUsersResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(*pbapiv1.GetUsersResponse)

	return &pbapiv1.GetUsersResponse{Message: resp.Message}, nil
}

func (s *UserAPIServer) CreateUser(ctx context.Context, r *pbapiv1.CreateUserRequest) (*pbapiv1.CreateUserResponse, error) {
	_, resp, err := s.createUser.ServeGRPC(ctx, r)

	if err != nil {
		return nil, err
	}

	return resp.(*pbapiv1.CreateUserResponse), nil
}

func decodeCreateUserRequest(_ context.Context, r interface{}) (interface{}, error) {
	fmt.Println("calling decodeCreateUserRequest")
	req := r.(*pbapiv1.CreateUserRequest)
	if req == nil {
		fmt.Println("request is nil")
		return nil, fmt.Errorf("request is nil")
	}
	return apiendpoints.CreateUserRequest{
		User: &models.User{
			ID:       req.Id,
			UserName: req.Username,
		},
	}, nil
}

func encodeCreateUserResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(apiendpoints.CreateUserResponse)

	return &pbapiv1.CreateUserResponse{
		Success: resp.Success,
	}, nil
}
