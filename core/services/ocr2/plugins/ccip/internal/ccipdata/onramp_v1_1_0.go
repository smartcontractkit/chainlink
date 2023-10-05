package ccipdata

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_onramp_1_1_0"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

var _ OnRampReader = &OnRampV1_1_0{}

// OnRampV1_1_0 The only difference that the plugins care about in 1.1 is that the dynamic config struct has changed.
type OnRampV1_1_0 struct {
	*OnRampV1_0_0
	onRamp *evm_2_evm_onramp_1_1_0.EVM2EVMOnRamp
}

func NewOnRampV1_1_0(lggr logger.Logger, sourceSelector, destSelector uint64, onRampAddress common.Address, sourceLP logpoller.LogPoller, source client.Client, finalityTags bool) (*OnRampV1_1_0, error) {
	onRamp, err := evm_2_evm_onramp_1_1_0.NewEVM2EVMOnRamp(onRampAddress, source)
	if err != nil {
		return nil, err
	}
	onRamp100, err := NewOnRampV1_0_0(lggr, sourceSelector, destSelector, onRampAddress, sourceLP, source, finalityTags)
	if err != nil {
		return nil, err
	}
	return &OnRampV1_1_0{
		OnRampV1_0_0: onRamp100,
		onRamp:       onRamp,
	}, nil
}

func (o *OnRampV1_1_0) RouterAddress() (common.Address, error) {
	config, err := o.onRamp.GetDynamicConfig(nil)
	if err != nil {
		return common.Address{}, err
	}
	return config.Router, nil
}
