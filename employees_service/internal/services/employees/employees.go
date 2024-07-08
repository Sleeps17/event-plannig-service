package employees_service

import (
	"context"
	"errors"
	"github.com/Sleeps17/events-planning-backend-service/employees_service/internal/domain/models"
)

var (
	ErrEmployeeAlreadyExists = errors.New("employee already exists")
	ErrEmployeeNotFound      = errors.New("employee not found")
)

type Service interface {
	Create(ctx context.Context, employee *models.Employee) (id uint64, err error)
	GetByID(ctx context.Context, employeeID uint64) (employee *models.Employee, err error)
	GetAll(ctx context.Context) (employees []*models.Employee, err error)
	Update(ctx context.Context, employee *models.Employee) (updatedEmployee *models.Employee, err error)
	Delete(ctx context.Context, employeeID uint64) (err error)
}
