package types

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"
)

func TestTransaction_GetMeta(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		meta     *sqlutil.JSON
		expected *TxMeta
		wantErr  bool
	}{
		{
			name: "successful parse with all fields",
			meta: func() *sqlutil.JSON {
				meta := TxMeta{
					JobID:               ptr(int32(123)),
					FailOnRevert:        null.BoolFrom(true),
					RequestID:           &common.Hash{1},
					RequestTxHash:       &common.Hash{2},
					RequestIDs:          []common.Hash{{3}, {4}},
					RequestTxHashes:     []common.Hash{{5}, {6}},
					MaxLink:             ptr("1000000000000000000"),
					SubID:               ptr(uint64(5)),
					GlobalSubID:         ptr("abc123"),
					MaxEth:              ptr("2000000000000000000"),
					ForceFulfilled:      ptr(true),
					UpkeepID:            ptr("7890"),
					WorkflowExecutionID: ptr("workflow1"),
					FwdrDestAddress:     &common.Address{7},
					MessageIDs:          []string{"msg1", "msg2"},
					SeqNumbers:          []uint64{1, 2},
					DualBroadcast:       ptr(true),
					DualBroadcastParams: ptr("params123"),
				}
				b, err := json.Marshal(meta)
				require.NoError(t, err)
				j := sqlutil.JSON(b)
				return &j
			}(),
			expected: &TxMeta{
				JobID:               ptr(int32(123)),
				FailOnRevert:        null.BoolFrom(true),
				RequestID:           &common.Hash{1},
				RequestTxHash:       &common.Hash{2},
				RequestIDs:          []common.Hash{{3}, {4}},
				RequestTxHashes:     []common.Hash{{5}, {6}},
				MaxLink:             ptr("1000000000000000000"),
				SubID:               ptr(uint64(5)),
				GlobalSubID:         ptr("abc123"),
				MaxEth:              ptr("2000000000000000000"),
				ForceFulfilled:      ptr(true),
				UpkeepID:            ptr("7890"),
				WorkflowExecutionID: ptr("workflow1"),
				FwdrDestAddress:     &common.Address{7},
				MessageIDs:          []string{"msg1", "msg2"},
				SeqNumbers:          []uint64{1, 2},
				DualBroadcast:       ptr(true),
				DualBroadcastParams: ptr("params123"),
			},
			wantErr: false,
		},
		{
			name:     "nil meta returns nil",
			meta:     nil,
			expected: nil,
			wantErr:  false,
		},
		{
			name: "invalid json returns error",
			meta: func() *sqlutil.JSON {
				j := sqlutil.JSON([]byte(`{invalid json`))
				return &j
			}(),
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tx := Transaction{
				Meta: tt.meta,
			}

			got, err := tx.GetMeta()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestTransaction_String(t *testing.T) {
	t.Parallel()

	addr1 := common.HexToAddress("0x1234567890123456789012345678901234567890")
	addr2 := common.HexToAddress("0x2234567890123456789012345678901234567890")
	chainID := big.NewInt(1)
	nonce := uint64(42)
	value := big.NewInt(1000000000000000000)
	data := []byte{1, 2, 3, 4}
	gasLimit := uint64(21000)

	tx := Transaction{
		ID:                1,
		IdempotencyKey:    ptr("key123"),
		ChainID:           chainID,
		Nonce:             &nonce,
		FromAddress:       addr1,
		ToAddress:         addr2,
		Value:             value,
		Data:              data,
		SpecifiedGasLimit: gasLimit,
	}

	str := tx.String()
	assert.Contains(t, str, "txID:1")
	assert.Contains(t, str, "IdempotencyKey:key123")
	assert.Contains(t, str, "ChainID:1")
	assert.Contains(t, str, "Nonce:42")
	assert.Contains(t, str, addr1.String())
	assert.Contains(t, str, addr2.String())
	assert.Contains(t, str, value.String())
}

func ptr[T any](t T) *T { return &t }
