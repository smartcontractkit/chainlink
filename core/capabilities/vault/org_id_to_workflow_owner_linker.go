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
		l.lggr.Errorw("failed to evaluate vault org-id-as-owner gate", "orgID", orgID, "workflowOwner", workflowOwner, "err", err)
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

		l.lggr.Errorw("missing workflow owner for org_id resolution")
		return LinkedVaultRequestIdentity{}, errors.New("workflow owner is required when org_id is missing")
	}
	if l.orgResolver == nil {
		if orgID != "" {
			l.lggr.Errorw("cannot verify org_id without org resolver", "orgID", orgID, "workflowOwner", workflowOwner)
			return LinkedVaultRequestIdentity{}, errors.New("org resolver is required when workflow owner is provided")
		}
		l.lggr.Errorw("cannot resolve org_id without org resolver", "workflowOwner", workflowOwner)
		return LinkedVaultRequestIdentity{}, errors.New("org resolver is required when org_id is missing")
	}

	if orgID != "" {
		resolvedOrgID, err := l.orgResolver.Get(ctx, workflowOwner)
		if err != nil {
			l.lggr.Errorw("failed to verify workflow owner against org_id", "orgID", orgID, "workflowOwner", workflowOwner, "err", err)
			return LinkedVaultRequestIdentity{}, fmt.Errorf("failed to verify org_id %q for workflow owner %q: %w", orgID, workflowOwner, err)
		}
		resolvedOrgID = strings.TrimSpace(resolvedOrgID)
		if resolvedOrgID == "" {
			l.lggr.Errorw("workflow owner resolved to empty org_id", "workflowOwner", workflowOwner)
			return LinkedVaultRequestIdentity{}, fmt.Errorf("resolved empty org_id for workflow owner %q", workflowOwner)
		}
		if resolvedOrgID != orgID {
			l.lggr.Errorw("workflow owner org_id verification failed", "requestOrgID", orgID, "resolvedOrgID", resolvedOrgID, "workflowOwner", workflowOwner)
			return LinkedVaultRequestIdentity{}, fmt.Errorf("workflow owner %q resolves to org_id %q, does not match request org_id %q", workflowOwner, resolvedOrgID, orgID)
		}

		return LinkedVaultRequestIdentity{OrgID: orgID, WorkflowOwner: workflowOwner}, nil
	}

	resolvedOrgID, resolveErr := l.orgResolver.Get(ctx, workflowOwner)
	if resolveErr != nil {
		l.lggr.Errorw("failed to resolve org_id from workflow owner", "workflowOwner", workflowOwner, "err", resolveErr)
		return LinkedVaultRequestIdentity{}, fmt.Errorf("failed to resolve org_id for workflow owner %q: %w", workflowOwner, resolveErr)
	}
	if strings.TrimSpace(resolvedOrgID) == "" {
		l.lggr.Errorw("resolved empty org_id from workflow owner", "workflowOwner", workflowOwner)
		return LinkedVaultRequestIdentity{}, fmt.Errorf("resolved empty org_id for workflow owner %q", workflowOwner)
	}

	return LinkedVaultRequestIdentity{
		OrgID:         strings.TrimSpace(resolvedOrgID),
		WorkflowOwner: workflowOwner,
	}, nil
}
