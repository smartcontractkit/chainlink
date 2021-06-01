package presenters

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/manyminds/api2go/jsonapi"
	"github.com/smartcontractkit/chainlink/core/services/keystore/ocrkey"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestOCRKeysBundleResource(t *testing.T) {
	timestamp := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

	var (
		ocrKeyBundleID = "7f993fb701b3410b1f6e8d4d93a7462754d24609b9b31a4fe64a0cb475a4d934"
		password       = "p4SsW0rD1!@#_"
	)

	ocrKeyBundleIDSha256, err := models.Sha256HashFromHex(ocrKeyBundleID)
	require.NoError(t, err)

	pk, err := ocrkey.NewKeyBundle()
	require.NoError(t, err)
	pkEncrypted, err := pk.Encrypt(password, utils.FastScryptParams)
	require.NoError(t, err)

	bundle := ocrkey.EncryptedKeyBundle{
		ID:                    ocrKeyBundleIDSha256,
		OnChainSigningAddress: pkEncrypted.OnChainSigningAddress,
		OffChainPublicKey:     pkEncrypted.OffChainPublicKey,
		ConfigPublicKey:       pkEncrypted.ConfigPublicKey,
		CreatedAt:             timestamp,
		UpdatedAt:             timestamp,
	}

	r := NewOCRKeysBundleResource(bundle)
	b, err := jsonapi.Marshal(r)
	require.NoError(t, err)

	expected := fmt.Sprintf(`
	{
		"data":{
			"type":"encryptedKeyBundles",
			"id":"%s",
			"attributes":{
				"onChainSigningAddress": "%s",
				"offChainPublicKey": "%s",
				"configPublicKey": "%s",
				"createdAt":"2000-01-01T00:00:00Z",
				"updatedAt":"2000-01-01T00:00:00Z",
				"deletedAt":null
			}
		}
	}`,
		ocrKeyBundleID,
		pkEncrypted.OnChainSigningAddress.String(),
		pkEncrypted.OffChainPublicKey.String(),
		pkEncrypted.ConfigPublicKey.String(),
	)

	assert.JSONEq(t, expected, string(b))

	// With a deleted field
	bundle.DeletedAt = gorm.DeletedAt(sql.NullTime{Time: timestamp, Valid: true})

	r = NewOCRKeysBundleResource(bundle)
	b, err = jsonapi.Marshal(r)
	require.NoError(t, err)

	expected = fmt.Sprintf(`
	{
		"data": {
			"type":"encryptedKeyBundles",
			"id":"%s",
			"attributes":{
				"onChainSigningAddress": "%s",
				"offChainPublicKey": "%s",
				"configPublicKey": "%s",
				"createdAt":"2000-01-01T00:00:00Z",
				"updatedAt":"2000-01-01T00:00:00Z",
				"deletedAt":"2000-01-01T00:00:00Z"
			}
		}
	}`,
		ocrKeyBundleID,
		pkEncrypted.OnChainSigningAddress.String(),
		pkEncrypted.OffChainPublicKey.String(),
		pkEncrypted.ConfigPublicKey.String(),
	)

	assert.JSONEq(t, expected, string(b))
}
