package cltest

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"flag"
	"fmt"
	"math/big"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"chainlink/core/adapters"
	"chainlink/core/eth"
	"chainlink/core/logger"
	"chainlink/core/store"
	strpkg "chainlink/core/store"
	"chainlink/core/store/assets"
	"chainlink/core/store/models"
	"chainlink/core/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
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

// NewJobWithExternalInitiator creates new Job with external inititaor
func NewJobWithExternalInitiator(ei *models.ExternalInitiator) models.JobSpec {
	j := NewJob()
	j.Initiators = []models.Initiator{{
		JobSpecID: j.ID,
		Type:      models.InitiatorExternal,
		InitiatorParams: models.InitiatorParams{
			Name: ei.Name,
		},
	}}
	return j
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

// NewTx returns a Tx using a specified from address and sentAt
func NewTx(from common.Address, sentAt uint64) *models.Tx {
	tx := &models.Tx{
		From:     from,
		Nonce:    0,
		Data:     []byte{},
		Value:    utils.NewBig(big.NewInt(0)),
		GasLimit: 250000,
		SentAt:   sentAt,
	}
	copy(tx.Hash[:], randomBytes(common.HashLength))
	return tx
}

// CreateTx creates a Tx from a specified address, and sentAt
func CreateTx(
	t testing.TB,
	store *strpkg.Store,
	from common.Address,
	sentAt uint64,
) *models.Tx {
	return CreateTxWithNonce(t, store, from, sentAt, 0)
}

// CreateTxWithNonce creates a Tx from a specified address, sentAt, and nonce
func CreateTxWithNonce(
	t testing.TB,
	store *strpkg.Store,
	from common.Address,
	sentAt uint64,
	nonce uint64,
) *models.Tx {
	data := make([]byte, 36)
	binary.LittleEndian.PutUint64(data, sentAt)
	ethTx := types.NewTransaction(nonce, common.Address{}, big.NewInt(0), 250000, big.NewInt(1), data)
	tx, err := store.CreateTx(null.String{}, ethTx, &from, sentAt)
	require.NoError(t, err)
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
func NewBridgeType(t testing.TB, info ...string) (*models.BridgeTypeAuthentication, *models.BridgeType) {
	btr := &models.BridgeTypeRequest{}

	if len(info) > 0 {
		btr.Name = models.MustNewTaskType(info[0])
	} else {
		btr.Name = models.MustNewTaskType("defaultFixtureBridgeType")
	}

	if len(info) > 1 {
		btr.URL = WebURL(t, info[1])
	} else {
		btr.URL = WebURL(t, "https://bridge.example.com/api")
	}

	bta, bt, err := models.NewBridgeType(btr)
	require.NoError(t, err)
	return bta, bt
}

// WebURL parses a url into a models.WebURL
func WebURL(t testing.TB, unparsed string) models.WebURL {
	parsed, err := url.Parse(unparsed)
	require.NoError(t, err)
	return models.WebURL(*parsed)
}

// JSONFromString create JSON from given body and arguments
func JSONFromString(t testing.TB, body string, args ...interface{}) models.JSON {
	return JSONFromBytes(t, []byte(fmt.Sprintf(body, args...)))
}

// JSONFromBytes creates JSON from a given byte array
func JSONFromBytes(t testing.TB, body []byte) models.JSON {
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

// NewRunLog create eth.Log for given jobid, address, block, and json
func NewRunLog(
	t *testing.T,
	jobID *models.ID,
	emitter common.Address,
	requester common.Address,
	blk int,
	json string,
) eth.Log {
	return eth.Log{
		Address:     emitter,
		BlockNumber: uint64(blk),
		Data:        StringToVersionedLogData20190207withoutIndexes(t, "internalID", requester, json),
		TxHash:      NewHash(),
		BlockHash:   NewHash(),
		Topics: []common.Hash{
			models.RunLogTopic20190207withoutIndexes,
			models.IDToTopic(jobID),
		},
	}
}

// NewServiceAgreementExecutionLog creates a log event for the given jobid,
// address, block, and json, to simulate a request for execution on a service
// agreement.
func NewServiceAgreementExecutionLog(
	t *testing.T,
	jobID *models.ID,
	logEmitter common.Address,
	executionRequester common.Address,
	blockHeight int,
	serviceAgreementJSON string,
) eth.Log {
	return eth.Log{
		Address:     logEmitter,
		BlockNumber: uint64(blockHeight),
		Data:        StringToVersionedLogData0(t, "internalID", serviceAgreementJSON),
		Topics: []common.Hash{
			models.ServiceAgreementExecutionLogTopic,
			models.IDToTopic(jobID),
			executionRequester.Hash(),
			NewLink(t, "1000000000000000000").ToHash(),
		},
	}
}

func NewLink(t *testing.T, amount string) *assets.Link {
	link := assets.NewLink(0)
	link, ok := link.SetString(amount, 10)
	assert.True(t, ok)
	return link
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
	case uint32:
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

func Int(val interface{}) *utils.Big {
	switch x := val.(type) {
	case int:
		return (*utils.Big)(big.NewInt(int64(x)))
	case uint32:
		return (*utils.Big)(big.NewInt(int64(x)))
	case uint64:
		return (*utils.Big)(big.NewInt(int64(x)))
	case int64:
		return (*utils.Big)(big.NewInt(x))
	default:
		logger.Panicf("Could not convert %v of type %T to utils.Big", val, val)
		return &utils.Big{}
	}
}

// MarkJobRunPendingBridge marks the jobrun as Pending Bridge Status
func MarkJobRunPendingBridge(jr models.JobRun, i int) models.JobRun {
	jr.Status = models.RunStatusPendingBridge
	jr.TaskRuns[i].Status = models.RunStatusPendingBridge
	return jr
}

type MockSigner struct{}

func (s MockSigner) SignHash(common.Hash) (models.Signature, error) {
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

func CreateJobRunWithStatus(t testing.TB, store *store.Store, j models.JobSpec, status models.RunStatus) models.JobRun {
	initr := j.Initiators[0]
	run := j.NewRun(initr)
	run.Status = status
	require.NoError(t, store.CreateJobRun(&run))
	return run
}

func BuildInitiatorRequests(t *testing.T, initrs []models.Initiator) []models.InitiatorRequest {
	bytes, err := json.Marshal(initrs)
	require.NoError(t, err)

	var dst []models.InitiatorRequest
	err = json.Unmarshal(bytes, &dst)
	require.NoError(t, err)
	return dst
}

func BuildTaskRequests(t *testing.T, initrs []models.TaskSpec) []models.TaskSpecRequest {
	bytes, err := json.Marshal(initrs)
	require.NoError(t, err)

	var dst []models.TaskSpecRequest
	err = json.Unmarshal(bytes, &dst)
	require.NoError(t, err)
	return dst
}

// CreateServiceAgreementViaWeb creates a service agreement from a fixture using /v2/service_agreements
func CreateServiceAgreementViaWeb(
	t *testing.T,
	app *TestApplication,
	path string,
	endAt time.Time,
) models.ServiceAgreement {
	client := app.NewHTTPClient()

	agreementWithoutOracle := MustJSONSet(t, string(MustReadFile(t, path)), "endAt", utils.ISO8601UTC(endAt))
	from := GetAccountAddress(t, app.ChainlinkApplication.GetStore())
	agreementWithOracle := MustJSONSet(t, agreementWithoutOracle, "oracles", []string{from.Hex()})

	resp, cleanup := client.Post("/v2/service_agreements", bytes.NewBufferString(agreementWithOracle))
	defer cleanup()

	AssertServerResponse(t, resp, http.StatusOK)
	responseSA := models.ServiceAgreement{}
	err := ParseJSONAPIResponse(t, resp, &responseSA)
	require.NoError(t, err)

	return FindServiceAgreement(t, app.Store, responseSA.ID)
}

func NewRunInput(value models.JSON) models.RunInput {
	jobRunID := models.NewID()
	return *models.NewRunInput(jobRunID, value, models.RunStatusUnstarted)
}

func NewRunInputWithString(t testing.TB, value string) models.RunInput {
	jobRunID := models.NewID()
	data := JSONFromString(t, value)
	return *models.NewRunInput(jobRunID, data, models.RunStatusUnstarted)
}

func NewRunInputWithResult(value interface{}) models.RunInput {
	jobRunID := models.NewID()
	return *models.NewRunInputWithResult(jobRunID, value, models.RunStatusUnstarted)
}
