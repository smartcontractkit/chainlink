package ccipdata

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/custom_token_pool"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_offramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_offramp_1_0_0"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/logpollerutil"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/rpclib"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/prices"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

const (
	EXEC_EXECUTION_STATE_CHANGES = "Exec execution state changes"
	EXEC_TOKEN_POOL_ADDED        = "Token pool added"
	EXEC_TOKEN_POOL_REMOVED      = "Token pool removed"
)

var (
	abiOffRampV1_0_0                                    = abihelpers.MustParseABI(evm_2_evm_offramp_1_0_0.EVM2EVMOffRampABI)
	abiCustomTokenPool                                  = abihelpers.MustParseABI(custom_token_pool.CustomTokenPoolABI)
	_                                     OffRampReader = &OffRampV1_0_0{}
	ExecutionStateChangedEventV1_0_0                    = abihelpers.MustGetEventID("ExecutionStateChanged", abiOffRampV1_0_0)
	ExecutionStateChangedSeqNrIndexV1_0_0               = 1
)

type ExecOnchainConfigV1_0_0 evm_2_evm_offramp_1_0_0.EVM2EVMOffRampDynamicConfig

func (d ExecOnchainConfigV1_0_0) AbiString() string {
	return `
	[
		{
			"components": [
				{"name": "permissionLessExecutionThresholdSeconds", "type": "uint32"},
				{"name": "router", "type": "address"},
				{"name": "priceRegistry", "type": "address"},
				{"name": "maxTokensLength", "type": "uint16"},
				{"name": "maxDataSize", "type": "uint32"}
			],
			"type": "tuple"
		}
	]`
}

func (d ExecOnchainConfigV1_0_0) Validate() error {
	if d.PermissionLessExecutionThresholdSeconds == 0 {
		return errors.New("must set PermissionLessExecutionThresholdSeconds")
	}
	if d.Router == (common.Address{}) {
		return errors.New("must set Router address")
	}
	if d.PriceRegistry == (common.Address{}) {
		return errors.New("must set PriceRegistry address")
	}
	if d.MaxTokensLength == 0 {
		return errors.New("must set MaxTokensLength")
	}
	if d.MaxDataSize == 0 {
		return errors.New("must set MaxDataSize")
	}
	return nil
}

func (d ExecOnchainConfigV1_0_0) PermissionLessExecutionThresholdDuration() time.Duration {
	return time.Duration(d.PermissionLessExecutionThresholdSeconds) * time.Second
}

type OffRampV1_0_0 struct {
	offRamp             *evm_2_evm_offramp_1_0_0.EVM2EVMOffRamp
	addr                common.Address
	lp                  logpoller.LogPoller
	lggr                logger.Logger
	ec                  client.Client
	evmBatchCaller      rpclib.EvmBatchCaller
	filters             []logpoller.Filter
	estimator           gas.EvmFeeEstimator
	executionReportArgs abi.Arguments
	eventIndex          int
	eventSig            common.Hash

	// Dynamic config
	configMu          sync.RWMutex
	gasPriceEstimator prices.GasPriceEstimatorExec
	offchainConfig    ExecOffchainConfig
	onchainConfig     ExecOnchainConfig
}

func (o *OffRampV1_0_0) GetStaticConfig(ctx context.Context) (OffRampStaticConfig, error) {
	if o.offRamp == nil {
		return OffRampStaticConfig{}, fmt.Errorf("offramp not initialized")
	}
	c, err := o.offRamp.GetStaticConfig(&bind.CallOpts{Context: ctx})
	if err != nil {
		return OffRampStaticConfig{}, fmt.Errorf("error while retrieving offramp config: %w", err)
	}
	return OffRampStaticConfig{
		CommitStore:         c.CommitStore,
		ChainSelector:       c.ChainSelector,
		SourceChainSelector: c.SourceChainSelector,
		OnRamp:              c.OnRamp,
		PrevOffRamp:         c.PrevOffRamp,
		ArmProxy:            c.ArmProxy,
	}, nil
}

func (o *OffRampV1_0_0) GetExecutionState(ctx context.Context, sequenceNumber uint64) (uint8, error) {
	return o.offRamp.GetExecutionState(&bind.CallOpts{Context: ctx}, sequenceNumber)
}

func (o *OffRampV1_0_0) GetSenderNonce(ctx context.Context, sender common.Address) (uint64, error) {
	return o.offRamp.GetSenderNonce(&bind.CallOpts{Context: ctx}, sender)
}

func (o *OffRampV1_0_0) CurrentRateLimiterState(ctx context.Context) (evm_2_evm_offramp.RateLimiterTokenBucket, error) {
	state, err := o.offRamp.CurrentRateLimiterState(&bind.CallOpts{Context: ctx})
	if err != nil {
		return *new(evm_2_evm_offramp.RateLimiterTokenBucket), err
	}
	return evm_2_evm_offramp.RateLimiterTokenBucket{
		Tokens:      state.Tokens,
		LastUpdated: state.LastUpdated,
		IsEnabled:   state.IsEnabled,
		Capacity:    state.Capacity,
		Rate:        state.Rate,
	}, nil
}

func (o *OffRampV1_0_0) GetDestinationToken(ctx context.Context, address common.Address) (common.Address, error) {
	return o.offRamp.GetDestinationToken(&bind.CallOpts{Context: ctx}, address)
}

func (o *OffRampV1_0_0) GetDestinationTokensFromSourceTokens(ctx context.Context, tokenAddresses []common.Address) ([]common.Address, error) {
	if len(tokenAddresses) == 0 {
		return []common.Address{}, nil
	}

	evmCalls := make([]rpclib.EvmCall, 0, len(tokenAddresses))
	for _, sourceTk := range tokenAddresses {
		evmCalls = append(evmCalls, rpclib.NewEvmCall(abiOffRampV1_0_0, "getDestinationToken", o.addr, sourceTk))
	}

	latestBlock, err := o.lp.LatestBlock(pg.WithParentCtx(ctx))
	if err != nil {
		return nil, fmt.Errorf("get latest block: %w", err)
	}

	results, err := o.evmBatchCaller.BatchCall(ctx, uint64(latestBlock), evmCalls)
	if err != nil {
		return nil, fmt.Errorf("batch call limit: %w", err)
	}

	destTokens, err := rpclib.ParseOutputs[common.Address](results, func(d rpclib.DataAndErr) (common.Address, error) {
		return rpclib.ParseOutput[common.Address](d, 0)
	})
	if err != nil {
		return nil, fmt.Errorf("parse outputs: %w", err)
	}

	seenDestTokens := make(map[common.Address]struct{})
	for _, destToken := range destTokens {
		if _, exists := seenDestTokens[destToken]; exists {
			return nil, fmt.Errorf("offRamp misconfig, destination token %s already exists", destToken)
		}
		seenDestTokens[destToken] = struct{}{}
	}

	return destTokens, nil
}

func (o *OffRampV1_0_0) GetTokenPoolsRateLimits(ctx context.Context, poolAddresses []common.Address) ([]TokenBucketRateLimit, error) {
	if len(poolAddresses) == 0 {
		return nil, nil
	}

	evmCalls := make([]rpclib.EvmCall, 0, len(poolAddresses))
	for _, poolAddress := range poolAddresses {
		evmCalls = append(evmCalls, rpclib.NewEvmCall(
			abiCustomTokenPool,
			"currentOffRampRateLimiterState",
			poolAddress,
			o.addr,
		))
	}

	latestBlock, err := o.lp.LatestBlock(pg.WithParentCtx(ctx))
	if err != nil {
		return nil, fmt.Errorf("get latest block: %w", err)
	}

	results, err := o.evmBatchCaller.BatchCall(ctx, uint64(latestBlock), evmCalls)
	if err != nil {
		return nil, fmt.Errorf("batch call limit: %w", err)
	}

	rateLimits, err := rpclib.ParseOutputs[TokenBucketRateLimit](results, func(d rpclib.DataAndErr) (TokenBucketRateLimit, error) {
		return rpclib.ParseOutput[TokenBucketRateLimit](d, 0)
	})
	if err != nil {
		return nil, fmt.Errorf("parse outputs: %w", err)
	}

	return rateLimits, nil
}

func (o *OffRampV1_0_0) GetSupportedTokens(ctx context.Context) ([]common.Address, error) {
	return o.offRamp.GetSupportedTokens(&bind.CallOpts{Context: ctx})
}

func (o *OffRampV1_0_0) GetPoolByDestToken(ctx context.Context, address common.Address) (common.Address, error) {
	return o.offRamp.GetPoolByDestToken(&bind.CallOpts{Context: ctx}, address)
}

func (o *OffRampV1_0_0) OffchainConfig() ExecOffchainConfig {
	o.configMu.RLock()
	defer o.configMu.RUnlock()
	return o.offchainConfig
}

func (o *OffRampV1_0_0) OnchainConfig() ExecOnchainConfig {
	o.configMu.RLock()
	defer o.configMu.RUnlock()
	return o.onchainConfig
}

func (o *OffRampV1_0_0) GasPriceEstimator() prices.GasPriceEstimatorExec {
	o.configMu.RLock()
	defer o.configMu.RUnlock()
	return o.gasPriceEstimator
}

func (o *OffRampV1_0_0) Address() common.Address {
	return o.addr
}

func (o *OffRampV1_0_0) ChangeConfig(onchainConfig []byte, offchainConfig []byte) (common.Address, common.Address, error) {
	onchainConfigParsed, err := abihelpers.DecodeAbiStruct[ExecOnchainConfigV1_0_0](onchainConfig)
	if err != nil {
		return common.Address{}, common.Address{}, err
	}

	offchainConfigParsed, err := ccipconfig.DecodeOffchainConfig[ExecOffchainConfig](offchainConfig)
	if err != nil {
		return common.Address{}, common.Address{}, err
	}
	destRouter, err := router.NewRouter(onchainConfigParsed.Router, o.ec)
	if err != nil {
		return common.Address{}, common.Address{}, err
	}
	destWrappedNative, err := destRouter.GetWrappedNative(nil)
	if err != nil {
		return common.Address{}, common.Address{}, err
	}
	o.configMu.Lock()
	o.offchainConfig = ExecOffchainConfig{
		SourceFinalityDepth:         offchainConfigParsed.SourceFinalityDepth,
		DestFinalityDepth:           offchainConfigParsed.DestFinalityDepth,
		DestOptimisticConfirmations: offchainConfigParsed.DestOptimisticConfirmations,
		BatchGasLimit:               offchainConfigParsed.BatchGasLimit,
		RelativeBoostPerWaitHour:    offchainConfigParsed.RelativeBoostPerWaitHour,
		MaxGasPrice:                 offchainConfigParsed.MaxGasPrice,
		InflightCacheExpiry:         offchainConfigParsed.InflightCacheExpiry,
		RootSnoozeTime:              offchainConfigParsed.RootSnoozeTime,
	}
	o.onchainConfig = ExecOnchainConfig{PermissionLessExecutionThresholdSeconds: time.Second * time.Duration(onchainConfigParsed.PermissionLessExecutionThresholdSeconds)}
	o.gasPriceEstimator = prices.NewExecGasPriceEstimator(o.estimator, big.NewInt(int64(offchainConfigParsed.MaxGasPrice)), 0)
	o.configMu.Unlock()

	o.lggr.Infow("Starting exec plugin",
		"offchainConfig", onchainConfigParsed,
		"onchainConfig", offchainConfigParsed)
	return onchainConfigParsed.PriceRegistry, destWrappedNative, nil
}

func (o *OffRampV1_0_0) GetDestinationTokens(ctx context.Context) ([]common.Address, error) {
	return o.offRamp.GetDestinationTokens(&bind.CallOpts{Context: ctx})
}

func (o *OffRampV1_0_0) Close(qopts ...pg.QOpt) error {
	return logpollerutil.UnregisterLpFilters(o.lp, o.filters, qopts...)
}

func (o *OffRampV1_0_0) GetExecutionStateChangesBetweenSeqNums(ctx context.Context, seqNumMin, seqNumMax uint64, confs int) ([]Event[ExecutionStateChanged], error) {
	logs, err := o.lp.IndexedLogsTopicRange(
		o.eventSig,
		o.addr,
		o.eventIndex,
		logpoller.EvmWord(seqNumMin),
		logpoller.EvmWord(seqNumMax),
		logpoller.Confirmations(confs),
		pg.WithParentCtx(ctx),
	)
	if err != nil {
		return nil, err
	}

	return parseLogs[ExecutionStateChanged](
		logs,
		o.lggr,
		func(log types.Log) (*ExecutionStateChanged, error) {
			sc, err := o.offRamp.ParseExecutionStateChanged(log)
			if err != nil {
				return nil, err
			}
			return &ExecutionStateChanged{SequenceNumber: sc.SequenceNumber}, nil
		},
	)
}

func encodeExecutionReportV1_0_0(args abi.Arguments, report ExecReport) ([]byte, error) {
	var msgs []evm_2_evm_offramp_1_0_0.InternalEVM2EVMMessage
	for _, msg := range report.Messages {
		var ta []evm_2_evm_offramp_1_0_0.ClientEVMTokenAmount
		for _, tokenAndAmount := range msg.TokenAmounts {
			ta = append(ta, evm_2_evm_offramp_1_0_0.ClientEVMTokenAmount{
				Token:  tokenAndAmount.Token,
				Amount: tokenAndAmount.Amount,
			})
		}
		msgs = append(msgs, evm_2_evm_offramp_1_0_0.InternalEVM2EVMMessage{
			SourceChainSelector: msg.SourceChainSelector,
			Sender:              msg.Sender,
			Receiver:            msg.Receiver,
			SequenceNumber:      msg.SequenceNumber,
			GasLimit:            msg.GasLimit,
			Strict:              msg.Strict,
			Nonce:               msg.Nonce,
			FeeToken:            msg.FeeToken,
			FeeTokenAmount:      msg.FeeTokenAmount,
			Data:                msg.Data,
			TokenAmounts:        ta,
			MessageId:           msg.MessageId,
		})
	}

	rep := evm_2_evm_offramp_1_0_0.InternalExecutionReport{
		Messages:          msgs,
		OffchainTokenData: report.OffchainTokenData,
		Proofs:            report.Proofs,
		ProofFlagBits:     report.ProofFlagBits,
	}
	return args.PackValues([]interface{}{&rep})
}

func (o *OffRampV1_0_0) EncodeExecutionReport(report ExecReport) ([]byte, error) {
	return encodeExecutionReportV1_0_0(o.executionReportArgs, report)
}

func decodeExecReportV1_0_0(args abi.Arguments, report []byte) (ExecReport, error) {
	unpacked, err := args.Unpack(report)
	if err != nil {
		return ExecReport{}, err
	}
	if len(unpacked) == 0 {
		return ExecReport{}, errors.New("assumptionViolation: expected at least one element")
	}

	erStruct, ok := unpacked[0].(struct {
		Messages []struct {
			SourceChainSelector uint64         `json:"sourceChainSelector"`
			SequenceNumber      uint64         `json:"sequenceNumber"`
			FeeTokenAmount      *big.Int       `json:"feeTokenAmount"`
			Sender              common.Address `json:"sender"`
			Nonce               uint64         `json:"nonce"`
			GasLimit            *big.Int       `json:"gasLimit"`
			Strict              bool           `json:"strict"`
			Receiver            common.Address `json:"receiver"`
			Data                []uint8        `json:"data"`
			TokenAmounts        []struct {
				Token  common.Address `json:"token"`
				Amount *big.Int       `json:"amount"`
			} `json:"tokenAmounts"`
			FeeToken  common.Address `json:"feeToken"`
			MessageId [32]uint8      `json:"messageId"`
		} `json:"messages"`
		OffchainTokenData [][][]uint8 `json:"offchainTokenData"`
		Proofs            [][32]uint8 `json:"proofs"`
		ProofFlagBits     *big.Int    `json:"proofFlagBits"`
	})

	if !ok {
		return ExecReport{}, fmt.Errorf("got %T", unpacked[0])
	}
	messages := []internal.EVM2EVMMessage{}
	for _, msg := range erStruct.Messages {
		var tokensAndAmounts []internal.TokenAmount
		for _, tokenAndAmount := range msg.TokenAmounts {
			tokensAndAmounts = append(tokensAndAmounts, internal.TokenAmount{
				Token:  tokenAndAmount.Token,
				Amount: tokenAndAmount.Amount,
			})
		}
		messages = append(messages, internal.EVM2EVMMessage{
			SequenceNumber:      msg.SequenceNumber,
			GasLimit:            msg.GasLimit,
			Nonce:               msg.Nonce,
			MessageId:           msg.MessageId,
			SourceChainSelector: msg.SourceChainSelector,
			Sender:              msg.Sender,
			Receiver:            msg.Receiver,
			Strict:              msg.Strict,
			FeeToken:            msg.FeeToken,
			FeeTokenAmount:      msg.FeeTokenAmount,
			Data:                msg.Data,
			TokenAmounts:        tokensAndAmounts,
			// TODO: Not needed for plugins, but should be recomputed for consistency.
			// Requires the offramp knowing about onramp version
			Hash: [32]byte{},
		})
	}

	// Unpack will populate with big.Int{false, <allocated empty nat>} for 0 values,
	// which is different from the expected big.NewInt(0). Rebuild to the expected value for this case.
	return ExecReport{
		Messages:          messages,
		OffchainTokenData: erStruct.OffchainTokenData,
		Proofs:            erStruct.Proofs,
		ProofFlagBits:     new(big.Int).SetBytes(erStruct.ProofFlagBits.Bytes()),
	}, nil

}

func (o *OffRampV1_0_0) DecodeExecutionReport(report []byte) (ExecReport, error) {
	return decodeExecReportV1_0_0(o.executionReportArgs, report)
}

func (o *OffRampV1_0_0) TokenEvents() []common.Hash {
	return []common.Hash{abihelpers.MustGetEventID("PoolAdded", abiOffRampV1_0_0), abihelpers.MustGetEventID("PoolRemoved", abiOffRampV1_0_0)}
}

func NewOffRampV1_0_0(lggr logger.Logger, addr common.Address, ec client.Client, lp logpoller.LogPoller, estimator gas.EvmFeeEstimator) (*OffRampV1_0_0, error) {
	offRamp, err := evm_2_evm_offramp_1_0_0.NewEVM2EVMOffRamp(addr, ec)
	if err != nil {
		return nil, err
	}

	executionStateChangedSequenceNumberIndex := 1
	executionReportArgs := abihelpers.MustGetMethodInputs("manuallyExecute", abiOffRampV1_0_0)[:1]
	var filters = []logpoller.Filter{
		{
			Name:      logpoller.FilterName(EXEC_EXECUTION_STATE_CHANGES, addr.String()),
			EventSigs: []common.Hash{ExecutionStateChangedEventV1_0_0},
			Addresses: []common.Address{addr},
		},
		{
			Name:      logpoller.FilterName(EXEC_TOKEN_POOL_ADDED, addr.String()),
			EventSigs: []common.Hash{abihelpers.MustGetEventID("PoolAdded", abiOffRampV1_0_0)},
			Addresses: []common.Address{addr},
		},
		{
			Name:      logpoller.FilterName(EXEC_TOKEN_POOL_REMOVED, addr.String()),
			EventSigs: []common.Hash{abihelpers.MustGetEventID("PoolRemoved", abiOffRampV1_0_0)},
			Addresses: []common.Address{addr},
		},
	}
	if err := logpollerutil.RegisterLpFilters(lp, filters); err != nil {
		return nil, err
	}
	return &OffRampV1_0_0{
		offRamp:             offRamp,
		ec:                  ec,
		addr:                addr,
		lggr:                lggr,
		lp:                  lp,
		filters:             filters,
		estimator:           estimator,
		executionReportArgs: executionReportArgs,
		eventSig:            ExecutionStateChangedEventV1_0_0,
		eventIndex:          executionStateChangedSequenceNumberIndex,
		configMu:            sync.RWMutex{},
		evmBatchCaller: rpclib.NewDynamicLimitedBatchCaller(
			lggr,
			ec,
			rpclib.DefaultRpcBatchSizeLimit,
			rpclib.DefaultRpcBatchBackOffMultiplier,
		),

		// values set on the fly after ChangeConfig is called
		gasPriceEstimator: prices.ExecGasPriceEstimator{},
		offchainConfig:    ExecOffchainConfig{},
		onchainConfig:     ExecOnchainConfig{},
	}, nil
}
