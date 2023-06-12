package config

import (
	"fmt"
	"math/big"

	"golang.org/x/exp/constraints"
)

// FriendlyNumber returns a string printing the integer or big.Int in both
// decimal and hexadecimal formats.
func FriendlyNumber[N constraints.Integer | *big.Int](n N) string {
	return fmt.Sprintf("#%[1]v (0x%[1]x)", n)
}
