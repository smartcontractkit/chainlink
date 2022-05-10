package keeper_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/onsi/gomega"
	"github.com/smartcontractkit/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/atomic"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/chains/evm"
	"github.com/smartcontractkit/chainlink/core/chains/evm/gas"
	gasmocks "github.com/smartcontractkit/chainlink/core/chains/evm/gas/mocks"
	evmmocks "github.com/smartcontractkit/chainlink/core/chains/evm/mocks"
	"github.com/smartcontractkit/chainlink/core/chains/evm/txmgr"
	txmmocks "github.com/smartcontractkit/chainlink/core/chains/evm/txmgr/mocks"
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keeper"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/utils"
	bigmath "github.com/smartcontractkit/chainlink/core/utils/big_math"
)

func newHead() evmtypes.Head {
	return evmtypes.NewHead(big.NewInt(20), utils.NewHash(), utils.NewHash(), 1000, utils.NewBigI(0))
}

func setup(t *testing.T) (
	*sqlx.DB,
	*configtest.TestGeneralConfig,
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
	cfg := cltest.NewTestGeneralConfig(t)
	cfg.Overrides.KeeperMaximumGracePeriod = null.IntFrom(0)
	cfg.Overrides.KeeperTurnLookBack = null.IntFrom(0)
	cfg.Overrides.KeeperTurnFlagEnabled = null.BoolFrom(true)
	cfg.Overrides.KeeperCheckUpkeepGasPriceFeatureEnabled = null.BoolFrom(true)
	db := pgtest.NewSqlxDB(t)
	keyStore := cltest.NewKeyStore(t, db, cfg)
	ethClient := cltest.NewEthClientMockWithDefaultChain(t)
	block := types.NewBlockWithHeader(&types.Header{
		Number: big.NewInt(1),
	})
	ethClient.On("BlockByNumber", mock.Anything, mock.Anything).Maybe().Return(block, nil)
	txm := new(txmmocks.TxManager)
	txm.Test(t)
	estimator := new(gasmocks.Estimator)
	estimator.Test(t)
	txm.On("GetGasEstimator").Return(estimator)
	estimator.On("GetLegacyGas", mock.Anything, mock.Anything).Maybe().Return(assets.GWei(60), uint64(0), nil)
	estimator.On("GetDynamicFee", mock.Anything).Maybe().Return(gas.DynamicFee{
		FeeCap: assets.GWei(60),
		TipCap: assets.GWei(60),
	}, uint64(60), nil)
	cc := evmtest.NewChainSet(t, evmtest.TestChainOpts{TxManager: txm, DB: db, Client: ethClient, KeyStore: keyStore.Eth(), GeneralConfig: cfg})
	jpv2 := cltest.NewJobPipelineV2(t, cfg, cc, db, keyStore, nil, nil)
	ch := evmtest.MustGetDefaultChain(t, cc)
	orm := keeper.NewORM(db, logger.TestLogger(t), ch.Config(), txmgr.SendEveryStrategy{})
	registry, job := cltest.MustInsertKeeperRegistry(t, db, orm, keyStore.Eth(), 0, 1, 20)
	lggr := logger.TestLogger(t)
	executer := keeper.NewUpkeepExecuter(job, orm, jpv2.Pr, ethClient, ch.HeadBroadcaster(), ch.TxManager().GetGasEstimator(), lggr, ch.Config())
	upkeep := cltest.MustInsertUpkeepForRegistry(t, db, ch.Config(), registry)
	err := executer.Start(testutils.Context(t))
	t.Cleanup(func() { txm.AssertExpectations(t); estimator.AssertExpectations(t); executer.Close() })
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

func Test_UpkeepExecuter_ErrorsIfStartedTwice(t *testing.T) {
	t.Parallel()
	_, _, _, executer, _, _, _, _, _, _, _, _ := setup(t)
	err := executer.Start(testutils.Context(t)) // already started in setup()
	require.Error(t, err)
}

func Test_UpkeepExecuter_PerformsUpkeep_Happy(t *testing.T) {
	t.Parallel()

	t.Run("runs upkeep on triggering block number", func(t *testing.T) {
		db, config, ethMock, executer, registry, upkeep, job, jpv2, txm, _, _, _ := setup(t)

		gasLimit := upkeep.ExecuteGas + config.KeeperRegistryPerformGasOverhead()
		gasPrice := bigmath.Div(bigmath.Mul(assets.GWei(60), 100+config.KeeperGasPriceBufferPercent()), 100)

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
				return bigmath.Equal(callArgs.GasPrice, gasPrice) &&
					callArgs.Gas == 650_000
			},
			checkUpkeepResponse,
		)

		head := newHead()
		executer.OnNewLongestChain(context.Background(), &head)
		ethTxCreated.AwaitOrFail(t)
		runs := cltest.WaitForPipelineComplete(t, 0, job.ID, 1, 5, jpv2.Jrm, time.Second, 100*time.Millisecond)
		require.Len(t, runs, 1)
		assert.False(t, runs[0].HasErrors())
		assert.False(t, runs[0].HasFatalErrors())
		waitLastRunHeight(t, db, upkeep, 20)

		ethMock.AssertExpectations(t)
		txm.AssertExpectations(t)
	})

	t.Run("runs upkeep on triggering block number on EIP1559 and non-EIP1559 chains", func(t *testing.T) {
		runTest := func(t *testing.T, eip1559 bool) {
			db, config, ethMock, executer, registry, upkeep, job, jpv2, txm, _, _, _ := setup(t)

			config.Overrides.GlobalEvmEIP1559DynamicFees = null.BoolFrom(eip1559)

			gasLimit := upkeep.ExecuteGas + config.KeeperRegistryPerformGasOverhead()
			gasPrice := bigmath.Div(bigmath.Mul(assets.GWei(60), 100+config.KeeperGasPriceBufferPercent()), 100)
			baseFeePerGas := utils.NewBig(big.NewInt(0).Mul(gasPrice, big.NewInt(2)))

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
					expectedGasPrice := bigmath.Div(
						bigmath.Mul(baseFeePerGas.ToInt(), 100+config.KeeperBaseFeeBufferPercent()),
						100,
					)

					return bigmath.Equal(callArgs.GasPrice, expectedGasPrice) &&
						650_000 == callArgs.Gas
				},
				checkUpkeepResponse,
			)

			head := newHead()
			head.BaseFeePerGas = baseFeePerGas

			executer.OnNewLongestChain(context.Background(), &head)
			ethTxCreated.AwaitOrFail(t)
			runs := cltest.WaitForPipelineComplete(t, 0, job.ID, 1, 5, jpv2.Jrm, time.Second, 100*time.Millisecond)
			require.Len(t, runs, 1)
			assert.False(t, runs[0].HasErrors())
			assert.False(t, runs[0].HasFatalErrors())
			waitLastRunHeight(t, db, upkeep, 20)

			ethMock.AssertExpectations(t)
			txm.AssertExpectations(t)
		}

		t.Run("EIP1559", func(t *testing.T) {
			runTest(t, true)
		})

		t.Run("non-EIP1559", func(t *testing.T) {
			runTest(t, false)
		})
	})

	t.Run("errors if submission key not found", func(t *testing.T) {
		_, config, ethMock, executer, registry, _, job, jpv2, _, keyStore, _, _ := setup(t)

		// replace expected key with random one
		_, err := keyStore.Eth().Create(&cltest.FixtureChainID)
		require.NoError(t, err)
		_, err = keyStore.Eth().Delete(job.KeeperSpec.FromAddress.Hex())
		require.NoError(t, err)

		gasPrice := bigmath.Div(bigmath.Mul(assets.GWei(60), 100+config.KeeperGasPriceBufferPercent()), 100)

		registryMock := cltest.NewContractMockReceiver(t, ethMock, keeper.Registry1_1ABI, registry.ContractAddress.Address())
		registryMock.MockMatchedResponse(
			"checkUpkeep",
			func(callArgs ethereum.CallMsg) bool {
				return bigmath.Equal(callArgs.GasPrice, gasPrice) &&
					callArgs.Gas == 650_000
			},
			checkUpkeepResponse,
		)

		head := newHead()
		executer.OnNewLongestChain(context.Background(), &head)
		runs := cltest.WaitForPipelineError(t, 0, job.ID, 1, 5, jpv2.Jrm, time.Second, 100*time.Millisecond)
		require.Len(t, runs, 1)
		assert.True(t, runs[0].HasErrors())
		assert.True(t, runs[0].HasFatalErrors())

		ethMock.AssertExpectations(t)
	})

	t.Run("errors if submission chain not found", func(t *testing.T) {
		db, _, ethMock, _, _, _, job, jpv2, _, _, ch, orm := setup(t)

		// change chain ID to non-configured chain
		job.KeeperSpec.EVMChainID = (*utils.Big)(big.NewInt(999))
		lggr := logger.TestLogger(t)
		executer := keeper.NewUpkeepExecuter(job, orm, jpv2.Pr, ethMock, ch.HeadBroadcaster(), ch.TxManager().GetGasEstimator(), lggr, ch.Config())
		err := executer.Start(testutils.Context(t))
		require.NoError(t, err)
		head := newHead()
		executer.OnNewLongestChain(context.Background(), &head)
		// TODO we want to see an errored run result once this is completed
		// https://app.shortcut.com/chainlinklabs/story/25397/remove-failearly-flag-from-eth-call-task
		cltest.AssertPipelineRunsStays(t, job.PipelineSpecID, db, 0)
		ethMock.AssertExpectations(t)
	})

	t.Run("triggers exactly one upkeep if heads are skipped but later heads arrive within range", func(t *testing.T) {
		db, config, ethMock, executer, registry, upkeep, job, jpv2, txm, _, _, _ := setup(t)

		etxs := []cltest.Awaiter{
			cltest.NewAwaiter(),
			cltest.NewAwaiter(),
		}
		gasLimit := upkeep.ExecuteGas + config.KeeperRegistryPerformGasOverhead()
		txm.On("CreateEthTransaction",
			mock.MatchedBy(func(newTx txmgr.NewTx) bool { return newTx.GasLimit == gasLimit }),
		).
			Once().
			Return(txmgr.EthTx{}, nil).
			Run(func(mock.Arguments) { etxs[0].ItHappened() })

		registryMock := cltest.NewContractMockReceiver(t, ethMock, keeper.Registry1_1ABI, registry.ContractAddress.Address())
		registryMock.MockResponse("checkUpkeep", checkUpkeepResponse)

		// turn falls somewhere between 20-39 (blockCountPerTurn=20)
		// heads 20 thru 35 were skipped (e.g. due to node reboot)
		head := cltest.Head(36)

		executer.OnNewLongestChain(context.Background(), head)
		runs := cltest.WaitForPipelineComplete(t, 0, job.ID, 1, 5, jpv2.Jrm, time.Second, 100*time.Millisecond)
		require.Len(t, runs, 1)
		assert.False(t, runs[0].HasErrors())
		etxs[0].AwaitOrFail(t)
		waitLastRunHeight(t, db, upkeep, 36)

		// heads 37, 38 etc do nothing
		for i := 37; i < 40; i++ {
			head = cltest.Head(i)
			executer.OnNewLongestChain(context.Background(), head)
		}

		// head 40 triggers a new run
		head = cltest.Head(40)

		txm.On("CreateEthTransaction",
			mock.MatchedBy(func(newTx txmgr.NewTx) bool { return newTx.GasLimit == gasLimit }),
		).
			Once().
			Return(txmgr.EthTx{}, nil).
			Run(func(mock.Arguments) { etxs[1].ItHappened() })

		executer.OnNewLongestChain(context.Background(), head)
		runs = cltest.WaitForPipelineComplete(t, 0, job.ID, 2, 5, jpv2.Jrm, time.Second, 100*time.Millisecond)
		require.Len(t, runs, 2)
		assert.False(t, runs[1].HasErrors())
		etxs[1].AwaitOrFail(t)
		waitLastRunHeight(t, db, upkeep, 40)

		ethMock.AssertExpectations(t)
	})
}

func Test_UpkeepExecuter_PerformsUpkeep_Error(t *testing.T) {
	t.Parallel()
	g := gomega.NewWithT(t)

	db, _, ethMock, executer, registry, _, _, _, _, _, _, _ := setup(t)

	wasCalled := atomic.NewBool(false)
	registryMock := cltest.NewContractMockReceiver(t, ethMock, keeper.Registry1_1ABI, registry.ContractAddress.Address())
	registryMock.MockRevertResponse("checkUpkeep").Run(func(args mock.Arguments) {
		wasCalled.Store(true)
	})

	head := newHead()
	executer.OnNewLongestChain(testutils.Context(t), &head)

	g.Eventually(wasCalled.Load).Should(gomega.Equal(true))
	cltest.AssertCountStays(t, db, "eth_txes", 0)
	ethMock.AssertExpectations(t)
}
