package config_test

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/dkg/config"
)

func TestValidatePluginConfig(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	db := pgtest.NewSqlxDB(t)
	kst := cltest.NewKeyStore(t, db)

	dkgEncryptKey, err := kst.DKGEncrypt().Create(ctx)
	require.NoError(t, err)
	dkgSignKey, err := kst.DKGSign().Create(ctx)
	require.NoError(t, err)

	encryptKeyBytes, err := dkgEncryptKey.PublicKey.MarshalBinary()
	require.NoError(t, err)
	encryptKey := hex.EncodeToString(encryptKeyBytes)

	signKeyBytes, err := dkgSignKey.PublicKey.MarshalBinary()
	require.NoError(t, err)
	signKey := hex.EncodeToString(signKeyBytes)

	pluginConfig := config.PluginConfig{
		EncryptionPublicKey: encryptKey,
		SigningPublicKey:    signKey,
	}
	t.Run("no error when keys are found", func(t *testing.T) {
		err = config.ValidatePluginConfig(pluginConfig, kst.DKGSign(), kst.DKGEncrypt())
		require.NoError(t, err)
	})

	t.Run("error when encryption key not found", func(t *testing.T) {
		pluginConfig = config.PluginConfig{
			EncryptionPublicKey: "wrongKey",
			SigningPublicKey:    signKey,
		}
		err = config.ValidatePluginConfig(pluginConfig, kst.DKGSign(), kst.DKGEncrypt())
		require.Error(t, err)
		require.Contains(t, err.Error(), "DKG encryption key: wrongKey not found in key store")
	})

	t.Run("error when sign key not found", func(t *testing.T) {
		pluginConfig = config.PluginConfig{
			EncryptionPublicKey: encryptKey,
			SigningPublicKey:    "wrongKey",
		}

		err = config.ValidatePluginConfig(pluginConfig, kst.DKGSign(), kst.DKGEncrypt())
		require.Error(t, err)
		require.Contains(t, err.Error(), "DKG sign key: wrongKey not found in key store")
	})
}
