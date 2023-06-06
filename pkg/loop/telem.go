package loop

import (
	"context"

	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-relay/pkg/loop/internal"
)

type GRPCOpts = internal.GRPCOpts

// SetupTelemetry initializes open telemetry and returns GRPCOpts with telemetry interceptors.
func SetupTelemetry(registerer prometheus.Registerer) GRPCOpts {
	otel.SetTracerProvider(sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	))
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	if registerer == nil {
		registerer = prometheus.DefaultRegisterer
	}
	return GRPCOpts{DialOpts: dialOptions(registerer), NewServer: newServerFn(registerer)}
}

var grpcpromBuckets = []float64{0.001, 0.01, 0.1, 0.3, 0.6, 1, 3, 6, 9, 20, 30, 60, 90, 120}

// dialOptions returns [grpc.DialOption]s to intercept and reports telemetry.
func dialOptions(r prometheus.Registerer) []grpc.DialOption {
	cm := grpcprom.NewClientMetrics(
		grpcprom.WithClientHandlingTimeHistogram(grpcprom.WithHistogramBuckets(grpcpromBuckets)),
	)
	r.MustRegister(cm)
	ctxExemplar := grpcprom.WithExemplarFromContext(exemplarFromContext)
	return []grpc.DialOption{
		// Order matters e.g. tracing interceptor have to create span first for the later exemplars to work.
		grpc.WithChainUnaryInterceptor(
			otelgrpc.UnaryClientInterceptor(),
			cm.UnaryClientInterceptor(ctxExemplar),
		),
		grpc.WithChainStreamInterceptor(
			otelgrpc.StreamClientInterceptor(),
			cm.StreamClientInterceptor(ctxExemplar),
		),
	}
}

// newServerFn return a func for constructing [*grpc.Server]s that intercepts and reports telemetry.
func newServerFn(r prometheus.Registerer) func(opts []grpc.ServerOption) *grpc.Server {
	srvMetrics := grpcprom.NewServerMetrics(
		grpcprom.WithServerHandlingTimeHistogram(grpcprom.WithHistogramBuckets(grpcpromBuckets)),
	)
	r.MustRegister(srvMetrics)
	ctxExemplar := grpcprom.WithExemplarFromContext(exemplarFromContext)
	interceptors := []grpc.ServerOption{
		// Order matters e.g. tracing interceptor have to create span first for the later exemplars to work.
		grpc.ChainUnaryInterceptor(
			otelgrpc.UnaryServerInterceptor(),
			srvMetrics.UnaryServerInterceptor(ctxExemplar),
		),
		grpc.ChainStreamInterceptor(
			otelgrpc.StreamServerInterceptor(),
			srvMetrics.StreamServerInterceptor(ctxExemplar),
		),
	}
	return func(opts []grpc.ServerOption) *grpc.Server {
		s := grpc.NewServer(append(opts, interceptors...)...)
		srvMetrics.InitializeMetrics(s)
		return s
	}
}

func exemplarFromContext(ctx context.Context) prometheus.Labels {
	if span := trace.SpanContextFromContext(ctx); span.IsSampled() {
		return prometheus.Labels{"traceID": span.TraceID().String()}
	}
	return nil
}
