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

	donTime := <-tp.donTimeStore.RequestDonTime(tp.workflowExecutionID, tp.timeRequestNum)

	if donTime.Err != nil {
		// TODO: Handle error or timeout; do we still want to increment timeRequestNum?
		return tp.GetNodeTime()
	}
	return fromUnixMilli(donTime.Timestamp)
}

func fromUnixMilli(ms int64) time.Time {
	return time.Unix(0, ms*int64(time.Millisecond)).UTC()
}
