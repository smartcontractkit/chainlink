package evm

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	pluginchain "github.com/smartcontractkit/ocr2keepers/pkg/chain"
	plugintypes "github.com/smartcontractkit/ocr2keepers/pkg/types"
	pluginutils "github.com/smartcontractkit/ocr2keepers/pkg/util"

	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated"
	registry "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_wrapper2_0"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

var (
	ErrFailedToGetLogs           = fmt.Errorf("failed to collect logs from log poller")
	ErrFailedToUnmarshalLogs     = fmt.Errorf("failed to unmarshal logs")
	ErrFailedToExtractCheckBlock = fmt.Errorf("error while fetching checkBlockNumber from report log")
)

type TransmitUnpacker interface {
	UnpackTransmitTxInput([]byte) ([]plugintypes.UpkeepResult, error)
}

type LogProvider struct {
	// provided dependencies
	rp             *RegistryPoller
	logger         logger.Logger
	client         evmclient.Client
	lookbackBlocks int64

	// properties initialized in constructor
	registry          *registry.KeeperRegistry
	packer            TransmitUnpacker
	txCheckBlockCache *pluginutils.Cache[string]
	cacheCleaner      *pluginutils.IntervalCacheCleaner[string]

	// run state properties
	sync     utils.StartStopOnce
	mu       sync.RWMutex
	runState int
	runError error
}

var _ plugintypes.PerformLogProvider = (*LogProvider)(nil)

func LogProviderFilterName(addr common.Address) string {
	return logpoller.FilterName("OCR2KeeperRegistry - LogProvider", addr)
}

func NewLogProvider(
	logger logger.Logger,
	rp *RegistryPoller,
	client evmclient.Client,
	lookbackBlocks int64,
) (*LogProvider, error) {
	var err error

	contract, err := registry.NewKeeperRegistry(common.HexToAddress("0x"), client)
	if err != nil {
		return nil, err
	}

	abi, err := abi.JSON(strings.NewReader(registry.KeeperRegistryABI))
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrABINotParsable, err)
	}

	return &LogProvider{
		rp:                rp,
		logger:            logger,
		client:            client,
		lookbackBlocks:    lookbackBlocks,
		registry:          contract,
		packer:            NewEvmRegistryPackerV2_0(abi),
		txCheckBlockCache: pluginutils.NewCache[string](time.Hour),
		cacheCleaner:      pluginutils.NewIntervalCacheCleaner[string](time.Minute),
	}, nil
}

// Name implements the job.ServiceCtx interface
func (c *LogProvider) Name() string {
	return c.logger.Name()
}

// Start implements the job.ServiceCtx interface
func (c *LogProvider) Start(ctx context.Context) error {
	return c.sync.StartOnce("AutomationLogProvider", func() error {
		c.mu.Lock()
		defer c.mu.Unlock()

		go c.cacheCleaner.Run(c.txCheckBlockCache)

		c.runState = 1

		return nil
	})
}

// Stop implements the job.ServiceCtx interface
func (c *LogProvider) Close() error {
	return c.sync.StopOnce("AutomationRegistry", func() error {
		c.mu.Lock()
		defer c.mu.Unlock()

		c.cacheCleaner.Stop()
		c.runState = 0
		c.runError = nil

		return nil
	})
}

// Ready implements the job.ServiceCtx interface
func (c *LogProvider) Ready() error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.runState == 1 {
		return nil
	}

	return c.sync.Ready()
}

// HealthReport implements the job.ServiceCtx interface
func (c *LogProvider) HealthReport() map[string]error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.runState > 1 {
		c.sync.SvcErrBuffer.Append(fmt.Errorf("failed run state: %w", c.runError))
	}

	return map[string]error{c.Name(): c.sync.Healthy()}
}

// PerformLogs provides a list of logs that indicate an upkeep was performed
func (c *LogProvider) PerformLogs(ctx context.Context) ([]plugintypes.PerformLog, error) {
	var (
		end  int64
		logs []logpoller.Log
		err  error
	)

	// always check the last lookback number of blocks and rebroadcast
	// this allows the plugin to make decisions based on event confirmations
	if end, logs, err = c.rp.GetLatest(
		ctx,
		c.lookbackBlocks,
		registry.KeeperRegistryUpkeepPerformed{}.Topic(),
	); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrFailedToGetLogs, err)
	}

	// maximum array size will be the length of the logs found in the query
	// above
	results := make([]plugintypes.PerformLog, 0, len(logs))

	for _, pollerLog := range logs {
		rawLog := pollerLog.ToGethLog()
		abilog, err := c.registry.ParseLog(rawLog)

		if err != nil {
			return nil, fmt.Errorf("%w (perform): %s", ErrFailedToUnmarshalLogs, err)
		}

		// ensure that the logs are of the expected type as returned from the
		// query above
		switch typedLog := abilog.(type) {
		case *registry.KeeperRegistryUpkeepPerformed:
			if typedLog == nil {
				continue
			}

			results = append(results, plugintypes.PerformLog{
				Key:             pluginchain.NewUpkeepKey(big.NewInt(int64(typedLog.CheckBlockNumber)), typedLog.Id),
				TransmitBlock:   pluginchain.BlockKey([]byte(fmt.Sprintf("%d", pollerLog.BlockNumber))),
				TransactionHash: rawLog.TxHash.Hex(),
				Confirmations:   end - pollerLog.BlockNumber,
			})
		}
	}

	// return log to subscribers
	return results, nil
}

// StaleReportLogs provides a list of logs that would indicate a report or an
// upkeep in a report might be stale. Use these logs to help unblock 'in-flight'
// upkeeps.
func (c *LogProvider) StaleReportLogs(ctx context.Context) ([]plugintypes.StaleReportLog, error) {
	var (
		end  int64
		logs []logpoller.Log
		err  error
	)

	// always check the last lookback number of blocks and rebroadcast
	// this allows the plugin to make decisions based on event confirmations

	// collect all log types in a single db request
	if end, logs, err = c.rp.GetLatest(
		ctx,
		c.lookbackBlocks,
		registry.KeeperRegistryReorgedUpkeepReport{}.Topic(),
		registry.KeeperRegistryStaleUpkeepReport{}.Topic(),
		registry.KeeperRegistryInsufficientFundsUpkeepReport{}.Topic(),
	); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrFailedToGetLogs, err)
	}

	return c.makeStaleReportLogs(ctx, end, logs)
}

func (c *LogProvider) makeStaleReportLogs(ctx context.Context, end int64, logs []logpoller.Log) ([]plugintypes.StaleReportLog, error) {
	results := make([]plugintypes.StaleReportLog, 0, len(logs))

	var (
		rawLog     gethtypes.Log
		abigenLog  generated.AbigenLog
		id         []byte
		name       string
		checkBlock plugintypes.BlockKey
		err        error
	)

	for _, pollerLog := range logs {
		rawLog = pollerLog.ToGethLog()

		if abigenLog, err = c.registry.ParseLog(rawLog); err != nil {
			return results, err
		}

		if name, id, err = idFromLog(abigenLog); err != nil {
			continue
		}

		if checkBlock, err = c.getCheckBlockNumberFromTxHash(ctx, pollerLog.TxHash, id); err != nil {
			c.logger.Error("%w (%s): %s", ErrFailedToExtractCheckBlock, name, err)
			continue
		}

		results = append(results, plugintypes.StaleReportLog{
			Key:             pluginchain.NewUpkeepKeyFromBlockAndID(checkBlock, id),
			TransmitBlock:   pluginchain.BlockKey([]byte(fmt.Sprintf("%d", pollerLog.BlockNumber))),
			TransactionHash: pollerLog.TxHash.Hex(),
			Confirmations:   end - pollerLog.BlockNumber,
		})
	}

	return results, nil
}

func idFromLog(abigenLog generated.AbigenLog) (string, []byte, error) {
	var (
		id   []byte
		name string
	)

	switch l := abigenLog.(type) {
	case *registry.KeeperRegistryInsufficientFundsUpkeepReport:
		if l == nil {
			return name, id, fmt.Errorf("nil log")
		}

		name = "insufficient funds"
		id = plugintypes.UpkeepIdentifier(l.Id.String())
	case *registry.KeeperRegistryReorgedUpkeepReport:
		if l == nil {
			return name, id, fmt.Errorf("nil log")
		}

		name = "reorged"
		id = plugintypes.UpkeepIdentifier(l.Id.String())
	case *registry.KeeperRegistryStaleUpkeepReport:
		if l == nil {
			return name, id, fmt.Errorf("nil log")
		}

		name = "stale upkeep"
		id = plugintypes.UpkeepIdentifier(l.Id.String())
	}

	return name, id, nil
}

// Fetches the checkBlockNumber for a particular transaction and an upkeep ID. Requires a RPC call to get txData
// so this function should not be used heavily
func (c *LogProvider) getCheckBlockNumberFromTxHash(ctx context.Context, txHash common.Hash, id plugintypes.UpkeepIdentifier) (bk plugintypes.BlockKey, e error) {
	defer func() {
		if r := recover(); r != nil {
			e = fmt.Errorf("recovered from panic in getCheckBlockNumberForUpkeep: %v", r)
		}
	}()

	// the cacheKey is the transaction hash plus the upkeep id
	cacheKey := txHash.String() + "|" + string(id)

	// Check if value already exists in cache for txHash, id pair
	if val, ok := c.txCheckBlockCache.Get(cacheKey); ok {
		return pluginchain.BlockKey(val), nil
	}

	var (
		tx  gethtypes.Transaction
		err error
	)

	// make an RPC call to get the transaction
	if err := c.client.CallContext(ctx, &tx, "eth_getTransactionByHash", txHash); err != nil {
		return nil, err
	}

	// get the data from the transaction and verify the data returned is more
	// than just a signature (the first 4 bytes)
	txData := tx.Data()
	if len(txData) < 4 {
		return nil, fmt.Errorf("error in getCheckBlockNumberForUpkeep, got invalid tx data %s", txData)
	}

	// decode the data into an array of upkeep results
	decodedReport, err := c.packer.UnpackTransmitTxInput(txData[4:]) // Remove first 4 bytes of function signature
	if err != nil {
		return nil, err
	}

	for _, upkeep := range decodedReport {
		bl, ui, err := upkeep.Key.BlockKeyAndUpkeepID()
		if err != nil {
			return nil, err
		}

		if string(ui) == string(id) {
			c.txCheckBlockCache.Set(cacheKey, bl.String(), pluginutils.DefaultCacheExpiration)
			return bl, nil
		}
	}

	return nil, fmt.Errorf("upkeep %s not found in tx hash %s", id, txHash)
}
