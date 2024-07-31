package v1_0_0

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/hashutil"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_onramp_1_0_0"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
)

func TestHasherV1_0_0(t *testing.T) {
	sourceChainSelector, destChainSelector := uint64(1), uint64(4)
	onRampAddress := common.HexToAddress("0x5550000000000000000000000000000000000001")
	onRampABI := abihelpers.MustParseABI(evm_2_evm_onramp_1_0_0.EVM2EVMOnRampABI)

	ramp, err := evm_2_evm_onramp_1_0_0.NewEVM2EVMOnRamp(onRampAddress, nil)
	require.NoError(t, err)
	hashingCtx := hashutil.NewKeccak()
	hasher := NewLeafHasher(sourceChainSelector, destChainSelector, onRampAddress, hashingCtx, ramp)

	message := evm_2_evm_onramp_1_0_0.InternalEVM2EVMMessage{
		SourceChainSelector: sourceChainSelector,
		Sender:              common.HexToAddress("0x1110000000000000000000000000000000000001"),
		Receiver:            common.HexToAddress("0x2220000000000000000000000000000000000001"),
		SequenceNumber:      1337,
		GasLimit:            big.NewInt(100),
		Strict:              false,
		Nonce:               1337,
		FeeToken:            common.Address{},
		FeeTokenAmount:      big.NewInt(1),
		Data:                []byte{},
		TokenAmounts:        []evm_2_evm_onramp_1_0_0.ClientEVMTokenAmount{{Token: common.HexToAddress("0x4440000000000000000000000000000000000001"), Amount: big.NewInt(12345678900)}},
		MessageId:           [32]byte{},
	}

	data, err := onRampABI.Events[CCIPSendRequestedEventName].Inputs.Pack(message)
	require.NoError(t, err)
	hash, err := hasher.HashLeaf(types.Log{Topics: []common.Hash{abihelpers.MustGetEventID("CCIPSendRequested", onRampABI)}, Data: data})
	require.NoError(t, err)

	// NOTE: Must match spec
	require.Equal(t, "26f282c6ac8231933b1799648d01ff6cec792a33fb37408b4d135968f9168ace", hex.EncodeToString(hash[:]))

	message = evm_2_evm_onramp_1_0_0.InternalEVM2EVMMessage{
		SourceChainSelector: sourceChainSelector,
		Sender:              common.HexToAddress("0x1110000000000000000000000000000000000001"),
		Receiver:            common.HexToAddress("0x2220000000000000000000000000000000000001"),
		SequenceNumber:      1337,
		GasLimit:            big.NewInt(100),
		Strict:              false,
		Nonce:               1337,
		FeeToken:            common.Address{},
		FeeTokenAmount:      big.NewInt(1e12),
		Data:                []byte("foo bar baz"),
		TokenAmounts: []evm_2_evm_onramp_1_0_0.ClientEVMTokenAmount{
			{Token: common.HexToAddress("0x4440000000000000000000000000000000000001"), Amount: big.NewInt(12345678900)},
			{Token: common.HexToAddress("0x6660000000000000000000000000000000000001"), Amount: big.NewInt(4204242)},
		},
		MessageId: [32]byte{},
	}

	data, err = onRampABI.Events[CCIPSendRequestedEventName].Inputs.Pack(message)
	require.NoError(t, err)
	hash, err = hasher.HashLeaf(types.Log{Topics: []common.Hash{abihelpers.MustGetEventID("CCIPSendRequested", onRampABI)}, Data: data})
	require.NoError(t, err)

	// NOTE: Must match spec
	require.Equal(t, "05cee92e7cb86a37b6536554828a5b21ff20ac3d4ef821ec47056f1d963313de", hex.EncodeToString(hash[:]))
}

func TestMetaDataHash(t *testing.T) {
	sourceChainSelector, destChainSelector := uint64(1), uint64(4)
	onRampAddress := common.HexToAddress("0x5550000000000000000000000000000000000001")
	ctx := hashutil.NewKeccak()
	hash := GetMetaDataHash(ctx, ctx.Hash([]byte(MetaDataHashPrefix)), sourceChainSelector, onRampAddress, destChainSelector)
	require.Equal(t, "1409948abde219f43870c3d6d1c16beabd8878eb5039a3fa765eb56e4b8ded9e", hex.EncodeToString(hash[:]))
}
