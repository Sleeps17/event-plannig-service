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

func (s *Service) Update(ctx context.Context, updated *models.Employee) (employee *models.Employee, err error) {
	upd, err := s.client.Update(ctx, &employeesv1.UpdateRequest{
		Id: updated.ID,
		UpdatedEmployee: &employeesv1.Employee{
			Id:           updated.ID,
			FirstName:    updated.FirstName,
			LastName:     updated.LastName,
			Patronymic:   updated.Patronymic,
			Email:        updated.Email,
			MobileNumber: updated.MobileNumber,
		},
	})

	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.Internal:
				return nil, employeeservice.ErrInternalServer
			case codes.NotFound:
				return nil, employeeservice.ErrEmployeeNotFound
			case codes.AlreadyExists:
				return nil, employeeservice.ErrEmployeeAlreadyExists
			}
		} else {
			return nil, authservice.ErrInternalServer
		}
	}

	return &models.Employee{
		ID:           upd.UpdatedEmployee.Id,
		FirstName:    upd.UpdatedEmployee.FirstName,
		LastName:     upd.UpdatedEmployee.LastName,
		Patronymic:   upd.UpdatedEmployee.Patronymic,
		Email:        upd.UpdatedEmployee.Email,
		MobileNumber: upd.UpdatedEmployee.MobileNumber,
	}, nil
}
