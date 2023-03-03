package chainlink

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"

	starkchain "github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/chain"
	"github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/ocr2"

	"github.com/smartcontractkit/chainlink-relay/pkg/logger"
	relaytypes "github.com/smartcontractkit/chainlink-relay/pkg/types"
)

var _ relaytypes.Relayer = (*relayer)(nil)

type relayer struct {
	chainSet starkchain.ChainSet
	ctx      context.Context

	lggr logger.Logger

	cancel func()
}

func NewRelayer(lggr logger.Logger, chainSet starkchain.ChainSet) *relayer {
	ctx, cancel := context.WithCancel(context.Background())
	return &relayer{
		chainSet: chainSet,
		ctx:      ctx,
		lggr:     lggr,
		cancel:   cancel,
	}
}

func (r *relayer) Name() string {
	return r.lggr.Name()
}

func (r *relayer) Start(context.Context) error {
	return nil
}

func (r *relayer) Close() error {
	r.cancel()
	return nil
}

func (r *relayer) Ready() error {
	return r.chainSet.Ready()
}

func (r *relayer) Healthy() error {
	return r.chainSet.Healthy()
}

func (r *relayer) HealthReport() map[string]error {
	return map[string]error{r.Name(): r.Healthy()}
}

func (r *relayer) NewConfigProvider(args relaytypes.RelayArgs) (relaytypes.ConfigProvider, error) {
	var relayConfig RelayConfig

	err := json.Unmarshal(args.RelayConfig, &relayConfig)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't unmarshal RelayConfig")
	}

	chain, err := r.chainSet.Chain(r.ctx, relayConfig.ChainID)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't initilize Chain")
	}

	reader, err := chain.Reader()
	if err != nil {
		return nil, errors.Wrap(err, "error in NewConfigProvider chain.Reader")
	}
	configProvider, err := ocr2.NewConfigProvider(relayConfig.ChainID, args.ContractID, reader, chain.Config(), r.lggr)
	if err != nil {
		return nil, errors.Wrap(err, "coudln't initialize ConfigProvider")
	}

	return configProvider, nil
}

func (r *relayer) NewMedianProvider(rargs relaytypes.RelayArgs, pargs relaytypes.PluginArgs) (relaytypes.MedianProvider, error) {
	var relayConfig RelayConfig

	err := json.Unmarshal(rargs.RelayConfig, &relayConfig)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't unmarshal RelayConfig")
	}

	chain, err := r.chainSet.Chain(r.ctx, relayConfig.ChainID)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't initilize Chain")
	}

	// todo: use pargs for median provider
	reader, err := chain.Reader()
	if err != nil {
		return nil, errors.Wrap(err, "error in NewMedianProvider chain.Reader")
	}
	medianProvider, err := ocr2.NewMedianProvider(relayConfig.ChainID, rargs.ContractID, pargs.TransmitterID, reader, chain.Config(), chain.TxManager(), r.lggr)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't initilize MedianProvider")
	}

	return medianProvider, nil
}

func (r *relayer) NewMercuryProvider(rargs relaytypes.RelayArgs, pargs relaytypes.PluginArgs) (relaytypes.MercuryProvider, error) {
	return nil, errors.New("mercury is not supported for starknet")
}
