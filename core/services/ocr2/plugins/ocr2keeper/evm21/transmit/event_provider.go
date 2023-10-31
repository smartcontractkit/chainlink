package transmit

import (
	"context"
	"encoding/hex"
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg/v3/types"

	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	iregistry21 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_keeper_registry_master_wrapper_2_1"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21/core"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

var _ ocr2keepers.TransmitEventProvider = &EventProvider{}

type logParser func(registry *iregistry21.IKeeperRegistryMaster, log logpoller.Log) (transmitEventLog, error)

type EventProvider struct {
	sync     utils.StartStopOnce
	mu       sync.RWMutex
	runState int
	runError error

	logger    logger.Logger
	logPoller logpoller.LogPoller
	registry  *iregistry21.IKeeperRegistryMaster
	client    evmclient.Client

	registryAddress common.Address
	lookbackBlocks  int64

	parseLog logParser
	cache    transmitEventCache
}

func EventProviderFilterName(addr common.Address) string {
	return logpoller.FilterName("KeepersRegistry TransmitEventProvider", addr)
}

func NewTransmitEventProvider(
	logger logger.Logger,
	logPoller logpoller.LogPoller,
	registryAddress common.Address,
	client evmclient.Client,
	lookbackBlocks int64,
) (*EventProvider, error) {
	var err error

	contract, err := iregistry21.NewIKeeperRegistryMaster(registryAddress, client)
	if err != nil {
		return nil, err
	}
	err = logPoller.RegisterFilter(logpoller.Filter{
		Name: EventProviderFilterName(contract.Address()),
		EventSigs: []common.Hash{
			// These are the events that are emitted when a node transmits a report
			iregistry21.IKeeperRegistryMasterUpkeepPerformed{}.Topic(),               // Happy path: report performed the upkeep
			iregistry21.IKeeperRegistryMasterReorgedUpkeepReport{}.Topic(),           // Report checkBlockNumber was reorged
			iregistry21.IKeeperRegistryMasterInsufficientFundsUpkeepReport{}.Topic(), // Upkeep didn't have sufficient funds when report reached chain, perform was aborted early
			// Report was too old when it reached the chain. For conditionals upkeep was already performed on a higher block than checkBlockNum
			// for logs upkeep was already performed for the particular log
			iregistry21.IKeeperRegistryMasterStaleUpkeepReport{}.Topic(),
		},
		Addresses: []common.Address{registryAddress},
	})
	if err != nil {
		return nil, err
	}

	return &EventProvider{
		logger:          logger,
		logPoller:       logPoller,
		registryAddress: registryAddress,
		lookbackBlocks:  lookbackBlocks,
		registry:        contract,
		client:          client,
		parseLog:        defaultLogParser,
		cache:           newTransmitEventCache(lookbackBlocks),
	}, nil
}

func (c *EventProvider) Name() string {
	return c.logger.Name()
}

func (c *EventProvider) Start(_ context.Context) error {
	return c.sync.StartOnce("AutomationTransmitEventProvider", func() error {
		c.mu.Lock()
		defer c.mu.Unlock()

		c.runState = 1
		return nil
	})
}

func (c *EventProvider) Close() error {
	return c.sync.StopOnce("AutomationRegistry", func() error {
		c.mu.Lock()
		defer c.mu.Unlock()

		c.runState = 0
		c.runError = nil
		return nil
	})
}

func (c *EventProvider) Ready() error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.runState == 1 {
		return nil
	}
	return c.sync.Ready()
}

func (c *EventProvider) HealthReport() map[string]error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.runState > 1 {
		c.sync.SvcErrBuffer.Append(fmt.Errorf("failed run state: %w", c.runError))
	}
	return map[string]error{c.Name(): c.sync.Healthy()}
}

func (c *EventProvider) GetLatestEvents(ctx context.Context) ([]ocr2keepers.TransmitEvent, error) {
	end, err := c.logPoller.LatestBlock(pg.WithParentCtx(ctx))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to get latest block from log poller", err)
	}

	// always check the last lookback number of blocks and rebroadcast
	// this allows the plugin to make decisions based on event confirmations
	logs, err := c.logPoller.LogsWithSigs(
		end.BlockNumber-c.lookbackBlocks,
		end.BlockNumber,
		[]common.Hash{
			iregistry21.IKeeperRegistryMasterUpkeepPerformed{}.Topic(),
			iregistry21.IKeeperRegistryMasterStaleUpkeepReport{}.Topic(),
			iregistry21.IKeeperRegistryMasterReorgedUpkeepReport{}.Topic(),
			iregistry21.IKeeperRegistryMasterInsufficientFundsUpkeepReport{}.Topic(),
		},
		c.registryAddress,
		pg.WithParentCtx(ctx),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to collect logs from log poller", err)
	}

	return c.processLogs(end.BlockNumber, logs...)
}

// processLogs will parse the unseen logs and return the corresponding transmit events.
func (c *EventProvider) processLogs(latestBlock int64, logs ...logpoller.Log) ([]ocr2keepers.TransmitEvent, error) {
	var vals []ocr2keepers.TransmitEvent

	for _, log := range logs {
		k := c.logKey(log)

		transmitEvent, ok := c.cache.get(ocr2keepers.BlockNumber(log.BlockNumber), k)
		if !ok {
			l, err := c.parseLog(c.registry, log)
			if err != nil {
				c.logger.Debugw("failed to parse log", "err", err)
				continue
			}
			id := l.Id()
			upkeepId := &ocr2keepers.UpkeepIdentifier{}
			ok := upkeepId.FromBigInt(id)
			if !ok {
				return nil, core.ErrInvalidUpkeepID
			}
			triggerW, err := core.UnpackTrigger(id, l.Trigger())
			if err != nil {
				return nil, fmt.Errorf("%w: failed to unpack trigger", err)
			}
			trigger := ocr2keepers.NewTrigger(
				ocr2keepers.BlockNumber(triggerW.BlockNum),
				triggerW.BlockHash,
			)
			switch core.GetUpkeepType(*upkeepId) {
			case ocr2keepers.LogTrigger:
				trigger.LogTriggerExtension = &ocr2keepers.LogTriggerExtension{}
				trigger.LogTriggerExtension.TxHash = triggerW.TxHash
				trigger.LogTriggerExtension.Index = triggerW.LogIndex
				trigger.LogTriggerExtension.BlockHash = triggerW.LogBlockHash
			default:
			}
			workID := core.UpkeepWorkID(*upkeepId, trigger)
			transmitEvent = ocr2keepers.TransmitEvent{
				Type:            l.TransmitEventType(),
				TransmitBlock:   ocr2keepers.BlockNumber(l.BlockNumber),
				TransactionHash: l.TxHash,
				WorkID:          workID,
				UpkeepID:        *upkeepId,
				CheckBlock:      trigger.BlockNumber,
			}
			c.cache.add(k, transmitEvent)
		}

		transmitEvent.Confirmations = latestBlock - int64(transmitEvent.TransmitBlock)

		vals = append(vals, transmitEvent)
	}

	return vals, nil
}

func (c *EventProvider) logKey(log logpoller.Log) string {
	logExt := ocr2keepers.LogTriggerExtension{
		TxHash:    log.TxHash,
		Index:     uint32(log.LogIndex),
		BlockHash: log.BlockHash,
	}
	logId := logExt.LogIdentifier()
	return hex.EncodeToString(logId)
}
