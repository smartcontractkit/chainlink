package observability

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_offramp"
)

type ObservedEVM2EVMOfframp struct {
	evm_2_evm_offramp.EVM2EVMOffRampInterface
	metric metricDetails
}

func NewObservedEvm2EvmOffRamp(address common.Address, pluginName string, client client.Client) (*ObservedEVM2EVMOfframp, error) {
	offRamp, err := evm_2_evm_offramp.NewEVM2EVMOffRamp(address, client)
	if err != nil {
		return nil, err
	}
	return &ObservedEVM2EVMOfframp{
		EVM2EVMOffRampInterface: offRamp,
		metric: metricDetails{
			histogram:  evm2evmOffRampHistogram,
			pluginName: pluginName,
			chainId:    client.ConfiguredChainID(),
		},
	}, nil
}

func (o *ObservedEVM2EVMOfframp) GetSupportedTokens(opts *bind.CallOpts) ([]common.Address, error) {
	return withObservedContract(o.metric, "GetSupportedTokens", func() ([]common.Address, error) {
		return o.EVM2EVMOffRampInterface.GetSupportedTokens(opts)
	})
}

func (o *ObservedEVM2EVMOfframp) GetDestinationTokens(opts *bind.CallOpts) ([]common.Address, error) {
	return withObservedContract(o.metric, "GetDestinationTokens", func() ([]common.Address, error) {
		return o.EVM2EVMOffRampInterface.GetDestinationTokens(opts)
	})
}

func (o *ObservedEVM2EVMOfframp) GetDestinationToken(opts *bind.CallOpts, sourceToken common.Address) (common.Address, error) {
	return withObservedContract(o.metric, "GetDestinationToken", func() (common.Address, error) {
		return o.EVM2EVMOffRampInterface.GetDestinationToken(opts, sourceToken)
	})
}

func (o *ObservedEVM2EVMOfframp) CurrentRateLimiterState(opts *bind.CallOpts) (evm_2_evm_offramp.RateLimiterTokenBucket, error) {
	return withObservedContract(o.metric, "CurrentRateLimiterState", func() (evm_2_evm_offramp.RateLimiterTokenBucket, error) {
		return o.EVM2EVMOffRampInterface.CurrentRateLimiterState(opts)
	})
}

func (o *ObservedEVM2EVMOfframp) GetPoolByDestToken(opts *bind.CallOpts, destToken common.Address) (common.Address, error) {
	return withObservedContract(o.metric, "GetPoolByDestToken", func() (common.Address, error) {
		return o.EVM2EVMOffRampInterface.GetPoolByDestToken(opts, destToken)
	})
}
