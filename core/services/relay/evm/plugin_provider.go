package evm

import (
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/services"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

type PluginProvider struct {
	services.ServiceCtx
	chainReader            types.ChainReader
	codec                  types.Codec
	transmitter            ocrtypes.ContractTransmitter
	contractConfigTracker  ocrtypes.ContractConfigTracker
	offchainConfigDigester ocrtypes.OffchainConfigDigester
}

var _ types.PluginProvider = (*PluginProvider)(nil)

func NewPluginProvider(
	chainReader types.ChainReader,
	codec types.Codec,
	transmitter ocrtypes.ContractTransmitter,
	contractConfigTracker ocrtypes.ContractConfigTracker,
	offchainConfigDigester ocrtypes.OffchainConfigDigester,
) *PluginProvider {
	return &PluginProvider{
		chainReader:            chainReader,
		codec:                  codec,
		transmitter:            transmitter,
		contractConfigTracker:  contractConfigTracker,
		offchainConfigDigester: offchainConfigDigester,
	}

}

func (p *PluginProvider) ContractTransmitter() ocrtypes.ContractTransmitter {
	return p.transmitter
}

func (p *PluginProvider) OffchainConfigDigester() ocrtypes.OffchainConfigDigester {
	return p.offchainConfigDigester
}

func (p *PluginProvider) ContractConfigTracker() ocrtypes.ContractConfigTracker {
	return p.contractConfigTracker
}

func (p *PluginProvider) ChainReader() types.ChainReader {
	return p.chainReader
}

func (p *PluginProvider) Codec() types.Codec {
	return p.codec
}
