package employees_repository_repository

import (
	"context"
	"errors"
	"github.com/Sleeps17/events-planning-backend-service/employees_service/internal/domain/models"
	repo "github.com/Sleeps17/events-planning-backend-service/employees_service/internal/repository"
	employeeservice "github.com/Sleeps17/events-planning-backend-service/employees_service/internal/services/employees"
)

const (
	defaultIDValue = 0
)

func (s *service) Create(ctx context.Context, employee *models.Employee) (uint64, error) {
	id, err := s.repo.Save(ctx, employee)
	if err != nil {
		if errors.Is(err, repo.ErrEmployeeAlreadyExists) {
			return defaultIDValue, employeeservice.ErrEmployeeAlreadyExists
		}

		return defaultIDValue, err
	}

	return id, nil
}
