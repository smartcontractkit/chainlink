package ocr2

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys"
	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/ocr2key"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	keystoremocks "github.com/smartcontractkit/chainlink/v2/core/services/keystore/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
)

func TestRegistryOCRSignerMatch(t *testing.T) {
	t.Parallel()

	evmKey, err := ocr2key.New(corekeys.EVM)
	require.NoError(t, err)
	solKey, err := ocr2key.New(corekeys.Solana)
	require.NoError(t, err)
	keyStore := keystoremocks.NewOCR2(t)
	keyStore.EXPECT().GetAll().Return([]ocr2key.KeyBundle{evmKey, solKey}, nil)

	signer, err := ocrcommon.MarshalMultichainPublicKey(map[string]ocrtypes.OnchainPublicKey{
		string(corekeys.EVM):    evmKey.PublicKey(),
		string(corekeys.Solana): solKey.PublicKey(),
	})
	require.NoError(t, err)

	match, ok, err := registryOCRSignerMatch(keyStore, &ocrtypes.ContractConfig{
		Signers: []ocrtypes.OnchainPublicKey{signer},
	})
	require.NoError(t, err)
	require.True(t, ok)
	assert.Equal(t, evmKey.ID(), match.PrimaryKeyBundle().ID())
	assert.Equal(t, solKey.ID(), match.KeyBundles["solana"].ID())

	_, ok, err = registryOCRSignerMatch(keyStore, nil)
	require.NoError(t, err)
	assert.False(t, ok)
}

func TestRegistryOCR2SpecRelayID(t *testing.T) {
	t.Parallel()

	transmitterID := "0x1234"
	spec := job.OCR2OracleSpec{
		Relay:       relay.NetworkEVM,
		ChainID:     "1337",
		RelayConfig: registryOCR2RelayConfig("1337", commontypes.DonTimePlugin, transmitterID),
	}

	relayID, err := spec.RelayID()
	require.NoError(t, err)
	assert.Equal(t, commontypes.RelayID{Network: relay.NetworkEVM, ChainID: "1337"}, relayID)
	assert.Equal(t, transmitterID, spec.RelayConfig["effectiveTransmitterID"])
	assert.Equal(t, []string{transmitterID}, spec.RelayConfig["sendingKeys"])
}

func TestNewServices_NilPeerWrapper(t *testing.T) {
	t.Parallel()

	d := &Delegate{
		lggr: logger.TestLogger(t),
	}

	_, err := d.NewServices(context.Background(), "dontime@1.0.0", 1, commontypes.DonTimePlugin, "", nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "libp2p peer was missing or not started")
}
