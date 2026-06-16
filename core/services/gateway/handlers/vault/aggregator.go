package vault

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"maps"
	"slices"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
)

type baseAggregator struct {
	capabilitiesRegistry capabilitiesRegistry
	// vaultHandlerDonID scopes registry lookup when several vault DONs exist.
	//
	// Source: gateway job TOML [[gatewayConfig.ShardedDONs]] DonName (see deployment/cre/jobs/pkg/gateway_job.go),
	// loaded as ShardedDONConfig.DonName and passed as DONConfig.DonId (handler_factory.shardedDONsToLegacy;
	// DonId is a legacy field name for that string, not the on-chain uint32 id).
	//
	// Matching: capabilities.DON.Name when non-empty (v2), else decimal capabilities.DON.ID string (v1 sync).
	vaultHandlerDonID string
}

func (a *baseAggregator) Aggregate(ctx context.Context, l logger.Logger, resps map[string]jsonrpc.Response[json.RawMessage], currResp *jsonrpc.Response[json.RawMessage]) (*jsonrpc.Response[json.RawMessage], error) {
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
	if len(dons) == 0 {
		return nil, fmt.Errorf("no DON found for vault capability %s", vaultcommon.CapabilityID)
	}
	if len(dons) == 1 {
		don := dons[0]
		return &don, nil
	}

	handlerDonID := strings.TrimSpace(a.vaultHandlerDonID)
	if handlerDonID == "" {
		return nil, fmt.Errorf("multiple DONs (%d) host vault capability %s but vault handler DonId is empty; set ShardedDONConfig.DonName so DONConfig.DonId matches the vault DON name or id in the registry (%s)",
			len(dons), vaultcommon.CapabilityID, summarizeVaultRegistryDONs(dons))
	}

	var matches []capabilities.DONWithNodes
	for i := range dons {
		d := dons[i]
		if vaultDONMatchesHandlerDonID(&d.DON, handlerDonID) {
			matches = append(matches, d)
		}
	}
	switch len(matches) {
	case 0:
		return nil, fmt.Errorf("multiple DONs (%d) host vault capability %s but none match vault handler DonId %q; registry has %s",
			len(dons), vaultcommon.CapabilityID, a.vaultHandlerDonID, summarizeVaultRegistryDONs(dons))
	case 1:
		d := matches[0]
		return &d, nil
	default:
		return nil, fmt.Errorf("%d DONs match vault handler DonId %q for vault capability %s", len(matches), a.vaultHandlerDonID, vaultcommon.CapabilityID)
	}
}

// vaultDONMatchesHandlerDonID reports whether don is the vault DON this handler is configured for.
// handlerDonID is vaultHandlerDonID (jobspec DonName / DONConfig.DonId); see struct comment.
func vaultDONMatchesHandlerDonID(don *capabilities.DON, handlerDonID string) bool {
	if don.Name != "" {
		return don.Name == handlerDonID
	}
	return strconv.FormatUint(uint64(don.ID), 10) == handlerDonID
}

func summarizeVaultRegistryDONs(dons []capabilities.DONWithNodes) string {
	var b strings.Builder
	for i, d := range dons {
		if i > 0 {
			b.WriteString("; ")
		}
		_, _ = fmt.Fprintf(&b, "name=%q id=%d", d.DON.Name, d.DON.ID)
	}
	return b.String()
}

func (a *baseAggregator) validateUsingQuorum(don capabilities.DON, resps map[string]jsonrpc.Response[json.RawMessage], l logger.Logger) (*jsonrpc.Response[json.RawMessage], error) {
	requiredQuorum := int(2*don.F + 1)

	if len(resps) < requiredQuorum {
		return nil, errInsufficientResponsesForQuorum
	}

	shaToCount := map[string]int{}
	maxShaToCount := 0
	for _, r := range resps {
		sha, err := a.sha(&r)
		if err != nil {
			l.Errorw("failed to compute digest of response during quorum validation, skipping...", "error", err)
			continue
		}
		shaToCount[sha]++
		if shaToCount[sha] > maxShaToCount {
			maxShaToCount = shaToCount[sha]
		}
	}

	var qualifiedDigests []string
	for sha, n := range shaToCount {
		if n >= requiredQuorum {
			qualifiedDigests = append(qualifiedDigests, sha)
		}
	}
	if len(qualifiedDigests) > 0 {
		slices.Sort(qualifiedDigests)
		want := qualifiedDigests[0]
		for _, k := range slices.Sorted(maps.Keys(resps)) {
			r := resps[k]
			sha, err := a.sha(&r)
			if err != nil {
				continue
			}
			if sha == want {
				out := r
				return &out, nil
			}
		}
	}

	remainingResponses := len(don.Members) - len(resps)
	if maxShaToCount+remainingResponses < requiredQuorum {
		l.Warnw("quorum unattainable for request", "requiredQuorum", requiredQuorum, "remainingResponses", remainingResponses, "maxShaToCount", maxShaToCount, "remainingResponses", remainingResponses, "allResponses", resps)
		return nil, errors.New(errQuorumUnobtainable.Error() + ". RequiredQuorum=" + strconv.Itoa(requiredQuorum) + ". maxShaToCount=" + strconv.Itoa(maxShaToCount) + " remainingResponses=" + strconv.Itoa(remainingResponses))
	}

	return nil, errInsufficientResponsesForQuorum
}

func (a *baseAggregator) unmarshal(r io.Reader, to any) error {
	d := json.NewDecoder(r)
	d.DisallowUnknownFields()
	return d.Decode(to)
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
		if resp.Error != nil {
			return nil, errors.New("response has an error, cannot validate signatures. Error: " + resp.Error.Error())
		}
		return nil, errors.New("response result and error both are is nil: cannot validate signatures")
	}

	r := &vaulttypes.SignedOCRResponse{}
	err := a.unmarshal(bytes.NewReader(*resp.Result), r)
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
