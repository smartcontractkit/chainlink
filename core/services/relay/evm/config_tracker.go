package evm

import (
	"context"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"

	evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"
	httypes "github.com/smartcontractkit/chainlink/core/chains/evm/headtracker/types"
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/config"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// ConfigTracker tracks the config of any contract implementing OCR2Abstract.
type ConfigTracker struct {
	utils.StartStopOnce
	lggr        logger.Logger
	client      evmclient.Client
	addr        common.Address
	contractABI abi.ABI
	chainType   config.ChainType

	latestBlockHeight   int64
	latestBlockHeightMu sync.RWMutex

	headBroadcaster  httypes.HeadBroadcaster
	unsubscribeHeads func()

	chStop chan struct{}
}

// NewConfigTracker builds a new config tracker
func NewConfigTracker(lggr logger.Logger, contractABI abi.ABI, client evmclient.Client, addr common.Address, chainType config.ChainType, headBroadcaster httypes.HeadBroadcaster) *ConfigTracker {
	return &ConfigTracker{
		client:              client,
		addr:                addr,
		contractABI:         contractABI,
		chainType:           chainType,
		latestBlockHeight:   -1,
		latestBlockHeightMu: sync.RWMutex{},
		lggr:                lggr,
		headBroadcaster:     headBroadcaster,
		unsubscribeHeads:    nil,
		chStop:              make(chan struct{}),
	}
}

// Start starts the config tracker in particular subscribing to the head broadcaster.
func (c *ConfigTracker) Start() error {
	return c.StartOnce("ConfigTracker", func() (err error) {
		var latestHead *evmtypes.Head
		latestHead, c.unsubscribeHeads = c.headBroadcaster.Subscribe(c)
		if latestHead != nil {
			c.setLatestBlockHeight(*latestHead)
		}

		return nil
	})
}

// Close cancels and requests and unsubscribes from the head broadcaster
func (c *ConfigTracker) Close() error {
	close(c.chStop)
	c.unsubscribeHeads()
	return nil
}

// Notify not implemented
func (c *ConfigTracker) Notify() <-chan struct{} {
	return nil
}

func callContract(ctx context.Context, addr common.Address, contractABI abi.ABI, method string, args []interface{}, caller contractReader) ([]interface{}, error) {
	input, err := contractABI.Pack(method, args...)
	if err != nil {
		return nil, err
	}
	output, err := caller.CallContract(ctx, ethereum.CallMsg{To: &addr, Data: input}, nil)
	if err != nil {
		return nil, err
	}
	return contractABI.Unpack(method, output)
}

// LatestConfigDetails queries an OCR2Abstract contract for the latest config details
func (c *ConfigTracker) LatestConfigDetails(ctx context.Context) (changedInBlock uint64, configDigest ocrtypes.ConfigDigest, err error) {
	latestConfigDetails, err := callContract(ctx, c.addr, c.contractABI, "latestConfigDetails", nil, c.client)
	if err != nil {
		return 0, ocrtypes.ConfigDigest{}, err
	}
	// Panic on these conversions erroring, would mean a broken contract.
	changedInBlock = uint64(*abi.ConvertType(latestConfigDetails[1], new(uint32)).(*uint32))
	configDigest = *abi.ConvertType(latestConfigDetails[2], new([32]byte)).(*[32]byte)
	return
}

// LatestConfig queries an OCR2Abstract contract for the latest config contents.
func (c *ConfigTracker) LatestConfig(ctx context.Context, changedInBlock uint64) (ocrtypes.ContractConfig, error) {
	topics, err := abi.MakeTopics([]interface{}{c.contractABI.Events["ConfigSet"].ID})
	if err != nil {
		return ocrtypes.ContractConfig{}, err
	}
	query := ethereum.FilterQuery{
		Addresses: []common.Address{c.addr},
		Topics:    topics,
		FromBlock: new(big.Int).SetUint64(changedInBlock),
		ToBlock:   new(big.Int).SetUint64(changedInBlock),
	}
	logs, err := c.client.FilterLogs(ctx, query)
	if err != nil {
		return ocrtypes.ContractConfig{}, err
	}
	if len(logs) == 0 {
		err = errors.New("Contract not configured yet")
		c.lggr.Warnw(err.Error())
		return ocrtypes.ContractConfig{}, err
	}
	return parseConfigSet(c.contractABI, logs[len(logs)-1])
}

func parseConfigSet(a abi.ABI, log types.Log) (ocrtypes.ContractConfig, error) {
	var changed struct {
		PreviousConfigBlockNumber uint32
		ConfigDigest              [32]byte
		ConfigCount               uint64
		Signers                   []common.Address
		Transmitters              []common.Address
		F                         uint8
		OnchainConfig             []byte
		OffchainConfigVersion     uint64
		OffchainConfig            []byte
	}
	// Use bound contract solely for its unpack log logic
	// which only uses the abi.
	err := bind.NewBoundContract(common.Address{}, a, nil, nil, nil).UnpackLog(&changed, "ConfigSet", log)
	if err != nil {
		return ocrtypes.ContractConfig{}, err
	}
	var transmitAccounts []ocrtypes.Account
	for _, addr := range changed.Transmitters {
		transmitAccounts = append(transmitAccounts, ocrtypes.Account(addr.Hex()))
	}
	var signers []ocrtypes.OnchainPublicKey
	for _, addr := range changed.Signers {
		addr := addr
		signers = append(signers, addr[:])
	}
	return ocrtypes.ContractConfig{
		ConfigDigest:          changed.ConfigDigest,
		ConfigCount:           changed.ConfigCount,
		Signers:               signers,
		Transmitters:          transmitAccounts,
		F:                     changed.F,
		OnchainConfig:         changed.OnchainConfig,
		OffchainConfigVersion: changed.OffchainConfigVersion,
		OffchainConfig:        changed.OffchainConfig,
	}, nil
}

// Connect conforms to HeadTrackable
func (c *ConfigTracker) Connect(*evmtypes.Head) error { return nil }

// OnNewLongestChain conformed to HeadTrackable and updates latestBlockHeight
func (c *ConfigTracker) OnNewLongestChain(_ context.Context, h *evmtypes.Head) {
	c.setLatestBlockHeight(*h)
}

func (c *ConfigTracker) setLatestBlockHeight(h evmtypes.Head) {
	var num int64
	if h.L1BlockNumber.Valid {
		num = h.L1BlockNumber.Int64
	} else {
		num = h.Number
	}
	c.latestBlockHeightMu.Lock()
	defer c.latestBlockHeightMu.Unlock()
	if num > c.latestBlockHeight {
		c.latestBlockHeight = num
	}
}

func (c *ConfigTracker) getLatestBlockHeight() int64 {
	c.latestBlockHeightMu.RLock()
	defer c.latestBlockHeightMu.RUnlock()
	return c.latestBlockHeight
}

// LatestBlockHeight returns the latest blockheight either from the cache or
// falling back to querying the node.
func (c *ConfigTracker) LatestBlockHeight(ctx context.Context) (blockHeight uint64, err error) {
	// We skip confirmation checking anyway on Optimism so there's no need to care
	// about the block height; we have no way of getting the L1 block height anyway
	if c.chainType != "" {
		return 0, nil
	}
	latestBlockHeight := c.getLatestBlockHeight()
	if latestBlockHeight >= 0 {
		return uint64(latestBlockHeight), nil
	}

	var cancel context.CancelFunc
	ctx, cancel = utils.WithCloseChan(ctx, c.chStop)
	defer cancel()

	c.lggr.Debugw("ConfigTracker: still waiting for first head, falling back to on-chain lookup")
	h, err := c.client.HeadByNumber(ctx, nil)
	if err != nil {
		return 0, err
	}
	if h == nil {
		return 0, errors.New("got nil head")
	}
	if h.L1BlockNumber.Valid {
		return uint64(h.L1BlockNumber.Int64), nil
	}
	return uint64(h.Number), nil
}
