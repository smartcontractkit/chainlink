package contracts

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/shopspring/decimal"
)

func newBigIntFromString(arg string) (*big.Int, error) {
	if arg == "0x" {
		// Oddly a legal value for zero
		arg = "0x0"
	}
	ret, ok := new(big.Int).SetString(arg, 0)
	if !ok {
		return nil, fmt.Errorf("cannot convert '%s' to big int", arg)
	}
	return ret, nil
}

var dec10 = decimal.NewFromInt(10)

func newDecimalFromString(arg string) (decimal.Decimal, error) {
	if strings.HasPrefix(arg, "0x") {
		// decimal package does not parse Hex values
		value, err := newBigIntFromString(arg)
		if err != nil {
			return decimal.Zero, fmt.Errorf("cannot convert '%s' to decimal", arg)
		}
		return decimal.NewFromString(value.Text(10))
	}
	return decimal.NewFromString(arg)
}
