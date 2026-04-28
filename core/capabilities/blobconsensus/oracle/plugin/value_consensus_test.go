package plugin_test

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	pbtypes "github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3/types"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/requests"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/blobconsensus/metrics"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/blobconsensus/oracle"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/blobconsensus/oracle/plugin"
	oracletypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/blobconsensus/oracle/types"

	"github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"
	valuespb "github.com/smartcontractkit/chainlink-protos/cre/go/values/pb"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	libocrTypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

type consensusPluginTest struct {
	requests                        []*oracle.ConsensusRequest
	verifyReport                    func(t *testing.T, report ocr3types.ReportPlus[[]byte], infos *structpb.Struct)
	expectedConsensusFailureMessage string
	expectedConsensusFailureCode    *oracletypes.ConsensusFailureCode
}

const (
	n                     = 7
	f                     = 2
	defaultMaxLengthBytes = 1000000 // 1 MB
	defaultMaxReportCount = 10000
)

// nillable observation and nillable default value, -1 indicates the value should be set as nil
func newSliceCr(t *testing.T, observation []byte, def []byte, metaData oracle.ConsensusRequestMetadata) *oracle.ConsensusRequest {
	observationVal, err := values.Wrap(observation)
	require.NoError(t, err, "failed to wrap nil value")

	defaultVal, err := values.Wrap(def)
	require.NoError(t, err, "failed to wrap nil value")

	simpleConsensusInputs := &sdk.SimpleConsensusInputs{
		Observation: &sdk.SimpleConsensusInputs_Value{Value: values.Proto(observationVal)},
		Default:     values.Proto(defaultVal),
		Descriptors: &sdk.ConsensusDescriptor{Descriptor_: &sdk.ConsensusDescriptor_Aggregation{Aggregation: sdk.AggregationType_AGGREGATION_TYPE_IDENTICAL}},
	}

	return oracle.NewConsensusRequest(serializeDeserialize(t, simpleConsensusInputs), time.Now(), time.Now().Add(1*time.Hour).UTC(), nil, metaData)
}

func Test_InsufficientIdenticalObservations(t *testing.T) {
	lggr := logger.Test(t)
	ctx := t.Context()

	md1 := newRequestMetaData()

	expectedFailureCode := oracletypes.ConsensusFailureCode_CONSENSUS_CALCULATION_FAILED

	reqToObservations := map[string]*consensusPluginTest{
		md1.RequestID(): {
			requests: []*oracle.ConsensusRequest{
				newIdenticalCr(t, 110, md1), newIdenticalCr(t, 110, md1),
				newIdenticalCr(t, 120, md1), newIdenticalCr(t, 120, md1),
				newIdenticalCr(t, 130, md1), newIdenticalCr(t, 130, md1),
				newIdenticalCr(t, 140, md1), newIdenticalCr(t, 140, md1),
			},
			expectedConsensusFailureMessage: "no values met f+1 threshold",
			expectedConsensusFailureCode:    &expectedFailureCode,
		},
	}

	runProtocolRoundTests(ctx, t, lggr, n, f, reqToObservations)
}

func Test_InsufficientIdenticalMapObservations(t *testing.T) {
	lggr := logger.Test(t)
	ctx := t.Context()

	md1 := newRequestMetaData()

	type testStruct struct {
		Field1 int
	}

	expectedFailureCode := oracletypes.ConsensusFailureCode_CONSENSUS_CALCULATION_FAILED

	reqToObservations := map[string]*consensusPluginTest{
		md1.RequestID(): {
			requests: []*oracle.ConsensusRequest{
				newIdenticalValueCr(t, mustWrap(t, testStruct{Field1: 100}), md1),
				newIdenticalValueCr(t, mustWrap(t, testStruct{Field1: 110}), md1),
				newIdenticalValueCr(t, mustWrap(t, testStruct{Field1: 120}), md1),
				newIdenticalValueCr(t, mustWrap(t, testStruct{Field1: 130}), md1),
				newIdenticalValueCr(t, mustWrap(t, testStruct{Field1: 140}), md1),
				newIdenticalValueCr(t, mustWrap(t, testStruct{Field1: 150}), md1),
				newIdenticalValueCr(t, mustWrap(t, testStruct{Field1: 160}), md1),
			},
			expectedConsensusFailureMessage: "no values met f+1 threshold",
			expectedConsensusFailureCode:    &expectedFailureCode,
		},
	}

	runProtocolRoundTests(ctx, t, lggr, n, f, reqToObservations)
}

func mustWrap(t *testing.T, v any) values.Value {
	val, err := values.Wrap(v)
	require.NoError(t, err)
	return val
}

func Test_SliceObservationAndDefaults(t *testing.T) {
	lggr := logger.Test(t)
	ctx := t.Context()

	md1 := newRequestMetaData()
	md2 := newRequestMetaData()

	reqToObservations := map[string]*consensusPluginTest{
		// Test with observations and defaults as byte slices
		md1.RequestID(): {
			requests: []*oracle.ConsensusRequest{
				newSliceCr(t, []byte("stuff"), []byte("otherstuff"), md1),
				newSliceCr(t, []byte("stuff"), []byte("otherstuff"), md1),
				newSliceCr(t, []byte("stuff"), []byte("otherstuff"), md1),
				newSliceCr(t, []byte("stuff"), []byte("otherstuff"), md1),
				newSliceCr(t, []byte("stuff"), []byte("otherstuff"), md1),
				newSliceCr(t, []byte("stuff"), []byte("otherstuff"), md1),
				newSliceCr(t, []byte("stuff"), []byte("otherstuff"), md1),
			},
			verifyReport: func(t *testing.T, report ocr3types.ReportPlus[[]byte], infos *structpb.Struct) {
				val, err := values.Wrap([]byte("stuff"))
				require.NoError(t, err)

				verifyValueConsensusReport(t, report, infos, val, "")
			},
		},

		// Test with just defaults as byte slices
		md2.RequestID(): {
			requests: []*oracle.ConsensusRequest{
				newSliceCr(t, []byte{}, []byte("otherstuff"), md2),
				newSliceCr(t, []byte{}, []byte("otherstuff"), md2),
				newSliceCr(t, []byte{}, []byte("otherstuff"), md2),
				newSliceCr(t, []byte{}, []byte("otherstuff"), md2),
				newSliceCr(t, []byte{}, []byte("otherstuff"), md2),
				newSliceCr(t, []byte{}, []byte("otherstuff"), md2),
				newSliceCr(t, []byte{}, []byte("otherstuff"), md2),
			},
			verifyReport: func(t *testing.T, report ocr3types.ReportPlus[[]byte], infos *structpb.Struct) {
				val, err := values.Wrap([]byte{})
				require.NoError(t, err)

				verifyValueConsensusReport(t, report, infos, val, "")
			},
		},

		// Test with a mixture of observations and defaults as byte slices
		md2.RequestID(): {
			requests: []*oracle.ConsensusRequest{
				newSliceCr(t, []byte("guff"), []byte("otherstuff"), md2),
				newSliceCr(t, []byte("somestuff"), []byte("otherstuff"), md2),
				newSliceCr(t, []byte("stuff"), []byte("otherstuff"), md2),
				newSliceCr(t, nil, []byte("otherstuff"), md2),
				newSliceCr(t, nil, []byte("otherstuff"), md2),
				newSliceCr(t, []byte("stuff"), []byte("otherstuff"), md2),
				newSliceCr(t, []byte("somestuff"), []byte("otherstuff"), md2),
			},
			verifyReport: func(t *testing.T, report ocr3types.ReportPlus[[]byte], infos *structpb.Struct) {
				val, err := values.Wrap([]byte("otherstuff"))
				require.NoError(t, err)

				verifyValueConsensusReport(t, report, infos, val, "")
			},
		},
	}

	runProtocolRoundTests(ctx, t, lggr, n, f, reqToObservations)
}

func Test_MismatchedLeaderConsensusDescriptor(t *testing.T) {
	lggr := logger.Test(t)
	ctx := t.Context()

	metaData := newRequestMetaData()

	protocolRoundTests := map[string]*consensusPluginTest{
		metaData.RequestID(): {
			requests: []*oracle.ConsensusRequest{
				newIdenticalCr(t, 110, metaData), newCr(t, 120, metaData), newCr(t, 130, metaData),
				newCr(t, 140, metaData), newCr(t, 150, metaData), newCr(t, 160, metaData),
				newCr(t, 170, metaData),
			},
			verifyReport: func(t *testing.T, report ocr3types.ReportPlus[[]byte], infos *structpb.Struct) {
				verifyValueConsensusReport(t, report, infos, values.NewInt64(140), "")
			},
		},
	}

	runProtocolRoundTests(ctx, t, lggr, n, f, protocolRoundTests)
}

func Test_MismatchedNonLeaderConsensusDescriptor(t *testing.T) {
	lggr := logger.Test(t)
	ctx := t.Context()

	metaData := newRequestMetaData()

	newCrIdenticalConsensus := func(observation int64, metaData oracle.ConsensusRequestMetadata) *oracle.ConsensusRequest {
		simpleConsensusInputs := &sdk.SimpleConsensusInputs{
			Observation: &sdk.SimpleConsensusInputs_Value{Value: values.Proto(values.NewInt64(observation))},
			Descriptors: &sdk.ConsensusDescriptor{Descriptor_: &sdk.ConsensusDescriptor_Aggregation{Aggregation: sdk.AggregationType_AGGREGATION_TYPE_IDENTICAL}},
		}

		return oracle.NewConsensusRequest(simpleConsensusInputs, time.Now().Add(1*time.Hour).UTC(), time.Now(), nil, metaData)
	}

	protocolRoundTests := map[string]*consensusPluginTest{
		metaData.RequestID(): {
			requests: []*oracle.ConsensusRequest{
				newCr(t, 110, metaData), newCr(t, 120, metaData), newCr(t, 130, metaData),
				newCr(t, 140, metaData), newCrIdenticalConsensus(150, metaData), newCr(t, 160, metaData),
				newCr(t, 170, metaData),
			},
			verifyReport: func(t *testing.T, report ocr3types.ReportPlus[[]byte], infos *structpb.Struct) {
				verifyValueConsensusReport(t, report, infos, values.NewInt64(130), "")
			},
		},
	}

	runProtocolRoundTests(ctx, t, lggr, n, f, protocolRoundTests)
}

func Test_MismatchedLeaderMetaData(t *testing.T) {
	lggr := logger.Test(t)
	ctx := t.Context()

	metaData := newRequestMetaData()
	leaderMetaData := metaData

	leaderMetaData.WorkflowDonID = 2

	protocolRoundTests := map[string]*consensusPluginTest{
		metaData.RequestID(): {
			requests: []*oracle.ConsensusRequest{
				newCr(t, 110, leaderMetaData), newCr(t, 120, metaData), newCr(t, 130, metaData),
				newCr(t, 140, metaData), newCr(t, 150, metaData), newCr(t, 160, metaData),
				newCr(t, 170, metaData),
			},
			verifyReport: func(t *testing.T, report ocr3types.ReportPlus[[]byte], infos *structpb.Struct) {
				verifyValueConsensusReport(t, report, infos, values.NewInt64(140), "")
			},
		},
	}

	runProtocolRoundTests(ctx, t, lggr, n, f, protocolRoundTests)
}

func Test_MismatchedNonLeaderMetaData(t *testing.T) {
	lggr := logger.Test(t)
	ctx := t.Context()

	metaData := newRequestMetaData()
	misMatchedMetaData := metaData

	misMatchedMetaData.WorkflowDonID = 2

	protocolRoundTests := map[string]*consensusPluginTest{
		metaData.RequestID(): {
			requests: []*oracle.ConsensusRequest{
				newCr(t, 110, metaData), newCr(t, 120, metaData), newCr(t, 130, metaData),
				newCr(t, 140, metaData), newCr(t, 150, misMatchedMetaData), newCr(t, 160, metaData),
				newCr(t, 170, metaData),
			},
			verifyReport: func(t *testing.T, report ocr3types.ReportPlus[[]byte], infos *structpb.Struct) {
				verifyValueConsensusReport(t, report, infos, values.NewInt64(130), "")
			},
		},
	}

	runProtocolRoundTests(ctx, t, lggr, n, f, protocolRoundTests)
}

func Test_ObservationDefaults(t *testing.T) {
	lggr := logger.Test(t)
	ctx := t.Context()

	md1 := newRequestMetaData()
	md2 := newRequestMetaData()
	md3 := newRequestMetaData()

	md1.KeyBundleID = "evm"

	reqToObservations := map[string]*consensusPluginTest{
		// Test a mixture of nil and non-nil observations with defaults
		md1.RequestID(): {
			requests: []*oracle.ConsensusRequest{
				newNillableCr(t, -1, 40, md1), newNillableCr(t, 20, 40, md1), newNillableCr(t, -1, 40, md1),
				newNillableCr(t, -1, 40, md1), newNillableCr(t, -1, 40, md1), newNillableCr(t, -1, 40, md1),
				newNillableCr(t, 70, 40, md1),
			},
			verifyReport: func(t *testing.T, report ocr3types.ReportPlus[[]byte], infos *structpb.Struct) {
				verifyValueConsensusReport(t, report, infos, values.NewInt64(40), "evm")
			},
		},

		// Test obs and default nil observations
		md2.RequestID(): {
			requests: []*oracle.ConsensusRequest{
				newNillableCr(t, 110, 100, md2), newNillableCr(t, 120, 100, md2), newNillableCr(t, -1, -1, md2),
				newNillableCr(t, -1, 100, md2), newNillableCr(t, 150, -1, md2), newNillableCr(t, 160, 100, md2),
				newNillableCr(t, 170, 100, md2),
			},
			verifyReport: func(t *testing.T, report ocr3types.ReportPlus[[]byte], infos *structpb.Struct) {
				verifyValueConsensusReport(t, report, infos, values.NewInt64(120), "")
			},
		},

		// Test insufficient non-nil observations but with sufficient matching defaults
		md3.RequestID(): {
			requests: []*oracle.ConsensusRequest{
				newNillableCr(t, 10, 40, md3), newNillableCr(t, -1, 40, md3), newNillableCr(t, 30, 40, md3),
				newNillableCr(t, -1, 40, md3), newNillableCr(t, -1, 40, md3), newNillableCr(t, -1, 40, md3),
				newNillableCr(t, -1, 40, md3),
			},
			verifyReport: func(t *testing.T, report ocr3types.ReportPlus[[]byte], infos *structpb.Struct) {
				verifyValueConsensusReport(t, report, infos, values.NewInt64(40), "")
			},
		},
	}

	runProtocolRoundTests(ctx, t, lggr, n, f, reqToObservations)
}

func Test_ReceivedAllObservationsFromAllNodes(t *testing.T) {
	lggr := logger.Test(t)
	ctx := t.Context()

	md1 := newRequestMetaData()
	md2 := newRequestMetaData()

	md1.KeyBundleID = "evm"

	reqToObservations := map[string]*consensusPluginTest{
		md1.RequestID(): {
			requests: []*oracle.ConsensusRequest{
				newCr(t, 10, md1), newCr(t, 20, md1), newCr(t, 30, md1),
				newCr(t, 40, md1), newCr(t, 50, md1), newCr(t, 60, md1),
				newCr(t, 70, md1),
			},
			verifyReport: func(t *testing.T, report ocr3types.ReportPlus[[]byte], infos *structpb.Struct) {
				verifyValueConsensusReport(t, report, infos, values.NewInt64(40), "evm")
			},
		},

		md2.RequestID(): {
			requests: []*oracle.ConsensusRequest{
				newCr(t, 110, md2), newCr(t, 120, md2), newCr(t, 130, md2),
				newCr(t, 140, md2), newCr(t, 150, md2), newCr(t, 160, md2),
				newCr(t, 170, md2),
			},
			verifyReport: func(t *testing.T, report ocr3types.ReportPlus[[]byte], infos *structpb.Struct) {
				verifyValueConsensusReport(t, report, infos, values.NewInt64(140), "")
			},
		},
	}

	runProtocolRoundTests(ctx, t, lggr, n, f, reqToObservations)
}

func Test_ReceivedObservationsWithMatchingDefaults(t *testing.T) {
	lggr := logger.Test(t)
	ctx := t.Context()

	md1 := newRequestMetaData()
	md1.KeyBundleID = "evm"

	reqToObservations := map[string]*consensusPluginTest{
		md1.RequestID(): {
			requests: []*oracle.ConsensusRequest{
				newCrWithObsAndDef(t, 10, 17, md1), newCrWithObsAndDef(t, 20, 17, md1), newCrWithObsAndDef(t, 30, 17, md1),
				newCrWithObsAndDef(t, 40, 17, md1), newCrWithObsAndDef(t, 50, 17, md1), newCrWithObsAndDef(t, 60, 17, md1),
				newCrWithObsAndDef(t, 70, 17, md1),
			},
			verifyReport: func(t *testing.T, report ocr3types.ReportPlus[[]byte], infos *structpb.Struct) {
				verifyValueConsensusReport(t, report, infos, values.NewInt64(40), "evm")
			},
		},
	}

	runProtocolRoundTests(ctx, t, lggr, n, f, reqToObservations)
}

// In this test some nodes have observations that match the default, and some have observations that do not match the default
// The consensus should be reached as there are sufficient observations with defaults that match the leader's default
func Test_ReceivedObservationsWithSomeMisMatchedDefaults_SufficientForConsensus(t *testing.T) {
	lggr := logger.Test(t)
	ctx := t.Context()

	md1 := newRequestMetaData()
	md1.KeyBundleID = "evm"

	reqToObservations := map[string]*consensusPluginTest{
		md1.RequestID(): {
			requests: []*oracle.ConsensusRequest{
				newCrWithObsAndDef(t, 10, 17, md1), newCrWithObsAndDef(t, 20, 17, md1), newCrWithObsAndDef(t, 30, 17, md1),
				newCrWithObsAndDef(t, 40, 16, md1), newCrWithObsAndDef(t, 50, 17, md1), newCrWithObsAndDef(t, 60, 17, md1),
				newCrWithObsAndDef(t, 70, 17, md1),
			},
			verifyReport: func(t *testing.T, report ocr3types.ReportPlus[[]byte], infos *structpb.Struct) {
				verifyValueConsensusReport(t, report, infos, values.NewInt64(30), "evm")
			},
		},
	}

	runProtocolRoundTests(ctx, t, lggr, n, f, reqToObservations)
}

// In this test some nodes have observations that match the default, and some have observations that do not match the default
// The consensus should not be reached as there are insufficient observations with defaults that match the leader's default
func Test_ReceivedObservationsWithSomeMisMatchedDefaults_InsufficientForConsensus(t *testing.T) {
	lggr := logger.Test(t)
	ctx := t.Context()

	md1 := newRequestMetaData()
	md1.KeyBundleID = "evm"

	expectedFailureCode := oracletypes.ConsensusFailureCode_FAILED_TO_CALCULATE_CONSENSUS_MDD

	reqToObservations := map[string]*consensusPluginTest{
		md1.RequestID(): {
			requests: []*oracle.ConsensusRequest{
				newCrWithObsAndDef(t, 10, 17, md1), newCrWithObsAndDef(t, 20, 12, md1), newCrWithObsAndDef(t, 30, 17, md1),
				newCrWithObsAndDef(t, 40, 16, md1), newCrWithObsAndDef(t, 50, 15, md1), newCrWithObsAndDef(t, 60, 11, md1),
				newCrWithObsAndDef(t, 70, 15, md1),
			},
			expectedConsensusFailureMessage: "failed to calculate consensus metadata, descriptor and default for request: no values met f+1 threshold",
			expectedConsensusFailureCode:    &expectedFailureCode,
		},
	}

	runProtocolRoundTests(ctx, t, lggr, n, f, reqToObservations)
}

func Test_MisMatchedDefaults_SufficientForConsensus_ReturnsDefault(t *testing.T) {
	lggr := logger.Test(t)
	ctx := t.Context()

	md1 := newRequestMetaData()
	md1.KeyBundleID = "evm"

	md2 := md1
	md2.WorkflowOwner = generateRandomHexString(20)

	md3 := md1
	md3.WorkflowOwner = generateRandomHexString(20)

	reqToObservations := map[string]*consensusPluginTest{
		md1.RequestID(): {
			requests: []*oracle.ConsensusRequest{
				newIdenticalCrWithDefault(t, 10, 14, md2), newIdenticalCrWithDefault(t, 20, 17, md3), newIdenticalCrWithDefault(t, 30, 17, md3),
				newIdenticalCrWithDefault(t, 40, 16, md1), newIdenticalCrWithDefault(t, 50, 17, md3), newIdenticalCrWithDefault(t, 60, 15, md1),
				newIdenticalCrWithDefault(t, 70, 19, md1),
			},
			verifyReport: func(t *testing.T, report ocr3types.ReportPlus[[]byte], infos *structpb.Struct) {
				verifyValueConsensusReport(t, report, infos, values.NewInt64(17), "evm")
			},
		},
	}

	runProtocolRoundTests(ctx, t, lggr, n, f, reqToObservations)
}

func Test_MisMatchedDefaults_InsufficientForConsensus(t *testing.T) {
	lggr := logger.Test(t)
	ctx := t.Context()

	md1 := newRequestMetaData()
	md1.KeyBundleID = "evm"

	expectedFailureCode := oracletypes.ConsensusFailureCode_FAILED_TO_CALCULATE_CONSENSUS_MDD

	reqToObservations := map[string]*consensusPluginTest{
		md1.RequestID(): {
			requests: []*oracle.ConsensusRequest{
				newCrWithObsAndDef(t, 10, 14, md1), newCrWithObsAndDef(t, 20, 15, md1), newCrWithObsAndDef(t, 30, 15, md1),
				newCrWithObsAndDef(t, 40, 16, md1), newCrWithObsAndDef(t, 50, 16, md1), newCrWithObsAndDef(t, 60, 17, md1),
				newCrWithObsAndDef(t, 70, 17, md1),
			},
			expectedConsensusFailureMessage: "failed to calculate consensus metadata, descriptor and default for request: no values met f+1 threshold",
			expectedConsensusFailureCode:    &expectedFailureCode,
		},
	}

	runProtocolRoundTests(ctx, t, lggr, n, f, reqToObservations)
}

func Test_MissingButSufficientObservations(t *testing.T) {
	lggr := logger.Test(t)
	ctx := t.Context()
	md3 := newRequestMetaData()
	md4 := newRequestMetaData()
	md5 := newRequestMetaData()

	reqToObservations := map[string]*consensusPluginTest{
		// Simulate some rounds where some nodes have not yet received the observation for req-3 and req-4
		md3.RequestID(): {
			requests: []*oracle.ConsensusRequest{
				newCr(t, 110, md3), newCr(t, 120, md3), newCr(t, 130, md3),
				newCr(t, 140, md3), newCr(t, 150, md3), nil, nil,
			},
			verifyReport: func(t *testing.T, report ocr3types.ReportPlus[[]byte], infos *structpb.Struct) {
				verifyValueConsensusReport(t, report, infos, values.NewInt64(130), "")
			},
		},
		md4.RequestID(): {
			requests: []*oracle.ConsensusRequest{
				newCr(t, 110, md4), nil, newCr(t, 130, md4),
				newCr(t, 140, md4), newCr(t, 150, md4), nil,
				newCr(t, 170, md4),
			},
			verifyReport: func(t *testing.T, report ocr3types.ReportPlus[[]byte], infos *structpb.Struct) {
				verifyValueConsensusReport(t, report, infos, values.NewInt64(140), "")
			},
		},
		md5.RequestID(): {
			requests: []*oracle.ConsensusRequest{
				newCr(t, 110, md5), nil, newCr(t, 130, md5),
				nil, newCr(t, 150, md5), nil,
				newCr(t, 170, md5),
			},
			verifyReport: nil,
		},
	}

	runProtocolRoundTests(ctx, t, lggr, n, f, reqToObservations)
}

func Test_InsufficientObservations(t *testing.T) {
	lggr := logger.Test(t)
	ctx := t.Context()

	md1 := newRequestMetaData()
	md6 := newRequestMetaData()

	md1.KeyBundleID = "evm"

	reqToObservations := map[string]*consensusPluginTest{
		md1.RequestID(): {
			requests: []*oracle.ConsensusRequest{
				newCr(t, 10, md1), newCr(t, 20, md1), newCr(t, 30, md1),
				newCr(t, 40, md1), newCr(t, 50, md1), newCr(t, 60, md1),
				newCr(t, 70, md1),
			},
			verifyReport: func(t *testing.T, report ocr3types.ReportPlus[[]byte], infos *structpb.Struct) {
				verifyValueConsensusReport(t, report, infos, values.NewInt64(40), "evm")
			},
		},

		// Simulate a round where there are insufficient observations for req-6
		md6.RequestID(): {
			requests: []*oracle.ConsensusRequest{
				newCr(t, 110, md6), nil, newCr(t, 130, md6),
				newCr(t, 140, md6), newCr(t, 150, md6), nil, nil,
			},
			verifyReport: nil,
		},
	}

	runProtocolRoundTests(ctx, t, lggr, n, f, reqToObservations)
}

func Test_LeaderHasNoMatchingRequest(t *testing.T) {
	lggr := logger.Test(t)
	ctx := t.Context()

	md1 := newRequestMetaData()
	md7 := newRequestMetaData()

	md1.KeyBundleID = "evm"

	reqToObservations := map[string]*consensusPluginTest{
		md1.RequestID(): {
			requests: []*oracle.ConsensusRequest{
				newCr(t, 10, md1), newCr(t, 20, md1), newCr(t, 30, md1),
				newCr(t, 40, md1), newCr(t, 50, md1), newCr(t, 60, md1),
				newCr(t, 70, md1),
			},
			verifyReport: func(t *testing.T, report ocr3types.ReportPlus[[]byte], infos *structpb.Struct) {
				verifyValueConsensusReport(t, report, infos, values.NewInt64(40), "evm")
			},
		},

		// Simulate a round where the leader has not yet received the observation for req-7
		md7.RequestID(): {
			requests: []*oracle.ConsensusRequest{
				nil, newCr(t, 120, md7), newCr(t, 130, md7),
				newCr(t, 140, md7), newCr(t, 150, md7), newCr(t, 160, md7),
				newCr(t, 170, md7),
			},
			verifyReport: nil,
		},
	}

	runProtocolRoundTests(ctx, t, lggr, n, f, reqToObservations)
}

func Test_WithOutcomeContext(t *testing.T) {
	lggr := logger.Test(t)
	ctx := t.Context()

	md1 := newRequestMetaData()

	md1.KeyBundleID = "evm"

	reqToObservations := map[string]*consensusPluginTest{
		md1.RequestID(): {
			requests: []*oracle.ConsensusRequest{
				newCr(t, 10, md1), newCr(t, 20, md1), newCr(t, 30, md1),
				newCr(t, 40, md1), newCr(t, 50, md1), newCr(t, 60, md1),
				newCr(t, 70, md1),
			},
			verifyReport: func(t *testing.T, report ocr3types.ReportPlus[[]byte], infos *structpb.Struct) {
				verifyValueConsensusReport(t, report, infos, values.NewInt64(40), "evm")
			},
		},
	}

	runProtocolRoundTests(ctx, t, lggr, n, f, reqToObservations)
}

func newRequestMetaData() oracle.ConsensusRequestMetadata {
	return oracle.ConsensusRequestMetadata{
		RequestMetadata: capabilities.RequestMetadata{
			WorkflowID:    "0039525c34de895c8fa68006bd63f6ce4a45ef1bc66377e791c6a8ae803dc0e4",
			WorkflowOwner: "1139525c34de895c8fa68006bd634387a9f1192a",

			WorkflowExecutionID:      generateRandomHexString(32),
			WorkflowName:             "a1b2c3d4e5f6a1b2c3d4",
			WorkflowDonID:            1,
			WorkflowDonConfigVersion: 1,
			ReferenceID:              "01",
			DecodedWorkflowName:      "test-workflow-decoded",
			SpendLimits:              nil,
		},
		KeyBundleID: "",
		ReportID:    generateRandomHexString(2),
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

func newIdenticalCr(t *testing.T, observation int64, metaData oracle.ConsensusRequestMetadata) *oracle.ConsensusRequest {
	simpleConsensusInputs := &sdk.SimpleConsensusInputs{
		Observation: &sdk.SimpleConsensusInputs_Value{Value: values.Proto(values.NewInt64(observation))},
		Descriptors: &sdk.ConsensusDescriptor{Descriptor_: &sdk.ConsensusDescriptor_Aggregation{Aggregation: sdk.AggregationType_AGGREGATION_TYPE_IDENTICAL}},
	}

	return oracle.NewConsensusRequest(serializeDeserialize(t, simpleConsensusInputs), time.Now(), time.Now().Add(1*time.Hour).UTC(), nil, metaData)
}

func newIdenticalValueCr(t *testing.T, observation values.Value, metaData oracle.ConsensusRequestMetadata) *oracle.ConsensusRequest {
	simpleConsensusInputs := &sdk.SimpleConsensusInputs{
		Observation: &sdk.SimpleConsensusInputs_Value{Value: values.Proto(observation)},
		Descriptors: &sdk.ConsensusDescriptor{Descriptor_: &sdk.ConsensusDescriptor_Aggregation{Aggregation: sdk.AggregationType_AGGREGATION_TYPE_IDENTICAL}},
	}

	return oracle.NewConsensusRequest(serializeDeserialize(t, simpleConsensusInputs), time.Now(), time.Now().Add(1*time.Hour).UTC(), nil, metaData)
}

func newIdenticalCrWithDefault(t *testing.T, observation int64, defaultObs int64, metaData oracle.ConsensusRequestMetadata) *oracle.ConsensusRequest {
	simpleConsensusInputs := &sdk.SimpleConsensusInputs{
		Observation: &sdk.SimpleConsensusInputs_Value{Value: values.Proto(values.NewInt64(observation))},
		Default:     values.Proto(values.NewInt64(defaultObs)),
		Descriptors: &sdk.ConsensusDescriptor{Descriptor_: &sdk.ConsensusDescriptor_Aggregation{Aggregation: sdk.AggregationType_AGGREGATION_TYPE_IDENTICAL}},
	}

	return oracle.NewConsensusRequest(serializeDeserialize(t, simpleConsensusInputs), time.Now(), time.Now().Add(1*time.Hour).UTC(), nil, metaData)
}

func newCr(t *testing.T, observation int64, metaData oracle.ConsensusRequestMetadata) *oracle.ConsensusRequest {
	simpleConsensusInputs := &sdk.SimpleConsensusInputs{
		Observation: &sdk.SimpleConsensusInputs_Value{Value: values.Proto(values.NewInt64(observation))},
		Descriptors: &sdk.ConsensusDescriptor{Descriptor_: &sdk.ConsensusDescriptor_Aggregation{Aggregation: sdk.AggregationType_AGGREGATION_TYPE_MEDIAN}},
	}

	return oracle.NewConsensusRequest(serializeDeserialize(t, simpleConsensusInputs), time.Now(), time.Now().Add(1*time.Hour).UTC(), nil, metaData)
}

func serializeDeserialize(t *testing.T, simpleConsensusInputs *sdk.SimpleConsensusInputs) *sdk.SimpleConsensusInputs {
	serialized, marshalErr := proto.Marshal(simpleConsensusInputs)
	require.NoError(t, marshalErr)

	deserialized := &sdk.SimpleConsensusInputs{}
	unmarshalErr := proto.Unmarshal(serialized, deserialized)
	require.NoError(t, unmarshalErr)
	return deserialized
}

func newCrWithError(t *testing.T, crErr error, metaData oracle.ConsensusRequestMetadata) *oracle.ConsensusRequest {
	simpleConsensusInputs := &sdk.SimpleConsensusInputs{
		Observation: &sdk.SimpleConsensusInputs_Error{
			Error: crErr.Error(),
		},
		Descriptors: &sdk.ConsensusDescriptor{Descriptor_: &sdk.ConsensusDescriptor_Aggregation{Aggregation: sdk.AggregationType_AGGREGATION_TYPE_MEDIAN}},
	}

	return oracle.NewConsensusRequest(serializeDeserialize(t, simpleConsensusInputs), time.Now(), time.Now().Add(1*time.Hour).UTC(), nil, metaData)
}

func newCrWithErrorAndDefault(t *testing.T, crErr error, def int64, metaData oracle.ConsensusRequestMetadata) *oracle.ConsensusRequest {
	simpleConsensusInputs := &sdk.SimpleConsensusInputs{
		Observation: &sdk.SimpleConsensusInputs_Error{
			Error: crErr.Error(),
		},
		Default:     values.Proto(values.NewInt64(def)),
		Descriptors: &sdk.ConsensusDescriptor{Descriptor_: &sdk.ConsensusDescriptor_Aggregation{Aggregation: sdk.AggregationType_AGGREGATION_TYPE_MEDIAN}},
	}

	return oracle.NewConsensusRequest(serializeDeserialize(t, simpleConsensusInputs), time.Now(), time.Now().Add(1*time.Hour).UTC(), nil, metaData)
}

func newCrWithObsAndDef(t *testing.T, observation int64, def int64, metaData oracle.ConsensusRequestMetadata) *oracle.ConsensusRequest {
	simpleConsensusInputs := &sdk.SimpleConsensusInputs{
		Observation: &sdk.SimpleConsensusInputs_Value{Value: values.Proto(values.NewInt64(observation))},
		Default:     values.Proto(values.NewInt64(def)),
		Descriptors: &sdk.ConsensusDescriptor{Descriptor_: &sdk.ConsensusDescriptor_Aggregation{Aggregation: sdk.AggregationType_AGGREGATION_TYPE_MEDIAN}},
	}

	return oracle.NewConsensusRequest(serializeDeserialize(t, simpleConsensusInputs), time.Now(), time.Now().Add(1*time.Hour).UTC(), nil, metaData)
}

// nillable observation and nillable default value, -1 indicates the value should be set as nil
func newNillableCr(t *testing.T, observation int64, def int64, metaData oracle.ConsensusRequestMetadata) *oracle.ConsensusRequest {
	observationVal, err := values.Wrap(nil)
	require.NoError(t, err, "failed to wrap nil value")

	var defaultVal values.Value

	if observation != -1 {
		observationVal, err = values.Wrap(values.NewInt64(observation))
		require.NoError(t, err, "failed to wrap observation value")
	}

	var simpleConsensusInputs *sdk.SimpleConsensusInputs
	if def != -1 {
		defaultVal, err = values.Wrap(values.NewInt64(def))
		require.NoError(t, err, "failed to wrap default value")
		simpleConsensusInputs = &sdk.SimpleConsensusInputs{
			Observation: &sdk.SimpleConsensusInputs_Value{Value: values.Proto(observationVal)},
			Default:     values.Proto(defaultVal),
			Descriptors: &sdk.ConsensusDescriptor{Descriptor_: &sdk.ConsensusDescriptor_Aggregation{Aggregation: sdk.AggregationType_AGGREGATION_TYPE_MEDIAN}},
		}
	} else {
		simpleConsensusInputs = &sdk.SimpleConsensusInputs{
			Observation: &sdk.SimpleConsensusInputs_Value{Value: values.Proto(observationVal)},
			Descriptors: &sdk.ConsensusDescriptor{Descriptor_: &sdk.ConsensusDescriptor_Aggregation{Aggregation: sdk.AggregationType_AGGREGATION_TYPE_MEDIAN}},
		}
	}

	return oracle.NewConsensusRequest(serializeDeserialize(t, simpleConsensusInputs), time.Now(), time.Now().Add(1*time.Hour).UTC(), nil, metaData)
}

type pluginAndRequestStore struct {
	plugin ocr3types.ReportingPlugin[[]byte]
	store  *requests.Store[*oracle.ConsensusRequest]
}

func runProtocolRoundTests(ctx context.Context, t *testing.T, lggr logger.Logger, n, f int, reqToObservations map[string]*consensusPluginTest) {
	pluginAndRequestStores := createPluginsAndStores(n, t, lggr, f, 5, 1000)

	addRequestsToAllStores(pluginAndRequestStores, reqToObservations, t)
	runProtocolRoundTestsWithPlugins(ctx, t, reqToObservations, pluginAndRequestStores, ocr3types.OutcomeContext{})
}

func createPluginsAndStores(n int, t *testing.T, lggr logger.Logger, f int, outcomeExpirySpan uint64, maxRequestOutcomeSize int) []pluginAndRequestStore {
	var pluginAndRequestStores []pluginAndRequestStore

	for i := 0; i < n; i++ {
		reportingPlugin, reqStore := createReportingPlugin(t, lggr, f, n, outcomeExpirySpan, maxRequestOutcomeSize)
		pluginAndRequestStores = append(pluginAndRequestStores, pluginAndRequestStore{
			plugin: reportingPlugin,
			store:  reqStore,
		})
	}
	return pluginAndRequestStores
}

// runProtocolRoundTestsWithPlugins simulates a single protocol round with the provided reporting plugins and request stores
// It verifies that all plugins reach the same outcome and that the reports generated are as expected according to the test
// and returns the outcome
func runProtocolRoundTestsWithPlugins(ctx context.Context, t *testing.T,
	reqToObservations map[string]*consensusPluginTest, pluginAndRequestStores []pluginAndRequestStore, previousOutcome ocr3types.OutcomeContext,
) ocr3types.Outcome {
	// Simulate a protocol round
	// Select the first reporting plugin as the leader, note that setting the observation to nil for the leader
	// will result in a nil outcome for that request
	leaderPlugin := pluginAndRequestStores[0].plugin

	query, err := leaderPlugin.Query(ctx, previousOutcome)
	require.NoError(t, err)

	var attributedObservations []libocrTypes.AttributedObservation
	for oracleIdx, plugin := range pluginAndRequestStores {
		observation, err := plugin.plugin.Observation(ctx, previousOutcome, query)

		fmt.Printf("Oracle %d observation: %v\n", oracleIdx, observation)

		require.NoError(t, err, "failed to get observation from reporting plugin")
		attributedObservations = append(attributedObservations, libocrTypes.AttributedObservation{
			Observation: observation,
			Observer:    commontypes.OracleID(oracleIdx), //nolint:gosec // G115
		})
	}

	for _, pluginAndStore := range pluginAndRequestStores {
		for _, obs := range attributedObservations {
			err := pluginAndStore.plugin.ValidateObservation(ctx, previousOutcome, query, obs)
			require.NoError(t, err, "failed to validate observation from reporting plugin")
		}
	}

	for _, pluginAndStore := range pluginAndRequestStores {
		quorumReached, err := pluginAndStore.plugin.ObservationQuorum(ctx, previousOutcome, query, attributedObservations)
		require.NoError(t, err, "failed to validate observation from reporting plugin")
		require.True(t, quorumReached, "quorum should be reached for observation")
	}

	var nodeOutcomes []ocr3types.Outcome
	for _, pluginAndStore := range pluginAndRequestStores {
		outcome, err := pluginAndStore.plugin.Outcome(ctx, previousOutcome, query, attributedObservations)
		if err != nil {
			continue
		}

		require.NoError(t, err, "failed to get outcome from reporting plugin")
		nodeOutcomes = append(nodeOutcomes, outcome)
	}

	// Verify that all outcomes are the same
	for i := 1; i < len(nodeOutcomes); i++ {
		require.True(t, bytes.Equal(nodeOutcomes[0], nodeOutcomes[i]), "outcomes should be equal across reporting plugins")
	}

	var allReports [][]ocr3types.ReportPlus[[]byte]
	for _, pluginAndStore := range pluginAndRequestStores {
		reports, err := pluginAndStore.plugin.Reports(ctx, 0, nodeOutcomes[0])
		require.NoError(t, err, "failed to report outcome from reporting plugin")

		outcome := &oracletypes.Outcome{}
		err = proto.Unmarshal(nodeOutcomes[0], outcome)
		require.NoError(t, err, "failed to unmarshal value from outcome")

		require.Len(t, reports, len(outcome.Outcomes), "reporting plugin returned wrong number of reports")
		allReports = append(allReports, reports)
	}

	// Verify all reports are the same for each request
	for i := 1; i < len(allReports); i++ {
		require.Equal(t, len(allReports[0]), len(allReports[i]), "number of reports should be equal across reporting plugins")

		for idx, reports := range allReports[0] {
			require.Equal(t, reports.ReportWithInfo.Report, allReports[i][idx].ReportWithInfo.Report, "reports should be equal across reporting plugins")
			require.Equal(t, reports.ReportWithInfo.Info, allReports[i][idx].ReportWithInfo.Info, "report info should be equal across reporting plugins")
		}
	}

	// Create a map to hold the request ID to expected result
	requestIDToOutcome := make(map[string]*consensusPluginTest)
	for reqID, obs := range reqToObservations {
		requestIDToOutcome[reqID] = obs
	}

	// Get reports and verify the value selected
	reports := allReports[0]
	receivedReportForRequestIDs := map[string]bool{}
	receivedFailureMessageForRequestIDs := map[string]bool{}
	receivedFailureCodeForRequestIDs := map[string]bool{}
	for _, report := range reports {
		var infos structpb.Struct
		err = proto.Unmarshal(report.ReportWithInfo.Info, &infos)
		require.NoError(t, err, "failed to unmarshal value from report")

		infoMap := infos.AsMap()

		if failureCodeVal, exists := infoMap[plugin.InfoConsensusFailureCode]; exists {
			reqID := infos.Fields[plugin.InfoRequestID].GetStringValue()
			expectedOutcome, ok := requestIDToOutcome[reqID]
			require.True(t, ok, "got report for a request without a test outcome %s", reqID)

			if expectedOutcome.expectedConsensusFailureCode == nil {
				require.FailNow(t, "not expecting failure code for request %s", reqID)
			}

			failureCodeStr := failureCodeVal.(string)

			codeInt, ok := oracletypes.ConsensusFailureCode_value[failureCodeStr]
			require.True(t, ok)

			failureCode := oracletypes.ConsensusFailureCode(codeInt)

			if *expectedOutcome.expectedConsensusFailureCode != failureCode {
				require.FailNow(t, "expected failure code %s but got %s for request %s", expectedOutcome.expectedConsensusFailureCode.String(), failureCode.String(), reqID)
			}

			receivedFailureCodeForRequestIDs[reqID] = true

			failureMessage, exists := infoMap[plugin.InfoConsensusFailureMessage]
			if !exists {
				require.FailNow(t, "expected failure message for request %s", reqID)
			}

			if len(expectedOutcome.expectedConsensusFailureMessage) == 0 {
				require.FailNow(t, "expected outcome failure message for request %s", reqID)
			}

			receivedFailureMessageForRequestIDs[reqID] = true

			require.Contains(t, failureMessage.(string), expectedOutcome.expectedConsensusFailureMessage)
		} else {
			serialisedValue := report.ReportWithInfo.Report[plugin.ReportMetaDataPrependLength:]
			actualProto := &valuespb.Value{}
			err := proto.Unmarshal(serialisedValue, actualProto)
			require.NoError(t, err, "failed to unmarshal value from report")

			reqID := infos.Fields[plugin.InfoRequestID].GetStringValue()

			receivedReportForRequestIDs[reqID] = true
			expectedOutcome, ok := requestIDToOutcome[reqID]
			require.True(t, ok, "got report for a request without a test outcome %s", reqID)

			if expectedOutcome.verifyReport != nil {
				expectedOutcome.verifyReport(t, report, &infos)
			} else {
				require.FailNow(t, "not expecting report for request %s", reqID)
			}
		}
	}

	// Verify all expected reports were received
	// Count expected reports (those with verifyReport set) and verify we got the right number
	expectedReportCount := 0
	for _, outcome := range requestIDToOutcome {
		if outcome.verifyReport != nil {
			expectedReportCount++
		}
	}
	actualReportCount := 0
	for reqID := range receivedReportForRequestIDs {
		if receivedReportForRequestIDs[reqID] {
			actualReportCount++
		}
	}
	require.Equal(t, expectedReportCount, actualReportCount, "expected %d success reports but got %d", expectedReportCount, actualReportCount)

	// Verify each expected report was received
	for reqID, outcome := range requestIDToOutcome {
		if outcome.verifyReport != nil {
			require.True(t, receivedReportForRequestIDs[reqID], "expected report for request ID %s was not received", reqID)
		}
	}

	// Verify all expected failure messages were received
	for reqID, outcome := range requestIDToOutcome {
		if len(outcome.expectedConsensusFailureMessage) > 0 {
			require.True(t, receivedFailureMessageForRequestIDs[reqID], "expected failure message for request ID %s was not received", reqID)
		}
	}

	// Veriofy all expected failure codes were received
	for reqID, outcome := range requestIDToOutcome {
		if outcome.expectedConsensusFailureCode != nil {
			require.True(t, receivedFailureCodeForRequestIDs[reqID], "expected failure code for request ID %s was not received", reqID)
		}
	}

	return nodeOutcomes[0]
}

func addRequestsToAllStores(pluginAndRequestStores []pluginAndRequestStore, reqToObservations map[string]*consensusPluginTest, t *testing.T) {
	for i := 0; i < len(pluginAndRequestStores); i++ {
		var pluginObs []*oracle.ConsensusRequest

		for _, obsData := range reqToObservations {
			observation := obsData.requests[i]
			if observation != nil {
				pluginObs = append(pluginObs, observation)
			}
		}
		for _, obs := range pluginObs {
			req := obs
			err := pluginAndRequestStores[i].store.Add(req)
			require.NoError(t, err, "failed to add request to store")
		}
	}
}

func removeRequestFromAllStores(pluginAndRequestStores []pluginAndRequestStore, requestID string) {
	for i := 0; i < len(pluginAndRequestStores); i++ {
		pluginAndRequestStores[i].store.Evict(requestID)
	}
}

func verifyValueConsensusReport(t *testing.T, report ocr3types.ReportPlus[[]byte], infos *structpb.Struct, expectedResult values.Value,
	expectedKeyBundleName string,
) {
	require.NotNil(t, report.ReportWithInfo, "report should not be nil")
	require.NotNil(t, report.ReportWithInfo.Report, "report value should not be nil")

	// Verify that the report info contains request ID
	require.NotNil(t, infos, "report info should not be nil")
	require.Contains(t, infos.Fields, plugin.InfoRequestID, "report info should contain request ID")

	if expectedKeyBundleName != "" {
		require.Contains(t, infos.Fields, "keyBundleName", "report info should contain key bundle name")
		keyBundleName := infos.Fields["keyBundleName"].GetStringValue()
		require.Equal(t, expectedKeyBundleName, keyBundleName, "keyBundle name should match expected value")
	}

	// Verify the value in the report matches the expected result
	serialisedValue := report.ReportWithInfo.Report[plugin.ReportMetaDataPrependLength:]
	actualProto := &valuespb.Value{}
	err := proto.Unmarshal(serialisedValue, actualProto)
	require.NoError(t, err, "failed to unmarshal value from report")

	expectedProto := values.Proto(expectedResult)
	fmt.Printf("Expected outcome: %s, Actual outcome: %s\n", expectedProto, actualProto)
	require.True(t, proto.Equal(actualProto, expectedProto), "expected outcome value to match expected value")
}

func createReportingPlugin(t *testing.T, lggr logger.Logger, f int, n int,
	outcomeExpirySpan uint64, maxRequestOutcomeSize int,
) (ocr3types.ReportingPlugin[[]byte], *requests.Store[*oracle.ConsensusRequest]) {
	reqStore := requests.NewStore[*oracle.ConsensusRequest]()

	metricsInstance, err := metrics.NewMetrics()
	require.NoError(t, err)

	rp, err := plugin.NewReportingPlugin(lggr, metricsInstance, f, n, reqStore, &pbtypes.ReportingPluginConfig{
		MaxQueryLengthBytes:              defaultMaxLengthBytes,
		MaxObservationLengthBytes:        defaultMaxLengthBytes,
		MaxOutcomeLengthBytes:            defaultMaxLengthBytes,
		MaxReportLengthBytes:             defaultMaxLengthBytes,
		MaxReportCount:                   defaultMaxReportCount,
		HistoricalOutcomeExpirySeqNrSpan: outcomeExpirySpan,
	}, "evm", maxRequestOutcomeSize)
	require.NoError(t, err)
	return plugin.AsOCR3ReportingPlugin(rp), reqStore
}

// Test_QuorumReachedWithSomeOversizedRequests verifies that consensus can be reached
// even when some nodes have oversized requests that are converted to error inputs.
// This simulates the scenario where some HTTP endpoints return oversized data near the cutoff,
// but a quorum of nodes get correctly sized data and can still reach consensus.
func Test_QuorumReachedWithSomeOversizedRequests(t *testing.T) {
	lggr := logger.Test(t)
	ctx := t.Context()

	md1 := newRequestMetaData()
	md1.KeyBundleID = "evm"

	// Create error input that simulates an oversized request error
	// This is what would be sent when validateRequestSize returns a user error
	oversizedError := "request size 1048577 bytes exceeds maximum allowed size: PerWorkflow.Consensus.ObservationSizeLimit limited for workflow[wf-id]: cannot use 1.000001mb, limit is 1mb"

	reqToObservations := map[string]*consensusPluginTest{
		md1.RequestID(): {
			requests: []*oracle.ConsensusRequest{
				// 5 nodes have valid observations (quorum)
				newCr(t, 10, md1),
				newCr(t, 20, md1),
				newCr(t, 30, md1),
				newCr(t, 40, md1),
				newCr(t, 50, md1),
				// 2 nodes have oversized requests converted to error inputs (below f+1 = 3)
				newCrWithError(t, errors.New(oversizedError), md1),
				newCrWithError(t, errors.New(oversizedError), md1),
			},
			verifyReport: func(t *testing.T, report ocr3types.ReportPlus[[]byte], infos *structpb.Struct) {
				// With 5 valid observations and median aggregation, the median should be 30
				verifyValueConsensusReport(t, report, infos, values.NewInt64(30), "evm")
			},
		},
	}

	runProtocolRoundTests(ctx, t, lggr, n, f, reqToObservations)
}
