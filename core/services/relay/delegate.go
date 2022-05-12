package relay

import (
	"context"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana"
	"github.com/smartcontractkit/chainlink-terra/pkg/terra"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/relay/evm"
	"github.com/smartcontractkit/chainlink/core/services/relay/types"
)

var (
	SupportedRelayers = map[types.Network]struct{}{
		types.EVM:    {},
		types.Solana: {},
		types.Terra:  {},
	}
	_ types.Relayer = &evm.Relayer{}
	_ types.Relayer = &solana.Relayer{}
	_ types.Relayer = &terra.Relayer{}
)

type Delegate struct {
	relayers map[types.Network]types.Relayer
	ks       keystore.Master
}

// NewDelegate creates a master relay Delegate which manages "relays" which are OCR2 median reporting plugins
// for various chains (evm and non-evm). nil Relayers will be disabled.
func NewDelegate(ks keystore.Master) *Delegate {
	d := &Delegate{
		ks:       ks,
		relayers: map[types.Network]types.Relayer{},
	}
	return d
}

// AddRelayer registers the relayer r, or a disabled placeholder if nil.
// NOT THREAD SAFE
func (d Delegate) AddRelayer(n types.Network, r types.Relayer) {
	d.relayers[n] = r
}

// Start starts all relayers it manages.
func (d Delegate) Start(ctx context.Context) error {
	var err error
	for _, r := range d.relayers {
		err = multierr.Combine(err, r.Start(ctx))
	}
	return err
}

// A Delegate relayer on close will close all relayers it manages.
func (d Delegate) Close() error {
	var err error
	for _, r := range d.relayers {
		err = multierr.Combine(err, r.Close())
	}
	return err
}

// A Delegate relayer is healthy if all relayers it manages are ready.
func (d Delegate) Ready() error {
	var err error
	for _, r := range d.relayers {
		err = multierr.Combine(err, r.Ready())
	}
	return err
}

// A Delegate relayer is healthy if all relayers it manages are healthy.
func (d Delegate) Healthy() error {
	var err error
	for _, r := range d.relayers {
		err = multierr.Combine(err, r.Healthy())
	}
	return err
}

func (d Delegate) NewConfigWatcher(relay types.Network, args types.ConfigWatcherArgs) (types.ConfigWatcher, error) {
	switch relay {
	case types.EVM:
		r, exists := d.relayers[types.EVM]
		if !exists {
			return nil, errors.New("no EVM relay found; is EVM enabled?")
		}
		return r.NewConfigWatcher(args)
	case types.Solana:
		r, exists := d.relayers[types.Solana]
		if !exists {
			return nil, errors.New("no Solana relay found; is Solana enabled?")
		}
		return r.NewConfigWatcher(args)
	case types.Terra:
		r, exists := d.relayers[types.Terra]
		if !exists {
			return nil, errors.New("no Terra relay found; is Terra enabled?")
		}
		return r.NewConfigWatcher(args)
	default:
		return nil, errors.Errorf("unknown relayer network type: %s", relay)
	}
}

func (d Delegate) NewMedianProvider(relay types.Network, args types.PluginArgs) (types.MedianProvider, error) {
	switch relay {
	case types.EVM:
		r, exists := d.relayers[types.EVM]
		if !exists {
			return nil, errors.New("no EVM relay found; is EVM enabled?")
		}
		return r.NewMedianProvider(args)
	case types.Solana:
		r, exists := d.relayers[types.Solana]
		if !exists {
			return nil, errors.New("no Solana relay found; is Solana enabled?")
		}
		return r.NewMedianProvider(args)
	case types.Terra:
		r, exists := d.relayers[types.Terra]
		if !exists {
			return nil, errors.New("no Terra relay found; is Terra enabled?")
		}
		return r.NewMedianProvider(args)
	default:
		return nil, errors.Errorf("unknown relayer network type: %s", relay)
	}
}
