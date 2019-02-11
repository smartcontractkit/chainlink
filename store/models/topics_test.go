package models_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/stretchr/testify/assert"
)

func TestTopicFiltersForRunLog(t *testing.T) {
	t.Parallel()

	jobID := "4a1eb0e8df314cb894024a38991cff0f"
	topics, err := models.TopicFiltersForRunLog([]common.Hash{models.RunLogTopic0}, jobID)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(topics))
	assert.Equal(
		t,
		[]common.Hash{models.RunLogTopic0},
		topics[models.RequestLogTopicSignature])

	assert.Equal(
		t,
		[]common.Hash{
			common.HexToHash("0x3461316562306538646633313463623839343032346133383939316366663066"),
			common.HexToHash("0x4a1eb0e8df314cb894024a38991cff0f00000000000000000000000000000000"),
		},
		topics[1])
}

func TestTopicFiltersForRunLog_Error(t *testing.T) {
	t.Parallel()

	jobID := "Q!1eb0e8df314cb894024a38991cff0f"
	topics, err := models.TopicFiltersForRunLog([]common.Hash{models.RunLogTopic0}, jobID)

	assert.Error(t, err)
	assert.Equal(t, [][]common.Hash{}, topics)
}

func TestRunTopic0(t *testing.T) {
	assert.Equal(t, "0x6d6db1f8fe19d95b1d0fa6a4bce7bb24fbf84597b35a33ff95521fac453c1529", models.RunLogTopic0.Hex())
}

func TestRunTopic20190123(t *testing.T) {
	assert.Equal(t, "0xe9cf09ba23a60c27cfb5ad84043dba257ed0ccea7f0095ff7054ec8088ce5871", models.RunLogTopic20190123.Hex())
}

func TestRunTopic20190207(t *testing.T) {
	assert.Equal(t, "0xd8d7ecc4800d25fa53ce0372f13a416d98907a7ef3d8d3bdd79cf4fe75529c65", models.RunLogTopic20190207.Hex())
}

func TestOracleFulfillmentFunctionID0(t *testing.T) {
	assert.Equal(t, "0x76005c26", models.OracleFulfillmentFunctionID0)
}

func TestOracleFulfillmentFunctionID20190123(t *testing.T) {
	assert.Equal(t, "0xeea57e70", models.OracleFulfillmentFunctionID20190123)
}
