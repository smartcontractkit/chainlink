package cltest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/packr"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
	"github.com/jinzhu/gorm"
	"github.com/manyminds/api2go/jsonapi"
	"github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink/core/cmd"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services"
	strpkg "github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/assets"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/smartcontractkit/chainlink/core/store/presenters"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/web"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
	null "gopkg.in/guregu/null.v3"
)

const (
	// RootDir the root directory for cltest
	RootDir = "/tmp/chainlink_test"
	// Username the test username
	Username = "testusername"
	// APIEmail of the API user
	APIEmail = "email@test.net"
	// Password the password
	Password = "password"
	// APISessionID ID for API user
	APISessionID = "session"
	// SessionSecret is the hardcoded secret solely used for test
	SessionSecret = "clsession_test_secret"
)

var storeCounter uint64

var minimumContractPayment = assets.NewLink(100)

func init() {
	gin.SetMode(gin.TestMode)
	gomega.SetDefaultEventuallyTimeout(3 * time.Second)
	logger.SetLogger(logger.CreateTestLogger())
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

// NewConfigWithWSServer return new config with specified wsserver
func NewConfigWithWSServer(t testing.TB, wsserver *httptest.Server) *TestConfig {
	t.Helper()

	count := atomic.AddUint64(&storeCounter, 1)
	rootdir := filepath.Join(RootDir, fmt.Sprintf("%d-%d", time.Now().UnixNano(), count))
	rawConfig := orm.NewConfig()
	rawConfig.Set("BRIDGE_RESPONSE_URL", "http://localhost:6688")
	rawConfig.Set("ETH_CHAIN_ID", 3)
	rawConfig.Set("CHAINLINK_DEV", true)
	rawConfig.Set("ETH_GAS_BUMP_THRESHOLD", 3)
	rawConfig.Set("LOG_LEVEL", orm.LogLevel{Level: zapcore.DebugLevel})
	rawConfig.Set("LOG_SQL", false)
	rawConfig.Set("LOG_SQL_MIGRATIONS", false)
	rawConfig.Set("MINIMUM_SERVICE_DURATION", "24h")
	rawConfig.Set("MIN_INCOMING_CONFIRMATIONS", 1)
	rawConfig.Set("MIN_OUTGOING_CONFIRMATIONS", 6)
	rawConfig.Set("MINIMUM_CONTRACT_PAYMENT", minimumContractPayment.Text(10))
	rawConfig.Set("ROOT", rootdir)
	rawConfig.Set("SESSION_TIMEOUT", "2m")
	rawConfig.SecretGenerator = mockSecretGenerator{}
	config := TestConfig{t: t, Config: rawConfig}
	config.SetEthereumServer(wsserver)
	return &config
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
	*services.ChainlinkApplication
	Config           *orm.Config
	Server           *httptest.Server
	wsServer         *httptest.Server
	connectedChannel chan struct{}
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
func NewApplication(t testing.TB) (*TestApplication, func()) {
	t.Helper()

	c, cfgCleanup := NewConfig(t)
	app, cleanup := NewApplicationWithConfig(t, c)
	return app, func() {
		cleanup()
		cfgCleanup()
	}
}

// NewApplicationWithKey creates a new TestApplication along with a new config
func NewApplicationWithKey(t testing.TB) (*TestApplication, func()) {
	t.Helper()

	config, cfgCleanup := NewConfig(t)
	app, cleanup := NewApplicationWithConfigAndKey(t, config)
	return app, func() {
		cleanup()
		cfgCleanup()
	}
}

// NewApplicationWithConfigAndKey creates a new TestApplication with the given testconfig
// it will also provide an unlocked account on the keystore
func NewApplicationWithConfigAndKey(t testing.TB, tc *TestConfig) (*TestApplication, func()) {
	t.Helper()

	app, cleanup := NewApplicationWithConfig(t, tc)
	app.ImportKey(key3cb8e3fd9d27e39a5e9e6852b0e96160061fd4ea)
	return app, cleanup
}

// NewApplicationWithConfig creates a New TestApplication with specified test config
func NewApplicationWithConfig(t testing.TB, tc *TestConfig) (*TestApplication, func()) {
	t.Helper()

	WipePostgresDatabase(t, tc.Config)
	ta := &TestApplication{t: t, connectedChannel: make(chan struct{}, 1)}
	app := services.NewApplication(tc.Config, func(app services.Application) {
		ta.connectedChannel <- struct{}{}
	}).(*services.ChainlinkApplication)
	ta.ChainlinkApplication = app
	ethMock := MockEthOnStore(t, app.Store)

	server := newServer(ta)
	tc.Config.Set("CLIENT_NODE_URL", server.URL)
	app.Store.Config = tc.Config

	ta.Config = tc.Config
	ta.Server = server
	ta.wsServer = tc.wsServer
	return ta, func() {
		if !ethMock.AllCalled() {
			panic("mock expectations set and not used on default TestApplication ethMock!!!")
		}
		require.NoError(t, ta.Stop())
	}
}

func newServer(app services.Application) *httptest.Server {
	engine := web.Router(app)
	return httptest.NewServer(engine)
}

func (ta *TestApplication) NewBox() packr.Box {
	ta.t.Helper()

	return packr.NewBox("../fixtures/operator_ui/dist")
}

func (ta *TestApplication) StartAndConnect() error {
	ta.t.Helper()

	err := ta.Start()
	if err != nil {
		return err
	}

	return ta.WaitForConnection()
}

// WaitForConnection wait for the StartAndConnect callback to be called
func (ta *TestApplication) WaitForConnection() error {
	select {
	case <-time.After(4 * time.Second):
		return errors.New("TestApplication#StartAndConnect() timed out")
	case <-ta.connectedChannel:
		return nil
	}
}

func (ta *TestApplication) MockStartAndConnect() (*EthMock, error) {
	ethMock := ta.MockEthClient()
	ethMock.Context("TestApplication#MockStartAndConnect()", func(ethMock *EthMock) {
		ethMock.Register("eth_chainId", ta.Config.ChainID())
		ethMock.Register("eth_getTransactionCount", `0x0`)
	})

	err := ta.Start()
	if err != nil {
		return ethMock, err
	}

	return ethMock, ta.WaitForConnection()
}

// Stop will stop the test application and perform cleanup
func (ta *TestApplication) Stop() error {
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

func (ta *TestApplication) MustSeedUserSession() models.User {
	mockUser := MustUser(APIEmail, Password)
	require.NoError(ta.t, ta.Store.SaveUser(&mockUser))
	session := NewSession(APISessionID)
	require.NoError(ta.t, ta.Store.SaveSession(&session))
	return mockUser
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

	ta.MustSeedUserSession()
	return HTTPClientCleaner{
		HTTPClient: NewMockAuthenticatedHTTPClient(ta.Config),
		t:          ta.t,
	}
}

// NewClientAndRenderer creates a new cmd.Client for the test application
func (ta *TestApplication) NewClientAndRenderer() (*cmd.Client, *RendererMock) {
	ta.MustSeedUserSession()
	r := &RendererMock{}
	client := &cmd.Client{
		Renderer:                       r,
		Config:                         ta.Config,
		AppFactory:                     seededAppFactory{ta.ChainlinkApplication},
		KeyStoreAuthenticator:          CallbackAuthenticator{func(*strpkg.Store, string) (string, error) { return Password, nil }},
		FallbackAPIInitializer:         &MockAPIInitializer{},
		Runner:                         EmptyRunner{},
		HTTP:                           NewMockAuthenticatedHTTPClient(ta.Config),
		CookieAuthenticator:            MockCookieAuthenticator{},
		FileSessionRequestBuilder:      &MockSessionRequestBuilder{},
		PromptingSessionRequestBuilder: &MockSessionRequestBuilder{},
		ChangePasswordPrompter:         &MockChangePasswordPrompter{},
	}
	return client, r
}

func (ta *TestApplication) NewAuthenticatingClient(prompter cmd.Prompter) *cmd.Client {
	ta.MustSeedUserSession()
	cookieAuth := cmd.NewSessionCookieAuthenticator(ta.Config, &cmd.MemoryCookieStore{})
	client := &cmd.Client{
		Renderer:                       &RendererMock{},
		Config:                         ta.Config,
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

// NewStoreWithConfig creates a new store with given config
func NewStoreWithConfig(config *TestConfig) (*strpkg.Store, func()) {
	WipePostgresDatabase(config.t, config.Config)
	s := strpkg.NewStore(config.Config)
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

func WipePostgresDatabase(t testing.TB, config orm.ConfigReader) {
	t.Helper()

	if strings.HasPrefix(strings.ToLower(orm.NormalizedDatabaseURL(config)), string(orm.DialectPostgres)) {
		db, err := gorm.Open(string(orm.DialectPostgres), orm.NormalizedDatabaseURL(config))
		if err != nil {
			t.Fatalf("unable to open postgres database for wiping: %+v", err)
			return
		}
		defer db.Close()

		if err := db.Exec(`
DO $$ DECLARE
    r RECORD;
BEGIN
    FOR r IN (SELECT tablename FROM pg_tables WHERE schemaname = current_schema()) LOOP
        EXECUTE 'DROP TABLE IF EXISTS ' || quote_ident(r.tablename) || ' CASCADE';
    END LOOP;
END $$;
		`).Error; err != nil {
			t.Fatalf("unable to wipe postgres database: %+v", err)
		}
	}
}

// NewJobSubscriber creates a new JobSubscriber
func NewJobSubscriber(t testing.TB) (*strpkg.Store, services.JobSubscriber, func()) {
	t.Helper()

	store, cl := NewStore(t)
	nl := services.NewJobSubscriber(store)
	return store, nl, func() {
		cl()
	}
}

// CommonJSON has an ID, and Name
type CommonJSON struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Digest string `json:"digest"`
}

// ParseCommonJSON will unmarshall given body into CommonJSON
func ParseCommonJSON(t testing.TB, body io.Reader) CommonJSON {
	t.Helper()

	b, err := ioutil.ReadAll(body)
	require.NoError(t, err)
	var respJSON CommonJSON
	json.Unmarshal(b, &respJSON)
	return respJSON
}

func ParseJSON(t testing.TB, body io.Reader) models.JSON {
	t.Helper()

	b, err := ioutil.ReadAll(body)
	require.NoError(t, err)
	return models.JSON{Result: gjson.ParseBytes(b)}
}

// ErrorsJSON has an errors attribute
type ErrorsJSON struct {
	Errors []string `json:"errors"`
}

// ParseErrorsJSON will unmarshall given body into ErrorsJSON
func ParseErrorsJSON(t testing.TB, body io.Reader) ErrorsJSON {
	t.Helper()

	b, err := ioutil.ReadAll(body)
	require.NoError(t, err)
	var respJSON ErrorsJSON
	json.Unmarshal(b, &respJSON)
	return respJSON
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

func CopyFile(t testing.TB, src, dst string) {
	t.Helper()

	from, err := os.Open(src)
	require.NoError(t, err)
	defer from.Close()

	to, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE, 0666)
	require.NoError(t, err)

	_, err = io.Copy(to, from)
	require.NoError(t, err)
	require.NoError(t, to.Close())
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

// ObserveLogs returns the observed logs
func ObserveLogs() (*observer.ObservedLogs, func()) {
	previousLogger := logger.GetLogger()
	core, observed := observer.New(zapcore.DebugLevel)
	logger.SetLogger(zap.New(core))
	return observed, func() {
		logger.SetLogger(previousLogger.Desugar())
	}
}

// ReadLogs returns the contents of the applications log file as a string
func ReadLogs(app *TestApplication) (string, error) {
	logFile := fmt.Sprintf("%s/log.jsonl", app.Store.Config.RootDir())
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

// FindJobRun returns JobRun for given JobRunID
func FindJobRun(t testing.TB, s *strpkg.Store, id *models.ID) models.JobRun {
	t.Helper()

	j, err := s.FindJobRun(id)
	require.NoError(t, err)

	return j
}

func FindExternalInitiator(t testing.TB, s *strpkg.Store, eia *models.ExternalInitiatorAuthentication) *models.ExternalInitiator {
	t.Helper()

	ei, err := s.FindExternalInitiator(eia)
	require.NoError(t, err)

	return ei
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
	AssertServerResponse(t, resp, 200)

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

	job.Tasks[0].Params, err = job.Tasks[0].Params.Merge(JSONFromString(t, `{"get":"%v"}`, url))
	assert.NoError(t, err)
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

	AssertServerResponse(t, resp, 200)
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
	AssertServerResponse(t, resp, 200)
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
) *models.ExternalInitiatorAuthentication {
	t.Helper()

	client := app.NewHTTPClient()
	resp, cleanup := client.Post(
		"/v2/external_initiators",
		bytes.NewBufferString(payload),
	)
	defer cleanup()
	AssertServerResponse(t, resp, 201)
	eia := &models.ExternalInitiatorAuthentication{}
	err := ParseJSONAPIResponse(t, resp, eia)
	require.NoError(t, err)

	return eia
}

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

// WaitForJobRunToPendConfirmations waits for a JobRun to reach PendingConfirmations Status
func WaitForJobRunToPendConfirmations(
	t testing.TB,
	store *strpkg.Store,
	jr models.JobRun,
) models.JobRun {
	t.Helper()

	return WaitForJobRunStatus(t, store, jr, models.RunStatusPendingConfirmations)
}

// WaitForJobRunToPendSleep waits for a JobRun to reach PendingBridge Status
func WaitForJobRunToPendSleep(
	t testing.TB,
	store *strpkg.Store,
	jr models.JobRun,
) models.JobRun {
	t.Helper()

	return WaitForJobRunStatus(t, store, jr, models.RunStatusPendingSleep)
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
		return jr.Status
	}).Should(gomega.Equal(status))
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
		return jr.Status
	}, duration).Should(gomega.Equal(status))
	return jr
}

// JobRunStaysPendingConfirmations tests if a JobRun will stay at the PendingConfirmations Status
func JobRunStaysPendingConfirmations(
	t testing.TB,
	store *strpkg.Store,
	jr models.JobRun,
) models.JobRun {
	t.Helper()

	return JobRunStays(t, store, jr, models.RunStatusPendingConfirmations)
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
		}).Should(gomega.HaveLen(want))
	} else {
		g.Eventually(func() []models.JobRun {
			jrs, err = store.JobRunsFor(j.ID)
			assert.NoError(t, err)
			return jrs
		}).Should(gomega.HaveLen(want))
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
		}).Should(gomega.BeNumerically(">=", want))
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
		}).Should(gomega.Equal(want))
	} else {
		g.Eventually(func() int {
			tas, count, err = store.TxAttempts(0, 1000)
			assert.NoError(t, err)
			return count
		}).Should(gomega.Equal(want))
	}
	return tas
}

// WaitForJobs waits for the wanted number of jobs.
func WaitForJobs(t testing.TB, store *strpkg.Store, want int) []models.JobSpec {
	t.Helper()
	g := gomega.NewGomegaWithT(t)

	var jobs []models.JobSpec
	if want == 0 {
		g.Consistently(func() []models.JobSpec {
			jobs = AllJobs(t, store)
			return jobs
		}).Should(gomega.HaveLen(want))
	} else {
		g.Eventually(func() []models.JobSpec {
			jobs = AllJobs(t, store)
			return jobs
		}).Should(gomega.HaveLen(want))
	}
	return jobs
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
		var count int
		assert.NoError(t, orm.DB.Model(&models.SyncEvent{}).Count(&count).Error)
		return count
	}).Should(gomega.Equal(want))
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
		var count int
		assert.NoError(t, orm.DB.Model(&models.SyncEvent{}).Count(&count).Error)
		return count
	}).Should(gomega.Equal(want))
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
	switch val.(type) {
	case int:
		return models.NewHead(big.NewInt(int64(val.(int))), NewHash())
	case uint64:
		return models.NewHead(big.NewInt(int64(val.(uint64))), NewHash())
	case int64:
		return models.NewHead(big.NewInt(val.(int64)), NewHash())
	case *big.Int:
		return models.NewHead(val.(*big.Int), NewHash())
	default:
		logger.Panicf("Could not convert %v of type %T to Head", val, val)
		return nil
	}
}

// NewBlockHeader return a new BlockHeader with given number
func NewBlockHeader(number int) *models.BlockHeader {
	return &models.BlockHeader{Number: BigHexInt(number)}
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

func hasHexPrefix(str string) bool {
	return len(str) >= 2 && str[0] == '0' && (str[1] == 'x' || str[1] == 'X')
}

func isHexCharacter(c byte) bool {
	return ('0' <= c && c <= '9') || ('a' <= c && c <= 'f') || ('A' <= c && c <= 'F')
}

func isHex(str string) bool {
	if len(str)%2 != 0 {
		return false
	}
	for _, c := range []byte(str) {
		if !isHexCharacter(c) {
			return false
		}
	}
	return true
}

// AssertValidHash checks that a string matches a specific hash format,
// includes a leading 0x and has a specific length (in bytes)
func AssertValidHash(t testing.TB, length int, hash string) {
	t.Helper()

	if !hasHexPrefix(hash) {
		assert.FailNowf(t, "Missing hash prefix", `"%+v" is missing hash prefix`, hash)
	}
	hash = hash[2:]
	hashlen := len(hash) / 2
	if hashlen != length {
		assert.FailNowf(t, "Wrong hash length", `"%+v" represents %d bytes, want %d`, hash, hashlen, length)
	}
	if !isHex(hash) {
		assert.FailNowf(t, "Invalid character", `"%+v" contains a non hexadecimal character`, hash)
	}
}

// AssertServerResponse is used to match against a client response, will print
// any errors returned if the request fails.
func AssertServerResponse(t testing.TB, resp *http.Response, expectedStatusCode int) {
	t.Helper()

	if resp.StatusCode == expectedStatusCode {
		return
	}

	t.Logf("expected status code %d got %d", expectedStatusCode, resp.StatusCode)

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
	err := store.ORM.DB.Find(&all).Error
	require.NoError(t, err)
	return all
}

func AllJobs(t testing.TB, store *strpkg.Store) []models.JobSpec {
	t.Helper()

	var all []models.JobSpec
	err := store.ORM.DB.Find(&all).Error
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
	err := store.ORM.DB.Order("created_at desc").First(&attempt).Count(&count).Error
	require.NoError(t, err)
	require.NotEqual(t, 0, count)
	return attempt
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

func MustSha256(in string) string {
	out, _ := utils.Sha256(in)
	return out
}

func MustParseURL(input string) *url.URL {
	u, err := url.Parse(input)
	if err != nil {
		logger.Panic(err)
	}
	return u
}
