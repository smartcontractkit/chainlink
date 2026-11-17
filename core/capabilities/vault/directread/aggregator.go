// Package directread aggregates Vault direct GetSecrets responses on the workflow DON.
package directread

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"google.golang.org/protobuf/types/known/anypb"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	caperrors "github.com/smartcontractkit/chainlink-common/pkg/capabilities/errors"
	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/contexts"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/executable/request"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

var ErrNotEnoughMatchingResponses = errors.New("not enough matching vault responses")

func NewAggregatorFactory(shareSizeLimit limits.BoundLimiter[commonconfig.Size]) request.AggregatorFactory {
	return func(req commoncap.CapabilityRequest, remoteDON commoncap.DON) request.ResponseAggregator {
		if req.Method != vaulttypes.MethodSecretsGet || req.Payload == nil {
			return nil
		}
		gsr := &vaultcommon.GetSecretsRequest{}
		if err := req.Payload.UnmarshalTo(gsr); err != nil || !gsr.GetSecretsDirectly {
			return nil
		}
		return NewAggregator(gsr, len(remoteDON.Members), int(remoteDON.F), shareSizeLimit)
	}
}

// Aggregator can't verify shares (they're encrypted to the workflow nodes), so it
// accepts a secret once 2F+1 nodes return the same identifier, ciphertext and
// vault public key
type Aggregator struct {
	requests []*vaultcommon.SecretRequest // per requested secret, from the original request

	numSecrets int
	n          int
	f          int
	quorum     int

	// shareSizeLimit caps the size of each decryption share in a response.
	shareSizeLimit limits.BoundLimiter[commonconfig.Size]

	replies     int
	peerErrors  map[string]int
	secretVotes []map[string][]*vaultcommon.SecretResponse // per secret, by voteKey
	votePKs     map[string]string                          // voteKey -> RawVaultPublicKey
}

var _ request.ResponseAggregator = (*Aggregator)(nil)

func NewAggregator(req *vaultcommon.GetSecretsRequest, n, f int, shareSizeLimit limits.BoundLimiter[commonconfig.Size]) *Aggregator {
	votes := make([]map[string][]*vaultcommon.SecretResponse, len(req.Requests))
	for i := range votes {
		votes[i] = map[string][]*vaultcommon.SecretResponse{}
	}
	requests := make([]*vaultcommon.SecretRequest, len(req.Requests))
	copy(requests, req.Requests)
	return &Aggregator{
		requests:   requests,
		numSecrets: len(req.Requests),
		n:          n,
		f:          f,
		quorum:     2*f + 1,

		shareSizeLimit: shareSizeLimit,

		peerErrors:  map[string]int{},
		secretVotes: votes,
		votePKs:     map[string]string{},
	}
}

func (a *Aggregator) OnResponse(_ p2ptypes.PeerID, resp commoncap.CapabilityResponse) (*commoncap.CapabilityResponse, error) {
	gsr := &vaultcommon.GetSecretsResponse{}
	if resp.Payload == nil {
		return a.OnError(p2ptypes.PeerID{}, types.Error_INTERNAL_ERROR, "empty response payload")
	}
	if err := resp.Payload.UnmarshalTo(gsr); err != nil {
		return a.OnError(p2ptypes.PeerID{}, types.Error_INTERNAL_ERROR, "failed to unmarshal GetSecretsResponse: "+err.Error())
	}
	// Vault nodes answer each requested secret in request order.
	if len(gsr.Responses) != a.numSecrets {
		return a.OnError(p2ptypes.PeerID{}, types.Error_INTERNAL_ERROR, fmt.Sprintf("expected %d secret responses, got %d", a.numSecrets, len(gsr.Responses)))
	}
	for i, sr := range gsr.Responses {
		if !idMatchesRequest(a.requests[i], sr.GetId()) {
			return a.OnError(p2ptypes.PeerID{}, types.Error_INTERNAL_ERROR, fmt.Sprintf("secret %d: response identifier does not match the requested identifier", i))
		}
		if err := a.validateSecretResponseShares(a.requests[i], sr); err != nil {
			return a.OnError(p2ptypes.PeerID{}, types.Error_INTERNAL_ERROR, fmt.Sprintf("secret %d: invalid shares in response: %s", i, err))
		}
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

func (a *Aggregator) OnError(_ p2ptypes.PeerID, _ types.Error, errMsg string) (*commoncap.CapabilityResponse, error) {
	a.replies++
	a.peerErrors[errMsg]++
	return a.decide()
}

// OnUnmarshalable counts the failure as an error reply so the aggregator can
// conclude once every peer has replied.
func (a *Aggregator) OnUnmarshalable(_ p2ptypes.PeerID, err error) (*commoncap.CapabilityResponse, error) {
	return a.OnError(p2ptypes.PeerID{}, types.Error_INTERNAL_ERROR, "failed to unmarshal capability response: "+err.Error())
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

func idMatchesRequest(req *vaultcommon.SecretRequest, id *vaultcommon.SecretIdentifier) bool {
	if id == nil || req.GetId() == nil {
		return false
	}
	return vaulttypes.KeyFor(id) == vaulttypes.KeyFor(req.GetId())
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
		out.Responses = append(out.Responses, mergeSecretResponses(a.requests[i].GetId(), a.secretVotes[i][k]))
	}
	payload, err := anypb.New(out)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal aggregated GetSecretsResponse: %w", err)
	}
	return &commoncap.CapabilityResponse{Payload: payload}, nil
}

// mergeSecretResponses combines the shares of every response in the winning
// bucket, labelling the result with the identifier from the original request.
func mergeSecretResponses(id *vaultcommon.SecretIdentifier, rs []*vaultcommon.SecretResponse) *vaultcommon.SecretResponse {
	first := rs[0]
	if first.GetData() == nil {
		return &vaultcommon.SecretResponse{Id: id, Result: &vaultcommon.SecretResponse_Error{Error: first.GetError()}}
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
		}
	}

	return &vaultcommon.SecretResponse{
		Id: id,
		Result: &vaultcommon.SecretResponse_Data{Data: &vaultcommon.SecretData{
			EncryptedValue:               first.GetData().GetEncryptedValue(),
			EncryptedDecryptionKeyShares: entries,
		}},
	}
}

// validateSecretResponseShares checks that a response carries exactly one
// decryption share per requested encryption key, each within the size limit.
func (a *Aggregator) validateSecretResponseShares(secretReq *vaultcommon.SecretRequest, sr *vaultcommon.SecretResponse) error {
	data := sr.GetData()
	if data == nil {
		return nil
	}
	if secretReq == nil {
		return errors.New("share data returned for a nil request")
	}
	if err := validateShareKeyLabels(secretReq, data); err != nil {
		return err
	}
	ctx := contexts.WithCRE(context.Background(), contexts.CRE{Owner: secretReq.GetId().GetOwner()})
	for _, es := range data.GetEncryptedDecryptionKeyShares() {
		if err := validateEncryptedSharesEntry(es); err != nil {
			return err
		}
		shareSize, err := encryptedShareSizeForLimit(es)
		if err != nil {
			return err
		}
		if err := a.shareSizeLimit.Check(ctx, commonconfig.Size(shareSize)*commonconfig.Byte); err != nil {
			if _, ok := errors.AsType[limits.ErrorBoundLimited[commonconfig.Size]](err); ok {
				return fmt.Errorf("share provided exceeds maximum size allowed: %w", err)
			}
			return errors.New("failed to check share size")
		}
	}
	return nil
}

func validateShareKeyLabels(secretReq *vaultcommon.SecretRequest, data *vaultcommon.SecretData) error {
	if len(data.EncryptedDecryptionKeyShares) != len(secretReq.EncryptionKeys) {
		return fmt.Errorf("expected %d encrypted share entries, got %d", len(secretReq.EncryptionKeys), len(data.EncryptedDecryptionKeyShares))
	}

	allowed := make(map[string]struct{}, len(secretReq.EncryptionKeys))
	for _, pk := range secretReq.EncryptionKeys {
		allowed[pk] = struct{}{}
	}

	seen := make(map[string]bool, len(secretReq.EncryptionKeys))
	for _, es := range data.EncryptedDecryptionKeyShares {
		if es == nil {
			return errors.New("nil encrypted share entry")
		}
		if _, ok := allowed[es.EncryptionKey]; !ok {
			return fmt.Errorf("unexpected encryption key in response: %s", es.EncryptionKey)
		}
		if seen[es.EncryptionKey] {
			return fmt.Errorf("duplicate encryption key in response: %s", es.EncryptionKey)
		}
		seen[es.EncryptionKey] = true
	}

	for _, pk := range secretReq.EncryptionKeys {
		if !seen[pk] {
			return fmt.Errorf("missing encryption key in response: %s", pk)
		}
	}

	return nil
}

func validateEncryptedSharesEntry(es *vaultcommon.EncryptedShares) error {
	if len(es.BinaryShares) == 1 && len(es.Shares) == 0 {
		return nil
	}
	return errors.New("response must have exactly 1 binary share per encryption key")
}

func encryptedShareSizeForLimit(es *vaultcommon.EncryptedShares) (int, error) {
	if len(es.BinaryShares) == 1 {
		return len(es.BinaryShares[0]), nil
	}
	return 0, errors.New("no share to measure")
}
