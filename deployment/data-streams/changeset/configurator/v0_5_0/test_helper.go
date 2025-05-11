package v0_5_0

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	ds "github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/testutil"
	dsutil "github.com/smartcontractkit/chainlink/deployment/data-streams/utils"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/view/v0_5"
)

func VerifyConfiguratorState(
	t *testing.T,
	inDs ds.MutableDataStore[ds.DefaultMetadata, ds.DefaultMetadata],
	chainSelector uint64,
	configuratorAddr common.Address,
	configDigest [32]byte,
	expectedConfig ConfiguratorConfig,
	expectedConfigCount uint64,
) {
	contractMetadata := testutil.MustGetContractMetaData[v0_5.ConfiguratorView](t, inDs, chainSelector, configuratorAddr.Hex())

	configIDString := dsutil.HexEncodeBytes(expectedConfig.ConfigID[:])
	configDigestString := dsutil.HexEncodeBytes(configDigest[:])

	// Retrieve the config state
	configState, err := contractMetadata.View.GetConfigState(configIDString, configDigestString)
	require.NoError(t, err, "Failed to get config state")

	// Verify basic configuration properties
	require.Equal(t, expectedConfigCount, configState.ConfigCount, "Configuration count mismatch")
	require.Equal(t, dsutil.HexEncodeBytes(expectedConfig.OffchainConfig), configState.EncodedOffchainConfig, "OffchainConfig mismatch")
	require.Equal(t, dsutil.HexEncodeBytes(expectedConfig.OnchainConfig), configState.EncodedOnchainConfig, "OnchainConfig mismatch")
	require.Equal(t, expectedConfig.F, configState.F, "F value mismatch")
	require.Equal(t, expectedConfig.OffchainConfigVersion, configState.OffchainConfigVersion, "OffchainConfigVersion mismatch")

	signersHex := dsutil.HexEncodeByteSlices(expectedConfig.Signers)
	require.Equal(t, signersHex, configState.Signers, "Signers mismatch")

	transmittersHex := dsutil.HexEncodeBytes32Slice(expectedConfig.OffchainTransmitters)
	require.Equal(t, transmittersHex, configState.OffchainTransmitters, "OffchainTransmitters mismatch")

	t.Log("All configurator state verifications passed")
}
