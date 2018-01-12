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
	"github.com/smartcontractkit/chainlink-go/logger"
	"github.com/smartcontractkit/chainlink-go/services"
	"github.com/smartcontractkit/chainlink-go/store"
	"github.com/smartcontractkit/chainlink-go/utils"
	"github.com/smartcontractkit/chainlink-go/web"
	"github.com/stretchr/testify/assert"
)

const RootDir = "/tmp/chainlink_test"
const Username = "testusername"
const Password = "password"

func init() {
	gomega.SetDefaultEventuallyTimeout(3 * time.Second)
}

func NewConfig() store.Config {
	return store.Config{
		RootDir:             path.Join(RootDir, fmt.Sprintf("%d", time.Now().UnixNano())),
		BasicAuthUsername:   Username,
		BasicAuthPassword:   Password,
		EthereumURL:         "http://example.com/api",
		ChainID:             3,
		EthMinConfirmations: 6,
		EthGasBumpWei:       big.NewInt(5000000000),
		EthGasBumpThreshold: 3,
		EthGasPriceDefault:  big.NewInt(20000000000),
		PollingSchedule:     "* * * * * *",
	}
}

type TestApplication struct {
	*services.Application
	Server   *httptest.Server
	WSServer *httptest.Server
}

func newWSServer() *httptest.Server {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		}
		conn, _ := upgrader.Upgrade(w, r, nil)
		msg := ""
		if err := conn.ReadJSON(&msg); err != nil {
			return
		}
	})
	server := httptest.NewServer(handler)
	return server
}

func configureWSServer(c store.Config) (*httptest.Server, store.Config) {
	wsserver := newWSServer()
	u, _ := url.Parse(wsserver.URL)
	u.Scheme = "ws"
	c.EthereumURL = u.String()
	return wsserver, c
}

func NewApplication() *TestApplication {
	ws, c := configureWSServer(NewConfig())
	a := NewApplicationWithConfig(c)
	a.WSServer = ws
	return a
}

func NewApplicationWithConfig(config store.Config) *TestApplication {
	ws, config := configureWSServer(config)
	app := services.NewApplication(config)
	return &TestApplication{
		Application: app,
		WSServer:    ws,
		Server:      newServer(app),
	}
}

func NewApplicationWithKeyStore() *TestApplication {
	app := NewApplication()
	if _, err := app.Store.KeyStore.NewAccount(Password); err != nil {
		logger.Fatal(err)
	}
	if err := app.Store.KeyStore.Unlock(Password); err != nil {
		logger.Fatal(err)
	}
	return app
}

func (self *TestApplication) NewServer() *httptest.Server {
	self.Server = newServer(self.Application)
	return self.Server
}

func newServer(app *services.Application) *httptest.Server {
	gin.SetMode(gin.TestMode)
	return httptest.NewServer(web.Router(app))
}

func (self *TestApplication) Stop() {
	self.Application.Stop()
	CleanUpStore(self.Store)
	if self.Server != nil {
		gin.SetMode(gin.DebugMode)
		self.Server.Close()
	}
	if self.WSServer != nil {
		self.WSServer.Close()
	}
}

func NewStore() *store.Store {
	return store.NewStore(NewConfig())
}

func CleanUpStore(store *store.Store) {
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

func AddPrivateKey(config store.Config, src string) {
	err := os.MkdirAll(config.KeysDir(), os.FileMode(0700))
	if err != nil {
		log.Fatal(err)
	}

	dst := config.KeysDir() + "/testwallet.json"
	copyFile(src, dst)
}

func BasicAuthPost(url string, contentType string, body io.Reader) (*http.Response, error) {
	return utils.BasicAuthPost(
		Username,
		Password,
		url,
		contentType,
		body)
}

func BasicAuthGet(url string) (*http.Response, error) {
	return utils.BasicAuthGet(Username, Password, url)
}
