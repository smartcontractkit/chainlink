package aptos

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
)

func TestDeployCache(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)
	cfg := memory.MemoryEnvironmentConfig{
		Nodes:       1,
		AptosChains: 1,
	}
	env := memory.NewMemoryEnvironment(t, lggr, zapcore.DebugLevel, cfg)

	chainSelector := env.AllChainSelectors()[0]
	fmt.Println(chainSelector)

	resp, err := commonChangesets.Apply(t, env, nil,
		commonChangesets.Configure(
			DeployPlatformChangeset,
			types.DeployConfig{
				ChainsToDeploy: []uint64{chainSelector},
			},
		),
	)
	require.NoError(t, err)
	require.NotNil(t, resp)

	addrs, err := resp.ExistingAddresses.AddressesForChain(chainSelector)
	require.NoError(t, err)
	require.Len(t, addrs, 1)
}
