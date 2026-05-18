package core

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"

	ocr2keepers "github.com/smartcontractkit/chainlink-common/pkg/types/automation"
)

func TestPackUnpackTrigger(t *testing.T) {
	tests := []struct {
		name    string
		id      []byte
		trigger triggerWrapper
		encoded []byte
		err     error
	}{
		{
			"happy flow log trigger",
			append([]byte{1}, common.LeftPadBytes([]byte{1}, 15)...),
			triggerWrapper{
				BlockNum:     1,
				BlockHash:    common.HexToHash("0x01111111"),
				LogIndex:     1,
				TxHash:       common.HexToHash("0x01111111"),
				LogBlockHash: common.HexToHash("0x01111abc"),
			},
			func() []byte {
				b, _ := hexutil.Decode("0x0000000000000000000000000000000000000000000000000000000001111abc0000000000000000000000000000000000000000000000000000000001111111000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000001111111")
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
			var idBytes [32]byte
			copy(idBytes[:], tc.id)
			id := ocr2keepers.UpkeepIdentifier(idBytes)

			encoded, err := PackTrigger(id.BigInt(), tc.trigger)
			if tc.err != nil {
				assert.EqualError(t, err, tc.err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.encoded, encoded)
				decoded, err := UnpackTrigger(id.BigInt(), encoded)
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
		var idBytes [32]byte
		copy(idBytes[:], uid)
		id := ocr2keepers.UpkeepIdentifier(idBytes)
		decoded, _ := hexutil.Decode("0x00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000001111111")
		_, err := UnpackTrigger(id.BigInt(), decoded)
		assert.EqualError(t, err, "unknown trigger type: 8")
	})
}
