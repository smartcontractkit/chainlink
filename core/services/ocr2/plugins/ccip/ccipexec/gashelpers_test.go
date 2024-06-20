package ccipexec

import (
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOverheadGas(t *testing.T) {
	// Only Data and TokenAmounts are used from the messages
	// And only the length is used so the contents doesn't matter.
	tests := []struct {
		dataLength     int
		numberOfTokens int
		want           uint64
	}{
		{
			dataLength:     0,
			numberOfTokens: 0,
			want:           119920,
		},
		{
			dataLength:     len([]byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0}),
			numberOfTokens: 1,
			want:           475948,
		},
	}

	for _, tc := range tests {
		got := overheadGas(tc.dataLength, tc.numberOfTokens)
		if !reflect.DeepEqual(tc.want, got) {
			t.Fatalf("expected: %v, got: %v", tc.want, got)
		}
	}
}

func TestMaxGasOverHeadGas(t *testing.T) {
	// Only Data and TokenAmounts are used from the messages
	// And only the length is used so the contents doesn't matter.
	tests := []struct {
		numMsgs        int
		dataLength     int
		numberOfTokens int
		want           uint64
	}{
		{
			numMsgs:        6,
			dataLength:     0,
			numberOfTokens: 0,
			want:           122992,
		},
		{
			numMsgs:        3,
			dataLength:     len([]byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0}),
			numberOfTokens: 1,
			want:           478508,
		},
	}

	for _, tc := range tests {
		got := maxGasOverHeadGas(tc.numMsgs, tc.dataLength, tc.numberOfTokens)
		if !reflect.DeepEqual(tc.want, got) {
			t.Fatalf("expected: %v, got: %v", tc.want, got)
		}
	}
}

func TestWaitBoostedFee(t *testing.T) {
	tests := []struct {
		name                     string
		sendTimeDiff             time.Duration
		fee                      *big.Int
		diff                     *big.Int
		relativeBoostPerWaitHour float64
	}{
		{
			"wait 10s",
			time.Second * 10,
			big.NewInt(6e18), // Fee:   6    LINK

			big.NewInt(1166666666665984), // Boost: 0.01 LINK
			0.07,
		},
		{
			"wait 5m",
			time.Minute * 5,
			big.NewInt(6e18),  // Fee:   6    LINK
			big.NewInt(35e15), // Boost: 0.35 LINK
			0.07,
		},
		{
			"wait 7m",
			time.Minute * 7,
			big.NewInt(6e18),  // Fee:   6    LINK
			big.NewInt(49e15), // Boost: 0.49 LINK
			0.07,
		},
		{
			"wait 12m",
			time.Minute * 12,
			big.NewInt(6e18),  // Fee:   6    LINK
			big.NewInt(84e15), // Boost: 0.84 LINK
			0.07,
		},
		{
			"wait 25m",
			time.Minute * 25,
			big.NewInt(6e18),               // Fee:   6 LINK
			big.NewInt(174999999999998976), // Boost: 1.75 LINK
			0.07,
		},
		{
			"wait 1h",
			time.Hour * 1,
			big.NewInt(6e18),   // Fee:   6 LINK
			big.NewInt(420e15), // Boost: 4.2 LINK
			0.07,
		},
		{
			"wait 5h",
			time.Hour * 5,
			big.NewInt(6e18),                // Fee:   6 LINK
			big.NewInt(2100000000000001024), // Boost: 21LINK
			0.07,
		},
		{
			"wait 24h",
			time.Hour * 24,
			big.NewInt(6e18), // Fee:   6 LINK
			big.NewInt(0).Mul(big.NewInt(10), big.NewInt(1008e15)), // Boost: 100LINK
			0.07,
		},
		{
			"high boost wait 10s",
			time.Second * 10,
			big.NewInt(5e18),
			big.NewInt(9722222222222336), // 1e16
			0.7,
		},
		{
			"high boost wait 5m",
			time.Minute * 5,
			big.NewInt(5e18),
			big.NewInt(291666666666667008), // 1e18
			0.7,
		},
		{
			"high boost wait 25m",
			time.Minute * 25,
			big.NewInt(5e18),
			big.NewInt(1458333333333334016), // 1e19
			0.7,
		},
		{
			"high boost wait 5h",
			time.Hour * 5,
			big.NewInt(5e18),
			big.NewInt(0).Mul(big.NewInt(10), big.NewInt(175e16)), // 1e20
			0.7,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			boosted := waitBoostedFee(tc.sendTimeDiff, tc.fee, tc.relativeBoostPerWaitHour)
			diff := big.NewInt(0).Sub(boosted, tc.fee)
			assert.Equal(t, diff, tc.diff)
			// we check that the actual diff is approximately equals to expected diff,
			// as we might get slightly different results locally vs. CI therefore normal Equal() would be unstable
			//diffUpperLimit := big.NewInt(0).Add(tc.diff, big.NewInt(1e9))
			//diffLowerLimit := big.NewInt(0).Add(tc.diff, big.NewInt(-1e9))
			//require.Equalf(t, -1, diff.Cmp(diffUpperLimit), "actual diff (%s) is larger than expected (%s)", diff.String(), diffUpperLimit.String())
			//require.Equal(t, 1, diff.Cmp(diffLowerLimit), "actual diff (%s) is smaller than expected (%s)", diff.String(), diffLowerLimit.String())
		})
	}
}
