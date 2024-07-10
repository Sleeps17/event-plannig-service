package employees

import (
	"context"
	"errors"
	"github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/domain/models"
)

var (
	ErrEmployeeAlreadyExists = errors.New("employee already exists")
	ErrEmployeeNotFound      = errors.New("employee not found")
	ErrInternalServer        = errors.New("internal server error")
)

type Service interface {
	Create(ctx context.Context, employee *models.Employee) (employeeID uint64, err error)
	GetByID(ctx context.Context, employeeID uint64) (employee *models.Employee, err error)
	GetAll(ctx context.Context) (employees []*models.Employee, err error)
	Update(ctx context.Context, updated *models.Employee) (employee *models.Employee, err error)
	Delete(ctx context.Context, employeeID uint64) (err error)
}
