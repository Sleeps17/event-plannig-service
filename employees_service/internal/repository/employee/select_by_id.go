package employee

import (
	"context"
	"errors"
	"fmt"
	"github.com/Sleeps17/events-planning-service-backend/employees_service/internal/domain/models"
	repo "github.com/Sleeps17/events-planning-service-backend/employees_service/internal/repository"
)

func (r *repository) SelectByID(ctx context.Context, employeeID uint64) (*models.Employee, error) {
	const op = "employee-repository.SelectByID"

	query := `SELECT * FROM employee_schema.employees WHERE id = $1`

	var employee models.Employee
	if err := r.pool.QueryRow(ctx, query, employeeID).Scan(
		&employee.ID,
		&employee.FirstName,
		&employee.LastName,
		&employee.Patronymic,
		&employee.Email,
		&employee.MobileNumber,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repo.ErrEmployeeNotFound
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &employee, nil
}
