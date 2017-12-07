package cltest

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http/httptest"
	"os"
	"path/filepath"
	"time"

	"github.com/araddon/dateparse"
	"github.com/gin-gonic/gin"
	"github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink-go/logger"
	"github.com/smartcontractkit/chainlink-go/models"
	"github.com/smartcontractkit/chainlink-go/services"
	"github.com/smartcontractkit/chainlink-go/web"
)

func init() {
	if err := os.RemoveAll(filepath.Dir(models.DBPath("test"))); err != nil {
		log.Println(err)
	}

	gomega.SetDefaultEventuallyTimeout(2 * time.Second)
	logger.SetLogger(logger.NewLogger("test"))
}

type TestStore struct {
	Server *httptest.Server
	services.Store
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

func Store() TestStore {
	orm := models.InitORM("test")
	store := services.Store{
		ORM:       orm,
		Scheduler: services.NewScheduler(orm),
	}
	return TestStore{
		Store: store,
	}
}

func (self TestStore) SetUpWeb() *httptest.Server {
	gin.SetMode(gin.TestMode)
	server := httptest.NewServer(web.Router(self.Store))
	self.Server = server
	return server
}

func (self TestStore) Close() {
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
