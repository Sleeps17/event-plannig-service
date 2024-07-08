package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
)

func (r *repository) IsAdmin(ctx context.Context, userID uint64) (bool, error) {
	const op = "user-repository.IsAdmin"

	query := `SELECT EXISTS(SELECT 1 FROM users_schema.admins WHERE user_id = $1)`

	var isAdmin bool
	if err := r.pool.QueryRow(ctx, query, userID).Scan(&isAdmin); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("%s: can't select admin: %w", op, err)
	}

	return isAdmin, nil
}
