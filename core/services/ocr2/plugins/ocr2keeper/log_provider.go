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

	contract, err := registry.NewKeeperRegistry(common.HexToAddress("0x"), client)
	if err != nil {
		return nil, err
	}

	// Add log filters for the log poller so that it can poll and find the logs that
	// we need. Not unregistering the filter later so we ignore the id
	_, err = logPoller.RegisterFilter(logpoller.Filter{
		EventSigs: []common.Hash{
			registry.KeeperRegistryUpkeepPerformed{}.Topic(),
		},
		Addresses: []common.Address{registryAddress},
	})
	if err != nil {
		return nil, err
	}

	return &LogProvider{
		logger:          logger,
		logPoller:       logPoller,
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

	performed, err := c.unmarshalLogs(logs)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to unmarshal logs", err)
	}

	vals := []plugintypes.PerformLog{}
	for _, p := range performed {
		// broadcast log to subscribers
		l := plugintypes.PerformLog{
			Key:             pluginutils.BlockAndIdToKey(big.NewInt(int64(p.CheckBlockNumber)), p.Id),
			TransmitBlock:   plugintypes.BlockKey([]byte(fmt.Sprintf("%d", p.BlockNumber))),
			TransactionHash: p.TxHash.Hex(),
			Confirmations:   end - p.BlockNumber,
		}
		vals = append(vals, l)
	}

	return vals, nil
}

func (c *LogProvider) unmarshalLogs(logs []logpoller.Log) ([]performed, error) {
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

type performed struct {
	logpoller.Log
	registry.KeeperRegistryUpkeepPerformed
}
