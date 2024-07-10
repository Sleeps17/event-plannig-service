package grpc

import (
	"context"
	"github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/domain/models"
	authservice "github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/services/auth"
	employeeservice "github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/services/employees"
	employeesv1 "github.com/Sleeps17/events-planning-service-backend/api_gateway/protos/gen/go/employees"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) GetByID(ctx context.Context, employeeID uint64) (employee *models.Employee, err error) {
	resp, err := s.client.GetByID(ctx, &employeesv1.GetByIDRequest{
		Id: employeeID,
	})

	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.Internal:
				return nil, employeeservice.ErrInternalServer
			case codes.NotFound:
				return nil, employeeservice.ErrEmployeeNotFound
			}
		} else {
			return nil, authservice.ErrInternalServer
		}
	}

	return &models.Employee{
		ID:           resp.Employee.Id,
		FirstName:    resp.Employee.FirstName,
		LastName:     resp.Employee.LastName,
		Patronymic:   resp.Employee.Patronymic,
		Email:        resp.Employee.Email,
		MobileNumber: resp.Employee.MobileNumber,
	}, nil
}
