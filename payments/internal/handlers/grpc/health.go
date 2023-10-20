package grpc

import (
	"context"
	"github.com/leta/order-management-system/payments/generated"
	"log"
)

func (s *GRPCServer) HealthCheck(ctx context.Context, in *generated.HealthCheckRequest) (*generated.HealthCheckResponse, error) {

	log.Printf("Received: Health check request")

	return &generated.HealthCheckResponse{
		Status: "OK",
	}, nil
}
