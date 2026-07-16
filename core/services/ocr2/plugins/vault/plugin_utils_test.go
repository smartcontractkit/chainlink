package vault

import (
	"testing"

	"github.com/stretchr/testify/require"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
)

func TestValidateEncryptedSharesEntry(t *testing.T) {
	t.Run("hex share", func(t *testing.T) {
		require.NoError(t, validateEncryptedSharesEntry(&vaultcommon.EncryptedShares{
			Shares: []string{"abcd"},
		}))
	})

	t.Run("binary share", func(t *testing.T) {
		require.NoError(t, validateEncryptedSharesEntry(&vaultcommon.EncryptedShares{
			BinaryShares: [][]byte{{1, 2, 3}},
		}))
	})

	t.Run("rejects empty", func(t *testing.T) {
		require.ErrorContains(t, validateEncryptedSharesEntry(&vaultcommon.EncryptedShares{}), "exactly 1 share")
	})

	t.Run("rejects both encodings", func(t *testing.T) {
		require.ErrorContains(t, validateEncryptedSharesEntry(&vaultcommon.EncryptedShares{
			Shares:       []string{"abcd"},
			BinaryShares: [][]byte{{1}},
		}), "exactly 1 share")
	})
}

func TestEncryptedShareSizeForLimit(t *testing.T) {
	t.Run("hex share", func(t *testing.T) {
		n, err := encryptedShareSizeForLimit(&vaultcommon.EncryptedShares{Shares: []string{"abcdef"}})
		require.NoError(t, err)
		require.Equal(t, 6, n)
	})

	t.Run("binary share", func(t *testing.T) {
		n, err := encryptedShareSizeForLimit(&vaultcommon.EncryptedShares{BinaryShares: [][]byte{{1, 2, 3, 4}}})
		require.NoError(t, err)
		require.Equal(t, 4, n)
	})

	t.Run("no share", func(t *testing.T) {
		_, err := encryptedShareSizeForLimit(&vaultcommon.EncryptedShares{})
		require.ErrorContains(t, err, "no share to measure")
	})
}

func TestAppendEncryptedShareEntry(t *testing.T) {
	t.Run("appends hex share", func(t *testing.T) {
		dst := &vaultcommon.EncryptedShares{EncryptionKey: "k"}
		appendEncryptedShareEntry(dst, &vaultcommon.EncryptedShares{Shares: []string{"a"}})
		require.Equal(t, []string{"a"}, dst.Shares)
		require.Empty(t, dst.BinaryShares)
	})

	t.Run("appends binary share", func(t *testing.T) {
		dst := &vaultcommon.EncryptedShares{EncryptionKey: "k"}
		appendEncryptedShareEntry(dst, &vaultcommon.EncryptedShares{BinaryShares: [][]byte{{1, 2}}})
		require.Len(t, dst.BinaryShares, 1)
		require.Equal(t, []byte{1, 2}, dst.BinaryShares[0])
		require.Empty(t, dst.Shares)
	})
}

func TestSecretRequestForID(t *testing.T) {
	t.Parallel()
	id := &vaultcommon.SecretIdentifier{Owner: "owner", Namespace: "main", Key: "secret"}
	req := &vaultcommon.GetSecretsRequest{
		Requests: []*vaultcommon.SecretRequest{{Id: id, EncryptionKeys: []string{"k1"}}},
	}
	got, err := secretRequestForID(req, id)
	require.NoError(t, err)
	require.Equal(t, "k1", got.EncryptionKeys[0])

	_, err = secretRequestForID(req, &vaultcommon.SecretIdentifier{Owner: "owner", Namespace: "main", Key: "missing"})
	require.ErrorContains(t, err, "no secret request")

	_, err = secretRequestForID(nil, id)
	require.ErrorContains(t, err, "GetSecrets request is nil")
}

func TestValidateGetSecretsShareLabels(t *testing.T) {
	t.Parallel()
	secretReq := &vaultcommon.SecretRequest{
		Id:             &vaultcommon.SecretIdentifier{Owner: "owner", Namespace: "main", Key: "secret"},
		EncryptionKeys: []string{"key-a", "key-b"},
	}

	t.Run("valid labels", func(t *testing.T) {
		t.Parallel()
		err := validateGetSecretsShareLabels(secretReq, &vaultcommon.SecretData{
			EncryptedDecryptionKeyShares: []*vaultcommon.EncryptedShares{
				{EncryptionKey: "key-a", Shares: []string{"s1"}},
				{EncryptionKey: "key-b", Shares: []string{"s2"}},
			},
		})
		require.NoError(t, err)
	})

	t.Run("bogus label", func(t *testing.T) {
		t.Parallel()
		err := validateGetSecretsShareLabels(secretReq, &vaultcommon.SecretData{
			EncryptedDecryptionKeyShares: []*vaultcommon.EncryptedShares{
				{EncryptionKey: "bogus", Shares: []string{"s1"}},
				{EncryptionKey: "key-b", Shares: []string{"s2"}},
			},
		})
		require.ErrorContains(t, err, "unexpected encryption key")
	})

	t.Run("missing label", func(t *testing.T) {
		t.Parallel()
		err := validateGetSecretsShareLabels(secretReq, &vaultcommon.SecretData{
			EncryptedDecryptionKeyShares: []*vaultcommon.EncryptedShares{
				{EncryptionKey: "key-a", Shares: []string{"s1"}},
			},
		})
		require.ErrorContains(t, err, "expected 2 encrypted share entries")
	})
}
