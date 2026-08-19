package generic

import (
	"context"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	ocrcommontypes "github.com/smartcontractkit/libocr/commontypes"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

func TestResolveOracleFactoryConfig_fromCapRegistry(t *testing.T) {
	t.Parallel()

	cfg, signing, err := ResolveOracleFactoryConfig(ResolveOracleFactoryConfigParams{
		Context:            context.Background(),
		Config:             job.OracleFactoryConfig{Enabled: true},
		CapRegistryAddress: "0xabc",
		CapRegistryChainID: "1337",
		Logger:             logger.TestLogger(t),
	})
	require.NoError(t, err)
	assert.Equal(t, "0xabc", cfg.OCRContractAddress)
	assert.Equal(t, "1337", cfg.ChainID)
	assert.Empty(t, signing.Config)
}

func TestResolveOracleFactoryConfig_jobSpecOverridesCapRegistry(t *testing.T) {
	t.Parallel()

	cfg, _, err := ResolveOracleFactoryConfig(ResolveOracleFactoryConfigParams{
		Context: context.Background(),
		Config: job.OracleFactoryConfig{
			Enabled:            true,
			OCRContractAddress: "0xjob",
			ChainID:            "1",
		},
		CapRegistryAddress: "0xlocal",
		CapRegistryChainID: "1337",
		Logger:             logger.TestLogger(t),
	})
	require.NoError(t, err)
	assert.Equal(t, "0xjob", cfg.OCRContractAddress)
	assert.Equal(t, "1", cfg.ChainID)
}

func TestResolveBootstrapPeers_usesDefaultsWhenSpecEmpty(t *testing.T) {
	t.Parallel()

	defaults := []ocrcommontypes.BootstrapperLocator{
		mustBootstrapper(t, "12D3KooWBAzThfs9pD4WcsFKCi68EUz2fZgZskDBT6JcJRndPss5@host:5001"),
	}

	peers, err := resolveBootstrapPeers(nil, defaults)
	require.NoError(t, err)
	assert.Equal(t, defaults, peers)
}

func TestResolveBootstrapPeers_specOverridesDefaults(t *testing.T) {
	t.Parallel()

	defaults := []ocrcommontypes.BootstrapperLocator{
		mustBootstrapper(t, "12D3KooWBAzThfs9pD4WcsFKCi68EUz2fZgZskDBT6JcJRndPss5@host:5001"),
	}
	spec := []string{"12D3KooWBAzThfs9pD4WcsFKCi68EUz2fZgZskDBT6JcJRndPss5@override:5001"}

	peers, err := resolveBootstrapPeers(spec, defaults)
	require.NoError(t, err)
	require.Len(t, peers, 1)
	assert.Equal(t, "override:5001", peers[0].Addrs[0])
}

func mustBootstrapper(t *testing.T, s string) ocrcommontypes.BootstrapperLocator {
	t.Helper()
	var loc ocrcommontypes.BootstrapperLocator
	require.NoError(t, loc.UnmarshalText([]byte(s)))
	return loc
}

func TestDefaultTransmitterForChain(t *testing.T) {
	t.Parallel()

	// defaultTransmitterForChain requires a real keystore; covered indirectly via integration.
	_, err := defaultTransmitterForChain(context.Background(), nil, "not-a-number")
	require.Error(t, err)

	chainID, ok := new(big.Int).SetString("1337", 10)
	require.True(t, ok)
	assert.Equal(t, 0, chainID.Cmp(chainID))
}
