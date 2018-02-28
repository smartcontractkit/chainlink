package models_test

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

func TestModels_HexToFunctionSelector(t *testing.T) {
	fid := models.HexToFunctionSelector("0xb3f98adc")
	assert.Equal(t, "0xb3f98adc", fid.String())
}

func TestModels_HexToFunctionSelectorOverflow(t *testing.T) {
	fid := models.HexToFunctionSelector("0xb3f98adc123456")
	assert.Equal(t, "0xb3f98adc", fid.String())
}

func TestModels_FunctionSelectorUnmarshalJSON(t *testing.T) {
	bytes := []byte(`"0xb3f98adc"`)
	var fid models.FunctionSelector
	err := json.Unmarshal(bytes, &fid)
	assert.Nil(t, err)
	assert.Equal(t, "0xb3f98adc", fid.String())
}

func TestModels_FunctionSelectorUnmarshalJSONError(t *testing.T) {
	bytes := []byte(`"0xb3f98adc123456"`)
	var fid models.FunctionSelector
	err := json.Unmarshal(bytes, &fid)
	assert.NotNil(t, err)
}

func TestModels_Header_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	var header models.BlockHeader

	data := cltest.LoadJSON("../../internal/fixtures/eth/subscription_new_heads.json")
	value := gjson.Get(string(data), "params.result")
	assert.Nil(t, json.Unmarshal([]byte(value.String()), &header))

	assert.Equal(t, cltest.BigHexInt(1263817), header.Number)
}

func TestModels_IndexableBlockNumber(t *testing.T) {
	tests := []struct {
		input      *big.Int
		want       string
		wantDigits int
	}{
		{big.NewInt(0), "0x0", 1},
		{big.NewInt(0xf), "0xf", 1},
		{big.NewInt(0x10), "0x10", 2},
	}
	for _, test := range tests {
		t.Run(test.want, func(t *testing.T) {
			t.Parallel()
			num := models.NewIndexableBlockNumber(test.input)
			assert.Equal(t, test.want, num.String())
			assert.Equal(t, test.wantDigits, num.Digits)
		})
	}
}
