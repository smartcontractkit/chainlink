package confidentialrelay

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"

	relaytypes "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/actions/confidentialrelay"
	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

var errUnknownMethod = errors.New("unknown relay method")

// maxLoggedNodeErrors caps how many distinct node error samples we surface in
// structured logs.
const maxLoggedNodeErrors = 5

// maxNodeErrorMessageLen is the max length of node error messages; longer
// messages are truncated.
const maxNodeErrorMessageLen = 200

// bundler assembles the gateway's response to the enclave. The gateway is a dumb
// fan-in: it forwards every per-node signed response it collected, in one bundle,
// without verifying signatures, checking signer membership, or deciding quorum.
//
// The enclave is the sole trust anchor. It groups the bundle by response hash,
// verifies each signature against the relay-DON signer set, and accepts the result
// backed by F+1 valid distinct signers. Because responses are keyed upstream by the
// gateway-authenticated node address (one response per node), a single compromised
// relay node can contribute at most its own one response to the bundle. It cannot
// manufacture quorum by fabricating signer identities; the worst it can do is add
// noise the enclave discards. This is what closes the prior liveness gap, where the
// gateway counted unverified, attacker-controlled signer identities toward quorum
// and could commit to a forged response, starving the enclave of the real one.
type bundler struct{}

// BundleSummary is the fan-in accounting for one Bundle call. Signed is the only
// count that can contribute to enclave quorum; the rest is diagnostic.
type BundleSummary struct {
	response    *jsonrpc.Response[json.RawMessage]
	signed      int
	errorCount  int
	undecodable int
	total       int
	nodeErrors  map[string]string
}

func (s *BundleSummary) Response() *jsonrpc.Response[json.RawMessage] { return s.response }
func (s *BundleSummary) Signed() int                                  { return s.signed }
func (s *BundleSummary) Error() int                                   { return s.errorCount }
func (s *BundleSummary) Undecodable() int                             { return s.undecodable }
func (s *BundleSummary) Total() int                                   { return s.total }

// NodeErrorsFormatted returns a stable, compact sample of node error messages
// for structured logs.
func (s *BundleSummary) NodeErrorsFormatted() string {
	if s == nil || len(s.nodeErrors) == 0 {
		return ""
	}
	addrs := make([]string, 0, len(s.nodeErrors))
	for addr := range s.nodeErrors {
		addrs = append(addrs, addr)
	}
	sort.Strings(addrs)
	parts := make([]string, 0, len(addrs))
	for _, addr := range addrs {
		parts = append(parts, addr+"="+s.nodeErrors[addr])
	}
	return strings.Join(parts, "; ")
}

func (s *BundleSummary) addError(nodeAddr string, wireErr *jsonrpc.WireError) {
	s.errorCount++
	if s.nodeErrors == nil {
		s.nodeErrors = make(map[string]string)
	}
	if len(s.nodeErrors) >= maxLoggedNodeErrors {
		return
	}
	if _, exists := s.nodeErrors[nodeAddr]; exists {
		return
	}
	msg := "nil result"
	if wireErr != nil {
		msg = sanitizeNodeErrorMessage(wireErr.Message)
		if msg == "" {
			msg = fmt.Sprintf("jsonrpc error code %d", wireErr.Code)
		}
	}
	s.nodeErrors[nodeAddr] = msg
}

func (s *BundleSummary) addUndecodable() { s.undecodable++ }

func (s *BundleSummary) setSignedResponse(resp *jsonrpc.Response[json.RawMessage], signed int) {
	s.response = resp
	s.signed = signed
}

func newBundleSummary(total int) *BundleSummary {
	return &BundleSummary{
		total:      total,
		nodeErrors: make(map[string]string),
	}
}

// Bundle builds the SignedXResponseBundle envelope from the per-node responses
// collected so far. Transport-level JSON-RPC errors, nil results, and responses
// that do not decode as a signed relay response are skipped (they carry no usable
// signature). Order is not significant: the enclave groups the bundle by response
// hash and verifies each signature independently.
func (b *bundler) Bundle(req jsonrpc.Request[json.RawMessage], resps map[string]jsonrpc.Response[json.RawMessage], l logger.Logger) (*BundleSummary, error) {
	summary := newBundleSummary(len(resps))

	switch req.Method {
	case relaytypes.MethodSecretsGet:
		out := make([]relaytypes.SignedSecretsResponseResult, 0, len(resps))
		for nodeAddr, r := range resps {
			if r.Error != nil || r.Result == nil {
				summary.addError(nodeAddr, r.Error)
				continue
			}
			var signed relaytypes.SignedSecretsResponseResult
			if err := json.Unmarshal(*r.Result, &signed); err != nil {
				summary.addUndecodable()
				l.Warnw("skipping undecodable secrets response from node", "nodeAddr", nodeAddr, "error", err)
				continue
			}
			out = append(out, signed)
		}
		encoded, err := json.Marshal(relaytypes.SignedSecretsResponseBundle{Responses: out})
		if err != nil {
			return nil, fmt.Errorf("failed to encode secrets response bundle: %w", err)
		}
		summary.setSignedResponse(wrapBundle(req, encoded), len(out))
		return summary, nil
	case relaytypes.MethodCapabilityExec:
		out := make([]relaytypes.SignedCapabilityResponseResult, 0, len(resps))
		for nodeAddr, r := range resps {
			if r.Error != nil || r.Result == nil {
				summary.addError(nodeAddr, r.Error)
				continue
			}
			var signed relaytypes.SignedCapabilityResponseResult
			if err := json.Unmarshal(*r.Result, &signed); err != nil {
				summary.addUndecodable()
				l.Warnw("skipping undecodable capability response from node", "nodeAddr", nodeAddr, "error", err)
				continue
			}
			out = append(out, signed)
		}
		encoded, err := json.Marshal(relaytypes.SignedCapabilityResponseBundle{Responses: out})
		if err != nil {
			return nil, fmt.Errorf("failed to encode capability response bundle: %w", err)
		}
		summary.setSignedResponse(wrapBundle(req, encoded), len(out))
		return summary, nil
	default:
		return nil, fmt.Errorf("%w: %q", errUnknownMethod, req.Method)
	}
}

func wrapBundle(req jsonrpc.Request[json.RawMessage], encoded json.RawMessage) *jsonrpc.Response[json.RawMessage] {
	return &jsonrpc.Response[json.RawMessage]{
		Version: req.Version,
		ID:      req.ID,
		Method:  req.Method,
		Result:  &encoded,
	}
}

func sanitizeNodeErrorMessage(msg string) string {
	msg = strings.TrimSpace(msg)
	if len(msg) <= maxNodeErrorMessageLen {
		return msg
	}
	return msg[:maxNodeErrorMessageLen-1] + ".."
}
