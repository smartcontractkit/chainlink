package ccipdata

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_offramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_onramp"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/hashlib"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

var (
	// Backwards compat for integration tests
	CCIPSendRequestEventSigV1_2_0 common.Hash
)

const (
	CCIPSendRequestSeqNumIndexV1_2_0 = 4
	CCIPSendRequestedEventNameV1_2_0 = "CCIPSendRequested"
	EVM2EVMOffRampEventNameV1_2_0    = "EVM2EVMMessage"
	MetaDataHashPrefixV1_2_0         = "EVM2EVMMessageHashV2"
)

func init() {
	onRampABI, err := abi.JSON(strings.NewReader(evm_2_evm_onramp.EVM2EVMOnRampABI))
	if err != nil {
		panic(err)
	}
	CCIPSendRequestEventSigV1_2_0 = abihelpers.MustGetEventID(CCIPSendRequestedEventNameV1_2_0, onRampABI)
}

// Backwards compat for integration tests
func DecodeOffRampMessageV1_2_0(b []byte) (*evm_2_evm_offramp.InternalEVM2EVMMessage, error) {
	offRampABI, err := abi.JSON(strings.NewReader(evm_2_evm_offramp.EVM2EVMOffRampABI))
	if err != nil {
		panic(err)
	}
	event, ok := offRampABI.Events[EVM2EVMOffRampEventNameV1_2_0]
	if !ok {
		panic("no such event")
	}
	unpacked, err := event.Inputs.Unpack(b)
	if err != nil {
		return nil, err
	}
	if len(unpacked) == 0 {
		return nil, fmt.Errorf("no message found when unpacking")
	}

	// Note must use unnamed type here
	receivedCp, ok := unpacked[0].(struct {
		SourceChainSelector uint64         `json:"sourceChainSelector"`
		Sender              common.Address `json:"sender"`
		Receiver            common.Address `json:"receiver"`
		SequenceNumber      uint64         `json:"sequenceNumber"`
		GasLimit            *big.Int       `json:"gasLimit"`
		Strict              bool           `json:"strict"`
		Nonce               uint64         `json:"nonce"`
		FeeToken            common.Address `json:"feeToken"`
		FeeTokenAmount      *big.Int       `json:"feeTokenAmount"`
		Data                []uint8        `json:"data"`
		TokenAmounts        []struct {
			Token  common.Address `json:"token"`
			Amount *big.Int       `json:"amount"`
		} `json:"tokenAmounts"`
		SourceTokenData [][]byte `json:"sourceTokenData"`
		MessageId       [32]byte `json:"messageId"`
	})
	if !ok {
		return nil, fmt.Errorf("invalid format have %T want %T", unpacked[0], receivedCp)
	}
	var tokensAndAmounts []evm_2_evm_offramp.ClientEVMTokenAmount
	for _, tokenAndAmount := range receivedCp.TokenAmounts {
		tokensAndAmounts = append(tokensAndAmounts, evm_2_evm_offramp.ClientEVMTokenAmount{
			Token:  tokenAndAmount.Token,
			Amount: tokenAndAmount.Amount,
		})
	}

	return &evm_2_evm_offramp.InternalEVM2EVMMessage{
		SourceChainSelector: receivedCp.SourceChainSelector,
		Sender:              receivedCp.Sender,
		Receiver:            receivedCp.Receiver,
		SequenceNumber:      receivedCp.SequenceNumber,
		GasLimit:            receivedCp.GasLimit,
		Strict:              receivedCp.Strict,
		Nonce:               receivedCp.Nonce,
		FeeToken:            receivedCp.FeeToken,
		FeeTokenAmount:      receivedCp.FeeTokenAmount,
		Data:                receivedCp.Data,
		TokenAmounts:        tokensAndAmounts,
		SourceTokenData:     receivedCp.SourceTokenData,
		MessageId:           receivedCp.MessageId,
	}, nil
}

type LeafHasherV1_2_0 struct {
	metaDataHash [32]byte
	ctx          hashlib.Ctx[[32]byte]
	onRamp       *evm_2_evm_onramp.EVM2EVMOnRamp
}

func NewLeafHasherV1_2_0(sourceChainSelector uint64, destChainSelector uint64, onRampId common.Address, ctx hashlib.Ctx[[32]byte], onRamp *evm_2_evm_onramp.EVM2EVMOnRamp) *LeafHasherV1_2_0 {
	return &LeafHasherV1_2_0{
		metaDataHash: getMetaDataHash(ctx, ctx.Hash([]byte(MetaDataHashPrefixV1_2_0)), sourceChainSelector, onRampId, destChainSelector),
		ctx:          ctx,
		onRamp:       onRamp,
	}
}

func (t *LeafHasherV1_2_0) HashLeaf(log types.Log) ([32]byte, error) {
	msg, err := t.onRamp.ParseCCIPSendRequested(log)
	if err != nil {
		return [32]byte{}, err
	}
	message := msg.Message
	encodedTokens, err := utils.ABIEncode(
		`[
{"components": [{"name":"token","type":"address"},{"name":"amount","type":"uint256"}], "type":"tuple[]"}]`, message.TokenAmounts)
	if err != nil {
		return [32]byte{}, err
	}

	bytesArray, err := abi.NewType("bytes[]", "bytes[]", nil)
	if err != nil {
		return [32]byte{}, err
	}

	encodedSourceTokenData, err := abi.Arguments{abi.Argument{Type: bytesArray}}.PackValues([]interface{}{message.SourceTokenData})
	if err != nil {
		return [32]byte{}, err
	}

	packedFixedSizeValues, err := utils.ABIEncode(
		`[
{"name": "sender", "type":"address"},
{"name": "receiver", "type":"address"},
{"name": "sequenceNumber", "type":"uint64"},
{"name": "gasLimit", "type":"uint256"},
{"name": "strict", "type":"bool"},
{"name": "nonce", "type":"uint64"},
{"name": "feeToken","type": "address"},
{"name": "feeTokenAmount","type": "uint256"}
]`,
		message.Sender,
		message.Receiver,
		message.SequenceNumber,
		message.GasLimit,
		message.Strict,
		message.Nonce,
		message.FeeToken,
		message.FeeTokenAmount,
	)
	if err != nil {
		return [32]byte{}, err
	}
	fixedSizeValuesHash := t.ctx.Hash(packedFixedSizeValues)

	packedValues, err := utils.ABIEncode(
		`[
{"name": "leafDomainSeparator","type":"bytes1"},
{"name": "metadataHash", "type":"bytes32"},
{"name": "fixedSizeValuesHash", "type":"bytes32"},
{"name": "dataHash", "type":"bytes32"},
{"name": "tokenAmountsHash", "type":"bytes32"},
{"name": "sourceTokenDataHash", "type":"bytes32"}
]`,
		leafDomainSeparator,
		t.metaDataHash,
		fixedSizeValuesHash,
		t.ctx.Hash(message.Data),
		t.ctx.Hash(encodedTokens),
		t.ctx.Hash(encodedSourceTokenData),
	)
	if err != nil {
		return [32]byte{}, err
	}
	return t.ctx.Hash(packedValues), nil
}

var _ OnRampReader = &OnRampV1_2_0{}

// Significant change in 1.2:
// - CCIPSendRequested event signature has changed
type OnRampV1_2_0 struct {
	onRamp                     *evm_2_evm_onramp.EVM2EVMOnRamp
	address                    common.Address
	lggr                       logger.Logger
	lp                         logpoller.LogPoller
	leafHasher                 LeafHasherInterface[[32]byte]
	client                     client.Client
	finalityTags               bool
	filterName                 string
	sendRequestedEventSig      common.Hash
	sendRequestedSeqNumberWord int
}

func (o *OnRampV1_2_0) Address() (common.Address, error) {
	return o.onRamp.Address(), nil
}

func (o *OnRampV1_2_0) GetDynamicConfig() (OnRampDynamicConfig, error) {
	if o.onRamp == nil {
		return OnRampDynamicConfig{}, fmt.Errorf("onramp not initialized")
	}
	config, err := o.onRamp.GetDynamicConfig(&bind.CallOpts{})
	if err != nil {
		return OnRampDynamicConfig{}, fmt.Errorf("get dynamic config: %w", err)
	}
	return OnRampDynamicConfig{
		Router:                            config.Router,
		MaxNumberOfTokensPerMsg:           config.MaxNumberOfTokensPerMsg,
		DestGasOverhead:                   config.DestGasOverhead,
		DestGasPerPayloadByte:             config.DestGasPerPayloadByte,
		DestDataAvailabilityOverheadGas:   config.DestDataAvailabilityOverheadGas,
		DestGasPerDataAvailabilityByte:    config.DestGasPerDataAvailabilityByte,
		DestDataAvailabilityMultiplierBps: config.DestDataAvailabilityMultiplierBps,
		PriceRegistry:                     config.PriceRegistry,
		MaxDataBytes:                      config.MaxDataBytes,
		MaxPerMsgGasLimit:                 config.MaxPerMsgGasLimit,
	}, nil
}

func (o *OnRampV1_2_0) logToMessage(log types.Log) (*internal.EVM2EVMMessage, error) {
	msg, err := o.onRamp.ParseCCIPSendRequested(log)
	if err != nil {
		return nil, err
	}
	h, err := o.leafHasher.HashLeaf(log)
	if err != nil {
		return nil, err
	}
	tokensAndAmounts := make([]internal.TokenAmount, len(msg.Message.TokenAmounts))
	for i, tokenAndAmount := range msg.Message.TokenAmounts {
		tokensAndAmounts[i] = internal.TokenAmount{
			Token:  tokenAndAmount.Token,
			Amount: tokenAndAmount.Amount,
		}
	}

	return &internal.EVM2EVMMessage{
		SequenceNumber:      msg.Message.SequenceNumber,
		GasLimit:            msg.Message.GasLimit,
		Nonce:               msg.Message.Nonce,
		MessageId:           msg.Message.MessageId,
		SourceChainSelector: msg.Message.SourceChainSelector,
		Sender:              msg.Message.Sender,
		Receiver:            msg.Message.Receiver,
		Strict:              msg.Message.Strict,
		FeeToken:            msg.Message.FeeToken,
		FeeTokenAmount:      msg.Message.FeeTokenAmount,
		Data:                msg.Message.Data,
		TokenAmounts:        tokensAndAmounts,
		SourceTokenData:     msg.Message.SourceTokenData, // Breaking change 1.2
		Hash:                h,
	}, nil
}

func (o *OnRampV1_2_0) GetSendRequestsGteSeqNum(ctx context.Context, seqNum uint64, confs int) ([]Event[internal.EVM2EVMMessage], error) {
	if !o.finalityTags {
		logs, err2 := o.lp.LogsDataWordGreaterThan(
			o.sendRequestedEventSig,
			o.address,
			o.sendRequestedSeqNumberWord,
			abihelpers.EvmWord(seqNum),
			confs,
			pg.WithParentCtx(ctx),
		)
		if err2 != nil {
			return nil, fmt.Errorf("logs data word greater than: %w", err2)
		}
		return parseLogs[internal.EVM2EVMMessage](logs, o.lggr, o.logToMessage)
	}
	latestFinalizedHash, err := latestFinalizedBlockHash(ctx, o.client)
	if err != nil {
		return nil, err
	}
	logs, err := o.lp.LogsUntilBlockHashDataWordGreaterThan(
		o.sendRequestedEventSig,
		o.address,
		o.sendRequestedSeqNumberWord,
		abihelpers.EvmWord(seqNum),
		latestFinalizedHash,
		pg.WithParentCtx(ctx),
	)
	if err != nil {
		return nil, fmt.Errorf("logs until block hash data word greater than: %w", err)
	}
	return parseLogs[internal.EVM2EVMMessage](logs, o.lggr, o.logToMessage)
}

func (o *OnRampV1_2_0) GetSendRequestsBetweenSeqNums(ctx context.Context, seqNumMin, seqNumMax uint64, confs int) ([]Event[internal.EVM2EVMMessage], error) {
	logs, err := o.lp.LogsDataWordRange(
		o.sendRequestedEventSig,
		o.address,
		o.sendRequestedSeqNumberWord,
		logpoller.EvmWord(seqNumMin),
		logpoller.EvmWord(seqNumMax),
		confs,
		pg.WithParentCtx(ctx))
	if err != nil {
		return nil, err
	}
	return parseLogs[internal.EVM2EVMMessage](logs, o.lggr, o.logToMessage)
}

func (o *OnRampV1_2_0) RouterAddress() (common.Address, error) {
	config, err := o.onRamp.GetDynamicConfig(nil)
	if err != nil {
		return common.Address{}, err
	}
	return config.Router, nil
}

func (o *OnRampV1_2_0) Close(qopts ...pg.QOpt) error {
	return o.lp.UnregisterFilter(o.filterName, qopts...)
}

func NewOnRampV1_2_0(
	lggr logger.Logger,
	sourceSelector,
	destSelector uint64,
	onRampAddress common.Address,
	sourceLP logpoller.LogPoller,
	source client.Client,
	finalityTags bool,
) (*OnRampV1_2_0, error) {
	onRamp, err := evm_2_evm_onramp.NewEVM2EVMOnRamp(onRampAddress, source)
	if err != nil {
		return nil, err
	}
	// Subscribe to the relevant logs
	// Note we can keep the same prefix across 1.0/1.1 and 1.2 because the onramp addresses will be different
	name := logpoller.FilterName(COMMIT_CCIP_SENDS, onRampAddress)
	if err = sourceLP.RegisterFilter(logpoller.Filter{
		Name:      name,
		EventSigs: []common.Hash{CCIPSendRequestEventSigV1_2_0},
		Addresses: []common.Address{onRampAddress},
	}); err != nil {
		return nil, err
	}
	return &OnRampV1_2_0{
		finalityTags:               finalityTags,
		lggr:                       lggr,
		client:                     source,
		lp:                         sourceLP,
		leafHasher:                 NewLeafHasherV1_2_0(sourceSelector, destSelector, onRampAddress, hashlib.NewKeccakCtx(), onRamp),
		onRamp:                     onRamp,
		filterName:                 name,
		address:                    onRampAddress,
		sendRequestedSeqNumberWord: CCIPSendRequestSeqNumIndexV1_2_0,
		sendRequestedEventSig:      CCIPSendRequestEventSigV1_2_0,
	}, nil
}
