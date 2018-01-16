package models_test

import (
	"encoding/json"
	"testing"

	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/stretchr/testify/assert"
)

func TestHexToFunctionID(t *testing.T) {
	fid := models.HexToFunctionID("0xb3f98adc")
	assert.Equal(t, "0xb3f98adc", fid.String())
}

func TestHexToFunctionIDOverflow(t *testing.T) {
	fid := models.HexToFunctionID("0xb3f98adc123456")
	assert.Equal(t, "0xb3f98adc", fid.String())
}

func TestFunctionIDUnmarshalJSON(t *testing.T) {
	bytes := []byte(`"0xb3f98adc"`)
	var fid models.FunctionID
	err := json.Unmarshal(bytes, &fid)
	assert.Nil(t, err)
	assert.Equal(t, "0xb3f98adc", fid.String())
}

func TestFunctionIDUnmarshalJSONError(t *testing.T) {
	bytes := []byte(`"0xb3f98adc123456"`)
	var fid models.FunctionID
	err := json.Unmarshal(bytes, &fid)
	assert.NotNil(t, err)
}
