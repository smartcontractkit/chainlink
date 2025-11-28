package changeset_test

import (
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	mcmsbindings "github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	testenv "github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"

	"github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/cre/mcms/changeset"
	"github.com/smartcontractkit/chainlink/deployment/cre/mcms/pkg"
)

func TestMCMSConfiguration(t *testing.T) {
	t.Parallel()
	// Test environment with a simulated EVM blockchain
	loader := testenv.NewLoader()
	env, err := loader.Load(t.Context(), testenv.WithEVMSimulatedN(t, 2))
	require.NoError(t, err)

	selectors := env.BlockChains.ListChainSelectors()
	require.Len(t, selectors, 2)
	// Create runtime instance
	r := runtime.NewFromEnvironment(*env)

	// Execute a changeset to deploy the MCMS contracts
	task := runtime.ChangesetTask(changeset.CsMCMSDeploy{}, changeset.DeployChangesetInput{
		ConfigID:       "test-config",
		Labels:         map[string]string{"foo": "bar"},
		ChainSelectors: selectors,
		ContractConfiguration: changeset.ContractConfiguration{
			Config: testMCMSCfg,
		},
	})

	err = r.Exec(task)
	require.NoError(t, err)

	// Verify deployment results - 5 mcms contracts per chain
	for _, selector := range selectors {
		addrs := r.State().DataStore.Addresses().Filter(datastore.AddressRefByChainSelector(selector))
		require.NoError(t, err)
		assert.Len(t, addrs, 5)
		chain, ok := env.BlockChains.EVMChains()[selector]
		require.True(t, ok, "chain selector %d not found", selector)
		for _, addr := range addrs {
			assert.Contains(t, addr.Labels.List(), "mcms_config=test-config")
			assert.Contains(t, addr.Labels.List(), "foo=bar")
			if strings.Contains(addr.Type.String(), "ManyChainMultiSig") {
				mcms, err2 := mcmsbindings.NewManyChainMultiSig(common.BytesToAddress(common.FromHex(addr.Address)), chain.Client)
				require.NoError(t, err2)
				config, err2 := mcms.GetConfig((&bind.CallOpts{Context: env.GetContext()}))
				require.NoError(t, err2)
				require.Equal(t, eoa, config.Signers[0].Addr)
			}
		}
	}

	// Execute a changeset to update the MCMS contracts
	task2 := runtime.ChangesetTask(changeset.CsMCMSConfigure{}, changeset.ConfigureChangesetInput{
		ChainSelectors: selectors,
		MCMSContractConfiguration: changeset.ContractConfiguration{
			Config: testMCMSCfg2,
		},
		// MCMSConfig: mcmsConfig, // TODO: Add MCMS once runtime.ChangesetTask supports them in Run()
	})

	err = r.Exec(task2)
	require.NoError(t, err)

	// Verify update configuration results - the signer should be changed to updatedEOA
	for _, selector := range selectors {
		addrs := r.State().DataStore.Addresses().Filter(datastore.AddressRefByChainSelector(selector))
		require.NoError(t, err)
		assert.Len(t, addrs, 5)
		chain, ok := env.BlockChains.EVMChains()[selector]
		require.True(t, ok, "chain selector %d not found", selector)
		for _, addr := range addrs {
			if strings.Contains(addr.Type.String(), "ManyChainMultiSig") {
				mcms, err := mcmsbindings.NewManyChainMultiSig(common.BytesToAddress(common.FromHex(addr.Address)), chain.Client)
				require.NoError(t, err)
				config, err := mcms.GetConfig((&bind.CallOpts{Context: env.GetContext()}))
				require.NoError(t, err)
				require.Equal(t, updatedEOA, config.Signers[0].Addr)
			}
		}
	}
}

var updatedEOA = common.HexToAddress("0xA01E9eD15b18D3688D0B84D88a98ed750D56999C")
var testMCMSCfg2 = types.MCMSWithTimelockConfigV2{
	Proposer:  pkg.MustGetMCMSConfig(1, []common.Address{updatedEOA}, nil),
	Bypasser:  pkg.MustGetMCMSConfig(1, []common.Address{updatedEOA}, nil),
	Canceller: pkg.MustGetMCMSConfig(1, []common.Address{updatedEOA}, nil),
}

// var mcmsConfig = proposalutils.TimelockConfig{
// 	MinDelay:   5 * time.Second,
// 	MCMSAction: mcmstypes.TimelockActionSchedule,
// }
