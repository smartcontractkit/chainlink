package changeset

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	cldfproposalutils "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"

	"github.com/smartcontractkit/chainlink/deployment/vault/changeset/types"
)

const (
	bridgeSrcSel  = uint64(5009297550715157269)                  // chain selector ethereum-mainnet
	bridgeDstSel  = uint64(3734403246176062136)                  // chain selector optimism-mainnet
	bridgeDstSel2 = uint64(15971525489660198786)                 // chain selector base-mainnet
	bridgeRecip   = "0x66598216D8E4d9AFE0F06d525B335b762229842f" // EtherSenderReceiver contract address
)

// --- fee buffer logic ---

func TestBufferedFee(t *testing.T) {
	t.Parallel()

	// Default buffer is 150% (15000 bps).
	require.Equal(t, big.NewInt(15_000), CCIPFeeBufferBps, "test assumes default 150% buffer")

	cases := []struct {
		fee  *big.Int
		want *big.Int
	}{
		{big.NewInt(0), big.NewInt(0)},
		{big.NewInt(100), big.NewInt(150)},
		{big.NewInt(1_000_000_000_000_000), big.NewInt(1_500_000_000_000_000)},
	}
	for _, c := range cases {
		require.Equal(t, c.want, bufferedFee(c.fee), "bufferedFee(%s)", c.fee)
	}
}

// --- test/remainder split logic ---

func TestSplitTransfersForTest(t *testing.T) {
	t.Parallel()

	full := big.NewInt(0).Add(TestBridgeAmount, big.NewInt(42)) // test + 42
	transfers := []types.CCIPBridgeTransfer{
		{ToChainSelector: bridgeDstSel, Amount: full, DestinationAddress: bridgeRecip},
	}

	test, remainder, err := splitTransfersForTest(transfers)
	require.NoError(t, err)
	require.Len(t, test, 1)
	require.Len(t, remainder, 1)

	require.Equal(t, TestBridgeAmount, test[0].Amount, "test amount should equal TestBridgeAmount")
	require.Equal(t, big.NewInt(42), remainder[0].Amount, "remainder = amount - test")

	// test + remainder == original
	sum := big.NewInt(0).Add(test[0].Amount, remainder[0].Amount)
	require.Equal(t, full, sum, "test + remainder must equal original amount")
	require.Equal(t, full, transfers[0].Amount, "original transfer slice must be untouched")

	// metadata
	require.Equal(t, bridgeDstSel, remainder[0].ToChainSelector)
	require.Equal(t, bridgeRecip, remainder[0].DestinationAddress)
}

func TestSplitTransfersForTest_AmountTooSmall(t *testing.T) {
	t.Parallel()

	for _, amt := range []*big.Int{
		big.NewInt(0).Set(TestBridgeAmount),                // equal -> remainder would be 0
		big.NewInt(0).Sub(TestBridgeAmount, big.NewInt(1)), // below test amount
	} {
		_, _, err := splitTransfersForTest([]types.CCIPBridgeTransfer{
			{ToChainSelector: bridgeDstSel, Amount: amt, DestinationAddress: bridgeRecip},
		})
		require.Error(t, err, "amount %s must be rejected", amt)
		require.Contains(t, err.Error(), "greater than test amount")
	}
}

// --- CCIP extraArgs v2 logic ---

func TestBuildExtraArgsV2(t *testing.T) {
	t.Parallel()

	got, err := buildExtraArgsV2(100_000)
	require.NoError(t, err)

	tag := common.FromHex(genericExtraArgsV2Tag)
	require.Equal(t, tag, got[:len(tag)], "must start with GENERIC_EXTRA_ARGS_V2_TAG")

	vals, err := extraArgsV2Args.Unpack(got[len(tag):])
	require.NoError(t, err)
	require.Len(t, vals, 2)
	require.Equal(t, big.NewInt(100_000), vals[0].(*big.Int))
	require.True(t, vals[1].(bool), "allowOutOfOrderExecution must be true")
}

// --- CCIP message construction logic ---

func TestBuildCCIPMessage(t *testing.T) {
	t.Parallel()

	weth := common.HexToAddress("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2")
	esr := common.HexToAddress(bridgeRecip)
	recipient := common.HexToAddress("0x1111111111111111111111111111111111111111")
	amount := big.NewInt(1_234)

	msg, err := buildCCIPMessage(weth, esr, recipient, amount, 100_000)
	require.NoError(t, err)

	// token = WETH, amount preserved
	require.Len(t, msg.TokenAmounts, 1)
	require.Equal(t, weth, msg.TokenAmounts[0].Token)
	require.Equal(t, amount, msg.TokenAmounts[0].Amount)

	// native fee token
	require.Equal(t, common.Address{}, msg.FeeToken, "feeToken must be address(0) for native fee")

	// receiver decodes to the destination EtherSenderReceiver
	addrTy, err := abi.NewType("address", "", nil)
	require.NoError(t, err)
	args := abi.Arguments{{Type: addrTy}}

	rcv, err := args.Unpack(msg.Receiver)
	require.NoError(t, err)
	require.Equal(t, esr, rcv[0].(common.Address), "receiver must abi-encode the EtherSenderReceiver")

	// data decodes to the final recipient
	dat, err := args.Unpack(msg.Data)
	require.NoError(t, err)
	require.Equal(t, recipient, dat[0].(common.Address), "data must abi-encode the final recipient")

	require.Equal(t, common.FromHex(genericExtraArgsV2Tag), msg.ExtraArgs[:4])
}

// --- static validation check ---

func TestValidateCCIPBridgeConfig_StaticErrors(t *testing.T) {
	t.Parallel()

	env, err := environment.New(t.Context())
	require.NoError(t, err)

	validTransfer := types.CCIPBridgeTransfer{
		ToChainSelector:    bridgeDstSel,
		Amount:             OneETH,
		DestinationAddress: bridgeRecip,
	}
	mcms := &cldfproposalutils.TimelockConfig{MinDelay: 0}

	cases := []struct {
		name   string
		cfg    types.CCIPBridgeConfig
		errSub string
	}{
		{
			name: "invalid token type",
			cfg: types.CCIPBridgeConfig{
				FromChainSelector: bridgeSrcSel,
				TokenType:         "weird",
				Transfers:         []types.CCIPBridgeTransfer{validTransfer},
				MCMSConfig:        mcms,
			},
			errSub: "token_type must be",
		},
		{
			name: "empty transfers",
			cfg: types.CCIPBridgeConfig{
				FromChainSelector: bridgeSrcSel,
				TokenType:         types.TokenTypeNative,
				Transfers:         nil,
				MCMSConfig:        mcms,
			},
			errSub: "at least one transfer",
		},
		{
			name: "nil mcms config",
			cfg: types.CCIPBridgeConfig{
				FromChainSelector: bridgeSrcSel,
				TokenType:         types.TokenTypeNative,
				Transfers:         []types.CCIPBridgeTransfer{validTransfer},
				MCMSConfig:        nil,
			},
			errSub: "MCMSConfig is required",
		},
		{
			name: "source chain not in environment",
			cfg: types.CCIPBridgeConfig{
				FromChainSelector: bridgeSrcSel,
				TokenType:         types.TokenTypeNative,
				Transfers:         []types.CCIPBridgeTransfer{validTransfer},
				MCMSConfig:        mcms,
			},
			errSub: "source chain",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			err := ValidateCCIPBridgeConfig(env.GetContext(), *env, c.cfg)
			require.Error(t, err)
			require.Contains(t, err.Error(), c.errSub)
		})
	}
}
