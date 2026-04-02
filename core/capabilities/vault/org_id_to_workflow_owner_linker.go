package vault

import (
	"context"
	"errors"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/services/orgresolver"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
)

// LinkedVaultRequestIdentity is the resolved identity forwarded to the vault OCR plugin.
type LinkedVaultRequestIdentity struct {
	OrgID         string
	WorkflowOwner string
}

// OrgIDToWorkflowOwnerLinker centralizes vault request identity resolution and verification.
type OrgIDToWorkflowOwnerLinker struct {
	orgResolver                    orgresolver.OrgResolver
	vaultOrgIDAsSecretOwnerEnabled limits.GateLimiter
}

// NewOrgIDToWorkflowOwnerLinker constructs the gated org/workflow-owner linker used by vault.
func NewOrgIDToWorkflowOwnerLinker(orgResolver orgresolver.OrgResolver, limitsFactory limits.Factory) (*OrgIDToWorkflowOwnerLinker, error) {
	vaultOrgIDAsSecretOwnerEnabled, err := limits.MakeGateLimiter(limitsFactory, cresettings.Default.VaultOrgIdAsSecretOwnerEnabled)
	if err != nil {
		return nil, fmt.Errorf("could not create vault org-id-as-owner gate limiter: %w", err)
	}

	return &OrgIDToWorkflowOwnerLinker{
		orgResolver:                    orgResolver,
		vaultOrgIDAsSecretOwnerEnabled: vaultOrgIDAsSecretOwnerEnabled,
	}, nil
}

// Close releases the gate limiter resources owned by the linker.
func (l *OrgIDToWorkflowOwnerLinker) Close() error {
	return l.vaultOrgIDAsSecretOwnerEnabled.Close()
}

// Link resolves or verifies the request identity from the caller-provided org and workflow owner.
func (l *OrgIDToWorkflowOwnerLinker) Link(ctx context.Context, orgID string, workflowOwner string) (LinkedVaultRequestIdentity, error) {
	enabled, err := l.vaultOrgIDAsSecretOwnerEnabled.Limit(ctx)
	if err != nil {
		return LinkedVaultRequestIdentity{}, fmt.Errorf("failed to evaluate vault org-id-as-owner gate: %w", err)
	}
	if !enabled {
		return LinkedVaultRequestIdentity{OrgID: orgID, WorkflowOwner: workflowOwner}, nil
	}

	if orgID == "" && workflowOwner == "" {
		return LinkedVaultRequestIdentity{}, errors.New("org_id and workflow owner cannot both be empty")
	}

	if workflowOwner == "" {
		return LinkedVaultRequestIdentity{OrgID: orgID, WorkflowOwner: workflowOwner}, nil
	}
	if l.orgResolver == nil {
		return LinkedVaultRequestIdentity{}, errors.New("org resolver is nil")
	}

	resolvedOrgID, err := l.orgResolver.Get(ctx, workflowOwner)
	if err != nil {
		if orgID != "" {
			return LinkedVaultRequestIdentity{}, fmt.Errorf("failed to verify org_id %q for workflow owner %q: %w", orgID, workflowOwner, err)
		}
		return LinkedVaultRequestIdentity{}, fmt.Errorf("failed to resolve org_id for workflow owner %q: %w", workflowOwner, err)
	}
	if resolvedOrgID == "" {
		return LinkedVaultRequestIdentity{}, fmt.Errorf("resolved empty org_id for workflow owner %q", workflowOwner)
	}
	if orgID != "" && resolvedOrgID != orgID {
		return LinkedVaultRequestIdentity{}, fmt.Errorf("workflow owner %q resolves to org_id %q, does not match request org_id %q", workflowOwner, resolvedOrgID, orgID)
	}

	return LinkedVaultRequestIdentity{
		OrgID:         resolvedOrgID,
		WorkflowOwner: workflowOwner,
	}, nil
}
