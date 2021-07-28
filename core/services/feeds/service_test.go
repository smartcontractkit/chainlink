package feeds_test

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"math/big"
	"testing"
	"time"

	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
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
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"
)

type TestService struct {
	feeds.Service
	orm         *mocks.ORM
	txm         *pgmocks.TransactionManager
	spawner     *jobmocks.Spawner
	fmsClient   *mocks.FeedsManagerClient
	csaKeystore *ksmocks.CSAKeystoreInterface
	ethKeystore *ksmocks.EthKeyStoreInterface
	cfg         *mocks.Config
}

func setupTestService(t *testing.T) *TestService {
	var (
		orm         = &mocks.ORM{}
		txm         = &pgmocks.TransactionManager{}
		spawner     = &jobmocks.Spawner{}
		fmsClient   = &mocks.FeedsManagerClient{}
		csaKeystore = &ksmocks.CSAKeystoreInterface{}
		ethKeystore = &ksmocks.EthKeyStoreInterface{}
		cfg         = &mocks.Config{}
	)

	t.Cleanup(func() {
		mock.AssertExpectationsForObjects(t,
			orm,
			txm,
			spawner,
			fmsClient,
			csaKeystore,
			ethKeystore,
			cfg,
		)
	})

	svc := feeds.NewService(orm, txm, spawner, csaKeystore, ethKeystore, cfg)
	svc.SetFMSClient(fmsClient)

	return &TestService{
		Service:     svc,
		orm:         orm,
		txm:         txm,
		spawner:     spawner,
		fmsClient:   fmsClient,
		csaKeystore: csaKeystore,
		ethKeystore: ethKeystore,
		cfg:         cfg,
	}
}

func Test_Service_RegisterManager(t *testing.T) {
	t.Parallel()

	_, privkey, err := ed25519.GenerateKey(nil)
	require.NoError(t, err)

	var (
		id        = int64(1)
		ms        = feeds.FeedsManager{}
		pubKeyHex = "0f17c3bf72de8beef6e2d17a14c0a972f5d7e0e66e70722373f12b88382d40f9"
	)

	var pubKey crypto.PublicKey
	_, err = hex.Decode([]byte(pubKeyHex), pubKey)
	require.NoError(t, err)
	key := csakey.Key{
		PublicKey: pubKey,
	}

	svc := setupTestService(t)

	svc.orm.On("CountManagers").Return(int64(0), nil)
	svc.orm.On("CreateManager", context.Background(), &ms).
		Return(id, nil)
	svc.csaKeystore.On("ListCSAKeys").Return([]csakey.Key{key}, nil)
	svc.csaKeystore.On("Unsafe_GetUnlockedPrivateKey", pubKey).Return([]byte(privkey), nil)
	// ListManagers runs in a goroutine so it might be called.
	svc.orm.On("ListManagers", context.Background()).Return([]feeds.FeedsManager{ms}, nil).Maybe()

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

	actual, err := svc.GetManager(id)
	require.NoError(t, err)

	assert.Equal(t, actual, &ms)
}

func Test_Service_CreateJobProposal(t *testing.T) {
	t.Parallel()

	var (
		id = int64(1)
		jp = feeds.JobProposal{}
	)
	svc := setupTestService(t)

	svc.orm.On("CreateJobProposal", context.Background(), &jp).
		Return(id, nil)

	actual, err := svc.CreateJobProposal(&jp)
	require.NoError(t, err)

	assert.Equal(t, actual, id)
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
		sendingKey = ethkey.Key{
			Address:   ethkey.EIP55AddressFromAddress(rawKey.Address),
			IsFunding: false,
		}
	)

	svc := setupTestService(t)

	// Mock fetching the information to send
	svc.orm.On("GetManager", ctx, feedsMgr.ID).Return(feedsMgr, nil)
	svc.ethKeystore.On("SendingKeys").Return([]ethkey.Key{sendingKey}, nil)
	svc.cfg.On("ChainID").Return(chainID)

	// Mock the send
	svc.fmsClient.On("UpdateNode", ctx, &proto.UpdateNodeRequest{
		JobTypes:           []proto.JobType{proto.JobType_JOB_TYPE_FLUX_MONITOR},
		ChainId:            chainID.Int64(),
		AccountAddresses:   []string{sendingKey.Address.String()},
		IsBootstrapPeer:    true,
		BootstrapMultiaddr: multiaddr,
	}).Return(&proto.UpdateNodeResponse{}, nil)

	err = svc.SyncNodeInfo(feedsMgr.ID)
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
		ctx = context.Background()
		jp  = &feeds.JobProposal{
			ID:         1,
			RemoteUUID: uuid.NewV4(),
			Status:     feeds.JobProposalStatusPending,
			Spec: `name = 'LINK / ETH | version 3 | contract 0x0000000000000000000000000000000000000000'
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

			// data source 2
			ds2 [type=bridge name="bridge-api1" requestData="{\\\"data\\\": {\\\"from\\\":\\\"LINK\\\",\\\"to\\\":\\\"ETH\\\"}}"];
			ds2_parse [type=jsonparse path="result"];
			ds2_multiply [type=multiply times=1000000000000000000];
			ds2 -> ds2_parse -> ds2_multiply -> answer1;

			answer1 [type=median index=0];
			"""
			`,
		}
		jb = job.Job{
			ID: int32(1),
		}
	)

	svc := setupTestService(t)

	svc.orm.On("GetJobProposal", ctx, jp.ID).Return(jp, nil)
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
		jp.ID,
		uuid.Must(uuid.FromString("00000000-0000-0000-0000-000000000001")),
		feeds.JobProposalStatusApproved,
	).Return(nil)
	svc.fmsClient.On("ApprovedJob",
		mock.MatchedBy(func(ctx context.Context) bool { return true }),
		&proto.ApprovedJobRequest{
			Uuid: jp.RemoteUUID.String(),
		},
	).Return(&proto.ApprovedJobResponse{}, nil)

	err := svc.ApproveJobProposal(ctx, jp.ID)
	require.NoError(t, err)
}

func Test_Service_RejectJobProposal(t *testing.T) {
	var (
		ctx = context.Background()
		jp  = &feeds.JobProposal{
			ID:         1,
			RemoteUUID: uuid.NewV4(),
			Status:     feeds.JobProposalStatusPending,
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
	svc.fmsClient.On("RejectedJob",
		mock.MatchedBy(func(ctx context.Context) bool { return true }),
		&proto.RejectedJobRequest{
			Uuid: jp.RemoteUUID.String(),
		},
	).Return(&proto.RejectedJobResponse{}, nil)

	err := svc.RejectJobProposal(ctx, jp.ID)
	require.NoError(t, err)
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
	_, privkey, err := ed25519.GenerateKey(nil)
	require.NoError(t, err)

	var (
		ms        = feeds.FeedsManager{}
		pubKeyHex = "0f17c3bf72de8beef6e2d17a14c0a972f5d7e0e66e70722373f12b88382d40f9"
	)

	var pubKey crypto.PublicKey
	_, err = hex.Decode([]byte(pubKeyHex), pubKey)
	require.NoError(t, err)
	key := csakey.Key{
		PublicKey: pubKey,
	}

	svc := setupTestService(t)

	svc.csaKeystore.On("ListCSAKeys").Return([]csakey.Key{key}, nil)
	svc.csaKeystore.On("Unsafe_GetUnlockedPrivateKey", pubKey).Return([]byte(privkey), nil)
	svc.orm.On("ListManagers", context.Background()).Return([]feeds.FeedsManager{ms}, nil)

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
