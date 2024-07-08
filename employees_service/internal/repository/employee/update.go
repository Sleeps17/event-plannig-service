package employee

import (
	"context"
	"errors"
	"fmt"
	"github.com/Sleeps17/events-planning-backend-service/employees_service/internal/domain/models"
	repo "github.com/Sleeps17/events-planning-backend-service/employees_service/internal/repository"
	"github.com/jackc/pgx/v5/pgconn"
)

func (r *repository) Update(ctx context.Context, updatedEmployee *models.Employee) error {
	const op = "employee-repository.Update"

	query := `
        UPDATE employees_schema.employees
        SET first_name = $1, last_name = $2, patronymic = $3, email = $4, mobile_number = $5
        WHERE id = $6
    `
	cmd, err := r.pool.Exec(ctx, query,
		updatedEmployee.FirstName,
		updatedEmployee.LastName,
		updatedEmployee.Patronymic,
		updatedEmployee.Email,
		updatedEmployee.MobileNumber,
		updatedEmployee.ID,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == uniqueViolationCode {
			return repo.ErrEmployeeAlreadyExists
		}

		if cmd.RowsAffected() == 0 {
			return repo.ErrEmployeeNotFound
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
