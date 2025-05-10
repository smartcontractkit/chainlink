package channel_config_store

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/metadata"

	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/testutil"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
)

func TestDeployChannelConfigStore(t *testing.T) {
	t.Parallel()

	e := testutil.NewMemoryEnv(t, false, 0)

	cc := DeployChannelConfigStoreConfig{
		ChainsToDeploy: []uint64{testutil.TestChain.Selector},
	}

	out, err := DeployChannelConfigStoreChangeset.Apply(e, cc)
	require.NoError(t, err)

	envDatastore, err := datastore.FromDefault[metadata.SerializedContractMetadata, datastore.DefaultMetadata](out.DataStore.Seal())
	require.NoError(t, err)

	// Verify Contract Is Deployed
	record, err := envDatastore.Addresses().Get(
		datastore.NewAddressRefKey(testutil.TestChain.Selector, datastore.ContractType(types.ChannelConfigStore), &deployment.Version1_0_0, ""),
	)
	require.NoError(t, err)
	require.NotNil(t, record)
}
