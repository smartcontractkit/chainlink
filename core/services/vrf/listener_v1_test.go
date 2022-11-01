package vrf

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/theodesp/go-heaps/pairing"

	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/solidity_vrf_coordinator_interface"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func TestConfirmedLogExtraction(t *testing.T) {
	lsn := listenerV1{}
	lsn.reqs = []request{
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
	lsn.latestHead = 0
	logs := lsn.extractConfirmedLogs()
	assert.Equal(t, 0, len(logs))     // None ready
	assert.Equal(t, 3, len(lsn.reqs)) // All pending
	lsn.latestHead = 2
	logs = lsn.extractConfirmedLogs()
	assert.Equal(t, 2, len(logs))     // 1 and 2 should be confirmed
	assert.Equal(t, 1, len(lsn.reqs)) // 3 is still pending
	assert.Equal(t, uint64(3), lsn.reqs[0].confirmedAtBlock)
	// Another block way in the future should clear it
	lsn.latestHead = 10
	logs = lsn.extractConfirmedLogs()
	assert.Equal(t, 1, len(logs))     // remaining log
	assert.Equal(t, 0, len(lsn.reqs)) // all processed
}

func TestResponsePruning(t *testing.T) {
	lsn := listenerV1{}
	lsn.latestHead = 10000
	lsn.respCount = map[[32]byte]uint64{
		utils.PadByteToHash(0x00): 1,
		utils.PadByteToHash(0x01): 1,
	}
	lsn.blockNumberToReqID = pairing.New()
	lsn.blockNumberToReqID.Insert(fulfilledReq{
		blockNumber: 1,
		reqID:       utils.PadByteToHash(0x00),
	})
	lsn.blockNumberToReqID.Insert(fulfilledReq{
		blockNumber: 2,
		reqID:       utils.PadByteToHash(0x01),
	})
	lsn.pruneConfirmedRequestCounts()
	assert.Equal(t, 2, len(lsn.respCount))
	lsn.latestHead = 10001
	lsn.pruneConfirmedRequestCounts()
	assert.Equal(t, 1, len(lsn.respCount))
	lsn.latestHead = 10002
	lsn.pruneConfirmedRequestCounts()
	assert.Equal(t, 0, len(lsn.respCount))
}
