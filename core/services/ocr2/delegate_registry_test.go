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
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
)

func TestRegistryOCRKeyBundle(t *testing.T) {
	t.Parallel()

	key, err := ocr2key.New(corekeys.EVM)
	require.NoError(t, err)
	keyStore := keystoremocks.NewOCR2(t)
	keyStore.EXPECT().GetAllOfType(corekeys.EVM).Return([]ocr2key.KeyBundle{key}, nil)

	got, err := registryOCRKeyBundle(keyStore, &ocrtypes.ContractConfig{
		Signers: []ocrtypes.OnchainPublicKey{key.PublicKey()},
	})
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, key.ID(), got.ID())
}

func TestRegistryOCR2SpecRelayID(t *testing.T) {
	t.Parallel()

	spec := job.OCR2OracleSpec{
		Relay:       relay.NetworkEVM,
		ChainID:     "1337",
		RelayConfig: job.JSONConfig{"chainID": "1337", "providerType": string(commontypes.DonTimePlugin)},
	}

	relayID, err := spec.RelayID()
	require.NoError(t, err)
	assert.Equal(t, commontypes.RelayID{Network: relay.NetworkEVM, ChainID: "1337"}, relayID)
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
