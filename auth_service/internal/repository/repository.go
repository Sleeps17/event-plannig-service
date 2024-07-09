package repository

import (
	"context"
	"errors"
	"github.com/Sleeps17/events-planning-service-backend/auth_service/internal/domain/models"
)

var (
	ErrUserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")
)

type UserRepository interface {
	Save(ctx context.Context, user *models.User) (uint64, error)
	SelectByLogin(ctx context.Context, login string) (*models.User, error)
	IsAdmin(ctx context.Context, userID uint64) (bool, error)
}

var (
	ErrAppExists   = errors.New("app already exists")
	ErrAppNotFound = errors.New("app not found")
)

type AppRepository interface {
	Save(ctx context.Context, app *models.App) (uint32, error)
	SelectByID(ctx context.Context, appID uint32) (*models.App, error)
}
