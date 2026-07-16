package ccipsolana

import (
	cryptorand "crypto/rand"
	"math/big"
	"math/rand"
	"testing"

	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-ccip/chains/solana/contracts/tests/config"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/latest/ccip_offramp"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/ccip"
	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipaptos"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipevm"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
)

func TestMessageHasher_EVM2SVM(t *testing.T) {
	registeredExtraDataCodecMap := map[string]ccipocr3.SourceChainExtraDataCodec{
		chainsel.FamilyAptos:  ccipaptos.ExtraDataDecoder{},
		chainsel.FamilyEVM:    ccipevm.ExtraDataDecoder{},
		chainsel.FamilySolana: ExtraDataDecoder{},
		chainsel.FamilySui:    ccipaptos.ExtraDataDecoder{},
	}
	var extraDataCodec = ccipocr3.ExtraDataCodecMap(registeredExtraDataCodecMap)
	any2AnyMsg, any2SolanaMsg, msgAccounts := createEVM2SolanaMessages(t)
	msgHasher := NewMessageHasherV1(logger.Test(t), extraDataCodec)
	actualHash, err := msgHasher.Hash(testutils.Context(t), any2AnyMsg)
	require.NoError(t, err)
	expectedHash, err := ccip.HashAnyToSVMMessage(any2SolanaMsg, any2AnyMsg.Header.OnRamp, msgAccounts)
	require.NoError(t, err)
	require.Equal(t, expectedHash, actualHash[:32])
}

func TestMessageHasher_InvalidReceiver(t *testing.T) {
	any2AnyMsg, _, _ := createEVM2SolanaMessages(t)

	// Set receiver to a []byte of 2 length
	any2AnyMsg.Receiver = []byte{0, 0}
	mockExtraDataCodec := mocks.NewSourceChainExtraDataCodec(t)
	mockExtraDataCodec.On("DecodeDestExecDataToMap", mock.Anything).Return(map[string]any{
		"destGasAmount": uint32(10),
	}, nil).Maybe()
	mockExtraDataCodec.On("DecodeExtraArgsToMap", mock.Anything).Return(map[string]any{
		"ComputeUnits":            uint32(1000),
		"AccountIsWritableBitmap": uint64(10),
		"Accounts": [][32]byte{
			[32]byte(config.CcipLogicReceiver.Bytes()),
			[32]byte(config.ReceiverTargetAccountPDA.Bytes()),
			[32]byte(solana.SystemProgramID.Bytes()),
		},
	}, nil).Maybe()

	registeredMockExtraDataCodecMap := map[string]ccipocr3.SourceChainExtraDataCodec{
		chainsel.FamilyEVM:    mockExtraDataCodec,
		chainsel.FamilySolana: mockExtraDataCodec,
	}

	edc := ccipocr3.ExtraDataCodecMap(registeredMockExtraDataCodecMap)
	msgHasher := NewMessageHasherV1(logger.Test(t), edc)
	_, err := msgHasher.Hash(testutils.Context(t), any2AnyMsg)
	require.Error(t, err)
}

func TestMessageHasher_InvalidDestinationTokenAddress(t *testing.T) {
	any2AnyMsg, _, _ := createEVM2SolanaMessages(t)

	// Set DestTokenAddress to a []byte of 2 length
	any2AnyMsg.TokenAmounts[0].DestTokenAddress = []byte{0, 0}
	mockExtraDataCodec := mocks.NewSourceChainExtraDataCodec(t)
	mockExtraDataCodec.On("DecodeDestExecDataToMap", mock.Anything).Return(map[string]any{
		"destGasAmount": uint32(10),
	}, nil).Maybe()
	mockExtraDataCodec.On("DecodeExtraArgsToMap", mock.Anything).Return(map[string]any{
		"ComputeUnits":            uint32(1000),
		"AccountIsWritableBitmap": uint64(10),
		"Accounts": [][32]byte{
			[32]byte(config.CcipLogicReceiver.Bytes()),
			[32]byte(config.ReceiverTargetAccountPDA.Bytes()),
			[32]byte(solana.SystemProgramID.Bytes()),
		},
	}, nil).Maybe()

	registeredMockExtraDataCodecMap := map[string]ccipocr3.SourceChainExtraDataCodec{
		chainsel.FamilyEVM:    mockExtraDataCodec,
		chainsel.FamilySolana: mockExtraDataCodec,
	}
	edc := ccipocr3.ExtraDataCodecMap(registeredMockExtraDataCodecMap)
	msgHasher := NewMessageHasherV1(logger.Test(t), edc)
	_, err := msgHasher.Hash(testutils.Context(t), any2AnyMsg)
	require.Error(t, err)
}

func TestParseExtraDataMap_NativeTypes(t *testing.T) {
	// Test with native Go types (in-process codec, no LOOP)
	account1 := [32]byte{0x01}
	account2 := [32]byte{0x02}
	receiver := [32]byte{0x03}

	input := map[string]any{
		"ComputeUnits":             uint32(2000000),
		"AccountIsWritableBitmap":  uint64(6),
		"AllowOutOfOrderExecution": true,
		"TokenReceiver":            receiver,
		"Accounts":                 [][32]byte{account1, account2},
	}

	ed, err := parseExtraDataMap(input)
	require.NoError(t, err)
	assert.Equal(t, uint32(2000000), ed.extraArgs.ComputeUnits)
	assert.Equal(t, uint64(6), ed.extraArgs.IsWritableBitmap)
	assert.Equal(t, solana.PublicKeyFromBytes(receiver[:]), ed.tokenReceiver)
	require.Len(t, ed.accounts, 2)
	assert.Equal(t, solana.PublicKeyFromBytes(account1[:]), ed.accounts[0])
	assert.Equal(t, solana.PublicKeyFromBytes(account2[:]), ed.accounts[1])
}

func TestParseExtraDataMap_LOOPConvertedTypes(t *testing.T) {
	// Test with types as they arrive after LOOP gRPC roundtrip:
	// uint32 -> int64, [32]byte -> []byte, [][32]byte -> []interface{}
	account1 := make([]byte, 32)
	account1[0] = 0x01
	account2 := make([]byte, 32)
	account2[0] = 0x02
	receiver := make([]byte, 32)
	receiver[0] = 0x03

	input := map[string]any{
		"ComputeUnits":             int64(2000000), // LOOP: uint32 -> int64
		"AccountIsWritableBitmap":  int64(6),       // LOOP: uint64 -> int64
		"AllowOutOfOrderExecution": true,
		"TokenReceiver":            receiver, // LOOP: [32]byte -> []byte
		"Accounts": []interface{}{ // LOOP: [][32]byte -> []interface{}
			account1,
			account2,
		},
	}

	ed, err := parseExtraDataMap(input)
	require.NoError(t, err)
	assert.Equal(t, uint32(2000000), ed.extraArgs.ComputeUnits)
	assert.Equal(t, uint64(6), ed.extraArgs.IsWritableBitmap)
	assert.Equal(t, solana.PublicKeyFromBytes(receiver), ed.tokenReceiver)
	require.Len(t, ed.accounts, 2)
	assert.Equal(t, solana.PublicKeyFromBytes(account1), ed.accounts[0])
	assert.Equal(t, solana.PublicKeyFromBytes(account2), ed.accounts[1])
}

func TestParseExtraDataMap_LOOPAccountsAsByteSlice(t *testing.T) {
	// Test [][]byte variant (alternative LOOP representation)
	account1 := make([]byte, 32)
	account1[0] = 0x01

	input := map[string]any{
		"ComputeUnits":            uint32(1000),
		"AccountIsWritableBitmap": uint64(0),
		"TokenReceiver":           [32]byte{},
		"Accounts":                [][]byte{account1},
	}

	ed, err := parseExtraDataMap(input)
	require.NoError(t, err)
	require.Len(t, ed.accounts, 1)
	assert.Equal(t, solana.PublicKeyFromBytes(account1), ed.accounts[0])
}

func TestParseExtraDataMap_InvalidTypes(t *testing.T) {
	t.Run("ComputeUnits out of range", func(t *testing.T) {
		input := map[string]any{
			"ComputeUnits": int64(5000000000), // exceeds uint32 max
		}
		_, err := parseExtraDataMap(input)
		require.ErrorContains(t, err, "ComputeUnits out of uint32 range")
	})

	t.Run("TokenReceiver wrong length", func(t *testing.T) {
		input := map[string]any{
			"TokenReceiver": []byte{0x01, 0x02}, // not 32 bytes
		}
		_, err := parseExtraDataMap(input)
		require.ErrorContains(t, err, "invalid length for TokenReceiver")
	})

	t.Run("Accounts element wrong length", func(t *testing.T) {
		input := map[string]any{
			"Accounts": []interface{}{
				[]byte{0x01, 0x02}, // not 32 bytes
			},
		}
		_, err := parseExtraDataMap(input)
		require.ErrorContains(t, err, "invalid length for Accounts[0]")
	})

	t.Run("Accounts element wrong type", func(t *testing.T) {
		input := map[string]any{
			"Accounts": []interface{}{
				"not bytes",
			},
		}
		_, err := parseExtraDataMap(input)
		require.ErrorContains(t, err, "invalid type for Accounts[0]")
	})
}

func createEVM2SolanaMessages(t *testing.T) (cciptypes.Message, ccip_offramp.Any2SVMRampMessage, []solana.PublicKey) {
	messageID := utils.RandomBytes32()
	sourceChain := uint64(5009297550715157269) // evm mainnet
	seqNum := rand.Uint64()
	nonce := rand.Uint64()
	destChain := rand.Uint64()

	messageData := make([]byte, rand.Intn(2048))
	_, err := cryptorand.Read(messageData)
	require.NoError(t, err)

	sender := abiEncodedAddress(t)
	receiver := solana.MustPublicKeyFromBase58("DS2tt4BX7YwCw7yrDNwbAdnYrxjeCPeGJbHmZEYC8RTb")
	tokenReceiver := solana.MustPublicKeyFromBase58("42Gia5bGsh8R2S44e37t9fsucap1qsgjr6GjBmWotgdF")
	extraArgs := ccip_offramp.Any2SVMRampExtraArgs{
		ComputeUnits:     uint32(10000),
		IsWritableBitmap: uint64(4),
	}
	abiEncodedExtraArgs := []byte{31, 59, 58, 186, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 32, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 39, 16, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 44, 230, 105, 156, 244, 184, 196, 235, 30, 58, 209, 82, 8, 202, 25, 73, 167, 169, 34, 150, 141, 129, 169, 150, 219, 160, 186, 44, 72, 156, 50, 170, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 160, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 44, 230, 105, 156, 244, 184, 196, 235, 30, 58, 209, 82, 8, 202, 25, 73, 167, 169, 34, 150, 141, 129, 169, 150, 219, 160, 186, 44, 72, 156, 50, 170}
	tokenAmount := cciptypes.NewBigInt(big.NewInt(rand.Int63()))
	destGasAmount, err := abiEncodeUint32(10)
	require.NoError(t, err)

	ccipTokenAmounts := make([]cciptypes.RampTokenAmount, 5)
	for z := range 5 {
		ccipTokenAmounts[z] = cciptypes.RampTokenAmount{
			SourcePoolAddress: cciptypes.UnknownAddress("DS2tt4BX7YwCw7yrDNwbAdnYrxjeCPeGJbHmZEYC8RTb"),
			DestTokenAddress:  receiver.Bytes(),
			Amount:            tokenAmount,
			DestExecData:      destGasAmount,
		}
	}

	solTokenAmounts := make([]ccip_offramp.Any2SVMTokenTransfer, 5)
	for z := range 5 {
		solTokenAmounts[z] = ccip_offramp.Any2SVMTokenTransfer{
			SourcePoolAddress: cciptypes.UnknownAddress("DS2tt4BX7YwCw7yrDNwbAdnYrxjeCPeGJbHmZEYC8RTb"),
			DestTokenAddress:  receiver,
			Amount:            ccip_offramp.CrossChainAmount{LeBytes: [32]uint8(encodeBigIntToFixedLengthLE(tokenAmount.Int, 32))},
			DestGasAmount:     uint32(10),
		}
	}

	any2SolanaMsg := ccip_offramp.Any2SVMRampMessage{
		Header: ccip_offramp.RampMessageHeader{
			MessageId:           messageID,
			SourceChainSelector: sourceChain,
			DestChainSelector:   destChain,
			SequenceNumber:      seqNum,
			Nonce:               nonce,
		},
		Sender:        sender,
		TokenReceiver: tokenReceiver,
		Data:          messageData,
		TokenAmounts:  solTokenAmounts,
		ExtraArgs:     extraArgs,
	}
	any2AnyMsg := cciptypes.Message{
		Header: cciptypes.RampMessageHeader{
			MessageID:           messageID,
			SourceChainSelector: cciptypes.ChainSelector(sourceChain),
			DestChainSelector:   cciptypes.ChainSelector(destChain),
			SequenceNumber:      cciptypes.SeqNum(seqNum),
			Nonce:               nonce,
			OnRamp:              abiEncodedAddress(t),
		},
		Sender:         sender,
		Receiver:       receiver.Bytes(),
		Data:           messageData,
		TokenAmounts:   ccipTokenAmounts,
		FeeToken:       []byte{},
		FeeTokenAmount: cciptypes.NewBigIntFromInt64(0),
		ExtraArgs:      abiEncodedExtraArgs,
	}

	msgAccounts := []solana.PublicKey{
		receiver,
		solana.MustPublicKeyFromBase58("42Gia5bGsh8R2S44e37t9fsucap1qsgjr6GjBmWotgdF"),
	}
	return any2AnyMsg, any2SolanaMsg, msgAccounts
}

func abiEncodedAddress(t *testing.T) []byte {
	addr := utils.RandomAddress()
	encoded, err := utils.ABIEncode(`[{"type": "address"}]`, addr)
	require.NoError(t, err)
	return encoded
}

func abiEncodeUint32(data uint32) ([]byte, error) {
	return utils.ABIEncode(`[{ "type": "uint32" }]`, data)
}

func TestToLittleEndian(t *testing.T) {
	mustSetString := func(s string) *big.Int {
		b, ok := big.NewInt(0).SetString(s, 10)
		if !ok {
			t.Fatalf("failed to set string %s", s)
		}
		return b
	}

	var tests = []struct {
		input    *big.Int
		expected []byte
	}{
		{
			input:    mustSetString("93632917990780833250"),
			expected: []uint8{0xe2, 0xd, 0xc6, 0xfb, 0xd2, 0xf2, 0x6a, 0x13, 0x5, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
		},
		{
			input:    mustSetString("9363291799078083325000000910750912570125"),
			expected: []uint8{0xd, 0x63, 0xbf, 0xfc, 0xcb, 0xfc, 0x27, 0x0, 0x4c, 0x7e, 0xe5, 0x81, 0xc7, 0x67, 0x28, 0x84, 0x1b, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
		},
	}

	for _, test := range tests {
		t.Run(test.input.String(), func(t *testing.T) {
			result := encodeBigIntToFixedLengthLE(test.input, 32)
			assert.Equal(t, test.expected, result, "expected %x, got %x", test.expected, result)
		})
	}
}
