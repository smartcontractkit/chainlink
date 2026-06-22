package vaultshare

import (
	"context"
	"errors"
	"fmt"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/executable/request"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

var (
	ErrFastPathReconstruction = errors.New("fast-path secret reconstruction failed after collecting peer responses")
)

var _ request.ResponseAggregator = (*Aggregator)(nil)

// NewAggregatorFactory returns a per-request aggregator factory gated on VaultFastPathGetSecretsEnabled.
func NewAggregatorFactory(f int) request.AggregatorFactory {
	gate, err := limits.MakeGateLimiter(limits.Factory{Settings: cresettings.DefaultGetter}, cresettings.Default.VaultFastPathGetSecretsEnabled)
	if err != nil {
		return func(context.Context, commoncap.CapabilityRequest) request.ResponseAggregator { return nil }
	}
	threshold := f + 1
	maxResponses := 2*f + 1
	return func(ctx context.Context, _ commoncap.CapabilityRequest) request.ResponseAggregator {
		if gate.AllowErr(ctx) != nil {
			return nil
		}
		return NewAggregator(threshold, maxResponses)
	}
}

type Aggregator struct {
	threshold      int
	maxResponses   int
	responseCount  int
	errorCount     map[string]int
	peerResponses  map[p2ptypes.PeerID]*vaultcommon.GetSecretsResponse
}

func NewAggregator(threshold, maxResponses int) *Aggregator {
	return &Aggregator{
		threshold:     threshold,
		maxResponses:  maxResponses,
		errorCount:    make(map[string]int),
		peerResponses: make(map[p2ptypes.PeerID]*vaultcommon.GetSecretsResponse),
	}
}

func (a *Aggregator) OnResponse(peerID p2ptypes.PeerID, msg *types.MessageBody) (*commoncap.CapabilityResponse, error) {
	if msg.Error != types.Error_OK {
		a.errorCount[msg.ErrorMsg]++
		a.responseCount++
		if a.errorCount[msg.ErrorMsg] >= a.threshold {
			return nil, errors.New(msg.ErrorMsg)
		}
		if a.responseCount >= a.maxResponses {
			return nil, ErrFastPathReconstruction
		}
		return nil, nil
	}

	resp, err := pb.UnmarshalCapabilityResponse(msg.Payload)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal capability response: %w", err)
	}
	if resp.Payload == nil {
		a.responseCount++
		if a.responseCount >= a.maxResponses {
			return nil, ErrFastPathReconstruction
		}
		return nil, nil
	}

	vaultResp := &vaultcommon.GetSecretsResponse{}
	if err := resp.Payload.UnmarshalTo(vaultResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal vault GetSecretsResponse: %w", err)
	}

	a.peerResponses[peerID] = vaultResp
	a.responseCount++

	if len(a.peerResponses) < a.threshold {
		if a.responseCount >= a.maxResponses {
			return nil, ErrFastPathReconstruction
		}
		return nil, nil
	}

	merged, err := mergeGetSecretsResponses(a.peerResponses)
	if err != nil {
		if a.responseCount >= a.maxResponses {
			return nil, ErrFastPathReconstruction
		}
		return nil, nil
	}

	anyPayload, err := anypb.New(merged)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal merged vault response: %w", err)
	}
	return &commoncap.CapabilityResponse{Payload: anyPayload}, nil
}

func (a *Aggregator) OnTimeout() (*commoncap.CapabilityResponse, error) {
	return nil, ErrFastPathReconstruction
}

func mergeGetSecretsResponses(peerResponses map[p2ptypes.PeerID]*vaultcommon.GetSecretsResponse) (*vaultcommon.GetSecretsResponse, error) {
	mergedByKey := make(map[string]*vaultcommon.SecretResponse)
	for _, resp := range peerResponses {
		for _, sr := range resp.Responses {
			if sr == nil || sr.Id == nil {
				continue
			}
			if sr.GetError() != "" {
				continue
			}
			key := secretKey(sr.Id)
			existing, ok := mergedByKey[key]
			if !ok {
				mergedByKey[key] = cloneSecretResponse(sr)
				continue
			}
			if err := mergeSecretResponseData(existing, sr); err != nil {
				return nil, err
			}
		}
	}
	if len(mergedByKey) == 0 {
		return nil, errors.New("no secret data in peer responses")
	}
	out := &vaultcommon.GetSecretsResponse{Responses: make([]*vaultcommon.SecretResponse, 0, len(mergedByKey))}
	for _, sr := range mergedByKey {
		out.Responses = append(out.Responses, sr)
	}
	return out, nil
}

func secretKey(id *vaultcommon.SecretIdentifier) string {
	return id.Owner + "::" + id.Namespace + "::" + id.Key
}

func cloneSecretResponse(sr *vaultcommon.SecretResponse) *vaultcommon.SecretResponse {
	cloned := &vaultcommon.SecretResponse{Id: sr.Id}
	if data := sr.GetData(); data != nil {
		cloned.Result = &vaultcommon.SecretResponse_Data{
			Data: &vaultcommon.SecretData{
				EncryptedValue: data.EncryptedValue,
			},
		}
		for _, share := range data.EncryptedDecryptionKeyShares {
			cloned.GetData().EncryptedDecryptionKeyShares = append(
				cloned.GetData().EncryptedDecryptionKeyShares,
				cloneEncryptedShares(share),
			)
		}
	}
	return cloned
}

func cloneEncryptedShares(share *vaultcommon.EncryptedShares) *vaultcommon.EncryptedShares {
	if share == nil {
		return nil
	}
	out := &vaultcommon.EncryptedShares{EncryptionKey: share.EncryptionKey}
	out.BinaryShares = append(out.BinaryShares, share.BinaryShares...)
	out.Shares = append(out.Shares, share.Shares...)
	return out
}

func mergeSecretResponseData(dst, src *vaultcommon.SecretResponse) error {
	dstData := dst.GetData()
	srcData := src.GetData()
	if dstData == nil || srcData == nil {
		return errors.New("missing secret data in peer response")
	}
	if dstData.EncryptedValue != "" && srcData.EncryptedValue != "" && dstData.EncryptedValue != srcData.EncryptedValue {
		return errors.New("inconsistent encrypted secret value across peer responses")
	}
	if dstData.EncryptedValue == "" {
		dstData.EncryptedValue = srcData.EncryptedValue
	}
	for _, srcShare := range srcData.EncryptedDecryptionKeyShares {
		if srcShare == nil {
			continue
		}
		var dstShare *vaultcommon.EncryptedShares
		for _, existing := range dstData.EncryptedDecryptionKeyShares {
			if existing != nil && existing.EncryptionKey == srcShare.EncryptionKey {
				dstShare = existing
				break
			}
		}
		if dstShare == nil {
			dstData.EncryptedDecryptionKeyShares = append(dstData.EncryptedDecryptionKeyShares, cloneEncryptedShares(srcShare))
			continue
		}
		dstShare.BinaryShares = append(dstShare.BinaryShares, srcShare.BinaryShares...)
		dstShare.Shares = append(dstShare.Shares, srcShare.Shares...)
	}
	return nil
}
