package ccipevm

import (
	"bytes"
	"context"
	cryptorand "crypto/rand"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"math/rand"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	agbinary "github.com/gagliardetto/binary"
	solanago "github.com/gagliardetto/solana-go"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/message_hasher"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/latest/fee_quoter"
	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-evm/pkg/assets"
	evmtestutils "github.com/smartcontractkit/chainlink-evm/pkg/testutils"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipaptos"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipsolana"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
)

var extraDataCodec = ccipocr3.ExtraDataCodecMap(map[string]ccipocr3.SourceChainExtraDataCodec{
	chainsel.FamilyAptos:  ccipaptos.ExtraDataDecoder{},
	chainsel.FamilyEVM:    ExtraDataDecoder{},
	chainsel.FamilySolana: ccipsolana.ExtraDataDecoder{},
	chainsel.FamilySui:    ccipaptos.ExtraDataDecoder{},
})

// NOTE: these test cases are only EVM <-> EVM.
// Update these cases once we have non-EVM examples.
func TestMessageHasher_EVM2EVM(t *testing.T) {
	ctx := testutils.Context(t)
	d := testSetup(t)

	testCases := []evmExtraArgs{
		{version: "v1", gasLimit: big.NewInt(rand.Int63()), chainSelector: 5009297550715157269},                  // ETH mainnet chain selector
		{version: "v2", gasLimit: big.NewInt(rand.Int63()), allowOOO: false, chainSelector: 5009297550715157269}, // ETH mainnet chain selector
		{version: "v2", gasLimit: big.NewInt(rand.Int63()), allowOOO: true, chainSelector: 5009297550715157269},  // ETH mainnet chain selector
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("tc_%d", i), func(tt *testing.T) {
			testHasherEVM2EVM(ctx, tt, d, tc, tc.chainSelector)
		})
	}
}

func testHasherEVM2EVM(ctx context.Context, t *testing.T, d *testSetupData, evmExtraArgs evmExtraArgs, sourceChainSelector uint64) {
	ccipMsg := createEVM2EVMMessage(t, d.contract, evmExtraArgs, sourceChainSelector)

	var tokenAmounts []message_hasher.InternalAny2EVMTokenTransfer
	for _, rta := range ccipMsg.TokenAmounts {
		destGasAmount, err := abiDecodeUint32(rta.DestExecData)
		require.NoError(t, err)

		tokenAmounts = append(tokenAmounts, message_hasher.InternalAny2EVMTokenTransfer{
			SourcePoolAddress: rta.SourcePoolAddress,
			DestTokenAddress:  common.BytesToAddress(rta.DestTokenAddress),
			ExtraData:         rta.ExtraData[:],
			Amount:            rta.Amount.Int,
			DestGasAmount:     destGasAmount,
		})
	}
	evmMsg := message_hasher.InternalAny2EVMRampMessage{
		Header: message_hasher.InternalRampMessageHeader{
			MessageId:           ccipMsg.Header.MessageID,
			SourceChainSelector: uint64(ccipMsg.Header.SourceChainSelector),
			DestChainSelector:   uint64(ccipMsg.Header.DestChainSelector),
			SequenceNumber:      uint64(ccipMsg.Header.SequenceNumber),
			Nonce:               ccipMsg.Header.Nonce,
		},
		Sender:       ccipMsg.Sender,
		Receiver:     common.BytesToAddress(ccipMsg.Receiver),
		GasLimit:     evmExtraArgs.gasLimit,
		Data:         ccipMsg.Data,
		TokenAmounts: tokenAmounts,
	}

	expectedHash, err := d.contract.Hash(&bind.CallOpts{Context: ctx}, evmMsg, ccipMsg.Header.OnRamp)
	require.NoError(t, err)

	evmMsgHasher := NewMessageHasherV1(logger.Test(t), extraDataCodec)
	actualHash, err := evmMsgHasher.Hash(ctx, ccipMsg)
	require.NoError(t, err)

	require.Equal(t, fmt.Sprintf("%x", expectedHash), strings.TrimPrefix(actualHash.String(), "0x"))
}

type evmExtraArgs struct {
	version       string
	gasLimit      *big.Int
	allowOOO      bool
	chainSelector uint64
}

func createEVM2EVMMessage(t *testing.T, messageHasher *message_hasher.MessageHasher, evmExtraArgs evmExtraArgs, sourceChainSelector uint64) cciptypes.Message {
	messageID := utils.RandomBytes32()

	sourceTokenData := make([]byte, rand.Intn(2048))
	_, err := cryptorand.Read(sourceTokenData)
	require.NoError(t, err)

	sourceChain := sourceChainSelector
	seqNum := rand.Uint64()
	nonce := rand.Uint64()
	destChain := rand.Uint64()

	var extraArgsBytes []byte
	if evmExtraArgs.version == "v1" {
		extraArgsBytes, err = messageHasher.EncodeEVMExtraArgsV1(nil, message_hasher.ClientEVMExtraArgsV1{
			GasLimit: evmExtraArgs.gasLimit,
		})
		require.NoError(t, err)
	} else if evmExtraArgs.version == "v2" {
		extraArgsBytes, err = messageHasher.EncodeEVMExtraArgsV2(nil, message_hasher.ClientGenericExtraArgsV2{
			GasLimit:                 evmExtraArgs.gasLimit,
			AllowOutOfOrderExecution: evmExtraArgs.allowOOO,
		})
		require.NoError(t, err)
	} else {
		require.FailNowf(t, "unknown extra args version", "version: %s", evmExtraArgs.version)
	}

	messageData := make([]byte, rand.Intn(2048))
	_, err = cryptorand.Read(messageData)
	require.NoError(t, err)

	numTokens := rand.Intn(10)
	var sourceTokenDatas [][]byte
	for range numTokens {
		sourceTokenDatas = append(sourceTokenDatas, sourceTokenData)
	}

	var tokenAmounts []cciptypes.RampTokenAmount
	for i := 0; i < len(sourceTokenDatas); i++ {
		extraData := utils.RandomBytes32()
		encodedDestExecData, err := utils.ABIEncode(`[{ "type": "uint32" }]`, rand.Uint32())
		require.NoError(t, err)
		tokenAmounts = append(tokenAmounts, cciptypes.RampTokenAmount{
			SourcePoolAddress: abiEncodedAddress(t),
			DestTokenAddress:  abiEncodedAddress(t),
			ExtraData:         extraData[:],
			Amount:            cciptypes.NewBigInt(big.NewInt(0).SetUint64(rand.Uint64())),
			DestExecData:      encodedDestExecData,
		})
	}

	return cciptypes.Message{
		Header: cciptypes.RampMessageHeader{
			MessageID:           messageID,
			SourceChainSelector: cciptypes.ChainSelector(sourceChain),
			DestChainSelector:   cciptypes.ChainSelector(destChain),
			SequenceNumber:      cciptypes.SeqNum(seqNum),
			Nonce:               nonce,
			OnRamp:              abiEncodedAddress(t),
		},
		Sender:         abiEncodedAddress(t),
		Receiver:       abiEncodedAddress(t),
		Data:           messageData,
		TokenAmounts:   tokenAmounts,
		FeeToken:       abiEncodedAddress(t),
		FeeTokenAmount: cciptypes.NewBigInt(big.NewInt(0).SetUint64(rand.Uint64())),
		ExtraArgs:      extraArgsBytes,
	}
}

func abiEncodedAddress(t *testing.T) []byte {
	addr := utils.RandomAddress()
	encoded, err := utils.ABIEncode(`[{"type": "address"}]`, addr)
	require.NoError(t, err)
	return encoded
}

type testSetupData struct {
	contractAddr common.Address
	contract     *message_hasher.MessageHasher
	sb           *backends.SimulatedBackend
	auth         *bind.TransactOpts
}

func testSetup(t *testing.T) *testSetupData {
	transactor := evmtestutils.MustNewSimTransactor(t)
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

func TestMessagerHasher_againstRmnSharedVector(t *testing.T) {
	transactor := evmtestutils.MustNewSimTransactor(t)
	backend := backends.NewSimulatedBackend(types.GenesisAlloc{
		transactor.From: {Balance: assets.Ether(1000).ToInt()},
	}, 30e6)

	msghasherAddr, _, _, err := message_hasher.DeployMessageHasher(transactor, backend)
	require.NoError(t, err)
	backend.Commit()

	msghasher, err := message_hasher.NewMessageHasher(msghasherAddr, backend)
	require.NoError(t, err)

	t.Run("vec1", func(t *testing.T) {
		const (
			messageID           = "c6f553ab71282f01324bbdbcc82e22a7e66efbcd108881ecc4cdbd728aed9b1e"
			onRampAddress       = "0000000000000000000000007a2088a1bfc9d81c55368ae168c2c02570cb814f"
			dataField           = "68656c6c6f"
			receiverAddress     = "677df0cb865368207999f2862ece576dc56d8df6"
			extraArgs           = "181dcf100000000000000000000000000000000000000000000000000000000000030d400000000000000000000000000000000000000000000000000000000000000000"
			senderAddress       = "f39fd6e51aad88f6f4ce6ab8827279cfffb92266"
			feeToken            = "9fe46736679d2d9a65f0992f2272de9f3c7fa6e0"
			sourceChainSelector = 3379446385462418246
			destChainSelector   = 12922642891491394802
			expectedMsgHash     = "0x1c61fef7a3dd153943419c1101031316ed7b7a3d75913c34cbe8628033f5924f"
		)

		var (
			msg = cciptypes.Message{
				Header: cciptypes.RampMessageHeader{
					MessageID:           cciptypes.Bytes32(common.Hex2Bytes(messageID)),
					SourceChainSelector: sourceChainSelector,
					DestChainSelector:   destChainSelector,
					SequenceNumber:      1,
					Nonce:               1,
					MsgHash:             cciptypes.Bytes32{},
					OnRamp:              common.HexToAddress(onRampAddress).Bytes(),
				},
				Sender:       common.HexToAddress(senderAddress).Bytes(),
				Data:         common.Hex2Bytes(dataField),
				Receiver:     common.Hex2Bytes(receiverAddress),
				ExtraArgs:    common.Hex2Bytes(extraArgs),
				FeeToken:     common.HexToAddress(feeToken).Bytes(),
				TokenAmounts: []cciptypes.RampTokenAmount{},
			}
			any2EVMMessage = ccipMsgToAny2EVMMessage(t, msg, sourceChainSelector)
		)

		onchainHash, err := msghasher.Hash(&bind.CallOpts{
			Context: t.Context(),
		}, any2EVMMessage, common.LeftPadBytes(msg.Header.OnRamp, 32))
		require.NoError(t, err)

		h := NewMessageHasherV1(logger.Test(t), extraDataCodec)
		msgH, err := h.Hash(t.Context(), msg)
		require.NoError(t, err)
		require.Equal(t, expectedMsgHash, msgH.String())
		require.Equal(t, onchainHash, [32]byte(msgH), "my hash and onchain hash should match")
	})

	t.Run("vec2", func(t *testing.T) {
		// source chain tx: https://sepolia.etherscan.io/tx/0x3b64b5cb2c972a3f5064801187f17360c2025fbcc51e11b67b25c7949daeec24#eventlog
		var (
			// header fields
			messageID           = mustBytes32FromString(t, "0xcdad95e113e35cf691295c1f42455d41062ba9a1b96a6280c1a5a678ef801721")
			destChainSelector   = cciptypes.ChainSelector(3478487238524512106)  // arb sepolia
			sourceChainSelector = cciptypes.ChainSelector(16015286601757825753) // sepolia
			sequenceNumber      = cciptypes.SeqNum(386)
			nonce               = uint64(1)
			// message fields
			// sender is parsed unpadded since its emitted unpadded from EVM.
			senderAddress = cciptypes.UnknownAddress(hexutil.MustDecode("0x269895AC2a2eC6e1Df37F68AcfbBDa53e62b71B1"))
			// onRampAddress is parsed padded because its set as a padded address in the offRamp
			onRampAddress = hexutil.MustDecode("0x00000000000000000000000089559ce6904d4c4B0f6aaB9065Ad02B1ed531Be4")
			dataField     = "0x"
			// receiver address is parsed padded because its emitted as padded from EVM.
			receiverAddress = cciptypes.UnknownAddress(hexutil.MustDecode("0x000000000000000000000000269895ac2a2ec6e1df37f68acfbbda53e62b71b1"))
			// extraArgs always abi-encoded
			extraArgs = hexutil.MustDecode("0x181dcf100000000000000000000000000000000000000000000000000000000000030d400000000000000000000000000000000000000000000000000000000000000000")
			// feeToken is parsed unpadded since its emitted unpadded from EVM.
			// however, it isn't used in the hash. its just set for completion.
			feeToken       = common.HexToAddress("0x097D90c9d3E0B50Ca60e1ae45F6A81010f9FB534")
			feeTokenAmount = big.NewInt(114310554250104)
			feeValueJuels  = big.NewInt(16499514422603741)
			tokenAmounts   = []cciptypes.RampTokenAmount{
				{
					// parsed unpadded since its emitted unpadded from EVM.
					SourcePoolAddress: cciptypes.UnknownAddress(hexutil.MustDecode("0xBBE734cAB186C0988CFBAfdFdbe442979a0c8697")),
					// parsed padded because its emitted padded from EVM.
					DestTokenAddress: cciptypes.UnknownAddress(hexutil.MustDecode("0x000000000000000000000000b8d6a6a41d5dd732aec3c438e91523b7613b963b")),
					// extra data always abi-encoded
					ExtraData: cciptypes.Bytes(hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000000012")),
					Amount:    cciptypes.NewBigInt(big.NewInt(100000000000000000)),
					// dest exec data always abi-encoded
					DestExecData: cciptypes.Bytes(hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000001e848")),
				},
			}

			msg = cciptypes.Message{
				Header: cciptypes.RampMessageHeader{
					MessageID:           messageID,
					SourceChainSelector: sourceChainSelector,
					DestChainSelector:   destChainSelector,
					SequenceNumber:      sequenceNumber,
					Nonce:               nonce,
					MsgHash:             cciptypes.Bytes32{},
					OnRamp:              onRampAddress,
				},
				Sender:         senderAddress,
				Data:           hexutil.MustDecode(dataField),
				Receiver:       receiverAddress,
				ExtraArgs:      extraArgs,
				FeeToken:       feeToken.Bytes(),
				FeeTokenAmount: cciptypes.NewBigInt(feeTokenAmount),
				FeeValueJuels:  cciptypes.NewBigInt(feeValueJuels),
				TokenAmounts:   tokenAmounts,
			}

			any2EVMMessage = ccipMsgToAny2EVMMessage(t, msg, sourceChainSelector)
		)

		const (
			rmnMsgHash = "0xb6ea678f918293745bfb8db05d79dcf08986c7da3e302ac5f6782618a6f11967"
		)

		h := NewMessageHasherV1(logger.Test(t), extraDataCodec)
		msgH, err := h.Hash(t.Context(), msg)
		require.NoError(t, err)

		msgHashOnchain, err := msghasher.Hash(&bind.CallOpts{
			Context: t.Context(),
		}, any2EVMMessage, onRampAddress)
		require.NoError(t, err)

		t.Logf("rmn hash: %s, onchain hash: %s, my hash: %s", rmnMsgHash, hexutil.Encode(msgHashOnchain[:]), msgH.String())
		require.Equal(t, msgHashOnchain, [32]byte(msgH), "my hash and onchain hash should match")
		require.Equal(t, rmnMsgHash, msgH.String(), "rmn hash and my hash should match")
	})

	t.Run("solana", func(t *testing.T) {
		key, err := solanago.NewRandomPrivateKey()
		require.NoError(t, err)

		// Borsh encoded extra args
		ea := fee_quoter.GenericExtraArgsV2{
			GasLimit:                 agbinary.Uint128{Lo: 5000, Hi: 0},
			AllowOutOfOrderExecution: false,
		}

		var extraArgsbuf bytes.Buffer
		encoder := agbinary.NewBorshEncoder(&extraArgsbuf)
		err = ea.MarshalWithEncoder(encoder)
		require.NoError(t, err)

		destGasAmount := uint32(10)
		destExecData := make([]byte, 4)
		binary.LittleEndian.PutUint32(destExecData, destGasAmount)

		// evmExtraArgsV2Tag from SVM on-chain contract
		// https://github.com/smartcontractkit/chainlink-ccip/blob/1b2ee24da54bddef8f3943dc84102686f2890f87/chains/solana/contracts/programs/ccip-router/src/extra_args.rs#L9
		evmExtraArgsV2 := hexutil.MustDecode("0x181dcf10")

		var (
			// header fields
			messageID           = mustBytes32FromString(t, "0xcdad95e113e35cf691295c1f42455d41062ba9a1b96a6280c1a5a678ef801721")
			sourceChainSelector = cciptypes.ChainSelector(124615329519749607)   // solana
			destChainSelector   = cciptypes.ChainSelector(16015286601757825753) // sepolia
			sequenceNumber      = cciptypes.SeqNum(386)
			nonce               = uint64(1)
			// message fields
			// sender is parsed unpadded since its emitted unpadded from SVM.
			senderAddress = cciptypes.UnknownAddress(key.PublicKey().Bytes())
			// onRampAddress is parsed padded because its set as a padded address in the offRamp
			onRampAddress = key.PublicKey().Bytes()
			dataField     = "0x"
			// receiver address is parsed padded because its emitted as padded from SVM.
			receiverAddress = cciptypes.UnknownAddress(hexutil.MustDecode("0x000000000000000000000000269895ac2a2ec6e1df37f68acfbbda53e62b71b1"))
			feeTokenAmount  = big.NewInt(114310554250104)
			feeValueJuels   = big.NewInt(16499514422603741)
			tokenAmounts    = []cciptypes.RampTokenAmount{
				{
					// parsed unpadded since its emitted unpadded from SVM.
					SourcePoolAddress: cciptypes.UnknownAddress(key.PublicKey().Bytes()),
					// parsed padded because its emitted padded from SVM.
					DestTokenAddress: cciptypes.UnknownAddress(hexutil.MustDecode("0x000000000000000000000000b8d6a6a41d5dd732aec3c438e91523b7613b963b")),
					// extra data always abi-encoded
					ExtraData: cciptypes.Bytes(hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000000012")),
					Amount:    cciptypes.NewBigInt(big.NewInt(100000000000000000)),
					// dest exec data always abi-encoded
					DestExecData: destExecData,
				},
			}
			msg = cciptypes.Message{
				Header: cciptypes.RampMessageHeader{
					MessageID:           messageID,
					SourceChainSelector: sourceChainSelector,
					DestChainSelector:   destChainSelector,
					SequenceNumber:      sequenceNumber,
					Nonce:               nonce,
					MsgHash:             cciptypes.Bytes32{},
					OnRamp:              onRampAddress,
				},
				Sender:         senderAddress,
				Data:           hexutil.MustDecode(dataField),
				Receiver:       receiverAddress,
				ExtraArgs:      append(evmExtraArgsV2, extraArgsbuf.Bytes()...),
				FeeToken:       key.PublicKey().Bytes(),
				FeeTokenAmount: cciptypes.NewBigInt(feeTokenAmount),
				FeeValueJuels:  cciptypes.NewBigInt(feeValueJuels),
				TokenAmounts:   tokenAmounts,
			}

			any2EVMMessage = ccipMsgToAny2EVMMessage(t, msg, sourceChainSelector)
		)

		//const (
		//	rmnMsgHash = "0xb6ea678f918293745bfb8db05d79dcf08986c7da3e302ac5f6782618a6f11967"
		//)

		h := NewMessageHasherV1(logger.Test(t), extraDataCodec)
		msgH, err := h.Hash(t.Context(), msg)
		require.NoError(t, err)

		msgHashOnchain, err := msghasher.Hash(&bind.CallOpts{
			Context: t.Context(),
		}, any2EVMMessage, onRampAddress)
		require.NoError(t, err)

		//t.Logf("rmn hash: %s, onchain hash: %s, my hash: %s", rmnMsgHash, hexutil.Encode(msgHashOnchain[:]), msgH.String())
		require.Equal(t, msgHashOnchain, [32]byte(msgH), "my hash and onchain hash should match")
		//require.Equal(t, rmnMsgHash, msgH.String(), "rmn hash and my hash should match")
	})

	t.Run("other vectors", func(t *testing.T) {
		// These test vectors are from real ccip transactions on sepolia.
		// onramp address: 0x89559ce6904d4c4b0f6aab9065ad02b1ed531be4
		// sequence numbers 386 to 419.
		var msgs []cciptypes.Message
		data, err := ioutil.ReadFile("msgs_test_vector.json")
		require.NoError(t, err)

		err = json.Unmarshal(data, &msgs)
		require.NoError(t, err)

		msgHasher := NewMessageHasherV1(logger.Test(t), extraDataCodec)

		for _, msg := range msgs {
			any2EVMMessage := ccipMsgToAny2EVMMessage(t, msg, msg.Header.SourceChainSelector)

			onchainHash, err := msghasher.Hash(&bind.CallOpts{
				Context: t.Context(),
			}, any2EVMMessage, common.LeftPadBytes(msg.Header.OnRamp, 32))
			require.NoError(t, err)

			myHash, err := msgHasher.Hash(t.Context(), msg)
			require.NoError(t, err)

			t.Logf("onchain hash: %s, my hash: %s", hexutil.Encode(onchainHash[:]), myHash.String())
			require.Equal(t, [32]byte(myHash), onchainHash, "my hash and onchain hash should match")
		}
	})
}

func ccipMsgToAny2EVMMessage(t *testing.T, msg cciptypes.Message, sourceSelector cciptypes.ChainSelector) message_hasher.InternalAny2EVMRampMessage {
	any2EVMMessage, err := tryCCIPMsgToAny2EVMMessage(msg, sourceSelector)
	require.NoError(t, err)
	return any2EVMMessage
}

// tryCCIPMsgToAny2EVMMessage is a version of ccipMsgToAny2EVMMessage that returns errors instead of calling t.Fatal.
// This is useful for fuzz testing where we expect some inputs to be invalid.
func tryCCIPMsgToAny2EVMMessage(msg cciptypes.Message, sourceSelector cciptypes.ChainSelector) (message_hasher.InternalAny2EVMRampMessage, error) {
	sourceChainFamily, err := chainsel.GetSelectorFamily(uint64(sourceSelector))
	if err != nil {
		return message_hasher.InternalAny2EVMRampMessage{}, fmt.Errorf("get source chain family: %w", err)
	}

	var tokenAmounts []message_hasher.InternalAny2EVMTokenTransfer
	for _, rta := range msg.TokenAmounts {
		decodedMap, err := extraDataCodec.DecodeTokenAmountDestExecData(rta.DestExecData, sourceSelector)
		if err != nil {
			return message_hasher.InternalAny2EVMRampMessage{}, fmt.Errorf("could not decode token amount dest exec data: %w", err)
		}
		gasAmount, err := extractDestGasAmountFromMap(decodedMap)
		if err != nil {
			return message_hasher.InternalAny2EVMRampMessage{}, fmt.Errorf("could not extract dest gas amount: %w", err)
		}
		destTokenAddress, err := abiDecodeAddress(rta.DestTokenAddress)
		if err != nil {
			return message_hasher.InternalAny2EVMRampMessage{}, fmt.Errorf("could not decode dest token address: %w", err)
		}

		var sourcePoolAddr []byte
		if sourceChainFamily == chainsel.FamilyEVM {
			// from https://github.com/smartcontractkit/chainlink/blob/e036012d5b562f5c30c5a87898239ba59aeb2f7b/contracts/src/v0.8/ccip/pools/TokenPool.sol#L84
			// remote pool addresses are abi-encoded addresses if the remote chain is EVM.
			sourcePoolAddr, err = abiEncodeAddress(common.BytesToAddress(rta.SourcePoolAddress))
			if err != nil {
				return message_hasher.InternalAny2EVMRampMessage{}, fmt.Errorf("abi encode source pool address: %w", err)
			}
		} else {
			sourcePoolAddr = rta.SourcePoolAddress
		}

		tokenAmounts = append(tokenAmounts, message_hasher.InternalAny2EVMTokenTransfer{
			SourcePoolAddress: sourcePoolAddr,
			DestTokenAddress:  destTokenAddress,
			ExtraData:         rta.ExtraData[:],
			Amount:            rta.Amount.Int,
			DestGasAmount:     gasAmount,
		})
	}

	decodedMap, err := extraDataCodec.DecodeExtraArgs(msg.ExtraArgs, sourceSelector)
	if err != nil {
		return message_hasher.InternalAny2EVMRampMessage{}, fmt.Errorf("could not decode extra args: %w", err)
	}
	gasLimit, err := parseExtraArgsMap(decodedMap)
	if err != nil {
		return message_hasher.InternalAny2EVMRampMessage{}, fmt.Errorf("could not parse extra args map: %w", err)
	}

	return message_hasher.InternalAny2EVMRampMessage{
		Header: message_hasher.InternalRampMessageHeader{
			MessageId:           msg.Header.MessageID,
			SourceChainSelector: uint64(msg.Header.SourceChainSelector),
			DestChainSelector:   uint64(msg.Header.DestChainSelector),
			SequenceNumber:      uint64(msg.Header.SequenceNumber),
			Nonce:               msg.Header.Nonce,
		},
		Sender:       common.LeftPadBytes(msg.Sender, 32),
		Data:         msg.Data,
		Receiver:     common.BytesToAddress(msg.Receiver),
		GasLimit:     gasLimit,
		TokenAmounts: tokenAmounts,
	}, nil
}

func mustBytes32FromString(t *testing.T, str string) cciptypes.Bytes32 {
	t.Helper()
	b, err := cciptypes.NewBytes32FromString(str)
	require.NoError(t, err)
	return b
}

func FuzzMessageHasher(f *testing.F) {
	transactor := evmtestutils.MustNewSimTransactor(f)
	//nolint:staticcheck // all gethwrappers still use legacy
	backend := backends.NewSimulatedBackend(types.GenesisAlloc{
		transactor.From: {Balance: assets.Ether(1000).ToInt()},
	}, 30e6)

	msghasherAddr, _, _, err := message_hasher.DeployMessageHasher(transactor, backend)
	require.NoError(f, err)
	backend.Commit()

	msghasher, err := message_hasher.NewMessageHasher(msghasherAddr, backend)
	require.NoError(f, err)

	// Seed with data from vec1 test case
	f.Add(
		common.Hex2Bytes("c6f553ab71282f01324bbdbcc82e22a7e66efbcd108881ecc4cdbd728aed9b1e"),
		uint64(3379446385462418246),
		uint64(12922642891491394802),
		uint64(1),
		uint64(1),
		common.HexToAddress("0000000000000000000000007a2088a1bfc9d81c55368ae168c2c02570cb814f").Bytes(),
		common.HexToAddress("f39fd6e51aad88f6f4ce6ab8827279cfffb92266").Bytes(),
		common.Hex2Bytes("68656c6c6f"),
		common.HexToAddress("677df0cb865368207999f2862ece576dc56d8df6").Bytes(),
		common.Hex2Bytes("181dcf100000000000000000000000000000000000000000000000000000000000030d400000000000000000000000000000000000000000000000000000000000000000"),
		// No token amounts for vec1
		[]byte{},
		[]byte{},
		[]byte{},
		[]byte{},
		[]byte{},
	)

	// Seed with data from vec2 test case, which includes a token amount.
	f.Add(
		hexutil.MustDecode("0xcdad95e113e35cf691295c1f42455d41062ba9a1b96a6280c1a5a678ef801721"),
		uint64(16015286601757825753),
		uint64(3478487238524512106),
		uint64(386),
		uint64(1),
		hexutil.MustDecode("0x00000000000000000000000089559ce6904d4c4B0f6aaB9065Ad02B1ed531Be4"),
		hexutil.MustDecode("0x269895AC2a2eC6e1Df37F68AcfbBDa53e62b71B1"),
		[]byte{}, // empty data field
		hexutil.MustDecode("0x000000000000000000000000269895ac2a2ec6e1df37f68acfbbda53e62b71b1"),
		hexutil.MustDecode("0x181dcf100000000000000000000000000000000000000000000000000000000000030d400000000000000000000000000000000000000000000000000000000000000000"),
		// Token amount from vec2
		[]byte(cciptypes.UnknownAddress(hexutil.MustDecode("0xBBE734cAB186C0988CFBAfdFdbe442979a0c8697"))),
		[]byte(cciptypes.UnknownAddress(hexutil.MustDecode("0x000000000000000000000000b8d6a6a41d5dd732aec3c438e91523b7613b963b"))),
		[]byte(cciptypes.Bytes(hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000000012"))),
		big.NewInt(100000000000000000).Bytes(),
		[]byte(cciptypes.Bytes(hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000001e848"))),
	)

	f.Fuzz(func(t *testing.T,
		messageID []byte,
		sourceChainSelector uint64,
		destChainSelector uint64,
		sequenceNumber uint64,
		nonce uint64,
		onRamp []byte,
		sender []byte,
		data []byte,
		receiver []byte,
		extraArgs []byte,
		// Fuzzed token amount fields
		tokenSourcePoolAddress []byte,
		tokenDestTokenAddress []byte,
		tokenExtraData []byte,
		tokenAmountBytes []byte,
		tokenDestExecData []byte,
	) {
		h := NewMessageHasherV1(logger.Test(t), extraDataCodec)
		if len(messageID) != 32 {
			t.Skip("messageID must be 32 bytes")
		}

		msg := cciptypes.Message{
			Header: cciptypes.RampMessageHeader{
				MessageID:           ([32]byte)(messageID),
				SourceChainSelector: cciptypes.ChainSelector(sourceChainSelector),
				DestChainSelector:   cciptypes.ChainSelector(destChainSelector),
				SequenceNumber:      cciptypes.SeqNum(sequenceNumber),
				Nonce:               nonce,
				OnRamp:              onRamp,
			},
			Sender:       sender,
			Data:         data,
			Receiver:     receiver,
			ExtraArgs:    extraArgs,
			TokenAmounts: []cciptypes.RampTokenAmount{},
		}

		// If any token data is present, try to add a token amount.
		// This allows the fuzzer to explore both messages with and without tokens.
		if len(tokenSourcePoolAddress) > 0 || len(tokenDestTokenAddress) > 0 || len(tokenExtraData) > 0 || len(tokenAmountBytes) > 0 || len(tokenDestExecData) > 0 {
			msg.TokenAmounts = []cciptypes.RampTokenAmount{
				{
					SourcePoolAddress: tokenSourcePoolAddress,
					DestTokenAddress:  tokenDestTokenAddress,
					ExtraData:         tokenExtraData,
					Amount:            cciptypes.NewBigInt(new(big.Int).SetBytes(tokenAmountBytes)),
					DestExecData:      tokenDestExecData,
				},
			}
		}

		// It's expected that many fuzzed inputs will be invalid for decoding.
		any2EVMMessage, err := tryCCIPMsgToAny2EVMMessage(msg, msg.Header.SourceChainSelector)
		if err != nil {
			t.Skipf("skipping invalid message: %v", err)
		}

		// The contract call can also revert for invalid inputs that pass our initial decoding.
		onchainHash, err := msghasher.Hash(&bind.CallOpts{
			Context: t.Context(),
		}, any2EVMMessage, common.LeftPadBytes(msg.Header.OnRamp, 32))
		if err != nil {
			t.Skipf("on-chain hashing failed: %v", err)
		}

		myHash, err := h.Hash(t.Context(), msg)
		if err != nil {
			t.Skipf("off-chain hashing failed: %v", err)
		}

		require.Equal(t, onchainHash, [32]byte(myHash), "my hash and onchain hash should match")
	})
}
