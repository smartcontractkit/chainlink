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

	multichainSigner, err := ocrcommon.MarshalMultichainPublicKey(map[string]ocrtypes.OnchainPublicKey{
		string(corekeys.EVM): []byte("b"),
	})
	require.NoError(t, err)
	got, ok = TransmitterForSigner(ocrtypes.ContractConfig{
		Signers:      []ocrtypes.OnchainPublicKey{multichainSigner},
		Transmitters: []ocrtypes.Account{"0xMultiChain"},
	}, []byte("b"))
	require.True(t, ok)
	assert.Equal(t, "0xMultiChain", got)

	got, ok = TransmitterForSigner(ocrtypes.ContractConfig{
		Signers:      []ocrtypes.OnchainPublicKey{[]byte("a")},
		Transmitters: []ocrtypes.Account{"736ea02dd58a4eff74565801cb9cf1d13ceb9134"},
	}, []byte("a"))
	require.True(t, ok)
	assert.Equal(t, "0x736ea02Dd58A4EFF74565801cB9Cf1D13CEB9134", got)
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

	multichainSigner, err := ocrcommon.MarshalMultichainPublicKey(map[string]ocrtypes.OnchainPublicKey{
		string(corekeys.EVM): kb2.PublicKey(),
	})
	require.NoError(t, err)
	got, ok = SelectOCRKeyBundleForConfig([]ocr2key.KeyBundle{kb1, kb2}, &ocrtypes.ContractConfig{
		Signers: []ocrtypes.OnchainPublicKey{multichainSigner},
	})
	require.True(t, ok)
	assert.Equal(t, kb2.ID(), got.ID())
}

func TestSelectOCRKeyBundlesForConfig(t *testing.T) {
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

	multichain := func(kbs map[string]ocr2key.KeyBundle) ocrtypes.OnchainPublicKey {
		pk, mErr := ocrcommon.MarshalMultichainKeyBundle(kbs)
		require.NoError(t, mErr)
		return pk
	}

	t.Run("nil config", func(t *testing.T) {
		t.Parallel()
		_, idx, ok := SelectOCRKeyBundlesForConfig(bundles, nil)
		assert.False(t, ok)
		assert.Equal(t, -1, idx)
	})

	t.Run("raw EVM signer", func(t *testing.T) {
		t.Parallel()
		got, idx, ok := SelectOCRKeyBundlesForConfig(bundles, &ocrtypes.ContractConfig{
			Signers: []ocrtypes.OnchainPublicKey{[]byte("other"), evm2.PublicKey()},
		})
		require.True(t, ok)
		assert.Equal(t, 1, idx)
		require.Len(t, got, 1)
		assert.Equal(t, evm2.ID(), got["evm"].ID())
	})

	t.Run("multi-family signer", func(t *testing.T) {
		t.Parallel()
		want := map[string]ocr2key.KeyBundle{"evm": evm2, "solana": sol, "aptos": aptos}
		other := multichain(map[string]ocr2key.KeyBundle{"evm": evm1, "solana": otherSol})
		got, idx, ok := SelectOCRKeyBundlesForConfig(bundles, &ocrtypes.ContractConfig{
			Signers: []ocrtypes.OnchainPublicKey{other, multichain(want)},
		})
		require.True(t, ok)
		assert.Equal(t, 1, idx)
		require.Len(t, got, 3)
		for family, kb := range want {
			assert.Equal(t, kb.ID(), got[family].ID(), family)
		}
		// The multichain keyring built from the matched bundles must reproduce the signer.
		assert.Equal(t, multichain(want), multichain(got))
	})

	t.Run("non-EVM only signer", func(t *testing.T) {
		t.Parallel()
		got, idx, ok := SelectOCRKeyBundlesForConfig(bundles, &ocrtypes.ContractConfig{
			Signers: []ocrtypes.OnchainPublicKey{multichain(map[string]ocr2key.KeyBundle{"solana": sol})},
		})
		require.True(t, ok)
		assert.Equal(t, 0, idx)
		require.Len(t, got, 1)
		assert.Equal(t, sol.ID(), got["solana"].ID())
	})

	t.Run("missing family bundle", func(t *testing.T) {
		t.Parallel()
		_, _, ok := SelectOCRKeyBundlesForConfig(bundles, &ocrtypes.ContractConfig{
			Signers: []ocrtypes.OnchainPublicKey{multichain(map[string]ocr2key.KeyBundle{"evm": evm1, "solana": otherSol})},
		})
		assert.False(t, ok)
	})

	t.Run("no match", func(t *testing.T) {
		t.Parallel()
		_, _, ok := SelectOCRKeyBundlesForConfig(bundles, &ocrtypes.ContractConfig{
			Signers: []ocrtypes.OnchainPublicKey{[]byte("nope"), nil},
		})
		assert.False(t, ok)
	})
}

func TestTransmitterAt(t *testing.T) {
	t.Parallel()

	cc := ocrtypes.ContractConfig{
		Transmitters: []ocrtypes.Account{"x", "736ea02dd58a4eff74565801cb9cf1d13ceb9134"},
	}
	got, ok := TransmitterAt(cc, 0)
	require.True(t, ok)
	assert.Equal(t, "x", got)

	got, ok = TransmitterAt(cc, 1)
	require.True(t, ok)
	assert.Equal(t, "0x736ea02Dd58A4EFF74565801cB9Cf1D13CEB9134", got)

	_, ok = TransmitterAt(cc, 2)
	assert.False(t, ok)
	_, ok = TransmitterAt(cc, -1)
	assert.False(t, ok)
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
