package plugin_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/blobconsensus/oracle/plugin"
	oracletypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/blobconsensus/oracle/types"

	"github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"
	valuespb "github.com/smartcontractkit/chainlink-protos/cre/go/values/pb"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	libocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

// Verifies Outcome() is deterministic for identical inputs even if the observations arrive in a different order.
func TestOutcomeDeterministicWithErrors(t *testing.T) {
	lggr := logger.Test(t)
	ctx := context.Background()

	const f, n = 1, 4
	reportingPlugin, _ := createReportingPlugin(t, lggr, f, n, 5, defaultMaxLengthBytes)

	md := newRequestMetaData()
	md.KeyBundleID = "evm"
	reqID := md.RequestID()

	makeObs := func(errMsg string, observerID uint8) libocrtypes.AttributedObservation {
		ro := &oracletypes.RequestObservation{
			Metadata: plugin.ToRequestMetaData(md),
			// individual errors are included in the aggregated outcome when consensus fails
			Input:      &sdk.SimpleConsensusInputs{Observation: &sdk.SimpleConsensusInputs_Error{Error: errMsg}},
			ReceivedAt: timestamppb.Now(),
		}
		obs := &oracletypes.Observation{
			Observations: map[string]*oracletypes.RequestObservation{
				reqID: ro,
			},
		}
		b, err := proto.Marshal(obs)
		require.NoError(t, err)
		return libocrtypes.AttributedObservation{
			Observation: b,
			Observer:    commontypes.OracleID(observerID),
		}
	}

	attributed := []libocrtypes.AttributedObservation{
		makeObs("first error", 0),
		makeObs("second error", 1),
		makeObs("third error", 2),
	}

	qBytes, err := proto.Marshal(&oracletypes.Query{RequestIDs: []string{reqID}})
	require.NoError(t, err)

	outctx := ocr3types.OutcomeContext{SeqNr: 1}

	outcomeA, err := reportingPlugin.Outcome(ctx, outctx, qBytes, attributed)
	require.NoError(t, err)

	// Reverse attributed observations to mimic nondeterministic ordering from map iteration.
	reversed := make([]libocrtypes.AttributedObservation, len(attributed))
	for i := range attributed {
		reversed[i] = attributed[len(attributed)-1-i]
	}

	outcomeB, err := reportingPlugin.Outcome(ctx, outctx, qBytes, reversed)
	require.NoError(t, err)

	require.Equal(t, outcomeA, outcomeB, "Outcome must be deterministic regardless of observation order")
}

// TestOutcomeDeterminism_FieldsMapMultipleFieldsFail verifies that Outcome() is
// deterministic when consensus fails for a FieldsMap descriptor with multiple fields
// failing aggregation and no defaults. handleFieldsMapAggregation iterates desc keys
// in sorted order so the first-failing field (and thus the failure message embedded
// in the outcome) is always the same.
func TestOutcomeDeterminism_FieldsMapMultipleFieldsFail(t *testing.T) {
	lggr := logger.Test(t)
	ctx := context.Background()

	const testF, testN = 1, 4
	reportingPlugin, _ := createReportingPlugin(t, lggr, testF, testN, 5, defaultMaxLengthBytes)

	md := newRequestMetaData()
	md.KeyBundleID = "evm"
	reqID := md.RequestID()

	fieldsMapDescriptor := &sdk.ConsensusDescriptor{
		Descriptor_: &sdk.ConsensusDescriptor_FieldsMap{
			FieldsMap: &sdk.FieldsMap{
				Fields: map[string]*sdk.ConsensusDescriptor{
					"Alpha":   {Descriptor_: &sdk.ConsensusDescriptor_Aggregation{Aggregation: sdk.AggregationType_AGGREGATION_TYPE_IDENTICAL}},
					"Beta":    {Descriptor_: &sdk.ConsensusDescriptor_Aggregation{Aggregation: sdk.AggregationType_AGGREGATION_TYPE_IDENTICAL}},
					"Gamma":   {Descriptor_: &sdk.ConsensusDescriptor_Aggregation{Aggregation: sdk.AggregationType_AGGREGATION_TYPE_IDENTICAL}},
					"Delta":   {Descriptor_: &sdk.ConsensusDescriptor_Aggregation{Aggregation: sdk.AggregationType_AGGREGATION_TYPE_IDENTICAL}},
					"Epsilon": {Descriptor_: &sdk.ConsensusDescriptor_Aggregation{Aggregation: sdk.AggregationType_AGGREGATION_TYPE_IDENTICAL}},
				},
			},
		},
	}

	makeObsWithFieldsMap := func(observerID uint8) libocrtypes.AttributedObservation {
		obsValue := valuespb.NewMapValue(map[string]*valuespb.Value{
			"Alpha":   values.Proto(values.NewString(fmt.Sprintf("alpha-%d", observerID))),
			"Beta":    values.Proto(values.NewString(fmt.Sprintf("beta-%d", observerID))),
			"Gamma":   values.Proto(values.NewString(fmt.Sprintf("gamma-%d", observerID))),
			"Delta":   values.Proto(values.NewString(fmt.Sprintf("delta-%d", observerID))),
			"Epsilon": values.Proto(values.NewString(fmt.Sprintf("epsilon-%d", observerID))),
		})

		ro := &oracletypes.RequestObservation{
			Metadata: plugin.ToRequestMetaData(md),
			Input: &sdk.SimpleConsensusInputs{
				Observation: &sdk.SimpleConsensusInputs_Value{Value: obsValue},
				Descriptors: fieldsMapDescriptor,
			},
			ReceivedAt: timestamppb.Now(),
		}
		obs := &oracletypes.Observation{
			Observations: map[string]*oracletypes.RequestObservation{reqID: ro},
		}
		b, err := proto.Marshal(obs)
		require.NoError(t, err)
		return libocrtypes.AttributedObservation{
			Observation: b,
			Observer:    commontypes.OracleID(observerID),
		}
	}

	attributed := []libocrtypes.AttributedObservation{
		makeObsWithFieldsMap(0),
		makeObsWithFieldsMap(1),
		makeObsWithFieldsMap(2),
		makeObsWithFieldsMap(3),
	}

	qBytes, err := proto.Marshal(&oracletypes.Query{RequestIDs: []string{reqID}})
	require.NoError(t, err)

	outctx := ocr3types.OutcomeContext{SeqNr: 1}

	seenOutcomes := map[string]bool{}
	for range 200 {
		outcome, err := reportingPlugin.Outcome(ctx, outctx, qBytes, attributed)
		require.NoError(t, err)
		seenOutcomes[string(outcome)] = true
	}

	require.Equal(t, 1, len(seenOutcomes),
		"Outcome() must be deterministic when multiple FieldsMap fields fail aggregation")
}
