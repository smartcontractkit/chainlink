package keeper_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/atomic"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	bptxmmocks "github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager/mocks"
	"github.com/smartcontractkit/chainlink/core/services/headtracker"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keeper"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func setup(t *testing.T) (
	*store.Store,
	*mocks.Client,
	*keeper.UpkeepExecuter,
	keeper.Registry,
	keeper.UpkeepRegistration,
	job.Job,
	cltest.JobPipelineV2TestHelper,
	*bptxmmocks.TxManager,
) {
	config := cltest.NewTestEVMConfig(t)
	config.GeneralConfig.Overrides.KeeperMaximumGracePeriod = null.IntFrom(0)
	store, strCleanup := cltest.NewStoreWithConfig(t, config)
	t.Cleanup(strCleanup)
	keyStore := cltest.NewKeyStore(t, store.DB)
	ethClient := cltest.NewEthClientMock(t)
	registry, job := cltest.MustInsertKeeperRegistry(t, store, keyStore.Eth())
	cfg := cltest.NewTestEVMConfig(t)
	jpv2 := cltest.NewJobPipelineV2(t, cfg, store.DB, nil, keyStore, nil)
	headBroadcaster := headtracker.NewHeadBroadcaster(logger.Default)
	txm := new(bptxmmocks.TxManager)
	orm := keeper.NewORM(store.DB, txm, store.Config, bulletprooftxmanager.SendEveryStrategy{})
	executer := keeper.NewUpkeepExecuter(job, orm, jpv2.Pr, ethClient, headBroadcaster, logger.Default, store.Config)
	upkeep := cltest.MustInsertUpkeepForRegistry(t, store, registry)
	err := executer.Start()
	t.Cleanup(func() { executer.Close() })
	require.NoError(t, err)
	return store, ethClient, executer, registry, upkeep, job, jpv2, txm
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
	_, _, executer, _, _, _, _, _ := setup(t)
	err := executer.Start() // already started in setup()
	require.Error(t, err)
}

func Test_UpkeepExecuter_PerformsUpkeep_Happy(t *testing.T) {
	t.Parallel()

	t.Run("runs upkeep on triggering block number", func(t *testing.T) {
		store, ethMock, executer, registry, upkeep, job, jpv2, txm := setup(t)

		gasLimit := upkeep.ExecuteGas + store.Config.KeeperRegistryPerformGasOverhead()
		ethTxCreated := cltest.NewAwaiter()
		txm.On("CreateEthTransaction",
			mock.Anything, mock.MatchedBy(func(newTx bulletprooftxmanager.NewTx) bool { return newTx.GasLimit == gasLimit }),
		).
			Once().
			Return(bulletprooftxmanager.EthTx{}, nil).
			Run(func(mock.Arguments) { ethTxCreated.ItHappened() })

		registryMock := cltest.NewContractMockReceiver(t, ethMock, keeper.RegistryABI, registry.ContractAddress.Address())
		registryMock.MockResponse("checkUpkeep", checkUpkeepResponse)

		head := models.NewHead(big.NewInt(20), utils.NewHash(), utils.NewHash(), 1000)
		executer.OnNewLongestChain(context.Background(), head)
		ethTxCreated.AwaitOrFail(t)
		assertLastRunHeight(t, store, upkeep, 20)
		runs := cltest.WaitForPipelineComplete(t, 0, job.ID, 1, 0, jpv2.Jrm, time.Second, 100*time.Millisecond)
		require.Len(t, runs, 1)
		_, ok := runs[0].Meta.Val.(map[string]interface{})["eth_tx_id"]
		assert.True(t, ok)

		ethMock.AssertExpectations(t)
		txm.AssertExpectations(t)
	})

	t.Run("triggers exactly one upkeep if heads are skipped but later heads arrive within range", func(t *testing.T) {
		store, ethMock, executer, registry, upkeep, job, jpv2, txm := setup(t)

		etxs := []cltest.Awaiter{
			cltest.NewAwaiter(),
			cltest.NewAwaiter(),
		}
		gasLimit := upkeep.ExecuteGas + store.Config.KeeperRegistryPerformGasOverhead()
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
		etxs[0].AwaitOrFail(t)
		runs := cltest.WaitForPipelineComplete(t, 0, job.ID, 1, 0, jpv2.Jrm, time.Second, 100*time.Millisecond)
		assertLastRunHeight(t, store, upkeep, 36)
		require.Len(t, runs, 1)
		_, ok := runs[0].Meta.Val.(map[string]interface{})["eth_tx_id"]
		assert.True(t, ok)

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
		etxs[1].AwaitOrFail(t)
		assertLastRunHeight(t, store, upkeep, 40)
		runs = cltest.WaitForPipelineComplete(t, 0, job.ID, 2, 0, jpv2.Jrm, time.Second, 100*time.Millisecond)
		require.Len(t, runs, 2)
		_, ok = runs[0].Meta.Val.(map[string]interface{})["eth_tx_id"]
		assert.True(t, ok)

		ethMock.AssertExpectations(t)
		txm.AssertExpectations(t)
	})
}

func Test_UpkeepExecuter_PerformsUpkeep_Error(t *testing.T) {
	t.Parallel()
	g := gomega.NewGomegaWithT(t)

	store, ethMock, executer, registry, _, _, _, _ := setup(t)

	wasCalled := atomic.NewBool(false)
	registryMock := cltest.NewContractMockReceiver(t, ethMock, keeper.RegistryABI, registry.ContractAddress.Address())
	registryMock.MockRevertResponse("checkUpkeep").Run(func(args mock.Arguments) {
		wasCalled.Store(true)
	})

	head := models.NewHead(big.NewInt(20), utils.NewHash(), utils.NewHash(), 1000)
	executer.OnNewLongestChain(context.TODO(), head)

	g.Eventually(wasCalled).Should(gomega.Equal(atomic.NewBool(true)))
	cltest.AssertCountStays(t, store, bulletprooftxmanager.EthTx{}, 0)
	ethMock.AssertExpectations(t)
}

func Test_UpkeepExecuter_ConstructCheckUpkeepCallMsg(t *testing.T) {
	store, _, executer, registry, upkeep, _, _, _ := setup(t)
	msg, err := executer.ExportedConstructCheckUpkeepCallMsg(upkeep)
	require.NoError(t, err)
	expectedGasLimit := upkeep.ExecuteGas + uint64(registry.CheckGas) + store.Config.KeeperRegistryCheckGasOverhead() + store.Config.KeeperRegistryPerformGasOverhead()
	require.Equal(t, expectedGasLimit, msg.Gas)
	require.Equal(t, registry.ContractAddress.Address(), *msg.To)
	require.Equal(t, utils.ZeroAddress, msg.From)
	require.Nil(t, msg.Value)
}
