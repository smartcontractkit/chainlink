package confidentialrelay

import (
	"encoding/json"
	"errors"
	"fmt"

	relaytypes "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/actions/confidentialrelay"
	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

var errUnknownMethod = errors.New("unknown relay method")

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

// Bundle builds the SignedXResponseBundle envelope from the per-node responses
// collected so far. Transport-level JSON-RPC errors, nil results, and responses
// that do not decode as a signed relay response are skipped (they carry no usable
// signature). Order is not significant: the enclave groups the bundle by response
// hash and verifies each signature independently. Returns the encoded bundle
// response and the count of signed responses it contains.
func (b *bundler) Bundle(req jsonrpc.Request[json.RawMessage], resps map[string]jsonrpc.Response[json.RawMessage], l logger.Logger) (*jsonrpc.Response[json.RawMessage], int, error) {
	switch req.Method {
	case relaytypes.MethodSecretsGet:
		out := make([]relaytypes.SignedSecretsResponseResult, 0, len(resps))
		for nodeAddr, r := range resps {
			if r.Error != nil || r.Result == nil {
				continue
			}
			var signed relaytypes.SignedSecretsResponseResult
			if err := json.Unmarshal(*r.Result, &signed); err != nil {
				l.Warnw("skipping undecodable secrets response from node", "nodeAddr", nodeAddr, "error", err)
				continue
			}
			out = append(out, signed)
		}
		encoded, err := json.Marshal(relaytypes.SignedSecretsResponseBundle{Responses: out})
		if err != nil {
			return nil, 0, fmt.Errorf("failed to encode secrets response bundle: %w", err)
		}
		return wrapBundle(req, encoded), len(out), nil
	case relaytypes.MethodCapabilityExec:
		out := make([]relaytypes.SignedCapabilityResponseResult, 0, len(resps))
		for nodeAddr, r := range resps {
			if r.Error != nil || r.Result == nil {
				continue
			}
			var signed relaytypes.SignedCapabilityResponseResult
			if err := json.Unmarshal(*r.Result, &signed); err != nil {
				l.Warnw("skipping undecodable capability response from node", "nodeAddr", nodeAddr, "error", err)
				continue
			}
			out = append(out, signed)
		}
		encoded, err := json.Marshal(relaytypes.SignedCapabilityResponseBundle{Responses: out})
		if err != nil {
			return nil, 0, fmt.Errorf("failed to encode capability response bundle: %w", err)
		}
		return wrapBundle(req, encoded), len(out), nil
	default:
		return nil, 0, fmt.Errorf("%w: %q", errUnknownMethod, req.Method)
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
