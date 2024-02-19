package v1_5_0

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
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_onramp"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/cciptypes"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcalc"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/v1_0_0"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/hashlib"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/logpollerutil"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

var (
	// Backwards compat for integration tests
	CCIPSendRequestEventSig common.Hash
)

const (
	CCIPSendRequestSeqNumIndex = 4
	CCIPSendRequestedEventName = "CCIPSendRequested"
	EVM2EVMOffRampEventName    = "EVM2EVMMessage"
	MetaDataHashPrefix         = "EVM2EVMMessageHashV2"
)

func init() {
	onRampABI, err := abi.JSON(strings.NewReader(evm_2_evm_onramp.EVM2EVMOnRampABI))
	if err != nil {
		panic(err)
	}
	CCIPSendRequestEventSig = abihelpers.MustGetEventID(CCIPSendRequestedEventName, onRampABI)
}

type LeafHasher struct {
	metaDataHash [32]byte
	ctx          hashlib.Ctx[[32]byte]
	onRamp       *evm_2_evm_onramp.EVM2EVMOnRamp
}

func NewLeafHasher(sourceChainSelector uint64, destChainSelector uint64, onRampId common.Address, ctx hashlib.Ctx[[32]byte], onRamp *evm_2_evm_onramp.EVM2EVMOnRamp) *LeafHasher {
	return &LeafHasher{
		metaDataHash: v1_0_0.GetMetaDataHash(ctx, ctx.Hash([]byte(MetaDataHashPrefix)), sourceChainSelector, onRampId, destChainSelector),
		ctx:          ctx,
		onRamp:       onRamp,
	}
}

func (t *LeafHasher) HashLeaf(log types.Log) ([32]byte, error) {
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
		v1_0_0.LeafDomainSeparator,
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

var _ ccipdata.OnRampReader = &OnRamp{}

// Significant change in 1.2:
// - CCIPSendRequested event signature has changed
type OnRamp struct {
	onRamp                     *evm_2_evm_onramp.EVM2EVMOnRamp
	address                    common.Address
	lggr                       logger.Logger
	lp                         logpoller.LogPoller
	leafHasher                 ccipdata.LeafHasherInterface[[32]byte]
	client                     client.Client
	sendRequestedEventSig      common.Hash
	sendRequestedSeqNumberWord int
	filters                    []logpoller.Filter
}

func (o *OnRamp) Address() (cciptypes.Address, error) {
	return ccipcalc.EvmAddrToGeneric(o.onRamp.Address()), nil
}

func (o *OnRamp) GetDynamicConfig() (cciptypes.OnRampDynamicConfig, error) {
	if o.onRamp == nil {
		return cciptypes.OnRampDynamicConfig{}, fmt.Errorf("onramp not initialized")
	}
	config, err := o.onRamp.GetDynamicConfig(&bind.CallOpts{})
	if err != nil {
		return cciptypes.OnRampDynamicConfig{}, fmt.Errorf("get dynamic config: %w", err)
	}
	return cciptypes.OnRampDynamicConfig{
		Router:                            ccipcalc.EvmAddrToGeneric(config.Router),
		MaxNumberOfTokensPerMsg:           config.MaxNumberOfTokensPerMsg,
		DestGasOverhead:                   config.DestGasOverhead,
		DestGasPerPayloadByte:             config.DestGasPerPayloadByte,
		DestDataAvailabilityOverheadGas:   config.DestDataAvailabilityOverheadGas,
		DestGasPerDataAvailabilityByte:    config.DestGasPerDataAvailabilityByte,
		DestDataAvailabilityMultiplierBps: config.DestDataAvailabilityMultiplierBps,
		PriceRegistry:                     ccipcalc.EvmAddrToGeneric(config.PriceRegistry),
		MaxDataBytes:                      config.MaxDataBytes,
		MaxPerMsgGasLimit:                 config.MaxPerMsgGasLimit,
	}, nil
}

func (o *OnRamp) logToMessage(log types.Log) (*cciptypes.EVM2EVMMessage, error) {
	msg, err := o.onRamp.ParseCCIPSendRequested(log)
	if err != nil {
		return nil, err
	}
	h, err := o.leafHasher.HashLeaf(log)
	if err != nil {
		return nil, err
	}
	tokensAndAmounts := make([]cciptypes.TokenAmount, len(msg.Message.TokenAmounts))
	for i, tokenAndAmount := range msg.Message.TokenAmounts {
		tokensAndAmounts[i] = cciptypes.TokenAmount{
			Token:  ccipcalc.EvmAddrToGeneric(tokenAndAmount.Token),
			Amount: tokenAndAmount.Amount,
		}
	}

	return &cciptypes.EVM2EVMMessage{
		SequenceNumber:      msg.Message.SequenceNumber,
		GasLimit:            msg.Message.GasLimit,
		Nonce:               msg.Message.Nonce,
		MessageID:           msg.Message.MessageId,
		SourceChainSelector: msg.Message.SourceChainSelector,
		Sender:              ccipcalc.EvmAddrToGeneric(msg.Message.Sender),
		Receiver:            ccipcalc.EvmAddrToGeneric(msg.Message.Receiver),
		Strict:              msg.Message.Strict,
		FeeToken:            ccipcalc.EvmAddrToGeneric(msg.Message.FeeToken),
		FeeTokenAmount:      msg.Message.FeeTokenAmount,
		Data:                msg.Message.Data,
		TokenAmounts:        tokensAndAmounts,
		SourceTokenData:     msg.Message.SourceTokenData, // Breaking change 1.2
		Hash:                h,
	}, nil
}

func (o *OnRamp) GetSendRequestsBetweenSeqNums(ctx context.Context, seqNumMin, seqNumMax uint64, finalized bool) ([]cciptypes.EVM2EVMMessageWithTxMeta, error) {
	logs, err := o.lp.LogsDataWordRange(
		o.sendRequestedEventSig,
		o.address,
		o.sendRequestedSeqNumberWord,
		logpoller.EvmWord(seqNumMin),
		logpoller.EvmWord(seqNumMax),
		ccipdata.LogsConfirmations(finalized),
		pg.WithParentCtx(ctx))
	if err != nil {
		return nil, err
	}

	parsedLogs, err := ccipdata.ParseLogs[cciptypes.EVM2EVMMessage](logs, o.lggr, o.logToMessage)
	if err != nil {
		return nil, err
	}

	res := make([]cciptypes.EVM2EVMMessageWithTxMeta, 0, len(logs))
	for _, log := range parsedLogs {
		res = append(res, cciptypes.EVM2EVMMessageWithTxMeta{
			TxMeta:         log.TxMeta,
			EVM2EVMMessage: log.Data,
		})
	}
	return res, nil
}

func (o *OnRamp) RouterAddress() (cciptypes.Address, error) {
	config, err := o.onRamp.GetDynamicConfig(nil)
	if err != nil {
		return "", err
	}
	return ccipcalc.EvmAddrToGeneric(config.Router), nil
}

func (o *OnRamp) Close(qopts ...pg.QOpt) error {
	return logpollerutil.UnregisterLpFilters(o.lp, o.filters, qopts...)
}

func (o *OnRamp) RegisterFilters(qopts ...pg.QOpt) error {
	return logpollerutil.RegisterLpFilters(o.lp, o.filters, qopts...)
}

func NewOnRamp(lggr logger.Logger, sourceSelector, destSelector uint64, onRampAddress common.Address, sourceLP logpoller.LogPoller, source client.Client) (*OnRamp, error) {
	onRamp, err := evm_2_evm_onramp.NewEVM2EVMOnRamp(onRampAddress, source)
	if err != nil {
		return nil, err
	}
	// Subscribe to the relevant logs
	// Note we can keep the same prefix across 1.0/1.1 and 1.2 because the onramp addresses will be different
	filters := []logpoller.Filter{
		{
			Name:      logpoller.FilterName(ccipdata.COMMIT_CCIP_SENDS, onRampAddress),
			EventSigs: []common.Hash{CCIPSendRequestEventSig},
			Addresses: []common.Address{onRampAddress},
		},
	}
	return &OnRamp{
		lggr:                       lggr,
		client:                     source,
		lp:                         sourceLP,
		leafHasher:                 NewLeafHasher(sourceSelector, destSelector, onRampAddress, hashlib.NewKeccakCtx(), onRamp),
		onRamp:                     onRamp,
		filters:                    filters,
		address:                    onRampAddress,
		sendRequestedSeqNumberWord: CCIPSendRequestSeqNumIndex,
		sendRequestedEventSig:      CCIPSendRequestEventSig,
	}, nil
}
