package client

import "context"

func (c *GrpcPaymentsClient) ProcessMpesaPayment(ctx context.Context, req *ProcessMpesaPaymentRequest) (*ProcessMpesaPaymentResponse, error) {
	return c.client.ProcessMpesaPayment(ctx, req)
}
