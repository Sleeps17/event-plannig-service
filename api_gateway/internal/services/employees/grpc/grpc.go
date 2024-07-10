package grpc

import (
	"fmt"
	"github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/config"
	employeesv1 "github.com/Sleeps17/events-planning-service-backend/api_gateway/protos/gen/go/employees"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Service struct {
	client employeesv1.EmployeesClient
}

func MustNew(cfg *config.GrpcServiceConfig) *Service {
	cc, err := grpc.NewClient(
		fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}

	return &Service{client: employeesv1.NewEmployeesClient(cc)}
}
