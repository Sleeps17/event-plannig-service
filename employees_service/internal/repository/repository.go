package repository

import (
	"context"
	"errors"
	"github.com/Sleeps17/events-planning-backend-service/employees_service/internal/domain/models"
)

var (
	ErrEmployeeAlreadyExists = errors.New("employee already exists")
	ErrEmployeeNotFound      = errors.New("employee not found")
)

type EmployeeRepository interface {
	Save(ctx context.Context, employee *models.Employee) (uint64, error)
	SelectByID(ctx context.Context, employeeID uint64) (*models.Employee, error)
	SelectAll(ctx context.Context) ([]*models.Employee, error)
	Update(ctx context.Context, updatedEmployee *models.Employee) error
	Delete(ctx context.Context, employeeID uint64) error
}
