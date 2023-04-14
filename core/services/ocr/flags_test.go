package ocr_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/fluxmonitorv2"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func TestFlags_IsLowered(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		getFlagsResult []bool
		expected       bool
	}{
		{"both lowered", []bool{false, false}, true},
		{"global lowered", []bool{false, true}, true},
		{"contract lowered", []bool{true, false}, true},
		{"both raised", []bool{true, true}, false},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			var (
				flagsContract = mocks.NewFlags(t)
				address       = testutils.NewAddress()
			)

			flags := fluxmonitorv2.ContractFlags{FlagsInterface: flagsContract}

			flagsContract.On("GetFlags", mock.Anything, mock.Anything).
				Run(func(args mock.Arguments) {
					require.Equal(t, []common.Address{
						utils.ZeroAddress,
						address,
					}, args.Get(1).([]common.Address))
				}).
				Return(tc.getFlagsResult, nil)

			result, err := flags.IsLowered(address)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}
