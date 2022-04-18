package ocr_test

import (
	"math"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/onsi/gomega"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/services/ocr"
	"github.com/smartcontractkit/chainlink/core/utils"
)

type configOverriderUni struct {
	overrider       *ocr.ConfigOverriderImpl
	contractAddress ethkey.EIP55Address
}

func newConfigOverriderUni(t *testing.T, pollITicker utils.TickerBase, flagsContract *mocks.Flags) (uni configOverriderUni) {
	var testLogger = logger.TestLogger(t)
	contractAddress := cltest.NewEIP55Address()

	flags := &ocr.ContractFlags{FlagsInterface: flagsContract}
	var err error
	uni.overrider, err = ocr.NewConfigOverriderImpl(
		testLogger,
		contractAddress,
		flags,
		pollITicker,
	)
	require.NoError(t, err)

	uni.contractAddress = contractAddress

	t.Cleanup(func() {
		flagsContract.AssertExpectations(t)
	})

	return uni
}

func TestIntegration_OCRConfigOverrider_EntersHibernation(t *testing.T) {
	g := gomega.NewWithT(t)

	flagsContract := new(mocks.Flags)
	flagsContract.Test(t)

	ticker := utils.NewPausableTicker(3 * time.Second)
	uni := newConfigOverriderUni(t, &ticker, flagsContract)

	// not hibernating, because one of the flags is lowered
	flagsContract.On("GetFlags", mock.Anything, mock.Anything).
		Run(checkFlagsAddress(t, uni.contractAddress)).
		Return([]bool{false, true}, nil).Once()

	// hibernating
	flagsContract.On("GetFlags", mock.Anything, mock.Anything).
		Run(checkFlagsAddress(t, uni.contractAddress)).
		Return([]bool{true, true}, nil)

	require.NoError(t, uni.overrider.Start(testutils.Context(t)))

	// not hibernating initially
	require.Nil(t, uni.overrider.ConfigOverride())

	expectedOverride := &ocrtypes.ConfigOverride{AlphaPPB: math.MaxUint64, DeltaC: uni.overrider.DeltaCFromAddress}

	// timeout needs to be longer than the poll interval of 3 seconds
	g.Eventually(func() *ocrtypes.ConfigOverride { return uni.overrider.ConfigOverride() }, 10*time.Second, 450*time.Millisecond).Should(gomega.Equal(expectedOverride))
}

func Test_OCRConfigOverrider(t *testing.T) {
	t.Parallel()

	t.Run("Before first tick returns nil override, later does return a specific override when hibernating", func(t *testing.T) {
		flagsContract := new(mocks.Flags)
		flagsContract.Test(t)

		ticker := NewFakeTicker()
		uni := newConfigOverriderUni(t, ticker, flagsContract)

		// not hibernating, because one of the flags is lowered
		flagsContract.On("GetFlags", mock.Anything, mock.Anything).
			Run(checkFlagsAddress(t, uni.contractAddress)).
			Return([]bool{false, true}, nil).Once()

		// hibernating
		flagsContract.On("GetFlags", mock.Anything, mock.Anything).
			Run(checkFlagsAddress(t, uni.contractAddress)).
			Return([]bool{true, true}, nil)

		require.NoError(t, uni.overrider.Start(testutils.Context(t)))

		// not hibernating initially
		require.Nil(t, uni.overrider.ConfigOverride())

		// update state by getting flags
		require.NoError(t, uni.overrider.ExportedUpdateFlagsStatus())

		expectedOverride := &ocrtypes.ConfigOverride{AlphaPPB: math.MaxUint64, DeltaC: uni.overrider.DeltaCFromAddress}
		require.Equal(t, expectedOverride, uni.overrider.ConfigOverride())
	})

	t.Run("Before first tick is hibernating, later exists hibernation", func(t *testing.T) {
		flagsContract := new(mocks.Flags)
		flagsContract.Test(t)

		ticker := NewFakeTicker()
		uni := newConfigOverriderUni(t, ticker, flagsContract)

		// hibernating
		flagsContract.On("GetFlags", mock.Anything, mock.Anything).
			Run(checkFlagsAddress(t, uni.contractAddress)).
			Return([]bool{true, true}, nil).Once()

		// not hibernating
		flagsContract.On("GetFlags", mock.Anything, mock.Anything).
			Run(checkFlagsAddress(t, uni.contractAddress)).
			Return([]bool{true, false}, nil)

		require.NoError(t, uni.overrider.Start(testutils.Context(t)))

		// initially enters hibernation
		expectedOverride := &ocrtypes.ConfigOverride{AlphaPPB: math.MaxUint64, DeltaC: uni.overrider.DeltaCFromAddress}
		require.Equal(t, expectedOverride, uni.overrider.ConfigOverride())

		// update state by getting flags
		require.NoError(t, uni.overrider.ExportedUpdateFlagsStatus())

		// should exit hibernation
		require.Nil(t, uni.overrider.ConfigOverride())
	})

	t.Run("Errors if flags contract is missing", func(t *testing.T) {
		var testLogger = logger.TestLogger(t)
		contractAddress := cltest.NewEIP55Address()
		flags := &ocr.ContractFlags{FlagsInterface: nil}
		_, err := ocr.NewConfigOverriderImpl(
			testLogger,
			contractAddress,
			flags,
			nil,
		)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Flags contract instance is missing, the contract does not exist")
	})

	t.Run("DeltaC should be stable per address", func(t *testing.T) {
		var testLogger = logger.TestLogger(t)
		flagsContract := new(mocks.Flags)
		flagsContract.Test(t)
		flagsContract.On("GetFlags", mock.Anything, mock.Anything).
			Return([]bool{true, true}, nil)
		flags := &ocr.ContractFlags{FlagsInterface: flagsContract}

		address1, err := ethkey.NewEIP55Address(common.BigToAddress(big.NewInt(10000)).Hex())
		require.NoError(t, err)

		address2, err := ethkey.NewEIP55Address(common.BigToAddress(big.NewInt(1234567890)).Hex())
		require.NoError(t, err)

		overrider1a, err := ocr.NewConfigOverriderImpl(testLogger, address1, flags, nil)
		require.NoError(t, err)

		overrider1b, err := ocr.NewConfigOverriderImpl(testLogger, address1, flags, nil)
		require.NoError(t, err)

		overrider2, err := ocr.NewConfigOverriderImpl(testLogger, address2, flags, nil)
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

type FakeTicker struct {
	ticks chan time.Time
}

func NewFakeTicker() *FakeTicker {
	return &FakeTicker{
		ticks: make(chan time.Time),
	}
}

func (t *FakeTicker) SimulateTick() {
	t.ticks <- time.Now()
}

func (t *FakeTicker) Ticks() <-chan time.Time {
	return t.ticks
}

func (t *FakeTicker) Pause()   {}
func (t *FakeTicker) Resume()  {}
func (t *FakeTicker) Destroy() {}
