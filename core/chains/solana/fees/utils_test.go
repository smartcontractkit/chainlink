package fees

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateFee(t *testing.T) {
	inputs := []struct {
		base, max, min uint64
		count          uint
		expected       uint64
	}{
		{0, 0, 0, 100, 0},    // test max
		{0, 10, 1, 0, 1},     // test min
		{0, 10, 0, 0, 0},     // test 0 count should return base
		{0, 10, 0, 1, 1},     // test 1 count on 0 base should return 1
		{0, 10, 0, 2, 2},     // test 2 count on 0 base should return 2
		{0, 10, 0, 3, 4},     // test 3 count on 0 base should return 4
		{0, 10, 0, 4, 8},     // test 4 count on 0 base should return 8
		{1, 10, 0, 0, 1},     // test 0 count on 1 base should return 1
		{1, 10, 0, 1, 2},     // test 1 count on 1 base should return 2
		{1, 100, 0, 64, 100}, // test 64 bcount on 1 base should return max (overflow)
	}

	for i, v := range inputs {
		t.Run(fmt.Sprintf("inputs[%d]", i), func(t *testing.T) {
			assert.Equal(t, v.expected, CalculateFee(v.base, v.max, v.min, v.count))
		})
	}
}
