package events_service

import (
	"context"
	"errors"
	"github.com/Sleeps17/event-plannig-service-backend/events-service/internal/domain/models"
	"time"
)

var (
	ErrRoomIsNotAvailable       = errors.New("room is not available")
	ErrEmployeesAreNotAvailable = errors.New("employees are not available")
	ErrEventNotFound            = errors.New("event not found")
)

type AddResponse struct {
	ID            uint64
	BusyEmployees []uint64
}

type UpdateResponse struct {
	UpdatedEvent  *models.Event
	BusyEmployees []uint64
}

type Service interface {
	GetByID(ctx context.Context, id uint64) (*models.Event, error)
	GetAllOfTheWeek(ctx context.Context, startDate, endDate time.Time) ([]*models.Event, error)
	Add(ctx context.Context, event *models.Event) (resp *AddResponse, err error)
	UpdateByID(ctx context.Context, updatedEvent *models.Event) (resp *UpdateResponse, err error)
	DeleteByID(ctx context.Context, id uint64) (err error)
}
