package pb_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
)

const (
	testWorkflowID  = "test-id-1"
	testConfigKey   = "test-key"
	testConfigValue = "test-value"
	testInputsKey   = "input-key"
	testInputsValue = "input-value"
	testError       = "test-error"
)

func TestMarshalUnmarshalRequest(t *testing.T) {
	req := capabilities.CapabilityRequest{
		Metadata: capabilities.RequestMetadata{
			WorkflowExecutionID: testWorkflowID,
		},
		Config: &values.Map{Underlying: map[string]values.Value{
			testConfigKey: &values.String{Underlying: testConfigValue},
		}},
		Inputs: &values.Map{Underlying: map[string]values.Value{
			testInputsKey: &values.String{Underlying: testInputsValue},
		}},
	}
	raw, err := pb.MarshalCapabilityRequest(req)
	require.NoError(t, err)

	unmarshaled, err := pb.UnmarshalCapabilityRequest(raw)
	require.NoError(t, err)

	require.Equal(t, req, unmarshaled)
}

func TestMarshalUnmarshalResponse(t *testing.T) {
	resp := capabilities.CapabilityResponse{
		Value: &values.String{Underlying: testConfigValue},
		Err:   errors.New(testError),
	}
	raw, err := pb.MarshalCapabilityResponse(resp)
	require.NoError(t, err)

	unmarshaled, err := pb.UnmarshalCapabilityResponse(raw)
	require.NoError(t, err)

	require.Equal(t, resp, unmarshaled)
}
