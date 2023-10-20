// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.23.4
// source: payments.proto

package generated

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// PaymentsClient is the client API for Payments service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PaymentsClient interface {
	HealthCheck(ctx context.Context, in *HealthCheckRequest, opts ...grpc.CallOption) (*HealthCheckResponse, error)
	ProcessMpesaPayment(ctx context.Context, in *MpesaPaymentRequest, opts ...grpc.CallOption) (*MpesaPaymentResponse, error)
}

type paymentsClient struct {
	cc grpc.ClientConnInterface
}

func NewPaymentsClient(cc grpc.ClientConnInterface) PaymentsClient {
	return &paymentsClient{cc}
}

func (c *paymentsClient) HealthCheck(ctx context.Context, in *HealthCheckRequest, opts ...grpc.CallOption) (*HealthCheckResponse, error) {
	out := new(HealthCheckResponse)
	err := c.cc.Invoke(ctx, "/payments.Payments/HealthCheck", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paymentsClient) ProcessMpesaPayment(ctx context.Context, in *MpesaPaymentRequest, opts ...grpc.CallOption) (*MpesaPaymentResponse, error) {
	out := new(MpesaPaymentResponse)
	err := c.cc.Invoke(ctx, "/payments.Payments/ProcessMpesaPayment", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PaymentsServer is the server API for Payments service.
// All implementations must embed UnimplementedPaymentsServer
// for forward compatibility
type PaymentsServer interface {
	HealthCheck(context.Context, *HealthCheckRequest) (*HealthCheckResponse, error)
	ProcessMpesaPayment(context.Context, *MpesaPaymentRequest) (*MpesaPaymentResponse, error)
	mustEmbedUnimplementedPaymentsServer()
}

// UnimplementedPaymentsServer must be embedded to have forward compatible implementations.
type UnimplementedPaymentsServer struct {
}

func (UnimplementedPaymentsServer) HealthCheck(context.Context, *HealthCheckRequest) (*HealthCheckResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HealthCheck not implemented")
}
func (UnimplementedPaymentsServer) ProcessMpesaPayment(context.Context, *MpesaPaymentRequest) (*MpesaPaymentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ProcessMpesaPayment not implemented")
}
func (UnimplementedPaymentsServer) mustEmbedUnimplementedPaymentsServer() {}

// UnsafePaymentsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PaymentsServer will
// result in compilation errors.
type UnsafePaymentsServer interface {
	mustEmbedUnimplementedPaymentsServer()
}

func RegisterPaymentsServer(s grpc.ServiceRegistrar, srv PaymentsServer) {
	s.RegisterService(&Payments_ServiceDesc, srv)
}

func _Payments_HealthCheck_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HealthCheckRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentsServer).HealthCheck(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/payments.Payments/HealthCheck",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentsServer).HealthCheck(ctx, req.(*HealthCheckRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Payments_ProcessMpesaPayment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MpesaPaymentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentsServer).ProcessMpesaPayment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/payments.Payments/ProcessMpesaPayment",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentsServer).ProcessMpesaPayment(ctx, req.(*MpesaPaymentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Payments_ServiceDesc is the grpc.ServiceDesc for Payments service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Payments_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "payments.Payments",
	HandlerType: (*PaymentsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "HealthCheck",
			Handler:    _Payments_HealthCheck_Handler,
		},
		{
			MethodName: "ProcessMpesaPayment",
			Handler:    _Payments_ProcessMpesaPayment_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "payments.proto",
}
