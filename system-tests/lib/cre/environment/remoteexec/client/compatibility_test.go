package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/remoteexec/agent"
)

func TestCheckCompatibilityStatusAcceptsSameMajor(t *testing.T) {
	err := checkCompatibilityStatus(&agent.AgentStatusResponse{
		ProtocolVersion: "1.3.0",
		Capabilities:    []string{"locks", "componentLogs"},
	}, []string{"locks"})
	require.NoError(t, err)
}

func TestCheckCompatibilityStatusRejectsDifferentMajor(t *testing.T) {
	err := checkCompatibilityStatus(&agent.AgentStatusResponse{
		ProtocolVersion: "2.0.0",
	}, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "incompatible protocol major versions")
}

func TestCheckCompatibilityStatusRejectsMissingRequiredCapability(t *testing.T) {
	err := checkCompatibilityStatus(&agent.AgentStatusResponse{
		ProtocolVersion: "1.0.0",
		Capabilities:    []string{"locks"},
	}, []string{"componentLogs"})
	require.Error(t, err)
	require.Contains(t, err.Error(), `required capability "componentLogs"`)
}

func TestCheckCompatibilityCallsStatusEndpoint(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_ = json.NewEncoder(w).Encode(agent.AgentStatusResponse{
			ProtocolVersion: "1.0.0",
			Capabilities:    []string{"locks", "componentLogs"},
		})
	}))
	defer server.Close()

	err := CheckCompatibility(context.Background(), &Runtime{AgentBaseURL: server.URL}, []string{"locks"})
	require.NoError(t, err)
}
