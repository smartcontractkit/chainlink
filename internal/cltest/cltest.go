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
	"os"
	"path"
	"reflect"
	"testing"
	"time"

	"github.com/araddon/dateparse"
	"github.com/gin-gonic/gin"
	"github.com/h2non/gock"
	"github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink-go/logger"
	"github.com/smartcontractkit/chainlink-go/services"
	"github.com/smartcontractkit/chainlink-go/store"
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
		EthGasBumpWei:      big.NewInt(5000000000),
		EthGasBumpThreshold: 3,
		PollingSchedule:     "* * * * * *",
	}
}

type TestApplication struct {
	*services.Application
	Server *httptest.Server
}

func NewApplication() *TestApplication {
	return NewApplicationWithConfig(NewConfig())
}

func NewApplicationWithConfig(config store.Config) *TestApplication {
	return &TestApplication{Application: services.NewApplication(config)}
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
	gin.SetMode(gin.TestMode)
	server := httptest.NewServer(web.Router(self.Application))
	self.Server = server
	return server
}

func (self *TestApplication) Stop() {
	self.Application.Stop()
	CleanUpStore(self.Store)
	if self.Server != nil {
		gin.SetMode(gin.DebugMode)
		self.Server.Close()
	}
}

func (self *TestApplication) MockEthClient() *EthMock {
	mock := NewMockGethRpc()
	eth := &store.EthClient{mock}
	self.Store.Eth.EthClient = eth
	return mock
}

func NewMockGethRpc() *EthMock {
	return &EthMock{}
}

type EthMock struct {
	Responses []MockResponse
}

type MockResponse struct {
	methodName string
	response   interface{}
	errMsg     string
	hasError   bool
}

func (self *EthMock) Register(method string, response interface{}) {
	res := MockResponse{
		methodName: method,
		response:   response,
	}
	self.Responses = append(self.Responses, res)
}

func (self *EthMock) RegisterError(method, errMsg string) {
	res := MockResponse{
		methodName: method,
		errMsg:     errMsg,
		hasError:   true,
	}
	self.Responses = append(self.Responses, res)
}

func (self *EthMock) AllCalled() bool {
	return len(self.Responses) == 0
}

func copyWithoutIndex(s []MockResponse, index int) []MockResponse {
	return append(s[:index], s[index+1:]...)
}

func (self *EthMock) Call(result interface{}, method string, args ...interface{}) error {
	for i, resp := range self.Responses {
		if resp.methodName == method {
			self.Responses = copyWithoutIndex(self.Responses, i)
			if resp.hasError {
				return fmt.Errorf(resp.errMsg)
			} else {
				ref := reflect.ValueOf(result)
				reflect.Indirect(ref).Set(reflect.ValueOf(resp.response))
				return nil
			}
		}
	}
	return fmt.Errorf("EthMock: Method %v not registered", method)
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

func TimeParse(s string) time.Time {
	t, err := dateparse.ParseAny(s)
	if err != nil {
		log.Fatal(err)
	}
	return t
}

func BasicAuthPost(url string, contentType string, body io.Reader) (*http.Response, error) {
	client := &http.Client{}
	request, _ := http.NewRequest("POST", url, body)
	request.Header.Set("Content-Type", contentType)
	request.SetBasicAuth(Username, Password)
	resp, err := client.Do(request)
	return resp, err
}

func BasicAuthGet(url string) (*http.Response, error) {
	client := &http.Client{}
	request, _ := http.NewRequest("GET", url, nil)
	request.SetBasicAuth(Username, Password)
	resp, err := client.Do(request)
	return resp, err
}
