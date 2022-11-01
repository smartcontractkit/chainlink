package config

import (
	"fmt"
	"math/big"
)

// FriendlyBigInt returns a string printing the integer in both
// decimal and hexadecimal formats.
func FriendlyBigInt(n *big.Int) string {
	return fmt.Sprintf("#%[1]v (0x%[1]x)", n)
}
