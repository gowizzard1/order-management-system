package grpc_test

import (
	"context"
	"fmt"
	ecom_grpc "github.com/leta/order-management-system/payments/internal/handlers/grpc"
	"net"
	"testing"
	"time"

	"github.com/leta/order-management-system/payments/internal/handlers/grpc"
	"github.com/leta/order-management-system/payments/internal/mock"
)

const (
	INVALID_PORT = "70000"
)

type TestGRPCServer struct {
	*grpc.GRPCServer

	// Add mock services here
	PaymentsService mock.PaymentsService
}

func NewTestGRPCServer(tb testing.TB) *TestGRPCServer {
	s := &TestGRPCServer{
		GRPCServer: ecom_grpc.NewGRPCServer(),
	}

	// Set mock services here
	s.GRPCServer.PaymentsService = &s.PaymentsService

	return s
}

func TestGRPCServer_Run(t *testing.T) {
	tests := []struct {
		name     string
		bindAddr string
		port     string
		wantErr  bool
	}{
		{
			name:     "valid bind address and port",
			bindAddr: "localhost",
			port:     getFreePort(t),
			wantErr:  false,
		},
		{
			name:     "invalid port",
			bindAddr: "localhost",
			port:     INVALID_PORT,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := ecom_grpc.NewGRPCServer()
			errCh := make(chan error)

			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			go func() {
				errCh <- server.Run(ctx, tt.bindAddr, tt.port)
			}()

			// Wait for server to respond or timeout after a specified time
			select {
			case err := <-errCh:
				if (err != nil) != tt.wantErr {
					t.Errorf("GRPCServer.Run() error = %v, wantErr %v", err, tt.wantErr)
				}
			case <-ctx.Done():
				if tt.wantErr {
					t.Errorf("GRPCServer.Run() expected error for port %v, but got none", tt.port)
				}
			}

			server.Stop()
		})
	}
}

// Helper function to get a free port
func getFreePort(t *testing.T) string {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("error resolving tcp address: %v", err)
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		t.Fatalf("error listening on tcp address: %v", err)
	}

	defer listener.Close()
	return fmt.Sprintf("%d", listener.Addr().(*net.TCPAddr).Port)
}
