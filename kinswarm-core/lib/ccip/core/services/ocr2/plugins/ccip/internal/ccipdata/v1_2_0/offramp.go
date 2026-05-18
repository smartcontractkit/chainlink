package v1_2_0

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/config"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_offramp_1_2_0"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcalc"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/v1_0_0"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/prices"
)

var (
	abiOffRamp                        = abihelpers.MustParseABI(evm_2_evm_offramp_1_2_0.EVM2EVMOffRampABI)
	_          ccipdata.OffRampReader = &OffRamp{}
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

// OffRamp In 1.2 we have a different estimator impl
type OffRamp struct {
	*v1_0_0.OffRamp
	offRampV120        evm_2_evm_offramp_1_2_0.EVM2EVMOffRampInterface
	feeEstimatorConfig ccipdata.FeeEstimatorConfigReader
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

func (o *OffRamp) GetRouter(ctx context.Context) (cciptypes.Address, error) {
	dynamicConfig, err := o.offRampV120.GetDynamicConfig(&bind.CallOpts{Context: ctx})
	if err != nil {
		return "", err
	}
	return ccipcalc.EvmAddrToGeneric(dynamicConfig.Router), nil
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

func NewOffRamp(
	lggr logger.Logger,
	addr common.Address,
	ec client.Client,
	lp logpoller.LogPoller,
	estimator gas.EvmFeeEstimator,
	destMaxGasPrice *big.Int,
	feeEstimatorConfig ccipdata.FeeEstimatorConfigReader,
) (*OffRamp, error) {
	v100, err := v1_0_0.NewOffRamp(lggr, addr, ec, lp, estimator, destMaxGasPrice, feeEstimatorConfig)
	if err != nil {
		return nil, err
	}

	offRamp, err := evm_2_evm_offramp_1_2_0.NewEVM2EVMOffRamp(addr, ec)
	if err != nil {
		return nil, err
	}

	v100.ExecutionReportArgs = abihelpers.MustGetMethodInputs("manuallyExecute", abiOffRamp)[:1]

	return &OffRamp{
		OffRamp:            v100,
		offRampV120:        offRamp,
		feeEstimatorConfig: feeEstimatorConfig,
	}, nil
}
