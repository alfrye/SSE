package apiendpoints

import (
	"context"

	"server-events/internal/services/usersvc"
	"server-events/internal/services/usersvc/models"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	CreateUser    endpoint.Endpoint
	StreamUpdates endpoint.Endpoint
}

type CreateUserRequest struct {
	User *models.User `jsoon:"user,omitempty"`
}

type CreateUserResponse struct {
	Success bool `json:"success,omitempty"`
}

type UpdateRequest struct {
}

type UpdateResponse struct {
	Message string `json:"message,omitempty"`
}

func NewEndpoints(u usersvc.UserSvc) *Endpoints {
	return &Endpoints{
		CreateUser:    MakeCreateUserEndpoint(u),
		StreamUpdates: MakeStreamUpdatesEndpoint(u),
	}

}

func MakeStreamUpdatesEndpoint(u usersvc.UserSvc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		_ = request.(UpdateRequest)
		resp, err := u.StreamUpdates(ctx)

		if err != nil {
			return nil, err
		}

		return UpdateResponse{
			Message: resp.Message,
		}, nil

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
