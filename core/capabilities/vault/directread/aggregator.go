// Package directread aggregates Vault direct GetSecrets responses on the workflow DON.
package directread

import (
	"errors"
	"fmt"
	"strings"

	"google.golang.org/protobuf/types/known/anypb"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	caperrors "github.com/smartcontractkit/chainlink-common/pkg/capabilities/errors"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/executable/request"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

var ErrNotEnoughMatchingResponses = errors.New("not enough matching vault responses")

func NewAggregatorFactory() request.AggregatorFactory {
	return func(req commoncap.CapabilityRequest, remoteDON commoncap.DON) request.ResponseAggregator {
		if req.Method != vaulttypes.MethodSecretsGet || req.Payload == nil {
			return nil
		}
		gsr := &vaultcommon.GetSecretsRequest{}
		if err := req.Payload.UnmarshalTo(gsr); err != nil || !gsr.GetSecretsDirectly {
			return nil
		}
		return NewAggregator(gsr, len(remoteDON.Members), int(remoteDON.F))
	}
}

// Aggregator can't verify shares (they're encrypted to the workflow nodes), so it
// accepts a secret once 2F+1 nodes return the same identifier, ciphertext and
// vault public key
type Aggregator struct {
	numSecrets int
	n          int
	f          int
	quorum     int

	replies     int
	peerErrors  map[string]int
	secretVotes []map[string][]*vaultcommon.SecretResponse // per secret, by voteKey
	votePKs     map[string]string                          // voteKey -> RawVaultPublicKey
}

var _ request.ResponseAggregator = (*Aggregator)(nil)

func NewAggregator(req *vaultcommon.GetSecretsRequest, n, f int) *Aggregator {
	votes := make([]map[string][]*vaultcommon.SecretResponse, len(req.Requests))
	for i := range votes {
		votes[i] = map[string][]*vaultcommon.SecretResponse{}
	}
	return &Aggregator{
		numSecrets:  len(req.Requests),
		n:           n,
		f:           f,
		quorum:      2*f + 1,
		peerErrors:  map[string]int{},
		secretVotes: votes,
		votePKs:     map[string]string{},
	}
}

func (a *Aggregator) OnResponse(_ p2ptypes.PeerID, resp commoncap.CapabilityResponse) (*commoncap.CapabilityResponse, error) {
	gsr := &vaultcommon.GetSecretsResponse{}
	if resp.Payload == nil {
		return a.OnError(p2ptypes.PeerID{}, "empty response payload")
	}
	if err := resp.Payload.UnmarshalTo(gsr); err != nil {
		return a.OnError(p2ptypes.PeerID{}, "failed to unmarshal GetSecretsResponse: "+err.Error())
	}
	// Vault nodes answer each requested secret in request order.
	if len(gsr.Responses) != a.numSecrets {
		return a.OnError(p2ptypes.PeerID{}, fmt.Sprintf("expected %d secret responses, got %d", a.numSecrets, len(gsr.Responses)))
	}

	a.replies++
	for i, sr := range gsr.Responses {
		k := voteKey(sr, gsr.RawVaultPublicKey)
		a.secretVotes[i][k] = append(a.secretVotes[i][k], sr)
		if sr.GetData() != nil {
			a.votePKs[k] = gsr.RawVaultPublicKey
		}
	}
	return a.decide()
}

func (a *Aggregator) OnError(_ p2ptypes.PeerID, errMsg string) (*commoncap.CapabilityResponse, error) {
	a.replies++
	a.peerErrors[errMsg]++
	return a.decide()
}

// Shares only combine within one DKG instance, so data votes include the key.
func voteKey(sr *vaultcommon.SecretResponse, publicKey string) string {
	id := "<nil>"
	if sr.GetId() != nil {
		id = vaulttypes.KeyFor(sr.GetId())
	}
	if sr.GetData() != nil {
		return id + "|data|" + publicKey + "|" + sr.GetData().GetEncryptedValue()
	}
	return id + "|error|" + sr.GetError()
}

func (a *Aggregator) decide() (*commoncap.CapabilityResponse, error) {
	pending := a.n - a.replies
	winners := make([]string, a.numSecrets)
	undecided := false
	for i, votes := range a.secretVotes {
		best, bestKey := 0, ""
		for k, rs := range votes {
			if len(rs) > best {
				best, bestKey = len(rs), k
			}
		}
		switch {
		case best >= a.quorum:
			winners[i] = bestKey
		case best+pending >= a.quorum:
			undecided = true
		default:
			return nil, a.failure(i)
		}
	}
	if undecided || (a.numSecrets == 0 && a.replies < a.quorum) {
		return nil, nil
	}
	return a.merge(winners)
}

func (a *Aggregator) failure(secret int) error {
	for msg, count := range a.peerErrors {
		if count >= a.f+1 {
			// F+1 identical errors include at least one honest node.
			return caperrors.DeserializeErrorFromString(msg)
		}
	}
	if len(a.secretVotes[secret]) > 1 {
		return fmt.Errorf("%w: %d different answers for secret %d from %d replies", vaultcommon.ErrSecretVersionSkew, len(a.secretVotes[secret]), secret, a.replies)
	}
	return fmt.Errorf("%w: need %d, have %d replies with %d errors (%s)", ErrNotEnoughMatchingResponses, a.quorum, a.replies, a.errorCount(), a.errorSummary())
}

func (a *Aggregator) errorCount() int {
	n := 0
	for _, c := range a.peerErrors {
		n += c
	}
	return n
}

func (a *Aggregator) errorSummary() string {
	parts := make([]string, 0, len(a.peerErrors))
	for msg, c := range a.peerErrors {
		parts = append(parts, fmt.Sprintf("%dx %q", c, msg))
	}
	return strings.Join(parts, ", ")
}

// merge matches the OCR path's shape: one entry per encryption key with every node's share.
func (a *Aggregator) merge(winners []string) (*commoncap.CapabilityResponse, error) {
	out := &vaultcommon.GetSecretsResponse{Responses: make([]*vaultcommon.SecretResponse, 0, a.numSecrets)}
	for i, k := range winners {
		if pk, ok := a.votePKs[k]; ok {
			out.RawVaultPublicKey = pk
		}
		out.Responses = append(out.Responses, mergeSecretResponses(a.secretVotes[i][k]))
	}
	payload, err := anypb.New(out)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal aggregated GetSecretsResponse: %w", err)
	}
	return &commoncap.CapabilityResponse{Payload: payload}, nil
}

func mergeSecretResponses(rs []*vaultcommon.SecretResponse) *vaultcommon.SecretResponse {
	first := rs[0]
	if first.GetData() == nil {
		return &vaultcommon.SecretResponse{Id: first.GetId(), Result: &vaultcommon.SecretResponse_Error{Error: first.GetError()}}
	}

	var entries []*vaultcommon.EncryptedShares
	byKey := map[string]*vaultcommon.EncryptedShares{}
	for _, r := range rs {
		for _, es := range r.GetData().GetEncryptedDecryptionKeyShares() {
			merged, ok := byKey[es.GetEncryptionKey()]
			if !ok {
				merged = &vaultcommon.EncryptedShares{EncryptionKey: es.GetEncryptionKey()}
				byKey[es.GetEncryptionKey()] = merged
				entries = append(entries, merged)
			}
			merged.BinaryShares = append(merged.BinaryShares, es.GetBinaryShares()...)
			merged.Shares = append(merged.Shares, es.GetShares()...)
		}
	}

	return &vaultcommon.SecretResponse{
		Id: first.GetId(),
		Result: &vaultcommon.SecretResponse_Data{Data: &vaultcommon.SecretData{
			EncryptedValue:               first.GetData().GetEncryptedValue(),
			EncryptedDecryptionKeyShares: entries,
		}},
	}
}
