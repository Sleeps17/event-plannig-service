package employee

import (
	"context"
	"fmt"
	repo "github.com/Sleeps17/events-planning-service-backend/employees_service/internal/repository"
)

func (r *repository) Delete(ctx context.Context, employeeID uint64) error {
	const op = "employee-repository.Delete"

	query := `DELETE FROM employees_schema.employees WHERE id = ?`

	cmd, err := r.pool.Exec(ctx, query, employeeID)
	if err != nil {
		if cmd.RowsAffected() == 0 {
			return repo.ErrEmployeeNotFound
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
