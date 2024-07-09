package events_repository

import (
	"context"
	"errors"
	"github.com/Sleeps17/events-planning-service-backend/events_service/internal/domain/models"
	"time"
)

var (
	ErrEventNotFound = errors.New("event not found")
)

type Repository interface {
	SelectByID(ctx context.Context, id uint64) (*models.Event, error)
	SelectAllBetweenTwoDates(ctx context.Context, startDate, endDate time.Time) ([]*models.Event, error)
	CheckRoomIsAvailable(ctx context.Context, roomID uint32, startDate, endDate time.Time) (bool, error)
	CheckEmployeesAreAvailable(ctx context.Context, employeesIDs []uint64, startDate, endDate time.Time) ([]uint64, error)
	Insert(ctx context.Context, event *models.Event) (id uint64, err error)
	UpdateByID(ctx context.Context, updatedEvent *models.Event) (err error)
	DeleteByID(ctx context.Context, id uint64) (err error)
}
