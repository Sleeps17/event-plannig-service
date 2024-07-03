package mock_rooms_service

import (
	"context"
	"github.com/Sleeps17/event-plannig-service-backend/events-service/internal/domain/models"
)

type Service struct{}

func New() *Service {
	return &Service{}
}

func (s *Service) GetByID(ctx context.Context, roomID uint32) (*models.Room, error) {
	return &models.Room{
		ID:       roomID,
		Name:     "Mock room",
		Capacity: 20,
	}, nil
}
