package v2

import (
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/workflowLib"
)

type TimeProvider struct {
	workflowExecutionID string
	timeRequestNum      int
	donTimeStore        *workflowLib.DonTimeStore
}

func NewTimeProvider(workflowExecutionID string) TimeProvider {
	return TimeProvider{
		workflowExecutionID: workflowExecutionID,
		timeRequestNum:      0,
		donTimeStore:        workflowLib.GetDonTimeStore(),
	}
}

func (tp *TimeProvider) GetNodeTime() time.Time {
	return fromUnixMilli(tp.donTimeStore.GetLastObservedDonTime())
}

// GetDONTime makes a request to the WorkflowLib plugin store for DON Time
func (tp *TimeProvider) GetDONTime() time.Time {
	defer func() {
		tp.timeRequestNum++
	}()

	donTimeResp := <-tp.donTimeStore.RequestDonTime(tp.workflowExecutionID, tp.timeRequestNum)

	if donTimeResp.Err != nil {
		// An error implies a timeout occured on this request, which means this node did not include the request
		// in its observation. Consensus may still have been reached for this DON Time.
		if donTime := tp.donTimeStore.GetDonTimeForSeqNum(tp.workflowExecutionID, tp.timeRequestNum); donTime != nil {
			return fromUnixMilli(*donTime)
		}
		// Consensus was not reached for this DON Time. Return the last observed DON Time and be a faulty node.
		return tp.GetNodeTime()
	}

	return fromUnixMilli(donTimeResp.Timestamp)
}

func fromUnixMilli(ms int64) time.Time {
	return time.Unix(0, ms*int64(time.Millisecond)).UTC()
}
