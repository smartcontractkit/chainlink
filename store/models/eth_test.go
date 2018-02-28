package models_test

import (
	"encoding/json"
	"testing"

	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

func TestHexToFunctionSelector(t *testing.T) {
	fid := models.HexToFunctionSelector("0xb3f98adc")
	assert.Equal(t, "0xb3f98adc", fid.String())
}

func TestHexToFunctionSelectorOverflow(t *testing.T) {
	fid := models.HexToFunctionSelector("0xb3f98adc123456")
	assert.Equal(t, "0xb3f98adc", fid.String())
}

func TestFunctionSelectorUnmarshalJSON(t *testing.T) {
	bytes := []byte(`"0xb3f98adc"`)
	var fid models.FunctionSelector
	err := json.Unmarshal(bytes, &fid)
	assert.Nil(t, err)
	assert.Equal(t, "0xb3f98adc", fid.String())
}

func TestFunctionSelectorUnmarshalJSONError(t *testing.T) {
	bytes := []byte(`"0xb3f98adc123456"`)
	var fid models.FunctionSelector
	err := json.Unmarshal(bytes, &fid)
	assert.NotNil(t, err)
}

func TestHeader_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	var header models.BlockHeader

	data := cltest.LoadJSON("../../internal/fixtures/eth/subscription_new_heads.json")
	value := gjson.Get(string(data), "params.result")
	assert.Nil(t, json.Unmarshal([]byte(value.String()), &header))

	assert.Equal(t, cltest.BigHexInt(1263817), header.Number)
}
