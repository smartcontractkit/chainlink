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
	"github.com/smartcontractkit/chainlink-go/logger"
	"github.com/smartcontractkit/chainlink-go/services"
	"github.com/smartcontractkit/chainlink-go/store"
	"github.com/smartcontractkit/chainlink-go/web"
	"github.com/stretchr/testify/assert"
)

const testRootDir = "/tmp/chainlink_test"
const testUsername = "testusername"
const testPassword = "testpassword"

func init() {
	gomega.SetDefaultEventuallyTimeout(2 * time.Second)
}

func NewConfig() store.Config {
	return store.Config{
		RootDir:           path.Join(testRootDir, fmt.Sprintf("%d", time.Now().UnixNano())),
		BasicAuthUsername: testUsername,
		BasicAuthPassword: testPassword,
		EthereumURL:       "http://example.com/api",
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
	app := &TestApplication{Application: services.NewApplication(config)}
	logger.SetLoggerDir(config.RootDir)
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

func CopyFile(src, dst string) {
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
	CopyFile(src, dst)
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
