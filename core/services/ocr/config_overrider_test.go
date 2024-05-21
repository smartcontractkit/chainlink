package ocr_test

import (
	"math"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting/types"

	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	evmutils "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type configOverriderUni struct {
	overrider       *ocr.ConfigOverriderImpl
	contractAddress types.EIP55Address
}

type deltaCConfig struct{}

func (d deltaCConfig) DeltaCOverride() time.Duration { return time.Hour * 24 * 7 }

func (d deltaCConfig) DeltaCJitterOverride() time.Duration { return time.Hour }

func newConfigOverriderUni(t *testing.T, pollITicker utils.TickerBase, flagsContract *mocks.Flags) (uni configOverriderUni) {
	var testLogger = logger.TestLogger(t)
	contractAddress := cltest.NewEIP55Address()

	flags := &ocr.ContractFlags{FlagsInterface: flagsContract}
	var err error
	uni.overrider, err = ocr.NewConfigOverriderImpl(
		testLogger,
		deltaCConfig{},
		contractAddress,
		flags,
		pollITicker,
	)
	require.NoError(t, err)

	uni.contractAddress = contractAddress

	return uni
}

func TestIntegration_OCRConfigOverrider_EntersHibernation(t *testing.T) {
	g := gomega.NewWithT(t)

	flagsContract := mocks.NewFlags(t)

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

	servicetest.Run(t, uni.overrider)

	// not hibernating initially
	require.Nil(t, uni.overrider.ConfigOverride())

	expectedOverride := &ocrtypes.ConfigOverride{AlphaPPB: math.MaxUint64, DeltaC: uni.overrider.DeltaCFromAddress}

	// timeout needs to be longer than the poll interval of 3 seconds
	g.Eventually(func() *ocrtypes.ConfigOverride { return uni.overrider.ConfigOverride() }, 10*time.Second, 450*time.Millisecond).Should(gomega.Equal(expectedOverride))
}

func Test_OCRConfigOverrider(t *testing.T) {
	t.Parallel()

	t.Run("Before first tick returns nil override, later does return a specific override when hibernating", func(t *testing.T) {
		flagsContract := mocks.NewFlags(t)

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

		servicetest.Run(t, uni.overrider)

		// not hibernating initially
		require.Nil(t, uni.overrider.ConfigOverride())

		// update state by getting flags
		require.NoError(t, uni.overrider.ExportedUpdateFlagsStatus())

		expectedOverride := &ocrtypes.ConfigOverride{AlphaPPB: math.MaxUint64, DeltaC: uni.overrider.DeltaCFromAddress}
		require.Equal(t, expectedOverride, uni.overrider.ConfigOverride())
	})

	t.Run("Before first tick is hibernating, later exists hibernation", func(t *testing.T) {
		flagsContract := mocks.NewFlags(t)

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

		servicetest.Run(t, uni.overrider)

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
			deltaCConfig{},
			contractAddress,
			flags,
			nil,
		)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Flags contract instance is missing, the contract does not exist")
	})

	t.Run("DeltaC should be stable per address", func(t *testing.T) {
		var testLogger = logger.TestLogger(t)
		flagsContract := mocks.NewFlags(t)
		flags := &ocr.ContractFlags{FlagsInterface: flagsContract}

		address1, err := types.NewEIP55Address(common.BigToAddress(big.NewInt(10000)).Hex())
		require.NoError(t, err)

		address2, err := types.NewEIP55Address(common.BigToAddress(big.NewInt(1234567890)).Hex())
		require.NoError(t, err)

		overrider1a, err := ocr.NewConfigOverriderImpl(testLogger, deltaCConfig{}, address1, flags, nil)
		require.NoError(t, err)

		overrider1b, err := ocr.NewConfigOverriderImpl(testLogger, deltaCConfig{}, address1, flags, nil)
		require.NoError(t, err)

		overrider2, err := ocr.NewConfigOverriderImpl(testLogger, deltaCConfig{}, address2, flags, nil)
		require.NoError(t, err)

		assert.Equal(t, cltest.MustParseDuration(t, "168h46m40s"), overrider1a.DeltaCFromAddress)
		assert.Equal(t, cltest.MustParseDuration(t, "168h46m40s"), overrider1b.DeltaCFromAddress)
		assert.Equal(t, cltest.MustParseDuration(t, "168h31m30s"), overrider2.DeltaCFromAddress)
	})
}

func checkFlagsAddress(t *testing.T, contractAddress types.EIP55Address) func(args mock.Arguments) {
	return func(args mock.Arguments) {
		require.Equal(t, []common.Address{
			evmutils.ZeroAddress,
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
