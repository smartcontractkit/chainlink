package logpoller

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
)

func Test_QueryArgs(t *testing.T) {
	tests := []struct {
		name      string
		queryArgs *queryArgs
		want      map[string]interface{}
		wantErr   bool
	}{
		{
			name:      "valid arguments",
			queryArgs: newQueryArgs(big.NewInt(20)).withAddress(utils.ZeroAddress),
			want: map[string]interface{}{
				"evm_chain_id": ubig.NewI(20),
				"address":      utils.ZeroAddress,
			},
		},
		{
			name:      "invalid topic index",
			queryArgs: newQueryArgs(big.NewInt(20)).withTopicIndex(0),
			wantErr:   true,
		},
		{
			name:      "custom argument",
			queryArgs: newEmptyArgs().withField("arg", "value"),
			want: map[string]interface{}{
				"arg": "value",
			},
		},
		{
			name:      "hash converted to bytes",
			queryArgs: newEmptyArgs().withField("hash", common.Hash{}),
			want: map[string]interface{}{
				"hash": make([]byte, 32),
			},
		},
		{
			name:      "hash array converted to bytes array",
			queryArgs: newEmptyArgs().withEventSigArray([]common.Hash{{}, {}}),
			want: map[string]interface{}{
				"event_sig_array": [][]byte{make([]byte, 32), make([]byte, 32)},
			},
		},
		{
			name:      "topic index incremented",
			queryArgs: newEmptyArgs().withTopicIndex(2),
			want: map[string]interface{}{
				"topic_index": 3,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args, err := tt.queryArgs.toArgs()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, args)
			}
		})
	}
}

func newEmptyArgs() *queryArgs {
	return &queryArgs{
		args: map[string]interface{}{},
		err:  []error{},
	}
}
