package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/Sleeps17/events-planning-service-backend/auth_service/internal/domain/models"
	_ "github.com/Sleeps17/events-planning-service-backend/auth_service/internal/domain/models"
	"github.com/Sleeps17/events-planning-service-backend/auth_service/internal/jwt"
	repo "github.com/Sleeps17/events-planning-service-backend/auth_service/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidAppID       = errors.New("invalid app id")
	ErrUserExists         = errors.New("user already exists")
)

type Auth struct {
	log         *slog.Logger
	usrProvider repo.UserRepository
	appProvider repo.AppRepository
	tokenTTL    time.Duration
}

func New(
	log *slog.Logger,
	userProvider repo.UserRepository,
	appProvider repo.AppRepository,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		log:         log,
		usrProvider: userProvider,
		appProvider: appProvider,
		tokenTTL:    tokenTTL,
	}
}

func (a *Auth) Login(ctx context.Context, login string, password string, appId uint32) (string, error) {
	a.log.Info("try log in user with login: " + login)

	user, err := a.usrProvider.SelectByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, repo.ErrUserNotFound) {
			a.log.Warn("user not found", slog.String("error", err.Error()))
			return "", fmt.Errorf("can't login user: %w", ErrInvalidCredentials)
		}

		a.log.Error("can't get user: %w", err.Error())

		return "", fmt.Errorf("can't get user: %w", err)
	}
	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.log.Warn("invalid credentials", slog.String("error", err.Error()))
		return "", fmt.Errorf("can't login user: %w", ErrInvalidCredentials)
	}

	app, err := a.appProvider.SelectByID(ctx, appId)
	if err != nil {
		if errors.Is(err, repo.ErrAppNotFound) {
			a.log.Warn("app not found", slog.String("error", err.Error()))
			return "", ErrInvalidAppID
		}
		a.log.Error("can't provide app", slog.String("error", err.Error()))
		return "", fmt.Errorf("can't provide app: %w", err)
	}

	token, err := jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		a.log.Error("can't generate jwt-token", slog.String("error", err.Error()))
		return "", fmt.Errorf("can't generate jwt-token: %w", err)
	}

	a.log.Info("user log in", slog.String("login", login))

	return token, nil
}

func (a *Auth) RegisterNewUser(ctx context.Context, login string, password string) (uint64, error) {

	a.log.Info("try to register user", slog.String("login", login))

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		a.log.Error("can't generate password hash", slog.String("error", err.Error()))
		return 0, fmt.Errorf("can't generate password hash: %d", err)
	}

	id, err := a.usrProvider.Save(ctx, &models.User{
		Login:    login,
		PassHash: passHash,
	})
	if err != nil {
		if errors.Is(err, repo.ErrUserExists) {
			a.log.Warn("user already exists", slog.String("login", login))
			return 0, ErrUserExists
		}
		a.log.Error("can't save new user", slog.String("login", login), slog.String("error", err.Error()))
		return 0, fmt.Errorf("can't save new user: %w", err)
	}

	a.log.Info("user registered", slog.String("login", login))

	return id, nil
}

func (a *Auth) IsAdmin(ctx context.Context, userID uint64) (bool, error) {
	a.log.Info("try to check if the user is admin")

	isAdmin, err := a.usrProvider.IsAdmin(ctx, userID)
	if err != nil {
		if errors.Is(err, repo.ErrUserNotFound) {
			a.log.Warn("user not found", slog.Int("userID", int(userID)))
			return false, ErrInvalidCredentials
		}
		a.log.Error("can't check if user is admin", slog.String("error", err.Error()))
		return false, fmt.Errorf("can't check if user is admin: %w", err)
	}

	a.log.Info("checked if user is admin", slog.Bool("Is_admin", isAdmin), slog.Int64("userID", int64(userID)))

	return isAdmin, nil
}
