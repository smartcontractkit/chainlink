package nodeconfig

import (
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

var update = flag.Bool("update", false, "update golden files")

func checkGolden(t *testing.T, name, got string) {
	t.Helper()
	path := filepath.Join("testdata", name)
	if *update {
		require.NoError(t, os.WriteFile(path, []byte(got), 0600))
		return
	}
	want, err := os.ReadFile(path)
	require.NoError(t, err)
	require.Equal(t, string(want), got, "golden mismatch for %s", name)
}

func TestGenerate_Worker(t *testing.T) {
	t.Parallel()

	got, err := Generate(Inputs{
		CapRegAddress:      "0x1234567890123456789012345678901234567890",
		WorkflowRegAddress: "0xabcdefabcdefabcdefabcdefabcdefabcdefabcd",
		RegistryChainID:    1337,
		Allowlist:          []string{"cron-trigger@1.0.0", "evm-1337"},
		BootstrapPeerID:    "12D3KooWQyADpbmd1QsrxEHGNRZMDtAERy88tJNGyVowFgAPQpMu",
		BootstrapHost:      "node-bt-0.default.svc.cluster.local",
		P2PPort:            5001,
	})
	require.NoError(t, err)
	require.NotContains(t, got, "[[EVM]]")
	checkGolden(t, "worker.toml", got)
}

func TestGenerate_Bootstrap(t *testing.T) {
	t.Parallel()

	got, err := Generate(Inputs{
		CapRegAddress:   "0x1234567890123456789012345678901234567890",
		RegistryChainID: 1337,
		IsBootstrapNode: true,
		BootstrapPeerID: "12D3KooWQyADpbmd1QsrxEHGNRZMDtAERy88tJNGyVowFgAPQpMu",
		P2PPort:         5001,
	})
	require.NoError(t, err)
	require.NotContains(t, got, "WorkflowRegistry")
	require.NotContains(t, got, "[[EVM]]")
	checkGolden(t, "bootstrap.toml", got)
}

func TestGenerate_Gateway(t *testing.T) {
	t.Parallel()

	got, err := Generate(Inputs{
		CapRegAddress:      "0x1234567890123456789012345678901234567890",
		WorkflowRegAddress: "0xabcdefabcdefabcdefabcdefabcdefabcdefabcd",
		RegistryChainID:    1337,
		IsGatewayNode:      true,
		P2PPort:            5001,
	})
	require.NoError(t, err)
	require.NotContains(t, got, "[[EVM]]")
	checkGolden(t, "gateway.toml", got)
}
