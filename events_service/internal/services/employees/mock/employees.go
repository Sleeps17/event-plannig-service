package mock_employees_service

import (
	"context"
	"fmt"
	"github.com/Sleeps17/events-planning-service-backend/events_service/internal/domain/models"
)

type Service struct{}

func New() *Service {
	return &Service{}
}

func (s *Service) GetByID(ctx context.Context, employeeID uint64) (*models.Employee, error) {
	return &models.Employee{
		ID:           employeeID,
		FirstName:    fmt.Sprintf("TestFirstName%d", employeeID),
		LastName:     fmt.Sprintf("TestLastName%d", employeeID),
		Patronymic:   fmt.Sprintf("TestPatronymic%d", employeeID),
		Email:        fmt.Sprintf("test%d@gmail.com", employeeID),
		MobileNumber: "+79136442470",
	}, nil
}

func (s *Service) GetByIDs(ctx context.Context, employeesIDs []uint64) ([]*models.Employee, error) {
	employees := make([]*models.Employee, len(employeesIDs))
	for i, employeeID := range employeesIDs {
		employees[i] = &models.Employee{
			ID:           employeeID,
			FirstName:    fmt.Sprintf("TestFirstName%d", employeeID),
			LastName:     fmt.Sprintf("TestLastName%d", employeeID),
			Patronymic:   fmt.Sprintf("TestPatronymic%d", employeeID),
			Email:        fmt.Sprintf("test%d@gmail.com", employeeID),
			MobileNumber: "+79136442470",
		}
	}
	return employees, nil
}
