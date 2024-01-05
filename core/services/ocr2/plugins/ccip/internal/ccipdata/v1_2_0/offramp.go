package v1_2_0

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_offramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/v1_0_0"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/prices"
)

var (
	abiOffRamp                        = abihelpers.MustParseABI(evm_2_evm_offramp.EVM2EVMOffRampABI)
	_          ccipdata.OffRampReader = &OffRamp{}
)

type ExecOnchainConfig evm_2_evm_offramp.EVM2EVMOffRampDynamicConfig

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
	if d.MaxDataBytes == 0 {
		return errors.New("must set MaxDataBytes")
	}
	if d.MaxPoolReleaseOrMintGas == 0 {
		return errors.New("must set MaxPoolReleaseOrMintGas")
	}
	return nil
}

func (d ExecOnchainConfig) PermissionLessExecutionThresholdDuration() time.Duration {
	return time.Duration(d.PermissionLessExecutionThresholdSeconds) * time.Second
}

// OffRamp In 1.2 we have a different estimator impl
type OffRamp struct {
	*v1_0_0.OffRamp

	offRamp             *evm_2_evm_offramp.EVM2EVMOffRamp
	executionReportArgs abi.Arguments
	lggr                logger.Logger
	ec                  client.Client
	estimator           gas.EvmFeeEstimator

	// Dynamic config
	configMu          sync.RWMutex
	gasPriceEstimator prices.GasPriceEstimatorExec
	offchainConfig    ccipdata.ExecOffchainConfig
	onchainConfig     ccipdata.ExecOnchainConfig
}

func (o *OffRamp) CurrentRateLimiterState(ctx context.Context) (evm_2_evm_offramp.RateLimiterTokenBucket, error) {
	return o.offRamp.CurrentRateLimiterState(&bind.CallOpts{Context: ctx})
}

func (o *OffRamp) ChangeConfig(onchainConfig []byte, offchainConfig []byte) (common.Address, common.Address, error) {
	onchainConfigParsed, err := abihelpers.DecodeAbiStruct[ExecOnchainConfig](onchainConfig)
	if err != nil {
		return common.Address{}, common.Address{}, err
	}

	offchainConfigParsed, err := ccipconfig.DecodeOffchainConfig[v1_0_0.ExecOffchainConfig](offchainConfig)
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
	o.offchainConfig = ccipdata.ExecOffchainConfig{
		DestOptimisticConfirmations: offchainConfigParsed.DestOptimisticConfirmations,
		BatchGasLimit:               offchainConfigParsed.BatchGasLimit,
		RelativeBoostPerWaitHour:    offchainConfigParsed.RelativeBoostPerWaitHour,
		InflightCacheExpiry:         offchainConfigParsed.InflightCacheExpiry,
		RootSnoozeTime:              offchainConfigParsed.RootSnoozeTime,
	}
	o.onchainConfig = ccipdata.ExecOnchainConfig{PermissionLessExecutionThresholdSeconds: time.Second * time.Duration(onchainConfigParsed.PermissionLessExecutionThresholdSeconds)}
	o.gasPriceEstimator = prices.NewDAGasPriceEstimator(o.estimator, big.NewInt(int64(offchainConfigParsed.MaxGasPrice)), 0, 0)
	o.configMu.Unlock()

	o.lggr.Infow("Starting exec plugin",
		"offchainConfig", onchainConfigParsed,
		"onchainConfig", offchainConfigParsed)
	return onchainConfigParsed.PriceRegistry, destWrappedNative, nil
}

func (o *OffRamp) OffchainConfig() ccipdata.ExecOffchainConfig {
	o.configMu.RLock()
	defer o.configMu.RUnlock()
	return o.offchainConfig
}

func (o *OffRamp) OnchainConfig() ccipdata.ExecOnchainConfig {
	o.configMu.RLock()
	defer o.configMu.RUnlock()
	return o.onchainConfig
}

func (o *OffRamp) GasPriceEstimator() prices.GasPriceEstimatorExec {
	o.configMu.RLock()
	defer o.configMu.RUnlock()
	return o.gasPriceEstimator
}

func EncodeExecutionReport(args abi.Arguments, report ccipdata.ExecReport) ([]byte, error) {
	var msgs []evm_2_evm_offramp.InternalEVM2EVMMessage
	for _, msg := range report.Messages {
		var ta []evm_2_evm_offramp.ClientEVMTokenAmount
		for _, tokenAndAmount := range msg.TokenAmounts {
			ta = append(ta, evm_2_evm_offramp.ClientEVMTokenAmount{
				Token:  tokenAndAmount.Token,
				Amount: tokenAndAmount.Amount,
			})
		}
		msgs = append(msgs, evm_2_evm_offramp.InternalEVM2EVMMessage{
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
			SourceTokenData:     msg.SourceTokenData,
			MessageId:           msg.MessageId,
		})
	}

	rep := evm_2_evm_offramp.InternalExecutionReport{
		Messages:          msgs,
		OffchainTokenData: report.OffchainTokenData,
		Proofs:            report.Proofs,
		ProofFlagBits:     report.ProofFlagBits,
	}
	return args.PackValues([]interface{}{&rep})
}

func (o *OffRamp) EncodeExecutionReport(report ccipdata.ExecReport) ([]byte, error) {
	return EncodeExecutionReport(o.executionReportArgs, report)
}

func DecodeExecReport(args abi.Arguments, report []byte) (ccipdata.ExecReport, error) {
	unpacked, err := args.Unpack(report)
	if err != nil {
		return ccipdata.ExecReport{}, err
	}
	if len(unpacked) == 0 {
		return ccipdata.ExecReport{}, errors.New("assumptionViolation: expected at least one element")
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
		return ccipdata.ExecReport{}, fmt.Errorf("got %T", unpacked[0])
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
			SourceTokenData:     msg.SourceTokenData,
			// TODO: Not needed for plugins, but should be recomputed for consistency.
			// Requires the offramp knowing about onramp version
			Hash: [32]byte{},
		})
	}

	// Unpack will populate with big.Int{false, <allocated empty nat>} for 0 values,
	// which is different from the expected big.NewInt(0). Rebuild to the expected value for this case.
	return ccipdata.ExecReport{
		Messages:          messages,
		OffchainTokenData: erStruct.OffchainTokenData,
		Proofs:            erStruct.Proofs,
		ProofFlagBits:     new(big.Int).SetBytes(erStruct.ProofFlagBits.Bytes()),
	}, nil

}

func (o *OffRamp) DecodeExecutionReport(report []byte) (ccipdata.ExecReport, error) {
	return DecodeExecReport(o.executionReportArgs, report)
}

func NewOffRamp(lggr logger.Logger, addr common.Address, ec client.Client, lp logpoller.LogPoller, estimator gas.EvmFeeEstimator) (*OffRamp, error) {
	v100, err := v1_0_0.NewOffRamp(lggr, addr, ec, lp, estimator)
	if err != nil {
		return nil, err
	}

	offRamp, err := evm_2_evm_offramp.NewEVM2EVMOffRamp(addr, ec)
	if err != nil {
		return nil, err
	}

	executionReportArgs := abihelpers.MustGetMethodInputs("manuallyExecute", abiOffRamp)[:1]

	return &OffRamp{
		OffRamp:             v100,
		offRamp:             offRamp,
		executionReportArgs: executionReportArgs,
		lggr:                lggr,
		ec:                  ec,
		estimator:           estimator,

		configMu: sync.RWMutex{},

		// values set on the fly after ChangeConfig is called
		gasPriceEstimator: prices.ExecGasPriceEstimator{},
		offchainConfig:    ccipdata.ExecOffchainConfig{},
		onchainConfig:     ccipdata.ExecOnchainConfig{},
	}, nil
}
