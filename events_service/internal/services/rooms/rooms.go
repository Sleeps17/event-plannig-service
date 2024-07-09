package rooms_service

import (
	"context"
	"github.com/Sleeps17/events-planning-service-backend/events_service/internal/domain/models"
)

type Service interface {
	GetByID(ctx context.Context, roomID uint32) (*models.Room, error)
}
