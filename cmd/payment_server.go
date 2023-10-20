package cmd

import (
	"context"
	"log"
	"os"

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

		s := ecom_grpc.NewGRPCServer()

		mpesaService := mpesa.NewMpesaService()

		// Setup orders service client
		conn, err := orders.ConnectToOrderService("localhost:50051")
		if err != nil {
			log.Fatalf("Failed to connect to orders service: %v", err)
		}
		orderClient := orders.NewGrpcOrderClient(conn)

		// Setup firebase client and firestore service
		firebase := db.NewFirebaseService()
		firestoreClient, err := firebase.GetApp().Firestore(ctx)
		if err != nil {
			log.Fatalf("failed to create firestore client: %v", err)
		}
		defer firestoreClient.Close()

		firestoreService := db.NewFirestoreService(firestoreClient)
		paymentRepository := db.NewPaymentsRepository(firestoreService)

		paymentService := mpesa.NewPaymentsService(mpesaService, orderClient, paymentRepository)

		// Register internal services
		s.PaymentsService = paymentService

		if err := s.Run(ctx, bindAddress, port); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}

		return
	},
}

func init() {
	rootCmd.AddCommand(serveCommand)
}
