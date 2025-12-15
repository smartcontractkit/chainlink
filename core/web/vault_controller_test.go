package web_test

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/smdkg/dkgocr/dkgocrtypes"
	"github.com/smartcontractkit/smdkg/dkgocr/tdh2shim"
	"github.com/smartcontractkit/smdkg/dummydkg"

	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/dkgrecipientkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/vault"
	"github.com/smartcontractkit/chainlink/v2/core/web"
)

func setupVaultControllerTests(t *testing.T) (cltest.HTTPClientCleaner, keystore.Master, vault.ORM) {
	t.Helper()
	ctx := testutils.Context(t)

	app := cltest.NewApplication(t)
	require.NoError(t, app.Start(ctx))

	client := app.NewHTTPClient(nil)
	err := app.GetKeyStore().DKGRecipient().EnsureKey(ctx)
	require.NoError(t, err)

	orm := vault.NewVaultORM(app.GetDB())

	return client, app.GetKeyStore(), orm
}

func TestVaultController_VerifyDKGResult_HappyPath(t *testing.T) {
	client, keystore, orm := setupVaultControllerTests(t)

	keys, err := keystore.DKGRecipient().GetAll()
	require.NoError(t, err)
	require.Len(t, keys, 1)

	keyrings := []dkgocrtypes.P256Keyring{keys[0]}
	instanceID := dkgocrtypes.InstanceID("test-instance-id")
	rp, err := dummydkg.NewResultPackage(instanceID, dkgocrtypes.ReportingPluginConfig{
		DealerPublicKeys:    []dkgocrtypes.P256ParticipantPublicKey{keys[0].PublicKey()},
		RecipientPublicKeys: []dkgocrtypes.P256ParticipantPublicKey{keys[0].PublicKey()},
		T:                   1,
	}, keyrings)
	require.NoError(t, err)

	rpb, err := rp.MarshalBinary()
	require.NoError(t, err)

	var configDigest types.ConfigDigest
	copy(configDigest[:], common.Hex2Bytes("1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"))
	signatures := []types.AttributedOnchainSignature{
		{
			Signature: common.Hex2Bytes("deadbeef"),
			Signer:    commontypes.OracleID(1),
		},
		{
			Signature: common.Hex2Bytes("cafebabe"),
			Signer:    commontypes.OracleID(2),
		},
	}
	err = orm.WriteResultPackage(t.Context(), instanceID, dkgocrtypes.ResultPackageDatabaseValue{
		ConfigDigest:            configDigest,
		SeqNr:                   1,
		ReportWithResultPackage: rpb,
		Signatures:              signatures,
	})
	require.NoError(t, err)

	tdh2PubKey, err := tdh2shim.TDH2PublicKeyFromDKGResult(rp)
	require.NoError(t, err)

	tdh2PubKeyBytes, err := tdh2PubKey.Marshal()
	require.NoError(t, err)

	bdata, err := json.Marshal(web.VerifyDKGResultRequest{
		InstanceID:      string(instanceID),
		MasterPublicKey: hex.EncodeToString(tdh2PubKeyBytes),
	})
	require.NoError(t, err)

	response, cleanup := client.Post("/v2/vault/dkg_results/verify", bytes.NewReader(bdata))
	t.Cleanup(cleanup)
	cltest.AssertServerResponse(t, response, http.StatusOK)
}

func TestVaultController_VerifyDKGResult_WrongKey(t *testing.T) {
	client, keystore, orm := setupVaultControllerTests(t)

	keys, err := keystore.DKGRecipient().GetAll()
	require.NoError(t, err)
	require.Len(t, keys, 1)

	// Use another key than the one we set the node up with.
	// This will fail verification.
	key, err := dkgrecipientkey.New()
	require.NoError(t, err)

	keyrings := []dkgocrtypes.P256Keyring{key}
	instanceID := dkgocrtypes.InstanceID("test-instance-id")
	rp, err := dummydkg.NewResultPackage(instanceID, dkgocrtypes.ReportingPluginConfig{
		DealerPublicKeys:    []dkgocrtypes.P256ParticipantPublicKey{key.PublicKey()},
		RecipientPublicKeys: []dkgocrtypes.P256ParticipantPublicKey{key.PublicKey()},
		T:                   1,
	}, keyrings)
	require.NoError(t, err)

	rpb, err := rp.MarshalBinary()
	require.NoError(t, err)

	var configDigest types.ConfigDigest
	copy(configDigest[:], common.Hex2Bytes("1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"))
	signatures := []types.AttributedOnchainSignature{
		{
			Signature: common.Hex2Bytes("deadbeef"),
			Signer:    commontypes.OracleID(1),
		},
		{
			Signature: common.Hex2Bytes("cafebabe"),
			Signer:    commontypes.OracleID(2),
		},
	}
	err = orm.WriteResultPackage(t.Context(), instanceID, dkgocrtypes.ResultPackageDatabaseValue{
		ConfigDigest:            configDigest,
		SeqNr:                   1,
		ReportWithResultPackage: rpb,
		Signatures:              signatures,
	})
	require.NoError(t, err)

	tdh2PubKey, err := tdh2shim.TDH2PublicKeyFromDKGResult(rp)
	require.NoError(t, err)

	tdh2PubKeyBytes, err := tdh2PubKey.Marshal()
	require.NoError(t, err)

	bdata, err := json.Marshal(web.VerifyDKGResultRequest{
		InstanceID:      string(instanceID),
		MasterPublicKey: hex.EncodeToString(tdh2PubKeyBytes),
	})
	require.NoError(t, err)

	response, cleanup := client.Post("/v2/vault/dkg_results/verify", bytes.NewReader(bdata))
	t.Cleanup(cleanup)
	cltest.AssertServerResponse(t, response, http.StatusBadRequest)
}

func TestVaultController_VerifyDKGResult_CantFindResultForInstanceID(t *testing.T) {
	client, _, _ := setupVaultControllerTests(t)

	bdata, err := json.Marshal(web.VerifyDKGResultRequest{
		InstanceID:      string("instance-id"),
		MasterPublicKey: "deadbeef",
	})
	require.NoError(t, err)

	response, cleanup := client.Post("/v2/vault/dkg_results/verify", bytes.NewReader(bdata))
	t.Cleanup(cleanup)
	cltest.AssertServerResponse(t, response, http.StatusNotFound)
}

func TestVaultController_VerifyDKGResult_MissingInstanceIDOrPublicKey(t *testing.T) {
	client, _, _ := setupVaultControllerTests(t)

	bdata, err := json.Marshal(web.VerifyDKGResultRequest{
		InstanceID: string("instance-id"),
	})
	require.NoError(t, err)

	response, cleanup := client.Post("/v2/vault/dkg_results/verify", bytes.NewReader(bdata))
	t.Cleanup(cleanup)
	cltest.AssertServerResponse(t, response, http.StatusBadRequest)

	bdata, err = json.Marshal(web.VerifyDKGResultRequest{
		MasterPublicKey: "foo",
	})
	require.NoError(t, err)

	response, cleanup = client.Post("/v2/vault/dkg_results/verify", bytes.NewReader(bdata))
	t.Cleanup(cleanup)
	cltest.AssertServerResponse(t, response, http.StatusBadRequest)
}

func TestVaultController_ExportDKGResult(t *testing.T) {
	client, keystore, orm := setupVaultControllerTests(t)

	keys, err := keystore.DKGRecipient().GetAll()
	require.NoError(t, err)
	require.Len(t, keys, 1)

	keyrings := []dkgocrtypes.P256Keyring{keys[0]}
	instanceID := dkgocrtypes.InstanceID("test-instance-id")
	rp, err := dummydkg.NewResultPackage(instanceID, dkgocrtypes.ReportingPluginConfig{
		DealerPublicKeys:    []dkgocrtypes.P256ParticipantPublicKey{keys[0].PublicKey()},
		RecipientPublicKeys: []dkgocrtypes.P256ParticipantPublicKey{keys[0].PublicKey()},
		T:                   1,
	}, keyrings)
	require.NoError(t, err)

	rpb, err := rp.MarshalBinary()
	require.NoError(t, err)

	var configDigest types.ConfigDigest
	copy(configDigest[:], common.Hex2Bytes("1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"))
	signatures := []types.AttributedOnchainSignature{
		{
			Signature: common.Hex2Bytes("deadbeef"),
			Signer:    commontypes.OracleID(1),
		},
		{
			Signature: common.Hex2Bytes("cafebabe"),
			Signer:    commontypes.OracleID(2),
		},
	}
	err = orm.WriteResultPackage(t.Context(), instanceID, dkgocrtypes.ResultPackageDatabaseValue{
		ConfigDigest:            configDigest,
		SeqNr:                   1,
		ReportWithResultPackage: rpb,
		Signatures:              signatures,
	})
	require.NoError(t, err)

	bdata, err := json.Marshal(web.ExportDKGResultRequest{
		InstanceID: string(instanceID),
	})
	require.NoError(t, err)

	resp, cleanup := client.Post("/v2/vault/dkg_results/export", bytes.NewReader(bdata))
	t.Cleanup(cleanup)
	cltest.AssertServerResponse(t, resp, http.StatusOK)
}
