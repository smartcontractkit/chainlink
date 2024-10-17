package evm

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/llo"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

// This is only used for the bootstrap job

var _ commontypes.ConfigProvider = (*lloConfigProvider)(nil)

type lloConfigProvider struct {
	services.Service
	eng *services.Engine

	lp              FilterRegisterer
	cps             []llo.ConfigPollerService
	digester        ocrtypes.OffchainConfigDigester
	runReplay       bool
	replayFromBlock uint64

	ms services.MultiStart
}

func (l *lloConfigProvider) OffchainConfigDigester() ocrtypes.OffchainConfigDigester {
	return l.digester
}
func (l *lloConfigProvider) ContractConfigTracker() ocrtypes.ContractConfigTracker {
	// FIXME: Only return Blue for now. This is a hack to make the bootstrap
	// job work, needs to support multiple config trackers here
	// MERC-5954
	return l.cps[0]
}

func newLLOConfigProvider(
	ctx context.Context,
	lggr logger.Logger,
	chain legacyevm.Chain,
	cc llo.ConfigCache,
	opts *types.RelayOpts,
) (commontypes.ConfigProvider, error) {
	if !common.IsHexAddress(opts.ContractID) {
		return nil, errors.New("invalid contractID, expected hex address")
	}

	configuratorAddress := common.HexToAddress(opts.ContractID)

	relayConfig, err := opts.RelayConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get relay config: %w", err)
	}
	donID := relayConfig.LLODONID
	if donID == 0 {
		return nil, errors.New("donID must be specified in relayConfig for LLO jobs")
	}

	lp := chain.LogPoller()

	cps, digester, err := newLLOConfigPollers(ctx, lggr, cc, lp, chain.Config().EVM().ChainID(), configuratorAddress, relayConfig)
	if err != nil {
		return nil, err
	}

	p := &lloConfigProvider{nil, nil, lp, cps, digester, opts.New, relayConfig.FromBlock, services.MultiStart{}}
	p.Service, p.eng = services.Config{
		Name:  "LLOConfigProvider",
		Start: p.start,
		Close: p.close,
	}.NewServiceEngine(lggr)
	return p, nil
}

func (l *lloConfigProvider) start(ctx context.Context) error {
	if l.runReplay && l.replayFromBlock != 0 {
		// Only replay if it's a brand new job.
		l.eng.Go(func(ctx context.Context) {
			l.eng.Infow("starting replay for config", "fromBlock", l.replayFromBlock)
			// #nosec G115
			if err := l.lp.Replay(ctx, int64(l.replayFromBlock)); err != nil {
				l.eng.Errorw("error replaying for config", "err", err)
			} else {
				l.eng.Infow("completed replaying for config", "replayFromBlock", l.replayFromBlock)
			}
		})
	}
	srvs := []services.StartClose{}
	for _, cp := range l.cps {
		srvs = append(srvs, cp)
	}
	err := l.ms.Start(ctx, srvs...)
	return err
}

func (l *lloConfigProvider) close() error {
	return l.ms.Close()
}
