package client

import (
	"context"
)

func (c *GrpcOrderClient) HealthCheck(ctx context.Context, req *HealthCheckRequest) (*HealthCheckResponse, error) {
	return c.client.HealthCheck(ctx, req)
}
