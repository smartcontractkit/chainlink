package offchainreporting_test

import (
	"math"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/services/offchainreporting"
	"github.com/smartcontractkit/chainlink/core/utils"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type configOverriderUni struct {
	overrider       *offchainreporting.ConfigOverriderImpl
	contractAddress ethkey.EIP55Address
}

func newConfigOverriderUni(t *testing.T, pollInterval time.Duration, flagsContract *mocks.Flags) (uni configOverriderUni) {
	contractAddress := cltest.NewEIP55Address()

	flags := &offchainreporting.ContractFlags{FlagsInterface: flagsContract}
	var err error
	uni.overrider, err = offchainreporting.NewConfigOverriderImpl(
		logger.Default,
		contractAddress,
		flags,
		pollInterval,
	)
	require.NoError(t, err)

	uni.contractAddress = contractAddress

	t.Cleanup(func() {
		flagsContract.AssertExpectations(t)
	})

	return uni
}

func Test_OCRConfigOverrider(t *testing.T) {
	t.Parallel()
	g := gomega.NewGomegaWithT(t)

	t.Run("Before first tick returns nil override, later does return a specific override when hibernating", func(t *testing.T) {
		flagsContract := new(mocks.Flags)
		flagsContract.Test(t)
		uni := newConfigOverriderUni(t, 2*time.Second, flagsContract)

		// not hibernating, because one of the flags is lowered
		flagsContract.On("GetFlags", mock.Anything, mock.Anything).
			Run(checkFlagsAddress(t, uni.contractAddress)).
			Return([]bool{false, true}, nil).Once()

		// hibernating
		flagsContract.On("GetFlags", mock.Anything, mock.Anything).
			Run(checkFlagsAddress(t, uni.contractAddress)).
			Return([]bool{true, true}, nil)

		require.NoError(t, uni.overrider.Start())
		g.Consistently(func() *ocrtypes.ConfigOverride { return uni.overrider.ConfigOverride() }, time.Second, 450*time.Millisecond).Should(gomega.BeNil())

		res := &ocrtypes.ConfigOverride{AlphaPPB: math.MaxUint64, DeltaC: uni.overrider.DeltaCFromAddress}
		g.Eventually(func() *ocrtypes.ConfigOverride { return uni.overrider.ConfigOverride() }, 5*time.Second, 450*time.Millisecond).Should(gomega.Equal(res))
	})

	t.Run("Before first tick is hibernating, later exists hibernation", func(t *testing.T) {
		flagsContract := new(mocks.Flags)
		flagsContract.Test(t)
		uni := newConfigOverriderUni(t, 2*time.Second, flagsContract)

		// hibernating
		flagsContract.On("GetFlags", mock.Anything, mock.Anything).
			Run(checkFlagsAddress(t, uni.contractAddress)).
			Return([]bool{true, true}, nil).Once()

		// not hibernating
		flagsContract.On("GetFlags", mock.Anything, mock.Anything).
			Run(checkFlagsAddress(t, uni.contractAddress)).
			Return([]bool{true, false}, nil)

		require.NoError(t, uni.overrider.Start())

		res := &ocrtypes.ConfigOverride{AlphaPPB: math.MaxUint64, DeltaC: uni.overrider.DeltaCFromAddress}
		g.Consistently(func() *ocrtypes.ConfigOverride { return uni.overrider.ConfigOverride() }, time.Second, 450*time.Millisecond).Should(gomega.Equal(res))

		g.Eventually(func() *ocrtypes.ConfigOverride { return uni.overrider.ConfigOverride() }, 5*time.Second, 450*time.Millisecond).Should(gomega.BeNil())
	})

	t.Run("Errors if flags contract is missing", func(t *testing.T) {
		contractAddress := cltest.NewEIP55Address()
		flags := &offchainreporting.ContractFlags{FlagsInterface: nil}
		_, err := offchainreporting.NewConfigOverriderImpl(
			logger.Default,
			contractAddress,
			flags,
			5*time.Second,
		)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Flags contract instance is missing, the contract does not exist")
	})

	t.Run("DeltaC should be stable per address", func(t *testing.T) {
		flagsContract := new(mocks.Flags)
		flagsContract.Test(t)
		flagsContract.On("GetFlags", mock.Anything, mock.Anything).
			Return([]bool{true, true}, nil)
		flags := &offchainreporting.ContractFlags{FlagsInterface: flagsContract}

		address1, err := ethkey.NewEIP55Address(common.BigToAddress(big.NewInt(10000)).Hex())
		require.NoError(t, err)

		address2, err := ethkey.NewEIP55Address(common.BigToAddress(big.NewInt(1234567890)).Hex())
		require.NoError(t, err)

		overrider1a, err := offchainreporting.NewConfigOverriderImpl(logger.Default, address1, flags, 5*time.Second)
		require.NoError(t, err)

		overrider1b, err := offchainreporting.NewConfigOverriderImpl(logger.Default, address1, flags, 5*time.Second)
		require.NoError(t, err)

		overrider2, err := offchainreporting.NewConfigOverriderImpl(logger.Default, address2, flags, 5*time.Second)
		require.NoError(t, err)

		require.Equal(t, overrider1a.DeltaCFromAddress, time.Duration(85600000000000))
		require.Equal(t, overrider1b.DeltaCFromAddress, time.Duration(85600000000000))
		require.Equal(t, overrider2.DeltaCFromAddress, time.Duration(84690000000000))
	})
}

func checkFlagsAddress(t *testing.T, contractAddress ethkey.EIP55Address) func(args mock.Arguments) {
	return func(args mock.Arguments) {
		require.Equal(t, []common.Address{
			utils.ZeroAddress,
			contractAddress.Address(),
		}, args.Get(1).([]common.Address))
	}
}
