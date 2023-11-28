package testutils

import (
	"context"
	"math/big"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
)

type MockRelayer struct {
	ChainStatus commontypes.ChainStatus
}

func (m MockRelayer) Name() string {
	panic("not implemented")
}

func (m MockRelayer) Start(ctx context.Context) error {
	panic("not implemented")
}

func (m MockRelayer) Close() error {
	panic("not implemented")
}

func (m MockRelayer) Ready() error {
	panic("not implemented")
}

func (m MockRelayer) HealthReport() map[string]error {
	panic("not implemented")
}

func (m MockRelayer) GetChainStatus(ctx context.Context) (commontypes.ChainStatus, error) {
	return m.ChainStatus, nil
}

func (m MockRelayer) ListNodeStatuses(ctx context.Context, pageSize int32, pageToken string) (stats []commontypes.NodeStatus, nextPageToken string, total int, err error) {
	panic("not implemented")
}

func (m MockRelayer) Transact(ctx context.Context, from, to string, amount *big.Int, balanceCheck bool) error {
	panic("not implemented")
}

func (m MockRelayer) NewConfigProvider(ctx context.Context, args commontypes.RelayArgs) (commontypes.ConfigProvider, error) {
	panic("not implemented")
}

func (m MockRelayer) NewPluginProvider(ctx context.Context, args commontypes.RelayArgs, args2 commontypes.PluginArgs) (commontypes.PluginProvider, error) {
	panic("not implemented")
}
