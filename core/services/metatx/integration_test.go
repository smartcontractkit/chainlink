package metatx_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core"
	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/meta_erc20"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/services/metatx"
	"github.com/stretchr/testify/require"
)

func TestMetaERC20(t *testing.T) {
	// Create a private key for holder1 that we can use to sign
	holder1Key := cltest.MustGenerateRandomKey(t)
	holder1Transactor, err := bind.NewKeyedTransactorWithChainID(holder1Key.ToEcdsaPrivKey(), testutils.SimulatedChainID)
	require.NoError(t, err)
	var (
		metaERC20Owner = testutils.MustNewSimTransactor(t)
		holder1        = holder1Transactor
		holder2        = testutils.MustNewSimTransactor(t)
		relay          = testutils.MustNewSimTransactor(t)
	)

	genesisData := core.GenesisAlloc{
		metaERC20Owner.From: {Balance: assets.Ether(1000).ToInt()},
		holder1.From:        {Balance: assets.Ether(1000).ToInt()},
		holder2.From:        {Balance: assets.Ether(1000).ToInt()},
		relay.From:          {Balance: assets.Ether(1000).ToInt()},
	}

	gasLimit := uint32(30e6)
	backend := cltest.NewSimulatedBackend(t, genesisData, gasLimit)

	_, _, metaWrapper, err := meta_erc20.DeployMetaERC20(metaERC20Owner, backend, assets.Ether(int64(1e18)).ToInt())
	require.NoError(t, err)

	backend.Commit()

	// transfer from owner to holder1
	_, err = metaWrapper.Transfer(metaERC20Owner, holder1.From, assets.Ether(1).ToInt())
	require.NoError(t, err)

	backend.Commit()

	holder1Bal, err := metaWrapper.BalanceOf(nil, holder1.From)
	require.NoError(t, err)
	require.Equal(t, assets.Ether(1).ToInt(), holder1Bal)

	// meta transfer from holder1 to holder2
	deadline := big.NewInt(int64(backend.Blockchain().CurrentHeader().Time + uint64(time.Hour)))
	v, r, s, err := metatx.SignMetaTransfer(
		metaWrapper,
		holder1Key.ToEcdsaPrivKey(),
		holder1.From,            // owner
		holder2.From,            // to
		assets.Ether(1).ToInt(), // amount
		deadline,
	)
	require.NoError(t, err)
	_, err = metaWrapper.MetaTransfer(relay, holder1.From, holder2.From, assets.Ether(1).ToInt(), deadline, v, r, s)
	require.NoError(t, err)

	backend.Commit()

	holder2Bal, err := metaWrapper.BalanceOf(nil, holder2.From)
	require.NoError(t, err)
	require.Equal(t, assets.Ether(1).ToInt(), holder2Bal)
}
