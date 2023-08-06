package core

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg"
	"github.com/test-go/testify/assert"
)

func TestPackUnpackTrigger(t *testing.T) {
	tests := []struct {
		name    string
		id      ocr2keepers.UpkeepIdentifier
		trigger triggerWrapper
		encoded []byte
		err     error
	}{
		{
			"happy flow log trigger",
			append([]byte{1}, common.LeftPadBytes([]byte{1}, 15)...),
			triggerWrapper{
				BlockNum:  1,
				BlockHash: common.HexToHash("0x01111111"),
				LogIndex:  1,
				TxHash:    common.HexToHash("0x01111111"),
			},
			func() []byte {
				b, _ := hexutil.Decode("0x0000000000000000000000000000000000000000000000000000000001111111000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000001111111")
				return b
			}(),
			nil,
		},
		{
			"happy flow conditional trigger",
			append([]byte{1}, common.LeftPadBytes([]byte{0}, 15)...),
			triggerWrapper{
				BlockNum:  1,
				BlockHash: common.HexToHash("0x01111111"),
			},
			func() []byte {
				b, _ := hexutil.Decode("0x00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000001111111")
				return b
			}(),
			nil,
		},
		{
			"invalid type",
			append([]byte{1}, common.LeftPadBytes([]byte{8}, 15)...),
			triggerWrapper{
				BlockNum:  1,
				BlockHash: common.HexToHash("0x01111111"),
			},
			[]byte{},
			fmt.Errorf("unknown trigger type: %d", 8),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			id, ok := big.NewInt(0).SetString(hexutil.Encode(tc.id)[2:], 16)
			assert.True(t, ok)

			encoded, err := PackTrigger(id, tc.trigger)
			if tc.err != nil {
				assert.EqualError(t, err, tc.err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.encoded, encoded)
				decoded, err := UnpackTrigger(id, encoded)
				assert.NoError(t, err)
				assert.Equal(t, tc.trigger.BlockNum, decoded.BlockNum)
			}
		})
	}

	t.Run("unpacking invalid trigger", func(t *testing.T) {
		_, err := UnpackTrigger(big.NewInt(0), []byte{1, 2, 3})
		assert.Error(t, err)
	})

	t.Run("unpacking unknown type", func(t *testing.T) {
		uid := append([]byte{1}, common.LeftPadBytes([]byte{8}, 15)...)
		id, ok := big.NewInt(0).SetString(hexutil.Encode(uid)[2:], 16)
		assert.True(t, ok)
		decoded, _ := hexutil.Decode("0x00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000001111111")
		_, err := UnpackTrigger(id, decoded)
		assert.EqualError(t, err, "unknown trigger type: 8")
	})
}
