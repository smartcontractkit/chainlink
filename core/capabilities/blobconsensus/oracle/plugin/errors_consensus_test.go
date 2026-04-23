package plugin_test

import (
	"errors"
	"testing"

	"google.golang.org/protobuf/types/known/structpb"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/blobconsensus/oracle"
	oracletypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/blobconsensus/oracle/types"

	"github.com/smartcontractkit/chainlink-protos/cre/go/values"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
)

func Test_ReceivedTooManyErrors(t *testing.T) {
	lggr := logger.Test(t)
	ctx := t.Context()

	md1 := newRequestMetaData()
	md1.KeyBundleID = "evm"

	expectedFailureCode := oracletypes.ConsensusFailureCode_RECEIVED_FPLUS1_ERRORS

	reqToObservations := map[string]*consensusPluginTest{
		md1.RequestID(): {requests: []*oracle.ConsensusRequest{
			newCr(t, 10, md1), newCrWithError(t, errors.New("its broken"), md1), newCr(t, 30, md1),
			newCrWithError(t, errors.New("its broken"), md1), newCr(t, 50, md1), newCr(t, 60, md1),
			newCrWithError(t, errors.New("its broken"), md1)},
			expectedConsensusFailureMessage: "consensus calculation failed: received 3 errors which is >= f+1 (3)",
			expectedConsensusFailureCode:    &expectedFailureCode,
		},
	}

	runProtocolRoundTests(ctx, t, lggr, n, f, reqToObservations)
}

func Test_ReceivedLargeErrors(t *testing.T) {
	lggr := logger.Test(t)
	ctx := t.Context()

	md1 := newRequestMetaData()
	md1.KeyBundleID = "evm"

	expectedFailureCode := oracletypes.ConsensusFailureCode_RECEIVED_FPLUS1_ERRORS

	reqToObservations := map[string]*consensusPluginTest{
		md1.RequestID(): {requests: []*oracle.ConsensusRequest{
			newCr(t, 10, md1), newCrWithError(t, errors.New(createLargeString(10000)), md1), newCr(t, 30, md1),
			newCrWithError(t, errors.New(createLargeString(10000)), md1), newCr(t, 50, md1), newCr(t, 60, md1),
			newCrWithError(t, errors.New(createLargeString(10000)), md1)},
			expectedConsensusFailureMessage: ":TRUNCATED",
			expectedConsensusFailureCode:    &expectedFailureCode},
	}

	runProtocolRoundTests(ctx, t, lggr, n, f, reqToObservations)
}

func createLargeString(size int) string {
	b := make([]byte, size)
	for i := range b {
		b[i] = 'a'
	}
	return string(b)
}

func Test_ReceivedTooManyErrorsWithDefault(t *testing.T) {
	lggr := logger.Test(t)
	ctx := t.Context()

	md1 := newRequestMetaData()

	md1.KeyBundleID = "evm"

	reqToObservations := map[string]*consensusPluginTest{
		md1.RequestID(): {requests: []*oracle.ConsensusRequest{
			newCrWithObsAndDef(t, 10, 20, md1), newCrWithErrorAndDefault(t, errors.New("its broken"), 20, md1), newCrWithObsAndDef(t, 30, 20, md1),
			newCrWithObsAndDef(t, 40, 20, md1), newCrWithErrorAndDefault(t, errors.New("its broken"), 20, md1), newCrWithObsAndDef(t, 60, 20, md1),
			newCrWithErrorAndDefault(t, errors.New("its broken"), 20, md1)},
			verifyReport: func(t *testing.T, report ocr3types.ReportPlus[[]byte], infos *structpb.Struct) {
				verifyValueConsensusReport(t, report, infos, values.NewInt64(20), "evm")
			}},
	}

	runProtocolRoundTests(ctx, t, lggr, n, f, reqToObservations)
}

func Test_ReceivedSufficientObservationsAndSomeErrors(t *testing.T) {
	lggr := logger.Test(t)
	ctx := t.Context()

	md1 := newRequestMetaData()

	md1.KeyBundleID = "evm"

	reqToObservations := map[string]*consensusPluginTest{
		md1.RequestID(): {requests: []*oracle.ConsensusRequest{
			newCr(t, 10, md1), newCrWithError(t, errors.New("its broken"), md1), newCr(t, 30, md1),
			newCr(t, 40, md1), newCr(t, 50, md1), newCr(t, 60, md1),
			newCr(t, 70, md1)},
			verifyReport: func(t *testing.T, report ocr3types.ReportPlus[[]byte], infos *structpb.Struct) {
				verifyValueConsensusReport(t, report, infos, values.NewInt64(40), "evm")
			}},
	}

	runProtocolRoundTests(ctx, t, lggr, n, f, reqToObservations)
}
