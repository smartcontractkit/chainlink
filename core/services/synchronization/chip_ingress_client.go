package synchronization

import (
	"fmt"
	"time"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"

	"github.com/smartcontractkit/chainlink-common/pkg/chipingress"
)

// ChipIngressClientConfig configures the gRPC client for the legacy telemetry chip-ingress endpoint.
type ChipIngressClientConfig struct {
	Endpoint           string // host:port
	InsecureConnection bool
	AuthHeaders        map[string]string // static CSA auth headers (keystore.BuildBeholderAuth)
	AuthHeadersTTL     time.Duration     // >0 selects rotating CSA tokens (must be >= 10m); same value Beholder uses
	AuthPublicKeyHex   string
	AuthKeySigner      chipingress.Signer
	MeterProvider      metric.MeterProvider
	TracerProvider     trace.TracerProvider
}

// NewChipIngressClient builds the chip-ingress client used for legacy telemetry.
func NewChipIngressClient(cfg ChipIngressClientConfig) (chipingress.Client, error) {
	auth, err := chipingress.NewHeaderProvider(chipingress.HeaderProviderConfig{
		AuthHeaders:        cfg.AuthHeaders,
		AuthHeadersTTL:     cfg.AuthHeadersTTL,
		AuthPublicKeyHex:   cfg.AuthPublicKeyHex,
		AuthKeySigner:      cfg.AuthKeySigner,
		InsecureConnection: cfg.InsecureConnection,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to build chip-ingress auth: %w", err)
	}

	// The transport option must be applied before WithTokenAuth: WithTokenAuth
	// reads the connection security at application time to decide whether its
	// per-RPC credentials require TLS.
	opts := []chipingress.Opt{}
	if cfg.InsecureConnection {
		opts = append(opts, chipingress.WithInsecureConnection())
	} else {
		opts = append(opts, chipingress.WithTLS())
	}
	if auth != nil {
		// NewHeaderProvider returns nil only when there are no headers and no TTL,
		// which cannot happen for the node (BuildBeholderAuth always returns headers),
		// but keep the nil check for safety.
		opts = append(opts, chipingress.WithTokenAuth(auth))
	}
	opts = append(opts,
		chipingress.WithNOPLookup(),
		chipingress.WithMeterProvider(cfg.MeterProvider),
		chipingress.WithTracerProvider(cfg.TracerProvider),
	)

	// grpc.NewClient connects lazily, so an unreachable endpoint does not fail
	// startup; a malformed host:port does (that is intended).
	return chipingress.NewClient(cfg.Endpoint, opts...)
}
