package allowlist

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/functions/generated/functions_allow_list"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	amocks "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/functions/allowlist/mocks"
)

func TestUpdateAllowedSendersInBatches(t *testing.T) {
	t.Run("OK-simple_update_in_batches", func(t *testing.T) {
		ctx := context.Background()
		config := OnchainAllowlistConfig{
			ContractAddress:           testutils.NewAddress(),
			ContractVersion:           1,
			BlockConfirmations:        1,
			UpdateFrequencySec:        2,
			UpdateTimeoutSec:          1,
			StoredAllowlistBatchSize:  2,
			OnchainAllowlistBatchSize: 10,
			FetchingDelayInRangeSec:   1,
		}

		// allowlistSize defines how big the mocked allowlist will be
		allowlistSize := 53
		// allowlist represents the actual allowlist the tos contract will return
		allowlist := make([]common.Address, 0, allowlistSize)
		// expectedAllowlist will be used to compare the actual status with what we actually want
		expectedAllowlist := make(map[common.Address]struct{}, 0)

		// we load both the expectedAllowlist and the allowlist the contract will return with some new addresses
		for i := 0; i < allowlistSize; i++ {
			addr := testutils.NewAddress()
			allowlist = append(allowlist, addr)
			expectedAllowlist[addr] = struct{}{}
		}

		tosContract := NewTosContractMock(allowlist)

		// with the orm mock we can validate the actual order in which the allowlist is fetched giving priority to newest addresses
		orm := amocks.NewORM(t)
		firstCall := orm.On("CreateAllowedSenders", context.Background(), allowlist[43:53]).Times(1).Return(nil)
		secondCall := orm.On("CreateAllowedSenders", context.Background(), allowlist[33:43]).Times(1).Return(nil).NotBefore(firstCall)
		thirdCall := orm.On("CreateAllowedSenders", context.Background(), allowlist[23:33]).Times(1).Return(nil).NotBefore(secondCall)
		forthCall := orm.On("CreateAllowedSenders", context.Background(), allowlist[13:23]).Times(1).Return(nil).NotBefore(thirdCall)
		fifthCall := orm.On("CreateAllowedSenders", context.Background(), allowlist[3:13]).Times(1).Return(nil).NotBefore(forthCall)
		orm.On("CreateAllowedSenders", context.Background(), allowlist[0:3]).Times(1).Return(nil).NotBefore(fifthCall)

		onchainAllowlist := &onchainAllowlist{
			config:             config,
			orm:                orm,
			blockConfirmations: big.NewInt(int64(config.BlockConfirmations)),
			lggr:               logger.TestLogger(t).Named("OnchainAllowlist"),
			stopCh:             make(services.StopChan),
		}

		// we set the onchain allowlist to an empty state before updating it in batches
		emptyMap := make(map[common.Address]struct{})
		onchainAllowlist.allowlist.Store(&emptyMap)

		err := onchainAllowlist.updateAllowedSendersInBatches(ctx, tosContract, big.NewInt(0))
		require.NoError(t, err)

		currentAllowlist := onchainAllowlist.allowlist.Load()
		require.Equal(t, &expectedAllowlist, currentAllowlist)
	})

	t.Run("OK-new_address_added_while_updating_in_batches", func(t *testing.T) {
		ctx := context.Background()
		config := OnchainAllowlistConfig{
			ContractAddress:           testutils.NewAddress(),
			ContractVersion:           1,
			BlockConfirmations:        1,
			UpdateFrequencySec:        2,
			UpdateTimeoutSec:          1,
			StoredAllowlistBatchSize:  2,
			OnchainAllowlistBatchSize: 10,
			FetchingDelayInRangeSec:   1,
		}

		// allowlistSize defines how big the initial mocked allowlist will be
		allowlistSize := 50
		// allowlist represents the actual allowlist the tos contract will return
		allowlist := make([]common.Address, 0)
		// expectedAllowlist will be used to compare the actual status with what we actually want
		expectedAllowlist := make(map[common.Address]struct{}, 0)

		// we load both the expectedAllowlist and the allowlist the contract will return with some new addresses
		for i := 0; i < allowlistSize; i++ {
			addr := testutils.NewAddress()
			allowlist = append(allowlist, addr)
			expectedAllowlist[addr] = struct{}{}
		}

		tosContract := NewTosContractMock(allowlist)

		// with the orm mock we can validate the actual order in which the allowlist is fetched giving priority to newest addresses
		orm := amocks.NewORM(t)
		firstCall := orm.On("CreateAllowedSenders", context.Background(), allowlist[40:50]).Times(1).Run(func(args mock.Arguments) {
			// after the first call we update the tosContract by adding a new address
			addr := testutils.NewAddress()
			allowlist = append(allowlist, addr)
			expectedAllowlist[addr] = struct{}{}
			*tosContract = *NewTosContractMock(allowlist)
		}).Return(nil)

		// this is the extra step that will fetch the new address we want to validate
		extraStepCall := orm.On("CreateAllowedSenders", context.Background(), allowlist[50:51]).Times(1).Return(nil).NotBefore(firstCall)

		secondCall := orm.On("CreateAllowedSenders", context.Background(), allowlist[30:40]).Times(1).Return(nil).NotBefore(extraStepCall)
		thirdCall := orm.On("CreateAllowedSenders", context.Background(), allowlist[20:30]).Times(1).Return(nil).NotBefore(secondCall)
		forthCall := orm.On("CreateAllowedSenders", context.Background(), allowlist[10:20]).Times(1).Return(nil).NotBefore(thirdCall)
		orm.On("CreateAllowedSenders", context.Background(), allowlist[0:10]).Times(1).Return(nil).NotBefore(forthCall)

		onchainAllowlist := &onchainAllowlist{
			config:             config,
			orm:                orm,
			blockConfirmations: big.NewInt(int64(config.BlockConfirmations)),
			lggr:               logger.TestLogger(t).Named("OnchainAllowlist"),
			stopCh:             make(services.StopChan),
		}

		// we set the onchain allowlist to an empty state before updating it in batches
		emptyMap := make(map[common.Address]struct{})
		onchainAllowlist.allowlist.Store(&emptyMap)

		err := onchainAllowlist.updateAllowedSendersInBatches(ctx, tosContract, big.NewInt(0))
		require.NoError(t, err)

		currentAllowlist := onchainAllowlist.allowlist.Load()
		require.Equal(t, &expectedAllowlist, currentAllowlist)
	})

	t.Run("OK-allowlist_size_smaller_than_batchsize", func(t *testing.T) {
		ctx := context.Background()
		config := OnchainAllowlistConfig{
			ContractAddress:           testutils.NewAddress(),
			ContractVersion:           1,
			BlockConfirmations:        1,
			UpdateFrequencySec:        2,
			UpdateTimeoutSec:          1,
			StoredAllowlistBatchSize:  2,
			OnchainAllowlistBatchSize: 100,
			FetchingDelayInRangeSec:   1,
		}

		// allowlistSize defines how big the mocked allowlist will be
		allowlistSize := 50
		// allowlist represents the actual allowlist the tos contract will return
		allowlist := make([]common.Address, 0, allowlistSize)
		// expectedAllowlist will be used to compare the actual status with what we actually want
		expectedAllowlist := make(map[common.Address]struct{}, 0)

		// we load both the expectedAllowlist and the allowlist the contract will return with some new addresses
		for i := 0; i < allowlistSize; i++ {
			addr := testutils.NewAddress()
			allowlist = append(allowlist, addr)
			expectedAllowlist[addr] = struct{}{}
		}

		tosContract := NewTosContractMock(allowlist)

		// with the orm mock we can validate the actual order in which the allowlist is fetched giving priority to newest addresses
		orm := amocks.NewORM(t)
		orm.On("CreateAllowedSenders", context.Background(), allowlist[0:50]).Times(1).Return(nil)

		onchainAllowlist := &onchainAllowlist{
			config:             config,
			orm:                orm,
			blockConfirmations: big.NewInt(int64(config.BlockConfirmations)),
			lggr:               logger.TestLogger(t).Named("OnchainAllowlist"),
			stopCh:             make(services.StopChan),
		}

		// we set the onchain allowlist to an empty state before updating it in batches
		emptyMap := make(map[common.Address]struct{})
		onchainAllowlist.allowlist.Store(&emptyMap)

		err := onchainAllowlist.updateAllowedSendersInBatches(ctx, tosContract, big.NewInt(0))
		require.NoError(t, err)

		currentAllowlist := onchainAllowlist.allowlist.Load()
		require.Equal(t, &expectedAllowlist, currentAllowlist)
	})
}

type tosContractMock struct {
	functions_allow_list.TermsOfServiceAllowListInterface

	onchainAllowlist []common.Address
}

func NewTosContractMock(onchainAllowlist []common.Address) *tosContractMock {
	return &tosContractMock{
		onchainAllowlist: onchainAllowlist,
	}
}

func (t *tosContractMock) GetAllowedSendersCount(opts *bind.CallOpts) (uint64, error) {
	return uint64(len(t.onchainAllowlist)), nil
}

func (t *tosContractMock) GetAllowedSendersInRange(opts *bind.CallOpts, allowedSenderIdxStart uint64, allowedSenderIdxEnd uint64) ([]common.Address, error) {
	// we replicate the onchain behaviour of including start and end indexes
	return t.onchainAllowlist[allowedSenderIdxStart : allowedSenderIdxEnd+1], nil
}
