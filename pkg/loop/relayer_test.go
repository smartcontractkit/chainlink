package loop_test

import (
	"context"
	"fmt"
	"math/big"
	"reflect"

	"github.com/smartcontractkit/chainlink-relay/pkg/types"
)

type staticRelayer struct{}

func (s staticRelayer) Start(ctx context.Context) error { return nil }

func (s staticRelayer) Close() error { return nil }

func (s staticRelayer) Ready() error { panic("unimplemented") }

func (s staticRelayer) Name() string { panic("unimplemented") }

func (s staticRelayer) HealthReport() map[string]error { panic("unimplemented") }

func (s staticRelayer) NewConfigProvider(ctx context.Context, r types.RelayArgs) (types.ConfigProvider, error) {
	if !equalRelayArgs(r, rargs) {
		return nil, fmt.Errorf("expected relay args:\n\t%v\nbut got:\n\t%v", rargs, r)
	}
	return staticConfigProvider{}, nil
}

func (s staticRelayer) NewMedianProvider(ctx context.Context, r types.RelayArgs, p types.PluginArgs) (types.MedianProvider, error) {
	if !equalRelayArgs(r, rargs) {
		return nil, fmt.Errorf("expected relay args:\n\t%v\nbut got:\n\t%v", rargs, r)
	}
	if !reflect.DeepEqual(pargs, p) {
		return nil, fmt.Errorf("expected plugin args %v but got %v", pargs, p)
	}
	return staticMedianProvider{}, nil
}

func (s staticRelayer) NewMercuryProvider(ctx context.Context, rargs types.RelayArgs, pargs types.PluginArgs) (types.MercuryProvider, error) {
	panic("unimplemented")
}

func (s staticRelayer) ChainStatus(ctx context.Context, id string) (types.ChainStatus, error) {
	if id != chainID {
		return types.ChainStatus{}, fmt.Errorf("expected id %s but got %s", chainID, id)
	}
	return chain, nil
}

func (s staticRelayer) ChainStatuses(ctx context.Context, o, l int) ([]types.ChainStatus, int, error) {
	if offset != o {
		return nil, -1, fmt.Errorf("expected offset %d but got %d", offset, o)
	}
	if limit != l {
		return nil, -1, fmt.Errorf("expected limit %d but got %d", limit, l)
	}
	return chains, count, nil
}

func (s staticRelayer) NodeStatuses(ctx context.Context, o, l int, cs ...string) ([]types.NodeStatus, int, error) {
	if offset != o {
		return nil, -1, fmt.Errorf("expected offset %d but got %d", offset, o)
	}
	if limit != l {
		return nil, -1, fmt.Errorf("expected limit %d but got %d", limit, l)
	}
	if !reflect.DeepEqual(chainIDs, cs) {
		return nil, -1, fmt.Errorf("expected chain IDs %v but got %v", chainIDs, cs)
	}
	return nodes, count, nil
}

func (s staticRelayer) SendTx(ctx context.Context, id, f, t string, a *big.Int, b bool) error {
	if id != chainID {
		return fmt.Errorf("expected id %s but got %s", chainID, id)
	}
	if f != from {
		return fmt.Errorf("expected from %s but got %s", from, f)
	}
	if t != to {
		return fmt.Errorf("expected to %s but got %s", to, t)
	}
	if amount.Cmp(a) != 0 {
		return fmt.Errorf("expected amount %s but got %s", amount, a)
	}
	if b != balanceCheck {
		return fmt.Errorf("expected balance check %t but got %t", balanceCheck, b)
	}
	return nil
}
