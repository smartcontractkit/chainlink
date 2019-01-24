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
	topics := models.TopicFiltersForRunLog(models.RunLogTopic, jobID)

	assert.Equal(t, 2, len(topics))
	assert.Equal(
		t,
		[]common.Hash{models.RunLogTopic},
		topics[models.RequestLogTopicSignature])

	assert.Equal(
		t,
		[]common.Hash{
			common.HexToHash("0x3461316562306538646633313463623839343032346133383939316366663066"),
			common.HexToHash("0x4a1eb0e8df314cb894024a38991cff0f00000000000000000000000000000000"),
		},
		topics[1])
}

func TestRunTopic(t *testing.T) {
	assert.Equal(t, common.HexToHash("0x6d6db1f8fe19d95b1d0fa6a4bce7bb24fbf84597b35a33ff95521fac453c1529"), models.RunLogTopic)
}

func TestOracleTopic(t *testing.T) {
	assert.Equal(t, common.HexToHash("0x574a42b2507013492566a555e07135cbe40e8085bf9dd794aa2028b1b23702c2"), models.OracleLogTopic)
}
