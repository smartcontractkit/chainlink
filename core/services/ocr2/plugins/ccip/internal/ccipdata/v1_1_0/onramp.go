package v1_1_0

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_onramp_1_1_0"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/v1_0_0"
)

var _ ccipdata.OnRampReader = &OnRamp{}

// OnRamp The only difference that the plugins care about in 1.1 is that the dynamic config struct has changed.
type OnRamp struct {
	*v1_0_0.OnRamp
	onRamp *evm_2_evm_onramp_1_1_0.EVM2EVMOnRamp
}

func NewOnRamp(lggr logger.Logger, sourceSelector, destSelector uint64, onRampAddress common.Address, sourceLP logpoller.LogPoller, source client.Client) (*OnRamp, error) {
	onRamp, err := evm_2_evm_onramp_1_1_0.NewEVM2EVMOnRamp(onRampAddress, source)
	if err != nil {
		return nil, err
	}
	onRamp100, err := v1_0_0.NewOnRamp(lggr, sourceSelector, destSelector, onRampAddress, sourceLP, source)
	if err != nil {
		return nil, err
	}
	return &OnRamp{
		OnRamp: onRamp100,
		onRamp: onRamp,
	}, nil
}

func (o *OnRamp) RouterAddress() (common.Address, error) {
	config, err := o.onRamp.GetDynamicConfig(nil)
	if err != nil {
		return common.Address{}, err
	}
	return config.Router, nil
}

func (o *OnRamp) GetDynamicConfig() (ccipdata.OnRampDynamicConfig, error) {
	if o.onRamp == nil {
		return ccipdata.OnRampDynamicConfig{}, fmt.Errorf("onramp not initialized")
	}
	legacyDynamicConfig, err := o.onRamp.GetDynamicConfig(nil)
	if err != nil {
		return ccipdata.OnRampDynamicConfig{}, err
	}
	return ccipdata.OnRampDynamicConfig{
		Router:                            legacyDynamicConfig.Router,
		MaxNumberOfTokensPerMsg:           legacyDynamicConfig.MaxTokensLength,
		DestGasOverhead:                   legacyDynamicConfig.DestGasOverhead,
		DestGasPerPayloadByte:             legacyDynamicConfig.DestGasPerPayloadByte,
		DestDataAvailabilityOverheadGas:   0,
		DestGasPerDataAvailabilityByte:    0,
		DestDataAvailabilityMultiplierBps: 0,
		PriceRegistry:                     legacyDynamicConfig.PriceRegistry,
		MaxDataBytes:                      legacyDynamicConfig.MaxDataSize,
		MaxPerMsgGasLimit:                 uint32(legacyDynamicConfig.MaxGasLimit),
	}, nil
}
