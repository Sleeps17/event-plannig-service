package grpc

import (
	"context"
	authservice "github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/services/auth"
	authv1 "github.com/Sleeps17/events-planning-service-backend/api_gateway/protos/gen/go/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	notAdmin = false
)

func (s *Service) IsAdmin(ctx context.Context, userID uint64, appID uint32) (isAdmin bool, err error) {
	resp, err := s.client.IsAdmin(ctx, &authv1.IsAdminRequest{UserId: userID, AppId: appID})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.Internal:
				return notAdmin, authservice.ErrInternalServer
			case codes.InvalidArgument:
				return notAdmin, authservice.ErrInvalidCredentials
			case codes.NotFound:
				return false, nil
			}
		} else {
			return notAdmin, authservice.ErrInternalServer
		}
	}

	return resp.IsAdmin, nil
}
