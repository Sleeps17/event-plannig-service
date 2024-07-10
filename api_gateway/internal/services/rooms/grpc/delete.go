package grpc

import (
	"context"
	roomservice "github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/services/rooms"
	roomsv1 "github.com/Sleeps17/events-planning-service-backend/api_gateway/protos/gen/go/rooms"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) Delete(ctx context.Context, roomID uint32) error {
	if _, err := s.client.Delete(
		ctx,
		&roomsv1.DeleteRequest{Id: roomID},
	); err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.Internal:
				return roomservice.ErrInternalServer
			case codes.NotFound:
				return roomservice.ErrRoomNotFound
			}
		} else {
			return roomservice.ErrInternalServer
		}
	}

	return nil
}
