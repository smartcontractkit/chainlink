package relay

import (
	"context"

	"github.com/smartcontractkit/chainlink/core/logger"

	"github.com/smartcontractkit/chainlink-relay/pkg/loop"
	"github.com/smartcontractkit/chainlink-relay/pkg/types"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"
)

type Network string

var (
	EVM             Network = "evm"
	Solana          Network = "solana"
	StarkNet        Network = "starknet"
	SupportedRelays         = map[Network]struct{}{
		EVM:      {},
		Solana:   {},
		StarkNet: {},
	}
)

var _ loop.ChainRelayer = (*chainRelayer)(nil)

// chainRelayer adapts a [types.Relayer] to [loop.ChainRelayer].
type chainRelayer struct {
	types.Relayer
	lggr logger.Logger
}

func NewChainRelayer(r types.Relayer, lggr logger.Logger) loop.ChainRelayer {
	return &chainRelayer{Relayer: r, lggr: lggr.Named("Relayer")}
}

func (c *chainRelayer) NewConfigProvider(ctx context.Context, rargs types.RelayArgs) (types.ConfigProvider, error) {
	return c.Relayer.NewConfigProvider(rargs)
}

func (c *chainRelayer) NewMedianPluginProvider(ctx context.Context, rargs types.RelayArgs, pargs types.PluginArgs, dataSource, juelsPerFeeCoin median.DataSource) (loop.ReportingPluginProvider, error) {
	provider, err := c.NewMedianProvider(rargs, pargs)
	if err != nil {
		return nil, err
	}
	lggr := c.lggr //TODO with key/vals from ctx! https://smartcontract-it.atlassian.net/browse/BCF-2113
	factory := median.NumericalMedianFactory{
		ContractTransmitter:       provider.MedianContract(),
		DataSource:                dataSource,
		JuelsPerFeeCoinDataSource: juelsPerFeeCoin,
		Logger: logger.NewOCRWrapper(lggr, true, func(msg string) {
			//TODO grpc service to save job errors: RecordError(msg) // implicit job ID
			// lggr.ErrorIf(d.jobORM.RecordError(jb.ID, msg), "unable to record error")
			// https://smartcontract-it.atlassian.net/browse/BCF-2115
			lggr.Error(msg)
		}),
		OnchainConfigCodec: provider.OnchainConfigCodec(),
		ReportCodec:        provider.ReportCodec(),
	}
	return ReportingPluginProvider{provider, factory}, nil
}

type ReportingPluginProvider struct {
	types.Plugin
	ocrtypes.ReportingPluginFactory
}
