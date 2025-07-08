package vault

import (
	"crypto/rand"
	"encoding/base64"
	"testing"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/tdh2/go/tdh2/tdh2easy"
	"go.uber.org/zap/zapcore"
	"golang.org/x/crypto/nacl/box"
	"google.golang.org/protobuf/proto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/requests"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestPlugin_Observation_NothingInBatch(t *testing.T) {
	lggr := logger.TestLogger(t)
	store := requests.NewStore[*Request]()
	r := &ReportingPlugin{
		lggr:                  lggr,
		store:                 store,
		batchSize:             10,
		publicKey:             nil,
		privateKeyShare:       nil,
		maxSecretsPerOwner:    1,
		maxCiphertextLenBytes: 1024,
		maxIdentifierLenBytes: 256,
	}

	seqNr := uint64(1)
	rdr := &kv{
		m: make(map[string]response),
	}
	data, err := r.Observation(t.Context(), seqNr, types.AttributedQuery{}, rdr, nil)
	require.NoError(t, err)

	obs := &vault.Observations{}
	err = proto.Unmarshal(data, obs)
	require.NoError(t, err)

	assert.Len(t, obs.Observations, 0)
}

func TestPlugin_Observation_GetSecretsRequest_SecretIdentifierInvalid(t *testing.T) {
	tcs := []struct {
		name     string
		id       *vault.SecretIdentifier
		maxIdLen int
		err      string
	}{
		{
			name: "nil id",
			id:   nil,
			err:  "invalid secret identifier: cannot be nil",
		},
		{
			name: "empty id",
			id:   &vault.SecretIdentifier{},
			err:  "invalid secret identifier: key cannot be empty",
		},
		{
			name: "empty id",
			id: &vault.SecretIdentifier{
				Key:       "hello",
				Namespace: "world",
			},
			err: "invalid secret identifier: owner cannot be empty",
		},
		{
			name:     "id is too long",
			maxIdLen: 10,
			id: &vault.SecretIdentifier{
				Owner:     "owner",
				Key:       "hello",
				Namespace: "world",
			},
			err: "invalid secret identifier: id exceeds maximum length of 10 bytes",
		},
		{
			name:     "contains key separator",
			maxIdLen: 10,
			id: &vault.SecretIdentifier{
				Owner:     ":",
				Key:       "hello",
				Namespace: "world",
			},
			err: "invalid secret identifier: id cannot contain `:`",
		},
	}

	for _, tc := range tcs {
		lggr := logger.TestLogger(t)
		store := requests.NewStore[*Request]()
		maxIdLen := 256
		if tc.maxIdLen > 0 {
			maxIdLen = tc.maxIdLen
		}
		r := &ReportingPlugin{
			lggr:                  lggr,
			store:                 store,
			batchSize:             10,
			publicKey:             nil,
			privateKeyShare:       nil,
			maxSecretsPerOwner:    1,
			maxCiphertextLenBytes: 1024,
			maxIdentifierLenBytes: maxIdLen,
		}

		seqNr := uint64(1)
		rdr := &kv{
			m: make(map[string]response),
		}
		p := &vault.GetSecretsRequest{
			Requests: []*vault.SecretRequest{
				{
					Id:             tc.id,
					EncryptionKeys: []string{"foo"},
				},
			},
		}
		err := store.Add(&Request{Payload: p})
		require.NoError(t, err)
		data, err := r.Observation(t.Context(), seqNr, types.AttributedQuery{}, rdr, nil)
		require.NoError(t, err)

		obs := &vault.Observations{}
		err = proto.Unmarshal(data, obs)
		require.NoError(t, err)

		assert.Len(t, obs.Observations, 1)
		o := obs.Observations[0]

		assert.Equal(t, o.RequestType, vault.RequestType_GET_SECRETS)
		assert.True(t, proto.Equal(o.GetGetSecretsRequest(), p))

		batchResp := o.GetGetSecretsResponse()
		assert.Len(t, p.Requests, 1)
		assert.Equal(t, len(p.Requests), len(batchResp.Responses))

		assert.True(t, proto.Equal(p.Requests[0].Id, batchResp.Responses[0].Id))
		resp := batchResp.Responses[0]
		assert.Contains(t, resp.GetError(), tc.err)
	}
}

func TestPlugin_Observation_GetSecretsRequest_SecretDoesNotExist(t *testing.T) {
	lggr := logger.TestLogger(t)
	store := requests.NewStore[*Request]()
	r := &ReportingPlugin{
		lggr:                  lggr,
		store:                 store,
		batchSize:             10,
		publicKey:             nil,
		privateKeyShare:       nil,
		maxSecretsPerOwner:    1,
		maxCiphertextLenBytes: 1024,
		maxIdentifierLenBytes: 100,
	}

	seqNr := uint64(1)
	rdr := &kv{
		m: make(map[string]response),
	}
	id := &vault.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "my-secret",
	}
	p := &vault.GetSecretsRequest{
		Requests: []*vault.SecretRequest{
			{
				Id:             id,
				EncryptionKeys: []string{"foo"},
			},
		},
	}
	err := store.Add(&Request{Payload: p})
	require.NoError(t, err)
	data, err := r.Observation(t.Context(), seqNr, types.AttributedQuery{}, rdr, nil)
	require.NoError(t, err)

	obs := &vault.Observations{}
	err = proto.Unmarshal(data, obs)
	require.NoError(t, err)

	assert.Len(t, obs.Observations, 1)
	o := obs.Observations[0]

	assert.Equal(t, o.RequestType, vault.RequestType_GET_SECRETS)
	assert.True(t, proto.Equal(o.GetGetSecretsRequest(), p))

	batchResp := o.GetGetSecretsResponse()
	assert.Len(t, p.Requests, 1)
	assert.Equal(t, len(p.Requests), len(batchResp.Responses))

	assert.True(t, proto.Equal(p.Requests[0].Id, batchResp.Responses[0].Id))
	resp := batchResp.Responses[0]
	assert.Contains(t, resp.GetError(), "key does not exist")
}

func TestPlugin_Observation_GetSecretsRequest_SecretExistsButIsIncorrect(t *testing.T) {
	lggr, observed := logger.TestLoggerObserved(t, zapcore.DebugLevel)
	store := requests.NewStore[*Request]()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := &ReportingPlugin{
		lggr:                  lggr,
		store:                 store,
		batchSize:             10,
		publicKey:             pk,
		privateKeyShare:       shares[0],
		maxSecretsPerOwner:    1,
		maxCiphertextLenBytes: 1024,
		maxIdentifierLenBytes: 100,
	}

	id := &vault.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "my-secret",
	}
	rdr := &kv{
		m: make(map[string]response),
	}

	err = newWriteStore(rdr).writeSecret(id, &vault.StoredSecret{
		EncryptedSecret: []byte("invalid-ciphertext"),
	})
	require.NoError(t, err)

	p := &vault.GetSecretsRequest{
		Requests: []*vault.SecretRequest{
			{
				Id:             id,
				EncryptionKeys: []string{"foo"},
			},
		},
	}
	err = store.Add(&Request{Payload: p})
	require.NoError(t, err)
	seqNr := uint64(1)
	data, err := r.Observation(t.Context(), seqNr, types.AttributedQuery{}, rdr, nil)
	require.NoError(t, err)

	obs := &vault.Observations{}
	err = proto.Unmarshal(data, obs)
	require.NoError(t, err)

	assert.Len(t, obs.Observations, 1)
	o := obs.Observations[0]

	assert.Equal(t, o.RequestType, vault.RequestType_GET_SECRETS)
	assert.True(t, proto.Equal(o.GetGetSecretsRequest(), p))

	batchResp := o.GetGetSecretsResponse()
	assert.Len(t, p.Requests, 1)
	assert.Equal(t, len(p.Requests), len(batchResp.Responses))

	assert.True(t, proto.Equal(p.Requests[0].Id, batchResp.Responses[0].Id))
	resp := batchResp.Responses[0]

	// The error returned is user-friendly
	assert.Contains(t, resp.GetError(), "failed to handle get secret request")

	// Inspect logs to get true source of error
	logs := observed.FilterMessage("failed to handle get secret request")
	assert.Equal(t, logs.Len(), 1)
	fields := logs.All()[0].ContextMap()
	errString := fields["error"]
	assert.Contains(t, errString, "failed to unmarshal ciphertext")
}

func TestPlugin_Observation_GetSecretsRequest_PublicKeyIsInvalid(t *testing.T) {
	lggr, _ := logger.TestLoggerObserved(t, zapcore.DebugLevel)
	store := requests.NewStore[*Request]()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := &ReportingPlugin{
		lggr:                  lggr,
		store:                 store,
		batchSize:             10,
		publicKey:             pk,
		privateKeyShare:       shares[0],
		maxSecretsPerOwner:    1,
		maxCiphertextLenBytes: 1024,
		maxIdentifierLenBytes: 100,
	}

	id := &vault.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "my-secret",
	}
	rdr := &kv{
		m: make(map[string]response),
	}

	plaintext := []byte("my-secret-value")
	ciphertext, err := tdh2easy.Encrypt(pk, plaintext)
	require.NoError(t, err)
	ciphertextBytes, err := ciphertext.Marshal()
	require.NoError(t, err)

	err = newWriteStore(rdr).writeSecret(id, &vault.StoredSecret{
		EncryptedSecret: ciphertextBytes,
	})
	require.NoError(t, err)

	p := &vault.GetSecretsRequest{
		Requests: []*vault.SecretRequest{
			{
				Id:             id,
				EncryptionKeys: []string{"foo"},
			},
		},
	}
	err = store.Add(&Request{Payload: p})
	require.NoError(t, err)
	seqNr := uint64(1)
	data, err := r.Observation(t.Context(), seqNr, types.AttributedQuery{}, rdr, nil)
	require.NoError(t, err)

	obs := &vault.Observations{}
	err = proto.Unmarshal(data, obs)
	require.NoError(t, err)

	assert.Len(t, obs.Observations, 1)
	o := obs.Observations[0]

	assert.Equal(t, o.RequestType, vault.RequestType_GET_SECRETS)
	assert.True(t, proto.Equal(o.GetGetSecretsRequest(), p))

	batchResp := o.GetGetSecretsResponse()
	assert.Len(t, p.Requests, 1)
	assert.Equal(t, len(p.Requests), len(batchResp.Responses))

	assert.True(t, proto.Equal(p.Requests[0].Id, batchResp.Responses[0].Id))
	resp := batchResp.Responses[0]

	assert.Contains(t, resp.GetError(), "failed to convert public key to bytes")
}

func TestPlugin_Observation_GetSecretsRequest_Success(t *testing.T) {
	lggr, _ := logger.TestLoggerObserved(t, zapcore.DebugLevel)
	store := requests.NewStore[*Request]()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := &ReportingPlugin{
		lggr:                  lggr,
		store:                 store,
		batchSize:             10,
		publicKey:             pk,
		privateKeyShare:       shares[0],
		maxSecretsPerOwner:    1,
		maxCiphertextLenBytes: 1024,
		maxIdentifierLenBytes: 100,
	}

	id := &vault.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "my-secret",
	}
	rdr := &kv{
		m: make(map[string]response),
	}

	plaintext := []byte("my-secret-value")
	ciphertext, err := tdh2easy.Encrypt(pk, plaintext)
	require.NoError(t, err)
	ciphertextBytes, err := ciphertext.Marshal()
	require.NoError(t, err)

	err = newWriteStore(rdr).writeSecret(id, &vault.StoredSecret{
		EncryptedSecret: ciphertextBytes,
	})
	require.NoError(t, err)

	pubK, privK, err := box.GenerateKey(rand.Reader)
	require.NoError(t, err)

	pks := base64.StdEncoding.EncodeToString(pubK[:])

	p := &vault.GetSecretsRequest{
		Requests: []*vault.SecretRequest{
			{
				Id:             id,
				EncryptionKeys: []string{pks},
			},
		},
	}
	err = store.Add(&Request{Payload: p})
	require.NoError(t, err)
	seqNr := uint64(1)
	data, err := r.Observation(t.Context(), seqNr, types.AttributedQuery{}, rdr, nil)
	require.NoError(t, err)

	obs := &vault.Observations{}
	err = proto.Unmarshal(data, obs)
	require.NoError(t, err)

	assert.Len(t, obs.Observations, 1)
	o := obs.Observations[0]

	assert.Equal(t, o.RequestType, vault.RequestType_GET_SECRETS)
	assert.True(t, proto.Equal(o.GetGetSecretsRequest(), p))

	batchResp := o.GetGetSecretsResponse()
	assert.Len(t, p.Requests, 1)
	assert.Equal(t, len(p.Requests), len(batchResp.Responses))

	assert.True(t, proto.Equal(p.Requests[0].Id, batchResp.Responses[0].Id))
	resp := batchResp.Responses[0]

	assert.Empty(t, resp.GetError())

	assert.Equal(t, base64.StdEncoding.EncodeToString(ciphertextBytes), resp.GetData().EncryptedValue)

	assert.Len(t, resp.GetData().EncryptedDecryptionKeyShares, 1)
	shareString := resp.GetData().EncryptedDecryptionKeyShares[0].Shares[0]

	share, err := base64.StdEncoding.DecodeString(shareString)
	require.NoError(t, err)
	msg, ok := box.OpenAnonymous(nil, share, pubK, privK)
	assert.True(t, ok)

	ds := &tdh2easy.DecryptionShare{}
	err = ds.Unmarshal(msg)
	require.NoError(t, err)

	ct := &tdh2easy.Ciphertext{}
	ctb, err := base64.StdEncoding.DecodeString(resp.GetData().EncryptedValue)
	require.NoError(t, err)
	err = ct.UnmarshalVerify(ctb, pk)
	require.NoError(t, err)

	gotSecret, err := tdh2easy.Aggregate(ct, []*tdh2easy.DecryptionShare{ds}, 3)
	require.NoError(t, err)

	assert.Equal(t, plaintext, gotSecret)
}

func TestPlugin_Observation_CreateSecretsRequest_SecretIdentifierInvalid(t *testing.T) {
	tcs := []struct {
		name     string
		id       *vault.SecretIdentifier
		maxIdLen int
		err      string
	}{
		{
			name: "nil id",
			id:   nil,
			err:  "invalid secret identifier: cannot be nil",
		},
		{
			name: "empty id",
			id:   &vault.SecretIdentifier{},
			err:  "invalid secret identifier: key cannot be empty",
		},
		{
			name: "empty id",
			id: &vault.SecretIdentifier{
				Key:       "hello",
				Namespace: "world",
			},
			err: "invalid secret identifier: owner cannot be empty",
		},
		{
			name:     "id is too long",
			maxIdLen: 10,
			id: &vault.SecretIdentifier{
				Owner:     "owner",
				Key:       "hello",
				Namespace: "world",
			},
			err: "invalid secret identifier: id exceeds maximum length of 10 bytes",
		},
		{
			name:     "contains key separator",
			maxIdLen: 10,
			id: &vault.SecretIdentifier{
				Owner:     ":",
				Key:       "hello",
				Namespace: "world",
			},
			err: "invalid secret identifier: id cannot contain `:`",
		},
	}

	for _, tc := range tcs {
		lggr := logger.TestLogger(t)
		store := requests.NewStore[*Request]()
		maxIdLen := 256
		if tc.maxIdLen > 0 {
			maxIdLen = tc.maxIdLen
		}
		r := &ReportingPlugin{
			lggr:                  lggr,
			store:                 store,
			batchSize:             10,
			publicKey:             nil,
			privateKeyShare:       nil,
			maxSecretsPerOwner:    1,
			maxCiphertextLenBytes: 1024,
			maxIdentifierLenBytes: maxIdLen,
		}

		seqNr := uint64(1)
		rdr := &kv{
			m: make(map[string]response),
		}
		p := &vault.CreateSecretsRequest{
			EncryptedSecrets: []*vault.EncryptedSecret{
				{
					Id:             tc.id,
					EncryptedValue: "foo",
				},
			},
		}
		err := store.Add(&Request{Payload: p})
		require.NoError(t, err)
		data, err := r.Observation(t.Context(), seqNr, types.AttributedQuery{}, rdr, nil)
		require.NoError(t, err)

		obs := &vault.Observations{}
		err = proto.Unmarshal(data, obs)
		require.NoError(t, err)

		assert.Len(t, obs.Observations, 1)
		o := obs.Observations[0]

		assert.Equal(t, o.RequestType, vault.RequestType_CREATE_SECRETS)
		assert.True(t, proto.Equal(o.GetCreateSecretsRequest(), p))

		batchResp := o.GetCreateSecretsResponse()
		assert.Len(t, p.EncryptedSecrets, 1)
		assert.Equal(t, len(p.EncryptedSecrets), len(batchResp.Responses))

		assert.True(t, proto.Equal(p.EncryptedSecrets[0].Id, batchResp.Responses[0].Id))
		resp := batchResp.Responses[0]
		assert.Contains(t, resp.GetError(), tc.err)
	}
}

func TestPlugin_Observation_CreateSecretsRequest_InvalidCiphertext(t *testing.T) {
	lggr := logger.TestLogger(t)
	store := requests.NewStore[*Request]()
	r := &ReportingPlugin{
		lggr:                  lggr,
		store:                 store,
		batchSize:             10,
		publicKey:             nil,
		privateKeyShare:       nil,
		maxSecretsPerOwner:    1,
		maxCiphertextLenBytes: 1024,
		maxIdentifierLenBytes: 100,
	}

	seqNr := uint64(1)
	rdr := &kv{
		m: make(map[string]response),
	}

	id := &vault.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "secret",
	}
	p := &vault.CreateSecretsRequest{
		EncryptedSecrets: []*vault.EncryptedSecret{
			{
				Id:             id,
				EncryptedValue: "foo",
			},
		},
	}
	err := store.Add(&Request{Payload: p})
	require.NoError(t, err)
	data, err := r.Observation(t.Context(), seqNr, types.AttributedQuery{}, rdr, nil)
	require.NoError(t, err)

	obs := &vault.Observations{}
	err = proto.Unmarshal(data, obs)
	require.NoError(t, err)

	assert.Len(t, obs.Observations, 1)
	o := obs.Observations[0]

	assert.Equal(t, o.RequestType, vault.RequestType_CREATE_SECRETS)
	assert.True(t, proto.Equal(o.GetCreateSecretsRequest(), p))

	batchResp := o.GetCreateSecretsResponse()
	assert.Len(t, p.EncryptedSecrets, 1)
	assert.Equal(t, len(p.EncryptedSecrets), len(batchResp.Responses))

	assert.True(t, proto.Equal(p.EncryptedSecrets[0].Id, batchResp.Responses[0].Id))
	resp := batchResp.Responses[0]
	assert.Contains(t, resp.GetError(), "invalid base64 encoding for ciphertext")
}

func TestPlugin_Observation_CreateSecretsRequest_InvalidCiphertext_TooLong(t *testing.T) {
	lggr := logger.TestLogger(t)
	store := requests.NewStore[*Request]()
	r := &ReportingPlugin{
		lggr:                  lggr,
		store:                 store,
		batchSize:             10,
		publicKey:             nil,
		privateKeyShare:       nil,
		maxSecretsPerOwner:    1,
		maxCiphertextLenBytes: 10,
		maxIdentifierLenBytes: 100,
	}

	seqNr := uint64(1)
	rdr := &kv{
		m: make(map[string]response),
	}

	id := &vault.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "secret",
	}
	ciphertext := []byte("a quick brown fox jumps over the lazy dog")
	p := &vault.CreateSecretsRequest{
		EncryptedSecrets: []*vault.EncryptedSecret{
			{
				Id:             id,
				EncryptedValue: base64.StdEncoding.EncodeToString(ciphertext),
			},
		},
	}
	err := store.Add(&Request{Payload: p})
	require.NoError(t, err)
	data, err := r.Observation(t.Context(), seqNr, types.AttributedQuery{}, rdr, nil)
	require.NoError(t, err)

	obs := &vault.Observations{}
	err = proto.Unmarshal(data, obs)
	require.NoError(t, err)

	assert.Len(t, obs.Observations, 1)
	o := obs.Observations[0]

	assert.Equal(t, o.RequestType, vault.RequestType_CREATE_SECRETS)
	assert.True(t, proto.Equal(o.GetCreateSecretsRequest(), p))

	batchResp := o.GetCreateSecretsResponse()
	assert.Len(t, p.EncryptedSecrets, 1)
	assert.Equal(t, len(p.EncryptedSecrets), len(batchResp.Responses))

	assert.True(t, proto.Equal(p.EncryptedSecrets[0].Id, batchResp.Responses[0].Id))
	resp := batchResp.Responses[0]
	assert.Contains(t, resp.GetError(), "ciphertext size exceeds maximum allowed size: 10 bytes")
}

func TestPlugin_Observation_CreateSecretsRequest_InvalidCiphertext_EncryptedWithWrongPublicKey(t *testing.T) {
	lggr := logger.TestLogger(t)
	store := requests.NewStore[*Request]()
	// Wrong key
	_, wrongPublicKey, _, err := tdh2easy.GenerateKeys(1, 3)
	// Right key
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := &ReportingPlugin{
		lggr:                  lggr,
		store:                 store,
		batchSize:             10,
		publicKey:             pk,
		privateKeyShare:       shares[0],
		maxSecretsPerOwner:    1,
		maxCiphertextLenBytes: 1024,
		maxIdentifierLenBytes: 100,
	}

	seqNr := uint64(1)
	rdr := &kv{
		m: make(map[string]response),
	}

	id := &vault.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "secret",
	}
	ct, err := tdh2easy.Encrypt(wrongPublicKey, []byte("my secret value"))
	require.NoError(t, err)

	ciphertextBytes, err := ct.Marshal()
	require.NoError(t, err)

	p := &vault.CreateSecretsRequest{
		EncryptedSecrets: []*vault.EncryptedSecret{
			{
				Id:             id,
				EncryptedValue: base64.StdEncoding.EncodeToString(ciphertextBytes),
			},
		},
	}
	err = store.Add(&Request{Payload: p})
	require.NoError(t, err)
	data, err := r.Observation(t.Context(), seqNr, types.AttributedQuery{}, rdr, nil)
	require.NoError(t, err)

	obs := &vault.Observations{}
	err = proto.Unmarshal(data, obs)
	require.NoError(t, err)

	assert.Len(t, obs.Observations, 1)
	o := obs.Observations[0]

	assert.Equal(t, o.RequestType, vault.RequestType_CREATE_SECRETS)
	assert.True(t, proto.Equal(o.GetCreateSecretsRequest(), p))

	batchResp := o.GetCreateSecretsResponse()
	assert.Len(t, p.EncryptedSecrets, 1)
	assert.Equal(t, len(p.EncryptedSecrets), len(batchResp.Responses))

	assert.True(t, proto.Equal(p.EncryptedSecrets[0].Id, batchResp.Responses[0].Id))
	resp := batchResp.Responses[0]
	assert.Contains(t, resp.GetError(), "failed to verify ciphertext")
}

func TestPlugin_Observation_CreateSecretsRequest_SecretExistsForKey(t *testing.T) {
	lggr := logger.TestLogger(t)
	store := requests.NewStore[*Request]()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := &ReportingPlugin{
		lggr:                  lggr,
		store:                 store,
		batchSize:             10,
		publicKey:             pk,
		privateKeyShare:       shares[0],
		maxSecretsPerOwner:    1,
		maxCiphertextLenBytes: 1024,
		maxIdentifierLenBytes: 100,
	}

	seqNr := uint64(1)
	rdr := &kv{
		m: make(map[string]response),
	}
	id := &vault.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "secret",
	}
	err = newWriteStore(rdr).writeSecret(id, &vault.StoredSecret{EncryptedSecret: []byte("already exists")})
	require.NoError(t, err)

	ct, err := tdh2easy.Encrypt(pk, []byte("my secret value"))
	require.NoError(t, err)

	ciphertextBytes, err := ct.Marshal()
	require.NoError(t, err)

	p := &vault.CreateSecretsRequest{
		EncryptedSecrets: []*vault.EncryptedSecret{
			{
				Id:             id,
				EncryptedValue: base64.StdEncoding.EncodeToString(ciphertextBytes),
			},
		},
	}
	err = store.Add(&Request{Payload: p})
	require.NoError(t, err)
	data, err := r.Observation(t.Context(), seqNr, types.AttributedQuery{}, rdr, nil)
	require.NoError(t, err)

	obs := &vault.Observations{}
	err = proto.Unmarshal(data, obs)
	require.NoError(t, err)

	assert.Len(t, obs.Observations, 1)
	o := obs.Observations[0]

	assert.Equal(t, o.RequestType, vault.RequestType_CREATE_SECRETS)
	assert.True(t, proto.Equal(o.GetCreateSecretsRequest(), p))

	batchResp := o.GetCreateSecretsResponse()
	assert.Len(t, p.EncryptedSecrets, 1)
	assert.Equal(t, len(p.EncryptedSecrets), len(batchResp.Responses))

	assert.True(t, proto.Equal(p.EncryptedSecrets[0].Id, batchResp.Responses[0].Id))
	resp := batchResp.Responses[0]
	assert.Contains(t, resp.GetError(), "key already exists")
}

func TestPlugin_Observation_CreateSecretsRequest_TooManySecretsForOwner(t *testing.T) {
	lggr := logger.TestLogger(t)
	store := requests.NewStore[*Request]()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := &ReportingPlugin{
		lggr:                  lggr,
		store:                 store,
		batchSize:             10,
		publicKey:             pk,
		privateKeyShare:       shares[0],
		maxSecretsPerOwner:    1,
		maxCiphertextLenBytes: 1024,
		maxIdentifierLenBytes: 100,
	}

	seqNr := uint64(1)
	rdr := &kv{
		m: make(map[string]response),
	}
	id := &vault.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "secret",
	}
	kvstore := newWriteStore(rdr)
	err = kvstore.writeMetadata(id.Owner, &vault.StoredMetadata{
		Keys: []string{"foo"},
	})
	require.NoError(t, err)

	ct, err := tdh2easy.Encrypt(pk, []byte("my secret value"))
	require.NoError(t, err)

	ciphertextBytes, err := ct.Marshal()
	require.NoError(t, err)

	p := &vault.CreateSecretsRequest{
		EncryptedSecrets: []*vault.EncryptedSecret{
			{
				Id:             id,
				EncryptedValue: base64.StdEncoding.EncodeToString(ciphertextBytes),
			},
		},
	}
	err = store.Add(&Request{Payload: p})
	require.NoError(t, err)
	data, err := r.Observation(t.Context(), seqNr, types.AttributedQuery{}, rdr, nil)
	require.NoError(t, err)

	obs := &vault.Observations{}
	err = proto.Unmarshal(data, obs)
	require.NoError(t, err)

	assert.Len(t, obs.Observations, 1)
	o := obs.Observations[0]

	assert.Equal(t, o.RequestType, vault.RequestType_CREATE_SECRETS)
	assert.True(t, proto.Equal(o.GetCreateSecretsRequest(), p))

	batchResp := o.GetCreateSecretsResponse()
	assert.Len(t, p.EncryptedSecrets, 1)
	assert.Equal(t, len(p.EncryptedSecrets), len(batchResp.Responses))

	assert.True(t, proto.Equal(p.EncryptedSecrets[0].Id, batchResp.Responses[0].Id))
	resp := batchResp.Responses[0]
	assert.Contains(t, resp.GetError(), "maximum number of secrets per owner reached")
}

func TestPlugin_Observation_CreateSecretsRequest_Success(t *testing.T) {
	lggr := logger.TestLogger(t)
	store := requests.NewStore[*Request]()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := &ReportingPlugin{
		lggr:                  lggr,
		store:                 store,
		batchSize:             10,
		publicKey:             pk,
		privateKeyShare:       shares[0],
		maxSecretsPerOwner:    1,
		maxCiphertextLenBytes: 1024,
		maxIdentifierLenBytes: 100,
	}

	seqNr := uint64(1)
	rdr := &kv{
		m: make(map[string]response),
	}
	id := &vault.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "secret",
	}
	ct, err := tdh2easy.Encrypt(pk, []byte("my secret value"))
	require.NoError(t, err)

	ciphertextBytes, err := ct.Marshal()
	require.NoError(t, err)

	p := &vault.CreateSecretsRequest{
		EncryptedSecrets: []*vault.EncryptedSecret{
			{
				Id:             id,
				EncryptedValue: base64.StdEncoding.EncodeToString(ciphertextBytes),
			},
		},
	}
	err = store.Add(&Request{Payload: p})
	require.NoError(t, err)
	data, err := r.Observation(t.Context(), seqNr, types.AttributedQuery{}, rdr, nil)
	require.NoError(t, err)

	obs := &vault.Observations{}
	err = proto.Unmarshal(data, obs)
	require.NoError(t, err)

	assert.Len(t, obs.Observations, 1)
	o := obs.Observations[0]

	assert.Equal(t, o.RequestType, vault.RequestType_CREATE_SECRETS)
	assert.True(t, proto.Equal(o.GetCreateSecretsRequest(), p))

	batchResp := o.GetCreateSecretsResponse()
	assert.Len(t, p.EncryptedSecrets, 1)
	assert.Equal(t, len(p.EncryptedSecrets), len(batchResp.Responses))

	assert.True(t, proto.Equal(p.EncryptedSecrets[0].Id, batchResp.Responses[0].Id))
	resp := batchResp.Responses[0]

	assert.Empty(t, resp.GetError())
}
