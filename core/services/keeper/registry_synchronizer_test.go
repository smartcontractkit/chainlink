package keeper_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/keeper_registry_wrapper"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keeper"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/stretchr/testify/require"
)

const syncInterval = 3 * time.Second

var regConfig = struct {
	PaymentPremiumPPB uint32
	BlockCountPerTurn *big.Int
	CheckGasLimit     uint32
	StalenessSeconds  *big.Int
	FallbackGasPrice  *big.Int
	FallbackLinkPrice *big.Int
}{
	PaymentPremiumPPB: 100,
	BlockCountPerTurn: big.NewInt(20),
	CheckGasLimit:     2_000_000,
	StalenessSeconds:  big.NewInt(3600),
	FallbackGasPrice:  big.NewInt(1000000),
	FallbackLinkPrice: big.NewInt(1000000),
}

var upkeep = struct {
	Target              common.Address
	ExecuteGas          uint32
	CheckData           []byte
	Balance             *big.Int
	LastKeeper          common.Address
	Admin               common.Address
	MaxValidBlocknumber uint64
}{
	Target:              cltest.NewAddress(),
	ExecuteGas:          2_000_000,
	CheckData:           common.Hex2Bytes("1234"),
	Balance:             big.NewInt(1000000000000000000),
	LastKeeper:          cltest.NewAddress(),
	Admin:               cltest.NewAddress(),
	MaxValidBlocknumber: 1_000_000_000,
}

func setupRegistrySync(t *testing.T) (*store.Store, *keeper.RegistrySynchronizer, *mocks.Client, job.Job, func()) {
	store, cleanup := cltest.NewStore(t)
	ethMock := new(mocks.Client)
	j := cltest.MustInsertKeeperJob(t, store, cltest.NewEIP55Address(), cltest.NewEIP55Address())
	contractAddress := j.KeeperSpec.ContractAddress
	contract, err := keeper_registry_wrapper.NewKeeperRegistry(
		contractAddress.Address(),
		ethMock,
	)
	require.NoError(t, err)

	synchronizer := keeper.NewRegistrySynchronizer(j, contract, store.DB, syncInterval)
	return store, synchronizer, ethMock, j, cleanup
}

func assertUpkeepIDs(t *testing.T, store *store.Store, expected []int32) {
	g := gomega.NewGomegaWithT(t)
	var upkeepIDs []int32
	err := store.DB.Model(keeper.UpkeepRegistration{}).Pluck("upkeep_id", &upkeepIDs).Error
	require.NoError(t, err)
	g.Expect(upkeepIDs).To(gomega.ContainElements(expected))
}

func Test_RegistrySynchronizer_Start(t *testing.T) {
	t.Parallel()
	_, synchronizer, _, _, cleanup := setupRegistrySync(t)
	defer cleanup()

	err := synchronizer.Start()
	require.NoError(t, err)
	defer synchronizer.Close()

	err = synchronizer.Start()
	require.Error(t, err)
}

func Test_RegistrySynchronizer_AddsAndRemovesUpkeeps(t *testing.T) {
	t.Parallel()
	store, synchronizer, ethMock, job, cleanup := setupRegistrySync(t)
	defer cleanup()

	contractAddress := job.KeeperSpec.ContractAddress.Address()
	fromAddress := job.KeeperSpec.FromAddress.Address()

	// 1st sync
	registryMock := cltest.NewContractMockReceiver(t, ethMock, keeper.RegistryABI, contractAddress)
	canceledUpkeeps := []*big.Int{big.NewInt(1)}
	registryMock.MockResponse("getConfig", regConfig).Once()
	registryMock.MockResponse("getKeeperList", []common.Address{fromAddress}).Once()
	registryMock.MockResponse("getCanceledUpkeepList", canceledUpkeeps).Once()
	registryMock.MockResponse("getUpkeepCount", big.NewInt(3)).Once()
	registryMock.MockResponse("getUpkeep", upkeep).Times(3) // sync all 3, then delete

	synchronizer.ExportedSyncRegistry()

	cltest.AssertCount(t, store, keeper.Registry{}, 1)
	cltest.AssertCount(t, store, keeper.UpkeepRegistration{}, 2)

	var registry keeper.Registry
	var upkeepRegistration keeper.UpkeepRegistration
	require.NoError(t, store.DB.First(&registry).Error)
	require.NoError(t, store.DB.First(&upkeepRegistration).Error)
	require.Equal(t, job.KeeperSpec.ContractAddress, registry.ContractAddress)
	require.Equal(t, job.KeeperSpec.FromAddress, registry.FromAddress)
	require.Equal(t, int32(20), registry.BlockCountPerTurn)
	require.Equal(t, int32(0), registry.KeeperIndex)
	require.Equal(t, int32(1), registry.NumKeepers)
	require.Equal(t, upkeep.CheckData, upkeepRegistration.CheckData)
	require.Equal(t, int32(upkeep.ExecuteGas), upkeepRegistration.ExecuteGas)

	assertUpkeepIDs(t, store, []int32{0, 2})
	ethMock.AssertExpectations(t)

	gomega.ContainElements()

	// 2nd sync
	canceledUpkeeps = []*big.Int{big.NewInt(0), big.NewInt(1), big.NewInt(3)}
	registryMock.MockResponse("getConfig", regConfig).Once()
	registryMock.MockResponse("getKeeperList", []common.Address{fromAddress}).Once()
	registryMock.MockResponse("getCanceledUpkeepList", canceledUpkeeps).Once()
	registryMock.MockResponse("getUpkeepCount", big.NewInt(5)).Once()
	registryMock.MockResponse("getUpkeep", upkeep).Times(2) // two new upkeeps to sync

	synchronizer.ExportedSyncRegistry()

	cltest.AssertCount(t, store, keeper.Registry{}, 1)
	cltest.AssertCount(t, store, keeper.UpkeepRegistration{}, 2)
	assertUpkeepIDs(t, store, []int32{2, 4})
	ethMock.AssertExpectations(t)
}
