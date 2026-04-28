package plugin_test

import (
	"bytes"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	ocrtypes "github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3/types"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/requests"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/blobconsensus/metrics"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/blobconsensus/oracle"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/blobconsensus/oracle/plugin"
	oracletypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/blobconsensus/oracle/types"
)

func Test_Report_MedianTimeStamp(t *testing.T) {
	lggr := logger.Test(t)
	ctx := t.Context()

	md1 := newRequestMetaData()
	md1.KeyBundleID = "evm"

	now := time.Now()
	timestamp1 := now.Add(1 * time.Second)
	timestamp2 := now.Add(2 * time.Second)
	timestamp3 := now.Add(3 * time.Second)

	reqToObservations := map[string]*consensusPluginTest{
		md1.RequestID(): {requests: []*oracle.ConsensusRequest{
			newRRts([]byte("somerandombytes"), md1, timestamp2), newRRts([]byte("somerandombytes"), md1, timestamp3), newRRts([]byte("somerandombytes"), md1, timestamp1),
			newRRts([]byte("somerandombytes"), md1, timestamp3), newRRts([]byte("somerandombytes"), md1, timestamp1), newRRts([]byte("somerandombytes"), md1, timestamp2),
			newRRts([]byte("somerandombytes"), md1, timestamp2)},
			verifyReport: func(t *testing.T, report ocr3types.ReportPlus[[]byte], infos *structpb.Struct) {
				meta, _, err := ocrtypes.Decode(report.ReportWithInfo.Report)
				require.NoError(t, err, "Failed to extract metadata fields from report")
				assert.Equal(t, uint32(timestamp2.Unix()), meta.Timestamp) // nolint
			}},
	}

	runProtocolRoundTests(ctx, t, lggr, n, f, reqToObservations)
}

// Test_MedianTimeStampsWithMismatchedObservationsIncludesAllTimestampsInCalculation checks that the median timestamp
// calculation includes all timestamps from the observations, even if those of observations that do not match the outcome.
// This matches the OCR1 behavior where the median is calculated from all timestamps, not just those that match the outcome.
func Test_Report_MedianTimeStampWithMismatchedObservationsIncludesAllTimestampsInCalculation(t *testing.T) {
	lggr := logger.Test(t)
	ctx := t.Context()

	md1 := newRequestMetaData()
	md1.KeyBundleID = "evm"

	now := time.Now()
	timestamp1 := now.Add(1 * time.Second)
	timestamp2 := now.Add(2 * time.Second)
	timestamp3 := now.Add(3 * time.Second)

	reqToObservations := map[string]*consensusPluginTest{
		md1.RequestID(): {requests: []*oracle.ConsensusRequest{
			newRRts([]byte("somerandombytes"), md1, timestamp1), newRRts([]byte("somerandombytes"), md1, timestamp1), newRRts([]byte("somerandombytes"), md1, timestamp1),
			newRRts([]byte("somerandombytes2"), md1, timestamp2), newRRts([]byte("somerandombytes2"), md1, timestamp2), newRRts([]byte("somerandombytes"), md1, timestamp3),
			newRRts([]byte("somerandombytes"), md1, timestamp3)},
			verifyReport: func(t *testing.T, report ocr3types.ReportPlus[[]byte], infos *structpb.Struct) {
				meta, _, err := ocrtypes.Decode(report.ReportWithInfo.Report)
				require.NoError(t, err, "Failed to extract metadata fields from report")
				assert.Equal(t, uint32(timestamp2.Unix()), meta.Timestamp) // nolint
			}},
	}

	runProtocolRoundTests(ctx, t, lggr, n, f, reqToObservations)
}

func Test_ReceivedIdenticalReportFromAllNodes(t *testing.T) {
	lggr := logger.Test(t)
	ctx := t.Context()

	md1 := newRequestMetaData()
	md1.KeyBundleID = "evm"

	reqToObservations := map[string]*consensusPluginTest{
		md1.RequestID(): {requests: []*oracle.ConsensusRequest{
			newRR([]byte("somerandombytes"), md1), newRR([]byte("somerandombytes"), md1), newRR([]byte("somerandombytes"), md1),
			newRR([]byte("somerandombytes"), md1), newRR([]byte("somerandombytes"), md1), newRR([]byte("somerandombytes"), md1),
			newRR([]byte("somerandombytes"), md1)},
			verifyReport: func(t *testing.T, report ocr3types.ReportPlus[[]byte], infos *structpb.Struct) {
				require.True(t, bytes.HasSuffix(report.ReportWithInfo.Report, []byte("somerandombytes")), "Report does not end with 'somerandombytes'")
				require.Equal(t, md1.RequestID(), infos.AsMap()[plugin.InfoRequestID], "RequestID does not match md1.RequestID")
				require.Equal(t, md1.KeyBundleID, infos.AsMap()["keyBundleName"], "KeyBundleID does not match md1.KeyBundleID")

				meta, _, err := ocrtypes.Decode(report.ReportWithInfo.Report)
				require.NoError(t, err, "Failed to extract metadata fields from report")

				require.Equal(t, md1.WorkflowExecutionID, meta.ExecutionID, "Metadata ExecutionID does not match")
				require.Equal(t, md1.WorkflowDonID, meta.DONID, "Metadata DONID does not match")
				require.Equal(t, md1.WorkflowDonConfigVersion, meta.DONConfigVersion, "Metadata DONConfigVersion does not match")
				require.Equal(t, md1.WorkflowID, meta.WorkflowID, "Metadata WorkflowID does not match")
				require.Equal(t, md1.WorkflowName, meta.WorkflowName, "Metadata WorkflowName does not match")
				require.Equal(t, md1.WorkflowOwner, meta.WorkflowOwner, "Metadata WorkflowOwner does not match")
				require.Equal(t, md1.ReportID, meta.ReportID, "Metadata WorkflowOwner does not match")
			}},
	}

	runProtocolRoundTests(ctx, t, lggr, n, f, reqToObservations)
}

func Test_ReceivedIdenticalReportFromSufficientNodes(t *testing.T) {
	lggr := logger.Test(t)
	ctx := t.Context()

	md1 := newRequestMetaData()
	md1.KeyBundleID = "evm"

	reqToObservations := map[string]*consensusPluginTest{
		md1.RequestID(): {requests: []*oracle.ConsensusRequest{
			newRR([]byte("somerandombytes2"), md1), newRR([]byte("somerandombytes"), md1), newRR([]byte("somerandombytes"), md1),
			newRR([]byte("somerandombytes"), md1), newRR([]byte("somerandombytes"), md1), newRR([]byte("somerandombytes"), md1),
			newRR([]byte("somerandombytes2"), md1)},
			verifyReport: func(t *testing.T, report ocr3types.ReportPlus[[]byte], infos *structpb.Struct) {
				require.True(t, bytes.HasSuffix(report.ReportWithInfo.Report, []byte("somerandombytes")), "Report does not end with 'somerandombytes'")
				require.Equal(t, md1.RequestID(), infos.AsMap()[plugin.InfoRequestID], "RequestID does not match md1.RequestID")
				require.Equal(t, md1.KeyBundleID, infos.AsMap()["keyBundleName"], "KeyBundleID does not match md1.KeyBundleID")

				meta, _, err := ocrtypes.Decode(report.ReportWithInfo.Report)
				require.NoError(t, err, "Failed to extract metadata fields from report")

				require.Equal(t, md1.WorkflowExecutionID, meta.ExecutionID, "Metadata ExecutionID does not match")
				require.Equal(t, md1.WorkflowDonID, meta.DONID, "Metadata DONID does not match")
				require.Equal(t, md1.WorkflowDonConfigVersion, meta.DONConfigVersion, "Metadata DONConfigVersion does not match")
				require.Equal(t, md1.WorkflowID, meta.WorkflowID, "Metadata WorkflowID does not match")
				require.Equal(t, md1.WorkflowName, meta.WorkflowName, "Metadata WorkflowName does not match")
				require.Equal(t, md1.WorkflowOwner, meta.WorkflowOwner, "Metadata WorkflowOwner does not match")
				require.Equal(t, md1.ReportID, meta.ReportID, "Metadata WorkflowOwner does not match")
			}},
	}

	runProtocolRoundTests(ctx, t, lggr, n, f, reqToObservations)
}

func Test_SufficientAndInsufficentReportsInSingleRound(t *testing.T) {
	lggr := logger.Test(t)
	ctx := t.Context()

	md1 := newRequestMetaData()
	md1.KeyBundleID = "evm"

	md2 := newRequestMetaData()
	md2.KeyBundleID = "evm"

	expectedFailureCode := oracletypes.ConsensusFailureCode_CONSENSUS_CALCULATION_FAILED

	reqToObservations := map[string]*consensusPluginTest{
		md1.RequestID(): {requests: []*oracle.ConsensusRequest{
			newRR([]byte("somerandombytes"), md1), newRR([]byte("somerandombytes"), md1), newRR([]byte("somerandombytes"), md1),
			newRR([]byte("somerandombytes"), md1), newRR([]byte("somerandombytes"), md1), newRR([]byte("somerandombytes"), md1),
			newRR([]byte("somerandombytes"), md1)},
			verifyReport: func(t *testing.T, report ocr3types.ReportPlus[[]byte], infos *structpb.Struct) {
				require.True(t, bytes.HasSuffix(report.ReportWithInfo.Report, []byte("somerandombytes")), "Report does not end with 'somerandombytes'")
				require.Equal(t, md1.RequestID(), infos.AsMap()[plugin.InfoRequestID], "RequestID does not match md1.RequestID")
			}},

		md2.RequestID(): {requests: []*oracle.ConsensusRequest{
			newRR([]byte("somerandombytes2"), md2), newRR([]byte("somerandombytes4"), md2), newRR([]byte("somerandombytes3"), md2),
			newRR([]byte("somerandombytes"), md2), newRR([]byte("somerandombytes"), md2), newRR([]byte("somerandombytes3"), md2),
			newRR([]byte("somerandombytes2"), md2)},
			expectedConsensusFailureMessage: "no values met f+1 threshold",
			expectedConsensusFailureCode:    &expectedFailureCode},
	}

	runProtocolRoundTests(ctx, t, lggr, n, f, reqToObservations)
}

func Test_ReceivedIdenticalMultipleQualifyingSetsOfIdenticalValues(t *testing.T) {
	lggr := logger.Test(t)
	ctx := t.Context()

	md1 := newRequestMetaData()
	md1.KeyBundleID = "evm"

	expectedFailureCode := oracletypes.ConsensusFailureCode_MORE_THAN_ONE_VALID_OUTCOME_FOR_IDENTICAL_CONSENSUS

	reqToObservations := map[string]*consensusPluginTest{
		md1.RequestID(): {requests: []*oracle.ConsensusRequest{
			newRR([]byte("somerandombytes2"), md1), newRR([]byte("somerandombytes"), md1), newRR([]byte("somerandombytes"), md1),
			newRR([]byte("somerandombytes2"), md1), newRR([]byte("somerandombytes"), md1), newRR([]byte("somerandombytes"), md1),
			newRR([]byte("somerandombytes2"), md1)},
			expectedConsensusFailureMessage: "not identical, multiple values with f+1 occurrences",
			expectedConsensusFailureCode:    &expectedFailureCode,
		},
	}

	runProtocolRoundTests(ctx, t, lggr, n, f, reqToObservations)
}

func newRR(rawBytes []byte, metaData oracle.ConsensusRequestMetadata) *oracle.ConsensusRequest {
	return newRRts(rawBytes, metaData, time.Now())
}

func newRRts(rawBytes []byte, metaData oracle.ConsensusRequestMetadata, recievedAt time.Time) *oracle.ConsensusRequest {
	simpleConsensusInputs := &sdk.SimpleConsensusInputs{
		Observation: &sdk.SimpleConsensusInputs_Value{Value: values.Proto(values.NewBytes(rawBytes))},
		Descriptors: &sdk.ConsensusDescriptor{Descriptor_: &sdk.ConsensusDescriptor_Aggregation{Aggregation: sdk.AggregationType_AGGREGATION_TYPE_IDENTICAL}},
	}

	return oracle.NewConsensusRequest(simpleConsensusInputs, recievedAt, time.Now().Add(1*time.Hour).UTC(), nil, metaData)
}

func Test_ReportTooLarge_ReturnsFailure(t *testing.T) {
	lggr := logger.Test(t)
	ctx := t.Context()

	// Create a reporting plugin with a very small max report length
	reqStore := requests.NewStore[*oracle.ConsensusRequest]()
	metricsInstance, err := metrics.NewMetrics()
	require.NoError(t, err)

	smallMaxReportLength := uint32(200) // Very small to trigger the error
	reportingPlugin, err := plugin.NewReportingPlugin(lggr, metricsInstance, f, n, reqStore, &ocrtypes.ReportingPluginConfig{
		MaxQueryLengthBytes:              1000000,
		MaxObservationLengthBytes:        1000000,
		MaxOutcomeLengthBytes:            1000000,
		MaxReportLengthBytes:             smallMaxReportLength,
		HistoricalOutcomeExpirySeqNrSpan: 5,
	}, "evm", 1000)
	require.NoError(t, err)

	// Create an outcome with a large report that exceeds the limit
	largeData := strings.Repeat("x", 210) // More than 200
	serialisedValue, err := proto.MarshalOptions{Deterministic: true}.Marshal(values.Proto(values.NewBytes([]byte(largeData))))
	require.NoError(t, err)

	outcome := &oracletypes.Outcome{
		Outcomes: []*oracletypes.ConsensusOutcome{
			{
				Outcome: &oracletypes.ConsensusOutcome_Success{
					Success: &oracletypes.ConsensusSuccessOutcome{
						Metadata: &oracletypes.RequestMetaData{
							RequestId:                "req-too-large",
							WorkflowExecutionId:      "0102030405060708091011121314151617181920212223242526272829303132",
							WorkflowId:               "0039525c34de895c8fa68006bd63f6ce4a45ef1bc66377e791c6a8ae803dc0e4",
							WorkflowOwner:            "1139525c34de895c8fa68006bd634387a9f1192a",
							WorkflowName:             "a1b2c3d4e5f6a1b2c3d4",
							WorkflowDonId:            1,
							WorkflowDonConfigVersion: 1,
							ReportId:                 "abcd",
							KeyBundleId:              "evm",
							RequestType:              oracletypes.RequestType_REPORT_GENERATION,
						},
						Outcome:   serialisedValue,
						Timestamp: timestamppb.Now(),
					},
				},
			},
		},
	}

	serialisedOutcome, err := proto.MarshalOptions{Deterministic: true}.Marshal(outcome)
	require.NoError(t, err)

	// Call Reports and verify we get a failure report instead of an error
	reports, err := reportingPlugin.Reports(ctx, 1, ocr3_1types.ReportsPlusPrecursor(serialisedOutcome))
	require.NoError(t, err, "Reports should not return an error for oversized reports")
	require.Len(t, reports, 1, "Should have one report")

	// Verify the report is a failure report (empty report body)
	require.Empty(t, reports[0].ReportWithInfo.Report, "Oversized report should have empty report body")

	// Parse the info to verify failure code
	infos := &structpb.Struct{}
	err = proto.Unmarshal(reports[0].ReportWithInfo.Info, infos)
	require.NoError(t, err)

	infoMap := infos.AsMap()
	require.Equal(t, "req-too-large", infoMap[plugin.InfoRequestID])
	require.Equal(t, oracletypes.ConsensusFailureCode_REPORT_TOO_LARGE.String(), infoMap[plugin.InfoConsensusFailureCode])
	require.Contains(t, infoMap[plugin.InfoConsensusFailureMessage], "report too large")
}

func Test_ReportWithinLimit_Succeeds(t *testing.T) {
	lggr := logger.Test(t)
	ctx := t.Context()

	// Create a reporting plugin with a reasonable max report length
	reqStore := requests.NewStore[*oracle.ConsensusRequest]()
	metricsInstance, err := metrics.NewMetrics()
	require.NoError(t, err)

	reportingPlugin, err := plugin.NewReportingPlugin(lggr, metricsInstance, f, n, reqStore, &ocrtypes.ReportingPluginConfig{
		MaxQueryLengthBytes:              1000000,
		MaxObservationLengthBytes:        1000000,
		MaxOutcomeLengthBytes:            1000000,
		MaxReportLengthBytes:             10000, // Larger limit
		HistoricalOutcomeExpirySeqNrSpan: 5,
	}, "evm", 1000)
	require.NoError(t, err)

	// Create an outcome with a small report that fits
	smallData := "small"
	serialisedValue, err := proto.MarshalOptions{Deterministic: true}.Marshal(values.Proto(values.NewBytes([]byte(smallData))))
	require.NoError(t, err)

	outcome := &oracletypes.Outcome{
		Outcomes: []*oracletypes.ConsensusOutcome{
			{
				Outcome: &oracletypes.ConsensusOutcome_Success{
					Success: &oracletypes.ConsensusSuccessOutcome{
						Metadata: &oracletypes.RequestMetaData{
							RequestId:                "req-small",
							WorkflowExecutionId:      "0102030405060708091011121314151617181920212223242526272829303132",
							WorkflowId:               "0039525c34de895c8fa68006bd63f6ce4a45ef1bc66377e791c6a8ae803dc0e4",
							WorkflowOwner:            "1139525c34de895c8fa68006bd634387a9f1192a",
							WorkflowName:             "a1b2c3d4e5f6a1b2c3d4",
							WorkflowDonId:            1,
							WorkflowDonConfigVersion: 1,
							ReportId:                 "abcd",
							KeyBundleId:              "evm",
							RequestType:              oracletypes.RequestType_REPORT_GENERATION,
						},
						Outcome:   serialisedValue,
						Timestamp: timestamppb.Now(),
					},
				},
			},
		},
	}

	serialisedOutcome, err := proto.MarshalOptions{Deterministic: true}.Marshal(outcome)
	require.NoError(t, err)

	// Call Reports and verify we get a success report
	reports, err := reportingPlugin.Reports(ctx, 1, ocr3_1types.ReportsPlusPrecursor(serialisedOutcome))
	require.NoError(t, err)
	require.Len(t, reports, 1)

	// Verify the report is NOT empty (successful report)
	require.NotEmpty(t, reports[0].ReportWithInfo.Report, "Successful report should have non-empty report body")

	// Parse the info to verify no failure code
	infos := &structpb.Struct{}
	err = proto.Unmarshal(reports[0].ReportWithInfo.Info, infos)
	require.NoError(t, err)

	infoMap := infos.AsMap()
	require.Equal(t, "req-small", infoMap[plugin.InfoRequestID])
	require.Nil(t, infoMap[plugin.InfoConsensusFailureCode], "Should not have failure code for successful report")
}

func Test_ReportCountLimit(t *testing.T) {
	lggr := logger.Test(t)
	ctx := t.Context()

	// Create a reporting plugin with a reasonable max report length
	reqStore := requests.NewStore[*oracle.ConsensusRequest]()
	metricsInstance, err := metrics.NewMetrics()
	require.NoError(t, err)

	maxReportCount := uint32(100)
	reportingPlugin, err := plugin.NewReportingPlugin(lggr, metricsInstance, f, n, reqStore, &ocrtypes.ReportingPluginConfig{
		MaxQueryLengthBytes:              1000000,
		MaxObservationLengthBytes:        1000000,
		MaxOutcomeLengthBytes:            1000000,
		MaxReportLengthBytes:             10000, // Larger limit
		MaxReportCount:                   maxReportCount,
		HistoricalOutcomeExpirySeqNrSpan: 5,
	}, "evm", 1000)
	require.NoError(t, err)

	// Create an outcome with a small report that fits
	smallData := "small"
	serialisedValue, err := proto.MarshalOptions{Deterministic: true}.Marshal(values.Proto(values.NewBytes([]byte(smallData))))
	require.NoError(t, err)

	outcome := &oracletypes.Outcome{}

	for i := 0; i < int(maxReportCount)+5; i++ {
		outcome.Outcomes = append(outcome.Outcomes, &oracletypes.ConsensusOutcome{
			Outcome: &oracletypes.ConsensusOutcome_Success{
				Success: &oracletypes.ConsensusSuccessOutcome{
					Metadata: &oracletypes.RequestMetaData{
						RequestId:                "req-small-" + strconv.Itoa(i),
						WorkflowExecutionId:      "0102030405060708091011121314151617181920212223242526272829303132",
						WorkflowId:               "0039525c34de895c8fa68006bd63f6ce4a45ef1bc66377e791c6a8ae803dc0e4",
						WorkflowOwner:            "1139525c34de895c8fa68006bd634387a9f1192a",
						WorkflowName:             "a1b2c3d4e5f6a1b2c3d4",
						WorkflowDonId:            1,
						WorkflowDonConfigVersion: 1,
						ReportId:                 "abcd",
						KeyBundleId:              "evm",
						RequestType:              oracletypes.RequestType_REPORT_GENERATION,
					},
					Outcome:   serialisedValue,
					Timestamp: timestamppb.Now(),
				},
			},
		})
	}

	serialisedOutcome, err := proto.MarshalOptions{Deterministic: true}.Marshal(outcome)
	require.NoError(t, err)

	// Call Reports and verify the number of reports is limited to the max count
	reports, err := reportingPlugin.Reports(ctx, 1, ocr3_1types.ReportsPlusPrecursor(serialisedOutcome))
	require.NoError(t, err)
	require.Len(t, reports, int(maxReportCount))
}
