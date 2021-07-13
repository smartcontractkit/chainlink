package feeds_test

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/lib/pq"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/keystest"
	"github.com/smartcontractkit/chainlink/core/services/feeds"
	"github.com/smartcontractkit/chainlink/core/services/feeds/mocks"
	"github.com/smartcontractkit/chainlink/core/services/feeds/proto"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/csakey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	ksmocks "github.com/smartcontractkit/chainlink/core/services/keystore/mocks"
	"github.com/smartcontractkit/chainlink/core/utils/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type TestService struct {
	feeds.Service
	orm         *mocks.ORM
	fmsClient   *mocks.FeedsManagerClient
	csaKeystore *ksmocks.CSAKeystoreInterface
	ethKeystore *ksmocks.EthKeyStoreInterface
	cfg         *mocks.Config
}

func setupTestService(t *testing.T) *TestService {
	orm := &mocks.ORM{}
	fmsClient := &mocks.FeedsManagerClient{}
	csaKeystore := &ksmocks.CSAKeystoreInterface{}
	ethKeystore := &ksmocks.EthKeyStoreInterface{}
	cfg := &mocks.Config{}

	t.Cleanup(func() {
		mock.AssertExpectationsForObjects(t,
			orm,
			fmsClient,
			csaKeystore,
			ethKeystore,
			cfg,
		)
	})

	svc := feeds.NewService(orm, csaKeystore, ethKeystore, cfg)
	svc.SetFMSClient(fmsClient)

	return &TestService{
		Service:     svc,
		orm:         orm,
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

func Test_Service_GetManagers(t *testing.T) {
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
		ctx      = context.Background()
		feedsMgr = &feeds.FeedsManager{
			ID:       1,
			JobTypes: pq.StringArray{feeds.JobTypeFluxMonitor},
		}
		chainID    = big.NewInt(1)
		fundingKey = ethkey.Key{
			Address:   ethkey.EIP55AddressFromAddress(rawKey.Address),
			IsFunding: true,
		}
	)

	svc := setupTestService(t)

	// Mock fetching the information to send
	svc.orm.On("GetManager", ctx, feedsMgr.ID).Return(feedsMgr, nil)
	svc.ethKeystore.On("FundingKeys").Return([]ethkey.Key{fundingKey}, nil)
	svc.cfg.On("ChainID").Return(chainID)

	// Mock the send
	svc.fmsClient.On("UpdateNode", ctx, &proto.UpdateNodeRequest{
		JobTypes:         []proto.JobType{proto.JobType_JOB_TYPE_FLUX_MONITOR},
		ChainId:          chainID.Int64(),
		FundingAddresses: []string{fundingKey.Address.String()},
	}).Return(&proto.UpdateNodeResponse{}, nil)

	err = svc.SyncNodeInfo(feedsMgr.ID)
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
