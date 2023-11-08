package main

import (
	"context"
	db "github.com/leta/order-management-system/orders/db/firebase"
	firebase2 "github.com/leta/order-management-system/orders/db/firebase"
	"github.com/leta/order-management-system/orders/internal/api/customers"
	"github.com/leta/order-management-system/orders/internal/api/orders"
	"github.com/leta/order-management-system/orders/internal/api/product"
	"log"
	"os"

	"github.com/leta/order-management-system/orders/internal/checkout"
	"github.com/leta/order-management-system/orders/internal/handlers"
	p "github.com/leta/order-management-system/payments/pkg/client"
)

const (
	BIND_ADDRESS = "BIND_ADDRESS"
	PORT         = "PORT"

	DEFAULT_BIND_ADDRESS = "localhost"
	DEFAULT_PORT         = "50051"
)

func main() {

	ctx := context.Background()

	log.Printf("Starting server")

	bindAddress := os.Getenv(BIND_ADDRESS)
	if bindAddress == "" {
		bindAddress = DEFAULT_BIND_ADDRESS
	}

	port := os.Getenv(PORT)
	if port == "" {
		port = DEFAULT_PORT
	}

	s := handlers.NewGRPCServer()

	firebase := firebase2.NewFirebaseService()
	firestoreClient, err := firebase.GetApp().Firestore(ctx)
	if err != nil {
		log.Fatalf("failed to create firestore client: %v", err)
	}
	defer firestoreClient.Close()

	firestoreService := db.NewFirestoreService(firestoreClient)

	productRepository := product.NewProductRepository(firestoreService)

	//prdSvs := product.NewProductService(productRepository)
	customerRepository := customers.NewCustomerRepository(firestoreService)
	customerSvc := customers.NewCustomerService(customerRepository)

	orderRepository := orders.NewOrderRepository(firestoreService)
	orderSvc := orders.NewOrdersService(*orderRepository)
	// Setup payments service client
	conn, err := p.ConnectToPaymentService("localhost:50051")
	if err != nil {
		log.Fatalf("Failed to connect to orders service: %v", err)
	}
	paymentsClient := p.NewGrpcPaymentsClient(conn)

	checkoutService := checkout.NewCheckoutService(
		productRepository, customerRepository, orderRepository, paymentsClient)

	s.ProductRepository = productRepository
	s.CustomerService = customerSvc
	s.OrderService = orderSvc
	s.CheckoutService = checkoutService

	if err := s.Run(ctx, bindAddress, port); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
