package keeper_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/services/keeper"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setup(t *testing.T) (
	*store.Store,
	*mocks.Client,
	keeper.UpkeepExecuter,
	keeper.Registry,
	func(),
) {
	store, strCleanup := cltest.NewStore(t)
	ethMock := new(mocks.Client)
	keeperORM := keeper.NewORM(store.ORM)
	executer := keeper.NewUpkeepExecuter(keeperORM, ethMock)
	registry := cltest.MustInsertKeeperRegistry(t, store)
	err := executer.Start()
	require.NoError(t, err)
	cleanup := func() { executer.Close(); strCleanup() }
	return store, ethMock, executer, registry, cleanup
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
	_, _, executer, _, cleanup := setup(t)
	defer cleanup()

	err := executer.Start() // already started in setup()
	require.Error(t, err)
}

func Test_UpkeepExecuter_PerformsUpkeep_Happy(t *testing.T) {
	store, ethMock, executer, registry, cleanup := setup(t)
	defer cleanup()

	upkeep := newUpkeep(registry, 0)
	err := store.DB.Create(&upkeep).Error
	require.NoError(t, err)

	registryMock := cltest.NewContractMockReceiver(t, ethMock, keeper.RegistryABI, registry.ContractAddress.Address())
	registryMock.MockResponse("checkUpkeep", checkUpkeepResponse)

	t.Run("runs upkeep on triggering block number", func(t *testing.T) {
		head := models.NewHead(big.NewInt(20), cltest.NewHash(), cltest.NewHash(), 1000)
		executer.OnNewLongestChain(context.TODO(), head)
		cltest.WaitForCount(t, store, models.EthTx{}, 1)
	})

	t.Run("skips upkeep on non-triggering block number", func(t *testing.T) {
		head := models.NewHead(big.NewInt(21), cltest.NewHash(), cltest.NewHash(), 1000)
		executer.OnNewLongestChain(context.TODO(), head)
		cltest.AssertCountStays(t, store, models.EthTx{}, 1)
	})

	ethMock.AssertExpectations(t)
}

func Test_UpkeepExecuter_PerformsUpkeep_Error(t *testing.T) {
	store, ethMock, executer, registry, cleanup := setup(t)
	defer cleanup()

	upkeep := newUpkeep(registry, 0)
	err := store.DB.Create(&upkeep).Error
	require.NoError(t, err)

	chUpkeepCalled := make(chan struct{})
	ethMock.
		On("CallContract", mock.Anything, mock.Anything, mock.Anything).
		Return(nil, errors.New("contract call revert")).
		Run(func(args mock.Arguments) {
			chUpkeepCalled <- struct{}{}
		})

	head := models.NewHead(big.NewInt(20), cltest.NewHash(), cltest.NewHash(), 1000)
	executer.OnNewLongestChain(context.TODO(), head)

	select {
	case <-time.NewTimer(5 * time.Second).C:
		t.Fatal("checkUpkeep never called")
	case <-chUpkeepCalled:
	}

	ethMock.AssertExpectations(t)
}
