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
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	regmocks "github.com/smartcontractkit/chainlink-common/pkg/types/core/mocks"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host"

	confworkflowtypes "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/actions/confidentialworkflow"
	capmocks "github.com/smartcontractkit/chainlink/v2/core/capabilities/mocks"
	workflowtypes "github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
	"github.com/smartcontractkit/chainlink/v2/core/utils/matches"

	sdkpb "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	storage_service "github.com/smartcontractkit/chainlink-protos/storage-service/go"
	valuespb "github.com/smartcontractkit/chainlink-protos/cre/go/values/pb"
	wfpb "github.com/smartcontractkit/chainlink-protos/workflows/go/v2"
)

// stubRetrieveURL returns a LocationRetrieverFunc that always returns the given
// URL, simulating the storage service's NodeService.DownloadArtifact presigned
// URL response without a real storage backend.
func stubRetrieveURL(url string) workflowtypes.LocationRetrieverFunc {
	return func(context.Context, *storage_service.DownloadArtifactRequest) (string, error) {
		return url, nil
	}
}

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
	t.Run("confidential true", func(t *testing.T) {
		data := []byte(`{"confidential":true}`)
		attrs, err := ParseWorkflowAttributes(data)
		require.NoError(t, err)
		assert.True(t, attrs.Confidential)
	})

	t.Run("empty data returns zero value", func(t *testing.T) {
		attrs, err := ParseWorkflowAttributes(nil)
		require.NoError(t, err)
		assert.False(t, attrs.Confidential)

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

	t.Run("unknown keys are ignored", func(t *testing.T) {
		// vault_don_secrets is a leftover key from a previous schema. Parsing
		// must tolerate it without failing.
		data := []byte(`{"confidential":true,"vault_don_secrets":[{"key":"X"}]}`)
		attrs, err := ParseWorkflowAttributes(data)
		require.NoError(t, err)
		assert.True(t, attrs.Confidential)
	})
}

func TestIsConfidential(t *testing.T) {
	t.Run("returns true for confidential", func(t *testing.T) {
		ok, err := IsConfidential([]byte(`{"confidential":true}`))
		require.NoError(t, err)
		assert.True(t, ok)
	})

	t.Run("returns false for non-confidential", func(t *testing.T) {
		ok, err := IsConfidential([]byte(`{"confidential":false}`))
		require.NoError(t, err)
		assert.False(t, ok)
	})

	t.Run("returns false for empty data", func(t *testing.T) {
		ok, err := IsConfidential(nil)
		require.NoError(t, err)
		assert.False(t, ok)
	})

	t.Run("returns error for malformed JSON", func(t *testing.T) {
		_, err := IsConfidential([]byte(`broken`))
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

	// Serialize the result into a ConfidentialWorkflowResponse, as the capability would.
	resultBytes, err := proto.Marshal(expectedResult)
	require.NoError(t, err)

	confResp := &confworkflowtypes.ConfidentialWorkflowResponse{
		ExecutionResult: resultBytes,
	}
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
			[]byte("fakehash"),
			"wf-123",
			"owner-abc",
			"my-workflow",
			"v1",
			stubRetrieveURL("https://example.com/binary.wasm"),
			nil,
			"",
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

		mod := NewConfidentialModule(
			capReg, nil, "wf", "owner", "name", "tag", stubRetrieveURL("https://example.com/x"), nil, "", lggr,
		)

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

		mod := NewConfidentialModule(
			capReg, nil, "wf", "owner", "name", "tag", stubRetrieveURL("https://example.com/x"), nil, "", lggr,
		)

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

		mod := NewConfidentialModule(
			capReg, nil, "wf", "owner", "name", "tag", stubRetrieveURL("https://example.com/x"), nil, "", lggr,
		)

		_, err := mod.Execute(ctx, execReq, &stubExecutionHelper{})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "returned nil payload")
	})

	t.Run("missing URL retriever errors", func(t *testing.T) {
		capReg := regmocks.NewCapabilitiesRegistry(t)
		mod := NewConfidentialModule(
			capReg, nil, "wf", "owner", "name", "tag", nil, nil, "", lggr,
		)

		_, err := mod.Execute(ctx, execReq, &stubExecutionHelper{})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "missing a URL retriever")
	})

	t.Run("retrieveURL error", func(t *testing.T) {
		capReg := regmocks.NewCapabilitiesRegistry(t)
		failingRetrieve := workflowtypes.LocationRetrieverFunc(
			func(context.Context, *storage_service.DownloadArtifactRequest) (string, error) {
				return "", errors.New("storage service down")
			},
		)
		mod := NewConfidentialModule(
			capReg, nil, "wf", "owner", "name", "tag", failingRetrieve, nil, "", lggr,
		)

		_, err := mod.Execute(ctx, execReq, &stubExecutionHelper{})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to mint pre-signed binary URL")
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
		mod := NewConfidentialModule(
			capReg,
			binaryHash,
			"wf-abc",
			"0xowner",
			"my-workflow",
			"v2",
			stubRetrieveURL("https://presigned.cloudfront.example.com/wasm?sig=abc"),
			nil,
			"",
			lggr,
		)

		_, err := mod.Execute(ctx, execReq, &stubExecutionHelper{executionID: "exec-xyz"})
		require.NoError(t, err)

		// Verify metadata.
		assert.Equal(t, "Execute", capturedReq.Method)
		assert.Equal(t, confidentialWorkflowsCapabilityID, capturedReq.CapabilityId)
		assert.Equal(t, "wf-abc", capturedReq.Metadata.WorkflowID)
		assert.Equal(t, "0xowner", capturedReq.Metadata.WorkflowOwner)
		assert.Equal(t, "my-workflow", capturedReq.Metadata.WorkflowName)
		assert.Equal(t, "v2", capturedReq.Metadata.WorkflowTag)
		assert.Equal(t, "exec-xyz", capturedReq.Metadata.WorkflowExecutionID)

		// Verify payload contents.
		var confReq confworkflowtypes.ConfidentialWorkflowRequest
		require.NoError(t, capturedReq.Payload.UnmarshalTo(&confReq))

		assert.Equal(t, "wf-abc", confReq.Execution.WorkflowId)
		// binary_url is the per-execution minted URL on the request (outside the
		// hashed WorkflowExecution); the in-hash Execution.binary_url is left unset.
		assert.Equal(t, "https://presigned.cloudfront.example.com/wasm?sig=abc", confReq.BinaryUrl)
		assert.Empty(t, confReq.Execution.BinaryUrl)
		assert.Equal(t, binaryHash, confReq.Execution.BinaryHash)

		// Verify the serialized ExecuteRequest round-trips.
		var roundTripped sdkpb.ExecuteRequest
		require.NoError(t, proto.Unmarshal(confReq.Execution.ExecuteRequest, &roundTripped))
		assert.Equal(t, execReq.GetConfig(), roundTripped.GetConfig())
	})
}

func TestConfidentialModule_InterfaceMethods(t *testing.T) {
	mod := &ConfidentialModule{}

	// These are no-ops but should not panic.
	mod.Start()
	mod.Close()
	assert.False(t, mod.IsLegacyDAG())
}
