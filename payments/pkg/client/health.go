package client

import (
	"context"
)

func (c *GrpcPaymentsClient) HealthCheck(ctx context.Context, req *HealthCheckRequest) (*HealthCheckResponse, error) {
	return c.client.HealthCheck(ctx, req)
}
