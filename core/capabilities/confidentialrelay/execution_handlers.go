package confidentialrelay

import (
	"sync"

	"github.com/smartcontractkit/chainlink-common/pkg/workflows/host"
)

type ExecutionHandlers struct {
	handlers sync.Map
}

func (e *ExecutionHandlers) AddExecution(workflowID, execID string, helper host.ExecutionHelperWithRawSecrets) {
	e.handlers.Store(wfexecID(workflowID, execID), helper)
}

func (e *ExecutionHandlers) RemoveExecution(workflowID, execID string) {
	e.handlers.Delete(wfexecID(workflowID, execID))
}

func (e *ExecutionHandlers) GetExecution(workflowID, execID string) (host.ExecutionHelperWithRawSecrets, bool) {
	value, ok := e.handlers.Load(wfexecID(workflowID, execID))
	if !ok {
		return nil, false
	}

	helper, ok := value.(host.ExecutionHelperWithRawSecrets)
	return helper, ok
}

func wfexecID(workflowID, execID string) string {
	return workflowID + "." + execID
}
