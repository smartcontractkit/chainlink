package llo_test

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	chainselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
	"github.com/smartcontractkit/chainlink-common/pkg/utils"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/channel_config_store"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/llo"
)

func Test_ChannelDefinitionCache_Integration(t *testing.T) {
	lggr, observedLogs := logger.TestLoggerObserved(t, zapcore.InfoLevel)
	db := pgtest.NewSqlxDB(t)
	ctx := testutils.Context(t)
	orm := llo.NewORM(db, testutils.SimulatedChainID)

	steve := testutils.MustNewSimTransactor(t) // config contract deployer and owner
	genesisData := core.GenesisAlloc{steve.From: {Balance: assets.Ether(1000).ToInt()}}
	backend := cltest.NewSimulatedBackend(t, genesisData, uint32(ethconfig.Defaults.Miner.GasCeil))
	backend.Commit() // ensure starting block number at least 1

	ethClient := client.NewSimulatedBackendClient(t, backend, testutils.SimulatedChainID)

	configStoreAddress, _, configStoreContract, err := channel_config_store.DeployChannelConfigStore(steve, backend)
	require.NoError(t, err)

	channel1 := rand.Uint32()
	channel2 := rand.Uint32()
	channel3 := rand.Uint32()

	chainSelector, err := chainselectors.SelectorFromChainId(testutils.SimulatedChainID.Uint64())
	require.NoError(t, err)

	streamIDs := []uint32{1, 2, 3}
	channel1Def := channel_config_store.IChannelConfigStoreChannelDefinition{
		ReportFormat:  uint32(llotypes.ReportFormatSolana),
		ChainSelector: chainSelector,
		StreamIDs:     streamIDs,
	}
	channel2Def := channel_config_store.IChannelConfigStoreChannelDefinition{
		ReportFormat:  uint32(llotypes.ReportFormatEVM),
		ChainSelector: chainSelector,
		StreamIDs:     streamIDs,
	}
	channel3Def := channel_config_store.IChannelConfigStoreChannelDefinition{
		ReportFormat:  uint32(llotypes.ReportFormatEVM),
		ChainSelector: chainSelector,
		StreamIDs:     append(streamIDs, 4),
	}

	require.NoError(t, utils.JustError(configStoreContract.AddChannel(steve, channel1, channel1Def)))
	require.NoError(t, utils.JustError(configStoreContract.AddChannel(steve, channel2, channel2Def)))

	h := backend.Commit()
	channel2Block, err := backend.BlockByHash(ctx, h)
	require.NoError(t, err)

	t.Run("with zero fromblock", func(t *testing.T) {
		lpOpts := logpoller.Opts{
			PollPeriod:               100 * time.Millisecond,
			FinalityDepth:            1,
			BackfillBatchSize:        3,
			RpcBatchSize:             2,
			KeepFinalizedBlocksDepth: 1000,
		}
		lp := logpoller.NewLogPoller(
			logpoller.NewORM(testutils.SimulatedChainID, db, lggr), ethClient, lggr, lpOpts)
		servicetest.Run(t, lp)
		cdc := llo.NewChannelDefinitionCache(lggr, orm, lp, configStoreAddress, 0)

		servicetest.Run(t, cdc)

		testutils.WaitForLogMessage(t, observedLogs, "Updated channel definitions")

		dfns := cdc.Definitions()

		require.Len(t, dfns, 2)
		require.Contains(t, dfns, channel1)
		require.Contains(t, dfns, channel2)
		assert.Equal(t, llotypes.ChannelDefinition{
			ReportFormat:  llotypes.ReportFormatSolana,
			ChainSelector: chainSelector,
			StreamIDs:     []uint32{1, 2, 3},
		}, dfns[channel1])
		assert.Equal(t, llotypes.ChannelDefinition{
			ReportFormat:  llotypes.ReportFormatEVM,
			ChainSelector: chainSelector,
			StreamIDs:     []uint32{1, 2, 3},
		}, dfns[channel2])

		// remove solana
		require.NoError(t, utils.JustError(configStoreContract.RemoveChannel(steve, channel1)))
		backend.Commit()
		testutils.WaitForLogMessageCount(t, observedLogs, "Updated channel definitions", 2)
		dfns = cdc.Definitions()

		require.Len(t, dfns, 1)
		assert.NotContains(t, dfns, channel1)
		require.Contains(t, dfns, channel2)

		assert.Equal(t, llotypes.ChannelDefinition{
			ReportFormat:  llotypes.ReportFormatEVM,
			ChainSelector: chainSelector,
			StreamIDs:     []uint32{1, 2, 3},
		}, dfns[channel2])

		// add channel3 with additional stream
		require.NoError(t, utils.JustError(configStoreContract.AddChannel(steve, channel3, channel3Def)))
		backend.Commit()
		testutils.WaitForLogMessageCount(t, observedLogs, "Updated channel definitions", 3)
		dfns = cdc.Definitions()

		require.Len(t, dfns, 2)
		require.Contains(t, dfns, channel2)
		require.Contains(t, dfns, channel3)

		assert.Equal(t, llotypes.ChannelDefinition{
			ReportFormat:  llotypes.ReportFormatEVM,
			ChainSelector: chainSelector,
			StreamIDs:     []uint32{1, 2, 3},
		}, dfns[channel2])
		assert.Equal(t, llotypes.ChannelDefinition{
			ReportFormat:  llotypes.ReportFormatEVM,
			ChainSelector: chainSelector,
			StreamIDs:     []uint32{1, 2, 3, 4},
		}, dfns[channel3])
	})

	t.Run("loads from ORM", func(t *testing.T) {
		// Override logpoller to always return no logs
		lpOpts := logpoller.Opts{
			PollPeriod:               100 * time.Millisecond,
			FinalityDepth:            1,
			BackfillBatchSize:        3,
			RpcBatchSize:             2,
			KeepFinalizedBlocksDepth: 1000,
		}
		lp := &mockLogPoller{
			LogPoller: logpoller.NewLogPoller(logpoller.NewORM(testutils.SimulatedChainID, db, lggr), ethClient, lggr, lpOpts),
			LatestBlockFn: func(ctx context.Context) (int64, error) {
				return 0, nil
			},
			LogsWithSigsFn: func(ctx context.Context, start, end int64, eventSigs []common.Hash, address common.Address) ([]logpoller.Log, error) {
				return []logpoller.Log{}, nil
			},
		}
		cdc := llo.NewChannelDefinitionCache(lggr, orm, lp, configStoreAddress, 0)

		servicetest.Run(t, cdc)

		dfns := cdc.Definitions()

		require.Len(t, dfns, 2)
		require.Contains(t, dfns, channel2)
		require.Contains(t, dfns, channel3)

		assert.Equal(t, llotypes.ChannelDefinition{
			ReportFormat:  llotypes.ReportFormatEVM,
			ChainSelector: chainSelector,
			StreamIDs:     []uint32{1, 2, 3},
		}, dfns[channel2])
		assert.Equal(t, llotypes.ChannelDefinition{
			ReportFormat:  llotypes.ReportFormatEVM,
			ChainSelector: chainSelector,
			StreamIDs:     []uint32{1, 2, 3, 4},
		}, dfns[channel3])
	})

	// clear out DB for next test
	pgtest.MustExec(t, db, `DELETE FROM channel_definitions`)

	t.Run("with non-zero fromBlock", func(t *testing.T) {
		lpOpts := logpoller.Opts{
			PollPeriod:               100 * time.Millisecond,
			FinalityDepth:            1,
			BackfillBatchSize:        3,
			RpcBatchSize:             2,
			KeepFinalizedBlocksDepth: 1000,
		}
		lp := logpoller.NewLogPoller(logpoller.NewORM(testutils.SimulatedChainID, db, lggr), ethClient, lggr, lpOpts)
		servicetest.Run(t, lp)
		cdc := llo.NewChannelDefinitionCache(lggr, orm, lp, configStoreAddress, channel2Block.Number().Int64()+1)

		// should only detect events from AFTER channel 2 was added
		servicetest.Run(t, cdc)

		testutils.WaitForLogMessageCount(t, observedLogs, "Updated channel definitions", 4)

		dfns := cdc.Definitions()

		require.Len(t, dfns, 1)
		require.Contains(t, dfns, channel3)

		assert.Equal(t, llotypes.ChannelDefinition{
			ReportFormat:  llotypes.ReportFormatEVM,
			ChainSelector: chainSelector,
			StreamIDs:     []uint32{1, 2, 3, 4},
		}, dfns[channel3])
	})
}

type mockLogPoller struct {
	logpoller.LogPoller
	LatestBlockFn  func(ctx context.Context) (int64, error)
	LogsWithSigsFn func(ctx context.Context, start, end int64, eventSigs []common.Hash, address common.Address) ([]logpoller.Log, error)
}

func (p *mockLogPoller) LogsWithSigs(ctx context.Context, start, end int64, eventSigs []common.Hash, address common.Address) ([]logpoller.Log, error) {
	return p.LogsWithSigsFn(ctx, start, end, eventSigs, address)
}
func (p *mockLogPoller) LatestBlock(ctx context.Context) (logpoller.LogPollerBlock, error) {
	block, err := p.LatestBlockFn(ctx)
	return logpoller.LogPollerBlock{BlockNumber: block}, err
}
