package rooms_service

import (
	"context"
	"github.com/Sleeps17/event-plannig-service-backend/events-service/internal/domain/models"
)

type Service interface {
	GetByID(ctx context.Context, roomID uint32) (*models.Room, error)
}
