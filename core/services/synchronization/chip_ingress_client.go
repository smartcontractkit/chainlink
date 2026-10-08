package synchronization

import (
	"fmt"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/chipingress"
)

// ChipIngressClientConfig configures the gRPC client for the legacy telemetry chip-ingress endpoint.
// Metrics and traces use the global OTel providers, which Beholder installs at startup.
type ChipIngressClientConfig struct {
	Endpoint           string // host:port
	InsecureConnection bool
	AuthHeaders        map[string]string // static CSA auth headers (keystore.BuildBeholderAuth)
	AuthHeadersTTL     time.Duration     // >0 selects rotating CSA tokens (must be >= 10m); same value Beholder uses
	AuthPublicKeyHex   string
	AuthKeySigner      chipingress.Signer
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

	// The transport option must come before WithTokenAuth, which reads it to
	// decide whether its per-RPC credentials require TLS.
	opts := []chipingress.Opt{chipingress.WithTLS()}
	if cfg.InsecureConnection {
		opts = []chipingress.Opt{chipingress.WithInsecureConnection()}
	}
	if auth != nil { // nil only when there are no headers and no TTL
		opts = append(opts, chipingress.WithTokenAuth(auth))
	}
	opts = append(opts, chipingress.WithNOPLookup())

	// grpc.NewClient connects lazily: an unreachable endpoint does not fail
	// startup, a malformed host:port does.
	return chipingress.NewClient(cfg.Endpoint, opts...)
}
