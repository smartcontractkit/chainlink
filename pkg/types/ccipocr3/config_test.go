package ccipocr3

import (
	"testing"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/assert"
)

func TestCommitPluginConfigValidate(t *testing.T) {
	testCases := []struct {
		name   string
		input  CommitPluginConfig
		expErr bool
	}{
		{
			name: "valid cfg",
			input: CommitPluginConfig{
				DestChain: ChainSelector(1),
				FChain: map[ChainSelector]int{
					ChainSelector(1): 1,
				},
				ObserverInfo: map[commontypes.OracleID]ObserverInfo{
					commontypes.OracleID(1): {
						Reads: []ChainSelector{ChainSelector(1)},
					},
				},
				PricedTokens: []types.Account{
					types.Account("0x123"),
					types.Account("0x124"),
				},
				NewMsgScanBatchSize: 256,
				TokenPricesObserver: true,
			},
			expErr: false,
		},
		{
			name: "dest chain is empty",
			input: CommitPluginConfig{
				FChain: map[ChainSelector]int{
					ChainSelector(1): 1,
				},
				ObserverInfo: map[commontypes.OracleID]ObserverInfo{
					commontypes.OracleID(1): {
						Reads: []ChainSelector{ChainSelector(1)},
					},
				},
				PricedTokens: []types.Account{
					types.Account("0x123"),
					types.Account("0x124"),
				},
				NewMsgScanBatchSize: 256,
				TokenPricesObserver: true,
			},
			expErr: true,
		},
		{
			name: "zero priced tokens",
			input: CommitPluginConfig{
				DestChain: ChainSelector(1),
				FChain: map[ChainSelector]int{
					ChainSelector(1): 1,
				},
				ObserverInfo: map[commontypes.OracleID]ObserverInfo{
					commontypes.OracleID(1): {
						Reads: []ChainSelector{ChainSelector(1)},
					},
				},
				NewMsgScanBatchSize: 256,
				TokenPricesObserver: true,
			},
			expErr: true,
		},
		{
			name: "empty batch scan size",
			input: CommitPluginConfig{
				DestChain: ChainSelector(1),
				FChain: map[ChainSelector]int{
					ChainSelector(1): 1,
				},
				ObserverInfo: map[commontypes.OracleID]ObserverInfo{
					commontypes.OracleID(1): {
						Reads: []ChainSelector{ChainSelector(1)},
					},
				},
				PricedTokens: []types.Account{
					types.Account("0x123"),
					types.Account("0x124"),
				},
				TokenPricesObserver: true,
			},
			expErr: true,
		},
		{
			name: "fChain not set for dest",
			input: CommitPluginConfig{
				DestChain: ChainSelector(1),
				FChain: map[ChainSelector]int{
					ChainSelector(2): 1,
				},
				ObserverInfo: map[commontypes.OracleID]ObserverInfo{
					commontypes.OracleID(1): {
						Reads: []ChainSelector{ChainSelector(1)},
					},
				},
				PricedTokens: []types.Account{
					types.Account("0x123"),
					types.Account("0x124"),
				},
				NewMsgScanBatchSize: 256,
				TokenPricesObserver: true,
			},
			expErr: true,
		},
		{
			name: "fChain not set for some chain",
			input: CommitPluginConfig{
				DestChain: ChainSelector(1),
				FChain: map[ChainSelector]int{
					ChainSelector(1): 1,
				},
				ObserverInfo: map[commontypes.OracleID]ObserverInfo{
					commontypes.OracleID(1): {
						Reads: []ChainSelector{ChainSelector(1), ChainSelector(123)}, // fChain not set for 123
					},
				},
				PricedTokens: []types.Account{
					types.Account("0x123"),
					types.Account("0x124"),
				},
				NewMsgScanBatchSize: 256,
				TokenPricesObserver: true,
			},
			expErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.input.Validate()
			if tc.expErr {
				assert.Error(t, actual)
				return
			}
			assert.NoError(t, actual)
		})
	}
}
