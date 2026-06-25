package confidentialrelay

import (
	"sync"

	"github.com/smartcontractkit/chainlink-common/pkg/workflows/host"
)

type ExecutionHandlers struct {
	handlers sync.Map
}

func (e *ExecutionHandlers) AddExecution(workflowID, execId string, helper host.ExecutionHelperWithRawSecrets) {
	e.handlers.Store(wfexecid(workflowID, execId), helper)
}

func (e *ExecutionHandlers) RemoveExecution(workflowID, execId string) {
	e.handlers.Delete(wfexecid(workflowID, execId))
}

func (e *ExecutionHandlers) GetExecution(workflowID, execId string) (host.ExecutionHelperWithRawSecrets, bool) {
	value, ok := e.handlers.Load(wfexecid(workflowID, execId))
	if !ok {
		return nil, false
	}

	helper, ok := value.(host.ExecutionHelperWithRawSecrets)
	return helper, ok
}

func wfexecid(workflowID, execId string) string {
	return workflowID + "." + execId
}
