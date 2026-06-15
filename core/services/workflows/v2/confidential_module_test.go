package v2

import (
	"context"
	"crypto/sha256"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	regmocks "github.com/smartcontractkit/chainlink-common/pkg/types/core/mocks"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/host"

	confworkflowtypes "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/actions/confidentialworkflow"
	capmocks "github.com/smartcontractkit/chainlink/v2/core/capabilities/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/utils/matches"

	sdkpb "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	valuespb "github.com/smartcontractkit/chainlink-protos/cre/go/values/pb"
	wfpb "github.com/smartcontractkit/chainlink-protos/workflows/go/v2"
)

// stubExecutionHelper implements host.ExecutionHelper for testing.
type stubExecutionHelper struct {
	executionID string
}

var _ host.ExecutionHelper = (*stubExecutionHelper)(nil)

func (s *stubExecutionHelper) CallCapability(context.Context, *sdkpb.CapabilityRequest) (*sdkpb.CapabilityResponse, error) {
	return nil, nil
}
func (s *stubExecutionHelper) GetSecrets(context.Context, *sdkpb.GetSecretsRequest) ([]*sdkpb.SecretResponse, error) {
	return nil, nil
}
func (s *stubExecutionHelper) GetWorkflowExecutionID() string { return s.executionID }
func (s *stubExecutionHelper) GetNodeTime() time.Time         { return time.Time{} }
func (s *stubExecutionHelper) GetDONTime() (time.Time, error) { return time.Time{}, nil }
func (s *stubExecutionHelper) EmitUserLog(string) error       { return nil }
func (s *stubExecutionHelper) EmitUserMetric(context.Context, *wfpb.WorkflowUserMetric) error {
	return nil
}

func TestParseWorkflowAttributes(t *testing.T) {
	t.Run("valid JSON with all fields", func(t *testing.T) {
		data := []byte(`{"confidential":true,"vault_don_secrets":[{"key":"API_KEY"},{"key":"SIGNING_KEY","namespace":"custom-ns"}]}`)
		attrs, err := ParseWorkflowAttributes(data)
		require.NoError(t, err)
		assert.True(t, attrs.Confidential)
		require.Len(t, attrs.VaultDonSecrets, 2)
		assert.Equal(t, "API_KEY", attrs.VaultDonSecrets[0].Key)
		assert.Empty(t, attrs.VaultDonSecrets[0].Namespace)
		assert.Equal(t, "SIGNING_KEY", attrs.VaultDonSecrets[1].Key)
		assert.Equal(t, "custom-ns", attrs.VaultDonSecrets[1].Namespace)
	})

	t.Run("empty data returns zero value", func(t *testing.T) {
		attrs, err := ParseWorkflowAttributes(nil)
		require.NoError(t, err)
		assert.False(t, attrs.Confidential)
		assert.Nil(t, attrs.VaultDonSecrets)

		attrs, err = ParseWorkflowAttributes([]byte{})
		require.NoError(t, err)
		assert.False(t, attrs.Confidential)
	})

	t.Run("non-confidential workflow", func(t *testing.T) {
		data := []byte(`{"confidential":false}`)
		attrs, err := ParseWorkflowAttributes(data)
		require.NoError(t, err)
		assert.False(t, attrs.Confidential)
	})

	t.Run("malformed JSON returns error", func(t *testing.T) {
		_, err := ParseWorkflowAttributes([]byte(`{not json}`))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse workflow attributes")
	})
}

func TestComputeBinaryHash(t *testing.T) {
	binary := []byte("hello world")
	hash := ComputeBinaryHash(binary)
	expected := sha256.Sum256(binary)
	assert.Equal(t, expected[:], hash)

	// Deterministic: same input produces same hash.
	assert.Equal(t, hash, ComputeBinaryHash(binary))
}

func TestConfidentialModule_Execute(t *testing.T) {
	ctx := context.Background()
	lggr := logger.Nop()

	// Build an ExecuteRequest to send through the module.
	execReq := &sdkpb.ExecuteRequest{
		Config: []byte("test-config"),
	}

	// Build the expected ExecutionResult that the enclave returns.
	expectedResult := &sdkpb.ExecutionResult{
		Result: &sdkpb.ExecutionResult_Value{
			Value: valuespb.NewStringValue("enclave-output"),
		},
	}

	confResp := &confworkflowtypes.ConfidentialWorkflowResponse{SdkExecutionResult: expectedResult}
	respPayload, err := anypb.New(confResp)
	require.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		capReg := regmocks.NewCapabilitiesRegistry(t)
		execCap := capmocks.NewExecutableCapability(t)

		capReg.EXPECT().GetExecutable(matches.AnyContext, confidentialWorkflowsCapabilityID).
			Return(execCap, nil).Once()

		execCap.EXPECT().Execute(matches.AnyContext, mock.MatchedBy(func(req capabilities.CapabilityRequest) bool {
			return req.Method == "Execute" &&
				req.CapabilityId == confidentialWorkflowsCapabilityID &&
				req.Metadata.WorkflowID == "wf-123" &&
				req.Metadata.WorkflowOwner == "owner-abc" &&
				req.Metadata.WorkflowExecutionID == "exec-456" &&
				req.Payload != nil
		})).Return(capabilities.CapabilityResponse{Payload: respPayload}, nil).Once()

		mod := NewConfidentialModule(
			capReg,
			"https://example.com/binary.wasm",
			[]byte("fakehash"),
			"wf-123",
			"owner-abc",
			"my-workflow",
			"v1",
			lggr,
		)

		result, err := mod.Execute(ctx, execReq, &stubExecutionHelper{executionID: "exec-456"})
		require.NoError(t, err)
		require.NotNil(t, result)

		val := result.GetValue()
		require.NotNil(t, val)
		assert.Equal(t, "enclave-output", val.GetStringValue())
	})

	t.Run("GetExecutable error", func(t *testing.T) {
		capReg := regmocks.NewCapabilitiesRegistry(t)
		capReg.EXPECT().GetExecutable(matches.AnyContext, confidentialWorkflowsCapabilityID).
			Return(nil, errors.New("capability not found")).Once()

		mod := NewConfidentialModule(capReg, "", nil, "wf", "owner", "name", "tag", lggr)

		_, err := mod.Execute(ctx, execReq, &stubExecutionHelper{})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get confidential-workflows capability")
	})

	t.Run("capability Execute error", func(t *testing.T) {
		capReg := regmocks.NewCapabilitiesRegistry(t)
		execCap := capmocks.NewExecutableCapability(t)

		capReg.EXPECT().GetExecutable(matches.AnyContext, confidentialWorkflowsCapabilityID).
			Return(execCap, nil).Once()
		execCap.EXPECT().Execute(matches.AnyContext, mock.Anything).
			Return(capabilities.CapabilityResponse{}, errors.New("enclave unavailable")).Once()

		mod := NewConfidentialModule(capReg, "", nil, "wf", "owner", "name", "tag", lggr)

		_, err := mod.Execute(ctx, execReq, &stubExecutionHelper{})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "confidential-workflows capability execution failed")
	})

	t.Run("nil payload in response", func(t *testing.T) {
		capReg := regmocks.NewCapabilitiesRegistry(t)
		execCap := capmocks.NewExecutableCapability(t)

		capReg.EXPECT().GetExecutable(matches.AnyContext, confidentialWorkflowsCapabilityID).
			Return(execCap, nil).Once()
		execCap.EXPECT().Execute(matches.AnyContext, mock.Anything).
			Return(capabilities.CapabilityResponse{Payload: nil}, nil).Once()

		mod := NewConfidentialModule(capReg, "", nil, "wf", "owner", "name", "tag", lggr)

		_, err := mod.Execute(ctx, execReq, &stubExecutionHelper{})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "returned nil payload")
	})

	t.Run("request fields are forwarded correctly", func(t *testing.T) {
		capReg := regmocks.NewCapabilitiesRegistry(t)
		execCap := capmocks.NewExecutableCapability(t)

		capReg.EXPECT().GetExecutable(matches.AnyContext, confidentialWorkflowsCapabilityID).
			Return(execCap, nil).Once()

		var capturedReq capabilities.CapabilityRequest
		execCap.EXPECT().Execute(matches.AnyContext, mock.Anything).
			Run(func(_ context.Context, req capabilities.CapabilityRequest) {
				capturedReq = req
			}).
			Return(capabilities.CapabilityResponse{Payload: respPayload}, nil).Once()

		binaryHash := ComputeBinaryHash([]byte("some-binary"))
		mod := NewConfidentialModule(capReg, "https://example.com/wasm", binaryHash, "wf-abc", "0xowner", "my-workflow", "v2", lggr)

		_, err := mod.Execute(ctx, execReq, &stubExecutionHelper{executionID: "exec-xyz"})
		require.NoError(t, err)

		assert.Equal(t, "Execute", capturedReq.Method)
		assert.Equal(t, confidentialWorkflowsCapabilityID, capturedReq.CapabilityId)
		assert.Equal(t, "wf-abc", capturedReq.Metadata.WorkflowID)
		assert.Equal(t, "0xowner", capturedReq.Metadata.WorkflowOwner)
		assert.Equal(t, "my-workflow", capturedReq.Metadata.WorkflowName)
		assert.Equal(t, "v2", capturedReq.Metadata.WorkflowTag)
		assert.Equal(t, "exec-xyz", capturedReq.Metadata.WorkflowExecutionID)

		var confReq confworkflowtypes.ConfidentialWorkflowRequest
		require.NoError(t, capturedReq.Payload.UnmarshalTo(&confReq))

		assert.Equal(t, "wf-abc", confReq.Execution.WorkflowId)
		assert.Equal(t, "https://example.com/wasm", confReq.Execution.BinaryUrl)
		assert.Equal(t, binaryHash, confReq.Execution.BinaryHash)

		assert.Equal(t, execReq.GetConfig(), confReq.Execution.SdkExecuteRequest.GetConfig())
	})
}

func TestConfidentialModule_Tee(t *testing.T) {
	ctx := context.Background()
	lggr := logger.Nop()

	buildRespPayload := func(t *testing.T, tees []*sdkpb.TeeTypeAndRegions) *anypb.Any {
		t.Helper()
		payload, err := anypb.New(&confworkflowtypes.ProvidedTeesResponse{Tee: tees})
		require.NoError(t, err)
		return payload
	}

	anyRegionsTee := func(regions ...string) *sdkpb.Tee {
		return &sdkpb.Tee{
			Item: &sdkpb.Tee_AnyRegions{
				AnyRegions: &sdkpb.Regions{Regions: regions},
			},
		}
	}

	t.Run("matching region returns true", func(t *testing.T) {
		capReg := regmocks.NewCapabilitiesRegistry(t)
		execCap := capmocks.NewExecutableCapability(t)

		capReg.EXPECT().GetExecutable(matches.AnyContext, confidentialWorkflowsCapabilityID).
			Return(execCap, nil).Once()
		execCap.EXPECT().Execute(matches.AnyContext, mock.MatchedBy(func(req capabilities.CapabilityRequest) bool {
			return req.Method == "ProvidedTees" &&
				req.CapabilityId == confidentialWorkflowsCapabilityID &&
				req.Metadata.WorkflowExecutionID == ""
		})).Return(capabilities.CapabilityResponse{Payload: buildRespPayload(t, []*sdkpb.TeeTypeAndRegions{
			{Type: sdkpb.TeeType_TEE_TYPE_AWS_NITRO, Regions: []string{"us-east-1", "eu-west-1"}},
		})}, nil).Once()

		mod := NewConfidentialModule(capReg, "https://example.com/binary.wasm", []byte("fakehash"), "wf-123", "owner-abc", "my-workflow", "v1", lggr)

		assert.True(t, mod.Tee(ctx, anyRegionsTee("us-east-1")))
	})

	t.Run("non-matching region returns false", func(t *testing.T) {
		capReg := regmocks.NewCapabilitiesRegistry(t)
		execCap := capmocks.NewExecutableCapability(t)

		capReg.EXPECT().GetExecutable(matches.AnyContext, confidentialWorkflowsCapabilityID).
			Return(execCap, nil).Once()
		execCap.EXPECT().Execute(matches.AnyContext, mock.Anything).
			Return(capabilities.CapabilityResponse{Payload: buildRespPayload(t, []*sdkpb.TeeTypeAndRegions{
				{Type: sdkpb.TeeType_TEE_TYPE_AWS_NITRO, Regions: []string{"us-east-1"}},
			})}, nil).Once()

		mod := NewConfidentialModule(capReg, "", nil, "wf", "owner", "name", "tag", lggr)
		assert.False(t, mod.Tee(ctx, anyRegionsTee("ap-southeast-1")))
	})

	t.Run("empty tees response returns false", func(t *testing.T) {
		capReg := regmocks.NewCapabilitiesRegistry(t)
		execCap := capmocks.NewExecutableCapability(t)

		capReg.EXPECT().GetExecutable(matches.AnyContext, confidentialWorkflowsCapabilityID).
			Return(execCap, nil).Once()
		execCap.EXPECT().Execute(matches.AnyContext, mock.Anything).
			Return(capabilities.CapabilityResponse{Payload: buildRespPayload(t, nil)}, nil).Once()

		mod := NewConfidentialModule(capReg, "", nil, "wf", "owner", "name", "tag", lggr)
		assert.False(t, mod.Tee(ctx, anyRegionsTee("us-east-1")))
	})

	t.Run("GetExecutable error returns false", func(t *testing.T) {
		capReg := regmocks.NewCapabilitiesRegistry(t)
		capReg.EXPECT().GetExecutable(matches.AnyContext, confidentialWorkflowsCapabilityID).
			Return(nil, errors.New("capability not found")).Once()

		mod := NewConfidentialModule(capReg, "", nil, "wf", "owner", "name", "tag", lggr)
		assert.False(t, mod.Tee(ctx, anyRegionsTee("us-east-1")))
	})

	t.Run("capability Execute error returns false", func(t *testing.T) {
		capReg := regmocks.NewCapabilitiesRegistry(t)
		execCap := capmocks.NewExecutableCapability(t)

		capReg.EXPECT().GetExecutable(matches.AnyContext, confidentialWorkflowsCapabilityID).
			Return(execCap, nil).Once()
		execCap.EXPECT().Execute(matches.AnyContext, mock.Anything).
			Return(capabilities.CapabilityResponse{}, errors.New("enclave unavailable")).Once()

		mod := NewConfidentialModule(capReg, "", nil, "wf", "owner", "name", "tag", lggr)
		assert.False(t, mod.Tee(ctx, anyRegionsTee("us-east-1")))
	})

	t.Run("nil payload returns false", func(t *testing.T) {
		capReg := regmocks.NewCapabilitiesRegistry(t)
		execCap := capmocks.NewExecutableCapability(t)

		capReg.EXPECT().GetExecutable(matches.AnyContext, confidentialWorkflowsCapabilityID).
			Return(execCap, nil).Once()
		execCap.EXPECT().Execute(matches.AnyContext, mock.Anything).
			Return(capabilities.CapabilityResponse{Payload: nil}, nil).Once()

		mod := NewConfidentialModule(capReg, "", nil, "wf", "owner", "name", "tag", lggr)
		assert.False(t, mod.Tee(ctx, anyRegionsTee("us-east-1")))
	})

	t.Run("request fields are correct", func(t *testing.T) {
		capReg := regmocks.NewCapabilitiesRegistry(t)
		execCap := capmocks.NewExecutableCapability(t)

		capReg.EXPECT().GetExecutable(matches.AnyContext, confidentialWorkflowsCapabilityID).
			Return(execCap, nil).Once()

		var capturedReq capabilities.CapabilityRequest
		execCap.EXPECT().Execute(matches.AnyContext, mock.Anything).
			Run(func(_ context.Context, req capabilities.CapabilityRequest) {
				capturedReq = req
			}).
			Return(capabilities.CapabilityResponse{Payload: buildRespPayload(t, nil)}, nil).Once()

		mod := NewConfidentialModule(capReg, "https://example.com/wasm", []byte("hash"), "wf-xyz", "0xowner", "my-workflow", "v3", lggr)
		_ = mod.Tee(ctx, anyRegionsTee("us-east-1"))

		assert.Equal(t, "ProvidedTees", capturedReq.Method)
		assert.Equal(t, confidentialWorkflowsCapabilityID, capturedReq.CapabilityId)
		assert.Equal(t, "wf-xyz", capturedReq.Metadata.WorkflowID)
		assert.Equal(t, "0xowner", capturedReq.Metadata.WorkflowOwner)
		assert.Equal(t, "my-workflow", capturedReq.Metadata.WorkflowName)
		assert.Equal(t, "v3", capturedReq.Metadata.WorkflowTag)
		assert.Empty(t, capturedReq.Metadata.WorkflowExecutionID)

		var emptyMsg emptypb.Empty
		require.NoError(t, capturedReq.Payload.UnmarshalTo(&emptyMsg))
	})

	t.Run("caches provider across calls", func(t *testing.T) {
		capReg := regmocks.NewCapabilitiesRegistry(t)
		execCap := capmocks.NewExecutableCapability(t)

		capReg.EXPECT().GetExecutable(matches.AnyContext, confidentialWorkflowsCapabilityID).
			Return(execCap, nil).Once()
		execCap.EXPECT().Execute(matches.AnyContext, mock.Anything).
			Return(capabilities.CapabilityResponse{Payload: buildRespPayload(t, []*sdkpb.TeeTypeAndRegions{
				{Type: sdkpb.TeeType_TEE_TYPE_AWS_NITRO, Regions: []string{"us-east-1"}},
			})}, nil).Once()

		mod := NewConfidentialModule(capReg, "", nil, "wf", "owner", "name", "tag", lggr)

		assert.True(t, mod.Tee(ctx, anyRegionsTee("us-east-1")))
		assert.True(t, mod.Tee(ctx, anyRegionsTee("us-east-1")))
		assert.False(t, mod.Tee(ctx, anyRegionsTee("eu-west-1")))
	})
}

func TestConfidentialModule_SetRequirements(t *testing.T) {
	ctx := context.Background()
	lggr := logger.Nop()

	execReq := &sdkpb.ExecuteRequest{
		Config: []byte("test-config"),
	}

	expectedResult := &sdkpb.ExecutionResult{
		Result: &sdkpb.ExecutionResult_Value{
			Value: valuespb.NewStringValue("enclave-output"),
		},
	}

	confResp := &confworkflowtypes.ConfidentialWorkflowResponse{SdkExecutionResult: expectedResult}
	respPayload, err := anypb.New(confResp)
	require.NoError(t, err)

	t.Run("requirements forwarded in execute", func(t *testing.T) {
		capReg := regmocks.NewCapabilitiesRegistry(t)
		execCap := capmocks.NewExecutableCapability(t)

		capReg.EXPECT().GetExecutable(matches.AnyContext, confidentialWorkflowsCapabilityID).
			Return(execCap, nil).Once()

		var capturedReq capabilities.CapabilityRequest
		execCap.EXPECT().Execute(matches.AnyContext, mock.Anything).
			Run(func(_ context.Context, req capabilities.CapabilityRequest) {
				capturedReq = req
			}).
			Return(capabilities.CapabilityResponse{Payload: respPayload}, nil).Once()

		mod := NewConfidentialModule(capReg, "", nil, "wf", "owner", "name", "tag", lggr)

		requirements := &sdkpb.Requirements{
			Tee: &sdkpb.Tee{
				Item: &sdkpb.Tee_AnyRegions{
					AnyRegions: &sdkpb.Regions{Regions: []string{"us-east-1"}},
				},
			},
		}
		mod.SetRequirements("exec-789", requirements)

		result, err := mod.Execute(ctx, execReq, &stubExecutionHelper{executionID: "exec-789"})
		require.NoError(t, err)
		require.NotNil(t, result)

		var confReq confworkflowtypes.ConfidentialWorkflowRequest
		require.NoError(t, capturedReq.Payload.UnmarshalTo(&confReq))
		require.NotNil(t, confReq.Execution.Requirements)
		assert.NotNil(t, confReq.Execution.Requirements.Tee)
	})

	t.Run("requirements consumed after execute", func(t *testing.T) {
		capReg := regmocks.NewCapabilitiesRegistry(t)
		execCap := capmocks.NewExecutableCapability(t)

		capReg.EXPECT().GetExecutable(matches.AnyContext, confidentialWorkflowsCapabilityID).
			Return(execCap, nil).Times(2)

		var secondReq capabilities.CapabilityRequest
		execCap.EXPECT().Execute(matches.AnyContext, mock.Anything).
			Return(capabilities.CapabilityResponse{Payload: respPayload}, nil).Once()
		execCap.EXPECT().Execute(matches.AnyContext, mock.Anything).
			Run(func(_ context.Context, req capabilities.CapabilityRequest) {
				secondReq = req
			}).
			Return(capabilities.CapabilityResponse{Payload: respPayload}, nil).Once()

		mod := NewConfidentialModule(capReg, "", nil, "wf", "owner", "name", "tag", lggr)
		mod.SetRequirements("exec-789", &sdkpb.Requirements{})

		_, err := mod.Execute(ctx, execReq, &stubExecutionHelper{executionID: "exec-789"})
		require.NoError(t, err)

		_, err = mod.Execute(ctx, execReq, &stubExecutionHelper{executionID: "exec-789"})
		require.NoError(t, err)

		var confReq confworkflowtypes.ConfidentialWorkflowRequest
		require.NoError(t, secondReq.Payload.UnmarshalTo(&confReq))
		assert.Nil(t, confReq.Execution.Requirements)
	})
}

func TestConfidentialModule_InterfaceMethods(t *testing.T) {
	mod := &ConfidentialModule{}

	// These are no-ops but should not panic.
	mod.Start()
	mod.Close()
	assert.False(t, mod.IsLegacyDAG())
}
