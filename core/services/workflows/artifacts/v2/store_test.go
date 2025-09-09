package v2

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"time"

	"github.com/jonboulle/clockwork"

	"testing"

	"github.com/smartcontractkit/chainlink-common/pkg/custmsg"
	storage_service "github.com/smartcontractkit/chainlink-protos/storage-service/go"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	ghcapabilities "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/workflowkey"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"

	"github.com/stretchr/testify/require"
)

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

func (m *mockFetcher) RetrieveURL(ctx context.Context, req *storage_service.DownloadArtifactRequest) (string, error) {
	return string(m.responseMap[req.Id+"-"+req.Type.String()].Body), m.responseMap[req.Id+"-"+req.Type.String()].Err
}

func Test_Store_DeleteWorkflowArtifacts(t *testing.T) {
	lggr := logger.TestLogger(t)
	db := pgtest.NewSqlxDB(t)
	orm := &orm{ds: db, lggr: lggr}

	workflowOwner := hex.EncodeToString([]byte("anOwner"))
	workflowName := "aName"
	workflowID := "anID"
	encryptionKey, err := workflowkey.New()
	require.NoError(t, err)

	_, err = orm.UpsertWorkflowSpec(testutils.Context(t), &job.WorkflowSpec{
		Workflow:      "",
		Config:        "",
		SecretsID:     sql.NullInt64{Int64: 0, Valid: false},
		WorkflowID:    workflowID,
		WorkflowOwner: workflowOwner,
		WorkflowName:  workflowName,
		BinaryURL:     "",
		ConfigURL:     "",
		CreatedAt:     time.Now(),
		SpecType:      job.DefaultSpecType,
	})
	require.NoError(t, err)

	fetcher := &mockFetcher{}

	h, err := NewStore(
		lggr,
		orm,
		fetcher.Fetch,
		fetcher.RetrieveURL,
		clockwork.NewFakeClock(),
		encryptionKey,
		custmsg.NewLabeler(),
		WithConfig(StoreConfig{
			ArtifactStorageHost: "example.com",
		}),
	)
	require.NoError(t, err)

	// Delete the workflow artifacts by ID
	err = h.DeleteWorkflowArtifacts(testutils.Context(t), workflowID)
	require.NoError(t, err)

	// Check that the workflow no longer exists
	_, err = orm.GetWorkflowSpec(testutils.Context(t), workflowID)
	require.ErrorIs(t, err, sql.ErrNoRows)
}

func Test_Store_FetchWorkflowArtifacts_WithStorage(t *testing.T) {
	lggr := logger.TestLogger(t)
	db := pgtest.NewSqlxDB(t)
	orm := &orm{ds: db, lggr: lggr}

	workflowID := "anID"
	encryptionKey, err := workflowkey.New()
	require.NoError(t, err)

	binaryURL := "http://storage.chain.link/" + workflowID + "/binary.wasm"
	binarySignedURL := binaryURL + "?auth=XXX"
	binaryData := "binary-data"
	binaryEncoded := base64.StdEncoding.EncodeToString([]byte(binaryData))
	configURL := "http://storage.chain.link/" + workflowID + "/config.yaml"
	configSignedURL := configURL + "?auth=XXX"
	configData := "config-data"
	fetcher := &mockFetcher{
		responseMap: map[string]mockFetchResp{
			workflowID + "-ARTIFACT_TYPE_BINARY": {Body: []byte(binarySignedURL)},
			binarySignedURL:                      {Body: []byte(binaryEncoded)},
			workflowID + "-ARTIFACT_TYPE_CONFIG": {Body: []byte(configSignedURL)},
			configSignedURL:                      {Body: []byte(configData)},
		},
	}

	h, err := NewStore(
		lggr,
		orm,
		fetcher.Fetch,
		fetcher.RetrieveURL,
		clockwork.NewFakeClock(),
		encryptionKey,
		custmsg.NewLabeler(),
		WithConfig(StoreConfig{
			ArtifactStorageHost: "storage.chain.link",
		}),
	)
	require.NoError(t, err)

	binary, config, err := h.FetchWorkflowArtifacts(testutils.Context(t), workflowID, binaryURL, configURL)
	require.NoError(t, err)
	require.Equal(t, []byte(binaryData), binary)
	require.Equal(t, []byte(configData), config)
}

func Test_Store_FetchWorkflowArtifacts_WithoutStorage(t *testing.T) {
	lggr := logger.TestLogger(t)
	db := pgtest.NewSqlxDB(t)
	orm := &orm{ds: db, lggr: lggr}

	workflowID := "anID"
	encryptionKey, err := workflowkey.New()
	require.NoError(t, err)

	binaryURL := "http://some-url.com/binary.wasm"
	binaryData := "binary-data"
	binaryEncoded := base64.StdEncoding.EncodeToString([]byte(binaryData))
	configURL := "http://some-url.com/config.yaml"
	configData := "config-data"
	fetcher := &mockFetcher{
		responseMap: map[string]mockFetchResp{
			binaryURL: {Body: []byte(binaryEncoded)},
			configURL: {Body: []byte(configData)},
		},
	}

	h, err := NewStore(
		lggr,
		orm,
		fetcher.Fetch,
		fetcher.RetrieveURL,
		clockwork.NewFakeClock(),
		encryptionKey,
		custmsg.NewLabeler(),
		WithConfig(StoreConfig{
			ArtifactStorageHost: "storage.chain.link",
		}),
	)
	require.NoError(t, err)

	binary, config, err := h.FetchWorkflowArtifacts(testutils.Context(t), workflowID, binaryURL, configURL)
	require.NoError(t, err)
	require.Equal(t, []byte(binaryData), binary)
	require.Equal(t, []byte(configData), config)
}

func Test_Store_FetchWorkflowArtifacts_SkipsRetrieving(t *testing.T) {
	lggr := logger.TestLogger(t)
	db := pgtest.NewSqlxDB(t)
	orm := &orm{ds: db, lggr: lggr}

	workflowID := "anID"
	encryptionKey, err := workflowkey.New()
	require.NoError(t, err)

	binaryURL := "http://example.com/id1/binary.wasm"
	binaryData := "binary-data"
	binaryEncoded := base64.StdEncoding.EncodeToString([]byte(binaryData))
	configURL := "http://example.com/id1/config.yaml"
	configData := "config-data"
	fetcher := &mockFetcher{
		responseMap: map[string]mockFetchResp{
			binaryURL: {Body: []byte(binaryEncoded)},
			configURL: {Body: []byte(configData)},
		},
	}

	h, err := NewStore(
		lggr,
		orm,
		fetcher.Fetch,
		nil, // No retrieval function provided, so it should skip retrieving
		clockwork.NewFakeClock(),
		encryptionKey,
		custmsg.NewLabeler(),
		WithConfig(StoreConfig{
			ArtifactStorageHost: "example.com",
		}),
	)
	require.NoError(t, err)

	binary, config, err := h.FetchWorkflowArtifacts(testutils.Context(t), workflowID, binaryURL, configURL)
	require.NoError(t, err)
	require.Equal(t, []byte(binaryData), binary)
	require.Equal(t, []byte(configData), config)
}
