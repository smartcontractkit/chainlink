package test

import (
	"context"

	libocr "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

type StaticPluginProvider struct{}

func (s StaticPluginProvider) Start(ctx context.Context) error { return nil }

func (s StaticPluginProvider) Close() error { return nil }

func (s StaticPluginProvider) Ready() error { panic("unimplemented") }

func (s StaticPluginProvider) Name() string { panic("unimplemented") }

func (s StaticPluginProvider) HealthReport() map[string]error { panic("unimplemented") }

func (s StaticPluginProvider) OffchainConfigDigester() libocr.OffchainConfigDigester {
	return staticOffchainConfigDigester{}
}

func (s StaticPluginProvider) ContractConfigTracker() libocr.ContractConfigTracker {
	return staticContractConfigTracker{}
}

func (s StaticPluginProvider) ContractTransmitter() libocr.ContractTransmitter {
	return staticContractTransmitter{}
}

func (s StaticPluginProvider) ChainReader() types.ChainReader {
	return staticChainReader{}
}
