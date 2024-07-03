package employees_service

import (
	"context"
	"github.com/Sleeps17/event-plannig-service-backend/events-service/internal/domain/models"
)

type Service interface {
	GetByID(ctx context.Context, employeeID uint64) (*models.Employee, error)
	GetByIDs(ctx context.Context, employeesIDs []uint64) ([]*models.Employee, error)
}
