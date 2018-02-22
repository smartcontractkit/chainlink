package cltest

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"net/url"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/tidwall/gjson"
	null "gopkg.in/guregu/null.v3"
)

func NewJob() models.Job {
	j := models.NewJob()
	j.Tasks = []models.Task{{Type: "NoOp"}}
	return j
}

func NewTask(taskType, json string) models.Task {
	params := JSONFromString(json)
	params, err := params.Add("type", taskType)
	mustNotErr(err)

	return models.Task{
		Type:   taskType,
		Params: params,
	}
}

func NewJobWithSchedule(sched string) models.Job {
	j := NewJob()
	j.Initiators = []models.Initiator{{Type: models.InitiatorCron, Schedule: models.Cron(sched)}}
	return j
}

func NewJobWithWebInitiator() models.Job {
	j := NewJob()
	j.Initiators = []models.Initiator{{Type: models.InitiatorWeb}}
	return j
}

func NewJobWithLogInitiator() models.Job {
	j := NewJob()
	j.Initiators = []models.Initiator{{
		Type:    models.InitiatorEthLog,
		Address: NewAddress(),
	}}
	return j
}

func NewTx(from common.Address, sentAt uint64) *models.Tx {
	return &models.Tx{
		From:     from,
		Nonce:    0,
		Data:     []byte{},
		Value:    big.NewInt(0),
		GasLimit: 250000,
	}
}

func CreateTxAndAttempt(
	store *store.Store,
	from common.Address,
	sentAt uint64,
) *models.Tx {
	tx := NewTx(from, sentAt)
	mustNotErr(store.Save(tx))
	_, err := store.AddAttempt(tx, tx.EthTx(big.NewInt(1)), sentAt)
	mustNotErr(err)
	return tx
}

func NewHash() common.Hash {
	b := make([]byte, 32)
	rand.Read(b)
	return common.BytesToHash(b)
}

func NewAddress() common.Address {
	b := make([]byte, 20)
	rand.Read(b)
	return common.BytesToAddress(b)
}

func NewBridgeType(info ...string) models.BridgeType {
	bt := models.BridgeType{}

	if len(info) > 0 {
		bt.Name = strings.ToLower(info[0])
	} else {
		bt.Name = strings.ToLower("defaultFixtureBridgeType")
	}

	if len(info) > 1 {
		bt.URL = WebURL(info[1])
	} else {
		bt.URL = WebURL("https://bridge.example.com/api")
	}

	return bt
}

func WebURL(unparsed string) models.WebURL {
	parsed, err := url.Parse(unparsed)
	mustNotErr(err)
	return models.WebURL{parsed}
}

func NullString(val interface{}) null.String {
	switch val.(type) {
	case string:
		return null.StringFrom(val.(string))
	case nil:
		return null.NewString("", false)
	default:
		panic("cannot create a null string of any type other than string or nil")
	}
}

func NullTime(val interface{}) null.Time {
	switch val.(type) {
	case string:
		return ParseNullableTime(val.(string))
	case nil:
		return null.NewTime(time.Unix(0, 0), false)
	default:
		panic("cannot create a null time of any type other than string or nil")
	}
}

func LogFromFixture(path string) ethtypes.Log {
	value := gjson.Get(string(LoadJSON(path)), "params.result")
	var el ethtypes.Log
	mustNotErr(json.Unmarshal([]byte(value.String()), &el))

	return el
}

func JSONFromFixture(path string) models.JSON {
	return JSONFromString(string(LoadJSON(path)))
}

func JSONFromString(body string, args ...interface{}) models.JSON {
	var j models.JSON
	str := fmt.Sprintf(body, args...)
	mustNotErr(json.Unmarshal([]byte(str), &j))
	return j
}

func NewRunLog(jobID string, addr common.Address, json string) ethtypes.Log {
	return ethtypes.Log{
		Address: addr,
		Data:    StringToRunLogPayload(json),
		Topics: []common.Hash{
			services.RunLogTopic,
			common.StringToHash("requestID"),
			common.StringToHash(jobID),
		},
	}
}

func BigHexInt(val uint64) hexutil.Big {
	return hexutil.Big(*big.NewInt(int64(val)))
}

func RunResultWithValue(val string) models.RunResult {
	data := models.JSON{}
	data, err := data.Add("value", val)
	if err != nil {
		return RunResultWithError(err)
	}

	return models.RunResult{Data: data}
}

func RunResultWithError(err error) models.RunResult {
	return models.RunResult{
		ErrorMessage: null.StringFrom(err.Error()),
	}
}
