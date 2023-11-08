package handlers

import (
	"context"
	"fmt"
	"github.com/leta/order-management-system/orders/generated"
	"github.com/leta/order-management-system/orders/internal/interfaces/api/customers"
	"github.com/leta/order-management-system/orders/internal/interfaces/api/orders"
	"github.com/leta/order-management-system/orders/internal/interfaces/api/product"
	"log"
	"net"
	"sync"

	"github.com/leta/order-management-system/orders/internal/service"
	"google.golang.org/grpc"
)

// GRPCServer struct represents the GRPC server for the orders management system.
type GRPCServer struct {
	generated.UnimplementedOrdersServer

	grpcServer *grpc.Server
	mu         sync.Mutex // synchronizes access to the grpcServer

	// Internal services &  repositories
	CheckoutService service.CheckoutService

	ProductRepository  product.RepositoryInterface
	CustomerService    customers.CustomerServiceInterface
	OrderService       orders.OrderServiceInterface
	OrderRepository    orders.OrderRepository
	CustomerRepository customers.CustomerRepositoryInterface
}

// NewGRPCServer creates a new instance of GRPCServer.
func NewGRPCServer() *GRPCServer {
	return &GRPCServer{}
}

// Run starts the GRPC server on the specified bind address and port.
func (s *GRPCServer) Run(ctx context.Context, bindAddress string, port string) error {
	lis, err := net.Listen("tcp", bindAddress+":"+port)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	s.mu.Lock()
	s.grpcServer = grpc.NewServer()
	s.mu.Unlock()

	generated.RegisterOrdersServer(s.grpcServer, s)

	log.Printf("Starting server on port %v", lis.Addr())

	return s.grpcServer.Serve(lis)
}

// Stop gracefully stops the GRPC server.
func (s *GRPCServer) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.grpcServer != nil {
		s.grpcServer.GracefulStop()
	}
}
