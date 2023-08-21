package evm

import (
	"context"
	"fmt"
	"math/big"
	"sync"

	"github.com/pkg/errors"
	"go.uber.org/multierr"
	"golang.org/x/exp/maps"

	relaytypes "github.com/smartcontractkit/chainlink-relay/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains"
	evmchain "github.com/smartcontractkit/chainlink/v2/core/chains/evm"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
)

// ErrNoChains indicates that no EVM chains have been started
var ErrNoChains = errors.New("no EVM chains loaded")

var _ legacyChainSet = &chainSet{}

type legacyChainSet interface {
	services.ServiceCtx
	chains.ChainStatuser
	chains.NodesStatuser

	Get(id *big.Int) (evmchain.Chain, error)

	Default() (evmchain.Chain, error)
	Chains() []evmchain.Chain
	ChainCount() int

	Configs() evmtypes.Configs

	SendTx(ctx context.Context, chainID, from, to string, amount *big.Int, balanceCheck bool) error
}

type EVMChainRelayerExtender interface {
	relay.RelayerExt
	Chain() evmchain.Chain
	Default() bool
}

type EVMChainRelayerExtenderSlicer interface {
	Slice() []EVMChainRelayerExtender
	Len() int
	ChainNodeConfigs() evmtypes.Configs
}

type ChainRelayerExtenders struct {
	exts []EVMChainRelayerExtender
	cfgs evmtypes.Configs
}

var _ EVMChainRelayerExtenderSlicer = &ChainRelayerExtenders{}

func NewLegacyChainsFromRelayerExtenders(exts EVMChainRelayerExtenderSlicer) *evmchain.LegacyChains {

	m := make(map[string]evmchain.Chain)
	var dflt evmchain.Chain
	for _, r := range exts.Slice() {
		m[r.Chain().ID().String()] = r.Chain()
		if r.Default() {
			dflt = r.Chain()
		}
	}
	l := evmchain.NewLegacyChains(exts.ChainNodeConfigs(), m)
	if dflt != nil {
		l.SetDefault(dflt)
	}
	return l
}

func newChainRelayerExtsFromSlice(exts []*ChainRelayerExt) *ChainRelayerExtenders {
	temp := make([]EVMChainRelayerExtender, len(exts))
	for i := range exts {
		temp[i] = exts[i]
	}
	return &ChainRelayerExtenders{
		exts: temp,
	}
}

func (c *ChainRelayerExtenders) ChainNodeConfigs() evmtypes.Configs {
	return c.cfgs
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
	// TODO remove this altogether. BFC-2440
	cs        *chainSet
	isDefault bool
}

var _ EVMChainRelayerExtender = &ChainRelayerExt{}

func (s *ChainRelayerExt) Chain() evmchain.Chain {
	return s.chain
}

func (s *ChainRelayerExt) Default() bool {
	return s.isDefault
}

var ErrCorruptEVMChain = errors.New("corrupt evm chain")

func (s *ChainRelayerExt) Start(ctx context.Context) error {
	if len(s.cs.chains) > 1 {
		err := fmt.Errorf("%w: internal error more than one chain (%d)", ErrCorruptEVMChain, len(s.cs.chains))
		panic(err)
	}
	return s.cs.Start(ctx)
}

func (s *ChainRelayerExt) Close() (err error) {
	return s.cs.Close()
}

func (s *ChainRelayerExt) Name() string {
	// we set each private chainSet logger to contain the chain id
	return s.cs.Name()
}

func (s *ChainRelayerExt) HealthReport() map[string]error {
	return s.cs.HealthReport()
}

func (s *ChainRelayerExt) Ready() (err error) {
	return s.cs.Ready()
}

var ErrInconsistentChainRelayerExtender = errors.New("inconsistent evm chain relayer extender")

func (s *ChainRelayerExt) ChainStatus(ctx context.Context, id string) (relaytypes.ChainStatus, error) {
	// TODO BCF-2441: update relayer interface
	// we need to implement the interface, but passing id doesn't really make sense because there is only
	// one chain here. check the id here to provide clear error reporting.
	if s.chain.ID().String() != id {
		return relaytypes.ChainStatus{}, fmt.Errorf("%w: given id %q does not match expected id %q", ErrInconsistentChainRelayerExtender, id, s.chain.ID())
	}
	return s.cs.ChainStatus(ctx, id)
}

func (s *ChainRelayerExt) ChainStatuses(ctx context.Context, offset, limit int) ([]relaytypes.ChainStatus, int, error) {
	stat, err := s.cs.ChainStatus(ctx, s.chain.ID().String())
	if err != nil {
		return nil, -1, err
	}
	return []relaytypes.ChainStatus{stat}, 1, nil

}

func (s *ChainRelayerExt) NodeStatuses(ctx context.Context, offset, limit int, chainIDs ...string) (nodes []relaytypes.NodeStatus, count int, err error) {
	if len(chainIDs) > 1 {
		return nil, -1, fmt.Errorf("single chain chain set only support one chain id. got %v", chainIDs)
	}
	cid := chainIDs[0]
	if cid != s.chain.ID().String() {
		return nil, -1, fmt.Errorf("unknown chain id %s. expected %s", cid, s.chain.ID())
	}
	return s.cs.NodeStatuses(ctx, offset, limit, chainIDs...)
}

func (s *ChainRelayerExt) SendTx(ctx context.Context, chainID, from, to string, amount *big.Int, balanceCheck bool) error {
	return s.cs.SendTx(ctx, chainID, from, to, amount, balanceCheck)
}

type chainSet struct {
	defaultID     *big.Int
	chains        map[string]evmchain.Chain
	startedChains []evmchain.Chain
	chainsMu      sync.RWMutex
	logger        logger.Logger
	opts          evmchain.ChainRelayExtenderConfig
}

func (cll *chainSet) Start(ctx context.Context) error {
	if !cll.opts.GeneralConfig.EVMEnabled() {
		cll.logger.Warn("EVM is disabled, no EVM-based chains will be started")
		return nil
	}
	if !cll.opts.GeneralConfig.EVMRPCEnabled() {
		cll.logger.Warn("EVM RPC connections are disabled. Chainlink will not connect to any EVM RPC node.")
	}
	var ms services.MultiStart
	for _, c := range cll.Chains() {
		if err := ms.Start(ctx, c); err != nil {
			return errors.Wrapf(err, "failed to start chain %q", c.ID().String())
		}
		cll.startedChains = append(cll.startedChains, c)
	}
	evmChainIDs := make([]*big.Int, len(cll.startedChains))
	for i, c := range cll.startedChains {
		evmChainIDs[i] = c.ID()
	}
	defChainID := "unspecified"
	if cll.defaultID != nil {
		defChainID = fmt.Sprintf("%q", cll.defaultID.String())
	}
	cll.logger.Infow(fmt.Sprintf("EVM: Started %d/%d chains, default chain ID is %s", len(cll.startedChains), len(cll.Chains()), defChainID), "startedEvmChainIDs", evmChainIDs)
	return nil
}

func (cll *chainSet) Close() (err error) {
	cll.logger.Debug("EVM: stopping")
	for _, c := range cll.startedChains {
		err = multierr.Combine(err, c.Close())
	}
	return
}

func (cll *chainSet) Name() string {
	return cll.logger.Name()
}

func (cll *chainSet) HealthReport() map[string]error {
	report := map[string]error{}
	for _, c := range cll.Chains() {
		maps.Copy(report, c.HealthReport())
	}
	return report
}

func (cll *chainSet) Ready() (err error) {
	for _, c := range cll.Chains() {
		err = multierr.Combine(err, c.Ready())
	}
	return
}

func (cll *chainSet) Get(id *big.Int) (evmchain.Chain, error) {
	if id == nil {
		if cll.defaultID == nil {
			cll.logger.Debug("Chain ID not specified, and default is nil")
			return nil, errors.New("chain ID not specified, and default is nil")
		}
		cll.logger.Debugf("Chain ID not specified, using default: %s", cll.defaultID.String())
		return cll.Default()
	}
	return cll.get(id.String())
}

func (cll *chainSet) get(id string) (evmchain.Chain, error) {
	cll.chainsMu.RLock()
	defer cll.chainsMu.RUnlock()
	c, exists := cll.chains[id]
	if exists {
		return c, nil
	}
	return nil, errors.Wrap(chains.ErrNotFound, fmt.Sprintf("failed to get chain with id %s", id))
}

func (cll *chainSet) ChainStatus(ctx context.Context, id string) (cfg relaytypes.ChainStatus, err error) {
	var cs []relaytypes.ChainStatus
	cs, _, err = cll.opts.OperationalConfigs.Chains(0, -1, id)
	if err != nil {
		return
	}
	l := len(cs)
	if l == 0 {
		err = chains.ErrNotFound
		return
	}
	if l > 1 {
		err = fmt.Errorf("multiple chains found: %d", len(cs))
		return
	}
	cfg = cs[0]
	return
}

func (cll *chainSet) ChainStatuses(ctx context.Context, offset, limit int) ([]relaytypes.ChainStatus, int, error) {
	return cll.opts.OperationalConfigs.Chains(offset, limit)
}

func (cll *chainSet) Default() (evmchain.Chain, error) {
	cll.chainsMu.RLock()
	length := len(cll.chains)
	cll.chainsMu.RUnlock()
	if length == 0 {
		return nil, errors.Wrap(ErrNoChains, "cannot get default EVM chain; no EVM chains are available")
	}
	if cll.defaultID == nil {
		// This is an invariant violation; if any chains are available then a
		// default should _always_ have been set in the constructor
		return nil, errors.New("no default chain ID specified")
	}

	return cll.Get(cll.defaultID)
}

func (cll *chainSet) Chains() (c []evmchain.Chain) {
	cll.chainsMu.RLock()
	defer cll.chainsMu.RUnlock()
	for _, chain := range cll.chains {
		c = append(c, chain)
	}
	return c
}

func (cll *chainSet) ChainCount() int {
	cll.chainsMu.RLock()
	defer cll.chainsMu.RUnlock()
	return len(cll.chains)
}

func (cll *chainSet) Configs() evmtypes.Configs {
	return cll.opts.OperationalConfigs
}

func (cll *chainSet) NodeStatuses(ctx context.Context, offset, limit int, chainIDs ...string) (nodes []relaytypes.NodeStatus, count int, err error) {
	nodes, count, err = cll.opts.OperationalConfigs.NodeStatusesPaged(offset, limit, chainIDs...)
	if err != nil {
		err = errors.Wrap(err, "GetNodesForChain failed to load nodes from DB")
		return
	}
	for i := range nodes {
		cll.addStateToNode(&nodes[i])
	}
	return
}

func (cll *chainSet) addStateToNode(n *relaytypes.NodeStatus) {
	cll.chainsMu.RLock()
	chain, exists := cll.chains[n.ChainID]
	cll.chainsMu.RUnlock()
	if !exists {
		// The EVM chain is disabled
		n.State = "Disabled"
		return
	}
	states := chain.Client().NodeStates()
	if states == nil {
		n.State = "Unknown"
		return
	}
	state, exists := states[n.Name]
	if exists {
		n.State = state
		return
	}
	// The node is in the DB and the chain is enabled but it's not running
	n.State = "NotLoaded"
}

func (cll *chainSet) SendTx(ctx context.Context, chainID, from, to string, amount *big.Int, balanceCheck bool) error {
	chain, err := cll.get(chainID)
	if err != nil {
		return err
	}

	return chain.SendTx(ctx, from, to, amount, balanceCheck)
}

func NewChainRelayerExtenders(ctx context.Context, opts evmchain.ChainRelayExtenderConfig) (*ChainRelayerExtenders, error) {
	if err := opts.Check(); err != nil {
		return nil, err
	}
	evmConfigs := opts.GeneralConfig.EVMConfigs()
	var enabled []*toml.EVMConfig
	for i := range evmConfigs {
		if evmConfigs[i].IsEnabled() {
			enabled = append(enabled, evmConfigs[i])
		}
	}

	defaultChainID := opts.GeneralConfig.DefaultChainID()
	if defaultChainID == nil && len(enabled) >= 1 {
		defaultChainID = enabled[0].ChainID.ToInt()
		if len(enabled) > 1 {
			opts.Logger.Debugf("Multiple chains present, default chain: %s", defaultChainID.String())
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
		cll := newChainSet(privOpts)

		cll.logger.Infow(fmt.Sprintf("Loading chain %s", cid), "evmChainID", cid)
		chain, err2 := evmchain.NewTOMLChain(ctx, enabled[i], privOpts)
		if err2 != nil {
			err = multierr.Combine(err, err2)
			continue
		}
		if _, exists := cll.chains[cid]; exists {
			return nil, errors.Errorf("duplicate chain with ID %s", cid)
		}
		cll.chains[cid] = chain

		s := &ChainRelayerExt{
			chain:     chain,
			cs:        cll,
			isDefault: (cid == defaultChainID.String()),
		}
		result = append(result, s)
	}
	return newChainRelayerExtsFromSlice(result), nil
}

func newChainSet(opts evmchain.ChainRelayExtenderConfig) *chainSet {
	return &chainSet{
		chains:        make(map[string]evmchain.Chain),
		startedChains: make([]evmchain.Chain, 0),
		logger:        opts.Logger.Named("ChainSet"),
		opts:          opts,
	}
}
