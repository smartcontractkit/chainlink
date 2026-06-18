package vault

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-common/pkg/workflows"
)

// DeriveJWTAuthorizedVaultWorkflowOwner derives the JWT-authorized Vault workflow-owner address using
// the same inputs as cre-platform-graphql/internal/service/account_service.go GetCreOrganizationInfo:
//
//	workflows.GenerateWorkflowOwnerAddress(strconv.FormatUint(tenantID, 10), orgID)
//
// tenantID must be non-zero (equivalent to GraphQL rejecting a missing tenant context).
func DeriveJWTAuthorizedVaultWorkflowOwner(orgID string, tenantID uint64, claimedWorkflowOwnerFromJWT string) (string, error) {
	if orgID == "" {
		return "", errors.New("org_id is required for JWT-derived vault workflow owner")
	}
	if tenantID == 0 {
		return "", ErrMissingTenantID
	}
	prefix := strconv.FormatUint(tenantID, 10)
	addr, err := workflows.GenerateWorkflowOwnerAddress(prefix, orgID)
	if err != nil {
		return "", fmt.Errorf("could not derive vault workflow owner address: %w", err)
	}
	derived := common.BytesToAddress(addr).Hex()
	if claimedWorkflowOwnerFromJWT != "" && !strings.EqualFold(strings.TrimSpace(claimedWorkflowOwnerFromJWT), derived) {
		return "", fmt.Errorf("JWT workflow_owner claim %q does not match derived workflow owner %q", claimedWorkflowOwnerFromJWT, derived)
	}
	return derived, nil
}
