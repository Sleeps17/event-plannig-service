package events

import (
	"context"
	"errors"
	"github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/domain/models"
	"time"
)

var (
	ErrRoomIsOccupied     = errors.New("room is occupied")
	ErrSomeWorkersAreBusy = errors.New("workers are busy")
	ErrEventNotFound      = errors.New("event not found")
	ErrInternalServer     = errors.New("internal server error")
)

type CreateResponse struct {
	ID            uint64
	BusyEmployees []*models.Employee
}

type UpdateResponse struct {
	UpdatedEvent  *models.Event
	BusyEmployees []*models.Employee
}

type Service interface {
	Create(ctx context.Context, event *models.Event) (resp *CreateResponse, err error)
	GetByID(ctx context.Context, id uint64) (event *models.Event, err error)
	GetAll(ctx context.Context, startDate, andDate time.Time) (events []*models.Event, err error)
	Update(ctx context.Context, event *models.Event) (*UpdateResponse, error)
	Delete(ctx context.Context, id uint64) (err error)
}
