package hasher_test

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/evm_2_evm_onramp"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/hasher"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/testhelpers"
)

func TestHasher(t *testing.T) {
	sourceChainSelector, destChainSelector := uint64(1), uint64(4)
	onRampAddress := common.HexToAddress("0x5550000000000000000000000000000000000001")

	hashingCtx := hasher.NewKeccakCtx()

	hasher := hasher.NewLeafHasher(sourceChainSelector, destChainSelector, onRampAddress, hashingCtx)

	message := evm_2_evm_onramp.InternalEVM2EVMMessage{
		SourceChainSelector: sourceChainSelector,
		SequenceNumber:      1337,
		FeeTokenAmount:      big.NewInt(1),
		Sender:              common.HexToAddress("0x1110000000000000000000000000000000000001"),
		Nonce:               1337,
		GasLimit:            big.NewInt(100),
		Strict:              false,
		Receiver:            common.HexToAddress("0x2220000000000000000000000000000000000001"),
		Data:                []byte{},
		TokenAmounts:        []evm_2_evm_onramp.ClientEVMTokenAmount{{Token: common.HexToAddress("0x4440000000000000000000000000000000000001"), Amount: big.NewInt(12345678900)}},
		FeeToken:            common.Address{},
		MessageId:           [32]byte{},
	}

	hash, err := hasher.HashLeaf(testhelpers.GenerateCCIPSendLog(t, message))
	require.NoError(t, err)

	// NOTE: Must match spec
	require.Equal(t, "26f282c6ac8231933b1799648d01ff6cec792a33fb37408b4d135968f9168ace", hex.EncodeToString(hash[:]))

	message = evm_2_evm_onramp.InternalEVM2EVMMessage{
		SourceChainSelector: sourceChainSelector,
		SequenceNumber:      1337,
		FeeTokenAmount:      big.NewInt(1e12),
		Sender:              common.HexToAddress("0x1110000000000000000000000000000000000001"),
		Nonce:               1337,
		GasLimit:            big.NewInt(100),
		Strict:              false,
		Receiver:            common.HexToAddress("0x2220000000000000000000000000000000000001"),
		Data:                []byte("foo bar baz"),
		TokenAmounts: []evm_2_evm_onramp.ClientEVMTokenAmount{
			{Token: common.HexToAddress("0x4440000000000000000000000000000000000001"), Amount: big.NewInt(12345678900)},
			{Token: common.HexToAddress("0x6660000000000000000000000000000000000001"), Amount: big.NewInt(4204242)},
		},
		FeeToken:  common.Address{},
		MessageId: [32]byte{},
	}

	hash, err = hasher.HashLeaf(testhelpers.GenerateCCIPSendLog(t, message))
	require.NoError(t, err)

	// NOTE: Must match spec
	require.Equal(t, "05cee92e7cb86a37b6536554828a5b21ff20ac3d4ef821ec47056f1d963313de", hex.EncodeToString(hash[:]))
}

func TestMetaDataHash(t *testing.T) {
	sourceChainSelector, destChainSelector := uint64(1), uint64(4)
	onRampAddress := common.HexToAddress("0x5550000000000000000000000000000000000001")
	ctx := hasher.NewKeccakCtx()
	hash := hasher.GetMetaDataHash(ctx, ctx.Hash([]byte("EVM2EVMSubscriptionMessagePlus")), sourceChainSelector, onRampAddress, destChainSelector)
	require.Equal(t, "e8b93c9d01a7a72ec6c7235e238701cf1511b267a31fdb78dd342649ee58c08d", hex.EncodeToString(hash[:]))
}
