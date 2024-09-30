package core

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel/metric"
	"log"
	"os"
	"time"

	"github.com/Masterminds/semver/v3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/exp/rand"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/beholder/pb"
	"github.com/smartcontractkit/chainlink/v2/core/build"
	"github.com/smartcontractkit/chainlink/v2/core/cmd"
	"github.com/smartcontractkit/chainlink/v2/core/recovery"
	"github.com/smartcontractkit/chainlink/v2/core/static"
)

func init() {
	// check version
	if static.Version == static.Unset {
		if !build.IsProd() {
			return
		}
		log.Println(`Version was unset on production build. Chainlink should be built with static.Version set to a valid semver for production builds.`)
	} else if _, err := semver.NewVersion(static.Version); err != nil {
		panic(fmt.Sprintf("Version invalid: %q is not valid semver", static.Version))
	}
}

func Main() (code int) {
	setupBeholder()
	go sendCustomMessages()
	go sendMetricTraces()

	recovery.ReportPanics(func() {
		app := cmd.NewApp(newProductionClient())
		if err := app.Run(os.Args); err != nil {
			fmt.Fprintf(os.Stderr, "Error running app: %v\n", err)
			code = 1
		}
	})
	return
}

// newProductionClient configures an instance of the CLI to be used in production.
func newProductionClient() *cmd.Shell {
	prompter := cmd.NewTerminalPrompter()
	return &cmd.Shell{
		Renderer:                       cmd.RendererTable{Writer: os.Stdout},
		AppFactory:                     cmd.ChainlinkAppFactory{},
		KeyStoreAuthenticator:          cmd.TerminalKeyStoreAuthenticator{Prompter: prompter},
		FallbackAPIInitializer:         cmd.NewPromptingAPIInitializer(prompter),
		Runner:                         cmd.ChainlinkRunner{},
		PromptingSessionRequestBuilder: cmd.NewPromptingSessionRequestBuilder(prompter),
		ChangePasswordPrompter:         cmd.NewChangePasswordPrompter(),
		PasswordPrompter:               cmd.NewPasswordPrompter(),
	}
}

func beholderDevConfig() beholder.Config {
	config := beholder.DefaultConfig()
	// Set the OTel exporter endpoint
	config.OtelExporterGRPCEndpoint = "localhost:4317"
	// Add some more Resource Attributes
	// Resource Attributes are static and are added to each emitted OTel data type
	config.ResourceAttributes = append(config.ResourceAttributes, []attribute.KeyValue{
		attribute.String("chain_id", "11155111"),
		attribute.String("node_id", "dev-node-id"),
	}...)
	// Emitter
	// Disable batching, should not be used in production
	config.EmitterBatchProcessor = false
	// Trace
	config.TraceSampleRatio = 1
	config.TraceBatchTimeout = 1 * time.Second
	// Metric
	config.MetricReaderInterval = 1 * time.Second
	// Log
	config.LogExportTimeout = 1 * time.Second
	// Disable batching, should not be used in production
	config.LogBatchProcessor = false
	return config
}

func setupBeholder() {
	config := beholderDevConfig()

	log.Printf("Beholder config: %#v", config)

	// Initialize beholder otel client which sets up OTel components
	otelClient, err := beholder.NewClient(context.Background(), config)
	if err != nil {
		log.Fatalf("Error creating Beholder client: %v", err)
	}
	// Handle OTel errors
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(e error) {
		log.Printf("otel error: %v", e)
	}))
	// Set global client so it will be accessible from anywhere through beholder/global functions
	beholder.SetClient(otelClient)
}

func sendCustomMessages() {
	// Define a custom protobuf payload to emit
	payload := &pb.TestCustomMessage{
		BoolVal:   true,
		IntVal:    42,
		FloatVal:  3.14,
		StringVal: "custom message from chainlink",
	}
	payloadBytes, err := proto.Marshal(payload)
	if err != nil {
		log.Fatalf("Failed to marshal protobuf")
	}

	// Emit the custom message anywhere from application logic
	for i := 0; ; i++ {
		log.Printf("Beholder: emitting custom message with ID: %d", i)
		err := beholder.GetEmitter().Emit(context.Background(), payloadBytes,
			"beholder_data_schema", "/custom-message/versions/1", // required
			"beholder_data_type", "custom_message",
			"message_ind", i,
		)
		if err != nil {
			log.Printf("Error emitting message: %v", err)
		}
		time.Sleep(1 * time.Second)
	}
}

func sendMetricTraces() {
	ctx := context.Background()

	// Define a new counter
	counter, err := beholder.GetMeter().Int64Counter("custom_message.count")
	if err != nil {
		log.Fatalf("failed to create new counter")
	}

	// Define a new gauge
	gauge, err := beholder.GetMeter().Int64Gauge("custom_message.gauge")
	if err != nil {
		log.Fatalf("failed to create new gauge")
	}

	for i := 0; ; i++ {
		log.Printf("Beholder: sending metric, trace  %d", i)
		// Use the counter and gauge for metrics within application logic
		labels := []attribute.KeyValue{
			attribute.String("application", "cl-node"),
			attribute.String("job", "demo-job"),
		}
		counter.Add(ctx, 1, metric.WithAttributes(labels...))
		gauge.Record(ctx, rand.Int63n(101), metric.WithAttributes(labels...))

		// Create a new trace span
		_, span := beholder.GetTracer().Start(ctx, "sendMetricTraces", trace.WithAttributes(
			attribute.String("app_name", "beholderdemo"),
			attribute.Int64("trace_ind", int64(i)),
		))
		span.End()
		time.Sleep(1 * time.Second)
	}
}
