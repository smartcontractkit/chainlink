package logprovider_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-automation/pkg/v3/types"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	ocr2keepers "github.com/smartcontractkit/chainlink-common/pkg/types/automation"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/log_upkeep_counter_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	evmregistry21 "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v21/core"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v21/logprovider"
)

func TestIntegration_LogEventProvider(t *testing.T) {
	tests := []struct {
		name          string
		bufferVersion logprovider.BufferVersion
		logLimit      uint32
	}{
		{
			name:          "default buffer",
			bufferVersion: logprovider.BufferVersionDefault,
			logLimit:      10,
		},
		{
			name:          "buffer v1",
			bufferVersion: logprovider.BufferVersionV1,
			logLimit:      10,
		},
	}

	for _, tc := range tests {
		bufferVersion, logLimit := tc.bufferVersion, tc.logLimit
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(testutils.Context(t))
			defer cancel()

			backend, stopMining, accounts := setupBackend(t)
			defer stopMining()
			carrol := accounts[2]

			db := setupDB(t)
			defer db.Close()

			opts := logprovider.NewOptions(200, big.NewInt(1))
			opts.ReadInterval = time.Second / 2
			opts.BufferVersion = bufferVersion
			opts.LogLimit = logLimit

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

			allPayloads := collectPayloads(ctx, t, logProvider, n, logsRounds/2)
			require.GreaterOrEqual(t, len(allPayloads), n,
				"failed to get logs after restart")

			t.Run("Restart", func(t *testing.T) {
				t.Log("restarting log provider")
				// assuming that our service was closed and restarted,
				// we should be able to backfill old logs and fetch new ones
				filterStore := logprovider.NewUpkeepFilterStore()
				logProvider2 := logprovider.NewLogProvider(logger.TestLogger(t), lp, big.NewInt(1), logprovider.NewLogEventsPacker(), filterStore, opts)

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
		})
	}
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
		upkeepID := evmregistry21.GenUpkeepID(types.LogTrigger, "111")
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
		upkeepID := evmregistry21.GenUpkeepID(types.LogTrigger, "222")
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
	tests := []struct {
		name          string
		bufferVersion logprovider.BufferVersion
		logLimit      uint32
	}{
		{
			name:          "default buffer",
			bufferVersion: logprovider.BufferVersionDefault,
			logLimit:      10,
		},
		{
			name:          "buffer v1",
			bufferVersion: logprovider.BufferVersionV1,
			logLimit:      10,
		},
	}

	for _, tc := range tests {
		bufferVersion, limitLow := tc.bufferVersion, tc.logLimit
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(testutils.Context(t), time.Second*60)
			defer cancel()

			backend, stopMining, accounts := setupBackend(t)
			defer stopMining()
			carrol := accounts[2]

			db := setupDB(t)
			defer db.Close()

			opts := logprovider.NewOptions(200, big.NewInt(1))
			opts.ReadInterval = time.Second / 4
			opts.BufferVersion = bufferVersion
			opts.LogLimit = limitLow

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

			allPayloads := collectPayloads(ctx, t, logProvider, n*rounds, 5)
			require.GreaterOrEqual(t, len(allPayloads), len(contracts), "failed to backfill logs")
		})
	}
}

func TestIntegration_LogRecoverer_Backfill(t *testing.T) {
	ctx := testutils.Context(t)

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

	var allProposals []ocr2keepers.UpkeepPayload
	for {
		poll(backend.Commit())
		proposals, err := recoverer.GetRecoveryProposals(ctx)
		require.NoError(t, err)
		allProposals = append(allProposals, proposals...)
		if len(allProposals) >= n {
			break // success
		}
		select {
		case <-ctx.Done():
			t.Fatalf("could not recover logs before timeout: %s", ctx.Err())
		case <-time.After(100 * time.Millisecond):
		}
	}
}

func collectPayloads(ctx context.Context, t *testing.T, logProvider logprovider.LogEventProvider, n, rounds int) []ocr2keepers.UpkeepPayload {
	allPayloads := make([]ocr2keepers.UpkeepPayload, 0)
	for ctx.Err() == nil && len(allPayloads) < n && rounds > 0 {
		logs, err := logProvider.GetLatestPayloads(ctx)
		require.NoError(t, err)
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
		latestPolled, lberr := lp.LatestBlock(ctx)
		require.NoError(t, lberr)
		if latestPolled.BlockNumber >= latestBlock {
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
) (
	ids []*big.Int,
	contractsAddrs []common.Address,
	contracts []*log_upkeep_counter_wrapper.LogUpkeepCounter,
) {
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
		b, err := ethClient.BlockByHash(ctx, backend.Commit())
		require.NoError(t, err)
		bn := b.Number()
		err = logProvider.RegisterFilter(ctx, logprovider.FilterOptions{
			UpkeepID:      id,
			TriggerConfig: newPlainLogTriggerConfig(upkeepAddr),
			UpdateBlock:   bn.Uint64(),
		})
		require.NoError(t, err)
	}
	return
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
	lorm := logpoller.NewORM(big.NewInt(1337), db, pollerLggr)
	lpOpts := logpoller.Opts{
		PollPeriod:               100 * time.Millisecond,
		FinalityDepth:            1,
		BackfillBatchSize:        2,
		RpcBatchSize:             2,
		KeepFinalizedBlocksDepth: 1000,
	}
	lp := logpoller.NewLogPoller(lorm, ethClient, pollerLggr, lpOpts)
	return lp, ethClient
}

func setup(lggr logger.Logger, poller logpoller.LogPoller, c evmclient.Client, stateStore evmregistry21.UpkeepStateReader, filterStore logprovider.UpkeepFilterStore, opts *logprovider.LogTriggersOptions) (logprovider.LogEventProvider, logprovider.LogRecoverer) {
	packer := logprovider.NewLogEventsPacker()
	if opts == nil {
		o := logprovider.NewOptions(200, big.NewInt(1))
		opts = &o
	}
	provider := logprovider.NewLogProvider(lggr, poller, big.NewInt(1), packer, filterStore, *opts)
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
	_, db := heavyweight.FullTestDBV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
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
