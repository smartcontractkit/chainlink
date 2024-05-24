package v1

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/theodesp/go-heaps/pairing"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/solidity_vrf_coordinator_interface"
)

func TestConfirmedLogExtraction(t *testing.T) {
	lsn := Listener{}
	lsn.Reqs = []request{
		{
			confirmedAtBlock: 2,
			req: &solidity_vrf_coordinator_interface.VRFCoordinatorRandomnessRequest{
				RequestID: utils.PadByteToHash(0x02),
			},
		},
		{
			confirmedAtBlock: 1,
			req: &solidity_vrf_coordinator_interface.VRFCoordinatorRandomnessRequest{
				RequestID: utils.PadByteToHash(0x01),
			},
		},
		{
			confirmedAtBlock: 3,
			req: &solidity_vrf_coordinator_interface.VRFCoordinatorRandomnessRequest{
				RequestID: utils.PadByteToHash(0x03),
			},
		},
	}
	// None are confirmed
	lsn.LatestHead = 0
	logs := lsn.extractConfirmedLogs()
	assert.Equal(t, 0, len(logs))     // None ready
	assert.Equal(t, 3, len(lsn.Reqs)) // All pending
	lsn.LatestHead = 2
	logs = lsn.extractConfirmedLogs()
	assert.Equal(t, 2, len(logs))     // 1 and 2 should be confirmed
	assert.Equal(t, 1, len(lsn.Reqs)) // 3 is still pending
	assert.Equal(t, uint64(3), lsn.Reqs[0].confirmedAtBlock)
	// Another block way in the future should clear it
	lsn.LatestHead = 10
	logs = lsn.extractConfirmedLogs()
	assert.Equal(t, 1, len(logs))     // remaining log
	assert.Equal(t, 0, len(lsn.Reqs)) // all processed
}

func TestResponsePruning(t *testing.T) {
	lsn := Listener{}
	lsn.LatestHead = 10000
	lsn.ResponseCount = map[[32]byte]uint64{
		utils.PadByteToHash(0x00): 1,
		utils.PadByteToHash(0x01): 1,
	}
	lsn.BlockNumberToReqID = pairing.New()
	lsn.BlockNumberToReqID.Insert(fulfilledReq{
		blockNumber: 1,
		reqID:       utils.PadByteToHash(0x00),
	})
	lsn.BlockNumberToReqID.Insert(fulfilledReq{
		blockNumber: 2,
		reqID:       utils.PadByteToHash(0x01),
	})
	lsn.pruneConfirmedRequestCounts()
	assert.Equal(t, 2, len(lsn.ResponseCount))
	lsn.LatestHead = 10001
	lsn.pruneConfirmedRequestCounts()
	assert.Equal(t, 1, len(lsn.ResponseCount))
	lsn.LatestHead = 10002
	lsn.pruneConfirmedRequestCounts()
	assert.Equal(t, 0, len(lsn.ResponseCount))
}
