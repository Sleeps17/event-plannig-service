package grpc

import (
	"context"
	"github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/domain/models"
	roomservice "github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/services/rooms"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Service) GetAll(ctx context.Context) ([]*models.Room, error) {
	resp, err := s.client.GetAll(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, roomservice.ErrInternalServer
	}

	var rooms []*models.Room
	for _, room := range resp.Rooms {
		rooms = append(rooms, &models.Room{
			ID:       room.Id,
			Name:     room.Name,
			Capacity: room.Capacity,
		})
	}

	return rooms, nil
}
