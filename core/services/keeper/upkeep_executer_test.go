package keeper_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/atomic"
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	bptxmmocks "github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager/mocks"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	gasmocks "github.com/smartcontractkit/chainlink/core/services/gas/mocks"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keeper"
	"github.com/smartcontractkit/chainlink/core/utils"
	bigmath "github.com/smartcontractkit/chainlink/core/utils/big_math"
)

func newHead() eth.Head {
	return eth.NewHead(big.NewInt(20), utils.NewHash(), utils.NewHash(), 1000, utils.NewBigI(0))
}

func setup(t *testing.T) (
	*gorm.DB,
	*configtest.TestGeneralConfig,
	*mocks.Client,
	*keeper.UpkeepExecuter,
	keeper.Registry,
	keeper.UpkeepRegistration,
	job.Job,
	cltest.JobPipelineV2TestHelper,
	*bptxmmocks.TxManager,
) {
	cfg := cltest.NewTestGeneralConfig(t)
	cfg.Overrides.KeeperMaximumGracePeriod = null.IntFrom(0)
	db := pgtest.NewGormDB(t)
	cfg.SetDB(db)
	keyStore := cltest.NewKeyStore(t, db)
	ethClient := cltest.NewEthClientMockWithDefaultChain(t)
	registry, job := cltest.MustInsertKeeperRegistry(t, db, keyStore.Eth())
	txm := new(bptxmmocks.TxManager)
	txm.Test(t)
	estimator := new(gasmocks.Estimator)
	estimator.Test(t)
	txm.On("GetGasEstimator").Return(estimator)
	estimator.On("GetLegacyGas", mock.Anything, mock.Anything).Maybe().Return(assets.GWei(60), uint64(0), nil)
	cc := evmtest.NewChainSet(t, evmtest.TestChainOpts{TxManager: txm, DB: db, Client: ethClient, KeyStore: keyStore.Eth(), GeneralConfig: cfg})
	jpv2 := cltest.NewJobPipelineV2(t, cfg, cc, db, keyStore)
	ch := evmtest.MustGetDefaultChain(t, cc)
	orm := keeper.NewORM(db, txm, ch.Config(), bulletprooftxmanager.SendEveryStrategy{})
	lggr := logger.TestLogger(t)
	executer := keeper.NewUpkeepExecuter(job, orm, jpv2.Pr, ethClient, ch.HeadBroadcaster(), ch.TxManager().GetGasEstimator(), lggr, ch.Config())
	upkeep := cltest.MustInsertUpkeepForRegistry(t, db, ch.Config(), registry)
	err := executer.Start()
	t.Cleanup(func() { txm.AssertExpectations(t); estimator.AssertExpectations(t); executer.Close() })
	require.NoError(t, err)
	return db, cfg, ethClient, executer, registry, upkeep, job, jpv2, txm
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
	_, _, _, executer, _, _, _, _, _ := setup(t)
	err := executer.Start() // already started in setup()
	require.Error(t, err)
}

func Test_UpkeepExecuter_PerformsUpkeep_Happy(t *testing.T) {
	t.Parallel()

	t.Run("runs upkeep on triggering block number", func(t *testing.T) {
		db, config, ethMock, executer, registry, upkeep, job, jpv2, txm := setup(t)

		gasLimit := upkeep.ExecuteGas + config.KeeperRegistryPerformGasOverhead()
		gasPrice := bigmath.Div(bigmath.Mul(assets.GWei(60), 100+config.KeeperGasPriceBufferPercent()), 100)

		ethTxCreated := cltest.NewAwaiter()
		txm.On("CreateEthTransaction",
			mock.Anything, mock.MatchedBy(func(newTx bulletprooftxmanager.NewTx) bool { return newTx.GasLimit == gasLimit }),
		).
			Once().
			Return(bulletprooftxmanager.EthTx{
				ID: 1,
			}, nil).
			Run(func(mock.Arguments) { ethTxCreated.ItHappened() })

		registryMock := cltest.NewContractMockReceiver(t, ethMock, keeper.RegistryABI, registry.ContractAddress.Address())
		registryMock.MockMatchedResponse(
			"checkUpkeep",
			func(callArgs ethereum.CallMsg) bool {
				return bigmath.Equal(callArgs.GasPrice, gasPrice) &&
					callArgs.Gas == 650_000
			},
			checkUpkeepResponse,
		)

		head := newHead()
		executer.OnNewLongestChain(context.Background(), head)
		ethTxCreated.AwaitOrFail(t)
		runs := cltest.WaitForPipelineComplete(t, 0, job.ID, 1, 5, jpv2.Jrm, time.Second, 100*time.Millisecond)
		require.Len(t, runs, 1)
		assert.False(t, runs[0].HasErrors())
		assert.False(t, runs[0].HasFatalErrors())
		waitLastRunHeight(t, db, upkeep, 20)

		ethMock.AssertExpectations(t)
		txm.AssertExpectations(t)
	})

	t.Run("triggers exactly one upkeep if heads are skipped but later heads arrive within range", func(t *testing.T) {
		db, config, ethMock, executer, registry, upkeep, job, jpv2, txm := setup(t)

		etxs := []cltest.Awaiter{
			cltest.NewAwaiter(),
			cltest.NewAwaiter(),
		}
		gasLimit := upkeep.ExecuteGas + config.KeeperRegistryPerformGasOverhead()
		txm.On("CreateEthTransaction",
			mock.Anything, mock.MatchedBy(func(newTx bulletprooftxmanager.NewTx) bool { return newTx.GasLimit == gasLimit }),
		).
			Once().
			Return(bulletprooftxmanager.EthTx{}, nil).
			Run(func(mock.Arguments) { etxs[0].ItHappened() })

		registryMock := cltest.NewContractMockReceiver(t, ethMock, keeper.RegistryABI, registry.ContractAddress.Address())
		registryMock.MockResponse("checkUpkeep", checkUpkeepResponse)

		// turn falls somewhere between 20-39 (blockCountPerTurn=20)
		// heads 20 thru 35 were skipped (e.g. due to node reboot)
		head := *cltest.Head(36)

		executer.OnNewLongestChain(context.Background(), head)
		runs := cltest.WaitForPipelineComplete(t, 0, job.ID, 1, 5, jpv2.Jrm, time.Second, 100*time.Millisecond)
		require.Len(t, runs, 1)
		assert.False(t, runs[0].HasErrors())
		etxs[0].AwaitOrFail(t)
		waitLastRunHeight(t, db, upkeep, 36)

		// heads 37, 38 etc do nothing
		for i := 37; i < 40; i++ {
			head = *cltest.Head(i)
			executer.OnNewLongestChain(context.Background(), head)
		}

		// head 40 triggers a new run
		head = *cltest.Head(40)

		txm.On("CreateEthTransaction",
			mock.Anything, mock.MatchedBy(func(newTx bulletprooftxmanager.NewTx) bool { return newTx.GasLimit == gasLimit }),
		).
			Once().
			Return(bulletprooftxmanager.EthTx{}, nil).
			Run(func(mock.Arguments) { etxs[1].ItHappened() })

		executer.OnNewLongestChain(context.Background(), head)
		runs = cltest.WaitForPipelineComplete(t, 0, job.ID, 2, 5, jpv2.Jrm, time.Second, 100*time.Millisecond)
		require.Len(t, runs, 2)
		assert.False(t, runs[1].HasErrors())
		etxs[1].AwaitOrFail(t)
		waitLastRunHeight(t, db, upkeep, 40)

		ethMock.AssertExpectations(t)
		txm.AssertExpectations(t)
	})
}

func Test_UpkeepExecuter_PerformsUpkeep_Error(t *testing.T) {
	t.Parallel()
	g := gomega.NewGomegaWithT(t)

	db, _, ethMock, executer, registry, _, _, _, _ := setup(t)

	wasCalled := atomic.NewBool(false)
	registryMock := cltest.NewContractMockReceiver(t, ethMock, keeper.RegistryABI, registry.ContractAddress.Address())
	registryMock.MockRevertResponse("checkUpkeep").Run(func(args mock.Arguments) {
		wasCalled.Store(true)
	})

	head := newHead()
	executer.OnNewLongestChain(context.TODO(), head)

	g.Eventually(wasCalled.Load).Should(gomega.Equal(true))
	cltest.AssertCountStays(t, db, bulletprooftxmanager.EthTx{}, 0)
	ethMock.AssertExpectations(t)
}
