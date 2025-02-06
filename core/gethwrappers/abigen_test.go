package gethwrappers

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient/simulated"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-integrations/evm/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/log_emitter"
)

// Test that the generated Deploy method fill all the required fields and returns the correct address.
// We perform this test using the generated LogEmitter wrapper.
func TestGeneratedDeployMethodAddressField(t *testing.T) {
	owner := testutils.MustNewSimTransactor(t)
	ec := simulated.NewBackend(types.GenesisAlloc{
		owner.From: {
			Balance: big.NewInt(0).Mul(big.NewInt(10), big.NewInt(1e18)),
		},
	}, simulated.WithBlockGasLimit(10e6)).Client()

	emitterAddr, _, emitter, err := log_emitter.DeployLogEmitter(owner, ec)
	require.NoError(t, err)
	require.Equal(t, emitterAddr, emitter.Address())
}
