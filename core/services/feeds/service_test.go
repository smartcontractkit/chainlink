package feeds_test

import (
	"context"
	"database/sql"
	"encoding/hex"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	configtest "github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/feeds"
	"github.com/smartcontractkit/chainlink/v2/core/services/feeds/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/feeds/proto"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	jobmocks "github.com/smartcontractkit/chainlink/v2/core/services/job/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/csakey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocrkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
	ksmocks "github.com/smartcontractkit/chainlink/v2/core/services/keystore/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/services/versioning"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
	"github.com/smartcontractkit/chainlink/v2/core/utils/crypto"

	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"
)

const FluxMonitorTestSpec = `
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

const OCR1TestSpec = `
type               = "offchainreporting"
schemaVersion      = 1
name              = "example OCR1 spec"
externalJobID       = "0EEC7E1D-D0D2-476C-A1A8-72DFB6633F46"
contractAddress    = "0x613a38AC1659769640aaE063C651F48E0250454C"
p2pBootstrapPeers  = [
	"/dns4/chain.link/tcp/1234/p2p/16Uiu2HAm58SP7UL8zsnpeuwHfytLocaqgnyaYKP8wu7qRdrixLju",
]
p2pv2Bootstrappers = []
keyBundleID        = "f5bf259689b26f1374efb3c9a9868796953a0f814bb2d39b968d0e61b58620a5"
transmitterAddress = "0x613a38AC1659769640aaE063C651F48E0250454C"
isBootstrapPeer		= false
observationSource = """
	// data source 1
	ds1          [type=memo value=<"{\\"USD\\": 1}">];
	ds1_parse    [type=jsonparse path="USD"];
	ds1_multiply [type=multiply times=3];

	ds2          [type=memo value=<"{\\"USD\\": 1}">];
	ds2_parse    [type=jsonparse path="USD"];
	ds2_multiply [type=multiply times=3];

	ds3          [type=fail msg="uh oh"];

	ds1 -> ds1_parse -> ds1_multiply -> answer;
	ds2 -> ds2_parse -> ds2_multiply -> answer;
	ds3 -> answer;

	answer [type=median index=0];
"""
`

const OCR2TestSpec = `
type               = "offchainreporting2"
pluginType         = "median"
schemaVersion      = 1
name              = "example OCR2 spec"
relay              = "evm"
contractID         = "0x613a38AC1659769640aaE063C651F48E0250454C"
externalJobID      = "0EEC7E1D-D0D2-476C-A1A8-72DFB6633F47"
observationSource  = """
ds1          [type=bridge name=voter_turnout];
ds1_parse    [type=jsonparse path="one,two"];
ds1_multiply [type=multiply times=1.23];
ds1 -> ds1_parse -> ds1_multiply -> answer1;
answer1      [type=median index=0];
"""
[relayConfig]
chainID = 1337
[pluginConfig]
juelsPerFeeCoinSource = """
ds1          [type=bridge name=voter_turnout];
ds1_parse    [type=jsonparse path="one,two"];
ds1_multiply [type=multiply times=1.23];
ds1 -> ds1_parse -> ds1_multiply -> answer1;
answer1      [type=median index=0];
"""
`
const BootstrapTestSpec = `
type				= "bootstrap"
schemaVersion		= 1
name              = "example Bootstrap spec"
contractID			= "0x613a38AC1659769640aaE063C651F48E0250454C"
relay				= "evm"
[relayConfig]
chainID 			= 1337
`

type TestService struct {
	feeds.Service
	orm          *mocks.ORM
	jobORM       *jobmocks.ORM
	connMgr      *mocks.ConnectionsManager
	spawner      *jobmocks.Spawner
	fmsClient    *mocks.FeedsManagerClient
	csaKeystore  *ksmocks.CSA
	p2pKeystore  *ksmocks.P2P
	ocr1Keystore *ksmocks.OCR
	ocr2Keystore *ksmocks.OCR2
	cc           evm.ChainSet
}

func setupTestService(t *testing.T) *TestService {
	t.Helper()

	return setupTestServiceCfg(t, nil)
}

func setupTestServiceCfg(t *testing.T, overrideCfg func(c *chainlink.Config, s *chainlink.Secrets)) *TestService {
	t.Helper()

	var (
		orm          = mocks.NewORM(t)
		jobORM       = jobmocks.NewORM(t)
		connMgr      = mocks.NewConnectionsManager(t)
		spawner      = jobmocks.NewSpawner(t)
		fmsClient    = mocks.NewFeedsManagerClient(t)
		csaKeystore  = ksmocks.NewCSA(t)
		p2pKeystore  = ksmocks.NewP2P(t)
		ocr1Keystore = ksmocks.NewOCR(t)
		ocr2Keystore = ksmocks.NewOCR2(t)
	)

	lggr := logger.TestLogger(t)

	db := pgtest.NewSqlxDB(t)
	gcfg := configtest.NewGeneralConfig(t, overrideCfg)
	keyStore := new(ksmocks.Master)
	scopedConfig := evmtest.NewChainScopedConfig(t, gcfg)
	ethKeyStore := cltest.NewKeyStore(t, db, gcfg).Eth()
	cc := evmtest.NewChainSet(t, evmtest.TestChainOpts{GeneralConfig: gcfg,
		HeadTracker: headtracker.NullTracker, KeyStore: ethKeyStore})
	keyStore.On("Eth").Return(ethKeyStore)
	keyStore.On("CSA").Return(csaKeystore)
	keyStore.On("P2P").Return(p2pKeystore)
	keyStore.On("OCR").Return(ocr1Keystore)
	keyStore.On("OCR2").Return(ocr2Keystore)
	svc := feeds.NewService(orm, jobORM, db, spawner, keyStore, scopedConfig, cc, lggr, "1.0.0")
	svc.SetConnectionsManager(connMgr)

	return &TestService{
		Service:      svc,
		orm:          orm,
		jobORM:       jobORM,
		connMgr:      connMgr,
		spawner:      spawner,
		fmsClient:    fmsClient,
		csaKeystore:  csaKeystore,
		p2pKeystore:  p2pKeystore,
		ocr1Keystore: ocr1Keystore,
		ocr2Keystore: ocr2Keystore,
		cc:           cc,
	}
}

func Test_Service_RegisterManager(t *testing.T) {
	t.Parallel()

	key := cltest.DefaultCSAKey

	var (
		id        = int64(1)
		pubKeyHex = "0f17c3bf72de8beef6e2d17a14c0a972f5d7e0e66e70722373f12b88382d40f9"
	)

	var pubKey crypto.PublicKey
	_, err := hex.Decode([]byte(pubKeyHex), pubKey)
	require.NoError(t, err)

	var (
		mgr = feeds.FeedsManager{
			Name:      "FMS",
			URI:       "localhost:8080",
			PublicKey: pubKey,
		}
		params = feeds.RegisterManagerParams{
			Name:      "FMS",
			URI:       "localhost:8080",
			PublicKey: pubKey,
		}
	)

	svc := setupTestService(t)

	svc.orm.On("CountManagers").Return(int64(0), nil)
	svc.orm.On("CreateManager", &mgr, mock.Anything).
		Return(id, nil)
	svc.orm.On("CreateBatchChainConfig", params.ChainConfigs, mock.Anything).
		Return([]int64{}, nil)
	svc.csaKeystore.On("GetAll").Return([]csakey.KeyV2{key}, nil)
	// ListManagers runs in a goroutine so it might be called.
	svc.orm.On("ListManagers", testutils.Context(t)).Return([]feeds.FeedsManager{mgr}, nil).Maybe()
	svc.connMgr.On("Connect", mock.IsType(feeds.ConnectOpts{}))

	actual, err := svc.RegisterManager(testutils.Context(t), params)
	// We need to stop the service because the manager will attempt to make a
	// connection
	svc.Close()
	require.NoError(t, err)

	assert.Equal(t, actual, id)
}

func Test_Service_ListManagers(t *testing.T) {
	t.Parallel()

	var (
		mgr  = feeds.FeedsManager{}
		mgrs = []feeds.FeedsManager{mgr}
	)
	svc := setupTestService(t)

	svc.orm.On("ListManagers").Return(mgrs, nil)
	svc.connMgr.On("IsConnected", mgr.ID).Return(false)

	actual, err := svc.ListManagers()
	require.NoError(t, err)

	assert.Equal(t, mgrs, actual)
}

func Test_Service_GetManager(t *testing.T) {
	t.Parallel()

	var (
		id  = int64(1)
		mgr = feeds.FeedsManager{ID: id}
	)
	svc := setupTestService(t)

	svc.orm.On("GetManager", id).
		Return(&mgr, nil)
	svc.connMgr.On("IsConnected", mgr.ID).Return(false)

	actual, err := svc.GetManager(id)
	require.NoError(t, err)

	assert.Equal(t, actual, &mgr)
}

func Test_Service_UpdateFeedsManager(t *testing.T) {
	key := cltest.DefaultCSAKey

	var (
		mgr = feeds.FeedsManager{ID: 1}
	)

	svc := setupTestService(t)

	svc.orm.On("UpdateManager", mgr, mock.Anything).Return(nil)
	svc.csaKeystore.On("GetAll").Return([]csakey.KeyV2{key}, nil)
	svc.connMgr.On("Disconnect", mgr.ID).Return(nil)
	svc.connMgr.On("Connect", mock.IsType(feeds.ConnectOpts{})).Return(nil)

	err := svc.UpdateManager(testutils.Context(t), mgr)
	require.NoError(t, err)
}

func Test_Service_ListManagersByIDs(t *testing.T) {
	t.Parallel()

	var (
		mgr  = feeds.FeedsManager{}
		mgrs = []feeds.FeedsManager{mgr}
	)
	svc := setupTestService(t)

	svc.orm.On("ListManagersByIDs", []int64{mgr.ID}).
		Return(mgrs, nil)
	svc.connMgr.On("IsConnected", mgr.ID).Return(false)

	actual, err := svc.ListManagersByIDs([]int64{mgr.ID})
	require.NoError(t, err)

	assert.Equal(t, mgrs, actual)
}

func Test_Service_CountManagers(t *testing.T) {
	t.Parallel()

	var (
		count = int64(1)
	)
	svc := setupTestService(t)

	svc.orm.On("CountManagers").
		Return(count, nil)

	actual, err := svc.CountManagers()
	require.NoError(t, err)

	assert.Equal(t, count, actual)
}

func Test_Service_CreateChainConfig(t *testing.T) {
	var (
		mgr         = feeds.FeedsManager{ID: 1}
		nodeVersion = &versioning.NodeVersion{
			Version: "1.0.0",
		}
		cfg = feeds.ChainConfig{
			FeedsManagerID: mgr.ID,
			ChainID:        "42",
			ChainType:      feeds.ChainTypeEVM,
			AccountAddress: "0x0000000000000000000000000000000000000000",
			AdminAddress:   "0x0000000000000000000000000000000000000001",
			FluxMonitorConfig: feeds.FluxMonitorConfig{
				Enabled: true,
			},
			OCR1Config: feeds.OCR1Config{
				Enabled: false,
			},
			OCR2Config: feeds.OCR2Config{
				Enabled: false,
			},
		}

		svc = setupTestService(t)
	)

	svc.orm.On("CreateChainConfig", cfg).Return(int64(1), nil)
	svc.orm.On("GetManager", mgr.ID).Return(&mgr, nil)
	svc.connMgr.On("GetClient", mgr.ID).Return(svc.fmsClient, nil)
	svc.orm.On("ListChainConfigsByManagerIDs", []int64{mgr.ID}).Return([]feeds.ChainConfig{cfg}, nil)
	svc.fmsClient.On("UpdateNode", mock.Anything, &proto.UpdateNodeRequest{
		Version: nodeVersion.Version,
		ChainConfigs: []*proto.ChainConfig{
			{
				Chain: &proto.Chain{
					Id:   cfg.ChainID,
					Type: proto.ChainType_CHAIN_TYPE_EVM,
				},
				AccountAddress:    cfg.AccountAddress,
				AdminAddress:      cfg.AdminAddress,
				FluxMonitorConfig: &proto.FluxMonitorConfig{Enabled: true},
				Ocr1Config:        &proto.OCR1Config{Enabled: false},
				Ocr2Config:        &proto.OCR2Config{Enabled: false},
			},
		},
	}).Return(&proto.UpdateNodeResponse{}, nil)

	actual, err := svc.CreateChainConfig(testutils.Context(t), cfg)
	require.NoError(t, err)
	assert.Equal(t, int64(1), actual)
}

func Test_Service_CreateChainConfig_InvalidAdminAddress(t *testing.T) {
	var (
		mgr = feeds.FeedsManager{ID: 1}
		cfg = feeds.ChainConfig{
			FeedsManagerID:    mgr.ID,
			ChainID:           "42",
			ChainType:         feeds.ChainTypeEVM,
			AccountAddress:    "0x0000000000000000000000000000000000000000",
			AdminAddress:      "0x00000000000",
			FluxMonitorConfig: feeds.FluxMonitorConfig{Enabled: false},
			OCR1Config:        feeds.OCR1Config{Enabled: false},
			OCR2Config:        feeds.OCR2Config{Enabled: false},
		}

		svc = setupTestService(t)
	)
	_, err := svc.CreateChainConfig(testutils.Context(t), cfg)
	require.Error(t, err)
	assert.Equal(t, "invalid admin address: 0x00000000000", err.Error())
}

func Test_Service_DeleteChainConfig(t *testing.T) {
	var (
		mgr         = feeds.FeedsManager{ID: 1}
		nodeVersion = &versioning.NodeVersion{
			Version: "1.0.0",
		}
		cfg = feeds.ChainConfig{
			ID:             1,
			FeedsManagerID: mgr.ID,
		}

		svc = setupTestService(t)
	)

	svc.orm.On("GetChainConfig", cfg.ID).Return(&cfg, nil)
	svc.orm.On("DeleteChainConfig", cfg.ID).Return(cfg.ID, nil)
	svc.orm.On("GetManager", mgr.ID).Return(&mgr, nil)
	svc.connMgr.On("GetClient", mgr.ID).Return(svc.fmsClient, nil)
	svc.orm.On("ListChainConfigsByManagerIDs", []int64{mgr.ID}).Return([]feeds.ChainConfig{}, nil)
	svc.fmsClient.On("UpdateNode", mock.Anything, &proto.UpdateNodeRequest{
		Version:      nodeVersion.Version,
		ChainConfigs: []*proto.ChainConfig{},
	}).Return(&proto.UpdateNodeResponse{}, nil)

	actual, err := svc.DeleteChainConfig(testutils.Context(t), cfg.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(1), actual)
}

func Test_Service_ListChainConfigsByManagerIDs(t *testing.T) {
	var (
		mgr = feeds.FeedsManager{ID: 1}
		cfg = feeds.ChainConfig{
			ID:             1,
			FeedsManagerID: mgr.ID,
		}
		ids = []int64{cfg.ID}

		svc = setupTestService(t)
	)

	svc.orm.On("ListChainConfigsByManagerIDs", ids).Return([]feeds.ChainConfig{cfg}, nil)

	actual, err := svc.ListChainConfigsByManagerIDs(ids)
	require.NoError(t, err)
	assert.Equal(t, []feeds.ChainConfig{cfg}, actual)
}

func Test_Service_UpdateChainConfig(t *testing.T) {
	var (
		mgr         = feeds.FeedsManager{ID: 1}
		nodeVersion = &versioning.NodeVersion{
			Version: "1.0.0",
		}
		cfg = feeds.ChainConfig{
			FeedsManagerID:    mgr.ID,
			ChainID:           "42",
			ChainType:         feeds.ChainTypeEVM,
			AccountAddress:    "0x0000000000000000000000000000000000000000",
			AdminAddress:      "0x0000000000000000000000000000000000000001",
			FluxMonitorConfig: feeds.FluxMonitorConfig{Enabled: false},
			OCR1Config:        feeds.OCR1Config{Enabled: false},
			OCR2Config:        feeds.OCR2Config{Enabled: false},
		}

		svc = setupTestService(t)
	)

	svc.orm.On("UpdateChainConfig", cfg).Return(int64(1), nil)
	svc.orm.On("GetChainConfig", cfg.ID).Return(&cfg, nil)
	svc.connMgr.On("GetClient", mgr.ID).Return(svc.fmsClient, nil)
	svc.orm.On("ListChainConfigsByManagerIDs", []int64{mgr.ID}).Return([]feeds.ChainConfig{cfg}, nil)
	svc.fmsClient.On("UpdateNode", mock.Anything, &proto.UpdateNodeRequest{
		Version: nodeVersion.Version,
		ChainConfigs: []*proto.ChainConfig{
			{
				Chain: &proto.Chain{
					Id:   cfg.ChainID,
					Type: proto.ChainType_CHAIN_TYPE_EVM,
				},
				AccountAddress:    cfg.AccountAddress,
				AdminAddress:      cfg.AdminAddress,
				FluxMonitorConfig: &proto.FluxMonitorConfig{Enabled: false},
				Ocr1Config:        &proto.OCR1Config{Enabled: false},
				Ocr2Config:        &proto.OCR2Config{Enabled: false},
			},
		},
	}).Return(&proto.UpdateNodeResponse{}, nil)

	actual, err := svc.UpdateChainConfig(testutils.Context(t), cfg)
	require.NoError(t, err)
	assert.Equal(t, int64(1), actual)
}

func Test_Service_UpdateChainConfig_InvalidAdminAddress(t *testing.T) {
	var (
		mgr = feeds.FeedsManager{ID: 1}
		cfg = feeds.ChainConfig{
			FeedsManagerID:    mgr.ID,
			ChainID:           "42",
			ChainType:         feeds.ChainTypeEVM,
			AccountAddress:    "0x0000000000000000000000000000000000000000",
			AdminAddress:      "0x00000000000",
			FluxMonitorConfig: feeds.FluxMonitorConfig{Enabled: false},
			OCR1Config:        feeds.OCR1Config{Enabled: false},
			OCR2Config:        feeds.OCR2Config{Enabled: false},
		}

		svc = setupTestService(t)
	)
	_, err := svc.UpdateChainConfig(testutils.Context(t), cfg)
	require.Error(t, err)
	assert.Equal(t, "invalid admin address: 0x00000000000", err.Error())
}

func Test_Service_ProposeJob(t *testing.T) {
	t.Parallel()

	var (
		idFluxMonitor         = int64(1)
		remoteUUIDFluxMonitor = uuid.NewV4()
		argsFluxMonitor       = &feeds.ProposeJobArgs{
			FeedsManagerID: 1,
			RemoteUUID:     remoteUUIDFluxMonitor,
			Spec:           FluxMonitorTestSpec,
			Version:        1,
		}
		jpFluxMonitor = feeds.JobProposal{
			FeedsManagerID: 1,
			Name:           null.StringFrom("example flux monitor spec"),
			RemoteUUID:     remoteUUIDFluxMonitor,
			Status:         feeds.JobProposalStatusPending,
		}
		specFluxMonitor = feeds.JobProposalSpec{
			Definition:    FluxMonitorTestSpec,
			Status:        feeds.SpecStatusPending,
			Version:       argsFluxMonitor.Version,
			JobProposalID: idFluxMonitor,
		}

		idOCR1         = int64(2)
		remoteUUIDOCR1 = uuid.NewV4()
		argsOCR1       = &feeds.ProposeJobArgs{
			FeedsManagerID: 1,
			RemoteUUID:     remoteUUIDOCR1,
			Spec:           OCR1TestSpec,
			Version:        1,
		}
		jpOCR1 = feeds.JobProposal{
			FeedsManagerID: 1,
			Name:           null.StringFrom("example OCR1 spec"),
			RemoteUUID:     remoteUUIDOCR1,
			Status:         feeds.JobProposalStatusPending,
		}
		specOCR1 = feeds.JobProposalSpec{
			Definition:    OCR1TestSpec,
			Status:        feeds.SpecStatusPending,
			Version:       argsOCR1.Version,
			JobProposalID: idOCR1,
		}

		idOCR2         = int64(3)
		remoteUUIDOCR2 = uuid.NewV4()
		argsOCR2       = &feeds.ProposeJobArgs{
			FeedsManagerID: 1,
			RemoteUUID:     remoteUUIDOCR2,
			Spec:           OCR2TestSpec,
			Version:        1,
		}
		jpOCR2 = feeds.JobProposal{
			FeedsManagerID: 1,
			Name:           null.StringFrom("example OCR2 spec"),
			RemoteUUID:     remoteUUIDOCR2,
			Status:         feeds.JobProposalStatusPending,
		}
		specOCR2 = feeds.JobProposalSpec{
			Definition:    OCR2TestSpec,
			Status:        feeds.SpecStatusPending,
			Version:       argsOCR2.Version,
			JobProposalID: idOCR2,
		}

		idBootstrap         = int64(4)
		remoteUUIDBootstrap = uuid.NewV4()
		argsBootstrap       = &feeds.ProposeJobArgs{
			FeedsManagerID: 1,
			RemoteUUID:     remoteUUIDBootstrap,
			Spec:           BootstrapTestSpec,
			Version:        1,
		}
		jpBootstrap = feeds.JobProposal{
			FeedsManagerID: 1,
			Name:           null.StringFrom("example Bootstrap spec"),
			RemoteUUID:     remoteUUIDBootstrap,
			Status:         feeds.JobProposalStatusPending,
		}
		specBootstrap = feeds.JobProposalSpec{
			Definition:    BootstrapTestSpec,
			Status:        feeds.SpecStatusPending,
			Version:       argsBootstrap.Version,
			JobProposalID: idBootstrap,
		}

		httpTimeout = models.MustMakeDuration(1 * time.Second)
	)

	testCases := []struct {
		name    string
		args    *feeds.ProposeJobArgs
		before  func(svc *TestService)
		wantID  int64
		wantErr string
	}{
		{
			name: "Create success (Flux Monitor)",
			before: func(svc *TestService) {
				svc.orm.On("GetJobProposalByRemoteUUID", jpFluxMonitor.RemoteUUID).Return(new(feeds.JobProposal), sql.ErrNoRows)
				svc.orm.On("UpsertJobProposal", &jpFluxMonitor, mock.Anything).Return(idFluxMonitor, nil)
				svc.orm.On("CreateSpec", specFluxMonitor, mock.Anything).Return(int64(100), nil)
				svc.orm.On("CountJobProposalsByStatus").Return(&feeds.JobProposalCounts{}, nil)
			},
			args:   argsFluxMonitor,
			wantID: idFluxMonitor,
		},
		{
			name: "Create success (OCR1)",
			before: func(svc *TestService) {
				svc.orm.On("GetJobProposalByRemoteUUID", jpOCR1.RemoteUUID).Return(new(feeds.JobProposal), sql.ErrNoRows)
				svc.orm.On("UpsertJobProposal", &jpOCR1, mock.Anything).Return(idOCR1, nil)
				svc.orm.On("CreateSpec", specOCR1, mock.Anything).Return(int64(100), nil)
				svc.orm.On("CountJobProposalsByStatus").Return(&feeds.JobProposalCounts{}, nil)
			},
			args:   argsOCR1,
			wantID: idOCR1,
		},
		{
			name: "Create success (OCR2)",
			before: func(svc *TestService) {
				svc.orm.On("GetJobProposalByRemoteUUID", jpOCR2.RemoteUUID).Return(new(feeds.JobProposal), sql.ErrNoRows)
				svc.orm.On("UpsertJobProposal", &jpOCR2, mock.Anything).Return(idOCR2, nil)
				svc.orm.On("CreateSpec", specOCR2, mock.Anything).Return(int64(100), nil)
				svc.orm.On("CountJobProposalsByStatus").Return(&feeds.JobProposalCounts{}, nil)
			},
			args:   argsOCR2,
			wantID: idOCR2,
		},
		{
			name: "Create success (Bootstrap)",
			before: func(svc *TestService) {
				svc.orm.On("GetJobProposalByRemoteUUID", jpBootstrap.RemoteUUID).Return(new(feeds.JobProposal), sql.ErrNoRows)
				svc.orm.On("UpsertJobProposal", &jpBootstrap, mock.Anything).Return(idBootstrap, nil)
				svc.orm.On("CreateSpec", specBootstrap, mock.Anything).Return(int64(102), nil)
				svc.orm.On("CountJobProposalsByStatus").Return(&feeds.JobProposalCounts{}, nil)
			},
			args:   argsBootstrap,
			wantID: idBootstrap,
		},
		{
			name: "Update success",
			before: func(svc *TestService) {
				svc.orm.
					On("GetJobProposalByRemoteUUID", jpFluxMonitor.RemoteUUID).
					Return(&feeds.JobProposal{
						FeedsManagerID: jpFluxMonitor.FeedsManagerID,
						RemoteUUID:     jpFluxMonitor.RemoteUUID,
						Status:         feeds.JobProposalStatusPending,
					}, nil)
				svc.orm.On("ExistsSpecByJobProposalIDAndVersion", jpFluxMonitor.ID, argsFluxMonitor.Version).Return(false, nil)
				svc.orm.On("UpsertJobProposal", &jpFluxMonitor, mock.Anything).Return(idFluxMonitor, nil)
				svc.orm.On("CreateSpec", specFluxMonitor, mock.Anything).Return(int64(100), nil)
				svc.orm.On("CountJobProposalsByStatus").Return(&feeds.JobProposalCounts{}, nil)
			},
			args:   argsFluxMonitor,
			wantID: idFluxMonitor,
		},
		{
			name:    "contains invalid job spec",
			args:    &feeds.ProposeJobArgs{},
			wantErr: "invalid job type",
		},
		{
			name:   "must be an ocr job to include bootstraps",
			before: func(svc *TestService) {},
			args: &feeds.ProposeJobArgs{
				Spec:       FluxMonitorTestSpec,
				Multiaddrs: pq.StringArray{"/dns4/example.com"},
			},
			wantErr: "only OCR job type supports multiaddr",
		},
		{
			name: "ensure an upsert validates the job proposal belongs to the feeds manager",
			before: func(svc *TestService) {
				svc.orm.
					On("GetJobProposalByRemoteUUID", jpFluxMonitor.RemoteUUID).
					Return(&feeds.JobProposal{
						FeedsManagerID: 2,
						RemoteUUID:     jpFluxMonitor.RemoteUUID,
					}, nil)
			},
			args:    argsFluxMonitor,
			wantErr: "cannot update a job proposal belonging to another feeds manager",
		},
		{
			name: "spec version already exists",
			before: func(svc *TestService) {
				svc.orm.
					On("GetJobProposalByRemoteUUID", jpFluxMonitor.RemoteUUID).
					Return(&feeds.JobProposal{
						FeedsManagerID: jpFluxMonitor.FeedsManagerID,
						RemoteUUID:     jpFluxMonitor.RemoteUUID,
						Status:         feeds.JobProposalStatusPending,
					}, nil)
				svc.orm.On("ExistsSpecByJobProposalIDAndVersion", jpFluxMonitor.ID, argsFluxMonitor.Version).Return(true, nil)
			},
			args:    argsFluxMonitor,
			wantErr: "proposed job spec version already exists",
		},
		{
			name: "upsert error",
			before: func(svc *TestService) {
				svc.orm.On("GetJobProposalByRemoteUUID", jpFluxMonitor.RemoteUUID).Return(new(feeds.JobProposal), sql.ErrNoRows)
				svc.orm.On("UpsertJobProposal", &jpFluxMonitor, mock.Anything).Return(int64(0), errors.New("orm error"))
			},
			args:    argsFluxMonitor,
			wantErr: "failed to upsert job proposal",
		},
		{
			name: "Create spec error",
			before: func(svc *TestService) {
				svc.orm.On("GetJobProposalByRemoteUUID", jpFluxMonitor.RemoteUUID).Return(new(feeds.JobProposal), sql.ErrNoRows)
				svc.orm.On("UpsertJobProposal", &jpFluxMonitor, mock.Anything).Return(idFluxMonitor, nil)
				svc.orm.On("CreateSpec", specFluxMonitor, mock.Anything).Return(int64(0), errors.New("orm error"))
			},
			args:    argsFluxMonitor,
			wantErr: "failed to create spec",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := setupTestServiceCfg(t, func(c *chainlink.Config, s *chainlink.Secrets) {
				c.JobPipeline.HTTPRequest.DefaultTimeout = &httpTimeout
			})
			if tc.before != nil {
				tc.before(svc)
			}

			actual, err := svc.ProposeJob(testutils.Context(t), tc.args)

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

func Test_Service_DeleteJob(t *testing.T) {
	t.Parallel()

	var (
		remoteUUID = uuid.NewV4()
		args       = &feeds.DeleteJobArgs{
			FeedsManagerID: 1,
			RemoteUUID:     remoteUUID,
		}

		approved = feeds.JobProposal{
			ID:             1,
			FeedsManagerID: 1,
			RemoteUUID:     remoteUUID,
			Status:         feeds.JobProposalStatusApproved,
		}

		httpTimeout = models.MustMakeDuration(1 * time.Second)
	)

	testCases := []struct {
		name    string
		args    *feeds.DeleteJobArgs
		before  func(svc *TestService)
		wantID  int64
		wantErr string
	}{
		{
			name: "Delete success",
			before: func(svc *TestService) {
				svc.orm.On("GetJobProposalByRemoteUUID", approved.RemoteUUID).Return(&approved, nil)
				svc.orm.On("DeleteProposal", approved.ID, mock.Anything).Return(nil)
			},
			args:   args,
			wantID: approved.ID,
		},
		{
			name: "Job proposal being deleted belongs to the feeds manager",
			before: func(svc *TestService) {
				svc.orm.
					On("GetJobProposalByRemoteUUID", approved.RemoteUUID).
					Return(&feeds.JobProposal{
						FeedsManagerID: 2,
						RemoteUUID:     approved.RemoteUUID,
						Status:         feeds.JobProposalStatusApproved,
					}, nil)
			},
			args:    args,
			wantErr: "cannot delete a job proposal belonging to another feeds manager",
		},
		{
			name: "Get proposal error",
			before: func(svc *TestService) {
				svc.orm.On("GetJobProposalByRemoteUUID", approved.RemoteUUID).Return(nil, errors.New("orm error"))
			},
			args:    args,
			wantErr: "GetJobProposalByRemoteUUID failed",
		},
		{
			name: "No proposal error",
			before: func(svc *TestService) {
				svc.orm.On("GetJobProposalByRemoteUUID", approved.RemoteUUID).Return(nil, sql.ErrNoRows)
			},
			args:    args,
			wantErr: "GetJobProposalByRemoteUUID did not find any proposals to delete",
		},
		{
			name: "Delete proposal error",
			before: func(svc *TestService) {
				svc.orm.On("GetJobProposalByRemoteUUID", approved.RemoteUUID).Return(&approved, nil)
				svc.orm.On("DeleteProposal", approved.ID, mock.Anything).Return(errors.New("orm error"))
			},
			args:    args,
			wantErr: "DeleteProposal failed",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := setupTestServiceCfg(t, func(c *chainlink.Config, s *chainlink.Secrets) {
				c.JobPipeline.HTTPRequest.DefaultTimeout = &httpTimeout
			})
			if tc.before != nil {
				tc.before(svc)
			}

			_, err := svc.DeleteJob(testutils.Context(t), tc.args)

			if tc.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_Service_RevokeJob(t *testing.T) {
	t.Parallel()

	var (
		remoteUUID = uuid.NewV4()
		args       = &feeds.RevokeJobArgs{
			FeedsManagerID: 1,
			RemoteUUID:     remoteUUID,
		}

		defn = `
name = 'LINK / ETH | version 3 | contract 0x0000000000000000000000000000000000000000'
type               = "offchainreporting2"
pluginType         = "median"
schemaVersion      = 1
relay              = "evm"
contractID         = "0x613a38AC1659769640aaE063C651F48E0250454C"
externalJobID 		 = '00000000-0000-0000-0000-000000000001'
observationSource  = """
// data source 1
ds1 [type=bridge name=\"bridge-api0\" requestData="{\\\"data\\": {\\\"from\\\":\\\"LINK\\\",\\\"to\\\":\\\"ETH\\\"}}"];
ds1_parse [type=jsonparse path="result"];
ds1_multiply [type=multiply times=1000000000000000000];
ds1 -> ds1_parse -> ds1_multiply -> answer1;

answer1 [type=median index=0];
"""
[relayConfig]
chainID = 0
[pluginConfig]
juelsPerFeeCoinSource = """
ds1          [type=bridge name=voter_turnout];
ds1_parse    [type=jsonparse path="one,two"];
ds1_multiply [type=multiply times=1.23];
ds1 -> ds1_parse -> ds1_multiply -> answer1;
answer1      [type=median index=0];
"""
`

		pendingProposal = &feeds.JobProposal{
			ID:             1,
			FeedsManagerID: 1,
			RemoteUUID:     remoteUUID,
			Status:         feeds.JobProposalStatusPending,
		}

		pendingSpec = &feeds.JobProposalSpec{
			ID:            20,
			Status:        feeds.SpecStatusPending,
			JobProposalID: pendingProposal.ID,
			Version:       1,
			Definition:    defn,
		}

		httpTimeout = models.MustMakeDuration(1 * time.Second)
	)

	testCases := []struct {
		name    string
		args    *feeds.RevokeJobArgs
		before  func(svc *TestService)
		wantID  int64
		wantErr string
	}{
		{
			name: "Revoke success",
			before: func(svc *TestService) {
				svc.orm.On("GetJobProposalByRemoteUUID", pendingProposal.RemoteUUID).Return(pendingProposal, nil)
				svc.orm.On("GetLatestSpec", pendingSpec.JobProposalID).Return(pendingSpec, nil)
				svc.orm.On("RevokeSpec", pendingSpec.ID, mock.Anything).Return(nil)
			},
			args:   args,
			wantID: pendingProposal.ID,
		},
		{
			name: "Job proposal being deleted belongs to the feeds manager",
			before: func(svc *TestService) {
				svc.orm.
					On("GetJobProposalByRemoteUUID", pendingProposal.RemoteUUID).
					Return(&feeds.JobProposal{
						FeedsManagerID: 2,
						RemoteUUID:     pendingProposal.RemoteUUID,
						Status:         feeds.JobProposalStatusApproved,
					}, nil)
			},
			args:    args,
			wantErr: "cannot revoke a job proposal belonging to another feeds manager",
		},
		{
			name: "Get proposal error",
			before: func(svc *TestService) {
				svc.orm.On("GetJobProposalByRemoteUUID", pendingProposal.RemoteUUID).Return(nil, errors.New("orm error"))
			},
			args:    args,
			wantErr: "GetJobProposalByRemoteUUID failed",
		},
		{
			name: "No proposal error",
			before: func(svc *TestService) {
				svc.orm.On("GetJobProposalByRemoteUUID", pendingProposal.RemoteUUID).Return(nil, sql.ErrNoRows)
			},
			args:    args,
			wantErr: "GetJobProposalByRemoteUUID did not find any proposals to revoke",
		},
		{
			name: "Get latest spec error",
			before: func(svc *TestService) {
				svc.orm.On("GetJobProposalByRemoteUUID", pendingProposal.RemoteUUID).Return(pendingProposal, nil)
				svc.orm.On("GetLatestSpec", pendingSpec.JobProposalID).Return(nil, sql.ErrNoRows)
			},
			args:    args,
			wantErr: "GetLatestSpec failed to get latest spec",
		},
		{
			name: "Not revokable error",
			before: func(svc *TestService) {
				svc.orm.On("GetJobProposalByRemoteUUID", pendingProposal.RemoteUUID).Return(pendingProposal, nil)
				svc.orm.On("GetLatestSpec", pendingSpec.JobProposalID).Return(&feeds.JobProposalSpec{
					ID:            20,
					Status:        feeds.SpecStatusApproved,
					JobProposalID: pendingProposal.ID,
					Version:       1,
					Definition:    defn,
				}, nil)
			},
			args:    args,
			wantErr: "only pending job proposals can be revoked",
		},
		{
			name: "Revoke proposal error",
			before: func(svc *TestService) {
				svc.orm.On("GetJobProposalByRemoteUUID", pendingProposal.RemoteUUID).Return(pendingProposal, nil)
				svc.orm.On("GetLatestSpec", pendingSpec.JobProposalID).Return(pendingSpec, nil)
				svc.orm.On("RevokeSpec", pendingSpec.ID, mock.Anything).Return(errors.New("orm error"))
			},
			args:    args,
			wantErr: "RevokeSpec failed",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := setupTestServiceCfg(t, func(c *chainlink.Config, s *chainlink.Secrets) {
				c.JobPipeline.HTTPRequest.DefaultTimeout = &httpTimeout
			})
			if tc.before != nil {
				tc.before(svc)
			}

			_, err := svc.RevokeJob(testutils.Context(t), tc.args)

			if tc.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_Service_SyncNodeInfo(t *testing.T) {
	p2pKey, err := p2pkey.NewV2()
	require.NoError(t, err)

	ocrKey, err := ocrkey.NewV2()
	require.NoError(t, err)

	var (
		multiaddr = "/dns4/chain.link/tcp/1234/p2p/16Uiu2HAm58SP7UL8zsnpeuwHfytLocaqgnyaYKP8wu7qRdrixLju"
		mgr       = &feeds.FeedsManager{ID: 1}
		ccfg      = feeds.ChainConfig{
			ID:             100,
			FeedsManagerID: mgr.ID,
			ChainID:        "42",
			ChainType:      feeds.ChainTypeEVM,
			AccountAddress: "0x0000",
			AdminAddress:   "0x0001",
			FluxMonitorConfig: feeds.FluxMonitorConfig{
				Enabled: true,
			},
			OCR1Config: feeds.OCR1Config{
				Enabled:     true,
				IsBootstrap: false,
				P2PPeerID:   null.StringFrom(p2pKey.PeerID().String()),
				KeyBundleID: null.StringFrom(ocrKey.GetID()),
			},
			OCR2Config: feeds.OCR2Config{
				Enabled:     true,
				IsBootstrap: true,
				Multiaddr:   null.StringFrom(multiaddr),
				Plugins: feeds.Plugins{
					Commit:  true,
					Execute: true,
					Median:  false,
					Mercury: true,
				},
			},
		}
		chainConfigs = []feeds.ChainConfig{ccfg}
		nodeVersion  = &versioning.NodeVersion{Version: "1.0.0"}
	)

	svc := setupTestService(t)

	svc.connMgr.On("GetClient", mgr.ID).Return(svc.fmsClient, nil)
	svc.orm.On("ListChainConfigsByManagerIDs", []int64{mgr.ID}).Return(chainConfigs, nil)

	// OCR1 key fetching
	svc.p2pKeystore.On("Get", p2pKey.PeerID()).Return(p2pKey, nil)
	svc.ocr1Keystore.On("Get", ocrKey.GetID()).Return(ocrKey, nil)

	svc.fmsClient.On("UpdateNode", mock.Anything, &proto.UpdateNodeRequest{
		Version: nodeVersion.Version,
		ChainConfigs: []*proto.ChainConfig{
			{
				Chain: &proto.Chain{
					Id:   ccfg.ChainID,
					Type: proto.ChainType_CHAIN_TYPE_EVM,
				},
				AccountAddress:    ccfg.AccountAddress,
				AdminAddress:      ccfg.AdminAddress,
				FluxMonitorConfig: &proto.FluxMonitorConfig{Enabled: true},
				Ocr1Config: &proto.OCR1Config{
					Enabled:     true,
					IsBootstrap: ccfg.OCR1Config.IsBootstrap,
					P2PKeyBundle: &proto.OCR1Config_P2PKeyBundle{
						PeerId:    p2pKey.PeerID().String(),
						PublicKey: p2pKey.PublicKeyHex(),
					},
					OcrKeyBundle: &proto.OCR1Config_OCRKeyBundle{
						BundleId:              ocrKey.GetID(),
						ConfigPublicKey:       ocrkey.ConfigPublicKey(ocrKey.PublicKeyConfig()).String(),
						OffchainPublicKey:     ocrKey.OffChainSigning.PublicKey().String(),
						OnchainSigningAddress: ocrKey.OnChainSigning.Address().String(),
					},
				},
				Ocr2Config: &proto.OCR2Config{
					Enabled:     true,
					IsBootstrap: ccfg.OCR2Config.IsBootstrap,
					Multiaddr:   multiaddr,
					Plugins: &proto.OCR2Config_Plugins{
						Commit:  ccfg.OCR2Config.Plugins.Commit,
						Execute: ccfg.OCR2Config.Plugins.Execute,
						Median:  ccfg.OCR2Config.Plugins.Median,
						Mercury: ccfg.OCR2Config.Plugins.Mercury,
					},
				},
			},
		},
	}).Return(&proto.UpdateNodeResponse{}, nil)

	err = svc.SyncNodeInfo(testutils.Context(t), mgr.ID)
	require.NoError(t, err)
}

func Test_Service_IsJobManaged(t *testing.T) {
	t.Parallel()

	svc := setupTestService(t)
	ctx := testutils.Context(t)
	jobID := int64(1)

	svc.orm.On("IsJobManaged", jobID, mock.Anything).Return(true, nil)

	isManaged, err := svc.IsJobManaged(ctx, jobID)
	require.NoError(t, err)
	assert.True(t, isManaged)
}

func Test_Service_ListJobProposals(t *testing.T) {
	t.Parallel()

	var (
		jp  = feeds.JobProposal{}
		jps = []feeds.JobProposal{jp}
	)
	svc := setupTestService(t)

	svc.orm.On("ListJobProposals").
		Return(jps, nil)

	actual, err := svc.ListJobProposals()
	require.NoError(t, err)

	assert.Equal(t, actual, jps)
}

func Test_Service_ListJobProposalsByManagersIDs(t *testing.T) {
	t.Parallel()

	var (
		jp    = feeds.JobProposal{}
		jps   = []feeds.JobProposal{jp}
		fmIDs = []int64{1}
	)
	svc := setupTestService(t)

	svc.orm.On("ListJobProposalsByManagersIDs", fmIDs).
		Return(jps, nil)

	actual, err := svc.ListJobProposalsByManagersIDs(fmIDs)
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

	svc.orm.On("GetJobProposal", id).
		Return(&ms, nil)

	actual, err := svc.GetJobProposal(id)
	require.NoError(t, err)

	assert.Equal(t, actual, &ms)
}

func Test_Service_CancelSpec(t *testing.T) {
	var (
		externalJobID = uuid.NewV4()
		jp            = &feeds.JobProposal{
			ID:             1,
			ExternalJobID:  uuid.NullUUID{UUID: externalJobID, Valid: true},
			RemoteUUID:     externalJobID,
			FeedsManagerID: 100,
		}
		spec = &feeds.JobProposalSpec{
			ID:            20,
			Status:        feeds.SpecStatusApproved,
			JobProposalID: jp.ID,
			Version:       1,
		}
		j = job.Job{
			ID:            1,
			ExternalJobID: externalJobID,
		}
	)

	testCases := []struct {
		name    string
		before  func(svc *TestService)
		specID  int64
		wantErr string
	}{
		{
			name: "success",
			before: func(svc *TestService) {
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.orm.On("GetSpec", spec.ID, mock.Anything).Return(spec, nil)
				svc.orm.On("GetJobProposal", jp.ID, mock.Anything).Return(jp, nil)

				svc.orm.On("CancelSpec",
					spec.ID,
					mock.Anything,
				).Return(nil)
				svc.jobORM.On("FindJobByExternalJobID", externalJobID, mock.Anything).Return(j, nil)
				svc.spawner.On("DeleteJob", j.ID, mock.Anything).Return(nil)

				svc.fmsClient.On("CancelledJob",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&proto.CancelledJobRequest{
						Uuid:    jp.RemoteUUID.String(),
						Version: int64(spec.Version),
					},
				).Return(&proto.CancelledJobResponse{}, nil)
				svc.orm.On("CountJobProposalsByStatus").Return(&feeds.JobProposalCounts{}, nil)
			},
			specID: spec.ID,
		},
		{
			name: "spec does not exist",
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", spec.ID, mock.Anything).Return(nil, errors.New("Not Found"))
			},
			specID:  spec.ID,
			wantErr: "orm: job proposal spec: Not Found",
		},
		{
			name: "must be an approved job proposal spec",
			before: func(svc *TestService) {
				pspec := &feeds.JobProposalSpec{
					ID:     spec.ID,
					Status: feeds.SpecStatusPending,
				}
				svc.orm.On("GetSpec", pspec.ID, mock.Anything).Return(pspec, nil)
			},
			specID:  spec.ID,
			wantErr: "must be an approved job proposal spec",
		},
		{
			name: "job proposal does not exist",
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", spec.ID, mock.Anything).Return(spec, nil)
				svc.orm.On("GetJobProposal", jp.ID, mock.Anything).Return(nil, errors.New("Not Found"))
			},
			specID:  spec.ID,
			wantErr: "orm: job proposal: Not Found",
		},
		{
			name: "rpc client not connected",
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", spec.ID, mock.Anything).Return(spec, nil)
				svc.orm.On("GetJobProposal", jp.ID, mock.Anything).Return(jp, nil)
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(nil, errors.New("Not Connected"))
			},
			specID:  spec.ID,
			wantErr: "fms rpc client: Not Connected",
		},
		{
			name: "cancel spec orm fails",
			before: func(svc *TestService) {
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.orm.On("GetSpec", spec.ID, mock.Anything).Return(spec, nil)
				svc.orm.On("GetJobProposal", jp.ID, mock.Anything).Return(jp, nil)

				svc.orm.On("CancelSpec",
					spec.ID,
					mock.Anything,
				).Return(errors.New("failure"))
			},
			specID:  spec.ID,
			wantErr: "failure",
		},
		{
			name: "find by external uuid orm fails",
			before: func(svc *TestService) {
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.orm.On("GetSpec", spec.ID, mock.Anything).Return(spec, nil)
				svc.orm.On("GetJobProposal", jp.ID, mock.Anything).Return(jp, nil)

				svc.orm.On("CancelSpec",
					spec.ID,
					mock.Anything,
				).Return(nil)
				svc.jobORM.On("FindJobByExternalJobID", externalJobID, mock.Anything).Return(job.Job{}, errors.New("failure"))
			},
			specID:  spec.ID,
			wantErr: "FindJobByExternalJobID failed: failure",
		},
		{
			name: "delete job fails",
			before: func(svc *TestService) {
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.orm.On("GetSpec", spec.ID, mock.Anything).Return(spec, nil)
				svc.orm.On("GetJobProposal", jp.ID, mock.Anything).Return(jp, nil)

				svc.orm.On("CancelSpec",
					spec.ID,
					mock.Anything,
				).Return(nil)
				svc.jobORM.On("FindJobByExternalJobID", externalJobID, mock.Anything).Return(j, nil)
				svc.spawner.On("DeleteJob", j.ID, mock.Anything).Return(errors.New("failure"))
			},
			specID:  spec.ID,
			wantErr: "DeleteJob failed: failure",
		},
		{
			name: "cancelled job rpc call fails",
			before: func(svc *TestService) {
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.orm.On("GetSpec", spec.ID, mock.Anything).Return(spec, nil)
				svc.orm.On("GetJobProposal", jp.ID, mock.Anything).Return(jp, nil)

				svc.orm.On("CancelSpec",
					spec.ID,
					mock.Anything,
				).Return(nil)
				svc.jobORM.On("FindJobByExternalJobID", externalJobID, mock.Anything).Return(j, nil)
				svc.spawner.On("DeleteJob", j.ID, mock.Anything).Return(nil)

				svc.fmsClient.On("CancelledJob",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&proto.CancelledJobRequest{
						Uuid:    jp.RemoteUUID.String(),
						Version: int64(spec.Version),
					},
				).Return(nil, errors.New("failure"))
			},
			specID:  spec.ID,
			wantErr: "failure",
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

			err := svc.CancelSpec(testutils.Context(t), tc.specID)

			if tc.wantErr != "" {
				require.Error(t, err)
				assert.EqualError(t, err, tc.wantErr)

				return
			}

			require.NoError(t, err)
		})
	}
}

func Test_Service_GetSpec(t *testing.T) {
	t.Parallel()

	var (
		id   = int64(1)
		spec = feeds.JobProposalSpec{ID: id}
	)
	svc := setupTestService(t)

	svc.orm.On("GetSpec", id).
		Return(&spec, nil)

	actual, err := svc.GetSpec(id)
	require.NoError(t, err)

	assert.Equal(t, &spec, actual)
}

func Test_Service_ListSpecsByJobProposalIDs(t *testing.T) {
	t.Parallel()

	var (
		id    = int64(1)
		jpID  = int64(200)
		spec  = feeds.JobProposalSpec{ID: id, JobProposalID: jpID}
		specs = []feeds.JobProposalSpec{spec}
	)
	svc := setupTestService(t)

	svc.orm.On("ListSpecsByJobProposalIDs", []int64{jpID}).
		Return(specs, nil)

	actual, err := svc.ListSpecsByJobProposalIDs([]int64{jpID})
	require.NoError(t, err)

	assert.Equal(t, specs, actual)
}

func Test_Service_ApproveSpec(t *testing.T) {
	var evmChainID *utils.Big
	address := ethkey.EIP55AddressFromAddress(common.Address{})

	var (
		ctx  = testutils.Context(t)
		defn = `
name = 'LINK / ETH | version 3 | contract 0x0000000000000000000000000000000000000000'
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

		jp = &feeds.JobProposal{
			ID:             1,
			FeedsManagerID: 100,
		}
		spec = &feeds.JobProposalSpec{
			ID:            20,
			Status:        feeds.SpecStatusPending,
			JobProposalID: jp.ID,
			Version:       1,
			Definition:    defn,
		}
		rejectedSpec = &feeds.JobProposalSpec{
			ID:            20,
			Status:        feeds.SpecStatusRejected,
			JobProposalID: jp.ID,
			Version:       1,
			Definition:    defn,
		}
		cancelledSpec = &feeds.JobProposalSpec{
			ID:            20,
			Status:        feeds.SpecStatusCancelled,
			JobProposalID: jp.ID,
			Version:       1,
			Definition:    defn,
		}
		j = job.Job{
			ID:            1,
			ExternalJobID: uuid.NewV4(),
		}
	)

	testCases := []struct {
		name        string
		httpTimeout *models.Duration
		before      func(svc *TestService)
		id          int64
		force       bool
		wantErr     string
	}{
		{
			name:        "pending job success",
			httpTimeout: models.MustNewDuration(1 * time.Minute),
			before: func(svc *TestService) {
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.orm.EXPECT().GetSpec(spec.ID, mock.Anything).Return(spec, nil)
				svc.orm.EXPECT().GetJobProposal(jp.ID, mock.Anything).Return(jp, nil)
				svc.jobORM.On("AssertBridgesExist", mock.IsType(pipeline.Pipeline{})).Return(nil)

				svc.jobORM.On("FindJobIDByAddress", address, evmChainID, mock.Anything).Return(int32(0), sql.ErrNoRows)

				svc.spawner.
					On("CreateJob",
						mock.MatchedBy(func(j *job.Job) bool {
							return j.Name.String == "LINK / ETH | version 3 | contract 0x0000000000000000000000000000000000000000"
						}),
						mock.Anything,
					).
					Run(func(args mock.Arguments) { (args.Get(0).(*job.Job)).ID = 1 }).
					Return(nil)
				svc.orm.On("ApproveSpec",
					spec.ID,
					uuid.Must(uuid.FromString("00000000-0000-0000-0000-000000000001")),
					mock.Anything,
				).Return(nil)
				svc.fmsClient.On("ApprovedJob",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&proto.ApprovedJobRequest{
						Uuid:    jp.RemoteUUID.String(),
						Version: int64(spec.Version),
					},
				).Return(&proto.ApprovedJobResponse{}, nil)
				svc.orm.On("CountJobProposalsByStatus").Return(&feeds.JobProposalCounts{}, nil)
			},
			id:    spec.ID,
			force: false,
		},
		{
			name:        "cancelled spec success when it is the latest spec",
			httpTimeout: models.MustNewDuration(1 * time.Minute),
			before: func(svc *TestService) {
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.orm.On("GetSpec", cancelledSpec.ID, mock.Anything).Return(cancelledSpec, nil)
				svc.orm.On("GetJobProposal", jp.ID, mock.Anything).Return(jp, nil)
				svc.orm.On("GetLatestSpec", cancelledSpec.JobProposalID).Return(cancelledSpec, nil)
				svc.jobORM.On("AssertBridgesExist", mock.IsType(pipeline.Pipeline{})).Return(nil)

				svc.jobORM.On("FindJobIDByAddress", address, evmChainID, mock.Anything).Return(int32(0), sql.ErrNoRows)

				svc.spawner.
					On("CreateJob",
						mock.MatchedBy(func(j *job.Job) bool {
							return j.Name.String == "LINK / ETH | version 3 | contract 0x0000000000000000000000000000000000000000"
						}),
						mock.Anything,
					).
					Run(func(args mock.Arguments) { (args.Get(0).(*job.Job)).ID = 1 }).
					Return(nil)
				svc.orm.On("ApproveSpec",
					cancelledSpec.ID,
					uuid.Must(uuid.FromString("00000000-0000-0000-0000-000000000001")),
					mock.Anything,
				).Return(nil)
				svc.fmsClient.On("ApprovedJob",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&proto.ApprovedJobRequest{
						Uuid:    jp.RemoteUUID.String(),
						Version: int64(spec.Version),
					},
				).Return(&proto.ApprovedJobResponse{}, nil)
				svc.orm.On("CountJobProposalsByStatus").Return(&feeds.JobProposalCounts{}, nil)
			},
			id:    cancelledSpec.ID,
			force: false,
		},
		{
			name: "revoked proposal fail",
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", spec.ID, mock.Anything).Return(spec, nil)
				svc.orm.On("GetJobProposal", jp.ID, mock.Anything).Return(&feeds.JobProposal{
					ID:     1,
					Status: feeds.JobProposalStatusRevoked,
				}, nil)
			},
			id:      spec.ID,
			force:   false,
			wantErr: "cannot approve spec for a revoked job proposal",
		},
		{
			name: "deleted proposal fail",
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", spec.ID, mock.Anything).Return(spec, nil)
				svc.orm.On("GetJobProposal", jp.ID, mock.Anything).Return(&feeds.JobProposal{
					ID:     jp.ID,
					Status: feeds.JobProposalStatusDeleted,
				}, nil)
			},
			id:      spec.ID,
			force:   false,
			wantErr: "cannot approve spec for a deleted job proposal",
		},
		{
			name: "approved spec fail",
			before: func(svc *TestService) {
				aspec := &feeds.JobProposalSpec{
					ID:            spec.ID,
					Status:        feeds.SpecStatusApproved,
					JobProposalID: jp.ID,
				}
				svc.orm.On("GetSpec", aspec.ID, mock.Anything).Return(aspec, nil)
				svc.orm.On("GetJobProposal", jp.ID, mock.Anything).Return(jp, nil)
			},
			id:      spec.ID,
			force:   false,
			wantErr: "cannot approve an approved spec",
		},
		{
			name: "rejected spec fail",
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", cancelledSpec.ID, mock.Anything).Return(rejectedSpec, nil)
				svc.orm.On("GetJobProposal", jp.ID, mock.Anything).Return(jp, nil)
			},
			id:      rejectedSpec.ID,
			force:   false,
			wantErr: "cannot approve a rejected spec",
		},
		{
			name: "cancelled spec failed not latest spec",
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", cancelledSpec.ID, mock.Anything).Return(cancelledSpec, nil)
				svc.orm.On("GetJobProposal", jp.ID, mock.Anything).Return(jp, nil)
				svc.orm.On("GetLatestSpec", cancelledSpec.JobProposalID).Return(&feeds.JobProposalSpec{
					ID:            21,
					Status:        feeds.SpecStatusPending,
					JobProposalID: jp.ID,
					Version:       2,
					Definition:    defn,
				}, nil)
			},
			id:      cancelledSpec.ID,
			force:   false,
			wantErr: "cannot approve a cancelled spec",
		},
		{
			name:        "already existing job replacement error",
			httpTimeout: models.MustNewDuration(1 * time.Minute),
			before: func(svc *TestService) {
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.orm.On("GetSpec", spec.ID, mock.Anything).Return(spec, nil)
				svc.orm.On("GetJobProposal", jp.ID, mock.Anything).Return(jp, nil)
				svc.jobORM.On("AssertBridgesExist", mock.IsType(pipeline.Pipeline{})).Return(nil)

				svc.jobORM.On("FindJobIDByAddress", address, evmChainID, mock.Anything).Return(j.ID, nil)
			},
			id:      spec.ID,
			force:   false,
			wantErr: "could not approve job proposal: a job for this contract address already exists - please use the 'force' option to replace it",
		},
		{
			name:        "already existing self managed job replacement success if forced",
			httpTimeout: models.MustNewDuration(1 * time.Minute),
			before: func(svc *TestService) {
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.orm.EXPECT().GetSpec(spec.ID, mock.Anything).Return(spec, nil)
				svc.orm.EXPECT().GetJobProposal(jp.ID, mock.Anything).Return(jp, nil)
				svc.jobORM.On("AssertBridgesExist", mock.IsType(pipeline.Pipeline{})).Return(nil)
				svc.jobORM.On("FindJobIDByAddress", address, evmChainID, mock.Anything).Return(j.ID, nil)
				svc.orm.EXPECT().GetApprovedSpec(jp.ID, mock.Anything).Return(nil, sql.ErrNoRows)

				svc.spawner.On("DeleteJob", j.ID, mock.Anything).Return(nil)

				svc.spawner.
					On("CreateJob",
						mock.MatchedBy(func(j *job.Job) bool {
							return j.Name.String == "LINK / ETH | version 3 | contract 0x0000000000000000000000000000000000000000"
						}),
						mock.Anything,
					).
					Run(func(args mock.Arguments) { (args.Get(0).(*job.Job)).ID = 1 }).
					Return(nil)
				svc.orm.On("ApproveSpec",
					spec.ID,
					uuid.Must(uuid.FromString("00000000-0000-0000-0000-000000000001")),
					mock.Anything,
				).Return(nil)
				svc.fmsClient.On("ApprovedJob",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&proto.ApprovedJobRequest{
						Uuid:    jp.RemoteUUID.String(),
						Version: int64(spec.Version),
					},
				).Return(&proto.ApprovedJobResponse{}, nil)
				svc.orm.On("CountJobProposalsByStatus").Return(&feeds.JobProposalCounts{}, nil)
			},
			id:    spec.ID,
			force: true,
		},
		{
			name:        "already existing FMS managed job replacement success if forced",
			httpTimeout: models.MustNewDuration(1 * time.Minute),
			before: func(svc *TestService) {
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.orm.EXPECT().GetSpec(spec.ID, mock.Anything).Return(spec, nil)
				svc.orm.EXPECT().GetJobProposal(jp.ID, mock.Anything).Return(jp, nil)
				svc.jobORM.On("AssertBridgesExist", mock.IsType(pipeline.Pipeline{})).Return(nil)
				svc.jobORM.On("FindJobIDByAddress", address, evmChainID, mock.Anything).Return(j.ID, nil)
				svc.orm.EXPECT().GetApprovedSpec(jp.ID, mock.Anything).Return(&feeds.JobProposalSpec{ID: 100}, nil)
				svc.orm.EXPECT().CancelSpec(int64(100), mock.Anything).Return(nil)

				svc.spawner.On("DeleteJob", j.ID, mock.Anything).Return(nil)

				svc.spawner.
					On("CreateJob",
						mock.MatchedBy(func(j *job.Job) bool {
							return j.Name.String == "LINK / ETH | version 3 | contract 0x0000000000000000000000000000000000000000"
						}),
						mock.Anything,
					).
					Run(func(args mock.Arguments) { (args.Get(0).(*job.Job)).ID = 1 }).
					Return(nil)
				svc.orm.On("ApproveSpec",
					spec.ID,
					uuid.Must(uuid.FromString("00000000-0000-0000-0000-000000000001")),
					mock.Anything,
				).Return(nil)
				svc.fmsClient.On("ApprovedJob",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&proto.ApprovedJobRequest{
						Uuid:    jp.RemoteUUID.String(),
						Version: int64(spec.Version),
					},
				).Return(&proto.ApprovedJobResponse{}, nil)
				svc.orm.On("CountJobProposalsByStatus").Return(&feeds.JobProposalCounts{}, nil)
			},
			id:    spec.ID,
			force: true,
		},
		{
			name: "spec does not exist",
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", spec.ID, mock.Anything).Return(nil, errors.New("Not Found"))
			},
			id:      spec.ID,
			force:   false,
			wantErr: "orm: job proposal spec: Not Found",
		},
		{
			name: "job proposal does not exist",
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", spec.ID, mock.Anything).Return(spec, nil)
				svc.orm.On("GetJobProposal", jp.ID, mock.Anything).Return(nil, errors.New("Not Found"))
			},
			id:      spec.ID,
			wantErr: "orm: job proposal: Not Found",
		},
		{
			name:        "bridges do not exist",
			httpTimeout: models.MustNewDuration(1 * time.Minute),
			before: func(svc *TestService) {
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.orm.On("GetSpec", spec.ID, mock.Anything).Return(spec, nil)
				svc.orm.On("GetJobProposal", jp.ID, mock.Anything).Return(jp, nil)
				svc.jobORM.On("AssertBridgesExist", mock.IsType(pipeline.Pipeline{})).Return(errors.New("bridges do not exist"))
			},
			id:      spec.ID,
			wantErr: "failed to approve job spec due to bridge check: bridges do not exist",
		},
		{
			name: "rpc client not connected",
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", spec.ID, mock.Anything).Return(spec, nil)
				svc.orm.On("GetJobProposal", jp.ID, mock.Anything).Return(jp, nil)
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(nil, errors.New("Not Connected"))
			},
			id:      spec.ID,
			force:   false,
			wantErr: "fms rpc client: Not Connected",
		},
		{
			name:        "Fetching the approved spec fails",
			httpTimeout: models.MustNewDuration(1 * time.Minute),
			before: func(svc *TestService) {
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.orm.EXPECT().GetSpec(spec.ID, mock.Anything).Return(spec, nil)
				svc.orm.EXPECT().GetJobProposal(jp.ID, mock.Anything).Return(jp, nil)
				svc.jobORM.On("AssertBridgesExist", mock.IsType(pipeline.Pipeline{})).Return(nil)
				svc.jobORM.On("FindJobIDByAddress", address, evmChainID, mock.Anything).Return(j.ID, nil)
				svc.orm.EXPECT().GetApprovedSpec(jp.ID, mock.Anything).Return(nil, errors.New("failure"))
			},
			id:      spec.ID,
			force:   true,
			wantErr: "could not approve job proposal: GetApprovedSpec failed: failure",
		},
		{
			name:        "spec cancellation fails",
			httpTimeout: models.MustNewDuration(1 * time.Minute),
			before: func(svc *TestService) {
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.orm.EXPECT().GetSpec(spec.ID, mock.Anything).Return(spec, nil)
				svc.orm.EXPECT().GetJobProposal(jp.ID, mock.Anything).Return(jp, nil)
				svc.jobORM.On("AssertBridgesExist", mock.IsType(pipeline.Pipeline{})).Return(nil)
				svc.jobORM.On("FindJobIDByAddress", address, evmChainID, mock.Anything).Return(j.ID, nil)
				svc.orm.EXPECT().GetApprovedSpec(jp.ID, mock.Anything).Return(&feeds.JobProposalSpec{ID: 100}, nil)
				svc.orm.EXPECT().CancelSpec(int64(100), mock.Anything).Return(errors.New("failure"))
			},
			id:      spec.ID,
			force:   true,
			wantErr: "could not approve job proposal: failure",
		},
		{
			name:        "create job error",
			httpTimeout: models.MustNewDuration(1 * time.Minute),
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", spec.ID, mock.Anything).Return(spec, nil)
				svc.orm.On("GetJobProposal", jp.ID, mock.Anything).Return(jp, nil)
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.jobORM.On("AssertBridgesExist", mock.IsType(pipeline.Pipeline{})).Return(nil)

				svc.jobORM.On("FindJobIDByAddress", address, evmChainID, mock.Anything).Return(int32(0), sql.ErrNoRows)

				svc.spawner.
					On("CreateJob",
						mock.MatchedBy(func(j *job.Job) bool {
							return j.Name.String == "LINK / ETH | version 3 | contract 0x0000000000000000000000000000000000000000"
						}),
						mock.Anything,
					).
					Return(errors.New("could not save"))
			},
			id:      spec.ID,
			force:   false,
			wantErr: "could not approve job proposal: could not save",
		},
		{
			name:        "approve spec orm error",
			httpTimeout: models.MustNewDuration(1 * time.Minute),
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", spec.ID, mock.Anything).Return(spec, nil)
				svc.orm.On("GetJobProposal", jp.ID, mock.Anything).Return(jp, nil)
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.jobORM.On("AssertBridgesExist", mock.IsType(pipeline.Pipeline{})).Return(nil)

				svc.jobORM.On("FindJobIDByAddress", address, evmChainID, mock.Anything).Return(int32(0), sql.ErrNoRows)

				svc.spawner.
					On("CreateJob",
						mock.MatchedBy(func(j *job.Job) bool {
							return j.Name.String == "LINK / ETH | version 3 | contract 0x0000000000000000000000000000000000000000"
						}),
						mock.Anything,
					).
					Run(func(args mock.Arguments) { (args.Get(0).(*job.Job)).ID = 1 }).
					Return(nil)
				svc.orm.On("ApproveSpec",
					spec.ID,
					uuid.Must(uuid.FromString("00000000-0000-0000-0000-000000000001")),
					mock.Anything,
				).Return(errors.New("failure"))
			},
			id:      spec.ID,
			force:   false,
			wantErr: "could not approve job proposal: failure",
		},
		{
			name:        "fms call error",
			httpTimeout: models.MustNewDuration(1 * time.Minute),
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", spec.ID, mock.Anything).Return(spec, nil)
				svc.orm.On("GetJobProposal", jp.ID, mock.Anything).Return(jp, nil)
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.jobORM.On("AssertBridgesExist", mock.IsType(pipeline.Pipeline{})).Return(nil)

				svc.jobORM.On("FindJobIDByAddress", address, evmChainID, mock.Anything).Return(int32(0), sql.ErrNoRows)

				svc.spawner.
					On("CreateJob",
						mock.MatchedBy(func(j *job.Job) bool {
							return j.Name.String == "LINK / ETH | version 3 | contract 0x0000000000000000000000000000000000000000"
						}),
						mock.Anything,
					).
					Run(func(args mock.Arguments) { (args.Get(0).(*job.Job)).ID = 1 }).
					Return(nil)
				svc.orm.On("ApproveSpec",
					spec.ID,
					uuid.Must(uuid.FromString("00000000-0000-0000-0000-000000000001")),
					mock.Anything,
				).Return(nil)
				svc.fmsClient.On("ApprovedJob",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&proto.ApprovedJobRequest{
						Uuid:    jp.RemoteUUID.String(),
						Version: int64(spec.Version),
					},
				).Return(nil, errors.New("failure"))
			},
			id:      spec.ID,
			force:   false,
			wantErr: "could not approve job proposal: failure",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			svc := setupTestServiceCfg(t, func(c *chainlink.Config, s *chainlink.Secrets) {
				if tc.httpTimeout != nil {
					c.JobPipeline.HTTPRequest.DefaultTimeout = tc.httpTimeout
				}
			})

			if tc.before != nil {
				tc.before(svc)
			}

			err := svc.ApproveSpec(ctx, tc.id, tc.force)

			if tc.wantErr != "" {
				require.Error(t, err)
				assert.EqualError(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_Service_ApproveSpec_OCR2(t *testing.T) {
	evmChainID := utils.NewBig(big.NewInt(0))
	address := ethkey.EIP55AddressFromAddress(common.HexToAddress("0x613a38AC1659769640aaE063C651F48E0250454C"))

	var (
		ctx  = testutils.Context(t)
		defn = `
name = 'LINK / ETH | version 3 | contract 0x0000000000000000000000000000000000000000'
type               = "offchainreporting2"
pluginType         = "median"
schemaVersion      = 1
relay              = "evm"
contractID         = "0x613a38AC1659769640aaE063C651F48E0250454C"
externalJobID 		 = '00000000-0000-0000-0000-000000000001'
observationSource  = """
// data source 1
ds1 [type=bridge name=\"bridge-api0\" requestData="{\\\"data\\": {\\\"from\\\":\\\"LINK\\\",\\\"to\\\":\\\"ETH\\\"}}"];
ds1_parse [type=jsonparse path="result"];
ds1_multiply [type=multiply times=1000000000000000000];
ds1 -> ds1_parse -> ds1_multiply -> answer1;

answer1 [type=median index=0];
"""
[relayConfig]
chainID = 0
[pluginConfig]
juelsPerFeeCoinSource = """
ds1          [type=bridge name=voter_turnout];
ds1_parse    [type=jsonparse path="one,two"];
ds1_multiply [type=multiply times=1.23];
ds1 -> ds1_parse -> ds1_multiply -> answer1;
answer1      [type=median index=0];
"""
`

		jp = &feeds.JobProposal{
			ID:             1,
			FeedsManagerID: 100,
		}
		spec = &feeds.JobProposalSpec{
			ID:            20,
			Status:        feeds.SpecStatusPending,
			JobProposalID: jp.ID,
			Version:       1,
			Definition:    defn,
		}
		rejectedSpec = &feeds.JobProposalSpec{
			ID:            20,
			Status:        feeds.SpecStatusRejected,
			JobProposalID: jp.ID,
			Version:       1,
			Definition:    defn,
		}
		cancelledSpec = &feeds.JobProposalSpec{
			ID:            20,
			Status:        feeds.SpecStatusCancelled,
			JobProposalID: jp.ID,
			Version:       1,
			Definition:    defn,
		}
		j = job.Job{
			ID:            1,
			ExternalJobID: uuid.NewV4(),
		}
	)

	testCases := []struct {
		name        string
		httpTimeout *models.Duration
		before      func(svc *TestService)
		id          int64
		force       bool
		wantErr     string
	}{
		{
			name:        "pending job success",
			httpTimeout: models.MustNewDuration(1 * time.Minute),
			before: func(svc *TestService) {
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.orm.On("GetSpec", spec.ID, mock.Anything).Return(spec, nil)
				svc.orm.On("GetJobProposal", jp.ID, mock.Anything).Return(jp, nil)
				svc.jobORM.On("AssertBridgesExist", mock.IsType(pipeline.Pipeline{})).Return(nil)
				svc.jobORM.On("FindJobIDByAddress", address, evmChainID, mock.Anything).Return(int32(0), sql.ErrNoRows)

				svc.spawner.
					On("CreateJob",
						mock.MatchedBy(func(j *job.Job) bool {
							return j.Name.String == "LINK / ETH | version 3 | contract 0x0000000000000000000000000000000000000000"
						}),
						mock.Anything,
					).
					Run(func(args mock.Arguments) { (args.Get(0).(*job.Job)).ID = 1 }).
					Return(nil)
				svc.orm.On("ApproveSpec",
					spec.ID,
					uuid.Must(uuid.FromString("00000000-0000-0000-0000-000000000001")),
					mock.Anything,
				).Return(nil)
				svc.fmsClient.On("ApprovedJob",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&proto.ApprovedJobRequest{
						Uuid:    jp.RemoteUUID.String(),
						Version: int64(spec.Version),
					},
				).Return(&proto.ApprovedJobResponse{}, nil)
				svc.orm.On("CountJobProposalsByStatus").Return(&feeds.JobProposalCounts{}, nil)
			},
			id:    spec.ID,
			force: false,
		},
		{
			name:        "cancelled spec success when it is the latest spec",
			httpTimeout: models.MustNewDuration(1 * time.Minute),
			before: func(svc *TestService) {
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.orm.On("GetSpec", cancelledSpec.ID, mock.Anything).Return(cancelledSpec, nil)
				svc.orm.On("GetJobProposal", jp.ID, mock.Anything).Return(jp, nil)
				svc.orm.On("GetLatestSpec", cancelledSpec.JobProposalID).Return(cancelledSpec, nil)
				svc.jobORM.On("AssertBridgesExist", mock.IsType(pipeline.Pipeline{})).Return(nil)

				svc.jobORM.On("FindJobIDByAddress", address, evmChainID, mock.Anything).Return(int32(0), sql.ErrNoRows)

				svc.spawner.
					On("CreateJob",
						mock.MatchedBy(func(j *job.Job) bool {
							return j.Name.String == "LINK / ETH | version 3 | contract 0x0000000000000000000000000000000000000000"
						}),
						mock.Anything,
					).
					Run(func(args mock.Arguments) { (args.Get(0).(*job.Job)).ID = 1 }).
					Return(nil)
				svc.orm.On("ApproveSpec",
					cancelledSpec.ID,
					uuid.Must(uuid.FromString("00000000-0000-0000-0000-000000000001")),
					mock.Anything,
				).Return(nil)
				svc.fmsClient.On("ApprovedJob",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&proto.ApprovedJobRequest{
						Uuid:    jp.RemoteUUID.String(),
						Version: int64(spec.Version),
					},
				).Return(&proto.ApprovedJobResponse{}, nil)
				svc.orm.On("CountJobProposalsByStatus").Return(&feeds.JobProposalCounts{}, nil)
			},
			id:    cancelledSpec.ID,
			force: false,
		},
		{
			name: "cancelled spec failed not latest spec",
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", cancelledSpec.ID, mock.Anything).Return(cancelledSpec, nil)
				svc.orm.On("GetJobProposal", jp.ID, mock.Anything).Return(jp, nil)
				svc.orm.On("GetLatestSpec", cancelledSpec.JobProposalID).Return(&feeds.JobProposalSpec{
					ID:            21,
					Status:        feeds.SpecStatusPending,
					JobProposalID: jp.ID,
					Version:       2,
					Definition:    defn,
				}, nil)
			},
			id:      cancelledSpec.ID,
			force:   false,
			wantErr: "cannot approve a cancelled spec",
		},
		{
			name: "rejected spec failed cannot be approved",
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", cancelledSpec.ID, mock.Anything).Return(rejectedSpec, nil)
				svc.orm.On("GetJobProposal", jp.ID, mock.Anything).Return(jp, nil)
			},
			id:      rejectedSpec.ID,
			force:   false,
			wantErr: "cannot approve a rejected spec",
		},
		{
			name:        "already existing job replacement error",
			httpTimeout: models.MustNewDuration(1 * time.Minute),
			before: func(svc *TestService) {
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.orm.On("GetSpec", spec.ID, mock.Anything).Return(spec, nil)
				svc.orm.On("GetJobProposal", jp.ID, mock.Anything).Return(jp, nil)
				svc.jobORM.On("AssertBridgesExist", mock.IsType(pipeline.Pipeline{})).Return(nil)
				svc.jobORM.On("FindJobIDByAddress", address, evmChainID, mock.Anything).Return(j.ID, nil)
			},
			id:      spec.ID,
			force:   false,
			wantErr: "could not approve job proposal: a job for this contract address already exists - please use the 'force' option to replace it",
		},
		{
			name:        "already existing self managed job replacement success if forced",
			httpTimeout: models.MustNewDuration(1 * time.Minute),
			before: func(svc *TestService) {
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.orm.On("GetSpec", spec.ID, mock.Anything).Return(spec, nil)
				svc.orm.On("GetJobProposal", jp.ID, mock.Anything).Return(jp, nil)
				svc.jobORM.On("AssertBridgesExist", mock.IsType(pipeline.Pipeline{})).Return(nil)
				svc.orm.EXPECT().GetApprovedSpec(jp.ID, mock.Anything).Return(nil, sql.ErrNoRows)
				svc.jobORM.On("FindJobIDByAddress", address, evmChainID, mock.Anything).Return(j.ID, nil)
				svc.spawner.On("DeleteJob", j.ID, mock.Anything).Return(nil)

				svc.spawner.
					On("CreateJob",
						mock.MatchedBy(func(j *job.Job) bool {
							return j.Name.String == "LINK / ETH | version 3 | contract 0x0000000000000000000000000000000000000000"
						}),
						mock.Anything,
					).
					Run(func(args mock.Arguments) { (args.Get(0).(*job.Job)).ID = 1 }).
					Return(nil)
				svc.orm.On("ApproveSpec",
					spec.ID,
					uuid.Must(uuid.FromString("00000000-0000-0000-0000-000000000001")),
					mock.Anything,
				).Return(nil)
				svc.fmsClient.On("ApprovedJob",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&proto.ApprovedJobRequest{
						Uuid:    jp.RemoteUUID.String(),
						Version: int64(spec.Version),
					},
				).Return(&proto.ApprovedJobResponse{}, nil)
				svc.orm.On("CountJobProposalsByStatus").Return(&feeds.JobProposalCounts{}, nil)
			},
			id:    spec.ID,
			force: true,
		},
		{
			name:        "already existing FMS managed job replacement success if forced",
			httpTimeout: models.MustNewDuration(1 * time.Minute),
			before: func(svc *TestService) {
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.orm.On("GetSpec", spec.ID, mock.Anything).Return(spec, nil)
				svc.orm.On("GetJobProposal", jp.ID, mock.Anything).Return(jp, nil)
				svc.jobORM.On("AssertBridgesExist", mock.IsType(pipeline.Pipeline{})).Return(nil)
				svc.orm.EXPECT().GetApprovedSpec(jp.ID, mock.Anything).Return(&feeds.JobProposalSpec{ID: 100}, nil)
				svc.orm.EXPECT().CancelSpec(int64(100), mock.Anything).Return(nil)
				svc.jobORM.On("FindJobIDByAddress", address, evmChainID, mock.Anything).Return(j.ID, nil)
				svc.spawner.On("DeleteJob", j.ID, mock.Anything).Return(nil)

				svc.spawner.
					On("CreateJob",
						mock.MatchedBy(func(j *job.Job) bool {
							return j.Name.String == "LINK / ETH | version 3 | contract 0x0000000000000000000000000000000000000000"
						}),
						mock.Anything,
					).
					Run(func(args mock.Arguments) { (args.Get(0).(*job.Job)).ID = 1 }).
					Return(nil)
				svc.orm.On("ApproveSpec",
					spec.ID,
					uuid.Must(uuid.FromString("00000000-0000-0000-0000-000000000001")),
					mock.Anything,
				).Return(nil)
				svc.fmsClient.On("ApprovedJob",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&proto.ApprovedJobRequest{
						Uuid:    jp.RemoteUUID.String(),
						Version: int64(spec.Version),
					},
				).Return(&proto.ApprovedJobResponse{}, nil)
				svc.orm.On("CountJobProposalsByStatus").Return(&feeds.JobProposalCounts{}, nil)
			},
			id:    spec.ID,
			force: true,
		},
		{
			name: "spec does not exist",
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", spec.ID, mock.Anything).Return(nil, errors.New("Not Found"))
			},
			id:      spec.ID,
			force:   false,
			wantErr: "orm: job proposal spec: Not Found",
		},
		{
			name: "cannot approve an approved spec",
			before: func(svc *TestService) {
				aspec := &feeds.JobProposalSpec{
					ID:            spec.ID,
					JobProposalID: jp.ID,
					Status:        feeds.SpecStatusApproved,
				}
				svc.orm.On("GetSpec", spec.ID, mock.Anything).Return(aspec, nil)
				svc.orm.On("GetJobProposal", jp.ID, mock.Anything).Return(jp, nil)
			},
			id:      spec.ID,
			force:   false,
			wantErr: "cannot approve an approved spec",
		},
		{
			name: "cannot approved a rejected spec",
			before: func(svc *TestService) {
				rspec := &feeds.JobProposalSpec{
					ID:            spec.ID,
					JobProposalID: jp.ID,
					Status:        feeds.SpecStatusRejected,
				}
				svc.orm.On("GetSpec", rspec.ID, mock.Anything).Return(rspec, nil)
				svc.orm.On("GetJobProposal", jp.ID, mock.Anything).Return(jp, nil)
			},
			id:      spec.ID,
			force:   false,
			wantErr: "cannot approve a rejected spec",
		},
		{
			name: "job proposal does not exist",
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", spec.ID, mock.Anything).Return(spec, nil)
				svc.orm.On("GetJobProposal", jp.ID, mock.Anything).Return(nil, errors.New("Not Found"))
			},
			id:      spec.ID,
			wantErr: "orm: job proposal: Not Found",
		},
		{
			name:        "bridges do not exist",
			httpTimeout: models.MustNewDuration(1 * time.Minute),
			before: func(svc *TestService) {
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.orm.On("GetSpec", spec.ID, mock.Anything).Return(spec, nil)
				svc.orm.On("GetJobProposal", jp.ID, mock.Anything).Return(jp, nil)
				svc.jobORM.On("AssertBridgesExist", mock.IsType(pipeline.Pipeline{})).Return(errors.New("bridges do not exist"))
			},
			id:      spec.ID,
			wantErr: "failed to approve job spec due to bridge check: bridges do not exist",
		},
		{
			name: "rpc client not connected",
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", spec.ID, mock.Anything).Return(spec, nil)
				svc.orm.On("GetJobProposal", jp.ID, mock.Anything).Return(jp, nil)
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(nil, errors.New("Not Connected"))
			},
			id:      spec.ID,
			force:   false,
			wantErr: "fms rpc client: Not Connected",
		},
		{
			name:        "create job error",
			httpTimeout: models.MustNewDuration(1 * time.Minute),
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", spec.ID, mock.Anything).Return(spec, nil)
				svc.orm.On("GetJobProposal", jp.ID, mock.Anything).Return(jp, nil)
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.jobORM.On("AssertBridgesExist", mock.IsType(pipeline.Pipeline{})).Return(nil)

				svc.jobORM.On("FindJobIDByAddress", address, evmChainID, mock.Anything).Return(int32(0), sql.ErrNoRows)

				svc.spawner.
					On("CreateJob",
						mock.MatchedBy(func(j *job.Job) bool {
							return j.Name.String == "LINK / ETH | version 3 | contract 0x0000000000000000000000000000000000000000"
						}),
						mock.Anything,
					).
					Return(errors.New("could not save"))
			},
			id:      spec.ID,
			force:   false,
			wantErr: "could not approve job proposal: could not save",
		},
		{
			name:        "approve spec orm error",
			httpTimeout: models.MustNewDuration(1 * time.Minute),
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", spec.ID, mock.Anything).Return(spec, nil)
				svc.orm.On("GetJobProposal", jp.ID, mock.Anything).Return(jp, nil)
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.jobORM.On("AssertBridgesExist", mock.IsType(pipeline.Pipeline{})).Return(nil)

				svc.jobORM.On("FindJobIDByAddress", address, evmChainID, mock.Anything).Return(int32(0), sql.ErrNoRows)

				svc.spawner.
					On("CreateJob",
						mock.MatchedBy(func(j *job.Job) bool {
							return j.Name.String == "LINK / ETH | version 3 | contract 0x0000000000000000000000000000000000000000"
						}),
						mock.Anything,
					).
					Run(func(args mock.Arguments) { (args.Get(0).(*job.Job)).ID = 1 }).
					Return(nil)
				svc.orm.On("ApproveSpec",
					spec.ID,
					uuid.Must(uuid.FromString("00000000-0000-0000-0000-000000000001")),
					mock.Anything,
				).Return(errors.New("failure"))
			},
			id:      spec.ID,
			force:   false,
			wantErr: "could not approve job proposal: failure",
		},
		{
			name:        "fms call error",
			httpTimeout: models.MustNewDuration(1 * time.Minute),
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", spec.ID, mock.Anything).Return(spec, nil)
				svc.orm.On("GetJobProposal", jp.ID, mock.Anything).Return(jp, nil)
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.jobORM.On("AssertBridgesExist", mock.IsType(pipeline.Pipeline{})).Return(nil)

				svc.jobORM.On("FindJobIDByAddress", address, evmChainID, mock.Anything).Return(int32(0), sql.ErrNoRows)

				svc.spawner.
					On("CreateJob",
						mock.MatchedBy(func(j *job.Job) bool {
							return j.Name.String == "LINK / ETH | version 3 | contract 0x0000000000000000000000000000000000000000"
						}),
						mock.Anything,
					).
					Run(func(args mock.Arguments) { (args.Get(0).(*job.Job)).ID = 1 }).
					Return(nil)
				svc.orm.On("ApproveSpec",
					spec.ID,
					uuid.Must(uuid.FromString("00000000-0000-0000-0000-000000000001")),
					mock.Anything,
				).Return(nil)
				svc.fmsClient.On("ApprovedJob",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&proto.ApprovedJobRequest{
						Uuid:    jp.RemoteUUID.String(),
						Version: int64(spec.Version),
					},
				).Return(nil, errors.New("failure"))
			},
			id:      spec.ID,
			force:   false,
			wantErr: "could not approve job proposal: failure",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			svc := setupTestServiceCfg(t, func(c *chainlink.Config, s *chainlink.Secrets) {
				if tc.httpTimeout != nil {
					c.JobPipeline.HTTPRequest.DefaultTimeout = tc.httpTimeout
				}
			})

			if tc.before != nil {
				tc.before(svc)
			}

			err := svc.ApproveSpec(ctx, tc.id, tc.force)

			if tc.wantErr != "" {
				require.Error(t, err)
				assert.EqualError(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_Service_ApproveSpec_Bootstrap(t *testing.T) {
	evmChainID := utils.NewBig(big.NewInt(0))
	address := ethkey.EIP55AddressFromAddress(common.HexToAddress("0x613a38AC1659769640aaE063C651F48E0250454C"))

	var (
		ctx  = testutils.Context(t)
		defn = `
name = 'LINK / ETH | version 3 | contract 0x0000000000000000000000000000000000000000'
type = 'bootstrap'
schemaVersion = 1
contractID = '0x613a38AC1659769640aaE063C651F48E0250454C'
externalJobID = '00000000-0000-0000-0000-000000000001'
relay = 'evm'

[relayConfig]
chainID = 0
`

		jp = &feeds.JobProposal{
			ID:             1,
			FeedsManagerID: 100,
		}
		spec = &feeds.JobProposalSpec{
			ID:            20,
			Status:        feeds.SpecStatusPending,
			JobProposalID: jp.ID,
			Version:       1,
			Definition:    defn,
		}
		rejectedSpec = &feeds.JobProposalSpec{
			ID:            20,
			Status:        feeds.SpecStatusRejected,
			JobProposalID: jp.ID,
			Version:       1,
			Definition:    defn,
		}
		cancelledSpec = &feeds.JobProposalSpec{
			ID:            20,
			Status:        feeds.SpecStatusCancelled,
			JobProposalID: jp.ID,
			Version:       1,
			Definition:    defn,
		}
		j = job.Job{
			ID:            1,
			ExternalJobID: uuid.NewV4(),
		}
	)

	testCases := []struct {
		name        string
		httpTimeout *models.Duration
		before      func(svc *TestService)
		id          int64
		force       bool
		wantErr     string
	}{
		{
			name:        "pending job success",
			httpTimeout: models.MustNewDuration(1 * time.Minute),
			before: func(svc *TestService) {
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.orm.On("GetSpec", spec.ID, mock.Anything).Return(spec, nil)
				svc.orm.On("GetJobProposal", jp.ID, mock.Anything).Return(jp, nil)
				svc.jobORM.On("AssertBridgesExist", mock.IsType(pipeline.Pipeline{})).Return(nil)
				svc.jobORM.On("FindJobIDByAddress", address, evmChainID, mock.Anything).Return(int32(0), sql.ErrNoRows)

				svc.spawner.
					On("CreateJob",
						mock.MatchedBy(func(j *job.Job) bool {
							return j.Name.String == "LINK / ETH | version 3 | contract 0x0000000000000000000000000000000000000000"
						}),
						mock.Anything,
					).
					Run(func(args mock.Arguments) { (args.Get(0).(*job.Job)).ID = 1 }).
					Return(nil)
				svc.orm.On("ApproveSpec",
					spec.ID,
					uuid.Must(uuid.FromString("00000000-0000-0000-0000-000000000001")),
					mock.Anything,
				).Return(nil)
				svc.fmsClient.On("ApprovedJob",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&proto.ApprovedJobRequest{
						Uuid:    jp.RemoteUUID.String(),
						Version: int64(spec.Version),
					},
				).Return(&proto.ApprovedJobResponse{}, nil)
				svc.orm.On("CountJobProposalsByStatus").Return(&feeds.JobProposalCounts{}, nil)
			},
			id:    spec.ID,
			force: false,
		},
		{
			name:        "cancelled spec success when it is the latest spec",
			httpTimeout: models.MustNewDuration(1 * time.Minute),
			before: func(svc *TestService) {
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.orm.On("GetSpec", cancelledSpec.ID, mock.Anything).Return(cancelledSpec, nil)
				svc.orm.On("GetJobProposal", jp.ID, mock.Anything).Return(jp, nil)
				svc.orm.On("GetLatestSpec", cancelledSpec.JobProposalID).Return(cancelledSpec, nil)
				svc.jobORM.On("AssertBridgesExist", mock.IsType(pipeline.Pipeline{})).Return(nil)

				svc.jobORM.On("FindJobIDByAddress", address, evmChainID, mock.Anything).Return(int32(0), sql.ErrNoRows)

				svc.spawner.
					On("CreateJob",
						mock.MatchedBy(func(j *job.Job) bool {
							return j.Name.String == "LINK / ETH | version 3 | contract 0x0000000000000000000000000000000000000000"
						}),
						mock.Anything,
					).
					Run(func(args mock.Arguments) { (args.Get(0).(*job.Job)).ID = 1 }).
					Return(nil)
				svc.orm.On("ApproveSpec",
					cancelledSpec.ID,
					uuid.Must(uuid.FromString("00000000-0000-0000-0000-000000000001")),
					mock.Anything,
				).Return(nil)
				svc.fmsClient.On("ApprovedJob",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&proto.ApprovedJobRequest{
						Uuid:    jp.RemoteUUID.String(),
						Version: int64(spec.Version),
					},
				).Return(&proto.ApprovedJobResponse{}, nil)
				svc.orm.On("CountJobProposalsByStatus").Return(&feeds.JobProposalCounts{}, nil)
			},
			id:    cancelledSpec.ID,
			force: false,
		},
		{
			name: "cancelled spec failed not latest spec",
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", cancelledSpec.ID, mock.Anything).Return(cancelledSpec, nil)
				svc.orm.On("GetJobProposal", jp.ID, mock.Anything).Return(jp, nil)
				svc.orm.On("GetLatestSpec", cancelledSpec.JobProposalID).Return(&feeds.JobProposalSpec{
					ID:            21,
					Status:        feeds.SpecStatusPending,
					JobProposalID: jp.ID,
					Version:       2,
					Definition:    defn,
				}, nil)
			},
			id:      cancelledSpec.ID,
			force:   false,
			wantErr: "cannot approve a cancelled spec",
		},
		{
			name: "rejected spec failed cannot be approved",
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", cancelledSpec.ID, mock.Anything).Return(rejectedSpec, nil)
				svc.orm.On("GetJobProposal", jp.ID, mock.Anything).Return(jp, nil)
			},
			id:      rejectedSpec.ID,
			force:   false,
			wantErr: "cannot approve a rejected spec",
		},
		{
			name:        "already existing job replacement error",
			httpTimeout: models.MustNewDuration(1 * time.Minute),
			before: func(svc *TestService) {
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.orm.On("GetSpec", spec.ID, mock.Anything).Return(spec, nil)
				svc.orm.On("GetJobProposal", jp.ID, mock.Anything).Return(jp, nil)
				svc.jobORM.On("AssertBridgesExist", mock.IsType(pipeline.Pipeline{})).Return(nil)
				svc.jobORM.On("FindJobIDByAddress", address, evmChainID, mock.Anything).Return(j.ID, nil)
			},
			id:      spec.ID,
			force:   false,
			wantErr: "could not approve job proposal: a job for this contract address already exists - please use the 'force' option to replace it",
		},
		{
			name:        "already existing self managed job replacement success if forced",
			httpTimeout: models.MustNewDuration(1 * time.Minute),
			before: func(svc *TestService) {
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.orm.On("GetSpec", spec.ID, mock.Anything).Return(spec, nil)
				svc.orm.On("GetJobProposal", jp.ID, mock.Anything).Return(jp, nil)
				svc.jobORM.On("AssertBridgesExist", mock.IsType(pipeline.Pipeline{})).Return(nil)
				svc.orm.EXPECT().GetApprovedSpec(jp.ID, mock.Anything).Return(nil, sql.ErrNoRows)
				svc.jobORM.On("FindJobIDByAddress", address, evmChainID, mock.Anything).Return(j.ID, nil)
				svc.spawner.On("DeleteJob", j.ID, mock.Anything).Return(nil)

				svc.spawner.
					On("CreateJob",
						mock.MatchedBy(func(j *job.Job) bool {
							return j.Name.String == "LINK / ETH | version 3 | contract 0x0000000000000000000000000000000000000000"
						}),
						mock.Anything,
					).
					Run(func(args mock.Arguments) { (args.Get(0).(*job.Job)).ID = 1 }).
					Return(nil)
				svc.orm.On("ApproveSpec",
					spec.ID,
					uuid.Must(uuid.FromString("00000000-0000-0000-0000-000000000001")),
					mock.Anything,
				).Return(nil)
				svc.fmsClient.On("ApprovedJob",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&proto.ApprovedJobRequest{
						Uuid:    jp.RemoteUUID.String(),
						Version: int64(spec.Version),
					},
				).Return(&proto.ApprovedJobResponse{}, nil)
				svc.orm.On("CountJobProposalsByStatus").Return(&feeds.JobProposalCounts{}, nil)
			},
			id:    spec.ID,
			force: true,
		},
		{
			name:        "already existing FMS managed job replacement success if forced",
			httpTimeout: models.MustNewDuration(1 * time.Minute),
			before: func(svc *TestService) {
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.orm.On("GetSpec", spec.ID, mock.Anything).Return(spec, nil)
				svc.orm.On("GetJobProposal", jp.ID, mock.Anything).Return(jp, nil)
				svc.jobORM.On("AssertBridgesExist", mock.IsType(pipeline.Pipeline{})).Return(nil)
				svc.orm.EXPECT().GetApprovedSpec(jp.ID, mock.Anything).Return(&feeds.JobProposalSpec{ID: 100}, nil)
				svc.orm.EXPECT().CancelSpec(int64(100), mock.Anything).Return(nil)
				svc.jobORM.On("FindJobIDByAddress", address, evmChainID, mock.Anything).Return(j.ID, nil)
				svc.spawner.On("DeleteJob", j.ID, mock.Anything).Return(nil)

				svc.spawner.
					On("CreateJob",
						mock.MatchedBy(func(j *job.Job) bool {
							return j.Name.String == "LINK / ETH | version 3 | contract 0x0000000000000000000000000000000000000000"
						}),
						mock.Anything,
					).
					Run(func(args mock.Arguments) { (args.Get(0).(*job.Job)).ID = 1 }).
					Return(nil)
				svc.orm.On("ApproveSpec",
					spec.ID,
					uuid.Must(uuid.FromString("00000000-0000-0000-0000-000000000001")),
					mock.Anything,
				).Return(nil)
				svc.fmsClient.On("ApprovedJob",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&proto.ApprovedJobRequest{
						Uuid:    jp.RemoteUUID.String(),
						Version: int64(spec.Version),
					},
				).Return(&proto.ApprovedJobResponse{}, nil)
				svc.orm.On("CountJobProposalsByStatus").Return(&feeds.JobProposalCounts{}, nil)
			},
			id:    spec.ID,
			force: true,
		},
		{
			name: "spec does not exist",
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", spec.ID, mock.Anything).Return(nil, errors.New("Not Found"))
			},
			id:      spec.ID,
			force:   false,
			wantErr: "orm: job proposal spec: Not Found",
		},
		{
			name: "cannot approve an approved spec",
			before: func(svc *TestService) {
				aspec := &feeds.JobProposalSpec{
					ID:            spec.ID,
					JobProposalID: jp.ID,
					Status:        feeds.SpecStatusApproved,
				}
				svc.orm.On("GetSpec", spec.ID, mock.Anything).Return(aspec, nil)
				svc.orm.On("GetJobProposal", jp.ID, mock.Anything).Return(jp, nil)
			},
			id:      spec.ID,
			force:   false,
			wantErr: "cannot approve an approved spec",
		},
		{
			name: "cannot approved a rejected spec",
			before: func(svc *TestService) {
				rspec := &feeds.JobProposalSpec{
					ID:            spec.ID,
					JobProposalID: jp.ID,
					Status:        feeds.SpecStatusRejected,
				}
				svc.orm.On("GetSpec", rspec.ID, mock.Anything).Return(rspec, nil)
				svc.orm.On("GetJobProposal", jp.ID, mock.Anything).Return(jp, nil)
			},
			id:      spec.ID,
			force:   false,
			wantErr: "cannot approve a rejected spec",
		},
		{
			name: "job proposal does not exist",
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", spec.ID, mock.Anything).Return(spec, nil)
				svc.orm.On("GetJobProposal", jp.ID, mock.Anything).Return(nil, errors.New("Not Found"))
			},
			id:      spec.ID,
			wantErr: "orm: job proposal: Not Found",
		},
		{
			name:        "bridges do not exist",
			httpTimeout: models.MustNewDuration(1 * time.Minute),
			before: func(svc *TestService) {
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.orm.On("GetSpec", spec.ID, mock.Anything).Return(spec, nil)
				svc.orm.On("GetJobProposal", jp.ID, mock.Anything).Return(jp, nil)
				svc.jobORM.On("AssertBridgesExist", mock.IsType(pipeline.Pipeline{})).Return(errors.New("bridges do not exist"))
			},
			id:      spec.ID,
			wantErr: "failed to approve job spec due to bridge check: bridges do not exist",
		},
		{
			name: "rpc client not connected",
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", spec.ID, mock.Anything).Return(spec, nil)
				svc.orm.On("GetJobProposal", jp.ID, mock.Anything).Return(jp, nil)
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(nil, errors.New("Not Connected"))
			},
			id:      spec.ID,
			force:   false,
			wantErr: "fms rpc client: Not Connected",
		},
		{
			name:        "create job error",
			httpTimeout: models.MustNewDuration(1 * time.Minute),
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", spec.ID, mock.Anything).Return(spec, nil)
				svc.orm.On("GetJobProposal", jp.ID, mock.Anything).Return(jp, nil)
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.jobORM.On("AssertBridgesExist", mock.IsType(pipeline.Pipeline{})).Return(nil)

				svc.jobORM.On("FindJobIDByAddress", address, evmChainID, mock.Anything).Return(int32(0), sql.ErrNoRows)

				svc.spawner.
					On("CreateJob",
						mock.MatchedBy(func(j *job.Job) bool {
							return j.Name.String == "LINK / ETH | version 3 | contract 0x0000000000000000000000000000000000000000"
						}),
						mock.Anything,
					).
					Return(errors.New("could not save"))
			},
			id:      spec.ID,
			force:   false,
			wantErr: "could not approve job proposal: could not save",
		},
		{
			name:        "approve spec orm error",
			httpTimeout: models.MustNewDuration(1 * time.Minute),
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", spec.ID, mock.Anything).Return(spec, nil)
				svc.orm.On("GetJobProposal", jp.ID, mock.Anything).Return(jp, nil)
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.jobORM.On("AssertBridgesExist", mock.IsType(pipeline.Pipeline{})).Return(nil)

				svc.jobORM.On("FindJobIDByAddress", address, evmChainID, mock.Anything).Return(int32(0), sql.ErrNoRows)

				svc.spawner.
					On("CreateJob",
						mock.MatchedBy(func(j *job.Job) bool {
							return j.Name.String == "LINK / ETH | version 3 | contract 0x0000000000000000000000000000000000000000"
						}),
						mock.Anything,
					).
					Run(func(args mock.Arguments) { (args.Get(0).(*job.Job)).ID = 1 }).
					Return(nil)
				svc.orm.On("ApproveSpec",
					spec.ID,
					uuid.Must(uuid.FromString("00000000-0000-0000-0000-000000000001")),
					mock.Anything,
				).Return(errors.New("failure"))
			},
			id:      spec.ID,
			force:   false,
			wantErr: "could not approve job proposal: failure",
		},
		{
			name:        "fms call error",
			httpTimeout: models.MustNewDuration(1 * time.Minute),
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", spec.ID, mock.Anything).Return(spec, nil)
				svc.orm.On("GetJobProposal", jp.ID, mock.Anything).Return(jp, nil)
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.jobORM.On("AssertBridgesExist", mock.IsType(pipeline.Pipeline{})).Return(nil)

				svc.jobORM.On("FindJobIDByAddress", address, evmChainID, mock.Anything).Return(int32(0), sql.ErrNoRows)

				svc.spawner.
					On("CreateJob",
						mock.MatchedBy(func(j *job.Job) bool {
							return j.Name.String == "LINK / ETH | version 3 | contract 0x0000000000000000000000000000000000000000"
						}),
						mock.Anything,
					).
					Run(func(args mock.Arguments) { (args.Get(0).(*job.Job)).ID = 1 }).
					Return(nil)
				svc.orm.On("ApproveSpec",
					spec.ID,
					uuid.Must(uuid.FromString("00000000-0000-0000-0000-000000000001")),
					mock.Anything,
				).Return(nil)
				svc.fmsClient.On("ApprovedJob",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&proto.ApprovedJobRequest{
						Uuid:    jp.RemoteUUID.String(),
						Version: int64(spec.Version),
					},
				).Return(nil, errors.New("failure"))
			},
			id:      spec.ID,
			force:   false,
			wantErr: "could not approve job proposal: failure",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			svc := setupTestServiceCfg(t, func(c *chainlink.Config, s *chainlink.Secrets) {
				if tc.httpTimeout != nil {
					c.JobPipeline.HTTPRequest.DefaultTimeout = tc.httpTimeout
				}
			})

			if tc.before != nil {
				tc.before(svc)
			}

			err := svc.ApproveSpec(ctx, tc.id, tc.force)

			if tc.wantErr != "" {
				require.Error(t, err)
				assert.EqualError(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_Service_RejectSpec(t *testing.T) {
	var (
		ctx = testutils.Context(t)
		jp  = &feeds.JobProposal{
			ID:             1,
			FeedsManagerID: 100,
		}
		spec = &feeds.JobProposalSpec{
			ID:            20,
			Status:        feeds.SpecStatusPending,
			JobProposalID: jp.ID,
			Version:       1,
		}
	)

	testCases := []struct {
		name    string
		before  func(svc *TestService)
		wantErr string
	}{
		{
			name: "Success",
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", spec.ID, mock.Anything).Return(spec, nil)
				svc.orm.On("GetJobProposal", jp.ID, mock.Anything).Return(jp, nil)
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.orm.On("RejectSpec",
					spec.ID,
					mock.Anything,
				).Return(nil)
				svc.fmsClient.On("RejectedJob",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&proto.RejectedJobRequest{
						Uuid:    jp.RemoteUUID.String(),
						Version: int64(spec.Version),
					},
				).Return(&proto.RejectedJobResponse{}, nil)
				svc.orm.On("CountJobProposalsByStatus").Return(&feeds.JobProposalCounts{}, nil)
			},
		},
		{
			name: "Fails to get spec",
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", spec.ID, mock.Anything).Return(nil, errors.New("failure"))
			},
			wantErr: "failure",
		},
		{
			name: "Cannot be a rejected proposal",
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", spec.ID, mock.Anything).Return(&feeds.JobProposalSpec{
					Status: feeds.SpecStatusRejected,
				}, nil)
			},
			wantErr: "must be a pending job proposal spec",
		},
		{
			name: "Fails to get proposal",
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", spec.ID, mock.Anything).Return(spec, nil)
				svc.orm.On("GetJobProposal", jp.ID, mock.Anything).Return(nil, errors.New("failure"))
			},
			wantErr: "failure",
		},
		{
			name: "FMS not connected",
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", spec.ID, mock.Anything).Return(spec, nil)
				svc.orm.On("GetJobProposal", jp.ID, mock.Anything).Return(jp, nil)
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(nil, errors.New("disconnected"))
			},
			wantErr: "disconnected",
		},
		{
			name: "Fails to update spec",
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", spec.ID, mock.Anything).Return(spec, nil)
				svc.orm.On("GetJobProposal", jp.ID, mock.Anything).Return(jp, nil)
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.orm.On("RejectSpec", mock.Anything, mock.Anything).Return(errors.New("failure"))
			},
			wantErr: "failure",
		},
		{
			name: "Fails to update spec",
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", spec.ID, mock.Anything).Return(spec, nil)
				svc.orm.On("GetJobProposal", jp.ID, mock.Anything).Return(jp, nil)
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.orm.On("RejectSpec", mock.Anything, mock.Anything).Return(nil)
				svc.fmsClient.
					On("RejectedJob",
						mock.MatchedBy(func(ctx context.Context) bool { return true }),
						&proto.RejectedJobRequest{
							Uuid:    jp.RemoteUUID.String(),
							Version: int64(spec.Version),
						}).
					Return(nil, errors.New("rpc failure"))
			},
			wantErr: "rpc failure",
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

			err := svc.RejectSpec(ctx, spec.ID)

			if tc.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_Service_UpdateSpecDefinition(t *testing.T) {
	var (
		ctx         = testutils.Context(t)
		specID      = int64(1)
		updatedSpec = "updated spec"
		spec        = &feeds.JobProposalSpec{
			ID:         specID,
			Status:     feeds.SpecStatusPending,
			Definition: "spec",
		}
	)

	testCases := []struct {
		name    string
		before  func(svc *TestService)
		specID  int64
		wantErr string
	}{
		{
			name: "success",
			before: func(svc *TestService) {
				svc.orm.
					On("GetSpec", specID, mock.Anything).
					Return(spec, nil)
				svc.orm.On("UpdateSpecDefinition",
					specID,
					updatedSpec,
					mock.Anything,
				).Return(nil)
			},
			specID: specID,
		},
		{
			name: "does not exist",
			before: func(svc *TestService) {
				svc.orm.
					On("GetSpec", specID, mock.Anything).
					Return(nil, sql.ErrNoRows)
			},
			specID:  specID,
			wantErr: "job proposal spec does not exist: sql: no rows in result set",
		},
		{
			name: "other get errors",
			before: func(svc *TestService) {
				svc.orm.
					On("GetSpec", specID, mock.Anything).
					Return(nil, errors.New("other db error"))
			},
			specID:  specID,
			wantErr: "database error: other db error",
		},
		{
			name: "cannot edit",
			before: func(svc *TestService) {
				spec := &feeds.JobProposalSpec{
					ID:     1,
					Status: feeds.SpecStatusApproved,
				}

				svc.orm.
					On("GetSpec", specID, mock.Anything).
					Return(spec, nil)
			},
			specID:  specID,
			wantErr: "must be a pending or cancelled spec",
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

			err := svc.UpdateSpecDefinition(ctx, tc.specID, updatedSpec)
			if tc.wantErr != "" {
				assert.EqualError(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
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

	tests := []struct {
		name       string
		beforeFunc func(svc *TestService)
	}{
		{
			name: "success with a feeds manager connection",
			beforeFunc: func(svc *TestService) {
				svc.csaKeystore.On("GetAll").Return([]csakey.KeyV2{key}, nil)
				svc.orm.On("ListManagers").Return([]feeds.FeedsManager{mgr}, nil)
				svc.connMgr.On("IsConnected", mgr.ID).Return(false)
				svc.connMgr.On("Connect", mock.IsType(feeds.ConnectOpts{}))
				svc.connMgr.On("Close")
				svc.orm.On("CountJobProposalsByStatus").Return(&feeds.JobProposalCounts{}, nil)
			},
		},
		{
			name: "success with no registered managers",
			beforeFunc: func(svc *TestService) {
				svc.csaKeystore.On("GetAll").Return([]csakey.KeyV2{key}, nil)
				svc.orm.On("ListManagers").Return([]feeds.FeedsManager{}, nil)
				svc.connMgr.On("Close")
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := setupTestService(t)

			if tt.beforeFunc != nil {
				tt.beforeFunc(svc)
			}

			require.NoError(t, svc.Start(testutils.Context(t)))

			svc.Close()
		})
	}
}
