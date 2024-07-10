package grpc

import (
	"context"
	authservice "github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/services/auth"
	authv1 "github.com/Sleeps17/events-planning-service-backend/api_gateway/protos/gen/go/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) Register(ctx context.Context, login, password string) error {
	_, err := s.client.Register(ctx, &authv1.RegisterRequest{
		Login:    login,
		Password: password,
	})

	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.AlreadyExists:
				return authservice.ErrUserAlreadyExists
			case codes.Internal:
				return authservice.ErrInternalServer
			case codes.InvalidArgument:
				return authservice.ErrInvalidCredentials
			}
		} else {
			return authservice.ErrInternalServer
		}
	}

	return nil
}
