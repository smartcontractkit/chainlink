package v1_5_0

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink-common/pkg/hashutil"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"

	"github.com/smartcontractkit/chainlink-evm/pkg/client"
	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/evm_2_evm_onramp"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/rmn_contract"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/cache"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcalc"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/logpollerutil"
)

var (
	// Backwards compat for integration tests
	CCIPSendRequestEventSig common.Hash
	ConfigSetEventSig       common.Hash
)

const (
	CCIPSendRequestSeqNumIndex = 4
	CCIPSendRequestedEventName = "CCIPSendRequested"
	ConfigSetEventName         = "ConfigSet"
)

func init() {
	onRampABI, err := abi.JSON(strings.NewReader(evm_2_evm_onramp.EVM2EVMOnRampABI))
	if err != nil {
		panic(err)
	}
	CCIPSendRequestEventSig = abihelpers.MustGetEventID(CCIPSendRequestedEventName, onRampABI)
	ConfigSetEventSig = abihelpers.MustGetEventID(ConfigSetEventName, onRampABI)
}

var _ ccipdata.OnRampReader = &OnRamp{}

type OnRamp struct {
	onRamp                     *evm_2_evm_onramp.EVM2EVMOnRamp
	address                    common.Address
	destChainSelectorBytes     [16]byte
	lggr                       logger.Logger
	lp                         logpoller.LogPoller
	leafHasher                 ccipdata.LeafHasherInterface[[32]byte]
	client                     client.Client
	sendRequestedEventSig      common.Hash
	sendRequestedSeqNumberWord int
	filters                    []logpoller.Filter
	cachedOnRampDynamicConfig  cache.AutoSync[cciptypes.OnRampDynamicConfig]
	// Static config can be cached, because it's never expected to change.
	// The only way to change that is through the contract's constructor (redeployment)
	cachedStaticConfig cache.OnceCtxFunction[evm_2_evm_onramp.EVM2EVMOnRampStaticConfig]
	cachedRmnContract  cache.OnceCtxFunction[*rmn_contract.RMNContract]
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
			Retention: ccipdata.CommitExecLogsRetention,
		},
		{
			Name:      logpoller.FilterName(ccipdata.CONFIG_CHANGED, onRampAddress),
			EventSigs: []common.Hash{ConfigSetEventSig},
			Addresses: []common.Address{onRampAddress},
			Retention: ccipdata.CacheEvictionLogsRetention,
		},
	}
	cachedStaticConfig := cache.OnceCtxFunction[evm_2_evm_onramp.EVM2EVMOnRampStaticConfig](func(ctx context.Context) (evm_2_evm_onramp.EVM2EVMOnRampStaticConfig, error) {
		return onRamp.GetStaticConfig(&bind.CallOpts{Context: ctx})
	})
	cachedRmnContract := cache.OnceCtxFunction[*rmn_contract.RMNContract](func(ctx context.Context) (*rmn_contract.RMNContract, error) {
		staticConfig, err := cachedStaticConfig(ctx)
		if err != nil {
			return nil, err
		}

		return rmn_contract.NewRMNContract(staticConfig.RmnProxy, source)
	})

	return &OnRamp{
		lggr:                       lggr,
		client:                     source,
		destChainSelectorBytes:     ccipcommon.SelectorToBytes(destSelector),
		lp:                         sourceLP,
		leafHasher:                 NewLeafHasher(sourceSelector, destSelector, onRampAddress, hashutil.NewKeccak(), onRamp),
		onRamp:                     onRamp,
		filters:                    filters,
		address:                    onRampAddress,
		sendRequestedSeqNumberWord: CCIPSendRequestSeqNumIndex,
		sendRequestedEventSig:      CCIPSendRequestEventSig,
		cachedOnRampDynamicConfig: cache.NewLogpollerEventsBased[cciptypes.OnRampDynamicConfig](
			sourceLP,
			[]common.Hash{ConfigSetEventSig},
			onRampAddress,
		),
		cachedStaticConfig: cache.CallOnceOnNoError(cachedStaticConfig),
		cachedRmnContract:  cache.CallOnceOnNoError(cachedRmnContract),
	}, nil
}

func (o *OnRamp) Address(context.Context) (cciptypes.Address, error) {
	return ccipcalc.EvmAddrToGeneric(o.onRamp.Address()), nil
}

func (o *OnRamp) GetDynamicConfig(ctx context.Context) (cciptypes.OnRampDynamicConfig, error) {
	return o.cachedOnRampDynamicConfig.Get(ctx, func(ctx context.Context) (cciptypes.OnRampDynamicConfig, error) {
		if o.onRamp == nil {
			return cciptypes.OnRampDynamicConfig{}, errors.New("onramp not initialized")
		}
		config, err := o.onRamp.GetDynamicConfig(&bind.CallOpts{})
		if err != nil {
			return cciptypes.OnRampDynamicConfig{}, fmt.Errorf("get dynamic config v1.5: %w", err)
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
	})
}

func (o *OnRamp) SourcePriceRegistryAddress(ctx context.Context) (cciptypes.Address, error) {
	c, err := o.GetDynamicConfig(ctx)
	if err != nil {
		return "", err
	}
	return c.PriceRegistry, nil
}

func (o *OnRamp) GetSendRequestsBetweenSeqNums(ctx context.Context, seqNumMin, seqNumMax uint64, finalized bool) ([]cciptypes.EVM2EVMMessageWithTxMeta, error) {
	logs, err := o.lp.LogsDataWordRange(
		ctx,
		o.sendRequestedEventSig,
		o.address,
		o.sendRequestedSeqNumberWord,
		logpoller.EvmWord(seqNumMin),
		logpoller.EvmWord(seqNumMax),
		ccipdata.LogsConfirmations(finalized),
	)
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

func (o *OnRamp) RouterAddress(context.Context) (cciptypes.Address, error) {
	config, err := o.onRamp.GetDynamicConfig(nil)
	if err != nil {
		return "", err
	}
	return ccipcalc.EvmAddrToGeneric(config.Router), nil
}

func (o *OnRamp) IsSourceChainHealthy(context.Context) (bool, error) {
	if err := o.lp.Healthy(); err != nil {
		return false, nil
	}
	return true, nil
}

func (o *OnRamp) IsSourceCursed(ctx context.Context) (bool, error) {
	arm, err := o.cachedRmnContract(ctx)
	if err != nil {
		return false, fmt.Errorf("initializing RMN contract through the RmnProxy: %w", err)
	}

	cursed, err := arm.IsCursed(&bind.CallOpts{Context: ctx}, o.destChainSelectorBytes)
	if err != nil {
		return false, fmt.Errorf("checking if source is cursed by RMN: %w", err)
	}
	return cursed, nil
}

func (o *OnRamp) Close() error {
	return logpollerutil.UnregisterLpFilters(context.Background(), o.lp, o.filters)
}

func (o *OnRamp) RegisterFilters(ctx context.Context) error {
	return logpollerutil.RegisterLpFilters(ctx, o.lp, o.filters)
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
