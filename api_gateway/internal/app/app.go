package app

import (
	"errors"
	"github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/config"
	httpserver "github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/http"
	authhandlers "github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/http/handlers/auth"
	employeeshandlers "github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/http/handlers/employees"
	eventshandlers "github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/http/handlers/events"
	roomshandlers "github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/http/handlers/rooms"
	authservice "github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/services/auth/grpc"
	employeeservice "github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/services/employees/grpc"
	eventservice "github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/services/events/grpc"
	roomservice "github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/services/rooms/grpc"
	"log/slog"
	"net/http"
)

type App struct {
	server *httpserver.Server
}

func MustNew(logger *slog.Logger, cfg *config.Config) *App {
	authService := authservice.MustNew(&cfg.Auth)
	employeesService := employeeservice.MustNew(&cfg.Employees)
	eventsService := eventservice.MustNew(&cfg.Events)
	roomsService := roomservice.MustNew(&cfg.Rooms)

	authHandler := authhandlers.New(logger, authService, employeesService)
	employeesHandler := employeeshandlers.New(logger, employeesService)
	eventsHandler := eventshandlers.New(logger, eventsService)
	roomsHandler := roomshandlers.New(logger, roomsService)

	srv := httpserver.New(
		cfg,
		authHandler,
		employeesHandler,
		eventsHandler,
		roomsHandler,
	)

	return &App{server: srv}
}

func (a *App) MustRun() {
	if err := a.server.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}

func (a *App) ShutDown() error {
	return a.server.ShutDown()
}
