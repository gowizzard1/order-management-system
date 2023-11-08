package client

import (
	"context"
	"github.com/leta/order-management-system/orders/generated"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type OrdersClient interface {
	UpdateOrderStatus(ctx context.Context, req *generated.UpdateOrderStatusRequest) (*generated.UpdateOrderStatusResponse, error)
}
type GrpcOrderClient struct {
	conn   *grpc.ClientConn
	client generated.OrdersClient
}

func NewGrpcOrderClient(conn *grpc.ClientConn) *GrpcOrderClient {

	client := generated.NewOrdersClient(conn)

	return &GrpcOrderClient{
		conn:   conn,
		client: client,
	}
}

func ConnectToOrderService(address string) (*grpc.ClientConn, error) {

	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	return conn, nil
}
