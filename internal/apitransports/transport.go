package apitransports

import (
	"context"
	"errors"
	"fmt"
	"server-events/internal/apiendpoints"
	"server-events/internal/services/usersvc/models"
	pbapiv1 "server-events/pkg/genproto/pb"

	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/go-kit/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserAPIServer struct {
	createUser grpctransport.Handler
	getUsers   grpctransport.Handler
	pbapiv1.UnimplementedUserSvcServer
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

func (s *UserAPIServer) GetUsers(req *pbapiv1.GetUsersRequest, stream pbapiv1.UserSvc_GetUsersServer) error {
	_, _, err := s.getUsers.ServeGRPC(stream.Context(), &apiendpoints.GetUserRequest{
		Req:    req,
		Stream: stream,
	})

	if err != nil {
		if errors.Is(err, context.Canceled) {
			return status.Errorf(codes.Canceled, "stream cancelled")
		}
		return status.Errorf(codes.Internal, "internal server error:%v", err)
	}

	return nil
}

// From example
// func (s *grpcServer) StreamData(req *pb.MyStreamingRequest, stream pb.YourService_StreamDataServer) error {
// 	_, err := s.streamData.ServeGRPC(stream.Context(), &service.StreamDataRequest{Req: req, Stream: stream})
// 	if err != nil {
// 		// Handle errors from your service layer
// 		if errors.Is(err, context.Canceled) {
// 			return status.Errorf(codes.Canceled, "stream cancelled")
// 		}
// 		// More sophisticated error handling for gRPC codes
// 		return status.Errorf(codes.Internal, "internal server error: %v", err)
// 	}
// 	return nil
// }

func decodeGetUsersRequest(_ context.Context, r interface{}) (interface{}, error) {
	fmt.Printf("Incoming request:%t\n", r)
	req := r.(*pbapiv1.GetUsersRequest)
	if req == nil {
		fmt.Println("request is nil")
		return nil, fmt.Errorf("request is nil")
	}
	return apiendpoints.GetUserRequest{
		Req: req,
	}, nil
}

// func decodeGRPCStreamDataRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
// 	// For streaming, the request is the initial message from the client.
// 	// The stream itself is passed implicitly.
// 	req := grpcReq.(*pb.MyStreamingRequest)
// 	return service.StreamDataRequest{Req: req}, nil
// }

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
	fmt.Printf("Incoming request:%t\n", r)
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
