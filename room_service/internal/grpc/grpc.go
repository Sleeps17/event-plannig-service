package grpcserver

import (
	"errors"
	"fmt"
	"github.com/Sleeps17/events-planning-service-backend/rooms_service/internal/config"
	"github.com/Sleeps17/events-planning-service-backend/rooms_service/internal/grpc/api"
	roomservice "github.com/Sleeps17/events-planning-service-backend/rooms_service/internal/services/rooms"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

type Server struct {
	logger *slog.Logger
	port   string
	srv    *grpc.Server
}

func New(
	cfg *config.GrpcConfig,
	logger *slog.Logger,
	roomsProvider roomservice.Service,
) *Server {
	grpcServer := grpc.NewServer()

	api.Register(
		grpcServer,
		logger,
		roomsProvider,
	)

	return &Server{
		srv:    grpcServer,
		logger: logger,
		port:   cfg.Port,
	}
}

func (s *Server) Run() error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%s", s.port))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	s.logger.Info(fmt.Sprintf("gRPC server is running on port: %s", s.port))
	if err := s.srv.Serve(l); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}

func (s *Server) Stop() {
	s.srv.GracefulStop()
}
