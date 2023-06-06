package cosmos

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func Test_sdkDecFromDecimal(t *testing.T) {
	tests := []string{
		"0.0",
		"0.1",
		"1.0",
		"0.000000000000000001",
	}
	for _, tt := range tests {
		t.Run(tt, func(t *testing.T) {
			val := decimal.RequireFromString(tt)
			exp := sdk.MustNewDecFromStr(tt)
			assert.Equal(t, exp, sdkDecFromDecimal(&val))
		})
	}
}
