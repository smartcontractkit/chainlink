package vault

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
)

type baseAggregator struct {
	capabilitiesRegistry capabilitiesRegistry
}

func (a *baseAggregator) Aggregate(ctx context.Context, l logger.Logger, resps map[string]*jsonrpc.Response[json.RawMessage], currResp *jsonrpc.Response[json.RawMessage]) (*jsonrpc.Response[json.RawMessage], error) {
	don, err := a.donForVaultCapability(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get DON for vault capability: %w", err)
	}

	currResp, err = a.validateUsingSignatures(don.DON, don.Nodes, currResp)
	if err == nil {
		return currResp, nil
	}

	l.Debugw("failed to validate signatures, falling back to quorum aggregation", "error", err)
	currResp, err = a.validateUsingQuorum(don.DON, resps, l)
	if err != nil {
		return nil, fmt.Errorf("failed to validate using quorum: %w", err)
	}

	return currResp, nil
}

func (a *baseAggregator) donForVaultCapability(ctx context.Context) (*capabilities.DONWithNodes, error) {
	dons, err := a.capabilitiesRegistry.DONsForCapability(ctx, vaultcommon.CapabilityID)
	if err != nil {
		return nil, err
	}
	// TODO: Support multiple vault capabilities in the capability registry.
	// For the initial Smartcon deployment there will be exactly one Vault capability
	// split across both DON families.
	if len(dons) != 1 {
		return nil, fmt.Errorf("expected exactly one DON for vault capability, found %d", len(dons))
	}

	don := dons[0]
	return &don, nil
}

func (a *baseAggregator) validateUsingQuorum(don capabilities.DON, resps map[string]*jsonrpc.Response[json.RawMessage], l logger.Logger) (*jsonrpc.Response[json.RawMessage], error) {
	requiredQuorum := int(2*don.F + 1)

	if len(resps) < requiredQuorum {
		return nil, errInsufficientResponsesForQuorum
	}

	shaToCount := map[string]int{}
	maxShaToCount := 0
	for _, r := range resps {
		sha, err := a.sha(r)
		if err != nil {
			l.Errorw("failed to compute digest of response during quorum validation, skipping...", "error", err)
			continue
		}
		shaToCount[sha]++
		if shaToCount[sha] > maxShaToCount {
			maxShaToCount = shaToCount[sha]
		}
		if shaToCount[sha] >= requiredQuorum {
			return r, nil
		}
	}

	remainingResponses := len(don.Members) - len(resps)
	if maxShaToCount+remainingResponses < requiredQuorum {
		return nil, errQuorumUnobtainable
	}

	return nil, errInsufficientResponsesForQuorum
}

// sha computes a hash of the response, taking into account that when a response
// contains signatures, they should be computed from the hash computation as they are not guaranteed
// to be identical.
func (a *baseAggregator) sha(resp *jsonrpc.Response[json.RawMessage]) (string, error) {
	// Case: No result so therefore no signatures. Early exit.
	if resp.Result == nil {
		return resp.Digest()
	}

	r := &vaulttypes.SignedOCRResponse{}
	err := json.Unmarshal(*resp.Result, r)
	if err != nil {
		return "", err
	}

	// Case: Result has no signatures. Early exit.
	if len(r.Signatures) == 0 {
		return resp.Digest()
	}

	// Case: We have signatures. In this case we copy the response,
	// zeroing out the signatures, and take the resulting digest.
	b, err := json.Marshal(resp)
	if err != nil {
		return "", err
	}

	copied := &jsonrpc.Response[json.RawMessage]{}
	err = json.Unmarshal(b, copied)
	if err != nil {
		return "", err
	}

	r.Signatures = nil
	rawMessage, err := json.Marshal(r)
	if err != nil {
		return "", err
	}
	copied.Result = (*json.RawMessage)(&rawMessage)
	return copied.Digest()
}

func (a *baseAggregator) validateUsingSignatures(don capabilities.DON, nodes []capabilities.Node, resp *jsonrpc.Response[json.RawMessage]) (*jsonrpc.Response[json.RawMessage], error) {
	if resp.Result == nil {
		return nil, errors.New("response result is nil: cannot validate signatures")
	}

	if resp.Method == vaulttypes.MethodSecretsGet {
		// SecretsGet responses are not signed.
		return resp, errors.New("cannot validate signatures for Get requests")
	}

	r := &vaulttypes.SignedOCRResponse{}
	err := json.Unmarshal(*resp.Result, r)
	if err != nil {
		return nil, err
	}

	signers := []common.Address{}
	for _, n := range nodes {
		signers = append(signers, common.BytesToAddress(n.Signer[0:20]))
	}

	err = vaulttypes.ValidateSignatures(r, signers, int(don.F+1))
	if err != nil {
		return nil, fmt.Errorf("failed to validate signatures: %w", err)
	}

	return resp, nil
}
