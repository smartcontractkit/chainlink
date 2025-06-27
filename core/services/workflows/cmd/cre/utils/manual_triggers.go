package utils

import (
	"context"
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/ethclient"
	evmserver "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/chain-capabilities/evm/server"
	crontrigger "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/triggers/cron/server"
	httptrigger "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/triggers/http/server"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/fakes"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type ManualTriggerCapabilitiesConfig struct {
	Client     *ethclient.Client
	PrivateKey *ecdsa.PrivateKey
}

type ManualTriggers struct {
	ManualCronTrigger     *fakes.ManualCronTriggerService
	ManualHTTPTrigger     *fakes.ManualHTTPTriggerService
	ManualEVMChainTrigger *fakes.FakeEVMChain
}

func NewManualTriggerCapabilities(
	ctx context.Context,
	lggr logger.Logger,
	registry *capabilities.Registry,
	cfg ManualTriggerCapabilitiesConfig,
) (*ManualTriggers, error) {
	// Cron
	manualCronTrigger := fakes.NewManualCronTriggerService(lggr)
	manualCronTriggerServer := crontrigger.NewCronServer(manualCronTrigger)
	if err := registry.Add(ctx, manualCronTriggerServer); err != nil {
		return nil, err
	}

	// HTTP
	manualHTTPTrigger := fakes.NewManualHTTPTriggerService(lggr)
	manualHTTPTriggerServer := httptrigger.NewHTTPServer(manualHTTPTrigger)
	if err := registry.Add(ctx, manualHTTPTriggerServer); err != nil {
		return nil, err
	}

	// EVM
	evm := fakes.NewFakeEvmChain(lggr, cfg.Client, cfg.PrivateKey)
	evmServer := evmserver.NewClientServer(evm)
	if err := registry.Add(ctx, evmServer); err != nil {
		return nil, err
	}

	return &ManualTriggers{
		ManualCronTrigger:     manualCronTrigger,
		ManualHTTPTrigger:     manualHTTPTrigger,
		ManualEVMChainTrigger: evm,
	}, nil
}

func (m *ManualTriggers) Start(ctx context.Context) error {
	err := m.ManualCronTrigger.Start(ctx)
	if err != nil {
		return err
	}

	err = m.ManualHTTPTrigger.Start(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (m *ManualTriggers) Close() error {
	err := m.ManualCronTrigger.Close()
	if err != nil {
		return err
	}

	err = m.ManualHTTPTrigger.Close()
	if err != nil {
		return err
	}

	return nil
}
