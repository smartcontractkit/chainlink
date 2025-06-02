package v1_2_0

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"

	"github.com/smartcontractkit/chainlink-evm/pkg/client"
	"github.com/smartcontractkit/chainlink-evm/pkg/gas"
	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
	evmtypes "github.com/smartcontractkit/chainlink-evm/pkg/types"

	evm_2_evm_offramp_1_2_0 "github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/evm_2_evm_offramp"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/cache"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcalc"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/logpollerutil"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/rpclib"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/prices"
)

const (
	ExecExecutionStateChanges = "Exec execution state changes"
	ExecTokenPoolAdded        = "Token pool added"
	ExecTokenPoolRemoved      = "Token pool removed"
)

var (
	abiOffRamp                                               = abihelpers.MustParseABI(evm_2_evm_offramp_1_2_0.EVM2EVMOffRampABI)
	_                                 ccipdata.OffRampReader = &OffRamp{}
	ExecutionStateChangedEvent                               = abihelpers.MustGetEventID("ExecutionStateChanged", abiOffRamp)
	PoolAddedEvent                                           = abihelpers.MustGetEventID("PoolAdded", abiOffRamp)
	PoolRemovedEvent                                         = abihelpers.MustGetEventID("PoolRemoved", abiOffRamp)
	ExecutionStateChangedSeqNrIndex                          = 1
	offrampPoolAddedPoolRemovedEvents                        = []common.Hash{PoolAddedEvent, PoolRemovedEvent}
)

type ExecOnchainConfig evm_2_evm_offramp_1_2_0.EVM2EVMOffRampDynamicConfig

func (d ExecOnchainConfig) AbiString() string {
	return `
	[
		{
			"components": [
				{"name": "permissionLessExecutionThresholdSeconds", "type": "uint32"},
				{"name": "router", "type": "address"},
				{"name": "priceRegistry", "type": "address"},
				{"name": "maxNumberOfTokensPerMsg", "type": "uint16"},
				{"name": "maxDataBytes", "type": "uint32"},
				{"name": "maxPoolReleaseOrMintGas", "type": "uint32"}
			],
			"type": "tuple"
		}
	]`
}

func (d ExecOnchainConfig) Validate() error {
	if d.PermissionLessExecutionThresholdSeconds == 0 {
		return errors.New("must set PermissionLessExecutionThresholdSeconds")
	}
	if d.Router == (common.Address{}) {
		return errors.New("must set Router address")
	}
	if d.PriceRegistry == (common.Address{}) {
		return errors.New("must set PriceRegistry address")
	}
	if d.MaxNumberOfTokensPerMsg == 0 {
		return errors.New("must set MaxNumberOfTokensPerMsg")
	}
	if d.MaxPoolReleaseOrMintGas == 0 {
		return errors.New("must set MaxPoolReleaseOrMintGas")
	}
	return nil
}

// JSONExecOffchainConfig is the configuration for nodes executing committed CCIP messages (v1.2).
// It comes from the OffchainConfig field of the corresponding OCR2 plugin configuration.
// NOTE: do not change the JSON format of this struct without consulting with the RDD people first.
type JSONExecOffchainConfig struct {
	// SourceFinalityDepth indicates how many confirmations a transaction should get on the source chain event before we consider it finalized.
	//
	// Deprecated: we now use the source chain finality instead.
	SourceFinalityDepth uint32
	// See [ccipdata.ExecOffchainConfig.DestOptimisticConfirmations]
	DestOptimisticConfirmations uint32
	// DestFinalityDepth indicates how many confirmations a transaction should get on the destination chain event before we consider it finalized.
	//
	// Deprecated: we now use the destination chain finality instead.
	DestFinalityDepth uint32
	// See [ccipdata.ExecOffchainConfig.BatchGasLimit]
	BatchGasLimit uint32
	// See [ccipdata.ExecOffchainConfig.RelativeBoostPerWaitHour]
	RelativeBoostPerWaitHour float64
	// See [ccipdata.ExecOffchainConfig.InflightCacheExpiry]
	InflightCacheExpiry config.Duration
	// See [ccipdata.ExecOffchainConfig.RootSnoozeTime]
	RootSnoozeTime config.Duration
	// See [ccipdata.ExecOffchainConfig.BatchingStrategyID]
	BatchingStrategyID uint32
	// See [ccipdata.ExecOffchainConfig.MessageVisibilityInterval]
	MessageVisibilityInterval config.Duration
}

func (c JSONExecOffchainConfig) Validate() error {
	if c.DestOptimisticConfirmations == 0 {
		return errors.New("must set DestOptimisticConfirmations")
	}
	if c.BatchGasLimit == 0 {
		return errors.New("must set BatchGasLimit")
	}
	if c.RelativeBoostPerWaitHour == 0 {
		return errors.New("must set RelativeBoostPerWaitHour")
	}
	if c.InflightCacheExpiry.Duration() == 0 {
		return errors.New("must set InflightCacheExpiry")
	}
	if c.RootSnoozeTime.Duration() == 0 {
		return errors.New("must set RootSnoozeTime")
	}

	return nil
}

type OffRamp struct {
	offRampV120 evm_2_evm_offramp_1_2_0.EVM2EVMOffRampInterface

	addr                    common.Address
	lp                      logpoller.LogPoller
	Logger                  logger.Logger
	Client                  client.Client
	evmBatchCaller          rpclib.EvmBatchCaller
	filters                 []logpoller.Filter
	Estimator               gas.EvmFeeEstimator
	DestMaxGasPrice         *big.Int
	ExecutionReportArgs     abi.Arguments
	eventIndex              int
	eventSig                common.Hash
	cachedOffRampTokens     cache.AutoSync[cciptypes.OffRampTokens]
	sourceToDestTokensCache sync.Map

	// Dynamic config
	// configMu guards all the dynamic config fields.
	configMu           sync.RWMutex
	gasPriceEstimator  prices.GasPriceEstimatorExec
	offchainConfig     cciptypes.ExecOffchainConfig
	onchainConfig      cciptypes.ExecOnchainConfig
	feeEstimatorConfig ccipdata.FeeEstimatorConfigReader
}

func (o *OffRamp) GetStaticConfig(ctx context.Context) (cciptypes.OffRampStaticConfig, error) {
	if o.offRampV120 == nil {
		return cciptypes.OffRampStaticConfig{}, errors.New("offramp not initialized")
	}
	c, err := o.offRampV120.GetStaticConfig(&bind.CallOpts{Context: ctx})
	if err != nil {
		return cciptypes.OffRampStaticConfig{}, fmt.Errorf("error while retrieving offramp config: %w", err)
	}
	return cciptypes.OffRampStaticConfig{
		CommitStore:         cciptypes.Address(c.CommitStore.String()),
		ChainSelector:       c.ChainSelector,
		SourceChainSelector: c.SourceChainSelector,
		OnRamp:              cciptypes.Address(c.OnRamp.String()),
		PrevOffRamp:         cciptypes.Address(c.PrevOffRamp.String()),
		ArmProxy:            cciptypes.Address(c.ArmProxy.String()),
	}, nil
}

func (o *OffRamp) GetSenderNonce(ctx context.Context, sender cciptypes.Address) (uint64, error) {
	evmAddr, err := ccipcalc.GenericAddrToEvm(sender)
	if err != nil {
		return 0, err
	}
	return o.offRampV120.GetSenderNonce(&bind.CallOpts{Context: ctx}, evmAddr)
}

func (o *OffRamp) ListSenderNonces(ctx context.Context, senders []cciptypes.Address) (map[cciptypes.Address]uint64, error) {
	if len(senders) == 0 {
		return make(map[cciptypes.Address]uint64), nil
	}

	evmSenders, err := ccipcalc.GenericAddrsToEvm(senders...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert generic addresses to evm addresses")
	}

	evmCalls := make([]rpclib.EvmCall, 0, len(evmSenders))
	for _, evmAddr := range evmSenders {
		evmCalls = append(evmCalls, rpclib.NewEvmCall(
			abiOffRamp,
			"getSenderNonce",
			o.addr,
			evmAddr,
		))
	}

	results, err := o.evmBatchCaller.BatchCall(ctx, 0, evmCalls)
	if err != nil {
		o.Logger.Errorw("error while batch fetching sender nonces", "err", err, "senders", evmSenders)
		return nil, err
	}

	nonces, err := rpclib.ParseOutputs[uint64](results, func(d rpclib.DataAndErr) (uint64, error) {
		return rpclib.ParseOutput[uint64](d, 0)
	})
	if err != nil {
		o.Logger.Errorw("error while parsing sender nonces", "err", err, "senders", evmSenders)
		return nil, err
	}

	if len(senders) != len(nonces) {
		o.Logger.Errorw("unexpected number of nonces returned", "senders", evmSenders, "nonces", nonces)
		return nil, errors.New("unexpected number of nonces returned")
	}

	senderNonce := make(map[cciptypes.Address]uint64, len(senders))
	for i, sender := range senders {
		senderNonce[sender] = nonces[i]
	}
	return senderNonce, nil
}

func (o *OffRamp) CurrentRateLimiterState(ctx context.Context) (cciptypes.TokenBucketRateLimit, error) {
	bucket, err := o.offRampV120.CurrentRateLimiterState(&bind.CallOpts{Context: ctx})
	if err != nil {
		return cciptypes.TokenBucketRateLimit{}, err
	}
	return cciptypes.TokenBucketRateLimit{
		Tokens:      bucket.Tokens,
		LastUpdated: bucket.LastUpdated,
		IsEnabled:   bucket.IsEnabled,
		Capacity:    bucket.Capacity,
		Rate:        bucket.Rate,
	}, nil
}

func (o *OffRamp) getDestinationTokensFromSourceTokens(ctx context.Context, tokenAddresses []cciptypes.Address) ([]cciptypes.Address, error) {
	destTokens := make([]cciptypes.Address, len(tokenAddresses))
	found := make(map[cciptypes.Address]bool)

	for i, tokenAddress := range tokenAddresses {
		if v, exists := o.sourceToDestTokensCache.Load(tokenAddress); exists {
			if destToken, isAddr := v.(cciptypes.Address); isAddr {
				destTokens[i] = destToken
				found[tokenAddress] = true
			} else {
				o.Logger.Errorf("source to dest cache contains invalid type %T", v)
			}
		}
	}

	if len(found) == len(tokenAddresses) {
		return destTokens, nil
	}

	evmAddrs, err := ccipcalc.GenericAddrsToEvm(tokenAddresses...)
	if err != nil {
		return nil, err
	}

	evmCalls := make([]rpclib.EvmCall, 0, len(tokenAddresses))
	for i, sourceTk := range tokenAddresses {
		if !found[sourceTk] {
			evmCalls = append(evmCalls, rpclib.NewEvmCall(abiOffRamp, "getDestinationToken", o.addr, evmAddrs[i]))
		}
	}

	results, err := o.evmBatchCaller.BatchCall(ctx, 0, evmCalls)
	if err != nil {
		return nil, fmt.Errorf("batch call limit: %w", err)
	}

	destTokensFromRPC, err := rpclib.ParseOutputs[common.Address](results, func(d rpclib.DataAndErr) (common.Address, error) {
		return rpclib.ParseOutput[common.Address](d, 0)
	})
	if err != nil {
		return nil, fmt.Errorf("parse outputs: %w", err)
	}

	j := 0
	for i, sourceToken := range tokenAddresses {
		if !found[sourceToken] {
			destTokens[i] = cciptypes.Address(destTokensFromRPC[j].String())
			o.sourceToDestTokensCache.Store(sourceToken, destTokens[i])
			j++
		}
	}

	seenDestTokens := mapset.NewSet[cciptypes.Address]()
	for _, destToken := range destTokens {
		if seenDestTokens.Contains(destToken) {
			return nil, fmt.Errorf("offRamp misconfig, destination token %s already exists", destToken)
		}
		seenDestTokens.Add(destToken)
	}

	return destTokens, nil
}

func (o *OffRamp) GetSourceToDestTokensMapping(ctx context.Context) (map[cciptypes.Address]cciptypes.Address, error) {
	tokens, err := o.GetTokens(ctx)
	if err != nil {
		return nil, err
	}

	destTokens, err := o.getDestinationTokensFromSourceTokens(ctx, tokens.SourceTokens)
	if err != nil {
		return nil, fmt.Errorf("get destination tokens from source tokens: %w", err)
	}

	srcToDstTokenMapping := make(map[cciptypes.Address]cciptypes.Address, len(tokens.SourceTokens))
	for i, sourceToken := range tokens.SourceTokens {
		srcToDstTokenMapping[sourceToken] = destTokens[i]
	}
	return srcToDstTokenMapping, nil
}

func (o *OffRamp) GetTokens(ctx context.Context) (cciptypes.OffRampTokens, error) {
	return o.cachedOffRampTokens.Get(ctx, func(ctx context.Context) (cciptypes.OffRampTokens, error) {
		destTokens, err := o.offRampV120.GetDestinationTokens(&bind.CallOpts{Context: ctx})
		if err != nil {
			return cciptypes.OffRampTokens{}, fmt.Errorf("get destination tokens: %w", err)
		}
		sourceTokens, err := o.offRampV120.GetSupportedTokens(&bind.CallOpts{Context: ctx})
		if err != nil {
			return cciptypes.OffRampTokens{}, err
		}

		return cciptypes.OffRampTokens{
			DestinationTokens: ccipcalc.EvmAddrsToGeneric(destTokens...),
			SourceTokens:      ccipcalc.EvmAddrsToGeneric(sourceTokens...),
		}, nil
	})
}

func (o *OffRamp) GetRouter(ctx context.Context) (cciptypes.Address, error) {
	dynamicConfig, err := o.offRampV120.GetDynamicConfig(&bind.CallOpts{Context: ctx})
	if err != nil {
		return "", err
	}
	return ccipcalc.EvmAddrToGeneric(dynamicConfig.Router), nil
}

func (o *OffRamp) OffchainConfig(ctx context.Context) (cciptypes.ExecOffchainConfig, error) {
	o.configMu.RLock()
	defer o.configMu.RUnlock()
	return o.offchainConfig, nil
}

func (o *OffRamp) OnchainConfig(ctx context.Context) (cciptypes.ExecOnchainConfig, error) {
	o.configMu.RLock()
	defer o.configMu.RUnlock()
	return o.onchainConfig, nil
}

func (o *OffRamp) GasPriceEstimator(ctx context.Context) (cciptypes.GasPriceEstimatorExec, error) {
	o.configMu.RLock()
	defer o.configMu.RUnlock()
	return o.gasPriceEstimator, nil
}

func (o *OffRamp) Address(ctx context.Context) (cciptypes.Address, error) {
	return cciptypes.Address(o.addr.String()), nil
}

func (o *OffRamp) UpdateDynamicConfig(onchainConfig cciptypes.ExecOnchainConfig, offchainConfig cciptypes.ExecOffchainConfig, gasPriceEstimator prices.GasPriceEstimatorExec) {
	o.configMu.Lock()
	o.onchainConfig = onchainConfig
	o.offchainConfig = offchainConfig
	o.gasPriceEstimator = gasPriceEstimator
	o.configMu.Unlock()
}

func (o *OffRamp) ChangeConfig(ctx context.Context, onchainConfigBytes []byte, offchainConfigBytes []byte) (cciptypes.Address, cciptypes.Address, error) {
	// Same as the v1.0.0 method, except for the ExecOnchainConfig type.
	onchainConfigParsed, err := abihelpers.DecodeAbiStruct[ExecOnchainConfig](onchainConfigBytes)
	if err != nil {
		return "", "", err
	}

	offchainConfigParsed, err := ccipconfig.DecodeOffchainConfig[JSONExecOffchainConfig](offchainConfigBytes)
	if err != nil {
		return "", "", err
	}
	destRouter, err := router.NewRouter(onchainConfigParsed.Router, o.Client)
	if err != nil {
		return "", "", err
	}
	destWrappedNative, err := destRouter.GetWrappedNative(nil)
	if err != nil {
		return "", "", err
	}
	offchainConfig := cciptypes.ExecOffchainConfig{
		DestOptimisticConfirmations: offchainConfigParsed.DestOptimisticConfirmations,
		BatchGasLimit:               offchainConfigParsed.BatchGasLimit,
		RelativeBoostPerWaitHour:    offchainConfigParsed.RelativeBoostPerWaitHour,
		InflightCacheExpiry:         offchainConfigParsed.InflightCacheExpiry,
		RootSnoozeTime:              offchainConfigParsed.RootSnoozeTime,
		MessageVisibilityInterval:   offchainConfigParsed.MessageVisibilityInterval,
		BatchingStrategyID:          offchainConfigParsed.BatchingStrategyID,
	}
	onchainConfig := cciptypes.ExecOnchainConfig{
		PermissionLessExecutionThresholdSeconds: time.Second * time.Duration(onchainConfigParsed.PermissionLessExecutionThresholdSeconds),
		Router:                                  cciptypes.Address(onchainConfigParsed.Router.String()),
	}
	priceEstimator := prices.NewDAGasPriceEstimator(o.Estimator, o.DestMaxGasPrice, 0, 0, o.feeEstimatorConfig)

	o.UpdateDynamicConfig(onchainConfig, offchainConfig, priceEstimator)

	o.Logger.Infow("Starting exec plugin",
		"offchainConfig", onchainConfigParsed,
		"onchainConfig", offchainConfigParsed)
	return cciptypes.Address(onchainConfigParsed.PriceRegistry.String()),
		cciptypes.Address(destWrappedNative.String()), nil
}

func (o *OffRamp) Close() error {
	return logpollerutil.UnregisterLpFilters(context.Background(), o.lp, o.filters)
}
func (o *OffRamp) RegisterFilters(ctx context.Context) error {
	return logpollerutil.RegisterLpFilters(ctx, o.lp, o.filters)
}

func (o *OffRamp) GetExecutionState(ctx context.Context, sequenceNumber uint64) (uint8, error) {
	return o.offRampV120.GetExecutionState(&bind.CallOpts{Context: ctx}, sequenceNumber)
}

func (o *OffRamp) GetExecutionStateChangesBetweenSeqNums(ctx context.Context, seqNumMin, seqNumMax uint64, confs int) ([]cciptypes.ExecutionStateChangedWithTxMeta, error) {
	latestBlock, err := o.lp.LatestBlock(ctx)
	if err != nil {
		return nil, fmt.Errorf("get lp latest block: %w", err)
	}

	logs, err := o.lp.IndexedLogsTopicRange(
		ctx,
		o.eventSig,
		o.addr,
		o.eventIndex,
		logpoller.EvmWord(seqNumMin),
		logpoller.EvmWord(seqNumMax),
		evmtypes.Confirmations(confs),
	)
	if err != nil {
		return nil, err
	}

	parsedLogs, err := ccipdata.ParseLogs[cciptypes.ExecutionStateChanged](
		logs,
		o.Logger,
		func(log types.Log) (*cciptypes.ExecutionStateChanged, error) {
			sc, err1 := o.offRampV120.ParseExecutionStateChanged(log)
			if err1 != nil {
				return nil, err1
			}

			return &cciptypes.ExecutionStateChanged{
				SequenceNumber: sc.SequenceNumber,
			}, nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("parse logs: %w", err)
	}

	res := make([]cciptypes.ExecutionStateChangedWithTxMeta, 0, len(parsedLogs))
	for _, log := range parsedLogs {
		res = append(res, cciptypes.ExecutionStateChangedWithTxMeta{
			TxMeta:                log.TxMeta.WithFinalityStatus(uint64(latestBlock.FinalizedBlockNumber)),
			ExecutionStateChanged: log.Data,
		})
	}
	return res, nil
}

func EncodeExecutionReport(ctx context.Context, args abi.Arguments, report cciptypes.ExecReport) ([]byte, error) {
	var msgs []evm_2_evm_offramp_1_2_0.InternalEVM2EVMMessage
	for _, msg := range report.Messages {
		var ta []evm_2_evm_offramp_1_2_0.ClientEVMTokenAmount
		for _, tokenAndAmount := range msg.TokenAmounts {
			evmAddrs, err := ccipcalc.GenericAddrsToEvm(tokenAndAmount.Token)
			if err != nil {
				return nil, err
			}
			ta = append(ta, evm_2_evm_offramp_1_2_0.ClientEVMTokenAmount{
				Token:  evmAddrs[0],
				Amount: tokenAndAmount.Amount,
			})
		}

		evmAddrs, err := ccipcalc.GenericAddrsToEvm(msg.Sender, msg.Receiver, msg.FeeToken)
		if err != nil {
			return nil, err
		}

		msgs = append(msgs, evm_2_evm_offramp_1_2_0.InternalEVM2EVMMessage{
			SourceChainSelector: msg.SourceChainSelector,
			Sender:              evmAddrs[0],
			Receiver:            evmAddrs[1],
			SequenceNumber:      msg.SequenceNumber,
			GasLimit:            msg.GasLimit,
			Strict:              msg.Strict,
			Nonce:               msg.Nonce,
			FeeToken:            evmAddrs[2],
			FeeTokenAmount:      msg.FeeTokenAmount,
			Data:                msg.Data,
			TokenAmounts:        ta,
			MessageId:           msg.MessageID,
			// NOTE: this field is new in v1.2.
			SourceTokenData: msg.SourceTokenData,
		})
	}

	rep := evm_2_evm_offramp_1_2_0.InternalExecutionReport{
		Messages:          msgs,
		OffchainTokenData: report.OffchainTokenData,
		Proofs:            report.Proofs,
		ProofFlagBits:     report.ProofFlagBits,
	}
	return args.PackValues([]interface{}{&rep})
}

func (o *OffRamp) EncodeExecutionReport(ctx context.Context, report cciptypes.ExecReport) ([]byte, error) {
	return EncodeExecutionReport(ctx, o.ExecutionReportArgs, report)
}

func DecodeExecReport(ctx context.Context, args abi.Arguments, report []byte) (cciptypes.ExecReport, error) {
	unpacked, err := args.Unpack(report)
	if err != nil {
		return cciptypes.ExecReport{}, err
	}
	if len(unpacked) == 0 {
		return cciptypes.ExecReport{}, errors.New("assumptionViolation: expected at least one element")
	}
	// Must be anonymous struct here
	erStruct, ok := unpacked[0].(struct {
		Messages []struct {
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
			SourceTokenData [][]uint8 `json:"sourceTokenData"`
			MessageId       [32]uint8 `json:"messageId"`
		} `json:"messages"`
		OffchainTokenData [][][]uint8 `json:"offchainTokenData"`
		Proofs            [][32]uint8 `json:"proofs"`
		ProofFlagBits     *big.Int    `json:"proofFlagBits"`
	})
	if !ok {
		return cciptypes.ExecReport{}, fmt.Errorf("got %T", unpacked[0])
	}
	messages := make([]cciptypes.EVM2EVMMessage, 0, len(erStruct.Messages))
	for _, msg := range erStruct.Messages {
		var tokensAndAmounts []cciptypes.TokenAmount
		for _, tokenAndAmount := range msg.TokenAmounts {
			tokensAndAmounts = append(tokensAndAmounts, cciptypes.TokenAmount{
				Token:  cciptypes.Address(tokenAndAmount.Token.String()),
				Amount: tokenAndAmount.Amount,
			})
		}
		messages = append(messages, cciptypes.EVM2EVMMessage{
			SequenceNumber:      msg.SequenceNumber,
			GasLimit:            msg.GasLimit,
			Nonce:               msg.Nonce,
			MessageID:           msg.MessageId,
			SourceChainSelector: msg.SourceChainSelector,
			Sender:              cciptypes.Address(msg.Sender.String()),
			Receiver:            cciptypes.Address(msg.Receiver.String()),
			Strict:              msg.Strict,
			FeeToken:            cciptypes.Address(msg.FeeToken.String()),
			FeeTokenAmount:      msg.FeeTokenAmount,
			Data:                msg.Data,
			TokenAmounts:        tokensAndAmounts,
			SourceTokenData:     msg.SourceTokenData,
			// TODO: Not needed for plugins, but should be recomputed for consistency.
			// Requires the offramp knowing about onramp version
			Hash: [32]byte{},
		})
	}

	// Unpack will populate with big.Int{false, <allocated empty nat>} for 0 values,
	// which is different from the expected big.NewInt(0). Rebuild to the expected value for this case.
	return cciptypes.ExecReport{
		Messages:          messages,
		OffchainTokenData: erStruct.OffchainTokenData,
		Proofs:            erStruct.Proofs,
		ProofFlagBits:     new(big.Int).SetBytes(erStruct.ProofFlagBits.Bytes()),
	}, nil
}

func (o *OffRamp) DecodeExecutionReport(ctx context.Context, report []byte) (cciptypes.ExecReport, error) {
	return DecodeExecReport(ctx, o.ExecutionReportArgs, report)
}

func NewOffRamp(lggr logger.Logger, addr common.Address, ec client.Client, lp logpoller.LogPoller, estimator gas.EvmFeeEstimator, destMaxGasPrice *big.Int, feeEstimatorConfig ccipdata.FeeEstimatorConfigReader) (*OffRamp, error) {
	offRamp, err := evm_2_evm_offramp_1_2_0.NewEVM2EVMOffRamp(addr, ec)
	if err != nil {
		return nil, err
	}

	executionStateChangedSequenceNumberIndex := 1
	executionReportArgs := abihelpers.MustGetMethodInputs("manuallyExecute", abiOffRamp)[:1]
	filters := []logpoller.Filter{
		{
			Name:      logpoller.FilterName(ExecExecutionStateChanges, addr.String()),
			EventSigs: []common.Hash{ExecutionStateChangedEvent},
			Addresses: []common.Address{addr},
			Retention: ccipdata.CommitExecLogsRetention,
		},
		{
			Name:      logpoller.FilterName(ExecTokenPoolAdded, addr.String()),
			EventSigs: []common.Hash{PoolAddedEvent},
			Addresses: []common.Address{addr},
			Retention: ccipdata.CacheEvictionLogsRetention,
		},
		{
			Name:      logpoller.FilterName(ExecTokenPoolRemoved, addr.String()),
			EventSigs: []common.Hash{PoolRemovedEvent},
			Addresses: []common.Address{addr},
			Retention: ccipdata.CacheEvictionLogsRetention,
		},
	}

	return &OffRamp{
		offRampV120:         offRamp,
		Client:              ec,
		addr:                addr,
		Logger:              lggr,
		lp:                  lp,
		filters:             filters,
		Estimator:           estimator,
		DestMaxGasPrice:     destMaxGasPrice,
		ExecutionReportArgs: executionReportArgs,
		eventSig:            ExecutionStateChangedEvent,
		eventIndex:          executionStateChangedSequenceNumberIndex,
		configMu:            sync.RWMutex{},
		evmBatchCaller: rpclib.NewDynamicLimitedBatchCaller(
			lggr,
			ec,
			rpclib.DefaultRpcBatchSizeLimit,
			rpclib.DefaultRpcBatchBackOffMultiplier,
			rpclib.DefaultMaxParallelRpcCalls,
		),
		cachedOffRampTokens: cache.NewLogpollerEventsBased[cciptypes.OffRampTokens](
			lp,
			offrampPoolAddedPoolRemovedEvents,
			offRamp.Address(),
		),
		// values set on the fly after ChangeConfig is called
		gasPriceEstimator:  prices.ExecGasPriceEstimator{},
		offchainConfig:     cciptypes.ExecOffchainConfig{},
		onchainConfig:      cciptypes.ExecOnchainConfig{},
		feeEstimatorConfig: feeEstimatorConfig,
	}, nil
}
