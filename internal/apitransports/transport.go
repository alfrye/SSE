package apitransports

import (
	"context"
	"fmt"
	"server-events/internal/apiendpoints"
	"server-events/internal/services/usersvc/models"
	pbapiv1 "server-events/pkg/genproto/pb"

	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/go-kit/log"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserAPIServer struct {
	createUser    grpctransport.Handler
	streamUpdates grpctransport.Handler
}

func NewUserAPIServer(ep *apiendpoints.Endpoints, logger log.Logger) pbapiv1.UserSvcServer {

	return &UserAPIServer{
		createUser: grpctransport.NewServer(
			ep.CreateUser,
			decodeCreateUserRequest,
			encodeCreateUserResponse,
		),
		streamUpdates: grpctransport.NewServer(
			ep.StreamUpdates,
			decodeStreamUpdateRequest,
			encodeStreamUpdateResponse,
		),
	}
}

func (s *UserAPIServer) StreamUpdates(ctx context.Context, r *emptypb.Empty) (*pbapiv1.UserSvc_StreamUpdatesServer, error) {
	_, resp, err := s.streamUpdates.ServeGRPC(ctx, r)

	if err != nil {
		return nil, err
	}

	return resp.(*pbapiv1.UserSvc_StreamUpdatesServer), nil
}

func decodeStreamUpdateRequest(_ context.Context, r interface{}) (interface{}, error) {
	fmt.Println("calling decodeStreamUpdateRequest")
	req := r.(*pbapiv1.UpdateRequest)
	if req == nil {
		fmt.Println("request is nil")
		return nil, fmt.Errorf("request is nil")
	}
	return apiendpoints.UpdateRequest{}, nil
}

func encodeStreamUpdateResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(apiendpoints.UpdateResponse)

	return &pbapiv1.UpdateResponse{
		Message: resp.Message,
	}, nil
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
