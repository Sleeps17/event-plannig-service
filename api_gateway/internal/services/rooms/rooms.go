package rooms

import (
	"context"
	"errors"
	"github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/domain/models"
)

var (
	ErrRoomAlreadyExists = errors.New("room already exists")
	ErrRoomNotFound      = errors.New("room not found")
	ErrInternalServer    = errors.New("internal server error")
)

type Service interface {
	Create(ctx context.Context, room *models.Room) (id uint32, err error)
	GetByID(ctx context.Context, roomID uint32) (*models.Room, error)
	GetAll(ctx context.Context) ([]*models.Room, error)
	Update(ctx context.Context, room *models.Room) (*models.Room, error)
	Delete(ctx context.Context, roomID uint32) error
}
