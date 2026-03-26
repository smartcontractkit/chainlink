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
	orgResolver         orgresolver.OrgResolver
	vaultJWTAuthEnabled limits.GateLimiter
	lggr                logger.Logger
}

// NewOrgIDToWorkflowOwnerLinker constructs the gated org/workflow-owner linker used by vault.
func NewOrgIDToWorkflowOwnerLinker(lggr logger.Logger, orgResolver orgresolver.OrgResolver, limitsFactory limits.Factory) (*OrgIDToWorkflowOwnerLinker, error) {
	vaultJWTAuthEnabled, err := limits.MakeGateLimiter(limitsFactory, cresettings.Default.VaultJWTAuthEnabled)
	if err != nil {
		return nil, fmt.Errorf("could not create vault JWT auth gate limiter: %w", err)
	}

	return &OrgIDToWorkflowOwnerLinker{
		orgResolver:         orgResolver,
		vaultJWTAuthEnabled: vaultJWTAuthEnabled,
		lggr:                logger.Named(lggr, "OrgIDToWorkflowOwnerLinker"),
	}, nil
}

// Close releases the gate limiter resources owned by the linker.
func (l *OrgIDToWorkflowOwnerLinker) Close() error {
	if l == nil || l.vaultJWTAuthEnabled == nil {
		return nil
	}

	return l.vaultJWTAuthEnabled.Close()
}

// Link resolves or verifies the request identity from the caller-provided org and workflow owner.
func (l *OrgIDToWorkflowOwnerLinker) Link(ctx context.Context, orgID string, workflowOwner string) (LinkedVaultRequestIdentity, error) {
	if l == nil || l.vaultJWTAuthEnabled == nil {
		return LinkedVaultRequestIdentity{OrgID: orgID, WorkflowOwner: workflowOwner}, nil
	}

	enabled, err := l.vaultJWTAuthEnabled.Limit(ctx)
	if err != nil {
		l.lggr.Errorw("failed to evaluate vault JWT auth gate", "orgID", orgID, "workflowOwner", workflowOwner, "err", err)
		return LinkedVaultRequestIdentity{}, fmt.Errorf("failed to evaluate vault JWT auth gate: %w", err)
	}
	if !enabled {
		l.lggr.Debugw("skipping vault identity linking because JWT auth gate is disabled", "orgID", orgID, "workflowOwner", workflowOwner)
		return LinkedVaultRequestIdentity{OrgID: orgID, WorkflowOwner: workflowOwner}, nil
	}

	orgID = strings.TrimSpace(orgID)
	workflowOwner = strings.TrimSpace(workflowOwner)
	if orgID != "" {
		if workflowOwner == "" {
			l.lggr.Debugw("using trusted org_id without workflow owner verification", "orgID", orgID)
			return LinkedVaultRequestIdentity{OrgID: orgID, WorkflowOwner: workflowOwner}, nil
		}
		if l.orgResolver == nil {
			l.lggr.Errorw("cannot verify org_id without org resolver", "orgID", orgID, "workflowOwner", workflowOwner)
			return LinkedVaultRequestIdentity{}, errors.New("org resolver is required when workflow owner is provided")
		}

		l.lggr.Debugw("verifying workflow owner against trusted org_id", "orgID", orgID, "workflowOwner", workflowOwner)
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

		l.lggr.Debugw("verified workflow owner against trusted org_id", "orgID", orgID, "workflowOwner", workflowOwner)
		return LinkedVaultRequestIdentity{OrgID: orgID, WorkflowOwner: workflowOwner}, nil
	}

	if workflowOwner == "" {
		l.lggr.Errorw("missing workflow owner for org_id resolution")
		return LinkedVaultRequestIdentity{}, errors.New("workflow owner is required when org_id is missing")
	}
	if l.orgResolver == nil {
		l.lggr.Errorw("cannot resolve org_id without org resolver", "workflowOwner", workflowOwner)
		return LinkedVaultRequestIdentity{}, errors.New("org resolver is required when org_id is missing")
	}

	l.lggr.Debugw("resolving org_id from workflow owner", "workflowOwner", workflowOwner)
	resolvedOrgID, resolveErr := l.orgResolver.Get(ctx, workflowOwner)
	if resolveErr != nil {
		l.lggr.Errorw("failed to resolve org_id from workflow owner", "workflowOwner", workflowOwner, "err", resolveErr)
		return LinkedVaultRequestIdentity{}, fmt.Errorf("failed to resolve org_id for workflow owner %q: %w", workflowOwner, resolveErr)
	}
	if strings.TrimSpace(resolvedOrgID) == "" {
		l.lggr.Errorw("resolved empty org_id from workflow owner", "workflowOwner", workflowOwner)
		return LinkedVaultRequestIdentity{}, fmt.Errorf("resolved empty org_id for workflow owner %q", workflowOwner)
	}

	l.lggr.Debugw("resolved org_id from workflow owner", "workflowOwner", workflowOwner, "orgID", strings.TrimSpace(resolvedOrgID))
	return LinkedVaultRequestIdentity{
		OrgID:         strings.TrimSpace(resolvedOrgID),
		WorkflowOwner: workflowOwner,
	}, nil
}
