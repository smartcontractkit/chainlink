package artifacts

import (
	"context"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"time"

	"github.com/jonboulle/clockwork"

	"testing"

	"github.com/smartcontractkit/chainlink-common/pkg/custmsg"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/secrets"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	ghcapabilities "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/workflowkey"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Handler_SecretsFor(t *testing.T) {
	lggr := logger.TestLogger(t)
	db := pgtest.NewSqlxDB(t)
	orm := &orm{ds: db, lggr: lggr}

	workflowOwner := hex.EncodeToString([]byte("anOwner"))
	workflowName := "aName"
	workflowID := "anID"
	decodedWorkflowName := "decodedName"
	encryptionKey, err := workflowkey.New()
	require.NoError(t, err)

	url := "http://example.com"
	hash := hex.EncodeToString([]byte(url))
	secretsPayload, err := generateSecrets(workflowOwner, map[string][]string{"Foo": {"Bar"}}, encryptionKey)
	require.NoError(t, err)
	secretsID, err := orm.Create(testutils.Context(t), url, hash, string(secretsPayload))
	require.NoError(t, err)

	_, err = orm.UpsertWorkflowSpec(testutils.Context(t), &job.WorkflowSpec{
		Workflow:      "",
		Config:        "",
		SecretsID:     sql.NullInt64{Int64: secretsID, Valid: true},
		WorkflowID:    workflowID,
		WorkflowOwner: workflowOwner,
		WorkflowName:  workflowName,
		BinaryURL:     "",
		ConfigURL:     "",
		CreatedAt:     time.Now(),
		SpecType:      job.DefaultSpecType,
	})
	require.NoError(t, err)

	fetcher := &mockFetcher{
		responseMap: map[string]mockFetchResp{
			url: {Err: errors.New("could not fetch")},
		},
	}

	require.NoError(t, err)
	h := NewStore(
		lggr,
		orm,
		fetcher.Fetch,
		clockwork.NewFakeClock(),
		encryptionKey,
		custmsg.NewLabeler(),
	)
	expectedSecrets := map[string]string{
		"Foo": "Bar",
	}
	decrypter := newMockDecrypter()
	decrypter.registerMock(secretsPayload, workflowOwner, expectedSecrets, nil)
	h.decryptSecrets = decrypter.decryptSecrets
	gotSecrets, err := h.SecretsFor(testutils.Context(t), workflowOwner, workflowName, decodedWorkflowName, workflowID)
	require.NoError(t, err)

	assert.Equal(t, expectedSecrets, gotSecrets)
}

func Test_Handler_SecretsFor_RefreshesSecrets(t *testing.T) {
	lggr := logger.TestLogger(t)
	db := pgtest.NewSqlxDB(t)
	orm := &orm{ds: db, lggr: lggr}

	workflowOwner := hex.EncodeToString([]byte("anOwner"))
	workflowName := "aName"
	workflowID := "anID"
	decodedWorkflowName := "decodedName"
	encryptionKey, err := workflowkey.New()
	require.NoError(t, err)

	secretsPayload, err := generateSecrets(workflowOwner, map[string][]string{"Foo": {"Bar"}}, encryptionKey)
	require.NoError(t, err)

	url := "http://example.com"
	hash := hex.EncodeToString([]byte(url))

	secretsID, err := orm.Create(testutils.Context(t), url, hash, string(secretsPayload))
	require.NoError(t, err)

	_, err = orm.UpsertWorkflowSpec(testutils.Context(t), &job.WorkflowSpec{
		Workflow:      "",
		Config:        "",
		SecretsID:     sql.NullInt64{Int64: secretsID, Valid: true},
		WorkflowID:    workflowID,
		WorkflowOwner: workflowOwner,
		WorkflowName:  workflowName,
		BinaryURL:     "",
		ConfigURL:     "",
		CreatedAt:     time.Now(),
		SpecType:      job.DefaultSpecType,
	})
	require.NoError(t, err)

	secretsPayload, err = generateSecrets(workflowOwner, map[string][]string{"Baz": {"Bar"}}, encryptionKey)
	require.NoError(t, err)
	fetcher := &mockFetcher{
		responseMap: map[string]mockFetchResp{
			url: {Body: secretsPayload},
		},
	}

	h := NewStore(
		lggr,
		orm,
		fetcher.Fetch,
		clockwork.NewFakeClock(),
		encryptionKey,
		custmsg.NewLabeler(),
	)

	expectedSecrets := map[string]string{
		"Baz": "Bar",
	}
	decrypter := newMockDecrypter()
	decrypter.registerMock(secretsPayload, workflowOwner, expectedSecrets, nil)
	h.decryptSecrets = decrypter.decryptSecrets

	gotSecrets, err := h.SecretsFor(testutils.Context(t), workflowOwner, workflowName, decodedWorkflowName, workflowID)
	require.NoError(t, err)

	assert.Equal(t, expectedSecrets, gotSecrets)
}

func Test_Handler_SecretsFor_RefreshLogic(t *testing.T) {
	lggr := logger.TestLogger(t)
	db := pgtest.NewSqlxDB(t)
	orm := &orm{ds: db, lggr: lggr}

	workflowOwner := hex.EncodeToString([]byte("anOwner"))
	workflowName := "aName"
	workflowID := "anID"
	decodedWorkflowName := "decodedName"
	encryptionKey, err := workflowkey.New()
	require.NoError(t, err)

	secretsPayload, err := generateSecrets(workflowOwner, map[string][]string{"Foo": {"Bar"}}, encryptionKey)
	require.NoError(t, err)

	url := "http://example.com"
	hash := hex.EncodeToString([]byte(url))

	secretsID, err := orm.Create(testutils.Context(t), url, hash, string(secretsPayload))
	require.NoError(t, err)

	_, err = orm.UpsertWorkflowSpec(testutils.Context(t), &job.WorkflowSpec{
		Workflow:      "",
		Config:        "",
		SecretsID:     sql.NullInt64{Int64: secretsID, Valid: true},
		WorkflowID:    workflowID,
		WorkflowOwner: workflowOwner,
		WorkflowName:  workflowName,
		BinaryURL:     "",
		ConfigURL:     "",
		CreatedAt:     time.Now(),
		SpecType:      job.DefaultSpecType,
	})
	require.NoError(t, err)

	fetcher := &mockFetcher{
		responseMap: map[string]mockFetchResp{
			url: {
				Body: secretsPayload,
			},
		},
	}
	clock := clockwork.NewFakeClock()

	h := NewStore(
		lggr,
		orm,
		fetcher.Fetch,
		clock,
		encryptionKey,
		custmsg.NewLabeler(),
	)

	expectedSecrets := map[string]string{
		"Foo": "Bar",
	}
	decrypter := newMockDecrypter()
	decrypter.registerMock(secretsPayload, workflowOwner, expectedSecrets, nil)
	h.decryptSecrets = decrypter.decryptSecrets

	gotSecrets, err := h.SecretsFor(testutils.Context(t), workflowOwner, workflowName, decodedWorkflowName, workflowID)
	require.NoError(t, err)

	assert.Equal(t, expectedSecrets, gotSecrets)

	// Now stub out an unparseable response, since we already fetched it recently above, we shouldn't need to refetch
	// SecretsFor should still succeed.
	fetcher.responseMap[url] = mockFetchResp{}

	gotSecrets, err = h.SecretsFor(testutils.Context(t), workflowOwner, workflowName, decodedWorkflowName, workflowID)
	require.NoError(t, err)

	assert.Equal(t, expectedSecrets, gotSecrets)

	secretsPayload, err = generateSecrets(workflowOwner, map[string][]string{"Baz": {"Bar"}}, encryptionKey)
	require.NoError(t, err)
	fetcher.responseMap[url] = mockFetchResp{
		Body: secretsPayload,
	}

	expectedSecrets = map[string]string{
		"Baz": "Bar",
	}
	decrypter.registerMock(secretsPayload, workflowOwner, expectedSecrets, nil)

	// Now advance so that we hit the freshness limit
	clock.Advance(48 * time.Hour)

	gotSecrets, err = h.SecretsFor(testutils.Context(t), workflowOwner, workflowName, decodedWorkflowName, workflowID)
	require.NoError(t, err)
	assert.Equal(t, expectedSecrets, gotSecrets)
}

func generateSecrets(workflowOwner string, secretsMap map[string][]string, encryptionKey workflowkey.Key) ([]byte, error) {
	sm, secretsEnvVars, err := secrets.EncryptSecretsForNodes(
		workflowOwner,
		secretsMap,
		map[string][32]byte{
			"p2pId": encryptionKey.PublicKey(),
		},
		secrets.SecretsConfig{},
	)
	if err != nil {
		return nil, err
	}
	return json.Marshal(secrets.EncryptedSecretsResult{
		EncryptedSecrets: sm,
		Metadata: secrets.Metadata{
			WorkflowOwner:          workflowOwner,
			EnvVarsAssignedToNodes: secretsEnvVars,
			NodePublicEncryptionKeys: map[string]string{
				"p2pId": encryptionKey.PublicKeyString(),
			},
		},
	})
}

type mockFetchResp struct {
	Body []byte
	Err  error
}

type mockFetcher struct {
	responseMap map[string]mockFetchResp
}

func (m *mockFetcher) Fetch(_ context.Context, mid string, req ghcapabilities.Request) ([]byte, error) {
	return m.responseMap[req.URL].Body, m.responseMap[req.URL].Err
}

type decryptSecretsOutput struct {
	output map[string]string
	err    error
}

type mockDecrypter struct {
	mocks map[string]decryptSecretsOutput
}

func (m *mockDecrypter) decryptSecrets(data []byte, owner string) (map[string]string, error) {
	input := string(data) + owner
	mock, exists := m.mocks[input]
	if exists {
		return mock.output, mock.err
	}
	return map[string]string{}, nil
}

func (m *mockDecrypter) registerMock(data []byte, owner string, output map[string]string, err error) {
	input := string(data) + owner
	m.mocks[input] = decryptSecretsOutput{output: output, err: err}
}

func newMockDecrypter() *mockDecrypter {
	return &mockDecrypter{
		mocks: map[string]decryptSecretsOutput{},
	}
}
