package cltest

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"math/big"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	cryptop2p "github.com/libp2p/go-libp2p-core/crypto"
	p2ppeer "github.com/libp2p/go-libp2p-core/peer"
	"github.com/manyminds/api2go/jsonapi"
	"github.com/onsi/gomega"
	uuid "github.com/satori/go.uuid"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting/types"
	"github.com/smartcontractkit/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
	"github.com/urfave/cli"
	"gopkg.in/guregu/null.v4"

	starkkey "github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/keys"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/auth"
	"github.com/smartcontractkit/chainlink/core/bridges"
	"github.com/smartcontractkit/chainlink/core/chains/evm"
	evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"
	evmconfig "github.com/smartcontractkit/chainlink/core/chains/evm/config"
	v2 "github.com/smartcontractkit/chainlink/core/chains/evm/config/v2"
	"github.com/smartcontractkit/chainlink/core/chains/evm/gas"
	httypes "github.com/smartcontractkit/chainlink/core/chains/evm/headtracker/types"
	evmMocks "github.com/smartcontractkit/chainlink/core/chains/evm/mocks"
	"github.com/smartcontractkit/chainlink/core/chains/evm/txmgr"
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/chains/solana"
	"github.com/smartcontractkit/chainlink/core/chains/starknet"
	"github.com/smartcontractkit/chainlink/core/chains/terra"
	"github.com/smartcontractkit/chainlink/core/cmd"
	"github.com/smartcontractkit/chainlink/core/config"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/configtest"
	configtest2 "github.com/smartcontractkit/chainlink/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/evmtest"
	clhttptest "github.com/smartcontractkit/chainlink/core/internal/testutils/httptest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/keystest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/logger/audit"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/csakey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/dkgencryptkey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/dkgsignkey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ocr2key"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ocrkey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/solkey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/terrakey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/vrfkey"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/services/webhook"
	clsessions "github.com/smartcontractkit/chainlink/core/sessions"
	"github.com/smartcontractkit/chainlink/core/static"
	"github.com/smartcontractkit/chainlink/core/store/dialects"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/web"
	webauth "github.com/smartcontractkit/chainlink/core/web/auth"
	webpresenters "github.com/smartcontractkit/chainlink/core/web/presenters"

	// Force import of pgtest to ensure that txdb is registered as a DB driver
	_ "github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
)

const (
	// APIKey of the fixture API user
	APIKey = "2d25e62eaf9143e993acaf48691564b2"
	// APISecret of the fixture API user.
	APISecret = "1eCP/w0llVkchejFaoBpfIGaLRxZK54lTXBCT22YLW+pdzE4Fafy/XO5LoJ2uwHi"
	// Collection of test fixture DB user emails per role
	APIEmailAdmin    = "apiuser@chainlink.test"
	APIEmailEdit     = "apiuser-edit@chainlink.test"
	APIEmailRun      = "apiuser-run@chainlink.test"
	APIEmailViewOnly = "apiuser-view-only@chainlink.test"
	// Password just a password we use everywhere for testing
	Password = testutils.Password
	// SessionSecret is the hardcoded secret solely used for test
	SessionSecret = "clsession_test_secret"
	// DefaultPeerID is the peer ID of the default p2p key
	DefaultPeerID = "12D3KooWPjceQrSwdWXPyLLeABRXmuqt69Rg3sBYbU1Nft9HyQ6X"
	// DefaultOCRKeyBundleID is the ID of the default ocr key bundle
	DefaultOCRKeyBundleID = "f5bf259689b26f1374efb3c9a9868796953a0f814bb2d39b968d0e61b58620a5"
	// DefaultOCR2KeyBundleID is the ID of the fixture ocr2 key bundle
	DefaultOCR2KeyBundleID = "92be59c45d0d7b192ef88d391f444ea7c78644f8607f567aab11d53668c27a4d"
	// Private key seed of test keys created with `big.NewInt(1)`, representations of value present in `scrub_logs` script
	KeyBigIntSeed = 1
)

var (
	DefaultP2PPeerID p2pkey.PeerID
	FixtureChainID   = *testutils.FixtureChainID
	source           rand.Source

	DefaultCSAKey        = csakey.MustNewV2XXXTestingOnly(big.NewInt(KeyBigIntSeed))
	DefaultOCRKey        = ocrkey.MustNewV2XXXTestingOnly(big.NewInt(KeyBigIntSeed))
	DefaultOCR2Key       = ocr2key.MustNewInsecure(keystest.NewRandReaderFromSeed(KeyBigIntSeed), "evm")
	DefaultP2PKey        = p2pkey.MustNewV2XXXTestingOnly(big.NewInt(KeyBigIntSeed))
	DefaultSolanaKey     = solkey.MustNewInsecure(keystest.NewRandReaderFromSeed(KeyBigIntSeed))
	DefaultTerraKey      = terrakey.MustNewInsecure(keystest.NewRandReaderFromSeed(KeyBigIntSeed))
	DefaultStarkNetKey   = starkkey.MustNewInsecure(keystest.NewRandReaderFromSeed(KeyBigIntSeed))
	DefaultVRFKey        = vrfkey.MustNewV2XXXTestingOnly(big.NewInt(KeyBigIntSeed))
	DefaultDKGSignKey    = dkgsignkey.MustNewXXXTestingOnly(big.NewInt(KeyBigIntSeed))
	DefaultDKGEncryptKey = dkgencryptkey.MustNewXXXTestingOnly(big.NewInt(KeyBigIntSeed))
)

func init() {
	gin.SetMode(gin.TestMode)

	gomega.SetDefaultEventuallyTimeout(testutils.DefaultWaitTimeout)
	gomega.SetDefaultEventuallyPollingInterval(DBPollingInterval)
	gomega.SetDefaultConsistentlyDuration(time.Second)
	gomega.SetDefaultConsistentlyPollingInterval(100 * time.Millisecond)

	logger.InitColor(true)
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		fmt.Printf("[gin] %-6s %-25s --> %s (%d handlers)\n", httpMethod, absolutePath, handlerName, nuHandlers)
	}

	// Seed the random number generator, otherwise separate modules will take
	// the same advisory locks when tested with `go test -p N` for N > 1
	seed := time.Now().UTC().UnixNano()
	fmt.Printf("cltest random seed: %v\n", seed)
	rand.Seed(seed)

	// Also seed the local source
	source = rand.NewSource(seed)
	defaultP2PPeerID, err := p2ppeer.Decode(configtest.DefaultPeerID)
	if err != nil {
		panic(err)
	}
	DefaultP2PPeerID = p2pkey.PeerID(defaultP2PPeerID)
}

func NewRandomPositiveInt64() int64 {
	id := rand.Int63()
	return id
}

func MustRandomBytes(t *testing.T, l int) (b []byte) {
	t.Helper()

	b = make([]byte, l)
	/* #nosec G404 */
	_, err := rand.Read(b)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

type JobPipelineV2TestHelper struct {
	Prm pipeline.ORM
	Jrm job.ORM
	Pr  pipeline.Runner
}

func NewJobPipelineV2(t testing.TB, cfg config.BasicConfig, cc evm.ChainSet, db *sqlx.DB, keyStore keystore.Master, restrictedHTTPClient, unrestrictedHTTPClient *http.Client) JobPipelineV2TestHelper {
	lggr := logger.TestLogger(t)
	prm := pipeline.NewORM(db, lggr, cfg)
	btORM := bridges.NewORM(db, lggr, cfg)
	jrm := job.NewORM(db, cc, prm, btORM, keyStore, lggr, cfg)
	pr := pipeline.NewRunner(prm, btORM, cfg, cc, keyStore.Eth(), keyStore.VRF(), lggr, restrictedHTTPClient, unrestrictedHTTPClient)
	return JobPipelineV2TestHelper{
		prm,
		jrm,
		pr,
	}
}

// NewEthBroadcaster creates a new txmgr.EthBroadcaster for use in testing.
func NewEthBroadcaster(t testing.TB, db *sqlx.DB, ethClient evmclient.Client, keyStore txmgr.KeyStore, config evmconfig.ChainScopedConfig, keyStates []ethkey.State, checkerFactory txmgr.TransmitCheckerFactory) *txmgr.EthBroadcaster {
	t.Helper()
	eventBroadcaster := NewEventBroadcaster(t, config.DatabaseURL())
	err := eventBroadcaster.Start(testutils.Context(t.(*testing.T)))
	require.NoError(t, err)
	t.Cleanup(func() { assert.NoError(t, eventBroadcaster.Close()) })
	lggr := logger.TestLogger(t)
	return txmgr.NewEthBroadcaster(db, ethClient, config, keyStore, eventBroadcaster,
		keyStates, gas.NewFixedPriceEstimator(config, lggr), nil, lggr,
		checkerFactory)
}

func NewEventBroadcaster(t testing.TB, dbURL url.URL) pg.EventBroadcaster {
	lggr := logger.TestLogger(t)
	return pg.NewEventBroadcaster(dbURL, 0, 0, lggr, uuid.NewV4())
}

func NewEthConfirmer(t testing.TB, db *sqlx.DB, ethClient evmclient.Client, config evmconfig.ChainScopedConfig, ks keystore.Eth, keyStates []ethkey.State, fn txmgr.ResumeCallback) *txmgr.EthConfirmer {
	t.Helper()
	lggr := logger.TestLogger(t)
	ec := txmgr.NewEthConfirmer(db, ethClient, config, ks, keyStates,
		gas.NewFixedPriceEstimator(config, lggr), fn, lggr)
	return ec
}

// TestApplication holds the test application and test servers
type TestApplication struct {
	t testing.TB
	*chainlink.ChainlinkApplication
	Logger  logger.Logger
	Server  *httptest.Server
	Started bool
	Backend *backends.SimulatedBackend
	Keys    []ethkey.KeyV2
}

// NewWSServer starts a websocket server which invokes callback for each message received.
// If chainID is set, then eth_chainId calls will be automatically handled.
func NewWSServer(t *testing.T, chainID *big.Int, callback testutils.JSONRPCHandler) string {
	server := testutils.NewWSServer(t, chainID, callback)
	return server.WSURL().String()
}

// Deprecated: use configtest/v2.NewTestGeneralConfig
// https://app.shortcut.com/chainlinklabs/story/33622/remove-legacy-config
func NewTestGeneralConfig(t testing.TB) *configtest.TestGeneralConfig {
	shutdownGracePeriod := testutils.DefaultWaitTimeout
	reaperInterval := time.Duration(0) // disable reaper
	overrides := configtest.GeneralConfigOverrides{
		Dialect:                   dialects.TransactionWrappedPostgres,
		AdvisoryLockID:            null.IntFrom(NewRandomPositiveInt64()),
		P2PEnabled:                null.BoolFrom(false),
		ShutdownGracePeriod:       &shutdownGracePeriod,
		JobPipelineReaperInterval: &reaperInterval,
	}
	return configtest.NewTestGeneralConfigWithOverrides(t, overrides)
}

// NewApplicationEVMDisabled creates a new application with default config but EVM disabled
// Useful for testing controllers
func NewApplicationEVMDisabled(t *testing.T) *TestApplication {
	t.Helper()

	c := configtest2.NewGeneralConfig(t, nil)

	return NewApplicationWithConfig(t, c)
}

// Deprecated: use NewApplicationEVMDisabled
// https://app.shortcut.com/chainlinklabs/story/33622/remove-legacy-config
func NewLegacyApplicationEVMDisabled(t *testing.T) *TestApplication {
	t.Helper()

	c := NewTestGeneralConfig(t)
	c.Overrides.EVMEnabled = null.BoolFrom(false)

	return NewApplicationWithConfig(t, c)
}

// NewApplication creates a New TestApplication along with a NewConfig
// It mocks the keystore with no keys or accounts by default
func NewApplication(t testing.TB, flagsAndDeps ...interface{}) *TestApplication {
	t.Helper()

	c := configtest2.NewGeneralConfig(t, nil)

	return NewApplicationWithConfig(t, c, flagsAndDeps...)
}

// NewApplicationWithKey creates a new TestApplication along with a new config
// It uses the native keystore and will load any keys that are in the database
func NewApplicationWithKey(t *testing.T, flagsAndDeps ...interface{}) *TestApplication {
	t.Helper()

	config := configtest2.NewGeneralConfig(t, nil)
	return NewApplicationWithConfigAndKey(t, config, flagsAndDeps...)
}

// NewApplicationWithConfigAndKey creates a new TestApplication with the given testorm
// it will also provide an unlocked account on the keystore
func NewApplicationWithConfigAndKey(t testing.TB, c config.GeneralConfig, flagsAndDeps ...interface{}) *TestApplication {
	t.Helper()

	app := NewApplicationWithConfig(t, c, flagsAndDeps...)
	require.NoError(t, app.KeyStore.Unlock(Password))
	chainID := *utils.NewBig(&FixtureChainID)
	for _, dep := range flagsAndDeps {
		switch v := dep.(type) {
		case ethkey.KeyV2:
			app.Keys = append(app.Keys, v)
		case p2pkey.KeyV2:
			require.NoError(t, app.GetKeyStore().P2P().Add(v))
		case evmtypes.DBChain:
			chainID = v.ID
		case *utils.Big:
			chainID = *v
		}
	}
	if len(app.Keys) == 0 {
		k, _ := MustInsertRandomKey(t, app.KeyStore.Eth(), 0, chainID)
		app.Keys = []ethkey.KeyV2{k}
	} else {
		id, ks := chainID.ToInt(), app.KeyStore.Eth()
		for _, k := range app.Keys {
			MustAddKeyToKeystore(t, k, id, ks)
		}
	}

	return app
}

const (
	UseRealExternalInitiatorManager = "UseRealExternalInitiatorManager"
)

// NewApplicationWithConfig creates a New TestApplication with specified test config.
// This should only be used in full integration tests. For controller tests, see NewApplicationEVMDisabled.
func NewApplicationWithConfig(t testing.TB, cfg config.GeneralConfig, flagsAndDeps ...interface{}) *TestApplication {
	t.Helper()
	testutils.SkipShortDB(t)

	var lggr logger.Logger
	for _, dep := range flagsAndDeps {
		argLggr, is := dep.(logger.Logger)
		if is {
			lggr = argLggr
			break
		}
	}
	if lggr == nil {
		lggr = logger.TestLogger(t)
	}

	var auditLogger audit.AuditLogger
	for _, dep := range flagsAndDeps {
		audLgger, is := dep.(audit.AuditLogger)
		if is {
			auditLogger = audLgger
			break
		}
	}

	if auditLogger == nil {
		auditLogger = audit.NoopLogger
	}

	var eventBroadcaster pg.EventBroadcaster = pg.NewNullEventBroadcaster()

	url := cfg.DatabaseURL()
	db, err := pg.NewConnection(url.String(), cfg.GetDatabaseDialectConfiguredOrDefault(), cfg)
	require.NoError(t, err)
	t.Cleanup(func() { assert.NoError(t, db.Close()) })

	var ethClient evmclient.Client
	var externalInitiatorManager webhook.ExternalInitiatorManager
	externalInitiatorManager = &webhook.NullExternalInitiatorManager{}
	var useRealExternalInitiatorManager bool
	var chainORM evmtypes.ORM
	for _, flag := range flagsAndDeps {
		switch dep := flag.(type) {
		case evmclient.Client:
			ethClient = dep
		case webhook.ExternalInitiatorManager:
			externalInitiatorManager = dep
		case evmtypes.DBChain:
			if chainORM != nil {
				panic("cannot set more than one chain")
			}
			chainORM = evmtest.NewMockORM([]evmtypes.DBChain{dep}, nil)
		case pg.EventBroadcaster:
			eventBroadcaster = dep
		default:
			switch flag {
			case UseRealExternalInitiatorManager:
				externalInitiatorManager = webhook.NewExternalInitiatorManager(db, clhttptest.NewTestLocalOnlyHTTPClient(), lggr, cfg)
			}

		}
	}
	if ethClient == nil {
		ethClient = evmclient.NewNullClient(cfg.DefaultChainID(), lggr)
	}

	keyStore := keystore.New(db, utils.FastScryptParams, lggr, cfg)
	if h, ok := cfg.(v2.HasEVMConfigs); ok {
		var ids []utils.Big
		for _, c := range h.EVMConfigs() {
			ids = append(ids, *c.ChainID)
		}
		if len(ids) > 0 {
			o := chainORM
			if o == nil {
				o = evm.NewORM(db, lggr, cfg)
			}
			if err = o.EnsureChains(ids); err != nil {
				t.Fatal(err)
			}
		}
	}
	mailMon := utils.NewMailboxMonitor(cfg.AppID().String())
	var chains chainlink.Chains
	chains.EVM, err = evm.LoadChainSet(testutils.Context(t), evm.ChainSetOpts{
		ORM:              chainORM,
		Config:           cfg,
		Logger:           lggr,
		DB:               db,
		KeyStore:         keyStore.Eth(),
		EventBroadcaster: eventBroadcaster,
		GenEthClient: func(_ *big.Int) evmclient.Client {
			if (ethClient.ChainID()).Cmp(cfg.DefaultChainID()) != 0 {
				t.Fatalf("expected eth client ChainID %d to match configured DefaultChainID %d", ethClient.ChainID(), cfg.DefaultChainID())
			}
			return ethClient
		},
		MailMon: mailMon,
	})
	if err != nil {
		lggr.Fatal(err)
	}
	if cfg.TerraEnabled() {
		terraLggr := lggr.Named("Terra")
		opts := terra.ChainSetOpts{
			Config:           cfg,
			Logger:           terraLggr,
			DB:               db,
			KeyStore:         keyStore.Terra(),
			EventBroadcaster: eventBroadcaster,
		}
		if newCfg, ok := cfg.(interface{ TerraConfigs() terra.TerraConfigs }); ok {
			cfgs := newCfg.TerraConfigs()
			opts.ORM = terra.NewORMImmut(cfgs)
			chains.Terra, err = terra.NewChainSetImmut(opts, cfgs)
			var ids []string
			for _, c := range cfgs {
				ids = append(ids, *c.ChainID)
			}
			if len(ids) > 0 {
				if err = terra.NewORM(db, terraLggr, cfg).EnsureChains(ids); err != nil {
					t.Fatal(err)
				}
			}
		} else {
			opts.ORM = terra.NewORM(db, terraLggr, cfg)
			chains.Terra, err = terra.NewChainSet(opts)
		}
		if err != nil {
			lggr.Fatal(err)
		}
	}
	if cfg.SolanaEnabled() {
		solLggr := lggr.Named("Solana")
		opts := solana.ChainSetOpts{
			Logger:   solLggr,
			DB:       db,
			KeyStore: keyStore.Solana(),
		}
		if newCfg, ok := cfg.(interface {
			SolanaConfigs() solana.SolanaConfigs
		}); ok {
			cfgs := newCfg.SolanaConfigs()
			opts.ORM = solana.NewORMImmut(cfgs)
			chains.Solana, err = solana.NewChainSetImmut(opts, cfgs)
			var ids []string
			for _, c := range cfgs {
				ids = append(ids, *c.ChainID)
			}
			if len(ids) > 0 {
				if err = solana.NewORM(db, solLggr, cfg).EnsureChains(ids); err != nil {
					t.Fatal(err)
				}
			}
		} else {
			opts.ORM = solana.NewORM(db, solLggr, cfg)
			chains.Solana, err = solana.NewChainSet(opts)
		}
		if err != nil {
			lggr.Fatal(err)
		}
	}
	if cfg.StarkNetEnabled() {
		starkLggr := lggr.Named("StarkNet")
		opts := starknet.ChainSetOpts{
			Config:   cfg,
			Logger:   starkLggr,
			KeyStore: keyStore.StarkNet(),
		}
		if newCfg, ok := cfg.(interface {
			StarknetConfigs() starknet.StarknetConfigs
		}); ok {
			cfgs := newCfg.StarknetConfigs()
			opts.ORM = starknet.NewORMImmut(cfgs)
			chains.StarkNet, err = starknet.NewChainSetImmut(opts, cfgs)
			var ids []string
			for _, c := range cfgs {
				ids = append(ids, *c.ChainID)
			}
			if len(ids) > 0 {
				if err = starknet.NewORM(db, starkLggr, cfg).EnsureChains(ids); err != nil {
					t.Fatal(err)
				}
			}
		} else {
			opts.ORM = starknet.NewORM(db, starkLggr, cfg)
			chains.StarkNet, err = starknet.NewChainSet(opts)
		}
		if err != nil {
			lggr.Fatal(err)
		}
	}
	c := clhttptest.NewTestLocalOnlyHTTPClient()
	appInstance, err := chainlink.NewApplication(chainlink.ApplicationOpts{
		Config:                   cfg,
		EventBroadcaster:         eventBroadcaster,
		MailMon:                  mailMon,
		SqlxDB:                   db,
		KeyStore:                 keyStore,
		Chains:                   chains,
		Logger:                   lggr,
		AuditLogger:              auditLogger,
		CloseLogger:              lggr.Sync,
		ExternalInitiatorManager: externalInitiatorManager,
		RestrictedHTTPClient:     c,
		UnrestrictedHTTPClient:   c,
		SecretGenerator:          MockSecretGenerator{},
	})
	require.NoError(t, err)
	app := appInstance.(*chainlink.ChainlinkApplication)
	ta := &TestApplication{
		t:                    t,
		ChainlinkApplication: app,
		Logger:               lggr,
	}
	ta.Server = httptest.NewServer(web.Router(t, app, nil))

	if !useRealExternalInitiatorManager {
		app.ExternalInitiatorManager = externalInitiatorManager
	}

	return ta
}

func NewEthMocksWithDefaultChain(t testing.TB) (c *evmMocks.Client) {
	testutils.SkipShortDB(t)
	c = NewEthMocks(t)
	c.On("ChainID").Return(&FixtureChainID).Maybe()
	return
}

func NewEthMocks(t testing.TB) *evmMocks.Client {
	return evmMocks.NewClient(t)
}

func NewEthMocksWithStartupAssertions(t testing.TB) *evmMocks.Client {
	testutils.SkipShort(t, "long test")
	c := NewEthMocks(t)
	c.On("Dial", mock.Anything).Maybe().Return(nil)
	c.On("SubscribeNewHead", mock.Anything, mock.Anything).Maybe().Return(EmptyMockSubscription(t), nil)
	c.On("SendTransaction", mock.Anything, mock.Anything).Maybe().Return(nil)
	c.On("HeadByNumber", mock.Anything, (*big.Int)(nil)).Maybe().Return(Head(0), nil)
	c.On("ChainID").Maybe().Return(&FixtureChainID)
	c.On("CallContract", mock.Anything, mock.Anything, mock.Anything).Maybe().Return([]byte{}, nil)
	c.On("SubscribeFilterLogs", mock.Anything, mock.Anything, mock.Anything).Maybe().Return(nil, errors.New("mocked"))
	c.On("CodeAt", mock.Anything, mock.Anything, mock.Anything).Maybe().Return([]byte{}, nil)
	c.On("Close").Maybe().Return()

	block := &types.Header{
		Number: big.NewInt(100),
	}
	c.On("HeaderByNumber", mock.Anything, mock.Anything).Maybe().Return(block, nil)

	return c
}

// NewEthMocksWithTransactionsOnBlocksAssertions sets an Eth mock with transactions on blocks
func NewEthMocksWithTransactionsOnBlocksAssertions(t testing.TB) *evmMocks.Client {
	testutils.SkipShort(t, "long test")
	c := NewEthMocks(t)
	c.On("Dial", mock.Anything).Maybe().Return(nil)
	c.On("SubscribeNewHead", mock.Anything, mock.Anything).Maybe().Return(EmptyMockSubscription(t), nil)
	c.On("SendTransaction", mock.Anything, mock.Anything).Maybe().Return(nil)
	c.On("HeadByNumber", mock.Anything, (*big.Int)(nil)).Maybe().Return(Head(2), nil)
	c.On("HeadByNumber", mock.Anything, big.NewInt(1)).Maybe().Return(Head(1), nil)
	c.On("HeadByNumber", mock.Anything, big.NewInt(0)).Maybe().Return(Head(0), nil)
	c.On("BatchCallContext", mock.Anything, mock.Anything).Maybe().Return(nil).Run(func(args mock.Arguments) {
		elems := args.Get(1).([]rpc.BatchElem)
		elems[0].Result = &evmtypes.Block{
			Number:       42,
			Hash:         utils.NewHash(),
			Transactions: LegacyTransactionsFromGasPrices(9001, 9002),
		}
		elems[1].Result = &evmtypes.Block{
			Number:       41,
			Hash:         utils.NewHash(),
			Transactions: LegacyTransactionsFromGasPrices(9003, 9004),
		}
	})
	c.On("ChainID").Maybe().Return(&FixtureChainID)
	c.On("Close").Maybe().Return()

	block := &types.Header{
		Number: big.NewInt(100),
	}
	c.On("HeaderByNumber", mock.Anything, mock.Anything).Maybe().Return(block, nil)

	return c
}

// Start starts the chainlink app and registers Stop to clean up at end of test.
func (ta *TestApplication) Start(ctx context.Context) error {
	ta.t.Helper()
	ta.Started = true
	err := ta.ChainlinkApplication.KeyStore.Unlock(Password)
	if err != nil {
		return err
	}

	err = ta.ChainlinkApplication.Start(ctx)
	if err != nil {
		return err
	}
	ta.t.Cleanup(func() { require.NoError(ta.t, ta.Stop()) })
	return nil
}

// Stop will stop the test application and perform cleanup
func (ta *TestApplication) Stop() error {
	ta.t.Helper()

	if !ta.Started {
		ta.t.Fatal("TestApplication Stop() called on an unstarted application")
	}

	// TODO: Here we double close, which is less than ideal.
	// We would prefer to invoke a method on an interface that
	// cleans up only in test.
	// FIXME: TestApplication probably needs to simply be removed
	err := ta.ChainlinkApplication.StopIfStarted()
	if ta.Server != nil {
		ta.Server.Close()
	}
	return err
}

func (ta *TestApplication) MustSeedNewSession(roleFixtureUserAPIEmail string) (id string) {
	session := NewSession()
	err := ta.GetSqlxDB().Get(&id, `INSERT INTO sessions (id, email, last_used, created_at) VALUES ($1, $2, $3, NOW()) RETURNING id`, session.ID, roleFixtureUserAPIEmail, session.LastUsed)
	require.NoError(ta.t, err)
	return id
}

// ImportKey adds private key to the application keystore and database
func (ta *TestApplication) Import(content string) {
	require.NoError(ta.t, ta.KeyStore.Unlock(Password))
	_, err := ta.KeyStore.Eth().Import([]byte(content), Password, &FixtureChainID)
	require.NoError(ta.t, err)
}

func (ta *TestApplication) NewHTTPClient(roleFixtureUserAPIEmail string) HTTPClientCleaner {
	ta.t.Helper()

	sessionID := ta.MustSeedNewSession(roleFixtureUserAPIEmail)

	return HTTPClientCleaner{
		HTTPClient: NewMockAuthenticatedHTTPClient(ta.Logger, ta.NewClientOpts(), sessionID),
		t:          ta.t,
	}
}

func (ta *TestApplication) NewClientOpts() cmd.ClientOpts {
	return cmd.ClientOpts{RemoteNodeURL: *MustParseURL(ta.t, ta.Server.URL), InsecureSkipVerify: true}
}

// NewClientAndRenderer creates a new cmd.Client for the test application
func (ta *TestApplication) NewClientAndRenderer() (*cmd.Client, *RendererMock) {
	sessionID := ta.MustSeedNewSession(APIEmailAdmin)
	r := &RendererMock{}
	lggr := logger.TestLogger(ta.t)
	client := &cmd.Client{
		Renderer:                       r,
		Config:                         ta.GetConfig(),
		Logger:                         lggr,
		AppFactory:                     seededAppFactory{ta.ChainlinkApplication},
		FallbackAPIInitializer:         NewMockAPIInitializer(ta.t),
		Runner:                         EmptyRunner{},
		HTTP:                           NewMockAuthenticatedHTTPClient(ta.Logger, ta.NewClientOpts(), sessionID),
		CookieAuthenticator:            MockCookieAuthenticator{t: ta.t},
		FileSessionRequestBuilder:      &MockSessionRequestBuilder{},
		PromptingSessionRequestBuilder: &MockSessionRequestBuilder{},
		ChangePasswordPrompter:         &MockChangePasswordPrompter{},
	}
	return client, r
}

func (ta *TestApplication) NewAuthenticatingClient(prompter cmd.Prompter) *cmd.Client {
	lggr := logger.TestLogger(ta.t)
	cookieAuth := cmd.NewSessionCookieAuthenticator(ta.NewClientOpts(), &cmd.MemoryCookieStore{}, lggr)
	client := &cmd.Client{
		Renderer:                       &RendererMock{},
		Config:                         ta.GetConfig(),
		Logger:                         lggr,
		AppFactory:                     seededAppFactory{ta.ChainlinkApplication},
		FallbackAPIInitializer:         NewMockAPIInitializer(ta.t),
		Runner:                         EmptyRunner{},
		HTTP:                           cmd.NewAuthenticatedHTTPClient(ta.Logger, ta.NewClientOpts(), cookieAuth, clsessions.SessionRequest{}),
		CookieAuthenticator:            cookieAuth,
		FileSessionRequestBuilder:      cmd.NewFileSessionRequestBuilder(lggr),
		PromptingSessionRequestBuilder: cmd.NewPromptingSessionRequestBuilder(prompter),
		ChangePasswordPrompter:         &MockChangePasswordPrompter{},
	}
	return client
}

// NewKeyStore returns a new, unlocked keystore
func NewKeyStore(t testing.TB, db *sqlx.DB, cfg pg.QConfig) keystore.Master {
	keystore := keystore.New(db, utils.FastScryptParams, logger.TestLogger(t), cfg)
	require.NoError(t, keystore.Unlock(Password))
	return keystore
}

func ParseJSON(t testing.TB, body io.Reader) models.JSON {
	t.Helper()

	b, err := io.ReadAll(body)
	require.NoError(t, err)
	return models.JSON{Result: gjson.ParseBytes(b)}
}

func ParseJSONAPIErrors(t testing.TB, body io.Reader) *models.JSONAPIErrors {
	t.Helper()

	b, err := io.ReadAll(body)
	require.NoError(t, err)
	var respJSON models.JSONAPIErrors
	err = json.Unmarshal(b, &respJSON)
	require.NoError(t, err)
	return &respJSON
}

// MustReadFile loads a file but should never fail
func MustReadFile(t testing.TB, file string) []byte {
	t.Helper()

	content, err := os.ReadFile(file)
	require.NoError(t, err)
	return content
}

type HTTPClientCleaner struct {
	HTTPClient cmd.HTTPClient
	t          testing.TB
}

func (r *HTTPClientCleaner) Get(path string, headers ...map[string]string) (*http.Response, func()) {
	resp, err := r.HTTPClient.Get(path, headers...)
	return bodyCleaner(r.t, resp, err)
}

func (r *HTTPClientCleaner) Post(path string, body io.Reader) (*http.Response, func()) {
	resp, err := r.HTTPClient.Post(path, body)
	return bodyCleaner(r.t, resp, err)
}

func (r *HTTPClientCleaner) Put(path string, body io.Reader) (*http.Response, func()) {
	resp, err := r.HTTPClient.Put(path, body)
	return bodyCleaner(r.t, resp, err)
}

func (r *HTTPClientCleaner) Patch(path string, body io.Reader, headers ...map[string]string) (*http.Response, func()) {
	resp, err := r.HTTPClient.Patch(path, body, headers...)
	return bodyCleaner(r.t, resp, err)
}

func (r *HTTPClientCleaner) Delete(path string) (*http.Response, func()) {
	resp, err := r.HTTPClient.Delete(path)
	return bodyCleaner(r.t, resp, err)
}

func bodyCleaner(t testing.TB, resp *http.Response, err error) (*http.Response, func()) {
	t.Helper()

	require.NoError(t, err)
	return resp, func() { require.NoError(t, resp.Body.Close()) }
}

// ParseResponseBody will parse the given response into a byte slice
func ParseResponseBody(t testing.TB, resp *http.Response) []byte {
	t.Helper()

	b, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	return b
}

// ParseJSONAPIResponse parses the response and returns the JSONAPI resource.
func ParseJSONAPIResponse(t testing.TB, resp *http.Response, resource interface{}) error {
	t.Helper()

	input := ParseResponseBody(t, resp)
	err := jsonapi.Unmarshal(input, resource)
	if err != nil {
		return fmt.Errorf("web: unable to unmarshal data, %+v", err)
	}

	return nil
}

// ParseJSONAPIResponseMeta parses the bytes of the root document and returns a
// map of *json.RawMessage's within the 'meta' key.
func ParseJSONAPIResponseMeta(input []byte) (map[string]*json.RawMessage, error) {
	var root map[string]*json.RawMessage
	err := json.Unmarshal(input, &root)
	if err != nil {
		return root, err
	}

	var meta map[string]*json.RawMessage
	err = json.Unmarshal(*root["meta"], &meta)
	return meta, err
}

// ParseJSONAPIResponseMetaCount parses the bytes of the root document and
// returns the value of the 'count' key from the 'meta' section.
func ParseJSONAPIResponseMetaCount(input []byte) (int, error) {
	meta, err := ParseJSONAPIResponseMeta(input)
	if err != nil {
		return -1, err
	}

	var metaCount int
	err = json.Unmarshal(*meta["count"], &metaCount)
	return metaCount, err
}

func CreateJobViaWeb(t testing.TB, app *TestApplication, request []byte) job.Job {
	t.Helper()

	client := app.NewHTTPClient(APIEmailAdmin)
	resp, cleanup := client.Post("/v2/jobs", bytes.NewBuffer(request))
	defer cleanup()
	AssertServerResponse(t, resp, http.StatusOK)

	var createdJob job.Job
	err := ParseJSONAPIResponse(t, resp, &createdJob)
	require.NoError(t, err)
	return createdJob
}

func CreateJobViaWeb2(t testing.TB, app *TestApplication, spec string) webpresenters.JobResource {
	t.Helper()

	client := app.NewHTTPClient(APIEmailAdmin)
	resp, cleanup := client.Post("/v2/jobs", bytes.NewBufferString(spec))
	defer cleanup()
	AssertServerResponse(t, resp, http.StatusOK)

	var jobResponse webpresenters.JobResource
	err := ParseJSONAPIResponse(t, resp, &jobResponse)
	require.NoError(t, err)
	return jobResponse
}

func DeleteJobViaWeb(t testing.TB, app *TestApplication, jobID int32) {
	t.Helper()

	client := app.NewHTTPClient(APIEmailAdmin)
	resp, cleanup := client.Delete(fmt.Sprintf("/v2/jobs/%v", jobID))
	defer cleanup()
	AssertServerResponse(t, resp, http.StatusNoContent)
}

func AwaitJobActive(t testing.TB, jobSpawner job.Spawner, jobID int32, waitFor time.Duration) {
	t.Helper()
	require.Eventually(t, func() bool {
		_, exists := jobSpawner.ActiveJobs()[jobID]
		return exists
	}, waitFor, 100*time.Millisecond)
}

func CreateJobRunViaExternalInitiatorV2(
	t testing.TB,
	app *TestApplication,
	jobID uuid.UUID,
	eia auth.Token,
	body string,
) webpresenters.PipelineRunResource {
	t.Helper()

	headers := make(map[string]string)
	headers[static.ExternalInitiatorAccessKeyHeader] = eia.AccessKey
	headers[static.ExternalInitiatorSecretHeader] = eia.Secret

	url := app.Server.URL + "/v2/jobs/" + jobID.String() + "/runs"
	bodyBuf := bytes.NewBufferString(body)
	resp, cleanup := UnauthenticatedPost(t, url, bodyBuf, headers)
	defer cleanup()
	AssertServerResponse(t, resp, 200)
	var pr webpresenters.PipelineRunResource
	err := ParseJSONAPIResponse(t, resp, &pr)
	require.NoError(t, err)

	// assert.Equal(t, j.ID, pr.JobSpecID)
	return pr
}

func CreateJobRunViaUser(
	t testing.TB,
	app *TestApplication,
	jobID uuid.UUID,
	body string,
) webpresenters.PipelineRunResource {
	t.Helper()

	bodyBuf := bytes.NewBufferString(body)
	client := app.NewHTTPClient(APIEmailAdmin)
	resp, cleanup := client.Post("/v2/jobs/"+jobID.String()+"/runs", bodyBuf)
	defer cleanup()
	AssertServerResponse(t, resp, 200)
	var pr webpresenters.PipelineRunResource
	err := ParseJSONAPIResponse(t, resp, &pr)
	require.NoError(t, err)

	return pr
}

// CreateExternalInitiatorViaWeb creates a bridgetype via web using /v2/bridge_types
func CreateExternalInitiatorViaWeb(
	t testing.TB,
	app *TestApplication,
	payload string,
) *webpresenters.ExternalInitiatorAuthentication {
	t.Helper()

	client := app.NewHTTPClient(APIEmailAdmin)
	resp, cleanup := client.Post(
		"/v2/external_initiators",
		bytes.NewBufferString(payload),
	)
	defer cleanup()
	AssertServerResponse(t, resp, http.StatusCreated)
	ei := &webpresenters.ExternalInitiatorAuthentication{}
	err := ParseJSONAPIResponse(t, resp, ei)
	require.NoError(t, err)

	return ei
}

const (
	// DBPollingInterval can't be too short to avoid DOSing the test database
	DBPollingInterval = 100 * time.Millisecond
	// AssertNoActionTimeout shouldn't be too long, or it will slow down tests
	AssertNoActionTimeout = 3 * time.Second
)

// WaitForSpecErrorV2 polls until the passed in jobID has count number
// of job spec errors.
func WaitForSpecErrorV2(t *testing.T, db *sqlx.DB, jobID int32, count int) []job.SpecError {
	t.Helper()

	g := gomega.NewWithT(t)
	var jse []job.SpecError
	g.Eventually(func() []job.SpecError {
		err := db.Select(&jse, `SELECT * FROM job_spec_errors WHERE job_id = $1`, jobID)
		assert.NoError(t, err)
		return jse
	}, testutils.WaitTimeout(t), DBPollingInterval).Should(gomega.HaveLen(count))
	return jse
}

func WaitForPipelineError(t testing.TB, nodeID int, jobID int32, expectedPipelineRuns int, expectedTaskRuns int, jo job.ORM, timeout, poll time.Duration) []pipeline.Run {
	t.Helper()
	return WaitForPipeline(t, nodeID, jobID, expectedPipelineRuns, expectedTaskRuns, jo, timeout, poll, pipeline.RunStatusErrored)
}
func WaitForPipelineComplete(t testing.TB, nodeID int, jobID int32, expectedPipelineRuns int, expectedTaskRuns int, jo job.ORM, timeout, poll time.Duration) []pipeline.Run {
	t.Helper()
	return WaitForPipeline(t, nodeID, jobID, expectedPipelineRuns, expectedTaskRuns, jo, timeout, poll, pipeline.RunStatusCompleted)
}

func WaitForPipeline(t testing.TB, nodeID int, jobID int32, expectedPipelineRuns int, expectedTaskRuns int, jo job.ORM, timeout, poll time.Duration, state pipeline.RunStatus) []pipeline.Run {
	t.Helper()

	var pr []pipeline.Run
	gomega.NewWithT(t).Eventually(func() bool {
		prs, _, err := jo.PipelineRuns(&jobID, 0, 1000)
		require.NoError(t, err)

		var matched []pipeline.Run
		for _, pr := range prs {
			if !pr.State.Finished() || pr.State != state {
				continue
			}

			// txdb effectively ignores transactionality of queries, so we need to explicitly expect a number of task runs
			// (if the read occurs mid-transaction and a job run is inserted but task runs not yet).
			if len(pr.PipelineTaskRuns) == expectedTaskRuns {
				matched = append(matched, pr)
			}
		}
		if len(matched) >= expectedPipelineRuns {
			pr = matched
			return true
		}
		return false
	}, timeout, poll).Should(
		gomega.BeTrue(),
		fmt.Sprintf(`expected at least %d runs with status "%s" on node %d for job %d, total runs %d`,
			expectedPipelineRuns,
			state,
			nodeID,
			jobID,
			len(pr),
		),
	)
	return pr
}

// AssertPipelineRunsStays asserts that the number of pipeline runs for a particular job remains at the provided values
func AssertPipelineRunsStays(t testing.TB, pipelineSpecID int32, db *sqlx.DB, want int) []pipeline.Run {
	t.Helper()
	g := gomega.NewWithT(t)

	var prs []pipeline.Run
	g.Consistently(func() []pipeline.Run {
		err := db.Select(&prs, `SELECT * FROM pipeline_runs WHERE pipeline_spec_id = $1`, pipelineSpecID)
		assert.NoError(t, err)
		return prs
	}, AssertNoActionTimeout, DBPollingInterval).Should(gomega.HaveLen(want))
	return prs
}

// AssertEthTxAttemptCountStays asserts that the number of tx attempts remains at the provided value
func AssertEthTxAttemptCountStays(t testing.TB, db *sqlx.DB, want int) []txmgr.EthTxAttempt {
	g := gomega.NewWithT(t)

	var txas []txmgr.EthTxAttempt
	var err error
	g.Consistently(func() []txmgr.EthTxAttempt {
		txas = make([]txmgr.EthTxAttempt, 0)
		err = db.Select(&txas, `SELECT * FROM eth_tx_attempts ORDER BY id ASC`)
		assert.NoError(t, err)
		return txas
	}, AssertNoActionTimeout, DBPollingInterval).Should(gomega.HaveLen(want))
	return txas
}

// Head given the value convert it into an Head
func Head(val interface{}) *evmtypes.Head {
	var h evmtypes.Head
	time := uint64(0)
	switch t := val.(type) {
	case int:
		h = evmtypes.NewHead(big.NewInt(int64(t)), utils.NewHash(), utils.NewHash(), time, utils.NewBig(&FixtureChainID))
	case uint64:
		h = evmtypes.NewHead(big.NewInt(int64(t)), utils.NewHash(), utils.NewHash(), time, utils.NewBig(&FixtureChainID))
	case int64:
		h = evmtypes.NewHead(big.NewInt(t), utils.NewHash(), utils.NewHash(), time, utils.NewBig(&FixtureChainID))
	case *big.Int:
		h = evmtypes.NewHead(t, utils.NewHash(), utils.NewHash(), time, utils.NewBig(&FixtureChainID))
	default:
		panic(fmt.Sprintf("Could not convert %v of type %T to Head", val, val))
	}
	return &h
}

// LegacyTransactionsFromGasPrices returns transactions matching the given gas prices
func LegacyTransactionsFromGasPrices(gasPrices ...int64) []evmtypes.Transaction {
	return LegacyTransactionsFromGasPricesTxType(0x0, gasPrices...)
}

func LegacyTransactionsFromGasPricesTxType(code evmtypes.TxType, gasPrices ...int64) []evmtypes.Transaction {
	txs := make([]evmtypes.Transaction, len(gasPrices))
	for i, gasPrice := range gasPrices {
		txs[i] = evmtypes.Transaction{Type: code, GasPrice: assets.NewWeiI(gasPrice), GasLimit: 42}
	}
	return txs
}

// DynamicFeeTransactionsFromTipCaps returns EIP-1559 transactions with the
// given TipCaps (FeeCap is arbitrary)
func DynamicFeeTransactionsFromTipCaps(tipCaps ...int64) []evmtypes.Transaction {
	return DynamicFeeTransactionsFromTipCapsTxType(0x02, tipCaps...)
}

func DynamicFeeTransactionsFromTipCapsTxType(code evmtypes.TxType, tipCaps ...int64) []evmtypes.Transaction {
	txs := make([]evmtypes.Transaction, len(tipCaps))
	for i, tipCap := range tipCaps {
		txs[i] = evmtypes.Transaction{Type: code, MaxPriorityFeePerGas: assets.NewWeiI(tipCap), GasLimit: 42, MaxFeePerGas: assets.GWei(5000)}
	}
	return txs
}

type TransactionReceipter interface {
	TransactionReceipt(context.Context, common.Hash) (*types.Receipt, error)
}

func RequireTxSuccessful(t testing.TB, client TransactionReceipter, txHash common.Hash) *types.Receipt {
	t.Helper()
	r, err := client.TransactionReceipt(context.Background(), txHash)
	require.NoError(t, err)
	require.NotNil(t, r)
	require.Equal(t, uint64(1), r.Status)
	return r
}

// AssertServerResponse is used to match against a client response, will print
// any errors returned if the request fails.
func AssertServerResponse(t testing.TB, resp *http.Response, expectedStatusCode int) {
	t.Helper()

	if resp.StatusCode == expectedStatusCode {
		return
	}

	t.Logf("expected status code %s got %s", http.StatusText(expectedStatusCode), http.StatusText(resp.StatusCode))

	if resp.StatusCode >= 300 && resp.StatusCode < 600 {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			assert.FailNowf(t, "Unable to read body", err.Error())
		}

		var result *models.JSONAPIErrors
		err = json.Unmarshal(b, &result)
		if err != nil {
			assert.FailNowf(t, fmt.Sprintf("Unable to unmarshal json from body '%s'", string(b)), err.Error())
		}

		assert.FailNowf(t, "Request failed", "Expected %d response, got %d with errors: %s", expectedStatusCode, resp.StatusCode, result.Errors)
	} else {
		assert.FailNowf(t, "Unexpected response", "Expected %d response, got %d", expectedStatusCode, resp.StatusCode)
	}
}

func DecodeSessionCookie(value string) (string, error) {
	var decrypted map[interface{}]interface{}
	codecs := securecookie.CodecsFromPairs([]byte(SessionSecret))
	err := securecookie.DecodeMulti(webauth.SessionName, value, &decrypted, codecs...)
	if err != nil {
		return "", err
	}
	value, ok := decrypted[webauth.SessionIDKey].(string)
	if !ok {
		return "", fmt.Errorf("decrypted[web.SessionIDKey] is not a string (%v)", value)
	}
	return value, nil
}

func MustGenerateSessionCookie(t testing.TB, value string) *http.Cookie {
	decrypted := map[interface{}]interface{}{webauth.SessionIDKey: value}
	codecs := securecookie.CodecsFromPairs([]byte(SessionSecret))
	encoded, err := securecookie.EncodeMulti(webauth.SessionName, decrypted, codecs...)
	if err != nil {
		logger.TestLogger(t).Panic(err)
	}
	return sessions.NewCookie(webauth.SessionName, encoded, &sessions.Options{})
}

func AssertError(t testing.TB, want bool, err error) {
	t.Helper()

	if want {
		assert.Error(t, err)
	} else {
		assert.NoError(t, err)
	}
}

func UnauthenticatedPost(t testing.TB, url string, body io.Reader, headers map[string]string) (*http.Response, func()) {
	t.Helper()
	return unauthenticatedHTTP(t, "POST", url, body, headers)
}

func UnauthenticatedGet(t testing.TB, url string, headers map[string]string) (*http.Response, func()) {
	t.Helper()
	return unauthenticatedHTTP(t, "GET", url, nil, headers)
}

func unauthenticatedHTTP(t testing.TB, method string, url string, body io.Reader, headers map[string]string) (*http.Response, func()) {
	t.Helper()

	client := clhttptest.NewTestLocalOnlyHTTPClient()
	request, err := http.NewRequest(method, url, body)
	require.NoError(t, err)
	request.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		request.Header.Add(key, value)
	}
	resp, err := client.Do(request)
	require.NoError(t, err)
	return resp, func() { resp.Body.Close() }
}

func MustParseDuration(t testing.TB, durationStr string) time.Duration {
	t.Helper()

	duration, err := time.ParseDuration(durationStr)
	require.NoError(t, err)
	return duration
}

func NewSession(optionalSessionID ...string) clsessions.Session {
	session := clsessions.NewSession()
	if len(optionalSessionID) > 0 {
		session.ID = optionalSessionID[0]
	}
	return session
}

func AllExternalInitiators(t testing.TB, db *sqlx.DB) []bridges.ExternalInitiator {
	t.Helper()

	var all []bridges.ExternalInitiator
	err := db.Select(&all, `SELECT * FROM external_initiators`)
	require.NoError(t, err)
	return all
}

type Awaiter chan struct{}

func NewAwaiter() Awaiter { return make(Awaiter) }

func (a Awaiter) ItHappened() { close(a) }

func (a Awaiter) AssertHappened(t *testing.T, expected bool) {
	t.Helper()
	select {
	case <-a:
		if !expected {
			t.Fatal("It happened")
		}
	default:
		if expected {
			t.Fatal("It didn't happen")
		}
	}
}

func (a Awaiter) AwaitOrFail(t testing.TB, durationParams ...time.Duration) {
	t.Helper()

	duration := 10 * time.Second
	if len(durationParams) > 0 {
		duration = durationParams[0]
	}

	select {
	case <-a:
	case <-time.After(duration):
		t.Fatal("Timed out waiting for Awaiter to get ItHappened")
	}
}

func CallbackOrTimeout(t testing.TB, msg string, callback func(), durationParams ...time.Duration) {
	t.Helper()

	duration := 100 * time.Millisecond
	if len(durationParams) > 0 {
		duration = durationParams[0]
	}

	done := make(chan struct{})
	go func() {
		callback()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(duration):
		t.Fatalf("CallbackOrTimeout: %s timed out", msg)
	}
}

func MustParseURL(t testing.TB, input string) *url.URL {
	return testutils.MustParseURL(t, input)
}

// EthereumLogIterator is the interface provided by gethwrapper representations of EVM
// logs.
type EthereumLogIterator interface{ Next() bool }

// GetLogs drains logs of EVM log representations. Since those log
// representations don't fit into a type hierarchy, this API is a bit awkward.
// It returns the logs as a slice of blank interface{}s, and if rv is non-nil,
// it must be a pointer to a slice for elements of the same type as the logs,
// in which case GetLogs will append the logs to it.
func GetLogs(t *testing.T, rv interface{}, logs EthereumLogIterator) []interface{} {
	v := reflect.ValueOf(rv)
	require.True(t, rv == nil ||
		v.Kind() == reflect.Ptr && v.Elem().Kind() == reflect.Slice,
		"must pass a slice to receive logs")
	var e reflect.Value
	if rv != nil {
		e = v.Elem()
	}
	var irv []interface{}
	for logs.Next() {
		log := reflect.Indirect(reflect.ValueOf(logs)).FieldByName("Event")
		if v.Kind() == reflect.Ptr {
			e.Set(reflect.Append(e, log))
		}
		irv = append(irv, log.Interface())
	}
	return irv
}

func MakeConfigDigest(t *testing.T) ocrtypes.ConfigDigest {
	t.Helper()
	b := make([]byte, 16)
	/* #nosec G404 */
	_, err := rand.Read(b)
	if err != nil {
		t.Fatal(err)
	}
	return MustBytesToConfigDigest(t, b)
}

func MustBytesToConfigDigest(t *testing.T, b []byte) ocrtypes.ConfigDigest {
	t.Helper()
	configDigest, err := ocrtypes.BytesToConfigDigest(b)
	if err != nil {
		t.Fatal(err)
	}
	return configDigest
}

// MockApplicationEthCalls mocks all calls made by the chainlink application as
// standard when starting and stopping
func MockApplicationEthCalls(t *testing.T, app *TestApplication, ethClient *evmMocks.Client, sub *evmMocks.Subscription) {
	t.Helper()

	// Start
	ethClient.On("Dial", mock.Anything).Return(nil)
	ethClient.On("SubscribeNewHead", mock.Anything, mock.Anything).Return(sub, nil).Maybe()
	ethClient.On("ChainID", mock.Anything).Return(app.GetConfig().DefaultChainID(), nil)
	ethClient.On("PendingNonceAt", mock.Anything, mock.Anything).Return(uint64(0), nil).Maybe()
	ethClient.On("HeadByNumber", mock.Anything, mock.Anything).Return(nil, nil).Maybe()
	ethClient.On("Close").Return().Maybe()
}

func BatchElemMatchesParams(req rpc.BatchElem, arg interface{}, method string) bool {
	return req.Method == method &&
		len(req.Args) == 1 && req.Args[0] == arg
}

func BatchElemMustMatchParams(t *testing.T, req rpc.BatchElem, hash common.Hash, method string) {
	t.Helper()
	if !BatchElemMatchesParams(req, hash, method) {
		t.Fatalf("Batch hash %v does not match expected %v", req.Args[0], hash)
	}
}

type SimulateIncomingHeadsArgs struct {
	StartBlock, EndBlock int64
	HeadTrackables       []httypes.HeadTrackable
	Blocks               *Blocks
}

// SimulateIncomingHeads spawns a goroutine which sends a stream of heads and closes the returned channel when finished.
func SimulateIncomingHeads(t *testing.T, args SimulateIncomingHeadsArgs) (done chan struct{}) {
	t.Helper()
	lggr := logger.TestLogger(t).Named("SimulateIncomingHeads")
	lggr.Infof("Simulating incoming heads from %v to %v...", args.StartBlock, args.EndBlock)

	if args.EndBlock > args.StartBlock {
		if l := 1 + args.EndBlock - args.StartBlock; l > int64(len(args.Blocks.Heads)) {
			t.Fatalf("invalid configuration: too few blocks %d for range length %d", len(args.Blocks.Heads), l)
		}
	}

	// Build the full chain of heads
	heads := args.Blocks.Heads
	ctx, cancel := context.WithTimeout(context.Background(), testutils.WaitTimeout(t))
	t.Cleanup(cancel)
	done = make(chan struct{})
	go func(t *testing.T) {
		defer close(done)
		ticker := time.NewTicker(250 * time.Millisecond)
		defer ticker.Stop()

		for current := args.StartBlock; ; current++ {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				_, exists := heads[current]
				if !exists {
					lggr.Infof("Out of heads: %d does not exist", current)
					return
				}

				lggr.Infof("Sending head: %d", current)
				for _, ht := range args.HeadTrackables {
					ht.OnNewLongestChain(ctx, heads[current])
				}
				if args.EndBlock >= 0 && current == args.EndBlock {
					return
				}
			}
		}
	}(t)
	return done
}

// Blocks - a helper logic to construct a range of linked heads
// and an ability to fork and create logs from them
type Blocks struct {
	t       *testing.T
	Hashes  []common.Hash
	mHashes map[int64]common.Hash
	Heads   map[int64]*evmtypes.Head
}

func (b *Blocks) LogOnBlockNum(i uint64, addr common.Address) types.Log {
	return RawNewRoundLog(b.t, addr, b.Hashes[i], i, 0, false)
}

func (b *Blocks) LogOnBlockNumRemoved(i uint64, addr common.Address) types.Log {
	return RawNewRoundLog(b.t, addr, b.Hashes[i], i, 0, true)
}

func (b *Blocks) LogOnBlockNumWithIndex(i uint64, logIndex uint, addr common.Address) types.Log {
	return RawNewRoundLog(b.t, addr, b.Hashes[i], i, logIndex, false)
}

func (b *Blocks) LogOnBlockNumWithIndexRemoved(i uint64, logIndex uint, addr common.Address) types.Log {
	return RawNewRoundLog(b.t, addr, b.Hashes[i], i, logIndex, true)
}

func (b *Blocks) LogOnBlockNumWithTopics(i uint64, logIndex uint, addr common.Address, topics []common.Hash) types.Log {
	return RawNewRoundLogWithTopics(b.t, addr, b.Hashes[i], i, logIndex, false, topics)
}

func (b *Blocks) HashesMap() map[int64]common.Hash {
	return b.mHashes
}

func (b *Blocks) Head(number uint64) *evmtypes.Head {
	return b.Heads[int64(number)]
}

func (b *Blocks) ForkAt(t *testing.T, blockNum int64, numHashes int) *Blocks {
	forked := NewBlocks(t, len(b.Heads)+numHashes)
	if _, exists := forked.Heads[blockNum]; !exists {
		t.Fatalf("Not enough length for block num: %v", blockNum)
	}

	for i := int64(0); i < blockNum; i++ {
		forked.Heads[i] = b.Heads[i]
	}

	forked.Heads[blockNum].ParentHash = b.Heads[blockNum].ParentHash
	forked.Heads[blockNum].Parent = b.Heads[blockNum].Parent
	return forked
}

func (b *Blocks) NewHead(number uint64) *evmtypes.Head {
	parentNumber := number - 1
	parent, ok := b.Heads[int64(parentNumber)]
	if !ok {
		b.t.Fatalf("Can't find parent block at index: %v", parentNumber)
	}
	head := &evmtypes.Head{
		Number:     parent.Number + 1,
		Hash:       utils.NewHash(),
		ParentHash: parent.Hash,
		Parent:     parent,
		Timestamp:  time.Unix(parent.Number+1, 0),
		EVMChainID: utils.NewBig(&FixtureChainID),
	}
	return head
}

func NewBlocks(t *testing.T, numHashes int) *Blocks {
	hashes := make([]common.Hash, 0)
	heads := make(map[int64]*evmtypes.Head)
	for i := int64(0); i < int64(numHashes); i++ {
		hash := utils.NewHash()
		hashes = append(hashes, hash)

		heads[i] = &evmtypes.Head{Hash: hash, Number: i, Timestamp: time.Unix(i, 0), EVMChainID: utils.NewBig(&FixtureChainID)}
		if i > 0 {
			parent := heads[i-1]
			heads[i].Parent = parent
			heads[i].ParentHash = parent.Hash
		}
	}

	hashesMap := make(map[int64]common.Hash)
	for i := 0; i < len(hashes); i++ {
		hashesMap[int64(i)] = hashes[i]
	}

	return &Blocks{
		t:       t,
		Hashes:  hashes,
		mHashes: hashesMap,
		Heads:   heads,
	}
}

// HeadBuffer - stores heads in sequence, with increasing timestamps
type HeadBuffer struct {
	t     *testing.T
	Heads []*evmtypes.Head
}

func NewHeadBuffer(t *testing.T) *HeadBuffer {
	return &HeadBuffer{
		t:     t,
		Heads: make([]*evmtypes.Head, 0),
	}
}

func (hb *HeadBuffer) Append(head *evmtypes.Head) {
	cloned := &evmtypes.Head{
		Number:     head.Number,
		Hash:       head.Hash,
		ParentHash: head.ParentHash,
		Parent:     head.Parent,
		Timestamp:  time.Unix(int64(len(hb.Heads)), 0),
		EVMChainID: head.EVMChainID,
	}
	hb.Heads = append(hb.Heads, cloned)
}

type HeadTrackableFunc func(context.Context, *evmtypes.Head)

func (fn HeadTrackableFunc) OnNewLongestChain(ctx context.Context, head *evmtypes.Head) {
	fn(ctx, head)
}

type testifyExpectationsAsserter interface {
	AssertExpectations(t mock.TestingT) bool
}

type fakeT struct{}

func (ft fakeT) Logf(format string, args ...interface{})   {}
func (ft fakeT) Errorf(format string, args ...interface{}) {}
func (ft fakeT) FailNow()                                  {}

func EventuallyExpectationsMet(t *testing.T, mock testifyExpectationsAsserter, timeout time.Duration, interval time.Duration) {
	t.Helper()

	chTimeout := time.After(timeout)
	for {
		var ft fakeT
		success := mock.AssertExpectations(ft)
		if success {
			return
		}
		select {
		case <-chTimeout:
			mock.AssertExpectations(t)
			t.FailNow()
		default:
			time.Sleep(interval)
		}
	}
}

func AssertCount(t *testing.T, db *sqlx.DB, tableName string, expected int64) {
	testutils.AssertCount(t, db, tableName, expected)
}

func WaitForCount(t *testing.T, db *sqlx.DB, tableName string, want int64) {
	t.Helper()
	g := gomega.NewWithT(t)
	var count int64
	var err error
	g.Eventually(func() int64 {
		err = db.Get(&count, fmt.Sprintf(`SELECT count(*) FROM %s;`, tableName))
		assert.NoError(t, err)
		return count
	}, testutils.WaitTimeout(t), DBPollingInterval).Should(gomega.Equal(want))
}

func AssertCountStays(t testing.TB, db *sqlx.DB, tableName string, want int64) {
	t.Helper()
	g := gomega.NewWithT(t)
	var count int64
	var err error
	g.Consistently(func() int64 {
		err = db.Get(&count, fmt.Sprintf(`SELECT count(*) FROM %q`, tableName))
		assert.NoError(t, err)
		return count
	}, AssertNoActionTimeout, DBPollingInterval).Should(gomega.Equal(want))
}

func AssertRecordEventually(t *testing.T, db *sqlx.DB, model interface{}, stmt string, check func() bool) {
	t.Helper()
	g := gomega.NewWithT(t)
	g.Eventually(func() bool {
		err := db.Get(model, stmt)
		require.NoError(t, err, "unable to find record in DB")
		return check()
	}, testutils.WaitTimeout(t), DBPollingInterval).Should(gomega.BeTrue())
}

func MustSendingKeyStates(t *testing.T, ethKeyStore keystore.Eth, chainID *big.Int) []ethkey.State {
	keys, err := ethKeyStore.EnabledKeysForChain(chainID)
	require.NoError(t, err)
	states, err := ethKeyStore.GetStatesForKeys(keys)
	require.NoError(t, err)
	return states
}

func MustRandomP2PPeerID(t *testing.T) p2ppeer.ID {
	reader := rand.New(source)
	p2pPrivkey, _, err := cryptop2p.GenerateEd25519Key(reader)
	require.NoError(t, err)
	id, err := p2ppeer.IDFromPrivateKey(p2pPrivkey)
	require.NoError(t, err)
	return id
}

func MustWebURL(t *testing.T, s string) *models.WebURL {
	uri, err := url.Parse(s)
	require.NoError(t, err)
	return (*models.WebURL)(uri)
}

func AssertPipelineTaskRunsSuccessful(t testing.TB, runs []pipeline.TaskRun) {
	t.Helper()
	for i, run := range runs {
		require.True(t, run.Error.IsZero(), fmt.Sprintf("pipeline.Task run failed (idx: %v, dotID: %v, error: '%v')", i, run.GetDotID(), run.Error.ValueOrZero()))
	}
}

func AssertPipelineTaskRunsErrored(t testing.TB, runs []pipeline.TaskRun) {
	t.Helper()
	for i, run := range runs {
		require.False(t, run.Error.IsZero(), fmt.Sprintf("expected pipeline.Task run to have failed, but it succeeded (idx: %v, dotID: %v, output: '%v')", i, run.GetDotID(), run.Output))
	}
}

func NewTestChainScopedConfig(t testing.TB) evmconfig.ChainScopedConfig {
	cfg := configtest2.NewGeneralConfig(t, nil)
	return evmtest.NewChainScopedConfig(t, cfg)
}

func MustGetStateForKey(t testing.TB, kst keystore.Eth, key ethkey.KeyV2) ethkey.State {
	states, err := kst.GetStatesForKeys([]ethkey.KeyV2{key})
	require.NoError(t, err)
	return states[0]
}

func NewTxmORM(t *testing.T, db *sqlx.DB, cfg pg.QConfig) txmgr.ORM {
	return txmgr.NewORM(db, logger.TestLogger(t), cfg)
}

// ClearDBTables deletes all rows from the given tables
func ClearDBTables(t *testing.T, db *sqlx.DB, tables ...string) {
	tx, err := db.Beginx()
	require.NoError(t, err)

	for _, table := range tables {
		_, err = tx.Exec(fmt.Sprintf("DELETE FROM %s", table))
		require.NoError(t, err)
	}

	err = tx.Commit()
	require.NoError(t, err)
}

// CopyFlagSetFromAction takes the given flagSet and applies the actions flag to it.
// `parentCommand` will filter the app commands and only applies the flags if the command/subcommand has a parent with that name, if left empty no filtering is done
func CopyFlagSetFromAction(action interface{}, flagSet *flag.FlagSet, parentCommand string) {
	cliApp := cmd.Client{}
	app := cmd.NewApp(&cliApp)

	foundName := parentCommand == ""

	for _, command := range app.Commands {
		flags := recursiveFindFlagsWithName(action, command, parentCommand, foundName)

		if flags != nil {
			for _, flag := range flags {
				flag.Apply(flagSet)
			}
		}
	}

}

func recursiveFindFlagsWithName(action interface{}, command cli.Command, parent string, foundName bool) []cli.Flag {

	if command.Action != nil {
		if getFunctionName(action) == getFunctionName(command.Action) && foundName {
			return command.Flags
		}
	}

	for _, subcommand := range command.Subcommands {
		if !foundName {
			foundName = strings.ToLower(subcommand.Name) == strings.ToLower(parent)
		}

		found := recursiveFindFlagsWithName(action, subcommand, parent, foundName)
		if found != nil {
			return found
		}
	}
	return nil
}

func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
