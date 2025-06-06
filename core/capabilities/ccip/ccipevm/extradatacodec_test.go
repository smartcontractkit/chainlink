package ccipevm

import (
	"math/big"
	"math/rand"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/message_hasher"
)

func Test_decodeExtraData(t *testing.T) {
	d := testSetup(t)
	gasLimit := big.NewInt(rand.Int63())
	extraDataDecoder := &ExtraDataDecoder{}

	t.Run("decode extra args into map evm v1", func(t *testing.T) {
		encoded, err := d.contract.EncodeEVMExtraArgsV1(nil, message_hasher.ClientEVMExtraArgsV1{
			GasLimit: gasLimit,
		})
		require.NoError(t, err)

		m, err := extraDataDecoder.DecodeExtraArgsToMap(encoded)
		require.NoError(t, err)
		require.Len(t, m, 1)

		gl, exist := m["gasLimit"]
		require.True(t, exist)
		require.Equal(t, gl, gasLimit)
	})

	t.Run("decode extra args into map evm v2", func(t *testing.T) {
		encoded, err := d.contract.EncodeEVMExtraArgsV2(nil, message_hasher.ClientGenericExtraArgsV2{
			GasLimit:                 gasLimit,
			AllowOutOfOrderExecution: true,
		})
		require.NoError(t, err)

		m, err := extraDataDecoder.DecodeExtraArgsToMap(encoded)
		require.NoError(t, err)
		require.Len(t, m, 2)

		gl, exist := m["gasLimit"]
		require.True(t, exist)
		require.Equal(t, gl, gasLimit)

		ooe, exist := m["allowOutOfOrderExecution"]
		require.True(t, exist)
		require.Equal(t, true, ooe)
	})

	t.Run("decode dest exec data into map", func(t *testing.T) {
		destGasAmount := uint32(10000)
		encoded, err := abiEncodeUint32(destGasAmount)
		require.NoError(t, err)
		m, err := extraDataDecoder.DecodeDestExecDataToMap(encoded)
		require.NoError(t, err)
		require.Len(t, m, 1)

		decoded, exist := m[evmDestExecDataKey]
		require.True(t, exist)
		require.Equal(t, destGasAmount, decoded)
	})
}

func TestExtraDataDecoder_DecodeExtraArgsToMap_SVMDestination(t *testing.T) {
	d := testSetup(t)
	extraDataDecoder := &ExtraDataDecoder{}

	t.Run("decode extra args into map svm", func(t *testing.T) {
		key, err := solana.NewRandomPrivateKey()
		require.NoError(t, err)
		cu := uint32(10000)
		bitmap := uint64(4)
		oooExec := true // enforced ooo exec for svm
		tokenReceiver := [32]byte(key.PublicKey().Bytes())
		accounts := [][32]byte{[32]byte(key.PublicKey().Bytes())}
		encoded, err := d.contract.EncodeSVMExtraArgsV1(nil, message_hasher.ClientSVMExtraArgsV1{
			ComputeUnits:             cu,
			AccountIsWritableBitmap:  bitmap,
			AllowOutOfOrderExecution: oooExec,
			TokenReceiver:            tokenReceiver,
			Accounts:                 accounts,
		})
		require.NoError(t, err)

		m, err := extraDataDecoder.DecodeExtraArgsToMap(encoded)
		require.NoError(t, err)
		require.Len(t, m, 5)

		cuDecoded, exist := m["computeUnits"]
		require.True(t, exist)
		require.Equal(t, cuDecoded, cu)

		bitmapDecoded, exist := m["accountIsWritableBitmap"]
		require.True(t, exist)
		require.Equal(t, bitmapDecoded, bitmap)

		ooeDecoded, exist := m["allowOutOfOrderExecution"]
		require.True(t, exist)
		require.Equal(t, ooeDecoded, oooExec)

		tokenReceiverDecoded, exist := m["tokenReceiver"]
		require.True(t, exist)
		require.Equal(t, tokenReceiverDecoded, tokenReceiver)

		accountsDecoded, exist := m["accounts"]
		require.True(t, exist)
		require.Equal(t, accountsDecoded, accounts)
	})

	t.Run("mostly empty except for out of order execution", func(t *testing.T) {
		extraData := hexutil.MustDecode("0x1f3b3aba00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
		m, err := extraDataDecoder.DecodeExtraArgsToMap(extraData)
		require.NoError(t, err)
		require.Equal(t, uint32(0), m["computeUnits"])
		require.Equal(t, true, m["allowOutOfOrderExecution"])
		require.Equal(t, uint64(0), m["accountIsWritableBitmap"])
		require.Equal(t, [][32]byte{}, m["accounts"])
		require.Equal(t, [32]byte{}, m["tokenReceiver"])
	})

	t.Run("real world use case", func(t *testing.T) {
		// See tx https://sepolia.etherscan.io/tx/0x9992d5c16e6b498ae94ccdcf942e23067980ee0c9d6b63bde6a1b931b12b12fc
		extraData := hexutil.MustDecode("0x1f3b3aba0000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000001388000000000000000000000000000000000000000000000000000000000000000030000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000000032b4c0e1c22b4ecbb7650755cdf3e6ae0815d71cd6828cf1ae9bb57fb37be51bd36b59f124bfe0e67767a2333b5f0d4956074bb4d2b221f2873812988278afd640000000000000000000000000000000000000000000000000000000000000000")
		m, err := extraDataDecoder.DecodeExtraArgsToMap(extraData)
		require.NoError(t, err)
		require.Equal(t, uint32(0x13880), m["computeUnits"])
		require.Equal(t, true, m["allowOutOfOrderExecution"])
		require.Equal(t, uint64(3), m["accountIsWritableBitmap"])
		require.Equal(t, [][32]uint8{
			{0x2b, 0x4c, 0xe, 0x1c, 0x22, 0xb4, 0xec, 0xbb, 0x76, 0x50, 0x75, 0x5c, 0xdf, 0x3e, 0x6a, 0xe0, 0x81, 0x5d, 0x71, 0xcd, 0x68, 0x28, 0xcf, 0x1a, 0xe9, 0xbb, 0x57, 0xfb, 0x37, 0xbe, 0x51, 0xbd},
			{0x36, 0xb5, 0x9f, 0x12, 0x4b, 0xfe, 0xe, 0x67, 0x76, 0x7a, 0x23, 0x33, 0xb5, 0xf0, 0xd4, 0x95, 0x60, 0x74, 0xbb, 0x4d, 0x2b, 0x22, 0x1f, 0x28, 0x73, 0x81, 0x29, 0x88, 0x27, 0x8a, 0xfd, 0x64},
			{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}},
			m["accounts"])
		require.Equal(t, [32]byte{}, m["tokenReceiver"])
	})
}
