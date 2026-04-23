package plugin_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/blobconsensus/oracle"
	oracletypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/blobconsensus/oracle/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink-protos/cre/go/values"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
)

func Test_DuplicateOutcomePrevention(t *testing.T) {
	lggr := logger.Test(t)
	ctx := t.Context()

	md1 := newRequestMetaData()
	md2 := newRequestMetaData()

	md1.KeyBundleID = "evm"

	reqToObservations := map[string]*consensusPluginTest{
		md1.RequestID(): {requests: []*oracle.ConsensusRequest{
			newCr(t, 10, md1), newCr(t, 20, md1), newCr(t, 30, md1),
			newCr(t, 40, md1), newCr(t, 50, md1), newCr(t, 60, md1),
			newCr(t, 70, md1)},
			verifyReport: func(t *testing.T, report ocr3types.ReportPlus[[]byte], infos *structpb.Struct) {
				verifyValueConsensusReport(t, report, infos, values.NewInt64(40), "evm")
			},
		},
		md2.RequestID(): {requests: []*oracle.ConsensusRequest{
			newCr(t, 110, md2), newCr(t, 120, md2), newCr(t, 130, md2),
			newCr(t, 140, md2), newCr(t, 150, md2), newCr(t, 160, md2),
			newCr(t, 170, md2)},
			verifyReport: func(t *testing.T, report ocr3types.ReportPlus[[]byte], infos *structpb.Struct) {
				verifyValueConsensusReport(t, report, infos, values.NewInt64(140), "")
			},
		},
	}

	pluginsAndStores := createPluginsAndStores(n, t, lggr, f, 5, 1000)

	addRequestsToAllStores(pluginsAndStores, reqToObservations, t)

	seqNr := uint64(1)
	outcome := runProtocolRoundTestsWithPlugins(ctx, t, reqToObservations, pluginsAndStores, ocr3types.OutcomeContext{
		SeqNr: seqNr,
	})

	outcome1 := reqToObservations[md1.RequestID()]
	outcome1.verifyReport = nil

	outcome2 := reqToObservations[md2.RequestID()]
	outcome2.verifyReport = nil

	seqNr++
	runProtocolRoundTestsWithPlugins(ctx, t, reqToObservations, pluginsAndStores, ocr3types.OutcomeContext{
		PreviousOutcome: outcome,
		SeqNr:           seqNr,
	})
}

func Test_HistoricalOutcomesAreRemovedOnExpiry(t *testing.T) {
	lggr := logger.Test(t)
	ctx := t.Context()

	successfulRequest := newRequestMetaData()
	pendingRequest := newRequestMetaData()

	successfulRequest.KeyBundleID = "evm"

	verifyReport1 := func(t *testing.T, report ocr3types.ReportPlus[[]byte], infos *structpb.Struct) {
		verifyValueConsensusReport(t, report, infos, values.NewInt64(40), "evm")
	}

	verifyReport2 := func(t *testing.T, report ocr3types.ReportPlus[[]byte], infos *structpb.Struct) {
		verifyValueConsensusReport(t, report, infos, values.NewInt64(140), "")
	}

	reqToObservations := map[string]*consensusPluginTest{
		// Consensus will succeed
		successfulRequest.RequestID(): {requests: []*oracle.ConsensusRequest{
			newCr(t, 10, successfulRequest), newCr(t, 20, successfulRequest), newCr(t, 30, successfulRequest),
			newCr(t, 40, successfulRequest), newCr(t, 50, successfulRequest), newCr(t, 60, successfulRequest),
			newCr(t, 70, successfulRequest)},
			verifyReport: verifyReport1,
		},

		// Consensus will remain pending as < 2f+1 observations are received initially
		pendingRequest.RequestID(): {requests: []*oracle.ConsensusRequest{
			newCr(t, 110, pendingRequest), nil, newCr(t, 130, pendingRequest),
			newCr(t, 140, pendingRequest), nil, newCr(t, 160, pendingRequest),
			nil},
			verifyReport: verifyReport2,
		},
	}

	historicalOutcomeExpirySpan := uint64(3)
	pluginsAndStores := createPluginsAndStores(n, t, lggr, f, historicalOutcomeExpirySpan, 1000)

	addRequestsToAllStores(pluginsAndStores, reqToObservations, t)

	postProtocolRound := func(t *testing.T, outcome *oracletypes.Outcome, requestIDToOutcome map[string]*oracletypes.ConsensusSuccessOutcome,
		requestIDToHistoricalOutcome map[string]uint64) {
		// If a request is successful post round remove it here
		for requestID := range requestIDToOutcome {
			removeRequestFromAllStores(pluginsAndStores, requestID)
		}
	}

	var previousOutcome []byte
	for seqNr := uint64(1); seqNr <= 10; seqNr++ {
		var verifyOutcome func(t *testing.T, outcome *oracletypes.Outcome, requestIDToOutcome map[string]*oracletypes.ConsensusSuccessOutcome,
			requestIDToHistoricalOutcome map[string]uint64)

		switch seqNr {
		case 1:

			outcome1 := reqToObservations[successfulRequest.RequestID()]
			outcome1.verifyReport = verifyReport1

			outcome2 := reqToObservations[pendingRequest.RequestID()]
			outcome2.verifyReport = nil

			verifyOutcome = func(t *testing.T, outcome *oracletypes.Outcome, requestIDToOutcome map[string]*oracletypes.ConsensusSuccessOutcome,
				requestIDToHistoricalOutcome map[string]uint64) {
				require.Len(t, outcome.Outcomes, 1)
				require.Len(t, outcome.HistoricalOutcomes, 1)

				require.NotNil(t, requestIDToHistoricalOutcome[successfulRequest.RequestID()])
			}
		case 2:

			// Submit the original successfully processed request here to ensure the that the plugin uses historical outcomes to prevent duplicate outcome
			addRequestsToAllStores(pluginsAndStores, map[string]*consensusPluginTest{successfulRequest.RequestID(): reqToObservations[successfulRequest.RequestID()]}, t)

			outcome1 := reqToObservations[successfulRequest.RequestID()]
			outcome1.verifyReport = nil

			outcome2 := reqToObservations[pendingRequest.RequestID()]
			outcome2.verifyReport = nil

			verifyOutcome = func(t *testing.T, outcome *oracletypes.Outcome,
				requestIDToOutcome map[string]*oracletypes.ConsensusSuccessOutcome,
				requestIDToHistoricalOutcome map[string]uint64) {
				require.Len(t, outcome.Outcomes, 0)
				require.Len(t, outcome.HistoricalOutcomes, 1)
			}

		case 3:
			// Simulate the observations arriving for the pending request
			updatedPendingRequestObservations := map[string]*consensusPluginTest{
				pendingRequest.RequestID(): {requests: []*oracle.ConsensusRequest{
					newCr(t, 110, pendingRequest), newCr(t, 120, pendingRequest), newCr(t, 130, pendingRequest),
					newCr(t, 140, pendingRequest), newCr(t, 150, pendingRequest), newCr(t, 160, pendingRequest),
					newCr(t, 170, pendingRequest)},
					verifyReport: func(t *testing.T, report ocr3types.ReportPlus[[]byte], infos *structpb.Struct) {
						verifyValueConsensusReport(t, report, infos, values.NewInt64(140), "")
					}},
			}

			removeRequestFromAllStores(pluginsAndStores, pendingRequest.RequestID())
			addRequestsToAllStores(pluginsAndStores, updatedPendingRequestObservations, t)

			outcome1 := reqToObservations[successfulRequest.RequestID()]
			outcome1.verifyReport = nil

			outcome2 := reqToObservations[pendingRequest.RequestID()]
			outcome2.verifyReport = verifyReport2

			verifyOutcome = func(t *testing.T, outcome *oracletypes.Outcome,
				requestIDToOutcome map[string]*oracletypes.ConsensusSuccessOutcome,
				requestIDToHistoricalOutcome map[string]uint64) {
				require.Len(t, outcome.Outcomes, 1)
				require.Len(t, outcome.HistoricalOutcomes, 2)
			}
		case 4:
			outcome1 := reqToObservations[successfulRequest.RequestID()]
			outcome1.verifyReport = nil

			outcome2 := reqToObservations[pendingRequest.RequestID()]
			outcome2.verifyReport = nil

			verifyOutcome = func(t *testing.T, outcome *oracletypes.Outcome,
				requestIDToOutcome map[string]*oracletypes.ConsensusSuccessOutcome,
				requestIDToHistoricalOutcome map[string]uint64) {
				require.Len(t, outcome.Outcomes, 0)
				require.Len(t, outcome.HistoricalOutcomes, 2)
			}

			// Simulate the expiry of the pending request and eventual receipt of the outcome for the successful request, both of which
			// would result in them being removed from the request store.
			removeRequestFromAllStores(pluginsAndStores, pendingRequest.RequestID())
			removeRequestFromAllStores(pluginsAndStores, successfulRequest.RequestID())

		case 5:
			// By this point the historical record of the first outcome should have expired, but the historical record of the second outcome should still exist
			outcome1 := reqToObservations[successfulRequest.RequestID()]
			outcome1.verifyReport = nil

			outcome2 := reqToObservations[pendingRequest.RequestID()]
			outcome2.verifyReport = nil
			verifyOutcome = func(t *testing.T, outcome *oracletypes.Outcome,
				requestIDToOutcome map[string]*oracletypes.ConsensusSuccessOutcome,
				requestIDToHistoricalOutcome map[string]uint64) {
				require.Len(t, outcome.Outcomes, 0)
				require.Len(t, outcome.HistoricalOutcomes, 1)
			}
		case 6, 7, 8, 9, 10:
			// Eventually by this round all historical outcomes should have expired
			outcome1 := reqToObservations[successfulRequest.RequestID()]
			outcome1.verifyReport = nil

			outcome2 := reqToObservations[pendingRequest.RequestID()]
			outcome2.verifyReport = nil
			verifyOutcome = func(t *testing.T, outcome *oracletypes.Outcome, requestIDToOutcome map[string]*oracletypes.ConsensusSuccessOutcome,
				requestIDToHistoricalOutcome map[string]uint64) {
				if seqNr == 10 {
					require.Len(t, outcome.Outcomes, 0)
					require.Len(t, outcome.HistoricalOutcomes, 0)
				}
			}
		}

		previousOutcome = runProtocolRoundTestsWithPlugins(ctx, t, reqToObservations, pluginsAndStores, ocr3types.OutcomeContext{
			SeqNr:           seqNr,
			PreviousOutcome: previousOutcome})

		requestsOutcome := &oracletypes.Outcome{}
		err := proto.Unmarshal(previousOutcome, requestsOutcome)
		require.NoError(t, err)

		requestIDToOutcome := make(map[string]*oracletypes.ConsensusSuccessOutcome)
		for _, ro := range requestsOutcome.Outcomes {
			switch v := ro.GetOutcome().(type) {
			case *oracletypes.ConsensusOutcome_Success:
				requestIDToOutcome[v.Success.Metadata.RequestId] = v.Success
			default:
				t.Fatalf("expected ConsensusSuccessOutcome, got %T", v)
			}
		}

		verifyOutcome(t, requestsOutcome, requestIDToOutcome, requestsOutcome.HistoricalOutcomes)
		postProtocolRound(t, requestsOutcome, requestIDToOutcome, requestsOutcome.HistoricalOutcomes)
	}
}
