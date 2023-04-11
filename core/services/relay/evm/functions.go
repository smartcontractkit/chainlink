package evm

import (
	"github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"

	relaytypes "github.com/smartcontractkit/chainlink-relay/pkg/types"
)

type functionsProvider struct {
	*configWatcher
	contractTransmitter ContractTransmitter
}

var (
	_ relaytypes.Plugin = (*functionsProvider)(nil)
)

func (p *functionsProvider) ContractTransmitter() types.ContractTransmitter {
	return p.contractTransmitter
}

func NewFunctionsProvider(chainSet evm.ChainSet, rargs relaytypes.RelayArgs, pargs relaytypes.PluginArgs, lggr logger.Logger, ethKeystore keystore.Eth) (relaytypes.Plugin, error) {
	configWatcher, err := newConfigProvider(lggr, chainSet, rargs)
	if err != nil {
		return nil, err
	}
	contractTransmitter, err := newContractTransmitter(lggr, rargs, pargs.TransmitterID, configWatcher, ethKeystore)
	if err != nil {
		return nil, err
	}
	return &functionsProvider{
		configWatcher:       configWatcher,
		contractTransmitter: contractTransmitter,
	}, nil
}
