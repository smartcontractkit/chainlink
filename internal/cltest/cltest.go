package cltest

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http/httptest"
	"os"
	"path"
	"time"

	"github.com/araddon/dateparse"
	"github.com/gin-gonic/gin"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink-go/logger"
	"github.com/smartcontractkit/chainlink-go/services"
	"github.com/smartcontractkit/chainlink-go/web"
)

const testRootDir = "../../tmp/test"

func init() {
	dir, err := homedir.Expand(testRootDir)
	if err != nil {
		logger.Fatal(err)
	}

	if err = os.RemoveAll(dir); err != nil {
		log.Println(err)
	}

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
	config := services.NewConfig(path.Join(testRootDir, fmt.Sprintf("%d", time.Now().UnixNano())))
	logger.SetLoggerDir(config.RootDir)
	store := services.NewStore(config)
	return &TestStore{
		Store: store,
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
	if self.Server != nil {
		gin.SetMode(gin.DebugMode)
		self.Server.Close()
	}
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
