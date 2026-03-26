package vault

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/smartcontractkit/chainlink-common/pkg/services/orgresolver"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
)

type LinkedVaultRequestIdentity struct {
	OrgID         string
	WorkflowOwner string
}

type OrgIdToWorkflowOwnerLinker struct {
	orgResolver         orgresolver.OrgResolver
	vaultJWTAuthEnabled limits.GateLimiter
}

func NewOrgIdToWorkflowOwnerLinker(orgResolver orgresolver.OrgResolver, limitsFactory limits.Factory) (*OrgIdToWorkflowOwnerLinker, error) {
	vaultJWTAuthEnabled, err := limits.MakeGateLimiter(limitsFactory, cresettings.Default.VaultJWTAuthEnabled)
	if err != nil {
		return nil, fmt.Errorf("could not create vault JWT auth gate limiter: %w", err)
	}

	return &OrgIdToWorkflowOwnerLinker{
		orgResolver:         orgResolver,
		vaultJWTAuthEnabled: vaultJWTAuthEnabled,
	}, nil
}

func (l *OrgIdToWorkflowOwnerLinker) Close() error {
	if l == nil || l.vaultJWTAuthEnabled == nil {
		return nil
	}

	return l.vaultJWTAuthEnabled.Close()
}

func (l *OrgIdToWorkflowOwnerLinker) Link(ctx context.Context, requestID string, orgID string, workflowOwner string) (LinkedVaultRequestIdentity, error) {
	if l == nil || l.vaultJWTAuthEnabled == nil {
		return LinkedVaultRequestIdentity{OrgID: orgID, WorkflowOwner: workflowOwner}, nil
	}

	enabled, err := l.vaultJWTAuthEnabled.Limit(ctx)
	if err != nil {
		return LinkedVaultRequestIdentity{}, fmt.Errorf("failed to evaluate vault JWT auth gate: %w", err)
	}
	if !enabled {
		return LinkedVaultRequestIdentity{OrgID: orgID, WorkflowOwner: workflowOwner}, nil
	}

	orgID = strings.TrimSpace(orgID)
	workflowOwner = strings.TrimSpace(workflowOwner)
	if orgID != "" {
		if workflowOwner == "" {
			return LinkedVaultRequestIdentity{OrgID: orgID, WorkflowOwner: workflowOwner}, nil
		}
		if l.orgResolver == nil {
			return LinkedVaultRequestIdentity{}, errors.New("org resolver is required when workflow owner is provided")
		}

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

	if workflowOwner == "" {
		workflowOwner = workflowOwnerFromRequestID(requestID)
	}
	if workflowOwner == "" {
		return LinkedVaultRequestIdentity{}, errors.New("workflow owner is required when org_id is missing")
	}
	if l.orgResolver == nil {
		return LinkedVaultRequestIdentity{}, errors.New("org resolver is required when org_id is missing")
	}

	resolvedOrgID, err := l.orgResolver.Get(ctx, workflowOwner)
	if err != nil {
		return LinkedVaultRequestIdentity{}, fmt.Errorf("failed to resolve org_id for workflow owner %q: %w", workflowOwner, err)
	}
	if strings.TrimSpace(resolvedOrgID) == "" {
		return LinkedVaultRequestIdentity{}, fmt.Errorf("resolved empty org_id for workflow owner %q", workflowOwner)
	}

	return LinkedVaultRequestIdentity{
		OrgID:         strings.TrimSpace(resolvedOrgID),
		WorkflowOwner: workflowOwner,
	}, nil
}

func workflowOwnerFromRequestID(requestID string) string {
	owner, _, ok := strings.Cut(requestID, vaulttypes.RequestIDSeparator)
	if !ok {
		return ""
	}

	return owner
}
