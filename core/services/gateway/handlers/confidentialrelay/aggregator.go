package confidentialrelay

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"sort"

	relaytypes "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/actions/confidentialrelay"
	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

var (
	errInsufficientResponsesForQuorum = errors.New("insufficient valid responses to reach quorum")
	errQuorumUnobtainable             = errors.New("quorum unobtainable")
	errUnknownMethod                  = errors.New("unknown relay method")
)

type aggregator struct{}

// Aggregate buckets per-node signed responses by their canonical logical hash and, once F+1
// unique signers vouch for the same hash, returns a single envelope carrying the shared logical
// result with the merged signature set. Transport-level JSON-RPC errors and undecodable
// responses are skipped without counting toward quorum.
//
// Defense-in-depth: when the gateway gains access to the relay-DON signer set (e.g., via the
// capability registry or an extended NodeConfig), each incoming RelayResponseSignature should
// be verified against the logical hash here before being merged into a bucket, mirroring the
// untrusted-host signature filter in confidential-compute (host.go). Trust anchor remains the
// enclave; verification at this layer is an early filter only.
func (a *aggregator) Aggregate(req jsonrpc.Request[json.RawMessage], resps map[string]jsonrpc.Response[json.RawMessage], donF int, donMembersCount int, l logger.Logger) (*jsonrpc.Response[json.RawMessage], error) {
	// F+1 (QuorumFPlusOne) is sufficient because each relay node calls the
	// target DON (Vault or capability) through CRE's standard capability
	// dispatch, which includes DON-level consensus. Every honest relay node
	// receives the same consensus-aggregated response and performs deterministic
	// translation, producing byte-identical outputs. F+1 matching responses
	// therefore guarantees at least one honest node vouched for the result.
	requiredQuorum := donF + 1

	if len(resps) < requiredQuorum {
		return nil, errInsufficientResponsesForQuorum
	}

	buckets := map[[32]byte]*responseBucket{}
	maxBucketSigs := 0
	contributingResponses := 0

	for nodeAddr, r := range resps {
		if r.Error != nil || r.Result == nil {
			continue
		}

		hash, sigs, raw, err := decodeSignedResponse(req, *r.Result)
		if err != nil {
			l.Warnw("failed to decode signed relay response, skipping", "nodeAddr", nodeAddr, "error", err)
			continue
		}

		b, ok := buckets[hash]
		if !ok {
			b = &responseBucket{result: raw, signers: map[string]struct{}{}}
			buckets[hash] = b
		}
		for _, sig := range sigs {
			key := string(sig.Signer)
			if _, dup := b.signers[key]; dup {
				continue
			}
			b.signers[key] = struct{}{}
			b.signatures = append(b.signatures, sig)
		}
		contributingResponses++
		if len(b.signatures) > maxBucketSigs {
			maxBucketSigs = len(b.signatures)
		}
	}

	var qualified [][32]byte
	for h, b := range buckets {
		if len(b.signatures) >= requiredQuorum {
			qualified = append(qualified, h)
		}
	}
	if len(qualified) > 0 {
		sort.Slice(qualified, func(i, j int) bool { return bytes.Compare(qualified[i][:], qualified[j][:]) < 0 })
		b := buckets[qualified[0]]
		merged, err := encodeMergedResponse(req, b.result, b.signatures)
		if err != nil {
			return nil, fmt.Errorf("failed to encode merged signed response: %w", err)
		}
		return &jsonrpc.Response[json.RawMessage]{
			Version: req.Version,
			ID:      req.ID,
			Method:  req.Method,
			Result:  &merged,
		}, nil
	}

	remainingResponses := donMembersCount - len(resps)
	if maxBucketSigs+remainingResponses < requiredQuorum {
		l.Warnw("quorum unattainable for request",
			"requiredQuorum", requiredQuorum,
			"remainingResponses", remainingResponses,
			"maxBucketSigs", maxBucketSigs,
			"contributingResponses", contributingResponses,
		)
		return nil, fmt.Errorf("%w: requiredQuorum=%d, maxBucketSigs=%d, remainingResponses=%d", errQuorumUnobtainable, requiredQuorum, maxBucketSigs, remainingResponses)
	}
	return nil, errInsufficientResponsesForQuorum
}

type responseBucket struct {
	// result holds the typed logical result for this hash, kept around so the merged envelope
	// can be re-encoded canonically without depending on byte-for-byte equality of the input.
	result     interface{}
	signatures []relaytypes.RelayResponseSignature
	signers    map[string]struct{}
}

// decodeSignedResponse unmarshals one node's signed response, computes the canonical logical
// hash for the request, and returns the per-node signatures plus the typed logical result for
// later re-encoding.
func decodeSignedResponse(req jsonrpc.Request[json.RawMessage], rawResult json.RawMessage) ([32]byte, []relaytypes.RelayResponseSignature, interface{}, error) {
	switch req.Method {
	case relaytypes.MethodSecretsGet:
		var params relaytypes.SecretsRequestParams
		if req.Params != nil {
			if err := json.Unmarshal(*req.Params, &params); err != nil {
				return [32]byte{}, nil, nil, fmt.Errorf("decode secrets request params: %w", err)
			}
		}
		var signed relaytypes.SignedSecretsResponseResult
		if err := json.Unmarshal(rawResult, &signed); err != nil {
			return [32]byte{}, nil, nil, fmt.Errorf("decode signed secrets response: %w", err)
		}
		hash, err := signed.Result.Hash(params)
		if err != nil {
			return [32]byte{}, nil, nil, fmt.Errorf("hash secrets response: %w", err)
		}
		return hash, signed.Signatures, signed.Result, nil
	case relaytypes.MethodCapabilityExec:
		var params relaytypes.CapabilityRequestParams
		if req.Params != nil {
			if err := json.Unmarshal(*req.Params, &params); err != nil {
				return [32]byte{}, nil, nil, fmt.Errorf("decode capability request params: %w", err)
			}
		}
		var signed relaytypes.SignedCapabilityResponseResult
		if err := json.Unmarshal(rawResult, &signed); err != nil {
			return [32]byte{}, nil, nil, fmt.Errorf("decode signed capability response: %w", err)
		}
		hash, err := signed.Result.Hash(params)
		if err != nil {
			return [32]byte{}, nil, nil, fmt.Errorf("hash capability response: %w", err)
		}
		return hash, signed.Signatures, signed.Result, nil
	default:
		return [32]byte{}, nil, nil, fmt.Errorf("%w: %q", errUnknownMethod, req.Method)
	}
}

// encodeMergedResponse builds the final aggregated envelope from the shared logical result and
// the merged signature set. Signatures are sorted by signer for determinism.
func encodeMergedResponse(req jsonrpc.Request[json.RawMessage], result interface{}, sigs []relaytypes.RelayResponseSignature) (json.RawMessage, error) {
	sort.Slice(sigs, func(i, j int) bool { return bytes.Compare(sigs[i].Signer, sigs[j].Signer) < 0 })

	switch req.Method {
	case relaytypes.MethodSecretsGet:
		typed, ok := result.(relaytypes.SecretsResponseResult)
		if !ok {
			return nil, fmt.Errorf("internal: expected SecretsResponseResult, got %T", result)
		}
		return json.Marshal(relaytypes.SignedSecretsResponseResult{Result: typed, Signatures: sigs})
	case relaytypes.MethodCapabilityExec:
		typed, ok := result.(relaytypes.CapabilityResponseResult)
		if !ok {
			return nil, fmt.Errorf("internal: expected CapabilityResponseResult, got %T", result)
		}
		return json.Marshal(relaytypes.SignedCapabilityResponseResult{Result: typed, Signatures: sigs})
	default:
		return nil, fmt.Errorf("%w: %q", errUnknownMethod, req.Method)
	}
}
