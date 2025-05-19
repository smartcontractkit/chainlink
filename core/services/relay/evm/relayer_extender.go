package evm

import (
	"errors"
	"fmt"

	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-evm/pkg/config/toml"
	"github.com/smartcontractkit/chainlink-evm/pkg/gas/rollups"
	"github.com/smartcontractkit/chainlink-evm/pkg/keys"
	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
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

func NewLegacyChains(
	lggr logger.Logger,
	ks keystore.Eth,
	chainOpts legacyevm.ChainOpts,
) (result []legacyevm.Chain, err error) {
	unique := make(map[string]struct{})
	var enabled []*toml.EVMConfig
	for i, cfg := range chainOpts.ChainConfigs {
		_, alreadyExists := unique[cfg.ChainID.String()]
		if alreadyExists {
			return nil, fmt.Errorf("duplicate chain definition for evm chain id %s", cfg.ChainID.String())
		}
		unique[cfg.ChainID.String()] = struct{}{}
		if chainOpts.ChainConfigs[i].IsEnabled() {
			enabled = append(enabled, chainOpts.ChainConfigs[i])
		}
	}
	newChainStore := chainOpts.GenChainStore
	if newChainStore == nil {
		newChainStore = keys.NewChainStore
	}

	// map with lazy initialization for the txm to access evm clients for different chain
	var clientsByChainID = make(map[string]rollups.DAClient)
	for i := range enabled {
		cid := enabled[i].ChainID.ToInt()
		opts := legacyevm.ChainRelayOpts{
			Logger:    logger.Named(lggr, cid.String()),
			KeyStore:  newChainStore(keystore.NewEthSigner(ks, cid), cid),
			ChainOpts: chainOpts,
		}

		opts.Logger.Infow(fmt.Sprintf("Loading chain %s", cid), "evmChainID", cid)
		chain, err2 := legacyevm.NewTOMLChain(enabled[i], opts, clientsByChainID)
		if err2 != nil {
			err = multierr.Combine(err, fmt.Errorf("failed to create chain %s: %w", cid, err2))
			continue
		}

		clientsByChainID[cid.String()] = chain.Client()
		result = append(result, chain)
	}
	return
}
func NewLegacyChainsAndConfig(
	lggr logger.Logger,
	ks keystore.Eth,
	chainOpts legacyevm.ChainOpts,
) (*LegacyChainsAndConfig, error) {
	result, err := NewLegacyChains(lggr, ks, chainOpts)
	// always return because it's accumulating errors
	return &LegacyChainsAndConfig{result, chainOpts.ChainConfigs}, err
}
