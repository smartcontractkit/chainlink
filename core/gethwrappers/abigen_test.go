package gethwrappers

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/log_emitter"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
)

// Test that the generated Deploy method fill all the required fields and returns the correct address.
// We perform this test using the generated LogEmitter wrapper.
func TestGeneratedDeployMethodAddressField(t *testing.T) {
	owner := testutils.MustNewSimTransactor(t)
	ec := backends.NewSimulatedBackendWithDatabase(rawdb.NewMemoryDatabase(), map[common.Address]core.GenesisAccount{
		owner.From: {
			Balance: big.NewInt(0).Mul(big.NewInt(10), big.NewInt(1e18)),
		},
	}, 10e6)
	emitterAddr, _, emitter, err := log_emitter.DeployLogEmitter(owner, ec)
	require.NoError(t, err)
	require.Equal(t, emitterAddr, emitter.Address())
}
