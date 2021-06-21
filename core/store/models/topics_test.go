package models_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/core/store/models"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJobSpecIDTopics(t *testing.T) {
	t.Parallel()

	jobID, err := models.NewJobIDFromString("4a1eb0e8df314cb894024a38991cff0f")
	require.NoError(t, err)

	topics := models.JobSpecIDTopics(jobID)
	assert.Equal(t, 2, len(topics))

	assert.Equal(
		t,
		[]common.Hash{
			common.HexToHash("0x4a1eb0e8df314cb894024a38991cff0f00000000000000000000000000000000"),
			common.HexToHash("0x3461316562306538646633313463623839343032346133383939316366663066"),
		},
		topics)
}

func TestRunLogTopic20190207withoutIndexes(t *testing.T) {
	assert.Equal(t, "0xd8d7ecc4800d25fa53ce0372f13a416d98907a7ef3d8d3bdd79cf4fe75529c65", models.RunLogTopic20190207withoutIndexes.Hex())
}

func TestOracleFulfillmentFunctionID20190128withoutCast(t *testing.T) {
	assert.Equal(t, "0x4ab0d190", models.OracleFulfillmentFunctionID20190128withoutCast)
}
