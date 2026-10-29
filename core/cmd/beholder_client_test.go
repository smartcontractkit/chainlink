package cmd

import (
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/mocks"
)

func TestNewBeholderClient_DomainEmittersFollowTelemetryIngressBeholderEnabled(t *testing.T) {
	t.Parallel()
	// A local gRPC server with no services answers OTel exports at once, so Close is quick.
	lis, err := (&net.ListenConfig{}).Listen(t.Context(), "tcp", "127.0.0.1:0")
	require.NoError(t, err)
	srv := grpc.NewServer()
	go func() { _ = srv.Serve(lis) }()
	t.Cleanup(srv.Stop)
	beholderTOML := fmt.Sprintf(`
[Telemetry]
Enabled = true
Endpoint = '%s'
InsecureConnection = true
ChipIngressEndpoint = '%s'
ChipIngressInsecureConnection = true
`, lis.Addr(), lis.Addr())
	for _, tt := range []struct {
		name    string
		toml    string
		enabled bool
	}{
		{name: "off", toml: beholderTOML},
		{name: "on", toml: beholderTOML + "\n[TelemetryIngress]\nBeholderEnabled = true\n", enabled: true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cfg, err := chainlink.GeneralConfigOpts{ConfigStrings: []string{tt.toml}}.New()
			require.NoError(t, err)

			ks := mocks.NewMaster(t)
			ks.On("CSA").Return(mocks.NewCSA(t)).Maybe()

			client, err := newBeholderClient(logger.TestLogger(t), ks, cfg.Tracing(), cfg, cfg.Telemetry(), "abc123", nil)
			require.NoError(t, err)
			require.NoError(t, client.Start(t.Context()))
			t.Cleanup(func() {
				// OTel shutdown may fail to flush to the absent collector; the domain emitters
				// must close cleanly.
				if closeErr := client.Close(); closeErr != nil {
					assert.NotContains(t, closeErr.Error(), "client connection is closing")
				}
			})

			_, err = client.ForDomain(beholder.DomainDataFeeds)
			if tt.enabled {
				require.NoError(t, err, "BeholderEnabled builds the domain emitters")
			} else {
				require.Error(t, err, "no domain emitters unless BeholderEnabled")
			}
		})
	}
}
