package offchainreporting_test

import (
	"math"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/logger"
	htmocks "github.com/smartcontractkit/chainlink/core/services/headtracker/mocks"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	logmocks "github.com/smartcontractkit/chainlink/core/services/log/mocks"
	"github.com/smartcontractkit/chainlink/core/services/offchainreporting"
	"github.com/smartcontractkit/chainlink/core/utils"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type configOverriderUni struct {
	lb              *logmocks.Broadcaster
	hb              *htmocks.HeadBroadcaster
	ec              *mocks.Client
	overrider       *offchainreporting.ConfigOverriderImpl
	contractAddress ethkey.EIP55Address
}

func newConfigOverriderUni(t *testing.T, pollInterval time.Duration, flagsContract *mocks.Flags) (uni configOverriderUni) {
	uni.lb = new(logmocks.Broadcaster)
	uni.hb = new(htmocks.HeadBroadcaster)
	uni.ec = cltest.NewEthClientMock(t)
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
		uni.lb.AssertExpectations(t)
		uni.hb.AssertExpectations(t)
		uni.ec.AssertExpectations(t)
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

		res := &ocrtypes.ConfigOverride{AlphaPPB: math.MaxUint64, DeltaC: 23*time.Hour + time.Duration(uni.contractAddress.Big().Int64()%3600)*time.Second}
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

		res := &ocrtypes.ConfigOverride{AlphaPPB: math.MaxUint64, DeltaC: 23*time.Hour + time.Duration(uni.contractAddress.Big().Int64()%3600)*time.Second}
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
}

func checkFlagsAddress(t *testing.T, contractAddress ethkey.EIP55Address) func(args mock.Arguments) {
	return func(args mock.Arguments) {
		require.Equal(t, []common.Address{
			utils.ZeroAddress,
			contractAddress.Address(),
		}, args.Get(1).([]common.Address))
	}
}
