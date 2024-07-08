package employees_repository_repository

import (
	"context"
	"errors"
	repo "github.com/Sleeps17/events-planning-backend-service/employees_service/internal/repository"
	employeeservice "github.com/Sleeps17/events-planning-backend-service/employees_service/internal/services/employees"
)

func (s *service) Delete(ctx context.Context, employeeID uint64) error {
	if err := s.repo.Delete(ctx, employeeID); err != nil {
		if errors.Is(err, repo.ErrEmployeeNotFound) {
			return employeeservice.ErrEmployeeNotFound
		}

		return err
	}

	return nil
}
