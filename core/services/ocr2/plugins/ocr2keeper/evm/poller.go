package evm

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	registry "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_wrapper2_0"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

var (
	ErrFailedToGetLatestBlock = fmt.Errorf("failed to get latest block from log poller")
)

type RegistryPoller struct {
	// provided dependencies
	logger logger.Logger
	p      logpoller.LogPoller
	addr   common.Address

	// properties initialized in constructor
	contract *registry.KeeperRegistry
	abi      abi.ABI

	// run state properties
	sync     utils.StartStopOnce
	mu       sync.RWMutex
	runState int
	runError error
}

func NewRegistryPoller(
	p logpoller.LogPoller,
	l logger.Logger,
	client evmclient.Client,
	addr common.Address,
) (*RegistryPoller, error) {
	contract, err := registry.NewKeeperRegistry(common.HexToAddress("0x"), client)
	if err != nil {
		return nil, err
	}

	abi, err := abi.JSON(strings.NewReader(registry.KeeperRegistryABI))
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrABINotParsable, err)
	}

	return &RegistryPoller{
		logger:   l,
		p:        p,
		addr:     addr,
		contract: contract,
		abi:      abi,
	}, nil
}

// Name implements the job.ServiceCtx interface
func (rp *RegistryPoller) Name() string {
	return RegistryEventFilterName(rp.addr)
}

// Start implements the job.ServiceCtx interface
func (rp *RegistryPoller) Start(ctx context.Context) error {
	return rp.sync.StartOnce("AutomationLogProvider", func() error {
		rp.mu.Lock()
		defer rp.mu.Unlock()

		rp.register()
		rp.runState = 1

		return nil
	})
}

// Stop implements the job.ServiceCtx interface
func (rp *RegistryPoller) Close() error {
	return rp.sync.StopOnce("AutomationRegistry", func() error {
		rp.mu.Lock()
		defer rp.mu.Unlock()

		rp.unregister()
		rp.runState = 0
		rp.runError = nil

		return nil
	})
}

// Ready implements the job.ServiceCtx interface
func (rp *RegistryPoller) Ready() error {
	rp.mu.RLock()
	defer rp.mu.RUnlock()

	if rp.runState == 1 {
		return nil
	}

	return rp.sync.Ready()
}

// HealthReport implements the job.ServiceCtx interface
func (rp *RegistryPoller) HealthReport() map[string]error {
	rp.mu.RLock()
	defer rp.mu.RUnlock()

	if rp.runState > 1 {
		rp.sync.SvcErrBuffer.Append(fmt.Errorf("failed run state: %w", rp.runError))
	}

	return map[string]error{rp.Name(): rp.sync.Healthy()}
}

func (rp *RegistryPoller) register() error {
	var upkeepStateEvents = []common.Hash{
		registry.KeeperRegistryUpkeepRegistered{}.Topic(),  // adds new upkeep id to registry
		registry.KeeperRegistryUpkeepReceived{}.Topic(),    // adds new upkeep id to registry via migration
		registry.KeeperRegistryUpkeepGasLimitSet{}.Topic(), // unpauses an upkeep
		registry.KeeperRegistryUpkeepUnpaused{}.Topic(),    // updates the gas limit for an upkeep
	}

	var upkeepActiveEvents = []common.Hash{
		registry.KeeperRegistryUpkeepPerformed{}.Topic(),
		registry.KeeperRegistryReorgedUpkeepReport{}.Topic(),
		registry.KeeperRegistryInsufficientFundsUpkeepReport{}.Topic(),
		registry.KeeperRegistryStaleUpkeepReport{}.Topic(),
	}

	return rp.p.RegisterFilter(logpoller.Filter{
		Name:      RegistryEventFilterName(rp.contract.Address()),
		EventSigs: append(upkeepStateEvents, upkeepActiveEvents...),
		Addresses: []common.Address{rp.addr},
	})
}

func (rp *RegistryPoller) unregister() {
	rp.p.UnregisterFilter(RegistryEventFilterName(rp.addr), nil)
}

// GetLatest is a convenience method to get the latest logs for the given event
// signatures with the provided lookback
func (rp *RegistryPoller) GetLatest(ctx context.Context, lookback int64, hashes ...common.Hash) (int64, []logpoller.Log, error) {
	end, err := rp.p.LatestBlock(pg.WithParentCtx(ctx))
	if err != nil {
		return 0, nil, fmt.Errorf("%w: %s", ErrFailedToGetLatestBlock, err)
	}

	logs, err := rp.p.LogsWithSigs(
		end-lookback,
		end,
		hashes,
		rp.addr,
		pg.WithParentCtx(ctx),
	)

	return end, logs, err
}

func RegistryEventFilterName(addr common.Address) string {
	return logpoller.FilterName("OCR2KeeperRegistry - Registry Events Filter", addr)
}
