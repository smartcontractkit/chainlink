package generic

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	ocrcommontypes "github.com/smartcontractkit/libocr/commontypes"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys"
	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/ocr2key"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

func TestResolveOracleFactoryConfig_fromCapRegistry(t *testing.T) {
	t.Parallel()

	cfg, signing, err := ResolveOracleFactoryConfig(context.Background(), ResolveOracleFactoryConfigParams{
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

	cfg, _, err := ResolveOracleFactoryConfig(context.Background(), ResolveOracleFactoryConfigParams{
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

func TestResolveOracleFactoryConfig_transmitterFromOCRConfig(t *testing.T) {
	t.Parallel()

	kb, err := ocr2key.New(corekeys.EVM)
	require.NoError(t, err)

	cc := &ocrtypes.ContractConfig{
		Signers: []ocrtypes.OnchainPublicKey{
			ocrtypes.OnchainPublicKey("other-signer"),
			kb.PublicKey(),
		},
		Transmitters: []ocrtypes.Account{"0xOther", "0xMine"},
	}

	cfg, _, err := ResolveOracleFactoryConfig(context.Background(), ResolveOracleFactoryConfigParams{
		Config:            job.OracleFactoryConfig{Enabled: true},
		OCRKeyBundle:      kb,
		OCRContractConfig: cc,
		Logger:            logger.TestLogger(t),
	})
	require.NoError(t, err)
	assert.Equal(t, "0xMine", cfg.TransmitterID)
}

func TestResolveOracleFactoryConfig_jobSpecTransmitterOverridesOCRConfig(t *testing.T) {
	t.Parallel()

	kb, err := ocr2key.New(corekeys.EVM)
	require.NoError(t, err)

	cc := &ocrtypes.ContractConfig{
		Signers:      []ocrtypes.OnchainPublicKey{kb.PublicKey()},
		Transmitters: []ocrtypes.Account{"0xFromOCR"},
	}

	cfg, _, err := ResolveOracleFactoryConfig(context.Background(), ResolveOracleFactoryConfigParams{
		Config:            job.OracleFactoryConfig{Enabled: true, TransmitterID: "0xFromSpec"},
		OCRKeyBundle:      kb,
		OCRContractConfig: cc,
		Logger:            logger.TestLogger(t),
	})
	require.NoError(t, err)
	assert.Equal(t, "0xFromSpec", cfg.TransmitterID)
}

func TestTransmitterForSigner(t *testing.T) {
	t.Parallel()

	cc := ocrtypes.ContractConfig{
		Signers:      []ocrtypes.OnchainPublicKey{[]byte("a"), []byte("b")},
		Transmitters: []ocrtypes.Account{"0xA", "0xB"},
	}

	got, ok := transmitterForSigner(cc, []byte("b"))
	require.True(t, ok)
	assert.Equal(t, "0xB", got)

	_, ok = transmitterForSigner(cc, []byte("missing"))
	assert.False(t, ok)

	// Signer present but transmitter list too short.
	short := ocrtypes.ContractConfig{
		Signers:      []ocrtypes.OnchainPublicKey{[]byte("a")},
		Transmitters: nil,
	}
	_, ok = transmitterForSigner(short, []byte("a"))
	assert.False(t, ok)
}

func TestSelectOCRKeyBundleForConfig(t *testing.T) {
	t.Parallel()

	kb1, err := ocr2key.New(corekeys.EVM)
	require.NoError(t, err)
	kb2, err := ocr2key.New(corekeys.EVM)
	require.NoError(t, err)

	cc := &ocrtypes.ContractConfig{
		Signers: []ocrtypes.OnchainPublicKey{kb2.PublicKey()},
	}

	got, ok := SelectOCRKeyBundleForConfig([]ocr2key.KeyBundle{kb1, kb2}, cc)
	require.True(t, ok)
	assert.Equal(t, kb2.ID(), got.ID())

	_, ok = SelectOCRKeyBundleForConfig([]ocr2key.KeyBundle{kb1, kb2}, nil)
	assert.False(t, ok)

	noMatch := &ocrtypes.ContractConfig{Signers: []ocrtypes.OnchainPublicKey{[]byte("nope")}}
	_, ok = SelectOCRKeyBundleForConfig([]ocr2key.KeyBundle{kb1, kb2}, noMatch)
	assert.False(t, ok)
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
}
