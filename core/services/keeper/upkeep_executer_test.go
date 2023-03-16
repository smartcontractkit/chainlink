package keeper_test

import (
	"math/big"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/onsi/gomega"
	"github.com/smartcontractkit/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/chains/evm"
	"github.com/smartcontractkit/chainlink/core/chains/evm/gas"
	gasmocks "github.com/smartcontractkit/chainlink/core/chains/evm/gas/mocks"
	evmmocks "github.com/smartcontractkit/chainlink/core/chains/evm/mocks"
	"github.com/smartcontractkit/chainlink/core/chains/evm/txmgr"
	txmmocks "github.com/smartcontractkit/chainlink/core/chains/evm/txmgr/mocks"
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/config"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	configtest "github.com/smartcontractkit/chainlink/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keeper"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func newHead() evmtypes.Head {
	return evmtypes.NewHead(big.NewInt(20), utils.NewHash(), utils.NewHash(), 1000, utils.NewBigI(0))
}

func mockEstimator(t *testing.T) (estimator *gasmocks.Estimator) {
	estimator = gasmocks.NewEstimator(t)
	estimator.On("GetLegacyGas", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe().Return(assets.GWei(60), uint32(0), nil)
	estimator.On("GetDynamicFee", mock.Anything, mock.Anything, mock.Anything).Maybe().Return(gas.DynamicFee{
		FeeCap: assets.GWei(60),
		TipCap: assets.GWei(60),
	}, uint32(60), nil)
	return
}

func setup(t *testing.T, estimator *gasmocks.Estimator, overrideFn func(c *chainlink.Config, s *chainlink.Secrets)) (
	*sqlx.DB,
	config.GeneralConfig,
	*evmmocks.Client,
	*keeper.UpkeepExecuter,
	keeper.Registry,
	keeper.UpkeepRegistration,
	job.Job,
	cltest.JobPipelineV2TestHelper,
	*txmmocks.TxManager,
	keystore.Master,
	evm.Chain,
	keeper.ORM,
) {
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.Keeper.TurnLookBack = ptr[int64](0)
		if fn := overrideFn; fn != nil {
			fn(c, s)
		}
	})
	db := pgtest.NewSqlxDB(t)
	keyStore := cltest.NewKeyStore(t, db, cfg)
	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	ethClient.On("HeadByNumber", mock.Anything, mock.Anything).Maybe().Return(&evmtypes.Head{Number: 1, Hash: utils.NewHash()}, nil)
	txm := txmmocks.NewTxManager(t)
	txm.On("GetGasEstimator").Return(estimator)
	cc := evmtest.NewChainSet(t, evmtest.TestChainOpts{TxManager: txm, DB: db, Client: ethClient, KeyStore: keyStore.Eth(), GeneralConfig: cfg})
	jpv2 := cltest.NewJobPipelineV2(t, cfg, cc, db, keyStore, nil, nil)
	ch := evmtest.MustGetDefaultChain(t, cc)
	orm := keeper.NewORM(db, logger.TestLogger(t), ch.Config(), txmgr.SendEveryStrategy{})
	registry, job := cltest.MustInsertKeeperRegistry(t, db, orm, keyStore.Eth(), 0, 1, 20)
	lggr := logger.TestLogger(t)
	executer := keeper.NewUpkeepExecuter(job, orm, jpv2.Pr, ethClient, ch.HeadBroadcaster(), ch.TxManager().GetGasEstimator(), lggr, ch.Config(), job.KeeperSpec.FromAddress.Address())
	upkeep := cltest.MustInsertUpkeepForRegistry(t, db, ch.Config(), registry)
	err := executer.Start(testutils.Context(t))
	t.Cleanup(func() { executer.Close() })
	require.NoError(t, err)
	return db, cfg, ethClient, executer, registry, upkeep, job, jpv2, txm, keyStore, ch, orm
}

var checkUpkeepResponse = struct {
	PerformData    []byte
	MaxLinkPayment *big.Int
	GasLimit       *big.Int
	GasWei         *big.Int
	LinkEth        *big.Int
}{
	PerformData:    common.Hex2Bytes("1234"),
	MaxLinkPayment: big.NewInt(0), // doesn't matter
	GasLimit:       big.NewInt(2_000_000),
	GasWei:         big.NewInt(0), // doesn't matter
	LinkEth:        big.NewInt(0), // doesn't matter
}

var checkPerformResponse = struct {
	Success bool
}{
	Success: true,
}

func Test_UpkeepExecuter_ErrorsIfStartedTwice(t *testing.T) {
	t.Parallel()
	_, _, _, executer, _, _, _, _, _, _, _, _ := setup(t, mockEstimator(t), nil)
	err := executer.Start(testutils.Context(t)) // already started in setup()
	require.Error(t, err)
}

func Test_UpkeepExecuter_PerformsUpkeep_Happy(t *testing.T) {
	taskRuns := 11

	t.Parallel()

	t.Run("runs upkeep on triggering block number", func(t *testing.T) {
		db, config, ethMock, executer, registry, upkeep, job, jpv2, txm, _, _, _ := setup(t, mockEstimator(t), nil)

		gasLimit := 5_000_000 + config.KeeperRegistryPerformGasOverhead()

		ethTxCreated := cltest.NewAwaiter()
		txm.On("CreateEthTransaction",
			mock.MatchedBy(func(newTx txmgr.NewTx) bool { return newTx.GasLimit == gasLimit }),
		).
			Once().
			Return(txmgr.EthTx{
				ID: 1,
			}, nil).
			Run(func(mock.Arguments) { ethTxCreated.ItHappened() })

		registryMock := cltest.NewContractMockReceiver(t, ethMock, keeper.Registry1_1ABI, registry.ContractAddress.Address())
		registryMock.MockMatchedResponse(
			"checkUpkeep",
			func(callArgs ethereum.CallMsg) bool {
				return callArgs.GasPrice == nil &&
					callArgs.Gas == 0
			},
			checkUpkeepResponse,
		)
		registryMock.MockMatchedResponse(
			"performUpkeep",
			func(callArgs ethereum.CallMsg) bool { return true },
			checkPerformResponse,
		)

		head := newHead()
		executer.OnNewLongestChain(testutils.Context(t), &head)
		ethTxCreated.AwaitOrFail(t)
		runs := cltest.WaitForPipelineComplete(t, 0, job.ID, 1, taskRuns, jpv2.Jrm, time.Second, 100*time.Millisecond)
		require.Len(t, runs, 1)
		assert.False(t, runs[0].HasErrors())
		assert.False(t, runs[0].HasFatalErrors())
		waitLastRunHeight(t, db, upkeep, 20)
	})

	t.Run("runs upkeep on triggering block number on EIP1559 and non-EIP1559 chains", func(t *testing.T) {
		runTest := func(t *testing.T, eip1559 bool) {
			db, config, ethMock, executer, registry, upkeep, job, jpv2, txm, _, _, _ := setup(t, mockEstimator(t), func(c *chainlink.Config, s *chainlink.Secrets) {
				c.EVM[0].GasEstimator.EIP1559DynamicFees = &eip1559
			})

			gasLimit := 5_000_000 + config.KeeperRegistryPerformGasOverhead()

			ethTxCreated := cltest.NewAwaiter()
			txm.On("CreateEthTransaction",
				mock.MatchedBy(func(newTx txmgr.NewTx) bool { return newTx.GasLimit == gasLimit }),
			).
				Once().
				Return(txmgr.EthTx{
					ID: 1,
				}, nil).
				Run(func(mock.Arguments) { ethTxCreated.ItHappened() })

			registryMock := cltest.NewContractMockReceiver(t, ethMock, keeper.Registry1_1ABI, registry.ContractAddress.Address())
			registryMock.MockMatchedResponse(
				"checkUpkeep",
				func(callArgs ethereum.CallMsg) bool {
					return callArgs.GasPrice == nil &&
						callArgs.Gas == 0
				},
				checkUpkeepResponse,
			)
			registryMock.MockMatchedResponse(
				"performUpkeep",
				func(callArgs ethereum.CallMsg) bool { return true },
				checkPerformResponse,
			)

			head := newHead()

			executer.OnNewLongestChain(testutils.Context(t), &head)
			ethTxCreated.AwaitOrFail(t)
			runs := cltest.WaitForPipelineComplete(t, 0, job.ID, 1, taskRuns, jpv2.Jrm, time.Second, 100*time.Millisecond)
			require.Len(t, runs, 1)
			assert.False(t, runs[0].HasErrors())
			assert.False(t, runs[0].HasFatalErrors())
			waitLastRunHeight(t, db, upkeep, 20)
		}

		t.Run("EIP1559", func(t *testing.T) {
			runTest(t, true)
		})

		t.Run("non-EIP1559", func(t *testing.T) {
			runTest(t, false)
		})
	})

	t.Run("errors if submission key not found", func(t *testing.T) {
		_, _, ethMock, executer, registry, _, job, jpv2, _, keyStore, _, _ := setup(t, mockEstimator(t), nil)

		// replace expected key with random one
		_, err := keyStore.Eth().Create(&cltest.FixtureChainID)
		require.NoError(t, err)
		_, err = keyStore.Eth().Delete(job.KeeperSpec.FromAddress.Hex())
		require.NoError(t, err)

		registryMock := cltest.NewContractMockReceiver(t, ethMock, keeper.Registry1_1ABI, registry.ContractAddress.Address())
		registryMock.MockMatchedResponse(
			"checkUpkeep",
			func(callArgs ethereum.CallMsg) bool {
				return callArgs.GasPrice == nil &&
					callArgs.Gas == 0
			},
			checkUpkeepResponse,
		)
		registryMock.MockMatchedResponse(
			"performUpkeep",
			func(callArgs ethereum.CallMsg) bool { return true },
			checkPerformResponse,
		)

		head := newHead()
		executer.OnNewLongestChain(testutils.Context(t), &head)
		runs := cltest.WaitForPipelineError(t, 0, job.ID, 1, taskRuns, jpv2.Jrm, time.Second, 100*time.Millisecond)
		require.Len(t, runs, 1)
		assert.True(t, runs[0].HasErrors())
		assert.True(t, runs[0].HasFatalErrors())
	})

	t.Run("errors if submission chain not found", func(t *testing.T) {
		db, _, ethMock, _, _, _, _, jpv2, _, keyStore, ch, orm := setup(t, mockEstimator(t), nil)

		registry, jb := cltest.MustInsertKeeperRegistry(t, db, orm, keyStore.Eth(), 0, 1, 20)
		// change chain ID to non-configured chain
		jb.KeeperSpec.EVMChainID = (*utils.Big)(big.NewInt(999))
		cltest.MustInsertUpkeepForRegistry(t, db, ch.Config(), registry)
		lggr := logger.TestLogger(t)
		executer := keeper.NewUpkeepExecuter(jb, orm, jpv2.Pr, ethMock, ch.HeadBroadcaster(), ch.TxManager().GetGasEstimator(), lggr, ch.Config(), jb.KeeperSpec.FromAddress.Address())
		err := executer.Start(testutils.Context(t))
		require.NoError(t, err)
		head := newHead()
		executer.OnNewLongestChain(testutils.Context(t), &head)
		// TODO we want to see an errored run result once this is completed
		// https://app.shortcut.com/chainlinklabs/story/25397/remove-failearly-flag-from-eth-call-task
		cltest.AssertPipelineRunsStays(t, jb.PipelineSpecID, db, 0)
	})

	t.Run("triggers if heads are skipped but later heads arrive within range", func(t *testing.T) {
		db, config, ethMock, executer, registry, upkeep, job, jpv2, txm, _, _, _ := setup(t, mockEstimator(t), nil)

		etxs := []cltest.Awaiter{
			cltest.NewAwaiter(),
			cltest.NewAwaiter(),
		}
		gasLimit := 5_000_000 + config.KeeperRegistryPerformGasOverhead()
		txm.On("CreateEthTransaction",
			mock.MatchedBy(func(newTx txmgr.NewTx) bool { return newTx.GasLimit == gasLimit }),
		).
			Once().
			Return(txmgr.EthTx{}, nil).
			Run(func(mock.Arguments) { etxs[0].ItHappened() })

		registryMock := cltest.NewContractMockReceiver(t, ethMock, keeper.Registry1_1ABI, registry.ContractAddress.Address())
		registryMock.MockResponse("checkUpkeep", checkUpkeepResponse)
		registryMock.MockMatchedResponse(
			"performUpkeep",
			func(callArgs ethereum.CallMsg) bool { return true },
			checkPerformResponse,
		)
		// turn falls somewhere between 20-39 (blockCountPerTurn=20)
		// heads 20 thru 35 were skipped (e.g. due to node reboot)
		head := cltest.Head(36)

		executer.OnNewLongestChain(testutils.Context(t), head)
		runs := cltest.WaitForPipelineComplete(t, 0, job.ID, 1, taskRuns, jpv2.Jrm, time.Second, 100*time.Millisecond)
		require.Len(t, runs, 1)
		assert.False(t, runs[0].HasErrors())
		etxs[0].AwaitOrFail(t)
		waitLastRunHeight(t, db, upkeep, 36)
	})
}

func Test_UpkeepExecuter_PerformsUpkeep_Error(t *testing.T) {
	t.Parallel()

	g := gomega.NewWithT(t)

	db, _, ethMock, executer, registry, _, _, _, _, _, _, _ := setup(t, mockEstimator(t), nil)

	var wasCalled atomic.Bool
	registryMock := cltest.NewContractMockReceiver(t, ethMock, keeper.Registry1_1ABI, registry.ContractAddress.Address())
	registryMock.MockRevertResponse("checkUpkeep").Run(func(args mock.Arguments) {
		wasCalled.Store(true)
	})

	head := newHead()
	executer.OnNewLongestChain(testutils.Context(t), &head)

	g.Eventually(wasCalled.Load).Should(gomega.Equal(true))
	cltest.AssertCountStays(t, db, "eth_txes", 0)
}

func ptr[T any](t T) *T { return &t }
