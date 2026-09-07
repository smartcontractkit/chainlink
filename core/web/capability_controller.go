package web

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

type CapabilityController struct {
	App chainlink.Application
}

type CapabilityRequestOuter struct {
	CapabilityName    string `json:"capabilityName" binding:"required"`
	CapabilityRequest []byte `json:"capabilityRequest" binding:"required"`
}

// ExecuteCapability executes a capability by name with the provided request
// Example:
//
//	"<application>/v2/execute_capability"
func (cc *CapabilityController) ExecuteCapability(c *gin.Context) {
	body := c.Request.Body
	if body == nil {
		jsonAPIError(c, http.StatusBadRequest, errors.New("missing request body"))
		return
	}

	capabilityRegistry := cc.App.GetCapabilitiesRegistry()
	if capabilityRegistry == nil {
		jsonAPIError(c, http.StatusInternalServerError, errors.New("capability registry not initialized"))
		return
	}
	var capabilityRequestOuter CapabilityRequestOuter
	if err := c.BindJSON(&capabilityRequestOuter); err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}

	capability, err := capabilityRegistry.GetExecutable(c.Request.Context(), capabilityRequestOuter.CapabilityName)
	if err != nil {
		jsonAPIError(c, http.StatusNotFound, err)
		return
	}

	capabilityRequest, err := pb.UnmarshalCapabilityRequest(capabilityRequestOuter.CapabilityRequest)
	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}

	resp, err := capability.Execute(c.Request.Context(), capabilityRequest)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	responseBytes, err := pb.MarshalCapabilityResponse(resp)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"capabilityResponse": responseBytes})
}
