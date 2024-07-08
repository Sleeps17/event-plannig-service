package api

import (
	"context"
	"errors"
	employeeservice "github.com/Sleeps17/events-planning-backend-service/employees_service/internal/services/employees"
	employeesv1 "github.com/Sleeps17/events-planning-backend-service/employees_service/protos/gen/go/employees"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"log/slog"
)

type api struct {
	employeesv1.UnimplementedEmployeesServer
	logger           *slog.Logger
	employeeProvider employeeservice.Service
}

func Register(srv *grpc.Server, logger *slog.Logger, employeeProvider employeeservice.Service) {
	employeesv1.RegisterEmployeesServer(srv, &api{
		logger:           logger,
		employeeProvider: employeeProvider,
	})
}

func (a *api) Create(ctx context.Context, request *employeesv1.CreateRequest) (*employeesv1.CreateResponse, error) {
	employee := fromRequest(request.Employee)

	a.logger.Info("try to handle create request", slog.Any("employeeID", employee.ID))

	id, err := a.employeeProvider.Create(ctx, employee)
	if err != nil {
		if errors.Is(err, employeeservice.ErrEmployeeAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, EmployeeAlreadyExistsMsg)
		}

		return nil, status.Error(codes.Internal, InternalErrorMsg)
	}

	return &employeesv1.CreateResponse{Id: id}, nil
}

func (a *api) GetByID(ctx context.Context, request *employeesv1.GetByIDRequest) (*employeesv1.GetByIDResponse, error) {
	id := request.GetId()

	employee, err := a.employeeProvider.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, employeeservice.ErrEmployeeNotFound) {
			return nil, status.Error(codes.NotFound, EmployeeNotFoundMsg)
		}

		return nil, status.Error(codes.Internal, InternalErrorMsg)
	}

	return &employeesv1.GetByIDResponse{Employee: toResponse(employee)}, nil
}

func (a *api) GetAll(ctx context.Context, _ *emptypb.Empty) (*employeesv1.GetAllResponse, error) {
	employees, err := a.employeeProvider.GetAll(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, InternalErrorMsg)
	}

	var transformedEmployees []*employeesv1.Employee
	for _, employee := range employees {
		transformedEmployees = append(transformedEmployees, toResponse(employee))
	}

	return &employeesv1.GetAllResponse{Employees: transformedEmployees}, nil
}

func (a *api) Update(ctx context.Context, request *employeesv1.UpdateRequest) (*employeesv1.UpdateResponse, error) {
	employee := fromRequest(request.UpdatedEmployee)
	employee.ID = request.GetId()

	updatedEmployee, err := a.employeeProvider.Update(ctx, employee)
	if err != nil {
		if errors.Is(err, employeeservice.ErrEmployeeNotFound) {
			return nil, status.Error(codes.NotFound, EmployeeNotFoundMsg)
		}

		if errors.Is(err, employeeservice.ErrEmployeeAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, EmployeeAlreadyExistsMsg)
		}

		return nil, status.Error(codes.Internal, InternalErrorMsg)
	}

	return &employeesv1.UpdateResponse{UpdatedEmployee: toResponse(updatedEmployee)}, nil
}

func (a *api) Delete(ctx context.Context, request *employeesv1.DeleteRequest) (*emptypb.Empty, error) {
	id := request.GetId()

	if err := a.employeeProvider.Delete(ctx, id); err != nil {
		if errors.Is(err, employeeservice.ErrEmployeeNotFound) {
			return nil, status.Error(codes.NotFound, EmployeeNotFoundMsg)
		}

		return nil, status.Error(codes.Internal, InternalErrorMsg)
	}

	return &emptypb.Empty{}, nil
}
