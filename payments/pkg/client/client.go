package client

import (
	"context"
	"github.com/leta/order-management-system/payments/generated"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var _ PaymentsClient = (*GrpcPaymentsClient)(nil)

type PaymentsClient interface {
	HealthCheck(ctx context.Context, req *HealthCheckRequest) (*HealthCheckResponse, error)
	ProcessMpesaPayment(ctx context.Context, req *ProcessMpesaPaymentRequest) (*ProcessMpesaPaymentResponse, error)
}

type GrpcPaymentsClient struct {
	conn   *grpc.ClientConn
	client generated.PaymentsClient
}

func NewGrpcPaymentsClient(conn *grpc.ClientConn) *GrpcPaymentsClient {

	client := generated.NewPaymentsClient(conn)

	return &GrpcPaymentsClient{
		conn:   conn,
		client: client,
	}
}

func ConnectToPaymentService(address string) (*grpc.ClientConn, error) {

	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	return conn, nil
}
