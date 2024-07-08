package employee

import (
	"context"
	"errors"
	"fmt"
	"github.com/Sleeps17/events-planning-backend-service/employees_service/internal/domain/models"
	repo "github.com/Sleeps17/events-planning-backend-service/employees_service/internal/repository"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	uniqueViolationCode = "23505"

	defaultIDValue = 0
)

func (r *repository) Save(ctx context.Context, employee *models.Employee) (uint64, error) {
	const op = "employee-repository.Save"

	query := `INSERT INTO employees_schema.employees(first_name, last_name, patronymic, email, mobile_numner) 
				VALUES($1, $2, $3, $4, $5) RETURNING id;`

	var id uint64
	if err := r.pool.QueryRow(
		ctx,
		query,
		employee.FirstName,
		employee.LastName,
		employee.Patronymic,
		employee.Email,
		employee.MobileNumber,
	).Scan(&id); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == uniqueViolationCode {
			return defaultIDValue, repo.ErrEmployeeAlreadyExists
		}

		return defaultIDValue, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}
