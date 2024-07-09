package api

import (
	"github.com/Sleeps17/events-planning-service-backend/employees_service/internal/domain/models"
	employeesv1 "github.com/Sleeps17/events-planning-service-backend/employees_service/protos/gen/go/employees"
)

func fromRequest(employee *employeesv1.Employee) *models.Employee {
	return &models.Employee{
		ID:           employee.Id,
		FirstName:    employee.FirstName,
		LastName:     employee.LastName,
		Patronymic:   employee.Patronymic,
		Email:        employee.Email,
		MobileNumber: employee.MobileNumber,
	}
}

func toResponse(employee *models.Employee) *employeesv1.Employee {
	return &employeesv1.Employee{
		Id:           employee.ID,
		FirstName:    employee.FirstName,
		LastName:     employee.LastName,
		Patronymic:   employee.Patronymic,
		Email:        employee.Email,
		MobileNumber: employee.MobileNumber,
	}
}
