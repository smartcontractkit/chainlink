package vault

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/nacl/box"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/tdh2/go/tdh2/tdh2easy"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/dkgrecipientkey"
	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault/vaultcrypto"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/requests"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	vaultcap "github.com/smartcontractkit/chainlink/v2/core/capabilities/vault"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaultutils"
)

type fakeReadTxn struct {
	*kv
	seqNr    uint64
	discards *int
}

func (f fakeReadTxn) SeqNr() uint64 { return f.seqNr }
func (f fakeReadTxn) Discard()      { *f.discards++ }

type fakeReadOnlyKV struct {
	kv       *kv
	openErr  error
	opens    int
	discards int
}

func (f *fakeReadOnlyKV) NewReadTransaction(context.Context) (ocr3_1types.KeyValueStateReadTransaction, error) {
	if f.openErr != nil {
		return nil, f.openErr
	}
	f.opens++
	return fakeReadTxn{kv: f.kv, seqNr: 7, discards: &f.discards}, nil
}

type directReadFixture struct {
	plugin    *ReportingPlugin
	kv        *fakeReadOnlyKV
	pk        *tdh2easy.PublicKey
	owner     string
	recipient struct{ pub, priv *[32]byte }
}

func newDirectReadFixture(t *testing.T) *directReadFixture {
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)

	f := &directReadFixture{
		kv:    &fakeReadOnlyKV{kv: &kv{m: make(map[string]response)}},
		pk:    pk,
		owner: common.HexToAddress("0x1111111111111111111111111111111111111111").Hex(),
	}
	f.recipient.pub, f.recipient.priv, err = box.GenerateKey(rand.Reader)
	require.NoError(t, err)

	f.plugin = newTestReportingPlugin(t, withKeys(pk, shares[0]), withMaxSecretsPerOwner(10))
	f.plugin.readOnlyKV = f.kv
	return f
}

func (f *directReadFixture) writeSecret(t *testing.T, key, value string) *vaultcommon.SecretIdentifier {
	id := &vaultcommon.SecretIdentifier{Owner: f.owner, Namespace: "main", Key: key}
	f.writeSecretLabelledFor(t, id, f.owner, value)
	return id
}

func (f *directReadFixture) writeSecretLabelledFor(t *testing.T, id *vaultcommon.SecretIdentifier, labelOwner, value string) {
	encrypted, err := vaultutils.EncryptSecretWithWorkflowOwner(value, f.pk, common.HexToAddress(labelOwner))
	require.NoError(t, err)
	ct, err := hex.DecodeString(encrypted)
	require.NoError(t, err)
	require.NoError(t, newTestWriteStore(t, f.kv.kv).WriteSecret(t.Context(), id, &vaultcommon.StoredSecret{EncryptedSecret: ct}))
}

func (f *directReadFixture) request(ids ...*vaultcommon.SecretIdentifier) *vaultcommon.GetSecretsRequest {
	req := &vaultcommon.GetSecretsRequest{GetSecretsDirectly: true}
	for _, id := range ids {
		req.Requests = append(req.Requests, &vaultcommon.SecretRequest{
			Id:             id,
			EncryptionKeys: []string{hex.EncodeToString(f.recipient.pub[:])},
		})
	}
	return req
}

func (f *directReadFixture) decrypt(t *testing.T, resp *vaultcommon.SecretResponse) string {
	ct, err := hex.DecodeString(resp.GetData().GetEncryptedValue())
	require.NoError(t, err)
	shares, err := vaultcrypto.SharesForKey(resp.GetData().GetEncryptedDecryptionKeyShares(), hex.EncodeToString(f.recipient.pub[:]))
	require.NoError(t, err)
	d := vaultcrypto.DecryptFunc(func(b []byte) ([]byte, error) {
		out, ok := box.OpenAnonymous(nil, b, f.recipient.pub, f.recipient.priv)
		if !ok {
			return nil, errors.New("failed to open box")
		}
		return out, nil
	})
	plaintext, _, err := vaultcrypto.DecryptSecret(ct, f.pk, 1, shares, d)
	require.NoError(t, err)
	return string(plaintext)
}

func TestGetSecretsDirect(t *testing.T) {
	t.Parallel()
	f := newDirectReadFixture(t)
	id := f.writeSecret(t, "present", "my-secret-value")
	missing := &vaultcommon.SecretIdentifier{Owner: f.owner, Namespace: "main", Key: "missing"}

	resp, err := f.plugin.GetSecretsDirect(t.Context(), f.request(id, missing))
	require.NoError(t, err)
	require.Len(t, resp.Responses, 2)

	assert.Equal(t, id.Key, resp.Responses[0].GetId().GetKey())
	assert.Empty(t, resp.Responses[0].GetError())
	assert.Equal(t, "my-secret-value", f.decrypt(t, resp.Responses[0]))

	assert.Equal(t, "key does not exist", resp.Responses[1].GetError())

	assert.Equal(t, 1, f.kv.opens, "all items are read from a single transaction")
	assert.Equal(t, 1, f.kv.discards, "the transaction is discarded")
}

func TestGetSecretsDirect_ItemValidationErrorsMatchOCRPath(t *testing.T) {
	t.Parallel()
	f := newDirectReadFixture(t)
	id := f.writeSecret(t, "dup", "value")

	resp, err := f.plugin.GetSecretsDirect(t.Context(), f.request(id, id, nil))
	require.NoError(t, err)
	require.Len(t, resp.Responses, 3)
	assert.Contains(t, resp.Responses[0].GetError(), "duplicate request for secret identifier")
	assert.Contains(t, resp.Responses[1].GetError(), "duplicate request for secret identifier")
	assert.Equal(t, "secret identifier cannot be nil", resp.Responses[2].GetError())
	assert.Equal(t, 1, f.kv.discards)
}

func TestGetSecretsDirect_WrongLabel(t *testing.T) {
	t.Parallel()
	f := newDirectReadFixture(t)
	// The stored ciphertext is labelled for a different owner than the identifier's.
	otherID := &vaultcommon.SecretIdentifier{Owner: f.owner, Namespace: "main", Key: "k"}
	f.writeSecretLabelledFor(t, otherID, "0x2222222222222222222222222222222222222222", "value")

	resp, err := f.plugin.GetSecretsDirect(t.Context(), f.request(otherID))
	require.NoError(t, err)
	require.Len(t, resp.Responses, 1)
	assert.NotEmpty(t, resp.Responses[0].GetError())
	assert.Nil(t, resp.Responses[0].GetData())
}

func TestGetSecretsDirect_NotReady(t *testing.T) {
	t.Parallel()
	f := newDirectReadFixture(t)
	f.plugin.readOnlyKV = nil

	_, err := f.plugin.GetSecretsDirect(t.Context(), f.request())
	require.ErrorIs(t, err, vaultcap.ErrDirectReadNotReady)
}

func TestGetSecretsDirect_OpenTransactionFails(t *testing.T) {
	t.Parallel()
	f := newDirectReadFixture(t)
	f.kv.openErr = errors.New("boom")

	_, err := f.plugin.GetSecretsDirect(t.Context(), f.request(f.writeSecret(t, "k", "v")))
	require.ErrorContains(t, err, "failed to open read transaction: boom")
	assert.Equal(t, 0, f.kv.discards)
}

func TestReportingPlugin_CloseClearsDirectReader(t *testing.T) {
	t.Parallel()
	holder := vaultcap.NewLazyDirectSecretsReader()

	old := newDirectReadFixture(t).plugin
	old.directReader = holder
	holder.Set(old)

	replacement := newDirectReadFixture(t).plugin
	replacement.directReader = holder
	holder.Set(replacement)

	require.NoError(t, old.Close())
	assert.Same(t, replacement, holder.Get(), "closing a stale instance must not clear its replacement")

	require.NoError(t, replacement.Close())
	assert.Nil(t, holder.Get())
}

func TestPlugin_ReportingPluginFactory_RegistersDirectReader(t *testing.T) {
	lggr := logger.Test(t)
	_, orm := setupORM(t)
	dkgrecipientKey, err := dkgrecipientkey.New()
	require.NoError(t, err)
	instanceID := "instanceID"
	_ = writeDKGPackage(t, orm, dkgrecipientKey, instanceID)

	rpf, err := NewReportingPluginFactory(lggr, requests.NewStore[*vaulttypes.Request](), orm, &dkgrecipientKey, vaultcap.NewLazyPublicKey(), limits.Factory{Settings: cresettings.DefaultGetter}, testRequestLifecycleTracker(t, lggr))
	require.NoError(t, err)
	holder := vaultcap.NewLazyDirectSecretsReader()
	rpf.SetDirectSecretsReader(holder)

	cfgb, err := proto.Marshal(&vaultcommon.ReportingPluginConfig{DKGInstanceID: &instanceID})
	require.NoError(t, err)
	onchainCfg := ocr3types.ReportingPluginConfig{OffchainConfig: cfgb, N: 10, F: 3}

	_, _, err = rpf.NewReportingPlugin(t.Context(), onchainCfg, nil, nil)
	require.NoError(t, err)
	assert.Nil(t, holder.Get(), "no read-only state, no direct reads")

	rp, _, err := rpf.NewReportingPlugin(t.Context(), onchainCfg, nil, &fakeReadOnlyKV{kv: &kv{m: map[string]response{}}})
	require.NoError(t, err)
	assert.Same(t, rp, holder.Get())

	require.NoError(t, rp.Close())
	assert.Nil(t, holder.Get())
}
