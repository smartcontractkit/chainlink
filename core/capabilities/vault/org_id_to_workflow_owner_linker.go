package vault

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
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
	lggr                           logger.Logger
}

// NewOrgIDToWorkflowOwnerLinker constructs the gated org/workflow-owner linker used by vault.
func NewOrgIDToWorkflowOwnerLinker(lggr logger.Logger, orgResolver orgresolver.OrgResolver, limitsFactory limits.Factory) (*OrgIDToWorkflowOwnerLinker, error) {
	vaultOrgIDAsSecretOwnerEnabled, err := limits.MakeGateLimiter(limitsFactory, cresettings.Default.VaultOrgIdAsSecretOwnerEnabled)
	if err != nil {
		return nil, fmt.Errorf("could not create vault org-id-as-owner gate limiter: %w", err)
	}

	return &OrgIDToWorkflowOwnerLinker{
		orgResolver:                    orgResolver,
		vaultOrgIDAsSecretOwnerEnabled: vaultOrgIDAsSecretOwnerEnabled,
		lggr:                           logger.Named(lggr, "OrgIDToWorkflowOwnerLinker"),
	}, nil
}

// Close releases the gate limiter resources owned by the linker.
func (l *OrgIDToWorkflowOwnerLinker) Close() error {
	if l == nil || l.vaultOrgIDAsSecretOwnerEnabled == nil {
		return nil
	}

	return l.vaultOrgIDAsSecretOwnerEnabled.Close()
}

// Link resolves or verifies the request identity from the caller-provided org and workflow owner.
func (l *OrgIDToWorkflowOwnerLinker) Link(ctx context.Context, orgID string, workflowOwner string) (LinkedVaultRequestIdentity, error) {
	if l == nil || l.vaultOrgIDAsSecretOwnerEnabled == nil {
		return LinkedVaultRequestIdentity{OrgID: orgID, WorkflowOwner: workflowOwner}, nil
	}

	enabled, err := l.vaultOrgIDAsSecretOwnerEnabled.Limit(ctx)
	if err != nil {
		return LinkedVaultRequestIdentity{}, fmt.Errorf("failed to evaluate vault org-id-as-owner gate: %w", err)
	}
	if !enabled {
		return LinkedVaultRequestIdentity{OrgID: orgID, WorkflowOwner: workflowOwner}, nil
	}

	orgID = strings.TrimSpace(orgID)
	workflowOwner = strings.TrimSpace(workflowOwner)
	if workflowOwner == "" {
		if orgID != "" {
			return LinkedVaultRequestIdentity{OrgID: orgID, WorkflowOwner: workflowOwner}, nil
		}

		return LinkedVaultRequestIdentity{}, errors.New("workflow owner is required when org_id is missing")
	}
	if l.orgResolver == nil {
		if orgID != "" {
			return LinkedVaultRequestIdentity{}, errors.New("org resolver is required when workflow owner is provided")
		}
		return LinkedVaultRequestIdentity{}, errors.New("org resolver is required when org_id is missing")
	}

	if orgID != "" {
		resolvedOrgID, err := l.orgResolver.Get(ctx, workflowOwner)
		if err != nil {
			return LinkedVaultRequestIdentity{}, fmt.Errorf("failed to verify org_id %q for workflow owner %q: %w", orgID, workflowOwner, err)
		}
		resolvedOrgID = strings.TrimSpace(resolvedOrgID)
		if resolvedOrgID == "" {
			return LinkedVaultRequestIdentity{}, fmt.Errorf("resolved empty org_id for workflow owner %q", workflowOwner)
		}
		if resolvedOrgID != orgID {
			return LinkedVaultRequestIdentity{}, fmt.Errorf("workflow owner %q resolves to org_id %q, does not match request org_id %q", workflowOwner, resolvedOrgID, orgID)
		}

		return LinkedVaultRequestIdentity{OrgID: orgID, WorkflowOwner: workflowOwner}, nil
	}

	resolvedOrgID, resolveErr := l.orgResolver.Get(ctx, workflowOwner)
	if resolveErr != nil {
		return LinkedVaultRequestIdentity{}, fmt.Errorf("failed to resolve org_id for workflow owner %q: %w", workflowOwner, resolveErr)
	}
	if strings.TrimSpace(resolvedOrgID) == "" {
		return LinkedVaultRequestIdentity{}, fmt.Errorf("resolved empty org_id for workflow owner %q", workflowOwner)
	}

	return LinkedVaultRequestIdentity{
		OrgID:         strings.TrimSpace(resolvedOrgID),
		WorkflowOwner: workflowOwner,
	}, nil
}
