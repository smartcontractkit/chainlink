package v1_0

import (
	"math/big"
	"testing"

	chainselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/initial/link_token"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
)

func TestLinkTokenView(t *testing.T) {
	selector := chainselectors.TEST_90000001.Selector
	env, err := environment.New(t.Context(),
		environment.WithEVMSimulated(t, []uint64{selector}),
	)
	require.NoError(t, err)

	chain := env.BlockChains.EVMChains()[selector]
	_, tx, lt, err := link_token.DeployLinkToken(chain.DeployerKey, chain.Client)
	require.NoError(t, err)
	_, err = chain.Confirm(tx)
	require.NoError(t, err)

	testLinkTokenViewWithChain(t, chain, lt)
}

func TestLinkTokenViewZk(t *testing.T) {
	// Timeouts in CI
	tests.SkipFlakey(t, "https://smartcontract-it.atlassian.net/browse/CCIP-6427")

	selector := chainselectors.TEST_90000050.Selector
	env, err := environment.New(t.Context(),
		environment.WithZKSyncContainer(t, []uint64{selector}),
	)
	require.NoError(t, err)

	chain := env.BlockChains.EVMChains()[selector]
	_, _, lt, err := link_token.DeployLinkTokenZk(nil, chain.ClientZkSyncVM, chain.DeployerKeyZkSyncVM, chain.Client)
	require.NoError(t, err)

	testLinkTokenViewWithChain(t, chain, lt)
}

func testLinkTokenViewWithChain(t *testing.T, chain cldf_evm.Chain, lt *link_token.LinkToken) {
	v, err := GenerateLinkTokenView(lt)
	require.NoError(t, err)

	assert.Equal(t, v.Owner, chain.DeployerKey.From)
	assert.Equal(t, "LinkToken 1.0.0", v.TypeAndVersion)
	assert.Equal(t, uint8(18), v.Decimals)
	// Initially nothing minted and no minters/burners.
	assert.Equal(t, "0", v.Supply.String())
	require.Empty(t, v.Minters)
	require.Empty(t, v.Burners)

	// Add some minters
	tx, err := lt.GrantMintAndBurnRoles(chain.DeployerKey, chain.DeployerKey.From)
	require.NoError(t, err)
	_, err = chain.Confirm(tx)
	require.NoError(t, err)
	tx, err = lt.Mint(chain.DeployerKey, chain.DeployerKey.From, big.NewInt(100))
	require.NoError(t, err)
	_, err = chain.Confirm(tx)
	require.NoError(t, err)

	v, err = GenerateLinkTokenView(lt)
	require.NoError(t, err)

	assert.Equal(t, "100", v.Supply.String())
	require.Len(t, v.Minters, 1)
	require.Equal(t, v.Minters[0].String(), chain.DeployerKey.From.String())
	require.Len(t, v.Burners, 1)
	require.Equal(t, v.Burners[0].String(), chain.DeployerKey.From.String())
}
