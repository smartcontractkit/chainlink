package testutil

import (
	"testing"

	"github.com/stretchr/testify/require"

	ds "github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/metadata"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/view/interfaces"
)

func MustGetContractMetaData[T interfaces.ContractView](
	t *testing.T,
	inDs ds.MutableDataStore[ds.DefaultMetadata, ds.DefaultMetadata],
	chainSelector uint64,
	contractAddress string,
) *metadata.GenericContractMetadata[T] {
	envDatastore, err := ds.FromDefault[metadata.SerializedContractMetadata, ds.DefaultMetadata](inDs.Seal())
	require.NoError(t, err)

	cm, err := envDatastore.ContractMetadata().Get(
		ds.NewContractMetadataKey(chainSelector, contractAddress),
	)
	require.NoError(t, err, "Failed to get contract metadata")

	contractMetadata, err := metadata.DeserializeMetadata[T](cm.Metadata)
	require.NoError(t, err, "Failed to convert contract metadata")
	require.NotNil(t, contractMetadata, "Failed to get contract metadata")
	return contractMetadata
}
