package values

import (
	"github.com/shopspring/decimal"

	"github.com/smartcontractkit/chainlink-common/pkg/values/pb"
)

type Decimal struct {
	Underlying decimal.Decimal
}

func NewDecimal(d decimal.Decimal) *Decimal {
	return &Decimal{Underlying: d}
}

func (d *Decimal) Proto() *pb.Value {
	return pb.NewDecimalValue(d.Underlying)
}

func (d *Decimal) Unwrap() (any, error) {
	return d.Underlying, nil
}
