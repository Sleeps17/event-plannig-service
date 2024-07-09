package grpcserver

import (
	"fmt"
	authgrpc "github.com/Sleeps17/events-planning-service-backend/auth_service/internal/grpc/auth"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

type Server struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(log *slog.Logger, authService authgrpc.Auth, port int) *Server {
	gRPCServer := grpc.NewServer()

	authgrpc.Register(gRPCServer, authService)

	return &Server{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (s *Server) Run() error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return fmt.Errorf("can't create listener: %w", err)
	}

	s.log.Info("gRPC server is running: ", slog.String("addr", l.Addr().String()))

	if err := s.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("can't serve gRPC server: %w", err)
	}

	return nil
}

func (s *Server) Stop() {
	s.log.Info("stopping gRPC server")

	s.gRPCServer.GracefulStop()
}
