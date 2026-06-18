package vault

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	linkingclient "github.com/smartcontractkit/chainlink-protos/linking-service/go/v1"
	vaultcap "github.com/smartcontractkit/chainlink/v2/core/capabilities/vault"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestTestJWTIssuer_WorksWithVaultJWTBasedAuth(t *testing.T) {
	issuer, err := NewTestJWTIssuer()
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, issuer.Close())
	})

	owner, err := vaultcap.DeriveJWTAuthorizedVaultWorkflowOwner("org-test", 1, "")
	require.NoError(t, err)

	params, err := json.Marshal(vaultcommon.ListSecretIdentifiersRequest{
		Namespace: "main",
		Owner:     owner,
	})
	require.NoError(t, err)

	req := jsonrpc.Request[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      "req-1",
		Method:  vaulttypes.MethodSecretsList,
		Params:  (*json.RawMessage)(&params),
	}

	requestDigest, err := ComputeRequestDigest(req)
	require.NoError(t, err)

	token, err := issuer.MintToken(JWTTokenClaims{
		KeyID:         DefaultJWTIssuerKeyID,
		Issuer:        issuer.LocalIssuerURL(),
		Audience:      "https://api.test.chain.link",
		TenantID:      1,
		OrgID:         "org-test",
		WorkflowOwner: owner,
		RequestDigest: requestDigest,
		Scopes:        []string{vaultcap.OAuthScopeVaultSecretsList},
	})
	require.NoError(t, err)

	req.Auth = token

	auth, err := vaultcap.NewJWTBasedAuth(vaultcap.JWTBasedAuthConfig{
		IssuerURL: issuer.LocalIssuerURL(),
		Audience:  "https://api.test.chain.link",
		TenantID:  1,
	}, limits.Factory{Settings: cresettings.DefaultGetter}, logger.TestLogger(t), vaultcap.WithJWTBasedAuthGateLimiter(limits.NewGateLimiter(true)))
	require.NoError(t, err)

	authResult, err := auth.AuthorizeRequest(t.Context(), req)
	require.NoError(t, err)
	require.Equal(t, "org-test", authResult.OrgID())
	require.Equal(t, owner, authResult.WorkflowOwner())
	require.Equal(t, requestDigest, authResult.Digest())
}

func TestTestLinkingService_ResolvesOwner(t *testing.T) {
	svc, err := NewTestLinkingService(map[string]string{
		"0xAbC": "org-123",
		"d6d4fc38c209f53caa5a311a0cb44259daa4e9e1": "org-456",
	})
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, svc.Close())
	})

	conn, err := grpc.NewClient(svc.LocalURL(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, conn.Close())
	})

	resp, err := linkingclient.NewLinkingServiceClient(conn).GetOrganizationFromWorkflowOwner(t.Context(), &linkingclient.GetOrganizationFromWorkflowOwnerRequest{
		WorkflowOwner: "0xabc",
	})
	require.NoError(t, err)
	require.Equal(t, "org-123", resp.GetOrganizationId())

	resp, err = linkingclient.NewLinkingServiceClient(conn).GetOrganizationFromWorkflowOwner(t.Context(), &linkingclient.GetOrganizationFromWorkflowOwnerRequest{
		WorkflowOwner: "0xD6d4fC38c209F53caa5a311a0cb44259dAA4E9e1",
	})
	require.NoError(t, err)
	require.Equal(t, "org-456", resp.GetOrganizationId())
}
