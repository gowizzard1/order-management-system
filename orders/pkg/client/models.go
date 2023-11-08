package client

import (
	"github.com/leta/order-management-system/orders/generated"
)

//type HealthCheckRequest = utils.HealthCheckRequest
//type HealthCheckResponse = utils.HealthCheckResponse

type UpdateOrderStatusRequest = generated.UpdateOrderStatusRequest
type UpdateOrderStatusResponse = generated.UpdateOrderStatusResponse

var OrderStatusPaid = generated.OrderStatus_PAID
var OrderStatusCancelled = generated.OrderStatus_CANCELLED
var OrderStatusFailed = generated.OrderStatus_FAILED
var OrderStatusPending = generated.OrderStatus_PENDING
