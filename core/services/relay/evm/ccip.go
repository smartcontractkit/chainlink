package evm

import (
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"

	relaytypes "github.com/smartcontractkit/chainlink-relay/pkg/types"
)

// CCIPCommitProvider provides all components needed for a CCIP Relay OCR2 plugin.
type CCIPCommitProvider interface {
	relaytypes.Plugin
}

// CCIPExecutionProvider provides all components needed for a CCIP Execution OCR2 plugin.
type CCIPExecutionProvider interface {
	relaytypes.Plugin
}

type ccipCommitProvider struct {
	*configWatcher
	contractTransmitter *contractTransmitter
}

func NewCCIPCommitProvider(lggr logger.Logger, chainSet evm.Chain, rargs relaytypes.RelayArgs, transmitterID string, ks keystore.Eth, eventBroadcaster pg.EventBroadcaster) (CCIPCommitProvider, error) {
	relayOpts := types.NewRelayOpts(rargs)
	configWatcher, err := newConfigProvider(lggr, chainSet, relayOpts, eventBroadcaster)
	if err != nil {
		return nil, err
	}
	contractTransmitter, err := newContractTransmitter(lggr, rargs, transmitterID, configWatcher, ks, ccip.CommitReportToEthTxMeta)
	if err != nil {
		return nil, err
	}
	return &ccipCommitProvider{
		configWatcher:       configWatcher,
		contractTransmitter: contractTransmitter,
	}, nil
}

func (c *ccipCommitProvider) ContractTransmitter() ocrtypes.ContractTransmitter {
	return c.contractTransmitter
}

type ccipExecutionProvider struct {
	*configWatcher
	contractTransmitter *contractTransmitter
}

var _ relaytypes.Plugin = (*ccipExecutionProvider)(nil)

func NewCCIPExecutionProvider(lggr logger.Logger, chainSet evm.Chain, rargs relaytypes.RelayArgs, transmitterID string, ks keystore.Eth, eventBroadcaster pg.EventBroadcaster) (CCIPExecutionProvider, error) {
	relayOpts := types.NewRelayOpts(rargs)

	configWatcher, err := newConfigProvider(lggr, chainSet, relayOpts, eventBroadcaster)
	if err != nil {
		return nil, err
	}
	contractTransmitter, err := newContractTransmitter(lggr, rargs, transmitterID, configWatcher, ks, ccip.ExecutionReportToEthTxMeta)
	if err != nil {
		return nil, err
	}
	return &ccipExecutionProvider{
		configWatcher:       configWatcher,
		contractTransmitter: contractTransmitter,
	}, nil
}

func (c *ccipExecutionProvider) ContractTransmitter() ocrtypes.ContractTransmitter {
	return c.contractTransmitter
}
