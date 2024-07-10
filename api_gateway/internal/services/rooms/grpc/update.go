package grpc

import (
	"context"
	"github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/domain/models"
	roomservice "github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/services/rooms"
	roomsv1 "github.com/Sleeps17/events-planning-service-backend/api_gateway/protos/gen/go/rooms"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) Update(ctx context.Context, room *models.Room) (*models.Room, error) {
	resp, err := s.client.Update(ctx, &roomsv1.UpdateRequest{
		Id: room.ID,
		UpdatedRoom: &roomsv1.Room{
			Id:       room.ID,
			Name:     room.Name,
			Capacity: room.Capacity,
		},
	})

	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.NotFound:
				return nil, roomservice.ErrRoomNotFound
			case codes.AlreadyExists:
				return nil, roomservice.ErrRoomAlreadyExists
			case codes.Internal:
				return nil, roomservice.ErrInternalServer
			}
		} else {
			return nil, roomservice.ErrInternalServer
		}
	}

	return &models.Room{
		ID:       resp.UpdatedRoom.Id,
		Name:     resp.UpdatedRoom.Name,
		Capacity: resp.UpdatedRoom.Capacity,
	}, nil
}
