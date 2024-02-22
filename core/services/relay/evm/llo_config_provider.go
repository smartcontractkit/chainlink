package evm

import (
	"github.com/ethereum/go-ethereum/common"
	pkgerrors "github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/llo"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

func newLLOConfigProvider(lggr logger.Logger, chain legacyevm.Chain, opts *types.RelayOpts) (*configWatcher, error) {
	if !common.IsHexAddress(opts.ContractID) {
		return nil, pkgerrors.Errorf("invalid contractID, expected hex address")
	}

	aggregatorAddress := common.HexToAddress(opts.ContractID)
	configDigester := llo.NewOffchainConfigDigester(chain.Config().EVM().ChainID(), aggregatorAddress)
	return newContractConfigProvider(lggr, chain, opts, aggregatorAddress, ChannelVerifierLogDecoder, configDigester)
}
