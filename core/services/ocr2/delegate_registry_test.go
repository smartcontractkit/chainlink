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

func TestRegistryOCRKeyBundles(t *testing.T) {
	t.Parallel()

	evmKey, err := ocr2key.New(corekeys.EVM)
	require.NoError(t, err)
	solKey, err := ocr2key.New(corekeys.Solana)
	require.NoError(t, err)
	otherEVMKey, err := ocr2key.New(corekeys.EVM)
	require.NoError(t, err)

	multichainSigner, err := ocrcommon.MarshalMultichainPublicKey(map[string]ocrtypes.OnchainPublicKey{
		string(corekeys.EVM):    evmKey.PublicKey(),
		string(corekeys.Solana): solKey.PublicKey(),
	})
	require.NoError(t, err)

	keyStore := keystoremocks.NewOCR2(t)
	keyStore.EXPECT().GetAll().Return([]ocr2key.KeyBundle{otherEVMKey, evmKey, solKey}, nil)

	got, idx, err := registryOCRKeyBundles(keyStore, &ocrtypes.ContractConfig{
		Signers: []ocrtypes.OnchainPublicKey{[]byte("other-signer"), multichainSigner},
	})
	require.NoError(t, err)
	assert.Equal(t, 1, idx)
	require.Len(t, got, 2)
	assert.Equal(t, evmKey.ID(), got[string(corekeys.EVM)].ID())
	assert.Equal(t, solKey.ID(), got[string(corekeys.Solana)].ID())
}

func TestRegistryOCRKeyBundles_NoMatch(t *testing.T) {
	t.Parallel()

	key, err := ocr2key.New(corekeys.EVM)
	require.NoError(t, err)
	keyStore := keystoremocks.NewOCR2(t)
	keyStore.EXPECT().GetAll().Return([]ocr2key.KeyBundle{key}, nil)

	_, _, err = registryOCRKeyBundles(keyStore, &ocrtypes.ContractConfig{
		Signers: []ocrtypes.OnchainPublicKey{[]byte("other-signer")},
	})
	require.ErrorContains(t, err, "none of the node's 1 OCR key bundles match the signers in the registry OCR config")

	_, _, err = registryOCRKeyBundles(keystoremocks.NewOCR2(t), nil)
	require.ErrorContains(t, err, "no OCR config in the Capabilities Registry")
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
