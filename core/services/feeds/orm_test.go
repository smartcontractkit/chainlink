package feeds_test

import (
	"database/sql"
	"testing"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	configtest "github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/feeds"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/testdata/testspecs"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
	"github.com/smartcontractkit/chainlink/v2/core/utils/crypto"
)

var (
	uri       = "http://192.168.0.1"
	name      = "Chainlink FMS"
	publicKey = crypto.PublicKey([]byte("11111111111111111111111111111111"))
)

type TestORM struct {
	feeds.ORM

	db *sqlx.DB
}

func setupORM(t *testing.T) *TestORM {
	t.Helper()

	var (
		db   = pgtest.NewSqlxDB(t)
		lggr = logger.TestLogger(t)
		orm  = feeds.NewORM(db, lggr, pgtest.NewQConfig(true))
	)

	return &TestORM{ORM: orm, db: db}
}

// Managers

func Test_ORM_CreateManager(t *testing.T) {
	t.Parallel()

	var (
		orm = setupORM(t)
		mgr = &feeds.FeedsManager{
			URI:       uri,
			Name:      name,
			PublicKey: publicKey,
		}
	)

	count, err := orm.CountManagers()
	require.NoError(t, err)
	require.Equal(t, int64(0), count)

	id, err := orm.CreateManager(mgr)
	require.NoError(t, err)

	count, err = orm.CountManagers()
	require.NoError(t, err)
	require.Equal(t, int64(1), count)

	assert.NotZero(t, id)
}

func Test_ORM_GetManager(t *testing.T) {
	t.Parallel()

	var (
		orm = setupORM(t)
		mgr = &feeds.FeedsManager{
			URI:       uri,
			Name:      name,
			PublicKey: publicKey,
		}
	)

	id, err := orm.CreateManager(mgr)
	require.NoError(t, err)

	actual, err := orm.GetManager(id)
	require.NoError(t, err)

	assert.Equal(t, id, actual.ID)
	assert.Equal(t, uri, actual.URI)
	assert.Equal(t, name, actual.Name)
	assert.Equal(t, publicKey, actual.PublicKey)

	_, err = orm.GetManager(-1)
	require.Error(t, err)
}

func Test_ORM_ListManagers(t *testing.T) {
	t.Parallel()

	var (
		orm = setupORM(t)
		mgr = &feeds.FeedsManager{
			URI:       uri,
			Name:      name,
			PublicKey: publicKey,
		}
	)

	id, err := orm.CreateManager(mgr)
	require.NoError(t, err)

	mgrs, err := orm.ListManagers()
	require.NoError(t, err)
	require.Len(t, mgrs, 1)

	actual := mgrs[0]
	assert.Equal(t, id, actual.ID)
	assert.Equal(t, uri, actual.URI)
	assert.Equal(t, name, actual.Name)
	assert.Equal(t, publicKey, actual.PublicKey)
}

func Test_ORM_ListManagersByIDs(t *testing.T) {
	t.Parallel()

	var (
		orm = setupORM(t)
		mgr = &feeds.FeedsManager{
			URI:       uri,
			Name:      name,
			PublicKey: publicKey,
		}
	)

	id, err := orm.CreateManager(mgr)
	require.NoError(t, err)

	mgrs, err := orm.ListManagersByIDs([]int64{id})
	require.NoError(t, err)
	require.Equal(t, 1, len(mgrs))

	actual := mgrs[0]
	assert.Equal(t, id, actual.ID)
	assert.Equal(t, uri, actual.URI)
	assert.Equal(t, name, actual.Name)
	assert.Equal(t, publicKey, actual.PublicKey)
}

func Test_ORM_UpdateManager(t *testing.T) {
	t.Parallel()

	var (
		orm = setupORM(t)
		mgr = &feeds.FeedsManager{
			URI:       uri,
			Name:      name,
			PublicKey: publicKey,
		}
	)

	id, err := orm.CreateManager(mgr)
	require.NoError(t, err)

	updatedMgr := feeds.FeedsManager{
		ID:        id,
		URI:       "127.0.0.1",
		Name:      "New Name",
		PublicKey: crypto.PublicKey([]byte("22222222222222222222222222222222")),
	}

	err = orm.UpdateManager(updatedMgr)
	require.NoError(t, err)

	actual, err := orm.GetManager(id)
	require.NoError(t, err)

	assert.Equal(t, updatedMgr.URI, actual.URI)
	assert.Equal(t, updatedMgr.Name, actual.Name)
	assert.Equal(t, updatedMgr.PublicKey, actual.PublicKey)
}

// Chain Config

func Test_ORM_CreateChainConfig(t *testing.T) {
	t.Parallel()

	var (
		orm  = setupORM(t)
		fmID = createFeedsManager(t, orm)
		cfg1 = feeds.ChainConfig{
			FeedsManagerID: fmID,
			ChainID:        "1",
			ChainType:      feeds.ChainTypeEVM,
			AccountAddress: "0x0001",
			AdminAddress:   "0x1001",
			FluxMonitorConfig: feeds.FluxMonitorConfig{
				Enabled: true,
			},
			OCR1Config: feeds.OCR1Config{
				Enabled:     true,
				IsBootstrap: false,
				P2PPeerID:   null.StringFrom("p2pkey"),
				KeyBundleID: null.StringFrom("ocrkey"),
			},
			OCR2Config: feeds.OCR2Config{
				Enabled:     true,
				IsBootstrap: true,
				Multiaddr:   null.StringFrom("dns/4"),
			},
		}
	)

	id, err := orm.CreateChainConfig(cfg1)
	require.NoError(t, err)

	actual, err := orm.GetChainConfig(id)
	require.NoError(t, err)

	assertChainConfigEqual(t, map[string]interface{}{
		"feedsManagerID":    cfg1.FeedsManagerID,
		"chainID":           cfg1.ChainID,
		"chainType":         cfg1.ChainType,
		"accountAddress":    cfg1.AccountAddress,
		"adminAddress":      cfg1.AdminAddress,
		"fluxMonitorConfig": cfg1.FluxMonitorConfig,
		"ocrConfig":         cfg1.OCR1Config,
		"ocr2Config":        cfg1.OCR2Config,
	}, *actual)
}

func Test_ORM_CreateBatchChainConfig(t *testing.T) {
	t.Parallel()

	var (
		orm  = setupORM(t)
		fmID = createFeedsManager(t, orm)
		cfg1 = feeds.ChainConfig{
			FeedsManagerID: fmID,
			ChainID:        "1",
			ChainType:      feeds.ChainTypeEVM,
			AccountAddress: "0x0001",
			AdminAddress:   "0x1001",
		}
		cfg2 = feeds.ChainConfig{
			FeedsManagerID: fmID,
			ChainID:        "42",
			ChainType:      "EVM",
			AccountAddress: "0x0002",
			AdminAddress:   "0x2002",
		}
	)

	ids, err := orm.CreateBatchChainConfig([]feeds.ChainConfig{cfg1, cfg2})
	require.NoError(t, err)

	assert.Len(t, ids, 2)

	actual, err := orm.GetChainConfig(ids[0])
	require.NoError(t, err)

	assertChainConfigEqual(t, map[string]interface{}{
		"feedsManagerID":    cfg1.FeedsManagerID,
		"chainID":           cfg1.ChainID,
		"chainType":         cfg1.ChainType,
		"accountAddress":    cfg1.AccountAddress,
		"adminAddress":      cfg1.AdminAddress,
		"fluxMonitorConfig": cfg1.FluxMonitorConfig,
		"ocrConfig":         cfg1.OCR1Config,
		"ocr2Config":        cfg1.OCR2Config,
	}, *actual)

	actual, err = orm.GetChainConfig(ids[1])
	require.NoError(t, err)

	assertChainConfigEqual(t, map[string]interface{}{
		"feedsManagerID":    cfg2.FeedsManagerID,
		"chainID":           cfg2.ChainID,
		"chainType":         cfg2.ChainType,
		"accountAddress":    cfg2.AccountAddress,
		"adminAddress":      cfg2.AdminAddress,
		"fluxMonitorConfig": cfg1.FluxMonitorConfig,
		"ocrConfig":         cfg1.OCR1Config,
		"ocr2Config":        cfg1.OCR2Config,
	}, *actual)

	// Test empty configs
	ids, err = orm.CreateBatchChainConfig([]feeds.ChainConfig{})
	require.NoError(t, err)
	require.Empty(t, ids)
}

func Test_ORM_DeleteChainConfig(t *testing.T) {
	t.Parallel()

	var (
		orm  = setupORM(t)
		fmID = createFeedsManager(t, orm)
		cfg1 = feeds.ChainConfig{
			FeedsManagerID: fmID,
			ChainID:        "1",
			ChainType:      feeds.ChainTypeEVM,
			AccountAddress: "0x0001",
			AdminAddress:   "0x1001",
		}
	)

	id, err := orm.CreateChainConfig(cfg1)
	require.NoError(t, err)

	_, err = orm.GetChainConfig(id)
	require.NoError(t, err)

	actual, err := orm.DeleteChainConfig(id)
	require.NoError(t, err)
	require.Equal(t, id, actual)

	_, err = orm.GetChainConfig(id)
	require.Error(t, err)
}

func Test_ORM_ListChainConfigsByManagerIDs(t *testing.T) {
	t.Parallel()

	var (
		orm  = setupORM(t)
		fmID = createFeedsManager(t, orm)
		cfg1 = feeds.ChainConfig{
			FeedsManagerID: fmID,
			ChainID:        "1",
			ChainType:      feeds.ChainTypeEVM,
			AccountAddress: "0x0001",
			AdminAddress:   "0x1001",
			FluxMonitorConfig: feeds.FluxMonitorConfig{
				Enabled: true,
			},
			OCR1Config: feeds.OCR1Config{
				Enabled:     true,
				IsBootstrap: false,
				P2PPeerID:   null.StringFrom("p2pkey"),
				KeyBundleID: null.StringFrom("ocrkey"),
			},
			OCR2Config: feeds.OCR2Config{
				Enabled:     true,
				IsBootstrap: true,
				Multiaddr:   null.StringFrom("dns/4"),
			},
		}
	)

	_, err := orm.CreateChainConfig(cfg1)
	require.NoError(t, err)

	actual, err := orm.ListChainConfigsByManagerIDs([]int64{fmID})
	require.NoError(t, err)
	require.Len(t, actual, 1)

	assertChainConfigEqual(t, map[string]interface{}{
		"feedsManagerID":    cfg1.FeedsManagerID,
		"chainID":           cfg1.ChainID,
		"chainType":         cfg1.ChainType,
		"accountAddress":    cfg1.AccountAddress,
		"adminAddress":      cfg1.AdminAddress,
		"fluxMonitorConfig": cfg1.FluxMonitorConfig,
		"ocrConfig":         cfg1.OCR1Config,
		"ocr2Config":        cfg1.OCR2Config,
	}, actual[0])
}

func Test_ORM_UpdateChainConfig(t *testing.T) {
	t.Parallel()

	var (
		orm  = setupORM(t)
		fmID = createFeedsManager(t, orm)
		cfg1 = feeds.ChainConfig{
			FeedsManagerID:    fmID,
			ChainID:           "1",
			ChainType:         feeds.ChainTypeEVM,
			AccountAddress:    "0x0001",
			AdminAddress:      "0x1001",
			FluxMonitorConfig: feeds.FluxMonitorConfig{Enabled: false},
			OCR1Config:        feeds.OCR1Config{Enabled: false},
			OCR2Config:        feeds.OCR2Config{Enabled: false},
		}
		updateCfg = feeds.ChainConfig{
			AccountAddress:    "0x0002",
			AdminAddress:      "0x1002",
			FluxMonitorConfig: feeds.FluxMonitorConfig{Enabled: true},
			OCR1Config: feeds.OCR1Config{
				Enabled:     true,
				IsBootstrap: false,
				P2PPeerID:   null.StringFrom("p2pkey"),
				KeyBundleID: null.StringFrom("ocrkey"),
			},
			OCR2Config: feeds.OCR2Config{
				Enabled:     true,
				IsBootstrap: true,
				Multiaddr:   null.StringFrom("dns/4"),
			},
		}
	)

	id, err := orm.CreateChainConfig(cfg1)
	require.NoError(t, err)

	updateCfg.ID = id

	id, err = orm.UpdateChainConfig(updateCfg)
	require.NoError(t, err)

	actual, err := orm.GetChainConfig(id)
	require.NoError(t, err)

	assertChainConfigEqual(t, map[string]interface{}{
		"feedsManagerID":    cfg1.FeedsManagerID,
		"chainID":           cfg1.ChainID,
		"chainType":         cfg1.ChainType,
		"accountAddress":    updateCfg.AccountAddress,
		"adminAddress":      updateCfg.AdminAddress,
		"fluxMonitorConfig": updateCfg.FluxMonitorConfig,
		"ocrConfig":         updateCfg.OCR1Config,
		"ocr2Config":        updateCfg.OCR2Config,
	}, *actual)
}

// Job Proposals

func Test_ORM_CreateJobProposal(t *testing.T) {
	t.Parallel()

	orm := setupORM(t)
	fmID := createFeedsManager(t, orm)

	jp := &feeds.JobProposal{
		Name:           null.StringFrom("jp1"),
		RemoteUUID:     uuid.New(),
		Status:         feeds.JobProposalStatusPending,
		FeedsManagerID: fmID,
	}

	count, err := orm.CountJobProposals()
	require.NoError(t, err)
	require.Equal(t, int64(0), count)

	id, err := orm.CreateJobProposal(jp)
	require.NoError(t, err)

	actual, err := orm.GetJobProposal(id)
	require.NoError(t, err)
	require.Equal(t, jp.Name, actual.Name)
	require.Equal(t, jp.RemoteUUID, actual.RemoteUUID)
	require.Equal(t, jp.Status, actual.Status)
	require.Equal(t, jp.FeedsManagerID, actual.FeedsManagerID)
	require.False(t, actual.PendingUpdate)
	require.NotEmpty(t, actual.CreatedAt)
	require.Equal(t, actual.CreatedAt.String(), actual.UpdatedAt.String())

	assert.NotZero(t, id)
}

func Test_ORM_GetJobProposal(t *testing.T) {
	t.Parallel()

	orm := setupORM(t)
	fmID := createFeedsManager(t, orm)
	remoteUUID := uuid.New()
	deletedUUID := uuid.New()
	name := null.StringFrom("jp1")

	jp := &feeds.JobProposal{
		Name:           name,
		RemoteUUID:     remoteUUID,
		Status:         feeds.JobProposalStatusPending,
		FeedsManagerID: fmID,
	}

	deletedJp := &feeds.JobProposal{
		Name:           name,
		RemoteUUID:     deletedUUID,
		Status:         feeds.JobProposalStatusDeleted,
		FeedsManagerID: fmID,
	}

	id, err := orm.CreateJobProposal(jp)
	require.NoError(t, err)

	_, err = orm.CreateJobProposal(deletedJp)
	require.NoError(t, err)

	assertJobEquals := func(actual *feeds.JobProposal) {
		assert.Equal(t, id, actual.ID)
		assert.Equal(t, name, actual.Name)
		assert.Equal(t, remoteUUID, actual.RemoteUUID)
		assert.Equal(t, jp.Status, actual.Status)
		assert.False(t, actual.ExternalJobID.Valid)
		assert.False(t, actual.PendingUpdate)
		assert.Equal(t, jp.FeedsManagerID, actual.FeedsManagerID)
	}

	t.Run("by id", func(t *testing.T) {
		actual, err := orm.GetJobProposal(id)
		require.NoError(t, err)

		assert.Equal(t, id, actual.ID)
		assertJobEquals(actual)

		_, err = orm.GetJobProposal(int64(0))
		require.Error(t, err)
	})

	t.Run("by remote uuid", func(t *testing.T) {
		actual, err := orm.GetJobProposalByRemoteUUID(remoteUUID)
		require.NoError(t, err)

		assertJobEquals(actual)

		_, err = orm.GetJobProposalByRemoteUUID(deletedUUID)
		require.Error(t, err)

		_, err = orm.GetJobProposalByRemoteUUID(uuid.New())
		require.Error(t, err)
	})
}

func Test_ORM_ListJobProposals(t *testing.T) {
	t.Parallel()

	orm := setupORM(t)
	fmID := createFeedsManager(t, orm)
	uuid := uuid.New()
	name := null.StringFrom("jp1")

	jp := &feeds.JobProposal{
		Name:           name,
		RemoteUUID:     uuid,
		Status:         feeds.JobProposalStatusPending,
		FeedsManagerID: fmID,
	}

	id, err := orm.CreateJobProposal(jp)
	require.NoError(t, err)

	jps, err := orm.ListJobProposals()
	require.NoError(t, err)
	require.Len(t, jps, 1)

	actual := jps[0]
	assert.Equal(t, id, actual.ID)
	assert.Equal(t, name, actual.Name)
	assert.Equal(t, uuid, actual.RemoteUUID)
	assert.Equal(t, jp.Status, actual.Status)
	assert.False(t, actual.ExternalJobID.Valid)
	assert.False(t, actual.PendingUpdate)
	assert.Equal(t, jp.FeedsManagerID, actual.FeedsManagerID)
}

func Test_ORM_CountJobProposalsByStatus(t *testing.T) {
	t.Parallel()

	var (
		orm  = setupORM(t)
		fmID = createFeedsManager(t, orm)

		// Set initial values for job proposal counts
		wantApproved, wantRejected int64
		wantPending, wantCancelled = int64(1), int64(1)
	)

	// Create a pending job proposal.
	_, err := orm.CreateJobProposal(&feeds.JobProposal{
		RemoteUUID:     uuid.New(),
		Status:         feeds.JobProposalStatusPending,
		FeedsManagerID: fmID,
	})
	require.NoError(t, err)

	// Create a cancelled job proposal.
	cancelledUUID := uuid.New()
	_, err = orm.CreateJobProposal(&feeds.JobProposal{
		RemoteUUID:     cancelledUUID,
		Status:         feeds.JobProposalStatusCancelled,
		FeedsManagerID: fmID,
	})
	require.NoError(t, err)

	// Get the initial count and assert against expected values.
	counts, err := orm.CountJobProposalsByStatus()
	require.NoError(t, err)

	assert.Equal(t, wantPending, counts.Pending)
	assert.Equal(t, wantApproved, counts.Approved)
	assert.Equal(t, wantRejected, counts.Rejected)
	assert.Equal(t, wantCancelled, counts.Cancelled)

	// Upsert the cancelled job proposal to rejected
	// which changes pending_update to TRUE, but leaves status as
	// cancelled.
	id, err := orm.UpsertJobProposal(&feeds.JobProposal{
		RemoteUUID:     cancelledUUID,
		Status:         feeds.JobProposalStatusRejected,
		FeedsManagerID: fmID,
	})
	require.NoError(t, err)

	// Assert that the upserted job proposal is now pending update.
	jp, err := orm.GetJobProposal(id)
	require.NoError(t, err)
	assert.Equal(t, true, jp.PendingUpdate)

	// Get final counts of job proposals and make assertions.
	counts, err = orm.CountJobProposalsByStatus()
	require.NoError(t, err)

	wantPending = 2 // One pending + one pending update
	wantCancelled = 0
	assert.Equal(t, wantPending, counts.Pending)
	assert.Equal(t, wantApproved, counts.Approved)
	assert.Equal(t, wantRejected, counts.Rejected)
	assert.Equal(t, wantCancelled, counts.Cancelled)
}

func Test_ORM_ListJobProposalByManagersIDs(t *testing.T) {
	t.Parallel()

	orm := setupORM(t)
	fmID := createFeedsManager(t, orm)
	uuid := uuid.New()
	name := null.StringFrom("jp1")

	jp := &feeds.JobProposal{
		Name:           name,
		RemoteUUID:     uuid,
		Status:         feeds.JobProposalStatusPending,
		FeedsManagerID: fmID,
	}

	id, err := orm.CreateJobProposal(jp)
	require.NoError(t, err)

	jps, err := orm.ListJobProposalsByManagersIDs([]int64{fmID})
	require.NoError(t, err)
	require.Len(t, jps, 1)

	actual := jps[0]
	assert.Equal(t, id, actual.ID)
	assert.Equal(t, name, actual.Name)
	assert.Equal(t, uuid, actual.RemoteUUID)
	assert.Equal(t, jp.Status, actual.Status)
	assert.False(t, actual.ExternalJobID.Valid)
	assert.False(t, actual.PendingUpdate)
	assert.Equal(t, jp.FeedsManagerID, actual.FeedsManagerID)
}

func Test_ORM_UpdateJobProposalStatus(t *testing.T) {
	t.Parallel()

	orm := setupORM(t)
	fmID := createFeedsManager(t, orm)
	jpID := createJobProposal(t, orm, feeds.JobProposalStatusPending, fmID)

	actualCreated, err := orm.GetJobProposal(jpID)
	require.NoError(t, err)

	err = orm.UpdateJobProposalStatus(jpID, feeds.JobProposalStatusRejected)
	require.NoError(t, err)

	actual, err := orm.GetJobProposal(jpID)
	require.NoError(t, err)

	assert.Equal(t, jpID, actual.ID)
	assert.Equal(t, feeds.JobProposalStatusRejected, actual.Status)
	assert.Equal(t, actualCreated.CreatedAt, actual.CreatedAt)
}

func Test_ORM_UpsertJobProposal(t *testing.T) {
	t.Parallel()

	var (
		orm           = setupORM(t)
		fmID          = createFeedsManager(t, orm)
		name          = null.StringFrom("jp1")
		externalJobID = uuid.NullUUID{UUID: uuid.New(), Valid: true}
	)

	jp := &feeds.JobProposal{
		Name:           name,
		RemoteUUID:     uuid.New(),
		Status:         feeds.JobProposalStatusPending,
		FeedsManagerID: fmID,
	}

	// The constraint chk_job_proposals_status_fsm ensures that approved job proposals must have an
	// externalJobID, deleted job proposals are ignored from the check, and all other statuses
	// should have a null externalJobID. We should test the transition between the statuses, moving
	// from pending to approved, and then approved to pending, and pending to deleted and so forth.

	// Create
	count, err := orm.CountJobProposals()
	require.NoError(t, err)
	require.Equal(t, int64(0), count)

	jpID, err := orm.UpsertJobProposal(jp)
	require.NoError(t, err)

	createdActual, err := orm.GetJobProposal(jpID)
	require.NoError(t, err)

	assert.False(t, createdActual.PendingUpdate)

	count, err = orm.CountJobProposals()
	require.NoError(t, err)
	require.Equal(t, int64(1), count)

	assert.NotZero(t, jpID)

	// Update
	jp.Multiaddrs = pq.StringArray{"dns/example.com"}
	jp.Name = null.StringFrom("jp1_updated")

	jpID, err = orm.UpsertJobProposal(jp)
	require.NoError(t, err)

	actual, err := orm.GetJobProposal(jpID)
	require.NoError(t, err)
	assert.Equal(t, jp.Name, actual.Name)
	assert.Equal(t, jp.Status, actual.Status)
	assert.Equal(t, jp.Multiaddrs, actual.Multiaddrs)

	// Ensure there is a difference in the created proposal and the upserted
	// proposal
	assert.NotEqual(t, createdActual.Multiaddrs, actual.Multiaddrs)
	assert.Equal(t, createdActual.CreatedAt, actual.CreatedAt) // CreatedAt does not change
	assert.True(t, actual.PendingUpdate)

	// Approve
	specID := createJobSpec(t, orm, int64(jpID))

	// Defer the FK requirement of an existing job for a job proposal.
	require.NoError(t, utils.JustError(orm.db.Exec(
		`SET CONSTRAINTS job_proposals_job_id_fkey DEFERRED`,
	)))

	err = orm.ApproveSpec(specID, externalJobID.UUID)
	require.NoError(t, err)

	actual, err = orm.GetJobProposal(jpID)
	require.NoError(t, err)

	// Assert that the job proposal is now approved.
	assert.Equal(t, feeds.JobProposalStatusApproved, actual.Status)
	assert.Equal(t, externalJobID, actual.ExternalJobID)

	// Update the proposal again
	jp.Multiaddrs = pq.StringArray{"dns/example1.com"}
	jp.Name = null.StringFrom("jp1_updated_again")
	jp.Status = feeds.JobProposalStatusPending

	_, err = orm.UpsertJobProposal(jp)
	require.NoError(t, err)

	actual, err = orm.GetJobProposal(jpID)
	require.NoError(t, err)

	assert.Equal(t, feeds.JobProposalStatusPending, actual.Status)
	assert.Equal(t, uuid.NullUUID{}, actual.ExternalJobID)
	assert.True(t, actual.PendingUpdate)

	// Delete the proposal
	err = orm.DeleteProposal(jpID)
	require.NoError(t, err)

	actual, err = orm.GetJobProposal(jpID)
	require.NoError(t, err)

	assert.Equal(t, feeds.JobProposalStatusDeleted, actual.Status)

	// Update deleted proposal
	jp.Status = feeds.JobProposalStatusRejected

	jpID, err = orm.UpsertJobProposal(jp)
	require.NoError(t, err)

	// Ensure the deleted proposal does not get updated
	actual, err = orm.GetJobProposal(jpID)
	require.NoError(t, err)
	assert.NotEqual(t, jp.Status, actual.Status)
	assert.Equal(t, feeds.JobProposalStatusDeleted, actual.Status)
}

// Job Proposal Specs

func Test_ORM_ApproveSpec(t *testing.T) {
	t.Parallel()

	var (
		orm           = setupORM(t)
		fmID          = createFeedsManager(t, orm)
		externalJobID = uuid.NullUUID{UUID: uuid.New(), Valid: true}
	)

	// Manually create the job proposal to set pending update
	jpID, err := orm.CreateJobProposal(&feeds.JobProposal{
		RemoteUUID:     uuid.New(),
		Status:         feeds.JobProposalStatusPending,
		FeedsManagerID: fmID,
		PendingUpdate:  true,
	})
	require.NoError(t, err)
	specID := createJobSpec(t, orm, int64(jpID))

	// Defer the FK requirement of an existing job for a job proposal.
	require.NoError(t, utils.JustError(orm.db.Exec(
		`SET CONSTRAINTS job_proposals_job_id_fkey DEFERRED`,
	)))

	err = orm.ApproveSpec(specID, externalJobID.UUID)
	require.NoError(t, err)

	actual, err := orm.GetSpec(specID)
	require.NoError(t, err)

	assert.Equal(t, specID, actual.ID)
	assert.Equal(t, feeds.SpecStatusApproved, actual.Status)

	actualJP, err := orm.GetJobProposal(jpID)
	require.NoError(t, err)

	assert.Equal(t, externalJobID, actualJP.ExternalJobID)
	assert.Equal(t, feeds.JobProposalStatusApproved, actualJP.Status)
	assert.False(t, actualJP.PendingUpdate)
}

func Test_ORM_CancelSpec(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name               string
		before             func(orm *TestORM) (int64, int64)
		wantSpecStatus     feeds.SpecStatus
		wantProposalStatus feeds.JobProposalStatus
		wantErr            string
	}{
		{
			name: "pending proposal",
			before: func(orm *TestORM) (int64, int64) {
				fmID := createFeedsManager(t, orm)
				jpID := createJobProposal(t, orm, feeds.JobProposalStatusPending, fmID)
				specID := createJobSpec(t, orm, int64(jpID))

				return jpID, specID
			},
			wantSpecStatus:     feeds.SpecStatusCancelled,
			wantProposalStatus: feeds.JobProposalStatusCancelled,
		},
		{
			name: "deleted proposal",
			before: func(orm *TestORM) (int64, int64) {
				fmID := createFeedsManager(t, orm)
				jpID := createJobProposal(t, orm, feeds.JobProposalStatusDeleted, fmID)
				specID := createJobSpec(t, orm, int64(jpID))

				return jpID, specID
			},
			wantSpecStatus:     feeds.SpecStatusCancelled,
			wantProposalStatus: feeds.JobProposalStatusDeleted,
		},
		{
			name: "not found",
			before: func(orm *TestORM) (int64, int64) {
				return 0, 0
			},
			wantErr: "sql: no rows in result set",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			orm := setupORM(t)

			jpID, specID := tc.before(orm)

			err := orm.CancelSpec(specID)

			if tc.wantErr != "" {
				require.EqualError(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)

				actual, err := orm.GetSpec(specID)
				require.NoError(t, err)

				assert.Equal(t, specID, actual.ID)
				assert.Equal(t, tc.wantSpecStatus, actual.Status)

				actualJP, err := orm.GetJobProposal(jpID)
				require.NoError(t, err)

				assert.Equal(t, tc.wantProposalStatus, actualJP.Status)
				assert.False(t, actualJP.PendingUpdate)
			}
		})
	}
}

func Test_ORM_DeleteProposal(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name                      string
		before                    func(orm *TestORM) int64
		wantProposalStatus        feeds.JobProposalStatus
		wantProposalPendingUpdate bool
		wantErr                   string
	}{
		{
			name: "pending proposal",
			before: func(orm *TestORM) int64 {
				fmID := createFeedsManager(t, orm)
				jpID := createJobProposal(t, orm, feeds.JobProposalStatusPending, fmID)

				return jpID
			},
			wantProposalPendingUpdate: false,
			wantProposalStatus:        feeds.JobProposalStatusDeleted,
		},
		{
			name: "approved proposal",
			before: func(orm *TestORM) int64 {
				fmID := createFeedsManager(t, orm)
				jpID := createJobProposal(t, orm, feeds.JobProposalStatusPending, fmID)
				specID := createJobSpec(t, orm, int64(jpID))

				externalJobID := uuid.NullUUID{UUID: uuid.New(), Valid: true}

				// Defer the FK requirement of an existing job for a job proposal.
				require.NoError(t, utils.JustError(orm.db.Exec(
					`SET CONSTRAINTS job_proposals_job_id_fkey DEFERRED`,
				)))

				err := orm.ApproveSpec(specID, externalJobID.UUID)
				require.NoError(t, err)

				return jpID
			},
			wantProposalPendingUpdate: true,
			wantProposalStatus:        feeds.JobProposalStatusDeleted,
		},
		{
			name: "cancelled proposal",
			before: func(orm *TestORM) int64 {
				fmID := createFeedsManager(t, orm)
				jpID := createJobProposal(t, orm, feeds.JobProposalStatusCancelled, fmID)

				return jpID
			},
			wantProposalPendingUpdate: false,
			wantProposalStatus:        feeds.JobProposalStatusDeleted,
		},
		{
			name: "rejected proposal",
			before: func(orm *TestORM) int64 {
				fmID := createFeedsManager(t, orm)
				jpID := createJobProposal(t, orm, feeds.JobProposalStatusRejected, fmID)

				return jpID
			},
			wantProposalPendingUpdate: false,
			wantProposalStatus:        feeds.JobProposalStatusDeleted,
		},
		{
			name: "not found",
			before: func(orm *TestORM) int64 {
				return 0
			},
			wantErr: "sql: no rows in result set",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			orm := setupORM(t)

			jpID := tc.before(orm)

			err := orm.DeleteProposal(jpID)

			if tc.wantErr != "" {
				require.EqualError(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)

				actual, err := orm.GetJobProposal(jpID)
				require.NoError(t, err)

				assert.Equal(t, jpID, actual.ID)
				assert.Equal(t, tc.wantProposalStatus, actual.Status)
				assert.Equal(t, tc.wantProposalPendingUpdate, actual.PendingUpdate)
			}
		})
	}
}

func Test_ORM_RevokeSpec(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name               string
		before             func(orm *TestORM) (int64, int64)
		wantProposalStatus feeds.JobProposalStatus
		wantSpecStatus     feeds.SpecStatus
		wantErr            string
	}{
		{
			name: "pending proposal",
			before: func(orm *TestORM) (int64, int64) {
				fmID := createFeedsManager(t, orm)
				jpID := createJobProposal(t, orm, feeds.JobProposalStatusPending, fmID)
				specID := createJobSpec(t, orm, int64(jpID))

				return jpID, specID
			},
			wantProposalStatus: feeds.JobProposalStatusRevoked,
			wantSpecStatus:     feeds.SpecStatusRevoked,
		},
		{
			name: "approved proposal",
			before: func(orm *TestORM) (int64, int64) {
				fmID := createFeedsManager(t, orm)
				jpID := createJobProposal(t, orm, feeds.JobProposalStatusPending, fmID)
				specID := createJobSpec(t, orm, int64(jpID))

				externalJobID := uuid.NullUUID{UUID: uuid.New(), Valid: true}

				// Defer the FK requirement of an existing job for a job proposal.
				require.NoError(t, utils.JustError(orm.db.Exec(
					`SET CONSTRAINTS job_proposals_job_id_fkey DEFERRED`,
				)))

				err := orm.ApproveSpec(specID, externalJobID.UUID)
				require.NoError(t, err)

				return jpID, specID
			},
			wantProposalStatus: feeds.JobProposalStatusApproved,
			wantSpecStatus:     feeds.SpecStatusApproved,
		},
		{
			name: "cancelled proposal",
			before: func(orm *TestORM) (int64, int64) {
				fmID := createFeedsManager(t, orm)
				jpID := createJobProposal(t, orm, feeds.JobProposalStatusCancelled, fmID)
				specID := createJobSpec(t, orm, int64(jpID))

				return jpID, specID
			},
			wantProposalStatus: feeds.JobProposalStatusRevoked,
			wantSpecStatus:     feeds.SpecStatusRevoked,
		},
		{
			name: "rejected proposal",
			before: func(orm *TestORM) (int64, int64) {
				fmID := createFeedsManager(t, orm)
				jpID := createJobProposal(t, orm, feeds.JobProposalStatusRejected, fmID)
				specID := createJobSpec(t, orm, int64(jpID))

				return jpID, specID
			},
			wantProposalStatus: feeds.JobProposalStatusRevoked,
			wantSpecStatus:     feeds.SpecStatusRevoked,
		},
		{
			name: "deleted proposal",
			before: func(orm *TestORM) (int64, int64) {
				fmID := createFeedsManager(t, orm)
				jpID := createJobProposal(t, orm, feeds.JobProposalStatusDeleted, fmID)
				specID := createJobSpec(t, orm, int64(jpID))

				return jpID, specID
			},
			wantProposalStatus: feeds.JobProposalStatusDeleted,
			wantSpecStatus:     feeds.SpecStatusRevoked,
		},
		{
			name: "not found",
			before: func(orm *TestORM) (int64, int64) {
				return 0, 0
			},
			wantErr: "sql: no rows in result set",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			orm := setupORM(t)

			jpID, specID := tc.before(orm)

			err := orm.RevokeSpec(specID)

			if tc.wantErr != "" {
				require.EqualError(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)

				actualJP, err := orm.GetJobProposal(jpID)
				require.NoError(t, err)

				assert.Equal(t, tc.wantProposalStatus, actualJP.Status)
				assert.False(t, actualJP.PendingUpdate)

				assert.Equal(t, jpID, actualJP.ID)
				assert.Equal(t, tc.wantProposalStatus, actualJP.Status)
			}
		})
	}
}

func Test_ORM_ExistsSpecByJobProposalIDAndVersion(t *testing.T) {
	t.Parallel()

	var (
		orm  = setupORM(t)
		fmID = createFeedsManager(t, orm)
		jpID = createJobProposal(t, orm, feeds.JobProposalStatusPending, fmID)
	)

	createJobSpec(t, orm, int64(jpID))

	exists, err := orm.ExistsSpecByJobProposalIDAndVersion(jpID, 1)
	require.NoError(t, err)
	require.True(t, exists)

	exists, err = orm.ExistsSpecByJobProposalIDAndVersion(jpID, 2)
	require.NoError(t, err)
	require.False(t, exists)
}

func Test_ORM_GetSpec(t *testing.T) {
	t.Parallel()

	var (
		orm    = setupORM(t)
		fmID   = createFeedsManager(t, orm)
		jpID   = createJobProposal(t, orm, feeds.JobProposalStatusPending, fmID)
		specID = createJobSpec(t, orm, int64(jpID))
	)

	actual, err := orm.GetSpec(specID)
	require.NoError(t, err)

	assert.Equal(t, "spec data", actual.Definition)
	assert.Equal(t, int32(1), actual.Version)
	assert.Equal(t, feeds.SpecStatusPending, actual.Status)
	assert.Equal(t, jpID, actual.JobProposalID)
}

func Test_ORM_GetApprovedSpec(t *testing.T) {
	t.Parallel()

	var (
		orm           = setupORM(t)
		fmID          = createFeedsManager(t, orm)
		jpID          = createJobProposal(t, orm, feeds.JobProposalStatusPending, fmID)
		specID        = createJobSpec(t, orm, int64(jpID))
		externalJobID = uuid.NullUUID{UUID: uuid.New(), Valid: true}
	)

	// Defer the FK requirement of a job proposal so we don't have to setup a
	// real job.
	require.NoError(t, utils.JustError(orm.db.Exec(
		`SET CONSTRAINTS job_proposals_job_id_fkey DEFERRED`,
	)))

	err := orm.ApproveSpec(specID, externalJobID.UUID)
	require.NoError(t, err)

	actual, err := orm.GetApprovedSpec(jpID)
	require.NoError(t, err)

	assert.Equal(t, specID, actual.ID)
	assert.Equal(t, feeds.SpecStatusApproved, actual.Status)

	err = orm.CancelSpec(specID)
	require.NoError(t, err)

	_, err = orm.GetApprovedSpec(jpID)
	require.Error(t, err)

	assert.ErrorIs(t, err, sql.ErrNoRows)
}

func Test_ORM_GetLatestSpec(t *testing.T) {
	t.Parallel()

	var (
		orm  = setupORM(t)
		fmID = createFeedsManager(t, orm)
		jpID = createJobProposal(t, orm, feeds.JobProposalStatusPending, fmID)
	)

	_ = createJobSpec(t, orm, int64(jpID))
	spec2ID, err := orm.CreateSpec(feeds.JobProposalSpec{
		Definition:    "spec data",
		Version:       2,
		Status:        feeds.SpecStatusPending,
		JobProposalID: jpID,
	})
	require.NoError(t, err)

	actual, err := orm.GetSpec(spec2ID)
	require.NoError(t, err)

	assert.Equal(t, spec2ID, actual.ID)
	assert.Equal(t, "spec data", actual.Definition)
	assert.Equal(t, int32(2), actual.Version)
	assert.Equal(t, feeds.SpecStatusPending, actual.Status)
	assert.Equal(t, jpID, actual.JobProposalID)
}

func Test_ORM_ListSpecsByJobProposalIDs(t *testing.T) {
	t.Parallel()

	var (
		orm  = setupORM(t)
		fmID = createFeedsManager(t, orm)

		jp1ID = createJobProposal(t, orm, feeds.JobProposalStatusPending, fmID)
		jp2ID = createJobProposal(t, orm, feeds.JobProposalStatusPending, fmID)
	)

	// Create the specs for the proposals
	createJobSpec(t, orm, int64(jp1ID))
	createJobSpec(t, orm, int64(jp2ID))

	specs, err := orm.ListSpecsByJobProposalIDs([]int64{jp1ID, jp2ID})
	require.NoError(t, err)
	require.Len(t, specs, 2)

	actual := specs[0]

	assert.Equal(t, "spec data", actual.Definition)
	assert.Equal(t, int32(1), actual.Version)
	assert.Equal(t, feeds.SpecStatusPending, actual.Status)
	assert.Equal(t, jp1ID, actual.JobProposalID)

	actual = specs[1]

	assert.Equal(t, "spec data", actual.Definition)
	assert.Equal(t, int32(1), actual.Version)
	assert.Equal(t, feeds.SpecStatusPending, actual.Status)
	assert.Equal(t, jp2ID, actual.JobProposalID)
}

func Test_ORM_RejectSpec(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name               string
		before             func(orm *TestORM) (int64, int64)
		wantSpecStatus     feeds.SpecStatus
		wantProposalStatus feeds.JobProposalStatus
		wantErr            string
	}{
		{
			name: "pending proposal",
			before: func(orm *TestORM) (int64, int64) {
				fmID := createFeedsManager(t, orm)
				jpID := createJobProposal(t, orm, feeds.JobProposalStatusPending, fmID)
				specID := createJobSpec(t, orm, int64(jpID))

				return jpID, specID
			},
			wantSpecStatus:     feeds.SpecStatusRejected,
			wantProposalStatus: feeds.JobProposalStatusRejected,
		},
		{
			name: "approved proposal",
			before: func(orm *TestORM) (int64, int64) {
				fmID := createFeedsManager(t, orm)
				jpID := createJobProposal(t, orm, feeds.JobProposalStatusPending, fmID)
				specID := createJobSpec(t, orm, int64(jpID))

				externalJobID := uuid.NullUUID{UUID: uuid.New(), Valid: true}

				// Defer the FK requirement of an existing job for a job proposal.
				require.NoError(t, utils.JustError(orm.db.Exec(
					`SET CONSTRAINTS job_proposals_job_id_fkey DEFERRED`,
				)))

				err := orm.ApproveSpec(specID, externalJobID.UUID)
				require.NoError(t, err)

				return jpID, specID
			},
			wantSpecStatus:     feeds.SpecStatusRejected,
			wantProposalStatus: feeds.JobProposalStatusApproved,
		},
		{
			name: "cancelled proposal",
			before: func(orm *TestORM) (int64, int64) {
				fmID := createFeedsManager(t, orm)
				jpID := createJobProposal(t, orm, feeds.JobProposalStatusCancelled, fmID)
				specID := createJobSpec(t, orm, int64(jpID))

				return jpID, specID
			},
			wantSpecStatus:     feeds.SpecStatusRejected,
			wantProposalStatus: feeds.JobProposalStatusRejected,
		},
		{
			name: "deleted proposal",
			before: func(orm *TestORM) (int64, int64) {
				fmID := createFeedsManager(t, orm)
				jpID := createJobProposal(t, orm, feeds.JobProposalStatusDeleted, fmID)
				specID := createJobSpec(t, orm, int64(jpID))

				return jpID, specID
			},
			wantSpecStatus:     feeds.SpecStatusRejected,
			wantProposalStatus: feeds.JobProposalStatusDeleted,
		},
		{
			name: "not found",
			before: func(orm *TestORM) (int64, int64) {
				return 0, 0
			},
			wantErr: "sql: no rows in result set",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			orm := setupORM(t)

			jpID, specID := tc.before(orm)

			err := orm.RejectSpec(specID)

			if tc.wantErr != "" {
				require.EqualError(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)

				actual, err := orm.GetSpec(specID)
				require.NoError(t, err)

				assert.Equal(t, specID, actual.ID)
				assert.Equal(t, tc.wantSpecStatus, actual.Status)

				actualJP, err := orm.GetJobProposal(jpID)
				require.NoError(t, err)

				assert.Equal(t, tc.wantProposalStatus, actualJP.Status)
				assert.False(t, actualJP.PendingUpdate)
			}
		})
	}
}

func Test_ORM_UpdateSpecDefinition(t *testing.T) {
	t.Parallel()

	var (
		orm    = setupORM(t)
		fmID   = createFeedsManager(t, orm)
		jpID   = createJobProposal(t, orm, feeds.JobProposalStatusPending, fmID)
		specID = createJobSpec(t, orm, int64(jpID))
	)

	prev, err := orm.GetSpec(specID)
	require.NoError(t, err)

	err = orm.UpdateSpecDefinition(specID, "updated spec")
	require.NoError(t, err)

	actual, err := orm.GetSpec(specID)
	require.NoError(t, err)

	assert.Equal(t, specID, actual.ID)
	require.NotEqual(t, prev.Definition, actual.Definition)
	require.Equal(t, "updated spec", actual.Definition)

	// Not found
	err = orm.UpdateSpecDefinition(-1, "updated spec")
	require.Error(t, err)
}

// Other

func Test_ORM_IsJobManaged(t *testing.T) {
	t.Parallel()

	var (
		orm           = setupORM(t)
		fmID          = createFeedsManager(t, orm)
		jpID          = createJobProposal(t, orm, feeds.JobProposalStatusPending, fmID)
		specID        = createJobSpec(t, orm, int64(jpID))
		externalJobID = uuid.NullUUID{UUID: uuid.New(), Valid: true}
	)

	j := createJob(t, orm.db, externalJobID.UUID)

	isManaged, err := orm.IsJobManaged(int64(j.ID))
	require.NoError(t, err)
	assert.False(t, isManaged)

	err = orm.ApproveSpec(specID, externalJobID.UUID)
	require.NoError(t, err)

	isManaged, err = orm.IsJobManaged(int64(j.ID))
	require.NoError(t, err)
	assert.True(t, isManaged)
}

// Helpers

func assertChainConfigEqual(t *testing.T, want map[string]interface{}, actual feeds.ChainConfig) {
	t.Helper()

	assert.Equal(t, want["feedsManagerID"], actual.FeedsManagerID)
	assert.Equal(t, want["chainID"], actual.ChainID)
	assert.Equal(t, want["chainType"], actual.ChainType)
	assert.Equal(t, want["accountAddress"], actual.AccountAddress)
	assert.Equal(t, want["adminAddress"], actual.AdminAddress)
	assert.Equal(t, want["fluxMonitorConfig"], actual.FluxMonitorConfig)
	assert.Equal(t, want["ocrConfig"], actual.OCR1Config)
	assert.Equal(t, want["ocr2Config"], actual.OCR2Config)
}

// createFeedsManager is a test helper to create a feeds manager
func createFeedsManager(t *testing.T, orm feeds.ORM) int64 {
	t.Helper()

	mgr := &feeds.FeedsManager{
		URI:       uri,
		Name:      name,
		PublicKey: publicKey,
	}

	id, err := orm.CreateManager(mgr)
	require.NoError(t, err)

	return id
}

func createJob(t *testing.T, db *sqlx.DB, externalJobID uuid.UUID) *job.Job {
	t.Helper()

	var (
		config      = configtest.NewGeneralConfig(t, nil)
		keyStore    = cltest.NewKeyStore(t, db, config)
		lggr        = logger.TestLogger(t)
		pipelineORM = pipeline.NewORM(db, lggr, config)
		bridgeORM   = bridges.NewORM(db, lggr, config)
		cc          = evmtest.NewChainSet(t, evmtest.TestChainOpts{DB: db, GeneralConfig: config, KeyStore: keyStore.Eth()})
		orm         = job.NewORM(db, cc, pipelineORM, bridgeORM, keyStore, lggr, config)
	)

	require.NoError(t, keyStore.OCR().Add(cltest.DefaultOCRKey))
	require.NoError(t, keyStore.P2P().Add(cltest.DefaultP2PKey))

	defer func() { assert.NoError(t, orm.Close()) }()

	_, bridge := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{}, config)
	_, bridge2 := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{}, config)

	_, address := cltest.MustInsertRandomKey(t, keyStore.Eth())
	jb, err := ocr.ValidatedOracleSpecToml(cc,
		testspecs.GenerateOCRSpec(testspecs.OCRSpecParams{
			JobID:              externalJobID.String(),
			TransmitterAddress: address.Hex(),
			DS1BridgeName:      bridge.Name.String(),
			DS2BridgeName:      bridge2.Name.String(),
		}).Toml(),
	)
	require.NoError(t, err)

	err = orm.CreateJob(&jb)
	require.NoError(t, err)

	return &jb
}

func createJobProposal(t *testing.T, orm feeds.ORM, status feeds.JobProposalStatus, fmID int64) int64 {
	t.Helper()

	id, err := orm.CreateJobProposal(&feeds.JobProposal{
		RemoteUUID:     uuid.New(),
		Status:         status,
		FeedsManagerID: fmID,
		PendingUpdate:  true,
	})
	require.NoError(t, err)

	return id
}

func createJobSpec(t *testing.T, orm feeds.ORM, jpID int64) int64 {
	t.Helper()

	id, err := orm.CreateSpec(feeds.JobProposalSpec{
		Definition:    "spec data",
		Version:       1,
		Status:        feeds.SpecStatusPending,
		JobProposalID: jpID,
	})
	require.NoError(t, err)

	return id
}
