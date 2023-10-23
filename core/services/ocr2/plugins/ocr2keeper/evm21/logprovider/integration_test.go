package logprovider_test

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"golang.org/x/time/rate"

	"github.com/smartcontractkit/sqlx"

	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg/v3/types"

	"github.com/smartcontractkit/chainlink/v2/core/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/log_upkeep_counter_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	kevmcore "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21/core"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21/logprovider"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

func TestIntegration_LogEventProvider(t *testing.T) {
	ctx, cancel := context.WithCancel(testutils.Context(t))
	defer cancel()

	backend, stopMining, accounts := setupBackend(t)
	defer stopMining()
	carrol := accounts[2]

	db := setupDB(t)
	defer db.Close()

	opts := logprovider.NewOptions(200)
	opts.ReadInterval = time.Second / 2
	lp, ethClient := setupDependencies(t, db, backend)
	filterStore := logprovider.NewUpkeepFilterStore()
	provider, _ := setup(logger.TestLogger(t), lp, nil, nil, filterStore, &opts)
	logProvider := provider.(logprovider.LogEventProviderTest)

	n := 10

	backend.Commit()
	lp.PollAndSaveLogs(ctx, 1) // Ensure log poller has a latest block

	ids, addrs, contracts := deployUpkeepCounter(ctx, t, n, ethClient, backend, carrol, logProvider)
	lp.PollAndSaveLogs(ctx, int64(n))

	go func() {
		if err := logProvider.Start(ctx); err != nil {
			t.Logf("error starting log provider: %s", err)
			t.Fail()
		}
	}()
	defer logProvider.Close()

	logsRounds := 10

	poll := pollFn(ctx, t, lp, ethClient)

	triggerEvents(ctx, t, backend, carrol, logsRounds, poll, contracts...)

	poll(backend.Commit())

	waitLogPoller(ctx, t, backend, lp, ethClient)

	waitLogProvider(ctx, t, logProvider, 3)

	allPayloads := collectPayloads(ctx, t, logProvider, n, 5)
	require.GreaterOrEqual(t, len(allPayloads), n,
		"failed to get logs after restart")

	t.Run("Restart", func(t *testing.T) {
		t.Log("restarting log provider")
		// assuming that our service was closed and restarted,
		// we should be able to backfill old logs and fetch new ones
		filterStore := logprovider.NewUpkeepFilterStore()
		logProvider2 := logprovider.NewLogProvider(logger.TestLogger(t), lp, logprovider.NewLogEventsPacker(), filterStore, opts)

		poll(backend.Commit())
		go func() {
			if err2 := logProvider2.Start(ctx); err2 != nil {
				t.Logf("error starting log provider: %s", err2)
				t.Fail()
			}
		}()
		defer logProvider2.Close()

		// re-register filters
		for i, id := range ids {
			err := logProvider2.RegisterFilter(ctx, logprovider.FilterOptions{
				UpkeepID:      id,
				TriggerConfig: newPlainLogTriggerConfig(addrs[i]),
				// using block number at which the upkeep was registered,
				// before we emitted any logs
				UpdateBlock: uint64(n),
			})
			require.NoError(t, err)
		}

		waitLogProvider(ctx, t, logProvider2, 2)

		t.Log("getting logs after restart")
		logsAfterRestart := collectPayloads(ctx, t, logProvider2, n, 5)
		require.GreaterOrEqual(t, len(logsAfterRestart), n,
			"failed to get logs after restart")
	})
}

func TestIntegration_LogEventProvider_UpdateConfig(t *testing.T) {
	ctx, cancel := context.WithCancel(testutils.Context(t))
	defer cancel()

	backend, stopMining, accounts := setupBackend(t)
	defer stopMining()
	carrol := accounts[2]

	db := setupDB(t)
	defer db.Close()

	opts := &logprovider.LogTriggersOptions{
		ReadInterval: time.Second / 2,
	}
	lp, ethClient := setupDependencies(t, db, backend)
	filterStore := logprovider.NewUpkeepFilterStore()
	provider, _ := setup(logger.TestLogger(t), lp, nil, nil, filterStore, opts)
	logProvider := provider.(logprovider.LogEventProviderTest)

	backend.Commit()
	lp.PollAndSaveLogs(ctx, 1) // Ensure log poller has a latest block
	_, addrs, contracts := deployUpkeepCounter(ctx, t, 1, ethClient, backend, carrol, logProvider)
	lp.PollAndSaveLogs(ctx, int64(5))
	require.Equal(t, 1, len(contracts))
	require.Equal(t, 1, len(addrs))

	t.Run("update filter config", func(t *testing.T) {
		upkeepID := kevmcore.GenUpkeepID(ocr2keepers.LogTrigger, "111")
		id := upkeepID.BigInt()
		cfg := newPlainLogTriggerConfig(addrs[0])
		b, err := ethClient.BlockByHash(ctx, backend.Commit())
		require.NoError(t, err)
		bn := b.Number()
		err = logProvider.RegisterFilter(ctx, logprovider.FilterOptions{
			UpkeepID:      id,
			TriggerConfig: cfg,
			UpdateBlock:   bn.Uint64(),
		})
		require.NoError(t, err)
		// old block
		err = logProvider.RegisterFilter(ctx, logprovider.FilterOptions{
			UpkeepID:      id,
			TriggerConfig: cfg,
			UpdateBlock:   bn.Uint64() - 1,
		})
		require.Error(t, err)
		// new block
		b, err = ethClient.BlockByHash(ctx, backend.Commit())
		require.NoError(t, err)
		bn = b.Number()
		err = logProvider.RegisterFilter(ctx, logprovider.FilterOptions{
			UpkeepID:      id,
			TriggerConfig: cfg,
			UpdateBlock:   bn.Uint64(),
		})
		require.NoError(t, err)
	})

	t.Run("register same log filter", func(t *testing.T) {
		upkeepID := kevmcore.GenUpkeepID(ocr2keepers.LogTrigger, "222")
		id := upkeepID.BigInt()
		cfg := newPlainLogTriggerConfig(addrs[0])
		b, err := ethClient.BlockByHash(ctx, backend.Commit())
		require.NoError(t, err)
		bn := b.Number()
		err = logProvider.RegisterFilter(ctx, logprovider.FilterOptions{
			UpkeepID:      id,
			TriggerConfig: cfg,
			UpdateBlock:   bn.Uint64(),
		})
		require.NoError(t, err)
	})
}

func TestIntegration_LogEventProvider_Backfill(t *testing.T) {
	ctx, cancel := context.WithTimeout(testutils.Context(t), time.Second*60)
	defer cancel()

	backend, stopMining, accounts := setupBackend(t)
	defer stopMining()
	carrol := accounts[2]

	db := setupDB(t)
	defer db.Close()

	opts := logprovider.NewOptions(200)
	opts.ReadInterval = time.Second / 4
	lp, ethClient := setupDependencies(t, db, backend)
	filterStore := logprovider.NewUpkeepFilterStore()
	provider, _ := setup(logger.TestLogger(t), lp, nil, nil, filterStore, &opts)
	logProvider := provider.(logprovider.LogEventProviderTest)

	n := 10

	backend.Commit()
	lp.PollAndSaveLogs(ctx, 1) // Ensure log poller has a latest block
	_, _, contracts := deployUpkeepCounter(ctx, t, n, ethClient, backend, carrol, logProvider)

	poll := pollFn(ctx, t, lp, ethClient)

	rounds := 8
	for i := 0; i < rounds; i++ {
		poll(backend.Commit())
		triggerEvents(ctx, t, backend, carrol, n, poll, contracts...)
		poll(backend.Commit())
	}

	waitLogPoller(ctx, t, backend, lp, ethClient)

	// starting the log provider should backfill logs
	go func() {
		if startErr := logProvider.Start(ctx); startErr != nil {
			t.Logf("error starting log provider: %s", startErr)
			t.Fail()
		}
	}()
	defer logProvider.Close()

	waitLogProvider(ctx, t, logProvider, 3)

	allPayloads := collectPayloads(ctx, t, logProvider, n, 5)
	require.GreaterOrEqual(t, len(allPayloads), len(contracts), "failed to backfill logs")
}

func TestIntegration_LogEventProvider_RateLimit(t *testing.T) {
	setupTest := func(
		t *testing.T,
		opts *logprovider.LogTriggersOptions,
	) (
		context.Context,
		*backends.SimulatedBackend,
		func(blockHash common.Hash),
		logprovider.LogEventProviderTest,
		[]*big.Int,
		func(),
	) {
		ctx, cancel := context.WithCancel(testutils.Context(t))
		backend, stopMining, accounts := setupBackend(t)
		userContractAccount := accounts[2]
		db := setupDB(t)

		deferFunc := func() {
			cancel()
			stopMining()
			_ = db.Close()
		}
		lp, ethClient := setupDependencies(t, db, backend)
		filterStore := logprovider.NewUpkeepFilterStore()
		provider, _ := setup(logger.TestLogger(t), lp, nil, nil, filterStore, opts)
		logProvider := provider.(logprovider.LogEventProviderTest)
		backend.Commit()
		lp.PollAndSaveLogs(ctx, 1) // Ensure log poller has a latest block

		rounds := 5
		numberOfUserContracts := 10
		poll := pollFn(ctx, t, lp, ethClient)

		// deployUpkeepCounter creates 'n' blocks and 'n' contracts
		ids, _, contracts := deployUpkeepCounter(
			ctx,
			t,
			numberOfUserContracts,
			ethClient,
			backend,
			userContractAccount,
			logProvider)

		// have log poller save logs for current blocks
		lp.PollAndSaveLogs(ctx, int64(numberOfUserContracts))

		for i := 0; i < rounds; i++ {
			triggerEvents(
				ctx,
				t,
				backend,
				userContractAccount,
				numberOfUserContracts,
				poll,
				contracts...)

			for dummyBlocks := 0; dummyBlocks < numberOfUserContracts; dummyBlocks++ {
				_ = backend.Commit()
			}

			poll(backend.Commit())
		}

		{
			// total block history at this point should be 566
			var minimumBlockCount int64 = 500
			latestBlock, _ := lp.LatestBlock()

			assert.GreaterOrEqual(t, latestBlock, minimumBlockCount, "to ensure the integrety of the test, the minimum block count before the test should be %d but got %d", minimumBlockCount, latestBlock)
		}

		require.NoError(t, logProvider.ReadLogs(ctx, ids...))

		return ctx, backend, poll, logProvider, ids, deferFunc
	}

	// polling for logs at approximately the same rate as a chain produces
	// blocks should not encounter rate limits
	t.Run("should allow constant polls within the rate and burst limit", func(t *testing.T) {
		ctx, backend, poll, logProvider, ids, deferFunc := setupTest(t, &logprovider.LogTriggersOptions{
			LookbackBlocks: 200,
			// BlockRateLimit is set low to ensure the test does not exceed the
			// rate limit
			BlockRateLimit: rate.Every(50 * time.Millisecond),
			// BlockLimitBurst is just set to a non-zero value
			BlockLimitBurst: 5,
		})

		defer deferFunc()

		// set the wait time between reads higher than the rate limit
		readWait := 50 * time.Millisecond
		timer := time.NewTimer(readWait)

		for i := 0; i < 4; i++ {
			<-timer.C

			// advance 1 block for every read
			poll(backend.Commit())

			err := logProvider.ReadLogs(ctx, ids...)
			if err != nil {
				assert.False(t, errors.Is(err, logprovider.ErrBlockLimitExceeded), "error should not contain block limit exceeded")
			}

			timer.Reset(readWait)
		}

		poll(backend.Commit())

		_, err := logProvider.GetLatestPayloads(ctx)

		require.NoError(t, err)
	})

	t.Run("should produce a rate limit error for over burst limit", func(t *testing.T) {
		ctx, backend, poll, logProvider, ids, deferFunc := setupTest(t, &logprovider.LogTriggersOptions{
			LookbackBlocks: 200,
			// BlockRateLimit is set low to ensure the test does not exceed the
			// rate limit
			BlockRateLimit: rate.Every(50 * time.Millisecond),
			// BlockLimitBurst is just set to a non-zero value
			BlockLimitBurst: 5,
		})

		defer deferFunc()

		// set the wait time between reads higher than the rate limit
		readWait := 50 * time.Millisecond
		timer := time.NewTimer(readWait)

		for i := 0; i < 4; i++ {
			<-timer.C

			// advance 4 blocks for every read
			for x := 0; x < 4; x++ {
				poll(backend.Commit())
			}

			err := logProvider.ReadLogs(ctx, ids...)
			if err != nil {
				assert.True(t, errors.Is(err, logprovider.ErrBlockLimitExceeded), "error should not contain block limit exceeded")
			}

			timer.Reset(readWait)
		}

		poll(backend.Commit())

		_, err := logProvider.GetLatestPayloads(ctx)

		require.NoError(t, err)
	})

	t.Run("should allow polling after lookback number of blocks have passed", func(t *testing.T) {
		ctx, backend, poll, logProvider, ids, deferFunc := setupTest(t, &logprovider.LogTriggersOptions{
			// BlockRateLimit is set low to ensure the test does not exceed the
			// rate limit
			BlockRateLimit: rate.Every(50 * time.Millisecond),
			// BlockLimitBurst is set low to ensure the test exceeds the burst limit
			BlockLimitBurst: 5,
			// LogBlocksLookback is set low to reduce the number of blocks required
			// to reset the block limiter to maxBurst
			LookbackBlocks: 50,
		})

		defer deferFunc()

		// simulate a burst in unpolled blocks
		for i := 0; i < 20; i++ {
			_ = backend.Commit()
		}

		poll(backend.Commit())

		// all entries should error at this point because there are too many
		// blocks to processes
		err := logProvider.ReadLogs(ctx, ids...)
		if err != nil {
			assert.True(t, errors.Is(err, logprovider.ErrBlockLimitExceeded), "error should not contain block limit exceeded")
		}

		// progress the chain by the same number of blocks as the lookback limit
		// to trigger the usage of maxBurst
		for i := 0; i < 50; i++ {
			_ = backend.Commit()
		}

		poll(backend.Commit())

		// all entries should reset to the maxBurst because they are beyond
		// the log lookback
		err = logProvider.ReadLogs(ctx, ids...)
		if err != nil {
			assert.True(t, errors.Is(err, logprovider.ErrBlockLimitExceeded), "error should not contain block limit exceeded")
		}

		poll(backend.Commit())

		_, err = logProvider.GetLatestPayloads(ctx)

		require.NoError(t, err)
	})
}

func TestIntegration_LogRecoverer_Backfill(t *testing.T) {
	t.Skip() // TODO: remove skip after removing constant timeouts
	ctx, cancel := context.WithTimeout(testutils.Context(t), time.Second*60)
	defer cancel()

	backend, stopMining, accounts := setupBackend(t)
	defer stopMining()
	carrol := accounts[2]

	db := setupDB(t)
	defer db.Close()

	lookbackBlocks := int64(200)
	opts := &logprovider.LogTriggersOptions{
		ReadInterval:   time.Second / 4,
		LookbackBlocks: lookbackBlocks,
	}
	lp, ethClient := setupDependencies(t, db, backend)
	filterStore := logprovider.NewUpkeepFilterStore()
	origDefaultRecoveryInterval := logprovider.RecoveryInterval
	logprovider.RecoveryInterval = time.Millisecond * 200
	defer func() {
		logprovider.RecoveryInterval = origDefaultRecoveryInterval
	}()
	provider, recoverer := setup(logger.TestLogger(t), lp, nil, &mockUpkeepStateStore{}, filterStore, opts)
	logProvider := provider.(logprovider.LogEventProviderTest)

	backend.Commit()
	lp.PollAndSaveLogs(ctx, 1) // Ensure log poller has a latest block

	n := 10
	_, _, contracts := deployUpkeepCounter(ctx, t, n, ethClient, backend, carrol, logProvider)

	poll := pollFn(ctx, t, lp, ethClient)

	rounds := 8
	for i := 0; i < rounds; i++ {
		triggerEvents(ctx, t, backend, carrol, n, poll, contracts...)
		poll(backend.Commit())
	}
	poll(backend.Commit())

	waitLogPoller(ctx, t, backend, lp, ethClient)

	// create dummy blocks
	var blockNumber int64
	for blockNumber < lookbackBlocks*4 {
		b, err := ethClient.BlockByHash(ctx, backend.Commit())
		require.NoError(t, err)
		bn := b.Number()
		blockNumber = bn.Int64()
	}
	// starting the log recoverer should backfill logs
	go func() {
		if startErr := recoverer.Start(ctx); startErr != nil {
			t.Logf("error starting log provider: %s", startErr)
			t.Fail()
		}
	}()
	defer recoverer.Close()

	lctx, lcancel := context.WithTimeout(ctx, time.Second*15)
	defer lcancel()
	var allProposals []ocr2keepers.UpkeepPayload
	for lctx.Err() == nil {
		poll(backend.Commit())
		proposals, err := recoverer.GetRecoveryProposals(ctx)
		require.NoError(t, err)
		allProposals = append(allProposals, proposals...)
		if len(allProposals) < n {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		break
	}
	require.NoError(t, lctx.Err(), "could not recover logs before timeout")
}

func collectPayloads(ctx context.Context, t *testing.T, logProvider logprovider.LogEventProvider, n, rounds int) []ocr2keepers.UpkeepPayload {
	allPayloads := make([]ocr2keepers.UpkeepPayload, 0)
	for ctx.Err() == nil && len(allPayloads) < n && rounds > 0 {
		logs, err := logProvider.GetLatestPayloads(ctx)
		require.NoError(t, err)
		require.LessOrEqual(t, len(logs), logprovider.AllowedLogsPerUpkeep, "failed to get all logs")
		allPayloads = append(allPayloads, logs...)
		rounds--
	}
	return allPayloads
}

// waitLogProvider waits until the provider reaches the given partition
func waitLogProvider(ctx context.Context, t *testing.T, logProvider logprovider.LogEventProviderTest, partition int) {
	t.Logf("waiting for log provider to reach partition %d", partition)
	for ctx.Err() == nil {
		currentPartition := logProvider.CurrentPartitionIdx()
		if currentPartition > uint64(partition) { // make sure we went over all items
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
}

// waitLogPoller waits until the log poller is familiar with the given block
func waitLogPoller(ctx context.Context, t *testing.T, backend *backends.SimulatedBackend, lp logpoller.LogPollerTest, ethClient *evmclient.SimulatedBackendClient) {
	t.Log("waiting for log poller to get updated")
	// let the log poller work
	b, err := ethClient.BlockByHash(ctx, backend.Commit())
	require.NoError(t, err)
	latestBlock := b.Number().Int64()
	for {
		latestPolled, lberr := lp.LatestBlock(pg.WithParentCtx(ctx))
		require.NoError(t, lberr)
		if latestPolled >= latestBlock {
			break
		}
		lp.PollAndSaveLogs(ctx, latestBlock)
	}
}

func pollFn(ctx context.Context, t *testing.T, lp logpoller.LogPollerTest, ethClient *evmclient.SimulatedBackendClient) func(blockHash common.Hash) {
	return func(blockHash common.Hash) {
		b, err := ethClient.BlockByHash(ctx, blockHash)
		require.NoError(t, err)
		bn := b.Number()
		lp.PollAndSaveLogs(ctx, bn.Int64())
	}
}

func triggerEvents(
	ctx context.Context,
	t *testing.T,
	backend *backends.SimulatedBackend,
	account *bind.TransactOpts,
	rounds int,
	poll func(blockHash common.Hash),
	contracts ...*log_upkeep_counter_wrapper.LogUpkeepCounter,
) {
	lctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var blockHash common.Hash
	for rounds > 0 && lctx.Err() == nil {
		rounds--
		for _, upkeepContract := range contracts {
			if lctx.Err() != nil {
				return
			}
			_, err := upkeepContract.Start(account)
			require.NoError(t, err)
			blockHash = backend.Commit()
		}
		poll(blockHash)
	}
}

func deployUpkeepCounter(
	ctx context.Context,
	t *testing.T,
	n int,
	ethClient *evmclient.SimulatedBackendClient,
	backend *backends.SimulatedBackend,
	account *bind.TransactOpts,
	logProvider logprovider.LogEventProvider,
) ([]*big.Int, []common.Address, []*log_upkeep_counter_wrapper.LogUpkeepCounter) {
	var ids []*big.Int
	var contracts []*log_upkeep_counter_wrapper.LogUpkeepCounter
	var contractsAddrs []common.Address
	for i := 0; i < n; i++ {
		upkeepAddr, _, upkeepContract, err := log_upkeep_counter_wrapper.DeployLogUpkeepCounter(
			account, backend,
			big.NewInt(100000),
		)
		require.NoError(t, err)
		backend.Commit()

		contracts = append(contracts, upkeepContract)
		contractsAddrs = append(contractsAddrs, upkeepAddr)

		// creating some dummy upkeepID to register filter
		upkeepID := ocr2keepers.UpkeepIdentifier(append(common.LeftPadBytes([]byte{1}, 16), upkeepAddr[:16]...))
		id := upkeepID.BigInt()
		ids = append(ids, id)
		b, err := ethClient.BlockByHash(context.Background(), backend.Commit())
		require.NoError(t, err)
		bn := b.Number()
		err = logProvider.RegisterFilter(ctx, logprovider.FilterOptions{
			UpkeepID:      id,
			TriggerConfig: newPlainLogTriggerConfig(upkeepAddr),
			UpdateBlock:   bn.Uint64(),
		})
		require.NoError(t, err)
	}
	return ids, contractsAddrs, contracts
}

func newPlainLogTriggerConfig(upkeepAddr common.Address) logprovider.LogTriggerConfig {
	return logprovider.LogTriggerConfig{
		ContractAddress: upkeepAddr,
		FilterSelector:  0,
		Topic0:          common.HexToHash("0x3d53a39550e04688065827f3bb86584cb007ab9ebca7ebd528e7301c9c31eb5d"),
	}
}

func setupDependencies(t *testing.T, db *sqlx.DB, backend *backends.SimulatedBackend) (logpoller.LogPollerTest, *evmclient.SimulatedBackendClient) {
	ethClient := evmclient.NewSimulatedBackendClient(t, backend, big.NewInt(1337))
	pollerLggr := logger.TestLogger(t)
	pollerLggr.SetLogLevel(zapcore.WarnLevel)
	lorm := logpoller.NewORM(big.NewInt(1337), db, pollerLggr, pgtest.NewQConfig(false))
	lp := logpoller.NewLogPoller(lorm, ethClient, pollerLggr, 100*time.Millisecond, false, 1, 2, 2, 1000)
	return lp, ethClient
}

func setup(lggr logger.Logger, poller logpoller.LogPoller, c client.Client, stateStore kevmcore.UpkeepStateReader, filterStore logprovider.UpkeepFilterStore, opts *logprovider.LogTriggersOptions) (logprovider.LogEventProvider, logprovider.LogRecoverer) {
	packer := logprovider.NewLogEventsPacker()
	if opts == nil {
		o := logprovider.NewOptions(200)
		opts = &o
	}
	provider := logprovider.NewLogProvider(lggr, poller, packer, filterStore, *opts)
	recoverer := logprovider.NewLogRecoverer(lggr, poller, c, stateStore, packer, filterStore, *opts)

	return provider, recoverer
}

func setupBackend(t *testing.T) (*backends.SimulatedBackend, func(), []*bind.TransactOpts) {
	sergey := testutils.MustNewSimTransactor(t) // owns all the link
	steve := testutils.MustNewSimTransactor(t)  // registry owner
	carrol := testutils.MustNewSimTransactor(t) // upkeep owner
	genesisData := core.GenesisAlloc{
		sergey.From: {Balance: assets.Ether(1000000000000000000).ToInt()},
		steve.From:  {Balance: assets.Ether(1000000000000000000).ToInt()},
		carrol.From: {Balance: assets.Ether(1000000000000000000).ToInt()},
	}
	backend := cltest.NewSimulatedBackend(t, genesisData, uint32(ethconfig.Defaults.Miner.GasCeil))
	stopMining := cltest.Mine(backend, 3*time.Second) // Should be greater than deltaRound since we cannot access old blocks on simulated blockchain
	return backend, stopMining, []*bind.TransactOpts{sergey, steve, carrol}
}

func ptr[T any](v T) *T { return &v }

func setupDB(t *testing.T) *sqlx.DB {
	_, db := heavyweight.FullTestDBV2(t, fmt.Sprintf("%s%d", "chainlink_test", 5432), func(c *chainlink.Config, s *chainlink.Secrets) {
		c.Feature.LogPoller = ptr(true)

		c.OCR.Enabled = ptr(false)
		c.OCR2.Enabled = ptr(true)

		c.EVM[0].Transactions.ForwardersEnabled = ptr(true)
		c.EVM[0].GasEstimator.Mode = ptr("FixedPrice")
	})
	return db
}

type mockUpkeepStateStore struct {
}

func (m *mockUpkeepStateStore) SelectByWorkIDs(ctx context.Context, workIDs ...string) ([]ocr2keepers.UpkeepState, error) {
	states := make([]ocr2keepers.UpkeepState, len(workIDs))
	for i := range workIDs {
		states[i] = ocr2keepers.UnknownState
	}
	return states, nil
}
