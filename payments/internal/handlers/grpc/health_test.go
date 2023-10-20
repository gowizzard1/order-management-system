package grpc_test

import (
	"context"
	"github.com/leta/order-management-system/payments/generated"
	"reflect"
	"testing"
)

func TestGRPCServer_HealthCheck(t *testing.T) {

	s := NewTestGRPCServer(t)

	type args struct {
		ctx context.Context
		in  *generated.HealthCheckRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *generated.HealthCheckResponse
		wantErr bool
	}{
		{
			name: "Test HealthCheck",
			args: args{
				ctx: context.Background(),
				in:  &generated.HealthCheckRequest{},
			},
			want: &generated.HealthCheckResponse{
				Status: "OK",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := s.HealthCheck(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("GRPCServer.HealthCheck() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GRPCServer.HealthCheck() = %v, want %v", got, tt.want)
			}
		})
	}
}
