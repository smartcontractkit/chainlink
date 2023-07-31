package evm

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg"

	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	iregistry21 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_keeper_registry_master_wrapper_2_1"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type TransmitEventProvider struct {
	sync     utils.StartStopOnce
	mu       sync.RWMutex
	runState int
	runError error

	logger          logger.Logger
	logPoller       logpoller.LogPoller
	registryAddress common.Address
	lookbackBlocks  int64
	registry        *iregistry21.IKeeperRegistryMaster
	client          evmclient.Client
}

func TransmitEventProviderFilterName(addr common.Address) string {
	return logpoller.FilterName("KeepersRegistry TransmitEventProvider", addr)
}

func NewTransmitEventProvider(
	logger logger.Logger,
	logPoller logpoller.LogPoller,
	registryAddress common.Address,
	client evmclient.Client,
	lookbackBlocks int64,
) (*TransmitEventProvider, error) {
	var err error

	contract, err := iregistry21.NewIKeeperRegistryMaster(registryAddress, client)
	if err != nil {
		return nil, err
	}
	// Add log filters for the log poller so that it can poll and find the logs that
	// we need.
	err = logPoller.RegisterFilter(logpoller.Filter{
		Name: TransmitEventProviderFilterName(contract.Address()),
		EventSigs: []common.Hash{
			iregistry21.IKeeperRegistryMasterUpkeepPerformed{}.Topic(),
			iregistry21.IKeeperRegistryMasterReorgedUpkeepReport{}.Topic(),
			iregistry21.IKeeperRegistryMasterInsufficientFundsUpkeepReport{}.Topic(),
			iregistry21.IKeeperRegistryMasterStaleUpkeepReport{}.Topic(),
		},
		Addresses: []common.Address{registryAddress},
	})
	if err != nil {
		return nil, err
	}

	return &TransmitEventProvider{
		logger:          logger,
		logPoller:       logPoller,
		registryAddress: registryAddress,
		lookbackBlocks:  lookbackBlocks,
		registry:        contract,
		client:          client,
	}, nil
}

func (c *TransmitEventProvider) Name() string {
	return c.logger.Name()
}

func (c *TransmitEventProvider) Start(ctx context.Context) error {
	return c.sync.StartOnce("AutomationTransmitEventProvider", func() error {
		c.mu.Lock()
		defer c.mu.Unlock()

		c.runState = 1
		return nil
	})
}

func (c *TransmitEventProvider) Close() error {
	return c.sync.StopOnce("AutomationRegistry", func() error {
		c.mu.Lock()
		defer c.mu.Unlock()

		c.runState = 0
		c.runError = nil
		return nil
	})
}

func (c *TransmitEventProvider) Ready() error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.runState == 1 {
		return nil
	}
	return c.sync.Ready()
}

func (c *TransmitEventProvider) HealthReport() map[string]error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.runState > 1 {
		c.sync.SvcErrBuffer.Append(fmt.Errorf("failed run state: %w", c.runError))
	}
	return map[string]error{c.Name(): c.sync.Healthy()}
}

func (c *TransmitEventProvider) Events(ctx context.Context) ([]ocr2keepers.TransmitEvent, error) {
	end, err := c.logPoller.LatestBlock(pg.WithParentCtx(ctx))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to get latest block from log poller", err)
	}

	// always check the last lookback number of blocks and rebroadcast
	// this allows the plugin to make decisions based on event confirmations
	logs, err := c.logPoller.LogsWithSigs(
		end-c.lookbackBlocks,
		end,
		[]common.Hash{
			iregistry21.IKeeperRegistryMasterUpkeepPerformed{}.Topic(),
			iregistry21.IKeeperRegistryMasterStaleUpkeepReport{}.Topic(),
			// TODO: enable once we have the corredponding types in ocr2keepers
			// iregistry21.IKeeperRegistryMasterReorgedUpkeepReport{}.Topic(),
			// iregistry21.IKeeperRegistryMasterInsufficientFundsUpkeepReport{}.Topic(),
		},
		c.registryAddress,
		pg.WithParentCtx(ctx),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to collect logs from log poller", err)
	}

	parsed, err := c.parseLogs(logs)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse logs", err)
	}

	return c.convertToTransmitEvents(parsed, end)
}

func (c *TransmitEventProvider) parseLogs(logs []logpoller.Log) ([]transmitEventLog, error) {
	results := []transmitEventLog{}

	for _, log := range logs {
		rawLog := log.ToGethLog()
		abilog, err := c.registry.ParseLog(rawLog)
		if err != nil {
			return nil, fmt.Errorf("%w: failed to parse log", err)
		}

		switch l := abilog.(type) {
		case *iregistry21.IKeeperRegistryMasterUpkeepPerformed:
			if l == nil {
				continue
			}
			results = append(results, transmitEventLog{
				Log:       log,
				Performed: l,
			})
		case *iregistry21.IKeeperRegistryMasterReorgedUpkeepReport:
			if l == nil {
				continue
			}
			results = append(results, transmitEventLog{
				Log:     log,
				Reorged: l,
			})
		case *iregistry21.IKeeperRegistryMasterStaleUpkeepReport:
			if l == nil {
				continue
			}
			results = append(results, transmitEventLog{
				Log:   log,
				Stale: l,
			})
		case *iregistry21.IKeeperRegistryMasterInsufficientFundsUpkeepReport:
			if l == nil {
				continue
			}
			results = append(results, transmitEventLog{
				Log:               log,
				InsufficientFunds: l,
			})
		default:
			c.logger.Debugw("skipping unknown log type", "l", l)
			continue
		}
	}

	return results, nil
}

func (c *TransmitEventProvider) convertToTransmitEvents(logs []transmitEventLog, latestBlock int64) ([]ocr2keepers.TransmitEvent, error) {
	vals := []ocr2keepers.TransmitEvent{}

	for _, l := range logs {
		upkeepId := ocr2keepers.UpkeepIdentifier(l.Id().Bytes())
		triggerID := l.TriggerID()
		vals = append(vals, ocr2keepers.TransmitEvent{
			Type:            l.TransmitEventType(),
			TransmitBlock:   BlockKeyHelper[int64]{}.MakeBlockKey(l.BlockNumber),
			Confirmations:   latestBlock - l.BlockNumber,
			TransactionHash: l.TxHash.Hex(),
			ID:              hex.EncodeToString(triggerID[:]),
			UpkeepID:        upkeepId,
		})
	}

	return vals, nil
}

// transmitEventLog is a wrapper around logpoller.Log and the parsed log
type transmitEventLog struct {
	logpoller.Log
	Performed         *iregistry21.IKeeperRegistryMasterUpkeepPerformed
	Stale             *iregistry21.IKeeperRegistryMasterStaleUpkeepReport
	Reorged           *iregistry21.IKeeperRegistryMasterReorgedUpkeepReport
	InsufficientFunds *iregistry21.IKeeperRegistryMasterInsufficientFundsUpkeepReport
}

func (l transmitEventLog) Id() *big.Int {
	switch {
	case l.Performed != nil:
		return l.Performed.Id
	case l.Stale != nil:
		return l.Stale.Id
	case l.Reorged != nil:
		return l.Reorged.Id
	case l.InsufficientFunds != nil:
		return l.InsufficientFunds.Id
	default:
		return nil
	}
}

func (l transmitEventLog) TriggerID() [32]byte {
	switch {
	case l.Performed != nil:
		return l.Performed.UpkeepTriggerID
	case l.Stale != nil:
		return l.Stale.UpkeepTriggerID
	case l.Reorged != nil:
		return l.Reorged.UpkeepTriggerID
	case l.InsufficientFunds != nil:
		return l.InsufficientFunds.UpkeepTriggerID
	default:
		return [32]byte{}
	}
}

func (l transmitEventLog) TransmitEventType() ocr2keepers.TransmitEventType {
	switch {
	case l.Performed != nil:
		return ocr2keepers.PerformEvent
	case l.Stale != nil:
		return ocr2keepers.StaleReportEvent
	case l.Reorged != nil:
		// TODO: use reorged event type
		return ocr2keepers.TransmitEventType(2)
	case l.InsufficientFunds != nil:
		// TODO: use insufficient funds event type
		return ocr2keepers.TransmitEventType(3)
	default:
		return ocr2keepers.TransmitEventType(0)
	}
}
