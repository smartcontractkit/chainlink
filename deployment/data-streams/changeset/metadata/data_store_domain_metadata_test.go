package metadata

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
)

func TestEnvMetadata(t *testing.T) {
	ds := datastore.NewMemoryDataStore[
		SerializedContractMetadata,
		DataStreamsMetadata,
	]()

	metaData, err := ds.EnvMetadataStore.Get()
	require.Error(t, err)
	require.EqualError(t, datastore.ErrEnvMetadataNotSet, err.Error())

	metaData.Metadata.DONs = []DonMetadata{
		{
			ID:                  "don1",
			ConfiguratorAddress: "0x1234567890abcdef",
			OffchainConfig: OffchainConfig{
				DeltaGrace:   "1",
				DeltaInitial: "1",
			},
			Streams: []int{1, 2},
		},
	}

	err = ds.EnvMetadataStore.Set(metaData)
	require.NoError(t, err)

	metadata, err := ds.EnvMetadataStore.Get()
	require.NoError(t, err)

	don, err := metadata.Metadata.GetDonByID("don1")
	require.NoError(t, err)
	require.Len(t, don.Streams, 2)
}
