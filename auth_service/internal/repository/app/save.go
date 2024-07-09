package app

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

func (r *repository) Save(ctx context.Context, app *models.App) (uint32, error) {
	const op = "app-repository.Save"

	query := `INSERT INTO apps_schema.apps(name, secret) VALUES ($1, $2) RETURNING id`

	var appID uint32
	if err := r.pool.QueryRow(ctx, query, app.Name, app.Secret).Scan(&appID); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == uniqueViolationCode {
			return defaultIDValue, repo.ErrAppExists
		}

		return appID, fmt.Errorf("%s: %w", op, err)
	}

	return appID, nil
}
