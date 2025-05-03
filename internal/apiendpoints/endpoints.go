package apiendpoints

import (
	"context"

	"server-events/internal/services/usersvc"
	"server-events/internal/services/usersvc/models"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	CreateUser endpoint.Endpoint
}

type CreateUserRequest struct {
	User *models.User `jsoon:"user,omitempty"`
}

type CreateUserResponse struct {
	Success bool `json:"success,omitempty"`
}

func NewEndpoints(u usersvc.UserSvc) *Endpoints {
	return &Endpoints{
		CreateUser: MakeCreateUserEndpoint(u),
	}

}

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
