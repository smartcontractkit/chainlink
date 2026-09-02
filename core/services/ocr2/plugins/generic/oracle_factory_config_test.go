package generic

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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
		Transmitter:        "0xTx",
		Logger:             logger.TestLogger(t),
	})
	require.NoError(t, err)
	assert.Equal(t, "0xabc", cfg.OCRContractAddress)
	assert.Equal(t, "1337", cfg.ChainID)
	assert.Empty(t, signing.Config)
}

func TestResolveOracleFactoryConfig_disabledSkipsResolution(t *testing.T) {
	t.Parallel()

	cfg, signing, err := ResolveOracleFactoryConfig(context.Background(), ResolveOracleFactoryConfigParams{
		Config:             job.OracleFactoryConfig{Enabled: false},
		CapRegistryAddress: "0xabc",
		CapRegistryChainID: "1337",
		Logger:             logger.TestLogger(t),
	})
	require.NoError(t, err)
	assert.Empty(t, cfg.OCRContractAddress)
	assert.Empty(t, cfg.ChainID)
	assert.Empty(t, cfg.TransmitterID)
	assert.Empty(t, signing.Config)
}

func TestResolveOracleFactoryConfig_jobSpecOverridesCapRegistry(t *testing.T) {
	t.Parallel()

	cfg, _, err := ResolveOracleFactoryConfig(context.Background(), ResolveOracleFactoryConfigParams{
		Config: job.OracleFactoryConfig{
			Enabled:            true,
			OCRContractAddress: "0xjob",
			ChainID:            "1",
			TransmitterID:      "0xTx",
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

	transmitter, ok := TransmitterForSigner(*cc, kb.PublicKey())
	require.True(t, ok)
	assert.Equal(t, "0xMine", transmitter)

	cfg, _, err := ResolveOracleFactoryConfig(context.Background(), ResolveOracleFactoryConfigParams{
		Config:      job.OracleFactoryConfig{Enabled: true},
		Transmitter: transmitter,
		Logger:      logger.TestLogger(t),
	})
	require.NoError(t, err)
	assert.Equal(t, "0xMine", cfg.TransmitterID)
}

func TestResolveOracleFactoryConfig_jobSpecTransmitterOverridesOCRConfig(t *testing.T) {
	t.Parallel()

	cfg, _, err := ResolveOracleFactoryConfig(context.Background(), ResolveOracleFactoryConfigParams{
		Config:      job.OracleFactoryConfig{Enabled: true, TransmitterID: "0xFromSpec"},
		Transmitter: "0xFromOCR",
		Logger:      logger.TestLogger(t),
	})
	require.NoError(t, err)
	assert.Equal(t, "0xFromSpec", cfg.TransmitterID)
}

func TestResolveOracleFactoryConfig_keyBundlesFillSigning(t *testing.T) {
	t.Parallel()

	kb, err := ocr2key.New(corekeys.EVM)
	require.NoError(t, err)

	cfg, signing, err := ResolveOracleFactoryConfig(context.Background(), ResolveOracleFactoryConfigParams{
		Config:        job.OracleFactoryConfig{Enabled: true},
		OCRKeyBundles: map[string]ocr2key.KeyBundle{"evm": kb},
		Logger:        logger.TestLogger(t),
	})
	require.NoError(t, err)
	assert.Equal(t, kb.ID(), cfg.OCRKeyBundleID)
	assert.Equal(t, "multi-chain", signing.StrategyName)
	assert.Equal(t, kb.ID(), signing.Config["evm"])
}

func TestTransmitterForSigner(t *testing.T) {
	t.Parallel()

	cc := ocrtypes.ContractConfig{
		Signers:      []ocrtypes.OnchainPublicKey{[]byte("a"), []byte("b")},
		Transmitters: []ocrtypes.Account{"0xA", "0xB"},
	}

	got, ok := TransmitterForSigner(cc, []byte("b"))
	require.True(t, ok)
	assert.Equal(t, "0xB", got)

	_, ok = TransmitterForSigner(cc, []byte("missing"))
	assert.False(t, ok)

	// Signer present but transmitter list too short.
	short := ocrtypes.ContractConfig{
		Signers:      []ocrtypes.OnchainPublicKey{[]byte("a")},
		Transmitters: nil,
	}
	_, ok = TransmitterForSigner(short, []byte("a"))
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

func TestDefaultTransmitterForChain_InvalidChainID(t *testing.T) {
	t.Parallel()

	// defaultTransmitterForChain requires a real keystore; covered indirectly via integration.
	_, err := defaultTransmitterForChain(context.Background(), nil, "not-a-number")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid chain_id")
}

func TestResolveOracleFactoryConfig_keystoreFallbackTransmitter(t *testing.T) {
	t.Parallel()

	// When no transmitter is provided from the on-chain OCR config, and the job spec
	// doesn't set one, ResolveOracleFactoryConfig should fall back to the keystore.
	// We can't test the full keystore path without a real keystore, but we can verify
	// that the function doesn't set a transmitter when EthKeystore is nil and ChainID is set.
	cfg, _, err := ResolveOracleFactoryConfig(context.Background(), ResolveOracleFactoryConfigParams{
		Config:             job.OracleFactoryConfig{Enabled: true},
		CapRegistryAddress: "0xabc",
		CapRegistryChainID: "1337",
		Transmitter:        "",  // no on-chain transmitter
		EthKeystore:        nil, // no keystore
		Logger:             logger.TestLogger(t),
	})
	require.NoError(t, err)
	assert.Equal(t, "0xabc", cfg.OCRContractAddress)
	assert.Equal(t, "1337", cfg.ChainID)
	assert.Empty(t, cfg.TransmitterID, "transmitter should be empty when no keystore fallback available")
}
