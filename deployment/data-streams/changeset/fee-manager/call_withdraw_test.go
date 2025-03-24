package fee_manager

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"

	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/testutil"
)

func TestWithdraw(t *testing.T) {
	res, err := NewDataStreamsEnvironment(t, NewDefaultOptions())
	require.NoError(t, err)

	linkTokenAddress := res.LinkTokenAddress
	feeManagerAddress := res.FeeManagerAddress
	e := res.Env

	chain := e.Chains[testutil.TestChain.Selector]
	require.NotNil(t, chain)

	recipient := common.HexToAddress("0x0fd8b81e3d1143ec7f1ce474827ab93c43523ea2")
	commonChangesets.MustFundAddressWithLink(t, e, chain, feeManagerAddress, 100)

	e, err = commonChangesets.Apply(t, e, nil,
		commonChangesets.Configure(
			WithdrawChangeset,
			FeeManagerWithdrawConfig{
				ConfigPerChain: map[uint64][]Withdraw{
					testutil.TestChain.Selector: {
						{FeeManagerAddress: feeManagerAddress,
							AssetAddress:     linkTokenAddress,
							RecipientAddress: recipient,
							Quantity:         big.NewInt(50)},
					},
				},
			},
		),
	)
	require.NoError(t, err)

	recipientBalance := commonChangesets.MaybeGetLinkBalance(t, e, chain, recipient)
	require.Equal(t, big.NewInt(50), recipientBalance)

	feeManagerBalance := commonChangesets.MaybeGetLinkBalance(t, e, chain, feeManagerAddress)
	require.Equal(t, big.NewInt(50), feeManagerBalance)
}
