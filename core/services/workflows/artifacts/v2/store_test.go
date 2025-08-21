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
	return string(m.responseMap[req.Id].Body), m.responseMap[req.Id].Err
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

	h := NewStore(
		lggr,
		orm,
		fetcher.Fetch,
		fetcher.RetrieveURL,
		clockwork.NewFakeClock(),
		encryptionKey,
		custmsg.NewLabeler(),
	)

	// Delete the workflow artifacts by ID
	err = h.DeleteWorkflowArtifacts(testutils.Context(t), workflowID)
	require.NoError(t, err)

	// Check that the workflow no longer exists
	_, err = orm.GetWorkflowSpec(testutils.Context(t), workflowID)
	require.ErrorIs(t, err, sql.ErrNoRows)
}

func Test_Store_FetchWorkflowArtifacts(t *testing.T) {
	lggr := logger.TestLogger(t)
	db := pgtest.NewSqlxDB(t)
	orm := &orm{ds: db, lggr: lggr}

	workflowID := "anID"
	encryptionKey, err := workflowkey.New()
	require.NoError(t, err)

	binaryID := "binary-1"
	binaryURL := "http://example.com/binary"
	binaryData := "binary-data"
	binaryEncoded := base64.StdEncoding.EncodeToString([]byte(binaryData))
	configID := "config-1"
	configURL := "http://example.com/config"
	configData := "config-data"
	fetcher := &mockFetcher{
		responseMap: map[string]mockFetchResp{
			binaryID:  {Body: []byte(binaryURL)},
			binaryURL: {Body: []byte(binaryEncoded)},
			configID:  {Body: []byte(configURL)},
			configURL: {Body: []byte(configData)},
		},
	}

	h := NewStore(
		lggr,
		orm,
		fetcher.Fetch,
		fetcher.RetrieveURL,
		clockwork.NewFakeClock(),
		encryptionKey,
		custmsg.NewLabeler(),
	)

	binary, config, err := h.FetchWorkflowArtifacts(testutils.Context(t), workflowID, binaryID, configID)
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

	binaryURL := "http://example.com/binary"
	binaryData := "binary-data"
	binaryEncoded := base64.StdEncoding.EncodeToString([]byte(binaryData))
	configURL := "http://example.com/config"
	configData := "config-data"
	fetcher := &mockFetcher{
		responseMap: map[string]mockFetchResp{
			binaryURL: {Body: []byte(binaryEncoded)},
			configURL: {Body: []byte(configData)},
		},
	}

	h := NewStore(
		lggr,
		orm,
		fetcher.Fetch,
		nil, // No retrieval function provided, so it should skip retrieving
		clockwork.NewFakeClock(),
		encryptionKey,
		custmsg.NewLabeler(),
	)

	binary, config, err := h.FetchWorkflowArtifacts(testutils.Context(t), workflowID, binaryURL, configURL)
	require.NoError(t, err)
	require.Equal(t, []byte(binaryData), binary)
	require.Equal(t, []byte(configData), config)
}
