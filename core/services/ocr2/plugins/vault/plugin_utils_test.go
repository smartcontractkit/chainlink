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
