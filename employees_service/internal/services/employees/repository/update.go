package employees_repository_repository

import (
	"context"
	"errors"
	"github.com/Sleeps17/events-planning-backend-service/employees_service/internal/domain/models"
	"github.com/Sleeps17/events-planning-backend-service/employees_service/internal/repository"
	employeeservice "github.com/Sleeps17/events-planning-backend-service/employees_service/internal/services/employees"
)

func (s *service) Update(ctx context.Context, employee *models.Employee) (*models.Employee, error) {
	if err := s.repo.Update(ctx, employee); err != nil {
		if errors.Is(err, repository.ErrEmployeeAlreadyExists) {
			return nil, employeeservice.ErrEmployeeAlreadyExists
		}

		if errors.Is(err, repository.ErrEmployeeNotFound) {
			return nil, employeeservice.ErrEmployeeNotFound
		}

		return nil, err
	}

	updatedEmployee, err := s.repo.SelectByID(ctx, employee.ID)
	if err != nil {
		return nil, err
	}

	return updatedEmployee, nil
}
