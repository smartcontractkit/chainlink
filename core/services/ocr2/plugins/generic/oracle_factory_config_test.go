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
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
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

	match, ok := MatchOCRSigner([]ocr2key.KeyBundle{kb}, cc)
	require.True(t, ok)
	transmitter, ok := TransmitterAt(*cc, match.Index)
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

func TestResolveOracleFactoryConfig_multiFamilyKeyBundlesFillSigning(t *testing.T) {
	t.Parallel()

	evmKB, err := ocr2key.New(corekeys.EVM)
	require.NoError(t, err)
	solKB, err := ocr2key.New(corekeys.Solana)
	require.NoError(t, err)

	cfg, signing, err := ResolveOracleFactoryConfig(context.Background(), ResolveOracleFactoryConfigParams{
		Config:        job.OracleFactoryConfig{Enabled: true},
		OCRKeyBundles: map[string]ocr2key.KeyBundle{"evm": evmKB, "solana": solKB},
		Logger:        logger.TestLogger(t),
	})
	require.NoError(t, err)
	assert.Equal(t, evmKB.ID(), cfg.OCRKeyBundleID)
	assert.Equal(t, map[string]string{"evm": evmKB.ID(), "solana": solKB.ID()}, signing.Config)
}

func TestTransmitterAt(t *testing.T) {
	t.Parallel()

	cc := ocrtypes.ContractConfig{
		Signers:      []ocrtypes.OnchainPublicKey{[]byte("a"), []byte("b")},
		Transmitters: []ocrtypes.Account{"0xA", "0xB"},
	}

	got, ok := TransmitterAt(cc, 1)
	require.True(t, ok)
	assert.Equal(t, "0xB", got)

	_, ok = TransmitterAt(cc, 2)
	assert.False(t, ok)
	_, ok = TransmitterAt(cc, -1)
	assert.False(t, ok)

	got, ok = TransmitterAt(ocrtypes.ContractConfig{
		Signers:      []ocrtypes.OnchainPublicKey{[]byte("a")},
		Transmitters: []ocrtypes.Account{"736ea02dd58a4eff74565801cb9cf1d13ceb9134"},
	}, 0)
	require.True(t, ok)
	assert.Equal(t, "0x736ea02Dd58A4EFF74565801cB9Cf1D13CEB9134", got)
}

func multichainSigner(t *testing.T, kbs ...ocr2key.KeyBundle) ocrtypes.OnchainPublicKey {
	t.Helper()
	pubKeys := map[string]ocrtypes.OnchainPublicKey{}
	for _, kb := range kbs {
		pubKeys[string(kb.ChainType())] = kb.PublicKey()
	}
	signer, err := ocrcommon.MarshalMultichainPublicKey(pubKeys)
	require.NoError(t, err)
	return signer
}

func TestMatchOCRSigner(t *testing.T) {
	t.Parallel()

	evm1, err := ocr2key.New(corekeys.EVM)
	require.NoError(t, err)
	evm2, err := ocr2key.New(corekeys.EVM)
	require.NoError(t, err)
	sol, err := ocr2key.New(corekeys.Solana)
	require.NoError(t, err)
	aptos, err := ocr2key.New(corekeys.Aptos)
	require.NoError(t, err)
	otherSol, err := ocr2key.New(corekeys.Solana)
	require.NoError(t, err)

	bundles := []ocr2key.KeyBundle{evm1, evm2, sol, aptos}

	t.Run("nil config", func(t *testing.T) {
		_, ok := MatchOCRSigner(bundles, nil)
		assert.False(t, ok)
	})

	t.Run("no match", func(t *testing.T) {
		_, ok := MatchOCRSigner(bundles, &ocrtypes.ContractConfig{
			Signers: []ocrtypes.OnchainPublicKey{[]byte("nope"), nil},
		})
		assert.False(t, ok)
	})

	t.Run("raw EVM signer picks registered bundle", func(t *testing.T) {
		match, ok := MatchOCRSigner(bundles, &ocrtypes.ContractConfig{
			Signers: []ocrtypes.OnchainPublicKey{[]byte("other"), evm2.PublicKey()},
		})
		require.True(t, ok)
		assert.Equal(t, 1, match.Index)
		require.Len(t, match.KeyBundles, 1)
		assert.Equal(t, evm2.ID(), match.KeyBundles["evm"].ID())
		assert.Equal(t, evm2.ID(), match.PrimaryKeyBundle().ID())
	})

	t.Run("multichain EVM-only signer", func(t *testing.T) {
		match, ok := MatchOCRSigner(bundles, &ocrtypes.ContractConfig{
			Signers: []ocrtypes.OnchainPublicKey{multichainSigner(t, evm2)},
		})
		require.True(t, ok)
		assert.Equal(t, 0, match.Index)
		require.Len(t, match.KeyBundles, 1)
		assert.Equal(t, evm2.ID(), match.KeyBundles["evm"].ID())
	})

	t.Run("multichain multi-family signer", func(t *testing.T) {
		match, ok := MatchOCRSigner(bundles, &ocrtypes.ContractConfig{
			Signers: []ocrtypes.OnchainPublicKey{
				multichainSigner(t, evm1, otherSol),
				multichainSigner(t, evm2, sol, aptos),
			},
		})
		require.True(t, ok)
		assert.Equal(t, 1, match.Index)
		require.Len(t, match.KeyBundles, 3)
		assert.Equal(t, evm2.ID(), match.KeyBundles["evm"].ID())
		assert.Equal(t, sol.ID(), match.KeyBundles["solana"].ID())
		assert.Equal(t, aptos.ID(), match.KeyBundles["aptos"].ID())
		assert.Equal(t, evm2.ID(), match.PrimaryKeyBundle().ID())
	})

	t.Run("multichain signer with a family the node lacks", func(t *testing.T) {
		_, ok := MatchOCRSigner(bundles, &ocrtypes.ContractConfig{
			Signers: []ocrtypes.OnchainPublicKey{multichainSigner(t, evm1, otherSol)},
		})
		assert.False(t, ok)
	})

	t.Run("multichain non-EVM signer", func(t *testing.T) {
		match, ok := MatchOCRSigner(bundles, &ocrtypes.ContractConfig{
			Signers: []ocrtypes.OnchainPublicKey{multichainSigner(t, sol, aptos)},
		})
		require.True(t, ok)
		require.Len(t, match.KeyBundles, 2)
		// No EVM bundle: falls back to the lowest-sorted family.
		assert.Equal(t, aptos.ID(), match.PrimaryKeyBundle().ID())
	})
}

func TestOCRSignerMatch_PrimaryKeyBundleEmpty(t *testing.T) {
	t.Parallel()

	assert.Nil(t, OCRSignerMatch{}.PrimaryKeyBundle())
}

func TestDefaultTransmitterForChain_InvalidChainID(t *testing.T) {
	t.Parallel()

	// DefaultTransmitterForChain requires a real keystore; covered indirectly via integration.
	_, err := DefaultTransmitterForChain(context.Background(), nil, "not-a-number")
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
