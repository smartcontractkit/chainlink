package external_adapter

import (
	"context"
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

	// expected:
	// {
	//     "data": {
	//         "token": "eth",
	//         "reserves": "platform",
	//         "supplyChains": [
	//             "5009297550715157269"
	//         ],
	//         "supplyChainBlocks": [
	//             0
	//         ]
	//     }
	// }

	// Setup test context, logger, and other dependencies
	ctx := testutils.Context(t)
	lggr := logger.TestLogger(t)
	runner := mocks.NewRunner(t)
	saver := ocrcommon.NewResultRunSaver(
		runner,
		lggr,
		1000,
		100,
	)

	job := job.Job{}
	spec := pipeline.Spec{}

	ea := NewExternalAdapter(runner, job, spec, saver, lggr)

	executedRun := &pipeline.Run{}

	// "data": {
	//         "mintables": {
	//             "5009297550715157269": {
	//                 "mintable": "0",
	//                 "block": 0
	//             }
	//         },
	//         "reserveInfo": {
	//             "reserveAmount": "10332550000000000000000",
	//             "timestamp": 1749483841486
	//         },
	//         "latestRelevantBlocks": {
	//             "5009297550715157269": 22667990
	//         },
	//         "supplyDetails": {
	//             "supply": "47550052000000000000000000",
	//             "premint": "0",
	//             "chains": {
	//                 "5009297550715157269": {
	//                     "latest_block": 22667990,
	//                     "response_block": 0,
	//                     "request_block": 22667990,
	//                     "mintable": "0",
	//                     "token_supply": "44153737311060787567559446",
	//                     "token_native_mint": "0",
	//                     "token_ccip_mint": "1637953921482741588493980",
	//                     "token_ccip_burn": "5034268610421954020934534",
	//                     "token_pre_mint": "0",
	//                     "aggregate_pre_mint": false
	//                 }
	//             }
	//         }
	//     },
	//     "statusCode": 200,
	//     "result": 0,
	//     "timestamps": {
	//         "providerDataRequestedUnixMs": 1749483841817,
	//         "providerDataReceivedUnixMs": 1749483841984
	//     },
	//     "meta": {
	//         "adapterName": "SECURE_MINT",
	//         "metrics": {
	//             "feedId": "{\"token\":\"eth\",\"reserves\":\"platform\",\"supplyChains\":[\"5009297550715157269\"],\"supplyChainBlocks\":[0]}"
	//         }
	//     }
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
	// 	oracleCreator.EXPECT().Create(mock.Anything, mock.Anything, mock.MatchedBy(func(cfg cctypes.OCR3ConfigWithMeta) bool {
	// 		return cfg.Config.PluginType == uint8(cctypes.PluginTypeCCIPExec)
	// 	})).
	// 		Return(mocks.NewCCIPOracle(t), nil)
	// },

	// 	mockConnector.EXPECT().AwaitConnection(mock.Anything, mock.Anything).Return(errors.New("gateway connection failed: timeout")).Run(func(ctx context.Context, gatewayID string) {
	// 		callCount++
	// 		if callCount == len(gateways) {
	// 			cancelFunc := ctx.Value(ctxKey("cancelFunc")).(context.CancelFunc)
	// 			cancelFunc()
	// 		}
	// 	})

	var pipelineVars pipeline.Vars
	runner.EXPECT().ExecuteRun(mock.Anything, mock.Anything, mock.Anything).Return(executedRun, results, nil).Run(func(ctx context.Context, spec pipeline.Spec, vars pipeline.Vars) {
		pipelineVars = vars
	})

	blocks := por.Blocks{
		1234567890: 1234567890,
	}
	payload, err := ea.GetPayload(ctx, blocks)
	if err != nil {
		t.Fatalf("GetPayload failed: %v", err)
	}

	// validate the blocks parameter serialized to json
	blocksJSON, err := pipelineVars.Get("ea_request")
	require.NoError(t, err)
	assert.JSONEq(t,
		`{
	        "token": "eth",
	        "reserves": "platform",
	        "supplyChains": [
	            "1234567890"
	        ],
	        "supplyChainBlocks": [
	            1234567890
	        ]
	    }`,
		blocksJSON.(string),
	)

	// Validate the payload
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
