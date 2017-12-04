package cltest

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"time"

	"github.com/araddon/dateparse"
	"github.com/gin-gonic/gin"
	"github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink-go/models"
	"github.com/smartcontractkit/chainlink-go/services"
	"github.com/smartcontractkit/chainlink-go/store"
	"github.com/smartcontractkit/chainlink-go/web"
)

func init() {
	gomega.SetDefaultEventuallyTimeout(2 * time.Second)
	services.SetLogger(services.NewLogger("test"))
}

var server *httptest.Server

type JobJSON struct {
	ID string `json:"id"`
}

func JobJSONFromResponse(resp *http.Response) JobJSON {
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var respJSON JobJSON
	json.Unmarshal(b, &respJSON)
	return respJSON
}

func Store() store.Store {
	os.Remove(models.DBPath("test"))
	orm := models.InitORM("test")
	return store.Store{
		ORM:       orm,
		Scheduler: services.NewScheduler(orm),
	}
}

func SetUpWeb(s store.Store) *httptest.Server {
	gin.SetMode(gin.TestMode)
	server = httptest.NewServer(web.Router(s))
	return server
}

func TearDownWeb() {
	gin.SetMode(gin.DebugMode)
	server.Close()
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
