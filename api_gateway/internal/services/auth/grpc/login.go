package grpc

import (
	"context"
	authservice "github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/services/auth"
	authv1 "github.com/Sleeps17/events-planning-service-backend/api_gateway/protos/gen/go/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	emptyToken = ""
)

func (s *Service) Login(ctx context.Context, login, password string, appID uint32) (token string, err error) {
	resp, err := s.client.Login(ctx, &authv1.LoginRequest{
		Login:    login,
		Password: password,
		AppId:    appID,
	})

	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.Internal:
				return emptyToken, authservice.ErrInternalServer
			case codes.InvalidArgument:
				return emptyToken, authservice.ErrInvalidCredentials
			}
		} else {
			return emptyToken, authservice.ErrInternalServer
		}
	}

	return resp.Token, nil
}
