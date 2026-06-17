package confidentialrelay

import (
	"sync"

	"github.com/smartcontractkit/chainlink-common/pkg/workflows/host"
)

type ExecutionHandlers struct {
	handlers sync.Map
}

func (e *ExecutionHandlers) AddExecution(workflowId, execId string, helper host.ExecutionHelperWithRawSecrets) {
	e.handlers.Store(wfexecid(workflowId, execId), helper)
}

func (e *ExecutionHandlers) RemoveExecution(workflowId, execId string) {
	e.handlers.Delete(wfexecid(workflowId, execId))
}

func (e *ExecutionHandlers) GetExecution(workflowId, execId string) (host.ExecutionHelperWithRawSecrets, bool) {
	value, ok := e.handlers.Load(wfexecid(workflowId, execId))
	if !ok {
		return nil, false
	}

	helper, ok := value.(host.ExecutionHelperWithRawSecrets)
	return helper, ok
}

func wfexecid(workflowId, execId string) string {
	return workflowId + "." + execId
}
