// Package proto contains placeholder types for the Arbiter gRPC service.
// TODO: Replace this file with the actual generated proto types.
//
//nolint:revive // This file mimics protobuf generated code naming conventions
package proto

import (
	"context"

	"google.golang.org/grpc"
)

// ReleaseStatus represents the status of a Helm release.
type ReleaseStatus int32

const (
	ReleaseStatus_RELEASE_STATUS_UNSPECIFIED       ReleaseStatus = 0
	ReleaseStatus_RELEASE_STATUS_INSTALLING        ReleaseStatus = 1
	ReleaseStatus_RELEASE_STATUS_INSTALL_SUCCEEDED ReleaseStatus = 2
	ReleaseStatus_RELEASE_STATUS_INSTALL_FAILED    ReleaseStatus = 3
	ReleaseStatus_RELEASE_STATUS_READY             ReleaseStatus = 4
	ReleaseStatus_RELEASE_STATUS_DEGRADED          ReleaseStatus = 5
	ReleaseStatus_RELEASE_STATUS_UNKNOWN           ReleaseStatus = 6
)

func (s ReleaseStatus) String() string {
	switch s {
	case ReleaseStatus_RELEASE_STATUS_INSTALLING:
		return "RELEASE_STATUS_INSTALLING"
	case ReleaseStatus_RELEASE_STATUS_INSTALL_SUCCEEDED:
		return "RELEASE_STATUS_INSTALL_SUCCEEDED"
	case ReleaseStatus_RELEASE_STATUS_INSTALL_FAILED:
		return "RELEASE_STATUS_INSTALL_FAILED"
	case ReleaseStatus_RELEASE_STATUS_READY:
		return "RELEASE_STATUS_READY"
	case ReleaseStatus_RELEASE_STATUS_DEGRADED:
		return "RELEASE_STATUS_DEGRADED"
	case ReleaseStatus_RELEASE_STATUS_UNKNOWN:
		return "RELEASE_STATUS_UNKNOWN"
	default:
		return "RELEASE_STATUS_UNSPECIFIED"
	}
}

// ShardReplica contains the status and metadata of a single shard replica.
type ShardReplica struct {
	Status  ReleaseStatus
	Message string
	Metrics map[string]float64
}

// ScaleIntentRequest is sent by the scaler to report current state and desired replicas.
type ScaleIntentRequest struct {
	CurrentReplicas     map[string]*ShardReplica
	DesiredReplicaCount int32
	Reason              string
}

// ScaleIntentResponse acknowledges receipt of the scale intent.
type ScaleIntentResponse struct {
	Status string
}

// GetScalingSpecRequest is used to query the current scaling specification.
type GetScalingSpecRequest struct {
	ScalableUnitName string
}

// ScalingSpecResponse contains the arbiter's view of the scaling state.
type ScalingSpecResponse struct {
	CurrentReplicaCount  int32
	DesiredReplicaCount  int32
	ApprovedReplicaCount int32
	LastScalingReason    string
}

// HealthCheckRequest is used to verify service health.
type HealthCheckRequest struct{}

// HealthCheckResponse indicates service health status.
type HealthCheckResponse struct {
	Status string
}

// ArbiterServiceServer is the server interface.
type ArbiterServiceServer interface {
	SubmitScaleIntent(context.Context, *ScaleIntentRequest) (*ScaleIntentResponse, error)
	GetScalingSpec(context.Context, *GetScalingSpecRequest) (*ScalingSpecResponse, error)
	HealthCheck(context.Context, *HealthCheckRequest) (*HealthCheckResponse, error)
}

// UnimplementedArbiterServiceServer provides default implementations.
type UnimplementedArbiterServiceServer struct{}

func (UnimplementedArbiterServiceServer) SubmitScaleIntent(context.Context, *ScaleIntentRequest) (*ScaleIntentResponse, error) {
	return nil, nil
}

func (UnimplementedArbiterServiceServer) GetScalingSpec(context.Context, *GetScalingSpecRequest) (*ScalingSpecResponse, error) {
	return nil, nil
}

func (UnimplementedArbiterServiceServer) HealthCheck(context.Context, *HealthCheckRequest) (*HealthCheckResponse, error) {
	return nil, nil
}

// RegisterArbiterServiceServer registers the server with grpc.
func RegisterArbiterServiceServer(s *grpc.Server, srv ArbiterServiceServer) {
	s.RegisterService(&ArbiterService_ServiceDesc, srv)
}

// ArbiterService_ServiceDesc is the service descriptor.
var ArbiterService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "arbiter.v1.ArbiterService",
	HandlerType: (*ArbiterServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SubmitScaleIntent",
			Handler:    _ArbiterService_SubmitScaleIntent_Handler,
		},
		{
			MethodName: "GetScalingSpec",
			Handler:    _ArbiterService_GetScalingSpec_Handler,
		},
		{
			MethodName: "HealthCheck",
			Handler:    _ArbiterService_HealthCheck_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "arbiter.proto",
}

func _ArbiterService_SubmitScaleIntent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ScaleIntentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ArbiterServiceServer).SubmitScaleIntent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/arbiter.v1.ArbiterService/SubmitScaleIntent",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ArbiterServiceServer).SubmitScaleIntent(ctx, req.(*ScaleIntentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ArbiterService_GetScalingSpec_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetScalingSpecRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ArbiterServiceServer).GetScalingSpec(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/arbiter.v1.ArbiterService/GetScalingSpec",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ArbiterServiceServer).GetScalingSpec(ctx, req.(*GetScalingSpecRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ArbiterService_HealthCheck_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HealthCheckRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ArbiterServiceServer).HealthCheck(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/arbiter.v1.ArbiterService/HealthCheck",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ArbiterServiceServer).HealthCheck(ctx, req.(*HealthCheckRequest))
	}
	return interceptor(ctx, in, info, handler)
}
