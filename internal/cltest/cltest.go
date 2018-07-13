package cltest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink/adapters"
	"github.com/smartcontractkit/chainlink/cmd"
	"github.com/smartcontractkit/chainlink/logger"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/utils"
	"github.com/smartcontractkit/chainlink/web"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
	null "gopkg.in/guregu/null.v3"
)

// RootDir the root directory for cltest
const RootDir = "/tmp/chainlink_test"

// Username the test username
const Username = "testusername"

// Password the password
const Password = "password"

var storeCounter uint64

func init() {
	gin.SetMode(gin.TestMode)
	gomega.SetDefaultEventuallyTimeout(3 * time.Second)
	logger.SetLogger(logger.CreateTestLogger())
}

// TestConfig struct with test store and wsServer
type TestConfig struct {
	store.Config
	wsServer *httptest.Server
}

// NewConfig returns a new TestConfig
func NewConfig() (*TestConfig, func()) {
	wsserver, cleanup := newWSServer()
	return NewConfigWithWSServer(wsserver), cleanup
}

// NewConfigWithWSServer return new config with specified wsserver
func NewConfigWithWSServer(wsserver *httptest.Server) *TestConfig {
	count := atomic.AddUint64(&storeCounter, 1)
	rootdir := path.Join(RootDir, fmt.Sprintf("%d-%d", time.Now().UnixNano(), count))
	config := TestConfig{
		Config: store.Config{
			LogLevel:                 store.LogLevel{Level: zapcore.DebugLevel},
			RootDir:                  rootdir,
			BasicAuthUsername:        Username,
			BasicAuthPassword:        Password,
			ChainID:                  3,
			MinOutgoingConfirmations: 6,
			MinIncomingConfirmations: 0,
			EthGasBumpWei:            *big.NewInt(5000000000),
			EthGasBumpThreshold:      3,
			EthGasPriceDefault:       *big.NewInt(20000000000),
			DatabasePollInterval:     store.Duration{Duration: time.Millisecond * 500},
			AllowOrigins:             "http://localhost:3000,http://localhost:6689",
			Dev:                      true,
		},
	}
	config.SetEthereumServer(wsserver)
	return &config
}

// SetEthereumServer sets the ethereum server for testconfig with given wsserver
func (tc *TestConfig) SetEthereumServer(wss *httptest.Server) {
	u, err := url.Parse(wss.URL)
	mustNotErr(err)
	u.Scheme = "ws"
	tc.EthereumURL = u.String()
	tc.wsServer = wss
}

// TestApplication holds the test application and test servers
type TestApplication struct {
	*services.ChainlinkApplication
	Server   *httptest.Server
	wsServer *httptest.Server
}

func newWSServer() (*httptest.Server, func()) {
	return NewWSServer("")
}

// NewWSServer returns a  new wsserver
func NewWSServer(msg string) (*httptest.Server, func()) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	var conn *websocket.Conn
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		conn, err = upgrader.Upgrade(w, r, nil)
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
		if conn != nil {
			mustNotErr(conn.Close())
		}
		server.Close()
	}
}

// NewApplication creates a New TestApplication along with a NewConfig
func NewApplication() (*TestApplication, func()) {
	c, cfgCleanup := NewConfig()
	app, cleanup := NewApplicationWithConfig(c)
	return app, func() {
		cleanup()
		cfgCleanup()
	}
}

// NewApplicationWithConfig creates a New TestApplication with specified test config
func NewApplicationWithConfig(tc *TestConfig) (*TestApplication, func()) {
	app := services.NewApplication(tc.Config).(*services.ChainlinkApplication)
	server := newServer(app)
	tc.Config.ClientNodeURL = server.URL
	app.Store.Config = tc.Config
	ethMock := MockEthOnStore(app.Store)
	ta := &TestApplication{
		ChainlinkApplication: app,
		Server:               server,
		wsServer:             tc.wsServer,
	}
	return ta, func() {
		if !ethMock.AllCalled() {
			panic("mock expectations set and not used on default TestApplication ethMock!!!")
		}
		ta.Stop()
	}
}

// NewApplicationWithKeyStore creates a new TestApplication along with a new config
func NewApplicationWithKeyStore() (*TestApplication, func()) {
	config, cfgCleanup := NewConfig()
	app, cleanup := NewApplicationWithConfigAndKeyStore(config)
	return app, func() {
		cleanup()
		cfgCleanup()
	}
}

// NewApplicationWithConfigAndKeyStore creates a new TestApplication with the given testconfig
// it will also provide an unlocked account on the keystore
func NewApplicationWithConfigAndKeyStore(tc *TestConfig) (*TestApplication, func()) {
	app, cleanup := NewApplicationWithConfig(tc)
	_, err := app.Store.KeyStore.NewAccount(Password)
	mustNotErr(err)
	mustNotErr(app.Store.KeyStore.Unlock(Password))
	return app, cleanup
}

func newServer(app *services.ChainlinkApplication) *httptest.Server {
	engine := web.Router(app)
	return httptest.NewServer(engine)
}

// Stop will stop the test application and perform cleanup
func (ta *TestApplication) Stop() error {
	ta.ChainlinkApplication.Stop()
	cleanUpStore(ta.Store)
	if ta.Server != nil {
		ta.Server.Close()
	}
	if ta.wsServer != nil {
		ta.wsServer.Close()
	}
	return nil
}

// NewStoreWithConfig creates a new store with given config
func NewStoreWithConfig(config *TestConfig) (*store.Store, func()) {
	s := store.NewStore(config.Config)
	return s, func() {
		cleanUpStore(s)
	}
}

// NewStore creates a new store
func NewStore() (*store.Store, func()) {
	c, cleanup := NewConfig()
	store, storeCleanup := NewStoreWithConfig(c)
	return store, func() {
		storeCleanup()
		cleanup()
	}
}

func cleanUpStore(store *store.Store) {
	logger.Sync()
	store.Close()
	go func() {
		if err := os.RemoveAll(store.Config.RootDir); err != nil {
			log.Println(err)
		}
	}()
}

// NewJobSubscriber creates a new JobSubscriber
func NewJobSubscriber() (*store.Store, services.JobSubscriber, func()) {
	store, cl := NewStore()
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
func ParseCommonJSON(body io.Reader) CommonJSON {
	b, err := ioutil.ReadAll(body)
	mustNotErr(err)
	var respJSON CommonJSON
	json.Unmarshal(b, &respJSON)
	return respJSON
}

// ErrorsJSON has an errors attribute
type ErrorsJSON struct {
	Errors []string `json:"errors"`
}

// ParseErrorsJSON will unmarshall given body into ErrorsJSON
func ParseErrorsJSON(body io.Reader) ErrorsJSON {
	b, err := ioutil.ReadAll(body)
	mustNotErr(err)
	var respJSON ErrorsJSON
	json.Unmarshal(b, &respJSON)
	return respJSON
}

// LoadJSON loads json from file and returns a byte slice
func LoadJSON(file string) []byte {
	content, err := ioutil.ReadFile(file)
	mustNotErr(err)
	return content
}

func copyFile(src, dst string) {
	from, err := os.Open(src)
	mustNotErr(err)
	defer from.Close()

	to, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE, 0666)
	mustNotErr(err)

	_, err = io.Copy(to, from)
	mustNotErr(err)
	mustNotErr(to.Close())
}

// AddPrivateKey adds private key from src to config
func AddPrivateKey(config *TestConfig, src string) {
	err := os.MkdirAll(config.KeysDir(), os.FileMode(0700))
	mustNotErr(err)

	dst := config.KeysDir() + "/testwallet.json"
	copyFile(src, dst)
}

// BasicAuthPost performs a POST request to the given url with specified contentType and body
// and returns the Response
func BasicAuthPost(url string, contentType string, body io.Reader) (*http.Response, func()) {
	resp, err := utils.BasicAuthPost(
		Username,
		Password,
		url,
		contentType,
		body)
	mustNotErr(err)
	return resp, func() { mustNotErr(resp.Body.Close()) }
}

// BasicAuthGet performs a GET request to given url and returns the Response
func BasicAuthGet(url string, headers ...map[string]string) (*http.Response, func()) {
	resp, err := utils.BasicAuthGet(Username, Password, url, headers...)
	mustNotErr(err)
	return resp, func() { mustNotErr(resp.Body.Close()) }
}

// BasicAuthPatch performs a PATCH request to the given url with specified contentType and body
// and returns the Response
func BasicAuthPatch(url string, contentType string, body io.Reader) (*http.Response, func()) {
	resp, err := utils.BasicAuthPatch(
		Username,
		Password,
		url,
		contentType,
		body)
	mustNotErr(err)
	return resp, func() { mustNotErr(resp.Body.Close()) }
}

// BasicAuthDelete performs a DELETE request to the given url with specified contentType and body
// and returns the Response
func BasicAuthDelete(url string, contentType string, body io.Reader) (*http.Response, func()) {
	resp, err := utils.BasicAuthDelete(
		Username,
		Password,
		url,
		contentType,
		body)
	mustNotErr(err)
	return resp, func() { mustNotErr(resp.Body.Close()) }
}

// ParseResponseBody will parse the given response into a byte slice
func ParseResponseBody(resp *http.Response) []byte {
	b, err := ioutil.ReadAll(resp.Body)
	mustNotErr(err)
	return b
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
func ObserveLogs() *observer.ObservedLogs {
	core, observed := observer.New(zapcore.DebugLevel)
	logger.SetLogger(zap.New(core))
	return observed
}

// ReadLogs returns the contents of the applications log file as a string
func ReadLogs(app *TestApplication) (string, error) {
	logFile := fmt.Sprintf("%s/log.jsonl", app.Store.Config.RootDir)
	b, err := ioutil.ReadFile(logFile)
	return string(b), err
}

// FixtureCreateJobViaWeb creates a job from a fixture using /v2/specs
func FixtureCreateJobViaWeb(t *testing.T, app *TestApplication, path string) models.JobSpec {
	resp, cleanup := BasicAuthPost(
		app.Server.URL+"/v2/specs",
		"application/json",
		bytes.NewBuffer(LoadJSON(path)),
	)
	defer cleanup()
	AssertServerResponse(t, resp, 200)

	return FindJob(app.Store, ParseCommonJSON(resp.Body).ID)
}

// FindJob returns JobSpec for given JobID
func FindJob(s *store.Store, id string) models.JobSpec {
	j, err := s.FindJob(id)
	mustNotErr(err)

	return j
}

// FindJobRun returns JobRun for given JobRunID
func FindJobRun(s *store.Store, id string) models.JobRun {
	j, err := s.FindJobRun(id)
	mustNotErr(err)

	return j
}

// FixtureCreateJobWithAssignmentViaWeb creates a job from a fixture using /v1/assignments
func FixtureCreateJobWithAssignmentViaWeb(t *testing.T, app *TestApplication, path string) models.JobSpec {
	resp, cleanup := BasicAuthPost(
		app.Server.URL+"/v1/assignments",
		"application/json",
		bytes.NewBuffer(LoadJSON(path)),
	)
	defer cleanup()
	AssertServerResponse(t, resp, 200)
	return FindJob(app.Store, ParseCommonJSON(resp.Body).ID)
}

// CreateJobSpecViaWeb creates a jobspec via web using /v2/specs
func CreateJobSpecViaWeb(t *testing.T, app *TestApplication, job models.JobSpec) models.JobSpec {
	marshaled, err := json.Marshal(&job)
	assert.NoError(t, err)
	resp, cleanup := BasicAuthPost(
		app.Server.URL+"/v2/specs",
		"application/json",
		bytes.NewBuffer(marshaled),
	)
	defer cleanup()
	AssertServerResponse(t, resp, 200)
	return FindJob(app.Store, ParseCommonJSON(resp.Body).ID)
}

// CreateJobRunViaWeb creates JobRun via web using /v2/specs/ID/runs
func CreateJobRunViaWeb(t *testing.T, app *TestApplication, j models.JobSpec, body ...string) models.JobRun {
	t.Helper()
	url := app.Server.URL + "/v2/specs/" + j.ID + "/runs"
	bodyBuffer := &bytes.Buffer{}
	if len(body) > 0 {
		bodyBuffer = bytes.NewBufferString(body[0])
	}
	resp, cleanup := BasicAuthPost(url, "application/json", bodyBuffer)
	defer cleanup()
	AssertServerResponse(t, resp, 200)
	jrID := ParseCommonJSON(resp.Body).ID

	jrs := []models.JobRun{}
	gomega.NewGomegaWithT(t).Eventually(func() []models.JobRun {
		assert.Nil(t, app.Store.Where("ID", jrID, &jrs))
		return jrs
	}).Should(gomega.HaveLen(1))
	jr := jrs[0]
	assert.Equal(t, j.ID, jr.JobID)

	return jr
}

// CreateHelloWorldJobViaWeb creates a HelloWorld JobSpec with the given MockServer Url
func CreateHelloWorldJobViaWeb(t *testing.T, app *TestApplication, url string) models.JobSpec {
	j := FixtureCreateJobViaWeb(t, app, "../internal/fixtures/web/hello_world_job.json")
	j.Tasks[0] = NewTask("httpget", fmt.Sprintf(`{"url":"%v"}`, url))
	return CreateJobSpecViaWeb(t, app, j)
}

// CreateMockAssignmentViaWeb creates a JobSpec with the given MockServer Url
func CreateMockAssignmentViaWeb(t *testing.T, app *TestApplication, url string) models.JobSpec {

	j := FixtureCreateJobWithAssignmentViaWeb(t, app, "../internal/fixtures/web/v1_format_job.json")
	j.Tasks[0] = NewTask("httpget", fmt.Sprintf(`{"url":"%v"}`, url))
	return CreateJobSpecViaWeb(t, app, j)
}

// UpdateJobRunViaWeb updates jobrun via web using /v2/runs/ID
func UpdateJobRunViaWeb(
	t *testing.T,
	app *TestApplication,
	jr models.JobRun,
	body string,
) models.JobRun {
	t.Helper()
	url := app.Server.URL + "/v2/runs/" + jr.ID
	resp, cleanup := BasicAuthPatch(url, "application/json", bytes.NewBufferString(body))
	defer cleanup()

	AssertServerResponse(t, resp, 200)
	jrID := ParseCommonJSON(resp.Body).ID
	assert.Nil(t, app.Store.One("ID", jrID, &jr))
	return jr
}

// CreateBridgeTypeViaWeb creates a bridgetype via web using /v2/bridge_types
func CreateBridgeTypeViaWeb(
	t *testing.T,
	app *TestApplication,
	payload string,
) models.BridgeType {
	resp, cleanup := BasicAuthPost(
		app.Server.URL+"/v2/bridge_types",
		"application/json",
		bytes.NewBufferString(payload),
	)
	defer cleanup()
	AssertServerResponse(t, resp, 200)
	name := ParseCommonJSON(resp.Body).Name
	bt, err := app.Store.FindBridge(name)
	assert.NoError(t, err)

	return bt
}

// NewClientAndRenderer creates a new cmd.Client with given config
func NewClientAndRenderer(config store.Config) (*cmd.Client, *RendererMock) {
	r := &RendererMock{}
	client := &cmd.Client{
		Renderer:        r,
		Config:          config,
		AppFactory:      EmptyAppFactory{},
		Auth:            CallbackAuthenticator{func(*store.Store, string) error { return nil }},
		UserInitializer: cltest.MockUserInitializer{},
		Runner:          EmptyRunner{},
	}
	return client, r
}

// WaitForJobRunToComplete waits for a JobRun to reach Completed Status
func WaitForJobRunToComplete(
	t *testing.T,
	store *store.Store,
	jr models.JobRun,
) models.JobRun {
	return WaitForJobRunStatus(t, store, jr, models.RunStatusCompleted)
}

// WaitForJobRunToPendBridge waits for a JobRun to reach PendingBridge Status
func WaitForJobRunToPendBridge(
	t *testing.T,
	store *store.Store,
	jr models.JobRun,
) models.JobRun {
	return WaitForJobRunStatus(t, store, jr, models.RunStatusPendingBridge)
}

// WaitForJobRunToPendConfirmations waits for a JobRun to reach PendingConfirmations Status
func WaitForJobRunToPendConfirmations(
	t *testing.T,
	store *store.Store,
	jr models.JobRun,
) models.JobRun {
	return WaitForJobRunStatus(t, store, jr, models.RunStatusPendingConfirmations)
}

// WaitForJobRunStatus waits for a JobRun to reach given status
func WaitForJobRunStatus(
	t *testing.T,
	store *store.Store,
	jr models.JobRun,
	status models.RunStatus,
) models.JobRun {
	t.Helper()
	gomega.NewGomegaWithT(t).Eventually(func() models.RunStatus {
		assert.Nil(t, store.One("ID", jr.ID, &jr))
		return jr.Status
	}).Should(gomega.Equal(status))
	return jr
}

// JobRunStays tests if a JobRun will consistently stay at the specified status
func JobRunStays(
	t *testing.T,
	store *store.Store,
	jr models.JobRun,
	status models.RunStatus,
) models.JobRun {
	t.Helper()
	gomega.NewGomegaWithT(t).Consistently(func() models.RunStatus {
		assert.Nil(t, store.One("ID", jr.ID, &jr))
		return jr.Status
	}).Should(gomega.Equal(status))
	return jr
}

// JobRunStaysPendingConfirmations tests if a JobRun will stay at the PendingConfirmations Status
func JobRunStaysPendingConfirmations(
	t *testing.T,
	store *store.Store,
	jr models.JobRun,
) models.JobRun {
	return JobRunStays(t, store, jr, models.RunStatusPendingConfirmations)
}

// WaitForRuns waits for the wanted number of runs then returns a slice of the JobRuns
func WaitForRuns(t *testing.T, j models.JobSpec, store *store.Store, want int) []models.JobRun {
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

// WaitForJobs waits for the wanted number of jobs.
func WaitForJobs(t *testing.T, store *store.Store, want int) []models.JobSpec {
	t.Helper()
	g := gomega.NewGomegaWithT(t)

	var jobs []models.JobSpec
	var err error
	if want == 0 {
		g.Consistently(func() []models.JobSpec {
			jobs, err = store.Jobs()
			assert.NoError(t, err)
			return jobs
		}).Should(gomega.HaveLen(want))
	} else {
		g.Eventually(func() []models.JobSpec {
			jobs, err = store.Jobs()
			assert.NoError(t, err)
			return jobs
		}).Should(gomega.HaveLen(want))
	}
	return jobs
}

// MustParseWebURL must parse the given url and return it
func MustParseWebURL(str string) models.WebURL {
	u, err := url.Parse(str)
	mustNotErr(err)
	return models.WebURL{URL: u}
}

// ParseISO8601 given the time string it Must parse the time and return it
func ParseISO8601(s string) time.Time {
	t, err := time.Parse(time.RFC3339Nano, s)
	mustNotErr(err)
	return t
}

// NullableTime will return a valid nullable time given time.Time
func NullableTime(t time.Time) null.Time {
	return null.Time{Time: t, Valid: true}
}

// ParseNullableTime given a time string parse it into a null.Time
func ParseNullableTime(s string) null.Time {
	return NullableTime(ParseISO8601(s))
}

// IndexableBlockNumber given the value convert it into an IndexableBlockNumber
func IndexableBlockNumber(val interface{}) *models.IndexableBlockNumber {
	switch val.(type) {
	case int:
		return models.NewIndexableBlockNumber(big.NewInt(int64(val.(int))), NewHash())
	case uint64:
		return models.NewIndexableBlockNumber(big.NewInt(int64(val.(uint64))), NewHash())
	case int64:
		return models.NewIndexableBlockNumber(big.NewInt(val.(int64)), NewHash())
	case *big.Int:
		return models.NewIndexableBlockNumber(val.(*big.Int), NewHash())
	default:
		logger.Panicf("Could not convert %v of type %T to IndexableBlockNumber", val, val)
		return nil
	}
}

// NewBlockHeader return a new BlockHeader with given number
func NewBlockHeader(number int) *models.BlockHeader {
	return &models.BlockHeader{Number: BigHexInt(number)}
}

func mustNotErr(err error) {
	if err != nil {
		logger.Panic(err)
	}
}

// UnwrapAdapter unwraps the adapter from given wrapped adapter
func UnwrapAdapter(wa adapters.AdapterWithMinConfs) adapters.Adapter {
	return wa.(adapters.MinConfsWrappedAdapter).Adapter
}

// GetAccountAddress returns Address of the account in the keystore of the passed in store
func GetAccountAddress(store *store.Store) common.Address {
	account, err := store.KeyStore.GetAccount()
	mustNotErr(err)

	return account.Address
}

// StringToHash give Keccak256 hash of string
func StringToHash(s string) common.Hash {
	return common.BytesToHash([]byte(s))
}

// AssertServerResponse is used to match against a client response, will print
// any errors returned if the request fails.
func AssertServerResponse(t *testing.T, resp *http.Response, expectedStatusCode int) {
	if resp.StatusCode == expectedStatusCode {
		return
	}

	if resp.StatusCode >= 300 && resp.StatusCode < 600 {
		var result map[string][]string
		err := json.Unmarshal(ParseResponseBody(resp), &result)
		mustNotErr(err)

		assert.FailNowf(t, "Request failed", "Expected %d response, got %d with errors: %s", expectedStatusCode, resp.StatusCode, result["errors"])
	} else {
		assert.FailNowf(t, "Unexpected response", "Expected %d response, got %d", expectedStatusCode, resp.StatusCode)
	}
}
