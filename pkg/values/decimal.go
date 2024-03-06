package values

import (
	"fmt"

	"github.com/shopspring/decimal"

	"github.com/smartcontractkit/chainlink-common/pkg/values/pb"
)

type Decimal struct {
	Underlying decimal.Decimal
}

func NewDecimal(d decimal.Decimal) *Decimal {
	return &Decimal{Underlying: d}
}

func (d *Decimal) proto() *pb.Value {
	return pb.NewDecimalValue(d.Underlying)
}

func (d *Decimal) Unwrap() (any, error) {
	return d.Underlying, nil
}

func (d *Decimal) UnwrapTo(to any) error {
	dv, ok := to.(*decimal.Decimal)
	if !ok {
		return fmt.Errorf("cannot unwrap to non-pointer type %T", to)
	}

	if dv == nil {
		return fmt.Errorf("cannot unwrap to nil pointer: %+v", to)
	}

	*dv = d.Underlying
	return nil
}
