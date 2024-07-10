package grpc

import (
	"context"
	"github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/domain/models"
	authservice "github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/services/auth"
	employeeservice "github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/services/employees"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Service) GetAll(ctx context.Context) ([]*models.Employee, error) {
	resp, err := s.client.GetAll(ctx, &emptypb.Empty{})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.Internal:
				return nil, employeeservice.ErrInternalServer
			}
		} else {
			return nil, authservice.ErrInternalServer
		}
	}

	var employees []*models.Employee
	for _, e := range resp.Employees {
		employees = append(employees, &models.Employee{
			ID:           e.Id,
			FirstName:    e.FirstName,
			LastName:     e.LastName,
			Patronymic:   e.Patronymic,
			Email:        e.Email,
			MobileNumber: e.MobileNumber,
		})
	}

	return employees, nil
}
