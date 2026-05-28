package vault

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/tdh2/go/tdh2/tdh2easy"
	"go.uber.org/zap/zapcore"
	"golang.org/x/crypto/nacl/box"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaultutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestPlugin_Observation_GetSecretsRequest_OmitsRequestWhenOptimizationsEnabled(t *testing.T) {
	r := newTestReportingPlugin(t, withVaultOptimizationsEnabled())

	id := &vaultcommon.SecretIdentifier{Owner: "owner", Namespace: "main", Key: "secret"}
	p := &vaultcommon.GetSecretsRequest{
		Requests: []*vaultcommon.SecretRequest{{Id: id}},
	}
	anyp, err := anypb.New(p)
	require.NoError(t, err)

	rdr := &kv{m: make(map[string]response)}
	require.NoError(t, newTestWriteStore(t, rdr).WritePendingQueue(t.Context(),
		[]*vaultcommon.StoredPendingQueueItem{{Id: "request-1", Item: anyp}},
	))

	data, err := r.Observation(t.Context(), 1, types.AttributedQuery{}, rdr, &blobber{})
	require.NoError(t, err)

	obs := &vaultcommon.Observations{}
	require.NoError(t, proto.Unmarshal(data, obs))
	require.Len(t, obs.Observations, 1)
	assert.Nil(t, obs.Observations[0].GetGetSecretsRequest())
}
func TestPlugin_Observation_GetSecretsRequest_SecretIdentifierInvalid(t *testing.T) {
	tcs := []struct {
		name            string
		id              *vaultcommon.SecretIdentifier
		maxIDLen        int
		maxOwnerLen     int
		maxNamespaceLen int
		maxKeyLen       int
		err             string
	}{
		{
			name: "nil id",
			id:   nil,
			err:  "secret identifier cannot be nil",
		},
		{
			name: "empty id",
			id:   &vaultcommon.SecretIdentifier{},
			err:  "key cannot be empty",
		},
		{
			name: "empty id",
			id: &vaultcommon.SecretIdentifier{
				Key:       "hello",
				Namespace: "world",
			},
			err: "owner cannot be empty",
		},
		{
			name:     "id is too long",
			maxIDLen: 10,
			id: &vaultcommon.SecretIdentifier{
				Owner:     "owner",
				Key:       "hello",
				Namespace: "world",
			},
			err: "owner exceeds maximum length of 3b",
		},
		{
			name:            "namespace exceeds maximum length",
			maxNamespaceLen: 3,
			id: &vaultcommon.SecretIdentifier{
				Owner:     "owner",
				Key:       "hello",
				Namespace: "world",
			},
			err: "namespace exceeds maximum length of 3b",
		},
		{
			name:      "key exceeds maximum length",
			maxKeyLen: 3,
			id: &vaultcommon.SecretIdentifier{
				Owner:     "owner",
				Key:       "hello",
				Namespace: "world",
			},
			err: "key exceeds maximum length of 3b",
		},
	}

	for _, tc := range tcs {
		ownerLen, namespaceLen, keyLen := 256, 256, 256
		if tc.maxIDLen > 0 {
			ownerLen = tc.maxIDLen / 3
			namespaceLen = tc.maxIDLen / 3
			keyLen = tc.maxIDLen / 3
		}
		if tc.maxOwnerLen > 0 {
			ownerLen = tc.maxOwnerLen
		}
		if tc.maxNamespaceLen > 0 {
			namespaceLen = tc.maxNamespaceLen
		}
		if tc.maxKeyLen > 0 {
			keyLen = tc.maxKeyLen
		}
		r := newTestReportingPlugin(t, withMaxIdentifierLengths(ownerLen, namespaceLen, keyLen))

		seqNr := uint64(1)
		rdr := &kv{
			m: make(map[string]response),
		}
		p := &vaultcommon.GetSecretsRequest{
			Requests: []*vaultcommon.SecretRequest{
				{
					Id:             tc.id,
					EncryptionKeys: []string{"foo"},
				},
			},
		}
		anyp, err := anypb.New(p)
		require.NoError(t, err)
		err = newTestWriteStore(t, rdr).WritePendingQueue(t.Context(),
			[]*vaultcommon.StoredPendingQueueItem{
				{Id: "request-1", Item: anyp},
			},
		)
		require.NoError(t, err)
		bf := &blobber{}
		data, err := r.Observation(t.Context(), seqNr, types.AttributedQuery{}, rdr, bf)
		require.NoError(t, err)

		obs := &vaultcommon.Observations{}
		err = proto.Unmarshal(data, obs)
		require.NoError(t, err)

		assert.Len(t, obs.Observations, 1)
		o := obs.Observations[0]

		assert.Equal(t, vaultcommon.RequestType_GET_SECRETS, o.RequestType)
		assert.True(t, proto.Equal(p, o.GetGetSecretsRequest()))

		batchResp := o.GetGetSecretsResponse()
		assert.Len(t, p.Requests, 1)
		assert.Len(t, p.Requests, len(batchResp.Responses))

		assert.True(t, proto.Equal(p.Requests[0].Id, batchResp.Responses[0].Id))
		resp := batchResp.Responses[0]
		assert.Contains(t, resp.GetError(), tc.err)
	}
}
func TestPlugin_Observation_GetSecretsRequest_ResponseUsesCanonicalIdentifier(t *testing.T) {
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(t, withKeys(pk, shares[0]))

	id := &vaultcommon.SecretIdentifier{
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

	createdID := &vaultcommon.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "my_secret",
	}
	err = newTestWriteStore(t, rdr).WriteSecret(t.Context(), createdID, &vaultcommon.StoredSecret{
		EncryptedSecret: ciphertextBytes,
	})
	require.NoError(t, err)

	pubK, _, err := box.GenerateKey(rand.Reader)
	require.NoError(t, err)

	pks := hex.EncodeToString(pubK[:])

	p := &vaultcommon.GetSecretsRequest{
		Requests: []*vaultcommon.SecretRequest{
			{
				Id:             id,
				EncryptionKeys: []string{pks},
			},
		},
	}
	anyp, err := anypb.New(p)
	require.NoError(t, err)
	err = newTestWriteStore(t, rdr).WritePendingQueue(t.Context(),
		[]*vaultcommon.StoredPendingQueueItem{
			{Id: "request-1", Item: anyp},
		},
	)
	require.NoError(t, err)
	seqNr := uint64(1)
	bf := &blobber{}
	data, err := r.Observation(t.Context(), seqNr, types.AttributedQuery{}, rdr, bf)
	require.NoError(t, err)

	obs := &vaultcommon.Observations{}
	err = proto.Unmarshal(data, obs)
	require.NoError(t, err)

	assert.Len(t, obs.Observations, 1)
	o := obs.Observations[0]

	assert.Equal(t, vaultcommon.RequestType_GET_SECRETS, o.RequestType)
	assert.True(t, proto.Equal(p, o.GetGetSecretsRequest()))

	batchResp := o.GetGetSecretsResponse()
	assert.Len(t, p.Requests, 1)
	assert.Len(t, p.Requests, len(batchResp.Responses))

	assert.True(t, proto.Equal(batchResp.Responses[0].Id, createdID))
}
func TestPlugin_Observation_GetSecretsRequest_WorkflowOwnerLabelAccepted(t *testing.T) {
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)

	r := newTestReportingPlugin(t, withKeys(pk, shares[0]))

	ownerAddr := common.HexToAddress("0x1111111111111111111111111111111111111111")
	owner := ownerAddr.Hex()
	id := &vaultcommon.SecretIdentifier{
		Owner:     owner,
		Namespace: "main",
		Key:       "my_secret",
	}
	rdr := &kv{m: make(map[string]response)}

	encrypted, err := vaultutils.EncryptSecretWithWorkflowOwner("my-secret-value", pk, ownerAddr)
	require.NoError(t, err)
	ciphertextBytes, err := hex.DecodeString(encrypted)
	require.NoError(t, err)

	err = newTestWriteStore(t, rdr).WriteSecret(t.Context(), id, &vaultcommon.StoredSecret{
		EncryptedSecret: ciphertextBytes,
	})
	require.NoError(t, err)

	pubK, _, err := box.GenerateKey(rand.Reader)
	require.NoError(t, err)

	p := &vaultcommon.GetSecretsRequest{
		Requests: []*vaultcommon.SecretRequest{
			{
				Id:             id,
				EncryptionKeys: []string{hex.EncodeToString(pubK[:])},
			},
		},
	}
	anyp, err := anypb.New(p)
	require.NoError(t, err)
	err = newTestWriteStore(t, rdr).WritePendingQueue(t.Context(),
		[]*vaultcommon.StoredPendingQueueItem{
			{Id: "request-1", Item: anyp},
		},
	)
	require.NoError(t, err)

	data, err := r.Observation(t.Context(), 1, types.AttributedQuery{}, rdr, &blobber{})
	require.NoError(t, err)

	obs := &vaultcommon.Observations{}
	err = proto.Unmarshal(data, obs)
	require.NoError(t, err)

	require.Len(t, obs.Observations, 1)
	batchResp := obs.Observations[0].GetGetSecretsResponse()
	require.Len(t, batchResp.Responses, 1)
	require.NotNil(t, batchResp.Responses[0].GetId())
	assert.Equal(t, owner, batchResp.Responses[0].GetId().GetOwner())
	assert.Empty(t, batchResp.Responses[0].GetError())
}
func TestPlugin_Observation_GetSecretsRequest_WrongOwnerLabelRejected(t *testing.T) {
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)

	r := newTestReportingPlugin(t, withKeys(pk, shares[0]))

	goodOwner := common.HexToAddress("0x1111111111111111111111111111111111111111")
	wrongOwner := common.HexToAddress("0x2222222222222222222222222222222222222222")

	id := &vaultcommon.SecretIdentifier{
		Owner:     wrongOwner.Hex(),
		Namespace: "main",
		Key:       "my_secret",
	}
	rdr := &kv{m: make(map[string]response)}

	encrypted, err := vaultutils.EncryptSecretWithWorkflowOwner("my-secret-value", pk, goodOwner)
	require.NoError(t, err)
	ciphertextBytes, err := hex.DecodeString(encrypted)
	require.NoError(t, err)

	err = newTestWriteStore(t, rdr).WriteSecret(t.Context(), id, &vaultcommon.StoredSecret{
		EncryptedSecret: ciphertextBytes,
	})
	require.NoError(t, err)

	pubK, _, err := box.GenerateKey(rand.Reader)
	require.NoError(t, err)

	p := &vaultcommon.GetSecretsRequest{
		Requests: []*vaultcommon.SecretRequest{
			{
				Id:             id,
				EncryptionKeys: []string{hex.EncodeToString(pubK[:])},
			},
		},
	}
	anyp, err := anypb.New(p)
	require.NoError(t, err)
	err = newTestWriteStore(t, rdr).WritePendingQueue(t.Context(),
		[]*vaultcommon.StoredPendingQueueItem{
			{Id: "request-1", Item: anyp},
		},
	)
	require.NoError(t, err)

	data, err := r.Observation(t.Context(), 1, types.AttributedQuery{}, rdr, &blobber{})
	require.NoError(t, err)

	obs := &vaultcommon.Observations{}
	err = proto.Unmarshal(data, obs)
	require.NoError(t, err)

	require.Len(t, obs.Observations, 1)
	batchResp := obs.Observations[0].GetGetSecretsResponse()
	require.Len(t, batchResp.Responses, 1)
	assert.Contains(t, batchResp.Responses[0].GetError(), "failed to handle get secret request")
}
func TestPlugin_Observation_GetSecretsRequest_SecretDoesNotExist(t *testing.T) {
	r := newTestReportingPlugin(t)

	seqNr := uint64(1)
	rdr := &kv{
		m: make(map[string]response),
	}
	id := &vaultcommon.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "my_secret",
	}
	p := &vaultcommon.GetSecretsRequest{
		Requests: []*vaultcommon.SecretRequest{
			{
				Id:             id,
				EncryptionKeys: []string{"foo"},
			},
		},
	}
	anyp, err := anypb.New(p)
	require.NoError(t, err)
	err = newTestWriteStore(t, rdr).WritePendingQueue(t.Context(),
		[]*vaultcommon.StoredPendingQueueItem{
			{Id: "request-1", Item: anyp},
		},
	)
	require.NoError(t, err)
	bf := &blobber{}
	data, err := r.Observation(t.Context(), seqNr, types.AttributedQuery{}, rdr, bf)
	require.NoError(t, err)

	obs := &vaultcommon.Observations{}
	err = proto.Unmarshal(data, obs)
	require.NoError(t, err)

	assert.Len(t, obs.Observations, 1)
	o := obs.Observations[0]

	assert.Equal(t, vaultcommon.RequestType_GET_SECRETS, o.RequestType)
	assert.True(t, proto.Equal(p, o.GetGetSecretsRequest()))

	batchResp := o.GetGetSecretsResponse()
	assert.Len(t, p.Requests, 1)
	assert.Len(t, p.Requests, len(batchResp.Responses))

	assert.True(t, proto.Equal(p.Requests[0].Id, batchResp.Responses[0].Id))
	resp := batchResp.Responses[0]
	assert.Contains(t, resp.GetError(), "key does not exist")
}
func TestPlugin_Observation_GetSecretsRequest_SecretExistsButIsIncorrect(t *testing.T) {
	lggr, observed := logger.TestLoggerObserved(t, zapcore.DebugLevel)
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)

	r := newTestReportingPlugin(t, withLggr(lggr), withKeys(pk, shares[0]))

	id := &vaultcommon.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "my_secret",
	}
	rdr := &kv{
		m: make(map[string]response),
	}

	err = newTestWriteStore(t, rdr).WriteSecret(t.Context(), id, &vaultcommon.StoredSecret{
		EncryptedSecret: []byte("invalid-ciphertext"),
	})
	require.NoError(t, err)

	p := &vaultcommon.GetSecretsRequest{
		Requests: []*vaultcommon.SecretRequest{
			{
				Id:             id,
				EncryptionKeys: []string{"foo"},
			},
		},
	}
	anyp, err := anypb.New(p)
	require.NoError(t, err)
	err = newTestWriteStore(t, rdr).WritePendingQueue(t.Context(),
		[]*vaultcommon.StoredPendingQueueItem{
			{Id: "request-1", Item: anyp},
		},
	)
	require.NoError(t, err)
	seqNr := uint64(1)
	data, err := r.Observation(t.Context(), seqNr, types.AttributedQuery{}, rdr, &blobber{})
	require.NoError(t, err)

	obs := &vaultcommon.Observations{}
	err = proto.Unmarshal(data, obs)
	require.NoError(t, err)

	assert.Len(t, obs.Observations, 1)
	o := obs.Observations[0]

	assert.Equal(t, vaultcommon.RequestType_GET_SECRETS, o.RequestType)
	assert.True(t, proto.Equal(p, o.GetGetSecretsRequest()))

	batchResp := o.GetGetSecretsResponse()
	assert.Len(t, p.Requests, 1)
	assert.Len(t, p.Requests, len(batchResp.Responses))

	assert.True(t, proto.Equal(p.Requests[0].Id, batchResp.Responses[0].Id))
	resp := batchResp.Responses[0]

	// The error returned is user-friendly
	assert.Contains(t, resp.GetError(), "failed to handle get secret request")

	// Inspect logs to get true source of error
	logs := observed.FilterMessage("failed to observe get secret request item")
	assert.Equal(t, 1, logs.Len())
	fields := logs.All()[0].ContextMap()
	errString := fields["error"]
	assert.Contains(t, errString, "failed to unmarshal ciphertext")
}
func TestPlugin_Observation_GetSecretsRequest_PublicKeyIsInvalid(t *testing.T) {
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(t, withKeys(pk, shares[0]))

	id := &vaultcommon.SecretIdentifier{
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

	err = newTestWriteStore(t, rdr).WriteSecret(t.Context(), id, &vaultcommon.StoredSecret{
		EncryptedSecret: ciphertextBytes,
	})
	require.NoError(t, err)

	p := &vaultcommon.GetSecretsRequest{
		Requests: []*vaultcommon.SecretRequest{
			{
				Id:             id,
				EncryptionKeys: []string{"foo"},
			},
		},
	}
	anyp, err := anypb.New(p)
	require.NoError(t, err)
	err = newTestWriteStore(t, rdr).WritePendingQueue(t.Context(),
		[]*vaultcommon.StoredPendingQueueItem{
			{Id: "request-1", Item: anyp},
		},
	)
	require.NoError(t, err)
	seqNr := uint64(1)
	data, err := r.Observation(t.Context(), seqNr, types.AttributedQuery{}, rdr, &blobber{})
	require.NoError(t, err)

	obs := &vaultcommon.Observations{}
	err = proto.Unmarshal(data, obs)
	require.NoError(t, err)

	assert.Len(t, obs.Observations, 1)
	o := obs.Observations[0]

	assert.Equal(t, vaultcommon.RequestType_GET_SECRETS, o.RequestType)
	assert.True(t, proto.Equal(p, o.GetGetSecretsRequest()))

	batchResp := o.GetGetSecretsResponse()
	assert.Len(t, p.Requests, 1)
	assert.Len(t, p.Requests, len(batchResp.Responses))

	assert.True(t, proto.Equal(p.Requests[0].Id, batchResp.Responses[0].Id))
	resp := batchResp.Responses[0]

	assert.Contains(t, resp.GetError(), "failed to convert public key to bytes")
}
func TestPlugin_Observation_GetSecretsRequest_SecretLabelIsInvalid(t *testing.T) {
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(t, withKeys(pk, shares[0]))

	id := &vaultcommon.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "my_secret",
	}
	rdr := &kv{
		m: make(map[string]response),
	}

	plaintext := []byte("my-secret-value")
	var label [32]byte
	ownerAddress := common.HexToAddress("0x0001020304050607080900010203040506070809")
	copy(label[12:], ownerAddress.Bytes()) // left-pad with 12 zero
	ciphertext, err := tdh2easy.EncryptWithLabel(pk, plaintext, label)
	require.NoError(t, err)
	ciphertextBytes, err := ciphertext.Marshal()
	require.NoError(t, err)

	err = newTestWriteStore(t, rdr).WriteSecret(t.Context(), id, &vaultcommon.StoredSecret{
		EncryptedSecret: ciphertextBytes,
	})
	require.NoError(t, err)

	pubK, _, err := box.GenerateKey(rand.Reader)
	require.NoError(t, err)

	pks := hex.EncodeToString(pubK[:])

	p := &vaultcommon.GetSecretsRequest{
		Requests: []*vaultcommon.SecretRequest{
			{
				Id:             id,
				EncryptionKeys: []string{pks},
			},
		},
	}
	anyp, err := anypb.New(p)
	require.NoError(t, err)
	err = newTestWriteStore(t, rdr).WritePendingQueue(t.Context(),
		[]*vaultcommon.StoredPendingQueueItem{
			{Id: "request-1", Item: anyp},
		},
	)
	require.NoError(t, err)
	seqNr := uint64(1)
	data, err := r.Observation(t.Context(), seqNr, types.AttributedQuery{}, rdr, &blobber{})
	require.NoError(t, err)

	obs := &vaultcommon.Observations{}
	err = proto.Unmarshal(data, obs)
	require.NoError(t, err)

	assert.Len(t, obs.Observations, 1)
	o := obs.Observations[0]

	assert.Equal(t, vaultcommon.RequestType_GET_SECRETS, o.RequestType)
	assert.True(t, proto.Equal(p, o.GetGetSecretsRequest()))

	batchResp := o.GetGetSecretsResponse()
	assert.Len(t, p.Requests, 1)
	assert.Len(t, p.Requests, len(batchResp.Responses))

	assert.True(t, proto.Equal(p.Requests[0].Id, batchResp.Responses[0].Id))
	resp := batchResp.Responses[0]

	assert.Contains(t, resp.GetError(), "failed to handle get secret request")
}
func TestPlugin_Observation_GetSecretsRequest_Success(t *testing.T) {
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(t, withKeys(pk, shares[0]))

	owner := "0x0001020304050607080900010203040506070809"
	id := &vaultcommon.SecretIdentifier{
		Owner:     owner,
		Namespace: "main",
		Key:       "my_secret",
	}
	rdr := &kv{
		m: make(map[string]response),
	}

	plaintext := []byte("my-secret-value")
	var label [32]byte
	ownerAddress := common.HexToAddress(owner)
	copy(label[12:], ownerAddress.Bytes()) // left-pad with 12 zero
	ciphertext, err := tdh2easy.EncryptWithLabel(pk, plaintext, label)
	require.NoError(t, err)
	ciphertextBytes, err := ciphertext.Marshal()
	require.NoError(t, err)

	err = newTestWriteStore(t, rdr).WriteSecret(t.Context(), id, &vaultcommon.StoredSecret{
		EncryptedSecret: ciphertextBytes,
	})
	require.NoError(t, err)

	pubK, privK, err := box.GenerateKey(rand.Reader)
	require.NoError(t, err)

	pks := hex.EncodeToString(pubK[:])

	p := &vaultcommon.GetSecretsRequest{
		Requests: []*vaultcommon.SecretRequest{
			{
				Id:             id,
				EncryptionKeys: []string{pks},
			},
		},
	}
	anyp, err := anypb.New(p)
	require.NoError(t, err)
	err = newTestWriteStore(t, rdr).WritePendingQueue(t.Context(),
		[]*vaultcommon.StoredPendingQueueItem{
			{Id: "request-1", Item: anyp},
		},
	)
	require.NoError(t, err)
	seqNr := uint64(1)
	data, err := r.Observation(t.Context(), seqNr, types.AttributedQuery{}, rdr, &blobber{})
	require.NoError(t, err)

	obs := &vaultcommon.Observations{}
	err = proto.Unmarshal(data, obs)
	require.NoError(t, err)

	assert.Len(t, obs.Observations, 1)
	o := obs.Observations[0]

	assert.Equal(t, vaultcommon.RequestType_GET_SECRETS, o.RequestType)
	assert.True(t, proto.Equal(p, o.GetGetSecretsRequest()))

	batchResp := o.GetGetSecretsResponse()
	assert.Len(t, p.Requests, 1)
	assert.Len(t, p.Requests, len(batchResp.Responses))

	assert.True(t, proto.Equal(p.Requests[0].Id, batchResp.Responses[0].Id))
	resp := batchResp.Responses[0]

	assert.Empty(t, resp.GetError())

	assert.Equal(t, hex.EncodeToString(ciphertextBytes), resp.GetData().EncryptedValue)

	assert.Len(t, resp.GetData().EncryptedDecryptionKeyShares, 1)
	shareString := resp.GetData().EncryptedDecryptionKeyShares[0].Shares[0]

	share, err := hex.DecodeString(shareString)
	require.NoError(t, err)
	msg, ok := box.OpenAnonymous(nil, share, pubK, privK)
	assert.True(t, ok)

	ds := &tdh2easy.DecryptionShare{}
	err = ds.Unmarshal(msg)
	require.NoError(t, err)

	ct := &tdh2easy.Ciphertext{}
	ctb, err := hex.DecodeString(resp.GetData().EncryptedValue)
	require.NoError(t, err)
	err = ct.UnmarshalVerify(ctb, pk)
	require.NoError(t, err)

	gotSecret, err := tdh2easy.Aggregate(ct, []*tdh2easy.DecryptionShare{ds}, 3)
	require.NoError(t, err)

	assert.Equal(t, plaintext, gotSecret)
}
func TestPlugin_Observation_GetSecretsRequest_BinarySharesWhenOptimizationsEnabled(t *testing.T) {
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(t, withKeys(pk, shares[0]), withVaultOptimizationsEnabled())

	owner := "0x0001020304050607080900010203040506070809"
	id := &vaultcommon.SecretIdentifier{
		Owner:     owner,
		Namespace: "main",
		Key:       "my_secret",
	}
	rdr := &kv{
		m: make(map[string]response),
	}

	plaintext := []byte("my-secret-value")
	var label [32]byte
	ownerAddress := common.HexToAddress(owner)
	copy(label[12:], ownerAddress.Bytes())
	ciphertext, err := tdh2easy.EncryptWithLabel(pk, plaintext, label)
	require.NoError(t, err)
	ciphertextBytes, err := ciphertext.Marshal()
	require.NoError(t, err)

	err = newTestWriteStore(t, rdr).WriteSecret(t.Context(), id, &vaultcommon.StoredSecret{
		EncryptedSecret: ciphertextBytes,
	})
	require.NoError(t, err)

	pubK, privK, err := box.GenerateKey(rand.Reader)
	require.NoError(t, err)
	pks := hex.EncodeToString(pubK[:])

	p := &vaultcommon.GetSecretsRequest{
		Requests: []*vaultcommon.SecretRequest{
			{
				Id:             id,
				EncryptionKeys: []string{pks},
			},
		},
		WorkflowOwner: owner,
	}
	anyp, err := anypb.New(p)
	require.NoError(t, err)
	err = newTestWriteStore(t, rdr).WritePendingQueue(t.Context(),
		[]*vaultcommon.StoredPendingQueueItem{
			{Id: "request-1", Item: anyp},
		},
	)
	require.NoError(t, err)

	data, err := r.Observation(t.Context(), 1, types.AttributedQuery{}, rdr, &blobber{})
	require.NoError(t, err)

	obs := &vaultcommon.Observations{}
	require.NoError(t, proto.Unmarshal(data, obs))
	require.Len(t, obs.Observations, 1)
	resp := obs.Observations[0].GetGetSecretsResponse().Responses[0]
	require.Empty(t, resp.GetError())

	encShares := resp.GetData().EncryptedDecryptionKeyShares
	require.Len(t, encShares, 1)
	require.Empty(t, encShares[0].Shares)
	require.Len(t, encShares[0].BinaryShares, 1)

	msg, ok := box.OpenAnonymous(nil, encShares[0].BinaryShares[0], pubK, privK)
	require.True(t, ok)

	ds := &tdh2easy.DecryptionShare{}
	require.NoError(t, ds.Unmarshal(msg))

	ct := &tdh2easy.Ciphertext{}
	require.NoError(t, ct.UnmarshalVerify(ciphertextBytes, pk))
	gotSecret, err := tdh2easy.Aggregate(ct, []*tdh2easy.DecryptionShare{ds}, 3)
	require.NoError(t, err)
	require.Equal(t, plaintext, gotSecret)

	err = r.ValidateObservation(
		t.Context(),
		1,
		types.AttributedQuery{},
		types.AttributedObservation{Observer: 0, Observation: data},
		rdr,
		&blobber{},
	)
	require.NoError(t, err)
}
func TestPlugin_StateTransition_GetSecretsRequest_ResponseSizeWithinLimit(t *testing.T) {
	_, pk, shares, err := tdh2easy.GenerateKeys(4, 10)
	require.NoError(t, err)

	numObservers := 10
	r := newTestReportingPlugin(
		t,
		withKeys(pk, shares[0]),
		withOnchainCfg(10, 3),
		withMaxCiphertextLengthBytes(2000),
		withMaxIdentifierLengths(64, 64, 64),
		withVaultOptimizationsEnabled(),
	)

	maxOwner := strings.Repeat("a", 64)
	maxNamespace := strings.Repeat("b", 64)

	numEncryptionKeys := 10
	encryptionKeys := make([]string, numEncryptionKeys)
	for i := range numEncryptionKeys {
		pubK, _, err2 := box.GenerateKey(rand.Reader)
		require.NoError(t, err2)
		encryptionKeys[i] = hex.EncodeToString(pubK[:])
	}

	plaintext := make([]byte, 1)
	_, err = rand.Read(plaintext)
	require.NoError(t, err)
	var label [32]byte
	copy(label[:], maxOwner[:32])
	ciphertext, err := tdh2easy.EncryptWithLabel(pk, plaintext, label)
	require.NoError(t, err)
	ciphertextBytes, err := ciphertext.Marshal()
	require.NoError(t, err)
	require.LessOrEqual(t, len(ciphertextBytes), 2000)
	encryptedValue := hex.EncodeToString(ciphertextBytes)

	// Create 10 observations from different observers, each with a distinct decryption share.
	aos := make([]types.AttributedObservation, numObservers)
	for i := range numObservers {
		aos[i] = types.AttributedObservation{
			Observer:    commontypes.OracleID(i),
			Observation: types.Observation(makeGetSecretsBinaryObservations(t, 10, maxOwner, maxNamespace, encryptionKeys, encryptedValue, ciphertext, shares[i])),
		}
	}

	kvStore := &kv{m: make(map[string]response)}
	reportPrecursor, err := r.StateTransition(
		t.Context(),
		1,
		types.AttributedQuery{},
		aos, kvStore, nil)
	require.NoError(t, err)

	twoFPlusOne := 2*r.onchainCfg.F + 1
	osOut := &vaultcommon.Outcomes{}
	err = proto.Unmarshal(reportPrecursor, osOut)
	require.NoError(t, err)
	for _, outcome := range osOut.Outcomes {
		for _, secretResp := range outcome.GetGetSecretsResponse().GetResponses() {
			data := secretResp.GetData()
			if data == nil {
				continue
			}
			for _, enc := range data.EncryptedDecryptionKeyShares {
				assert.Empty(t, enc.Shares)
				assert.Len(t, enc.BinaryShares, twoFPlusOne,
					"expected at most 2f+1 binary shares per encryption key, got %d (N=%d)", len(enc.BinaryShares), numObservers)
			}
		}
	}

	t.Logf("StateTransition response size: %d bytes (%.2f KB)", len(reportPrecursor), float64(len(reportPrecursor))/1024.0)
	maxResponseSize := 512 * 1024
	assert.LessOrEqual(t, len(reportPrecursor), maxResponseSize,
		"StateTransition response size %d exceeds 512KB limit", len(reportPrecursor))
}
func TestPlugin_StateTransition_GetSecretsRequest_CombinesShares(t *testing.T) {
	lggr, observed := logger.TestLoggerObserved(t, zapcore.DebugLevel)
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(t, withLggr(lggr), withKeys(pk, shares[0]), withOnchainCfg(4, 1))

	seqNr := uint64(1)
	kv := &kv{
		m: make(map[string]response),
	}

	id := &vaultcommon.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "secret",
	}
	req := &vaultcommon.GetSecretsRequest{
		Requests: []*vaultcommon.SecretRequest{
			{
				Id: id,
			},
		},
	}
	resp1 := &vaultcommon.GetSecretsResponse{
		Responses: []*vaultcommon.SecretResponse{
			{
				Id: id,
				Result: &vaultcommon.SecretResponse_Data{
					Data: &vaultcommon.SecretData{
						EncryptedValue: "encrypted-value",
						EncryptedDecryptionKeyShares: []*vaultcommon.EncryptedShares{
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
	resp2 := &vaultcommon.GetSecretsResponse{
		Responses: []*vaultcommon.SecretResponse{
			{
				Id: id,
				Result: &vaultcommon.SecretResponse_Data{
					Data: &vaultcommon.SecretData{
						EncryptedValue: "encrypted-value",
						EncryptedDecryptionKeyShares: []*vaultcommon.EncryptedShares{
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
	resp3 := &vaultcommon.GetSecretsResponse{
		Responses: []*vaultcommon.SecretResponse{
			{
				Id: id,
				Result: &vaultcommon.SecretResponse_Data{
					Data: &vaultcommon.SecretData{
						EncryptedValue: "encrypted-value",
						EncryptedDecryptionKeyShares: []*vaultcommon.EncryptedShares{
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
			{Observer: 0, Observation: types.Observation(obsb1)},
			{Observer: 1, Observation: types.Observation(obsb2)},
			{Observer: 2, Observation: types.Observation(obsb3)},
		}, kv, nil)
	require.NoError(t, err)

	os := &vaultcommon.Outcomes{}
	err = proto.Unmarshal(reportPrecursor, os)
	require.NoError(t, err)

	assert.Len(t, os.Outcomes, 1)

	o := os.Outcomes[0]

	expectedResp := &vaultcommon.GetSecretsResponse{
		Responses: []*vaultcommon.SecretResponse{
			{
				Id: id,
				Result: &vaultcommon.SecretResponse_Data{
					Data: &vaultcommon.SecretData{
						EncryptedValue: "encrypted-value",
						EncryptedDecryptionKeyShares: []*vaultcommon.EncryptedShares{
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
func TestPlugin_StateTransition_GetSecretsRequest_CombinesBinaryShares(t *testing.T) {
	lggr, observed := logger.TestLoggerObserved(t, zapcore.DebugLevel)
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(t, withLggr(lggr), withKeys(pk, shares[0]), withOnchainCfg(4, 1))

	seqNr := uint64(1)
	kv := &kv{
		m: make(map[string]response),
	}

	id := &vaultcommon.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "secret",
	}
	req := &vaultcommon.GetSecretsRequest{
		Requests: []*vaultcommon.SecretRequest{
			{
				Id: id,
			},
		},
	}
	share1 := []byte("encrypted-share-1")
	share2 := []byte("encrypted-share-2")
	share3 := []byte("encrypted-share-3")
	resp1 := &vaultcommon.GetSecretsResponse{
		Responses: []*vaultcommon.SecretResponse{
			{
				Id: id,
				Result: &vaultcommon.SecretResponse_Data{
					Data: &vaultcommon.SecretData{
						EncryptedValue: "encrypted-value",
						EncryptedDecryptionKeyShares: []*vaultcommon.EncryptedShares{
							{
								EncryptionKey: "my-encryption-key",
								BinaryShares:  [][]byte{share1},
							},
						},
					},
				},
			},
		},
	}
	resp2 := &vaultcommon.GetSecretsResponse{
		Responses: []*vaultcommon.SecretResponse{
			{
				Id: id,
				Result: &vaultcommon.SecretResponse_Data{
					Data: &vaultcommon.SecretData{
						EncryptedValue: "encrypted-value",
						EncryptedDecryptionKeyShares: []*vaultcommon.EncryptedShares{
							{
								EncryptionKey: "my-encryption-key",
								BinaryShares:  [][]byte{share2},
							},
						},
					},
				},
			},
		},
	}
	resp3 := &vaultcommon.GetSecretsResponse{
		Responses: []*vaultcommon.SecretResponse{
			{
				Id: id,
				Result: &vaultcommon.SecretResponse_Data{
					Data: &vaultcommon.SecretData{
						EncryptedValue: "encrypted-value",
						EncryptedDecryptionKeyShares: []*vaultcommon.EncryptedShares{
							{
								EncryptionKey: "my-encryption-key",
								BinaryShares:  [][]byte{share3},
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
			{Observer: 0, Observation: types.Observation(obsb1)},
			{Observer: 1, Observation: types.Observation(obsb2)},
			{Observer: 2, Observation: types.Observation(obsb3)},
		}, kv, nil)
	require.NoError(t, err)

	os := &vaultcommon.Outcomes{}
	err = proto.Unmarshal(reportPrecursor, os)
	require.NoError(t, err)

	require.Len(t, os.Outcomes, 1)
	o := os.Outcomes[0]

	expectedResp := &vaultcommon.GetSecretsResponse{
		Responses: []*vaultcommon.SecretResponse{
			{
				Id: id,
				Result: &vaultcommon.SecretResponse_Data{
					Data: &vaultcommon.SecretData{
						EncryptedValue: "encrypted-value",
						EncryptedDecryptionKeyShares: []*vaultcommon.EncryptedShares{
							{
								EncryptionKey: "my-encryption-key",
								BinaryShares:  [][]byte{share1, share2, share3},
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
func TestPlugin_StateTransition_GetSecretsRequest_CapsSharesAtTwoFPlusOne(t *testing.T) {
	lggr, observed := logger.TestLoggerObserved(t, zapcore.DebugLevel)
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(t, withLggr(lggr), withKeys(pk, shares[0]), withOnchainCfg(4, 1), withVaultOptimizationsEnabled())

	seqNr := uint64(1)
	kv := &kv{
		m: make(map[string]response),
	}

	id := &vaultcommon.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "secret",
	}
	req := &vaultcommon.GetSecretsRequest{
		Requests: []*vaultcommon.SecretRequest{
			{
				Id: id,
			},
		},
	}
	makeResp := func(share string) *vaultcommon.GetSecretsResponse {
		return &vaultcommon.GetSecretsResponse{
			Responses: []*vaultcommon.SecretResponse{
				{
					Id: id,
					Result: &vaultcommon.SecretResponse_Data{
						Data: &vaultcommon.SecretData{
							EncryptedValue: "encrypted-value",
							EncryptedDecryptionKeyShares: []*vaultcommon.EncryptedShares{
								{
									EncryptionKey: "my-encryption-key",
									Shares:        []string{share},
								},
							},
						},
					},
				},
			},
		}
	}

	obsb1 := marshalObservations(t, observation{id, req, makeResp("encrypted-share-1")})
	obsb2 := marshalObservations(t, observation{id, req, makeResp("encrypted-share-2")})
	obsb3 := marshalObservations(t, observation{id, req, makeResp("encrypted-share-3")})
	obsb4 := marshalObservations(t, observation{id, req, makeResp("encrypted-share-4")})
	reportPrecursor, err := r.StateTransition(
		t.Context(),
		seqNr,
		types.AttributedQuery{},
		[]types.AttributedObservation{
			{Observer: 0, Observation: types.Observation(obsb1)},
			{Observer: 1, Observation: types.Observation(obsb2)},
			{Observer: 2, Observation: types.Observation(obsb3)},
			{Observer: 3, Observation: types.Observation(obsb4)},
		}, kv, nil)
	require.NoError(t, err)

	os := &vaultcommon.Outcomes{}
	err = proto.Unmarshal(reportPrecursor, os)
	require.NoError(t, err)

	assert.Len(t, os.Outcomes, 1)

	o := os.Outcomes[0]

	twoFPlusOne := 2*r.onchainCfg.F + 1
	got := o.GetGetSecretsResponse().Responses[0].GetData().EncryptedDecryptionKeyShares[0].Shares
	assert.Len(t, got, twoFPlusOne)
	expectedResp := &vaultcommon.GetSecretsResponse{
		Responses: []*vaultcommon.SecretResponse{
			{
				Id: id,
				Result: &vaultcommon.SecretResponse_Data{
					Data: &vaultcommon.SecretData{
						EncryptedValue: "encrypted-value",
						EncryptedDecryptionKeyShares: []*vaultcommon.EncryptedShares{
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
func TestPlugin_StateTransition_GetSecretsRequest_DoesNotCapSharesWhenOptimizationsDisabled(t *testing.T) {
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(t, withKeys(pk, shares[0]), withOnchainCfg(4, 1))

	id := &vaultcommon.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "secret",
	}
	req := &vaultcommon.GetSecretsRequest{
		Requests: []*vaultcommon.SecretRequest{{Id: id}},
	}
	makeResp := func(share string) *vaultcommon.GetSecretsResponse {
		return &vaultcommon.GetSecretsResponse{
			Responses: []*vaultcommon.SecretResponse{
				{
					Id: id,
					Result: &vaultcommon.SecretResponse_Data{
						Data: &vaultcommon.SecretData{
							EncryptedValue: "encrypted-value",
							EncryptedDecryptionKeyShares: []*vaultcommon.EncryptedShares{
								{EncryptionKey: "my-encryption-key", Shares: []string{share}},
							},
						},
					},
				},
			},
		}
	}

	reportPrecursor, err := r.StateTransition(
		t.Context(),
		1,
		types.AttributedQuery{},
		[]types.AttributedObservation{
			{Observer: 0, Observation: types.Observation(marshalObservations(t, observation{id, req, makeResp("share-1")}))},
			{Observer: 1, Observation: types.Observation(marshalObservations(t, observation{id, req, makeResp("share-2")}))},
			{Observer: 2, Observation: types.Observation(marshalObservations(t, observation{id, req, makeResp("share-3")}))},
			{Observer: 3, Observation: types.Observation(marshalObservations(t, observation{id, req, makeResp("share-4")}))},
		},
		&kv{m: make(map[string]response)},
		nil,
	)
	require.NoError(t, err)

	os := &vaultcommon.Outcomes{}
	require.NoError(t, proto.Unmarshal(reportPrecursor, os))
	require.Len(t, os.Outcomes, 1)

	got := os.Outcomes[0].GetGetSecretsResponse().Responses[0].GetData().EncryptedDecryptionKeyShares[0].Shares
	assert.Len(t, got, 4)
	require.NotNil(t, os.Outcomes[0].GetGetSecretsRequest())
}
func TestPlugin_StateTransition_GetSecretsRequest_OmitsOutcomeRequestWhenOptimizationsEnabled(t *testing.T) {
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(t, withKeys(pk, shares[0]), withOnchainCfg(4, 1), withVaultOptimizationsEnabled())

	id := &vaultcommon.SecretIdentifier{Owner: "owner", Namespace: "main", Key: "secret"}
	req := &vaultcommon.GetSecretsRequest{Requests: []*vaultcommon.SecretRequest{{Id: id}}}
	resp := &vaultcommon.GetSecretsResponse{
		Responses: []*vaultcommon.SecretResponse{
			{
				Id: id,
				Result: &vaultcommon.SecretResponse_Data{
					Data: &vaultcommon.SecretData{
						EncryptedValue: "encrypted-value",
						EncryptedDecryptionKeyShares: []*vaultcommon.EncryptedShares{
							{EncryptionKey: "my-encryption-key", Shares: []string{"share-1"}},
						},
					},
				},
			},
		},
	}
	obsb := marshalObservations(t, observation{id, req, resp})

	reportPrecursor, err := r.StateTransition(
		t.Context(),
		1,
		types.AttributedQuery{},
		[]types.AttributedObservation{
			{Observer: 0, Observation: types.Observation(obsb)},
			{Observer: 1, Observation: types.Observation(obsb)},
			{Observer: 2, Observation: types.Observation(obsb)},
		},
		&kv{m: make(map[string]response)},
		nil,
	)
	require.NoError(t, err)

	os := &vaultcommon.Outcomes{}
	require.NoError(t, proto.Unmarshal(reportPrecursor, os))
	require.Len(t, os.Outcomes, 1)
	assert.Nil(t, os.Outcomes[0].Request)
}
