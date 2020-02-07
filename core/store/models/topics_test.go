package models_test

import (
	"testing"

	"chainlink/core/store/models"

	"github.com/stretchr/testify/assert"
)

func TestRunLogTopic0original(t *testing.T) {
	assert.Equal(t, "0x6d6db1f8fe19d95b1d0fa6a4bce7bb24fbf84597b35a33ff95521fac453c1529", models.RunLogTopic0original.Hex())
}

func TestRunLogTopic20190123withFulfillmentParams(t *testing.T) {
	assert.Equal(t, "0xe9cf09ba23a60c27cfb5ad84043dba257ed0ccea7f0095ff7054ec8088ce5871", models.RunLogTopic20190123withFullfillmentParams.Hex())
}

func TestRunLogTopic20190207withoutIndexes(t *testing.T) {
	assert.Equal(t, "0xd8d7ecc4800d25fa53ce0372f13a416d98907a7ef3d8d3bdd79cf4fe75529c65", models.RunLogTopic20190207withoutIndexes.Hex())
}

func TestOracleFulfillmentFunctionID0original(t *testing.T) {
	assert.Equal(t, "0x76005c26", models.OracleFullfillmentFunctionID0original)
}

func TestOracleFulfillmentFunctionID20190123withFulfillmentParams(t *testing.T) {
	assert.Equal(t, "0xeea57e70", models.OracleFulfillmentFunctionID20190123withFulfillmentParams)
}

func TestOracleFulfillmentFunctionID20190128withoutCast(t *testing.T) {
	assert.Equal(t, "0x4ab0d190", models.OracleFulfillmentFunctionID20190128withoutCast)
}
