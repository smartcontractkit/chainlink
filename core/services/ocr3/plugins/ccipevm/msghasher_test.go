package ccipevm

import (
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
	"github.com/ethereum/go-ethereum/crypto"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/message_hasher"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMessageHasher_e2e(t *testing.T) {
	// Deploy messageHasher contract
	ctx := testutils.Context(t)
	d := testSetup(t)

	// Setup random msg data
	metadataHash := utils.RandomBytes32()

	sourceTokenData := make([]byte, rand.Intn(2048))
	_, err := cryptorand.Read(sourceTokenData)
	assert.NoError(t, err)

	sourceChain := rand.Uint64()
	seqNum := rand.Uint64()
	chainFeeLimit := rand.Uint64()
	nonce := rand.Uint64()
	strict := rand.Intn(2) == 1
	feeTokenAmount := rand.Uint64()

	data := make([]byte, rand.Intn(2048))
	_, err = cryptorand.Read(data)
	assert.NoError(t, err)

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
	ccipMsg := cciptypes.CCIPMsg{
		CCIPMsgBaseDetails: cciptypes.CCIPMsgBaseDetails{
			SourceChain: cciptypes.ChainSelector(sourceChain),
			SeqNum:      cciptypes.SeqNum(seqNum),
		},
		ChainFeeLimit:   cciptypes.NewBigInt(big.NewInt(0).SetUint64(chainFeeLimit)),
		Nonce:           nonce,
		Sender:          types.Account(utils.RandomAddress().String()),
		Receiver:        types.Account(utils.RandomAddress().String()),
		Strict:          strict,
		FeeToken:        types.Account(utils.RandomAddress().String()),
		FeeTokenAmount:  cciptypes.NewBigInt(big.NewInt(0).SetUint64(feeTokenAmount)),
		Data:            data,
		TokenAmounts:    tokenAmounts,
		SourceTokenData: sourceTokenDatas,
	}

	evmTokenAmounts := make([]message_hasher.ClientEVMTokenAmount, 0, len(ccipMsg.TokenAmounts))
	for _, ta := range ccipMsg.TokenAmounts {
		evmTokenAmounts = append(evmTokenAmounts, message_hasher.ClientEVMTokenAmount{
			Token:  common.HexToAddress(string(ta.Token)),
			Amount: ta.Amount,
		})
	}
	evmMsg := message_hasher.InternalEVM2EVMMessage{
		SourceChainSelector: uint64(ccipMsg.SourceChain),
		Sender:              common.HexToAddress(string(ccipMsg.Sender)),
		Receiver:            common.HexToAddress(string(ccipMsg.Receiver)),
		SequenceNumber:      uint64(ccipMsg.SeqNum),
		GasLimit:            ccipMsg.ChainFeeLimit.Int,
		Strict:              ccipMsg.Strict,
		Nonce:               ccipMsg.Nonce,
		FeeToken:            common.HexToAddress(string(ccipMsg.FeeToken)),
		FeeTokenAmount:      ccipMsg.FeeTokenAmount.Int,
		Data:                ccipMsg.Data,
		TokenAmounts:        evmTokenAmounts,
		SourceTokenData:     ccipMsg.SourceTokenData,
	}

	h, err := d.contract.Hash(&bind.CallOpts{Context: ctx}, evmMsg, metadataHash)
	assert.NoError(t, err)

	evmMsgHasher := NewMessageHasherV1(metadataHash)
	h2, err := evmMsgHasher.Hash(ctx, ccipMsg)
	assert.NoError(t, err)

	assert.Equal(t, fmt.Sprintf("%x", h), strings.TrimPrefix(h2.String(), "0x"))
}

type testSetupData struct {
	contractAddr common.Address
	contract     *message_hasher.MessageHasher
	sb           *backends.SimulatedBackend
	auth         *bind.TransactOpts
}

const chainID = 1337

func testSetup(t *testing.T) *testSetupData {
	// Generate a new key pair for the simulated account
	privateKey, err := crypto.GenerateKey()
	assert.NoError(t, err)
	// Set up the genesis account with balance
	blnc, ok := big.NewInt(0).SetString("999999999999999999999999999999999999", 10)
	assert.True(t, ok)
	alloc := map[common.Address]core.GenesisAccount{crypto.PubkeyToAddress(privateKey.PublicKey): {Balance: blnc}}
	simulatedBackend := backends.NewSimulatedBackend(alloc, 0)
	// Create a transactor

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(chainID))
	assert.NoError(t, err)
	auth.GasLimit = uint64(0)

	// Deploy the contract
	address, _, _, err := message_hasher.DeployMessageHasher(auth, simulatedBackend)
	assert.NoError(t, err)
	simulatedBackend.Commit()

	// Setup contract client
	contract, err := message_hasher.NewMessageHasher(address, simulatedBackend)
	assert.NoError(t, err)

	return &testSetupData{
		contractAddr: address,
		contract:     contract,
		sb:           simulatedBackend,
		auth:         auth,
	}
}

func TestMessageHasher_Hash(t *testing.T) {
	ctx := testutils.Context(t)

	largeNumber, ok := big.NewInt(0).SetString("1000000000000000000", 10) //1e18
	require.True(t, ok)

	msgData, err := hex.DecodeString("64617461")
	require.NoError(t, err)

	sourceTokenData1, err := hex.DecodeString("000000000000000000000000000000000000000000000000000000000000002000" +
		"000000000000000000000000000000000000000000000000000000000000600000000000000000000000000000000000000000000" +
		"0000000000000000000a000000000000000000000000000000000000000000000000000000000000000e00000000000000000000" +
		"0000000000000000000000000000000000000000000200000000000000000000000009e7218a11a2cda657ae50bd9cc5f953174aa" +
		"e2a50000000000000000000000000000000000000000000000000000000000000020000000000000000000000000e4eebe19216af8" +
		"6b9a996f53bd80b8365f832be80000000000000000000000000000000000000000000000000000000000000000")
	require.NoError(t, err)

	sourceTokenData2, err := hex.DecodeString("000000000000000000000000000000000000000000000000000000000000002" +
		"00000000000000000000000000000000000000000000000000000000000000060000000000000000000000000000000000000000" +
		"00000000000000000000000a000000000000000000000000000000000000000000000000000000000000000e00000000000000000" +
		"000000000000000000000000000000000000000000000020000000000000000000000000e2c2bb2f43b91f65b5519708e34031039" +
		"4c72d8f0000000000000000000000000000000000000000000000000000000000000020000000000000000000000000bd2f7046d10" +
		"59abfe5316b48f050684a4676710f0000000000000000000000000000000000000000000000000000000000000000")
	require.NoError(t, err)

	// metadataHash used in this test is copied from on-chain tests
	// keccak256(abi.encode(Internal.EVM_2_EVM_MESSAGE_HASH, i_chainSelector, destChainSelector, address(this)))
	metadataHash := [32]byte{39, 130, 244, 70, 94, 31, 113, 169, 251, 136, 123, 6, 255, 77, 50, 91, 73,
		144, 94, 70, 13, 16, 47, 1, 171, 201, 40, 185, 144, 12, 103, 129}

	testCases := []struct {
		name   string
		msg    cciptypes.CCIPMsg
		exp    string
		expErr bool
	}{
		{
			name: "empty msg",
			msg: cciptypes.CCIPMsg{
				ChainFeeLimit:  cciptypes.NewBigIntFromInt64(0),
				FeeTokenAmount: cciptypes.NewBigIntFromInt64(0),
			},
			exp:    "0x3682d9965b91efa44c6274446a362dca2ea526bf3858bf54d52ec56be716f6be",
			expErr: false,
		},
		{
			name: "base msg",
			msg: cciptypes.CCIPMsg{
				CCIPMsgBaseDetails: cciptypes.CCIPMsgBaseDetails{
					SourceChain: 1,
					SeqNum:      1,
				},
				ChainFeeLimit:   cciptypes.NewBigIntFromInt64(400000),
				Nonce:           1,
				Sender:          "0x00007e64E1fB0C487F25dd6D3601ff6aF8d32e4e",
				Receiver:        "0x00007e64E1fB0C487F25dd6D3601ff6aF8d32e4e",
				Strict:          false,
				FeeToken:        "0xcE4ec7b524851E51d5C55eeFbBb8E58E8Ce2515F",
				FeeTokenAmount:  cciptypes.NewBigIntFromInt64(1234567890),
				Data:            []byte{},
				TokenAmounts:    []cciptypes.TokenAmount{},
				SourceTokenData: [][]byte{},
				Metadata:        cciptypes.CCIPMsgMetadata{},
			},
			exp:    "0x23bf76c493e9bf58346b7cac0e9f357f5879f3d673819e5f27fc443cf9c907b9",
			expErr: false,
		},
		{
			name: "full msg",
			msg: cciptypes.CCIPMsg{
				CCIPMsgBaseDetails: cciptypes.CCIPMsgBaseDetails{
					SourceChain: 1,
					SeqNum:      1,
				},
				ChainFeeLimit:  cciptypes.NewBigIntFromInt64(400000),
				Nonce:          1,
				Sender:         "0x00007e64E1fB0C487F25dd6D3601ff6aF8d32e4e",
				Receiver:       "0x00007e64E1fB0C487F25dd6D3601ff6aF8d32e4e",
				Strict:         false,
				FeeToken:       "0xcE4ec7b524851E51d5C55eeFbBb8E58E8Ce2515F",
				FeeTokenAmount: cciptypes.NewBigIntFromInt64(1234567890),
				Data:           msgData,
				TokenAmounts: []cciptypes.TokenAmount{
					{
						Token:  "0xcE4ec7b524851E51d5C55eeFbBb8E58E8Ce2515F",
						Amount: largeNumber,
					},
				},
				SourceTokenData: [][]byte{sourceTokenData1},
				Metadata:        cciptypes.CCIPMsgMetadata{},
			},
			exp:    "0xe04ade4e6a1121155ca1e89f17c9df6c9236fdcfdf38b97594594d0540345d60",
			expErr: false,
		},
		{
			name: "full msg 2 - two source token data items",
			msg: cciptypes.CCIPMsg{
				CCIPMsgBaseDetails: cciptypes.CCIPMsgBaseDetails{
					SourceChain: 1,
					SeqNum:      1,
				},
				ChainFeeLimit:  cciptypes.NewBigIntFromInt64(400000),
				Nonce:          1,
				Sender:         "0x00007e64E1fB0C487F25dd6D3601ff6aF8d32e4e",
				Receiver:       "0x00007e64E1fB0C487F25dd6D3601ff6aF8d32e4e",
				Strict:         false,
				FeeToken:       "0xcE4ec7b524851E51d5C55eeFbBb8E58E8Ce2515F",
				FeeTokenAmount: cciptypes.NewBigIntFromInt64(1234567890),
				Data:           msgData,
				TokenAmounts: []cciptypes.TokenAmount{
					{
						Token:  "0xcE4ec7b524851E51d5C55eeFbBb8E58E8Ce2515F",
						Amount: largeNumber,
					},
					{
						Token:  "0x3c78e47de47B765dcEE2F30F31B3CF5F10B42d1F",
						Amount: largeNumber,
					},
				},
				SourceTokenData: [][]byte{sourceTokenData1, sourceTokenData2},
				Metadata:        cciptypes.CCIPMsgMetadata{},
			},
			exp:    "0x30123234e5d9e0cd94610e83be2e0128167b9ab072e8e9450f1f7704b9901589",
			expErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := NewMessageHasherV1(metadataHash)
			hash, err := m.Hash(ctx, tc.msg)
			if tc.expErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.exp, hash.String())
		})
	}
}
