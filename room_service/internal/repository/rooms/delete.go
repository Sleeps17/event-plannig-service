package rooms

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
)

func (r *repository) Delete(ctx context.Context, id uint64) error {
	const op = "room-repository.Delete"

	query := `DELETE FROM rooms_schema.rooms WHERE id = $1`

	if _, err := r.pool.Exec(ctx, query, id); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return fmt.Errorf("%s: database error: %w", op, pgErr)
		}

		return fmt.Errorf("%s: can't delete room: %w", op, err)
	}

	return nil
}
