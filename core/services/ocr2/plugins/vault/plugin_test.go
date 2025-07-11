package vault

import (
	"crypto/rand"
	"encoding/base64"
	"testing"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
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
		lggr:                           lggr,
		store:                          store,
		batchSize:                      10,
		publicKey:                      nil,
		privateKeyShare:                nil,
		maxSecretsPerOwner:             1,
		maxCiphertextLenBytes:          1024,
		maxIdentifierOwnerLenBytes:     100,
		maxIdentifierNamespaceLenBytes: 100,
		maxIdentifierKeyLenBytes:       100,
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
			err: "invalid secret identifier: owner exceeds maximum length of 3 bytes",
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
			lggr:                           lggr,
			store:                          store,
			batchSize:                      10,
			publicKey:                      nil,
			privateKeyShare:                nil,
			maxSecretsPerOwner:             1,
			maxCiphertextLenBytes:          1024,
			maxIdentifierOwnerLenBytes:     maxIdLen / 3,
			maxIdentifierNamespaceLenBytes: maxIdLen / 3,
			maxIdentifierKeyLenBytes:       maxIdLen / 3,
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

func TestPlugin_Observation_GetSecretsRequest_FillsInNamespace(t *testing.T) {
	lggr, _ := logger.TestLoggerObserved(t, zapcore.DebugLevel)
	store := requests.NewStore[*Request]()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := &ReportingPlugin{
		lggr:                           lggr,
		store:                          store,
		batchSize:                      10,
		publicKey:                      pk,
		privateKeyShare:                shares[0],
		maxSecretsPerOwner:             1,
		maxCiphertextLenBytes:          1024,
		maxIdentifierOwnerLenBytes:     100,
		maxIdentifierNamespaceLenBytes: 100,
		maxIdentifierKeyLenBytes:       100,
	}

	id := &vault.SecretIdentifier{
		Owner:     "owner",
		Namespace: "",
		Key:       "my_secret",
	}
	rdr := &kv{
		m: make(map[string]response),
	}

	plaintext := []byte("my-secret-value")
	ciphertext, err := tdh2easy.Encrypt(pk, plaintext)
	require.NoError(t, err)
	ciphertextBytes, err := ciphertext.Marshal()
	require.NoError(t, err)

	createdId := &vault.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "my_secret",
	}
	err = newWriteStore(rdr).writeSecret(createdId, &vault.StoredSecret{
		EncryptedSecret: ciphertextBytes,
	})
	require.NoError(t, err)

	pubK, _, err := box.GenerateKey(rand.Reader)
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

	assert.Equal(t, batchResp.Responses[0].Id, createdId)
}

func TestPlugin_Observation_GetSecretsRequest_SecretDoesNotExist(t *testing.T) {
	lggr := logger.TestLogger(t)
	store := requests.NewStore[*Request]()
	r := &ReportingPlugin{
		lggr:                           lggr,
		store:                          store,
		batchSize:                      10,
		publicKey:                      nil,
		privateKeyShare:                nil,
		maxSecretsPerOwner:             1,
		maxCiphertextLenBytes:          1024,
		maxIdentifierOwnerLenBytes:     100,
		maxIdentifierNamespaceLenBytes: 100,
		maxIdentifierKeyLenBytes:       100,
	}

	seqNr := uint64(1)
	rdr := &kv{
		m: make(map[string]response),
	}
	id := &vault.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "my_secret",
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
		lggr:                           lggr,
		store:                          store,
		batchSize:                      10,
		publicKey:                      pk,
		privateKeyShare:                shares[0],
		maxSecretsPerOwner:             1,
		maxCiphertextLenBytes:          1024,
		maxIdentifierOwnerLenBytes:     100,
		maxIdentifierNamespaceLenBytes: 100,
		maxIdentifierKeyLenBytes:       100,
	}

	id := &vault.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "my_secret",
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
		lggr:                           lggr,
		store:                          store,
		batchSize:                      10,
		publicKey:                      pk,
		privateKeyShare:                shares[0],
		maxSecretsPerOwner:             1,
		maxCiphertextLenBytes:          1024,
		maxIdentifierOwnerLenBytes:     100,
		maxIdentifierNamespaceLenBytes: 100,
		maxIdentifierKeyLenBytes:       100,
	}

	id := &vault.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "my_secret",
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
		lggr:                           lggr,
		store:                          store,
		batchSize:                      10,
		publicKey:                      pk,
		privateKeyShare:                shares[0],
		maxSecretsPerOwner:             1,
		maxCiphertextLenBytes:          1024,
		maxIdentifierOwnerLenBytes:     100,
		maxIdentifierNamespaceLenBytes: 100,
		maxIdentifierKeyLenBytes:       100,
	}

	id := &vault.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "my_secret",
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
			err: "invalid secret identifier: owner exceeds maximum length of 3 bytes",
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
			lggr:                           lggr,
			store:                          store,
			batchSize:                      10,
			publicKey:                      nil,
			privateKeyShare:                nil,
			maxSecretsPerOwner:             1,
			maxCiphertextLenBytes:          1024,
			maxIdentifierOwnerLenBytes:     maxIdLen / 3,
			maxIdentifierNamespaceLenBytes: maxIdLen / 3,
			maxIdentifierKeyLenBytes:       maxIdLen / 3,
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
		lggr:                           lggr,
		store:                          store,
		batchSize:                      10,
		publicKey:                      nil,
		privateKeyShare:                nil,
		maxSecretsPerOwner:             1,
		maxCiphertextLenBytes:          1024,
		maxIdentifierOwnerLenBytes:     100,
		maxIdentifierNamespaceLenBytes: 100,
		maxIdentifierKeyLenBytes:       100,
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
		lggr:                           lggr,
		store:                          store,
		batchSize:                      10,
		publicKey:                      nil,
		privateKeyShare:                nil,
		maxSecretsPerOwner:             1,
		maxCiphertextLenBytes:          10,
		maxIdentifierOwnerLenBytes:     100,
		maxIdentifierNamespaceLenBytes: 100,
		maxIdentifierKeyLenBytes:       100,
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
		lggr:                           lggr,
		store:                          store,
		batchSize:                      10,
		publicKey:                      pk,
		privateKeyShare:                shares[0],
		maxSecretsPerOwner:             1,
		maxCiphertextLenBytes:          1024,
		maxIdentifierOwnerLenBytes:     100,
		maxIdentifierNamespaceLenBytes: 100,
		maxIdentifierKeyLenBytes:       100,
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
		lggr:                           lggr,
		store:                          store,
		batchSize:                      10,
		publicKey:                      pk,
		privateKeyShare:                shares[0],
		maxSecretsPerOwner:             1,
		maxCiphertextLenBytes:          1024,
		maxIdentifierOwnerLenBytes:     100,
		maxIdentifierNamespaceLenBytes: 100,
		maxIdentifierKeyLenBytes:       100,
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
		lggr:                           lggr,
		store:                          store,
		batchSize:                      10,
		publicKey:                      pk,
		privateKeyShare:                shares[0],
		maxSecretsPerOwner:             1,
		maxCiphertextLenBytes:          1024,
		maxIdentifierOwnerLenBytes:     100,
		maxIdentifierNamespaceLenBytes: 100,
		maxIdentifierKeyLenBytes:       100,
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
		lggr:                           lggr,
		store:                          store,
		batchSize:                      10,
		publicKey:                      pk,
		privateKeyShare:                shares[0],
		maxSecretsPerOwner:             1,
		maxCiphertextLenBytes:          1024,
		maxIdentifierOwnerLenBytes:     100,
		maxIdentifierNamespaceLenBytes: 100,
		maxIdentifierKeyLenBytes:       100,
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

type observation struct {
	id   *vault.SecretIdentifier
	req  proto.Message
	resp proto.Message
}

func marshalObservations(t *testing.T, observations ...observation) []byte {
	obs := &vault.Observations{
		Observations: []*vault.Observation{},
	}
	for _, ob := range observations {
		o := &vault.Observation{
			Id: keyFor(ob.id),
		}
		switch ob.req.(type) {
		case *vault.GetSecretsRequest:
			o.RequestType = vault.RequestType_GET_SECRETS
			o.Request = &vault.Observation_GetSecretsRequest{
				GetSecretsRequest: ob.req.(*vault.GetSecretsRequest),
			}
		case *vault.CreateSecretsRequest:
			o.RequestType = vault.RequestType_CREATE_SECRETS
			o.Request = &vault.Observation_CreateSecretsRequest{
				CreateSecretsRequest: ob.req.(*vault.CreateSecretsRequest),
			}
		}

		switch ob.resp.(type) {
		case *vault.GetSecretsResponse:
			o.Response = &vault.Observation_GetSecretsResponse{
				GetSecretsResponse: ob.resp.(*vault.GetSecretsResponse),
			}
		case *vault.CreateSecretsResponse:
			o.Response = &vault.Observation_CreateSecretsResponse{
				CreateSecretsResponse: ob.resp.(*vault.CreateSecretsResponse),
			}
		}

		obs.Observations = append(obs.Observations, o)
	}

	b, err := proto.Marshal(obs)
	require.NoError(t, err)
	return b
}

func TestPlugin_StateTransition_InsufficientObservations(t *testing.T) {
	lggr, observed := logger.TestLoggerObserved(t, zapcore.DebugLevel)
	store := requests.NewStore[*Request]()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := &ReportingPlugin{
		lggr: lggr,
		config: ocr3types.ReportingPluginConfig{
			N: 4,
			F: 1,
		},
		store:                          store,
		batchSize:                      10,
		publicKey:                      pk,
		privateKeyShare:                shares[0],
		maxSecretsPerOwner:             1,
		maxCiphertextLenBytes:          1024,
		maxIdentifierOwnerLenBytes:     100,
		maxIdentifierNamespaceLenBytes: 100,
		maxIdentifierKeyLenBytes:       100,
	}

	seqNr := uint64(1)
	kv := &kv{
		m: make(map[string]response),
	}

	id1 := &vault.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "secret",
	}
	req := &vault.GetSecretsRequest{
		Requests: []*vault.SecretRequest{
			{
				Id: id1,
			},
		},
	}
	resp := &vault.GetSecretsResponse{
		Responses: []*vault.SecretResponse{
			{
				Id: id1,
				Result: &vault.SecretResponse_Error{
					Error: "key does not exist",
				},
			},
		},
	}

	obs1b := marshalObservations(t, observation{id1, req, resp})

	reportPrecursor, err := r.StateTransition(
		t.Context(),
		seqNr,
		types.AttributedQuery{},
		[]types.AttributedObservation{
			{Observation: types.Observation(obs1b)},
		}, kv, nil)
	require.NoError(t, err)

	os := &vault.Outcomes{}
	err = proto.Unmarshal(reportPrecursor, os)
	require.NoError(t, err)

	assert.Len(t, os.Outcomes, 0)

	assert.Equal(t, 1, observed.FilterMessage("insufficient observations found for id").Len())
}

func TestPlugin_StateTransition_InvalidObservations(t *testing.T) {
	lggr, observed := logger.TestLoggerObserved(t, zapcore.DebugLevel)
	store := requests.NewStore[*Request]()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := &ReportingPlugin{
		lggr: lggr,
		config: ocr3types.ReportingPluginConfig{
			N: 4,
			F: 1,
		},
		store:                          store,
		batchSize:                      10,
		publicKey:                      pk,
		privateKeyShare:                shares[0],
		maxSecretsPerOwner:             1,
		maxCiphertextLenBytes:          1024,
		maxIdentifierOwnerLenBytes:     100,
		maxIdentifierNamespaceLenBytes: 100,
		maxIdentifierKeyLenBytes:       100,
	}

	seqNr := uint64(1)
	kv := &kv{
		m: make(map[string]response),
	}

	id1 := &vault.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "secret",
	}
	req := &vault.GetSecretsRequest{
		Requests: []*vault.SecretRequest{
			{
				Id: id1,
			},
		},
	}
	resp := &vault.CreateSecretsResponse{}

	// Request and response don't match
	obsb := marshalObservations(t, observation{id1, req, resp})
	reportPrecursor, err := r.StateTransition(
		t.Context(),
		seqNr,
		types.AttributedQuery{},
		[]types.AttributedObservation{
			{Observation: types.Observation(obsb)},
		}, kv, nil)
	require.NoError(t, err)

	os := &vault.Outcomes{}
	err = proto.Unmarshal(reportPrecursor, os)
	require.NoError(t, err)

	assert.Len(t, os.Outcomes, 0)

	assert.Equal(t, 1, observed.FilterMessage("invalid observation").Len())

	// Invalid observation -- data can't be unmarshaled
	obsb = marshalObservations(t, observation{id1, req, resp})
	reportPrecursor, err = r.StateTransition(
		t.Context(),
		seqNr,
		types.AttributedQuery{},
		[]types.AttributedObservation{
			{Observation: types.Observation([]byte("hello world"))},
		}, kv, nil)
	require.NoError(t, err)

	os = &vault.Outcomes{}
	err = proto.Unmarshal(reportPrecursor, os)
	require.NoError(t, err)

	assert.Len(t, os.Outcomes, 0)

	assert.Equal(t, 1, observed.FilterMessage("invalid observation").Len())
}

func TestPlugin_StateTransition_ShasDontMatch(t *testing.T) {
	lggr, observed := logger.TestLoggerObserved(t, zapcore.DebugLevel)
	store := requests.NewStore[*Request]()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := &ReportingPlugin{
		lggr: lggr,
		config: ocr3types.ReportingPluginConfig{
			N: 4,
			F: 1,
		},
		store:                          store,
		batchSize:                      10,
		publicKey:                      pk,
		privateKeyShare:                shares[0],
		maxSecretsPerOwner:             1,
		maxCiphertextLenBytes:          1024,
		maxIdentifierOwnerLenBytes:     100,
		maxIdentifierNamespaceLenBytes: 100,
		maxIdentifierKeyLenBytes:       100,
	}

	seqNr := uint64(1)
	kv := &kv{
		m: make(map[string]response),
	}

	id := &vault.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "secret",
	}
	req := &vault.GetSecretsRequest{
		Requests: []*vault.SecretRequest{
			{
				Id: id,
			},
		},
	}
	resp1 := &vault.GetSecretsResponse{
		Responses: []*vault.SecretResponse{
			{
				Id: id,
				Result: &vault.SecretResponse_Error{
					Error: "key does not exist",
				},
			},
		},
	}
	resp2 := &vault.GetSecretsResponse{
		Responses: []*vault.SecretResponse{
			{
				Id: id,
				Result: &vault.SecretResponse_Error{
					Error: "something else",
				},
			},
		},
	}

	obsb := marshalObservations(t, observation{id, req, resp1}, observation{id, req, resp2}, observation{id, req, resp1})
	reportPrecursor, err := r.StateTransition(
		t.Context(),
		seqNr,
		types.AttributedQuery{},
		[]types.AttributedObservation{
			{Observation: types.Observation(obsb)},
		}, kv, nil)
	require.NoError(t, err)

	os := &vault.Outcomes{}
	err = proto.Unmarshal(reportPrecursor, os)
	require.NoError(t, err)

	assert.Len(t, os.Outcomes, 0)

	assert.Equal(t, 1, observed.FilterMessage("insufficient observations found for id").Len())
}

func TestPlugin_StateTransition_AggregatesValidationErrors(t *testing.T) {
	lggr, observed := logger.TestLoggerObserved(t, zapcore.DebugLevel)
	store := requests.NewStore[*Request]()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := &ReportingPlugin{
		lggr: lggr,
		config: ocr3types.ReportingPluginConfig{
			N: 4,
			F: 1,
		},
		store:                          store,
		batchSize:                      10,
		publicKey:                      pk,
		privateKeyShare:                shares[0],
		maxSecretsPerOwner:             1,
		maxCiphertextLenBytes:          1024,
		maxIdentifierOwnerLenBytes:     100,
		maxIdentifierNamespaceLenBytes: 100,
		maxIdentifierKeyLenBytes:       100,
	}

	seqNr := uint64(1)
	kv := &kv{
		m: make(map[string]response),
	}

	id := &vault.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "secret",
	}
	req := &vault.GetSecretsRequest{
		Requests: []*vault.SecretRequest{
			{
				Id: id,
			},
		},
	}
	resp := &vault.GetSecretsResponse{
		Responses: []*vault.SecretResponse{
			{
				Id: id,
				Result: &vault.SecretResponse_Error{
					Error: "key does not exist",
				},
			},
		},
	}

	obsb := marshalObservations(t, observation{id, req, resp}, observation{id, req, resp}, observation{id, req, resp})
	reportPrecursor, err := r.StateTransition(
		t.Context(),
		seqNr,
		types.AttributedQuery{},
		[]types.AttributedObservation{
			{Observation: types.Observation(obsb)},
		}, kv, nil)
	require.NoError(t, err)

	os := &vault.Outcomes{}
	err = proto.Unmarshal(reportPrecursor, os)
	require.NoError(t, err)

	assert.Len(t, os.Outcomes, 1)

	o := os.Outcomes[0]
	assert.True(t, proto.Equal(req, o.GetGetSecretsRequest()))
	assert.True(t, proto.Equal(resp, o.GetGetSecretsResponse()))

	assert.Equal(t, 1, observed.FilterMessage("sufficient observations for sha").Len())
}

func TestPlugin_StateTransition_GetSecretsRequest_CombinesShares(t *testing.T) {
	lggr, observed := logger.TestLoggerObserved(t, zapcore.DebugLevel)
	store := requests.NewStore[*Request]()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := &ReportingPlugin{
		lggr: lggr,
		config: ocr3types.ReportingPluginConfig{
			N: 4,
			F: 1,
		},
		store:                          store,
		batchSize:                      10,
		publicKey:                      pk,
		privateKeyShare:                shares[0],
		maxSecretsPerOwner:             1,
		maxCiphertextLenBytes:          1024,
		maxIdentifierOwnerLenBytes:     100,
		maxIdentifierNamespaceLenBytes: 100,
		maxIdentifierKeyLenBytes:       100,
	}

	seqNr := uint64(1)
	kv := &kv{
		m: make(map[string]response),
	}

	id := &vault.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "secret",
	}
	req := &vault.GetSecretsRequest{
		Requests: []*vault.SecretRequest{
			{
				Id: id,
			},
		},
	}
	resp1 := &vault.GetSecretsResponse{
		Responses: []*vault.SecretResponse{
			{
				Id: id,
				Result: &vault.SecretResponse_Data{
					Data: &vault.SecretData{
						EncryptedValue: "encrypted-value",
						EncryptedDecryptionKeyShares: []*vault.EncryptedShares{
							{
								EncryptionKey: "my-encryption-key",
								Shares:        []string{"encrypted-share-1"},
							},
						},
					},
				},
			},
		},
	}
	resp2 := &vault.GetSecretsResponse{
		Responses: []*vault.SecretResponse{
			{
				Id: id,
				Result: &vault.SecretResponse_Data{
					Data: &vault.SecretData{
						EncryptedValue: "encrypted-value",
						EncryptedDecryptionKeyShares: []*vault.EncryptedShares{
							{
								EncryptionKey: "my-encryption-key",
								Shares:        []string{"encrypted-share-2"},
							},
						},
					},
				},
			},
		},
	}
	resp3 := &vault.GetSecretsResponse{
		Responses: []*vault.SecretResponse{
			{
				Id: id,
				Result: &vault.SecretResponse_Data{
					Data: &vault.SecretData{
						EncryptedValue: "encrypted-value",
						EncryptedDecryptionKeyShares: []*vault.EncryptedShares{
							{
								EncryptionKey: "my-encryption-key",
								Shares:        []string{"encrypted-share-3"},
							},
						},
					},
				},
			},
		},
	}

	obsb := marshalObservations(t, observation{id, req, resp1}, observation{id, req, resp2}, observation{id, req, resp3})
	reportPrecursor, err := r.StateTransition(
		t.Context(),
		seqNr,
		types.AttributedQuery{},
		[]types.AttributedObservation{
			{Observation: types.Observation(obsb)},
		}, kv, nil)
	require.NoError(t, err)

	os := &vault.Outcomes{}
	err = proto.Unmarshal(reportPrecursor, os)
	require.NoError(t, err)

	assert.Len(t, os.Outcomes, 1)

	o := os.Outcomes[0]
	assert.True(t, proto.Equal(req, o.GetGetSecretsRequest()))

	expectedResp := &vault.GetSecretsResponse{
		Responses: []*vault.SecretResponse{
			{
				Id: id,
				Result: &vault.SecretResponse_Data{
					Data: &vault.SecretData{
						EncryptedValue: "encrypted-value",
						EncryptedDecryptionKeyShares: []*vault.EncryptedShares{
							{
								EncryptionKey: "my-encryption-key",
								Shares:        []string{"encrypted-share-1", "encrypted-share-2", "encrypted-share-3"},
							},
						},
					},
				},
			},
		},
	}
	assert.True(t, proto.Equal(expectedResp, o.GetGetSecretsResponse()), o.GetGetSecretsResponse())

	assert.Equal(t, 1, observed.FilterMessage("sufficient observations for sha").Len())
}

func TestPlugin_StateTransition_CreateSecretsRequest_WritesSecrets(t *testing.T) {
	lggr, observed := logger.TestLoggerObserved(t, zapcore.DebugLevel)
	store := requests.NewStore[*Request]()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := &ReportingPlugin{
		lggr: lggr,
		config: ocr3types.ReportingPluginConfig{
			N: 4,
			F: 1,
		},
		store:                          store,
		batchSize:                      10,
		publicKey:                      pk,
		privateKeyShare:                shares[0],
		maxSecretsPerOwner:             1,
		maxCiphertextLenBytes:          1024,
		maxIdentifierOwnerLenBytes:     100,
		maxIdentifierNamespaceLenBytes: 100,
		maxIdentifierKeyLenBytes:       100,
	}

	seqNr := uint64(1)
	kv := &kv{
		m: make(map[string]response),
	}
	rs := newReadStore(kv)

	id := &vault.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "secret",
	}
	value := []byte("encrypted-value")
	enc := base64.StdEncoding.EncodeToString(value)
	req := &vault.CreateSecretsRequest{
		EncryptedSecrets: []*vault.EncryptedSecret{
			{
				Id:             id,
				EncryptedValue: enc,
			},
		},
	}
	resp := &vault.CreateSecretsResponse{
		Responses: []*vault.CreateSecretResponse{
			{
				Id:      id,
				Success: false,
				Error:   "",
			},
		},
	}

	obsb := marshalObservations(t, observation{id, req, resp}, observation{id, req, resp}, observation{id, req, resp})
	reportPrecursor, err := r.StateTransition(
		t.Context(),
		seqNr,
		types.AttributedQuery{},
		[]types.AttributedObservation{
			{Observation: types.Observation(obsb)},
		}, kv, nil)
	require.NoError(t, err)

	os := &vault.Outcomes{}
	err = proto.Unmarshal(reportPrecursor, os)
	require.NoError(t, err)

	assert.Len(t, os.Outcomes, 1)

	o := os.Outcomes[0]
	assert.True(t, proto.Equal(req, o.GetCreateSecretsRequest()))

	expectedResp := &vault.CreateSecretsResponse{
		Responses: []*vault.CreateSecretResponse{
			{
				Id:      id,
				Success: true,
				Error:   "",
			},
		},
	}
	assert.True(t, proto.Equal(expectedResp, o.GetCreateSecretsResponse()), o.GetCreateSecretsResponse())

	ss, err := rs.getSecret(id)
	require.NoError(t, err)

	assert.Equal(t, ss.EncryptedSecret, []byte("encrypted-value"))

	assert.Equal(t, 1, observed.FilterMessage("sufficient observations for sha").Len())
}

func TestPlugin_Reports(t *testing.T) {
	value := "encrypted-value"
	id := &vault.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "secret",
	}
	req := &vault.CreateSecretsRequest{
		EncryptedSecrets: []*vault.EncryptedSecret{
			{
				Id:             id,
				EncryptedValue: value,
			},
		},
	}
	resp := &vault.CreateSecretsResponse{
		Responses: []*vault.CreateSecretResponse{
			{
				Id:      id,
				Success: false,
				Error:   "",
			},
		},
	}
	id2 := &vault.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "secret2",
	}
	req2 := &vault.CreateSecretsRequest{
		EncryptedSecrets: []*vault.EncryptedSecret{
			{
				Id:             id2,
				EncryptedValue: value,
			},
		},
	}
	resp2 := &vault.CreateSecretsResponse{
		Responses: []*vault.CreateSecretResponse{
			{
				Id:      id2,
				Success: false,
				Error:   "",
			},
		},
	}
	expectedOutcome1 := &vault.Outcome{
		Id:          keyFor(id),
		RequestType: vault.RequestType_CREATE_SECRETS,
		Request: &vault.Outcome_CreateSecretsRequest{
			CreateSecretsRequest: req,
		},
		Response: &vault.Outcome_CreateSecretsResponse{
			CreateSecretsResponse: resp,
		},
	}

	expectedOutcome2 := &vault.Outcome{
		Id:          keyFor(id2),
		RequestType: vault.RequestType_CREATE_SECRETS,
		Request: &vault.Outcome_CreateSecretsRequest{
			CreateSecretsRequest: req2,
		},
		Response: &vault.Outcome_CreateSecretsResponse{
			CreateSecretsResponse: resp2,
		},
	}
	os := &vault.Outcomes{
		Outcomes: []*vault.Outcome{
			expectedOutcome1,
			expectedOutcome2,
		},
	}

	osb, err := proto.Marshal(os)
	require.NoError(t, err)

	lggr, _ := logger.TestLoggerObserved(t, zapcore.DebugLevel)
	store := requests.NewStore[*Request]()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := &ReportingPlugin{
		lggr: lggr,
		config: ocr3types.ReportingPluginConfig{
			N: 4,
			F: 1,
		},
		store:                          store,
		batchSize:                      10,
		publicKey:                      pk,
		privateKeyShare:                shares[0],
		maxSecretsPerOwner:             1,
		maxCiphertextLenBytes:          1024,
		maxIdentifierOwnerLenBytes:     100,
		maxIdentifierNamespaceLenBytes: 100,
		maxIdentifierKeyLenBytes:       100,
	}

	rs, err := r.Reports(t.Context(), uint64(1), osb)
	require.NoError(t, err)

	assert.Len(t, rs, 2)

	o1b := rs[0]
	o1 := &vault.Outcome{}
	err = proto.Unmarshal(o1b.ReportWithInfo.Report, o1)
	require.NoError(t, err)
	assert.True(t, proto.Equal(o1, expectedOutcome1))

	o2b := rs[1]
	o2 := &vault.Outcome{}
	err = proto.Unmarshal(o2b.ReportWithInfo.Report, o2)
	require.NoError(t, err)
	assert.True(t, proto.Equal(o2, expectedOutcome2))
}
