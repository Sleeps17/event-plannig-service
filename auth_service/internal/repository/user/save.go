package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/Sleeps17/events-planning-service-backend/auth_service/internal/domain/models"
	repo "github.com/Sleeps17/events-planning-service-backend/auth_service/internal/repository"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	uniqueViolationCode = "23505"

	defaultIDValue = 0
)

func (r *repository) Save(ctx context.Context, user *models.User) (uint64, error) {
	const op = "user-repository.Save"

	query := `INSERT INTO users_schema.users(login, pass_hash) VALUES($1, $2) RETURNING id`

	var userId uint64
	if err := r.pool.QueryRow(ctx, query, user.Login, user.PassHash).Scan(&userId); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == uniqueViolationCode {
			return defaultIDValue, repo.ErrUserExists
		}

		return defaultIDValue, fmt.Errorf("%s: can't add user: %w", op, err)
	}

	return userId, nil
}
