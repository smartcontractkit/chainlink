package evm

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"

	pluginutils "github.com/smartcontractkit/chainlink-automation/pkg/util"
	ocr2keepers "github.com/smartcontractkit/chainlink-automation/pkg/v2"
	"github.com/smartcontractkit/chainlink-automation/pkg/v2/encoding"

	"github.com/smartcontractkit/chainlink-common/pkg/services"

	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	registry "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_wrapper2_0"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type TransmitUnpacker interface {
	UnpackTransmitTxInput([]byte) ([]ocr2keepers.UpkeepResult, error)
}

type LogProvider struct {
	sync              services.StateMachine
	mu                sync.RWMutex
	runState          int
	runError          error
	logger            logger.Logger
	logPoller         logpoller.LogPoller
	registryAddress   common.Address
	lookbackBlocks    int64
	registry          *registry.KeeperRegistry
	client            evmclient.Client
	packer            TransmitUnpacker
	txCheckBlockCache *pluginutils.Cache[string]
	cacheCleaner      *pluginutils.IntervalCacheCleaner[string]
}

func LogProviderFilterName(addr common.Address) string {
	return logpoller.FilterName("OCR2KeeperRegistry - LogProvider", addr)
}

func NewLogProvider(
	ctx context.Context,
	logger logger.Logger,
	logPoller logpoller.LogPoller,
	registryAddress common.Address,
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

	// Add log filters for the log poller so that it can poll and find the logs that
	// we need.
	err = logPoller.RegisterFilter(ctx, logpoller.Filter{
		Name: LogProviderFilterName(contract.Address()),
		EventSigs: []common.Hash{
			registry.KeeperRegistryUpkeepPerformed{}.Topic(),
			registry.KeeperRegistryReorgedUpkeepReport{}.Topic(),
			registry.KeeperRegistryInsufficientFundsUpkeepReport{}.Topic(),
			registry.KeeperRegistryStaleUpkeepReport{}.Topic(),
		},
		Addresses: []common.Address{registryAddress},
	})
	if err != nil {
		return nil, err
	}

	return &LogProvider{
		logger:            logger.Named("AutomationLogProvider"),
		logPoller:         logPoller,
		registryAddress:   registryAddress,
		lookbackBlocks:    lookbackBlocks,
		registry:          contract,
		client:            client,
		packer:            NewEvmRegistryPackerV2_0(abi),
		txCheckBlockCache: pluginutils.NewCache[string](time.Hour),
		cacheCleaner:      pluginutils.NewIntervalCacheCleaner[string](time.Minute),
	}, nil
}

func (c *LogProvider) Name() string {
	return c.logger.Name()
}

func (c *LogProvider) Start(ctx context.Context) error {
	return c.sync.StartOnce("AutomationLogProvider", func() error {
		c.mu.Lock()
		defer c.mu.Unlock()

		go c.cacheCleaner.Run(c.txCheckBlockCache)
		c.runState = 1
		return nil
	})
}

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

func (c *LogProvider) Ready() error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.runState == 1 {
		return nil
	}
	return c.sync.Ready()
}

func (c *LogProvider) HealthReport() map[string]error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.runState > 1 {
		c.sync.SvcErrBuffer.Append(fmt.Errorf("failed run state: %w", c.runError))
	}
	return map[string]error{c.Name(): c.sync.Healthy()}
}

func (c *LogProvider) PerformLogs(ctx context.Context) ([]ocr2keepers.PerformLog, error) {
	end, err := c.logPoller.LatestBlock(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to get latest block from log poller", err)
	}

	// always check the last lookback number of blocks and rebroadcast
	// this allows the plugin to make decisions based on event confirmations
	logs, err := c.logPoller.LogsWithSigs(
		ctx,
		end.BlockNumber-c.lookbackBlocks,
		end.BlockNumber,
		[]common.Hash{
			registry.KeeperRegistryUpkeepPerformed{}.Topic(),
		},
		c.registryAddress,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to collect logs from log poller", err)
	}

	performed, err := c.unmarshalPerformLogs(logs)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to unmarshal logs", err)
	}

	vals := []ocr2keepers.PerformLog{}
	for _, p := range performed {
		// broadcast log to subscribers
		l := ocr2keepers.PerformLog{
			Key:             UpkeepKeyHelper[uint32]{}.MakeUpkeepKey(p.CheckBlockNumber, p.Id),
			TransmitBlock:   BlockKeyHelper[int64]{}.MakeBlockKey(p.BlockNumber),
			TransactionHash: p.TxHash.Hex(),
			Confirmations:   end.BlockNumber - p.BlockNumber,
		}
		vals = append(vals, l)
	}

	return vals, nil
}

func (c *LogProvider) StaleReportLogs(ctx context.Context) ([]ocr2keepers.StaleReportLog, error) {
	end, err := c.logPoller.LatestBlock(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to get latest block from log poller", err)
	}

	// always check the last lookback number of blocks and rebroadcast
	// this allows the plugin to make decisions based on event confirmations

	// ReorgedUpkeepReportLogs
	logs, err := c.logPoller.LogsWithSigs(
		ctx,
		end.BlockNumber-c.lookbackBlocks,
		end.BlockNumber,
		[]common.Hash{
			registry.KeeperRegistryReorgedUpkeepReport{}.Topic(),
		},
		c.registryAddress,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to collect logs from log poller", err)
	}
	reorged, err := c.unmarshalReorgUpkeepLogs(logs)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to unmarshal reorg logs", err)
	}

	// StaleUpkeepReportLogs
	logs, err = c.logPoller.LogsWithSigs(
		ctx,
		end.BlockNumber-c.lookbackBlocks,
		end.BlockNumber,
		[]common.Hash{
			registry.KeeperRegistryStaleUpkeepReport{}.Topic(),
		},
		c.registryAddress,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to collect logs from log poller", err)
	}
	staleUpkeep, err := c.unmarshalStaleUpkeepLogs(logs)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to unmarshal stale upkeep logs", err)
	}

	// InsufficientFundsUpkeepReportLogs
	logs, err = c.logPoller.LogsWithSigs(
		ctx,
		end.BlockNumber-c.lookbackBlocks,
		end.BlockNumber,
		[]common.Hash{
			registry.KeeperRegistryInsufficientFundsUpkeepReport{}.Topic(),
		},
		c.registryAddress,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to collect logs from log poller", err)
	}
	insufficientFunds, err := c.unmarshalInsufficientFundsUpkeepLogs(logs)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to unmarshal insufficient fund upkeep logs", err)
	}

	vals := []ocr2keepers.StaleReportLog{}
	for _, r := range reorged {
		upkeepId := ocr2keepers.UpkeepIdentifier(r.Id.String())
		checkBlockNumber, err := c.getCheckBlockNumberFromTxHash(r.TxHash, upkeepId)
		if err != nil {
			c.logger.Error("error while fetching checkBlockNumber from reorged report log: %w", err)
			continue
		}
		l := ocr2keepers.StaleReportLog{
			Key:             encoding.BasicEncoder{}.MakeUpkeepKey(checkBlockNumber, upkeepId),
			TransmitBlock:   BlockKeyHelper[int64]{}.MakeBlockKey(r.BlockNumber),
			TransactionHash: r.TxHash.Hex(),
			Confirmations:   end.BlockNumber - r.BlockNumber,
		}
		vals = append(vals, l)
	}
	for _, r := range staleUpkeep {
		upkeepId := ocr2keepers.UpkeepIdentifier(r.Id.String())
		checkBlockNumber, err := c.getCheckBlockNumberFromTxHash(r.TxHash, upkeepId)
		if err != nil {
			c.logger.Error("error while fetching checkBlockNumber from stale report log: %w", err)
			continue
		}
		l := ocr2keepers.StaleReportLog{
			Key:             encoding.BasicEncoder{}.MakeUpkeepKey(checkBlockNumber, upkeepId),
			TransmitBlock:   BlockKeyHelper[int64]{}.MakeBlockKey(r.BlockNumber),
			TransactionHash: r.TxHash.Hex(),
			Confirmations:   end.BlockNumber - r.BlockNumber,
		}
		vals = append(vals, l)
	}
	for _, r := range insufficientFunds {
		upkeepId := ocr2keepers.UpkeepIdentifier(r.Id.String())
		checkBlockNumber, err := c.getCheckBlockNumberFromTxHash(r.TxHash, upkeepId)
		if err != nil {
			c.logger.Error("error while fetching checkBlockNumber from insufficient funds report log: %w", err)
			continue
		}
		l := ocr2keepers.StaleReportLog{
			Key:             encoding.BasicEncoder{}.MakeUpkeepKey(checkBlockNumber, upkeepId),
			TransmitBlock:   BlockKeyHelper[int64]{}.MakeBlockKey(r.BlockNumber),
			TransactionHash: r.TxHash.Hex(),
			Confirmations:   end.BlockNumber - r.BlockNumber,
		}
		vals = append(vals, l)
	}

	return vals, nil
}

func (c *LogProvider) unmarshalPerformLogs(logs []logpoller.Log) ([]performed, error) {
	results := []performed{}

	for _, log := range logs {
		rawLog := log.ToGethLog()
		abilog, err := c.registry.ParseLog(rawLog)
		if err != nil {
			return results, err
		}

		switch l := abilog.(type) {
		case *registry.KeeperRegistryUpkeepPerformed:
			if l == nil {
				continue
			}

			r := performed{
				Log:                           log,
				KeeperRegistryUpkeepPerformed: *l,
			}

			results = append(results, r)
		}
	}

	return results, nil
}

func (c *LogProvider) unmarshalReorgUpkeepLogs(logs []logpoller.Log) ([]reorged, error) {
	results := []reorged{}

	for _, log := range logs {
		rawLog := log.ToGethLog()
		abilog, err := c.registry.ParseLog(rawLog)
		if err != nil {
			return results, err
		}

		switch l := abilog.(type) {
		case *registry.KeeperRegistryReorgedUpkeepReport:
			if l == nil {
				continue
			}

			r := reorged{
				Log:                               log,
				KeeperRegistryReorgedUpkeepReport: *l,
			}

			results = append(results, r)
		}
	}

	return results, nil
}

func (c *LogProvider) unmarshalStaleUpkeepLogs(logs []logpoller.Log) ([]staleUpkeep, error) {
	results := []staleUpkeep{}

	for _, log := range logs {
		rawLog := log.ToGethLog()
		abilog, err := c.registry.ParseLog(rawLog)
		if err != nil {
			return results, err
		}

		switch l := abilog.(type) {
		case *registry.KeeperRegistryStaleUpkeepReport:
			if l == nil {
				continue
			}

			r := staleUpkeep{
				Log:                             log,
				KeeperRegistryStaleUpkeepReport: *l,
			}

			results = append(results, r)
		}
	}

	return results, nil
}

func (c *LogProvider) unmarshalInsufficientFundsUpkeepLogs(logs []logpoller.Log) ([]insufficientFunds, error) {
	results := []insufficientFunds{}

	for _, log := range logs {
		rawLog := log.ToGethLog()
		abilog, err := c.registry.ParseLog(rawLog)
		if err != nil {
			return results, err
		}

		switch l := abilog.(type) {
		case *registry.KeeperRegistryInsufficientFundsUpkeepReport:
			if l == nil {
				continue
			}

			r := insufficientFunds{
				Log: log,
				KeeperRegistryInsufficientFundsUpkeepReport: *l,
			}

			results = append(results, r)
		}
	}

	return results, nil
}

// Fetches the checkBlockNumber for a particular transaction and an upkeep ID. Requires a RPC call to get txData
// so this function should not be used heavily
func (c *LogProvider) getCheckBlockNumberFromTxHash(txHash common.Hash, id ocr2keepers.UpkeepIdentifier) (bk ocr2keepers.BlockKey, e error) {
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
		// TODO: the log provider should be in the evm package for isolation
		res, ok := upkeep.(EVMAutomationUpkeepResult20)
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

type performed struct {
	logpoller.Log
	registry.KeeperRegistryUpkeepPerformed
}

type reorged struct {
	logpoller.Log
	registry.KeeperRegistryReorgedUpkeepReport
}

type staleUpkeep struct {
	logpoller.Log
	registry.KeeperRegistryStaleUpkeepReport
}

type insufficientFunds struct {
	logpoller.Log
	registry.KeeperRegistryInsufficientFundsUpkeepReport
}
