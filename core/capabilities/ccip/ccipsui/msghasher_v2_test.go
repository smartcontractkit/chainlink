package ccipsui

import (
	"context"
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	ccipocr3 "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMessageHasherV2_Deterministic_ParityWithMove(t *testing.T) {
	messageID := hexTo32Bytes(t, "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
	receiver := hexTo32Bytes(t, "0000000000000000000000000000000000000000000000000000000000001234")
	sender, err := hex.DecodeString("8765432109fedcba8765432109fedcba87654321")
	require.NoError(t, err)
	tokenReceiver := hexTo32Bytes(t, "0000000000000000000000000000000000000000000000000000000000005678")
	objectId := hexTo32Bytes(t, "0000000000000000000000000000000000000000000000000000000000aabbcc")

	metadataHash, err := computeMetadataHash(uint64(1000), uint64(2000), []byte("onramp"))
	require.NoError(t, err)

	hash1, err := computeMessageDataHashV2(
		metadataHash, messageID, receiver, uint64(1), big.NewInt(200000), tokenReceiver, uint64(0),
		sender, []byte("test payload"), []any2SuiTokenTransfer{}, [][32]byte{objectId},
	)
	require.NoError(t, err)

	hash2, err := computeMessageDataHashV2(
		metadataHash, messageID, receiver, uint64(1), big.NewInt(200000), tokenReceiver, uint64(0),
		sender, []byte("test payload"), []any2SuiTokenTransfer{}, [][32]byte{objectId},
	)
	require.NoError(t, err)

	assert.Equal(t, hash1, hash2)

	expectedHashHex := "1463b1b58f28f74dd73d4447da139d065051ddbb292549847a8c315d19148fc1"
	assert.Equal(t, expectedHashHex, hex.EncodeToString(hash1[:]))
}

func TestMessageHasherV2_DifferentObjectIds_DifferentHash(t *testing.T) {
	messageID := hexTo32Bytes(t, "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
	receiver := hexTo32Bytes(t, "0000000000000000000000000000000000000000000000000000000000001234")
	sender, err := hex.DecodeString("8765432109fedcba8765432109fedcba87654321")
	require.NoError(t, err)
	tokenReceiver := hexTo32Bytes(t, "0000000000000000000000000000000000000000000000000000000000000000")
	objectIdA := hexTo32Bytes(t, "0000000000000000000000000000000000000000000000000000000000001111")
	objectIdB := hexTo32Bytes(t, "0000000000000000000000000000000000000000000000000000000000002222")

	metadataHash, err := computeMetadataHash(uint64(123456789), uint64(987654321), []byte("source-onramp-address"))
	require.NoError(t, err)

	hashA, err := computeMessageDataHashV2(
		metadataHash, messageID, receiver, uint64(42), big.NewInt(500000), tokenReceiver, uint64(0),
		sender, []byte("sample message data"), []any2SuiTokenTransfer{}, [][32]byte{objectIdA},
	)
	require.NoError(t, err)

	hashB, err := computeMessageDataHashV2(
		metadataHash, messageID, receiver, uint64(42), big.NewInt(500000), tokenReceiver, uint64(0),
		sender, []byte("sample message data"), []any2SuiTokenTransfer{}, [][32]byte{objectIdB},
	)
	require.NoError(t, err)

	assert.NotEqual(t, hashA, hashB)
}

func TestExtractReceiverObjectIdsFromMap(t *testing.T) {
	idA := hexTo32Bytes(t, "0000000000000000000000000000000000000000000000000000000000001111")

	tests := []struct {
		name      string
		input     map[string]any
		expected  [][32]byte
		expectErr bool
	}{
		{
			name:     "missing key defaults to empty",
			input:    map[string]any{"gasLimit": big.NewInt(1)},
			expected: [][32]byte{},
		},
		{
			name:     "[][32]byte slice",
			input:    map[string]any{"receiverObjectIds": [][32]byte{idA}},
			expected: [][32]byte{idA},
		},
		{
			name:     "[][]byte slice",
			input:    map[string]any{"receiverObjectIds": [][]byte{idA[:]}},
			expected: [][32]byte{idA},
		},
		{
			name:      "invalid length",
			input:     map[string]any{"receiverObjectIds": [][]byte{{0x01, 0x02}}},
			expectErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := extractReceiverObjectIdsFromMap(tc.input)
			if tc.expectErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestMessageHasherV2_Hash_UsesReceiverObjectIdsFromExtraArgs(t *testing.T) {
	lggr := logger.Test(t)
	hasher := NewMessageHasherV2(lggr, mockExtraDataCodec{
		extraArgs: map[string]any{
			"gasLimit":          big.NewInt(200000),
			"tokenReceiver":     hexTo32Bytes(t, "0000000000000000000000000000000000000000000000000000000000005678"),
			"receiverObjectIds": [][32]byte{hexTo32Bytes(t, "0000000000000000000000000000000000000000000000000000000000aabbcc")},
		},
	})

	msg := ccipocr3.Message{
		Header: ccipocr3.RampMessageHeader{
			MessageID:           hexTo32Bytes(t, "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"),
			SourceChainSelector: 1000,
			DestChainSelector:   2000,
			SequenceNumber:      1,
			Nonce:               0,
			OnRamp:              []byte("onramp"),
		},
		Sender:   mustHexDecode(t, "8765432109fedcba8765432109fedcba87654321"),
		Data:     []byte("test payload"),
		Receiver: mustLeftPad32(t, "1234"),
	}

	hashWithObject, err := hasher.Hash(context.Background(), msg)
	require.NoError(t, err)

	hasherEmpty := NewMessageHasherV2(lggr, mockExtraDataCodec{
		extraArgs: map[string]any{
			"gasLimit":          big.NewInt(200000),
			"tokenReceiver":     hexTo32Bytes(t, "0000000000000000000000000000000000000000000000000000000000005678"),
			"receiverObjectIds": [][32]byte{},
		},
	})
	hashEmpty, err := hasherEmpty.Hash(context.Background(), msg)
	require.NoError(t, err)

	assert.NotEqual(t, hashWithObject, hashEmpty)
	assert.Equal(t, "0x1463b1b58f28f74dd73d4447da139d065051ddbb292549847a8c315d19148fc1", hexutil.Encode(hashWithObject[:]))
}

type mockExtraDataCodec struct {
	extraArgs map[string]any
}

func (m mockExtraDataCodec) DecodeExtraArgs(_ ccipocr3.Bytes, _ ccipocr3.ChainSelector) (map[string]any, error) {
	return m.extraArgs, nil
}

func (m mockExtraDataCodec) DecodeTokenAmountDestExecData(_ ccipocr3.Bytes, _ ccipocr3.ChainSelector) (map[string]any, error) {
	return map[string]any{}, nil
}

func hexTo32Bytes(t *testing.T, hexStr string) [32]byte {
	t.Helper()
	bytes, err := hex.DecodeString(hexStr)
	require.NoError(t, err)
	require.Len(t, bytes, 32)

	var result [32]byte
	copy(result[:], bytes)
	return result
}

func mustHexDecode(t *testing.T, s string) []byte {
	t.Helper()
	b, err := hex.DecodeString(s)
	require.NoError(t, err)
	return b
}

func mustLeftPad32(t *testing.T, hexSuffix string) ccipocr3.UnknownAddress {
	t.Helper()
	b, err := hex.DecodeString(hexSuffix)
	require.NoError(t, err)
	var addr [32]byte
	copy(addr[32-len(b):], b)
	return ccipocr3.UnknownAddress(addr[:])
}
