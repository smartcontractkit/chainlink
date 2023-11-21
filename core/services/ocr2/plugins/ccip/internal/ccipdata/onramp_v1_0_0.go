package ccipdata

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_onramp_1_0_0"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/hashlib"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

const (
	CCIPSendRequestedEventNameV1_0_0 = "CCIPSendRequested"
	MetaDataHashPrefixV1_0_0         = "EVM2EVMMessageEvent"
)

var leafDomainSeparator = [1]byte{0x00}

type LeafHasherV1_0_0 struct {
	metaDataHash [32]byte
	ctx          hashlib.Ctx[[32]byte]
	onRamp       *evm_2_evm_onramp_1_0_0.EVM2EVMOnRamp
}

func getMetaDataHash[H hashlib.Hash](ctx hashlib.Ctx[H], prefix [32]byte, sourceChainSelector uint64, onRampId common.Address, destChainSelector uint64) H {
	paddedOnRamp := onRampId.Hash()
	return ctx.Hash(utils.ConcatBytes(prefix[:], math.U256Bytes(big.NewInt(0).SetUint64(sourceChainSelector)), math.U256Bytes(big.NewInt(0).SetUint64(destChainSelector)), paddedOnRamp[:]))
}

func NewLeafHasherV1_0_0(sourceChainSelector uint64, destChainSelector uint64, onRampId common.Address, ctx hashlib.Ctx[[32]byte], onRamp *evm_2_evm_onramp_1_0_0.EVM2EVMOnRamp) *LeafHasherV1_0_0 {
	return &LeafHasherV1_0_0{
		metaDataHash: getMetaDataHash(ctx, ctx.Hash([]byte(MetaDataHashPrefixV1_0_0)), sourceChainSelector, onRampId, destChainSelector),
		ctx:          ctx,
		onRamp:       onRamp,
	}
}

func (t *LeafHasherV1_0_0) HashLeaf(log types.Log) ([32]byte, error) {
	message, err := t.onRamp.ParseCCIPSendRequested(log)
	if err != nil {
		return [32]byte{}, err
	}
	encodedTokens, err := utils.ABIEncode(
		`[
{"components": [{"name":"token","type":"address"},{"name":"amount","type":"uint256"}], "type":"tuple[]"}]`, message.Message.TokenAmounts)
	if err != nil {
		return [32]byte{}, err
	}

	packedValues, err := utils.ABIEncode(
		`[
{"name": "leafDomainSeparator","type":"bytes1"},
{"name": "metadataHash", "type":"bytes32"},
{"name": "sequenceNumber", "type":"uint64"},
{"name": "nonce", "type":"uint64"},
{"name": "sender", "type":"address"},
{"name": "receiver", "type":"address"},
{"name": "dataHash", "type":"bytes32"},
{"name": "tokenAmountsHash", "type":"bytes32"},
{"name": "gasLimit", "type":"uint256"},
{"name": "strict", "type":"bool"},
{"name": "feeToken","type": "address"},
{"name": "feeTokenAmount","type": "uint256"}
]`,
		leafDomainSeparator,
		t.metaDataHash,
		message.Message.SequenceNumber,
		message.Message.Nonce,
		message.Message.Sender,
		message.Message.Receiver,
		t.ctx.Hash(message.Message.Data),
		t.ctx.Hash(encodedTokens),
		message.Message.GasLimit,
		message.Message.Strict,
		message.Message.FeeToken,
		message.Message.FeeTokenAmount,
	)
	if err != nil {
		return [32]byte{}, err
	}
	return t.ctx.Hash(packedValues), nil
}

var _ OnRampReader = &OnRampV1_0_0{}

type OnRampV1_0_0 struct {
	address                    common.Address
	onRamp                     *evm_2_evm_onramp_1_0_0.EVM2EVMOnRamp
	lp                         logpoller.LogPoller
	lggr                       logger.Logger
	client                     client.Client
	leafHasher                 LeafHasherInterface[[32]byte]
	filterName                 string
	sendRequestedEventSig      common.Hash
	sendRequestedSeqNumberWord int
}

func (o *OnRampV1_0_0) Address() (common.Address, error) {
	return o.onRamp.Address(), nil
}

func (o *OnRampV1_0_0) GetDynamicConfig() (OnRampDynamicConfig, error) {
	if o.onRamp == nil {
		return OnRampDynamicConfig{}, fmt.Errorf("onramp not initialized")
	}
	legacyDynamicConfig, err := o.onRamp.GetDynamicConfig(nil)
	if err != nil {
		return OnRampDynamicConfig{}, err
	}
	return OnRampDynamicConfig{
		Router:                            legacyDynamicConfig.Router,
		MaxNumberOfTokensPerMsg:           legacyDynamicConfig.MaxTokensLength,
		DestGasOverhead:                   0,
		DestGasPerPayloadByte:             0,
		DestDataAvailabilityOverheadGas:   0,
		DestGasPerDataAvailabilityByte:    0,
		DestDataAvailabilityMultiplierBps: 0,
		PriceRegistry:                     legacyDynamicConfig.PriceRegistry,
		MaxDataBytes:                      legacyDynamicConfig.MaxDataSize,
		MaxPerMsgGasLimit:                 uint32(legacyDynamicConfig.MaxGasLimit),
	}, nil
}

func (o *OnRampV1_0_0) GetLastUSDCMessagePriorToLogIndexInTx(ctx context.Context, logIndex int64, txHash common.Hash) ([]byte, error) {
	return nil, errors.New("USDC not supported in < 1.2.0")
}

func NewOnRampV1_0_0(lggr logger.Logger, sourceSelector, destSelector uint64, onRampAddress common.Address, sourceLP logpoller.LogPoller, source client.Client) (*OnRampV1_0_0, error) {
	onRamp, err := evm_2_evm_onramp_1_0_0.NewEVM2EVMOnRamp(onRampAddress, source)
	if err != nil {
		return nil, err
	}
	onRampABI := abihelpers.MustParseABI(evm_2_evm_onramp_1_0_0.EVM2EVMOnRampABI)
	// Subscribe to the relevant logs
	name := logpoller.FilterName(COMMIT_CCIP_SENDS, onRampAddress)
	eventSig := abihelpers.MustGetEventID(CCIPSendRequestedEventNameV1_0_0, onRampABI)
	err = sourceLP.RegisterFilter(logpoller.Filter{
		Name:      name,
		EventSigs: []common.Hash{eventSig},
		Addresses: []common.Address{onRampAddress},
	})
	if err != nil {
		return nil, err
	}
	return &OnRampV1_0_0{
		lggr:       lggr,
		address:    onRampAddress,
		onRamp:     onRamp,
		client:     source,
		lp:         sourceLP,
		leafHasher: NewLeafHasherV1_0_0(sourceSelector, destSelector, onRampAddress, hashlib.NewKeccakCtx(), onRamp),
		filterName: name,
		// offset || sourceChainID || seqNum || ...
		sendRequestedSeqNumberWord: 2,
		sendRequestedEventSig:      eventSig,
	}, nil
}

func (o *OnRampV1_0_0) logToMessage(log types.Log) (*internal.EVM2EVMMessage, error) {
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
		SourceTokenData:     make([][]byte, len(msg.Message.TokenAmounts)), // Always empty in 1.0
		Hash:                h,
	}, nil
}

func (o *OnRampV1_0_0) GetSendRequestsGteSeqNum(ctx context.Context, seqNum uint64) ([]Event[internal.EVM2EVMMessage], error) {
	logs, err := o.lp.LogsDataWordGreaterThan(
		o.sendRequestedEventSig,
		o.address,
		o.sendRequestedSeqNumberWord,
		abihelpers.EvmWord(seqNum),
		logpoller.Finalized,
		pg.WithParentCtx(ctx),
	)
	if err != nil {
		return nil, fmt.Errorf("logs data word greater than: %w", err)
	}
	return parseLogs[internal.EVM2EVMMessage](logs, o.lggr, o.logToMessage)
}

func (o *OnRampV1_0_0) RouterAddress() (common.Address, error) {
	config, err := o.onRamp.GetDynamicConfig(nil)
	if err != nil {
		return common.Address{}, err
	}
	return config.Router, nil
}

func (o *OnRampV1_0_0) GetSendRequestsBetweenSeqNums(ctx context.Context, seqNumMin, seqNumMax uint64) ([]Event[internal.EVM2EVMMessage], error) {
	logs, err := o.lp.LogsDataWordRange(
		o.sendRequestedEventSig,
		o.address,
		o.sendRequestedSeqNumberWord,
		logpoller.EvmWord(seqNumMin),
		logpoller.EvmWord(seqNumMax),
		logpoller.Finalized,
		pg.WithParentCtx(ctx))
	if err != nil {
		return nil, err
	}
	return parseLogs[internal.EVM2EVMMessage](logs, o.lggr, o.logToMessage)
}

func (o *OnRampV1_0_0) Close(qopts ...pg.QOpt) error {
	return o.lp.UnregisterFilter(o.filterName, qopts...)
}
