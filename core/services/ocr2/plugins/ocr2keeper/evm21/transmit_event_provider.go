package evm

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg"
	pluginutils "github.com/smartcontractkit/ocr2keepers/pkg/util"

	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	iregistry21 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_keeper_registry_master_wrapper_2_1"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_log_automation"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type TransmitUnpacker interface {
	UnpackTransmitTxInput([]byte) ([]ocr2keepers.UpkeepResult, error)
}

type TransmitEventProvider struct {
	sync              utils.StartStopOnce
	mu                sync.RWMutex
	runState          int
	runError          error
	logger            logger.Logger
	logPoller         logpoller.LogPoller
	registryAddress   common.Address
	lookbackBlocks    int64
	registry          *iregistry21.IKeeperRegistryMaster
	client            evmclient.Client
	packer            TransmitUnpacker
	txCheckBlockCache *pluginutils.Cache[string]
	cacheCleaner      *pluginutils.IntervalCacheCleaner[string]
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

	contract, err := iregistry21.NewIKeeperRegistryMaster(common.HexToAddress("0x"), client)
	if err != nil {
		return nil, err
	}

	keeperABI, err := abi.JSON(strings.NewReader(iregistry21.IKeeperRegistryMasterABI))
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrABINotParsable, err)
	}
	logDataABI, err := abi.JSON(strings.NewReader(i_log_automation.ILogAutomationABI))
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrABINotParsable, err)
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
		logger:            logger,
		logPoller:         logPoller,
		registryAddress:   registryAddress,
		lookbackBlocks:    lookbackBlocks,
		registry:          contract,
		client:            client,
		packer:            NewEvmRegistryPackerV2_1(keeperABI, logDataABI),
		txCheckBlockCache: pluginutils.NewCache[string](time.Hour),
		cacheCleaner:      pluginutils.NewIntervalCacheCleaner[string](time.Minute),
	}, nil
}

func (c *TransmitEventProvider) Name() string {
	return c.logger.Name()
}

func (c *TransmitEventProvider) Start(ctx context.Context) error {
	return c.sync.StartOnce("AutomationTransmitEventProvider", func() error {
		c.mu.Lock()
		defer c.mu.Unlock()

		go c.cacheCleaner.Run(c.txCheckBlockCache)
		c.runState = 1
		return nil
	})
}

func (c *TransmitEventProvider) Close() error {
	return c.sync.StopOnce("AutomationRegistry", func() error {
		c.mu.Lock()
		defer c.mu.Unlock()

		c.cacheCleaner.Stop()
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

	performed, err := c.unmarshalTransmitLogs(logs)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to unmarshal logs", err)
	}

	return c.convertToTransmitEvents(performed, end)
}

func (c *TransmitEventProvider) unmarshalTransmitLogs(logs []logpoller.Log) ([]transmitEventLog, error) {
	results := []transmitEventLog{}

	for _, log := range logs {
		rawLog := log.ToGethLog()
		abilog, err := c.registry.ParseLog(rawLog)
		if err != nil {
			return results, err
		}

		switch l := abilog.(type) {
		case *iregistry21.IKeeperRegistryMasterUpkeepPerformed:
			if l == nil {
				continue
			}

			r := transmitEventLog{
				Log:       log,
				Performed: l,
			}

			results = append(results, r)
		case *iregistry21.IKeeperRegistryMasterReorgedUpkeepReport:
			if l == nil {
				continue
			}

			r := transmitEventLog{
				Log:     log,
				Reorged: l,
			}

			results = append(results, r)
		case *iregistry21.IKeeperRegistryMasterStaleUpkeepReport:
			if l == nil {
				continue
			}

			r := transmitEventLog{
				Log:   log,
				Stale: l,
			}

			results = append(results, r)
		case *iregistry21.IKeeperRegistryMasterInsufficientFundsUpkeepReport:
			if l == nil {
				continue
			}

			r := transmitEventLog{
				Log:               log,
				InsufficientFunds: l,
			}

			results = append(results, r)
		}
	}

	return results, nil
}

func (c *TransmitEventProvider) convertToTransmitEvents(logs []transmitEventLog, latestBlock int64) ([]ocr2keepers.TransmitEvent, error) {
	var err error
	vals := []ocr2keepers.TransmitEvent{}

	for _, l := range logs {
		var checkBlockNumber ocr2keepers.BlockKey
		upkeepId := ocr2keepers.UpkeepIdentifier(l.Id().Bytes())
		switch getUpkeepType(upkeepId) {
		case conditionTrigger:
			checkBlockNumber, err = c.getCheckBlockNumberFromTxHash(l.TxHash, upkeepId)
			if err != nil {
				c.logger.Error("error while fetching checkBlockNumber from perform report log: %w", err)
				continue
			}
		default:
		}
		triggerID := l.TriggerID()
		vals = append(vals, ocr2keepers.TransmitEvent{
			Type:            l.TransmitEventType(),
			TransmitBlock:   BlockKeyHelper[int64]{}.MakeBlockKey(l.BlockNumber),
			Confirmations:   latestBlock - l.BlockNumber,
			TransactionHash: l.TxHash.Hex(),
			ID:              hex.EncodeToString(triggerID[:]),
			UpkeepID:        upkeepId,
			CheckBlock:      checkBlockNumber,
		})
	}

	return vals, nil
}

// Fetches the checkBlockNumber for a particular transaction and an upkeep ID. Requires a RPC call to get txData
// so this function should not be used heavily
func (c *TransmitEventProvider) getCheckBlockNumberFromTxHash(txHash common.Hash, id ocr2keepers.UpkeepIdentifier) (bk ocr2keepers.BlockKey, e error) {
	defer func() {
		if r := recover(); r != nil {
			e = fmt.Errorf("recovered from panic in getCheckBlockNumberForUpkeep: %v", r)
		}
	}()

	// Check if value already exists in cache for txHash, id pair
	cacheKey := txHash.String() + "|" + string(id)
	if val, ok := c.txCheckBlockCache.Get(cacheKey); ok {
		return ocr2keepers.BlockKey(val), nil
	}

	var tx gethtypes.Transaction
	err := c.client.CallContext(context.Background(), &tx, "eth_getTransactionByHash", txHash)
	if err != nil {
		return "", err
	}

	txData := tx.Data()
	if len(txData) < 4 {
		return "", fmt.Errorf("error in getCheckBlockNumberForUpkeep, got invalid tx data %s", txData)
	}

	decodedReport, err := c.packer.UnpackTransmitTxInput(txData[4:]) // Remove first 4 bytes of function signature
	if err != nil {
		return "", err
	}

	for _, upkeep := range decodedReport {
		res, ok := upkeep.(EVMAutomationUpkeepResult21)
		if !ok {
			return "", fmt.Errorf("unexpected type")
		}

		if res.ID.String() == string(id) {
			bl := fmt.Sprintf("%d", res.Block)

			c.txCheckBlockCache.Set(cacheKey, bl, pluginutils.DefaultCacheExpiration)

			return ocr2keepers.BlockKey(bl), nil
		}
	}

	return "", fmt.Errorf("upkeep %s not found in tx hash %s", id, txHash)
}

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
		return l.Performed.TriggerID
	case l.Stale != nil:
		return l.Stale.TriggerID
	case l.Reorged != nil:
		return l.Reorged.TriggerID
	case l.InsufficientFunds != nil:
		return l.InsufficientFunds.TriggerID
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
