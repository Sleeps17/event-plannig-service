package grpcapp

import (
	"errors"
	"fmt"
	"github.com/Sleeps17/event-plannig-service-backend/events-service/internal/config"
	"github.com/Sleeps17/event-plannig-service-backend/events-service/internal/grpc/api"
	employeeservice "github.com/Sleeps17/event-plannig-service-backend/events-service/internal/services/employees"
	eventservice "github.com/Sleeps17/event-plannig-service-backend/events-service/internal/services/events"
	roomservice "github.com/Sleeps17/event-plannig-service-backend/events-service/internal/services/rooms"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

type Server struct {
	srv    *grpc.Server
	logger *slog.Logger
	port   string
}

func New(
	cfg *config.GrpcConfig,
	logger *slog.Logger,
	eventsProvider eventservice.Service,
	roomsProvider roomservice.Service,
	employeesProvider employeeservice.Service,
) *Server {
	grpcServer := grpc.NewServer()

	api.Register(
		grpcServer,
		logger,
		eventsProvider,
		roomsProvider,
		employeesProvider,
	)

	return &Server{
		srv:    grpcServer,
		logger: logger,
		port:   cfg.Port,
	}
}

func (s *Server) MustRun() {
	l, err := net.Listen("tcp", fmt.Sprintf(":%s", s.port))
	if err != nil {
		panic(fmt.Errorf("failed to listen: %w", err))
	}

	s.logger.Info(fmt.Sprintf("gRPC server is running on port: %s", s.port))
	if err := s.srv.Serve(l); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		panic(fmt.Errorf("failed to serve: %w", err))
	}
}

func (s *Server) Stop() {
	s.logger.Info("gRPC server is stopping")
	s.srv.GracefulStop()
}
