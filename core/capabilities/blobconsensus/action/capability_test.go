package action

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	caperrors "github.com/smartcontractkit/chainlink-common/pkg/capabilities/errors"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/blobconsensus/libs/testutils"
)

func Test_SimpleConsensus(t *testing.T) {
	lggr := logger.Test(t)
	ctx := t.Context()

	capability, err := NewConsensusCapability(lggr, clockwork.NewRealClock(), time.Minute, limits.Factory{Logger: lggr})
	require.NoError(t, err)

	oracleFactory := testutils.NewOracleFactory(t, lggr)

	err = capability.Initialise(ctx, core.StandardCapabilitiesDependencies{
		OracleFactory: oracleFactory,
	})
	require.NoError(t, err)

	servicetest.Run(t, capability)

	metadata := newRequestMetaData()

	input := &sdk.SimpleConsensusInputs{
		Observation: &sdk.SimpleConsensusInputs_Value{
			Value: values.Proto(values.NewInt64(10)),
		},
		Descriptors: &sdk.ConsensusDescriptor{
			Descriptor_: &sdk.ConsensusDescriptor_Aggregation{
				Aggregation: sdk.AggregationType_AGGREGATION_TYPE_MEDIAN,
			},
		},
		Default: nil,
	}

	result, err := capability.Simple(ctx, metadata, input)
	require.NoError(t, err)

	expectedResult, err := values.Wrap(10)
	require.NoError(t, err)
	expectedProto := values.Proto(expectedResult)

	require.True(t, proto.Equal(result.Response, expectedProto))
}

func Test_SimpleConsensus_Error(t *testing.T) {
	lggr := logger.Test(t)
	ctx := t.Context()

	capability, err := NewConsensusCapability(lggr, clockwork.NewRealClock(), time.Minute, limits.Factory{Logger: lggr})
	require.NoError(t, err)

	oracleFactory := testutils.NewOracleFactory(t, lggr)

	err = capability.Initialise(ctx, core.StandardCapabilitiesDependencies{
		OracleFactory: oracleFactory,
	})
	require.NoError(t, err)

	servicetest.Run(t, capability)

	metadata := newRequestMetaData()

	input := &sdk.SimpleConsensusInputs{
		Observation: &sdk.SimpleConsensusInputs_Error{
			Error: "its broken",
		},
		Descriptors: &sdk.ConsensusDescriptor{
			Descriptor_: &sdk.ConsensusDescriptor_Aggregation{
				Aggregation: sdk.AggregationType_AGGREGATION_TYPE_MEDIAN,
			},
		},
		Default: nil,
	}

	_, err = capability.Simple(ctx, metadata, input)

	var capError caperrors.Error
	ok := errors.As(err, &capError)
	require.True(t, ok, "error should be of type caperrors.Error")
	require.Equal(t, caperrors.ConsensusFailed, capError.Code())
	require.Equal(t, caperrors.OriginUser, capError.Origin())
	require.Equal(t, caperrors.VisibilityPublic, capError.Visibility())
}

func Test_Report(t *testing.T) {
	lggr := logger.Test(t)
	ctx := t.Context()

	capability, err := NewConsensusCapability(lggr, clockwork.NewRealClock(), time.Minute, limits.Factory{Logger: lggr})
	require.NoError(t, err)

	oracleFactory := testutils.NewOracleFactory(t, lggr)

	err = capability.Initialise(ctx, core.StandardCapabilitiesDependencies{
		OracleFactory: oracleFactory,
	})
	require.NoError(t, err)

	servicetest.Run(t, capability)

	metadata := newRequestMetaData()

	input := &sdk.ReportRequest{
		EncodedPayload: []byte("somerandom-payload"),
		EncoderName:    "evm",
		SigningAlgo:    "ecdsa",
		HashingAlgo:    "keccak256",
	}

	result, err := capability.Report(ctx, metadata, input)
	require.NoError(t, err)

	require.True(t, strings.HasSuffix(string(result.Response.RawReport), "somerandom-payload"))
}

func Test_ReportSupportsAptosSigningAndHashing(t *testing.T) {
	lggr := logger.Test(t)
	ctx := t.Context()

	capability, err := NewConsensusCapability(lggr, clockwork.NewRealClock(), time.Minute, limits.Factory{Logger: lggr})
	require.NoError(t, err)

	oracleFactory := testutils.NewOracleFactory(t, lggr)

	err = capability.Initialise(ctx, core.StandardCapabilitiesDependencies{
		OracleFactory: oracleFactory,
	})
	require.NoError(t, err)

	servicetest.Run(t, capability)

	metadata := newRequestMetaData()

	input := &sdk.ReportRequest{
		EncodedPayload: []byte("somerandom-payload"),
		EncoderName:    "aptos",
		SigningAlgo:    "ed25519",
		HashingAlgo:    "blake2b_256",
	}

	result, err := capability.Report(ctx, metadata, input)
	require.NoError(t, err)
	require.True(t, strings.HasSuffix(string(result.Response.RawReport), "somerandom-payload"))
}

func Test_ReportRequiresValidSigningAlgo(t *testing.T) {
	lggr := logger.Test(t)
	ctx := t.Context()

	capability, err := NewConsensusCapability(lggr, clockwork.NewRealClock(), time.Minute, limits.Factory{Logger: lggr})
	require.NoError(t, err)

	oracleFactory := testutils.NewOracleFactory(t, lggr)

	err = capability.Initialise(ctx, core.StandardCapabilitiesDependencies{
		OracleFactory: oracleFactory,
	})
	require.NoError(t, err)

	servicetest.Run(t, capability)

	metadata := newRequestMetaData()

	input := &sdk.ReportRequest{
		EncodedPayload: []byte("somerandom-payload"),
		EncoderName:    "evm",
		SigningAlgo:    "invalid-signing-algo",
		HashingAlgo:    "keccak256",
	}

	_, err = capability.Report(ctx, metadata, input)
	require.Error(t, err)
	require.ErrorContains(t, err, "unsupported signing algorithm")
}

func Test_ReportRequiresValidHashingAlgo(t *testing.T) {
	lggr := logger.Test(t)
	ctx := t.Context()

	capability, err := NewConsensusCapability(lggr, clockwork.NewRealClock(), time.Minute, limits.Factory{Logger: lggr})
	require.NoError(t, err)

	oracleFactory := testutils.NewOracleFactory(t, lggr)

	err = capability.Initialise(ctx, core.StandardCapabilitiesDependencies{
		OracleFactory: oracleFactory,
	})
	require.NoError(t, err)

	servicetest.Run(t, capability)

	metadata := newRequestMetaData()

	input := &sdk.ReportRequest{
		EncodedPayload: []byte("somerandom-payload"),
		EncoderName:    "evm",
		SigningAlgo:    "ecdsa",
		HashingAlgo:    "invalid-hashing-algo",
	}

	_, err = capability.Report(ctx, metadata, input)
	require.Error(t, err)
	require.ErrorContains(t, err, "unsupported hashing algorithm")
}

func Test_ReportRequiresValidEncoderName(t *testing.T) {
	lggr := logger.Test(t)
	ctx := t.Context()

	capability, err := NewConsensusCapability(lggr, clockwork.NewRealClock(), time.Minute, limits.Factory{Logger: lggr})
	require.NoError(t, err)

	oracleFactory := testutils.NewOracleFactory(t, lggr)

	err = capability.Initialise(ctx, core.StandardCapabilitiesDependencies{
		OracleFactory: oracleFactory,
	})
	require.NoError(t, err)

	servicetest.Run(t, capability)

	metadata := newRequestMetaData()

	input := &sdk.ReportRequest{
		EncodedPayload: []byte("somerandom-payload"),
		EncoderName:    "invalid-encoder-name",
		SigningAlgo:    "ecdsa",
		HashingAlgo:    "keccak256",
	}

	_, err = capability.Report(ctx, metadata, input)
	require.Error(t, err)
	require.ErrorContains(t, err, "unsupported encoder name")
}

func Test_SimpleInputsSizeValidation_UserErrorSentToConsensus(t *testing.T) {
	lggr := logger.Test(t)
	ctx := t.Context()

	capability, err := NewConsensusCapability(lggr, clockwork.NewRealClock(), time.Minute, limits.Factory{Logger: lggr})
	require.NoError(t, err)

	oracleFactory := testutils.NewOracleFactory(t, lggr)

	capConfig := &ConsensusCapabilityConfig{
		MaxRequestSizeBytes: 2,
	}

	capConfigJSON, err := json.Marshal(capConfig)
	require.NoError(t, err)

	err = capability.Initialise(ctx, core.StandardCapabilitiesDependencies{
		Config:        string(capConfigJSON),
		OracleFactory: oracleFactory,
	})
	require.NoError(t, err)

	servicetest.Run(t, capability)

	metadata := capabilities.RequestMetadata{
		WorkflowID:               "wf-id",
		WorkflowOwner:            "",
		WorkflowExecutionID:      "wex-id",
		WorkflowName:             "",
		WorkflowDonID:            0,
		WorkflowDonConfigVersion: 0,
		ReferenceID:              "1",
		DecodedWorkflowName:      "",
	}

	input := &sdk.SimpleConsensusInputs{
		Observation: &sdk.SimpleConsensusInputs_Value{Value: values.Proto(values.NewInt64(34))},
		Descriptors: &sdk.ConsensusDescriptor{Descriptor_: &sdk.ConsensusDescriptor_Aggregation{Aggregation: sdk.AggregationType_AGGREGATION_TYPE_IDENTICAL}},
	}

	// User error should be sent as error input to consensus, which will result in a consensus failure
	// if >= f+1 nodes have the same error, or consensus success if quorum has valid observations
	_, err = capability.Simple(ctx, metadata, input)
	require.Error(t, err)
	// The error should be a consensus failure since all nodes will have the same oversized request
	var capErr caperrors.Error
	errors.As(err, &capErr)
	require.Equal(t, caperrors.ConsensusFailed, capErr.Code())
}

func Test_SimpleRequestSizeValidation(t *testing.T) {
	lggr := logger.Test(t)
	ctx := t.Context()

	capability, err := NewConsensusCapability(lggr, clockwork.NewRealClock(), time.Minute, limits.Factory{Logger: lggr})
	require.NoError(t, err)

	oracleFactory := testutils.NewOracleFactory(t, lggr)

	// Set MaxRequestSizeBytes to 2MB
	capConfig := &ConsensusCapabilityConfig{
		MaxRequestSizeBytes: 2000000, // 2MB
	}

	capConfigJSON, err := json.Marshal(capConfig)
	require.NoError(t, err)

	err = capability.Initialise(ctx, core.StandardCapabilitiesDependencies{
		Config:        string(capConfigJSON),
		OracleFactory: oracleFactory,
	})
	require.NoError(t, err)

	servicetest.Run(t, capability)

	metadata := capabilities.RequestMetadata{
		WorkflowID:               "wf-id",
		WorkflowOwner:            "",
		WorkflowExecutionID:      "wex-id",
		WorkflowName:             "",
		WorkflowDonID:            0,
		WorkflowDonConfigVersion: 0,
		ReferenceID:              "1",
		DecodedWorkflowName:      "",
	}

	// Create a request that exceeds the limit
	largeData := make([]byte, 2100000)
	for i := range largeData {
		largeData[i] = 'x'
	}
	largeValue, err := values.Wrap(string(largeData))
	require.NoError(t, err)

	input := &sdk.SimpleConsensusInputs{
		Observation: &sdk.SimpleConsensusInputs_Value{Value: values.Proto(largeValue)},
		Descriptors: &sdk.ConsensusDescriptor{Descriptor_: &sdk.ConsensusDescriptor_Aggregation{Aggregation: sdk.AggregationType_AGGREGATION_TYPE_IDENTICAL}},
	}

	// User error should be sent as error input to consensus, which will result in a consensus failure
	// since all nodes will have the same oversized request
	_, err = capability.Simple(ctx, metadata, input)
	require.Error(t, err)
	var capErr caperrors.Error
	errors.As(err, &capErr)
	require.Equal(t, caperrors.ConsensusFailed, capErr.Code())
}

func Test_ToReportID(t *testing.T) {
	tests := []struct {
		name            string
		stepReferenceID string
		expectedID      string
		expectError     bool
	}{
		{
			name:            "Valid step reference ID",
			stepReferenceID: "0",
			expectedID:      "0000",
			expectError:     false,
		},
		{
			name:            "Valid step reference ID",
			stepReferenceID: "1",
			expectedID:      "0001",
			expectError:     false,
		},
		{
			name:            "Valid step reference ID",
			stepReferenceID: "26",
			expectedID:      "001a",
			expectError:     false,
		},
		{
			name:            "Valid step reference ID",
			stepReferenceID: "614",
			expectedID:      "0266",
			expectError:     false,
		},

		{
			name:            "Valid step reference ID",
			stepReferenceID: "65535", // Exceeds 2 bytes when encoded as hex
			expectedID:      "ffff",
			expectError:     false,
		},
		{
			name:            "Empty step reference ID",
			stepReferenceID: "",
			expectedID:      "",
			expectError:     true,
		},
		{
			name:            "Non-numeric step reference ID",
			stepReferenceID: "abc",
			expectedID:      "",
			expectError:     true,
		},
		{
			name:            "Step reference ID too large",
			stepReferenceID: "65536", // Exceeds 2 bytes when encoded as hex
			expectedID:      "",
			expectError:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := toReportID(tt.stepReferenceID)
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedID, result)
			}
		})
	}
}

func newRequestMetaData() capabilities.RequestMetadata {
	return capabilities.RequestMetadata{
		WorkflowID:    "0039525c34de895c8fa68006bd63f6ce4a45ef1bc66377e791c6a8ae803dc0e4",
		WorkflowOwner: "1139525c34de895c8fa68006bd634387a9f1192a",

		WorkflowExecutionID:      generateRandomHexString(32),
		WorkflowName:             "a1b2c3d4e5f6a1b2c3d4",
		WorkflowDonID:            1,
		WorkflowDonConfigVersion: 1,
		ReferenceID:              "01",
		DecodedWorkflowName:      "test-workflow-decoded",
		SpendLimits:              nil,
	}
}

func generateRandomHexString(byteLength int) string {
	randomBytes := make([]byte, byteLength)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(fmt.Sprintf("failed to generate random bytes: %v", err))
	}
	return hex.EncodeToString(randomBytes)
}
