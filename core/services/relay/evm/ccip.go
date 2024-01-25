package evm

import (
	"github.com/ethereum/go-ethereum/common"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/ccipcommit"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/ccipexec"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

// CCIPCommitProvider provides all components needed for a CCIP Relay OCR2 plugin.
type CCIPCommitProvider interface {
	commontypes.Plugin
}

// CCIPExecutionProvider provides all components needed for a CCIP Execution OCR2 plugin.
type CCIPExecutionProvider interface {
	commontypes.Plugin
}

type ccipCommitProvider struct {
	*configWatcher
	contractTransmitter *contractTransmitter
}

func NewCCIPCommitProvider(lggr logger.Logger, chainSet legacyevm.Chain, rargs commontypes.RelayArgs, transmitterID string, ks keystore.Eth) (CCIPCommitProvider, error) {
	relayOpts := types.NewRelayOpts(rargs)
	configWatcher, err := newConfigProvider(lggr, chainSet, relayOpts)
	if err != nil {
		return nil, err
	}
	address := common.HexToAddress(relayOpts.ContractID)
	typ, ver, err := ccipconfig.TypeAndVersion(address, chainSet.Client())
	if err != nil {
		return nil, err
	}
	fn, err := ccipcommit.CommitReportToEthTxMeta(typ, ver)
	if err != nil {
		return nil, err
	}
	contractTransmitter, err := newContractTransmitter(lggr, rargs, transmitterID, ks, configWatcher, configTransmitterOpts{}, fn)
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

func (c *ccipCommitProvider) ChainReader() commontypes.ChainReader {
	return nil
}

func (c *ccipCommitProvider) Codec() commontypes.Codec {
	return nil
}

type ccipExecutionProvider struct {
	*configWatcher
	contractTransmitter *contractTransmitter
}

var _ commontypes.Plugin = (*ccipExecutionProvider)(nil)

func NewCCIPExecutionProvider(lggr logger.Logger, chainSet legacyevm.Chain, rargs commontypes.RelayArgs, transmitterID string, ks keystore.Eth) (CCIPExecutionProvider, error) {
	relayOpts := types.NewRelayOpts(rargs)

	configWatcher, err := newConfigProvider(lggr, chainSet, relayOpts)
	if err != nil {
		return nil, err
	}
	address := common.HexToAddress(relayOpts.ContractID)
	typ, ver, err := ccipconfig.TypeAndVersion(address, chainSet.Client())
	if err != nil {
		return nil, err
	}
	fn, err := ccipexec.ExecReportToEthTxMeta(typ, ver)
	if err != nil {
		return nil, err
	}
	contractTransmitter, err := newContractTransmitter(lggr, rargs, transmitterID, ks, configWatcher, configTransmitterOpts{}, fn)
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

func (c *ccipExecutionProvider) ChainReader() commontypes.ChainReader {
	return nil
}

func (c *ccipExecutionProvider) Codec() commontypes.Codec {
	return nil
}
