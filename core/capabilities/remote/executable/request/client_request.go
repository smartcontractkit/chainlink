package request

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/ocr2key"
	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	caperrors "github.com/smartcontractkit/chainlink-common/pkg/capabilities/errors"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/validation"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

// errRemoteCapabilityExecuteError preserves the legacy "TRANSPORT : ErrorMsg" string from the
// remote executable client while wrapping a deserialized caperrors.Error so callers can
// errors.As into caperrors.Error after RPC (see capability_executor metrics).
type errRemoteCapabilityExecuteError struct {
	s    string
	wrap caperrors.Error
}

func (e *errRemoteCapabilityExecuteError) Error() string { return e.s }

func (e *errRemoteCapabilityExecuteError) Unwrap() error { return e.wrap }

func newRemoteCapabilityExecuteError(transportErr types.Error, errMsg string) error {
	return &errRemoteCapabilityExecuteError{
		s:    fmt.Sprintf("%s : %s", transportErr, errMsg),
		wrap: caperrors.DeserializeErrorFromString(errMsg),
	}
}

func newRemoteCapabilityExecuteErrorWithMessage(display, errMsg string) error {
	return &errRemoteCapabilityExecuteError{
		s:    display,
		wrap: caperrors.DeserializeErrorFromString(errMsg),
	}
}

// newRemoteCapabilityExecuteErrorFromCapError wraps a capability error produced locally
// (not deserialized from a remote payload). Unwrap returns caperrors.Error so
// capability_executor can errors.As into caperrors.Error.
func newRemoteCapabilityExecuteErrorFromCapError(capErr caperrors.Error) error {
	return &errRemoteCapabilityExecuteError{
		s:    capErr.Error(),
		wrap: capErr,
	}
}

// ErrResponseQuorumUnreachable is returned when enough peer responses have been
// received that the required number of matching payloads (configurable; defaults to F+1) are no longer mathematically possible.
var ErrResponseQuorumUnreachable = errors.New("response quorum unreachable: not enough matching capability responses")

type clientResponse struct {
	Result []byte
	Err    error
}

// messageFields are the fields of the outgoing MessageBody that the caller controls.
type messageFields struct {
	CapabilityID     string
	CapabilityDonID  uint32
	CallerDonID      uint32
	Method           string
	Payload          []byte
	MessageID        []byte
	CapabilityMethod string
}

type ClientRequest struct {
	id                       string
	dispatcher               types.Dispatcher
	peers                    []p2ptypes.PeerID
	message                  messageFields
	responseCh               chan clientResponse
	createdAt                time.Time
	responseIDCount          map[[32]byte]int
	meteringResponses        map[[32]byte][]commoncap.MeteringNodeDetail
	errorCount               map[string]int
	totalErrorCount          int
	payloadNotAvailableCount int
	responseReceived         map[p2ptypes.PeerID]bool
	lggr                     logger.Logger
	signers                  [][]byte
	workflowExecutionID      string
	referenceID              string

	requiredResponseConfirmations int
	remoteNodeCount               int
	remoteDonF                    int

	requestTimeout time.Duration

	respSent bool
	mux      sync.Mutex
}

// NewClientExecuteRequest creates a new client request for a remote executable capability.
// The request is not sent until Send is called.
func NewClientExecuteRequest(ctx context.Context, lggr logger.Logger, req commoncap.CapabilityRequest,
	remoteCapabilityInfo commoncap.CapabilityInfo, localDonInfo commoncap.DON, dispatcher types.Dispatcher,
	requestTimeout time.Duration, capMethodName string,
	signers [][]byte, minResponsesToAggregate uint32,
) (*ClientRequest, error) {
	rawRequest, err := proto.MarshalOptions{Deterministic: true}.Marshal(pb.CapabilityRequestToProto(req))
	if err != nil {
		return nil, fmt.Errorf("failed to marshal capability request: %w", err)
	}

	workflowExecutionID := req.Metadata.WorkflowExecutionID
	if err = validation.ValidateWorkflowOrExecutionID(workflowExecutionID); err != nil {
		return nil, fmt.Errorf("workflow execution ID is invalid: %w", err)
	}

	// the requestID must be delineated by the workflow execution ID and the reference ID
	// to ensure that it supports parallel step execution
	requestID := types.MethodExecute + ":" + workflowExecutionID + ":" + req.Metadata.ReferenceID

	lggr = logger.With(lggr, "requestId", requestID) // cap ID and method name included in the parent logger

	if remoteCapabilityInfo.DON == nil {
		return nil, errors.New("remote capability info missing DON")
	}

	msg := messageFields{
		CapabilityID:     remoteCapabilityInfo.ID,
		CapabilityDonID:  remoteCapabilityInfo.DON.ID,
		CallerDonID:      localDonInfo.ID,
		Method:           types.MethodExecute,
		Payload:          rawRequest,
		MessageID:        []byte(requestID),
		CapabilityMethod: capMethodName,
	}

	return newClientRequest(lggr, requestID, remoteCapabilityInfo, dispatcher, requestTimeout, msg, workflowExecutionID, req.Metadata.ReferenceID, signers, minResponsesToAggregate)
}

var (
	// Extra time the workflow DON client waits beyond requestTimeout for P2P responses
	// after remote capability nodes hit their execution deadline.
	defaultResponseAggregationGrace = 10 * time.Second
)

// requiredConfirmations returns the minimum number of matching peer responses needed before
// the workflow DON accepts a read result. When minResponsesToAggregate is non-zero it is used
// directly; otherwise the default of F+1 is used (backward compatible).
func requiredConfirmations(f uint8, minResponsesToAggregate uint32) int {
	if minResponsesToAggregate > 0 {
		return int(minResponsesToAggregate)
	}
	return int(f) + 1
}

func newClientRequest(lggr logger.Logger, requestID string, remoteCapabilityInfo commoncap.CapabilityInfo,
	dispatcher types.Dispatcher, requestTimeout time.Duration, message messageFields,
	workflowExecutionID, stepRef string, signers [][]byte, minResponsesToAggregate uint32,
) (*ClientRequest, error) {
	remoteCapabilityDonInfo := remoteCapabilityInfo.DON
	if remoteCapabilityDonInfo == nil {
		return nil, errors.New("remote capability info missing DON")
	}

	responseReceived := make(map[p2ptypes.PeerID]bool, len(remoteCapabilityDonInfo.Members))
	for _, peerID := range remoteCapabilityDonInfo.Members {
		responseReceived[peerID] = false
	}

	return &ClientRequest{
		id:                            requestID,
		dispatcher:                    dispatcher,
		peers:                         remoteCapabilityDonInfo.Members,
		message:                       message,
		createdAt:                     time.Now(),
		requestTimeout:                requestTimeout,
		requiredResponseConfirmations: requiredConfirmations(remoteCapabilityDonInfo.F, minResponsesToAggregate),
		remoteNodeCount:               len(remoteCapabilityDonInfo.Members),
		remoteDonF:                    int(remoteCapabilityDonInfo.F),
		responseIDCount:               make(map[[32]byte]int),
		meteringResponses:             make(map[[32]byte][]commoncap.MeteringNodeDetail),
		errorCount:                    make(map[string]int),
		responseReceived:              responseReceived,
		responseCh:                    make(chan clientResponse, 1),
		lggr:                          lggr,
		signers:                       signers,
		workflowExecutionID:           workflowExecutionID,
		referenceID:                   stepRef,
	}, nil
}

// Send dispatches the request to all remote capability peers. Callers must register the
// request to receive responses before calling Send, since peers may reply immediately.
func (c *ClientRequest) Send() {
	for _, peerID := range c.peers {
		// the dispatcher mutates the message, so each peer gets its own
		message := &types.MessageBody{
			CapabilityId:     c.message.CapabilityID,
			CapabilityDonId:  c.message.CapabilityDonID,
			CallerDonId:      c.message.CallerDonID,
			Method:           c.message.Method,
			Payload:          c.message.Payload,
			MessageId:        c.message.MessageID,
			CapabilityMethod: c.message.CapabilityMethod,
		}
		if err := c.dispatcher.Send(peerID, message); err != nil {
			c.lggr.Errorw("failed to send message", "peerID", peerID, "error", err)
		}
	}
}

func (c *ClientRequest) ID() string {
	return c.id
}

func (c *ClientRequest) ResponseChan() <-chan clientResponse {
	return c.responseCh
}

func (c *ClientRequest) Expired() bool {
	return time.Since(c.createdAt) > c.requestTimeout+defaultResponseAggregationGrace
}

func (c *ClientRequest) Cancel(err error) {
	c.mux.Lock()
	defer c.mux.Unlock()
	if !c.respSent {
		c.sendResponse(clientResponse{Err: err})
	}
}

func (c *ClientRequest) OnMessage(_ context.Context, msg *types.MessageBody) error {
	c.mux.Lock()
	defer c.mux.Unlock()

	if c.respSent {
		return nil
	}

	if msg.Sender == nil {
		return errors.New("sender missing from message")
	}

	c.lggr.Debugw("OnMessage called for client request")

	sender, err := remote.ToPeerID(msg.Sender)
	if err != nil {
		return fmt.Errorf("failed to convert message sender to PeerID: %w", err)
	}

	received, expected := c.responseReceived[sender]
	if !expected {
		return fmt.Errorf("response from peer %s not expected", sender)
	}

	if received {
		return fmt.Errorf("response from peer %s already received", sender)
	}

	c.responseReceived[sender] = true

	if msg.Error == types.Error_OK {
		resp, err := pb.UnmarshalCapabilityResponse(msg.Payload)
		if err != nil {
			return fmt.Errorf("failed to unmarshal capability response: %w", err)
		}

		// metering reports per node are aggregated into a single array of values. for any single node message, the
		// metering values are extracted from the CapabilityResponse, added to an array, and the CapabilityResponse
		// is marshalled without the metering value to get the hash. each node could have a different metering value
		// which would result in different hashes. removing the metering detail allows for direct comparison of results.
		responseID, err := c.getMessageHash(resp)
		if err != nil {
			return fmt.Errorf("failed to get message hash: %w", err)
		}

		lggr := logger.With(c.lggr, "responseID", hex.EncodeToString(responseID[:]), "requiredCount", c.requiredResponseConfirmations, "peer", sender)

		nodeReports, exists := c.meteringResponses[responseID]
		if !exists {
			nodeReports = make([]commoncap.MeteringNodeDetail, 0)
		}

		rpt, err := commoncap.ExtractMeteringFromMetadata(sender, resp.Metadata)
		if err != nil {
			lggr.Warnw("invalid metering detail", "err", err)
		} else {
			nodeReports = append(nodeReports, rpt)
		}

		c.responseIDCount[responseID]++
		c.meteringResponses[responseID] = nodeReports

		if len(c.responseIDCount) > 1 {
			lggr.Warnw("received multiple unique responses for the same request", "count for responseID", len(c.responseIDCount))
		}

		if c.responseIDCount[responseID] == c.requiredResponseConfirmations || c.hasValidAttestation(resp) {
			payload, err := c.encodePayloadWithMetadata(msg, commoncap.ResponseMetadata{Metering: nodeReports})
			if err != nil {
				return fmt.Errorf("failed to encode payload with metadata: %w", err)
			}

			c.sendResponse(clientResponse{Result: payload})
		} else {
			c.trySendQuorumUnreachableError()
		}
	} else {
		c.lggr.Debugw("received error from peer", "error", msg.Error, "errorMsg", msg.ErrorMsg, "peer", sender)
		if commoncap.ErrResponsePayloadNotAvailable.Is(errors.New(msg.ErrorMsg)) {
			c.payloadNotAvailableCount++
			if c.payloadNotAvailableCount == c.remoteNodeCount-c.requiredResponseConfirmations+1 {
				// return an error to indicate unexpected state, but do not send an error as we might still receive a response with valid attestation.
				return fmt.Errorf("unexpected state: received %d payload not available responses, while max allowed is %d. This means a bug in the code, please investigate",
					c.payloadNotAvailableCount, c.remoteNodeCount-c.requiredResponseConfirmations)
			}
			return nil
		}

		c.errorCount[msg.ErrorMsg]++
		c.totalErrorCount++

		if len(c.errorCount) > 1 {
			c.lggr.Warnw("received multiple different errors for the same request", "numDifferentErrors", len(c.errorCount))
		}

		if c.errorCount[msg.ErrorMsg] == c.requiredResponseConfirmations {
			c.sendResponse(clientResponse{Err: newRemoteCapabilityExecuteError(msg.Error, msg.ErrorMsg)})
		} else if c.totalErrorCount == c.remoteNodeCount-c.requiredResponseConfirmations+1 {
			c.sendResponse(clientResponse{Err: newRemoteCapabilityExecuteErrorWithMessage(
				fmt.Sprintf("received %d errors, last error %s : %s", c.totalErrorCount, msg.Error, msg.ErrorMsg),
				msg.ErrorMsg,
			)})
		}
	}
	return nil
}

func (c *ClientRequest) responsesReceivedCount() int {
	n := 0
	for _, received := range c.responseReceived {
		if received {
			n++
		}
	}
	return n
}

func (c *ClientRequest) maxMatchingResponseCount() int {
	currentMax := 0
	for _, count := range c.responseIDCount {
		if count > currentMax {
			currentMax = count
		}
	}
	return currentMax
}

// quorumStillPossible reports whether requiredResponseConfirmations identical OK
// responses can still be reached. Each pending peer can add at most one vote to
// the current leading payload hash.
func (c *ClientRequest) quorumStillPossible(pending int) bool {
	return c.maxMatchingResponseCount()+pending >= c.requiredResponseConfirmations
}

func (c *ClientRequest) trySendQuorumUnreachableError() {
	pending := c.pending()
	if c.respSent || c.quorumStillPossible(pending) {
		return
	}
	received := c.responsesReceivedCount()
	unique := len(c.responseIDCount)
	maxMatch := c.maxMatchingResponseCount()

	detail := fmt.Errorf(
		"received %d/%d peer responses with %d unique payloads; best match count %d, need %d (%d responses pending)",
		received, c.remoteNodeCount, unique, maxMatch, c.requiredResponseConfirmations, pending,
	)
	capErr := caperrors.NewPublicSystemError(
		fmt.Errorf("%w: %w", ErrResponseQuorumUnreachable, detail),
		caperrors.ConsensusFailed,
	)
	c.lggr.Warnw("response quorum unreachable, failing early",
		"responsesReceived", received,
		"remoteNodeCount", c.remoteNodeCount,
		"uniqueResponseCount", unique,
		"maxMatchingResponseCount", maxMatch,
		"requiredResponseConfirmations", c.requiredResponseConfirmations,
		"pendingResponses", pending,
	)
	c.sendResponse(clientResponse{Err: newRemoteCapabilityExecuteErrorFromCapError(capErr)})
}

func (c *ClientRequest) pending() int {
	received := c.responsesReceivedCount()
	pending := c.remoteNodeCount - received
	return pending
}

func (c *ClientRequest) hasValidAttestation(resp commoncap.CapabilityResponse) bool {
	if resp.OCRAttestation == nil {
		return false
	}

	// report_attestation.go in libocr only ever collects F+1 signatures before declaring a report complete
	// (it stops as soon as it has one signature past the DON's F threshold). That means an OCRAttestation
	// can never carry more than F+1 signatures, regardless of how the capability DON is configured.
	// When requiredResponseConfirmations (driven by minResponsesToAggregate) is set above F+1, verifyAttestation
	// would therefore always fail and log a spurious "this is most likely a bug" error on every response. In that
	// case we deliberately skip attestation verification altogether and fall back to the identical-response quorum check instead.
	if c.requiredResponseConfirmations > c.remoteDonF+1 {
		return false
	}

	err := c.verifyAttestation(resp)
	if err != nil {
		// This is error is caused by a known bug.
		// The issue is not critical, as we can fallback to the quorum aggregation and only causes read failures
		// if F+1 nodes have original RPC response used in report and F of those nodes are malicious or down.
		// TODO: PLEX-3350 - fix validation of OCR attestation.
		c.lggr.Debugw("Attestation is present, but not valid.", "error", err)
		return false
	}

	return true
}

func (c *ClientRequest) verifyAttestation(resp commoncap.CapabilityResponse) error {
	attestation := resp.OCRAttestation
	if attestation == nil {
		return errors.New("attestation is missing")
	}

	if len(attestation.Sigs) < c.requiredResponseConfirmations {
		return fmt.Errorf("not enough signatures: got %d, need at least %d", len(attestation.Sigs), c.requiredResponseConfirmations)
	}

	if len(c.signers) < c.requiredResponseConfirmations {
		return fmt.Errorf("number of configured OCR signers is less than required confirmations: got %d, need at least %d", len(c.signers), c.requiredResponseConfirmations)
	}

	reportData, err := commoncap.ResponseToReportData(c.workflowExecutionID, c.referenceID, resp.Payload.Value, resp.Metadata)
	if err != nil {
		return fmt.Errorf("failed to convert response to report data: %w", err)
	}
	sigData := ocr2key.ReportToSigData3(attestation.ConfigDigest, attestation.SequenceNumber, reportData[:])
	signed := make([]bool, len(c.signers))
	for _, sig := range attestation.Sigs {
		if int(sig.Signer) >= len(c.signers) {
			return fmt.Errorf("invalid signer index: %d", sig.Signer)
		}

		if signed[sig.Signer] {
			return fmt.Errorf("duplicate signature from signer index: %d", sig.Signer)
		}

		if !ocr2key.EvmVerifyBlob(c.signers[sig.Signer], sigData, sig.Signature) {
			return fmt.Errorf("invalid signature from signer index: %d", sig.Signer)
		}

		signed[sig.Signer] = true
	}

	return nil
}

func (c *ClientRequest) sendResponse(response clientResponse) {
	c.responseCh <- response
	close(c.responseCh)
	c.respSent = true
	if response.Err != nil {
		c.lggr.Warnw("received error response", "error", remote.SanitizeLogString(response.Err.Error()))
		return
	}
	c.lggr.Debugw("received OK response")
}

func (c *ClientRequest) getMessageHash(msg commoncap.CapabilityResponse) ([32]byte, error) {
	// clear metadata to ensure it doesn't affect the hash, as different nodes might have different metadata (e.g. different metering values)
	// since msg is passed as value, this won't affect the original message
	msg.Metadata = commoncap.ResponseMetadata{}
	msg.OCRAttestation = nil
	payload, err := pb.MarshalCapabilityResponse(msg)
	if err != nil {
		return [32]byte{}, err
	}

	return sha256.Sum256(payload), nil
}

func (c *ClientRequest) encodePayloadWithMetadata(msg *types.MessageBody, metadata commoncap.ResponseMetadata) ([]byte, error) {
	resp, err := pb.UnmarshalCapabilityResponse(msg.Payload)
	if err != nil {
		return nil, err
	}

	resp.Metadata = metadata

	return pb.MarshalCapabilityResponse(resp)
}
