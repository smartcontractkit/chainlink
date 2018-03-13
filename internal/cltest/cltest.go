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
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/h2non/gock"
	"github.com/onsi/gomega"
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

const RootDir = "/tmp/chainlink_test"
const Username = "testusername"
const Password = "password"

var storeCounter uint64 = 0

func init() {
	gin.SetMode(gin.TestMode)
	gomega.SetDefaultEventuallyTimeout(3 * time.Second)
}

type TestConfig struct {
	store.Config
	wsServer *httptest.Server
}

func NewConfig() (*TestConfig, func()) {
	wsserver := newWSServer()
	return NewConfigWithWSServer(wsserver), func() { wsserver.Close() }
}

func NewConfigWithWSServer(wsserver *httptest.Server) *TestConfig {
	count := atomic.AddUint64(&storeCounter, 1)
	rootdir := path.Join(RootDir, fmt.Sprintf("%d-%d", time.Now().UnixNano(), count))
	config := TestConfig{
		Config: store.Config{
			LogLevel:            store.LogLevel{zapcore.DebugLevel},
			RootDir:             rootdir,
			BasicAuthUsername:   Username,
			BasicAuthPassword:   Password,
			ChainID:             3,
			EthMinConfirmations: 6,
			EthGasBumpWei:       *big.NewInt(5000000000),
			EthGasBumpThreshold: 3,
			EthGasPriceDefault:  *big.NewInt(20000000000),
		},
	}
	config.SetEthereumServer(wsserver)
	return &config
}

func (tc *TestConfig) SetEthereumServer(wss *httptest.Server) {
	u, _ := url.Parse(wss.URL)
	u.Scheme = "ws"
	tc.EthereumURL = u.String()
	tc.wsServer = wss
}

type TestApplication struct {
	*services.ChainlinkApplication
	Server   *httptest.Server
	wsServer *httptest.Server
}

func newWSServer() *httptest.Server {
	return NewWSServer("")
}

func NewWSServer(msg string) *httptest.Server {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, _ := upgrader.Upgrade(w, r, nil)
		conn.WriteMessage(websocket.BinaryMessage, []byte(msg))
	})
	server := httptest.NewServer(handler)
	return server
}

func NewApplication() (*TestApplication, func()) {
	c, _ := NewConfig()
	return NewApplicationWithConfig(c)
}

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

func NewApplicationWithKeyStore() (*TestApplication, func()) {
	app, cleanup := NewApplication()
	_, err := app.Store.KeyStore.NewAccount(Password)
	mustNotErr(err)
	mustNotErr(app.Store.KeyStore.Unlock(Password))
	return app, cleanup
}

func newServer(app *services.ChainlinkApplication) *httptest.Server {
	return httptest.NewServer(web.Router(app))
}

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

func NewStoreWithConfig(config *TestConfig) (*store.Store, func()) {
	s := store.NewStore(config.Config)
	return s, func() {
		cleanUpStore(s)
		if config.wsServer != nil {
			config.wsServer.Close()
		}
	}
}

func NewStore() (*store.Store, func()) {
	c, _ := NewConfig()
	return NewStoreWithConfig(c)
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

func NewEthereumListener() (*services.EthereumListener, func()) {
	store, cl := NewStore()
	nl := &services.EthereumListener{Store: store}
	return nl, func() {
		cl()
	}
}

func CloseGock(t *testing.T) {
	assert.True(t, gock.IsDone(), "Not all gock requests were fulfilled")
	gock.DisableNetworking()
	gock.Off()
}

type CommonJSON struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func ParseCommonJSON(body io.Reader) CommonJSON {
	b, err := ioutil.ReadAll(body)
	mustNotErr(err)
	var respJSON CommonJSON
	json.Unmarshal(b, &respJSON)
	return respJSON
}

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

func AddPrivateKey(config *TestConfig, src string) {
	err := os.MkdirAll(config.KeysDir(), os.FileMode(0700))
	mustNotErr(err)

	dst := config.KeysDir() + "/testwallet.json"
	copyFile(src, dst)
}

func BasicAuthPost(url string, contentType string, body io.Reader) *http.Response {
	resp, err := utils.BasicAuthPost(
		Username,
		Password,
		url,
		contentType,
		body)
	mustNotErr(err)
	return resp
}

func BasicAuthGet(url string) *http.Response {
	resp, err := utils.BasicAuthGet(Username, Password, url)
	mustNotErr(err)
	return resp
}

func BasicAuthPatch(url string, contentType string, body io.Reader) *http.Response {
	resp, err := utils.BasicAuthPatch(
		Username,
		Password,
		url,
		contentType,
		body)
	mustNotErr(err)
	return resp
}

func ParseResponseBody(resp *http.Response) []byte {
	b, err := ioutil.ReadAll(resp.Body)
	mustNotErr(err)
	mustNotErr(resp.Body.Close())
	return b
}

func ObserveLogs() *observer.ObservedLogs {
	core, observed := observer.New(zapcore.DebugLevel)
	logger.SetLogger(logger.NewLogger(zap.New(core)))
	return observed
}

func FixtureCreateJobViaWeb(t *testing.T, app *TestApplication, path string) models.JobSpec {
	resp := BasicAuthPost(
		app.Server.URL+"/v2/specs",
		"application/json",
		bytes.NewBuffer(LoadJSON(path)),
	)
	defer resp.Body.Close()
	CheckStatusCode(t, resp, 200)
	j, err := app.Store.FindJob(ParseCommonJSON(resp.Body).ID)
	assert.Nil(t, err)

	return j
}

func CreateJobSpecViaWeb(t *testing.T, app *TestApplication, job models.JobSpec) models.JobSpec {
	marshaled, err := json.Marshal(&job)
	assert.Nil(t, err)
	resp := BasicAuthPost(
		app.Server.URL+"/v2/specs",
		"application/json",
		bytes.NewBuffer(marshaled),
	)
	defer resp.Body.Close()
	CheckStatusCode(t, resp, 200)
	j, err := app.Store.FindJob(ParseCommonJSON(resp.Body).ID)
	assert.Nil(t, err)

	return j
}

func CreateJobRunViaWeb(t *testing.T, app *TestApplication, j models.JobSpec, body ...string) models.JobRun {
	t.Helper()
	url := app.Server.URL + "/v2/specs/" + j.ID + "/runs"
	bodyBuffer := &bytes.Buffer{}
	if len(body) > 0 {
		bodyBuffer = bytes.NewBufferString(body[0])
	}
	resp := BasicAuthPost(url, "application/json", bodyBuffer)
	defer resp.Body.Close()
	CheckStatusCode(t, resp, 200)
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

func UpdateJobRunViaWeb(
	t *testing.T,
	app *TestApplication,
	jr models.JobRun,
	body string,
) models.JobRun {
	t.Helper()
	url := app.Server.URL + "/v2/runs/" + jr.ID
	resp := BasicAuthPatch(url, "application/json", bytes.NewBufferString(body))
	defer resp.Body.Close()

	CheckStatusCode(t, resp, 200)
	jrID := ParseCommonJSON(resp.Body).ID
	assert.Nil(t, app.Store.One("ID", jrID, &jr))
	return jr
}

func CreateBridgeTypeViaWeb(
	t *testing.T,
	app *TestApplication,
	payload string,
) models.BridgeType {
	resp := BasicAuthPost(
		app.Server.URL+"/v2/bridge_types",
		"application/json",
		bytes.NewBufferString(payload),
	)
	defer resp.Body.Close()
	CheckStatusCode(t, resp, 200)
	var bt models.BridgeType
	name := ParseCommonJSON(resp.Body).Name
	assert.Nil(t, app.Store.One("Name", name, &bt))

	return bt
}

func NewClientAndRenderer(config store.Config) (*cmd.Client, *RendererMock) {
	r := &RendererMock{}
	client := &cmd.Client{
		r,
		config,
		EmptyAppFactory{},
		CallbackAuthenticator{func(*store.Store, string) {}},
		EmptyRunner{},
	}
	return client, r
}

func CheckStatusCode(t *testing.T, resp *http.Response, expected int) {
	assert.Equal(t, expected, resp.StatusCode)
	if resp.StatusCode != expected {
		buf, err := ioutil.ReadAll(resp.Body)
		assert.Nil(t, err)
		resp.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
		fmt.Printf("\n\nERROR unexpected HTML Response: %v\n\n", string(buf))
	}
}

func WaitForJobRunToComplete(
	t *testing.T,
	store *store.Store,
	jr models.JobRun,
) models.JobRun {
	return WaitForJobRunStatus(t, store, jr, models.StatusCompleted)
}

func WaitForJobRunToPend(
	t *testing.T,
	store *store.Store,
	jr models.JobRun,
) models.JobRun {
	return WaitForJobRunStatus(t, store, jr, models.StatusPending)
}

func WaitForJobRunStatus(
	t *testing.T,
	store *store.Store,
	jr models.JobRun,
	status string,
) models.JobRun {
	t.Helper()
	gomega.NewGomegaWithT(t).Eventually(func() string {
		assert.Nil(t, store.One("ID", jr.ID, &jr))
		return jr.Status
	}).Should(gomega.Equal(status))
	return jr
}

func StringToRunLogData(str string) hexutil.Bytes {
	length := len([]byte(str))
	lenHex := utils.RemoveHexPrefix(hexutil.EncodeUint64(uint64(length)))
	if len(lenHex) < 64 {
		lenHex = strings.Repeat("0", 64-len(lenHex)) + lenHex
	}

	data := utils.RemoveHexPrefix(utils.StringToHex(str))
	prefix := "0x0000000000000000000000000000000000000000000000000000000000000020"

	var endPad string
	if length%32 != 0 {
		endPad = strings.Repeat("00", (32 - (length % 32)))
	}
	return hexutil.MustDecode(prefix + lenHex + data + endPad)
}

func WaitForRuns(t *testing.T, j models.JobSpec, store *store.Store, want int) []models.JobRun {
	t.Helper()
	g := gomega.NewGomegaWithT(t)

	var jrs []models.JobRun
	var err error
	if want == 0 {
		g.Consistently(func() []models.JobRun {
			jrs, err = store.JobRunsFor(j.ID)
			assert.Nil(t, err)
			return jrs
		}).Should(gomega.HaveLen(want))
	} else {
		g.Eventually(func() []models.JobRun {
			jrs, err = store.JobRunsFor(j.ID)
			assert.Nil(t, err)
			return jrs
		}).Should(gomega.HaveLen(want))
	}
	return jrs
}

func MustParseWebURL(str string) models.WebURL {
	u, err := url.Parse(str)
	mustNotErr(err)
	return models.WebURL{u}
}

func ParseISO8601(s string) time.Time {
	t, err := time.Parse(time.RFC3339Nano, s)
	mustNotErr(err)
	return t
}

func NullableTime(t time.Time) null.Time {
	return null.Time{Time: t, Valid: true}
}

func ParseNullableTime(s string) null.Time {
	return NullableTime(ParseISO8601(s))
}

func IndexableBlockNumber(n int64) *models.IndexableBlockNumber {
	return models.NewIndexableBlockNumber(big.NewInt(n))
}

func mustNotErr(err error) {
	if err != nil {
		logger.Panic(err)
	}
}
