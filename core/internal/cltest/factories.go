package cltest

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"flag"
	"fmt"
	"math/big"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/smartcontractkit/chainlink/core/adapters"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/store"
	strpkg "github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/assets"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/sjson"
	"github.com/urfave/cli"
	null "gopkg.in/guregu/null.v3"
)

// NewJob return new NoOp JobSpec
func NewJob() models.JobSpec {
	j := models.NewJob()
	j.Tasks = []models.TaskSpec{{Type: adapters.TaskTypeNoOp}}
	return j
}

// NewTask given the tasktype and json params return a TaskSpec
func NewTask(t *testing.T, taskType string, json ...string) models.TaskSpec {
	if len(json) == 0 {
		json = append(json, ``)
	}
	params := JSONFromString(t, json[0])
	params, err := params.Add("type", taskType)
	require.NoError(t, err)

	return models.TaskSpec{
		Type:   models.MustNewTaskType(taskType),
		Params: params,
	}
}

// NewJobWithSchedule create new job with the given schedule
func NewJobWithSchedule(sched string) models.JobSpec {
	j := NewJob()
	j.Initiators = []models.Initiator{{
		JobSpecID: j.ID,
		Type:      models.InitiatorCron,
		InitiatorParams: models.InitiatorParams{
			Schedule: models.Cron(sched),
		},
	}}
	return j
}

// NewJobWithWebInitiator create new Job with web inititaor
func NewJobWithWebInitiator() models.JobSpec {
	j := NewJob()
	j.Initiators = []models.Initiator{{
		JobSpecID: j.ID,
		Type:      models.InitiatorWeb,
	}}
	return j
}

// NewJobWithLogInitiator create new Job with ethlog inititaor
func NewJobWithLogInitiator() models.JobSpec {
	j := NewJob()
	j.Initiators = []models.Initiator{{
		JobSpecID: j.ID,
		Type:      models.InitiatorEthLog,
		InitiatorParams: models.InitiatorParams{
			Address: NewAddress(),
		},
	}}
	return j
}

// NewJobWithRunLogInitiator creates a new JobSpec with the RunLog initiator
func NewJobWithRunLogInitiator() models.JobSpec {
	j := NewJob()
	j.Initiators = []models.Initiator{{
		JobSpecID: j.ID,
		Type:      models.InitiatorRunLog,
		InitiatorParams: models.InitiatorParams{
			Address: NewAddress(),
		},
	}}
	return j
}

// NewJobWithSALogInitiator creates new JobSpec with the ServiceAgreement
// initiator
func NewJobWithSALogInitiator() models.JobSpec {
	j := NewJobWithRunLogInitiator()
	j.Initiators[0].Type = models.InitiatorServiceAgreementExecutionLog
	return j
}

// NewJobWithRunAtInitiator create new Job with RunAt inititaor
func NewJobWithRunAtInitiator(t time.Time) models.JobSpec {
	j := NewJob()
	j.Initiators = []models.Initiator{{
		JobSpecID: j.ID,
		Type:      models.InitiatorRunAt,
		InitiatorParams: models.InitiatorParams{
			Time: models.NewAnyTime(t),
		},
	}}
	return j
}

// NewTx create a tx given from address and sentat
func NewTx(from common.Address, sentAt uint64) *models.Tx {
	return &models.Tx{
		From:     from,
		Nonce:    0,
		Data:     []byte{},
		Value:    models.NewBig(big.NewInt(0)),
		GasLimit: 250000,
	}
}

// CreateTxAndAttempt create tx attempt with given store, from address, and sentat
func CreateTxAndAttempt(
	store *strpkg.Store,
	from common.Address,
	sentAt uint64,
) *models.Tx {
	tx := NewTx(from, sentAt)
	b := make([]byte, 36)
	binary.LittleEndian.PutUint64(b, uint64(sentAt))
	tx.Data = b
	mustNotErr(store.SaveTx(tx))
	_, err := store.AddTxAttempt(tx, tx.EthTx(big.NewInt(1)), sentAt)
	mustNotErr(err)
	return tx
}

// NewHash return random Keccak256
func NewHash() common.Hash {
	return common.BytesToHash(randomBytes(32))
}

// NewAddress return a random new address
func NewAddress() common.Address {
	return common.BytesToAddress(randomBytes(20))
}

func randomBytes(n int) []byte {
	b := make([]byte, n)
	rand.Read(b)
	return b
}

// NewBridgeType create new bridge type given info slice
func NewBridgeType(info ...string) (*models.BridgeTypeAuthentication, *models.BridgeType) {
	btr := &models.BridgeTypeRequest{}

	if len(info) > 0 {
		btr.Name = models.MustNewTaskType(info[0])
	} else {
		btr.Name = models.MustNewTaskType("defaultFixtureBridgeType")
	}

	if len(info) > 1 {
		btr.URL = WebURL(info[1])
	} else {
		btr.URL = WebURL("https://bridge.example.com/api")
	}

	bta, bt, err := models.NewBridgeType(btr)
	mustNotErr(err)
	return bta, bt
}

// NewBridgeTypeWithConfirmations creates a new bridge type with given default confs and info slice
func NewBridgeTypeWithConfirmations(confirmations uint64, info ...string) *models.BridgeType {
	_, bt := NewBridgeType(info...)
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

// JSONFromString create JSON from given body and arguments
func JSONFromString(t *testing.T, body string, args ...interface{}) models.JSON {
	return JSONFromBytes(t, []byte(fmt.Sprintf(body, args...)))
}

// JSONFromBytes creates JSON from a given byte array
func JSONFromBytes(t *testing.T, body []byte) models.JSON {
	j, err := models.ParseJSON(body)
	require.NoError(t, err)
	return j
}

// MustJSONSet uses sjson.Set to set a path in a JSON string and returns the string
// See https://github.com/tidwall/sjson
func MustJSONSet(t *testing.T, json, path string, value interface{}) string {
	json, err := sjson.Set(json, path, value)
	require.NoError(t, err)
	return json
}

// MustJSONDel uses sjson.Delete to remove a path from a JSON string and returns the string
func MustJSONDel(t *testing.T, json, path string) string {
	json, err := sjson.Delete(json, path)
	require.NoError(t, err)
	return json
}

// NewRunLog create models.Log for given jobid, address, block, and json
func NewRunLog(
	t *testing.T,
	jobID string,
	emitter common.Address,
	requester common.Address,
	blk int,
	json string,
) models.Log {
	return models.Log{
		Address:     emitter,
		BlockNumber: uint64(blk),
		Data:        StringToVersionedLogData20190207withoutIndexes(t, "internalID", requester, json),
		Topics: []common.Hash{
			models.RunLogTopic20190207withoutIndexes,
			StringToHash(jobID),
		},
	}
}

// NewServiceAgreementExecutionLog creates a log event for the given jobid,
// address, block, and json, to simulate a request for execution on a service
// agreement.
func NewServiceAgreementExecutionLog(
	t *testing.T,
	jobID string,
	logEmitter common.Address,
	executionRequester common.Address,
	blockHeight int,
	serviceAgreementJSON string,
) models.Log {
	return models.Log{
		Address:     logEmitter,
		BlockNumber: uint64(blockHeight),
		Data:        StringToVersionedLogData0(t, "internalID", serviceAgreementJSON),
		Topics: []common.Hash{
			models.ServiceAgreementExecutionLogTopic,
			StringToHash(jobID),
			executionRequester.Hash(),
			assets.NewLink(1000000000).ToHash(),
		},
	}
}

func StringToVersionedLogData0(t *testing.T, internalID, str string) []byte {
	buf := bytes.NewBuffer(hexutil.MustDecode(StringToHash(internalID).Hex()))
	buf.Write(utils.EVMWordUint64(1))
	buf.Write(utils.EVMWordUint64(common.HashLength * 3))

	cbor, err := JSONFromString(t, str).CBOR()
	require.NoError(t, err)
	buf.Write(utils.EVMWordUint64(uint64(len(cbor))))
	paddedLength := common.HashLength * ((len(cbor) / common.HashLength) + 1)
	buf.Write(common.RightPadBytes(cbor, paddedLength))

	return buf.Bytes()
}

func StringToVersionedLogData20190123withFulfillmentParams(t *testing.T, internalID, str string) []byte {
	requestID := hexutil.MustDecode(StringToHash(internalID).Hex())
	buf := bytes.NewBuffer(requestID)

	version := utils.EVMWordUint64(1)
	buf.Write(version)

	dataLocation := utils.EVMWordUint64(common.HashLength * 6)
	buf.Write(dataLocation)

	callbackAddr := utils.EVMWordUint64(0)
	buf.Write(callbackAddr)

	callbackFunc := utils.EVMWordUint64(0)
	buf.Write(callbackFunc)

	expiration := utils.EVMWordUint64(4000000000)
	buf.Write(expiration)

	cbor, err := JSONFromString(t, str).CBOR()
	require.NoError(t, err)
	buf.Write(utils.EVMWordUint64(uint64(len(cbor))))
	paddedLength := common.HashLength * ((len(cbor) / common.HashLength) + 1)
	buf.Write(common.RightPadBytes(cbor, paddedLength))

	return buf.Bytes()
}

func StringToVersionedLogData20190207withoutIndexes(
	t *testing.T,
	internalID string,
	requester common.Address,
	str string,
) []byte {
	requesterBytes := requester.Hash().Bytes()
	buf := bytes.NewBuffer(requesterBytes)

	requestID := hexutil.MustDecode(StringToHash(internalID).Hex())
	buf.Write(requestID)

	payment := hexutil.MustDecode(minimumContractPayment.ToHash().Hex())
	buf.Write(payment)

	callbackAddr := utils.EVMWordUint64(0)
	buf.Write(callbackAddr)

	callbackFunc := utils.EVMWordUint64(0)
	buf.Write(callbackFunc)

	expiration := utils.EVMWordUint64(4000000000)
	buf.Write(expiration)

	version := utils.EVMWordUint64(1)
	buf.Write(version)

	dataLocation := utils.EVMWordUint64(common.HashLength * 8)
	buf.Write(dataLocation)

	cbor, err := JSONFromString(t, str).CBOR()
	require.NoError(t, err)
	buf.Write(utils.EVMWordUint64(uint64(len(cbor))))
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

func Int(val interface{}) *models.Big {
	switch x := val.(type) {
	case int:
		return (*models.Big)(big.NewInt(int64(x)))
	case uint64:
		return (*models.Big)(big.NewInt(int64(x)))
	case int64:
		return (*models.Big)(big.NewInt(x))
	default:
		logger.Panicf("Could not convert %v of type %T to models.Big", val, val)
		return &models.Big{}
	}
}

// NewBigHexInt creates new BigHexInt from value
func NewBigHexInt(val interface{}) *hexutil.Big {
	rval := BigHexInt(val)
	return &rval
}

// RunResultWithResult creates a RunResult with given result
func RunResultWithResult(val string) models.RunResult {
	data := models.JSON{}
	data, err := data.Add("result", val)
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

func NewJobRunner(s *strpkg.Store) (services.JobRunner, func()) {
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

func EmptyCLIContext() *cli.Context {
	set := flag.NewFlagSet("test", 0)
	return cli.NewContext(nil, set, nil)
}

func CreateJobRunWithStatus(store *store.Store, j models.JobSpec, status models.RunStatus) models.JobRun {
	initr := j.Initiators[0]
	run := j.NewRun(initr)
	run.Status = status
	mustNotErr(store.CreateJobRun(&run))
	return run
}
