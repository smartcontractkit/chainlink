package cltest

import (
	"bytes"
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
	"github.com/smartcontractkit/chainlink/adapters"
	"github.com/smartcontractkit/chainlink/logger"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/utils"
	"github.com/tidwall/gjson"
	null "gopkg.in/guregu/null.v3"
)

// NewJob return new NoOp JobSpec
func NewJob() models.JobSpec {
	j := models.NewJob()
	j.Tasks = []models.TaskSpec{{Type: adapters.TaskTypeNoOp}}
	return j
}

// NewTask given the tasktype and json params return a TaskSpec
func NewTask(taskType string, json ...string) models.TaskSpec {
	if len(json) == 0 {
		json = append(json, ``)
	}
	params := JSONFromString(json[0])
	params, err := params.Add("type", taskType)
	mustNotErr(err)

	return models.TaskSpec{
		Type:   models.MustNewTaskType(taskType),
		Params: params,
	}
}

// NewJobWithSchedule create new job with the given schedule
func NewJobWithSchedule(sched string) (models.JobSpec, models.Initiator) {
	j := NewJob()
	j.Initiators = []models.Initiator{{
		Type: models.InitiatorCron,
		InitiatorParams: models.InitiatorParams{
			Schedule: models.Cron(sched),
		},
	}}
	return j, j.Initiators[0]
}

// NewJobWithWebInitiator create new Job with web inititaor
func NewJobWithWebInitiator() (models.JobSpec, models.Initiator) {
	j := NewJob()
	j.Initiators = []models.Initiator{{Type: models.InitiatorWeb}}
	return j, j.Initiators[0]
}

// NewJobWithLogInitiator create new Job with ethlog inititaor
func NewJobWithLogInitiator() (models.JobSpec, models.Initiator) {
	j := NewJob()
	j.Initiators = []models.Initiator{{
		Type: models.InitiatorEthLog,
		InitiatorParams: models.InitiatorParams{
			Address: NewAddress(),
		},
	}}
	return j, j.Initiators[0]
}

// NewJobWithRunAtInitiator create new Job with RunAt inititaor
func NewJobWithRunAtInitiator(t time.Time) (models.JobSpec, models.Initiator) {
	j := NewJob()
	j.Initiators = []models.Initiator{{
		Type: models.InitiatorRunAt,
		InitiatorParams: models.InitiatorParams{
			Time: models.Time{Time: t},
		},
	}}
	return j, j.Initiators[0]
}

// NewTx create a tx given from address and sentat
func NewTx(from common.Address, sentAt uint64) *models.Tx {
	return &models.Tx{
		From:     from,
		Nonce:    0,
		Data:     []byte{},
		Value:    big.NewInt(0),
		GasLimit: 250000,
	}
}

// CreateTxAndAttempt create tx attempt with given store, from address, and sentat
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

// NewHash return random Keccak256
func NewHash() common.Hash {
	b := make([]byte, 32)
	rand.Read(b)
	return common.BytesToHash(b)
}

// NewAddress return a random new address
func NewAddress() common.Address {
	b := make([]byte, 20)
	rand.Read(b)
	return common.BytesToAddress(b)
}

// NewBridgeType create new bridge type given info slice
func NewBridgeType(info ...string) models.BridgeType {
	bt := models.BridgeType{}

	if len(info) > 0 {
		bt.Name = models.MustNewTaskType(info[0])
	} else {
		bt.Name = models.MustNewTaskType("defaultFixtureBridgeType")
	}

	if len(info) > 1 {
		bt.URL = WebURL(info[1])
	} else {
		bt.URL = WebURL("https://bridge.example.com/api")
	}

	bt.IncomingToken = utils.NewBytes32ID()
	bt.OutgoingToken = utils.NewBytes32ID()

	return bt
}

// NewBridgeTypeWithConfirmations creates a new bridge type with given default confs and info slice
func NewBridgeTypeWithConfirmations(confirmations uint64, info ...string) models.BridgeType {
	bt := NewBridgeType(info...)
	bt.Confirmations = confirmations

	return bt
}

// WebURL parses a url into a models.WebURL
func WebURL(unparsed string) models.WebURL {
	parsed, err := url.Parse(unparsed)
	mustNotErr(err)
	return models.WebURL(*parsed)
}

// NullString creates null.String from given value
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

// NullTime creates a null.Time from given value
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

// LogFromFixture create ethtypes.log from file path
func LogFromFixture(path string) ethtypes.Log {
	value := gjson.Get(string(LoadJSON(path)), "params.result")
	var el ethtypes.Log
	mustNotErr(json.Unmarshal([]byte(value.String()), &el))

	return el
}

// JSONFromFixture create models.JSON from file path
func JSONFromFixture(path string) models.JSON {
	return JSONFromString(string(LoadJSON(path)))
}

// JSONResultFromFixture create model.JSON with params.result found in the given file path
func JSONResultFromFixture(path string) models.JSON {
	res := gjson.Get(string(LoadJSON(path)), "params.result")
	return JSONFromString(res.String())
}

// JSONFromString create JSON from given body and arguments
func JSONFromString(body string, args ...interface{}) models.JSON {
	j, err := models.ParseJSON([]byte(fmt.Sprintf(body, args...)))
	mustNotErr(err)
	return j
}

type EasyJSON struct {
	models.JSON
}

func (ejs EasyJSON) Add(key string, val interface{}) EasyJSON {
	ejs = ejs.Delete(key)

	var err error
	ejs.JSON, err = ejs.JSON.Add(key, val)
	mustNotErr(err)

	return ejs
}

func (ejs EasyJSON) Delete(key string) EasyJSON {
	var err error
	ejs.JSON, err = ejs.JSON.Delete(key)
	mustNotErr(err)
	return ejs
}

func EasyJSONFromFixture(path string) EasyJSON {
	return EasyJSON{JSON: JSONFromFixture(path)}
}

func EasyJSONFromString(body string, args ...interface{}) EasyJSON {
	return EasyJSON{JSON: JSONFromString(body, args...)}
}

// NewRunLog create ethtypes.Log for given jobid, address, block, and json
func NewRunLog(
	jobID string,
	emitter common.Address,
	requester common.Address,
	blk int,
	json string,
) ethtypes.Log {
	return ethtypes.Log{
		Address:     emitter,
		BlockNumber: uint64(blk),
		Data:        StringToVersionedLogData("internalID", json),
		Topics: []common.Hash{
			services.RunLogTopic,
			StringToHash(jobID),
			requester.Hash(),
			minimumContractPayment.ToHash(),
		},
	}
}

// StringToVersionedLogData encodes a string to the log data field.
func StringToVersionedLogData(internalID, str string) []byte {
	buf := bytes.NewBuffer(hexutil.MustDecode(StringToHash(internalID).Hex()))
	buf.Write(hexutil.MustDecode(utils.EVMHexNumber(1)))
	buf.Write(hexutil.MustDecode(utils.EVMHexNumber(common.HashLength * 3)))

	cbor, err := JSONFromString(str).CBOR()
	mustNotErr(err)
	buf.Write(hexutil.MustDecode(utils.EVMHexNumber(len(cbor))))
	paddedLength := common.HashLength * ((len(cbor) / common.HashLength) + 1)
	buf.Write(common.RightPadBytes(cbor, paddedLength))

	return buf.Bytes()
}

// BigHexInt create hexutil.Big value from given value
func BigHexInt(val interface{}) hexutil.Big {
	switch x := val.(type) {
	case int:
		return hexutil.Big(*big.NewInt(int64(x)))
	case uint64:
		return hexutil.Big(*big.NewInt(int64(x)))
	case int64:
		return hexutil.Big(*big.NewInt(x))
	default:
		logger.Panicf("Could not convert %v of type %T to hexutil.Big", val, val)
		return hexutil.Big{}
	}
}

func Int(val interface{}) *models.Int {
	switch x := val.(type) {
	case int:
		return (*models.Int)(big.NewInt(int64(x)))
	case uint64:
		return (*models.Int)(big.NewInt(int64(x)))
	case int64:
		return (*models.Int)(big.NewInt(x))
	default:
		logger.Panicf("Could not convert %v of type %T to models.Int", val, val)
		return &models.Int{}
	}
}

// NewBigHexInt creates new BigHexInt from value
func NewBigHexInt(val interface{}) *hexutil.Big {
	rval := BigHexInt(val)
	return &rval
}

// RunResultWithValue creates a runresult with given value
func RunResultWithValue(val string) models.RunResult {
	data := models.JSON{}
	data, err := data.Add("value", val)
	if err != nil {
		return RunResultWithError(err)
	}

	return models.RunResult{Data: data}
}

// RunResultWithData creates a run result with a given data JSON object
func RunResultWithData(val string) models.RunResult {
	data, err := models.ParseJSON([]byte(val))
	if err != nil {
		return RunResultWithError(err)
	}
	return models.RunResult{Data: data}
}

// RunResultWithError creates a runresult with given error
func RunResultWithError(err error) models.RunResult {
	return models.RunResult{
		Status:       models.RunStatusErrored,
		ErrorMessage: null.StringFrom(err.Error()),
	}
}

// MarkJobRunPendingBridge marks the jobrun as Pending Bridge Status
func MarkJobRunPendingBridge(jr models.JobRun, i int) models.JobRun {
	jr.Status = models.RunStatusPendingBridge
	jr.Result.Status = models.RunStatusPendingBridge
	jr.TaskRuns[i].Status = models.RunStatusPendingBridge
	jr.TaskRuns[i].Result.Status = models.RunStatusPendingBridge
	return jr
}

func NewJobRunner(s *store.Store) (services.JobRunner, func()) {
	rm := services.NewJobRunner(s)
	return rm, func() { rm.Stop() }
}

type MockSigner struct{}

func (s MockSigner) Sign(input []byte) (models.Signature, error) {
	return models.NewSignature("0xb7a987222fc36c4c8ed1b91264867a422769998aadbeeb1c697586a04fa2b616025b5ca936ec5bdb150999e298b6ecf09251d3c4dd1306dedec0692e7037584800")
}

func ServiceAgreementFromString(str string) (models.ServiceAgreement, error) {
	us, err := models.NewUnsignedServiceAgreementFromRequest(strings.NewReader(str))
	if err != nil {
		return models.ServiceAgreement{}, err
	}
	return models.BuildServiceAgreement(us, MockSigner{})
}
