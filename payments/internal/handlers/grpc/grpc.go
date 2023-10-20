package grpc

import (
	"github.com/leta/order-management-system/payments/internal/service"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func LogError(err error) {
	log.Printf("[grpc] error: %s", err)
}

func GRPCErrorStatusCode(err error) error {
	code := service.ErrorCode(err)

	switch code {
	case service.INVALID_ERROR:
		return status.Error(codes.InvalidArgument, service.ErrorMessage(err))
	case service.NOT_FOUND_ERROR:
		return status.Error(codes.NotFound, service.ErrorMessage(err))
	case service.ALREADY_EXISTS_ERROR:
		return status.Error(codes.AlreadyExists, service.ErrorMessage(err))
	case service.INTERNAL_ERROR:
		return status.Error(codes.Internal, service.ErrorMessage(err))
	default:
		return status.Error(codes.Unknown, service.ErrorMessage(err))
	}

}
