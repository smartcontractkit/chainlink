package ccipdata

import (
	"context"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_onramp"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/hashlib"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/logpollerutil"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

var (
	// Backwards compat for integration tests
	CCIPSendRequestEventSigV1_3_0 common.Hash
)

const (
	CCIPSendRequestSeqNumIndexV1_3_0 = 4
	CCIPSendRequestedEventNameV1_3_0 = "CCIPSendRequested"
	EVM2EVMOffRampEventNameV1_3_0    = "EVM2EVMMessage"
	MetaDataHashPrefixV1_3_0         = "EVM2EVMMessageHashV2"
)

func init() {
	onRampABI, err := abi.JSON(strings.NewReader(evm_2_evm_onramp.EVM2EVMOnRampABI))
	if err != nil {
		panic(err)
	}
	CCIPSendRequestEventSigV1_3_0 = abihelpers.MustGetEventID(CCIPSendRequestedEventNameV1_3_0, onRampABI)
}

type LeafHasherV1_3_0 struct {
	metaDataHash [32]byte
	ctx          hashlib.Ctx[[32]byte]
	onRamp       *evm_2_evm_onramp.EVM2EVMOnRamp
}

func NewLeafHasherV1_3_0(sourceChainSelector uint64, destChainSelector uint64, onRampId common.Address, ctx hashlib.Ctx[[32]byte], onRamp *evm_2_evm_onramp.EVM2EVMOnRamp) *LeafHasherV1_3_0 {
	return &LeafHasherV1_3_0{
		metaDataHash: getMetaDataHash(ctx, ctx.Hash([]byte(MetaDataHashPrefixV1_3_0)), sourceChainSelector, onRampId, destChainSelector),
		ctx:          ctx,
		onRamp:       onRamp,
	}
}

func (t *LeafHasherV1_3_0) HashLeaf(log types.Log) ([32]byte, error) {
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

var _ OnRampReader = &OnRampV1_3_0{}

// Significant change in 1.2:
// - CCIPSendRequested event signature has changed
type OnRampV1_3_0 struct {
	onRamp                     *evm_2_evm_onramp.EVM2EVMOnRamp
	address                    common.Address
	lggr                       logger.Logger
	lp                         logpoller.LogPoller
	leafHasher                 LeafHasherInterface[[32]byte]
	client                     client.Client
	sendRequestedEventSig      common.Hash
	sendRequestedSeqNumberWord int
	filters                    []logpoller.Filter
}

func (o *OnRampV1_3_0) Address() (common.Address, error) {
	return o.onRamp.Address(), nil
}

func (o *OnRampV1_3_0) GetDynamicConfig() (OnRampDynamicConfig, error) {
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

func (o *OnRampV1_3_0) logToMessage(log types.Log) (*internal.EVM2EVMMessage, error) {
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

func (o *OnRampV1_3_0) GetSendRequestsBetweenSeqNums(ctx context.Context, seqNumMin, seqNumMax uint64, finalized bool) ([]Event[internal.EVM2EVMMessage], error) {
	logs, err := o.lp.LogsDataWordRange(
		o.sendRequestedEventSig,
		o.address,
		o.sendRequestedSeqNumberWord,
		logpoller.EvmWord(seqNumMin),
		logpoller.EvmWord(seqNumMax),
		logsConfirmations(finalized),
		pg.WithParentCtx(ctx))
	if err != nil {
		return nil, err
	}
	return parseLogs[internal.EVM2EVMMessage](logs, o.lggr, o.logToMessage)
}

func (o *OnRampV1_3_0) RouterAddress() (common.Address, error) {
	config, err := o.onRamp.GetDynamicConfig(nil)
	if err != nil {
		return common.Address{}, err
	}
	return config.Router, nil
}

func (o *OnRampV1_3_0) Close(qopts ...pg.QOpt) error {
	return logpollerutil.UnregisterLpFilters(o.lp, o.filters, qopts...)
}

func (o *OnRampV1_3_0) RegisterFilters(qopts ...pg.QOpt) error {
	return logpollerutil.RegisterLpFilters(o.lp, o.filters, qopts...)
}

func NewOnRampV1_3_0(lggr logger.Logger, sourceSelector, destSelector uint64, onRampAddress common.Address, sourceLP logpoller.LogPoller, source client.Client) (*OnRampV1_3_0, error) {
	onRamp, err := evm_2_evm_onramp.NewEVM2EVMOnRamp(onRampAddress, source)
	if err != nil {
		return nil, err
	}
	// Subscribe to the relevant logs
	// Note we can keep the same prefix across 1.0/1.1 and 1.2 because the onramp addresses will be different
	filters := []logpoller.Filter{
		{
			Name:      logpoller.FilterName(COMMIT_CCIP_SENDS, onRampAddress),
			EventSigs: []common.Hash{CCIPSendRequestEventSigV1_3_0},
			Addresses: []common.Address{onRampAddress},
		},
	}
	return &OnRampV1_3_0{
		lggr:                       lggr,
		client:                     source,
		lp:                         sourceLP,
		leafHasher:                 NewLeafHasherV1_3_0(sourceSelector, destSelector, onRampAddress, hashlib.NewKeccakCtx(), onRamp),
		onRamp:                     onRamp,
		filters:                    filters,
		address:                    onRampAddress,
		sendRequestedSeqNumberWord: CCIPSendRequestSeqNumIndexV1_3_0,
		sendRequestedEventSig:      CCIPSendRequestEventSigV1_3_0,
	}, nil
}
