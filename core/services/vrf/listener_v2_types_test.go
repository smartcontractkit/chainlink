package vrf

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/chains/evm/log/mocks"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
)

func Test_BatchFulfillments_AddRun(t *testing.T) {
	batchLimit := uint64(2500)
	bfs := newBatchFulfillments(batchLimit)
	for i := 0; i < 4; i++ {
		bfs.addRun(vrfPipelineResult{
			gasLimit: 500,
			req: pendingRequest{
				req: &vrf_coordinator_v2.VRFCoordinatorV2RandomWordsRequested{
					RequestId: big.NewInt(1),
				},
				lb: &mocks.Broadcast{},
			},
			run: pipeline.NewRun(pipeline.Spec{}, pipeline.Vars{}),
		})
		require.Len(t, bfs.fulfillments, 1)
	}

	require.Equal(t, uint64(2000), bfs.fulfillments[0].totalGasLimit)

	// This addition should create and add a new batch
	bfs.addRun(vrfPipelineResult{
		gasLimit: 500,
		req: pendingRequest{
			req: &vrf_coordinator_v2.VRFCoordinatorV2RandomWordsRequested{
				RequestId: big.NewInt(1),
			},
			lb: &mocks.Broadcast{},
		},
		run: pipeline.NewRun(pipeline.Spec{}, pipeline.Vars{}),
	})
	require.Len(t, bfs.fulfillments, 2)
}
