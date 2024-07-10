package grpc

import (
	"context"
	"github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/domain/models"
	roomservice "github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/services/rooms"
	roomsv1 "github.com/Sleeps17/events-planning-service-backend/api_gateway/protos/gen/go/rooms"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) GetByID(ctx context.Context, roomID uint32) (*models.Room, error) {
	resp, err := s.client.GetByID(ctx, &roomsv1.GetRequest{Id: roomID})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.NotFound:
				return nil, roomservice.ErrRoomNotFound
			case codes.Internal:
				return nil, roomservice.ErrInternalServer
			}
		} else {
			return nil, roomservice.ErrInternalServer
		}
	}

	return &models.Room{
		ID:       resp.Room.Id,
		Name:     resp.Room.Name,
		Capacity: resp.Room.Capacity,
	}, nil
}
