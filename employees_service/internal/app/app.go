package app

import (
	"context"
	"errors"
	"github.com/Sleeps17/events-planning-backend-service/employees_service/internal/config"
	"github.com/Sleeps17/events-planning-backend-service/employees_service/internal/database"
	grpcserver "github.com/Sleeps17/events-planning-backend-service/employees_service/internal/grpc"
	employeesrepository "github.com/Sleeps17/events-planning-backend-service/employees_service/internal/repository/employee"
	employeeservice "github.com/Sleeps17/events-planning-backend-service/employees_service/internal/services/employees/repository"
	"google.golang.org/grpc"
	"log/slog"
)

type App struct {
	server *grpcserver.Server
}

func MustNew(logger *slog.Logger, cfg *config.Config) *App {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Storage.Timeout)
	defer cancel()
	db := database.MustInit(ctx, &cfg.Storage)
	defer db.Close()
	logger.Info("database was successfully inited")

	employeesRepo := employeesrepository.New(db.Pool)
	employeesProvider := employeeservice.New(employeesRepo)

	grpcsrv := grpcserver.New(&cfg.Server, logger, employeesProvider)

	return &App{server: grpcsrv}
}

func (a *App) MustRun() {
	if err := a.server.Run(); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		panic(err)
	}
}

func (a *App) ShutDown() {
	a.server.Stop()
}
