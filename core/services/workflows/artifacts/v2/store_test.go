package v2

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/workflowkey"
	"github.com/smartcontractkit/chainlink-common/pkg/contexts"
	"github.com/smartcontractkit/chainlink-common/pkg/custmsg"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	storage_service "github.com/smartcontractkit/chainlink-protos/storage-service/go"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	ghcapabilities "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
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
	t.Parallel()
	lggr := logger.TestLogger(t)
	db := pgtest.NewSqlxDB(t)
	orm := &orm{ds: db, lggr: lggr}

	workflowOwner := hex.EncodeToString([]byte("owner-" + uuid.New().String()[:8]))
	workflowName := "name-" + uuid.New().String()[:8]
	workflowID := "id-" + uuid.New().String()[:8]
	encryptionKey, err := workflowkey.New()
	require.NoError(t, err)

	_, err = orm.UpsertWorkflowSpec(t.Context(), &job.WorkflowSpec{
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
		limits.Factory{Logger: lggr},
		WithConfig(StoreConfig{
			ArtifactStorageHost: "example.com",
		}),
	)
	require.NoError(t, err)

	// Delete the workflow artifacts by ID
	_, err = h.DeleteWorkflowArtifacts(t.Context(), workflowID)
	require.NoError(t, err)

	// Check that the workflow no longer exists
	_, err = orm.GetWorkflowSpec(t.Context(), workflowID)
	require.ErrorIs(t, err, sql.ErrNoRows)
}

func Test_Store_DeleteWorkflowArtifactsBatch(t *testing.T) {
	t.Parallel()
	lggr := logger.TestLogger(t)
	db := pgtest.NewSqlxDB(t)
	orm := &orm{ds: db, lggr: lggr}

	encryptionKey, err := workflowkey.New()
	require.NoError(t, err)

	prefix := uuid.New().String()[:8]
	wf1 := "wf-1-" + prefix
	wf2 := "wf-2-" + prefix
	wf3 := "wf-3-" + prefix
	owner := "owner-" + prefix

	ctx := t.Context()
	for _, id := range []string{wf1, wf2, wf3} {
		_, err = orm.UpsertWorkflowSpec(ctx, &job.WorkflowSpec{
			Workflow:      "",
			Config:        "",
			WorkflowID:    id,
			WorkflowOwner: owner,
			WorkflowName:  "name-" + id,
			CreatedAt:     time.Now(),
			SpecType:      job.DefaultSpecType,
		})
		require.NoError(t, err)
	}

	fetcher := &mockFetcher{}
	h, err := NewStore(
		lggr, orm, fetcher.Fetch, fetcher.RetrieveURL,
		clockwork.NewFakeClock(), encryptionKey, custmsg.NewLabeler(),
		limits.Factory{Logger: lggr},
		WithConfig(StoreConfig{ArtifactStorageHost: "example.com"}),
	)
	require.NoError(t, err)

	require.NoError(t, h.DeleteWorkflowArtifactsBatch(ctx, []string{wf1, wf3}))

	_, err = orm.GetWorkflowSpec(ctx, wf1)
	require.ErrorIs(t, err, sql.ErrNoRows)
	_, err = orm.GetWorkflowSpec(ctx, wf3)
	require.ErrorIs(t, err, sql.ErrNoRows)

	got, err := orm.GetWorkflowSpec(ctx, wf2)
	require.NoError(t, err)
	require.Equal(t, "name-"+wf2, got.WorkflowName)
}

func Test_Store_FetchWorkflowArtifacts_WithStorage(t *testing.T) {
	t.Parallel()
	lggr := logger.TestLogger(t)
	db := pgtest.NewSqlxDB(t)
	orm := &orm{ds: db, lggr: lggr}

	workflowID := "id-" + uuid.New().String()[:8]
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
		limits.Factory{Logger: lggr},
		WithConfig(StoreConfig{
			ArtifactStorageHost: "storage.chain.link",
		}),
	)
	require.NoError(t, err)

	ctx := contexts.WithCRE(t.Context(), contexts.CRE{Workflow: workflowID})
	binary, config, err := h.FetchWorkflowArtifacts(ctx, workflowID, binaryURL, configURL)
	require.NoError(t, err)
	require.Equal(t, []byte(binaryData), binary)
	require.Equal(t, []byte(configData), config)
}

func Test_Store_FetchWorkflowArtifacts_WithoutStorage(t *testing.T) {
	t.Parallel()
	lggr := logger.TestLogger(t)
	db := pgtest.NewSqlxDB(t)
	orm := &orm{ds: db, lggr: lggr}

	workflowID := "id-" + uuid.New().String()[:8]
	encryptionKey, err := workflowkey.New()
	require.NoError(t, err)

	binaryURL := fmt.Sprintf("http://some-url-%s.com/binary.wasm", workflowID)
	binaryData := "binary-data"
	binaryEncoded := base64.StdEncoding.EncodeToString([]byte(binaryData))
	configURL := fmt.Sprintf("http://some-url-%s.com/config.yaml", workflowID)
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
		limits.Factory{Logger: lggr},
		WithConfig(StoreConfig{
			ArtifactStorageHost: "storage.chain.link",
		}),
	)
	require.NoError(t, err)

	ctx := contexts.WithCRE(t.Context(), contexts.CRE{Workflow: workflowID})
	binary, config, err := h.FetchWorkflowArtifacts(ctx, workflowID, binaryURL, configURL)
	require.NoError(t, err)
	require.Equal(t, []byte(binaryData), binary)
	require.Equal(t, []byte(configData), config)
}

func Test_Store_FetchWorkflowArtifacts_SkipsRetrieving(t *testing.T) {
	t.Parallel()
	lggr := logger.TestLogger(t)
	db := pgtest.NewSqlxDB(t)
	orm := &orm{ds: db, lggr: lggr}

	workflowID := "id-" + uuid.New().String()[:8]
	encryptionKey, err := workflowkey.New()
	require.NoError(t, err)

	binaryURL := fmt.Sprintf("http://example-%s.com/id1/binary.wasm", workflowID)
	binaryData := "binary-data"
	binaryEncoded := base64.StdEncoding.EncodeToString([]byte(binaryData))
	configURL := fmt.Sprintf("http://example-%s.com/id1/config.yaml", workflowID)
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
		limits.Factory{Logger: lggr},
		WithConfig(StoreConfig{
			ArtifactStorageHost: "example.com",
		}),
	)
	require.NoError(t, err)

	ctx := contexts.WithCRE(t.Context(), contexts.CRE{Workflow: workflowID})
	binary, config, err := h.FetchWorkflowArtifacts(ctx, workflowID, binaryURL, configURL)
	require.NoError(t, err)
	require.Equal(t, []byte(binaryData), binary)
	require.Equal(t, []byte(configData), config)
}
