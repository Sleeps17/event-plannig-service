package repository

import (
	"context"
	"errors"
	"github.com/Sleeps17/events-planning-service-backend/rooms_service/internal/domain/models"
)

var (
	ErrRoomExists   = errors.New("room already exists")
	ErrRoomNotFound = errors.New("room not found")
)

type Repository interface {
	Insert(ctx context.Context, room *models.Room) (uint64, error)
	SelectByID(ctx context.Context, id uint64) (*models.Room, error)
	SelectAll(ctx context.Context) ([]*models.Room, error)
	Update(ctx context.Context, room *models.Room) error
	Delete(ctx context.Context, id uint64) error
}
