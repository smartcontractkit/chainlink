package ea

import (
	"context"
	"encoding/json"
	"math/big"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/types/core/securemint"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	sm_config "github.com/smartcontractkit/chainlink/v2/core/services/ocr3/securemint/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline/mocks"
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

	config := &sm_config.SecureMintConfig{
		Token:          "eth",
		Reserves:       "platform",
		ChainSelectors: []string{"5009297550715157269"},
	}
	job := job.Job{}
	spec := pipeline.Spec{}
	executedRun := &pipeline.Run{}

	ea, err := NewExternalAdapter(config, runner, job, spec, saver, lggr)
	require.NoError(t, err)

	results := pipeline.TaskRunResults{
		{
			Task: &pipeline.AnyTask{},
			Result: pipeline.Result{
				Value: map[string]any{ // outer `data` field is already stripped off in the parse step of the pipeline
					"mintables": map[string]any{
						"5009297550715157269": map[string]any{
							"mintable": "10",
							"block":    8,
						},
					},
					"reserveInfo": map[string]any{
						"reserveAmount": "10332550000000000000000",
						"timestamp":     1749483841486,
					},
					"latestBlocks": map[string]any{
						"5009297550715157269": 23,
					},
				},
				Error: nil,
			},
		},
	}

	// Capture the 'ea_request' parameter from the pipeline run
	var eaRequest any
	runner.EXPECT().ExecuteRun(mock.Anything, mock.Anything, mock.Anything).Return(executedRun, results, nil).Run(func(_ context.Context, _ pipeline.Spec, vars pipeline.Vars) {
		var err error
		eaRequest, err = vars.Get("ea_request")
		require.NoError(t, err)
	})

	payload, err := ea.GetPayload(ctx, securemint.Blocks{1234567890: 1234567890, 5009297550715157269: 10})
	require.NoError(t, err, "GetPayload should not return an error")

	// Validate the 'ea_request' parameter serialized to json
	eaRequestJSON, err := json.Marshal(eaRequest)
	require.NoError(t, err, "Failed to marshal ea_request to JSON")
	assert.JSONEq(t,
		`{
	        "reserves": "platform",
	        "supplyChains": [
	            "5009297550715157269"
	        ],
	        "supplyChainBlocks": [
	            10
	        ],
			"token": "eth"
	    }`,
		string(eaRequestJSON),
	)

	// Validate the resulting payload
	amount, ok := big.NewInt(10).SetString("10332550000000000000000", 10)
	require.True(t, ok, "Failed to parse reserve amount from string")
	expectedPayload := securemint.ExternalAdapterPayload{
		Mintables: securemint.Mintables{
			5009297550715157269: {
				Block:    8,
				Mintable: big.NewInt(10),
			},
		},
		ReserveInfo: securemint.ReserveInfo{
			ReserveAmount: amount,
			Timestamp:     time.UnixMilli(1749483841486),
		},
		LatestBlocks: securemint.Blocks{
			5009297550715157269: 23,
		},
	}
	assert.Equal(t, expectedPayload, payload)
}
