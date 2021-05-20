package keeper_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/services/headtracker"
	"github.com/stretchr/testify/assert"

	"github.com/ethereum/go-ethereum/common"
	"github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keeper"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/atomic"
)

func setup(t *testing.T) (
	*store.Store,
	*mocks.Client,
	*keeper.UpkeepExecuter,
	keeper.Registry,
	keeper.UpkeepRegistration,
	job.Job,
	cltest.JobPipelineV2TestHelper,
) {
	config, cfgCleanup := cltest.NewConfig(t)
	t.Cleanup(cfgCleanup)
	config.Set("KEEPER_MAXIMUM_GRACE_PERIOD", 0)
	store, strCleanup := cltest.NewStoreWithConfig(t, config)
	t.Cleanup(strCleanup)
	ethMock := new(mocks.Client)
	registry, job := cltest.MustInsertKeeperRegistry(t, store)
	jpv2 := cltest.NewJobPipelineV2(t, store.DB)
	headBroadcaster := headtracker.NewHeadBroadcaster()
	executer := keeper.NewUpkeepExecuter(job, store.DB, jpv2.Pr, ethMock, headBroadcaster, store.Config)
	upkeep := cltest.MustInsertUpkeepForRegistry(t, store, registry)
	err := executer.Start()
	t.Cleanup(func() { executer.Close() })
	require.NoError(t, err)
	return store, ethMock, executer, registry, upkeep, job, jpv2
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
	_, _, executer, _, _, _, _ := setup(t)
	err := executer.Start() // already started in setup()
	require.Error(t, err)
}

func Test_UpkeepExecuter_PerformsUpkeep_Happy(t *testing.T) {
	t.Parallel()

	t.Run("runs upkeep on triggering block number", func(t *testing.T) {
		store, ethMock, executer, registry, upkeep, job, jpv2 := setup(t)

		registryMock := cltest.NewContractMockReceiver(t, ethMock, keeper.RegistryABI, registry.ContractAddress.Address())
		registryMock.MockResponse("checkUpkeep", checkUpkeepResponse)

		head := models.NewHead(big.NewInt(20), cltest.NewHash(), cltest.NewHash(), 1000)
		executer.OnNewLongestChain(context.Background(), head)
		cltest.WaitForCount(t, store, models.EthTx{}, 1)
		assertLastRunHeight(t, store, upkeep, 20)
		runs := cltest.WaitForPipelineComplete(t, 0, job.ID, 1, 0, jpv2.Jrm, time.Second, 100*time.Millisecond)
		require.Len(t, runs, 1)
		_, ok := runs[0].Meta.Val.(map[string]interface{})["eth_tx_id"]
		assert.True(t, ok)

		ethMock.AssertExpectations(t)
	})

	t.Run("triggers exactly one upkeep if heads are skipped but later heads arrive within range", func(t *testing.T) {
		store, ethMock, executer, registry, upkeep, job, jpv2 := setup(t)

		registryMock := cltest.NewContractMockReceiver(t, ethMock, keeper.RegistryABI, registry.ContractAddress.Address())
		registryMock.MockResponse("checkUpkeep", checkUpkeepResponse)

		// turn falls somewhere between 20-39 (blockCountPerTurn=20)
		// heads 20 thru 35 were skipped (e.g. due to node reboot)
		head := *cltest.Head(36)

		executer.OnNewLongestChain(context.Background(), head)
		cltest.WaitForCount(t, store, models.EthTx{}, 1)
		runs := cltest.WaitForPipelineComplete(t, 0, job.ID, 1, 0, jpv2.Jrm, time.Second, 100*time.Millisecond)
		assertLastRunHeight(t, store, upkeep, 36)
		require.Len(t, runs, 1)
		_, ok := runs[0].Meta.Val.(map[string]interface{})["eth_tx_id"]
		assert.True(t, ok)

		// heads 37, 38 etc do nothing
		for i := 37; i < 40; i++ {
			head = *cltest.Head(i)
			executer.OnNewLongestChain(context.Background(), head)
			cltest.AssertCountStays(t, store, models.EthTx{}, 1)
		}

		// head 40 triggers a new run
		head = *cltest.Head(40)

		executer.OnNewLongestChain(context.Background(), head)
		cltest.WaitForCount(t, store, models.EthTx{}, 2)
		assertLastRunHeight(t, store, upkeep, 40)
		runs = cltest.WaitForPipelineComplete(t, 0, job.ID, 1, 0, jpv2.Jrm, time.Second, 100*time.Millisecond)
		require.Len(t, runs, 2)
		_, ok = runs[0].Meta.Val.(map[string]interface{})["eth_tx_id"]
		assert.True(t, ok)

		ethMock.AssertExpectations(t)
	})

}

func Test_UpkeepExecuter_PerformsUpkeep_Error(t *testing.T) {
	t.Parallel()
	g := gomega.NewGomegaWithT(t)

	store, ethMock, executer, registry, _, _, _ := setup(t)

	wasCalled := atomic.NewBool(false)
	registryMock := cltest.NewContractMockReceiver(t, ethMock, keeper.RegistryABI, registry.ContractAddress.Address())
	registryMock.MockRevertResponse("checkUpkeep").Run(func(args mock.Arguments) {
		wasCalled.Store(true)
	})

	head := models.NewHead(big.NewInt(20), cltest.NewHash(), cltest.NewHash(), 1000)
	executer.OnNewLongestChain(context.TODO(), head)

	g.Eventually(wasCalled).Should(gomega.Equal(atomic.NewBool(true)))
	cltest.AssertCountStays(t, store, models.EthTx{}, 0)
	ethMock.AssertExpectations(t)
}
