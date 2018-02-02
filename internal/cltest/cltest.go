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

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/h2non/gock"
	"github.com/onsi/gomega"
	. "github.com/onsi/gomega"
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
			RootDir:             rootdir,
			BasicAuthUsername:   Username,
			BasicAuthPassword:   Password,
			ChainID:             3,
			EthMinConfirmations: 6,
			EthGasBumpWei:       *big.NewInt(5000000000),
			EthGasBumpThreshold: 3,
			EthGasPriceDefault:  *big.NewInt(20000000000),
			PollingSchedule:     "* * * * * *",
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
	ta := &TestApplication{
		ChainlinkApplication: app,
		Server:               server,
		wsServer:             tc.wsServer,
	}
	return ta, func() {
		ta.Stop()
	}
}

func NewApplicationWithKeyStore() (*TestApplication, func()) {
	app, cleanup := NewApplication()
	if _, err := app.Store.KeyStore.NewAccount(Password); err != nil {
		logger.Fatal(err)
	}
	if err := app.Store.KeyStore.Unlock(Password); err != nil {
		logger.Fatal(err)
	}
	return app, cleanup
}

func newServer(app *services.ChainlinkApplication) *httptest.Server {
	return httptest.NewServer(web.Router(app))
}

func (ta *TestApplication) Stop() {
	ta.ChainlinkApplication.Stop()
	cleanUpStore(ta.Store)
	if ta.Server != nil {
		ta.Server.Close()
	}
	if ta.wsServer != nil {
		ta.wsServer.Close()
	}
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
	if err != nil {
		log.Fatal(err)
	}
	var respJSON CommonJSON
	json.Unmarshal(b, &respJSON)
	return respJSON
}

func LoadJSON(file string) []byte {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	return content
}

func copyFile(src, dst string) {
	from, err := os.Open(src)
	if err != nil {
		log.Fatal(err)
	}
	defer from.Close()

	to, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}

	_, err = io.Copy(to, from)
	if err != nil {
		log.Fatal(err)
	}

	if err = to.Close(); err != nil {
		log.Fatal(err)
	}
}

func AddPrivateKey(config *TestConfig, src string) {
	err := os.MkdirAll(config.KeysDir(), os.FileMode(0700))
	if err != nil {
		log.Fatal(err)
	}

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
	if err != nil {
		log.Fatal(err)
	}
	return resp
}

func BasicAuthGet(url string) *http.Response {
	resp, err := utils.BasicAuthGet(Username, Password, url)
	if err != nil {
		log.Fatal(err)
	}
	return resp
}

func ParseResponseBody(resp *http.Response) []byte {
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	if err = resp.Body.Close(); err != nil {
		log.Fatal(err)
	}
	return b
}

func ObserveLogs() *observer.ObservedLogs {
	core, observed := observer.New(zapcore.DebugLevel)
	logger.SetLogger(logger.NewLogger(zap.New(core)))
	return observed
}

func FixtureCreateJobViaWeb(t *testing.T, app *TestApplication, path string) *models.Job {
	resp := BasicAuthPost(
		app.Server.URL+"/v2/jobs",
		"application/json",
		bytes.NewBuffer(LoadJSON(path)),
	)
	defer resp.Body.Close()
	CheckStatusCode(t, resp, 200)
	j, err := app.Store.FindJob(ParseCommonJSON(resp.Body).ID)
	assert.Nil(t, err)

	return j
}

func CreateJobRunViaWeb(t *testing.T, app *TestApplication, j *models.Job) *models.JobRun {
	url := app.Server.URL + "/v2/jobs/" + j.ID + "/runs"
	resp := BasicAuthPost(url, "application/json", &bytes.Buffer{})
	defer resp.Body.Close()
	CheckStatusCode(t, resp, 200)
	jrID := ParseCommonJSON(resp.Body).ID

	jrs := []*models.JobRun{}
	Eventually(func() []*models.JobRun {
		assert.Nil(t, app.Store.Where("ID", jrID, &jrs))
		return jrs
	}).Should(HaveLen(1))
	jr := jrs[0]
	assert.Equal(t, j.ID, jr.JobID)

	return jr
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
	app *TestApplication,
	jr *models.JobRun,
) *models.JobRun {
	Eventually(func() string {
		assert.Nil(t, app.Store.One("ID", jr.ID, jr))
		return jr.Status
	}).Should(Equal(models.StatusCompleted))
	return jr
}
