package cltest

import (
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
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/h2non/gock"
	"github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink/logger"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/utils"
	"github.com/smartcontractkit/chainlink/web"
	"github.com/stretchr/testify/assert"
)

const RootDir = "/tmp/chainlink_test"
const Username = "testusername"
const Password = "password"

func init() {
	gomega.SetDefaultEventuallyTimeout(3 * time.Second)
}

type TestConfig struct {
	store.Config
	wsServer *httptest.Server
}

func NewConfig() *TestConfig {
	config := TestConfig{
		Config: store.Config{
			RootDir:             path.Join(RootDir, fmt.Sprintf("%d", time.Now().UnixNano())),
			BasicAuthUsername:   Username,
			BasicAuthPassword:   Password,
			ChainID:             3,
			EthMinConfirmations: 6,
			EthGasBumpWei:       big.NewInt(5000000000),
			EthGasBumpThreshold: 3,
			EthGasPriceDefault:  big.NewInt(20000000000),
			PollingSchedule:     "* * * * * *",
		},
	}
	config.SetEthereumServer(newWSServer())
	return &config
}

func (self *TestConfig) SetEthereumServer(wss *httptest.Server) {
	u, _ := url.Parse(wss.URL)
	u.Scheme = "ws"
	self.EthereumURL = u.String()
	self.wsServer = wss
}

type TestApplication struct {
	*services.Application
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
	return NewApplicationWithConfig(NewConfig())
}

func NewApplicationWithConfig(config *TestConfig) (*TestApplication, func()) {
	app := services.NewApplication(config.Config)
	ta := &TestApplication{
		Application: app,
		Server:      newServer(app),
		wsServer:    config.wsServer,
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

func newServer(app *services.Application) *httptest.Server {
	gin.SetMode(gin.TestMode)
	return httptest.NewServer(web.Router(app))
}

func (self *TestApplication) Stop() {
	self.Application.Stop()
	cleanUpStore(self.Store)
	if self.Server != nil {
		gin.SetMode(gin.DebugMode)
		self.Server.Close()
	}
	if self.wsServer != nil {
		self.wsServer.Close()
	}
}

func NewStoreWithConfig(config *TestConfig) (*store.Store, func()) {
	s := store.NewStore(config.Config)
	return s, func() {
		cleanUpStore(s)
	}
}

func NewStore() (*store.Store, func()) {
	return NewStoreWithConfig(NewConfig())
}

func cleanUpStore(store *store.Store) {
	store.Close()
	if err := os.RemoveAll(store.Config.RootDir); err != nil {
		log.Println(err)
	}
}

func CloseGock(t *testing.T) {
	assert.True(t, gock.IsDone(), "Not all gock requests were fulfilled")
	gock.DisableNetworking()
	gock.Off()
}

type JobJSON struct {
	ID string `json:"id"`
}

func JobJSONFromResponse(body io.Reader) JobJSON {
	b, err := ioutil.ReadAll(body)
	if err != nil {
		log.Fatal(err)
	}
	var respJSON JobJSON
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
	defer to.Close()

	_, err = io.Copy(to, from)
	if err != nil {
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
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return b
}
