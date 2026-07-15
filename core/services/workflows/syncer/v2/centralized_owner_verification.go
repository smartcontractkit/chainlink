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

// isCentralizedWorkflowSource reports whether workflow metadata originated from the
// centralized (remote) gRPC workflow registry, as opposed to the on-chain workflow
// registry contract (source prefix "contract:") or the local file source (source
// prefix "file:"). Owner/org verification only applies to the remote gRPC registry:
// the file source is an operator-local dev/test mechanism with no organization ID to
// verify against, and its contents are written by the node operator directly.
func isCentralizedWorkflowSource(source string) bool {
	return strings.HasPrefix(source, "grpc:")
}

// verifyCentralizedOwnerOrgMapping checks that a workflow owner claimed by a
// centralized source is consistent with the organization ID. Legitimate centralized
// CRE workflow owners are deterministically derived from (tenantID, orgID) via
// workflows.GenerateWorkflowOwnerAddress.
//
// This guards against a malicious centralized (gRPC) source claiming an arbitrary
// workflow owner - in particular an owner that actually belongs to an on-chain
// registry account - in order to hijack that owner's identity. Because a valid
// centralized owner must equal the address derived from (tenantID, orgID), and on-chain
// owners are arbitrary externally-owned/contract addresses that do not lie in this
// derived space, the derivation cannot be satisfied for an on-chain owner (doing so
// would require a keccak preimage). The check therefore confines a centralized source
// to owners within its own derived namespace, making it infeasible to spoof an on-chain
// owner through the gRPC source.
func verifyCentralizedOwnerOrgMapping(lggr logger.Logger, source, ownerHex, orgID string, tenantID uint64) error {
	if strings.TrimSpace(orgID) == "" {
		// Verification is enabled but the centralized source provided no organization ID,
		// so ownership cannot be verified. Fail closed rather than deriving an owner from an
		// empty org (which would surface as a misleading owner mismatch), and avoid letting a
		// source bypass verification by simply omitting the organization ID.
		logger.Sugared(lggr).Criticalw(
			"centralized workflow source did not provide an organization ID; cannot verify workflow owner",
			"source", source,
			"claimedWorkflowOwner", ownerHex,
			"tenantID", tenantID,
		)
		return fmt.Errorf("%w: source=%s claimedOwner=%s: missing organization ID", ErrCentralizedOwnerOrgMismatch, source, ownerHex)
	}

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
