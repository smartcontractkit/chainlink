package loop

import (
	"context"
	"fmt"
	"math/big"
	"os/exec"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/goplugin"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
)

var _ Relayer = (*RelayerService)(nil)

// RelayerService is a [types.Service] that maintains an internal [Relayer].
type RelayerService struct {
	goplugin.PluginService[*GRPCPluginRelayer, Relayer]
}

// NewRelayerService returns a new [*RelayerService].
// cmd must return a new exec.Cmd each time it is called.
func NewRelayerService(lggr logger.Logger, grpcOpts GRPCOpts, cmd func() *exec.Cmd, config string, keystore core.Keystore) *RelayerService {
	newService := func(ctx context.Context, instance any) (Relayer, error) {
		plug, ok := instance.(PluginRelayer)
		if !ok {
			return nil, fmt.Errorf("expected PluginRelayer but got %T", instance)
		}
		r, err := plug.NewRelayer(ctx, config, keystore)
		if err != nil {
			return nil, fmt.Errorf("failed to create Relayer: %w", err)
		}
		return r, nil
	}
	stopCh := make(chan struct{})
	lggr = logger.Named(lggr, "RelayerService")
	var rs RelayerService
	broker := BrokerConfig{StopCh: stopCh, Logger: lggr, GRPCOpts: grpcOpts}
	rs.Init(PluginRelayerName, &GRPCPluginRelayer{BrokerConfig: broker}, newService, lggr, cmd, stopCh)
	return &rs
}

func (r *RelayerService) NewContractReader(ctx context.Context, contractReaderConfig []byte) (types.ContractReader, error) {
	if err := r.WaitCtx(ctx); err != nil {
		return nil, err
	}
	return r.Service.NewContractReader(ctx, contractReaderConfig)
}

func (r *RelayerService) NewConfigProvider(ctx context.Context, args types.RelayArgs) (types.ConfigProvider, error) {
	if err := r.WaitCtx(ctx); err != nil {
		return nil, err
	}
	return r.Service.NewConfigProvider(ctx, args)
}

func (r *RelayerService) NewPluginProvider(ctx context.Context, rargs types.RelayArgs, pargs types.PluginArgs) (types.PluginProvider, error) {
	if err := r.WaitCtx(ctx); err != nil {
		return nil, err
	}
	return r.Service.NewPluginProvider(ctx, rargs, pargs)
}

func (r *RelayerService) NewLLOProvider(ctx context.Context, rargs types.RelayArgs, pargs types.PluginArgs) (types.LLOProvider, error) {
	if err := r.WaitCtx(ctx); err != nil {
		return nil, err
	}
	return r.Service.NewLLOProvider(ctx, rargs, pargs)
}

func (r *RelayerService) GetChainStatus(ctx context.Context) (types.ChainStatus, error) {
	if err := r.WaitCtx(ctx); err != nil {
		return types.ChainStatus{}, err
	}
	return r.Service.GetChainStatus(ctx)
}

func (r *RelayerService) ListNodeStatuses(ctx context.Context, pageSize int32, pageToken string) (nodes []types.NodeStatus, nextPageToken string, total int, err error) {
	if err := r.WaitCtx(ctx); err != nil {
		return nil, "", -1, err
	}
	return r.Service.ListNodeStatuses(ctx, pageSize, pageToken)
}

func (r *RelayerService) Transact(ctx context.Context, from, to string, amount *big.Int, balanceCheck bool) error {
	if err := r.WaitCtx(ctx); err != nil {
		return err
	}
	return r.Service.Transact(ctx, from, to, amount, balanceCheck)
}
