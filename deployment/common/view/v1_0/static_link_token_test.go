package v1_0

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/link_token_interface"
)

func TestStaticLinkTokenView(t *testing.T) {
	selector := chain_selectors.TEST_90000001.Selector
	env, err := environment.New(t.Context(),
		environment.WithEVMSimulated(t, []uint64{selector}),
	)
	require.NoError(t, err)

	chain := env.BlockChains.EVMChains()[selector]
	_, tx, lt, err := link_token_interface.DeployLinkToken(chain.DeployerKey, chain.Client)
	require.NoError(t, err)
	_, err = chain.Confirm(tx)
	require.NoError(t, err)
	v, err := GenerateStaticLinkTokenView(lt)
	require.NoError(t, err)

	assert.Equal(t, v.Owner, common.HexToAddress("0x0")) // Ownerless
	assert.Equal(t, "StaticLinkToken 1.0.0", v.TypeAndVersion)
	assert.Equal(t, uint8(18), v.Decimals)
	assert.Equal(t, "1000000000000000000000000000", v.Supply.String())
}
