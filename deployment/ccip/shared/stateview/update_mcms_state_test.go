package stateview

import (
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"

	"github.com/smartcontractkit/chainlink/deployment"
	evmstate "github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/evm"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
)

func TestUpdateMCMSStateWithAddressFromDatastoreForChain(t *testing.T) {
	selector := chain_selectors.TEST_90000001.Selector
	env, err := environment.New(t.Context(),
		environment.WithEVMSimulated(t, []uint64{selector}),
	)
	require.NoError(t, err)

	// CLLCCIP addresses
	cllccipTimelock := common.HexToAddress("0x0000000000000000000000000000000000100001")
	cllccipCallProxy := common.HexToAddress("0x0000000000000000000000000000000000100002")
	cllccipProposer := common.HexToAddress("0x0000000000000000000000000000000000100003")
	cllccipBypasser := common.HexToAddress("0x0000000000000000000000000000000000100004")
	cllccipCanceller := common.HexToAddress("0x0000000000000000000000000000000000100005")

	// RMNMCMS addresses
	rmnTimelock := common.HexToAddress("0x0000000000000000000000000000000000200001")
	rmnCallProxy := common.HexToAddress("0x0000000000000000000000000000000000200002")
	rmnProposer := common.HexToAddress("0x0000000000000000000000000000000000200003")
	rmnBypasser := common.HexToAddress("0x0000000000000000000000000000000000200004")
	rmnCanceller := common.HexToAddress("0x0000000000000000000000000000000000200005")

	// Populate DataStore with both MCMS with different qualifiers
	store := datastore.NewMemoryDataStore()
	addBundle := func(qualifier string, timelock, callProxy, proposer, bypasser, canceller common.Address) {
		for _, ref := range []datastore.AddressRef{
			{ChainSelector: selector, Address: timelock.Hex(), Type: datastore.ContractType(types.RBACTimelock), Version: &deployment.Version1_0_0, Qualifier: qualifier},
			{ChainSelector: selector, Address: callProxy.Hex(), Type: datastore.ContractType(types.CallProxy), Version: &deployment.Version1_0_0, Qualifier: qualifier},
			{ChainSelector: selector, Address: proposer.Hex(), Type: datastore.ContractType(types.ProposerManyChainMultisig), Version: &deployment.Version1_0_0, Qualifier: qualifier},
			{ChainSelector: selector, Address: bypasser.Hex(), Type: datastore.ContractType(types.BypasserManyChainMultisig), Version: &deployment.Version1_0_0, Qualifier: qualifier},
			{ChainSelector: selector, Address: canceller.Hex(), Type: datastore.ContractType(types.CancellerManyChainMultisig), Version: &deployment.Version1_0_0, Qualifier: qualifier},
		} {
			require.NoError(t, store.Addresses().Add(ref))
		}
	}
	addBundle("CLLCCIP", cllccipTimelock, cllccipCallProxy, cllccipProposer, cllccipBypasser, cllccipCanceller)
	addBundle("RMNMCMS", rmnTimelock, rmnCallProxy, rmnProposer, rmnBypasser, rmnCanceller)

	state := CCIPOnChainState{
		Chains: map[uint64]evmstate.CCIPChainState{
			selector: {ABIByAddress: make(map[string]string)},
		},
		evmMu: &sync.RWMutex{},
	}

	cldfEnv := cldf.Environment{
		BlockChains: env.BlockChains,
		DataStore:   store.Seal(),
	}

	// state should contain the CLLCCIP bundle
	require.NoError(t, state.UpdateMCMSStateWithAddressFromDatastoreForChain(cldfEnv, selector, "CLLCCIP"))

	chainState, ok := state.EVMChainState(selector)
	require.True(t, ok)
	require.Equal(t, cllccipTimelock, chainState.Timelock.Address())
	require.Equal(t, cllccipCallProxy, chainState.CallProxy.Address())
	require.Equal(t, cllccipProposer, chainState.ProposerMcm.Address())
	require.Equal(t, cllccipBypasser, chainState.BypasserMcm.Address())
	require.Equal(t, cllccipCanceller, chainState.CancellerMcm.Address())

	// state should now contain the RMN bundle
	require.NoError(t, state.UpdateMCMSStateWithAddressFromDatastoreForChain(cldfEnv, selector, "RMNMCMS"))

	chainState, ok = state.EVMChainState(selector)
	require.True(t, ok)
	require.Equal(t, rmnTimelock, chainState.Timelock.Address())
	require.Equal(t, rmnCallProxy, chainState.CallProxy.Address())
	require.Equal(t, rmnProposer, chainState.ProposerMcm.Address())
	require.Equal(t, rmnBypasser, chainState.BypasserMcm.Address())
	require.Equal(t, rmnCanceller, chainState.CancellerMcm.Address())

	// the contracts for each qualifier should be different
	require.NotEqual(t, cllccipTimelock, rmnTimelock, "qualifiers should isolate timelocks")
	require.NotEqual(t, cllccipProposer, rmnProposer, "qualifiers should isolate proposers")
}
