package evm

import (
	"context"
	"fmt"
	"math/big"

	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink-relay/pkg/loop"
	relaytypes "github.com/smartcontractkit/chainlink-relay/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm"
	evmchain "github.com/smartcontractkit/chainlink/v2/core/chains/evm"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"
)

// ErrNoChains indicates that no EVM chains have been started
var ErrNoChains = errors.New("no EVM chains loaded")

type EVMChainRelayerExtender interface {
	loop.RelayerExt
	Chain() evmchain.Chain
}

type EVMChainRelayerExtenderSlicer interface {
	Slice() []EVMChainRelayerExtender
	Len() int
	AppConfig() evmchain.AppConfig
}

type ChainRelayerExtenders struct {
	exts []EVMChainRelayerExtender
	cfg  evmchain.AppConfig
}

var _ EVMChainRelayerExtenderSlicer = &ChainRelayerExtenders{}

func NewLegacyChainsFromRelayerExtenders(exts EVMChainRelayerExtenderSlicer) *evmchain.LegacyChains {
	m := make(map[string]evmchain.Chain)
	for _, r := range exts.Slice() {
		m[r.Chain().ID().String()] = r.Chain()
	}
	return evmchain.NewLegacyChains(m, exts.AppConfig().EVMConfigs())
}

func newChainRelayerExtsFromSlice(exts []*ChainRelayerExt, appConfig evm.AppConfig) *ChainRelayerExtenders {
	temp := make([]EVMChainRelayerExtender, len(exts))
	for i := range exts {
		temp[i] = exts[i]
	}
	return &ChainRelayerExtenders{
		exts: temp,
		cfg:  appConfig,
	}
}

func (c *ChainRelayerExtenders) AppConfig() evmchain.AppConfig {
	return c.cfg
}

func (c *ChainRelayerExtenders) Slice() []EVMChainRelayerExtender {
	return c.exts
}

func (c *ChainRelayerExtenders) Len() int {
	return len(c.exts)
}

// implements OneChain
type ChainRelayerExt struct {
	chain evmchain.Chain
}

var _ EVMChainRelayerExtender = &ChainRelayerExt{}

func (s *ChainRelayerExt) GetChainStatus(ctx context.Context) (relaytypes.ChainStatus, error) {
	return s.chain.GetChainStatus(ctx)
}

func (s *ChainRelayerExt) ListNodeStatuses(ctx context.Context, pageSize int32, pageToken string) (stats []relaytypes.NodeStatus, nextPageToken string, total int, err error) {
	return s.chain.ListNodeStatuses(ctx, pageSize, pageToken)
}

func (s *ChainRelayerExt) Transact(ctx context.Context, from, to string, amount *big.Int, balanceCheck bool) error {
	return s.chain.Transact(ctx, from, to, amount, balanceCheck)
}

func (s *ChainRelayerExt) ID() string {
	return s.chain.ID().String()
}

func (s *ChainRelayerExt) Chain() evmchain.Chain {
	return s.chain
}

var ErrCorruptEVMChain = errors.New("corrupt evm chain")

func (s *ChainRelayerExt) Start(ctx context.Context) error {
	return s.chain.Start(ctx)
}

func (s *ChainRelayerExt) Close() (err error) {
	return s.chain.Close()
}

func (s *ChainRelayerExt) Name() string {
	return s.chain.Name()
}

func (s *ChainRelayerExt) HealthReport() map[string]error {
	return s.chain.HealthReport()
}

func (s *ChainRelayerExt) Ready() (err error) {
	return s.chain.Ready()
}

func NewChainRelayerExtenders(ctx context.Context, opts evmchain.ChainRelayExtenderConfig) (*ChainRelayerExtenders, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
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

	var result []*ChainRelayerExt
	var err error
	for i := range enabled {

		cid := enabled[i].ChainID.String()
		privOpts := evmchain.ChainRelayExtenderConfig{
			Logger:    opts.Logger.Named(cid),
			ChainOpts: opts.ChainOpts,
			KeyStore:  opts.KeyStore,
		}

		privOpts.Logger.Infow(fmt.Sprintf("Loading chain %s", cid), "evmChainID", cid)
		chain, err2 := evmchain.NewTOMLChain(ctx, enabled[i], privOpts)
		if err2 != nil {
			err = multierr.Combine(err, fmt.Errorf("failed to create chain %s: %w", cid, err2))
			continue
		}

		s := &ChainRelayerExt{
			chain: chain,
		}
		result = append(result, s)
	}
	// always return because it's accumulating errors
	return newChainRelayerExtsFromSlice(result, opts.AppConfig), err
}
