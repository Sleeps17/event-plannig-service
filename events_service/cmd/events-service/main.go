package main

import (
	"context"
	"fmt"
	"github.com/Sleeps17/events-planning-service-backend/events_service/internal/config"
	"github.com/Sleeps17/events-planning-service-backend/events_service/internal/database"
	grpcapp "github.com/Sleeps17/events-planning-service-backend/events_service/internal/grpc"
	eventsrepository "github.com/Sleeps17/events-planning-service-backend/events_service/internal/repository/events/postgresql"
	employeeservice "github.com/Sleeps17/events-planning-service-backend/events_service/internal/services/employees/mock"
	eventservice "github.com/Sleeps17/events-planning-service-backend/events_service/internal/services/events/repository"
	roomservice "github.com/Sleeps17/events-planning-service-backend/events_service/internal/services/rooms/mock"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoad()

	fmt.Println(cfg)

	logger := setupLogger(cfg.Env)
	logger.Info("logger was successfully setup", slog.String("env", cfg.Env))

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Storage.Timeout)
	defer cancel()
	db := database.MustInit(ctx, &cfg.Storage)
	defer db.Close()
	logger.Info("database was successfully inited")

	eventsRepo := eventsrepository.New(db.Pool)
	eventProvider := eventservice.New(eventsRepo)
	roomsProvider := roomservice.New()
	employeesProvider := employeeservice.New()

	app := grpcapp.New(
		&cfg.Server,
		logger,
		eventProvider,
		roomsProvider,
		employeesProvider,
	)
	go app.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	sig := <-stop

	logger.Info("application stopping", slog.String("signal", sig.String()))
	app.Stop()
}

const (
	localEnv = "local"
	devEnv   = "dev"
	prodEnv  = "prod"
)

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case localEnv:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case devEnv:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	case prodEnv:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn}),
		)
	default:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}),
		)
	}

	return log
}
