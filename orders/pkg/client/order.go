package client

import (
	"context"
)

func (c *GrpcOrderClient) UpdateOrderStatus(
	ctx context.Context, req *UpdateOrderStatusRequest) (*UpdateOrderStatusResponse, error) {
	return c.client.UpdateOrderStatus(ctx, req)
}
