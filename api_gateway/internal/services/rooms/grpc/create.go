package grpc

import (
	"context"
	"github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/domain/models"
	roomservice "github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/services/rooms"
	roomsv1 "github.com/Sleeps17/events-planning-service-backend/api_gateway/protos/gen/go/rooms"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	emptyID = 0
)

func (s *Service) Create(ctx context.Context, room *models.Room) (id uint32, err error) {
	resp, err := s.client.Create(ctx, &roomsv1.CreateRequest{
		Room: &roomsv1.Room{
			Name:     room.Name,
			Capacity: room.Capacity,
		},
	})

	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.AlreadyExists:
				return emptyID, roomservice.ErrRoomAlreadyExists
			case codes.Internal:
				return emptyID, roomservice.ErrInternalServer
			}
		} else {
			return emptyID, roomservice.ErrInternalServer
		}
	}

	return resp.Id, nil
}
