package usersvc

import (
	"context"
	"fmt"
	"time"

	"server-events/internal/services/usersvc/models"
	pbapiv1 "server-events/pkg/genproto/pb"

	"github.com/go-kit/log"
)

type (
	UserSvc interface {
		CreateUser(ctx context.Context, user *models.User) (bool, error)
		GetUsers(ctx context.Context, req *pbapiv1.GetUsersRequest, stream pbapiv1.UserSvc_GetUsersServer) error
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
  fmt.Println("test)
	return true, nil

}

// func (s *userSvc) GetUsers(userID string, stream chan<- string) error {
// 	for i := 0; i < 5; i++ {
// 		stream <- "Message " + userID
// 		time.Sleep(1 * time.Second)
// 	}

// 	close(stream)
// 	return nil
// }

func (s *userSvc) GetUsers(ctx context.Context, req *pbapiv1.GetUsersRequest, stream pbapiv1.UserSvc_GetUsersServer) error {
	fmt.Printf("Received streaming request with input: %s\n", req.UserId)

	// Simulate streaming data
	for i := 0; i < 5; i++ {
		select {
		case <-ctx.Done():
			fmt.Println("Context cancelled, stopping stream.")
			return ctx.Err()
		default:
			chunk := fmt.Sprintf("Chunk %d for %s", i+1, req.UserId)
			resp := &pbapiv1.GetUsersResponse{Message: chunk}
			if err := stream.Send(resp); err != nil {
				fmt.Printf("Failed to send stream chunk: %v\n", err)
				return err
			}
			fmt.Printf("Sent chunk: %s\n", chunk)
			time.Sleep(500 * time.Millisecond) // Simulate work
		}
	}
	return nil
}
