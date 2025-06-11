package ea

import (
	"context"
	"encoding/json"
	"math/big"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline/mocks"
	"github.com/smartcontractkit/por_mock_ocr3plugin/por"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_GetPayload(t *testing.T) {

	// Setup test context, logger, and other dependencies
	ctx := testutils.Context(t)
	lggr := logger.NullLogger
	runner := mocks.NewRunner(t)
	saver := ocrcommon.NewResultRunSaver(
		runner,
		lggr,
		1000,
		100,
	)

	job := job.Job{}
	spec := pipeline.Spec{}
	executedRun := &pipeline.Run{}

	ea := NewExternalAdapter(runner, job, spec, saver, lggr)

	results := pipeline.TaskRunResults{
		{
			Task: &pipeline.AnyTask{},
			Result: pipeline.Result{
				Value: map[string]any{ // outer `data` field is already stripped off in the parse step of the pipeline
					"mintables": map[string]any{
						"1234567890": map[string]any{
							"mintable": "10",
							"block":    8,
						},
					},
					"reserveInfo": map[string]any{
						"reserveAmount": "10332550000000000000000",
						"timestamp":     1749483841486,
					},
					"latestRelevantBlocks": map[string]any{
						"1234567890": 23,
					},
				},
				Error: nil,
			},
		},
	}

	// capture the 'ea_request' parameter from the pipeline run
	var eaRequest any
	runner.EXPECT().ExecuteRun(mock.Anything, mock.Anything, mock.Anything).Return(executedRun, results, nil).Run(func(_ context.Context, _ pipeline.Spec, vars pipeline.Vars) {
		var err error
		eaRequest, err = vars.Get("ea_request")
		require.NoError(t, err)
	})

	payload, err := ea.GetPayload(ctx, por.Blocks{1234567890: 1234567890})
	require.NoError(t, err, "GetPayload should not return an error")

	// validate the 'ea_request' parameter serialized to json
	eaRequestJSON, err := json.Marshal(eaRequest)
	require.NoError(t, err, "Failed to marshal ea_request to JSON")
	assert.JSONEq(t,
		`{
	        "reserves": "platform",
	        "supplyChains": [
	            "1234567890"
	        ],
	        "supplyChainBlocks": [
	            1234567890
	        ],
			"token": "eth"
	    }`,
		string(eaRequestJSON),
	)

	// Validate the resulting payload
	amount, ok := big.NewInt(10).SetString("10332550000000000000000", 10)
	require.True(t, ok, "Failed to parse reserve amount from string")
	expectedPayload := por.ExternalAdapterPayload{
		Mintables: por.Mintables{
			1234567890: {
				Block:    8,
				Mintable: big.NewInt(10),
			},
		},
		ReserveInfo: por.ReserveInfo{
			ReserveAmount: amount,
			Timestamp:     time.UnixMilli(1749483841486),
		},
		LatestRelevantBlocks: por.Blocks{
			1234567890: 23,
		},
	}
	assert.Equal(t, expectedPayload, payload)
}
