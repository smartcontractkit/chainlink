package evm

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm"
)

// ErrNoChains indicates that no EVM chains have been started
var ErrNoChains = errors.New("no EVM chains loaded")

type LegacyChainsAndConfig struct {
	rs  []legacyevm.Chain
	cfg toml.EVMConfigs
}

func (r *LegacyChainsAndConfig) NewLegacyChains() *legacyevm.LegacyChains {
	m := make(map[string]legacyevm.Chain)
	for _, r := range r.Slice() {
		m[r.ID().String()] = r
	}
	return legacyevm.NewLegacyChains(m, r.cfg)
}

func (r *LegacyChainsAndConfig) Slice() []legacyevm.Chain {
	return r.rs
}

func (r *LegacyChainsAndConfig) Len() int {
	return len(r.rs)
}

func NewLegacyChains(ctx context.Context, opts legacyevm.ChainRelayOpts) (result []legacyevm.Chain, err error) {
	if err = opts.Validate(); err != nil {
		return
	}

	unique := make(map[string]struct{})

	evmConfigs := opts.AppConfig.EVMConfigs()
	var enabled []*toml.EVMConfig
	for i, cfg := range evmConfigs {
		_, alreadyExists := unique[cfg.ChainID.String()]
		if alreadyExists {
			return nil, fmt.Errorf("duplicate chain definition for evm chain id %s", cfg.ChainID.String())
		}
		unique[cfg.ChainID.String()] = struct{}{}
		if evmConfigs[i].IsEnabled() {
			enabled = append(enabled, evmConfigs[i])
		}
	}

	for i := range enabled {
		cid := enabled[i].ChainID.String()
		privOpts := legacyevm.ChainRelayOpts{
			Logger:    opts.Logger.Named(cid),
			ChainOpts: opts.ChainOpts,
			KeyStore:  opts.KeyStore,
		}

		privOpts.Logger.Infow(fmt.Sprintf("Loading chain %s", cid), "evmChainID", cid)
		chain, err2 := legacyevm.NewTOMLChain(ctx, enabled[i], privOpts)
		if err2 != nil {
			err = multierr.Combine(err, fmt.Errorf("failed to create chain %s: %w", cid, err2))
			continue
		}

		result = append(result, chain)
	}
	return
}
func NewLegacyChainsAndConfig(ctx context.Context, opts legacyevm.ChainRelayOpts) (*LegacyChainsAndConfig, error) {
	result, err := NewLegacyChains(ctx, opts)
	// always return because it's accumulating errors
	return &LegacyChainsAndConfig{result, opts.AppConfig.EVMConfigs()}, err
}
