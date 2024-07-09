package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/Sleeps17/events-planning-service-backend/auth_service/internal/domain/models"
	repo "github.com/Sleeps17/events-planning-service-backend/auth_service/internal/repository"
)

func (r *repository) SelectByLogin(ctx context.Context, login string) (*models.User, error) {
	const op = "postgresql.User"

	query := `SELECT id, login, pass_hash FROM users_schema.users WHERE login = $1`

	var user models.User
	row := r.pool.QueryRow(ctx, query, login)

	if err := row.Scan(&user.ID, &user.Login, &user.PassHash); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repo.ErrUserNotFound
		}

		return nil, fmt.Errorf("%s: can't select user: %w", op, err)
	}

	return &user, nil
}
