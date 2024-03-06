package values

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_DecimalUnwrapTo(t *testing.T) {
	dv := decimal.NewFromFloat(1.00)
	tr := NewDecimal(dv)

	var dec decimal.Decimal
	err := tr.UnwrapTo(&dec)
	require.NoError(t, err)

	assert.Equal(t, dv, dec)

	var s string
	err = tr.UnwrapTo(&s)
	require.Error(t, err)

	decn := (*decimal.Decimal)(nil)
	err = tr.UnwrapTo(decn)
	assert.ErrorContains(t, err, "unwrap to nil pointer")
}
