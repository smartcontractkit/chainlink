package chainlink

import (
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	ksmocks "github.com/smartcontractkit/chainlink/v2/core/services/keystore/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/telemetry"
)

type nopLogExporter struct{}

func (nopLogExporter) Export(context.Context, []sdklog.Record) error { return nil }
func (nopLogExporter) Shutdown(context.Context) error                { return nil }
func (nopLogExporter) ForceFlush(context.Context) error              { return nil }

func TestNewTelemetryManager(t *testing.T) {
	t.Parallel()
	telemetryIngress := func(t *testing.T, beholderEnabled bool) GeneralConfig {
		t.Helper()
		toml := ""
		if beholderEnabled {
			toml = "[TelemetryIngress]\nBeholderEnabled = true\n"
		}
		cfg, err := GeneralConfigOpts{ConfigStrings: []string{toml}}.New()
		require.NoError(t, err)
		return cfg
	}
	lggr := logger.Sugared(logger.TestLogger(t))

	t.Run("flag off: the WSRPC manager", func(t *testing.T) {
		t.Parallel()
		m, err := newTelemetryManager(telemetryIngress(t, false).TelemetryIngress(), nil, ksmocks.NewCSA(t), "abc123", lggr)
		require.NoError(t, err)
		assert.IsType(t, &telemetry.Manager{}, m)
	})
	t.Run("flag on without a Beholder client: an error, not a panic", func(t *testing.T) {
		t.Parallel()
		var noClient *beholder.Client
		_, err := newTelemetryManager(telemetryIngress(t, true).TelemetryIngress(), noClient, ksmocks.NewCSA(t), "abc123", lggr)
		require.Error(t, err)
	})
	t.Run("flag on with Beholder telemetry off (the no-op client): an error", func(t *testing.T) {
		t.Parallel()
		_, err := newTelemetryManager(telemetryIngress(t, true).TelemetryIngress(), beholder.NewNoopClient(), ksmocks.NewCSA(t), "abc123", lggr)
		require.Error(t, err, "the no-op client has no domain emitters")
	})
	t.Run("flag on: the domain manager on Beholder's domain emitters", func(t *testing.T) {
		t.Parallel()
		// A local gRPC server with no services answers OTel exports at once, so Close is quick.
		lis, err := (&net.ListenConfig{}).Listen(t.Context(), "tcp", "127.0.0.1:0")
		require.NoError(t, err)
		srv := grpc.NewServer()
		go func() { _ = srv.Serve(lis) }()
		t.Cleanup(srv.Stop)

		client, err := beholder.NewGRPCClient(beholder.Config{
			OtelExporterGRPCEndpoint:         lis.Addr().String(),
			InsecureConnection:               true,
			ChipIngressEmitterGRPCEndpoint:   lis.Addr().String(),
			ChipIngressInsecureConnection:    true,
			ChipIngressDomainEmittersEnabled: true,
		}, func(...otlploggrpc.Option) (sdklog.Exporter, error) { return nopLogExporter{}, nil })
		require.NoError(t, err)
		require.NoError(t, client.Start(t.Context()))
		t.Cleanup(func() { _ = client.Close() })

		m, err := newTelemetryManager(telemetryIngress(t, true).TelemetryIngress(), client, ksmocks.NewCSA(t), "abc123", lggr)
		require.NoError(t, err)
		assert.IsType(t, &telemetry.DomainManager{}, m)
	})
}
