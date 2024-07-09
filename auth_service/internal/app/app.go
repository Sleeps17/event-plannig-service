package app

import (
	"context"
	"errors"
	grpcserver "github.com/Sleeps17/events-planning-service-backend/auth_service/internal/app/grpc"
	"github.com/Sleeps17/events-planning-service-backend/auth_service/internal/config"
	"github.com/Sleeps17/events-planning-service-backend/auth_service/internal/database"
	apprepository "github.com/Sleeps17/events-planning-service-backend/auth_service/internal/repository/app"
	userrepository "github.com/Sleeps17/events-planning-service-backend/auth_service/internal/repository/user"
	"github.com/Sleeps17/events-planning-service-backend/auth_service/internal/services/auth"
	"google.golang.org/grpc"
	"log/slog"
)

type App struct {
	GRPCSrv *grpcserver.Server
}

func MustNew(log *slog.Logger, cfg *config.Config) *App {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Storage.Timeout)
	defer cancel()

	db := database.MustInit(ctx, &cfg.Storage)

	userProvider := userrepository.New(db.Pool)
	appProvider := apprepository.New(db.Pool)

	authService := auth.New(log, userProvider, appProvider, cfg.TokenTTL)
	grpcServer := grpcserver.New(log, authService, cfg.Server.Port)

	return &App{
		GRPCSrv: grpcServer,
	}
}

func (a *App) MustRun() {
	if err := a.GRPCSrv.Run(); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		panic(err)
	}
}
