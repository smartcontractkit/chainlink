package cltest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/big"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/auth"
	"github.com/smartcontractkit/chainlink/core/cmd"
	"github.com/smartcontractkit/chainlink/core/gracefulpanic"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	strpkg "github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/smartcontractkit/chainlink/core/store/presenters"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/web"

	"github.com/DATA-DOG/go-txdb"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/packr"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
	"github.com/jinzhu/gorm"
	"github.com/manyminds/api2go/jsonapi"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
	"go.uber.org/zap/zapcore"
	null "gopkg.in/guregu/null.v3"
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
	Password = "password"
	// SessionSecret is the hardcoded secret solely used for test
	SessionSecret = "clsession_test_secret"
	// DefaultKey is the address of the fixture key
	DefaultKey = "0x3cb8e3FD9d27e39a5e9e6852b0e96160061fd4ea"
	// AllowUnstarted enable an application that can be used in tests without being started
	AllowUnstarted = "allow_unstarted"
)

var storeCounter uint64

var minimumContractPayment = assets.NewLink(100)

func init() {
	gin.SetMode(gin.TestMode)
	gomega.SetDefaultEventuallyTimeout(3 * time.Second)
	lvl := logLevelFromEnv()
	logger.SetLogger(CreateTestLogger(lvl))
	// Register txdb as dialect wrapping postgres
	// See: DialectTransactionWrappedPostgres
	config := orm.NewConfig()

	parsed, err := url.Parse(config.DatabaseURL())
	if err != nil || parsed.Path == "" {
		msg := fmt.Sprintf("invalid DATABASE_URL: `%s`. You must set DATABASE_URL env var to point to your test database. Note that the test database MUST end in `_test` to differentiate from a possible production DB. HINT: Try DATABASE_URL=postgresql://postgres@localhost:5432/chainlink_test?sslmode=disable", config.DatabaseURL())
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
	txdb.Register("cloudsqlpostgres", "postgres", config.DatabaseURL(), txdb.SavePointOption(nil))

	// Seed the random number generator, otherwise separate modules will take
	// the same advisory locks when tested with `go test -p N` for N > 1
	seed := time.Now().UTC().UnixNano()
	logger.Debugf("Using seed: %v", seed)
	rand.Seed(seed)
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

	wsserver, cleanup := newWSServer()
	return NewConfigWithWSServer(t, wsserver), cleanup
}

func NewRandomInt64() int64 {
	id := rand.Int63()
	return id
}

// NewTestConfig returns a test configuration
func NewTestConfig(t testing.TB, options ...interface{}) *TestConfig {
	t.Helper()

	count := atomic.AddUint64(&storeCounter, 1)
	rootdir := filepath.Join(RootDir, fmt.Sprintf("%d-%d", time.Now().UnixNano(), count))
	rawConfig := orm.NewConfig()

	rawConfig.Dialect = orm.DialectTransactionWrappedPostgres
	for _, opt := range options {
		switch v := opt.(type) {
		case orm.DialectName:
			rawConfig.Dialect = v
		}
	}

	uniqueRandomID := NewRandomInt64()
	// Unique advisory lock is required otherwise all tests will block each other
	rawConfig.AdvisoryLockID = uniqueRandomID

	rawConfig.Set("BRIDGE_RESPONSE_URL", "http://localhost:6688")
	rawConfig.Set("ETH_CHAIN_ID", 3)
	rawConfig.Set("CHAINLINK_DEV", true)
	rawConfig.Set("ETH_GAS_BUMP_THRESHOLD", 3)
	rawConfig.Set("MIGRATE_DATABASE", false)
	rawConfig.Set("MINIMUM_SERVICE_DURATION", "24h")
	rawConfig.Set("MIN_INCOMING_CONFIRMATIONS", 1)
	rawConfig.Set("MIN_OUTGOING_CONFIRMATIONS", 6)
	rawConfig.Set("MINIMUM_CONTRACT_PAYMENT", minimumContractPayment.Text(10))
	rawConfig.Set("ROOT", rootdir)
	rawConfig.Set("SESSION_TIMEOUT", "2m")
	rawConfig.SecretGenerator = mockSecretGenerator{}
	config := TestConfig{t: t, Config: rawConfig}
	return &config
}

// NewConfigWithWSServer return new config with specified wsserver
func NewConfigWithWSServer(t testing.TB, wsserver *httptest.Server) *TestConfig {
	t.Helper()

	config := NewTestConfig(t)
	config.SetEthereumServer(wsserver)
	return config
}

// SetEthereumServer sets the ethereum server for testconfig with given wsserver
func (tc *TestConfig) SetEthereumServer(wss *httptest.Server) {
	u, err := url.Parse(wss.URL)
	require.NoError(tc.t, err)
	u.Scheme = "ws"
	tc.Set("ETH_URL", u.String())
	tc.wsServer = wss
}

// TestApplication holds the test application and test servers
type TestApplication struct {
	t testing.TB
	*chainlink.ChainlinkApplication
	Config           *TestConfig
	Server           *httptest.Server
	wsServer         *httptest.Server
	connectedChannel chan struct{}
	Started          bool
	EthMock          *EthMock
	Backend          *backends.SimulatedBackend
	allowUnstarted   bool
}

func newWSServer() (*httptest.Server, func()) {
	return NewWSServer("")
}

// NewWSServer returns a  new wsserver
func NewWSServer(msg string) (*httptest.Server, func()) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		conn, err := upgrader.Upgrade(w, r, nil)
		logger.PanicIf(err)
		for {
			_, _, err = conn.ReadMessage()
			if err != nil {
				break
			}
			err = conn.WriteMessage(websocket.BinaryMessage, []byte(msg))
			if err != nil {
				break
			}
		}
	})
	server := httptest.NewServer(handler)
	return server, func() {
		server.Close()
	}
}

// NewApplication creates a New TestApplication along with a NewConfig
// It mocks the keystore with no keys or accounts by default
func NewApplication(t testing.TB, flags ...string) (*TestApplication, func()) {
	t.Helper()

	c, cfgCleanup := NewConfig(t)

	app, cleanup := NewApplicationWithConfig(t, c, flags...)
	kst := new(mocks.KeyStoreInterface)
	kst.On("Accounts").Return([]accounts.Account{})
	app.Store.KeyStore = kst

	return app, func() {
		cleanup()
		cfgCleanup()
	}
}

// NewApplicationWithKey creates a new TestApplication along with a new config
// It uses the native keystore and will load any keys that are in the database
func NewApplicationWithKey(t testing.TB, flags ...string) (*TestApplication, func()) {
	t.Helper()

	config, cfgCleanup := NewConfig(t)
	app, cleanup := NewApplicationWithConfigAndKey(t, config, flags...)
	return app, func() {
		cleanup()
		cfgCleanup()
	}
}

// NewApplicationWithConfigAndKey creates a new TestApplication with the given testconfig
// it will also provide an unlocked account on the keystore
func NewApplicationWithConfigAndKey(t testing.TB, tc *TestConfig, flags ...string) (*TestApplication, func()) {
	t.Helper()

	app, cleanup := NewApplicationWithConfig(t, tc, flags...)
	app.Store.KeyStore.Unlock(Password)

	return app, cleanup
}

// NewApplicationWithConfig creates a New TestApplication with specified test config
func NewApplicationWithConfig(t testing.TB, tc *TestConfig, flags ...string) (*TestApplication, func()) {
	t.Helper()

	ta := &TestApplication{t: t, connectedChannel: make(chan struct{}, 1)}
	app := chainlink.NewApplication(tc.Config, func(app chainlink.Application) {
		ta.connectedChannel <- struct{}{}
	}).(*chainlink.ChainlinkApplication)
	ta.ChainlinkApplication = app
	ta.EthMock = MockEthOnStore(t, app.Store, flags...)
	server := newServer(ta)
	tc.Config.Set("CLIENT_NODE_URL", server.URL)
	app.Store.Config = tc.Config

	for _, flag := range flags {
		if flag == AllowUnstarted {
			ta.allowUnstarted = true
		}
	}

	ta.Config = tc
	ta.Server = server
	ta.wsServer = tc.wsServer
	return ta, func() {
		require.NoError(t, ta.Stop())
		require.True(t, ta.EthMock.AllCalled(), ta.EthMock.Remaining())
	}
}

func NewApplicationWithConfigAndKeyOnSimulatedBlockchain(
	t testing.TB, tc *TestConfig, backend *backends.SimulatedBackend,
	flags ...string) (app *TestApplication, cleanup func()) {
	chainId := int(backend.Blockchain().Config().ChainID.Int64())
	tc.Config.Set("ETH_CHAIN_ID", chainId)
	app, appCleanup := NewApplicationWithConfigAndKey(t, tc, flags...)
	var client SimulatedBackendClient
	if txm, ok := app.Store.TxManager.(*strpkg.EthTxManager); ok {
		client = SimulatedBackendClient{b: backend, t: t, chainId: chainId}
		txm.Client = &client
	} else {
		log.Panic("SimulatedBackend only works on EthTxManager")
	}
	// Clean out the mock registrations, since we don't need those...
	app.EthMock.Responses = app.EthMock.Responses[:0]
	app.EthMock.Subscriptions = app.EthMock.Subscriptions[:0]
	return app, func() { appCleanup(); client.Close() }
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

	return ta.ChainlinkApplication.Start()
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
	require.NoError(ta.t, ta.ChainlinkApplication.Stop())
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
	_, err := ta.Store.KeyStore.Import([]byte(content), Password, Password)
	require.NoError(ta.t, err)
	require.NoError(ta.t, ta.Store.KeyStore.Unlock(Password))
}

func (ta *TestApplication) AddUnlockedKey() {
	_, err := ta.Store.KeyStore.NewAccount(Password)
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
		HTTP:                           cmd.NewAuthenticatedHTTPClient(ta.Config, cookieAuth),
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
func NewStoreWithConfig(config *TestConfig) (*strpkg.Store, func()) {
	s := strpkg.NewInsecureStore(config.Config, gracefulpanic.NewSignal())
	return s, func() {
		cleanUpStore(config.t, s)
	}
}

// NewStore creates a new store
func NewStore(t testing.TB) (*strpkg.Store, func()) {
	t.Helper()

	c, cleanup := NewConfig(t)
	store, storeCleanup := NewStoreWithConfig(c)
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

// FindJob returns JobSpec for given JobID
func FindJob(t testing.TB, s *strpkg.Store, id *models.ID) models.JobSpec {
	t.Helper()

	j, err := s.FindJob(id)
	require.NoError(t, err)

	return j
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
	headers[web.ExternalInitiatorAccessKeyHeader] = eia.AccessKey
	headers[web.ExternalInitiatorSecretHeader] = eia.Secret

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

	buffer := MustReadFile(t, "testdata/hello_world_job.json")

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
	bta *models.BridgeTypeAuthentication,
	body string,
) models.JobRun {
	t.Helper()

	client := app.NewHTTPClient()
	headers := map[string]string{"Authorization": "Bearer " + bta.IncomingToken}
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
) *models.BridgeTypeAuthentication {
	t.Helper()

	client := app.NewHTTPClient()
	resp, cleanup := client.Post(
		"/v2/bridge_types",
		bytes.NewBufferString(payload),
	)
	defer cleanup()
	AssertServerResponse(t, resp, http.StatusOK)
	bt := &models.BridgeTypeAuthentication{}
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
	DBWaitTimeout = 5 * time.Second
	// DBPollingInterval can't be too short to avoid DOSing the test database
	DBPollingInterval = 100 * time.Millisecond
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

// WaitForJobRunStatus waits for a JobRun to reach given status
func WaitForJobRunStatus(
	t testing.TB,
	store *strpkg.Store,
	jr models.JobRun,
	status models.RunStatus,
) models.JobRun {
	t.Helper()

	var err error
	gomega.NewGomegaWithT(t).Eventually(func() models.RunStatus {
		jr, err = store.Unscoped().FindJobRun(jr.ID)
		assert.NoError(t, err)
		return jr.GetStatus()
	}, DBWaitTimeout, DBPollingInterval).Should(gomega.Equal(status))
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
		}, DBWaitTimeout, DBPollingInterval).Should(gomega.HaveLen(want))
	} else {
		g.Eventually(func() []models.JobRun {
			jrs, err = store.JobRunsFor(j.ID)
			assert.NoError(t, err)
			return jrs
		}, DBWaitTimeout, DBPollingInterval).Should(gomega.HaveLen(want))
	}
	return jrs
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

func WaitForTxAttemptCount(t testing.TB, store *strpkg.Store, want int) []models.TxAttempt {
	t.Helper()
	g := gomega.NewGomegaWithT(t)

	var tas []models.TxAttempt
	var count int
	var err error
	if want == 0 {
		g.Consistently(func() int {
			tas, count, err = store.TxAttempts(0, 1000)
			assert.NoError(t, err)
			return count
		}, DBWaitTimeout, DBPollingInterval).Should(gomega.Equal(want))
	} else {
		g.Eventually(func() int {
			tas, count, err = store.TxAttempts(0, 1000)
			assert.NoError(t, err)
			return count
		}, DBWaitTimeout, DBPollingInterval).Should(gomega.Equal(want))
	}
	return tas
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
	}, DBWaitTimeout, DBPollingInterval).Should(gomega.Equal(want))
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
	return null.Time{Time: t, Valid: true}
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

// EmptyBlock returns a new empty ethereum block
func EmptyBlock() models.Block {
	return models.Block{}
}

// BlockWithTransactions returns a new ethereum block with transactions
// matching the given gas prices
func BlockWithTransactions(gasPrices ...uint64) models.Block {
	txs := make([]models.Transaction, len(gasPrices))
	for i, gasPrice := range gasPrices {
		txs[i].GasPrice = hexutil.Uint64(gasPrice)
	}
	return models.Block{
		Transactions: txs,
	}
}

// GetAccountAddress returns Address of the account in the keystore of the passed in store
func GetAccountAddress(t testing.TB, store *strpkg.Store) common.Address {
	t.Helper()

	account, err := store.KeyStore.GetFirstAccount()
	require.NoError(t, err)

	return account.Address
}

// GetAccountAddresses returns the Address of all registered accounts
func GetAccountAddresses(store *strpkg.Store) []common.Address {
	accounts := store.KeyStore.GetAccounts()

	addresses := []common.Address{}
	for _, account := range accounts {
		addresses = append(addresses, account.Address)
	}
	return addresses
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
	err := store.RawDB(func(db *gorm.DB) error {
		return db.Find(&all).Error
	})
	require.NoError(t, err)
	return all
}

func AllJobs(t testing.TB, store *strpkg.Store) []models.JobSpec {
	t.Helper()

	var all []models.JobSpec
	err := store.ORM.RawDB(func(db *gorm.DB) error {
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

func GetLastTxAttempt(t testing.TB, store *strpkg.Store) models.TxAttempt {
	t.Helper()

	var attempt models.TxAttempt
	var count int
	err := store.ORM.RawDB(func(db *gorm.DB) error {
		return db.Order("created_at desc").First(&attempt).Count(&count).Error
	})
	require.NoError(t, err)
	require.NotEqual(t, 0, count)
	return attempt
}

func GetLastTx(t testing.TB, store *strpkg.Store) models.Tx {
	t.Helper()

	var tx models.Tx
	var count int
	err := store.ORM.RawDB(func(db *gorm.DB) error {
		return db.Order("created_at desc").First(&tx).Count(&count).Error
	})
	require.NoError(t, err)
	require.NotEqual(t, 0, count)
	return tx
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

func MakeRoundStateReturnData(
	roundID uint64,
	eligible bool,
	answer, startAt, timeout, availableFunds, paymentAmount, oracleCount uint64,
) string {
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
	return hexutil.Encode(data)
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

// ChainlinkEthLogFromGethLog returns a copy of l as an models.Log. (They have
// identical fields, but the field tags differ, and the types differ slightly.)
func ChainlinkEthLogFromGethLog(l types.Log) models.Log {
	return models.Log{
		Address:     l.Address,
		Topics:      l.Topics,
		Data:        models.UntrustedBytes(l.Data),
		BlockNumber: l.BlockNumber,
		TxHash:      l.TxHash,
		TxIndex:     l.TxIndex,
		BlockHash:   l.BlockHash,
		Index:       l.Index,
		Removed:     l.Removed,
	}
}

func FindJobRun(t *testing.T, store *strpkg.Store, id *models.ID) models.JobRun {
	jr, err := store.FindJobRun(id)
	require.NoError(t, err)
	return jr
}

func MustHexToUint64(t *testing.T, hex string) uint64 {
	res, err := utils.HexToUint64(hex)
	if err != nil {
		t.Fatal(err)
	}
	return res
}

func MustDefaultKey(t *testing.T, s *strpkg.Store) models.Key {
	k, err := s.KeyByAddress(common.HexToAddress(DefaultKey))
	require.NoError(t, err)
	return k
}
