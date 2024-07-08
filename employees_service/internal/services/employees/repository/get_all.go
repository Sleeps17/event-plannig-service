package employees_repository_repository

import (
	"context"
	"github.com/Sleeps17/events-planning-backend-service/employees_service/internal/domain/models"
)

func (s *service) GetAll(ctx context.Context) ([]*models.Employee, error) {
	employees, err := s.repo.SelectAll(ctx)
	if err != nil {
		return nil, err
	}

	return employees, nil
}
