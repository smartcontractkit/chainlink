package streams_test

import (
	"crypto/ed25519"
	"crypto/rand"
	"math/big"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	ocrTypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/datastreams"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/triggers"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote"
	remotetypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/streams"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

const (
	triggerID          = "streams-trigger@1.0.0"
	workflowID         = "workflowID1"
	baseTriggerEventID = "triggerEvent_"
	basePrice          = 2000000000
	baseTimestamp      = 1000000000
)

type feed struct {
	feedID    [32]byte
	feedIDStr string
	reports   []report
}

type report struct {
	rawReport  []byte
	reportCtx  []byte
	signatures [][]byte
}

type node struct {
	peerID p2ptypes.PeerID
	bundle ocr2key.KeyBundle
}

// Integration/load test that combines Trigger Subscriber, Streams Trigger Aggregator and Streams Codec.
// It measures time needed to receive and process trigger events from multiple nodes and produce a local aggregated event.
// For more meaningful measurements, increase the values of parameters P and T.
func TestStreamsTrigger(t *testing.T) {
	N := 31 // trigger DON nodes
	F := 10 // faulty nodes
	R := 5  // different reports per feed (i.e. prices and timestamps)
	P := 2  // feeds
	T := 2  // test iterations

	nodes := newNodes(t, N)
	feeds := newFeedsWithSignedReports(t, nodes, N, P, R)

	allowedSigners := make([][]byte, N)
	for i := 0; i < N; i++ {
		allowedSigners[i] = nodes[i].bundle.PublicKey() // bad name - see comment on evmKeyring.PublicKey
	}
	lggr := logger.TestLogger(t)
	ctx := testutils.Context(t)
	codec := streams.NewCodec(lggr)
	agg := triggers.NewMercuryRemoteAggregator(codec, allowedSigners, F, lggr)

	capInfo := capabilities.CapabilityInfo{
		ID: triggerID,
	}
	capMembers := make([]p2ptypes.PeerID, N)
	for i := 0; i < N; i++ {
		capMembers[i] = nodes[i].peerID
	}
	capDonInfo := capabilities.DON{
		Members: capMembers,
		F:       uint8(F),
	}
	config := &capabilities.RemoteTriggerConfig{
		MinResponsesToAggregate: uint32(F + 1),
	}
	subscriber := remote.NewTriggerSubscriber(config, capInfo, capDonInfo, capabilities.DON{}, nil, agg, lggr)

	// register trigger
	req := capabilities.TriggerRegistrationRequest{
		Metadata: capabilities.RequestMetadata{
			WorkflowID: workflowID,
		},
	}
	triggerEventCallbackCh, err := subscriber.RegisterTrigger(ctx, req)
	require.NoError(t, err)

	// send and process all trigger events
	startTs := time.Now().UnixMilli()
	processingTime := int64(0)
	for c := 0; c < T; c++ {
		triggerEventID := baseTriggerEventID + strconv.Itoa(c)
		for i := 0; i < N; i++ { // every node ...
			reportList := make([]datastreams.FeedReport, P)
			for j := 0; j < P; j++ { //  ... sends reports for every feed ...
				reportIdx := (i + j) % R
				signatures := make([][]byte, F+1)
				for k := 0; k < F+1; k++ { // ... each signed by F+1 nodes
					signatures[k] = feeds[j].reports[reportIdx].signatures[(i+k)%N]
				}
				signedStreamsReport := datastreams.FeedReport{
					FeedID:        feeds[j].feedIDStr,
					FullReport:    feeds[j].reports[reportIdx].rawReport,
					ReportContext: feeds[j].reports[reportIdx].reportCtx,
					Signatures:    signatures,
				}
				reportList[j] = signedStreamsReport
			}

			msg := newTriggerEvent(t, reportList, triggerEventID, nodes[i].peerID)

			processingStart := time.Now().UnixMilli()
			subscriber.Receive(ctx, msg)
			processingTime += time.Now().UnixMilli() - processingStart
		}

		response := <-triggerEventCallbackCh
		validateLatestReports(t, response.Event.Outputs, P, basePrice+R-1, baseTimestamp+R-1)
	}
	totalTime := time.Now().UnixMilli() - startTs
	lggr.Infow("elapsed", "totalMs", totalTime, "processingMs", processingTime)
}

func newNodes(t *testing.T, N int) []node {
	nodes := make([]node, N)
	for i := 0; i < N; i++ {
		bundle, err := ocr2key.New(chaintype.EVM)
		require.NoError(t, err)
		nodes[i].bundle = bundle
		nodes[i].peerID = newPeerID(t)
	}
	return nodes
}

func newPeerID(t *testing.T) ragetypes.PeerID {
	_, privKey, err := ed25519.GenerateKey(rand.Reader)
	require.NoError(t, err)
	peerID, err := ragetypes.PeerIDFromPrivateKey(privKey)
	require.NoError(t, err)
	return peerID
}

func newFeedsWithSignedReports(t *testing.T, nodes []node, N, P, R int) []feed {
	feeds := make([]feed, P)
	for i := 0; i < P; i++ {
		id, idStr := newFeedID(t)
		feeds[i].feedID = id
		feeds[i].feedIDStr = idStr
		feeds[i].reports = make([]report, R)
		for j := 0; j < R; j++ {
			report := newReport(t, id, big.NewInt(int64(basePrice+j)), int64(baseTimestamp+j))
			feeds[i].reports[j].rawReport = report
			reportCtx := ocrTypes.ReportContext{ReportTimestamp: ocrTypes.ReportTimestamp{Epoch: uint32(baseTimestamp + j)}}
			feeds[i].reports[j].reportCtx = rawReportContext(reportCtx)
			feeds[i].reports[j].signatures = make([][]byte, N)
			for k := 0; k < N; k++ {
				signature, err := nodes[k].bundle.Sign(reportCtx, report)
				require.NoError(t, err)
				feeds[i].reports[j].signatures[k] = signature
			}
		}
	}
	return feeds
}

func newTriggerEvent(t *testing.T, reportList []datastreams.FeedReport, triggerEventID string, sender ragetypes.PeerID) *remotetypes.MessageBody {
	outputs, err := values.WrapMap(&datastreams.StreamsTriggerEvent{
		Timestamp: 10,
		Payload:   reportList,
	})
	require.NoError(t, err)

	triggerEvent := capabilities.TriggerEvent{
		TriggerType: triggerID,
		ID:          triggerEventID,
		Outputs:     outputs,
	}

	marshaled, err := pb.MarshalTriggerResponse(capabilities.TriggerResponse{Event: triggerEvent})

	require.NoError(t, err)
	msg := &remotetypes.MessageBody{
		Sender: sender[:],
		Method: remotetypes.MethodTriggerEvent,
		Metadata: &remotetypes.MessageBody_TriggerEventMetadata{
			TriggerEventMetadata: &remotetypes.TriggerEventMetadata{
				WorkflowIds:    []string{workflowID},
				TriggerEventId: triggerEventID,
			},
		},
		Payload: marshaled,
	}
	return msg
}

func validateLatestReports(t *testing.T, wrapped values.Value, expectedFeedsLen int, expectedPrice int, expectedTimestamp int) {
	triggerEvent := datastreams.StreamsTriggerEvent{}
	require.NoError(t, wrapped.UnwrapTo(&triggerEvent))
	require.Equal(t, expectedFeedsLen, len(triggerEvent.Payload))
	priceBig := big.NewInt(int64(expectedPrice))
	for _, report := range triggerEvent.Payload {
		require.Equal(t, priceBig.Bytes(), report.BenchmarkPrice)
		require.Equal(t, int64(expectedTimestamp), report.ObservationTimestamp)
	}
}
