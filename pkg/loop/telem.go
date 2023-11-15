package loop

import (
	"context"
	"net"
	"os"
	"runtime/debug"

	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/smartcontractkit/chainlink-relay/pkg/loop/internal"
)

type GRPCOpts = internal.GRPCOpts

type TracingConfig struct {
	// NodeAttributes are the attributes to attach to traces.
	NodeAttributes map[string]string

	// Enables tracing; requires a collector to be provided
	Enabled bool

	// Collector is the address of the OTEL collector to send traces to.
	CollectorTarget string

	// SamplingRatio is the ratio of traces to sample. 1.0 means sample all traces.
	SamplingRatio float64

	// TLSCertPath is the path to the TLS certificate to use when connecting to the collector.
	TLSCertPath string

	// OnDialError is called when the dialer fails, providing an opportunity to log.
	OnDialError func(error)
}

// NewGRPCOpts initializes open telemetry and returns GRPCOpts with telemetry interceptors.
// It is called from the host and each plugin - intended as there is bidirectional communication
func NewGRPCOpts(registerer prometheus.Registerer) GRPCOpts {
	if registerer == nil {
		registerer = prometheus.DefaultRegisterer
	}
	return GRPCOpts{DialOpts: dialOptions(registerer), NewServer: newServerFn(registerer)}
}

// SetupTracing initializes open telemetry with the provided config.
// It sets the global trace provider and opens a connection to the configured collector.
func SetupTracing(config TracingConfig) (err error) {
	if !config.Enabled {
		return nil
	}

	ctx := context.Background()

	var creds credentials.TransportCredentials
	if config.TLSCertPath != "" {
		creds, err = credentials.NewClientTLSFromFile(config.TLSCertPath, "")
		if err != nil {
			return err
		}
	} else {
		creds = insecure.NewCredentials()
	}

	conn, err := grpc.DialContext(ctx, config.CollectorTarget,
		// Note the potential use of insecure transport here. TLS is recommended in production.
		grpc.WithTransportCredentials(creds),
		grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
			conn, err2 := net.Dial("tcp", s)
			if err2 != nil {
				config.OnDialError(err2)
			}
			return conn, err2
		}))
	if err != nil {
		return err
	}

	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return err
	}

	var version string
	var service string
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		version = "unknown"
		service = "cl-node"
	} else {
		version = buildInfo.Main.Version
		service = buildInfo.Main.Path
	}

	attributes := []attribute.KeyValue{
		semconv.ServiceNameKey.String(service),
		semconv.ProcessPIDKey.Int(os.Getpid()),
		semconv.ServiceVersionKey.String(version),
	}

	for k, v := range config.NodeAttributes {
		attributes = append(attributes, attribute.String(k, v))
	}

	resource := resource.NewWithAttributes(
		semconv.SchemaURL,
		attributes...,
	)

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(resource),
		sdktrace.WithSampler(
			sdktrace.ParentBased(
				sdktrace.TraceIDRatioBased(config.SamplingRatio),
			),
		),
	)

	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return nil
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
