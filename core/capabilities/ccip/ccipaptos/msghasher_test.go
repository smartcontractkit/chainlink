package ccipaptos

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/require"
)

// the equivalent test with the same values exists in the Aptos offramp.move contract, see test_calculate_message_hash
func TestComputeMessageDataHash(t *testing.T) {
	expectedHashStr := "0xc8d6cf666864a60dd6ecd89e5c294734c53b3218d3f83d2d19a3c3f9e200e00d"

	metadataHashBytes, err := hexutil.Decode("0xaabbccddeeff00112233445566778899aabbccddeeff00112233445566778899")
	require.NoError(t, err)
	var metadataHash [32]byte
	copy(metadataHash[:], metadataHashBytes)

	messageIDBytes, err := hexutil.Decode("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
	require.NoError(t, err)
	var messageID [32]byte
	copy(messageID[:], messageIDBytes)

	receiverBytes, err := hexutil.Decode("0x0000000000000000000000000000000000000000000000000000000000001234")
	require.NoError(t, err)
	var receiver [32]byte
	copy(receiver[:], receiverBytes)

	sequenceNumber := uint64(42)
	nonce := uint64(123)
	gasLimit := big.NewInt(500000)

	sender, err := hexutil.Decode("0x8765432109fedcba8765432109fedcba87654321")
	require.NoError(t, err)

	data := []byte("sample message data")
	srcPool1, err := hexutil.Decode("0xabcdef1234567890abcdef1234567890abcdef12")
	require.NoError(t, err)
	destToken1Bytes, err := hexutil.Decode("0x0000000000000000000000000000000000000000000000000000000000005678")
	require.NoError(t, err)
	var destToken1 [32]byte
	copy(destToken1[:], destToken1Bytes)
	extraData1, err := hexutil.Decode("0x00112233")
	require.NoError(t, err)
	token1 := any2AptosTokenTransfer{
		SourcePoolAddress: srcPool1,
		DestTokenAddress:  destToken1,
		DestGasAmount:     10000,
		ExtraData:         extraData1,
		Amount:            big.NewInt(1000000),
	}
	srcPool2, err := hexutil.Decode("0x123456789abcdef123456789abcdef123456789a")
	require.NoError(t, err)
	destToken2Bytes, err := hexutil.Decode("0x0000000000000000000000000000000000000000000000000000000000009abc")
	require.NoError(t, err)
	var destToken2 [32]byte
	copy(destToken2[:], destToken2Bytes)
	extraData2, err := hexutil.Decode("0xffeeddcc")
	require.NoError(t, err)
	token2 := any2AptosTokenTransfer{
		SourcePoolAddress: srcPool2,
		DestTokenAddress:  destToken2,
		DestGasAmount:     20000,
		ExtraData:         extraData2,
		Amount:            big.NewInt(5000000),
	}

	tokens := []any2AptosTokenTransfer{token1, token2}

	computedHash, err := computeMessageDataHash(metadataHash, messageID, receiver, sequenceNumber, gasLimit, nonce, sender, data, tokens)
	require.NoError(t, err)

	require.Equal(t, expectedHashStr, hexutil.Encode(computedHash[:]), "Computed hash does not match expected hash")
}

// the equivalent test with the same values exists in the Aptos offramp.move contract, see test_calculate_metadata_hash
func TestComputeMetadataHash(t *testing.T) {
	expectedHashStr := "0x812acb01df318f85be452cf6664891cf5481a69dac01e0df67102a295218dd17"
	expectedHashAlternateStr := "0x6caf8756ae02ee4f12b83b38e0f21b5e43e90d203bd06729486fd4a0fc8bcc5e"

	sourceChainSelector := uint64(123456789)
	destinationChainSelector := uint64(987654321)
	onRamp := []byte("source-onramp-address")

	metadataHash, err := computeMetadataHash(sourceChainSelector, destinationChainSelector, onRamp)
	require.NoError(t, err)
	require.Equal(t, expectedHashStr, hexutil.Encode(metadataHash[:]), "Computed hash does not match expected hash")

	metadataHashAlternate, err := computeMetadataHash(sourceChainSelector+1, destinationChainSelector, onRamp)
	require.NoError(t, err)
	require.Equal(t, expectedHashAlternateStr, hexutil.Encode(metadataHashAlternate[:]), "Alternate computed hash does not match expected alternate hash")
}
