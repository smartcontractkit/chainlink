package models_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTopicFiltersForRunLog(t *testing.T) {
	t.Parallel()

	jobID, err := models.NewIDFromString("4a1eb0e8df314cb894024a38991cff0f")
	require.NoError(t, err)

	topics := models.TopicFiltersForRunLog([]common.Hash{models.RunLogTopic0original}, jobID)
	assert.Equal(t, 2, len(topics))
	assert.Equal(
		t,
		[]common.Hash{models.RunLogTopic0original},
		topics[models.RequestLogTopicSignature])

	assert.Equal(
		t,
		[]common.Hash{
			common.HexToHash("0x4a1eb0e8df314cb894024a38991cff0f00000000000000000000000000000000"),
			common.HexToHash("0x3461316562306538646633313463623839343032346133383939316366663066"),
		},
		topics[1])
}

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
