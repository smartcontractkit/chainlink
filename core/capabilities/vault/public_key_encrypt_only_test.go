package vault

import (
	"encoding/hex"
	"encoding/json"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/tdh2/go/tdh2/tdh2easy"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/requests"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/registry"
	"github.com/smartcontractkit/chainlink-common/pkg/settings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestEncryptOnlyPublicKey(t *testing.T) {
	t.Parallel()

	_, pk, _, err := tdh2easy.GenerateKeys(2, 3)
	require.NoError(t, err)
	full, err := pk.Marshal()
	require.NoError(t, err)

	// Sanity: the full key carries HArray.
	var fm map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(full, &fm))
	require.Contains(t, fm, "HArray")

	out, err := encryptOnlyPublicKey(full)
	require.NoError(t, err)

	var m map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(out, &m))
	assert.NotContains(t, m, "HArray", "HArray must be stripped")
	assert.Contains(t, m, "Group")
	assert.Contains(t, m, "G_bar")
	assert.Contains(t, m, "H")

	// The stripped key still unmarshals and can encrypt.
	encPk := tdh2easy.PublicKey{}
	require.NoError(t, encPk.Unmarshal(out))
	_, err = tdh2easy.Encrypt(&encPk, []byte("hello"))
	require.NoError(t, err)

	// Invalid JSON input errors.
	_, err = encryptOnlyPublicKey([]byte("not json"))
	require.Error(t, err)
}

func TestCapability_GetPublicKey_EncryptOnlyGate(t *testing.T) {
	t.Parallel()

	_, pk, _, err := tdh2easy.GenerateKeys(2, 3)
	require.NoError(t, err)
	lazy := NewLazyPublicKey()
	lazy.Set(pk)

	newCap := func(t *testing.T, getter settings.Getter) *Capability {
		lggr := logger.TestLogger(t)
		clock := clockwork.NewFakeClock()
		expiry := 10 * time.Second
		store := requests.NewStore[*vaulttypes.Request]()
		handler := requests.NewHandler(lggr, store, clock, expiry)
		reg := registry.NewRegistry(lggr)
		lf := limits.Factory{Settings: getter}
		c, cerr := NewCapability(lggr, clock, expiry, handler, reg, lazy, lf, newTestRequestLifecycleTracker(t))
		require.NoError(t, cerr)
		return c
	}

	hasHArray := func(t *testing.T, hexKey string) bool {
		pkb, derr := hex.DecodeString(hexKey)
		require.NoError(t, derr)
		var m map[string]json.RawMessage
		require.NoError(t, json.Unmarshal(pkb, &m))
		_, ok := m["HArray"]
		return ok
	}

	t.Run("gate closed returns full key with HArray", func(t *testing.T) {
		t.Parallel()
		c := newCap(t, cresettings.DefaultGetter)
		resp, gerr := c.GetPublicKey(t.Context(), &vault.GetPublicKeyRequest{})
		require.NoError(t, gerr)
		assert.True(t, hasHArray(t, resp.PublicKey), "full key must retain HArray when gate closed")
	})

	t.Run("gate open returns encrypt-only key without HArray", func(t *testing.T) {
		t.Parallel()
		getter, gerr := settings.NewJSONGetter([]byte(`{"global":{"VaultPublicKeyEncryptOnlyEnabled":"true"}}`))
		require.NoError(t, gerr)
		c := newCap(t, getter)
		resp, rerr := c.GetPublicKey(t.Context(), &vault.GetPublicKeyRequest{})
		require.NoError(t, rerr)
		assert.False(t, hasHArray(t, resp.PublicKey), "encrypt-only key must not carry HArray when gate open")

		// The returned key is still usable to encrypt.
		pkb, derr := hex.DecodeString(resp.PublicKey)
		require.NoError(t, derr)
		encPk := tdh2easy.PublicKey{}
		require.NoError(t, encPk.Unmarshal(pkb))
		_, eerr := tdh2easy.Encrypt(&encPk, []byte("hello"))
		require.NoError(t, eerr)
	})
}
