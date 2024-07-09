package employee

import (
	"context"
	"errors"
	"fmt"
	"github.com/Sleeps17/events-planning-service-backend/employees_service/internal/domain/models"
)

func (r *repository) SelectAll(ctx context.Context) ([]*models.Employee, error) {
	const op = "employee-repository.SelectAll"

	query := `SELECT * FROM employee_schema.employees`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []*models.Employee{}, nil
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var employees []*models.Employee
	for rows.Next() {
		var employee models.Employee
		if err := rows.Scan(
			&employee.ID,
			&employee.FirstName,
			&employee.LastName,
			&employee.Patronymic,
			&employee.Email,
			&employee.MobileNumber,
		); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		employees = append(employees, &employee)
	}

	return employees, nil
}
