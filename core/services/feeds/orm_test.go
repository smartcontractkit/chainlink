package feeds_test

import (
	"context"
	"testing"

	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/services/feeds"
	"github.com/smartcontractkit/chainlink/core/utils/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
)

var (
	uri                       = "http://192.168.0.1"
	name                      = "Chainlink FMS"
	publicKey                 = crypto.PublicKey([]byte("11111111111111111111111111111111"))
	jobTypes                  = pq.StringArray{feeds.JobTypeFluxMonitor, feeds.JobTypeOffchainReporting}
	ocrBootstrapPeerMultiaddr = null.StringFrom("/dns4/ocr-bootstrap.chain.link/tcp/0000/p2p/7777777")
)

type TestORM struct {
	feeds.ORM

	db *gorm.DB
}

func setupORM(t *testing.T) *TestORM {
	t.Helper()

	db := pgtest.NewGormDB(t)
	orm := feeds.NewORM(db)

	return &TestORM{ORM: orm, db: db}
}

func Test_ORM_CreateManager(t *testing.T) {
	t.Parallel()

	orm := setupORM(t)
	mgr := &feeds.FeedsManager{
		URI:                       uri,
		Name:                      name,
		PublicKey:                 publicKey,
		JobTypes:                  jobTypes,
		IsOCRBootstrapPeer:        true,
		OCRBootstrapPeerMultiaddr: ocrBootstrapPeerMultiaddr,
	}

	count, err := orm.CountManagers(context.Background())
	require.NoError(t, err)
	require.Equal(t, int64(0), count)

	id, err := orm.CreateManager(context.Background(), mgr)
	require.NoError(t, err)

	count, err = orm.CountManagers(context.Background())
	require.NoError(t, err)
	require.Equal(t, int64(1), count)

	assert.NotZero(t, id)
}

func Test_ORM_ListManagers(t *testing.T) {
	t.Parallel()

	orm := setupORM(t)
	mgr := &feeds.FeedsManager{
		URI:                       uri,
		Name:                      name,
		PublicKey:                 publicKey,
		JobTypes:                  jobTypes,
		IsOCRBootstrapPeer:        true,
		OCRBootstrapPeerMultiaddr: ocrBootstrapPeerMultiaddr,
	}

	id, err := orm.CreateManager(context.Background(), mgr)
	require.NoError(t, err)

	mgrs, err := orm.ListManagers(context.Background())
	require.NoError(t, err)
	require.Len(t, mgrs, 1)

	actual := mgrs[0]
	assert.Equal(t, id, actual.ID)
	assert.Equal(t, uri, actual.URI)
	assert.Equal(t, name, actual.Name)
	assert.Equal(t, publicKey, actual.PublicKey)
	assert.Equal(t, jobTypes, actual.JobTypes)
	assert.True(t, actual.IsOCRBootstrapPeer)
	assert.Equal(t, ocrBootstrapPeerMultiaddr, actual.OCRBootstrapPeerMultiaddr)
}

func Test_ORM_GetManager(t *testing.T) {
	t.Parallel()

	orm := setupORM(t)
	mgr := &feeds.FeedsManager{
		URI:                       uri,
		Name:                      name,
		PublicKey:                 publicKey,
		JobTypes:                  jobTypes,
		IsOCRBootstrapPeer:        true,
		OCRBootstrapPeerMultiaddr: ocrBootstrapPeerMultiaddr,
	}

	id, err := orm.CreateManager(context.Background(), mgr)
	require.NoError(t, err)

	actual, err := orm.GetManager(context.Background(), id)
	require.NoError(t, err)

	assert.Equal(t, id, actual.ID)
	assert.Equal(t, uri, actual.URI)
	assert.Equal(t, name, actual.Name)
	assert.Equal(t, publicKey, actual.PublicKey)
	assert.Equal(t, jobTypes, actual.JobTypes)
	assert.True(t, actual.IsOCRBootstrapPeer)
	assert.Equal(t, ocrBootstrapPeerMultiaddr, actual.OCRBootstrapPeerMultiaddr)

	actual, err = orm.GetManager(context.Background(), -1)
	require.Nil(t, actual)
	require.Error(t, err)
}

func Test_ORM_UpdateManager(t *testing.T) {
	t.Parallel()

	orm := setupORM(t)
	mgr := &feeds.FeedsManager{
		URI:                uri,
		Name:               name,
		PublicKey:          publicKey,
		JobTypes:           jobTypes,
		IsOCRBootstrapPeer: false,
	}

	id, err := orm.CreateManager(context.Background(), mgr)
	require.NoError(t, err)

	updatedMgr := feeds.FeedsManager{
		ID:                        id,
		URI:                       "127.0.0.1",
		Name:                      "New Name",
		PublicKey:                 crypto.PublicKey([]byte("22222222222222222222222222222222")),
		JobTypes:                  pq.StringArray{feeds.JobTypeFluxMonitor},
		IsOCRBootstrapPeer:        true,
		OCRBootstrapPeerMultiaddr: ocrBootstrapPeerMultiaddr,
	}

	err = orm.UpdateManager(context.Background(), updatedMgr)
	require.NoError(t, err)

	actual, err := orm.GetManager(context.Background(), id)
	require.NoError(t, err)

	assert.Equal(t, updatedMgr.URI, actual.URI)
	assert.Equal(t, updatedMgr.Name, actual.Name)
	assert.Equal(t, updatedMgr.PublicKey, actual.PublicKey)
	assert.Equal(t, updatedMgr.JobTypes, actual.JobTypes)
	assert.True(t, actual.IsOCRBootstrapPeer)
	assert.Equal(t, ocrBootstrapPeerMultiaddr, actual.OCRBootstrapPeerMultiaddr)
}

func Test_ORM_CreateJobProposal(t *testing.T) {
	t.Parallel()

	orm := setupORM(t)
	fmID := createFeedsManager(t, orm)

	jp := &feeds.JobProposal{
		RemoteUUID:     uuid.NewV4(),
		Spec:           "",
		Status:         feeds.JobProposalStatusPending,
		FeedsManagerID: fmID,
	}

	count, err := orm.CountJobProposals(context.Background())
	require.NoError(t, err)
	require.Equal(t, int64(0), count)

	id, err := orm.CreateJobProposal(context.Background(), jp)
	require.NoError(t, err)

	count, err = orm.CountJobProposals(context.Background())
	require.NoError(t, err)
	require.Equal(t, int64(1), count)

	assert.NotZero(t, id)
}

func Test_ORM_UpsertJobProposal(t *testing.T) {
	t.Parallel()

	orm := setupORM(t)
	fmID := createFeedsManager(t, orm)
	ctx := context.Background()

	jp := &feeds.JobProposal{
		RemoteUUID:     uuid.NewV4(),
		Spec:           "",
		Status:         feeds.JobProposalStatusPending,
		FeedsManagerID: fmID,
	}

	// Create
	count, err := orm.CountJobProposals(ctx)
	require.NoError(t, err)
	require.Equal(t, int64(0), count)

	id, err := orm.UpsertJobProposal(ctx, jp)
	require.NoError(t, err)

	createdActual, err := orm.GetJobProposal(ctx, id)
	require.NoError(t, err)

	count, err = orm.CountJobProposals(ctx)
	require.NoError(t, err)
	require.Equal(t, int64(1), count)

	assert.NotZero(t, id)

	// Update
	jp.Spec = "updated"
	jp.Status = feeds.JobProposalStatusRejected
	jp.Multiaddrs = pq.StringArray{"dns/example.com"}

	id, err = orm.UpsertJobProposal(ctx, jp)
	require.NoError(t, err)

	actual, err := orm.GetJobProposal(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, jp.Spec, actual.Spec)
	assert.Equal(t, jp.Status, actual.Status)
	assert.Equal(t, jp.Multiaddrs, actual.Multiaddrs)

	// Ensure there is a difference in the created proposal and the upserted
	// proposal
	assert.NotEqual(t, createdActual.Spec, actual.Spec)
	assert.NotEqual(t, createdActual.Status, actual.Status)
	assert.NotEqual(t, createdActual.Multiaddrs, actual.Multiaddrs)
	assert.NotEqual(t, createdActual.UpdatedAt, actual.UpdatedAt)
	assert.Equal(t, createdActual.CreatedAt, actual.CreatedAt) // CreatedAt does not change
}

func Test_ORM_ListJobProposals(t *testing.T) {
	t.Parallel()

	orm := setupORM(t)
	fmID := createFeedsManager(t, orm)
	uuid := uuid.NewV4()

	jp := &feeds.JobProposal{
		RemoteUUID:     uuid,
		Spec:           "",
		Status:         feeds.JobProposalStatusPending,
		FeedsManagerID: fmID,
	}

	id, err := orm.CreateJobProposal(context.Background(), jp)
	require.NoError(t, err)

	jps, err := orm.ListJobProposals(context.Background())
	require.NoError(t, err)
	require.Len(t, jps, 1)

	actual := jps[0]
	assert.Equal(t, id, actual.ID)
	assert.Equal(t, uuid, actual.RemoteUUID)
	assert.Equal(t, jp.Status, actual.Status)
	assert.False(t, actual.ExternalJobID.Valid)
	assert.Equal(t, jp.FeedsManagerID, actual.FeedsManagerID)
}

func Test_ORM_UpdateJobProposalSpec(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	orm := setupORM(t)
	fmID := createFeedsManager(t, orm)

	jp := &feeds.JobProposal{
		RemoteUUID:     uuid.NewV4(),
		Spec:           "",
		Status:         feeds.JobProposalStatusPending,
		FeedsManagerID: fmID,
	}

	id, err := orm.CreateJobProposal(ctx, jp)
	require.NoError(t, err)

	err = orm.UpdateJobProposalSpec(ctx, id, "updated spec")
	require.NoError(t, err)

	actual, err := orm.GetJobProposal(context.Background(), id)
	require.NoError(t, err)

	assert.Equal(t, id, actual.ID)
	assert.Equal(t, "updated spec", actual.Spec)
}

func Test_ORM_GetJobProposal(t *testing.T) {
	t.Parallel()

	orm := setupORM(t)
	fmID := createFeedsManager(t, orm)
	remoteUUID := uuid.NewV4()

	jp := &feeds.JobProposal{
		RemoteUUID:     remoteUUID,
		Spec:           "",
		Status:         feeds.JobProposalStatusPending,
		FeedsManagerID: fmID,
	}

	id, err := orm.CreateJobProposal(context.Background(), jp)
	require.NoError(t, err)

	assertJobEquals := func(actual *feeds.JobProposal) {
		assert.Equal(t, id, actual.ID)
		assert.Equal(t, remoteUUID, actual.RemoteUUID)
		assert.Equal(t, jp.Status, actual.Status)
		assert.False(t, actual.ExternalJobID.Valid)
		assert.Equal(t, jp.FeedsManagerID, actual.FeedsManagerID)
	}

	t.Run("by id", func(t *testing.T) {
		actual, err := orm.GetJobProposal(context.Background(), id)
		require.NoError(t, err)

		assert.Equal(t, id, actual.ID)
		assertJobEquals(actual)

		actual, err = orm.GetJobProposal(context.Background(), int64(0))
		require.Nil(t, actual)
		require.Error(t, err)
	})

	t.Run("by remote uuid", func(t *testing.T) {
		actual, err := orm.GetJobProposalByRemoteUUID(context.Background(), remoteUUID)
		require.NoError(t, err)

		assertJobEquals(actual)

		actual, err = orm.GetJobProposalByRemoteUUID(context.Background(), uuid.NewV4())
		require.Nil(t, actual)
		require.Error(t, err)
	})
}

func Test_ORM_UpdateJobProposalStatus(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	orm := setupORM(t)
	fmID := createFeedsManager(t, orm)

	jp := &feeds.JobProposal{
		RemoteUUID:     uuid.NewV4(),
		Spec:           "",
		Status:         feeds.JobProposalStatusPending,
		FeedsManagerID: fmID,
	}

	id, err := orm.CreateJobProposal(ctx, jp)
	require.NoError(t, err)

	err = orm.UpdateJobProposalStatus(ctx, id, feeds.JobProposalStatusRejected)
	require.NoError(t, err)

	actual, err := orm.GetJobProposal(context.Background(), id)
	require.NoError(t, err)

	assert.Equal(t, id, actual.ID)
	assert.Equal(t, feeds.JobProposalStatusRejected, actual.Status)
}

func Test_ORM_ApproveJobProposal(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	orm := setupORM(t)
	fmID := createFeedsManager(t, orm)
	externalJobID := uuid.NullUUID{UUID: uuid.NewV4(), Valid: true}

	jp := &feeds.JobProposal{
		RemoteUUID:     uuid.NewV4(),
		Spec:           "",
		Status:         feeds.JobProposalStatusPending,
		FeedsManagerID: fmID,
	}

	// Defer the FK requirement of a job proposal.
	require.NoError(t, orm.db.Exec(
		`SET CONSTRAINTS job_proposals_job_id_fkey DEFERRED`,
	).Error)

	id, err := orm.CreateJobProposal(ctx, jp)
	require.NoError(t, err)

	err = orm.ApproveJobProposal(ctx, id, externalJobID.UUID, feeds.JobProposalStatusApproved)
	require.NoError(t, err)

	actual, err := orm.GetJobProposal(context.Background(), id)
	require.NoError(t, err)

	assert.Equal(t, id, actual.ID)
	assert.Equal(t, externalJobID, actual.ExternalJobID)
	assert.Equal(t, feeds.JobProposalStatusApproved, actual.Status)
}

// createFeedsManager is a test helper to create a feeds manager
func createFeedsManager(t *testing.T, orm feeds.ORM) int64 {
	mgr := &feeds.FeedsManager{
		URI:                uri,
		Name:               name,
		PublicKey:          publicKey,
		JobTypes:           jobTypes,
		IsOCRBootstrapPeer: false,
	}

	id, err := orm.CreateManager(context.Background(), mgr)
	require.NoError(t, err)

	return id
}
