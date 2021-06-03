package cltest

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/services/gasupdater"
	httypes "github.com/smartcontractkit/chainlink/core/services/headtracker/types"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/static"
	"github.com/smartcontractkit/chainlink/core/store/dialects"

	p2ppeer "github.com/libp2p/go-libp2p-core/peer"
	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/auth"
	"github.com/smartcontractkit/chainlink/core/cmd"
	"github.com/smartcontractkit/chainlink/core/gracefulpanic"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/flux_aggregator_wrapper"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	strpkg "github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/smartcontractkit/chainlink/core/store/presenters"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/web"
	webpresenters "github.com/smartcontractkit/chainlink/core/web/presenters"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting/types"

	"github.com/DATA-DOG/go-txdb"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/packr"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
	"github.com/manyminds/api2go/jsonapi"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
	"go.uber.org/zap/zapcore"
	null "gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
)

const (
	// RootDir the root directory for cltest
	RootDir = "/tmp/chainlink_test"
	// APIKey of the fixture API user
	APIKey = "2d25e62eaf9143e993acaf48691564b2"
	// APISecret of the fixture API user.
	APISecret = "1eCP/w0llVkchejFaoBpfIGaLRxZK54lTXBCT22YLW+pdzE4Fafy/XO5LoJ2uwHi"
	// APIEmail is the email of the fixture API user
	APIEmail = "apiuser@chainlink.test"
	// Password just a password we use everywhere for testing
	Password = "p4SsW0rD1!@#_"
	// SessionSecret is the hardcoded secret solely used for test
	SessionSecret = "clsession_test_secret"
	// DefaultKeyAddress is the ETH address of the fixture key
	DefaultKeyAddress = "0xF67D0290337bca0847005C7ffD1BC75BA9AAE6e4"
	// DefaultKeyFixtureFileName is the filename of the fixture key
	DefaultKeyFixtureFileName = "testkey-0xF67D0290337bca0847005C7ffD1BC75BA9AAE6e4.json"
	// DefaultKeyJSON is the JSON for the default key encrypted with fast scrypt and password 'password' (used for fixture file)
	DefaultKeyJSON = `{"address":"F67D0290337bca0847005C7ffD1BC75BA9AAE6e4","crypto":{"cipher":"aes-128-ctr","ciphertext":"9c3565050ba4e10ea388bcd17d77c141441ce1be5db339f0201b9ed733d780c6","cipherparams":{"iv":"f968fc947495646ee8b5dbaadb242ec0"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"33ad88742a983dfeb8adcc9a39fdde4cb47f7e23ea2ef80b35723d940959e3fd"},"mac":"b3747959cbbb9b26f861ab82d69154b4ec8108bbac017c1341f6fd3295beceaf"},"id":"8c79a654-96b1-45d9-8978-3efa07578011","version":3}`
	// AllowUnstarted enable an application that can be used in tests without being started
	AllowUnstarted = "allow_unstarted"
	// DefaultPeerID is the peer ID of the fixture p2p key
	DefaultPeerID = "12D3KooWApUJaQB2saFjyEUfq6BmysnsSnhLnY5CF9tURYVKgoXK"
	// A peer ID without an associated p2p key.
	NonExistentPeerID = "12D3KooWAdCzaesXyezatDzgGvCngqsBqoUqnV9PnVc46jsVt2i9"
	// DefaultOCRKeyBundleID is the ID of the fixture ocr key bundle
	DefaultOCRKeyBundleID = "7f993fb701b3410b1f6e8d4d93a7462754d24609b9b31a4fe64a0cb475a4d934"
)

var (
	DefaultP2PPeerID     models.PeerID
	NonExistentP2PPeerID models.PeerID
	// DefaultOCRKeyBundleIDSha256 is the ID of the fixture ocr key bundle
	DefaultOCRKeyBundleIDSha256 models.Sha256Hash
	FluxAggAddress              = common.HexToAddress("0x3cCad4715152693fE3BC4460591e3D3Fbd071b42")
	storeCounter                uint64
	minimumContractPayment      = assets.NewLink(100)
)

func init() {
	gin.SetMode(gin.TestMode)
	gomega.SetDefaultEventuallyTimeout(3 * time.Second)
	lvl := logLevelFromEnv()
	logger.SetLogger(logger.CreateTestLogger(lvl))
	// Register txdb as dialect wrapping postgres
	// See: DialectTransactionWrappedPostgres
	config := orm.NewConfig()

	parsed := config.DatabaseURL()
	if parsed.Path == "" {
		msg := fmt.Sprintf("invalid DATABASE_URL: `%s`. You must set DATABASE_URL env var to point to your test database. Note that the test database MUST end in `_test` to differentiate from a possible production DB. HINT: Try DATABASE_URL=postgresql://postgres@localhost:5432/chainlink_test?sslmode=disable", parsed.String())
		panic(msg)
	}
	if !strings.HasSuffix(parsed.Path, "_test") {
		msg := fmt.Sprintf("cannot run tests against database named `%s`. Note that the test database MUST end in `_test` to differentiate from a possible production DB. HINT: Try DATABASE_URL=postgresql://postgres@localhost:5432/chainlink_test?sslmode=disable", parsed.Path[1:])
		panic(msg)
	}
	// Disable SavePoints because they cause random errors for reasons I cannot fathom
	// Perhaps txdb's built-in transaction emulation is broken in some subtle way?
	// NOTE: That this will cause transaction BEGIN/ROLLBACK to effectively be
	// a no-op, this should have no negative impact on normal test operation.
	// If you MUST test BEGIN/ROLLBACK behaviour, you will have to configure your
	// store to use the raw DialectPostgres dialect and setup a one-use database.
	// See BootstrapThrowawayORM() as a convenience function to help you do this.
	txdb.Register(string(dialects.TransactionWrappedPostgres), string(dialects.Postgres), parsed.String(), txdb.SavePointOption(nil))

	// Seed the random number generator, otherwise separate modules will take
	// the same advisory locks when tested with `go test -p N` for N > 1
	seed := time.Now().UTC().UnixNano()
	logger.Debugf("Using seed: %v", seed)
	rand.Seed(seed)

	defaultP2PPeerID, err := p2ppeer.Decode(DefaultPeerID)
	if err != nil {
		panic(err)
	}
	DefaultP2PPeerID = models.PeerID(defaultP2PPeerID)
	nonExistentP2PPeerID, err := p2ppeer.Decode(NonExistentPeerID)
	if err != nil {
		panic(err)
	}
	NonExistentP2PPeerID = models.PeerID(nonExistentP2PPeerID)
	DefaultOCRKeyBundleIDSha256, err = models.Sha256HashFromHex(DefaultOCRKeyBundleID)
	if err != nil {
		panic(err)
	}
}

func logLevelFromEnv() zapcore.Level {
	lvl := zapcore.ErrorLevel
	if env := os.Getenv("LOG_LEVEL"); env != "" {
		_ = lvl.Set(env)
	}
	return lvl
}

// TestConfig struct with test store and wsServer
type TestConfig struct {
	t testing.TB
	*orm.Config
	wsServer *httptest.Server
}

// NewConfig returns a new TestConfig
func NewConfig(t testing.TB) (*TestConfig, func()) {
	t.Helper()

	wsserver, url, cleanup := newWSServer()
	config := NewConfigWithWSServer(t, url, wsserver)
	// Tests almost always want to request to localhost so its easier to set this here
	config.Set("DEFAULT_HTTP_ALLOW_UNRESTRICTED_NETWORK_ACCESS", true)
	// Disable gas updater for application tests
	config.Set("GAS_UPDATER_ENABLED", false)
	// Disable tx re-sending for application tests
	config.Set("ETH_TX_RESEND_AFTER_THRESHOLD", 0)
	// Limit ETH_FINALITY_DEPTH to avoid useless extra work backfilling heads
	config.Set("ETH_FINALITY_DEPTH", 1)
	// Disable the EthTxReaper
	config.Set("ETH_TX_REAPER_THRESHOLD", 0)
	// Set low sampling interval to remain within test head waiting timeouts
	config.Set("ETH_HEAD_TRACKER_SAMPLING_INTERVAL", "100ms")
	return config, cleanup
}

func NewRandomInt64() int64 {
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

func MustJobIDFromString(t *testing.T, s string) models.JobID {
	t.Helper()
	id, err := models.NewJobIDFromString(s)
	require.NoError(t, err)
	return id
}

// NewTestConfig returns a test configuration
func NewTestConfig(t testing.TB, options ...interface{}) *TestConfig {
	t.Helper()

	count := atomic.AddUint64(&storeCounter, 1)
	rootdir := filepath.Join(RootDir, fmt.Sprintf("%d-%d", time.Now().UnixNano(), count))
	rawConfig := orm.NewConfig()

	rawConfig.Dialect = dialects.TransactionWrappedPostgres
	for _, opt := range options {
		switch v := opt.(type) {
		case dialects.DialectName:
			rawConfig.Dialect = v
		}
	}

	// Unique advisory lock is required otherwise all tests will block each other
	rawConfig.AdvisoryLockID = NewRandomInt64()

	rawConfig.Set("BRIDGE_RESPONSE_URL", "http://localhost:6688")
	rawConfig.Set("ETH_CHAIN_ID", eth.NullClientChainID)
	rawConfig.Set("CHAINLINK_DEV", true)
	rawConfig.Set("ETH_GAS_BUMP_THRESHOLD", 3)
	rawConfig.Set("MIGRATE_DATABASE", false)
	rawConfig.Set("MINIMUM_SERVICE_DURATION", "24h")
	rawConfig.Set("MIN_INCOMING_CONFIRMATIONS", 1)
	rawConfig.Set("MIN_OUTGOING_CONFIRMATIONS", 6)
	rawConfig.Set("MINIMUM_CONTRACT_PAYMENT", minimumContractPayment.Text(10))
	rawConfig.Set("ROOT", rootdir)
	rawConfig.Set("SESSION_TIMEOUT", "2m")
	rawConfig.Set("INSECURE_FAST_SCRYPT", "true")
	rawConfig.Set("BALANCE_MONITOR_ENABLED", "false")
	rawConfig.Set("P2P_LISTEN_PORT", "12345")
	rawConfig.Set("P2P_PEER_ID", DefaultP2PPeerID.String())
	rawConfig.Set("DATABASE_TIMEOUT", "5s")
	rawConfig.Set("GLOBAL_LOCK_RETRY_INTERVAL", "10ms")
	rawConfig.Set("ORM_MAX_OPEN_CONNS", "5")
	rawConfig.Set("ORM_MAX_IDLE_CONNS", "2")
	rawConfig.Set("ETH_TX_REAPER_THRESHOLD", 0)
	rawConfig.SecretGenerator = mockSecretGenerator{}
	config := TestConfig{t: t, Config: rawConfig}
	return &config
}

// NewConfigWithWSServer return new config with specified wsserver
func NewConfigWithWSServer(t testing.TB, url string, wsserver *httptest.Server) *TestConfig {
	t.Helper()

	config := NewTestConfig(t)
	config.Set("ETH_URL", url)
	config.wsServer = wsserver
	return config
}

type JobPipelineV2TestHelper struct {
	Prm pipeline.ORM
	Eb  postgres.EventBroadcaster
	Jrm job.ORM
	Pr  pipeline.Runner
}

func NewJobPipelineV2(t testing.TB, db *gorm.DB) JobPipelineV2TestHelper {
	config, cleanupCfg := NewConfig(t)
	t.Cleanup(cleanupCfg)
	prm, eb, cleanup := NewPipelineORM(t, config, db)
	jrm := job.NewORM(db, config.Config, prm, eb, &postgres.NullAdvisoryLocker{})
	t.Cleanup(cleanup)
	pr := pipeline.NewRunner(prm, config.Config)
	return JobPipelineV2TestHelper{
		prm,
		eb,
		jrm,
		pr,
	}
}

func NewPipelineORM(t testing.TB, config *TestConfig, db *gorm.DB) (pipeline.ORM, postgres.EventBroadcaster, func()) {
	t.Helper()
	eventBroadcaster := postgres.NewEventBroadcaster(config.DatabaseURL(), 0, 0)
	eventBroadcaster.Start()
	return pipeline.NewORM(db, config, eventBroadcaster), eventBroadcaster, func() {
		eventBroadcaster.Close()
	}
}

func NewEthBroadcaster(t testing.TB, store *strpkg.Store, config *TestConfig, keys ...models.Key) (*bulletprooftxmanager.EthBroadcaster, func()) {
	t.Helper()
	eventBroadcaster := postgres.NewEventBroadcaster(config.DatabaseURL(), 0, 0)
	eventBroadcaster.Start()
	return bulletprooftxmanager.NewEthBroadcaster(store.DB, store.EthClient, config, store.KeyStore, &postgres.NullAdvisoryLocker{}, eventBroadcaster, keys), func() {
		eventBroadcaster.Close()
	}
}

// TestApplication holds the test application and test servers
type TestApplication struct {
	t testing.TB
	*chainlink.ChainlinkApplication
	Config           *TestConfig
	Logger           *logger.Logger
	Server           *httptest.Server
	wsServer         *httptest.Server
	connectedChannel chan struct{}
	Started          bool
	Backend          *backends.SimulatedBackend
	Key              models.Key
	allowUnstarted   bool
}

func newWSServer() (*httptest.Server, string, func()) {
	return NewWSServer("", nil)
}

// NewWSServer returns a  new wsserver
func NewWSServer(msg string, callback func(data []byte)) (*httptest.Server, string, func()) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		logger.PanicIf(err)
		for {
			_, data, err := conn.ReadMessage()
			if err != nil {
				break
			}

			if callback != nil {
				callback(data)
			}

			err = conn.WriteMessage(websocket.BinaryMessage, []byte(msg))
			if err != nil {
				break
			}
		}
	})
	server := httptest.NewServer(handler)

	u, err := url.Parse(server.URL)
	logger.PanicIf(err)
	u.Scheme = "ws"

	return server, u.String(), func() {
		server.Close()
	}
}

// NewApplication creates a New TestApplication along with a NewConfig
// It mocks the keystore with no keys or accounts by default
func NewApplication(t testing.TB, flagsAndDeps ...interface{}) (*TestApplication, func()) {
	t.Helper()

	c, cfgCleanup := NewConfig(t)

	app, cleanup := NewApplicationWithConfig(t, c, flagsAndDeps...)

	return app, func() {
		cleanup()
		cfgCleanup()
	}
}

// NewApplicationWithKey creates a new TestApplication along with a new config
// It uses the native keystore and will load any keys that are in the database
func NewApplicationWithKey(t testing.TB, flagsAndDeps ...interface{}) (*TestApplication, func()) {
	t.Helper()

	config, cfgCleanup := NewConfig(t)
	app, cleanup := NewApplicationWithConfigAndKey(t, config, flagsAndDeps...)
	return app, func() {
		cleanup()
		cfgCleanup()
	}
}

// NewApplicationWithConfigAndKey creates a new TestApplication with the given testconfig
// it will also provide an unlocked account on the keystore
func NewApplicationWithConfigAndKey(t testing.TB, tc *TestConfig, flagsAndDeps ...interface{}) (*TestApplication, func()) {
	t.Helper()

	app, cleanup := NewApplicationWithConfig(t, tc, flagsAndDeps...)
	for _, dep := range flagsAndDeps {
		switch v := dep.(type) {
		case models.Key:
			MustAddKeyToKeystore(t, &v, app.Store)
			app.Key = v
		}
	}
	if app.Key.Address.Address() == utils.ZeroAddress {
		app.Key, _ = MustAddRandomKeyToKeystore(t, app.Store, 0)
	}
	require.NoError(t, app.Store.KeyStore.Unlock(Password))

	return app, cleanup
}

// NewApplicationWithConfig creates a New TestApplication with specified test config
func NewApplicationWithConfig(t testing.TB, tc *TestConfig, flagsAndDeps ...interface{}) (*TestApplication, func()) {
	t.Helper()

	var ethClient eth.Client = &eth.NullClient{}
	var advisoryLocker postgres.AdvisoryLocker = &postgres.NullAdvisoryLocker{}
	var externalInitiatorManager chainlink.ExternalInitiatorManager = &services.NullExternalInitiatorManager{}
	var ks strpkg.KeyStoreInterface

	for _, flag := range flagsAndDeps {
		switch dep := flag.(type) {
		case eth.Client:
			ethClient = dep
		case postgres.AdvisoryLocker:
			advisoryLocker = dep
		case chainlink.ExternalInitiatorManager:
			externalInitiatorManager = dep
		case strpkg.KeyStoreInterface:
			ks = dep
		}
	}

	ta := &TestApplication{t: t, connectedChannel: make(chan struct{}, 1)}

	var keyStoreGenerator strpkg.KeyStoreGenerator
	if ks == nil {
		keyStoreGenerator = strpkg.InsecureKeyStoreGen
	} else {
		keyStoreGenerator = func(*gorm.DB, *orm.Config) strpkg.KeyStoreInterface {
			return ks
		}
	}

	appInstance, err := chainlink.NewApplication(tc.Config, ethClient, advisoryLocker, keyStoreGenerator, externalInitiatorManager, func(app chainlink.Application) {
		ta.connectedChannel <- struct{}{}
	})
	require.NoError(t, err)
	app := appInstance.(*chainlink.ChainlinkApplication)
	ta.ChainlinkApplication = app
	server := newServer(ta)

	tc.Config.Set("CLIENT_NODE_URL", server.URL)

	app.Store.Config = tc.Config

	for _, flag := range flagsAndDeps {
		if flag == AllowUnstarted {
			ta.allowUnstarted = true
		}
	}

	ta.Config = tc
	ta.Server = server
	ta.wsServer = tc.wsServer
	return ta, func() {
		ta.StopIfStarted()
	}
}

func NewEthMocks(t testing.TB) (*mocks.Client, *mocks.Subscription, func()) {
	c := new(mocks.Client)
	s := new(mocks.Subscription)
	var assertMocksCalled func()
	switch tt := t.(type) {
	case *testing.T:
		assertMocksCalled = func() {
			c.AssertExpectations(tt)
			s.AssertExpectations(tt)
		}
	case *testing.B:
		assertMocksCalled = func() {}
	}
	return c, s, assertMocksCalled
}

func NewEthMocksWithStartupAssertions(t testing.TB) (*mocks.Client, *mocks.Subscription, func()) {
	c, s, assertMocksCalled := NewEthMocks(t)
	c.On("Dial", mock.Anything).Maybe().Return(nil)
	c.On("ChainID", mock.Anything).Maybe().Return(NewTestConfig(t).ChainID(), nil)
	c.On("PendingNonceAt", mock.Anything, mock.Anything).Return(uint64(0), nil).Maybe()
	c.On("NonceAt", mock.Anything, mock.Anything, mock.Anything).Return(uint64(0), nil).Maybe()
	c.On("EthSubscribe", mock.Anything, mock.Anything, "newHeads").Maybe().Return(EmptyMockSubscription(), nil)
	c.On("SubscribeNewHead", mock.Anything, mock.Anything).Maybe().Return(EmptyMockSubscription(), nil)
	c.On("SendTransaction", mock.Anything, mock.Anything).Maybe().Return(nil)
	s.On("Err").Return(nil).Maybe()
	s.On("Unsubscribe").Return(nil).Maybe()
	return c, s, assertMocksCalled
}

func newServer(app chainlink.Application) *httptest.Server {
	engine := web.Router(app)
	return httptest.NewServer(engine)
}

func (ta *TestApplication) NewBox() packr.Box {
	ta.t.Helper()

	return packr.NewBox("../fixtures/operator_ui/dist")
}

func (ta *TestApplication) Start() error {
	ta.t.Helper()
	ta.Started = true
	ta.ChainlinkApplication.Store.KeyStore.Unlock(Password)

	err := ta.ChainlinkApplication.Start()
	return err
}

func (ta *TestApplication) StartAndConnect() error {
	ta.t.Helper()

	err := ta.Start()
	if err != nil {
		return err
	}

	return ta.waitForConnection()
}

// waitForConnection wait for the StartAndConnect callback to be called
func (ta *TestApplication) waitForConnection() error {
	select {
	case <-time.After(4 * time.Second):
		return errors.New("TestApplication#StartAndConnect() timed out")
	case <-ta.connectedChannel:
		return nil
	}
}

// Stop will stop the test application and perform cleanup
func (ta *TestApplication) Stop() error {
	ta.t.Helper()

	if !ta.Started {
		if ta.allowUnstarted {
			return nil
		}
		ta.t.Fatal("TestApplication Stop() called on an unstarted application")
	}

	// TODO: Here we double close, which is less than ideal.
	// We would prefer to invoke a method on an interface that
	// cleans up only in test.
	ta.ChainlinkApplication.StopIfStarted()
	cleanUpStore(ta.t, ta.Store)
	if ta.Server != nil {
		ta.Server.Close()
	}
	if ta.wsServer != nil {
		ta.wsServer.Close()
	}
	return nil
}

func (ta *TestApplication) MustSeedNewSession() string {
	session := NewSession()
	require.NoError(ta.t, ta.Store.SaveSession(&session))
	return session.ID
}

// ImportKey adds private key to the application disk keystore, not database.
func (ta *TestApplication) ImportKey(content string) {
	_, err := ta.Store.KeyStore.ImportKey([]byte(content), Password)
	require.NoError(ta.t, err)
	require.NoError(ta.t, ta.Store.KeyStore.Unlock(Password))
}

func (ta *TestApplication) NewHTTPClient() HTTPClientCleaner {
	ta.t.Helper()

	sessionID := ta.MustSeedNewSession()

	return HTTPClientCleaner{
		HTTPClient: NewMockAuthenticatedHTTPClient(ta.Config, sessionID),
		t:          ta.t,
	}
}

// NewClientAndRenderer creates a new cmd.Client for the test application
func (ta *TestApplication) NewClientAndRenderer() (*cmd.Client, *RendererMock) {
	sessionID := ta.MustSeedNewSession()
	r := &RendererMock{}
	client := &cmd.Client{
		Renderer:                       r,
		Config:                         ta.Config.Config,
		AppFactory:                     seededAppFactory{ta.ChainlinkApplication},
		KeyStoreAuthenticator:          CallbackAuthenticator{func(*strpkg.Store, string) (string, error) { return Password, nil }},
		FallbackAPIInitializer:         &MockAPIInitializer{},
		Runner:                         EmptyRunner{},
		HTTP:                           NewMockAuthenticatedHTTPClient(ta.Config, sessionID),
		CookieAuthenticator:            MockCookieAuthenticator{},
		FileSessionRequestBuilder:      &MockSessionRequestBuilder{},
		PromptingSessionRequestBuilder: &MockSessionRequestBuilder{},
		ChangePasswordPrompter:         &MockChangePasswordPrompter{},
	}
	return client, r
}

func (ta *TestApplication) NewAuthenticatingClient(prompter cmd.Prompter) *cmd.Client {
	cookieAuth := cmd.NewSessionCookieAuthenticator(ta.Config.Config, &cmd.MemoryCookieStore{})
	client := &cmd.Client{
		Renderer:                       &RendererMock{},
		Config:                         ta.Config.Config,
		AppFactory:                     seededAppFactory{ta.ChainlinkApplication},
		KeyStoreAuthenticator:          CallbackAuthenticator{func(*strpkg.Store, string) (string, error) { return Password, nil }},
		FallbackAPIInitializer:         &MockAPIInitializer{},
		Runner:                         EmptyRunner{},
		HTTP:                           cmd.NewAuthenticatedHTTPClient(ta.Config, cookieAuth, models.SessionRequest{}),
		CookieAuthenticator:            cookieAuth,
		FileSessionRequestBuilder:      cmd.NewFileSessionRequestBuilder(),
		PromptingSessionRequestBuilder: cmd.NewPromptingSessionRequestBuilder(prompter),
		ChangePasswordPrompter:         &MockChangePasswordPrompter{},
	}
	return client
}

func (ta *TestApplication) MustCreateJobRun(txHashBytes []byte, blockHashBytes []byte) *models.JobRun {
	job := NewJobWithWebInitiator()
	err := ta.Store.CreateJob(&job)
	require.NoError(ta.t, err)

	jr := NewJobRun(job)
	txHash := common.BytesToHash(txHashBytes)
	jr.RunRequest.TxHash = &txHash
	blockHash := common.BytesToHash(blockHashBytes)
	jr.RunRequest.BlockHash = &blockHash

	err = ta.Store.CreateJobRun(&jr)
	require.NoError(ta.t, err)

	return &jr
}

// NewStoreWithConfig creates a new store with given config
func NewStoreWithConfig(t testing.TB, config *TestConfig, flagsAndDeps ...interface{}) (*strpkg.Store, func()) {
	t.Helper()

	var advisoryLocker postgres.AdvisoryLocker = &postgres.NullAdvisoryLocker{}
	for _, flag := range flagsAndDeps {
		switch dep := flag.(type) {
		case postgres.AdvisoryLocker:
			advisoryLocker = dep
		}
	}
	s, err := strpkg.NewInsecureStore(config.Config, &eth.NullClient{}, advisoryLocker, gracefulpanic.NewSignal())
	if err != nil {
		require.NoError(t, err)
	}
	s.Config.SetRuntimeStore(s.ORM)
	return s, func() {
		cleanUpStore(config.t, s)
	}
}

// NewStore creates a new store
func NewStore(t testing.TB, flagsAndDeps ...interface{}) (*strpkg.Store, func()) {
	t.Helper()

	c, cleanup := NewConfig(t)
	store, storeCleanup := NewStoreWithConfig(t, c, flagsAndDeps...)
	return store, func() {
		storeCleanup()
		cleanup()
	}
}

func cleanUpStore(t testing.TB, store *strpkg.Store) {
	t.Helper()

	defer func() {
		if err := os.RemoveAll(store.Config.RootDir()); err != nil {
			logger.Warn("unable to clear test store:", err)
		}
	}()
	logger.Sync()
	require.NoError(t, store.Close())
}

func ParseJSON(t testing.TB, body io.Reader) models.JSON {
	t.Helper()

	b, err := ioutil.ReadAll(body)
	require.NoError(t, err)
	return models.JSON{Result: gjson.ParseBytes(b)}
}

func ParseJSONAPIErrors(t testing.TB, body io.Reader) *models.JSONAPIErrors {
	t.Helper()

	b, err := ioutil.ReadAll(body)
	require.NoError(t, err)
	var respJSON models.JSONAPIErrors
	json.Unmarshal(b, &respJSON)
	return &respJSON
}

// MustReadFile loads a file but should never fail
func MustReadFile(t testing.TB, file string) []byte {
	t.Helper()

	content, err := ioutil.ReadFile(file)
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

	b, err := ioutil.ReadAll(resp.Body)
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

// ReadLogs returns the contents of the applications log file as a string
func ReadLogs(config orm.ConfigReader) (string, error) {
	logFile := fmt.Sprintf("%s/log.jsonl", config.RootDir())
	b, err := ioutil.ReadFile(logFile)
	return string(b), err
}

func FindServiceAgreement(t testing.TB, s *strpkg.Store, id string) models.ServiceAgreement {
	t.Helper()

	sa, err := s.FindServiceAgreement(id)
	require.NoError(t, err)

	return sa
}

// CreateJobSpecViaWeb creates a jobspec via web using /v2/specs
func CreateJobSpecViaWeb(t testing.TB, app *TestApplication, job models.JobSpec) models.JobSpec {
	t.Helper()

	marshaled, err := json.Marshal(&job)
	assert.NoError(t, err)
	return CreateSpecViaWeb(t, app, string(marshaled))
}

// CreateJobSpecViaWeb creates a jobspec via web using /v2/specs
func CreateSpecViaWeb(t testing.TB, app *TestApplication, spec string) models.JobSpec {
	t.Helper()

	client := app.NewHTTPClient()
	resp, cleanup := client.Post("/v2/specs", bytes.NewBufferString(spec))
	defer cleanup()
	AssertServerResponse(t, resp, http.StatusOK)

	var createdJob models.JobSpec
	err := ParseJSONAPIResponse(t, resp, &createdJob)
	require.NoError(t, err)
	return createdJob
}

func CreateJobViaWeb(t testing.TB, app *TestApplication, request []byte) job.Job {
	t.Helper()

	client := app.NewHTTPClient()
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

	client := app.NewHTTPClient()
	resp, cleanup := client.Post("/v2/jobs", bytes.NewBufferString(spec))
	defer cleanup()
	AssertServerResponse(t, resp, http.StatusOK)

	var jobResponse webpresenters.JobResource
	err := ParseJSONAPIResponse(t, resp, &jobResponse)
	require.NoError(t, err)
	return jobResponse
}

// CreateJobRunViaWeb creates JobRun via web using /v2/specs/ID/runs
func CreateJobRunViaWeb(t testing.TB, app *TestApplication, j models.JobSpec, body ...string) models.JobRun {
	t.Helper()

	bodyBuffer := &bytes.Buffer{}
	if len(body) > 0 {
		bodyBuffer = bytes.NewBufferString(body[0])
	}
	client := app.NewHTTPClient()
	resp, cleanup := client.Post("/v2/specs/"+j.ID.String()+"/runs", bodyBuffer)
	defer cleanup()
	AssertServerResponse(t, resp, http.StatusOK)
	var jr models.JobRun
	err := ParseJSONAPIResponse(t, resp, &jr)
	require.NoError(t, err)

	assert.Equal(t, j.ID, jr.JobSpecID)
	return jr
}

func CreateJobRunViaExternalInitiator(
	t testing.TB,
	app *TestApplication,
	j models.JobSpec,
	eia auth.Token,
	body string,
) models.JobRun {
	t.Helper()

	headers := make(map[string]string)
	headers[static.ExternalInitiatorAccessKeyHeader] = eia.AccessKey
	headers[static.ExternalInitiatorSecretHeader] = eia.Secret

	url := app.Config.ClientNodeURL() + "/v2/specs/" + j.ID.String() + "/runs"
	bodyBuf := bytes.NewBufferString(body)
	resp, cleanup := UnauthenticatedPost(t, url, bodyBuf, headers)
	defer cleanup()
	AssertServerResponse(t, resp, 200)
	var jr models.JobRun
	err := ParseJSONAPIResponse(t, resp, &jr)
	require.NoError(t, err)

	assert.Equal(t, j.ID, jr.JobSpecID)
	return jr
}

// CreateHelloWorldJobViaWeb creates a HelloWorld JobSpec with the given MockServer Url
func CreateHelloWorldJobViaWeb(t testing.TB, app *TestApplication, url string) models.JobSpec {
	t.Helper()

	buffer := []byte(`
{
  "initiators": [{ "type": "web" }],
  "tasks": [
    { "type": "HTTPGetWithUnrestrictedNetworkAccess", "params": {
		"get": "https://bitstamp.net/api/ticker/",
        "headers": {
          "Key1": ["value"],
          "Key2": ["value", "value"]
        }
      }
    },
    { "type": "JsonParse", "params": { "path": ["last"] }},
    { "type": "EthBytes32" },
    {
      "type": "EthTx", "params": {
        "address": "0x356a04bce728ba4c62a30294a55e6a8600a320b3",
        "functionSelector": "0x609ff1bd"
      }
    }
  ]
}
`)

	var job models.JobSpec
	err := json.Unmarshal(buffer, &job)
	require.NoError(t, err)

	data, err := models.Merge(job.Tasks[0].Params, JSONFromString(t, `{"get":"%v"}`, url))
	require.NoError(t, err)
	job.Tasks[0].Params = data
	return CreateJobSpecViaWeb(t, app, job)
}

// UpdateJobRunViaWeb updates jobrun via web using /v2/runs/ID
func UpdateJobRunViaWeb(
	t testing.TB,
	app *TestApplication,
	jr models.JobRun,
	bridgeResource *webpresenters.BridgeResource,
	body string,
) models.JobRun {
	t.Helper()

	client := app.NewHTTPClient()
	headers := map[string]string{"Authorization": "Bearer " + bridgeResource.IncomingToken}
	resp, cleanup := client.Patch("/v2/runs/"+jr.ID.String(), bytes.NewBufferString(body), headers)
	defer cleanup()

	AssertServerResponse(t, resp, http.StatusOK)
	var respJobRun presenters.JobRun
	assert.NoError(t, ParseJSONAPIResponse(t, resp, &respJobRun))
	assert.Equal(t, jr.ID, respJobRun.ID)
	jr = respJobRun.JobRun
	return jr
}

// CreateBridgeTypeViaWeb creates a bridgetype via web using /v2/bridge_types
func CreateBridgeTypeViaWeb(
	t testing.TB,
	app *TestApplication,
	payload string,
) *webpresenters.BridgeResource {
	t.Helper()

	client := app.NewHTTPClient()
	resp, cleanup := client.Post(
		"/v2/bridge_types",
		bytes.NewBufferString(payload),
	)
	defer cleanup()
	AssertServerResponse(t, resp, http.StatusOK)
	bt := &webpresenters.BridgeResource{}
	err := ParseJSONAPIResponse(t, resp, bt)
	require.NoError(t, err)

	return bt
}

// CreateExternalInitiatorViaWeb creates a bridgetype via web using /v2/bridge_types
func CreateExternalInitiatorViaWeb(
	t testing.TB,
	app *TestApplication,
	payload string,
) *presenters.ExternalInitiatorAuthentication {
	t.Helper()

	client := app.NewHTTPClient()
	resp, cleanup := client.Post(
		"/v2/external_initiators",
		bytes.NewBufferString(payload),
	)
	defer cleanup()
	AssertServerResponse(t, resp, http.StatusCreated)
	ei := &presenters.ExternalInitiatorAuthentication{}
	err := ParseJSONAPIResponse(t, resp, ei)
	require.NoError(t, err)

	return ei
}

const (
	// DBWaitTimeout is how long we wait by default for something to appear in
	// the DB. It needs to be fairly long because integration
	// tests rely on it.
	DBWaitTimeout = 20 * time.Second
	// DBPollingInterval can't be too short to avoid DOSing the test database
	DBPollingInterval = 100 * time.Millisecond
	// AssertNoActionTimeout shouldn't be too long, or it will slow down tests
	AssertNoActionTimeout = 3 * time.Second
)

// WaitForJobRunToComplete waits for a JobRun to reach Completed Status
func WaitForJobRunToComplete(
	t testing.TB,
	store *strpkg.Store,
	jr models.JobRun,
) models.JobRun {
	t.Helper()

	return WaitForJobRunStatus(t, store, jr, models.RunStatusCompleted)
}

// WaitForJobRunToPendBridge waits for a JobRun to reach PendingBridge Status
func WaitForJobRunToPendBridge(
	t testing.TB,
	store *strpkg.Store,
	jr models.JobRun,
) models.JobRun {
	t.Helper()

	return WaitForJobRunStatus(t, store, jr, models.RunStatusPendingBridge)
}

// WaitForJobRunToPendIncomingConfirmations waits for a JobRun to reach PendingIncomingConfirmations Status
func WaitForJobRunToPendIncomingConfirmations(
	t testing.TB,
	store *strpkg.Store,
	jr models.JobRun,
) models.JobRun {
	t.Helper()
	return WaitForJobRunStatus(t, store, jr, models.RunStatusPendingIncomingConfirmations)
}

// WaitForJobRunToPendOutgoingConfirmations waits for a JobRun to reach PendingOutgoingConfirmations Status
func WaitForJobRunToPendOutgoingConfirmations(
	t testing.TB,
	store *strpkg.Store,
	jr models.JobRun,
) models.JobRun {
	t.Helper()
	return WaitForJobRunStatus(t, store, jr, models.RunStatusPendingOutgoingConfirmations)
}

func SendBlocksUntilComplete(
	t testing.TB,
	store *strpkg.Store,
	jr models.JobRun,
	blockCh chan<- *models.Head,
	start int64,
	ethClient *mocks.Client,
) models.JobRun {
	t.Helper()

	var err error
	block := start
	gomega.NewGomegaWithT(t).Eventually(func() models.RunStatus {
		h := models.NewHead(big.NewInt(block), NewHash(), NewHash(), 0)
		blockCh <- &h
		block++
		jr, err = store.Unscoped().FindJobRun(jr.ID)
		assert.NoError(t, err)
		st := jr.GetStatus()
		return st
	}, DBWaitTimeout, DBPollingInterval).Should(gomega.Equal(models.RunStatusCompleted))
	return jr
}

// WaitForJobRunStatus waits for a JobRun to reach given status
func WaitForJobRunStatus(
	t testing.TB,
	store *strpkg.Store,
	jr models.JobRun,
	wantStatus models.RunStatus,

) models.JobRun {
	t.Helper()

	var err error
	gomega.NewGomegaWithT(t).Eventually(func() models.RunStatus {
		jr, err = store.Unscoped().FindJobRun(jr.ID)
		assert.NoError(t, err)
		st := jr.GetStatus()
		if wantStatus != models.RunStatusErrored {
			if st == models.RunStatusErrored {
				t.Fatalf("waiting for job run status %s but got %s, error was: '%s'", wantStatus, models.RunStatusErrored, jr.Result.ErrorMessage.String)
			}
		}
		return st
	}, DBWaitTimeout, DBPollingInterval).Should(gomega.Equal(wantStatus))
	return jr
}

// JobRunStays tests if a JobRun will consistently stay at the specified status
func JobRunStays(
	t testing.TB,
	store *strpkg.Store,
	jr models.JobRun,
	status models.RunStatus,
	optionalDuration ...time.Duration,
) models.JobRun {
	t.Helper()

	duration := time.Second
	if len(optionalDuration) > 0 {
		duration = optionalDuration[0]
	}

	var err error
	gomega.NewGomegaWithT(t).Consistently(func() models.RunStatus {
		jr, err = store.FindJobRun(jr.ID)
		assert.NoError(t, err)
		return jr.GetStatus()
	}, duration, DBPollingInterval).Should(gomega.Equal(status))
	return jr
}

// JobRunStaysPendingIncomingConfirmations tests if a JobRun will stay at the PendingIncomingConfirmations Status
func JobRunStaysPendingIncomingConfirmations(
	t testing.TB,
	store *strpkg.Store,
	jr models.JobRun,
) models.JobRun {
	t.Helper()

	return JobRunStays(t, store, jr, models.RunStatusPendingIncomingConfirmations)
}

// Polls until the passed in jobID has count number
// of job spec errors.
func WaitForSpecError(t *testing.T, store *strpkg.Store, jobID models.JobID, count int) []models.JobSpecError {
	t.Helper()
	g := gomega.NewGomegaWithT(t)
	var jse []models.JobSpecError
	g.Eventually(func() []models.JobSpecError {
		err := store.DB.
			Where("job_spec_id = ?", jobID.String()).
			Find(&jse).Error
		assert.NoError(t, err)
		return jse
	}, DBWaitTimeout, DBPollingInterval).Should(gomega.HaveLen(count))
	return jse
}

// WaitForSpecErrorV2 polls until the passed in jobID has count number
// of job spec errors.
func WaitForSpecErrorV2(t *testing.T, store *strpkg.Store, jobID int32, count int) []job.SpecError {
	t.Helper()

	g := gomega.NewGomegaWithT(t)
	var jse []job.SpecError
	g.Eventually(func() []job.SpecError {
		err := store.DB.
			Where("job_id = ?", jobID).
			Find(&jse).Error
		assert.NoError(t, err)
		return jse
	}, DBWaitTimeout, DBPollingInterval).Should(gomega.HaveLen(count))
	return jse
}

// WaitForRuns waits for the wanted number of runs then returns a slice of the JobRuns
func WaitForRuns(t testing.TB, j models.JobSpec, store *strpkg.Store, want int) []models.JobRun {
	t.Helper()
	g := gomega.NewGomegaWithT(t)

	var jrs []models.JobRun
	var err error
	if want == 0 {
		g.Consistently(func() []models.JobRun {
			jrs, err = store.JobRunsFor(j.ID)
			assert.NoError(t, err)
			return jrs
		}, AssertNoActionTimeout, DBPollingInterval).Should(gomega.HaveLen(want))
	} else {
		g.Eventually(func() []models.JobRun {
			jrs, err = store.JobRunsFor(j.ID)
			assert.NoError(t, err)
			return jrs
		}, DBWaitTimeout, DBPollingInterval).Should(gomega.HaveLen(want))
	}
	return jrs
}

func WaitForPipelineRuns(t testing.TB, nodeID int, jobID int32, jo job.ORM, want int, timeout, poll time.Duration) []pipeline.Run {
	t.Helper()

	var err error
	prs := []pipeline.Run{}
	gomega.NewGomegaWithT(t).Eventually(func() []pipeline.Run {
		prs, _, err = jo.PipelineRunsByJobID(jobID, 0, 1000)
		assert.NoError(t, err)
		return prs
	}, timeout, poll).Should(gomega.HaveLen(want))

	return prs
}

func WaitForPipelineComplete(t testing.TB, nodeID int, jobID int32, count int, expectedTaskRuns int, jo job.ORM, timeout, poll time.Duration) []pipeline.Run {
	t.Helper()

	var pr []pipeline.Run
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		prs, _, err := jo.PipelineRunsByJobID(jobID, 0, 1000)
		assert.NoError(t, err)
		var completed []pipeline.Run

		for i := range prs {
			if !prs[i].Outputs.Null {
				if !prs[i].Errors.HasError() {
					// txdb effectively ignores transactionality of queries, so we need to explicitly expect a number of task runs
					// (if the read occurrs mid-transaction and a job run in inserted but task runs not yet).
					if len(prs[i].PipelineTaskRuns) == expectedTaskRuns {
						completed = append(completed, prs[i])
					}
				}
			}
		}
		if len(completed) >= count {
			pr = completed
			return true
		}
		return false
	}, timeout, poll).Should(gomega.BeTrue(), fmt.Sprintf("job %d on node %d not complete with %d runs", jobID, nodeID, count))
	return pr
}

// AssertRunsStays asserts that the number of job runs for a particular job remains at the provided values
func AssertRunsStays(t testing.TB, j models.JobSpec, store *strpkg.Store, want int) []models.JobRun {
	t.Helper()
	g := gomega.NewGomegaWithT(t)

	var jrs []models.JobRun
	var err error
	g.Consistently(func() []models.JobRun {
		jrs, err = store.JobRunsFor(j.ID)
		assert.NoError(t, err)
		return jrs
	}, AssertNoActionTimeout, DBPollingInterval).Should(gomega.HaveLen(want))
	return jrs
}

// AssertPipelineRunsStays asserts that the number of pipeline runs for a particular job remains at the provided values
func AssertPipelineRunsStays(t testing.TB, pipelineSpecID int32, store *strpkg.Store, want int) []pipeline.Run {
	t.Helper()
	g := gomega.NewGomegaWithT(t)

	var prs []pipeline.Run
	g.Consistently(func() []pipeline.Run {
		err := store.DB.
			Where("pipeline_spec_id = ?", pipelineSpecID).
			Find(&prs).Error
		assert.NoError(t, err)
		return prs
	}, AssertNoActionTimeout, DBPollingInterval).Should(gomega.HaveLen(want))
	return prs
}

// WaitForRunsAtLeast waits for at least the passed number of runs to start.
func WaitForRunsAtLeast(t testing.TB, j models.JobSpec, store *strpkg.Store, want int) {
	t.Helper()
	g := gomega.NewGomegaWithT(t)

	if want == 0 {
		t.Fatal("must want more than 0 runs when waiting")
	} else {
		g.Eventually(func() int {
			jrs, err := store.JobRunsFor(j.ID)
			require.NoError(t, err)
			return len(jrs)
		}, DBWaitTimeout, DBPollingInterval).Should(gomega.BeNumerically(">=", want))
	}
}

func WaitForEthTxAttemptsForEthTx(t testing.TB, store *strpkg.Store, ethTx models.EthTx) []models.EthTxAttempt {
	t.Helper()
	g := gomega.NewGomegaWithT(t)

	var attempts []models.EthTxAttempt
	var err error
	g.Eventually(func() int {
		err = store.DB.Order("created_at desc").Where("eth_tx_id = ?", ethTx.ID).Find(&attempts).Error
		assert.NoError(t, err)
		return len(attempts)
	}, DBWaitTimeout, DBPollingInterval).Should(gomega.BeNumerically(">", 0))
	return attempts
}

func WaitForEthTxAttemptCount(t testing.TB, store *strpkg.Store, want int) []models.EthTxAttempt {
	t.Helper()
	g := gomega.NewGomegaWithT(t)

	var txas []models.EthTxAttempt
	var err error
	g.Eventually(func() []models.EthTxAttempt {
		err = store.DB.Find(&txas).Error
		assert.NoError(t, err)
		return txas
	}, DBWaitTimeout, DBPollingInterval).Should(gomega.HaveLen(want))
	return txas
}

// AssertEthTxAttemptCountStays asserts that the number of tx attempts remains at the provided value
func AssertEthTxAttemptCountStays(t testing.TB, store *strpkg.Store, want int) []models.EthTxAttempt {
	t.Helper()
	g := gomega.NewGomegaWithT(t)

	var txas []models.EthTxAttempt
	var err error
	g.Consistently(func() []models.EthTxAttempt {
		err = store.DB.Find(&txas).Error
		assert.NoError(t, err)
		return txas
	}, AssertNoActionTimeout, DBPollingInterval).Should(gomega.HaveLen(want))
	return txas
}

// WaitForSyncEventCount checks if the sync event count eventually reaches
// the amound specified in parameter want.
func WaitForSyncEventCount(
	t testing.TB,
	orm *orm.ORM,
	want int,
) {
	t.Helper()
	gomega.NewGomegaWithT(t).Eventually(func() int {
		count, err := orm.CountOf(&models.SyncEvent{})
		assert.NoError(t, err)
		return count
	}, DBWaitTimeout, DBPollingInterval).Should(gomega.Equal(want))
}

// AssertSyncEventCountStays ensures that the event sync count stays consistent
// for a period of time
func AssertSyncEventCountStays(
	t testing.TB,
	orm *orm.ORM,
	want int,
) {
	t.Helper()
	gomega.NewGomegaWithT(t).Consistently(func() int {
		count, err := orm.CountOf(&models.SyncEvent{})
		assert.NoError(t, err)
		return count
	}, AssertNoActionTimeout, DBPollingInterval).Should(gomega.Equal(want))
}

// ParseISO8601 given the time string it Must parse the time and return it
func ParseISO8601(t testing.TB, s string) time.Time {
	t.Helper()

	tm, err := time.Parse(time.RFC3339Nano, s)
	require.NoError(t, err)
	return tm
}

// NullableTime will return a valid nullable time given time.Time
func NullableTime(t time.Time) null.Time {
	return null.TimeFrom(t)
}

// ParseNullableTime given a time string parse it into a null.Time
func ParseNullableTime(t testing.TB, s string) null.Time {
	t.Helper()

	return NullableTime(ParseISO8601(t, s))
}

// Head given the value convert it into an Head
func Head(val interface{}) *models.Head {
	var h models.Head
	time := uint64(0)
	switch t := val.(type) {
	case int:
		h = models.NewHead(big.NewInt(int64(t)), NewHash(), NewHash(), time)
	case uint64:
		h = models.NewHead(big.NewInt(int64(t)), NewHash(), NewHash(), time)
	case int64:
		h = models.NewHead(big.NewInt(t), NewHash(), NewHash(), time)
	case *big.Int:
		h = models.NewHead(t, NewHash(), NewHash(), time)
	default:
		logger.Panicf("Could not convert %v of type %T to Head", val, val)
	}
	return &h
}

// TransactionsFromGasPrices returns transactions matching the given gas prices
func TransactionsFromGasPrices(gasPrices ...int64) []gasupdater.Transaction {
	txs := make([]gasupdater.Transaction, len(gasPrices))
	for i, gasPrice := range gasPrices {
		txs[i] = gasupdater.Transaction{GasPrice: big.NewInt(gasPrice), GasLimit: 42}
	}
	return txs
}

// BlockWithTransactions returns a new ethereum block with transactions
// matching the given gas prices
func BlockWithTransactions(gasPrices ...int64) *types.Block {
	txs := make([]*types.Transaction, len(gasPrices))
	for i, gasPrice := range gasPrices {
		txs[i] = types.NewTransaction(0, common.Address{}, nil, 0, big.NewInt(gasPrice), nil)
	}
	return types.NewBlock(&types.Header{}, txs, nil, nil, new(trie.Trie))
}

func StringToHash(s string) common.Hash {
	return common.BytesToHash([]byte(s))
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
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			assert.FailNowf(t, "Unable to read body", err.Error())
		}

		var result map[string][]string
		err = json.Unmarshal(b, &result)
		if err != nil {
			assert.FailNowf(t, fmt.Sprintf("Unable to unmarshal json from body '%s'", string(b)), err.Error())
		}

		assert.FailNowf(t, "Request failed", "Expected %d response, got %d with errors: %s", expectedStatusCode, resp.StatusCode, result["errors"])
	} else {
		assert.FailNowf(t, "Unexpected response", "Expected %d response, got %d", expectedStatusCode, resp.StatusCode)
	}
}

func DecodeSessionCookie(value string) (string, error) {
	var decrypted map[interface{}]interface{}
	codecs := securecookie.CodecsFromPairs([]byte(SessionSecret))
	err := securecookie.DecodeMulti(web.SessionName, value, &decrypted, codecs...)
	if err != nil {
		return "", err
	}
	value, ok := decrypted[web.SessionIDKey].(string)
	if !ok {
		return "", fmt.Errorf("decrypted[web.SessionIDKey] is not a string (%v)", value)
	}
	return value, nil
}

func MustGenerateSessionCookie(value string) *http.Cookie {
	decrypted := map[interface{}]interface{}{web.SessionIDKey: value}
	codecs := securecookie.CodecsFromPairs([]byte(SessionSecret))
	encoded, err := securecookie.EncodeMulti(web.SessionName, decrypted, codecs...)
	if err != nil {
		logger.Panic(err)
	}
	return sessions.NewCookie(web.SessionName, encoded, &sessions.Options{})
}

func NormalizedJSON(t testing.TB, input []byte) string {
	t.Helper()

	normalized, err := utils.NormalizedJSON(input)
	require.NoError(t, err)
	return normalized
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

	client := http.Client{}
	request, err := http.NewRequest("POST", url, body)
	require.NoError(t, err)
	request.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		request.Header.Add(key, value)
	}
	resp, err := client.Do(request)
	require.NoError(t, err)
	return resp, func() { resp.Body.Close() }
}

func UnauthenticatedPatch(t testing.TB, url string, body io.Reader, headers map[string]string) (*http.Response, func()) {
	t.Helper()

	client := http.Client{}
	request, err := http.NewRequest("PATCH", url, body)
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

func NewSession(optionalSessionID ...string) models.Session {
	session := models.NewSession()
	if len(optionalSessionID) > 0 {
		session.ID = optionalSessionID[0]
	}
	return session
}

func AllExternalInitiators(t testing.TB, store *strpkg.Store) []models.ExternalInitiator {
	t.Helper()

	var all []models.ExternalInitiator
	err := store.RawDBWithAdvisoryLock(func(db *gorm.DB) error {
		return db.Find(&all).Error
	})
	require.NoError(t, err)
	return all
}

func AllJobs(t testing.TB, store *strpkg.Store) []models.JobSpec {
	t.Helper()

	var all []models.JobSpec
	err := store.ORM.RawDBWithAdvisoryLock(func(db *gorm.DB) error {
		return db.Find(&all).Error
	})
	require.NoError(t, err)
	return all
}

func MustAllJobsWithStatus(t testing.TB, store *strpkg.Store, statuses ...models.RunStatus) []*models.JobRun {
	t.Helper()

	var runs []*models.JobRun
	err := store.UnscopedJobRunsWithStatus(func(jr *models.JobRun) {
		runs = append(runs, jr)
	}, statuses...)
	require.NoError(t, err)
	return runs
}

func GetLastEthTxAttempt(t testing.TB, store *strpkg.Store) models.EthTxAttempt {
	t.Helper()

	var txa models.EthTxAttempt
	var count int64
	err := store.ORM.RawDBWithAdvisoryLock(func(db *gorm.DB) error {
		return db.Order("created_at desc").First(&txa).Count(&count).Error
	})
	require.NoError(t, err)
	require.NotEqual(t, 0, count)
	return txa
}

type Awaiter chan struct{}

func NewAwaiter() Awaiter { return make(Awaiter) }

func (a Awaiter) ItHappened() { close(a) }

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
		t.Fatal(fmt.Sprintf("CallbackOrTimeout: %s timed out", msg))
	}
}

func MustParseURL(input string) *url.URL {
	u, err := url.Parse(input)
	if err != nil {
		logger.Panic(err)
	}
	return u
}

func MustResultString(t *testing.T, input models.RunResult) string {
	result := input.Data.Get("result")
	require.Equal(t, gjson.String, result.Type, fmt.Sprintf("result type %s is not string", result.Type))
	return result.String()
}

// GenericEncode eth encodes values based on the provided types
func GenericEncode(types []string, values ...interface{}) ([]byte, error) {
	if len(values) != len(types) {
		return nil, errors.New("must include same number of values as types")
	}
	var args abi.Arguments
	for _, t := range types {
		ty, _ := abi.NewType(t, "", nil)
		args = append(args, abi.Argument{Type: ty})
	}
	out, err := args.PackValues(values)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func MustGenericEncode(types []string, values ...interface{}) []byte {
	if len(values) != len(types) {
		panic("must include same number of values as types")
	}
	var args abi.Arguments
	for _, t := range types {
		ty, _ := abi.NewType(t, "", nil)
		args = append(args, abi.Argument{Type: ty})
	}
	out, err := args.PackValues(values)
	if err != nil {
		panic(err)
	}
	return out
}

func MakeRoundStateReturnData(
	roundID uint64,
	eligible bool,
	answer, startAt, timeout, availableFunds, paymentAmount, oracleCount uint64,
) []byte {
	var data []byte
	if eligible {
		data = append(data, utils.EVMWordUint64(1)...)
	} else {
		data = append(data, utils.EVMWordUint64(0)...)
	}
	data = append(data, utils.EVMWordUint64(roundID)...)
	data = append(data, utils.EVMWordUint64(answer)...)
	data = append(data, utils.EVMWordUint64(startAt)...)
	data = append(data, utils.EVMWordUint64(timeout)...)
	data = append(data, utils.EVMWordUint64(availableFunds)...)
	data = append(data, utils.EVMWordUint64(oracleCount)...)
	data = append(data, utils.EVMWordUint64(paymentAmount)...)
	return data
}

var fluxAggregatorABI = eth.MustGetABI(flux_aggregator_wrapper.FluxAggregatorABI)

func MockFluxAggCall(client *mocks.Client, address common.Address, funcName string) *mock.Call {
	funcSig := hexutil.Encode(fluxAggregatorABI.Methods[funcName].ID)
	if len(funcSig) != 10 {
		panic(fmt.Sprintf("Unable to find FluxAgg function with name %s", funcName))
	}
	return client.On(
		"CallContract",
		mock.Anything,
		mock.MatchedBy(func(callArgs ethereum.CallMsg) bool {
			return *callArgs.To == address &&
				hexutil.Encode(callArgs.Data)[0:10] == funcSig
		}),
		mock.Anything)
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
func MockApplicationEthCalls(t *testing.T, app *TestApplication, ethClient *mocks.Client) (verify func()) {
	t.Helper()

	// Start
	ethClient.On("Dial", mock.Anything).Return(nil)
	sub := new(mocks.Subscription)
	sub.On("Err").Return(nil)
	ethClient.On("SubscribeNewHead", mock.Anything, mock.Anything).Return(sub, nil)
	ethClient.On("ChainID", mock.Anything).Return(app.Store.Config.ChainID(), nil)
	ethClient.On("PendingNonceAt", mock.Anything, mock.Anything).Return(uint64(0), nil).Maybe()

	// Stop
	sub.On("Unsubscribe").Return(nil)

	return func() {
		ethClient.AssertExpectations(t)
	}
}

func MockSubscribeToLogsCh(ethClient *mocks.Client, sub *mocks.Subscription) chan chan<- types.Log {
	logsCh := make(chan chan<- types.Log, 1)
	ethClient.On("SubscribeFilterLogs", mock.Anything, mock.Anything, mock.Anything).
		Return(sub, nil).
		Run(func(args mock.Arguments) { // context.Context, ethereum.FilterQuery, chan<- types.Log
			logsCh <- args.Get(2).(chan<- types.Log)
		})
	return logsCh
}

func MustNewJSONSerializable(t *testing.T, s string) pipeline.JSONSerializable {
	t.Helper()

	js := new(pipeline.JSONSerializable)
	err := js.UnmarshalJSON([]byte(s))
	require.NoError(t, err)
	return *js
}

func BatchElemMatchesHash(req rpc.BatchElem, hash common.Hash) bool {
	return req.Method == "eth_getTransactionReceipt" &&
		len(req.Args) == 1 && req.Args[0] == hash
}

func BatchElemMustMatchHash(t *testing.T, req rpc.BatchElem, hash common.Hash) {
	t.Helper()
	if !BatchElemMatchesHash(req, hash) {
		t.Fatalf("Batch hash %v does not match expected %v", req.Args[0], hash)
	}
}

type SimulateIncomingHeadsArgs struct {
	StartBlock, EndBlock int64
	BackfillDepth        int64
	Interval             time.Duration
	Timeout              time.Duration
	HeadTrackables       []httypes.HeadTrackable
	Hashes               map[int64]common.Hash
}

func SimulateIncomingHeads(t *testing.T, args SimulateIncomingHeadsArgs) (func(), chan struct{}) {
	t.Helper()

	if args.BackfillDepth == 0 {
		t.Fatal("BackfillDepth must be > 0")
	}

	// Build the full chain of heads
	heads := make(map[int64]*models.Head)
	first := args.StartBlock - args.BackfillDepth
	if first < 0 {
		first = 0
	}
	last := args.EndBlock
	if last == 0 {
		last = args.StartBlock + 300 // If no .EndBlock is provided, assume we want 300 heads
	}
	for i := first; i <= last; i++ {
		// If a particular block should have a particular
		// hash, use that. Otherwise, generate a random one.
		var hash common.Hash
		if args.Hashes != nil {
			if h, exists := args.Hashes[i]; exists {
				hash = h
			}
		}
		if hash == (common.Hash{}) {
			hash = NewHash()
		}
		heads[i] = &models.Head{Hash: hash, Number: i}
		if i > first {
			heads[i].Parent = heads[i-1]
		}
	}

	if args.Timeout == 0 {
		args.Timeout = 60 * time.Second
	}
	if args.Interval == 0 {
		args.Interval = 250 * time.Millisecond
	}
	ctx, cancel := context.WithTimeout(context.Background(), args.Timeout)
	defer cancel()
	chTimeout := time.After(args.Timeout)

	chDone := make(chan struct{})
	go func() {
		current := int64(args.StartBlock)
		for {
			select {
			case <-chDone:
				return
			case <-chTimeout:
				return
			default:
				// Trim chain to backfill depth
				ptr := heads[current]
				for i := int64(0); i < args.BackfillDepth && ptr.Parent != nil; i++ {
					ptr = ptr.Parent
				}
				ptr.Parent = nil

				for _, ht := range args.HeadTrackables {
					ht.OnNewLongestChain(ctx, *heads[current])
				}
				if args.EndBlock >= 0 && current == args.EndBlock {
					chDone <- struct{}{}
					return
				}
				current++
				time.Sleep(args.Interval)
			}
		}
	}()
	var once sync.Once
	cleanup := func() {
		once.Do(func() {
			close(chDone)
			cancel()
		})
	}
	return cleanup, chDone
}

type HeadTrackableFunc func(context.Context, models.Head)

func (HeadTrackableFunc) Connect(*models.Head) error { return nil }
func (HeadTrackableFunc) Disconnect()                {}
func (fn HeadTrackableFunc) OnNewLongestChain(ctx context.Context, head models.Head) {
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

func AssertCount(t *testing.T, store *strpkg.Store, model interface{}, expected int64) {
	t.Helper()
	var count int64
	err := store.DB.Unscoped().Model(model).Count(&count).Error
	require.NoError(t, err)
	require.Equal(t, expected, count)
}

func WaitForCount(t testing.TB, store *strpkg.Store, model interface{}, want int64) {
	t.Helper()
	g := gomega.NewGomegaWithT(t)
	var count int64
	var err error
	g.Eventually(func() int64 {
		err = store.DB.Model(model).Count(&count).Error
		assert.NoError(t, err)
		return count
	}, DBWaitTimeout, DBPollingInterval).Should(gomega.Equal(want))
}

func AssertCountStays(t testing.TB, store *strpkg.Store, model interface{}, want int64) {
	t.Helper()
	g := gomega.NewGomegaWithT(t)
	var count int64
	var err error
	g.Consistently(func() int64 {
		err = store.DB.Model(model).Count(&count).Error
		assert.NoError(t, err)
		return count
	}, AssertNoActionTimeout, DBPollingInterval).Should(gomega.Equal(want))
}

func AssertRecordEventually(t *testing.T, store *strpkg.Store, model interface{}, check func() bool) {
	t.Helper()
	g := gomega.NewGomegaWithT(t)
	g.Eventually(func() bool {
		err := store.DB.Find(model).Error
		require.NoError(t, err, "unable to find record in DB")
		return check()
	}, DBWaitTimeout, DBPollingInterval).Should(gomega.BeTrue())
}

func MustSendingKeys(t *testing.T, ks strpkg.KeyStoreInterface) (keys []models.Key) {
	var err error
	keys, err = ks.SendingKeys()
	require.NoError(t, err)
	return keys
}
