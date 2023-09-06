package evm

import (
	"context"
	"fmt"
	"math/big"

	"github.com/pkg/errors"
	"go.uber.org/multierr"

	relaytypes "github.com/smartcontractkit/chainlink-relay/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm"
	evmchain "github.com/smartcontractkit/chainlink/v2/core/chains/evm"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
)

// ErrNoChains indicates that no EVM chains have been started
var ErrNoChains = errors.New("no EVM chains loaded")

type EVMChainRelayerExtender interface {
	relay.RelayerExt
	Chain() evmchain.Chain
	// compatibility remove after BCF-2441
	ChainStatus(ctx context.Context, id string) (relaytypes.ChainStatus, error)
	ChainStatuses(ctx context.Context, offset, limit int) ([]relaytypes.ChainStatus, int, error)
	NodeStatuses(ctx context.Context, offset, limit int, chainIDs ...string) (nodes []relaytypes.NodeStatus, count int, err error)
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

var ErrInconsistentChainRelayerExtender = errors.New("inconsistent evm chain relayer extender")

// Legacy interface remove after BFC-2441, BCF-2564

func (s *ChainRelayerExt) ChainStatus(ctx context.Context, id string) (relaytypes.ChainStatus, error) {
	if s.chain.ID().String() != id {
		return relaytypes.ChainStatus{}, fmt.Errorf("%w: given id %q does not match expected id %q", ErrInconsistentChainRelayerExtender, id, s.chain.ID())
	}
	return s.chain.GetChainStatus(ctx)
}

func (s *ChainRelayerExt) ChainStatuses(ctx context.Context, offset, limit int) ([]relaytypes.ChainStatus, int, error) {
	stat, err := s.chain.GetChainStatus(ctx)
	if err != nil {
		return nil, -1, err
	}
	return []relaytypes.ChainStatus{stat}, 1, nil

}

func (s *ChainRelayerExt) NodeStatuses(ctx context.Context, offset, limit int, chainIDs ...string) (nodes []relaytypes.NodeStatus, total int, err error) {
	if len(chainIDs) > 1 {
		return nil, -1, fmt.Errorf("single chain chain set only support one chain id. got %v", chainIDs)
	}
	cid := chainIDs[0]
	if cid != s.chain.ID().String() {
		return nil, -1, fmt.Errorf("unknown chain id %s. expected %s", cid, s.chain.ID())
	}
	nodes, _, total, err = s.ListNodeStatuses(ctx, int32(limit), "")
	if err != nil {
		return nil, -1, err
	}
	if len(nodes) < offset {
		return []relaytypes.NodeStatus{}, -1, fmt.Errorf("out of range")
	}
	if limit <= 0 {
		limit = len(nodes)
	} else if len(nodes) < limit {
		limit = len(nodes)
	}
	return nodes[offset:limit], total, nil

}

func NewChainRelayerExtenders(ctx context.Context, opts evmchain.ChainRelayExtenderConfig) (*ChainRelayerExtenders, error) {
	if err := opts.Check(); err != nil {
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
			Logger:        opts.Logger.Named(cid),
			RelayerConfig: opts.RelayerConfig,
			DB:            opts.DB,
			KeyStore:      opts.KeyStore,
		}

		privOpts.Logger.Infow(fmt.Sprintf("Loading chain %s", cid), "evmChainID", cid)
		chain, err2 := evmchain.NewTOMLChain(ctx, enabled[i], privOpts)
		if err2 != nil {
			err = multierr.Combine(err, err2)
			continue
		}

		s := &ChainRelayerExt{
			chain: chain,
		}
		result = append(result, s)
	}
	return newChainRelayerExtsFromSlice(result, opts.AppConfig), nil
}
