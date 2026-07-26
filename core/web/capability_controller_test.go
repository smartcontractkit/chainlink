package web_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	registrymock "github.com/smartcontractkit/chainlink-common/pkg/types/core/mocks"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities"
	capmock "github.com/smartcontractkit/chainlink/v2/core/capabilities/mocks"
	appmocks "github.com/smartcontractkit/chainlink/v2/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/web"
)

func TestCapabilityController_ExecuteCapability_MissingBody(t *testing.T) {
	t.Parallel()
	mockApp := appmocks.NewApplication(t)

	controller := web.CapabilityController{App: mockApp}

	w := httptest.NewRecorder()
	engine := gin.New()
	engine.POST("/v2/capabilities/execute", controller.ExecuteCapability)

	req, err := http.NewRequestWithContext(t.Context(), "POST", "/v2/capabilities/execute", nil)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCapabilityController_ExecuteCapability_RegistryNotInitialized(t *testing.T) {
	t.Parallel()
	mockApp := appmocks.NewApplication(t)
	requestBody := web.CapabilityRequestOuter{
		CapabilityName:    "test-capability",
		CapabilityRequest: []byte(`{"test": "request"}`),
	}
	mockApp.EXPECT().GetCapabilitiesRegistry().Return(nil)

	controller := web.CapabilityController{App: mockApp}

	w := httptest.NewRecorder()
	engine := gin.New()
	engine.POST("/v2/capabilities/execute", controller.ExecuteCapability)

	reqJSON, err := json.Marshal(requestBody)
	require.NoError(t, err)
	req, err := http.NewRequestWithContext(t.Context(), "POST", "/v2/capabilities/execute", bytes.NewBuffer(reqJSON))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCapabilityController_ExecuteCapability_MissingRequiredFields(t *testing.T) {
	t.Parallel()
	mockApp := appmocks.NewApplication(t)
	mockApp.EXPECT().GetCapabilitiesRegistry().Return(nil)

	controller := web.CapabilityController{App: mockApp}

	w := httptest.NewRecorder()
	engine := gin.New()
	engine.POST("/v2/capabilities/execute", controller.ExecuteCapability)

	invalidRequest := `{"capabilityName": ""}` // missing capabilityRequest
	req, err := http.NewRequestWithContext(t.Context(), "POST", "/v2/capabilities/execute", bytes.NewBufferString(invalidRequest))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCapabilityController_ExecuteCapability(t *testing.T) {
	t.Parallel()
	mockApp := appmocks.NewApplication(t)
	mockRegistry := registrymock.NewCapabilitiesRegistry(t)
	mockApp.EXPECT().GetCapabilitiesRegistry().Return(&capabilities.Registry{
		CapabilitiesRegistryBase: mockRegistry,
	})

	executableCap := capmock.NewExecutableCapability(t)
	capabilityName := "test-capability"
	mockRegistry.EXPECT().GetExecutable(mock.Anything, capabilityName).Return(executableCap, nil)
	expectedResponse := commoncap.CapabilityResponse{
		Value: &values.Map{Underlying: map[string]values.Value{
			"result": &values.String{Underlying: "success"},
		}},
	}
	executableCap.EXPECT().Execute(mock.Anything, mock.AnythingOfType("capabilities.CapabilityRequest")).Return(expectedResponse, nil)

	controller := web.CapabilityController{App: mockApp}

	inputsMap, err := values.NewMap(map[string]any{
		"test": "input",
	})
	require.NoError(t, err)

	configMap, err := values.NewMap(map[string]any{
		"config": "value",
	})
	require.NoError(t, err)

	capabilityRequest := commoncap.CapabilityRequest{
		Metadata: commoncap.RequestMetadata{
			WorkflowID:          "test-workflow",
			WorkflowExecutionID: "test-execution",
		},
		Config: configMap,
		Inputs: inputsMap,
	}

	capabilityRequestBytes, err := pb.MarshalCapabilityRequest(capabilityRequest)
	require.NoError(t, err)
	requestBody := web.CapabilityRequestOuter{
		CapabilityName:    capabilityName,
		CapabilityRequest: capabilityRequestBytes,
	}
	reqJSON, err := json.Marshal(requestBody)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	engine := gin.New()
	engine.POST("/v2/capabilities/execute", controller.ExecuteCapability)

	req, err := http.NewRequestWithContext(t.Context(), "POST", "/v2/capabilities/execute", bytes.NewBuffer(reqJSON))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]any
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response, "capabilityResponse")
}
