package ccipsui

import (
	"context"
	"encoding/binary"
	"math/big"
	"testing"

	"github.com/aptos-labs/aptos-go-sdk/bcs"
	"github.com/ethereum/go-ethereum/accounts/abi"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_3/message_hasher"
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipevm"
)

func TestExecutePluginCodecV2_IncludesReceiverObjectIds(t *testing.T) {
	ctx := context.Background()
	const evmSourceSelector = ccipocr3.ChainSelector(5009297550715157269)

	extraDataCodec := ccipocr3.ExtraDataCodecMap(map[string]ccipocr3.SourceChainExtraDataCodec{
		chainsel.FamilyEVM: ccipevm.ExtraDataDecoder{},
	})

	objectID := hexTo32Bytes(t, "0000000000000000000000000000000000000000000000000000000000aabbcc")
	tokenReceiver := hexTo32Bytes(t, "0000000000000000000000000000000000000000000000000000000000005678")
	receiverPackage := hexTo32Bytes(t, "0000000000000000000000000000000000000000000000000000000000001234")

	extraArgs, err := ccipevm.SerializeClientSUIExtraArgsV1(message_hasher.ClientSuiExtraArgsV1{
		GasLimit:                 big.NewInt(200000),
		AllowOutOfOrderExecution: true,
		TokenReceiver:            tokenReceiver,
		ReceiverObjectIds:        [][32]byte{objectID},
	})
	require.NoError(t, err)

	report := singleMessageExecuteReport(t, evmSourceSelector, receiverPackage, extraArgs)

	v1Codec := NewExecutePluginCodecV1(extraDataCodec)
	v2Codec := NewExecutePluginCodecV2(extraDataCodec)

	v1Bytes, err := v1Codec.Encode(ctx, report)
	require.NoError(t, err)

	v2Bytes, err := v2Codec.Encode(ctx, report)
	require.NoError(t, err)

	assert.Greater(t, len(v2Bytes), len(v1Bytes))

	des := bcs.NewDeserializer(v2Bytes)
	des.U64() // source_chain_selector
	des.ReadFixedBytes(32)
	des.U64()
	des.U64()
	des.U64()
	des.U64()
	des.ReadBytes()
	des.ReadBytes()
	des.Struct(&aptosAccount{})
	des.U256()
	des.ReadFixedBytes(32)

	ids := deserializeReceiverObjectIDs(des)
	require.NoError(t, des.Error())
	assert.Equal(t, [][32]byte{objectID}, ids)
}

func TestExecutePluginCodecV2_EmptyReceiverObjectIds(t *testing.T) {
	ctx := context.Background()
	const evmSourceSelector = ccipocr3.ChainSelector(5009297550715157269)

	extraDataCodec := ccipocr3.ExtraDataCodecMap(map[string]ccipocr3.SourceChainExtraDataCodec{
		chainsel.FamilyEVM: ccipevm.ExtraDataDecoder{},
	})

	tokenReceiver := hexTo32Bytes(t, "0000000000000000000000000000000000000000000000000000000000005678")
	receiverPackage := hexTo32Bytes(t, "0000000000000000000000000000000000000000000000000000000000001234")

	extraArgs, err := ccipevm.SerializeClientSUIExtraArgsV1(message_hasher.ClientSuiExtraArgsV1{
		GasLimit:                 big.NewInt(0),
		AllowOutOfOrderExecution: true,
		TokenReceiver:            tokenReceiver,
		ReceiverObjectIds:        [][32]byte{},
	})
	require.NoError(t, err)

	report := singleMessageExecuteReport(t, evmSourceSelector, receiverPackage, extraArgs)

	v2Bytes, err := NewExecutePluginCodecV2(extraDataCodec).Encode(ctx, report)
	require.NoError(t, err)

	des := bcs.NewDeserializer(v2Bytes)
	skipToReceiverObjectIDs(des)
	ids := deserializeReceiverObjectIDs(des)
	require.NoError(t, des.Error())
	assert.Empty(t, ids)
}

func TestExecutePluginCodecV2_EncodeDecode_Roundtrip(t *testing.T) {
	ctx := context.Background()
	const evmSourceSelector = ccipocr3.ChainSelector(5009297550715157269)

	extraDataCodec := ccipocr3.ExtraDataCodecMap(map[string]ccipocr3.SourceChainExtraDataCodec{
		chainsel.FamilyEVM: ccipevm.ExtraDataDecoder{},
	})

	objectID := hexTo32Bytes(t, "0000000000000000000000000000000000000000000000000000000000001111")
	tokenReceiver := hexTo32Bytes(t, "0000000000000000000000000000000000000000000000000000000000005678")
	receiverPackage := hexTo32Bytes(t, "0000000000000000000000000000000000000000000000000000000000001234")

	extraArgs, err := ccipevm.SerializeClientSUIExtraArgsV1(message_hasher.ClientSuiExtraArgsV1{
		GasLimit:                 big.NewInt(500000),
		AllowOutOfOrderExecution: false,
		TokenReceiver:            tokenReceiver,
		ReceiverObjectIds:        [][32]byte{objectID},
	})
	require.NoError(t, err)

	original := singleMessageExecuteReport(t, evmSourceSelector, receiverPackage, extraArgs)
	codec := NewExecutePluginCodecV2(extraDataCodec)

	encoded, err := codec.Encode(ctx, original)
	require.NoError(t, err)

	decoded, err := codec.Decode(ctx, encoded)
	require.NoError(t, err)

	require.Len(t, decoded.ChainReports, 1)
	require.Len(t, decoded.ChainReports[0].Messages, 1)

	origMsg := original.ChainReports[0].Messages[0]
	gotMsg := decoded.ChainReports[0].Messages[0]

	assert.Equal(t, origMsg.Header.MessageID, gotMsg.Header.MessageID)
	assert.Equal(t, origMsg.Header.SourceChainSelector, gotMsg.Header.SourceChainSelector)
	assert.Equal(t, origMsg.Header.DestChainSelector, gotMsg.Header.DestChainSelector)
	assert.Equal(t, origMsg.Header.SequenceNumber, gotMsg.Header.SequenceNumber)
	assert.Equal(t, origMsg.Header.Nonce, gotMsg.Header.Nonce)
	assert.Equal(t, origMsg.Sender, gotMsg.Sender)
	assert.Equal(t, origMsg.Data, gotMsg.Data)
	assert.Equal(t, ccipocr3.UnknownAddress(tokenReceiver[:]), gotMsg.Receiver)

	require.Len(t, gotMsg.TokenAmounts, 1)
	assert.Equal(t, origMsg.TokenAmounts[0].SourcePoolAddress, gotMsg.TokenAmounts[0].SourcePoolAddress)
	assert.Equal(t, origMsg.TokenAmounts[0].DestTokenAddress, gotMsg.TokenAmounts[0].DestTokenAddress)
	assert.Equal(t, origMsg.TokenAmounts[0].ExtraData, gotMsg.TokenAmounts[0].ExtraData)
	assert.Equal(t, origMsg.TokenAmounts[0].Amount, gotMsg.TokenAmounts[0].Amount)

	destGasAmount := binary.LittleEndian.Uint32(gotMsg.TokenAmounts[0].DestExecData)
	assert.Equal(t, uint32(10000), destGasAmount)

	assert.Equal(t, original.ChainReports[0].Proofs, decoded.ChainReports[0].Proofs)
}

func singleMessageExecuteReport(
	t *testing.T,
	sourceSelector ccipocr3.ChainSelector,
	receiverPackage [32]byte,
	extraArgs []byte,
) ccipocr3.ExecutePluginReport {
	t.Helper()

	destExecData, err := encodeABIUint32(10000)
	require.NoError(t, err)

	return ccipocr3.ExecutePluginReport{
		ChainReports: []ccipocr3.ExecutePluginReportSingleChain{{
			SourceChainSelector: sourceSelector,
			Messages: []ccipocr3.Message{{
				Header: ccipocr3.RampMessageHeader{
					MessageID:           hexTo32Bytes(t, "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"),
					SourceChainSelector: sourceSelector,
					DestChainSelector:   2000,
					SequenceNumber:      42,
					Nonce:               0,
				},
				Sender:    mustHexDecode(t, "8765432109fedcba8765432109fedcba87654321"),
				Data:      []byte("test payload"),
				Receiver:  ccipocr3.UnknownAddress(receiverPackage[:]),
				ExtraArgs: extraArgs,
				TokenAmounts: []ccipocr3.RampTokenAmount{{
					SourcePoolAddress: []byte("source-pool"),
					DestTokenAddress:  receiverPackage[:],
					ExtraData:         []byte{0x01},
					Amount:            ccipocr3.NewBigInt(big.NewInt(1000)),
					DestExecData:      destExecData,
				}},
			}},
			OffchainTokenData: [][][]byte{{{0xab, 0xcd}}},
			Proofs:            []ccipocr3.Bytes32{hexTo32Bytes(t, "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")},
		}},
	}
}

func skipToReceiverObjectIDs(des *bcs.Deserializer) {
	des.U64()
	des.ReadFixedBytes(32)
	des.U64()
	des.U64()
	des.U64()
	des.U64()
	des.ReadBytes()
	des.ReadBytes()
	des.Struct(&aptosAccount{})
	des.U256()
	des.ReadFixedBytes(32)
}

type aptosAccount struct {
	data [32]byte
}

func (a *aptosAccount) UnmarshalBCS(des *bcs.Deserializer) {
	copy(a.data[:], des.ReadFixedBytes(32))
}

func encodeABIUint32(v uint32) ([]byte, error) {
	args := abi.Arguments{{Type: mustABIType("uint32")}}
	return args.Pack(v)
}

func mustABIType(t string) abi.Type {
	typ, err := abi.NewType(t, "", nil)
	if err != nil {
		panic(err)
	}
	return typ
}
