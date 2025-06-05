package dummy

import (
	"context"
	"encoding/json"
	"math/big"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/services/llo/bm"
	"github.com/smartcontractkit/chainlink/v2/core/services/llo/channeldefinitions"
	"github.com/smartcontractkit/chainlink/v2/core/services/llo/retirement"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/llo/config"
)

// The dummy relayer is a simple reference implementation that doesn't actually
// connect to any chain. It's useful for testing and as a reference for
// implementing a new relayer.

type RelayConfig struct {
	ConfigTracker ConfigTrackerCfg
}

type relayer struct {
	lggr    logger.Logger
	chainID string
}

func NewRelayer(lggr logger.Logger, chainID string) loop.Relayer {
	return &relayer{lggr, chainID}
}

func (r *relayer) NewContractWriter(ctx context.Context, chainWriterConfig []byte) (types.ContractWriter, error) {
	return nil, nil
}

func (r *relayer) NewContractReader(ctx context.Context, contractReaderConfig []byte) (types.ContractReader, error) {
	return nil, nil
}

func (r *relayer) EVM() (types.EVMService, error) {
	return nil, nil
}

func (r *relayer) NewConfigProvider(_ context.Context, rargs types.RelayArgs) (types.ConfigProvider, error) {
	var cfg RelayConfig
	if err := json.Unmarshal(rargs.RelayConfig, &cfg); err != nil {
		return nil, err
	}
	return NewConfigProvider(r.lggr, cfg)
}
func (r *relayer) NewPluginProvider(context.Context, types.RelayArgs, types.PluginArgs) (types.PluginProvider, error) {
	return nil, nil
}
func (r *relayer) NewLLOProvider(ctx context.Context, rargs types.RelayArgs, pargs types.PluginArgs) (types.LLOProvider, error) {
	cp, err := r.NewConfigProvider(ctx, rargs)
	if err != nil {
		return nil, err
	}
	transmitter := bm.NewTransmitter(r.lggr, pargs.TransmitterID)
	pluginCfg := new(config.PluginConfig)
	if err = pluginCfg.Unmarshal(pargs.PluginConfig); err != nil {
		return nil, err
	}
	cdc, err := channeldefinitions.NewStaticChannelDefinitionCache(r.lggr, pluginCfg.ChannelDefinitions)
	if err != nil {
		return nil, err
	}
	src := retirement.NewNeverShouldRetireCache()
	return NewLLOProvider(r.lggr, cp, transmitter, cdc, src), nil
}
func (r *relayer) LatestHead(_ context.Context) (types.Head, error) {
	return types.Head{}, nil
}
func (r *relayer) GetChainStatus(ctx context.Context) (types.ChainStatus, error) {
	return types.ChainStatus{}, nil
}
func (r *relayer) ListNodeStatuses(ctx context.Context, pageSize int32, pageToken string) (stats []types.NodeStatus, nextPageToken string, total int, err error) {
	return nil, "", 0, nil
}
func (r *relayer) Transact(ctx context.Context, from, to string, amount *big.Int, balanceCheck bool) error {
	return nil
}
func (r *relayer) Replay(ctx context.Context, fromBlock string, args map[string]any) error {
	return nil
}
func (r *relayer) Name() string                { return r.lggr.Name() }
func (r *relayer) Start(context.Context) error { return nil }
func (r *relayer) Close() error                { return nil }
func (r *relayer) Ready() error                { return nil }
func (r *relayer) HealthReport() map[string]error {
	return map[string]error{r.lggr.Name(): nil}
}
