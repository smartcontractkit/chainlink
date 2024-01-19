package values

import (
	"github.com/shopspring/decimal"

	"github.com/smartcontractkit/chainlink-common/pkg/values/pb"
)

type Decimal struct {
	Underlying decimal.Decimal
}

func NewDecimal(d decimal.Decimal) (*Decimal, error) {
	return &Decimal{Underlying: d}, nil
}

func (d *Decimal) Proto() (*pb.Value, error) {
	return pb.NewDecimalValue(d.Underlying)
}

func (d *Decimal) Unwrap() (any, error) {
	return d.Underlying, nil
}
