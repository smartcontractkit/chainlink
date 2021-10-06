package feeds_test

import (
	"context"
	"database/sql"
	"encoding/hex"
	"math/big"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/chains/evm"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/keystest"
	"github.com/smartcontractkit/chainlink/core/services/feeds"
	"github.com/smartcontractkit/chainlink/core/services/feeds/mocks"
	"github.com/smartcontractkit/chainlink/core/services/feeds/proto"
	"github.com/smartcontractkit/chainlink/core/services/job"
	jobmocks "github.com/smartcontractkit/chainlink/core/services/job/mocks"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/csakey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	ksmocks "github.com/smartcontractkit/chainlink/core/services/keystore/mocks"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	pgmocks "github.com/smartcontractkit/chainlink/core/services/postgres/mocks"
	"github.com/smartcontractkit/chainlink/core/services/versioning"
	verMocks "github.com/smartcontractkit/chainlink/core/services/versioning/mocks"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils/crypto"

	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"
)

const TestSpec = `
type              = "fluxmonitor"
schemaVersion     = 1
name              = "example flux monitor spec"
contractAddress   = "0x3cCad4715152693fE3BC4460591e3D3Fbd071b42"
externalJobID     = "0EEC7E1D-D0D2-476C-A1A8-72DFB6633F47"
threshold = 0.5
absoluteThreshold = 0.0 # optional

idleTimerPeriod = "1s"
idleTimerDisabled = false

pollTimerPeriod = "1m"
pollTimerDisabled = false

observationSource = """
ds1  [type=http method=GET url="https://api.coindesk.com/v1/bpi/currentprice.json"];
jp1  [type=jsonparse path="bpi,USD,rate_float"];
ds1 -> jp1 -> answer1;
answer1 [type=median index=0];
"""
`

type TestService struct {
	feeds.Service
	orm         *mocks.ORM
	jobORM      *jobmocks.ORM
	verORM      *verMocks.ORM
	connMgr     *mocks.ConnectionsManager
	txm         *pgmocks.TransactionManager
	spawner     *jobmocks.Spawner
	fmsClient   *mocks.FeedsManagerClient
	csaKeystore *ksmocks.CSA
	ethKeystore *ksmocks.Eth
	cfg         *mocks.Config
	cc          evm.ChainSet
}

func setupTestService(t *testing.T) *TestService {
	var (
		orm         = &mocks.ORM{}
		jobORM      = &jobmocks.ORM{}
		verORM      = &verMocks.ORM{}
		connMgr     = &mocks.ConnectionsManager{}
		txm         = &pgmocks.TransactionManager{}
		spawner     = &jobmocks.Spawner{}
		fmsClient   = &mocks.FeedsManagerClient{}
		csaKeystore = &ksmocks.CSA{}
		ethKeystore = &ksmocks.Eth{}
		cfg         = &mocks.Config{}
	)

	t.Cleanup(func() {
		mock.AssertExpectationsForObjects(t,
			orm,
			jobORM,
			verORM,
			connMgr,
			txm,
			spawner,
			fmsClient,
			csaKeystore,
			ethKeystore,
			cfg,
		)
	})

	gcfg := configtest.NewTestGeneralConfig(t)
	gcfg.Overrides.EthereumDisabled = null.BoolFrom(true)
	cc := evmtest.NewChainSet(t, evmtest.TestChainOpts{GeneralConfig: gcfg})
	svc := feeds.NewService(orm, jobORM, verORM, txm, spawner, csaKeystore, ethKeystore, cfg, cc)
	svc.SetConnectionsManager(connMgr)

	return &TestService{
		Service:     svc,
		orm:         orm,
		jobORM:      jobORM,
		verORM:      verORM,
		connMgr:     connMgr,
		txm:         txm,
		spawner:     spawner,
		fmsClient:   fmsClient,
		csaKeystore: csaKeystore,
		ethKeystore: ethKeystore,
		cfg:         cfg,
		cc:          cc,
	}
}

func Test_Service_RegisterManager(t *testing.T) {
	t.Parallel()

	key := cltest.DefaultCSAKey

	var (
		id        = int64(1)
		ms        = feeds.FeedsManager{}
		pubKeyHex = "0f17c3bf72de8beef6e2d17a14c0a972f5d7e0e66e70722373f12b88382d40f9"
	)

	var pubKey crypto.PublicKey
	_, err := hex.Decode([]byte(pubKeyHex), pubKey)
	require.NoError(t, err)

	svc := setupTestService(t)

	svc.orm.On("CountManagers", context.Background()).Return(int64(0), nil)
	svc.orm.On("CreateManager", context.Background(), &ms).
		Return(id, nil)
	svc.csaKeystore.On("GetAll").Return([]csakey.KeyV2{key}, nil)
	// ListManagers runs in a goroutine so it might be called.
	svc.orm.On("ListManagers", context.Background()).Return([]feeds.FeedsManager{ms}, nil).Maybe()
	svc.connMgr.On("Connect", mock.IsType(feeds.ConnectOpts{}))

	actual, err := svc.RegisterManager(&ms)
	// We need to stop the service because the manager will attempt to make a
	// connection
	defer svc.Close()
	require.NoError(t, err)

	assert.Equal(t, actual, id)
}

func Test_Service_ListManagers(t *testing.T) {
	t.Parallel()

	var (
		ms  = feeds.FeedsManager{}
		mss = []feeds.FeedsManager{ms}
	)
	svc := setupTestService(t)

	svc.orm.On("ListManagers", context.Background()).
		Return(mss, nil)
	svc.connMgr.On("IsConnected", ms.ID).Return(false)

	actual, err := svc.ListManagers()
	require.NoError(t, err)

	assert.Equal(t, actual, mss)
}

func Test_Service_GetManager(t *testing.T) {
	t.Parallel()

	var (
		id = int64(1)
		ms = feeds.FeedsManager{ID: id}
	)
	svc := setupTestService(t)

	svc.orm.On("GetManager", context.Background(), id).
		Return(&ms, nil)
	svc.connMgr.On("IsConnected", ms.ID).Return(false)

	actual, err := svc.GetManager(id)
	require.NoError(t, err)

	assert.Equal(t, actual, &ms)
}

func Test_Service_CreateJobProposal(t *testing.T) {
	t.Parallel()

	var (
		id = int64(1)
		jp = feeds.JobProposal{
			FeedsManagerID: 1,
			RemoteUUID:     uuid.NewV4(),
			Status:         "pending",
			Spec:           TestSpec,
		}
	)
	svc := setupTestService(t)

	svc.cfg.On("DefaultHTTPTimeout").Return(models.MustMakeDuration(1 * time.Second))
	svc.orm.On("CreateJobProposal", context.Background(), &jp).
		Return(id, nil)

	actual, err := svc.CreateJobProposal(&jp)
	require.NoError(t, err)

	assert.Equal(t, actual, id)
}

func Test_Service_ProposeJob(t *testing.T) {
	t.Parallel()

	var (
		ctx = context.Background()
		id  = int64(1)
		jp  = feeds.JobProposal{
			FeedsManagerID: 1,
			RemoteUUID:     uuid.NewV4(),
			Status:         "pending",
			Spec:           TestSpec,
		}
		httpTimeout = models.MustMakeDuration(1 * time.Second)
	)

	testCases := []struct {
		name     string
		proposal feeds.JobProposal
		before   func(svc *TestService)
		wantID   int64
		wantErr  string
	}{
		{
			name: "Create success",
			before: func(svc *TestService) {
				svc.cfg.On("DefaultHTTPTimeout").Return(httpTimeout)
				svc.orm.On("GetJobProposalByRemoteUUID", ctx, jp.RemoteUUID).Return(nil, sql.ErrNoRows)
				svc.orm.On("UpsertJobProposal", ctx, &jp).Return(id, nil)
			},
			wantID:   id,
			proposal: jp,
		},
		{
			name: "Update success",
			before: func(svc *TestService) {
				svc.cfg.On("DefaultHTTPTimeout").Return(httpTimeout)
				svc.orm.
					On("GetJobProposalByRemoteUUID", ctx, jp.RemoteUUID).
					Return(&feeds.JobProposal{
						FeedsManagerID: jp.FeedsManagerID,
						RemoteUUID:     jp.RemoteUUID,
						Status:         feeds.JobProposalStatusPending,
					}, nil)
				svc.orm.On("UpsertJobProposal", ctx, &jp).Return(id, nil)
			},
			wantID:   id,
			proposal: jp,
		},
		{
			name: "Updates the status of a rejected job proposal",
			before: func(svc *TestService) {
				svc.cfg.On("DefaultHTTPTimeout").Return(httpTimeout)
				svc.orm.
					On("GetJobProposalByRemoteUUID", ctx, jp.RemoteUUID).
					Return(&feeds.JobProposal{
						FeedsManagerID: jp.FeedsManagerID,
						RemoteUUID:     jp.RemoteUUID,
						Status:         feeds.JobProposalStatusRejected,
					}, nil)
				svc.orm.On("UpsertJobProposal", ctx, &jp).Return(id, nil)
			},
			wantID:   id,
			proposal: jp,
		},
		{
			name:     "contains invalid job spec",
			proposal: feeds.JobProposal{Spec: ""},
			wantErr:  "invalid job type",
		},
		{
			name: "must be an ocr job to include bootstraps",
			proposal: feeds.JobProposal{
				RemoteUUID: uuid.NewV4(),
				Status:     "pending",
				Spec:       TestSpec,
				Multiaddrs: pq.StringArray{"/dns4/example.com"},
			},
			before: func(svc *TestService) {
				svc.cfg.On("DefaultHTTPTimeout").Return(httpTimeout)
			},
			wantErr: "only OCR job type supports multiaddr",
		},
		{
			name:     "ensure an upsert validates the job propsal belongs to the feeds manager",
			proposal: jp,
			before: func(svc *TestService) {
				svc.cfg.On("DefaultHTTPTimeout").Return(httpTimeout)
				svc.orm.
					On("GetJobProposalByRemoteUUID", ctx, jp.RemoteUUID).
					Return(&feeds.JobProposal{
						FeedsManagerID: 2,
						RemoteUUID:     jp.RemoteUUID,
						Status:         feeds.JobProposalStatusPending,
					}, nil)
			},
			wantErr: "cannot update a job proposal belonging to another feeds manager",
		},
		{
			name:     "ensure an upsert does not occur on an approved job proposal",
			proposal: jp,
			before: func(svc *TestService) {
				svc.cfg.On("DefaultHTTPTimeout").Return(httpTimeout)
				svc.orm.
					On("GetJobProposalByRemoteUUID", ctx, jp.RemoteUUID).
					Return(&feeds.JobProposal{
						FeedsManagerID: jp.FeedsManagerID,
						RemoteUUID:     jp.RemoteUUID,
						Status:         feeds.JobProposalStatusApproved,
					}, nil)
			},
			wantErr: "cannot repropose a job that has already been approved",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := setupTestService(t)
			if tc.before != nil {
				tc.before(svc)
			}

			actual, err := svc.ProposeJob(&tc.proposal)

			if tc.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantErr)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.wantID, actual)
			}
		})
	}

}

func Test_Service_SyncNodeInfo(t *testing.T) {
	rawKey, err := keystest.NewKey()
	require.NoError(t, err)
	var (
		ctx       = context.Background()
		multiaddr = "/dns4/chain.link/tcp/1234/p2p/16Uiu2HAm58SP7UL8zsnpeuwHfytLocaqgnyaYKP8wu7qRdrixLju"
		feedsMgr  = &feeds.FeedsManager{
			ID:                        1,
			JobTypes:                  pq.StringArray{feeds.JobTypeFluxMonitor},
			IsOCRBootstrapPeer:        true,
			OCRBootstrapPeerMultiaddr: null.StringFrom(multiaddr),
		}
		chainID    = big.NewInt(1)
		sendingKey = ethkey.KeyV2{
			Address: ethkey.EIP55AddressFromAddress(rawKey.Address),
		}
		nodeVersion = &versioning.NodeVersion{
			Version: "1.0.0",
		}
	)

	svc := setupTestService(t)

	// Mock fetching the information to send
	svc.orm.On("GetManager", ctx, feedsMgr.ID).Return(feedsMgr, nil)
	svc.ethKeystore.On("SendingKeys").Return([]ethkey.KeyV2{sendingKey}, nil)
	svc.cfg.On("ChainID").Return(chainID)
	svc.connMgr.On("GetClient", feedsMgr.ID).Return(svc.fmsClient, nil)
	svc.connMgr.On("IsConnected", feedsMgr.ID).Return(false, nil)
	svc.verORM.On("FindLatestNodeVersion").Return(nodeVersion, nil)

	// Mock the send
	svc.fmsClient.On("UpdateNode", ctx, &proto.UpdateNodeRequest{
		JobTypes:           []proto.JobType{proto.JobType_JOB_TYPE_FLUX_MONITOR},
		ChainId:            chainID.Int64(),
		ChainIds:           []int64{chainID.Int64()},
		AccountAddresses:   []string{sendingKey.Address.String()},
		IsBootstrapPeer:    true,
		BootstrapMultiaddr: multiaddr,
		Version:            nodeVersion.Version,
	}).Return(&proto.UpdateNodeResponse{}, nil)

	err = svc.SyncNodeInfo(feedsMgr.ID)
	require.NoError(t, err)
}

func Test_Service_UpdateFeedsManager(t *testing.T) {
	key := cltest.DefaultCSAKey

	var (
		ctx = context.Background()
		mgr = feeds.FeedsManager{
			ID: 1,
		}
	)

	svc := setupTestService(t)

	ctx = mockTransactWithContext(ctx, svc.txm)
	svc.orm.On("UpdateManager", ctx, mgr).Return(nil)
	svc.csaKeystore.On("GetAll").Return([]csakey.KeyV2{key}, nil)
	svc.connMgr.On("Disconnect", mgr.ID).Return(nil)
	svc.connMgr.On("Connect", mock.IsType(feeds.ConnectOpts{})).Return(nil)

	err := svc.UpdateFeedsManager(ctx, mgr)
	require.NoError(t, err)
}

func Test_Service_ListJobProposals(t *testing.T) {
	t.Parallel()

	var (
		jp  = feeds.JobProposal{}
		jps = []feeds.JobProposal{jp}
	)
	svc := setupTestService(t)

	svc.orm.On("ListJobProposals", context.Background()).
		Return(jps, nil)

	actual, err := svc.ListJobProposals()
	require.NoError(t, err)

	assert.Equal(t, actual, jps)
}

func Test_Service_GetJobProposal(t *testing.T) {
	t.Parallel()

	var (
		id = int64(1)
		ms = feeds.JobProposal{ID: id}
	)
	svc := setupTestService(t)

	svc.orm.On("GetJobProposal", context.Background(), id).
		Return(&ms, nil)

	actual, err := svc.GetJobProposal(id)
	require.NoError(t, err)

	assert.Equal(t, actual, &ms)
}

func Test_Service_ApproveJobProposal(t *testing.T) {
	var (
		ctx  = context.Background()
		spec = `name = 'LINK / ETH | version 3 | contract 0x0000000000000000000000000000000000000000'
schemaVersion = 1
contractAddress = '0x0000000000000000000000000000000000000000'
type = 'fluxmonitor'
externalJobID = '00000000-0000-0000-0000-000000000001'
threshold = 1.0
idleTimerPeriod = '4h'
idleTimerDisabled = false
pollingTimerPeriod = '1m'
pollingTimerDisabled = false
observationSource = """
// data source 1
ds1 [type=bridge name=\"bridge-api0\" requestData="{\\\"data\\": {\\\"from\\\":\\\"LINK\\\",\\\"to\\\":\\\"ETH\\\"}}"];
ds1_parse [type=jsonparse path="result"];
ds1_multiply [type=multiply times=1000000000000000000];
ds1 -> ds1_parse -> ds1_multiply -> answer1;

answer1 [type=median index=0];
"""
`
		pendingProposal = &feeds.JobProposal{
			ID:             1,
			RemoteUUID:     uuid.NewV4(),
			Status:         feeds.JobProposalStatusPending,
			FeedsManagerID: 2,
			Spec:           spec,
		}
		cancelledProposal = &feeds.JobProposal{
			ID:             1,
			RemoteUUID:     uuid.NewV4(),
			Status:         feeds.JobProposalStatusCancelled,
			FeedsManagerID: 2,
			Spec:           spec,
		}
		jb = job.Job{
			ID: int32(1),
		}
	)

	testCases := []struct {
		name    string
		before  func(svc *TestService)
		id      int64
		wantErr string
	}{
		{
			name: "pending job success",
			id:   pendingProposal.ID,
			before: func(svc *TestService) {
				svc.orm.On("GetJobProposal", ctx, pendingProposal.ID).Return(pendingProposal, nil)
				svc.connMgr.On("GetClient", pendingProposal.FeedsManagerID).Return(svc.fmsClient, nil)
				ctx = mockTransactWithContext(ctx, svc.txm)

				svc.cfg.On("DefaultHTTPTimeout").Return(models.MakeDuration(1 * time.Minute))
				svc.spawner.
					On("CreateJob",
						ctx,
						mock.MatchedBy(func(j job.Job) bool {
							return true
						}),
						null.StringFrom("LINK / ETH | version 3 | contract 0x0000000000000000000000000000000000000000"),
					).
					Return(jb, nil)
				svc.orm.On("ApproveJobProposal",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					pendingProposal.ID,
					uuid.Must(uuid.FromString("00000000-0000-0000-0000-000000000001")),
					feeds.JobProposalStatusApproved,
				).Return(nil)
				svc.fmsClient.On("ApprovedJob",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&proto.ApprovedJobRequest{
						Uuid: pendingProposal.RemoteUUID.String(),
					},
				).Return(&proto.ApprovedJobResponse{}, nil)
			},
		},
		{
			name: "cancelled job success",
			id:   cancelledProposal.ID,
			before: func(svc *TestService) {
				svc.orm.On("GetJobProposal", ctx, cancelledProposal.ID).Return(cancelledProposal, nil)
				svc.connMgr.On("GetClient", cancelledProposal.FeedsManagerID).Return(svc.fmsClient, nil)
				ctx = mockTransactWithContext(ctx, svc.txm)

				svc.cfg.On("DefaultHTTPTimeout").Return(models.MakeDuration(1 * time.Minute))
				svc.spawner.
					On("CreateJob",
						ctx,
						mock.MatchedBy(func(j job.Job) bool {
							return true
						}),
						null.StringFrom("LINK / ETH | version 3 | contract 0x0000000000000000000000000000000000000000"),
					).
					Return(jb, nil)
				svc.orm.On("ApproveJobProposal",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					cancelledProposal.ID,
					uuid.Must(uuid.FromString("00000000-0000-0000-0000-000000000001")),
					feeds.JobProposalStatusApproved,
				).Return(nil)
				svc.fmsClient.On("ApprovedJob",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&proto.ApprovedJobRequest{
						Uuid: cancelledProposal.RemoteUUID.String(),
					},
				).Return(&proto.ApprovedJobResponse{}, nil)
			},
		},
		{
			name: "job proposal does not exist",
			id:   int64(1),
			before: func(svc *TestService) {
				svc.orm.On("GetJobProposal", ctx, int64(1)).Return(nil, errors.New("Not Found"))
			},
			wantErr: "job proposal error: Not Found",
		},
		{
			name: "FMS client not connected",
			id:   int64(1),
			before: func(svc *TestService) {
				svc.orm.On("GetJobProposal", ctx, cancelledProposal.ID).Return(pendingProposal, nil)
				svc.connMgr.On("GetClient", pendingProposal.FeedsManagerID).Return(nil, errors.New("Not Connected"))
			},
			wantErr: "fms rpc client is not connected: Not Connected",
		},
		{
			name: "job proposal already approved",
			id:   int64(1),
			before: func(svc *TestService) {
				jp := &feeds.JobProposal{
					ID:             1,
					RemoteUUID:     uuid.NewV4(),
					Status:         feeds.JobProposalStatusApproved,
					FeedsManagerID: 2,
					Spec:           spec,
				}
				svc.orm.On("GetJobProposal", ctx, jp.ID).Return(jp, nil)
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
			},
			wantErr: "must be a pending or cancelled job proposal",
		},
		{
			name: "orm error",
			id:   int64(1),
			before: func(svc *TestService) {
				svc.orm.On("GetJobProposal", ctx, pendingProposal.ID).Return(pendingProposal, nil)
				svc.connMgr.On("GetClient", pendingProposal.FeedsManagerID).Return(svc.fmsClient, nil)
				ctx = mockTransactWithContext(ctx, svc.txm)

				svc.cfg.On("DefaultHTTPTimeout").Return(models.MakeDuration(1 * time.Minute))
				svc.spawner.
					On("CreateJob",
						ctx,
						mock.MatchedBy(func(j job.Job) bool {
							return true
						}),
						null.StringFrom("LINK / ETH | version 3 | contract 0x0000000000000000000000000000000000000000"),
					).
					Return(job.Job{}, errors.New("could not save"))
			},
			wantErr: "could not approve job proposal: could not save",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			svc := setupTestService(t)

			if tc.before != nil {
				tc.before(svc)
			}

			err := svc.ApproveJobProposal(ctx, tc.id)

			if tc.wantErr != "" {
				require.Error(t, err)
				assert.EqualError(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_Service_RejectJobProposal(t *testing.T) {
	var (
		ctx = context.Background()
		jp  = &feeds.JobProposal{
			ID:             1,
			RemoteUUID:     uuid.NewV4(),
			Status:         feeds.JobProposalStatusPending,
			FeedsManagerID: 2,
		}
	)

	svc := setupTestService(t)

	svc.orm.On("GetJobProposal", ctx, jp.ID).Return(jp, nil)
	ctx = mockTransactWithContext(ctx, svc.txm)
	svc.orm.On("UpdateJobProposalStatus",
		mock.MatchedBy(func(ctx context.Context) bool { return true }),
		jp.ID,
		feeds.JobProposalStatusRejected,
	).Return(nil)
	svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
	svc.fmsClient.On("RejectedJob",
		mock.MatchedBy(func(ctx context.Context) bool { return true }),
		&proto.RejectedJobRequest{
			Uuid: jp.RemoteUUID.String(),
		},
	).Return(&proto.RejectedJobResponse{}, nil)

	err := svc.RejectJobProposal(ctx, jp.ID)
	require.NoError(t, err)
}

func Test_Service_CancelJobProposal(t *testing.T) {
	var (
		externalJobID = uuid.NewV4()
		jp            = &feeds.JobProposal{
			ID:             1,
			ExternalJobID:  uuid.NullUUID{UUID: externalJobID, Valid: true},
			RemoteUUID:     externalJobID,
			Status:         feeds.JobProposalStatusApproved,
			FeedsManagerID: 2,
		}
		j = job.Job{
			ID:            1,
			ExternalJobID: externalJobID,
		}
	)

	testCases := []struct {
		name     string
		beforeFn func(svc *TestService)
		wantErr  string
	}{
		{
			name: "success",
			beforeFn: func(svc *TestService) {
				ctx := mockTransactWithContext(context.Background(), svc.txm)

				svc.orm.On("GetJobProposal", ctx, jp.ID).Return(jp, nil)
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.orm.On("CancelJobProposal",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					jp.ID,
				).Return(nil)
				svc.jobORM.On("FindJobByExternalJobID", ctx, externalJobID).Return(j, nil)
				svc.spawner.On("DeleteJob", ctx, j.ID).Return(nil)

				svc.fmsClient.On("CancelledJob",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&proto.CancelledJobRequest{
						Uuid: jp.RemoteUUID.String(),
					},
				).Return(&proto.CancelledJobResponse{}, nil)
			},
		},
		{
			name: "must be an approved job proposal",
			beforeFn: func(svc *TestService) {
				svc.orm.On("GetJobProposal", context.Background(), jp.ID).Return(&feeds.JobProposal{
					ID:             1,
					ExternalJobID:  uuid.NullUUID{UUID: externalJobID, Valid: true},
					RemoteUUID:     externalJobID,
					Status:         feeds.JobProposalStatusPending,
					FeedsManagerID: 2,
				}, nil)
			},
			wantErr: "must be a approved job proposal",
		},
		{
			name: "rpc client not connected",
			beforeFn: func(svc *TestService) {
				svc.orm.On("GetJobProposal", context.Background(), jp.ID).Return(jp, nil)
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(nil, errors.New("not connected"))
			},
			wantErr: "fms rpc client: not connected",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := setupTestService(t)

			tc.beforeFn(svc)

			err := svc.CancelJobProposal(context.Background(), jp.ID)

			if tc.wantErr != "" {
				require.Error(t, err)
				assert.EqualError(t, err, tc.wantErr)

				return
			}

			require.NoError(t, err)
		})
	}
}

func Test_Service_IsJobManaged(t *testing.T) {
	t.Parallel()

	svc := setupTestService(t)
	ctx := context.Background()
	jobID := int64(1)

	svc.orm.On("IsJobManaged", ctx, jobID).Return(true, nil)

	isManaged, err := svc.IsJobManaged(ctx, jobID)
	require.NoError(t, err)
	assert.True(t, isManaged)
}

func Test_Service_UpdateJobProposalSpec(t *testing.T) {
	var (
		ctx = context.Background()
		jp  = &feeds.JobProposal{
			ID:         1,
			RemoteUUID: uuid.NewV4(),
			Status:     feeds.JobProposalStatusPending,
			Spec:       "spec",
		}
		updatedSpec = "updated spec"
	)

	svc := setupTestService(t)

	svc.orm.On("GetJobProposal", ctx, jp.ID).Return(jp, nil)
	svc.orm.On("UpdateJobProposalSpec",
		mock.MatchedBy(func(ctx context.Context) bool { return true }),
		jp.ID,
		updatedSpec,
	).Return(nil)

	err := svc.UpdateJobProposalSpec(ctx, jp.ID, updatedSpec)
	require.NoError(t, err)
}

func Test_Service_StartStop(t *testing.T) {
	key := cltest.DefaultCSAKey

	var (
		mgr = feeds.FeedsManager{
			ID:  1,
			URI: "localhost:2000",
		}
		pubKeyHex = "0f17c3bf72de8beef6e2d17a14c0a972f5d7e0e66e70722373f12b88382d40f9"
	)

	var pubKey crypto.PublicKey
	_, err := hex.Decode([]byte(pubKeyHex), pubKey)
	require.NoError(t, err)

	svc := setupTestService(t)

	svc.csaKeystore.On("GetAll").Return([]csakey.KeyV2{key}, nil)
	svc.orm.On("ListManagers", context.Background()).Return([]feeds.FeedsManager{mgr}, nil)
	svc.connMgr.On("IsConnected", mgr.ID).Return(false)
	svc.connMgr.On("Connect", mock.IsType(feeds.ConnectOpts{}))
	svc.connMgr.On("Close")

	err = svc.Start()
	require.NoError(t, err)

	svc.Close()
}

func mockTransactWithContext(ctx context.Context, txm *pgmocks.TransactionManager) context.Context {
	call := txm.On("TransactWithContext",
		mock.MatchedBy(func(ctx context.Context) bool { return true }),
		mock.Anything,
	)
	call.Run(func(args mock.Arguments) {
		arg := args.Get(1).(postgres.TxFn)
		err := arg(ctx)
		call.Return(err)
	})

	return ctx
}
