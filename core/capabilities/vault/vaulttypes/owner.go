package vaulttypes

import (
	"context"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/smartcontractkit/chainlink-common/pkg/services/orgresolver"
)

// IsWorkflowOwnerAddress returns true if the owner string looks like a normalized
// Ethereum address (40 hex chars, possibly with 0x prefix). This is used to
// distinguish between a workflow owner (on-chain address) and an org ID.
func IsWorkflowOwnerAddress(owner string) bool {
	normalized := strings.TrimPrefix(strings.ToLower(owner), "0x")
	if len(normalized) != 40 {
		return false
	}
	_, err := hex.DecodeString(normalized)
	return err == nil
}

// ResolveOwnerToOrgID resolves an owner to an org ID. If the owner is a workflow
// owner (on-chain address), it uses the OrgResolver to find the corresponding org ID.
// If the owner is already an org ID (not an address), it is returned as-is.
// If orgResolver is nil and the owner is a workflow address, an error is returned.
func ResolveOwnerToOrgID(ctx context.Context, orgResolver orgresolver.OrgResolver, owner string) (string, error) {
	if !IsWorkflowOwnerAddress(owner) {
		return owner, nil
	}

	if orgResolver == nil {
		return "", fmt.Errorf("cannot resolve workflow owner %q to org ID: org resolver is not configured", owner)
	}

	orgID, err := orgResolver.Get(ctx, owner)
	if err != nil {
		return "", fmt.Errorf("failed to resolve org ID for workflow owner %q: %w", owner, err)
	}

	if orgID == "" {
		return "", fmt.Errorf("org resolver returned empty org ID for workflow owner %q", owner)
	}

	return orgID, nil
}
