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

const (
	emptyID = 0
)

func (s *Service) Create(ctx context.Context, employee *models.Employee) (employeeID uint64, err error) {
	resp, err := s.client.Create(ctx, &employeesv1.CreateRequest{
		Employee: &employeesv1.Employee{
			Id:           employee.ID,
			FirstName:    employee.FirstName,
			LastName:     employee.LastName,
			Patronymic:   employee.Patronymic,
			Email:        employee.Email,
			MobileNumber: employee.MobileNumber,
		},
	})

	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.Internal:
				return emptyID, employeeservice.ErrInternalServer
			case codes.AlreadyExists:
				return emptyID, employeeservice.ErrEmployeeAlreadyExists
			}
		} else {
			return emptyID, authservice.ErrInternalServer
		}
	}

	return resp.Id, nil
}
