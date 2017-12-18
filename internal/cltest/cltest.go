package cltest

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"
	"time"

	"github.com/araddon/dateparse"
	"github.com/gin-gonic/gin"
	"github.com/h2non/gock"
	"github.com/onsi/gomega"
	configlib "github.com/smartcontractkit/chainlink-go/config"
	"github.com/smartcontractkit/chainlink-go/logger"
	"github.com/smartcontractkit/chainlink-go/services"
	"github.com/smartcontractkit/chainlink-go/web"
	"github.com/stretchr/testify/assert"
)

const testRootDir = "/tmp/chainlink_test"
const testUsername = "testusername"
const testPassword = "testpassword"

func init() {
	gomega.SetDefaultEventuallyTimeout(2 * time.Second)
}

type TestStore struct {
	*services.Store
	Server *httptest.Server
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

func Store() *TestStore {
	config := NewConfig()
	if err := os.MkdirAll(config.RootDir, os.FileMode(0700)); err != nil {
		log.Fatal(err)
	}
	logger.SetLoggerDir(config.RootDir)
	store := services.NewStore(config)
	return &TestStore{
		Store: store,
	}
}

func NewConfig() configlib.Config {
	return configlib.Config{
		RootDir:           path.Join(testRootDir, fmt.Sprintf("%d", time.Now().UnixNano())),
		BasicAuthUsername: testUsername,
		BasicAuthPassword: testPassword,
		EthereumURL:       "http://example.com/api",
	}
}

func (self *TestStore) SetUpWeb() *httptest.Server {
	gin.SetMode(gin.TestMode)
	server := httptest.NewServer(web.Router(self.Store))
	self.Server = server
	return server
}

func (self *TestStore) Close() {
	self.Store.Close()
	if err := os.RemoveAll(self.Config.RootDir); err != nil {
		log.Println(err)
	}

	if self.Server != nil {
		gin.SetMode(gin.DebugMode)
		self.Server.Close()
	}
}

func CloseGock(t *testing.T) {
	assert.True(t, gock.IsDone(), "Not all gock requests were fulfilled")
	gock.DisableNetworking()
	gock.Off()
}

func LoadJSON(file string) []byte {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	return content
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
	request.SetBasicAuth(testUsername, testPassword)
	resp, err := client.Do(request)
	return resp, err
}

func BasicAuthGet(url string) (*http.Response, error) {
	client := &http.Client{}
	request, _ := http.NewRequest("GET", url, nil)
	request.SetBasicAuth(testUsername, testPassword)
	resp, err := client.Do(request)
	return resp, err
}
