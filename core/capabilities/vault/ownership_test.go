package vault_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	vault "github.com/smartcontractkit/chainlink/v2/core/capabilities/vault"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
)

func TestAuthorizedOwnerFromRequestID(t *testing.T) {
	require.Equal(t, "0xabc", vault.AuthorizedOwnerFromRequestID("0xabc"+vaulttypes.RequestIDSeparator+"req-1"))
	require.Empty(t, vault.AuthorizedOwnerFromRequestID("req-1"))
}

func TestStampAuthorizedOwnerOnCreate(t *testing.T) {
	req := &vaultcommon.CreateSecretsRequest{
		EncryptedSecrets: []*vaultcommon.EncryptedSecret{
			{Id: &vaultcommon.SecretIdentifier{Owner: "0xother", Key: "k"}, EncryptedValue: "cipher"},
		},
	}
	vault.StampAuthorizedOwnerOnCreate(req, "0xauthorized")
	require.Equal(t, "0xauthorized", req.EncryptedSecrets[0].Id.Owner)
}

func TestStampAuthorizedOwnerOnUpdate(t *testing.T) {
	req := &vaultcommon.UpdateSecretsRequest{
		EncryptedSecrets: []*vaultcommon.EncryptedSecret{
			{Id: &vaultcommon.SecretIdentifier{Owner: "0xother", Key: "k"}, EncryptedValue: "cipher"},
		},
	}
	vault.StampAuthorizedOwnerOnUpdate(req, "0xauthorized")
	require.Equal(t, "0xauthorized", req.EncryptedSecrets[0].Id.Owner)
}

func TestStampAuthorizedOwnerOnDelete(t *testing.T) {
	req := &vaultcommon.DeleteSecretsRequest{
		Ids: []*vaultcommon.SecretIdentifier{{Owner: "0xother", Key: "k"}},
	}
	vault.StampAuthorizedOwnerOnDelete(req, "0xauthorized")
	require.Equal(t, "0xauthorized", req.Ids[0].Owner)
}

func TestStampAuthorizedOwnerOnList(t *testing.T) {
	req := &vaultcommon.ListSecretIdentifiersRequest{Owner: "0xother", Namespace: "ns"}
	vault.StampAuthorizedOwnerOnList(req, "0xauthorized")
	require.Equal(t, "0xauthorized", req.Owner)
}

func TestStampAuthorizedOwnerFromRequestID_Create(t *testing.T) {
	authorized := "0xauthorized"
	req := &vaultcommon.CreateSecretsRequest{
		RequestId: authorized + vaulttypes.RequestIDSeparator + "req-1",
		EncryptedSecrets: []*vaultcommon.EncryptedSecret{
			{Id: &vaultcommon.SecretIdentifier{Owner: "0xforeign", Key: "k"}, EncryptedValue: "cipher"},
		},
	}
	vault.StampAuthorizedOwnerFromRequestID(req.RequestId, req)
	require.Equal(t, authorized, req.EncryptedSecrets[0].Id.Owner)
}

func TestStampAuthorizedOwnerFromRequestID_NoPrefixNoOp(t *testing.T) {
	req := &vaultcommon.CreateSecretsRequest{
		RequestId: "req-1",
		EncryptedSecrets: []*vaultcommon.EncryptedSecret{
			{Id: &vaultcommon.SecretIdentifier{Owner: "0xother", Key: "k"}, EncryptedValue: "cipher"},
		},
	}
	vault.StampAuthorizedOwnerFromRequestID(req.RequestId, req)
	require.Equal(t, "0xother", req.EncryptedSecrets[0].Id.Owner)
}
