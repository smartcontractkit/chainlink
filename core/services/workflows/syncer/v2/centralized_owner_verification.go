package v2

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/settings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	pkgworkflows "github.com/smartcontractkit/chainlink-common/pkg/workflows"
)

// ErrCentralizedOwnerOrgMismatch is returned when a workflow owner claimed by a
// centralized (off-chain) source does not match the owner deterministically derived
// from its organization ID.
var ErrCentralizedOwnerOrgMismatch = errors.New("centralized workflow owner does not match owner derived from its organization ID")

// ErrTenantIDNotConfigured is returned when cresettings.TenantID is unset (zero).
var ErrTenantIDNotConfigured = errors.New("TenantID cresetting is not configured")

// tenantIDFromSettings reads TenantID from cresettings. Zero is treated as unset.
func tenantIDFromSettings(ctx context.Context, getter settings.Getter) (uint64, error) {
	tenantID, err := cresettings.Default.TenantID.GetOrDefault(ctx, getter)
	if err != nil {
		return 0, err
	}
	if tenantID == 0 {
		return 0, ErrTenantIDNotConfigured
	}
	return tenantID, nil
}

// isCentralizedWorkflowSource reports whether workflow metadata originated from a
// centralized / off-chain source (gRPC workflow registry or local file source), as
// opposed to the on-chain workflow registry contract (source prefix "contract:").
func isCentralizedWorkflowSource(source string) bool {
	return strings.HasPrefix(source, "grpc:") || strings.HasPrefix(source, "file:")
}

// verifyCentralizedOwnerOrgMapping checks that a workflow owner claimed by a
// centralized source is consistent with the organization ID. Legitimate centralized
// CRE workflow owners are deterministically derived from (tenantID, orgID) via
// workflows.GenerateWorkflowOwnerAddress.
func verifyCentralizedOwnerOrgMapping(lggr logger.Logger, source, ownerHex, orgID string, tenantID uint64) error {
	derived, err := pkgworkflows.GenerateWorkflowOwnerAddress(strconv.FormatUint(tenantID, 10), orgID)
	if err != nil {
		return fmt.Errorf("failed to derive expected workflow owner for centralized source: %w", err)
	}
	derivedHex := hex.EncodeToString(derived)

	claimedOwner := strings.ToLower(strings.TrimPrefix(strings.TrimSpace(ownerHex), "0x"))
	if strings.ToLower(derivedHex) != claimedOwner {
		const criticalMismatchMsg = "centralized workflow owner does not match owner derived from its organization ID: possible data corruption or malicious workflow registry"
		logger.Sugared(lggr).Criticalw(
			criticalMismatchMsg,
			"source", source,
			"claimedWorkflowOwner", ownerHex,
			"organizationID", orgID,
			"derivedWorkflowOwner", derivedHex,
			"tenantID", tenantID,
		)
		return fmt.Errorf("%w: source=%s claimedOwner=%s organizationID=%s derivedOwner=%s",
			ErrCentralizedOwnerOrgMismatch, source, ownerHex, orgID, derivedHex)
	}
	return nil
}

// maybeVerifyCentralizedOwnerOrgMapping runs owner/org verification when the gate is open.
func maybeVerifyCentralizedOwnerOrgMapping(
	ctx context.Context,
	lggr logger.Logger,
	source, ownerHex, orgID string,
	gate limits.GateLimiter,
	getter settings.Getter,
) error {
	if !isCentralizedWorkflowSource(source) {
		return nil
	}

	gateErr := gate.AllowErr(ctx)
	if gateErr != nil {
		if errors.Is(gateErr, limits.ErrorNotAllowed{}) {
			return nil
		}
		return fmt.Errorf("failed to evaluate limit CentralizedWorkflowOwnerVerificationEnabled: %w", gateErr)
	}

	tenantID, err := tenantIDFromSettings(ctx, getter)
	if err != nil {
		return fmt.Errorf("failed to resolve tenant ID from cresettings: %w", err)
	}
	return verifyCentralizedOwnerOrgMapping(lggr, source, ownerHex, orgID, tenantID)
}
