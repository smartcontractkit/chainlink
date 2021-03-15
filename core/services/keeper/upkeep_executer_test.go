package keeper_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/onsi/gomega"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
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
	*keeper.UpkeepExecutor,
	keeper.Registry,
	func(),
) {
	store, strCleanup := cltest.NewStore(t)
	ethMock := new(mocks.Client)
	executor := keeper.NewUpkeepExecutor(store.DB, ethMock)
	registry := cltest.MustInsertKeeperRegistry(t, store)
	err := executor.Start()
	require.NoError(t, err)
	cleanup := func() { executor.Close(); strCleanup() }
	return store, ethMock, executor, registry, cleanup
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

func Test_UpkeepExecutor_ErrorsIfStartedTwice(t *testing.T) {
	t.Parallel()
	_, _, executor, _, cleanup := setup(t)
	defer cleanup()

	err := executor.Start() // already started in setup()
	require.Error(t, err)
}

func Test_UpkeepExecutor_PerformsUpkeep_Happy(t *testing.T) {
	t.Parallel()
	store, ethMock, executor, registry, cleanup := setup(t)
	defer cleanup()

	upkeep := newUpkeep(registry, 0)
	err := store.DB.Create(&upkeep).Error
	require.NoError(t, err)

	registryMock := cltest.NewContractMockReceiver(t, ethMock, keeper.RegistryABI, registry.ContractAddress.Address())
	registryMock.MockResponse("checkUpkeep", checkUpkeepResponse)

	t.Run("runs upkeep on triggering block number", func(t *testing.T) {
		head := models.NewHead(big.NewInt(20), cltest.NewHash(), cltest.NewHash(), 1000)
		executor.OnNewLongestChain(context.TODO(), head)
		cltest.WaitForCount(t, store, models.EthTx{}, 1)
	})

	t.Run("skips upkeep on non-triggering block number", func(t *testing.T) {
		head := models.NewHead(big.NewInt(21), cltest.NewHash(), cltest.NewHash(), 1000)
		executor.OnNewLongestChain(context.TODO(), head)
		cltest.AssertCountStays(t, store, models.EthTx{}, 1)
	})

	ethMock.AssertExpectations(t)
}

func Test_UpkeepExecutor_PerformsUpkeep_Error(t *testing.T) {
	t.Parallel()
	g := gomega.NewGomegaWithT(t)
	store, ethMock, executor, registry, cleanup := setup(t)
	defer cleanup()

	upkeep := newUpkeep(registry, 0)
	err := store.DB.Create(&upkeep).Error
	require.NoError(t, err)

	wasCalled := atomic.NewBool(false)
	ethMock.
		On("CallContract", mock.Anything, mock.Anything, mock.Anything).
		Return(nil, errors.New("contract call revert")).
		Run(func(args mock.Arguments) {
			wasCalled.Store(true)
		})

	head := models.NewHead(big.NewInt(20), cltest.NewHash(), cltest.NewHash(), 1000)
	executor.OnNewLongestChain(context.TODO(), head)

	g.Eventually(wasCalled).Should(gomega.Equal(atomic.NewBool(true)))
	ethMock.AssertExpectations(t)
}
