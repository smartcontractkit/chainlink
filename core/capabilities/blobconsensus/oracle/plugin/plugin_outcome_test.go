package plugin_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/blobconsensus/oracle"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/blobconsensus/oracle/plugin"
	oracletypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/blobconsensus/oracle/types"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	libocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

func testMetaData() oracle.ConsensusRequestMetadata {
	return oracle.ConsensusRequestMetadata{
		RequestMetadata: capabilities.RequestMetadata{
			WorkflowID:               "0039525c34de895c8fa68006bd63f6ce4a45ef1bc66377e791c6a8ae803dc0e4",
			WorkflowOwner:            "1139525c34de895c8fa68006bd634387a9f1192a",
			WorkflowExecutionID:      "0102030405060708091011121314151617181920212223242526272829303132",
			WorkflowName:             "a1b2c3d4e5f6a1b2c3d4",
			WorkflowDonID:            1,
			WorkflowDonConfigVersion: 1,
			ReferenceID:              "01",
		},
		KeyBundleID: "evm",
		ReportID:    "aabb",
	}
}

const expectedMetadataString = "requestId=0102030405060708091011121314151617181920212223242526272829303132-01" +
	" workflowExecutionId=0102030405060708091011121314151617181920212223242526272829303132" +
	" workflowStepReference=01" +
	" workflowId=0039525c34de895c8fa68006bd63f6ce4a45ef1bc66377e791c6a8ae803dc0e4" +
	" workflowOwner=1139525c34de895c8fa68006bd634387a9f1192a" +
	" workflowName=a1b2c3d4e5f6a1b2c3d4" +
	" workflowDonId=1" +
	" workflowDonConfigVersion=1" +
	" reportId=aabb" +
	" keyBundleId=evm"

// makeOutcomeTestObs builds a single AttributedObservation for direct Outcome() tests.
// If isError is true the observation carries an error string; otherwise it carries an int64 value.
func makeOutcomeTestObs(
	t *testing.T,
	reqID string,
	md oracle.ConsensusRequestMetadata,
	descriptorAgg sdk.AggregationType,
	observerID uint8,
	isError bool,
	removeLibUseInFailureMessageFormattingFlag bool,
) libocrtypes.AttributedObservation {
	t.Helper()

	var simpleInputs *sdk.SimpleConsensusInputs
	if isError {
		simpleInputs = &sdk.SimpleConsensusInputs{
			Observation: &sdk.SimpleConsensusInputs_Error{
				Error: fmt.Sprintf("error from observer %d", observerID),
			},
			Descriptors: &sdk.ConsensusDescriptor{
				Descriptor_: &sdk.ConsensusDescriptor_Aggregation{Aggregation: descriptorAgg},
			},
		}
	} else {
		simpleInputs = &sdk.SimpleConsensusInputs{
			Observation: &sdk.SimpleConsensusInputs_Value{
				Value: values.Proto(values.NewInt64(int64(observerID) * 10)),
			},
			Descriptors: &sdk.ConsensusDescriptor{
				Descriptor_: &sdk.ConsensusDescriptor_Aggregation{Aggregation: descriptorAgg},
			},
		}
	}

	ro := &oracletypes.RequestObservation{
		Metadata:   plugin.ToRequestMetaData(md),
		Input:      simpleInputs,
		ReceivedAt: timestamppb.New(time.Now()),
		RemoveLibUseInFailureMessageFormattingFlag: removeLibUseInFailureMessageFormattingFlag,
	}

	obsProto := &oracletypes.Observation{
		Observations: map[string]*oracletypes.RequestObservation{reqID: ro},
	}
	b, err := proto.Marshal(obsProto)
	require.NoError(t, err)

	return libocrtypes.AttributedObservation{
		Observation: b,
		Observer:    commontypes.OracleID(observerID),
	}
}

// extractSingleFailureMessage unmarshals Outcome bytes, asserts exactly one failed
// consensus outcome is present, and returns its FailureMessage.
func extractSingleFailureMessage(t *testing.T, outcomeBytes ocr3types.Outcome) string {
	t.Helper()
	outcome := &oracletypes.Outcome{}
	require.NoError(t, proto.Unmarshal(outcomeBytes, outcome))
	require.Len(t, outcome.Outcomes, 1, "expected exactly one consensus outcome")
	failure := outcome.Outcomes[0].GetFailure()
	require.NotNil(t, failure, "expected a failed consensus outcome, got a success")
	return failure.FailureMessage
}

// Test_Outcome_PlusOneErrors checks that when every observation carries
// RemoveLibUseInFailureMessageFormatting=true and f+1 errors are received, Outcome() embeds the per-field
// metadata string ("Consensus metadata: requestId=...") and the descriptor type
// string ("Descriptor type: AGGREGATION_TYPE_MEDIAN") instead of the verbose proto dump.
func Test_Outcome_PlusOneErrors(t *testing.T) {
	t.Parallel()

	lggr := logger.Test(t)
	ctx := context.Background()

	const testF, testN = 2, 7
	reportingPlugin, _ := createReportingPlugin(t, lggr, testF, testN, 5, defaultMaxLengthBytes)

	md := testMetaData()
	reqID := md.RequestID()

	// 2f+1 = 5 observations: 3 errors (= f+1) and 2 values, all with RemoveLibUseInFailureMessageFormatting=true.
	attributed := []libocrtypes.AttributedObservation{
		makeOutcomeTestObs(t, reqID, md, sdk.AggregationType_AGGREGATION_TYPE_MEDIAN, 0, true, true),
		makeOutcomeTestObs(t, reqID, md, sdk.AggregationType_AGGREGATION_TYPE_MEDIAN, 1, true, true),
		makeOutcomeTestObs(t, reqID, md, sdk.AggregationType_AGGREGATION_TYPE_MEDIAN, 2, true, true),
		makeOutcomeTestObs(t, reqID, md, sdk.AggregationType_AGGREGATION_TYPE_MEDIAN, 3, false, true),
		makeOutcomeTestObs(t, reqID, md, sdk.AggregationType_AGGREGATION_TYPE_MEDIAN, 4, false, true),
	}

	qBytes, err := proto.Marshal(&oracletypes.Query{RequestIDs: []string{reqID}})
	require.NoError(t, err)

	outcomeBytes, err := reportingPlugin.Outcome(ctx, ocr3types.OutcomeContext{SeqNr: 1}, qBytes, attributed)
	require.NoError(t, err)

	msg := extractSingleFailureMessage(t, outcomeBytes)

	// Reduced path: individual metadata fields formatted as key=value pairs.
	assert.Contains(t, msg, "Consensus metadata: "+expectedMetadataString)
	assert.Contains(t, msg, "Descriptor type: AGGREGATION_TYPE_MEDIAN")
	// Old verbose proto-dump key must be absent.
	assert.NotContains(t, msg, "Consensus metadata, descriptor and default:")
}

// Test_Outcome_AggregationFailure checks that when aggregation itself fails
// (not enough identical values) and all observations have RemoveLibUseInFailureMessageFormatting=true, the failure
// message uses the reduced metadata format including the descriptor type string.
func Test_Outcome_AggregationFailure(t *testing.T) {
	t.Parallel()

	lggr := logger.Test(t)
	ctx := context.Background()

	const testF, testN = 2, 7
	reportingPlugin, _ := createReportingPlugin(t, lggr, testF, testN, 5, defaultMaxLengthBytes)

	md := testMetaData()
	reqID := md.RequestID()

	// Five distinct values with IDENTICAL aggregation: no value reaches the f+1=3 threshold,
	// so CalculateOutcomeForObservations returns an aggregation error.
	attributed := []libocrtypes.AttributedObservation{
		makeOutcomeTestObs(t, reqID, md, sdk.AggregationType_AGGREGATION_TYPE_IDENTICAL, 0, false, true),
		makeOutcomeTestObs(t, reqID, md, sdk.AggregationType_AGGREGATION_TYPE_IDENTICAL, 1, false, true),
		makeOutcomeTestObs(t, reqID, md, sdk.AggregationType_AGGREGATION_TYPE_IDENTICAL, 2, false, true),
		makeOutcomeTestObs(t, reqID, md, sdk.AggregationType_AGGREGATION_TYPE_IDENTICAL, 3, false, true),
		makeOutcomeTestObs(t, reqID, md, sdk.AggregationType_AGGREGATION_TYPE_IDENTICAL, 4, false, true),
	}

	qBytes, err := proto.Marshal(&oracletypes.Query{RequestIDs: []string{reqID}})
	require.NoError(t, err)

	outcomeBytes, err := reportingPlugin.Outcome(ctx, ocr3types.OutcomeContext{SeqNr: 1}, qBytes, attributed)
	require.NoError(t, err)

	msg := extractSingleFailureMessage(t, outcomeBytes)

	// Reduced path applied to the aggregation-failure branch as well.
	assert.Contains(t, msg, "Consensus metadata: "+expectedMetadataString)
	assert.Contains(t, msg, "Descriptor type: AGGREGATION_TYPE_IDENTICAL")
	assert.NotContains(t, msg, "Consensus metadata, descriptor and default:")
}

// Test_Outcome_RemoveLibUseInFailureMessageFormatting asserts the exact FailureMessage strings
// produced when every observation has RemoveLibUseInFailureMessageFormatting=true (structured
// metadata lines, descriptor type name, and bracket-formatted error lists without JSON/lib paths).
func Test_Outcome_RemoveLibUseInFailureMessageFormatting(t *testing.T) {
	t.Parallel()

	lggr := logger.Test(t)
	ctx := context.Background()

	const testF, testN = 2, 7
	reportingPlugin, _ := createReportingPlugin(t, lggr, testF, testN, 5, defaultMaxLengthBytes)

	md := testMetaData()
	reqID := md.RequestID()

	t.Run("f_plus_one_errors", func(t *testing.T) {
		attributed := []libocrtypes.AttributedObservation{
			makeOutcomeTestObs(t, reqID, md, sdk.AggregationType_AGGREGATION_TYPE_MEDIAN, 0, true, true),
			makeOutcomeTestObs(t, reqID, md, sdk.AggregationType_AGGREGATION_TYPE_MEDIAN, 1, true, true),
			makeOutcomeTestObs(t, reqID, md, sdk.AggregationType_AGGREGATION_TYPE_MEDIAN, 2, true, true),
			makeOutcomeTestObs(t, reqID, md, sdk.AggregationType_AGGREGATION_TYPE_MEDIAN, 3, false, true),
			makeOutcomeTestObs(t, reqID, md, sdk.AggregationType_AGGREGATION_TYPE_MEDIAN, 4, false, true),
		}

		qBytes, err := proto.Marshal(&oracletypes.Query{RequestIDs: []string{reqID}})
		require.NoError(t, err)

		outcomeBytes, err := reportingPlugin.Outcome(ctx, ocr3types.OutcomeContext{SeqNr: 1}, qBytes, attributed)
		require.NoError(t, err)

		msg := extractSingleFailureMessage(t, outcomeBytes)
		want := fmt.Sprintf(
			"consensus calculation failed: received 3 errors which is >= f+1 (3) for requestID %s; Consensus metadata: %s; Descriptor type: AGGREGATION_TYPE_MEDIAN; Errors received: [error from observer 0,error from observer 1,error from observer 2]",
			reqID,
			expectedMetadataString,
		)
		assert.Equal(t, want, msg)
	})

	t.Run("aggregation_failure", func(t *testing.T) {
		attributed := []libocrtypes.AttributedObservation{
			makeOutcomeTestObs(t, reqID, md, sdk.AggregationType_AGGREGATION_TYPE_IDENTICAL, 0, false, true),
			makeOutcomeTestObs(t, reqID, md, sdk.AggregationType_AGGREGATION_TYPE_IDENTICAL, 1, false, true),
			makeOutcomeTestObs(t, reqID, md, sdk.AggregationType_AGGREGATION_TYPE_IDENTICAL, 2, false, true),
			makeOutcomeTestObs(t, reqID, md, sdk.AggregationType_AGGREGATION_TYPE_IDENTICAL, 3, false, true),
			makeOutcomeTestObs(t, reqID, md, sdk.AggregationType_AGGREGATION_TYPE_IDENTICAL, 4, false, true),
		}

		qBytes, err := proto.Marshal(&oracletypes.Query{RequestIDs: []string{reqID}})
		require.NoError(t, err)

		outcomeBytes, err := reportingPlugin.Outcome(ctx, ocr3types.OutcomeContext{SeqNr: 1}, qBytes, attributed)
		require.NoError(t, err)

		msg := extractSingleFailureMessage(t, outcomeBytes)
		want := fmt.Sprintf(
			"consensus calculation failed: no values met f+1 threshold; Consensus metadata: %s; Descriptor type: AGGREGATION_TYPE_IDENTICAL; Values received: [0,10,20,30,40]; Errors received: []",
			expectedMetadataString,
		)
		assert.Equal(t, want, msg)
	})
}
