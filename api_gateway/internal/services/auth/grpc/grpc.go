package grpc

import (
	"fmt"
	"github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/config"
	authv1 "github.com/Sleeps17/events-planning-service-backend/api_gateway/protos/gen/go/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Service struct {
	client authv1.AuthClient
}

func MustNew(cfg *config.GrpcServiceConfig) *Service {
	cc, err := grpc.NewClient(
		fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}

	return &Service{client: authv1.NewAuthClient(cc)}
}
