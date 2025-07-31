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
		lggr:  lggr,
		store: store,
		cfg: &ReportingPluginConfig{
			BatchSize:                      10,
			PublicKey:                      nil,
			PrivateKeyShare:                nil,
			MaxSecretsPerOwner:             1,
			MaxCiphertextLenBytes:          1024,
			MaxIdentifierOwnerLenBytes:     100,
			MaxIdentifierNamespaceLenBytes: 100,
			MaxIdentifierKeyLenBytes:       100,
		},
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

	assert.Empty(t, obs.Observations)
}

func TestPlugin_Observation_GetSecretsRequest_SecretIdentifierInvalid(t *testing.T) {
	tcs := []struct {
		name     string
		id       *vault.SecretIdentifier
		maxIDLen int
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
			maxIDLen: 10,
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
		maxIDLen := 256
		if tc.maxIDLen > 0 {
			maxIDLen = tc.maxIDLen
		}
		r := &ReportingPlugin{
			lggr:  lggr,
			store: store,
			cfg: &ReportingPluginConfig{
				BatchSize:                      10,
				PublicKey:                      nil,
				PrivateKeyShare:                nil,
				MaxSecretsPerOwner:             1,
				MaxCiphertextLenBytes:          1024,
				MaxIdentifierOwnerLenBytes:     maxIDLen / 3,
				MaxIdentifierNamespaceLenBytes: maxIDLen / 3,
				MaxIdentifierKeyLenBytes:       maxIDLen / 3,
			},
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

		assert.Equal(t, vault.RequestType_GET_SECRETS, o.RequestType)
		assert.True(t, proto.Equal(o.GetGetSecretsRequest(), p))

		batchResp := o.GetGetSecretsResponse()
		assert.Len(t, p.Requests, 1)
		assert.Len(t, p.Requests, len(batchResp.Responses))

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
		lggr:  lggr,
		store: store,
		cfg: &ReportingPluginConfig{
			BatchSize:                      10,
			PublicKey:                      pk,
			PrivateKeyShare:                shares[0],
			MaxSecretsPerOwner:             1,
			MaxCiphertextLenBytes:          1024,
			MaxIdentifierOwnerLenBytes:     100,
			MaxIdentifierNamespaceLenBytes: 100,
			MaxIdentifierKeyLenBytes:       100,
		},
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

	createdID := &vault.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "my_secret",
	}
	err = NewWriteStore(rdr).WriteSecret(createdID, &vault.StoredSecret{
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

	assert.Equal(t, vault.RequestType_GET_SECRETS, o.RequestType)
	assert.True(t, proto.Equal(o.GetGetSecretsRequest(), p))

	batchResp := o.GetGetSecretsResponse()
	assert.Len(t, p.Requests, 1)
	assert.Len(t, p.Requests, len(batchResp.Responses))

	assert.True(t, proto.Equal(batchResp.Responses[0].Id, createdID))
}

func TestPlugin_Observation_GetSecretsRequest_SecretDoesNotExist(t *testing.T) {
	lggr := logger.TestLogger(t)
	store := requests.NewStore[*Request]()
	r := &ReportingPlugin{
		lggr:  lggr,
		store: store,
		cfg: &ReportingPluginConfig{
			BatchSize:                      10,
			PublicKey:                      nil,
			PrivateKeyShare:                nil,
			MaxSecretsPerOwner:             1,
			MaxCiphertextLenBytes:          1024,
			MaxIdentifierOwnerLenBytes:     100,
			MaxIdentifierNamespaceLenBytes: 100,
			MaxIdentifierKeyLenBytes:       100,
		},
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

	assert.Equal(t, vault.RequestType_GET_SECRETS, o.RequestType)
	assert.True(t, proto.Equal(o.GetGetSecretsRequest(), p))

	batchResp := o.GetGetSecretsResponse()
	assert.Len(t, p.Requests, 1)
	assert.Len(t, p.Requests, len(batchResp.Responses))

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
		lggr:  lggr,
		store: store,
		cfg: &ReportingPluginConfig{
			BatchSize:                      10,
			PublicKey:                      pk,
			PrivateKeyShare:                shares[0],
			MaxSecretsPerOwner:             1,
			MaxCiphertextLenBytes:          1024,
			MaxIdentifierOwnerLenBytes:     100,
			MaxIdentifierNamespaceLenBytes: 100,
			MaxIdentifierKeyLenBytes:       100,
		},
	}

	id := &vault.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "my_secret",
	}
	rdr := &kv{
		m: make(map[string]response),
	}

	err = NewWriteStore(rdr).WriteSecret(id, &vault.StoredSecret{
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

	assert.Equal(t, vault.RequestType_GET_SECRETS, o.RequestType)
	assert.True(t, proto.Equal(o.GetGetSecretsRequest(), p))

	batchResp := o.GetGetSecretsResponse()
	assert.Len(t, p.Requests, 1)
	assert.Len(t, p.Requests, len(batchResp.Responses))

	assert.True(t, proto.Equal(p.Requests[0].Id, batchResp.Responses[0].Id))
	resp := batchResp.Responses[0]

	// The error returned is user-friendly
	assert.Contains(t, resp.GetError(), "failed to handle get secret request")

	// Inspect logs to get true source of error
	logs := observed.FilterMessage("failed to handle get secret request")
	assert.Equal(t, 1, logs.Len())
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
		lggr:  lggr,
		store: store,
		cfg: &ReportingPluginConfig{
			BatchSize:                      10,
			PublicKey:                      pk,
			PrivateKeyShare:                shares[0],
			MaxSecretsPerOwner:             1,
			MaxCiphertextLenBytes:          1024,
			MaxIdentifierOwnerLenBytes:     100,
			MaxIdentifierNamespaceLenBytes: 100,
			MaxIdentifierKeyLenBytes:       100,
		},
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

	err = NewWriteStore(rdr).WriteSecret(id, &vault.StoredSecret{
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

	assert.Equal(t, vault.RequestType_GET_SECRETS, o.RequestType)
	assert.True(t, proto.Equal(o.GetGetSecretsRequest(), p))

	batchResp := o.GetGetSecretsResponse()
	assert.Len(t, p.Requests, 1)
	assert.Len(t, p.Requests, len(batchResp.Responses))

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
		lggr:  lggr,
		store: store,
		cfg: &ReportingPluginConfig{
			BatchSize:                      10,
			PublicKey:                      pk,
			PrivateKeyShare:                shares[0],
			MaxSecretsPerOwner:             1,
			MaxCiphertextLenBytes:          1024,
			MaxIdentifierOwnerLenBytes:     100,
			MaxIdentifierNamespaceLenBytes: 100,
			MaxIdentifierKeyLenBytes:       100,
		},
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

	err = NewWriteStore(rdr).WriteSecret(id, &vault.StoredSecret{
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

	assert.Equal(t, vault.RequestType_GET_SECRETS, o.RequestType)
	assert.True(t, proto.Equal(o.GetGetSecretsRequest(), p))

	batchResp := o.GetGetSecretsResponse()
	assert.Len(t, p.Requests, 1)
	assert.Len(t, p.Requests, len(batchResp.Responses))

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
		maxIDLen int
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
			maxIDLen: 10,
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
		maxIDLen := 256
		if tc.maxIDLen > 0 {
			maxIDLen = tc.maxIDLen
		}
		r := &ReportingPlugin{
			lggr:  lggr,
			store: store,
			cfg: &ReportingPluginConfig{
				BatchSize:                      10,
				PublicKey:                      nil,
				PrivateKeyShare:                nil,
				MaxSecretsPerOwner:             1,
				MaxCiphertextLenBytes:          1024,
				MaxIdentifierOwnerLenBytes:     maxIDLen / 3,
				MaxIdentifierNamespaceLenBytes: maxIDLen / 3,
				MaxIdentifierKeyLenBytes:       maxIDLen / 3,
			},
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

		assert.Equal(t, vault.RequestType_CREATE_SECRETS, o.RequestType)
		assert.True(t, proto.Equal(o.GetCreateSecretsRequest(), p))

		batchResp := o.GetCreateSecretsResponse()
		assert.Len(t, p.EncryptedSecrets, 1)
		assert.Len(t, p.EncryptedSecrets, len(batchResp.Responses))

		assert.True(t, proto.Equal(p.EncryptedSecrets[0].Id, batchResp.Responses[0].Id))
		resp := batchResp.Responses[0]
		assert.Contains(t, resp.GetError(), tc.err)
	}
}

func TestPlugin_Observation_CreateSecretsRequest_DisallowsDuplicateRequests(t *testing.T) {
	lggr := logger.TestLogger(t)
	store := requests.NewStore[*Request]()
	r := &ReportingPlugin{
		lggr:  lggr,
		store: store,
		cfg: &ReportingPluginConfig{
			BatchSize:                      10,
			PublicKey:                      nil,
			PrivateKeyShare:                nil,
			MaxSecretsPerOwner:             1,
			MaxCiphertextLenBytes:          1024,
			MaxIdentifierOwnerLenBytes:     30,
			MaxIdentifierNamespaceLenBytes: 30,
			MaxIdentifierKeyLenBytes:       30,
		},
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
	p := &vault.CreateSecretsRequest{
		EncryptedSecrets: []*vault.EncryptedSecret{
			{
				Id:             id,
				EncryptedValue: "foo",
			},
			{
				Id:             id,
				EncryptedValue: "bla",
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

	assert.Equal(t, vault.RequestType_CREATE_SECRETS, o.RequestType)
	assert.True(t, proto.Equal(o.GetCreateSecretsRequest(), p))

	batchResp := o.GetCreateSecretsResponse()
	assert.Len(t, p.EncryptedSecrets, 2)
	assert.Len(t, p.EncryptedSecrets, len(batchResp.Responses))

	assert.True(t, proto.Equal(p.EncryptedSecrets[0].Id, batchResp.Responses[0].Id))
	resp := batchResp.Responses[0]
	assert.Contains(t, resp.GetError(), "duplicate create request for secret identifier")

	assert.True(t, proto.Equal(p.EncryptedSecrets[1].Id, batchResp.Responses[1].Id))
	resp = batchResp.Responses[1]
	assert.Contains(t, resp.GetError(), "duplicate create request for secret identifier")
}

func TestPlugin_StateTransition_CreateSecretsRequest_CorrectlyTracksLimits(t *testing.T) {
	lggr := logger.TestLogger(t)
	store := requests.NewStore[*Request]()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := &ReportingPlugin{
		lggr:  lggr,
		store: store,
		cfg: &ReportingPluginConfig{
			BatchSize:                      10,
			PublicKey:                      pk,
			PrivateKeyShare:                shares[0],
			MaxSecretsPerOwner:             1,
			MaxCiphertextLenBytes:          1024,
			MaxIdentifierOwnerLenBytes:     30,
			MaxIdentifierNamespaceLenBytes: 30,
			MaxIdentifierKeyLenBytes:       30,
		},
	}

	seqNr := uint64(1)
	rdr := &kv{
		m: make(map[string]response),
	}

	ct, err := tdh2easy.Encrypt(pk, []byte("my secret value"))
	require.NoError(t, err)

	ciphertextBytes, err := ct.Marshal()
	require.NoError(t, err)

	id1 := &vault.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "my_secret",
	}
	req1 := &vault.CreateSecretsRequest{
		RequestId: "req1",
		EncryptedSecrets: []*vault.EncryptedSecret{
			{
				Id:             id1,
				EncryptedValue: base64.StdEncoding.EncodeToString(ciphertextBytes),
			},
		},
	}
	resp1 := &vault.CreateSecretsResponse{
		Responses: []*vault.CreateSecretResponse{
			{
				Id:      id1,
				Success: false,
			},
		},
	}

	id2 := &vault.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "my_secret2",
	}
	req2 := &vault.CreateSecretsRequest{
		RequestId: "req2",
		EncryptedSecrets: []*vault.EncryptedSecret{
			{
				Id:             id2,
				EncryptedValue: base64.StdEncoding.EncodeToString(ciphertextBytes),
			},
		},
	}
	resp2 := &vault.CreateSecretsResponse{
		Responses: []*vault.CreateSecretResponse{
			{
				Id:      id2,
				Success: false,
			},
		},
	}

	obs := marshalObservations(t, observation{id1, req1, resp1}, observation{id2, req2, resp2})

	reportPrecursor, err := r.StateTransition(
		t.Context(),
		seqNr,
		types.AttributedQuery{},
		[]types.AttributedObservation{
			{Observation: obs},
			{Observation: obs},
			{Observation: obs},
		},
		rdr,
		nil,
	)
	require.NoError(t, err)

	os := &vault.Outcomes{}
	err = proto.Unmarshal(reportPrecursor, os)
	require.NoError(t, err)

	assert.Len(t, os.Outcomes, 2)

	o1 := os.Outcomes[0]
	assert.Equal(t, vault.RequestType_CREATE_SECRETS, o1.RequestType)
	assert.Len(t, o1.GetCreateSecretsResponse().Responses, 1)
	r1 := o1.GetCreateSecretsResponse().Responses[0]
	assert.True(t, r1.Success)

	o2 := os.Outcomes[1]
	assert.Equal(t, vault.RequestType_CREATE_SECRETS, o2.RequestType)
	assert.Len(t, o2.GetCreateSecretsResponse().Responses, 1)
	r2 := o2.GetCreateSecretsResponse().Responses[0]
	assert.False(t, r2.Success)
	assert.Contains(t, r2.GetError(), "owner has reached maximum number of secrets")
}

func TestPlugin_Observation_CreateSecretsRequest_InvalidCiphertext(t *testing.T) {
	lggr := logger.TestLogger(t)
	store := requests.NewStore[*Request]()
	r := &ReportingPlugin{
		lggr:  lggr,
		store: store,
		cfg: &ReportingPluginConfig{
			BatchSize:                      10,
			PublicKey:                      nil,
			PrivateKeyShare:                nil,
			MaxSecretsPerOwner:             1,
			MaxCiphertextLenBytes:          1024,
			MaxIdentifierOwnerLenBytes:     100,
			MaxIdentifierNamespaceLenBytes: 100,
			MaxIdentifierKeyLenBytes:       100,
		},
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

	assert.Equal(t, vault.RequestType_CREATE_SECRETS, o.RequestType)
	assert.True(t, proto.Equal(o.GetCreateSecretsRequest(), p))

	batchResp := o.GetCreateSecretsResponse()
	assert.Len(t, p.EncryptedSecrets, 1)
	assert.Len(t, p.EncryptedSecrets, len(batchResp.Responses))

	assert.True(t, proto.Equal(p.EncryptedSecrets[0].Id, batchResp.Responses[0].Id))
	resp := batchResp.Responses[0]
	assert.Contains(t, resp.GetError(), "invalid base64 encoding for ciphertext")
}

func TestPlugin_Observation_CreateSecretsRequest_InvalidCiphertext_TooLong(t *testing.T) {
	lggr := logger.TestLogger(t)
	store := requests.NewStore[*Request]()
	r := &ReportingPlugin{
		lggr:  lggr,
		store: store,
		cfg: &ReportingPluginConfig{
			BatchSize:                      10,
			PublicKey:                      nil,
			PrivateKeyShare:                nil,
			MaxSecretsPerOwner:             1,
			MaxCiphertextLenBytes:          10,
			MaxIdentifierOwnerLenBytes:     100,
			MaxIdentifierNamespaceLenBytes: 100,
			MaxIdentifierKeyLenBytes:       100,
		},
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

	assert.Equal(t, vault.RequestType_CREATE_SECRETS, o.RequestType)
	assert.True(t, proto.Equal(o.GetCreateSecretsRequest(), p))

	batchResp := o.GetCreateSecretsResponse()
	assert.Len(t, p.EncryptedSecrets, 1)
	assert.Len(t, p.EncryptedSecrets, len(batchResp.Responses))

	assert.True(t, proto.Equal(p.EncryptedSecrets[0].Id, batchResp.Responses[0].Id))
	resp := batchResp.Responses[0]
	assert.Contains(t, resp.GetError(), "ciphertext size exceeds maximum allowed size: 10 bytes")
}

func TestPlugin_Observation_CreateSecretsRequest_InvalidCiphertext_EncryptedWithWrongPublicKey(t *testing.T) {
	lggr := logger.TestLogger(t)
	store := requests.NewStore[*Request]()
	// Wrong key
	_, wrongPublicKey, _, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	// Right key
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := &ReportingPlugin{
		lggr:  lggr,
		store: store,
		cfg: &ReportingPluginConfig{
			BatchSize:                      10,
			PublicKey:                      pk,
			PrivateKeyShare:                shares[0],
			MaxSecretsPerOwner:             1,
			MaxCiphertextLenBytes:          1024,
			MaxIdentifierOwnerLenBytes:     100,
			MaxIdentifierNamespaceLenBytes: 100,
			MaxIdentifierKeyLenBytes:       100,
		},
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

	assert.Equal(t, vault.RequestType_CREATE_SECRETS, o.RequestType)
	assert.True(t, proto.Equal(o.GetCreateSecretsRequest(), p))

	batchResp := o.GetCreateSecretsResponse()
	assert.Len(t, p.EncryptedSecrets, 1)
	assert.Len(t, p.EncryptedSecrets, len(batchResp.Responses))

	assert.True(t, proto.Equal(p.EncryptedSecrets[0].Id, batchResp.Responses[0].Id))
	resp := batchResp.Responses[0]
	assert.Contains(t, resp.GetError(), "failed to verify ciphertext")
}

func TestPlugin_StateTransition_CreateSecretsRequest_TooManySecretsForOwner(t *testing.T) {
	lggr := logger.TestLogger(t)
	store := requests.NewStore[*Request]()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := &ReportingPlugin{
		lggr:  lggr,
		store: store,
		cfg: &ReportingPluginConfig{
			BatchSize:                      10,
			PublicKey:                      pk,
			PrivateKeyShare:                shares[0],
			MaxSecretsPerOwner:             1,
			MaxCiphertextLenBytes:          1024,
			MaxIdentifierOwnerLenBytes:     100,
			MaxIdentifierNamespaceLenBytes: 100,
			MaxIdentifierKeyLenBytes:       100,
		},
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
	kvstore := NewWriteStore(rdr)
	err = kvstore.WriteMetadata(id.Owner, &vault.StoredMetadata{
		SecretIdentifiers: []*vault.SecretIdentifier{
			{
				Owner:     "owner",
				Namespace: "main",
				Key:       "secret2",
			},
		},
	})
	require.NoError(t, err)

	ct, err := tdh2easy.Encrypt(pk, []byte("my secret value"))
	require.NoError(t, err)

	ciphertextBytes, err := ct.Marshal()
	require.NoError(t, err)

	req := &vault.CreateSecretsRequest{
		EncryptedSecrets: []*vault.EncryptedSecret{
			{
				Id:             id,
				EncryptedValue: base64.StdEncoding.EncodeToString(ciphertextBytes),
			},
		},
	}
	resp := &vault.CreateSecretsResponse{
		Responses: []*vault.CreateSecretResponse{
			{
				Id:      id,
				Success: false,
			},
		},
	}
	data, err := r.StateTransition(
		t.Context(),
		seqNr,
		types.AttributedQuery{},
		[]types.AttributedObservation{
			{
				Observation: marshalObservations(t, observation{id, req, resp}),
			},
		},
		rdr,
		nil,
	)
	require.NoError(t, err)

	os := &vault.Outcomes{}
	err = proto.Unmarshal(data, os)
	require.NoError(t, err)

	assert.Len(t, os.Outcomes, 1)
	o := os.Outcomes[0]

	assert.Len(t, o.GetCreateSecretsResponse().Responses, 1)
	assert.Contains(t, o.GetCreateSecretsResponse().Responses[0].Error, "owner has reached maximum number of secrets")
}

func TestPlugin_StateTransition_CreateSecretsRequest_SecretExistsForKey(t *testing.T) {
	lggr := logger.TestLogger(t)
	store := requests.NewStore[*Request]()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := &ReportingPlugin{
		lggr:  lggr,
		store: store,
		cfg: &ReportingPluginConfig{
			BatchSize:                      10,
			PublicKey:                      pk,
			PrivateKeyShare:                shares[0],
			MaxSecretsPerOwner:             1,
			MaxCiphertextLenBytes:          1024,
			MaxIdentifierOwnerLenBytes:     100,
			MaxIdentifierNamespaceLenBytes: 100,
			MaxIdentifierKeyLenBytes:       100,
		},
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
	kvstore := NewWriteStore(rdr)
	err = kvstore.WriteSecret(id, &vault.StoredSecret{
		EncryptedSecret: []byte("some-ciphertext"),
	})
	require.NoError(t, err)

	ct, err := tdh2easy.Encrypt(pk, []byte("my secret value"))
	require.NoError(t, err)

	ciphertextBytes, err := ct.Marshal()
	require.NoError(t, err)

	req := &vault.CreateSecretsRequest{
		EncryptedSecrets: []*vault.EncryptedSecret{
			{
				Id:             id,
				EncryptedValue: base64.StdEncoding.EncodeToString(ciphertextBytes),
			},
		},
	}
	resp := &vault.CreateSecretsResponse{
		Responses: []*vault.CreateSecretResponse{
			{
				Id:      id,
				Success: false,
			},
		},
	}
	data, err := r.StateTransition(
		t.Context(),
		seqNr,
		types.AttributedQuery{},
		[]types.AttributedObservation{
			{
				Observation: marshalObservations(t, observation{id, req, resp}),
			},
		},
		rdr,
		nil,
	)
	require.NoError(t, err)

	os := &vault.Outcomes{}
	err = proto.Unmarshal(data, os)
	require.NoError(t, err)

	assert.Len(t, os.Outcomes, 1)
	o := os.Outcomes[0]

	assert.Len(t, o.GetCreateSecretsResponse().Responses, 1)
	assert.Contains(t, o.GetCreateSecretsResponse().Responses[0].Error, "key already exists")
}

func TestPlugin_Observation_CreateSecretsRequest_Success(t *testing.T) {
	lggr := logger.TestLogger(t)
	store := requests.NewStore[*Request]()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := &ReportingPlugin{
		lggr:  lggr,
		store: store,
		cfg: &ReportingPluginConfig{
			BatchSize:                      10,
			PublicKey:                      pk,
			PrivateKeyShare:                shares[0],
			MaxSecretsPerOwner:             1,
			MaxCiphertextLenBytes:          1024,
			MaxIdentifierOwnerLenBytes:     100,
			MaxIdentifierNamespaceLenBytes: 100,
			MaxIdentifierKeyLenBytes:       100,
		},
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

	assert.Equal(t, vault.RequestType_CREATE_SECRETS, o.RequestType)
	assert.True(t, proto.Equal(o.GetCreateSecretsRequest(), p))

	batchResp := o.GetCreateSecretsResponse()
	assert.Len(t, p.EncryptedSecrets, 1)
	assert.Len(t, p.EncryptedSecrets, len(batchResp.Responses))

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
		switch tr := ob.req.(type) {
		case *vault.GetSecretsRequest:
			o.RequestType = vault.RequestType_GET_SECRETS
			o.Request = &vault.Observation_GetSecretsRequest{
				GetSecretsRequest: tr,
			}
		case *vault.CreateSecretsRequest:
			o.RequestType = vault.RequestType_CREATE_SECRETS
			o.Request = &vault.Observation_CreateSecretsRequest{
				CreateSecretsRequest: tr,
			}
		}

		switch tr := ob.resp.(type) {
		case *vault.GetSecretsResponse:
			o.Response = &vault.Observation_GetSecretsResponse{
				GetSecretsResponse: tr,
			}
		case *vault.CreateSecretsResponse:
			o.Response = &vault.Observation_CreateSecretsResponse{
				CreateSecretsResponse: tr,
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
		onchainCfg: ocr3types.ReportingPluginConfig{
			N: 4,
			F: 1,
		},
		store: store,
		cfg: &ReportingPluginConfig{
			BatchSize:                      10,
			PublicKey:                      pk,
			PrivateKeyShare:                shares[0],
			MaxSecretsPerOwner:             1,
			MaxCiphertextLenBytes:          1024,
			MaxIdentifierOwnerLenBytes:     100,
			MaxIdentifierNamespaceLenBytes: 100,
			MaxIdentifierKeyLenBytes:       100,
		},
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

	assert.Empty(t, os.Outcomes, 0)

	assert.Equal(t, 1, observed.FilterMessage("insufficient observations found for id").Len())
}

func TestPlugin_ValidateObservations_InvalidObservations(t *testing.T) {
	lggr, _ := logger.TestLoggerObserved(t, zapcore.DebugLevel)
	store := requests.NewStore[*Request]()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := &ReportingPlugin{
		lggr: lggr,
		onchainCfg: ocr3types.ReportingPluginConfig{
			N: 4,
			F: 1,
		},
		store: store,
		cfg: &ReportingPluginConfig{
			BatchSize:                      10,
			PublicKey:                      pk,
			PrivateKeyShare:                shares[0],
			MaxSecretsPerOwner:             1,
			MaxCiphertextLenBytes:          1024,
			MaxIdentifierOwnerLenBytes:     100,
			MaxIdentifierNamespaceLenBytes: 100,
			MaxIdentifierKeyLenBytes:       100,
		},
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
	err = r.ValidateObservation(
		t.Context(),
		seqNr,
		types.AttributedQuery{},
		types.AttributedObservation{Observation: types.Observation(obsb)},
		kv,
		nil,
	)
	require.ErrorContains(t, err, "GetSecrets observation must have both request and response")

	// Invalid observation -- data can't be unmarshaled
	err = r.ValidateObservation(
		t.Context(),
		seqNr,
		types.AttributedQuery{},
		types.AttributedObservation{Observation: types.Observation([]byte("hello world"))},
		kv,
		nil,
	)

	require.ErrorContains(t, err, "failed to unmarshal observations")

	// Invalid observation -- a single observation set has observations for multiple request ids
	correctResp := &vault.GetSecretsResponse{
		Responses: []*vault.SecretResponse{
			{
				Id: id1,
			},
		},
	}
	obsb = marshalObservations(t, observation{id1, req, correctResp}, observation{id1, req, correctResp})
	err = r.ValidateObservation(
		t.Context(),
		seqNr,
		types.AttributedQuery{},
		types.AttributedObservation{Observation: types.Observation(obsb)},
		kv,
		nil,
	)
	assert.ErrorContains(t, err, "invalid observation: a single observation cannot contain duplicate observations for the same request id")
}

func TestPlugin_StateTransition_ShasDontMatch(t *testing.T) {
	lggr, observed := logger.TestLoggerObserved(t, zapcore.DebugLevel)
	store := requests.NewStore[*Request]()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := &ReportingPlugin{
		lggr: lggr,
		onchainCfg: ocr3types.ReportingPluginConfig{
			N: 4,
			F: 1,
		},
		store: store,
		cfg: &ReportingPluginConfig{
			BatchSize:                      10,
			PublicKey:                      pk,
			PrivateKeyShare:                shares[0],
			MaxSecretsPerOwner:             1,
			MaxCiphertextLenBytes:          1024,
			MaxIdentifierOwnerLenBytes:     100,
			MaxIdentifierNamespaceLenBytes: 100,
			MaxIdentifierKeyLenBytes:       100,
		},
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

	assert.Empty(t, os.Outcomes)

	assert.Equal(t, 1, observed.FilterMessage("insufficient observations found for id").Len())
}

func TestPlugin_StateTransition_AggregatesValidationErrors(t *testing.T) {
	lggr, observed := logger.TestLoggerObserved(t, zapcore.DebugLevel)
	store := requests.NewStore[*Request]()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := &ReportingPlugin{
		lggr: lggr,
		onchainCfg: ocr3types.ReportingPluginConfig{
			N: 4,
			F: 1,
		},
		store: store,
		cfg: &ReportingPluginConfig{
			BatchSize:                      10,
			PublicKey:                      pk,
			PrivateKeyShare:                shares[0],
			MaxSecretsPerOwner:             1,
			MaxCiphertextLenBytes:          1024,
			MaxIdentifierOwnerLenBytes:     100,
			MaxIdentifierNamespaceLenBytes: 100,
			MaxIdentifierKeyLenBytes:       100,
		},
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

	obsb := marshalObservations(t, observation{id, req, resp})
	reportPrecursor, err := r.StateTransition(
		t.Context(),
		seqNr,
		types.AttributedQuery{},
		[]types.AttributedObservation{
			{Observation: types.Observation(obsb)},
			{Observation: types.Observation(obsb)},
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
		onchainCfg: ocr3types.ReportingPluginConfig{
			N: 4,
			F: 1,
		},
		store: store,
		cfg: &ReportingPluginConfig{
			BatchSize:                      10,
			PublicKey:                      pk,
			PrivateKeyShare:                shares[0],
			MaxSecretsPerOwner:             1,
			MaxCiphertextLenBytes:          1024,
			MaxIdentifierOwnerLenBytes:     100,
			MaxIdentifierNamespaceLenBytes: 100,
			MaxIdentifierKeyLenBytes:       100,
		},
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

	obsb1 := marshalObservations(t, observation{id, req, resp1})
	obsb2 := marshalObservations(t, observation{id, req, resp2})
	obsb3 := marshalObservations(t, observation{id, req, resp3})
	reportPrecursor, err := r.StateTransition(
		t.Context(),
		seqNr,
		types.AttributedQuery{},
		[]types.AttributedObservation{
			{Observation: types.Observation(obsb1)},
			{Observation: types.Observation(obsb2)},
			{Observation: types.Observation(obsb3)},
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
		onchainCfg: ocr3types.ReportingPluginConfig{
			N: 4,
			F: 1,
		},
		store: store,
		cfg: &ReportingPluginConfig{
			BatchSize:                      10,
			PublicKey:                      pk,
			PrivateKeyShare:                shares[0],
			MaxSecretsPerOwner:             1,
			MaxCiphertextLenBytes:          1024,
			MaxIdentifierOwnerLenBytes:     100,
			MaxIdentifierNamespaceLenBytes: 100,
			MaxIdentifierKeyLenBytes:       100,
		},
	}

	seqNr := uint64(1)
	kv := &kv{
		m: make(map[string]response),
	}
	rs := NewReadStore(kv)

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

	obsb := marshalObservations(t, observation{id, req, resp})
	reportPrecursor, err := r.StateTransition(
		t.Context(),
		seqNr,
		types.AttributedQuery{},
		[]types.AttributedObservation{
			{Observation: types.Observation(obsb)},
			{Observation: types.Observation(obsb)},
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

	ss, err := rs.GetSecret(id)
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
				Success: true,
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

	id2 := &vault.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "secret2",
	}
	req2 := &vault.GetSecretsRequest{
		Requests: []*vault.SecretRequest{
			{
				Id: id2,
			},
		},
	}
	resp2 := &vault.GetSecretsResponse{
		Responses: []*vault.SecretResponse{
			{
				Id:     id2,
				Result: &vault.SecretResponse_Data{Data: &vault.SecretData{EncryptedValue: value}},
			},
		},
	}
	expectedOutcome2 := &vault.Outcome{
		Id:          keyFor(id2),
		RequestType: vault.RequestType_GET_SECRETS,
		Request: &vault.Outcome_GetSecretsRequest{
			GetSecretsRequest: req2,
		},
		Response: &vault.Outcome_GetSecretsResponse{
			GetSecretsResponse: resp2,
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
		onchainCfg: ocr3types.ReportingPluginConfig{
			N: 4,
			F: 1,
		},
		store: store,
		cfg: &ReportingPluginConfig{
			BatchSize:                      10,
			PublicKey:                      pk,
			PrivateKeyShare:                shares[0],
			MaxSecretsPerOwner:             1,
			MaxCiphertextLenBytes:          1024,
			MaxIdentifierOwnerLenBytes:     100,
			MaxIdentifierNamespaceLenBytes: 100,
			MaxIdentifierKeyLenBytes:       100,
		},
	}

	rs, err := r.Reports(t.Context(), uint64(1), osb)
	require.NoError(t, err)

	assert.Len(t, rs, 2)

	o1 := rs[0]
	info1 := &vault.ReportInfo{}
	err = proto.Unmarshal(o1.ReportWithInfo.Info, info1)
	require.NoError(t, err)
	assert.True(t, proto.Equal(&vault.ReportInfo{
		Id:          keyFor(id),
		Format:      vault.ReportFormat_REPORT_FORMAT_JSON,
		RequestType: vault.RequestType_CREATE_SECRETS,
	}, info1))

	expectedBytes, err := ToCanonicalJSON(resp)
	require.NoError(t, err)
	assert.Equal(t, expectedBytes, []byte(o1.ReportWithInfo.Report))

	o2 := rs[1]
	info2 := &vault.ReportInfo{}
	err = proto.Unmarshal(o2.ReportWithInfo.Info, info2)
	require.NoError(t, err)
	assert.True(t, proto.Equal(&vault.ReportInfo{
		Id:          keyFor(id2),
		Format:      vault.ReportFormat_REPORT_FORMAT_PROTOBUF,
		RequestType: vault.RequestType_GET_SECRETS,
	}, info2))

	o2r := &vault.GetSecretsResponse{}
	err = proto.Unmarshal(o2.ReportWithInfo.Report, o2r)
	require.NoError(t, err)
	assert.True(t, proto.Equal(resp2, o2r))
}
