package grpc

import (
	"context"
	eventservice "github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/services/events"
	eventsv1 "github.com/Sleeps17/events-planning-service-backend/api_gateway/protos/gen/go/events"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) Delete(ctx context.Context, id uint64) error {
	if _, err := s.client.Delete(
		ctx,
		&eventsv1.DeleteRequest{Id: id},
	); err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.NotFound:
				return eventservice.ErrEventNotFound
			case codes.Internal:
				return eventservice.ErrInternalServer
			}
		} else {
			return eventservice.ErrInternalServer
		}
	}

	return nil
}
