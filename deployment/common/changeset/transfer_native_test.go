package changeset

import (
	"math/big"
	"testing"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"

	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
)

func Test_TransferFunds_VerifyPreconditions(t *testing.T) {
	selector := chain_selectors.TEST_90000001.Selector
	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithEVMSimulated(t, []uint64{selector}),
	))
	require.NoError(t, err)

	testCases := []struct {
		name    string
		input   TransferNativeInput
		wantErr bool
	}{
		{
			name: "valid transfer",
			input: TransferNativeInput{
				ChainSel: chain_selectors.TEST_90000001.Selector,
				Address:  "0x1234567890123456789012345678901234567890",
				Amount:   big.NewInt(1000000000000000000), // 1 ETH
			},
			wantErr: false,
		},
		{
			name: "empty address",
			input: TransferNativeInput{
				ChainSel: chain_selectors.TEST_90000001.Selector,
				Address:  "",
				Amount:   big.NewInt(1000000000000000000),
			},
			wantErr: true,
		},
		{
			name: "invalid address format",
			input: TransferNativeInput{
				ChainSel: chain_selectors.TEST_90000001.Selector,
				Address:  "not-a-valid-address",
				Amount:   big.NewInt(1000000000000000000),
			},
			wantErr: true,
		},
		{
			name: "zero amount",
			input: TransferNativeInput{
				ChainSel: chain_selectors.TEST_90000001.Selector,
				Address:  "0x1234567890123456789012345678901234567890",
				Amount:   big.NewInt(0),
			},
			wantErr: true,
		},
		{
			name: "invalid chain selector",
			input: TransferNativeInput{
				ChainSel: 0, // Invalid chain selector
				Address:  "0x1234567890123456789012345678901234567890",
				Amount:   big.NewInt(100),
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tf := TransferNative{}
			err := tf.VerifyPreconditions(rt.Environment(), tc.input)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_TransferFundsChangeset(t *testing.T) {
	selector := chain_selectors.TEST_90000001.Selector
	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithEVMSimulated(t, []uint64{selector}),
	))
	require.NoError(t, err)

	tf := TransferNative{}

	t.Run("happy path", func(t *testing.T) {
		addr := utils.RandomAddress()
		transferVal := big.NewInt(1_000_000_000) // transfer 1gwei

		input := TransferNativeInput{
			ChainSel: chain_selectors.TEST_90000001.Selector,
			Address:  addr.Hex(),
			Amount:   transferVal,
		}

		err = tf.VerifyPreconditions(rt.Environment(), input)
		require.NoError(t, err)

		_, err = tf.Apply(rt.Environment(), input)
		require.NoError(t, err)

		chain, ok := rt.Environment().BlockChains.EVMChains()[input.ChainSel]
		require.True(t, ok)

		bal, err := chain.Client.BalanceAt(t.Context(), addr, nil)
		require.NoError(t, err)
		require.Equal(t, transferVal, bal)
	})

	t.Run("insufficient funds", func(t *testing.T) {
		chain, ok := rt.Environment().BlockChains.EVMChains()[chain_selectors.TEST_90000001.Selector]
		require.True(t, ok)
		bal, err := chain.Client.BalanceAt(t.Context(), chain.DeployerKey.From, nil)
		require.NoError(t, err)

		addr := utils.RandomAddress()
		transferVal := bal // transfer entire balance, leaving no funds for gas

		input := TransferNativeInput{
			ChainSel: chain_selectors.TEST_90000001.Selector,
			Address:  addr.Hex(),
			Amount:   transferVal,
		}

		err = tf.VerifyPreconditions(rt.Environment(), input)
		require.NoError(t, err)

		_, err = tf.Apply(rt.Environment(), input)
		require.Error(t, err)
	})
}
