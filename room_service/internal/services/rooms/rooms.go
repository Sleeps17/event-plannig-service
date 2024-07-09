package rooms

import (
	"context"
	"errors"
	"github.com/Sleeps17/events-planning-service-backend/rooms_service/internal/domain/models"
)

var (
	ErrRoomExists   = errors.New("room already exists")
	ErrRoomNotFound = errors.New("room not found")
)

type Service interface {
	Add(ctx context.Context, room *models.Room) (uint64, error)
	GetByID(ctx context.Context, id uint64) (*models.Room, error)
	GetAll(ctx context.Context) ([]*models.Room, error)
	Update(ctx context.Context, updatedRoom *models.Room) error
	DeleteByID(ctx context.Context, id uint64) error
}
