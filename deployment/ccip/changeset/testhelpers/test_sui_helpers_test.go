package testhelpers

import (
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/require"
)

// Test_MakeSuiSourceSVMExtraArgsV1_GoldenVector asserts the Sui-source SVMExtraArgsV1 builder
// produces the exact BCS byte vector the ccipaptos decoder (which the Solana dest MessageHasher
// reaches via the source-family-dispatching codec bundle) round-trips. The expected bytes are the
// golden vector from core/capabilities/ccip/ccipaptos/extradatadecoder_test.go
// ("decode extra args into map svm"), proving encoder and decoder agree and guarding the
// length-prefix-vs-FixedBytes distinction for the vector<u8> token_receiver/accounts fields.
func Test_MakeSuiSourceSVMExtraArgsV1_GoldenVector(t *testing.T) {
	computeUnits := uint32(100000)
	bitmap := uint64(255)
	allowOOO := true

	tokenReceiverBytes := hexutil.MustDecode("0x1234567890123456789012345678901234567890123456789012345678901234")
	require.Len(t, tokenReceiverBytes, 32)
	var tokenReceiver [32]byte
	copy(tokenReceiver[:], tokenReceiverBytes)

	account0Bytes := hexutil.MustDecode("0x1234567890123456789012345678901212345678901234567890123456789012")
	require.Len(t, account0Bytes, 32)
	var account0 [32]byte
	copy(account0[:], account0Bytes)

	// Expected golden vector (tag + BCS fields), identical to the ccipaptos decoder test fixture:
	//   1f3b3aba                       SVMExtraArgsV1 tag
	//   a0860100                       u32 LE 100000
	//   ff00000000000000               u64 LE 255
	//   01                             bool true
	//   20 + 32B                       tokenReceiver vector<u8>
	//   01 + 20 + 32B                  accounts: 1 x vector<u8>(32)
	expected := hexutil.MustDecode(
		"0x1f3b3abaa0860100ff000000000000000120123456789012345678901234567890123456789012345678901234567890123401201234567890123456789012345678901212345678901234567890123456789012",
	)

	got := MakeSuiSourceSVMExtraArgsV1(computeUnits, bitmap, allowOOO, tokenReceiver, [][32]byte{account0})
	require.Equal(t, expected, got, "MakeSuiSourceSVMExtraArgsV1 must match the ccipaptos decoder golden vector")
}

// Test_MakeSuiSourceSVMExtraArgsV1_MultipleAccounts covers the multi-account sequence case
// against the second ccipaptos decoder fixture ("svm with multiple accounts").
func Test_MakeSuiSourceSVMExtraArgsV1_MultipleAccounts(t *testing.T) {
	computeUnits := uint32(100000)
	bitmap := uint64(255)
	allowOOO := true

	tokenReceiverBytes := hexutil.MustDecode("0x1234567890123456789012345678901234567890123456789012345678901234")
	var tokenReceiver [32]byte
	copy(tokenReceiver[:], tokenReceiverBytes)

	account0Bytes := hexutil.MustDecode("0x1234567890123456789012345678901212345678901234567890123456789012")
	account1Bytes := hexutil.MustDecode("0x9ab25d7fff22ac56789012345678901212345678901234567890123456789012")
	var account0, account1 [32]byte
	copy(account0[:], account0Bytes)
	copy(account1[:], account1Bytes)

	expected := hexutil.MustDecode(
		"0x1f3b3abaa0860100ff000000000000000120123456789012345678901234567890123456789012345678901234567890123402201234567890123456789012345678901212345678901234567890123456789012209ab25d7fff22ac56789012345678901212345678901234567890123456789012",
	)

	got := MakeSuiSourceSVMExtraArgsV1(computeUnits, bitmap, allowOOO, tokenReceiver, [][32]byte{account0, account1})
	require.Equal(t, expected, got, "multi-account vector must match the ccipaptos decoder fixture")
}
