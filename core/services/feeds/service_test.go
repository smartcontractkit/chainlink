package feeds_test

import (
	"context"
	"database/sql"
	"encoding/hex"
	"fmt"
	"maps"
	"math/big"
	"slices"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
	"gopkg.in/guregu/null.v4"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	proto "github.com/smartcontractkit/chainlink-protos/orchestrator/feedsmanager"

	"github.com/smartcontractkit/chainlink-evm/pkg/chains/legacyevm"
	"github.com/smartcontractkit/chainlink-evm/pkg/heads"
	"github.com/smartcontractkit/chainlink-evm/pkg/types"
	evmbig "github.com/smartcontractkit/chainlink-evm/pkg/utils/big"

	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/feeds"
	"github.com/smartcontractkit/chainlink/v2/core/services/feeds/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	jobmocks "github.com/smartcontractkit/chainlink/v2/core/services/job/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/csakey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocrkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/workflowkey"
	ksmocks "github.com/smartcontractkit/chainlink/v2/core/services/keystore/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/services/versioning"
	"github.com/smartcontractkit/chainlink/v2/core/testdata/testspecs"
	"github.com/smartcontractkit/chainlink/v2/core/utils/crypto"
)

const FluxMonitorTestSpecTemplate = `
type              = "fluxmonitor"
schemaVersion     = 1
name              = "%s"
contractAddress   = "0x3cCad4715152693fE3BC4460591e3D3Fbd071b42"
externalJobID     = "%s"
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

const OCR1TestSpecTemplate = `
type               = "offchainreporting"
schemaVersion      = 1
name              = "%s"
externalJobID       = "%s"
evmChainID 		   = 0
contractAddress    = "0x613a38AC1659769640aaE063C651F48E0250454C"
p2pv2Bootstrappers = ["12D3KooWHfYFQ8hGttAYbMCevQVESEQhzJAqFZokMVtom8bNxwGq@127.0.0.1:5001"]
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

const OCR2TestSpecTemplate = `
type               = "offchainreporting2"
pluginType         = "median"
schemaVersion      = 1
name              = "%s"
relay              = "evm"
contractID         = "0x613a38AC1659769640aaE063C651F48E0250454C"
externalJobID      = "%s"
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
gasPriceSubunitsSource = """
ds1          [type=bridge name=voter_turnout];
ds1_parse    [type=jsonparse path="one,two"];
ds1_multiply [type=multiply times=1.23];
ds1 -> ds1_parse -> ds1_multiply -> answer1;
answer1      [type=median index=0];
"""
[pluginConfig.juelsPerFeeCoinCache]
updateInterval = "1m"
`

const StreamTestSpecTemplate = `
name = '%s'
type = 'stream'
schemaVersion = 1
externalJobID = '%s'
streamID = %d
observationSource = """
ds1_payload [type=bridge name=\"bridge-ncfx\" timeout=\"50s\" requestData=\"{\\\"data\\\":{\\\"endpoint\\\":\\\"cryptolwba\\\",\\\"from\\\":\\\"SEI\\\",\\\"to\\\":\\\"USD\\\"}}\"];
ds1_benchmark [type=jsonparse path=\"data,mid\"];
ds1_bid [type=jsonparse path=\"data,bid\"];
ds1_ask [type=jsonparse path=\"data,ask\"];
ds2_payload [type=bridge name=\"bridge-tiingo\" timeout=\"50s\" requestData=\"{\\\"data\\\":{\\\"endpoint\\\":\\\"cryptolwba\\\",\\\"from\\\":\\\"SEI\\\",\\\"to\\\":\\\"USD\\\"}}\"];
ds2_benchmark [type=jsonparse path=\"data,mid\"];
ds2_bid [type=jsonparse path=\"data,bid\"];
ds2_ask [type=jsonparse path=\"data,ask\"];
ds3_payload [type=bridge name=\"bridge-gsr\" timeout=\"50s\" requestData=\"{\\\"data\\\":{\\\"endpoint\\\":\\\"cryptolwba\\\",\\\"from\\\":\\\"SEI\\\",\\\"to\\\":\\\"USD\\\"}}\"];
ds3_benchmark [type=jsonparse path=\"data,mid\"];
ds3_bid [type=jsonparse path=\"data,bid\"];
ds3_ask [type=jsonparse path=\"data,ask\"];
ds1_payload -> ds1_benchmark -> benchmark_price;
ds2_payload -> ds2_benchmark -> benchmark_price;
ds3_payload -> ds3_benchmark -> benchmark_price;
benchmark_price [type=median allowedFaults=2 index=0];
ds1_payload -> ds1_bid -> bid_price;
ds2_payload -> ds2_bid -> bid_price;
ds3_payload -> ds3_bid -> bid_price;
bid_price [type=median allowedFaults=2 index=1];
ds1_payload -> ds1_ask -> ask_price;
ds2_payload -> ds2_ask -> ask_price;
ds3_payload -> ds3_ask -> ask_price;
ask_price [type=median allowedFaults=2 index=2];
"""
`

const BootstrapTestSpecTemplate = `
type				= "bootstrap"
schemaVersion		= 1
name              = "%s"
contractID			= "0x613a38AC1659769640aaE063C651F48E0250454C"
relay				= "evm"
[relayConfig]
chainID 			= 1337
`

type TestService struct {
	feeds.Service
	orm              *mocks.ORM
	jobORM           *jobmocks.ORM
	connMgr          *mocks.ConnectionsManager
	spawner          *jobmocks.Spawner
	fmsClient        *mocks.FeedsManagerClient
	csaKeystore      *ksmocks.CSA
	p2pKeystore      *ksmocks.P2P
	ocr1Keystore     *ksmocks.OCR
	ocr2Keystore     *ksmocks.OCR2
	workflowKeystore *ksmocks.Workflow
	legacyChains     legacyevm.LegacyChainContainer
	logs             *observer.ObservedLogs
}

func setupTestService(t *testing.T, opts ...feeds.ServiceOption) *TestService {
	t.Helper()

	return setupTestServiceCfg(t, nil, opts...)
}

func setupTestServiceCfg(
	t *testing.T, overrideCfg func(c *chainlink.Config, s *chainlink.Secrets), opts ...feeds.ServiceOption,
) *TestService {
	t.Helper()

	var (
		orm              = mocks.NewORM(t)
		jobORM           = jobmocks.NewORM(t)
		connMgr          = mocks.NewConnectionsManager(t)
		spawner          = jobmocks.NewSpawner(t)
		fmsClient        = mocks.NewFeedsManagerClient(t)
		csaKeystore      = ksmocks.NewCSA(t)
		p2pKeystore      = ksmocks.NewP2P(t)
		ocr1Keystore     = ksmocks.NewOCR(t)
		ocr2Keystore     = ksmocks.NewOCR2(t)
		workflowKeystore = ksmocks.NewWorkflow(t)
	)

	lggr, observedLogs := logger.TestLoggerObserved(t, zap.DebugLevel)

	db := pgtest.NewSqlxDB(t)
	gcfg := configtest.NewGeneralConfig(t, overrideCfg)
	keyStore := new(ksmocks.Master)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	legacyChains := evmtest.NewLegacyChains(t, evmtest.TestChainOpts{
		ChainConfigs:   gcfg.EVMConfigs(),
		DatabaseConfig: gcfg.Database(),
		FeatureConfig:  gcfg.Feature(),
		ListenerConfig: gcfg.Database().Listener(),
		KeyStore:       ethKeyStore,
		DB:             db,
		HeadTracker:    heads.NullTracker,
	})
	keyStore.On("Eth").Return(ethKeyStore)
	keyStore.On("CSA").Return(csaKeystore)
	keyStore.On("P2P").Return(p2pKeystore)
	keyStore.On("OCR").Return(ocr1Keystore)
	keyStore.On("OCR2").Return(ocr2Keystore)
	keyStore.On("Workflow").Return(workflowKeystore)
	svc := feeds.NewService(orm, jobORM, db, spawner, keyStore, gcfg, gcfg.Feature(), gcfg.Insecure(),
		gcfg.JobPipeline(), gcfg.OCR(), gcfg.OCR2(), legacyChains, lggr, "1.0.0", nil, opts...)
	svc.SetConnectionsManager(connMgr)

	return &TestService{
		Service:          svc,
		orm:              orm,
		jobORM:           jobORM,
		connMgr:          connMgr,
		spawner:          spawner,
		fmsClient:        fmsClient,
		csaKeystore:      csaKeystore,
		p2pKeystore:      p2pKeystore,
		ocr1Keystore:     ocr1Keystore,
		ocr2Keystore:     ocr2Keystore,
		workflowKeystore: workflowKeystore,
		legacyChains:     legacyChains,
		logs:             observedLogs,
	}
}

func Test_Service_RegisterManager(t *testing.T) {
	t.Parallel()

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

	svc.orm.On("CountManagers", mock.Anything).Return(int64(0), nil)
	svc.orm.On("CreateManager", mock.Anything, &mgr, mock.Anything).
		Return(id, nil)
	svc.orm.On("CreateBatchChainConfig", mock.Anything, params.ChainConfigs, mock.Anything).
		Return([]int64{}, nil)
	// ListManagers runs in a goroutine so it might be called.
	svc.orm.On("ListManagers", testutils.Context(t)).Return([]feeds.FeedsManager{mgr}, nil).Maybe()
	transactCall := svc.orm.On("Transact", mock.Anything, mock.Anything)
	transactCall.Run(func(args mock.Arguments) {
		fn := args[1].(func(orm feeds.ORM) error)
		transactCall.ReturnArguments = mock.Arguments{fn(svc.orm)}
	})
	svc.connMgr.On("Connect", mock.IsType(feeds.ConnectOpts{}))

	actual, err := svc.RegisterManager(testutils.Context(t), params)
	require.NoError(t, err)

	assert.Equal(t, actual, id)
}

func Test_Service_RegisterManager_MultiFeedsManager(t *testing.T) {
	t.Parallel()

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

	svc := setupTestServiceCfg(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		var multiFeedsManagers = true
		c.Feature.MultiFeedsManagers = &multiFeedsManagers
	})
	ctx := testutils.Context(t)

	svc.orm.On("ManagerExists", ctx, params.PublicKey).Return(false, nil)
	svc.orm.On("CreateManager", mock.Anything, &mgr, mock.Anything).
		Return(id, nil)
	svc.orm.On("CreateBatchChainConfig", mock.Anything, params.ChainConfigs, mock.Anything).
		Return([]int64{}, nil)
	// ListManagers runs in a goroutine so it might be called.
	svc.orm.On("ListManagers", ctx).Return([]feeds.FeedsManager{mgr}, nil).Maybe()
	transactCall := svc.orm.On("Transact", mock.Anything, mock.Anything)
	transactCall.Run(func(args mock.Arguments) {
		fn := args[1].(func(orm feeds.ORM) error)
		transactCall.ReturnArguments = mock.Arguments{fn(svc.orm)}
	})
	svc.connMgr.On("Connect", mock.IsType(feeds.ConnectOpts{}))

	actual, err := svc.RegisterManager(ctx, params)
	require.NoError(t, err)

	assert.Equal(t, actual, id)
}

func Test_Service_RegisterManager_InvalidCreateManager(t *testing.T) {
	t.Parallel()

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

	svc.orm.On("CountManagers", mock.Anything).Return(int64(0), nil)
	svc.orm.On("CreateManager", mock.Anything, &mgr, mock.Anything).
		Return(id, errors.New("orm error"))
	// ListManagers runs in a goroutine so it might be called.
	svc.orm.On("ListManagers", testutils.Context(t)).Return([]feeds.FeedsManager{mgr}, nil).Maybe()

	transactCall := svc.orm.On("Transact", mock.Anything, mock.Anything)
	transactCall.Run(func(args mock.Arguments) {
		fn := args[1].(func(orm feeds.ORM) error)
		transactCall.ReturnArguments = mock.Arguments{fn(svc.orm)}
	})
	_, err = svc.RegisterManager(testutils.Context(t), params)
	require.Error(t, err)
	assert.Equal(t, "orm error", err.Error())
}

func Test_Service_RegisterManager_DuplicateFeedsManager(t *testing.T) {
	t.Parallel()

	var pubKeyHex = "0f17c3bf72de8beef6e2d17a14c0a972f5d7e0e66e70722373f12b88382d40f9"
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

	svc := setupTestServiceCfg(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		var multiFeedsManagers = true
		c.Feature.MultiFeedsManagers = &multiFeedsManagers
	})
	ctx := testutils.Context(t)

	svc.orm.On("ManagerExists", ctx, params.PublicKey).Return(true, nil)
	// ListManagers runs in a goroutine so it might be called.
	svc.orm.On("ListManagers", ctx).Return([]feeds.FeedsManager{mgr}, nil).Maybe()

	_, err = svc.RegisterManager(ctx, params)
	require.Error(t, err)

	assert.Equal(t, "manager was previously registered using the same public key", err.Error())
}

func Test_Service_ListManagers(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	var (
		mgr  = feeds.FeedsManager{}
		mgrs = []feeds.FeedsManager{mgr}
	)
	svc := setupTestService(t)

	svc.orm.On("ListManagers", mock.Anything).Return(mgrs, nil)
	svc.connMgr.On("IsConnected", mgr.ID).Return(false)

	actual, err := svc.ListManagers(ctx)
	require.NoError(t, err)

	assert.Equal(t, mgrs, actual)
}

func Test_Service_GetManager(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	var (
		id  = int64(1)
		mgr = feeds.FeedsManager{ID: id}
	)
	svc := setupTestService(t)

	svc.orm.On("GetManager", mock.Anything, id).
		Return(&mgr, nil)
	svc.connMgr.On("IsConnected", mgr.ID).Return(false)

	actual, err := svc.GetManager(ctx, id)
	require.NoError(t, err)

	assert.Equal(t, actual, &mgr)
}

func Test_Service_UpdateFeedsManager(t *testing.T) {

	var (
		mgr = feeds.FeedsManager{ID: 1}
	)

	svc := setupTestService(t)

	svc.orm.On("UpdateManager", mock.Anything, mgr, mock.Anything).Return(nil)
	svc.connMgr.On("Disconnect", mgr.ID).Return(nil)
	svc.connMgr.On("Connect", mock.IsType(feeds.ConnectOpts{})).Return(nil)

	err := svc.UpdateManager(testutils.Context(t), mgr)
	require.NoError(t, err)
}

func Test_Service_EnableFeedsManager(t *testing.T) {
	mgr := feeds.FeedsManager{ID: 1}

	svc := setupTestService(t)

	svc.orm.On("EnableManager", mock.Anything, mgr.ID).Return(&mgr, nil)
	svc.connMgr.On("IsConnected", mgr.ID).Return(false)
	svc.connMgr.On("Disconnect", mgr.ID).Return(nil)
	svc.connMgr.On("Connect", mock.IsType(feeds.ConnectOpts{})).Return(nil)

	actual, err := svc.EnableManager(testutils.Context(t), 1)
	require.NoError(t, err)
	require.NotNil(t, actual)
}

func Test_Service_DisableFeedsManager(t *testing.T) {
	mgr := feeds.FeedsManager{ID: 1}

	svc := setupTestService(t)

	svc.orm.On("DisableManager", mock.Anything, mgr.ID).Return(&mgr, nil)
	svc.connMgr.On("IsConnected", mgr.ID).Return(false)
	svc.connMgr.On("Disconnect", mgr.ID).Return(nil)

	actual, err := svc.DisableManager(testutils.Context(t), 1)
	require.NoError(t, err)
	require.NotNil(t, actual)
}

func Test_Service_ListManagersByIDs(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	var (
		mgr  = feeds.FeedsManager{}
		mgrs = []feeds.FeedsManager{mgr}
	)
	svc := setupTestService(t)

	svc.orm.On("ListManagersByIDs", mock.Anything, []int64{mgr.ID}).
		Return(mgrs, nil)
	svc.connMgr.On("IsConnected", mgr.ID).Return(false)

	actual, err := svc.ListManagersByIDs(ctx, []int64{mgr.ID})
	require.NoError(t, err)

	assert.Equal(t, mgrs, actual)
}

func Test_Service_CreateChainConfig(t *testing.T) {
	tests := []struct {
		name              string
		chainType         feeds.ChainType
		expectedID        int64
		expectedChainType proto.ChainType
	}{
		{
			name:              "EVM Chain Type",
			chainType:         feeds.ChainTypeEVM,
			expectedID:        int64(1),
			expectedChainType: proto.ChainType_CHAIN_TYPE_EVM,
		},
		{
			name:              "Solana Chain Type",
			chainType:         feeds.ChainTypeSolana,
			expectedID:        int64(1),
			expectedChainType: proto.ChainType_CHAIN_TYPE_SOLANA,
		},
		{
			name:              "Starknet Chain Type",
			chainType:         feeds.ChainTypeStarknet,
			expectedID:        int64(1),
			expectedChainType: proto.ChainType_CHAIN_TYPE_STARKNET,
		},
		{
			name:              "Aptos Chain Type",
			chainType:         feeds.ChainTypeAptos,
			expectedID:        int64(1),
			expectedChainType: proto.ChainType_CHAIN_TYPE_APTOS,
		},
		{
			name:              "Tron Chain Type",
			chainType:         feeds.ChainTypeTron,
			expectedID:        int64(1),
			expectedChainType: proto.ChainType_CHAIN_TYPE_TRON,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				mgr         = feeds.FeedsManager{ID: 1}
				nodeVersion = &versioning.NodeVersion{
					Version: "1.0.0",
				}
				cfg = feeds.ChainConfig{
					FeedsManagerID:          mgr.ID,
					ChainID:                 "42",
					ChainType:               tt.chainType,
					AccountAddress:          "0x0000000000000000000000000000000000000000",
					AccountAddressPublicKey: null.StringFrom("0x0000000000000000000000000000000000000002"),
					AdminAddress:            "0x0000000000000000000000000000000000000001",
					FluxMonitorConfig: feeds.FluxMonitorConfig{
						Enabled: true,
					},
					OCR1Config: feeds.OCR1Config{
						Enabled: false,
					},
					OCR2Config: feeds.OCR2ConfigModel{
						Enabled: false,
					},
				}

				svc = setupTestService(t)
			)

			p2pKey, err := p2pkey.NewV2()
			require.NoError(t, err)
			svc.p2pKeystore.On("GetAll").Return([]p2pkey.KeyV2{p2pKey}, nil)

			workflowKey, err := workflowkey.New()
			require.NoError(t, err)
			svc.workflowKeystore.On("EnsureKey", mock.Anything).Return(nil)
			svc.workflowKeystore.On("GetAll").Return([]workflowkey.Key{workflowKey}, nil)

			svc.orm.On("CreateChainConfig", mock.Anything, cfg).Return(int64(1), nil)
			svc.orm.On("GetManager", mock.Anything, mgr.ID).Return(&mgr, nil)
			svc.connMgr.On("GetClient", mgr.ID).Return(svc.fmsClient, nil)
			svc.orm.On("ListChainConfigsByManagerIDs", mock.Anything, []int64{mgr.ID}).Return([]feeds.ChainConfig{cfg}, nil)
			wkID := workflowKey.ID()
			svc.fmsClient.On("UpdateNode", mock.Anything, &proto.UpdateNodeRequest{
				Version: nodeVersion.Version,
				ChainConfigs: []*proto.ChainConfig{
					{
						Chain: &proto.Chain{
							Id:   cfg.ChainID,
							Type: tt.expectedChainType,
						},
						AccountAddress:          cfg.AccountAddress,
						AccountAddressPublicKey: &cfg.AccountAddressPublicKey.String,
						AdminAddress:            cfg.AdminAddress,
						FluxMonitorConfig:       &proto.FluxMonitorConfig{Enabled: true},
						Ocr1Config:              &proto.OCR1Config{Enabled: false},
						Ocr2Config:              &proto.OCR2Config{Enabled: false},
					},
				},
				WorkflowKey:   &wkID,
				P2PKeyBundles: []*proto.P2PKeyBundle{{PeerId: p2pKey.PeerID().String(), PublicKey: p2pKey.PublicKeyHex()}},
			}).Return(&proto.UpdateNodeResponse{}, nil)

			actual, err := svc.CreateChainConfig(testutils.Context(t), cfg)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedID, actual)
			waitSyncNodeInfoCall(t, svc.logs)
		})
	}
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
			OCR2Config:        feeds.OCR2ConfigModel{Enabled: false},
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

	workflowKey, err := workflowkey.New()
	require.NoError(t, err)
	svc.workflowKeystore.On("EnsureKey", mock.Anything).Return(nil)
	svc.workflowKeystore.On("GetAll").Return([]workflowkey.Key{workflowKey}, nil)
	svc.p2pKeystore.On("GetAll").Return([]p2pkey.KeyV2{}, nil)

	svc.orm.On("GetChainConfig", mock.Anything, cfg.ID).Return(&cfg, nil)
	svc.orm.On("DeleteChainConfig", mock.Anything, cfg.ID).Return(cfg.ID, nil)
	svc.orm.On("GetManager", mock.Anything, mgr.ID).Return(&mgr, nil)
	svc.connMgr.On("GetClient", mgr.ID).Return(svc.fmsClient, nil)
	svc.orm.On("ListChainConfigsByManagerIDs", mock.Anything, []int64{mgr.ID}).Return([]feeds.ChainConfig{}, nil)
	wkID := workflowKey.ID()
	svc.fmsClient.On("UpdateNode", mock.Anything, &proto.UpdateNodeRequest{
		Version:       nodeVersion.Version,
		ChainConfigs:  []*proto.ChainConfig{},
		WorkflowKey:   &wkID,
		P2PKeyBundles: []*proto.P2PKeyBundle{},
	}).Return(&proto.UpdateNodeResponse{}, nil)

	actual, err := svc.DeleteChainConfig(testutils.Context(t), cfg.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(1), actual)
	waitSyncNodeInfoCall(t, svc.logs)
}

func Test_Service_ListChainConfigsByManagerIDs(t *testing.T) {
	ctx := testutils.Context(t)
	var (
		mgr = feeds.FeedsManager{ID: 1}
		cfg = feeds.ChainConfig{
			ID:             1,
			FeedsManagerID: mgr.ID,
		}
		ids = []int64{cfg.ID}

		svc = setupTestService(t)
	)

	svc.orm.On("ListChainConfigsByManagerIDs", mock.Anything, ids).Return([]feeds.ChainConfig{cfg}, nil)

	actual, err := svc.ListChainConfigsByManagerIDs(ctx, ids)
	require.NoError(t, err)
	assert.Equal(t, []feeds.ChainConfig{cfg}, actual)
}

func Test_Service_UpdateChainConfig(t *testing.T) {
	tests := []struct {
		name              string
		chainType         feeds.ChainType
		expectedChainType proto.ChainType
	}{
		{
			name:              "EVM Chain Type",
			chainType:         feeds.ChainTypeEVM,
			expectedChainType: proto.ChainType_CHAIN_TYPE_EVM,
		},
		{
			name:              "Solana Chain Type",
			chainType:         feeds.ChainTypeSolana,
			expectedChainType: proto.ChainType_CHAIN_TYPE_SOLANA,
		},
		{
			name:              "Starknet Chain Type",
			chainType:         feeds.ChainTypeStarknet,
			expectedChainType: proto.ChainType_CHAIN_TYPE_STARKNET,
		},
		{
			name:              "Aptos Chain Type",
			chainType:         feeds.ChainTypeAptos,
			expectedChainType: proto.ChainType_CHAIN_TYPE_APTOS,
		},
		{
			name:              "Tron Chain Type",
			chainType:         feeds.ChainTypeTron,
			expectedChainType: proto.ChainType_CHAIN_TYPE_TRON,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				mgr         = feeds.FeedsManager{ID: 1}
				nodeVersion = &versioning.NodeVersion{
					Version: "1.0.0",
				}
				cfg = feeds.ChainConfig{
					FeedsManagerID:          mgr.ID,
					ChainID:                 "42",
					ChainType:               tt.chainType,
					AccountAddress:          "0x0000000000000000000000000000000000000000",
					AccountAddressPublicKey: null.StringFrom("0x0000000000000000000000000000000000000002"),
					AdminAddress:            "0x0000000000000000000000000000000000000001",
					FluxMonitorConfig:       feeds.FluxMonitorConfig{Enabled: false},
					OCR1Config:              feeds.OCR1Config{Enabled: false},
					OCR2Config:              feeds.OCR2ConfigModel{Enabled: false},
				}

				svc = setupTestService(t)
			)

			workflowKey, err := workflowkey.New()
			require.NoError(t, err)
			svc.workflowKeystore.On("EnsureKey", mock.Anything).Return(nil)
			svc.workflowKeystore.On("GetAll").Return([]workflowkey.Key{workflowKey}, nil)
			svc.p2pKeystore.On("GetAll").Return([]p2pkey.KeyV2{}, nil)

			svc.orm.On("UpdateChainConfig", mock.Anything, cfg).Return(int64(1), nil)
			svc.orm.On("GetChainConfig", mock.Anything, cfg.ID).Return(&cfg, nil)
			svc.connMgr.On("GetClient", mgr.ID).Return(svc.fmsClient, nil)
			svc.orm.On("ListChainConfigsByManagerIDs", mock.Anything, []int64{mgr.ID}).Return([]feeds.ChainConfig{cfg}, nil)
			wkID := workflowKey.ID()
			svc.fmsClient.On("UpdateNode", mock.Anything, &proto.UpdateNodeRequest{
				Version: nodeVersion.Version,
				ChainConfigs: []*proto.ChainConfig{
					{
						Chain: &proto.Chain{
							Id:   cfg.ChainID,
							Type: tt.expectedChainType,
						},
						AccountAddress:          cfg.AccountAddress,
						AdminAddress:            cfg.AdminAddress,
						AccountAddressPublicKey: &cfg.AccountAddressPublicKey.String,
						FluxMonitorConfig:       &proto.FluxMonitorConfig{Enabled: false},
						Ocr1Config:              &proto.OCR1Config{Enabled: false},
						Ocr2Config:              &proto.OCR2Config{Enabled: false},
					},
				},
				WorkflowKey:   &wkID,
				P2PKeyBundles: []*proto.P2PKeyBundle{},
			}).Return(&proto.UpdateNodeResponse{}, nil)

			actual, err := svc.UpdateChainConfig(testutils.Context(t), cfg)
			require.NoError(t, err)
			assert.Equal(t, int64(1), actual)
			waitSyncNodeInfoCall(t, svc.logs)
		})
	}
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
			OCR2Config:        feeds.OCR2ConfigModel{Enabled: false},
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
		remoteUUIDFluxMonitor = uuid.New()
		nameAndExternalJobID  = uuid.New()
		spec                  = fmt.Sprintf(FluxMonitorTestSpecTemplate, nameAndExternalJobID, nameAndExternalJobID)
		argsFluxMonitor       = &feeds.ProposeJobArgs{
			FeedsManagerID: 1,
			RemoteUUID:     remoteUUIDFluxMonitor,
			Spec:           spec,
			Version:        1,
		}
		jpFluxMonitor = feeds.JobProposal{
			FeedsManagerID: 1,
			Name:           null.StringFrom(nameAndExternalJobID.String()),
			RemoteUUID:     remoteUUIDFluxMonitor,
			Status:         feeds.JobProposalStatusPending,
		}
		specFluxMonitor = feeds.JobProposalSpec{
			Definition:    spec,
			Status:        feeds.SpecStatusPending,
			Version:       argsFluxMonitor.Version,
			JobProposalID: idFluxMonitor,
		}

		idOCR1                   = int64(2)
		remoteUUIDOCR1           = uuid.New()
		ocr1NameAndExternalJobID = uuid.New()
		ocr1Spec                 = fmt.Sprintf(OCR1TestSpecTemplate, ocr1NameAndExternalJobID, ocr1NameAndExternalJobID)
		argsOCR1                 = &feeds.ProposeJobArgs{
			FeedsManagerID: 1,
			RemoteUUID:     remoteUUIDOCR1,
			Spec:           ocr1Spec,
			Version:        1,
		}
		jpOCR1 = feeds.JobProposal{
			FeedsManagerID: 1,
			Name:           null.StringFrom(ocr1NameAndExternalJobID.String()),
			RemoteUUID:     remoteUUIDOCR1,
			Status:         feeds.JobProposalStatusPending,
		}
		specOCR1 = feeds.JobProposalSpec{
			Definition:    ocr1Spec,
			Status:        feeds.SpecStatusPending,
			Version:       argsOCR1.Version,
			JobProposalID: idOCR1,
		}

		idOCR2                   = int64(3)
		remoteUUIDOCR2           = uuid.New()
		ocr2NameAndExternalJobID = uuid.New()
		ocr2Spec                 = fmt.Sprintf(OCR2TestSpecTemplate, ocr2NameAndExternalJobID, ocr2NameAndExternalJobID)
		argsOCR2                 = &feeds.ProposeJobArgs{
			FeedsManagerID: 1,
			RemoteUUID:     remoteUUIDOCR2,
			Spec:           ocr2Spec,
			Version:        1,
		}
		jpOCR2 = feeds.JobProposal{
			FeedsManagerID: 1,
			Name:           null.StringFrom(ocr2NameAndExternalJobID.String()),
			RemoteUUID:     remoteUUIDOCR2,
			Status:         feeds.JobProposalStatusPending,
		}
		specOCR2 = feeds.JobProposalSpec{
			Definition:    ocr2Spec,
			Status:        feeds.SpecStatusPending,
			Version:       argsOCR2.Version,
			JobProposalID: idOCR2,
		}

		idBootstrap         = int64(4)
		remoteUUIDBootstrap = uuid.New()
		bootstrapName       = uuid.New()
		bootstrapSpec       = fmt.Sprintf(BootstrapTestSpecTemplate, bootstrapName)
		argsBootstrap       = &feeds.ProposeJobArgs{
			FeedsManagerID: 1,
			RemoteUUID:     remoteUUIDBootstrap,
			Spec:           bootstrapSpec,
			Version:        1,
		}
		jpBootstrap = feeds.JobProposal{
			FeedsManagerID: 1,
			Name:           null.StringFrom(bootstrapName.String()),
			RemoteUUID:     remoteUUIDBootstrap,
			Status:         feeds.JobProposalStatusPending,
		}
		specBootstrap = feeds.JobProposalSpec{
			Definition:    bootstrapSpec,
			Status:        feeds.SpecStatusPending,
			Version:       argsBootstrap.Version,
			JobProposalID: idBootstrap,
		}

		httpTimeout = *commonconfig.MustNewDuration(1 * time.Second)

		// variables for workflow spec
		wfJobSpec           = testspecs.DefaultWorkflowJobSpec(t)
		proposalIDWF        = int64(11)
		jobProposalSpecIdWF = int64(101)
		jobIDWF             = int32(1001)
		remoteUUIDWF        = uuid.New()
		argsWF              = &feeds.ProposeJobArgs{
			FeedsManagerID: 1,
			RemoteUUID:     remoteUUIDWF,
			Spec:           wfJobSpec.Toml(),
			Version:        1,
		}
		jpWF = feeds.JobProposal{
			FeedsManagerID: 1,
			Name:           null.StringFrom("test-spec"),
			RemoteUUID:     remoteUUIDWF,
			Status:         feeds.JobProposalStatusPending,
		}
		acceptedjpWF = feeds.JobProposal{
			ID:             13,
			FeedsManagerID: 1,
			Name:           null.StringFrom("test-spec"),
			RemoteUUID:     remoteUUIDWF,
			Status:         feeds.JobProposalStatusPending,
		}
		proposalSpecWF = feeds.JobProposalSpec{
			Definition:    wfJobSpec.Toml(),
			Status:        feeds.SpecStatusPending,
			Version:       1,
			JobProposalID: proposalIDWF,
		}
		autoApprovableProposalSpecWF = feeds.JobProposalSpec{
			ID:            jobProposalSpecIdWF,
			Definition:    wfJobSpec.Toml(),
			Status:        feeds.SpecStatusPending,
			Version:       1,
			JobProposalID: proposalIDWF,
		}
	)

	testCases := []struct {
		name    string
		args    *feeds.ProposeJobArgs
		before  func(svc *TestService)
		wantID  int64
		wantErr string
	}{
		{
			name: "Auto approve new WF spec",
			before: func(svc *TestService) {
				svc.orm.On("GetJobProposalByRemoteUUID", mock.Anything, argsWF.RemoteUUID).Return(new(feeds.JobProposal), sql.ErrNoRows)
				svc.orm.On("UpsertJobProposal", mock.Anything, &jpWF).Return(proposalIDWF, nil)
				svc.orm.On("CreateSpec", mock.Anything, proposalSpecWF).Return(jobProposalSpecIdWF, nil)
				svc.orm.On("CountJobProposalsByStatus", mock.Anything).Return(&feeds.JobProposalCounts{}, nil)
				transactCall := svc.orm.On("Transact", mock.Anything, mock.Anything)
				transactCall.Run(func(args mock.Arguments) {
					fn := args[1].(func(orm feeds.ORM) error)
					transactCall.ReturnArguments = mock.Arguments{fn(svc.orm)}
				})
				// Auto approve is really a call to ApproveJobProposal and so we have to mock that as well
				svc.connMgr.On("GetClient", argsWF.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.orm.EXPECT().GetSpec(mock.Anything, jobProposalSpecIdWF).Return(&autoApprovableProposalSpecWF, nil)
				svc.orm.EXPECT().GetJobProposal(mock.Anything, autoApprovableProposalSpecWF.JobProposalID).Return(&acceptedjpWF, nil)
				svc.jobORM.On("AssertBridgesExist", mock.Anything, mock.IsType(pipeline.Pipeline{})).Return(nil)

				svc.jobORM.On("FindJobByExternalJobID", mock.Anything, mock.Anything).Return(job.Job{}, sql.ErrNoRows)
				svc.orm.On("WithDataSource", mock.Anything).Return(feeds.ORM(svc.orm))
				svc.jobORM.On("WithDataSource", mock.Anything).Return(job.ORM(svc.jobORM))
				svc.jobORM.On("FindJobIDByWorkflow", mock.Anything, mock.Anything).Return(int32(0), sql.ErrNoRows) // no existing job
				svc.spawner.
					On("CreateJob",
						mock.Anything,
						mock.Anything,
						mock.MatchedBy(func(j *job.Job) bool {
							match := j.WorkflowSpec.Workflow == wfJobSpec.Job().WorkflowSpec.Workflow
							if !match {
								t.Logf("got wf spec %s want %s", j.WorkflowSpec.Workflow, wfJobSpec.Job().WorkflowSpec.Workflow)
							}
							return match
						}),
					).
					Run(func(args mock.Arguments) { (args.Get(2).(*job.Job)).ID = 1 }).
					Return(nil)
				svc.orm.On("ApproveSpec",
					mock.Anything,
					jobProposalSpecIdWF,
					mock.IsType(uuid.UUID{}),
				).Return(nil)
				svc.fmsClient.On("ApprovedJob",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&proto.ApprovedJobRequest{
						Uuid:    jpWF.RemoteUUID.String(),
						Version: int64(proposalSpecWF.Version),
					},
				).Return(&proto.ApprovedJobResponse{}, nil)
			},
			args:   argsWF,
			wantID: proposalIDWF,
		},

		{
			name: "Auto approve existing WF spec found by FindJobIDByWorkflow",
			before: func(svc *TestService) {
				svc.orm.On("GetJobProposalByRemoteUUID", mock.Anything, argsWF.RemoteUUID).Return(new(feeds.JobProposal), sql.ErrNoRows)
				svc.orm.On("UpsertJobProposal", mock.Anything, &jpWF).Return(proposalIDWF, nil)
				svc.orm.On("CreateSpec", mock.Anything, proposalSpecWF).Return(jobProposalSpecIdWF, nil)
				svc.orm.On("CountJobProposalsByStatus", mock.Anything).Return(&feeds.JobProposalCounts{}, nil)
				transactCall := svc.orm.On("Transact", mock.Anything, mock.Anything)
				transactCall.Run(func(args mock.Arguments) {
					fn := args[1].(func(orm feeds.ORM) error)
					transactCall.ReturnArguments = mock.Arguments{fn(svc.orm)}
				})
				// Auto approve is really a call to ApproveJobProposal and so we have to mock that as well
				svc.connMgr.On("GetClient", argsWF.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.orm.EXPECT().GetSpec(mock.Anything, jobProposalSpecIdWF).Return(&autoApprovableProposalSpecWF, nil)
				svc.orm.EXPECT().GetJobProposal(mock.Anything, autoApprovableProposalSpecWF.JobProposalID).Return(&acceptedjpWF, nil)
				svc.jobORM.On("AssertBridgesExist", mock.Anything, mock.IsType(pipeline.Pipeline{})).Return(nil)

				svc.jobORM.On("FindJobByExternalJobID", mock.Anything, mock.Anything).Return(job.Job{}, sql.ErrNoRows)
				svc.orm.On("WithDataSource", mock.Anything).Return(feeds.ORM(svc.orm))
				svc.jobORM.On("WithDataSource", mock.Anything).Return(job.ORM(svc.jobORM))
				svc.jobORM.On("FindJobIDByWorkflow", mock.Anything, mock.Anything).Return(jobIDWF, sql.ErrNoRows)
				svc.orm.On("GetApprovedSpec", mock.Anything, acceptedjpWF.ID).Return(&autoApprovableProposalSpecWF, nil)
				svc.orm.On("CancelSpec", mock.Anything, autoApprovableProposalSpecWF.ID).Return(nil)
				svc.spawner.On("DeleteJob", mock.Anything, mock.Anything, jobIDWF).Return(nil)
				svc.spawner.
					On("CreateJob",
						mock.Anything,
						mock.Anything,
						mock.MatchedBy(func(j *job.Job) bool {
							match := j.WorkflowSpec.Workflow == wfJobSpec.Job().WorkflowSpec.Workflow
							if !match {
								t.Logf("got wf spec %s want %s", j.WorkflowSpec.Workflow, wfJobSpec.Job().WorkflowSpec.Workflow)
							}
							return match
						}),
					).
					Run(func(args mock.Arguments) { (args.Get(2).(*job.Job)).ID = 1 }).
					Return(nil)
				svc.orm.On("ApproveSpec",
					mock.Anything,
					jobProposalSpecIdWF,
					mock.IsType(uuid.UUID{}),
				).Return(nil)
				svc.fmsClient.On("ApprovedJob",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&proto.ApprovedJobRequest{
						Uuid:    jpWF.RemoteUUID.String(),
						Version: int64(proposalSpecWF.Version),
					},
				).Return(&proto.ApprovedJobResponse{}, nil)
			},
			args:   argsWF,
			wantID: proposalIDWF,
		},

		{
			name: "Auto approve WF spec: error creating job for new spec",
			before: func(svc *TestService) {
				svc.orm.On("GetJobProposalByRemoteUUID", mock.Anything, argsWF.RemoteUUID).Return(new(feeds.JobProposal), sql.ErrNoRows)
				svc.orm.On("UpsertJobProposal", mock.Anything, &jpWF).Return(proposalIDWF, nil)
				svc.orm.On("CreateSpec", mock.Anything, proposalSpecWF).Return(jobProposalSpecIdWF, nil)
				transactCall := svc.orm.On("Transact", mock.Anything, mock.Anything)
				transactCall.Run(func(args mock.Arguments) {
					fn := args[1].(func(orm feeds.ORM) error)
					transactCall.ReturnArguments = mock.Arguments{fn(svc.orm)}
				})
				// Auto approve is really a call to ApproveJobProposal and so we have to mock that as well
				svc.connMgr.On("GetClient", argsWF.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.orm.EXPECT().GetSpec(mock.Anything, jobProposalSpecIdWF).Return(&proposalSpecWF, nil)
				svc.orm.EXPECT().GetJobProposal(mock.Anything, proposalSpecWF.JobProposalID).Return(&jpWF, nil)
				svc.jobORM.On("AssertBridgesExist", mock.Anything, mock.IsType(pipeline.Pipeline{})).Return(nil)

				svc.jobORM.On("FindJobByExternalJobID", mock.Anything, mock.Anything).Return(job.Job{}, sql.ErrNoRows)
				svc.orm.On("WithDataSource", mock.Anything).Return(feeds.ORM(svc.orm))
				svc.jobORM.On("WithDataSource", mock.Anything).Return(job.ORM(svc.jobORM))
				svc.jobORM.On("FindJobIDByWorkflow", mock.Anything, mock.Anything).Return(int32(0), sql.ErrNoRows) // no existing job
				svc.spawner.
					On("CreateJob",
						mock.Anything,
						mock.Anything,
						mock.MatchedBy(func(j *job.Job) bool {
							match := j.WorkflowSpec.Workflow == wfJobSpec.Job().WorkflowSpec.Workflow
							if !match {
								t.Logf("got wf spec %s want %s", j.WorkflowSpec.Workflow, wfJobSpec.Job().WorkflowSpec.Workflow)
							}
							return match
						}),
					).
					Run(func(args mock.Arguments) { (args.Get(2).(*job.Job)).ID = 1 }).
					Return(errors.New("error creating job"))
			},
			args:    argsWF,
			wantID:  0,
			wantErr: "error creating job",
		},

		{
			name: "Create success (Flux Monitor)",
			before: func(svc *TestService) {
				svc.orm.On("GetJobProposalByRemoteUUID", mock.Anything, jpFluxMonitor.RemoteUUID).Return(new(feeds.JobProposal), sql.ErrNoRows)
				svc.orm.On("UpsertJobProposal", mock.Anything, &jpFluxMonitor).Return(idFluxMonitor, nil)
				svc.orm.On("CreateSpec", mock.Anything, specFluxMonitor).Return(int64(100), nil)
				svc.orm.On("CountJobProposalsByStatus", mock.Anything).Return(&feeds.JobProposalCounts{}, nil)
				transactCall := svc.orm.On("Transact", mock.Anything, mock.Anything)
				transactCall.Run(func(args mock.Arguments) {
					fn := args[1].(func(orm feeds.ORM) error)
					transactCall.ReturnArguments = mock.Arguments{fn(svc.orm)}
				})
			},
			args:   argsFluxMonitor,
			wantID: idFluxMonitor,
		},
		{
			name: "Create success (OCR1)",
			before: func(svc *TestService) {
				svc.orm.On("GetJobProposalByRemoteUUID", mock.Anything, jpOCR1.RemoteUUID).Return(new(feeds.JobProposal), sql.ErrNoRows)
				svc.orm.On("UpsertJobProposal", mock.Anything, &jpOCR1).Return(idOCR1, nil)
				svc.orm.On("CreateSpec", mock.Anything, specOCR1).Return(int64(100), nil)
				svc.orm.On("CountJobProposalsByStatus", mock.Anything).Return(&feeds.JobProposalCounts{}, nil)
				transactCall := svc.orm.On("Transact", mock.Anything, mock.Anything)
				transactCall.Run(func(args mock.Arguments) {
					fn := args[1].(func(orm feeds.ORM) error)
					transactCall.ReturnArguments = mock.Arguments{fn(svc.orm)}
				})
			},
			args:   argsOCR1,
			wantID: idOCR1,
		},
		{
			name: "Create success (OCR2)",
			before: func(svc *TestService) {
				svc.orm.On("GetJobProposalByRemoteUUID", mock.Anything, jpOCR2.RemoteUUID).Return(new(feeds.JobProposal), sql.ErrNoRows)
				svc.orm.On("UpsertJobProposal", mock.Anything, &jpOCR2).Return(idOCR2, nil)
				svc.orm.On("CreateSpec", mock.Anything, specOCR2).Return(int64(100), nil)
				svc.orm.On("CountJobProposalsByStatus", mock.Anything).Return(&feeds.JobProposalCounts{}, nil)
				transactCall := svc.orm.On("Transact", mock.Anything, mock.Anything)
				transactCall.Run(func(args mock.Arguments) {
					fn := args[1].(func(orm feeds.ORM) error)
					transactCall.ReturnArguments = mock.Arguments{fn(svc.orm)}
				})
			},
			args:   argsOCR2,
			wantID: idOCR2,
		},
		{
			name: "Create success (Bootstrap)",
			before: func(svc *TestService) {
				svc.orm.On("GetJobProposalByRemoteUUID", mock.Anything, jpBootstrap.RemoteUUID).Return(new(feeds.JobProposal), sql.ErrNoRows)
				svc.orm.On("UpsertJobProposal", mock.Anything, &jpBootstrap).Return(idBootstrap, nil)
				svc.orm.On("CreateSpec", mock.Anything, specBootstrap).Return(int64(102), nil)
				svc.orm.On("CountJobProposalsByStatus", mock.Anything).Return(&feeds.JobProposalCounts{}, nil)
				transactCall := svc.orm.On("Transact", mock.Anything, mock.Anything)
				transactCall.Run(func(args mock.Arguments) {
					fn := args[1].(func(orm feeds.ORM) error)
					transactCall.ReturnArguments = mock.Arguments{fn(svc.orm)}
				})
			},
			args:   argsBootstrap,
			wantID: idBootstrap,
		},
		{
			name: "Update success",
			before: func(svc *TestService) {
				svc.orm.
					On("GetJobProposalByRemoteUUID", mock.Anything, jpFluxMonitor.RemoteUUID).
					Return(&feeds.JobProposal{
						FeedsManagerID: jpFluxMonitor.FeedsManagerID,
						RemoteUUID:     jpFluxMonitor.RemoteUUID,
						Status:         feeds.JobProposalStatusPending,
					}, nil)
				svc.orm.On("ExistsSpecByJobProposalIDAndVersion", mock.Anything, jpFluxMonitor.ID, argsFluxMonitor.Version).Return(false, nil)
				svc.orm.On("UpsertJobProposal", mock.Anything, &jpFluxMonitor).Return(idFluxMonitor, nil)
				svc.orm.On("CreateSpec", mock.Anything, specFluxMonitor).Return(int64(100), nil)
				svc.orm.On("CountJobProposalsByStatus", mock.Anything).Return(&feeds.JobProposalCounts{}, nil)
				transactCall := svc.orm.On("Transact", mock.Anything, mock.Anything)
				transactCall.Run(func(args mock.Arguments) {
					fn := args[1].(func(orm feeds.ORM) error)
					transactCall.ReturnArguments = mock.Arguments{fn(svc.orm)}
				})
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
				Spec:       spec,
				Multiaddrs: pq.StringArray{"/dns4/example.com"},
			},
			wantErr: "only OCR job type supports multiaddr",
		},
		{
			name: "ensure an upsert validates the job proposal belongs to the feeds manager",
			before: func(svc *TestService) {
				svc.orm.
					On("GetJobProposalByRemoteUUID", mock.Anything, jpFluxMonitor.RemoteUUID).
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
					On("GetJobProposalByRemoteUUID", mock.Anything, jpFluxMonitor.RemoteUUID).
					Return(&feeds.JobProposal{
						FeedsManagerID: jpFluxMonitor.FeedsManagerID,
						RemoteUUID:     jpFluxMonitor.RemoteUUID,
						Status:         feeds.JobProposalStatusPending,
					}, nil)
				svc.orm.On("ExistsSpecByJobProposalIDAndVersion", mock.Anything, jpFluxMonitor.ID, argsFluxMonitor.Version).Return(true, nil)
			},
			args:    argsFluxMonitor,
			wantErr: "version conflict",
		},
		{
			name: "upsert error",
			before: func(svc *TestService) {
				svc.orm.On("GetJobProposalByRemoteUUID", mock.Anything, jpFluxMonitor.RemoteUUID).Return(new(feeds.JobProposal), sql.ErrNoRows)
				svc.orm.On("UpsertJobProposal", mock.Anything, &jpFluxMonitor).Return(int64(0), errors.New("orm error"))
				transactCall := svc.orm.On("Transact", mock.Anything, mock.Anything)
				transactCall.Run(func(args mock.Arguments) {
					fn := args[1].(func(orm feeds.ORM) error)
					transactCall.ReturnArguments = mock.Arguments{fn(svc.orm)}
				})
			},
			args:    argsFluxMonitor,
			wantErr: "failed to upsert job proposal",
		},
		{
			name: "Create spec error",
			before: func(svc *TestService) {
				svc.orm.On("GetJobProposalByRemoteUUID", mock.Anything, jpFluxMonitor.RemoteUUID).Return(new(feeds.JobProposal), sql.ErrNoRows)
				svc.orm.On("UpsertJobProposal", mock.Anything, &jpFluxMonitor).Return(idFluxMonitor, nil)
				svc.orm.On("CreateSpec", mock.Anything, specFluxMonitor).Return(int64(0), errors.New("orm error"))
				transactCall := svc.orm.On("Transact", mock.Anything, mock.Anything)
				transactCall.Run(func(args mock.Arguments) {
					fn := args[1].(func(orm feeds.ORM) error)
					transactCall.ReturnArguments = mock.Arguments{fn(svc.orm)}
				})
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
				c.OCR.Enabled = testutils.Ptr(true)
				c.OCR2.Enabled = testutils.Ptr(true)
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
		remoteUUID = uuid.New()
		args       = &feeds.DeleteJobArgs{
			FeedsManagerID: 1,
			RemoteUUID:     remoteUUID,
		}

		approved = feeds.JobProposal{
			ID:             321,
			FeedsManagerID: 1,
			RemoteUUID:     remoteUUID,
			ExternalJobID:  uuid.NullUUID{UUID: uuid.New(), Valid: true},
			Status:         feeds.JobProposalStatusApproved,
		}

		wfSpecID    = int32(4321)
		workflowJob = job.Job{
			ID:             1,
			WorkflowSpecID: &wfSpecID,
		}
		jobProposalSpec = &feeds.JobProposalSpec{
			ID:            20,
			Status:        feeds.SpecStatusApproved,
			JobProposalID: approved.ID,
			Version:       1,
		}

		httpTimeout = *commonconfig.MustNewDuration(1 * time.Second)
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
				svc.orm.On("GetJobProposalByRemoteUUID", mock.Anything, approved.RemoteUUID).Return(&approved, nil)
				svc.orm.On("DeleteProposal", mock.Anything, approved.ID).Return(nil)
				svc.orm.On("CountJobProposalsByStatus", mock.Anything).Return(&feeds.JobProposalCounts{}, nil)
				svc.jobORM.On("FindJobByExternalJobID", mock.Anything, approved.ExternalJobID.UUID).Return(job.Job{}, sql.ErrNoRows)
			},
			args:   args,
			wantID: approved.ID,
		},
		{
			name: "Job proposal being deleted belongs to the feeds manager",
			before: func(svc *TestService) {
				svc.orm.
					On("GetJobProposalByRemoteUUID", mock.Anything, approved.RemoteUUID).
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
				svc.orm.On("GetJobProposalByRemoteUUID", mock.Anything, approved.RemoteUUID).Return(nil, errors.New("orm error"))
			},
			args:    args,
			wantErr: "GetJobProposalByRemoteUUID failed",
		},
		{
			name: "No proposal error",
			before: func(svc *TestService) {
				svc.orm.On("GetJobProposalByRemoteUUID", mock.Anything, approved.RemoteUUID).Return(nil, sql.ErrNoRows)
			},
			args:    args,
			wantErr: "GetJobProposalByRemoteUUID did not find any proposals to delete",
		},
		{
			name: "Delete proposal error",
			before: func(svc *TestService) {
				svc.orm.On("GetJobProposalByRemoteUUID", mock.Anything, approved.RemoteUUID).Return(&approved, nil)
				svc.orm.On("DeleteProposal", mock.Anything, approved.ID).Return(errors.New("orm error"))
			},
			args:    args,
			wantErr: "DeleteProposal failed",
		},
		{
			name: "Delete workflow-spec with auto-cancellation",
			before: func(svc *TestService) {
				svc.orm.On("GetJobProposalByRemoteUUID", mock.Anything, approved.RemoteUUID).Return(&approved, nil)
				svc.orm.On("DeleteProposal", mock.Anything, approved.ID).Return(nil)
				svc.orm.On("CountJobProposalsByStatus", mock.Anything).Return(&feeds.JobProposalCounts{}, nil)
				svc.jobORM.On("FindJobByExternalJobID", mock.Anything, approved.ExternalJobID.UUID).Return(workflowJob, nil)
				svc.orm.On("GetApprovedSpec", mock.Anything, approved.ID).Return(jobProposalSpec, nil)

				// mocks for CancelSpec()
				svc.orm.On("GetSpec", mock.Anything, jobProposalSpec.ID).Return(jobProposalSpec, nil)
				svc.orm.On("GetJobProposal", mock.Anything, approved.ID).Return(&approved, nil)
				svc.connMgr.On("GetClient", mock.Anything).Return(svc.fmsClient, nil)

				svc.orm.On("CancelSpec", mock.Anything, jobProposalSpec.ID).Return(nil)
				svc.jobORM.On("FindJobByExternalJobID", mock.Anything, approved.ExternalJobID.UUID).Return(workflowJob, nil)
				svc.spawner.On("DeleteJob", mock.Anything, mock.Anything, workflowJob.ID).Return(nil)

				svc.fmsClient.On("CancelledJob",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&proto.CancelledJobRequest{
						Uuid:    approved.RemoteUUID.String(),
						Version: int64(jobProposalSpec.Version),
					},
				).Return(&proto.CancelledJobResponse{}, nil)
				svc.orm.On("CountJobProposalsByStatus", mock.Anything).Return(&feeds.JobProposalCounts{}, nil)
				svc.orm.On("WithDataSource", mock.Anything).Return(feeds.ORM(svc.orm))
				svc.jobORM.On("WithDataSource", mock.Anything).Return(job.ORM(svc.jobORM))
			},
			args:   args,
			wantID: approved.ID,
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
		remoteUUID = uuid.New()
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
# omit gasPriceSubunitsSource intentionally 
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

		httpTimeout = *commonconfig.MustNewDuration(1 * time.Second)
	)

	testCases := []struct {
		name    string
		args    *feeds.RevokeJobArgs
		before  func(svc *TestService)
		wantID  int64
		wantErr string
	}{
		{
			name: "Revoke success when latest spec status is pending",
			before: func(svc *TestService) {
				svc.orm.On("GetJobProposalByRemoteUUID", mock.Anything, pendingProposal.RemoteUUID).Return(pendingProposal, nil)
				svc.orm.On("GetLatestSpec", mock.Anything, pendingSpec.JobProposalID).Return(pendingSpec, nil)
				svc.orm.On("RevokeSpec", mock.Anything, pendingSpec.ID).Return(nil)
				svc.orm.On("CountJobProposalsByStatus", mock.Anything).Return(&feeds.JobProposalCounts{}, nil)
			},
			args:   args,
			wantID: pendingProposal.ID,
		},
		{
			name: "Revoke success when latest spec status is cancelled",
			before: func(svc *TestService) {
				svc.orm.On("GetJobProposalByRemoteUUID", mock.Anything, pendingProposal.RemoteUUID).Return(pendingProposal, nil)
				svc.orm.On("GetLatestSpec", mock.Anything, pendingSpec.JobProposalID).Return(&feeds.JobProposalSpec{
					ID:            20,
					Status:        feeds.SpecStatusCancelled,
					JobProposalID: pendingProposal.ID,
					Version:       1,
					Definition:    defn,
				}, nil)
				svc.orm.On("RevokeSpec", mock.Anything, pendingSpec.ID).Return(nil)
				svc.orm.On("CountJobProposalsByStatus", mock.Anything).Return(&feeds.JobProposalCounts{}, nil)
			},
			args:   args,
			wantID: pendingProposal.ID,
		},
		{
			name: "Job proposal being revoked belongs to the feeds manager",
			before: func(svc *TestService) {
				svc.orm.
					On("GetJobProposalByRemoteUUID", mock.Anything, pendingProposal.RemoteUUID).
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
				svc.orm.On("GetJobProposalByRemoteUUID", mock.Anything, pendingProposal.RemoteUUID).Return(nil, errors.New("orm error"))
			},
			args:    args,
			wantErr: "GetJobProposalByRemoteUUID failed",
		},
		{
			name: "No proposal error",
			before: func(svc *TestService) {
				svc.orm.On("GetJobProposalByRemoteUUID", mock.Anything, pendingProposal.RemoteUUID).Return(nil, sql.ErrNoRows)
			},
			args:    args,
			wantErr: "GetJobProposalByRemoteUUID did not find any proposals to revoke",
		},
		{
			name: "Get latest spec error",
			before: func(svc *TestService) {
				svc.orm.On("GetJobProposalByRemoteUUID", mock.Anything, pendingProposal.RemoteUUID).Return(pendingProposal, nil)
				svc.orm.On("GetLatestSpec", mock.Anything, pendingSpec.JobProposalID).Return(nil, sql.ErrNoRows)
			},
			args:    args,
			wantErr: "GetLatestSpec failed to get latest spec",
		},
		{
			name: "Not revokable due to spec status approved",
			before: func(svc *TestService) {
				svc.orm.On("GetJobProposalByRemoteUUID", mock.Anything, pendingProposal.RemoteUUID).Return(pendingProposal, nil)
				svc.orm.On("GetLatestSpec", mock.Anything, pendingSpec.JobProposalID).Return(&feeds.JobProposalSpec{
					ID:            20,
					Status:        feeds.SpecStatusApproved,
					JobProposalID: pendingProposal.ID,
					Version:       1,
					Definition:    defn,
				}, nil)
			},
			args:    args,
			wantErr: "only pending job specs can be revoked",
		},
		{
			name: "Not revokable due to spec status rejected",
			before: func(svc *TestService) {
				svc.orm.On("GetJobProposalByRemoteUUID", mock.Anything, pendingProposal.RemoteUUID).Return(pendingProposal, nil)
				svc.orm.On("GetLatestSpec", mock.Anything, pendingSpec.JobProposalID).Return(&feeds.JobProposalSpec{
					ID:            20,
					Status:        feeds.SpecStatusRejected,
					JobProposalID: pendingProposal.ID,
					Version:       1,
					Definition:    defn,
				}, nil)
			},
			args:    args,
			wantErr: "only pending job specs can be revoked",
		},
		{
			name: "Not revokable due to spec status already revoked",
			before: func(svc *TestService) {
				svc.orm.On("GetJobProposalByRemoteUUID", mock.Anything, pendingProposal.RemoteUUID).Return(pendingProposal, nil)
				svc.orm.On("GetLatestSpec", mock.Anything, pendingSpec.JobProposalID).Return(&feeds.JobProposalSpec{
					ID:            20,
					Status:        feeds.SpecStatusRevoked,
					JobProposalID: pendingProposal.ID,
					Version:       1,
					Definition:    defn,
				}, nil)
			},
			args:    args,
			wantErr: "only pending job specs can be revoked",
		},
		{
			name: "Not revokable due to proposal status deleted",
			before: func(svc *TestService) {
				svc.orm.On("GetJobProposalByRemoteUUID", mock.Anything, pendingProposal.RemoteUUID).Return(&feeds.JobProposal{
					ID:             1,
					FeedsManagerID: 1,
					RemoteUUID:     remoteUUID,
					Status:         feeds.JobProposalStatusDeleted,
				}, nil)
				svc.orm.On("GetLatestSpec", mock.Anything, pendingSpec.JobProposalID).Return(pendingSpec, nil)
			},
			args:    args,
			wantErr: "only pending job specs can be revoked",
		},
		{
			name: "Revoke proposal error",
			before: func(svc *TestService) {
				svc.orm.On("GetJobProposalByRemoteUUID", mock.Anything, pendingProposal.RemoteUUID).Return(pendingProposal, nil)
				svc.orm.On("GetLatestSpec", mock.Anything, pendingSpec.JobProposalID).Return(pendingSpec, nil)
				svc.orm.On("RevokeSpec", mock.Anything, pendingSpec.ID).Return(errors.New("orm error"))
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
				c.OCR2.Enabled = testutils.Ptr(true)
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
	tests := []struct {
		name      string
		chainType feeds.ChainType
		protoType proto.ChainType
	}{
		{
			name:      "EVM Chain Type",
			chainType: feeds.ChainTypeEVM,
			protoType: proto.ChainType_CHAIN_TYPE_EVM,
		},
		{
			name:      "Solana Chain Type",
			chainType: feeds.ChainTypeSolana,
			protoType: proto.ChainType_CHAIN_TYPE_SOLANA,
		},
		{
			name:      "Starknet Chain Type",
			chainType: feeds.ChainTypeStarknet,
			protoType: proto.ChainType_CHAIN_TYPE_STARKNET,
		},
		{
			name:      "Aptos Chain Type",
			chainType: feeds.ChainTypeAptos,
			protoType: proto.ChainType_CHAIN_TYPE_APTOS,
		},
		{
			name:      "Tron Chain Type",
			chainType: feeds.ChainTypeTron,
			protoType: proto.ChainType_CHAIN_TYPE_TRON,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p2pKey1 := p2pkey.MustNewV2XXXTestingOnly(big.NewInt(1))
			p2pKey2 := p2pkey.MustNewV2XXXTestingOnly(big.NewInt(2))
			p2pKey3 := p2pkey.MustNewV2XXXTestingOnly(big.NewInt(3))

			ocrKey, err := ocrkey.NewV2()
			require.NoError(t, err)

			workflowKey, err := workflowkey.New()
			require.NoError(t, err)

			var (
				multiaddr     = "/dns4/chain.link/tcp/1234/p2p/16Uiu2HAm58SP7UL8zsnpeuwHfytLocaqgnyaYKP8wu7qRdrixLju"
				mgr           = &feeds.FeedsManager{ID: 1}
				forwarderAddr = "0x0002"
				ccfg          = feeds.ChainConfig{
					ID:             100,
					FeedsManagerID: mgr.ID,
					ChainID:        "42",
					ChainType:      tt.chainType,
					AccountAddress: "0x0000",
					AdminAddress:   "0x0001",
					FluxMonitorConfig: feeds.FluxMonitorConfig{
						Enabled: true,
					},
					OCR1Config: feeds.OCR1Config{
						Enabled:     true,
						IsBootstrap: false,
						P2PPeerID:   null.StringFrom(p2pKey1.PeerID().String()),
						KeyBundleID: null.StringFrom(ocrKey.GetID()),
					},
					OCR2Config: feeds.OCR2ConfigModel{
						Enabled:          true,
						IsBootstrap:      true,
						Multiaddr:        null.StringFrom(multiaddr),
						ForwarderAddress: null.StringFrom(forwarderAddr),
						Plugins: feeds.Plugins{
							Commit:     true,
							Execute:    true,
							Median:     false,
							Mercury:    true,
							Rebalancer: true,
						},
					},
				}
				chainConfigs = []feeds.ChainConfig{ccfg}
				nodeVersion  = &versioning.NodeVersion{Version: "1.0.0"}
			)

			svc := setupTestService(t)

			svc.connMgr.On("GetClient", mgr.ID).Return(svc.fmsClient, nil)
			svc.orm.On("ListChainConfigsByManagerIDs", mock.Anything, []int64{mgr.ID}).Return(chainConfigs, nil)

			// OCR1 key fetching
			svc.p2pKeystore.On("Get", p2pKey1.PeerID()).Return(p2pKey1, nil)
			svc.ocr1Keystore.On("Get", ocrKey.GetID()).Return(ocrKey, nil)

			svc.workflowKeystore.On("EnsureKey", mock.Anything).Return(nil)
			svc.workflowKeystore.On("GetAll").Return([]workflowkey.Key{workflowKey}, nil)
			svc.p2pKeystore.On("GetAll").Return([]p2pkey.KeyV2{p2pKey1, p2pKey2, p2pKey3}, nil)
			wkID := workflowKey.ID()
			svc.fmsClient.On("UpdateNode", mock.Anything, &proto.UpdateNodeRequest{
				Version: nodeVersion.Version,
				ChainConfigs: []*proto.ChainConfig{
					{
						Chain: &proto.Chain{
							Id:   ccfg.ChainID,
							Type: tt.protoType,
						},
						AccountAddress:    ccfg.AccountAddress,
						AdminAddress:      ccfg.AdminAddress,
						FluxMonitorConfig: &proto.FluxMonitorConfig{Enabled: true},
						Ocr1Config: &proto.OCR1Config{
							Enabled:     true,
							IsBootstrap: ccfg.OCR1Config.IsBootstrap,
							P2PKeyBundle: &proto.OCR1Config_P2PKeyBundle{
								PeerId:    p2pKey1.PeerID().String(),
								PublicKey: p2pKey1.PublicKeyHex(),
							},
							OcrKeyBundle: &proto.OCR1Config_OCRKeyBundle{
								BundleId:              ocrKey.GetID(),
								ConfigPublicKey:       ocrkey.ConfigPublicKey(ocrKey.PublicKeyConfig()).String(),
								OffchainPublicKey:     ocrKey.OffChainSigning.PublicKey().String(),
								OnchainSigningAddress: ocrKey.OnChainSigning.Address().String(),
							},
						},
						Ocr2Config: &proto.OCR2Config{
							Enabled:          true,
							IsBootstrap:      ccfg.OCR2Config.IsBootstrap,
							Multiaddr:        multiaddr,
							ForwarderAddress: &forwarderAddr,
							Plugins: &proto.OCR2Config_Plugins{
								Commit:     ccfg.OCR2Config.Plugins.Commit,
								Execute:    ccfg.OCR2Config.Plugins.Execute,
								Median:     ccfg.OCR2Config.Plugins.Median,
								Mercury:    ccfg.OCR2Config.Plugins.Mercury,
								Rebalancer: ccfg.OCR2Config.Plugins.Rebalancer,
							},
						},
					},
				},
				WorkflowKey: &wkID,
				P2PKeyBundles: []*proto.P2PKeyBundle{
					{PeerId: p2pKey1.PeerID().String(), PublicKey: p2pKey1.PublicKeyHex()},
					{PeerId: p2pKey2.PeerID().String(), PublicKey: p2pKey2.PublicKeyHex()},
					{PeerId: p2pKey3.PeerID().String(), PublicKey: p2pKey3.PublicKeyHex()},
				},
			}).Return(&proto.UpdateNodeResponse{}, nil)

			err = svc.SyncNodeInfo(testutils.Context(t), mgr.ID)
			require.NoError(t, err)
		})
	}
}

func Test_Service_syncNodeInfoWithRetry(t *testing.T) {
	t.Parallel()

	mgr := feeds.FeedsManager{ID: 1}
	nodeVersion := &versioning.NodeVersion{Version: "1.0.0"}
	cfg := feeds.ChainConfig{
		FeedsManagerID:          mgr.ID,
		ChainID:                 "42",
		ChainType:               feeds.ChainTypeEVM,
		AccountAddress:          "0x0000000000000000000000000000000000000000",
		AccountAddressPublicKey: null.StringFrom("0x0000000000000000000000000000000000000002"),
		AdminAddress:            "0x0000000000000000000000000000000000000001",
		FluxMonitorConfig:       feeds.FluxMonitorConfig{Enabled: true},
		OCR1Config:              feeds.OCR1Config{Enabled: false},
		OCR2Config:              feeds.OCR2ConfigModel{Enabled: false},
	}
	workflowKey, err := workflowkey.New()
	require.NoError(t, err)

	request := func() *proto.UpdateNodeRequest {
		return &proto.UpdateNodeRequest{
			Version: nodeVersion.Version,
			ChainConfigs: []*proto.ChainConfig{
				{
					Chain: &proto.Chain{
						Id:   cfg.ChainID,
						Type: proto.ChainType_CHAIN_TYPE_EVM,
					},
					AccountAddress:          cfg.AccountAddress,
					AccountAddressPublicKey: &cfg.AccountAddressPublicKey.String,
					AdminAddress:            cfg.AdminAddress,
					FluxMonitorConfig:       &proto.FluxMonitorConfig{Enabled: true},
					Ocr1Config:              &proto.OCR1Config{Enabled: false},
					Ocr2Config:              &proto.OCR2Config{Enabled: false},
				},
			},
			WorkflowKey:   func(s string) *string { return &s }(workflowKey.ID()),
			P2PKeyBundles: []*proto.P2PKeyBundle{},
		}
	}
	successResponse := func() *proto.UpdateNodeResponse {
		return &proto.UpdateNodeResponse{ChainConfigErrors: map[string]*proto.ChainConfigError{}}
	}
	failureResponse := func(chainID string) *proto.UpdateNodeResponse {
		return &proto.UpdateNodeResponse{
			ChainConfigErrors: map[string]*proto.ChainConfigError{chainID: {Message: "error chain " + chainID}},
		}
	}

	tests := []struct {
		name     string
		setup    func(t *testing.T, svc *TestService)
		run      func(svc *TestService) (any, error)
		wantLogs []string
	}{
		{
			name: "create chain",
			setup: func(t *testing.T, svc *TestService) {
				svc.workflowKeystore.On("EnsureKey", mock.Anything).Return(nil)
				svc.workflowKeystore.EXPECT().GetAll().Return([]workflowkey.Key{workflowKey}, nil)
				svc.p2pKeystore.EXPECT().GetAll().Return([]p2pkey.KeyV2{}, nil)
				svc.orm.EXPECT().CreateChainConfig(mock.Anything, cfg).Return(int64(1), nil)
				svc.orm.EXPECT().GetManager(mock.Anything, mgr.ID).Return(&mgr, nil)
				svc.orm.EXPECT().ListChainConfigsByManagerIDs(mock.Anything, []int64{mgr.ID}).Return([]feeds.ChainConfig{cfg}, nil)
				svc.connMgr.EXPECT().GetClient(mgr.ID).Return(svc.fmsClient, nil)
				svc.fmsClient.EXPECT().UpdateNode(mock.Anything, request()).Return(nil, errors.New("error-0")).Once()
				svc.fmsClient.EXPECT().UpdateNode(mock.Anything, request()).Return(failureResponse("1"), nil).Once()
				svc.fmsClient.EXPECT().UpdateNode(mock.Anything, request()).Return(failureResponse("2"), nil).Once()
				svc.fmsClient.EXPECT().UpdateNode(mock.Anything, request()).Return(successResponse(), nil).Once()
			},
			run: func(svc *TestService) (any, error) {
				return svc.CreateChainConfig(testutils.Context(t), cfg)
			},
			wantLogs: []string{
				`failed to sync node info attempt="0" err="SyncNodeInfo.UpdateNode call failed: error-0"`,
				`failed to sync node info attempt="1" err="SyncNodeInfo.UpdateNode call partially failed: error chain 1"`,
				`failed to sync node info attempt="2" err="SyncNodeInfo.UpdateNode call partially failed: error chain 2"`,
				`successfully synced node info`,
			},
		},
		{
			name: "update chain",
			setup: func(t *testing.T, svc *TestService) {
				svc.workflowKeystore.On("EnsureKey", mock.Anything).Return(nil)
				svc.workflowKeystore.EXPECT().GetAll().Return([]workflowkey.Key{workflowKey}, nil)
				svc.p2pKeystore.EXPECT().GetAll().Return([]p2pkey.KeyV2{}, nil)
				svc.orm.EXPECT().UpdateChainConfig(mock.Anything, cfg).Return(int64(1), nil)
				svc.orm.EXPECT().GetChainConfig(mock.Anything, cfg.ID).Return(&cfg, nil)
				svc.orm.EXPECT().ListChainConfigsByManagerIDs(mock.Anything, []int64{mgr.ID}).Return([]feeds.ChainConfig{cfg}, nil)
				svc.connMgr.EXPECT().GetClient(mgr.ID).Return(svc.fmsClient, nil)
				svc.fmsClient.EXPECT().UpdateNode(mock.Anything, request()).Return(failureResponse("3"), nil).Once()
				svc.fmsClient.EXPECT().UpdateNode(mock.Anything, request()).Return(nil, errors.New("error-4")).Once()
				svc.fmsClient.EXPECT().UpdateNode(mock.Anything, request()).Return(failureResponse("5"), nil).Once()
				svc.fmsClient.EXPECT().UpdateNode(mock.Anything, request()).Return(successResponse(), nil).Once()
			},
			run: func(svc *TestService) (any, error) {
				return svc.UpdateChainConfig(testutils.Context(t), cfg)
			},
			wantLogs: []string{
				`failed to sync node info attempt="0" err="SyncNodeInfo.UpdateNode call partially failed: error chain 3"`,
				`failed to sync node info attempt="1" err="SyncNodeInfo.UpdateNode call failed: error-4"`,
				`failed to sync node info attempt="2" err="SyncNodeInfo.UpdateNode call partially failed: error chain 5"`,
				`successfully synced node info`,
			},
		},
		{
			name: "delete chain",
			setup: func(t *testing.T, svc *TestService) {
				svc.workflowKeystore.On("EnsureKey", mock.Anything).Return(nil)
				svc.workflowKeystore.EXPECT().GetAll().Return([]workflowkey.Key{workflowKey}, nil)
				svc.p2pKeystore.EXPECT().GetAll().Return([]p2pkey.KeyV2{}, nil)
				svc.orm.EXPECT().GetChainConfig(mock.Anything, cfg.ID).Return(&cfg, nil)
				svc.orm.EXPECT().DeleteChainConfig(mock.Anything, cfg.ID).Return(cfg.ID, nil)
				svc.orm.EXPECT().GetManager(mock.Anything, mgr.ID).Return(&mgr, nil)
				svc.orm.EXPECT().ListChainConfigsByManagerIDs(mock.Anything, []int64{mgr.ID}).Return([]feeds.ChainConfig{cfg}, nil)
				svc.connMgr.EXPECT().GetClient(mgr.ID).Return(svc.fmsClient, nil)
				svc.fmsClient.EXPECT().UpdateNode(mock.Anything, request()).Return(failureResponse("6"), nil).Once()
				svc.fmsClient.EXPECT().UpdateNode(mock.Anything, request()).Return(failureResponse("7"), nil).Once()
				svc.fmsClient.EXPECT().UpdateNode(mock.Anything, request()).Return(nil, errors.New("error-8")).Once()
				svc.fmsClient.EXPECT().UpdateNode(mock.Anything, request()).Return(successResponse(), nil).Once()
			},
			run: func(svc *TestService) (any, error) {
				return svc.DeleteChainConfig(testutils.Context(t), cfg.ID)
			},
			wantLogs: []string{
				`failed to sync node info attempt="0" err="SyncNodeInfo.UpdateNode call partially failed: error chain 6"`,
				`failed to sync node info attempt="1" err="SyncNodeInfo.UpdateNode call partially failed: error chain 7"`,
				`failed to sync node info attempt="2" err="SyncNodeInfo.UpdateNode call failed: error-8"`,
				`successfully synced node info`,
			},
		},
		{
			name: "more errors than MaxAttempts",
			setup: func(t *testing.T, svc *TestService) {
				svc.workflowKeystore.On("EnsureKey", mock.Anything).Return(nil)
				svc.workflowKeystore.EXPECT().GetAll().Return([]workflowkey.Key{workflowKey}, nil)
				svc.p2pKeystore.EXPECT().GetAll().Return([]p2pkey.KeyV2{}, nil)
				svc.orm.EXPECT().CreateChainConfig(mock.Anything, cfg).Return(int64(1), nil)
				svc.orm.EXPECT().GetManager(mock.Anything, mgr.ID).Return(&mgr, nil)
				svc.orm.EXPECT().ListChainConfigsByManagerIDs(mock.Anything, []int64{mgr.ID}).Return([]feeds.ChainConfig{cfg}, nil)
				svc.connMgr.EXPECT().GetClient(mgr.ID).Return(svc.fmsClient, nil)
				svc.fmsClient.EXPECT().UpdateNode(mock.Anything, request()).Return(failureResponse("9"), nil).Once()
				svc.fmsClient.EXPECT().UpdateNode(mock.Anything, request()).Return(failureResponse("10"), nil).Once()
				svc.fmsClient.EXPECT().UpdateNode(mock.Anything, request()).Return(nil, errors.New("error-11")).Once()
				svc.fmsClient.EXPECT().UpdateNode(mock.Anything, request()).Return(failureResponse("12"), nil).Once()
			},
			run: func(svc *TestService) (any, error) {
				return svc.CreateChainConfig(testutils.Context(t), cfg)
			},
			wantLogs: []string{
				`failed to sync node info attempt="0" err="SyncNodeInfo.UpdateNode call partially failed: error chain 9"`,
				`failed to sync node info attempt="1" err="SyncNodeInfo.UpdateNode call partially failed: error chain 10"`,
				`failed to sync node info attempt="2" err="SyncNodeInfo.UpdateNode call failed: error-11"`,
				`failed to sync node info attempt="3" err="SyncNodeInfo.UpdateNode call partially failed: error chain 12"`,
				`failed to sync node info; aborting err="SyncNodeInfo.UpdateNode call partially failed: error chain 12"`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := setupTestService(t, feeds.WithSyncMinDelay(5*time.Millisecond),
				feeds.WithSyncMaxDelay(50*time.Millisecond), feeds.WithSyncMaxAttempts(4))

			tt.setup(t, svc)
			_, err := tt.run(svc)

			require.NoError(t, err)
			assert.EventuallyWithT(t, func(collect *assert.CollectT) {
				assert.Equal(collect, tt.wantLogs, logMessages(svc.logs.All()))
			}, 1*time.Second, 50*time.Millisecond)
		})
	}
}

func Test_Service_IsJobManaged(t *testing.T) {
	t.Parallel()

	svc := setupTestService(t)
	ctx := testutils.Context(t)
	jobID := int64(1)

	svc.orm.On("IsJobManaged", mock.Anything, jobID).Return(true, nil)

	isManaged, err := svc.IsJobManaged(ctx, jobID)
	require.NoError(t, err)
	assert.True(t, isManaged)
}

func Test_Service_ListJobProposalsByManagersIDs(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	var (
		jp    = feeds.JobProposal{}
		jps   = []feeds.JobProposal{jp}
		fmIDs = []int64{1}
	)
	svc := setupTestService(t)

	svc.orm.On("ListJobProposalsByManagersIDs", mock.Anything, fmIDs).
		Return(jps, nil)

	actual, err := svc.ListJobProposalsByManagersIDs(ctx, fmIDs)
	require.NoError(t, err)

	assert.Equal(t, actual, jps)
}

func Test_Service_GetJobProposal(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	var (
		id = int64(1)
		ms = feeds.JobProposal{ID: id}
	)
	svc := setupTestService(t)

	svc.orm.On("GetJobProposal", mock.Anything, id).
		Return(&ms, nil)

	actual, err := svc.GetJobProposal(ctx, id)
	require.NoError(t, err)

	assert.Equal(t, actual, &ms)
}

func Test_Service_CancelSpec(t *testing.T) {
	var (
		externalJobID = uuid.New()
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
				svc.orm.On("GetSpec", mock.Anything, spec.ID).Return(spec, nil)
				svc.orm.On("GetJobProposal", mock.Anything, jp.ID).Return(jp, nil)

				svc.orm.On("CancelSpec", mock.Anything, spec.ID).Return(nil)
				svc.jobORM.On("FindJobByExternalJobID", mock.Anything, externalJobID).Return(j, nil)
				svc.spawner.On("DeleteJob", mock.Anything, mock.Anything, j.ID).Return(nil)

				svc.fmsClient.On("CancelledJob",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&proto.CancelledJobRequest{
						Uuid:    jp.RemoteUUID.String(),
						Version: int64(spec.Version),
					},
				).Return(&proto.CancelledJobResponse{}, nil)
				svc.orm.On("CountJobProposalsByStatus", mock.Anything).Return(&feeds.JobProposalCounts{}, nil)
				svc.orm.On("WithDataSource", mock.Anything).Return(feeds.ORM(svc.orm))
				svc.jobORM.On("WithDataSource", mock.Anything).Return(job.ORM(svc.jobORM))
			},
			specID: spec.ID,
		},
		{
			name: "success without external job id",
			before: func(svc *TestService) {
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.orm.On("GetSpec", mock.Anything, spec.ID).Return(spec, nil)
				svc.orm.On("GetJobProposal", mock.Anything, jp.ID).Return(&feeds.JobProposal{
					ID:             1,
					RemoteUUID:     externalJobID,
					FeedsManagerID: 100,
				}, nil)

				svc.orm.On("CancelSpec", mock.Anything, spec.ID).Return(nil)
				svc.fmsClient.On("CancelledJob",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&proto.CancelledJobRequest{
						Uuid:    jp.RemoteUUID.String(),
						Version: int64(spec.Version),
					},
				).Return(&proto.CancelledJobResponse{}, nil)
				svc.orm.On("CountJobProposalsByStatus", mock.Anything).Return(&feeds.JobProposalCounts{}, nil)
				svc.orm.On("WithDataSource", mock.Anything).Return(feeds.ORM(svc.orm))
				svc.jobORM.On("WithDataSource", mock.Anything).Return(job.ORM(svc.jobORM))
			},
			specID: spec.ID,
		},
		{
			name: "success without jobs",
			before: func(svc *TestService) {
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.orm.On("GetSpec", mock.Anything, spec.ID).Return(spec, nil)
				svc.orm.On("GetJobProposal", mock.Anything, jp.ID).Return(jp, nil)

				svc.orm.On("CancelSpec", mock.Anything, spec.ID).Return(nil)
				svc.jobORM.On("FindJobByExternalJobID", mock.Anything, externalJobID).Return(job.Job{}, sql.ErrNoRows)
				svc.fmsClient.On("CancelledJob",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&proto.CancelledJobRequest{
						Uuid:    jp.RemoteUUID.String(),
						Version: int64(spec.Version),
					},
				).Return(&proto.CancelledJobResponse{}, nil)
				svc.orm.On("CountJobProposalsByStatus", mock.Anything).Return(&feeds.JobProposalCounts{}, nil)
				svc.orm.On("WithDataSource", mock.Anything).Return(feeds.ORM(svc.orm))
				svc.jobORM.On("WithDataSource", mock.Anything).Return(job.ORM(svc.jobORM))
			},
			specID: spec.ID,
		},
		{
			name: "spec does not exist",
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", mock.Anything, spec.ID).Return(nil, errors.New("Not Found"))
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
				svc.orm.On("GetSpec", mock.Anything, pspec.ID, mock.Anything).Return(pspec, nil)
			},
			specID:  spec.ID,
			wantErr: "must be an approved job proposal spec",
		},
		{
			name: "job proposal does not exist",
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", mock.Anything, spec.ID).Return(spec, nil)
				svc.orm.On("GetJobProposal", mock.Anything, jp.ID).Return(nil, errors.New("Not Found"))
			},
			specID:  spec.ID,
			wantErr: "orm: job proposal: Not Found",
		},
		{
			name: "rpc client not connected",
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", mock.Anything, spec.ID).Return(spec, nil)
				svc.orm.On("GetJobProposal", mock.Anything, jp.ID).Return(jp, nil)
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(nil, errors.New("Not Connected"))
			},
			specID:  spec.ID,
			wantErr: "fms rpc client: Not Connected",
		},
		{
			name: "cancel spec orm fails",
			before: func(svc *TestService) {
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.orm.On("GetSpec", mock.Anything, spec.ID).Return(spec, nil)
				svc.orm.On("GetJobProposal", mock.Anything, jp.ID).Return(jp, nil)
				svc.orm.On("CancelSpec", mock.Anything, spec.ID).Return(errors.New("failure"))
				svc.orm.On("WithDataSource", mock.Anything).Return(feeds.ORM(svc.orm))
				svc.jobORM.On("WithDataSource", mock.Anything).Return(job.ORM(svc.jobORM))
			},
			specID:  spec.ID,
			wantErr: "failure",
		},
		{
			name: "find by external uuid orm fails",
			before: func(svc *TestService) {
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.orm.On("GetSpec", mock.Anything, spec.ID).Return(spec, nil)
				svc.orm.On("GetJobProposal", mock.Anything, jp.ID).Return(jp, nil)

				svc.orm.On("CancelSpec", mock.Anything, spec.ID).Return(nil)
				svc.jobORM.On("FindJobByExternalJobID", mock.Anything, externalJobID).Return(job.Job{}, errors.New("failure"))
				svc.orm.On("WithDataSource", mock.Anything).Return(feeds.ORM(svc.orm))
				svc.jobORM.On("WithDataSource", mock.Anything).Return(job.ORM(svc.jobORM))
			},
			specID:  spec.ID,
			wantErr: "FindJobByExternalJobID failed: failure",
		},
		{
			name: "delete job fails",
			before: func(svc *TestService) {
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.orm.On("GetSpec", mock.Anything, spec.ID).Return(spec, nil)
				svc.orm.On("GetJobProposal", mock.Anything, jp.ID).Return(jp, nil)

				svc.orm.On("CancelSpec", mock.Anything, spec.ID).Return(nil)
				svc.jobORM.On("FindJobByExternalJobID", mock.Anything, externalJobID).Return(j, nil)
				svc.spawner.On("DeleteJob", mock.Anything, mock.Anything, j.ID).Return(errors.New("failure"))
				svc.orm.On("WithDataSource", mock.Anything).Return(feeds.ORM(svc.orm))
				svc.jobORM.On("WithDataSource", mock.Anything).Return(job.ORM(svc.jobORM))
			},
			specID:  spec.ID,
			wantErr: "DeleteJob failed: failure",
		},
		{
			name: "cancelled job rpc call fails",
			before: func(svc *TestService) {
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.orm.On("GetSpec", mock.Anything, spec.ID).Return(spec, nil)
				svc.orm.On("GetJobProposal", mock.Anything, jp.ID).Return(jp, nil)

				svc.orm.On("CancelSpec", mock.Anything, spec.ID).Return(nil)
				svc.jobORM.On("FindJobByExternalJobID", mock.Anything, externalJobID).Return(j, nil)
				svc.spawner.On("DeleteJob", mock.Anything, mock.Anything, j.ID).Return(nil)

				svc.fmsClient.On("CancelledJob",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&proto.CancelledJobRequest{
						Uuid:    jp.RemoteUUID.String(),
						Version: int64(spec.Version),
					},
				).Return(nil, errors.New("failure"))
				svc.orm.On("WithDataSource", mock.Anything).Return(feeds.ORM(svc.orm))
				svc.jobORM.On("WithDataSource", mock.Anything).Return(job.ORM(svc.jobORM))
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
	ctx := testutils.Context(t)

	var (
		id   = int64(1)
		spec = feeds.JobProposalSpec{ID: id}
	)
	svc := setupTestService(t)

	svc.orm.On("GetSpec", mock.Anything, id).
		Return(&spec, nil)

	actual, err := svc.GetSpec(ctx, id)
	require.NoError(t, err)

	assert.Equal(t, &spec, actual)
}

func Test_Service_ListSpecsByJobProposalIDs(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	var (
		id    = int64(1)
		jpID  = int64(200)
		spec  = feeds.JobProposalSpec{ID: id, JobProposalID: jpID}
		specs = []feeds.JobProposalSpec{spec}
	)
	svc := setupTestService(t)

	svc.orm.On("ListSpecsByJobProposalIDs", mock.Anything, []int64{jpID}).
		Return(specs, nil)

	actual, err := svc.ListSpecsByJobProposalIDs(ctx, []int64{jpID})
	require.NoError(t, err)

	assert.Equal(t, specs, actual)
}

func Test_Service_ApproveSpec(t *testing.T) {
	var evmChainID *evmbig.Big
	address := types.EIP55AddressFromAddress(common.Address{})
	externalJobID := uuid.New()

	var (
		ctx  = testutils.Context(t)
		defn = `
name = 'LINK / ETH | version 3 | contract 0x0000000000000000000000000000000000000000'
schemaVersion = 1
contractAddress = '0x0000000000000000000000000000000000000000'
externalJobID = '%s'
type = 'fluxmonitor'
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
			Definition:    fmt.Sprintf(defn, externalJobID),
		}
		spec2 = &feeds.JobProposalSpec{
			ID:            20,
			Status:        feeds.SpecStatusPending,
			JobProposalID: jp.ID,
			Version:       1,
			Definition:    fmt.Sprintf(defn, uuid.Nil),
		}
		rejectedSpec = &feeds.JobProposalSpec{
			ID:            20,
			Status:        feeds.SpecStatusRejected,
			JobProposalID: jp.ID,
			Version:       1,
			Definition:    fmt.Sprintf(defn, externalJobID),
		}
		cancelledSpec = &feeds.JobProposalSpec{
			ID:            20,
			Status:        feeds.SpecStatusCancelled,
			JobProposalID: jp.ID,
			Version:       1,
			Definition:    fmt.Sprintf(defn, externalJobID),
		}
		j = job.Job{
			ID:            1,
			ExternalJobID: externalJobID,
		}
	)

	testCases := []struct {
		name        string
		httpTimeout *commonconfig.Duration
		before      func(svc *TestService)
		id          int64
		force       bool
		wantErr     string
	}{
		{
			name:        "pending job success for new proposals",
			httpTimeout: commonconfig.MustNewDuration(1 * time.Minute),
			before: func(svc *TestService) {
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.orm.EXPECT().GetSpec(mock.Anything, spec.ID).Return(spec, nil)
				svc.orm.EXPECT().GetJobProposal(mock.Anything, jp.ID).Return(jp, nil)
				svc.jobORM.On("AssertBridgesExist", mock.Anything, mock.IsType(pipeline.Pipeline{})).Return(nil)

				svc.jobORM.On("FindJobByExternalJobID", mock.Anything, externalJobID).Return(job.Job{}, sql.ErrNoRows)
				svc.jobORM.On("FindJobIDByAddress", mock.Anything, address, evmChainID, mock.Anything).Return(int32(0), sql.ErrNoRows)

				svc.spawner.
					On("CreateJob",
						mock.Anything,
						mock.Anything,
						mock.MatchedBy(func(j *job.Job) bool {
							return j.Name.String == "LINK / ETH | version 3 | contract 0x0000000000000000000000000000000000000000"
						}),
					).
					Run(func(args mock.Arguments) { (args.Get(2).(*job.Job)).ID = 1 }).
					Return(nil)
				svc.orm.On("ApproveSpec",
					mock.Anything,
					spec.ID,
					mock.IsType(uuid.UUID{}),
				).Return(nil)
				svc.fmsClient.On("ApprovedJob",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&proto.ApprovedJobRequest{
						Uuid:    jp.RemoteUUID.String(),
						Version: int64(spec.Version),
					},
				).Return(&proto.ApprovedJobResponse{}, nil)
				svc.orm.On("CountJobProposalsByStatus", mock.Anything).Return(&feeds.JobProposalCounts{}, nil)
				svc.orm.On("WithDataSource", mock.Anything).Return(feeds.ORM(svc.orm))
				svc.jobORM.On("WithDataSource", mock.Anything).Return(job.ORM(svc.jobORM))
			},
			id:    spec.ID,
			force: false,
		},
		{
			name:        "cancelled spec success when it is the latest spec",
			httpTimeout: commonconfig.MustNewDuration(1 * time.Minute),
			before: func(svc *TestService) {
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.orm.On("GetSpec", mock.Anything, cancelledSpec.ID).Return(cancelledSpec, nil)
				svc.orm.On("GetJobProposal", mock.Anything, jp.ID).Return(jp, nil)
				svc.orm.On("GetLatestSpec", mock.Anything, cancelledSpec.JobProposalID).Return(cancelledSpec, nil)
				svc.jobORM.On("AssertBridgesExist", mock.Anything, mock.IsType(pipeline.Pipeline{})).Return(nil)

				svc.jobORM.On("FindJobByExternalJobID", mock.Anything, externalJobID).Return(job.Job{}, sql.ErrNoRows)
				svc.jobORM.On("FindJobIDByAddress", mock.Anything, address, evmChainID, mock.Anything).Return(int32(0), sql.ErrNoRows)

				svc.spawner.
					On("CreateJob",
						mock.Anything,
						mock.Anything,
						mock.MatchedBy(func(j *job.Job) bool {
							return j.Name.String == "LINK / ETH | version 3 | contract 0x0000000000000000000000000000000000000000"
						}),
					).
					Run(func(args mock.Arguments) { (args.Get(2).(*job.Job)).ID = 1 }).
					Return(nil)
				svc.orm.On("ApproveSpec",
					mock.Anything,
					cancelledSpec.ID,
					externalJobID,
				).Return(nil)
				svc.fmsClient.On("ApprovedJob",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&proto.ApprovedJobRequest{
						Uuid:    jp.RemoteUUID.String(),
						Version: int64(spec.Version),
					},
				).Return(&proto.ApprovedJobResponse{}, nil)
				svc.orm.On("CountJobProposalsByStatus", mock.Anything).Return(&feeds.JobProposalCounts{}, nil)
				svc.orm.On("WithDataSource", mock.Anything).Return(feeds.ORM(svc.orm))
				svc.jobORM.On("WithDataSource", mock.Anything).Return(job.ORM(svc.jobORM))
			},
			id:    cancelledSpec.ID,
			force: false,
		},
		{
			name:        "pending job fail due to spec missing external job id",
			httpTimeout: commonconfig.MustNewDuration(1 * time.Minute),
			before: func(svc *TestService) {
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.orm.EXPECT().GetSpec(mock.Anything, spec.ID).Return(spec2, nil)
				svc.orm.EXPECT().GetJobProposal(mock.Anything, jp.ID).Return(jp, nil)
			},
			id:      spec.ID,
			force:   false,
			wantErr: "failed to approve job spec due to missing ExternalJobID in spec",
		},
		{
			name: "failed due to proposal being revoked",
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", mock.Anything, spec.ID).Return(spec, nil)
				svc.orm.On("GetJobProposal", mock.Anything, jp.ID).Return(&feeds.JobProposal{
					ID:     1,
					Status: feeds.JobProposalStatusRevoked,
				}, nil)
			},
			id:      spec.ID,
			force:   false,
			wantErr: "cannot approve spec for a revoked job proposal",
		},
		{
			name: "failed due to proposal being deleted",
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", mock.Anything, spec.ID).Return(spec, nil)
				svc.orm.On("GetJobProposal", mock.Anything, jp.ID).Return(&feeds.JobProposal{
					ID:     jp.ID,
					Status: feeds.JobProposalStatusDeleted,
				}, nil)
			},
			id:      spec.ID,
			force:   false,
			wantErr: "cannot approve spec for a deleted job proposal",
		},
		{
			name: "failed due to spec already approved",
			before: func(svc *TestService) {
				aspec := &feeds.JobProposalSpec{
					ID:            spec.ID,
					Status:        feeds.SpecStatusApproved,
					JobProposalID: jp.ID,
				}
				svc.orm.On("GetSpec", mock.Anything, aspec.ID, mock.Anything).Return(aspec, nil)
				svc.orm.On("GetJobProposal", mock.Anything, jp.ID).Return(jp, nil)
			},
			id:      spec.ID,
			force:   false,
			wantErr: "cannot approve an approved spec",
		},
		{
			name: "rejected spec fail",
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", mock.Anything, cancelledSpec.ID, mock.Anything).Return(rejectedSpec, nil)
				svc.orm.On("GetJobProposal", mock.Anything, jp.ID).Return(jp, nil)
			},
			id:      rejectedSpec.ID,
			force:   false,
			wantErr: "cannot approve a rejected spec",
		},
		{
			name: "cancelled spec failed not latest spec",
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", mock.Anything, cancelledSpec.ID, mock.Anything).Return(cancelledSpec, nil)
				svc.orm.On("GetJobProposal", mock.Anything, jp.ID).Return(jp, nil)
				svc.orm.On("GetLatestSpec", mock.Anything, cancelledSpec.JobProposalID).Return(&feeds.JobProposalSpec{
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
			name:        "already existing job replacement (found via external job id) error",
			httpTimeout: commonconfig.MustNewDuration(1 * time.Minute),
			before: func(svc *TestService) {
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.orm.On("GetSpec", mock.Anything, spec.ID).Return(spec, nil)
				svc.orm.On("GetJobProposal", mock.Anything, jp.ID).Return(jp, nil)
				svc.jobORM.On("AssertBridgesExist", mock.Anything, mock.IsType(pipeline.Pipeline{})).Return(nil)

				svc.jobORM.On("FindJobByExternalJobID", mock.Anything, externalJobID).Return(j, nil)
				svc.orm.On("WithDataSource", mock.Anything).Return(feeds.ORM(svc.orm))
				svc.jobORM.On("WithDataSource", mock.Anything).Return(job.ORM(svc.jobORM))
			},
			id:      spec.ID,
			force:   false,
			wantErr: "could not approve job proposal: a job for this contract address already exists - please use the 'force' option to replace it",
		},
		{
			name:        "already existing job replacement error",
			httpTimeout: commonconfig.MustNewDuration(1 * time.Minute),
			before: func(svc *TestService) {
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.orm.On("GetSpec", mock.Anything, spec.ID).Return(spec, nil)
				svc.orm.On("GetJobProposal", mock.Anything, jp.ID).Return(jp, nil)
				svc.jobORM.On("AssertBridgesExist", mock.Anything, mock.IsType(pipeline.Pipeline{})).Return(nil)

				svc.jobORM.On("FindJobByExternalJobID", mock.Anything, externalJobID).Return(job.Job{}, sql.ErrNoRows)
				svc.jobORM.On("FindJobIDByAddress", mock.Anything, address, evmChainID, mock.Anything).Return(j.ID, nil)
				svc.orm.On("WithDataSource", mock.Anything).Return(feeds.ORM(svc.orm))
				svc.jobORM.On("WithDataSource", mock.Anything).Return(job.ORM(svc.jobORM))
			},
			id:      spec.ID,
			force:   false,
			wantErr: "could not approve job proposal: a job for this contract address already exists - please use the 'force' option to replace it",
		},
		{
			name:        "already existing self managed job replacement success if forced (via external job id)",
			httpTimeout: commonconfig.MustNewDuration(1 * time.Minute),
			before: func(svc *TestService) {
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.orm.EXPECT().GetSpec(mock.Anything, spec.ID).Return(spec, nil)
				svc.orm.EXPECT().GetJobProposal(mock.Anything, jp.ID).Return(jp, nil)
				svc.jobORM.On("AssertBridgesExist", mock.Anything, mock.IsType(pipeline.Pipeline{})).Return(nil)
				svc.jobORM.On("FindJobByExternalJobID", mock.Anything, externalJobID).Return(j, nil)
				svc.orm.EXPECT().GetApprovedSpec(mock.Anything, jp.ID).Return(nil, sql.ErrNoRows)

				svc.spawner.On("DeleteJob", mock.Anything, mock.Anything, j.ID).Return(nil)

				svc.spawner.
					On("CreateJob",
						mock.Anything,
						mock.Anything,
						mock.MatchedBy(func(j *job.Job) bool {
							return j.Name.String == "LINK / ETH | version 3 | contract 0x0000000000000000000000000000000000000000"
						}),
					).
					Run(func(args mock.Arguments) { (args.Get(2).(*job.Job)).ID = 1 }).
					Return(nil)
				svc.orm.On("ApproveSpec",
					mock.Anything,
					spec.ID,
					externalJobID,
				).Return(nil)
				svc.fmsClient.On("ApprovedJob",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&proto.ApprovedJobRequest{
						Uuid:    jp.RemoteUUID.String(),
						Version: int64(spec.Version),
					},
				).Return(&proto.ApprovedJobResponse{}, nil)
				svc.orm.On("CountJobProposalsByStatus", mock.Anything).Return(&feeds.JobProposalCounts{}, nil)
				svc.orm.On("WithDataSource", mock.Anything).Return(feeds.ORM(svc.orm))
				svc.jobORM.On("WithDataSource", mock.Anything).Return(job.ORM(svc.jobORM))
			},
			id:    spec.ID,
			force: true,
		},
		{
			name:        "already existing self managed job replacement success if forced",
			httpTimeout: commonconfig.MustNewDuration(1 * time.Minute),
			before: func(svc *TestService) {
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.orm.EXPECT().GetSpec(mock.Anything, spec.ID).Return(spec, nil)
				svc.orm.EXPECT().GetJobProposal(mock.Anything, jp.ID).Return(jp, nil)
				svc.jobORM.On("AssertBridgesExist", mock.Anything, mock.IsType(pipeline.Pipeline{})).Return(nil)
				svc.jobORM.On("FindJobByExternalJobID", mock.Anything, externalJobID).Return(job.Job{}, sql.ErrNoRows)
				svc.jobORM.On("FindJobIDByAddress", mock.Anything, address, evmChainID, mock.Anything).Return(j.ID, nil)
				svc.orm.EXPECT().GetApprovedSpec(mock.Anything, jp.ID).Return(nil, sql.ErrNoRows)

				svc.spawner.On("DeleteJob", mock.Anything, mock.Anything, j.ID).Return(nil)

				svc.spawner.
					On("CreateJob",
						mock.Anything,
						mock.Anything,
						mock.MatchedBy(func(j *job.Job) bool {
							return j.Name.String == "LINK / ETH | version 3 | contract 0x0000000000000000000000000000000000000000"
						}),
					).
					Run(func(args mock.Arguments) { (args.Get(2).(*job.Job)).ID = 1 }).
					Return(nil)
				svc.orm.On("ApproveSpec",
					mock.Anything,
					spec.ID,
					externalJobID,
				).Return(nil)
				svc.fmsClient.On("ApprovedJob",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&proto.ApprovedJobRequest{
						Uuid:    jp.RemoteUUID.String(),
						Version: int64(spec.Version),
					},
				).Return(&proto.ApprovedJobResponse{}, nil)
				svc.orm.On("CountJobProposalsByStatus", mock.Anything).Return(&feeds.JobProposalCounts{}, nil)
				svc.orm.On("WithDataSource", mock.Anything).Return(feeds.ORM(svc.orm))
				svc.jobORM.On("WithDataSource", mock.Anything).Return(job.ORM(svc.jobORM))
			},
			id:    spec.ID,
			force: true,
		},
		{
			name:        "already existing FMS managed job replacement success if forced",
			httpTimeout: commonconfig.MustNewDuration(1 * time.Minute),
			before: func(svc *TestService) {
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.orm.EXPECT().GetSpec(mock.Anything, spec.ID).Return(spec, nil)
				svc.orm.EXPECT().GetJobProposal(mock.Anything, jp.ID).Return(jp, nil)
				svc.jobORM.On("AssertBridgesExist", mock.Anything, mock.IsType(pipeline.Pipeline{})).Return(nil)
				svc.jobORM.On("FindJobByExternalJobID", mock.Anything, externalJobID).Return(job.Job{}, sql.ErrNoRows)
				svc.jobORM.On("FindJobIDByAddress", mock.Anything, address, evmChainID, mock.Anything).Return(j.ID, nil)
				svc.orm.EXPECT().GetApprovedSpec(mock.Anything, jp.ID).Return(&feeds.JobProposalSpec{ID: 100}, nil)
				svc.orm.EXPECT().CancelSpec(mock.Anything, int64(100)).Return(nil)

				svc.spawner.On("DeleteJob", mock.Anything, mock.Anything, j.ID).Return(nil)

				svc.spawner.
					On("CreateJob",
						mock.Anything,
						mock.Anything,
						mock.MatchedBy(func(j *job.Job) bool {
							return j.Name.String == "LINK / ETH | version 3 | contract 0x0000000000000000000000000000000000000000"
						}),
					).
					Run(func(args mock.Arguments) { (args.Get(2).(*job.Job)).ID = 1 }).
					Return(nil)
				svc.orm.On("ApproveSpec",
					mock.Anything,
					spec.ID,
					externalJobID,
				).Return(nil)
				svc.fmsClient.On("ApprovedJob",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&proto.ApprovedJobRequest{
						Uuid:    jp.RemoteUUID.String(),
						Version: int64(spec.Version),
					},
				).Return(&proto.ApprovedJobResponse{}, nil)
				svc.orm.On("CountJobProposalsByStatus", mock.Anything).Return(&feeds.JobProposalCounts{}, nil)
				svc.orm.On("WithDataSource", mock.Anything).Return(feeds.ORM(svc.orm))
				svc.jobORM.On("WithDataSource", mock.Anything).Return(job.ORM(svc.jobORM))
			},
			id:    spec.ID,
			force: true,
		},
		{
			name: "spec does not exist",
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", mock.Anything, spec.ID).Return(nil, errors.New("Not Found"))
			},
			id:      spec.ID,
			force:   false,
			wantErr: "orm: job proposal spec: Not Found",
		},
		{
			name: "job proposal does not exist",
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", mock.Anything, spec.ID).Return(spec, nil)
				svc.orm.On("GetJobProposal", mock.Anything, jp.ID).Return(nil, errors.New("Not Found"))
			},
			id:      spec.ID,
			wantErr: "orm: job proposal: Not Found",
		},
		{
			name:        "bridges do not exist",
			httpTimeout: commonconfig.MustNewDuration(1 * time.Minute),
			before: func(svc *TestService) {
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.orm.On("GetSpec", mock.Anything, spec.ID).Return(spec, nil)
				svc.orm.On("GetJobProposal", mock.Anything, jp.ID).Return(jp, nil)
				svc.jobORM.On("AssertBridgesExist", mock.Anything, mock.IsType(pipeline.Pipeline{})).Return(errors.New("bridges do not exist"))
			},
			id:      spec.ID,
			wantErr: "failed to approve job spec due to bridge check: bridges do not exist",
		},
		{
			name: "rpc client not connected",
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", mock.Anything, spec.ID).Return(spec, nil)
				svc.orm.On("GetJobProposal", mock.Anything, jp.ID).Return(jp, nil)
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(nil, errors.New("Not Connected"))
			},
			id:      spec.ID,
			force:   false,
			wantErr: "fms rpc client: Not Connected",
		},
		{
			name:        "Fetching the approved spec fails (via external job id)",
			httpTimeout: commonconfig.MustNewDuration(1 * time.Minute),
			before: func(svc *TestService) {
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.orm.EXPECT().GetSpec(mock.Anything, spec.ID).Return(spec, nil)
				svc.orm.EXPECT().GetJobProposal(mock.Anything, jp.ID).Return(jp, nil)
				svc.jobORM.On("AssertBridgesExist", mock.Anything, mock.IsType(pipeline.Pipeline{})).Return(nil)
				svc.jobORM.On("FindJobByExternalJobID", mock.Anything, externalJobID).Return(j, nil)
				svc.orm.EXPECT().GetApprovedSpec(mock.Anything, jp.ID).Return(nil, errors.New("failure"))
				svc.orm.On("WithDataSource", mock.Anything).Return(feeds.ORM(svc.orm))
				svc.jobORM.On("WithDataSource", mock.Anything).Return(job.ORM(svc.jobORM))
			},
			id:      spec.ID,
			force:   true,
			wantErr: "could not approve job proposal: GetApprovedSpec failed: failure",
		},
		{
			name:        "Fetching the approved spec fails",
			httpTimeout: commonconfig.MustNewDuration(1 * time.Minute),
			before: func(svc *TestService) {
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.orm.EXPECT().GetSpec(mock.Anything, spec.ID).Return(spec, nil)
				svc.orm.EXPECT().GetJobProposal(mock.Anything, jp.ID).Return(jp, nil)
				svc.jobORM.On("AssertBridgesExist", mock.Anything, mock.IsType(pipeline.Pipeline{})).Return(nil)
				svc.jobORM.On("FindJobByExternalJobID", mock.Anything, externalJobID).Return(job.Job{}, sql.ErrNoRows)
				svc.jobORM.On("FindJobIDByAddress", mock.Anything, address, evmChainID, mock.Anything).Return(j.ID, nil)
				svc.orm.EXPECT().GetApprovedSpec(mock.Anything, jp.ID).Return(nil, errors.New("failure"))
				svc.orm.On("WithDataSource", mock.Anything).Return(feeds.ORM(svc.orm))
				svc.jobORM.On("WithDataSource", mock.Anything).Return(job.ORM(svc.jobORM))
			},
			id:      spec.ID,
			force:   true,
			wantErr: "could not approve job proposal: GetApprovedSpec failed: failure",
		},
		{
			name:        "spec cancellation fails (via external job id)",
			httpTimeout: commonconfig.MustNewDuration(1 * time.Minute),
			before: func(svc *TestService) {
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.orm.EXPECT().GetSpec(mock.Anything, spec.ID).Return(spec, nil)
				svc.orm.EXPECT().GetJobProposal(mock.Anything, jp.ID).Return(jp, nil)
				svc.jobORM.On("AssertBridgesExist", mock.Anything, mock.IsType(pipeline.Pipeline{})).Return(nil)
				svc.jobORM.On("FindJobByExternalJobID", mock.Anything, externalJobID).Return(j, nil)
				svc.orm.EXPECT().GetApprovedSpec(mock.Anything, jp.ID).Return(&feeds.JobProposalSpec{ID: 100}, nil)
				svc.orm.EXPECT().CancelSpec(mock.Anything, int64(100)).Return(errors.New("failure"))
				svc.orm.On("WithDataSource", mock.Anything).Return(feeds.ORM(svc.orm))
				svc.jobORM.On("WithDataSource", mock.Anything).Return(job.ORM(svc.jobORM))
			},
			id:      spec.ID,
			force:   true,
			wantErr: "could not approve job proposal: failure",
		},
		{
			name:        "spec cancellation fails",
			httpTimeout: commonconfig.MustNewDuration(1 * time.Minute),
			before: func(svc *TestService) {
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.orm.EXPECT().GetSpec(mock.Anything, spec.ID).Return(spec, nil)
				svc.orm.EXPECT().GetJobProposal(mock.Anything, jp.ID).Return(jp, nil)
				svc.jobORM.On("AssertBridgesExist", mock.Anything, mock.IsType(pipeline.Pipeline{})).Return(nil)
				svc.jobORM.On("FindJobByExternalJobID", mock.Anything, externalJobID).Return(job.Job{}, sql.ErrNoRows)
				svc.jobORM.On("FindJobIDByAddress", mock.Anything, address, evmChainID, mock.Anything).Return(j.ID, nil)
				svc.orm.EXPECT().GetApprovedSpec(mock.Anything, jp.ID).Return(&feeds.JobProposalSpec{ID: 100}, nil)
				svc.orm.EXPECT().CancelSpec(mock.Anything, int64(100)).Return(errors.New("failure"))
				svc.orm.On("WithDataSource", mock.Anything).Return(feeds.ORM(svc.orm))
				svc.jobORM.On("WithDataSource", mock.Anything).Return(job.ORM(svc.jobORM))
			},
			id:      spec.ID,
			force:   true,
			wantErr: "could not approve job proposal: failure",
		},
		{
			name:        "create job error",
			httpTimeout: commonconfig.MustNewDuration(1 * time.Minute),
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", mock.Anything, spec.ID).Return(spec, nil)
				svc.orm.On("GetJobProposal", mock.Anything, jp.ID).Return(jp, nil)
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.jobORM.On("AssertBridgesExist", mock.Anything, mock.IsType(pipeline.Pipeline{})).Return(nil)

				svc.jobORM.On("FindJobByExternalJobID", mock.Anything, externalJobID).Return(job.Job{}, sql.ErrNoRows)
				svc.jobORM.On("FindJobIDByAddress", mock.Anything, address, evmChainID, mock.Anything).Return(int32(0), sql.ErrNoRows)

				svc.spawner.
					On("CreateJob",
						mock.Anything,
						mock.Anything,
						mock.MatchedBy(func(j *job.Job) bool {
							return j.Name.String == "LINK / ETH | version 3 | contract 0x0000000000000000000000000000000000000000"
						}),
					).
					Return(errors.New("could not save"))
				svc.orm.On("WithDataSource", mock.Anything).Return(feeds.ORM(svc.orm))
				svc.jobORM.On("WithDataSource", mock.Anything).Return(job.ORM(svc.jobORM))
			},
			id:      spec.ID,
			force:   false,
			wantErr: "could not approve job proposal: could not save",
		},
		{
			name:        "approve spec orm error",
			httpTimeout: commonconfig.MustNewDuration(1 * time.Minute),
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", mock.Anything, spec.ID).Return(spec, nil)
				svc.orm.On("GetJobProposal", mock.Anything, jp.ID).Return(jp, nil)
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.jobORM.On("AssertBridgesExist", mock.Anything, mock.IsType(pipeline.Pipeline{})).Return(nil)

				svc.jobORM.On("FindJobByExternalJobID", mock.Anything, externalJobID).Return(job.Job{}, sql.ErrNoRows)
				svc.jobORM.On("FindJobIDByAddress", mock.Anything, address, evmChainID, mock.Anything).Return(int32(0), sql.ErrNoRows)

				svc.spawner.
					On("CreateJob",
						mock.Anything,
						mock.Anything,
						mock.MatchedBy(func(j *job.Job) bool {
							return j.Name.String == "LINK / ETH | version 3 | contract 0x0000000000000000000000000000000000000000"
						}),
					).
					Run(func(args mock.Arguments) { (args.Get(2).(*job.Job)).ID = 1 }).
					Return(nil)
				svc.orm.On("ApproveSpec",
					mock.Anything,
					spec.ID,
					externalJobID,
				).Return(errors.New("failure"))
				svc.orm.On("WithDataSource", mock.Anything).Return(feeds.ORM(svc.orm))
				svc.jobORM.On("WithDataSource", mock.Anything).Return(job.ORM(svc.jobORM))
			},
			id:      spec.ID,
			force:   false,
			wantErr: "could not approve job proposal: failure",
		},
		{
			name:        "fms call error",
			httpTimeout: commonconfig.MustNewDuration(1 * time.Minute),
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", mock.Anything, spec.ID).Return(spec, nil)
				svc.orm.On("GetJobProposal", mock.Anything, jp.ID).Return(jp, nil)
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.jobORM.On("AssertBridgesExist", mock.Anything, mock.IsType(pipeline.Pipeline{})).Return(nil)

				svc.jobORM.On("FindJobByExternalJobID", mock.Anything, externalJobID).Return(job.Job{}, sql.ErrNoRows)
				svc.jobORM.On("FindJobIDByAddress", mock.Anything, address, evmChainID, mock.Anything).Return(int32(0), sql.ErrNoRows)

				svc.spawner.
					On("CreateJob",
						mock.Anything,
						mock.Anything,
						mock.MatchedBy(func(j *job.Job) bool {
							return j.Name.String == "LINK / ETH | version 3 | contract 0x0000000000000000000000000000000000000000"
						}),
					).
					Run(func(args mock.Arguments) { (args.Get(2).(*job.Job)).ID = 1 }).
					Return(nil)
				svc.orm.On("ApproveSpec",
					mock.Anything,
					spec.ID,
					externalJobID,
				).Return(nil)
				svc.fmsClient.On("ApprovedJob",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&proto.ApprovedJobRequest{
						Uuid:    jp.RemoteUUID.String(),
						Version: int64(spec.Version),
					},
				).Return(nil, errors.New("failure"))
				svc.orm.On("WithDataSource", mock.Anything).Return(feeds.ORM(svc.orm))
				svc.jobORM.On("WithDataSource", mock.Anything).Return(job.ORM(svc.jobORM))
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
				c.OCR2.Enabled = testutils.Ptr(true)
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
	address := "0x613a38AC1659769640aaE063C651F48E0250454C"
	feedIDHex := "0x0000000000000000000000000000000000000000000000000000000000000001"
	feedID := common.HexToHash(feedIDHex)
	externalJobID := uuid.New()

	var (
		ctx  = testutils.Context(t)
		defn = `
name = 'LINK / ETH | version 3 | contract 0x0000000000000000000000000000000000000000'
type               = "offchainreporting2"
pluginType         = "median"
schemaVersion      = 1
relay              = "evm"
contractID         = "0x613a38AC1659769640aaE063C651F48E0250454C"
externalJobID      = '%s'
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
gasPriceSubunitsSource = """
ds1          [type=bridge name=voter_turnout];
ds1_parse    [type=jsonparse path="one,two"];
ds1_multiply [type=multiply times=1.23];
ds1 -> ds1_parse -> ds1_multiply -> answer1;
answer1      [type=median index=0];
"""
[pluginConfig.juelsPerFeeCoinCache]
updateInterval = "30s"
`
		defn2 = `
name = 'LINK / ETH | version 3 | contract 0x0000000000000000000000000000000000000000'
type               = "offchainreporting2"
pluginType         = "median"
schemaVersion      = 1
relay              = "evm"
contractID         = "0x613a38AC1659769640aaE063C651F48E0250454C"
externalJobID      = '%s'
feedID             = '%s'
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
# intentionally do not set gasPriceSubunitsSource for this pipeline example to cover case when none is set
[pluginConfig.juelsPerFeeCoinCache]
updateInterval = "20m"
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
			Definition:    fmt.Sprintf(defn, externalJobID.String()),
		}
		rejectedSpec = &feeds.JobProposalSpec{
			ID:            20,
			Status:        feeds.SpecStatusRejected,
			JobProposalID: jp.ID,
			Version:       1,
			Definition:    fmt.Sprintf(defn, externalJobID.String()),
		}
		cancelledSpec = &feeds.JobProposalSpec{
			ID:            20,
			Status:        feeds.SpecStatusCancelled,
			JobProposalID: jp.ID,
			Version:       1,
			Definition:    fmt.Sprintf(defn, externalJobID.String()),
		}
		j = job.Job{
			ID:            1,
			ExternalJobID: externalJobID,
		}
	)

	testCases := []struct {
		name        string
		httpTimeout *commonconfig.Duration
		before      func(svc *TestService)
		id          int64
		force       bool
		wantErr     string
	}{
		{
			name:        "pending job success",
			httpTimeout: commonconfig.MustNewDuration(1 * time.Minute),
			before: func(svc *TestService) {
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.orm.On("GetSpec", mock.Anything, spec.ID).Return(spec, nil)
				svc.orm.On("GetJobProposal", mock.Anything, jp.ID).Return(jp, nil)
				svc.jobORM.On("AssertBridgesExist", mock.Anything, mock.IsType(pipeline.Pipeline{})).Return(nil)
				svc.jobORM.On("FindJobByExternalJobID", mock.Anything, externalJobID).Return(job.Job{}, sql.ErrNoRows)
				svc.jobORM.On("FindOCR2JobIDByAddress", mock.Anything, address, (*common.Hash)(nil)).Return(int32(0), sql.ErrNoRows)

				svc.spawner.
					On("CreateJob",
						mock.Anything,
						mock.Anything,
						mock.MatchedBy(func(j *job.Job) bool {
							return j.Name.String == "LINK / ETH | version 3 | contract 0x0000000000000000000000000000000000000000"
						}),
					).
					Run(func(args mock.Arguments) { (args.Get(2).(*job.Job)).ID = 1 }).
					Return(nil)
				svc.orm.On("ApproveSpec",
					mock.Anything,
					spec.ID,
					externalJobID,
				).Return(nil)
				svc.fmsClient.On("ApprovedJob",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&proto.ApprovedJobRequest{
						Uuid:    jp.RemoteUUID.String(),
						Version: int64(spec.Version),
					},
				).Return(&proto.ApprovedJobResponse{}, nil)
				svc.orm.On("CountJobProposalsByStatus", mock.Anything).Return(&feeds.JobProposalCounts{}, nil)
				svc.orm.On("WithDataSource", mock.Anything).Return(feeds.ORM(svc.orm))
				svc.jobORM.On("WithDataSource", mock.Anything).Return(job.ORM(svc.jobORM))
			},
			id:    spec.ID,
			force: false,
		},
		{
			name:        "cancelled spec success when it is the latest spec",
			httpTimeout: commonconfig.MustNewDuration(1 * time.Minute),
			before: func(svc *TestService) {
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.orm.On("GetSpec", mock.Anything, cancelledSpec.ID, mock.Anything).Return(cancelledSpec, nil)
				svc.orm.On("GetJobProposal", mock.Anything, jp.ID).Return(jp, nil)
				svc.orm.On("GetLatestSpec", mock.Anything, cancelledSpec.JobProposalID).Return(cancelledSpec, nil)
				svc.jobORM.On("AssertBridgesExist", mock.Anything, mock.IsType(pipeline.Pipeline{})).Return(nil)

				svc.jobORM.On("FindJobByExternalJobID", mock.Anything, externalJobID).Return(job.Job{}, sql.ErrNoRows)
				svc.jobORM.On("FindOCR2JobIDByAddress", mock.Anything, address, (*common.Hash)(nil)).Return(int32(0), sql.ErrNoRows)

				svc.spawner.
					On("CreateJob",
						mock.Anything,
						mock.Anything,
						mock.MatchedBy(func(j *job.Job) bool {
							return j.Name.String == "LINK / ETH | version 3 | contract 0x0000000000000000000000000000000000000000"
						}),
					).
					Run(func(args mock.Arguments) { (args.Get(2).(*job.Job)).ID = 1 }).
					Return(nil)
				svc.orm.On("ApproveSpec",
					mock.Anything,
					cancelledSpec.ID,
					externalJobID,
				).Return(nil)
				svc.fmsClient.On("ApprovedJob",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&proto.ApprovedJobRequest{
						Uuid:    jp.RemoteUUID.String(),
						Version: int64(spec.Version),
					},
				).Return(&proto.ApprovedJobResponse{}, nil)
				svc.orm.On("CountJobProposalsByStatus", mock.Anything).Return(&feeds.JobProposalCounts{}, nil)
				svc.orm.On("WithDataSource", mock.Anything).Return(feeds.ORM(svc.orm))
				svc.jobORM.On("WithDataSource", mock.Anything).Return(job.ORM(svc.jobORM))
			},
			id:    cancelledSpec.ID,
			force: false,
		},
		{
			name: "cancelled spec failed not latest spec",
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", mock.Anything, cancelledSpec.ID, mock.Anything).Return(cancelledSpec, nil)
				svc.orm.On("GetJobProposal", mock.Anything, jp.ID).Return(jp, nil)
				svc.orm.On("GetLatestSpec", mock.Anything, cancelledSpec.JobProposalID).Return(&feeds.JobProposalSpec{
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
				svc.orm.On("GetSpec", mock.Anything, cancelledSpec.ID, mock.Anything).Return(rejectedSpec, nil)
				svc.orm.On("GetJobProposal", mock.Anything, jp.ID).Return(jp, nil)
			},
			id:      rejectedSpec.ID,
			force:   false,
			wantErr: "cannot approve a rejected spec",
		},
		{
			name:        "already existing job replacement error",
			httpTimeout: commonconfig.MustNewDuration(1 * time.Minute),
			before: func(svc *TestService) {
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.orm.On("GetSpec", mock.Anything, spec.ID).Return(spec, nil)
				svc.orm.On("GetJobProposal", mock.Anything, jp.ID).Return(jp, nil)
				svc.jobORM.On("AssertBridgesExist", mock.Anything, mock.IsType(pipeline.Pipeline{})).Return(nil)
				svc.jobORM.On("FindJobByExternalJobID", mock.Anything, externalJobID).Return(job.Job{}, sql.ErrNoRows)
				svc.jobORM.On("FindOCR2JobIDByAddress", mock.Anything, address, (*common.Hash)(nil)).Return(j.ID, nil)
				svc.orm.On("WithDataSource", mock.Anything).Return(feeds.ORM(svc.orm))
				svc.jobORM.On("WithDataSource", mock.Anything).Return(job.ORM(svc.jobORM))
			},
			id:      spec.ID,
			force:   false,
			wantErr: "could not approve job proposal: a job for this contract address already exists - please use the 'force' option to replace it",
		},
		{
			name:        "already existing self managed job replacement success if forced without feedID",
			httpTimeout: commonconfig.MustNewDuration(1 * time.Minute),
			before: func(svc *TestService) {
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.orm.On("GetSpec", mock.Anything, spec.ID).Return(spec, nil)
				svc.orm.On("GetJobProposal", mock.Anything, jp.ID).Return(jp, nil)
				svc.jobORM.On("AssertBridgesExist", mock.Anything, mock.IsType(pipeline.Pipeline{})).Return(nil)
				svc.orm.EXPECT().GetApprovedSpec(mock.Anything, jp.ID).Return(nil, sql.ErrNoRows)
				svc.jobORM.On("FindJobByExternalJobID", mock.Anything, externalJobID).Return(job.Job{}, sql.ErrNoRows)
				svc.jobORM.On("FindOCR2JobIDByAddress", mock.Anything, address, (*common.Hash)(nil)).Return(j.ID, nil)
				svc.spawner.On("DeleteJob", mock.Anything, mock.Anything, j.ID).Return(nil)

				svc.spawner.
					On("CreateJob",
						mock.Anything,
						mock.Anything,
						mock.MatchedBy(func(j *job.Job) bool {
							return j.Name.String == "LINK / ETH | version 3 | contract 0x0000000000000000000000000000000000000000"
						}),
					).
					Run(func(args mock.Arguments) { (args.Get(2).(*job.Job)).ID = 1 }).
					Return(nil)
				svc.orm.On("ApproveSpec",
					mock.Anything,
					spec.ID,
					externalJobID,
				).Return(nil)
				svc.fmsClient.On("ApprovedJob",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&proto.ApprovedJobRequest{
						Uuid:    jp.RemoteUUID.String(),
						Version: int64(spec.Version),
					},
				).Return(&proto.ApprovedJobResponse{}, nil)
				svc.orm.On("CountJobProposalsByStatus", mock.Anything).Return(&feeds.JobProposalCounts{}, nil)
				svc.orm.On("WithDataSource", mock.Anything).Return(feeds.ORM(svc.orm))
				svc.jobORM.On("WithDataSource", mock.Anything).Return(job.ORM(svc.jobORM))
			},
			id:    spec.ID,
			force: true,
		},
		{
			name:        "already existing self managed job replacement success if forced with feedID",
			httpTimeout: commonconfig.MustNewDuration(1 * time.Minute),
			before: func(svc *TestService) {
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.orm.On("GetSpec", mock.Anything, spec.ID).Return(&feeds.JobProposalSpec{
					ID:            20,
					Status:        feeds.SpecStatusPending,
					JobProposalID: jp.ID,
					Version:       1,
					Definition:    fmt.Sprintf(defn2, externalJobID.String(), &feedID),
				}, nil)
				svc.orm.On("GetJobProposal", mock.Anything, jp.ID).Return(jp, nil)
				svc.jobORM.On("AssertBridgesExist", mock.Anything, mock.IsType(pipeline.Pipeline{})).Return(nil)
				svc.orm.EXPECT().GetApprovedSpec(mock.Anything, jp.ID).Return(nil, sql.ErrNoRows)
				svc.jobORM.On("FindJobByExternalJobID", mock.Anything, externalJobID).Return(job.Job{}, sql.ErrNoRows)
				svc.jobORM.On("FindOCR2JobIDByAddress", mock.Anything, address, &feedID).Return(j.ID, nil)
				svc.spawner.On("DeleteJob", mock.Anything, mock.Anything, j.ID).Return(nil)

				svc.spawner.
					On("CreateJob",
						mock.Anything,
						mock.Anything,
						mock.MatchedBy(func(j *job.Job) bool {
							return j.Name.String == "LINK / ETH | version 3 | contract 0x0000000000000000000000000000000000000000"
						}),
					).
					Run(func(args mock.Arguments) { (args.Get(2).(*job.Job)).ID = 1 }).
					Return(nil)
				svc.orm.On("ApproveSpec",
					mock.Anything,
					spec.ID,
					externalJobID,
				).Return(nil)
				svc.fmsClient.On("ApprovedJob",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&proto.ApprovedJobRequest{
						Uuid:    jp.RemoteUUID.String(),
						Version: int64(spec.Version),
					},
				).Return(&proto.ApprovedJobResponse{}, nil)
				svc.orm.On("CountJobProposalsByStatus", mock.Anything).Return(&feeds.JobProposalCounts{}, nil)
				svc.orm.On("WithDataSource", mock.Anything).Return(feeds.ORM(svc.orm))
				svc.jobORM.On("WithDataSource", mock.Anything).Return(job.ORM(svc.jobORM))
			},
			id:    spec.ID,
			force: true,
		},
		{
			name:        "already existing FMS managed job replacement success if forced",
			httpTimeout: commonconfig.MustNewDuration(1 * time.Minute),
			before: func(svc *TestService) {
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.orm.On("GetSpec", mock.Anything, spec.ID).Return(spec, nil)
				svc.orm.On("GetJobProposal", mock.Anything, jp.ID).Return(jp, nil)
				svc.jobORM.On("AssertBridgesExist", mock.Anything, mock.IsType(pipeline.Pipeline{})).Return(nil)
				svc.orm.EXPECT().GetApprovedSpec(mock.Anything, jp.ID).Return(&feeds.JobProposalSpec{ID: 100}, nil)
				svc.orm.EXPECT().CancelSpec(mock.Anything, int64(100)).Return(nil)
				svc.jobORM.On("FindJobByExternalJobID", mock.Anything, externalJobID).Return(job.Job{}, sql.ErrNoRows)
				svc.jobORM.On("FindOCR2JobIDByAddress", mock.Anything, address, (*common.Hash)(nil)).Return(j.ID, nil)
				svc.spawner.On("DeleteJob", mock.Anything, mock.Anything, j.ID).Return(nil)

				svc.spawner.
					On("CreateJob",
						mock.Anything,
						mock.Anything,
						mock.MatchedBy(func(j *job.Job) bool {
							return j.Name.String == "LINK / ETH | version 3 | contract 0x0000000000000000000000000000000000000000"
						}),
					).
					Run(func(args mock.Arguments) { (args.Get(2).(*job.Job)).ID = 1 }).
					Return(nil)
				svc.orm.On("ApproveSpec",
					mock.Anything,
					spec.ID,
					externalJobID,
				).Return(nil)
				svc.fmsClient.On("ApprovedJob",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&proto.ApprovedJobRequest{
						Uuid:    jp.RemoteUUID.String(),
						Version: int64(spec.Version),
					},
				).Return(&proto.ApprovedJobResponse{}, nil)
				svc.orm.On("CountJobProposalsByStatus", mock.Anything).Return(&feeds.JobProposalCounts{}, nil)
				svc.orm.On("WithDataSource", mock.Anything).Return(feeds.ORM(svc.orm))
				svc.jobORM.On("WithDataSource", mock.Anything).Return(job.ORM(svc.jobORM))
			},
			id:    spec.ID,
			force: true,
		},
		{
			name: "spec does not exist",
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", mock.Anything, spec.ID).Return(nil, errors.New("Not Found"))
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
				svc.orm.On("GetSpec", mock.Anything, spec.ID).Return(aspec, nil)
				svc.orm.On("GetJobProposal", mock.Anything, jp.ID).Return(jp, nil)
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
				svc.orm.On("GetSpec", mock.Anything, rspec.ID, mock.Anything).Return(rspec, nil)
				svc.orm.On("GetJobProposal", mock.Anything, jp.ID).Return(jp, nil)
			},
			id:      spec.ID,
			force:   false,
			wantErr: "cannot approve a rejected spec",
		},
		{
			name: "job proposal does not exist",
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", mock.Anything, spec.ID).Return(spec, nil)
				svc.orm.On("GetJobProposal", mock.Anything, jp.ID).Return(nil, errors.New("Not Found"))
			},
			id:      spec.ID,
			wantErr: "orm: job proposal: Not Found",
		},
		{
			name:        "bridges do not exist",
			httpTimeout: commonconfig.MustNewDuration(1 * time.Minute),
			before: func(svc *TestService) {
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.orm.On("GetSpec", mock.Anything, spec.ID).Return(spec, nil)
				svc.orm.On("GetJobProposal", mock.Anything, jp.ID).Return(jp, nil)
				svc.jobORM.On("AssertBridgesExist", mock.Anything, mock.IsType(pipeline.Pipeline{})).Return(errors.New("bridges do not exist"))
			},
			id:      spec.ID,
			wantErr: "failed to approve job spec due to bridge check: bridges do not exist",
		},
		{
			name: "rpc client not connected",
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", mock.Anything, spec.ID).Return(spec, nil)
				svc.orm.On("GetJobProposal", mock.Anything, jp.ID).Return(jp, nil)
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(nil, errors.New("Not Connected"))
			},
			id:      spec.ID,
			force:   false,
			wantErr: "fms rpc client: Not Connected",
		},
		{
			name:        "create job error",
			httpTimeout: commonconfig.MustNewDuration(1 * time.Minute),
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", mock.Anything, spec.ID).Return(spec, nil)
				svc.orm.On("GetJobProposal", mock.Anything, jp.ID).Return(jp, nil)
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.jobORM.On("AssertBridgesExist", mock.Anything, mock.IsType(pipeline.Pipeline{})).Return(nil)

				svc.jobORM.On("FindJobByExternalJobID", mock.Anything, externalJobID).Return(job.Job{}, sql.ErrNoRows)
				svc.jobORM.On("FindOCR2JobIDByAddress", mock.Anything, address, (*common.Hash)(nil)).Return(int32(0), sql.ErrNoRows)

				svc.spawner.
					On("CreateJob",
						mock.Anything,
						mock.Anything,
						mock.MatchedBy(func(j *job.Job) bool {
							return j.Name.String == "LINK / ETH | version 3 | contract 0x0000000000000000000000000000000000000000"
						}),
					).
					Return(errors.New("could not save"))
				svc.orm.On("WithDataSource", mock.Anything).Return(feeds.ORM(svc.orm))
				svc.jobORM.On("WithDataSource", mock.Anything).Return(job.ORM(svc.jobORM))
			},
			id:      spec.ID,
			force:   false,
			wantErr: "could not approve job proposal: could not save",
		},
		{
			name:        "approve spec orm error",
			httpTimeout: commonconfig.MustNewDuration(1 * time.Minute),
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", mock.Anything, spec.ID).Return(spec, nil)
				svc.orm.On("GetJobProposal", mock.Anything, jp.ID).Return(jp, nil)
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.jobORM.On("AssertBridgesExist", mock.Anything, mock.IsType(pipeline.Pipeline{})).Return(nil)

				svc.jobORM.On("FindJobByExternalJobID", mock.Anything, externalJobID).Return(job.Job{}, sql.ErrNoRows)
				svc.jobORM.On("FindOCR2JobIDByAddress", mock.Anything, address, (*common.Hash)(nil)).Return(int32(0), sql.ErrNoRows)

				svc.spawner.
					On("CreateJob",
						mock.Anything,
						mock.Anything,
						mock.MatchedBy(func(j *job.Job) bool {
							return j.Name.String == "LINK / ETH | version 3 | contract 0x0000000000000000000000000000000000000000"
						}),
					).
					Run(func(args mock.Arguments) { (args.Get(2).(*job.Job)).ID = 1 }).
					Return(nil)
				svc.orm.On("ApproveSpec",
					mock.Anything,
					spec.ID,
					externalJobID,
				).Return(errors.New("failure"))
				svc.orm.On("WithDataSource", mock.Anything).Return(feeds.ORM(svc.orm))
				svc.jobORM.On("WithDataSource", mock.Anything).Return(job.ORM(svc.jobORM))
			},
			id:      spec.ID,
			force:   false,
			wantErr: "could not approve job proposal: failure",
		},
		{
			name:        "fms call error",
			httpTimeout: commonconfig.MustNewDuration(1 * time.Minute),
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", mock.Anything, spec.ID).Return(spec, nil)
				svc.orm.On("GetJobProposal", mock.Anything, jp.ID).Return(jp, nil)
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.jobORM.On("AssertBridgesExist", mock.Anything, mock.IsType(pipeline.Pipeline{})).Return(nil)

				svc.jobORM.On("FindJobByExternalJobID", mock.Anything, externalJobID).Return(job.Job{}, sql.ErrNoRows)
				svc.jobORM.On("FindOCR2JobIDByAddress", mock.Anything, address, (*common.Hash)(nil)).Return(int32(0), sql.ErrNoRows)

				svc.spawner.
					On("CreateJob",
						mock.Anything,
						mock.Anything,
						mock.MatchedBy(func(j *job.Job) bool {
							return j.Name.String == "LINK / ETH | version 3 | contract 0x0000000000000000000000000000000000000000"
						}),
					).
					Run(func(args mock.Arguments) { (args.Get(2).(*job.Job)).ID = 1 }).
					Return(nil)
				svc.orm.On("ApproveSpec",
					mock.Anything,
					spec.ID,
					externalJobID,
				).Return(nil)
				svc.fmsClient.On("ApprovedJob",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&proto.ApprovedJobRequest{
						Uuid:    jp.RemoteUUID.String(),
						Version: int64(spec.Version),
					},
				).Return(nil, errors.New("failure"))
				svc.orm.On("WithDataSource", mock.Anything).Return(feeds.ORM(svc.orm))
				svc.jobORM.On("WithDataSource", mock.Anything).Return(job.ORM(svc.jobORM))
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
				c.OCR2.Enabled = testutils.Ptr(true)
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

func Test_Service_ApproveSpec_Stream(t *testing.T) {
	externalJobID := uuid.New()
	streamName := "LINK / ETH | version 3 | contract 0x0000000000000000000000000000000000000000"
	streamID := uint32(1009001032)

	var (
		ctx = testutils.Context(t)

		jp = &feeds.JobProposal{
			ID:             1,
			FeedsManagerID: 100,
		}
		spec = &feeds.JobProposalSpec{
			ID:            20,
			Status:        feeds.SpecStatusPending,
			JobProposalID: jp.ID,
			Version:       1,
			Definition:    fmt.Sprintf(StreamTestSpecTemplate, streamName, externalJobID.String(), streamID),
		}
		rejectedSpec = &feeds.JobProposalSpec{
			ID:            20,
			Status:        feeds.SpecStatusRejected,
			JobProposalID: jp.ID,
			Version:       1,
			Definition:    fmt.Sprintf(StreamTestSpecTemplate, streamName, externalJobID.String(), streamID),
		}
		cancelledSpec = &feeds.JobProposalSpec{
			ID:            20,
			Status:        feeds.SpecStatusCancelled,
			JobProposalID: jp.ID,
			Version:       1,
			Definition:    fmt.Sprintf(StreamTestSpecTemplate, streamName, externalJobID.String(), streamID),
		}
		j = job.Job{
			ID:            1,
			ExternalJobID: externalJobID,
		}
	)

	testCases := []struct {
		name        string
		httpTimeout *commonconfig.Duration
		before      func(svc *TestService)
		id          int64
		force       bool
		wantErr     string
	}{
		{
			name:        "pending job success",
			httpTimeout: commonconfig.MustNewDuration(1 * time.Minute),
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", mock.Anything, spec.ID).Return(spec, nil)
				svc.orm.On("GetJobProposal", mock.Anything, jp.ID).Return(jp, nil)
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.jobORM.On("AssertBridgesExist", mock.Anything, mock.IsType(pipeline.Pipeline{})).Return(nil)
				svc.jobORM.On("FindJobByExternalJobID", mock.Anything, externalJobID).Return(job.Job{}, sql.ErrNoRows)
				svc.jobORM.On("FindJobIDByStreamID", mock.Anything, mock.Anything).Return(int32(0), sql.ErrNoRows)

				svc.spawner.
					On("CreateJob",
						mock.Anything,
						mock.Anything,
						mock.MatchedBy(func(j *job.Job) bool {
							return j.Name.String == streamName
						}),
					).
					Run(func(args mock.Arguments) { (args.Get(2).(*job.Job)).ID = 1 }).
					Return(nil)
				svc.orm.On("ApproveSpec",
					mock.Anything,
					spec.ID,
					externalJobID,
				).Return(nil)
				svc.fmsClient.On("ApprovedJob",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&proto.ApprovedJobRequest{
						Uuid:    jp.RemoteUUID.String(),
						Version: int64(spec.Version),
					},
				).Return(&proto.ApprovedJobResponse{}, nil)
				svc.orm.On("CountJobProposalsByStatus", mock.Anything).Return(&feeds.JobProposalCounts{}, nil)
				svc.orm.On("WithDataSource", mock.Anything).Return(feeds.ORM(svc.orm))
				svc.jobORM.On("WithDataSource", mock.Anything).Return(job.ORM(svc.jobORM))
			},
			id:    spec.ID,
			force: false,
		},
		{
			name:        "cancelled spec success when it is the latest spec",
			httpTimeout: commonconfig.MustNewDuration(1 * time.Minute),
			before: func(svc *TestService) {
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.orm.On("GetSpec", mock.Anything, cancelledSpec.ID, mock.Anything).Return(cancelledSpec, nil)
				svc.orm.On("GetJobProposal", mock.Anything, jp.ID).Return(jp, nil)
				svc.orm.On("GetLatestSpec", mock.Anything, cancelledSpec.JobProposalID).Return(cancelledSpec, nil)
				svc.jobORM.On("AssertBridgesExist", mock.Anything, mock.IsType(pipeline.Pipeline{})).Return(nil)

				svc.jobORM.On("FindJobByExternalJobID", mock.Anything, externalJobID).Return(job.Job{}, sql.ErrNoRows)
				svc.jobORM.On("FindJobIDByStreamID", mock.Anything, mock.Anything).Return(int32(0), sql.ErrNoRows)

				svc.spawner.
					On("CreateJob",
						mock.Anything,
						mock.Anything,
						mock.MatchedBy(func(j *job.Job) bool {
							return j.Name.String == streamName
						}),
					).
					Run(func(args mock.Arguments) { (args.Get(2).(*job.Job)).ID = 1 }).
					Return(nil)
				svc.orm.On("ApproveSpec",
					mock.Anything,
					cancelledSpec.ID,
					externalJobID,
				).Return(nil)
				svc.fmsClient.On("ApprovedJob",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&proto.ApprovedJobRequest{
						Uuid:    jp.RemoteUUID.String(),
						Version: int64(spec.Version),
					},
				).Return(&proto.ApprovedJobResponse{}, nil)
				svc.orm.On("CountJobProposalsByStatus", mock.Anything).Return(&feeds.JobProposalCounts{}, nil)
				svc.orm.On("WithDataSource", mock.Anything).Return(feeds.ORM(svc.orm))
				svc.jobORM.On("WithDataSource", mock.Anything).Return(job.ORM(svc.jobORM))
			},
			id:    cancelledSpec.ID,
			force: false,
		},
		{
			name: "cancelled spec failed not latest spec",
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", mock.Anything, cancelledSpec.ID, mock.Anything).Return(cancelledSpec, nil)
				svc.orm.On("GetJobProposal", mock.Anything, jp.ID).Return(jp, nil)
				svc.orm.On("GetLatestSpec", mock.Anything, cancelledSpec.JobProposalID).Return(&feeds.JobProposalSpec{
					ID:            21,
					Status:        feeds.SpecStatusPending,
					JobProposalID: jp.ID,
					Version:       2,
					Definition:    StreamTestSpecTemplate,
				}, nil)
			},
			id:      cancelledSpec.ID,
			force:   false,
			wantErr: "cannot approve a cancelled spec",
		},
		{
			name: "rejected spec failed cannot be approved",
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", mock.Anything, cancelledSpec.ID, mock.Anything).Return(rejectedSpec, nil)
				svc.orm.On("GetJobProposal", mock.Anything, jp.ID).Return(jp, nil)
			},
			id:      rejectedSpec.ID,
			force:   false,
			wantErr: "cannot approve a rejected spec",
		},
		{
			name:        "already existing job replacement error",
			httpTimeout: commonconfig.MustNewDuration(1 * time.Minute),
			before: func(svc *TestService) {
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.orm.On("GetSpec", mock.Anything, spec.ID).Return(spec, nil)
				svc.orm.On("GetJobProposal", mock.Anything, jp.ID).Return(jp, nil)
				svc.jobORM.On("AssertBridgesExist", mock.Anything, mock.IsType(pipeline.Pipeline{})).Return(nil)
				svc.jobORM.On("FindJobByExternalJobID", mock.Anything, externalJobID).Return(job.Job{}, sql.ErrNoRows)
				svc.jobORM.On("FindJobIDByStreamID", mock.Anything, mock.Anything).Return(j.ID, nil)
				svc.orm.On("WithDataSource", mock.Anything).Return(feeds.ORM(svc.orm))
				svc.jobORM.On("WithDataSource", mock.Anything).Return(job.ORM(svc.jobORM))
			},
			id:      spec.ID,
			force:   false,
			wantErr: "could not approve job proposal: a job for this contract address already exists - please use the 'force' option to replace it",
		},
		{
			name:        "already existing self managed job replacement success if forced without feedID",
			httpTimeout: commonconfig.MustNewDuration(1 * time.Minute),
			before: func(svc *TestService) {
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.orm.On("GetSpec", mock.Anything, spec.ID).Return(spec, nil)
				svc.orm.On("GetJobProposal", mock.Anything, jp.ID).Return(jp, nil)
				svc.jobORM.On("AssertBridgesExist", mock.Anything, mock.IsType(pipeline.Pipeline{})).Return(nil)
				svc.orm.EXPECT().GetApprovedSpec(mock.Anything, jp.ID).Return(nil, sql.ErrNoRows)
				svc.jobORM.On("FindJobByExternalJobID", mock.Anything, externalJobID).Return(job.Job{}, sql.ErrNoRows)
				svc.jobORM.On("FindJobIDByStreamID", mock.Anything, mock.Anything).Return(j.ID, nil)
				svc.spawner.On("DeleteJob", mock.Anything, mock.Anything, j.ID).Return(nil)

				svc.spawner.
					On("CreateJob",
						mock.Anything,
						mock.Anything,
						mock.MatchedBy(func(j *job.Job) bool {
							return j.Name.String == streamName
						}),
					).
					Run(func(args mock.Arguments) { (args.Get(2).(*job.Job)).ID = 1 }).
					Return(nil)
				svc.orm.On("ApproveSpec",
					mock.Anything,
					spec.ID,
					externalJobID,
				).Return(nil)
				svc.fmsClient.On("ApprovedJob",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&proto.ApprovedJobRequest{
						Uuid:    jp.RemoteUUID.String(),
						Version: int64(spec.Version),
					},
				).Return(&proto.ApprovedJobResponse{}, nil)
				svc.orm.On("CountJobProposalsByStatus", mock.Anything).Return(&feeds.JobProposalCounts{}, nil)
				svc.orm.On("WithDataSource", mock.Anything).Return(feeds.ORM(svc.orm))
				svc.jobORM.On("WithDataSource", mock.Anything).Return(job.ORM(svc.jobORM))
			},
			id:    spec.ID,
			force: true,
		},
		{
			name:        "already existing self managed job replacement success if forced with feedID",
			httpTimeout: commonconfig.MustNewDuration(1 * time.Minute),
			before: func(svc *TestService) {
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.orm.On("GetSpec", mock.Anything, spec.ID).Return(&feeds.JobProposalSpec{
					ID:            20,
					Status:        feeds.SpecStatusPending,
					JobProposalID: jp.ID,
					Version:       1,
					Definition:    fmt.Sprintf(StreamTestSpecTemplate, streamName, externalJobID.String(), streamID),
				}, nil)
				svc.orm.On("GetJobProposal", mock.Anything, jp.ID).Return(jp, nil)
				svc.jobORM.On("AssertBridgesExist", mock.Anything, mock.IsType(pipeline.Pipeline{})).Return(nil)
				svc.orm.EXPECT().GetApprovedSpec(mock.Anything, jp.ID).Return(nil, sql.ErrNoRows)
				svc.jobORM.On("FindJobByExternalJobID", mock.Anything, externalJobID).Return(job.Job{}, sql.ErrNoRows)
				svc.jobORM.On("FindJobIDByStreamID", mock.Anything, mock.Anything).Return(j.ID, nil)
				svc.spawner.On("DeleteJob", mock.Anything, mock.Anything, j.ID).Return(nil)

				svc.spawner.
					On("CreateJob",
						mock.Anything,
						mock.Anything,
						mock.MatchedBy(func(j *job.Job) bool {
							return j.Name.String == streamName
						}),
					).
					Run(func(args mock.Arguments) { (args.Get(2).(*job.Job)).ID = 1 }).
					Return(nil)
				svc.orm.On("ApproveSpec",
					mock.Anything,
					spec.ID,
					externalJobID,
				).Return(nil)
				svc.fmsClient.On("ApprovedJob",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&proto.ApprovedJobRequest{
						Uuid:    jp.RemoteUUID.String(),
						Version: int64(spec.Version),
					},
				).Return(&proto.ApprovedJobResponse{}, nil)
				svc.orm.On("CountJobProposalsByStatus", mock.Anything).Return(&feeds.JobProposalCounts{}, nil)
				svc.orm.On("WithDataSource", mock.Anything).Return(feeds.ORM(svc.orm))
				svc.jobORM.On("WithDataSource", mock.Anything).Return(job.ORM(svc.jobORM))
			},
			id:    spec.ID,
			force: true,
		},
		{
			name:        "already existing FMS managed job replacement success if forced",
			httpTimeout: commonconfig.MustNewDuration(1 * time.Minute),
			before: func(svc *TestService) {
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.orm.On("GetSpec", mock.Anything, spec.ID).Return(spec, nil)
				svc.orm.On("GetJobProposal", mock.Anything, jp.ID).Return(jp, nil)
				svc.jobORM.On("AssertBridgesExist", mock.Anything, mock.IsType(pipeline.Pipeline{})).Return(nil)
				svc.orm.EXPECT().GetApprovedSpec(mock.Anything, jp.ID).Return(&feeds.JobProposalSpec{ID: 100}, nil)
				svc.orm.EXPECT().CancelSpec(mock.Anything, int64(100)).Return(nil)
				svc.jobORM.On("FindJobByExternalJobID", mock.Anything, externalJobID).Return(job.Job{}, sql.ErrNoRows)
				svc.jobORM.On("FindJobIDByStreamID", mock.Anything, mock.Anything).Return(j.ID, nil)
				svc.spawner.On("DeleteJob", mock.Anything, mock.Anything, j.ID).Return(nil)

				svc.spawner.
					On("CreateJob",
						mock.Anything,
						mock.Anything,
						mock.MatchedBy(func(j *job.Job) bool {
							return j.Name.String == streamName
						}),
					).
					Run(func(args mock.Arguments) { (args.Get(2).(*job.Job)).ID = 1 }).
					Return(nil)
				svc.orm.On("ApproveSpec",
					mock.Anything,
					spec.ID,
					externalJobID,
				).Return(nil)
				svc.fmsClient.On("ApprovedJob",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&proto.ApprovedJobRequest{
						Uuid:    jp.RemoteUUID.String(),
						Version: int64(spec.Version),
					},
				).Return(&proto.ApprovedJobResponse{}, nil)
				svc.orm.On("CountJobProposalsByStatus", mock.Anything).Return(&feeds.JobProposalCounts{}, nil)
				svc.orm.On("WithDataSource", mock.Anything).Return(feeds.ORM(svc.orm))
				svc.jobORM.On("WithDataSource", mock.Anything).Return(job.ORM(svc.jobORM))
			},
			id:    spec.ID,
			force: true,
		},
		{
			name: "spec does not exist",
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", mock.Anything, spec.ID).Return(nil, errors.New("Not Found"))
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
				svc.orm.On("GetSpec", mock.Anything, spec.ID).Return(aspec, nil)
				svc.orm.On("GetJobProposal", mock.Anything, jp.ID).Return(jp, nil)
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
				svc.orm.On("GetSpec", mock.Anything, rspec.ID, mock.Anything).Return(rspec, nil)
				svc.orm.On("GetJobProposal", mock.Anything, jp.ID).Return(jp, nil)
			},
			id:      spec.ID,
			force:   false,
			wantErr: "cannot approve a rejected spec",
		},
		{
			name: "job proposal does not exist",
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", mock.Anything, spec.ID).Return(spec, nil)
				svc.orm.On("GetJobProposal", mock.Anything, jp.ID).Return(nil, errors.New("Not Found"))
			},
			id:      spec.ID,
			wantErr: "orm: job proposal: Not Found",
		},
		{
			name:        "bridges do not exist",
			httpTimeout: commonconfig.MustNewDuration(1 * time.Minute),
			before: func(svc *TestService) {
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.orm.On("GetSpec", mock.Anything, spec.ID).Return(spec, nil)
				svc.orm.On("GetJobProposal", mock.Anything, jp.ID).Return(jp, nil)
				svc.jobORM.On("AssertBridgesExist", mock.Anything, mock.IsType(pipeline.Pipeline{})).Return(errors.New("bridges do not exist"))
			},
			id:      spec.ID,
			wantErr: "failed to approve job spec due to bridge check: bridges do not exist",
		},
		{
			name: "rpc client not connected",
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", mock.Anything, spec.ID).Return(spec, nil)
				svc.orm.On("GetJobProposal", mock.Anything, jp.ID).Return(jp, nil)
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(nil, errors.New("Not Connected"))
			},
			id:      spec.ID,
			force:   false,
			wantErr: "fms rpc client: Not Connected",
		},
		{
			name:        "create job error",
			httpTimeout: commonconfig.MustNewDuration(1 * time.Minute),
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", mock.Anything, spec.ID).Return(spec, nil)
				svc.orm.On("GetJobProposal", mock.Anything, jp.ID).Return(jp, nil)
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.jobORM.On("AssertBridgesExist", mock.Anything, mock.IsType(pipeline.Pipeline{})).Return(nil)

				svc.jobORM.On("FindJobByExternalJobID", mock.Anything, externalJobID).Return(job.Job{}, sql.ErrNoRows)
				svc.jobORM.On("FindJobIDByStreamID", mock.Anything, mock.Anything).Return(int32(0), sql.ErrNoRows)

				svc.spawner.
					On("CreateJob",
						mock.Anything,
						mock.Anything,
						mock.MatchedBy(func(j *job.Job) bool {
							return j.Name.String == streamName
						}),
					).
					Return(errors.New("could not save"))
				svc.orm.On("WithDataSource", mock.Anything).Return(feeds.ORM(svc.orm))
				svc.jobORM.On("WithDataSource", mock.Anything).Return(job.ORM(svc.jobORM))
			},
			id:      spec.ID,
			force:   false,
			wantErr: "could not approve job proposal: could not save",
		},
		{
			name:        "approve spec orm error",
			httpTimeout: commonconfig.MustNewDuration(1 * time.Minute),
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", mock.Anything, spec.ID).Return(spec, nil)
				svc.orm.On("GetJobProposal", mock.Anything, jp.ID).Return(jp, nil)
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.jobORM.On("AssertBridgesExist", mock.Anything, mock.IsType(pipeline.Pipeline{})).Return(nil)

				svc.jobORM.On("FindJobByExternalJobID", mock.Anything, externalJobID).Return(job.Job{}, sql.ErrNoRows)
				svc.jobORM.On("FindJobIDByStreamID", mock.Anything, mock.Anything).Return(int32(0), sql.ErrNoRows)

				svc.spawner.
					On("CreateJob",
						mock.Anything,
						mock.Anything,
						mock.MatchedBy(func(j *job.Job) bool {
							return j.Name.String == streamName
						}),
					).
					Run(func(args mock.Arguments) { (args.Get(2).(*job.Job)).ID = 1 }).
					Return(nil)
				svc.orm.On("ApproveSpec",
					mock.Anything,
					spec.ID,
					externalJobID,
				).Return(errors.New("failure"))
				svc.orm.On("WithDataSource", mock.Anything).Return(feeds.ORM(svc.orm))
				svc.jobORM.On("WithDataSource", mock.Anything).Return(job.ORM(svc.jobORM))
			},
			id:      spec.ID,
			force:   false,
			wantErr: "could not approve job proposal: failure",
		},
		{
			name:        "fms call error",
			httpTimeout: commonconfig.MustNewDuration(1 * time.Minute),
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", mock.Anything, spec.ID).Return(spec, nil)
				svc.orm.On("GetJobProposal", mock.Anything, jp.ID).Return(jp, nil)
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.jobORM.On("AssertBridgesExist", mock.Anything, mock.IsType(pipeline.Pipeline{})).Return(nil)

				svc.jobORM.On("FindJobByExternalJobID", mock.Anything, externalJobID).Return(job.Job{}, sql.ErrNoRows)
				svc.jobORM.On("FindJobIDByStreamID", mock.Anything, mock.Anything).Return(int32(0), sql.ErrNoRows)

				svc.spawner.
					On("CreateJob",
						mock.Anything,
						mock.Anything,
						mock.MatchedBy(func(j *job.Job) bool {
							return j.Name.String == streamName
						}),
					).
					Run(func(args mock.Arguments) { (args.Get(2).(*job.Job)).ID = 1 }).
					Return(nil)
				svc.orm.On("ApproveSpec",
					mock.Anything,
					spec.ID,
					externalJobID,
				).Return(nil)
				svc.fmsClient.On("ApprovedJob",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&proto.ApprovedJobRequest{
						Uuid:    jp.RemoteUUID.String(),
						Version: int64(spec.Version),
					},
				).Return(nil, errors.New("failure"))
				svc.orm.On("WithDataSource", mock.Anything).Return(feeds.ORM(svc.orm))
				svc.jobORM.On("WithDataSource", mock.Anything).Return(job.ORM(svc.jobORM))
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
				c.OCR2.Enabled = testutils.Ptr(true)
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
	address := "0x613a38AC1659769640aaE063C651F48E0250454C"
	feedIDHex := "0x0000000000000000000000000000000000000000000000000000000000000001"
	feedID := common.HexToHash(feedIDHex)
	externalJobID := uuid.New()

	var (
		ctx  = testutils.Context(t)
		defn = `
name = 'LINK / ETH | version 3 | contract 0x0000000000000000000000000000000000000000'
type = 'bootstrap'
schemaVersion = 1
contractID = '0x613a38AC1659769640aaE063C651F48E0250454C'
externalJobID = '%s'
relay = 'evm'

[relayConfig]
chainID = 0
`
		defn2 = `
name = 'LINK / ETH | version 3 | contract 0x0000000000000000000000000000000000000000'
type = 'bootstrap'
schemaVersion = 1
contractID = '0x613a38AC1659769640aaE063C651F48E0250454C'
externalJobID = '%s'
feedID = '%s'
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
			Definition:    fmt.Sprintf(defn, externalJobID.String()),
		}
		rejectedSpec = &feeds.JobProposalSpec{
			ID:            20,
			Status:        feeds.SpecStatusRejected,
			JobProposalID: jp.ID,
			Version:       1,
			Definition:    fmt.Sprintf(defn, externalJobID.String()),
		}
		cancelledSpec = &feeds.JobProposalSpec{
			ID:            20,
			Status:        feeds.SpecStatusCancelled,
			JobProposalID: jp.ID,
			Version:       1,
			Definition:    fmt.Sprintf(defn, externalJobID.String()),
		}
		j = job.Job{
			ID:            1,
			ExternalJobID: externalJobID,
		}
	)

	testCases := []struct {
		name        string
		httpTimeout *commonconfig.Duration
		before      func(svc *TestService)
		id          int64
		force       bool
		wantErr     string
	}{
		{
			name:        "pending job success",
			httpTimeout: commonconfig.MustNewDuration(1 * time.Minute),
			before: func(svc *TestService) {
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.orm.On("GetSpec", mock.Anything, spec.ID).Return(spec, nil)
				svc.orm.On("GetJobProposal", mock.Anything, jp.ID).Return(jp, nil)
				svc.jobORM.On("AssertBridgesExist", mock.Anything, mock.IsType(pipeline.Pipeline{})).Return(nil)
				svc.jobORM.On("FindJobByExternalJobID", mock.Anything, externalJobID).Return(job.Job{}, sql.ErrNoRows)
				svc.jobORM.On("FindOCR2JobIDByAddress", mock.Anything, address, (*common.Hash)(nil)).Return(int32(0), sql.ErrNoRows)

				svc.spawner.
					On("CreateJob",
						mock.Anything,
						mock.Anything,
						mock.MatchedBy(func(j *job.Job) bool {
							return j.Name.String == "LINK / ETH | version 3 | contract 0x0000000000000000000000000000000000000000"
						}),
					).
					Run(func(args mock.Arguments) { (args.Get(2).(*job.Job)).ID = 1 }).
					Return(nil)
				svc.orm.On("ApproveSpec",
					mock.Anything,
					spec.ID,
					externalJobID,
				).Return(nil)
				svc.fmsClient.On("ApprovedJob",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&proto.ApprovedJobRequest{
						Uuid:    jp.RemoteUUID.String(),
						Version: int64(spec.Version),
					},
				).Return(&proto.ApprovedJobResponse{}, nil)
				svc.orm.On("CountJobProposalsByStatus", mock.Anything).Return(&feeds.JobProposalCounts{}, nil)
				svc.orm.On("WithDataSource", mock.Anything).Return(feeds.ORM(svc.orm))
				svc.jobORM.On("WithDataSource", mock.Anything).Return(job.ORM(svc.jobORM))
			},
			id:    spec.ID,
			force: false,
		},
		{
			name:        "cancelled spec success when it is the latest spec",
			httpTimeout: commonconfig.MustNewDuration(1 * time.Minute),
			before: func(svc *TestService) {
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.orm.On("GetSpec", mock.Anything, cancelledSpec.ID, mock.Anything).Return(cancelledSpec, nil)
				svc.orm.On("GetJobProposal", mock.Anything, jp.ID).Return(jp, nil)
				svc.orm.On("GetLatestSpec", mock.Anything, cancelledSpec.JobProposalID).Return(cancelledSpec, nil)
				svc.jobORM.On("AssertBridgesExist", mock.Anything, mock.IsType(pipeline.Pipeline{})).Return(nil)

				svc.jobORM.On("FindJobByExternalJobID", mock.Anything, externalJobID).Return(job.Job{}, sql.ErrNoRows)
				svc.jobORM.On("FindOCR2JobIDByAddress", mock.Anything, address, (*common.Hash)(nil)).Return(int32(0), sql.ErrNoRows)

				svc.spawner.
					On("CreateJob",
						mock.Anything,
						mock.Anything,
						mock.MatchedBy(func(j *job.Job) bool {
							return j.Name.String == "LINK / ETH | version 3 | contract 0x0000000000000000000000000000000000000000"
						}),
					).
					Run(func(args mock.Arguments) { (args.Get(2).(*job.Job)).ID = 1 }).
					Return(nil)
				svc.orm.On("ApproveSpec",
					mock.Anything,
					cancelledSpec.ID,
					externalJobID,
				).Return(nil)
				svc.fmsClient.On("ApprovedJob",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&proto.ApprovedJobRequest{
						Uuid:    jp.RemoteUUID.String(),
						Version: int64(spec.Version),
					},
				).Return(&proto.ApprovedJobResponse{}, nil)
				svc.orm.On("CountJobProposalsByStatus", mock.Anything).Return(&feeds.JobProposalCounts{}, nil)
				svc.orm.On("WithDataSource", mock.Anything).Return(feeds.ORM(svc.orm))
				svc.jobORM.On("WithDataSource", mock.Anything).Return(job.ORM(svc.jobORM))
			},
			id:    cancelledSpec.ID,
			force: false,
		},
		{
			name: "cancelled spec failed not latest spec",
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", mock.Anything, cancelledSpec.ID, mock.Anything).Return(cancelledSpec, nil)
				svc.orm.On("GetJobProposal", mock.Anything, jp.ID).Return(jp, nil)
				svc.orm.On("GetLatestSpec", mock.Anything, cancelledSpec.JobProposalID).Return(&feeds.JobProposalSpec{
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
				svc.orm.On("GetSpec", mock.Anything, cancelledSpec.ID, mock.Anything).Return(rejectedSpec, nil)
				svc.orm.On("GetJobProposal", mock.Anything, jp.ID).Return(jp, nil)
			},
			id:      rejectedSpec.ID,
			force:   false,
			wantErr: "cannot approve a rejected spec",
		},
		{
			name:        "already existing job replacement error",
			httpTimeout: commonconfig.MustNewDuration(1 * time.Minute),
			before: func(svc *TestService) {
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.orm.On("GetSpec", mock.Anything, spec.ID).Return(spec, nil)
				svc.orm.On("GetJobProposal", mock.Anything, jp.ID).Return(jp, nil)
				svc.jobORM.On("AssertBridgesExist", mock.Anything, mock.IsType(pipeline.Pipeline{})).Return(nil)
				svc.jobORM.On("FindJobByExternalJobID", mock.Anything, externalJobID).Return(job.Job{}, sql.ErrNoRows)
				svc.jobORM.On("FindOCR2JobIDByAddress", mock.Anything, address, (*common.Hash)(nil)).Return(j.ID, nil)
				svc.orm.On("WithDataSource", mock.Anything).Return(feeds.ORM(svc.orm))
				svc.jobORM.On("WithDataSource", mock.Anything).Return(job.ORM(svc.jobORM))
			},
			id:      spec.ID,
			force:   false,
			wantErr: "could not approve job proposal: a job for this contract address already exists - please use the 'force' option to replace it",
		},
		{
			name:        "already existing self managed job replacement success if forced without feedID",
			httpTimeout: commonconfig.MustNewDuration(1 * time.Minute),
			before: func(svc *TestService) {
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.orm.On("GetSpec", mock.Anything, spec.ID).Return(spec, nil)
				svc.orm.On("GetJobProposal", mock.Anything, jp.ID).Return(jp, nil)
				svc.jobORM.On("AssertBridgesExist", mock.Anything, mock.IsType(pipeline.Pipeline{})).Return(nil)
				svc.orm.EXPECT().GetApprovedSpec(mock.Anything, jp.ID).Return(nil, sql.ErrNoRows)
				svc.jobORM.On("FindJobByExternalJobID", mock.Anything, externalJobID).Return(job.Job{}, sql.ErrNoRows)
				svc.jobORM.On("FindOCR2JobIDByAddress", mock.Anything, address, (*common.Hash)(nil)).Return(j.ID, nil)
				svc.spawner.On("DeleteJob", mock.Anything, mock.Anything, j.ID).Return(nil)

				svc.spawner.
					On("CreateJob",
						mock.Anything,
						mock.Anything,
						mock.MatchedBy(func(j *job.Job) bool {
							return j.Name.String == "LINK / ETH | version 3 | contract 0x0000000000000000000000000000000000000000"
						}),
					).
					Run(func(args mock.Arguments) { (args.Get(2).(*job.Job)).ID = 1 }).
					Return(nil)
				svc.orm.On("ApproveSpec",
					mock.Anything,
					spec.ID,
					externalJobID,
				).Return(nil)
				svc.fmsClient.On("ApprovedJob",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&proto.ApprovedJobRequest{
						Uuid:    jp.RemoteUUID.String(),
						Version: int64(spec.Version),
					},
				).Return(&proto.ApprovedJobResponse{}, nil)
				svc.orm.On("CountJobProposalsByStatus", mock.Anything).Return(&feeds.JobProposalCounts{}, nil)
				svc.orm.On("WithDataSource", mock.Anything).Return(feeds.ORM(svc.orm))
				svc.jobORM.On("WithDataSource", mock.Anything).Return(job.ORM(svc.jobORM))
			},
			id:    spec.ID,
			force: true,
		},
		{
			name:        "already existing self managed job replacement success if forced with feedID",
			httpTimeout: commonconfig.MustNewDuration(1 * time.Minute),
			before: func(svc *TestService) {
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.orm.On("GetSpec", mock.Anything, spec.ID).Return(&feeds.JobProposalSpec{
					ID:            20,
					Status:        feeds.SpecStatusPending,
					JobProposalID: jp.ID,
					Version:       1,
					Definition:    fmt.Sprintf(defn2, externalJobID.String(), feedID),
				}, nil)
				svc.orm.On("GetJobProposal", mock.Anything, jp.ID).Return(jp, nil)
				svc.jobORM.On("AssertBridgesExist", mock.Anything, mock.IsType(pipeline.Pipeline{})).Return(nil)
				svc.orm.EXPECT().GetApprovedSpec(mock.Anything, jp.ID).Return(nil, sql.ErrNoRows)
				svc.jobORM.On("FindJobByExternalJobID", mock.Anything, externalJobID).Return(job.Job{}, sql.ErrNoRows)
				svc.jobORM.On("FindOCR2JobIDByAddress", mock.Anything, address, &feedID).Return(j.ID, nil)
				svc.spawner.On("DeleteJob", mock.Anything, mock.Anything, j.ID).Return(nil)

				svc.spawner.
					On("CreateJob",
						mock.Anything,
						mock.Anything,
						mock.MatchedBy(func(j *job.Job) bool {
							return j.Name.String == "LINK / ETH | version 3 | contract 0x0000000000000000000000000000000000000000"
						}),
					).
					Run(func(args mock.Arguments) { (args.Get(2).(*job.Job)).ID = 1 }).
					Return(nil)
				svc.orm.On("ApproveSpec",
					mock.Anything,
					spec.ID,
					externalJobID,
				).Return(nil)
				svc.fmsClient.On("ApprovedJob",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&proto.ApprovedJobRequest{
						Uuid:    jp.RemoteUUID.String(),
						Version: int64(spec.Version),
					},
				).Return(&proto.ApprovedJobResponse{}, nil)
				svc.orm.On("CountJobProposalsByStatus", mock.Anything).Return(&feeds.JobProposalCounts{}, nil)
				svc.orm.On("WithDataSource", mock.Anything).Return(feeds.ORM(svc.orm))
				svc.jobORM.On("WithDataSource", mock.Anything).Return(job.ORM(svc.jobORM))
			},
			id:    spec.ID,
			force: true,
		},
		{
			name:        "already existing FMS managed job replacement success if forced",
			httpTimeout: commonconfig.MustNewDuration(1 * time.Minute),
			before: func(svc *TestService) {
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.orm.On("GetSpec", mock.Anything, spec.ID).Return(spec, nil)
				svc.orm.On("GetJobProposal", mock.Anything, jp.ID).Return(jp, nil)
				svc.jobORM.On("AssertBridgesExist", mock.Anything, mock.IsType(pipeline.Pipeline{})).Return(nil)
				svc.orm.EXPECT().GetApprovedSpec(mock.Anything, jp.ID).Return(&feeds.JobProposalSpec{ID: 100}, nil)
				svc.orm.EXPECT().CancelSpec(mock.Anything, int64(100)).Return(nil)
				svc.jobORM.On("FindJobByExternalJobID", mock.Anything, externalJobID).Return(job.Job{}, sql.ErrNoRows)
				svc.jobORM.On("FindOCR2JobIDByAddress", mock.Anything, address, (*common.Hash)(nil)).Return(j.ID, nil)
				svc.spawner.On("DeleteJob", mock.Anything, mock.Anything, j.ID).Return(nil)

				svc.spawner.
					On("CreateJob",
						mock.Anything,
						mock.Anything,
						mock.MatchedBy(func(j *job.Job) bool {
							return j.Name.String == "LINK / ETH | version 3 | contract 0x0000000000000000000000000000000000000000"
						}),
					).
					Run(func(args mock.Arguments) { (args.Get(2).(*job.Job)).ID = 1 }).
					Return(nil)
				svc.orm.On("ApproveSpec",
					mock.Anything,
					spec.ID,
					externalJobID,
				).Return(nil)
				svc.fmsClient.On("ApprovedJob",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&proto.ApprovedJobRequest{
						Uuid:    jp.RemoteUUID.String(),
						Version: int64(spec.Version),
					},
				).Return(&proto.ApprovedJobResponse{}, nil)
				svc.orm.On("CountJobProposalsByStatus", mock.Anything).Return(&feeds.JobProposalCounts{}, nil)
				svc.orm.On("WithDataSource", mock.Anything).Return(feeds.ORM(svc.orm))
				svc.jobORM.On("WithDataSource", mock.Anything).Return(job.ORM(svc.jobORM))
			},
			id:    spec.ID,
			force: true,
		},
		{
			name: "spec does not exist",
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", mock.Anything, spec.ID).Return(nil, errors.New("Not Found"))
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
				svc.orm.On("GetSpec", mock.Anything, spec.ID).Return(aspec, nil)
				svc.orm.On("GetJobProposal", mock.Anything, jp.ID).Return(jp, nil)
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
				svc.orm.On("GetSpec", mock.Anything, rspec.ID).Return(rspec, nil)
				svc.orm.On("GetJobProposal", mock.Anything, jp.ID).Return(jp, nil)
			},
			id:      spec.ID,
			force:   false,
			wantErr: "cannot approve a rejected spec",
		},
		{
			name: "job proposal does not exist",
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", mock.Anything, spec.ID).Return(spec, nil)
				svc.orm.On("GetJobProposal", mock.Anything, jp.ID).Return(nil, errors.New("Not Found"))
			},
			id:      spec.ID,
			wantErr: "orm: job proposal: Not Found",
		},
		{
			name:        "bridges do not exist",
			httpTimeout: commonconfig.MustNewDuration(1 * time.Minute),
			before: func(svc *TestService) {
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.orm.On("GetSpec", mock.Anything, spec.ID).Return(spec, nil)
				svc.orm.On("GetJobProposal", mock.Anything, jp.ID).Return(jp, nil)
				svc.jobORM.On("AssertBridgesExist", mock.Anything, mock.IsType(pipeline.Pipeline{})).Return(errors.New("bridges do not exist"))
			},
			id:      spec.ID,
			wantErr: "failed to approve job spec due to bridge check: bridges do not exist",
		},
		{
			name: "rpc client not connected",
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", mock.Anything, spec.ID).Return(spec, nil)
				svc.orm.On("GetJobProposal", mock.Anything, jp.ID).Return(jp, nil)
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(nil, errors.New("Not Connected"))
			},
			id:      spec.ID,
			force:   false,
			wantErr: "fms rpc client: Not Connected",
		},
		{
			name:        "create job error",
			httpTimeout: commonconfig.MustNewDuration(1 * time.Minute),
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", mock.Anything, spec.ID).Return(spec, nil)
				svc.orm.On("GetJobProposal", mock.Anything, jp.ID).Return(jp, nil)
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.jobORM.On("AssertBridgesExist", mock.Anything, mock.IsType(pipeline.Pipeline{})).Return(nil)

				svc.jobORM.On("FindJobByExternalJobID", mock.Anything, externalJobID).Return(job.Job{}, sql.ErrNoRows)
				svc.jobORM.On("FindOCR2JobIDByAddress", mock.Anything, address, (*common.Hash)(nil)).Return(int32(0), sql.ErrNoRows)

				svc.spawner.
					On("CreateJob",
						mock.Anything,
						mock.Anything,
						mock.MatchedBy(func(j *job.Job) bool {
							return j.Name.String == "LINK / ETH | version 3 | contract 0x0000000000000000000000000000000000000000"
						}),
					).
					Return(errors.New("could not save"))
				svc.orm.On("WithDataSource", mock.Anything).Return(feeds.ORM(svc.orm))
				svc.jobORM.On("WithDataSource", mock.Anything).Return(job.ORM(svc.jobORM))
			},
			id:      spec.ID,
			force:   false,
			wantErr: "could not approve job proposal: could not save",
		},
		{
			name:        "approve spec orm error",
			httpTimeout: commonconfig.MustNewDuration(1 * time.Minute),
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", mock.Anything, spec.ID).Return(spec, nil)
				svc.orm.On("GetJobProposal", mock.Anything, jp.ID).Return(jp, nil)
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.jobORM.On("AssertBridgesExist", mock.Anything, mock.IsType(pipeline.Pipeline{})).Return(nil)

				svc.jobORM.On("FindJobByExternalJobID", mock.Anything, externalJobID).Return(job.Job{}, sql.ErrNoRows)
				svc.jobORM.On("FindOCR2JobIDByAddress", mock.Anything, address, (*common.Hash)(nil), mock.Anything).Return(int32(0), sql.ErrNoRows)

				svc.spawner.
					On("CreateJob",
						mock.Anything,
						mock.Anything,
						mock.MatchedBy(func(j *job.Job) bool {
							return j.Name.String == "LINK / ETH | version 3 | contract 0x0000000000000000000000000000000000000000"
						}),
					).
					Run(func(args mock.Arguments) { (args.Get(2).(*job.Job)).ID = 1 }).
					Return(nil)
				svc.orm.On("ApproveSpec",
					mock.Anything,
					spec.ID,
					externalJobID,
				).Return(errors.New("failure"))
				svc.orm.On("WithDataSource", mock.Anything).Return(feeds.ORM(svc.orm))
				svc.jobORM.On("WithDataSource", mock.Anything).Return(job.ORM(svc.jobORM))
			},
			id:      spec.ID,
			force:   false,
			wantErr: "could not approve job proposal: failure",
		},
		{
			name:        "fms call error",
			httpTimeout: commonconfig.MustNewDuration(1 * time.Minute),
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", mock.Anything, spec.ID).Return(spec, nil)
				svc.orm.On("GetJobProposal", mock.Anything, jp.ID).Return(jp, nil)
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.jobORM.On("AssertBridgesExist", mock.Anything, mock.IsType(pipeline.Pipeline{})).Return(nil)

				svc.jobORM.On("FindJobByExternalJobID", mock.Anything, externalJobID).Return(job.Job{}, sql.ErrNoRows)
				svc.jobORM.On("FindOCR2JobIDByAddress", mock.Anything, address, (*common.Hash)(nil)).Return(int32(0), sql.ErrNoRows)

				svc.spawner.
					On("CreateJob",
						mock.Anything,
						mock.Anything,
						mock.MatchedBy(func(j *job.Job) bool {
							return j.Name.String == "LINK / ETH | version 3 | contract 0x0000000000000000000000000000000000000000"
						}),
					).
					Run(func(args mock.Arguments) { (args.Get(2).(*job.Job)).ID = 1 }).
					Return(nil)
				svc.orm.On("ApproveSpec",
					mock.Anything,
					spec.ID,
					externalJobID,
				).Return(nil)
				svc.fmsClient.On("ApprovedJob",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&proto.ApprovedJobRequest{
						Uuid:    jp.RemoteUUID.String(),
						Version: int64(spec.Version),
					},
				).Return(nil, errors.New("failure"))
				svc.orm.On("WithDataSource", mock.Anything).Return(feeds.ORM(svc.orm))
				svc.jobORM.On("WithDataSource", mock.Anything).Return(job.ORM(svc.jobORM))
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
				c.OCR2.Enabled = testutils.Ptr(true)
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
				svc.orm.On("GetSpec", mock.Anything, spec.ID).Return(spec, nil)
				svc.orm.On("GetJobProposal", mock.Anything, jp.ID).Return(jp, nil)
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.orm.On("RejectSpec",
					mock.Anything,
					spec.ID,
				).Return(nil)
				svc.fmsClient.On("RejectedJob",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					&proto.RejectedJobRequest{
						Uuid:    jp.RemoteUUID.String(),
						Version: int64(spec.Version),
					},
				).Return(&proto.RejectedJobResponse{}, nil)
				svc.orm.On("CountJobProposalsByStatus", mock.Anything).Return(&feeds.JobProposalCounts{}, nil)
				transactCall := svc.orm.On("Transact", mock.Anything, mock.Anything)
				transactCall.Run(func(args mock.Arguments) {
					fn := args[1].(func(orm feeds.ORM) error)
					transactCall.ReturnArguments = mock.Arguments{fn(svc.orm)}
				})
			},
		},
		{
			name: "Fails to get spec",
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", mock.Anything, spec.ID).Return(nil, errors.New("failure"))
			},
			wantErr: "failure",
		},
		{
			name: "Cannot be a rejected proposal",
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", mock.Anything, spec.ID).Return(&feeds.JobProposalSpec{
					Status: feeds.SpecStatusRejected,
				}, nil)
			},
			wantErr: "must be a pending job proposal spec",
		},
		{
			name: "Fails to get proposal",
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", mock.Anything, spec.ID).Return(spec, nil)
				svc.orm.On("GetJobProposal", mock.Anything, jp.ID).Return(nil, errors.New("failure"))
			},
			wantErr: "failure",
		},
		{
			name: "FMS not connected",
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", mock.Anything, spec.ID).Return(spec, nil)
				svc.orm.On("GetJobProposal", mock.Anything, jp.ID).Return(jp, nil)
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(nil, errors.New("disconnected"))
			},
			wantErr: "disconnected",
		},
		{
			name: "Fails to update spec",
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", mock.Anything, spec.ID).Return(spec, nil)
				svc.orm.On("GetJobProposal", mock.Anything, jp.ID).Return(jp, nil)
				svc.connMgr.On("GetClient", jp.FeedsManagerID).Return(svc.fmsClient, nil)
				svc.orm.On("RejectSpec", mock.Anything, mock.Anything).Return(errors.New("failure"))
				transactCall := svc.orm.On("Transact", mock.Anything, mock.Anything)
				transactCall.Run(func(args mock.Arguments) {
					fn := args[1].(func(orm feeds.ORM) error)
					transactCall.ReturnArguments = mock.Arguments{fn(svc.orm)}
				})
			},
			wantErr: "failure",
		},
		{
			name: "Fails to update spec",
			before: func(svc *TestService) {
				svc.orm.On("GetSpec", mock.Anything, spec.ID).Return(spec, nil)
				svc.orm.On("GetJobProposal", mock.Anything, jp.ID).Return(jp, nil)
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
				transactCall := svc.orm.On("Transact", mock.Anything, mock.Anything)
				transactCall.Run(func(args mock.Arguments) {
					fn := args[1].(func(orm feeds.ORM) error)
					transactCall.ReturnArguments = mock.Arguments{fn(svc.orm)}
				})
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
					On("GetSpec", mock.Anything, specID, mock.Anything).
					Return(spec, nil)
				svc.orm.On("UpdateSpecDefinition", mock.Anything,
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
					On("GetSpec", mock.Anything, specID, mock.Anything).
					Return(nil, sql.ErrNoRows)
			},
			specID:  specID,
			wantErr: "job proposal spec does not exist: sql: no rows in result set",
		},
		{
			name: "other get errors",
			before: func(svc *TestService) {
				svc.orm.
					On("GetSpec", mock.Anything, specID, mock.Anything).
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
					On("GetSpec", mock.Anything, specID, mock.Anything).
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
		mgr2 = feeds.FeedsManager{
			ID:  2,
			URI: "localhost:2001",
		}
		pubKeyHex = "0f17c3bf72de8beef6e2d17a14c0a972f5d7e0e66e70722373f12b88382d40f9"
	)

	var pubKey crypto.PublicKey
	_, err := hex.Decode([]byte(pubKeyHex), pubKey)
	require.NoError(t, err)

	tests := []struct {
		name                     string
		enableMultiFeedsManagers bool
		beforeFunc               func(svc *TestService)
	}{
		{
			name: "success with a feeds manager connection",
			beforeFunc: func(svc *TestService) {
				svc.csaKeystore.On("EnsureKey", mock.Anything).Return(nil)
				svc.csaKeystore.On("GetAll").Return([]csakey.KeyV2{key}, nil)
				svc.orm.On("ListManagers", mock.Anything).Return([]feeds.FeedsManager{mgr}, nil)
				svc.connMgr.On("IsConnected", mgr.ID).Return(false)
				svc.connMgr.On("Connect", mock.IsType(feeds.ConnectOpts{}))
				svc.connMgr.On("Close")
				svc.orm.On("CountJobProposalsByStatus", mock.Anything).Return(&feeds.JobProposalCounts{}, nil)
			},
		},
		{
			name:                     "success with multiple feeds managers connection",
			enableMultiFeedsManagers: true,
			beforeFunc: func(svc *TestService) {
				svc.csaKeystore.On("EnsureKey", mock.Anything).Return(nil)
				svc.csaKeystore.On("GetAll").Return([]csakey.KeyV2{key}, nil)
				svc.orm.On("ListManagers", mock.Anything).Return([]feeds.FeedsManager{mgr, mgr2}, nil)
				svc.connMgr.On("IsConnected", mgr.ID).Return(false)
				svc.connMgr.On("IsConnected", mgr2.ID).Return(false)
				svc.connMgr.On("Connect", mock.IsType(feeds.ConnectOpts{})).Twice()
				svc.connMgr.On("Close")
				svc.orm.On("CountJobProposalsByStatus", mock.Anything).Return(&feeds.JobProposalCounts{}, nil)
			},
		},
		{
			name: "success with no registered managers",
			beforeFunc: func(svc *TestService) {
				svc.csaKeystore.On("EnsureKey", mock.Anything).Return(nil)
				svc.csaKeystore.On("GetAll").Return([]csakey.KeyV2{key}, nil)
				svc.orm.On("ListManagers", mock.Anything).Return([]feeds.FeedsManager{}, nil)
				svc.connMgr.On("Close")
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := setupTestServiceCfg(t, func(c *chainlink.Config, s *chainlink.Secrets) {
				c.Feature.MultiFeedsManagers = &tt.enableMultiFeedsManagers
			})

			if tt.beforeFunc != nil {
				tt.beforeFunc(svc)
			}

			servicetest.Run(t, svc)
		})
	}
}

func logMessages(logEntries []observer.LoggedEntry) []string {
	messages := make([]string, 0, len(logEntries))
	for _, entry := range logEntries {
		messageWithContext := entry.Message
		contextMap := entry.ContextMap()
		for _, key := range slices.Sorted(maps.Keys(contextMap)) {
			if key == "version" || key == "errVerbose" {
				continue
			}
			messageWithContext += fmt.Sprintf(" %v=\"%v\"", key, entry.ContextMap()[key])
		}

		messages = append(messages, messageWithContext)
	}

	return messages
}

func waitSyncNodeInfoCall(t *testing.T, logs *observer.ObservedLogs) {
	assert.EventuallyWithT(t, func(collect *assert.CollectT) {
		assert.Contains(collect, logMessages(logs.All()), "successfully synced node info")
	}, 1*time.Second, 5*time.Millisecond)
}
