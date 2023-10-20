package cmd

import (
	"context"
	"log"
	"os"

	"github.com/gorilla/handlers"
	"github.com/spf13/cobra"
)

var serveCommand = &cobra.Command{
	Use:   "serve",
	Short: "serve the api",
	RunE: func(cmd *cobra.Command, args []string) (err error) {

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

		firebase := db.NewFirebaseService()
		firestoreClient, err := firebase.GetApp().Firestore(ctx)
		if err != nil {
			log.Fatalf("failed to create firestore client: %v", err)
		}
		defer firestoreClient.Close()

		firestoreService := db.NewFirestoreService(firestoreClient)

		ProductRepository := db.NewProductService(firestoreService)
		customerRepository := db.NewCustomerService(firestoreService)
		orderRepository := db.NewOrderRepository(firestoreService)

		// Setup payments service client
		conn, err := payments.ConnectToPaymentService("localhost:50051")
		if err != nil {
			log.Fatalf("Failed to connect to orders service: %v", err)
		}
		paymentsClient := payments.NewGrpcPaymentsClient(conn)

		checkoutService := checkout.NewCheckoutService(
			ProductRepository, customerRepository, orderRepository, paymentsClient)

		s.ProductRepository = ProductRepository
		s.CustomerRepository = customerRepository
		s.OrderRepository = orderRepository
		s.CheckoutService = checkoutService

		if err := s.Run(ctx, bindAddress, port); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}

		return
	},
}

func init() {
	rootCmd.AddCommand(serveCommand)
}
