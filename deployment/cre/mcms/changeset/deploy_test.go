package changeset_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	testenv "github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/cre/mcms/changeset"
	"github.com/smartcontractkit/chainlink/deployment/cre/mcms/pkg"
)

func TestMCMSDeployment(t *testing.T) {
	t.Parallel()
	// Test environment with a simulated EVM blockchain
	loader := testenv.NewLoader()
	env, err := loader.Load(t.Context(), testenv.WithEVMSimulatedN(t, 2))
	require.NoError(t, err)

	selectors := env.BlockChains.ListChainSelectors()
	require.Len(t, selectors, 2)
	// Create runtime instance
	r := runtime.NewFromEnvironment(*env)

	// Execute a changeset
	task := runtime.ChangesetTask(changeset.CsMCMSDeploy{}, changeset.DeployChangesetInput{
		ConfigID:       "test-config",
		Labels:         map[string]string{"foo": "bar"},
		ChainSelectors: selectors,
		ContractConfiguration: changeset.ContractConfiguration{
			Config: testMCMSCfg,
		},
		//
	})

	err = r.Exec(task)
	require.NoError(t, err)

	// Verify deployment results - 5 mcms contracts per chain
	for _, selector := range selectors {
		addrs := r.State().DataStore.Addresses().Filter(datastore.AddressRefByChainSelector(selector))
		require.NoError(t, err)
		assert.Len(t, addrs, 5)
		for _, addr := range addrs {
			assert.Contains(t, addr.Labels.List(), "mcms_config=test-config")
			assert.Contains(t, addr.Labels.List(), "foo=bar")
		}
	}
}

var d = 5 * time.Second
var eoa = common.HexToAddress("0xA01E9eD15b18D3688D0B84D88a98ed750D56999B")
var testMCMSCfg = types.MCMSWithTimelockConfigV2{
	Proposer: pkg.MustGetMCMSConfig(1, []common.Address{
		eoa,
	}, nil),
	Bypasser:         pkg.MustGetMCMSConfig(1, []common.Address{eoa}, nil),
	Canceller:        pkg.MustGetMCMSConfig(1, []common.Address{eoa}, nil),
	TimelockMinDelay: big.NewInt(int64(d.Seconds())),
}
