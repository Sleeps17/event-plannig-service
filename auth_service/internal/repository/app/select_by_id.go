package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/Sleeps17/events-planning-service-backend/auth-service/internal/domain/models"
	repo "github.com/Sleeps17/events-planning-service-backend/auth-service/internal/repository"
	"github.com/jackc/pgx/v5"
)

func (r *repository) SelectByID(ctx context.Context, appID uint32) (*models.App, error) {
	const op = "app-repository.SelectByID"

	query := `SELECT id, name, secret FROM apps_schema.apps WHERE id = $1`

	var app models.App
	if err := r.pool.QueryRow(ctx, query, appID).Scan(&app.ID, &app.Name, &app.Secret); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repo.ErrAppNotFound
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &app, nil
}
