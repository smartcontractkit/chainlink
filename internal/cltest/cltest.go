package cltest

import (
	"encoding/json"
	"github.com/araddon/dateparse"
	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink-go/models"
	"github.com/smartcontractkit/chainlink-go/scheduler"
	"github.com/smartcontractkit/chainlink-go/web"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"time"
)

var server *httptest.Server
var sched *scheduler.Scheduler

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

func SetUpDB() {
	models.InitDBTest()
}

func TearDownDB() {
	models.CloseDB()
}

func SetUpWeb() *httptest.Server {
	gin.SetMode(gin.TestMode)
	server = httptest.NewServer(web.Router())
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
