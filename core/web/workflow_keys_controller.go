package web

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

// WorkflowKeysController exposes the workflow key
type WorkflowKeysController struct {
	App chainlink.Application
}

// Index lists workflow keys
// Example:
// "GET <application>/keys/workflow"
func (wfkc *WorkflowKeysController) Index(c *gin.Context) {
	keys, err := wfkc.App.GetKeyStore().Workflow().GetAll()
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	jsonAPIResponse(c, presenters.NewWorkflowKeyResources(keys), "workflowKey")
}
