package parseutil

import (
	"math/big"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

func ParseBigIntFromAny(val any) (*big.Int, error) {
	if val == nil {
		return nil, errors.Errorf("nil value passed")
	}

	switch v := val.(type) {
	case decimal.Decimal:
		return ParseBigIntFromString(v.String())
	case *decimal.Decimal:
		return ParseBigIntFromString(v.String())
	case *big.Int:
		return v, nil
	case string:
		return ParseBigIntFromString(v)
	case int:
		return big.NewInt(int64(v)), nil
	case int64:
		return big.NewInt(v), nil
	case float64:
		i := new(big.Int)
		big.NewFloat(v).Int(i)
		return i, nil
	default:
		return nil, errors.Errorf("unsupported big int type %T", val)
	}
}

func ParseBigIntFromString(v string) (*big.Int, error) {
	valBigInt, success := new(big.Int).SetString(v, 10)
	if !success {
		return nil, errors.Errorf("unable to convert to integer %s", v)
	}

	return valBigInt, nil
}
