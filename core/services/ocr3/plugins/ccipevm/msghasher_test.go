package ccipevm

import (
	"context"
	cryptorand "crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"math/rand"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/message_hasher"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"

	"github.com/stretchr/testify/require"
)

func TestMessageHasher_e2e(t *testing.T) {
	ctx := testutils.Context(t)
	d := testSetup(t)

	// low budget "fuzz" test.
	// TODO: should actually write a real fuzz test.
	for i := 0; i < 5; i++ {
		testHasher(ctx, t, d)
	}
}

func testHasher(ctx context.Context, t *testing.T, d *testSetupData) {
	destChainSelector := rand.Uint64()
	onRampAddress := testutils.NewAddress().Bytes()
	ccipMsg := createCCIPMsg(t)

	evmTokenAmounts := make([]message_hasher.ClientEVMTokenAmount, 0, len(ccipMsg.TokenAmounts))
	for _, ta := range ccipMsg.TokenAmounts {
		evmTokenAmounts = append(evmTokenAmounts, message_hasher.ClientEVMTokenAmount{
			Token:  common.HexToAddress(string(ta.Token)),
			Amount: ta.Amount,
		})
	}
	evmMsg := message_hasher.InternalAny2EVMRampMessage{
		Header: message_hasher.InternalRampMessageHeader{
			MessageId:           mustMessageID(t, ccipMsg.ID),
			SourceChainSelector: uint64(ccipMsg.SourceChain),
			DestChainSelector:   destChainSelector,
			SequenceNumber:      uint64(ccipMsg.SeqNum),
			Nonce:               ccipMsg.Nonce,
		},
		Sender:          common.HexToAddress(string(ccipMsg.Sender)).Bytes(),
		Receiver:        common.HexToAddress(string(ccipMsg.Receiver)),
		GasLimit:        ccipMsg.ChainFeeLimit.Int,
		Data:            ccipMsg.Data,
		TokenAmounts:    evmTokenAmounts,
		SourceTokenData: ccipMsg.SourceTokenData,
	}

	expectedHash, err := d.contract.Hash(&bind.CallOpts{Context: ctx}, evmMsg, onRampAddress)
	require.NoError(t, err)

	evmMsgHasher := NewMessageHasherV1(onRampAddress, cciptypes.ChainSelector(destChainSelector))
	actualHash, err := evmMsgHasher.Hash(ctx, ccipMsg)
	require.NoError(t, err)

	require.Equal(t, fmt.Sprintf("%x", expectedHash), strings.TrimPrefix(actualHash.String(), "0x"))
}

// TODO: fix this once messageID is part of CCIPMsg
func createCCIPMsg(t *testing.T) cciptypes.CCIPMsg {
	// Setup random msg data
	messageID := utils.RandomBytes32()

	sourceTokenData := make([]byte, rand.Intn(2048))
	_, err := cryptorand.Read(sourceTokenData)
	require.NoError(t, err)

	sourceChain := rand.Uint64()
	seqNum := rand.Uint64()
	chainFeeLimit := rand.Uint64()
	nonce := rand.Uint64()

	messageData := make([]byte, rand.Intn(2048))
	_, err = cryptorand.Read(messageData)
	require.NoError(t, err)

	sourceTokenDatas := make([][]byte, rand.Intn(10))
	for i := range sourceTokenDatas {
		sourceTokenDatas[i] = sourceTokenData
	}

	numTokenAmounts := rand.Intn(50)
	tokenAmounts := make([]cciptypes.TokenAmount, 0, numTokenAmounts)
	for i := 0; i < numTokenAmounts; i++ {
		tokenAmounts = append(tokenAmounts, cciptypes.TokenAmount{
			Token:  types.Account(utils.RandomAddress().String()),
			Amount: big.NewInt(0).SetUint64(rand.Uint64()),
		})
	}
	return cciptypes.CCIPMsg{
		CCIPMsgBaseDetails: cciptypes.CCIPMsgBaseDetails{
			ID:          hex.EncodeToString(messageID[:]),
			SourceChain: cciptypes.ChainSelector(sourceChain),
			SeqNum:      cciptypes.SeqNum(seqNum),
		},
		ChainFeeLimit: cciptypes.NewBigInt(big.NewInt(0).SetUint64(chainFeeLimit)),
		Nonce:         nonce,
		Sender:        types.Account(utils.RandomAddress().String()),
		Receiver:      types.Account(utils.RandomAddress().String()),
		// TODO: remove this field if not needed, not used by the hasher
		// Strict: strict,
		// NOTE: not used by the hasher
		// FeeToken:        types.Account(utils.RandomAddress().String()),
		// FeeTokenAmount:  cciptypes.NewBigInt(big.NewInt(0).SetUint64(feeTokenAmount)),
		Data:            messageData,
		TokenAmounts:    tokenAmounts,
		SourceTokenData: sourceTokenDatas,
	}
}

func mustMessageID(t *testing.T, msgIDHex string) [32]byte {
	msgID, err := hex.DecodeString(msgIDHex)
	require.NoError(t, err)
	require.Len(t, msgID, 32)
	var msgID32 [32]byte
	copy(msgID32[:], msgID)
	return msgID32
}

type testSetupData struct {
	contractAddr common.Address
	contract     *message_hasher.MessageHasher
	sb           *backends.SimulatedBackend
	auth         *bind.TransactOpts
}

func testSetup(t *testing.T) *testSetupData {
	transactor := testutils.MustNewSimTransactor(t)
	simulatedBackend := backends.NewSimulatedBackend(core.GenesisAlloc{
		transactor.From: {Balance: assets.Ether(1000).ToInt()},
	}, 30e6)

	// Deploy the contract
	address, _, _, err := message_hasher.DeployMessageHasher(transactor, simulatedBackend)
	require.NoError(t, err)
	simulatedBackend.Commit()

	// Setup contract client
	contract, err := message_hasher.NewMessageHasher(address, simulatedBackend)
	require.NoError(t, err)

	return &testSetupData{
		contractAddr: address,
		contract:     contract,
		sb:           simulatedBackend,
		auth:         transactor,
	}
}
