package grpc

import (
	"context"
	authservice "github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/services/auth"
	employeeservice "github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/services/employees"
	employeesv1 "github.com/Sleeps17/events-planning-service-backend/api_gateway/protos/gen/go/employees"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) Delete(ctx context.Context, employeeID uint64) (err error) {
	if _, err := s.client.Delete(
		ctx,
		&employeesv1.DeleteRequest{Id: emptyID},
	); err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.Internal:
				return employeeservice.ErrInternalServer
			case codes.NotFound:
				return employeeservice.ErrEmployeeNotFound
			}
		} else {
			return authservice.ErrInternalServer
		}
	}

	return nil
}
