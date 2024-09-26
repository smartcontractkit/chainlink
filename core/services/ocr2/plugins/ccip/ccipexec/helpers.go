package ccipexec

import (
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/pkg/errors"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
)

// helper struct to hold the commitReport and the related send requests
type commitReportWithSendRequests struct {
	commitReport         cciptypes.CommitStoreReport
	sendRequestsWithMeta []cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta
}

func (r *commitReportWithSendRequests) validate() error {
	// make sure that number of messages is the expected
	if exp := int(r.commitReport.Interval.Max - r.commitReport.Interval.Min + 1); len(r.sendRequestsWithMeta) != exp {
		return errors.Errorf(
			"unexpected missing sendRequestsWithMeta in committed root %x have %d want %d", r.commitReport.MerkleRoot, len(r.sendRequestsWithMeta), exp)
	}

	return nil
}

// uniqueSenders returns slice of unique senders based on the send requests. Order is preserved based on the order of the send requests (by sequence number).
func (r *commitReportWithSendRequests) uniqueSenders() []cciptypes.Address {
	orderedUniqueSenders := make([]cciptypes.Address, 0, len(r.sendRequestsWithMeta))
	visitedSenders := mapset.NewSet[cciptypes.Address]()

	for _, req := range r.sendRequestsWithMeta {
		if !visitedSenders.Contains(req.Sender) {
			orderedUniqueSenders = append(orderedUniqueSenders, req.Sender)
			visitedSenders.Add(req.Sender)
		}
	}
	return orderedUniqueSenders
}

func (r *commitReportWithSendRequests) allRequestsAreExecutedAndFinalized() bool {
	for _, req := range r.sendRequestsWithMeta {
		if !req.Executed || !req.Finalized {
			return false
		}
	}
	return true
}

// checks if the send request fits the commit report interval
func (r *commitReportWithSendRequests) sendReqFits(sendReq cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta) bool {
	return sendReq.SequenceNumber >= r.commitReport.Interval.Min &&
		sendReq.SequenceNumber <= r.commitReport.Interval.Max
}
