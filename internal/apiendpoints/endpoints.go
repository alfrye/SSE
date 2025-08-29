package apiendpoints

import (
	"context"
	"fmt"

	"server-events/internal/services/usersvc"
	"server-events/internal/services/usersvc/models"
	pbapiv1 "server-events/pkg/genproto/pb"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	CreateUser endpoint.Endpoint
	GetUsers   endpoint.Endpoint
}

type CreateUserRequest struct {
	User *models.User `jsoon:"user,omitempty"`
}

type CreateUserResponse struct {
	Success bool `json:"success,omitempty"`
}

type GetUserRequest struct {
	UserID string
}

type StreamDataRequest struct {
	Req    *pbapiv1.GetUsersRequest
	Stream pbapiv1.UserSvc_GetUsersServer // Pass the stream directly
}
type GetUsersResponse struct {
	Message string
}

func NewEndpoints(u usersvc.UserSvc) *Endpoints {
	return &Endpoints{
		CreateUser: MakeCreateUserEndpoint(u),
		GetUsers:   MakeGetUsersEndpoint(u),
	}

}

func MakeGetUsersEndpoint(u usersvc.UserSvc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		fmt.Println("Executed make get suers endpoints")
		req := request.(GetUserRequest)
		result, err := u.GetUsers(ctx, req.UserID)
		return GetUsersResponse{Message: result}, nil
	}
}

//Example
// func MakeStreamDataEndpoint(s YourService) endpoint.Endpoint {
// 	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
// 		req := request.(streamDataRequest)
// 		return nil, s.StreamData(ctx, req.Req, req.Stream) // Call the service method directly
// 	}
// }

func MakeCreateUserEndpoint(u usersvc.UserSvc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateUserRequest)

		resp, err := u.CreateUser(ctx, req.User)

		if err != nil {
			return CreateUserResponse{
				Success: false,
			}, nil
		}

		return CreateUserResponse{
			Success: resp,
		}, nil

	}
}
