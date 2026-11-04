package request

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/ocr2key"
	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	caperrors "github.com/smartcontractkit/chainlink-common/pkg/capabilities/errors"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

// ErrResponseQuorumUnreachable is returned when enough peer responses have been
// received that the required number of matching payloads (configurable; defaults to F+1) are no longer mathematically possible.
var ErrResponseQuorumUnreachable = errors.New("response quorum unreachable: not enough matching capability responses")

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

func newRemoteCapabilityExecuteErrorWithMessage(display string, errMsg string) error {
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

// requiredConfirmations returns the minimum number of matching peer responses needed before
// the workflow DON accepts a read result. When minResponsesToAggregate is non-zero it is used
// directly; otherwise the default of F+1 is used (backward compatible).
func requiredConfirmations(f uint8, minResponsesToAggregate uint32) int {
	if minResponsesToAggregate > 0 {
		return int(minResponsesToAggregate)
	}
	return int(f) + 1
}

// quorumAggregator is the default ResponseAggregator. An OK response resolves
// the request once requiredResponseConfirmations peers return identical
// payloads (or a response carries a valid OCR attestation, checked first);
// requiredResponseConfirmations identical errors resolve the request with that
// error. Unmarshalable payloads and payload-not-available responses never
// resolve the request.
type quorumAggregator struct {
	lggr                logger.Logger
	signers             [][]byte
	workflowExecutionID string
	referenceID         string

	requiredResponseConfirmations int
	remoteNodeCount               int
	remoteDonF                    int

	// replies counts every peer reply of any kind (OK response, error,
	// payload not available, unmarshalable payload).
	replies int

	responseIDCount          map[[32]byte]int
	meteringResponses        map[[32]byte][]commoncap.MeteringNodeDetail
	errorCount               map[string]int
	totalErrorCount          int
	payloadNotAvailableCount int
}

var _ ResponseAggregator = (*quorumAggregator)(nil)

func newQuorumAggregator(lggr logger.Logger, signers [][]byte, workflowExecutionID, referenceID string, remoteNodeCount int, remoteDonF uint8, minResponsesToAggregate uint32) *quorumAggregator {
	return &quorumAggregator{
		lggr:                          lggr,
		signers:                       signers,
		workflowExecutionID:           workflowExecutionID,
		referenceID:                   referenceID,
		requiredResponseConfirmations: requiredConfirmations(remoteDonF, minResponsesToAggregate),
		remoteNodeCount:               remoteNodeCount,
		remoteDonF:                    int(remoteDonF),
		responseIDCount:               make(map[[32]byte]int),
		meteringResponses:             make(map[[32]byte][]commoncap.MeteringNodeDetail),
		errorCount:                    make(map[string]int),
	}
}

// OnUnmarshalable surfaces the error to the OnMessage caller without counting
// toward any error quorum.
func (a *quorumAggregator) OnUnmarshalable(_ p2ptypes.PeerID, err error) (*commoncap.CapabilityResponse, error) {
	a.replies++
	return nil, observationError{fmt.Errorf("failed to unmarshal capability response: %w", err)}
}

func (a *quorumAggregator) OnError(_ p2ptypes.PeerID, transportErr types.Error, errMsg string) (*commoncap.CapabilityResponse, error) {
	a.replies++

	if commoncap.ErrResponsePayloadNotAvailable.Is(errors.New(errMsg)) {
		a.payloadNotAvailableCount++
		if a.payloadNotAvailableCount == a.remoteNodeCount-a.requiredResponseConfirmations+1 {
			// return an error to indicate unexpected state, but do not send an error as we might still receive a response with valid attestation.
			return nil, observationError{fmt.Errorf("unexpected state: received %d payload not available responses, while max allowed is %d. This means a bug in the code, please investigate",
				a.payloadNotAvailableCount, a.remoteNodeCount-a.requiredResponseConfirmations)}
		}
		return nil, nil
	}

	a.errorCount[errMsg]++
	a.totalErrorCount++

	if len(a.errorCount) > 1 {
		a.lggr.Warnw("received multiple different errors for the same request", "numDifferentErrors", len(a.errorCount))
	}

	if a.errorCount[errMsg] == a.requiredResponseConfirmations {
		return nil, newRemoteCapabilityExecuteError(transportErr, errMsg)
	} else if a.totalErrorCount == a.remoteNodeCount-a.requiredResponseConfirmations+1 {
		return nil, newRemoteCapabilityExecuteErrorWithMessage(
			fmt.Sprintf("received %d errors, last error %s : %s", a.totalErrorCount, transportErr, errMsg),
			errMsg,
		)
	}
	return nil, nil
}

func (a *quorumAggregator) OnResponse(sender p2ptypes.PeerID, resp commoncap.CapabilityResponse) (*commoncap.CapabilityResponse, error) {
	a.replies++

	responseID, err := a.getMessageHash(resp)
	if err != nil {
		return nil, observationError{fmt.Errorf("failed to get message hash: %w", err)}
	}

	lggr := logger.With(a.lggr, "responseID", hex.EncodeToString(responseID[:]), "requiredCount", a.requiredResponseConfirmations, "peer", sender)

	// metering reports per node are aggregated into a single array of values. for any single node message, the
	// metering values are extracted from the CapabilityResponse, added to an array, and the CapabilityResponse
	// is marshalled without the metering value to get the hash. each node could have a different metering value
	// which would result in different hashes. removing the metering detail allows for direct comparison of results.
	nodeReports, exists := a.meteringResponses[responseID]
	if !exists {
		nodeReports = make([]commoncap.MeteringNodeDetail, 0)
	}

	rpt, err := commoncap.ExtractMeteringFromMetadata(sender, resp.Metadata)
	if err != nil {
		lggr.Warnw("invalid metering detail", "err", err)
	} else {
		nodeReports = append(nodeReports, rpt)
	}

	a.responseIDCount[responseID]++
	a.meteringResponses[responseID] = nodeReports

	if len(a.responseIDCount) > 1 {
		lggr.Warnw("received multiple unique responses for the same request", "count for responseID", len(a.responseIDCount))
	}

	if a.responseIDCount[responseID] == a.requiredResponseConfirmations || a.hasValidAttestation(resp) {
		resp.Metadata = commoncap.ResponseMetadata{Metering: nodeReports}
		return &resp, nil
	}

	if err := a.trySendQuorumUnreachableError(); err != nil {
		return nil, err
	}
	return nil, nil
}

func (a *quorumAggregator) maxMatchingResponseCount() int {
	currentMax := 0
	for _, count := range a.responseIDCount {
		if count > currentMax {
			currentMax = count
		}
	}
	return currentMax
}

// quorumStillPossible reports whether requiredResponseConfirmations identical OK
// responses can still be reached. Each pending peer can add at most one vote to
// the current leading payload hash.
func (a *quorumAggregator) quorumStillPossible(pending int) bool {
	return a.maxMatchingResponseCount()+pending >= a.requiredResponseConfirmations
}

func (a *quorumAggregator) trySendQuorumUnreachableError() error {
	pending := a.pending()
	if a.quorumStillPossible(pending) {
		return nil
	}
	received := a.replies
	unique := len(a.responseIDCount)
	maxMatch := a.maxMatchingResponseCount()

	detail := fmt.Errorf(
		"received %d/%d peer responses with %d unique payloads; best match count %d, need %d (%d responses pending)",
		received, a.remoteNodeCount, unique, maxMatch, a.requiredResponseConfirmations, pending,
	)
	capErr := caperrors.NewPublicSystemError(
		fmt.Errorf("%w: %w", ErrResponseQuorumUnreachable, detail),
		caperrors.ConsensusFailed,
	)
	a.lggr.Warnw("response quorum unreachable, failing early",
		"responsesReceived", received,
		"remoteNodeCount", a.remoteNodeCount,
		"uniqueResponseCount", unique,
		"maxMatchingResponseCount", maxMatch,
		"requiredResponseConfirmations", a.requiredResponseConfirmations,
		"pendingResponses", pending,
	)
	return newRemoteCapabilityExecuteErrorFromCapError(capErr)
}

func (a *quorumAggregator) pending() int {
	return a.remoteNodeCount - a.replies
}

func (a *quorumAggregator) hasValidAttestation(resp commoncap.CapabilityResponse) bool {
	if resp.OCRAttestation == nil {
		return false
	}

	// report_attestation.go in libocr only ever collects F+1 signatures before declaring a report complete
	// (it stops as soon as it has one signature past the DON's F threshold). That means an OCRAttestation
	// can never carry more than F+1 signatures, regardless of how the capability DON is configured.
	// When requiredResponseConfirmations (driven by minResponsesToAggregate) is set above F+1, verifyAttestation
	// would therefore always fail and log a spurious "this is most likely a bug" error on every response. In that
	// case we deliberately skip attestation verification altogether and fall back to the identical-response quorum check instead.
	if a.requiredResponseConfirmations > a.remoteDonF+1 {
		return false
	}

	err := a.verifyAttestation(resp)
	if err != nil {
		// This is error is caused by a known bug.
		// The issue is not critical, as we can fallback to the quorum aggregation and only causes read failures
		// if F+1 nodes have original RPC response used in report and F of those nodes are malicious or down.
		// TODO: PLEX-3350 - fix validation of OCR attestation.
		a.lggr.Debugw("Attestation is present, but not valid.", "error", err)
		return false
	}

	return true
}

func (a *quorumAggregator) verifyAttestation(resp commoncap.CapabilityResponse) error {
	attestation := resp.OCRAttestation
	if attestation == nil {
		return errors.New("attestation is missing")
	}

	if len(attestation.Sigs) < a.requiredResponseConfirmations {
		return fmt.Errorf("not enough signatures: got %d, need at least %d", len(attestation.Sigs), a.requiredResponseConfirmations)
	}

	if len(a.signers) < a.requiredResponseConfirmations {
		return fmt.Errorf("number of configured OCR signers is less than required confirmations: got %d, need at least %d", len(a.signers), a.requiredResponseConfirmations)
	}

	reportData, err := commoncap.ResponseToReportData(a.workflowExecutionID, a.referenceID, resp.Payload.Value, resp.Metadata)
	if err != nil {
		return fmt.Errorf("failed to convert response to report data: %w", err)
	}
	sigData := ocr2key.ReportToSigData3(attestation.ConfigDigest, attestation.SequenceNumber, reportData[:])
	signed := make([]bool, len(a.signers))
	for _, sig := range attestation.Sigs {
		if int(sig.Signer) >= len(a.signers) {
			return fmt.Errorf("invalid signer index: %d", sig.Signer)
		}

		if signed[sig.Signer] {
			return fmt.Errorf("duplicate signature from signer index: %d", sig.Signer)
		}

		if !ocr2key.EvmVerifyBlob(a.signers[sig.Signer], sigData, sig.Signature) {
			return fmt.Errorf("invalid signature from signer index: %d", sig.Signer)
		}

		signed[sig.Signer] = true
	}

	return nil
}

func (a *quorumAggregator) getMessageHash(msg commoncap.CapabilityResponse) ([32]byte, error) {
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
