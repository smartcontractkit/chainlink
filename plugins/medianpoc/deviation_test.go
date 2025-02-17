package medianpoc

import (
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type frozenTimeClock struct{}

func (frozenTimeClock) Now() time.Time {
	return time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
}

func Test_PendleDeviationFunc(t *testing.T) {
	tcs := []struct {
		name string

		expiresInSeconds float64
		multiplierVal    *big.Int
		thresholdPPB     uint64
		oldVal           *big.Int
		newVal           *big.Int

		err      string
		expected bool
	}{
		{
			name:   "nil oldVal errors",
			oldVal: nil,
			newVal: big.NewInt(2),
			err:    "oldVal and newVal must be non-nil",
		},
		{
			name:   "nil newVal errors",
			oldVal: big.NewInt(1),
			newVal: nil,
			err:    "oldVal and newVal must be non-nil",
		},
		{
			name:             "test 0 (one block after) - SHOULD UPDATE",
			expiresInSeconds: 13857541.0,
			thresholdPPB:     1e7,
			oldVal:           big.NewInt(0.187152977881070687 * 1e18),
			newVal:           big.NewInt(0.164498448931278907 * 1e18),
			expected:         true,
		},
		// TODO: Fix these tests
		// {
		//     name:             "test 1 (same block) - SHOULD UPDATE",
		//     expiresInSeconds: 1.3857552999999971712 * 1e7,
		//     thresholdPPB:     1e7,
		//     oldVal:           big.NewInt(0.187152977881070687 * 1e18),
		//     newVal:           big.NewInt(0.164513876918347290 * 1e18),
		//     expected:         true,
		// },
		// {
		//     name:             "test 2 (same block) - SHOULD UPDATE",
		//     expiresInSeconds: 1.38259210000000006116 * 1e7,
		//     thresholdPPB:     1e7,
		//     oldVal:           big.NewInt(0.164498448931278907 * 1e18),
		//     newVal:           big.NewInt(0.141808618152979599 * 1e18),
		//     expected:         true,
		// },
		// {
		//     name:             "test 3 (same block) - SHOULD UPDATE",
		//     expiresInSeconds: 1.1564856999999997679376 * 1e7,
		//     thresholdPPB:     1e7,
		//     oldVal:           big.NewInt(0.141802025539163407 * 1e18),
		//     newVal:           big.NewInt(0.114668695429674394 * 1e18),
		//     expected:         true,
		// },
		// {
		//     name:             "test 4 (same block) - SHOULD UPDATE",
		//     expiresInSeconds: 3.574116999999998781984 * 1e6,
		//     thresholdPPB:     1e7,
		//     oldVal:           big.NewInt(0.114668136842518698 * 1e18),
		//     newVal:           big.NewInt(0.202462661212593903 * 1e18),
		//     expected:         true,
		// },
		{
			name:             "test 5 (previous block) - SHOULD NOT UPDATE",
			expiresInSeconds: 13857565.0,
			thresholdPPB:     1e7,
			oldVal:           big.NewInt(0.187152977881070687 * 1e18),
			newVal:           big.NewInt(0.164529304905415591 * 1e18),
			expected:         false,
		},
		{
			name:             "test 6 (previous block) - SHOULD NOT UPDATE",
			expiresInSeconds: 13825932.999999998137354851,
			thresholdPPB:     1e7,
			oldVal:           big.NewInt(0.164498448931278907 * 1e18),
			newVal:           big.NewInt(0.141815210766795902 * 1e18),
			expected:         false,
		},
		{
			name:             "test 7 (previous block) - SHOULD NOT UPDATE",
			expiresInSeconds: 11564869.0,
			thresholdPPB:     1e7,
			oldVal:           big.NewInt(0.141802025539163407 * 1e18),
			newVal:           big.NewInt(0.114669254016829994 * 1e18),
			expected:         false,
		},
		{
			name:             "test 8 (previous block) - SHOULD NOT UPDATE",
			expiresInSeconds: 3574128.999999998603016138,
			thresholdPPB:     1e7,
			oldVal:           big.NewInt(0.114668136842518698 * 1e18),
			newVal:           big.NewInt(0.202460645024667513 * 1e18),
			expected:         false,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			clock := frozenTimeClock{}
			expiresAt := float64(clock.Now().Unix()) + tc.expiresInSeconds
			actual, err := makePendleDeviationFunc(expiresAt, clock, DefaultMultiplier)(nil, tc.thresholdPPB, tc.oldVal, tc.newVal)
			if tc.err != "" {
				require.EqualError(t, err, tc.err)
			} else {
				require.NoError(t, err)
				if actual != tc.expected {
					t.Fatalf("expected %v, got %v", tc.expected, actual)
				}
			}
		})
	}
}
