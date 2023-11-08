package main

import (
	"context"
	"github.com/leta/order-management-system/payments/internal/api/payments"
	"log"
	"os"

	o "github.com/leta/order-management-system/orders/pkg/client"
	db "github.com/leta/order-management-system/payments/db/firebase"
	"github.com/leta/order-management-system/payments/internal/handlers/grpc"
	"github.com/leta/order-management-system/payments/internal/mpesa"
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

	s := grpc.NewGRPCServer()

	mpesaService := mpesa.NewMpesaService()

	// Setup orders service client
	conn, err := o.ConnectToOrderService("localhost:50051")
	if err != nil {
		log.Fatalf("Failed to connect to orders service: %v", err)
	}
	orderClient := o.NewGrpcOrderClient(conn)

	// Setup firebase client and firestore service
	firebase := db.NewFirebaseService()
	firestoreClient, err := firebase.GetApp().Firestore(ctx)
	if err != nil {
		log.Fatalf("failed to create firestore client: %v", err)
	}
	defer firestoreClient.Close()

	firestoreService := db.NewFirestoreService(firestoreClient)
	paymentRepository := payments.NewPaymentsRepository(firestoreService)

	paymentService := mpesa.NewPaymentsService(mpesaService, orderClient, paymentRepository)

	// Register internal services
	s.PaymentsService = paymentService

	if err := s.Run(ctx, bindAddress, port); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
