package employees_repository_repository

import (
	"context"
	"errors"
	"github.com/Sleeps17/events-planning-service-backend/employees_service/internal/domain/models"
	"github.com/Sleeps17/events-planning-service-backend/employees_service/internal/repository"
	employeeservice "github.com/Sleeps17/events-planning-service-backend/employees_service/internal/services/employees"
)

func (s *service) GetByID(ctx context.Context, employeeID uint64) (*models.Employee, error) {
	employee, err := s.repo.SelectByID(ctx, employeeID)
	if err != nil {
		if errors.Is(err, repository.ErrEmployeeNotFound) {
			return nil, employeeservice.ErrEmployeeNotFound
		}

		return nil, err
	}

	return employee, nil
}
