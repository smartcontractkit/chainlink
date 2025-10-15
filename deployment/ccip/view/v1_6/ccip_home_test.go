package v1_6

import (
	"encoding/json"
	"testing"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	capabilities_registry "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/ccip_home"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
)

func TestCCIPHomeView(t *testing.T) {
	t.Parallel()

	selector := chain_selectors.TEST_90000001.Selector
	e, err := environment.New(t.Context(),
		environment.WithEVMSimulated(t, []uint64{selector}),
		environment.WithLogger(logger.Test(t)),
	)
	require.NoError(t, err)

	chain := e.BlockChains.EVMChains()[selector]

	_, tx, cr, err := capabilities_registry.DeployCapabilitiesRegistry(
		chain.DeployerKey, chain.Client)
	require.NoError(t, err)
	_, err = cldf.ConfirmIfNoError(chain, tx, err)
	require.NoError(t, err)

	_, tx, ch, err := ccip_home.DeployCCIPHome(
		chain.DeployerKey, chain.Client, cr.Address())
	_, err = cldf.ConfirmIfNoError(chain, tx, err)
	require.NoError(t, err)

	v, err := GenerateCCIPHomeView(cr, ch)
	require.NoError(t, err)
	assert.Equal(t, "CCIPHome 1.6.0", v.TypeAndVersion)

	_, err = json.MarshalIndent(v, "", "  ")
	require.NoError(t, err)
}
