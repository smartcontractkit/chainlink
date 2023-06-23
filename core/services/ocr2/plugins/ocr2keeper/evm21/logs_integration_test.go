package evm_test

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/sqlx"

	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg"

	"github.com/smartcontractkit/chainlink/v2/core/assets"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	iregistry21 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_keeper_registry_master_wrapper_2_1"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_log_automation"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/log_upkeep_counter_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	kevm21 "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21"
)

func TestIntegration_LogEventProvider(t *testing.T) {
	ctx, cancel := context.WithCancel(testutils.Context(t))
	defer cancel()

	backend, stopMining, accounts := setupBackend(t)
	defer stopMining()
	carrol := accounts[2]

	db := setupDB(t)
	defer db.Close()

	opts := &kevm21.LogEventProviderOptions{
		FetchInterval:   time.Second / 2,
		FetchPartitions: 2,
	}
	logProvider, lp, ethClient := setupLogProvider(t, db, backend, opts)

	n := 10

	ids, addrs, contracts := deployUpkeepCounter(t, n, backend, carrol, logProvider)
	lp.PollAndSaveLogs(ctx, int64(n))

	go func() {
		if err := logProvider.Start(ctx); err != nil {
			t.Logf("error starting log provider: %s", err)
			t.Fail()
		}
	}()
	defer logProvider.Close()

	logsRounds := 10
	pollerTimeout := time.Second * 5

	poll := pollFn(ctx, t, lp, ethClient)

	triggerEvents(ctx, t, backend, carrol, logsRounds, poll, contracts...)

	poll(backend.Commit())
	// let it time to poll
	<-time.After(pollerTimeout)

	logs, _ := logProvider.GetLogs()
	require.NoError(t, logProvider.Close())

	require.GreaterOrEqual(t, len(logs), n, "failed to get all logs")
	t.Run("Restart", func(t *testing.T) {
		// assuming that our service was closed and restarted,
		// we should be able to backfill old logs and fetch new ones
		require.NoError(t, logProvider.Close())

		poll(backend.Commit())

		go func() {
			if err := logProvider.Start(ctx); err != nil {
				t.Logf("error starting log provider: %s", err)
				t.Fail()
			}
		}()
		defer logProvider.Close()

		for i, addr := range addrs {
			id := ids[i]
			require.NoError(t, logProvider.RegisterFilter(id, newPlainLogTriggerConfig(addr)))
		}
		logsAfterRestart, _ := logProvider.GetLogs()
		require.GreaterOrEqual(t, len(logsAfterRestart), 0,
			"logs should have been marked visited")

		triggerEvents(ctx, t, backend, carrol, logsRounds, poll, contracts...)
		// let it time to poll
		poll(backend.Commit())

		<-time.After(pollerTimeout)

		logsAfterRestart, _ = logProvider.GetLogs()
		require.NoError(t, logProvider.Close())
		require.GreaterOrEqual(t, len(logsAfterRestart), n,
			"failed to get logs after restart")
	})
}

func TestIntegration_LogEventProvider_RateLimit(t *testing.T) {
	ctx, cancel := context.WithCancel(testutils.Context(t))
	defer cancel()

	backend, stopMining, accounts := setupBackend(t)
	defer stopMining()
	carrol := accounts[2]

	db := setupDB(t)
	defer db.Close()

	opts := &kevm21.LogEventProviderOptions{
		FetchInterval:   time.Second / 2,
		FetchPartitions: 2,
	}
	logProvider, lp, ethClient := setupLogProvider(t, db, backend, opts)

	n := 10

	ids, _, contracts := deployUpkeepCounter(t, n, backend, carrol, logProvider)
	lp.PollAndSaveLogs(ctx, int64(n))
	poll := pollFn(ctx, t, lp, ethClient)

	var wg sync.WaitGroup

	for i := 0; i < 3; i++ {
		triggerEvents(ctx, t, backend, carrol, 25, poll, contracts...)
		poll(backend.Commit())
		fetchWorkers := 20
		limitErrs := int32(0)
		for i := 0; i < fetchWorkers; i++ {
			idsCp := make([]*big.Int, len(ids))
			copy(idsCp, ids)
			wg.Add(1)
			go func(i int, ids []*big.Int) {
				defer wg.Done()
				err := logProvider.FetchLogs(ctx, true, ids...)
				if err != nil {
					require.True(t, strings.Contains(err.Error(), "block limit exceeded"))
					atomic.AddInt32(&limitErrs, 1)
				}
			}(i, idsCp)
		}
		_ = backend.Commit()
		poll(backend.Commit())

		wg.Wait()
		require.GreaterOrEqual(t, 1, 1)
	}

	logs, err := logProvider.GetLogs()
	require.NoError(t, err)
	require.NoError(t, logProvider.Close())

	require.GreaterOrEqual(t, len(logs), 1, "failed to fetch logs")
	require.LessOrEqual(t, len(logs), n, "failed to apply limits")
}

func TestIntegration_LogEventProvider_Backfill(t *testing.T) {
	ctx, cancel := context.WithCancel(testutils.Context(t))
	defer cancel()

	backend, stopMining, accounts := setupBackend(t)
	defer stopMining()
	carrol := accounts[2]

	db := setupDB(t)
	defer db.Close()

	logProvider, lp, ethClient := setupLogProvider(t, db, backend, &kevm21.LogEventProviderOptions{
		FetchInterval:     time.Second / 2,
		LogBlocksLookback: 512,
		FetchPartitions:   2,
	})

	n := 20
	pollerTimeout := time.Second * 2

	_, _, contracts := deployUpkeepCounter(t, n, backend, carrol, logProvider)
	lp.PollAndSaveLogs(ctx, int64(n))

	logsRounds := 50
	poll := pollFn(ctx, t, lp, ethClient)

	go func() {
		if err := logProvider.Start(ctx); err != nil {
			t.Logf("error starting log provider: %s", err)
			t.Fail()
		}
	}()
	defer logProvider.Close()

	triggerEvents(ctx, t, backend, carrol, logsRounds, poll, contracts...)

	_ = backend.Commit()
	poll(backend.Commit())

	<-time.After(pollerTimeout * 2) // let the log poller work

	logs, err := logProvider.GetLogs()
	require.NoError(t, err)
	require.NoError(t, logProvider.Close())

	require.GreaterOrEqual(t, len(logs), n, "failed to backfill logs")
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
	t *testing.T,
	n int,
	backend *backends.SimulatedBackend,
	account *bind.TransactOpts,
	logProvider kevm21.LogEventProvider,
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
		id := big.NewInt(0).SetBytes(upkeepID)
		ids = append(ids, id)
		err = logProvider.RegisterFilter(id, newPlainLogTriggerConfig(upkeepAddr))
		require.NoError(t, err)
	}
	return ids, contractsAddrs, contracts
}

func newPlainLogTriggerConfig(upkeepAddr common.Address) kevm21.LogTriggerConfig {
	return kevm21.LogTriggerConfig{
		ContractAddress: upkeepAddr,
		FilterSelector:  0,
		Topic0:          common.HexToHash("0x3d53a39550e04688065827f3bb86584cb007ab9ebca7ebd528e7301c9c31eb5d"),
	}
}

func setupLogProvider(t *testing.T, db *sqlx.DB, backend *backends.SimulatedBackend, opts *kevm21.LogEventProviderOptions) (kevm21.LogEventProviderTest, logpoller.LogPollerTest, *evmclient.SimulatedBackendClient) {
	ethClient := evmclient.NewSimulatedBackendClient(t, backend, big.NewInt(1337))
	pollerLggr := logger.TestLogger(t)
	pollerLggr.SetLogLevel(zapcore.WarnLevel)
	lorm := logpoller.NewORM(big.NewInt(1337), db, pollerLggr, pgtest.NewQConfig(false))
	lp := logpoller.NewLogPoller(lorm, ethClient, pollerLggr, 100*time.Millisecond, 1, 2, 2, 1000)

	lggr := logger.TestLogger(t)
	keeperRegistryABI, err := abi.JSON(strings.NewReader(iregistry21.IKeeperRegistryMasterABI))
	require.NoError(t, err)
	logDataABI, err := abi.JSON(strings.NewReader(i_log_automation.ILogAutomationABI))
	require.NoError(t, err)
	logProvider := kevm21.NewLogEventProvider(lggr, lp, kevm21.NewEvmRegistryPackerV2_1(keeperRegistryABI, logDataABI), opts)

	return logProvider, lp, ethClient
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
