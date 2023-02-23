package ocr2keeper

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	pluginutils "github.com/smartcontractkit/ocr2keepers/pkg/chain"
	plugintypes "github.com/smartcontractkit/ocr2keepers/pkg/types"

	evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/core/chains/evm/logpoller"
	registry "github.com/smartcontractkit/chainlink/core/gethwrappers/generated/keeper_registry_wrapper2_0"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pg"
)

type LogProvider struct {
	logger          logger.Logger
	logPoller       logpoller.LogPoller
	FilterName      string
	registryAddress common.Address
	lookbackBlocks  int64
	registry        *registry.KeeperRegistry
}

var _ plugintypes.PerformLogProvider = (*LogProvider)(nil)

func NewLogProvider(
	logger logger.Logger,
	logPoller logpoller.LogPoller,
	registryAddress common.Address,
	client evmclient.Client,
	lookbackBlocks int64,
) (*LogProvider, error) {
	var err error
	filterName := logpoller.FilterName("OCR2KeeperRegistry", registryAddress)

	contract, err := registry.NewKeeperRegistry(common.HexToAddress("0x"), client)
	if err != nil {
		return nil, err
	}

	// Add log filters for the log poller so that it can poll and find the logs that
	// we need.
	err = logPoller.RegisterFilter(logpoller.Filter{
		Name: filterName,
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
		logger:          logger,
		logPoller:       logPoller,
		FilterName:      filterName,
		registryAddress: registryAddress,
		lookbackBlocks:  lookbackBlocks,
		registry:        contract,
	}, nil
}

func (c *LogProvider) PerformLogs(ctx context.Context) ([]plugintypes.PerformLog, error) {
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
			registry.KeeperRegistryUpkeepPerformed{}.Topic(),
		},
		c.registryAddress,
		pg.WithParentCtx(ctx),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to collect logs from log poller", err)
	}

	performed, err := c.unmarshalPerformLogs(logs)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to unmarshal logs", err)
	}

	vals := []plugintypes.PerformLog{}
	for _, p := range performed {
		// broadcast log to subscribers
		l := plugintypes.PerformLog{
			Key:             pluginutils.NewUpkeepKey(big.NewInt(int64(p.CheckBlockNumber)), p.Id),
			TransmitBlock:   pluginutils.BlockKey([]byte(fmt.Sprintf("%d", p.BlockNumber))),
			TransactionHash: p.TxHash.Hex(),
			Confirmations:   end - p.BlockNumber,
		}
		vals = append(vals, l)
	}

	return vals, nil
}

func (c *LogProvider) StaleReportLogs(ctx context.Context) ([]plugintypes.StaleReportLog, error) {
	end, err := c.logPoller.LatestBlock(pg.WithParentCtx(ctx))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to get latest block from log poller", err)
	}

	// always check the last lookback number of blocks and rebroadcast
	// this allows the plugin to make decisions based on event confirmations

	// ReorgedUpkeepReportLogs
	logs, err := c.logPoller.LogsWithSigs(
		end-c.lookbackBlocks,
		end,
		[]common.Hash{
			registry.KeeperRegistryReorgedUpkeepReport{}.Topic(),
		},
		c.registryAddress,
		pg.WithParentCtx(ctx),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to collect logs from log poller", err)
	}

	reorged, err := c.unmarshalReorgLogs(logs)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to unmarshal reorg logs", err)
	}

	// StaleUpkeepReportLogs
	logs, err = c.logPoller.LogsWithSigs(
		end-c.lookbackBlocks,
		end,
		[]common.Hash{
			registry.KeeperRegistryStaleUpkeepReport{}.Topic(),
		},
		c.registryAddress,
		pg.WithParentCtx(ctx),
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
		end-c.lookbackBlocks,
		end,
		[]common.Hash{
			registry.KeeperRegistryInsufficientFundsUpkeepReport{}.Topic(),
		},
		c.registryAddress,
		pg.WithParentCtx(ctx),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to collect logs from log poller", err)
	}

	insufficientFunds, err := c.unmarshalInsufficientFundsUpkeepLogs(logs)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to unmarshal insufficient fund upkeep logs", err)
	}

	vals := []plugintypes.StaleReportLog{}
	for _, r := range reorged {
		l := plugintypes.StaleReportLog{
			UpkeepId:        plugintypes.UpkeepIdentifier(r.Id.String()),
			TransmitBlock:   pluginutils.BlockKey([]byte(fmt.Sprintf("%d", r.BlockNumber))),
			TransactionHash: r.TxHash.Hex(),
			Confirmations:   end - r.BlockNumber,
		}
		vals = append(vals, l)
	}
	for _, r := range staleUpkeep {
		l := plugintypes.StaleReportLog{
			UpkeepId:        plugintypes.UpkeepIdentifier(r.Id.String()),
			TransmitBlock:   pluginutils.BlockKey([]byte(fmt.Sprintf("%d", r.BlockNumber))),
			TransactionHash: r.TxHash.Hex(),
			Confirmations:   end - r.BlockNumber,
		}
		vals = append(vals, l)
	}
	for _, r := range insufficientFunds {
		l := plugintypes.StaleReportLog{
			UpkeepId:        plugintypes.UpkeepIdentifier(r.Id.String()),
			TransmitBlock:   pluginutils.BlockKey([]byte(fmt.Sprintf("%d", r.BlockNumber))),
			TransactionHash: r.TxHash.Hex(),
			Confirmations:   end - r.BlockNumber,
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

func (c *LogProvider) unmarshalReorgLogs(logs []logpoller.Log) ([]reorged, error) {
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
