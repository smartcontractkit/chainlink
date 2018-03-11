package models_test

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
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
	tests := []struct {
		name       string
		path       string
		wantNumber hexutil.Big
		wantHash   string
	}{
		{
			"parity",
			"../../internal/fixtures/eth/subscription_new_heads_parity.json",
			cltest.BigHexInt(1263817),
			"0xf8e4691ceab8052d1cb478c6c5e0d9b122e747ad838023633f63bd5e81ec5114",
		},
		{
			"geth",
			"../../internal/fixtures/eth/subscription_new_heads_geth.json",
			cltest.BigHexInt(1263817),
			"0xf8e4691ceab8052d1cb478c6c5e0d9b122e747ad838023633f63bd5e81ec5fff",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var header models.BlockHeader

			data := cltest.LoadJSON(test.path)
			value := gjson.Get(string(data), "params.result")
			assert.Nil(t, json.Unmarshal([]byte(value.String()), &header))

			assert.Equal(t, test.wantNumber, header.Number)
			assert.Equal(t, test.wantHash, header.Hash().String())
		})
	}
}

func TestModels_IndexableBlockNumber_New(t *testing.T) {
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

func TestModels_IndexableBlockNumber_GreaterThan(t *testing.T) {
	tests := []struct {
		name    string
		left    *models.IndexableBlockNumber
		right   *models.IndexableBlockNumber
		greater bool
	}{
		{"nil nil", nil, nil, false},
		{"present nil", cltest.IndexableBlockNumber(1), nil, false},
		{"nil present", cltest.IndexableBlockNumber(2), cltest.IndexableBlockNumber(1), false},
		{"less", cltest.IndexableBlockNumber(1), cltest.IndexableBlockNumber(2), false},
		{"equal", cltest.IndexableBlockNumber(2), cltest.IndexableBlockNumber(2), false},
		{"greater", cltest.IndexableBlockNumber(2), cltest.IndexableBlockNumber(1), true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.greater, test.left.GreaterThan(test.right))
		})
	}
}

func TestModels_IndexableBlockNumber_NextInt(t *testing.T) {
	tests := []struct {
		name string
		bn   *models.IndexableBlockNumber
		want *big.Int
	}{
		{"nil", nil, big.NewInt(0)},
		{"one", cltest.IndexableBlockNumber(1), big.NewInt(2)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.want, test.bn.NextInt())
		})
	}
}
