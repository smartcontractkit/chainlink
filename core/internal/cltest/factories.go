package cltest

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"flag"
	"fmt"
	"math/big"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/adapters"
	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/fluxmonitor"
	strpkg "github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"github.com/urfave/cli"
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

// NewJobWithExternalInitiator creates new Job with external initiator
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

// NewJobWithWebInitiator create new Job with web initiator
func NewJobWithWebInitiator() models.JobSpec {
	j := NewJob()
	j.Initiators = []models.Initiator{{
		JobSpecID: j.ID,
		Type:      models.InitiatorWeb,
	}}
	return j
}

// NewJobWithLogInitiator create new Job with ethlog initiator
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

// NewJobWithRunAtInitiator create new Job with RunAt initiator
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

// NewJobWithFluxMonitorInitiator create new Job with FluxMonitor initiator
func NewJobWithFluxMonitorInitiator() models.JobSpec {
	j := NewJob()
	j.Initiators = []models.Initiator{{
		JobSpecID: j.ID,
		Type:      models.InitiatorFluxMonitor,
		InitiatorParams: models.InitiatorParams{
			Address:           NewAddress(),
			RequestData:       models.JSON{Result: gjson.Parse(`{"data":{"coin":"ETH","market":"USD"}}`)},
			Feeds:             models.JSON{Result: gjson.Parse(`["https://lambda.staging.devnet.tools/bnc/call"]`)},
			Threshold:         0.5,
			AbsoluteThreshold: 0.01,
			IdleTimer: models.IdleTimerConfig{
				Duration: models.MustMakeDuration(time.Minute),
			},
			PollTimer: models.PollTimerConfig{
				Period: models.MustMakeDuration(time.Minute),
			},
			Precision: 2,
		},
	}}
	return j
}

// NewJobWithFluxMonitorInitiator create new Job with FluxMonitor initiator
func NewJobWithFluxMonitorInitiatorWithBridge() models.JobSpec {
	j := NewJob()
	j.Initiators = []models.Initiator{{
		JobSpecID: j.ID,
		Type:      models.InitiatorFluxMonitor,
		InitiatorParams: models.InitiatorParams{
			Address:           NewAddress(),
			RequestData:       models.JSON{Result: gjson.Parse(`{"data":{"coin":"ETH","market":"USD"}}`)},
			Feeds:             models.JSON{Result: gjson.Parse(`[{"bridge":"testbridge"}]`)},
			Threshold:         0.5,
			AbsoluteThreshold: 0.01,
			Precision:         2,
		},
	}}
	return j
}

// NewJobWithRandomnessLog create new Job with VRF initiator
func NewJobWithRandomnessLog() models.JobSpec {
	j := NewJob()
	j.Initiators = []models.Initiator{{
		JobSpecID: j.ID,
		Type:      models.InitiatorRandomnessLog,
		InitiatorParams: models.InitiatorParams{
			Address: NewAddress(),
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

func NewTransaction(nonce uint64, sentAtV ...uint64) *models.Tx {
	from := common.HexToAddress("0xf208000000000000000000000000000000000000")
	to := common.HexToAddress("0x7000000000000000000000000000000000000000")

	value := new(big.Int).Exp(big.NewInt(10), big.NewInt(36), nil)
	gasLimit := uint64(50000)
	data := hexutil.MustDecode("0xda7ada7a")

	sentAt := uint64(0)
	if len(sentAtV) > 0 {
		sentAt = sentAtV[0]
	}

	transaction := types.NewTransaction(nonce, to, value, gasLimit, new(big.Int), data)
	return &models.Tx{
		From:        from,
		SentAt:      sentAt,
		To:          *transaction.To(),
		Nonce:       transaction.Nonce(),
		Data:        transaction.Data(),
		Value:       utils.NewBig(transaction.Value()),
		GasLimit:    transaction.Gas(),
		GasPrice:    utils.NewBig(transaction.GasPrice()),
		Hash:        transaction.Hash(),
		SignedRawTx: hexutil.MustDecode("0xcafe11"),
	}
}

// CreateTx creates a Tx from a specified address, and sentAt
func CreateTx(
	t testing.TB,
	store *strpkg.Store,
	from common.Address,
	sentAt uint64,
) *models.Tx {
	return CreateTxWithNonceAndGasPrice(t, store, from, sentAt, 0, 1)
}

// CreateTxWithNonceAndGasPrice creates a Tx from a specified address, sentAt, nonce and gas price
func CreateTxWithNonceAndGasPrice(
	t testing.TB,
	store *strpkg.Store,
	from common.Address,
	sentAt uint64,
	nonce uint64,
	gasPrice int64,
) *models.Tx {
	return CreateTxWithNonceGasPriceAndRecipient(t, store, from, common.Address{}, sentAt, nonce, gasPrice)
}

// CreateTxWithNonceGasPriceAndRecipient creates a Tx from a specified sender, recipient, sentAt, nonce and gas price
func CreateTxWithNonceGasPriceAndRecipient(
	t testing.TB,
	store *strpkg.Store,
	from common.Address,
	to common.Address,
	sentAt uint64,
	nonce uint64,
	gasPrice int64,
) *models.Tx {
	data := make([]byte, 36)
	binary.LittleEndian.PutUint64(data, sentAt)

	transaction := types.NewTransaction(nonce, to, big.NewInt(0), 250000, big.NewInt(gasPrice), data)
	tx := &models.Tx{
		From:        from,
		SentAt:      sentAt,
		To:          *transaction.To(),
		Nonce:       transaction.Nonce(),
		Data:        transaction.Data(),
		Value:       utils.NewBig(transaction.Value()),
		GasLimit:    transaction.Gas(),
		GasPrice:    utils.NewBig(transaction.GasPrice()),
		Hash:        transaction.Hash(),
		SignedRawTx: hexutil.MustDecode("0xcafe22"),
	}

	tx, err := store.CreateTx(tx)
	require.NoError(t, err)
	_, err = store.AddTxAttempt(tx, tx)
	require.NoError(t, err)
	return tx
}

func AddTxAttempt(
	t testing.TB,
	store *strpkg.Store,
	tx *models.Tx,
	etx *types.Transaction,
	blkNum uint64,
) *models.TxAttempt {
	transaction := types.NewTransaction(tx.Nonce, common.Address{}, big.NewInt(0), 250000, big.NewInt(1), tx.Data)

	newTxAttempt := &models.Tx{
		From:        tx.From,
		SentAt:      blkNum,
		To:          *transaction.To(),
		Nonce:       transaction.Nonce(),
		Data:        transaction.Data(),
		Value:       utils.NewBig(transaction.Value()),
		GasLimit:    transaction.Gas(),
		GasPrice:    utils.NewBig(transaction.GasPrice()),
		Hash:        transaction.Hash(),
		SignedRawTx: []byte{byte(len(tx.Attempts))},
	}

	txAttempt, err := store.AddTxAttempt(tx, newTxAttempt)
	require.NoError(t, err)
	return txAttempt
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

// NewRunLog create models.Log for given jobid, address, block, and json
func NewRunLog(
	t *testing.T,
	jobID *models.ID,
	emitter common.Address,
	requester common.Address,
	blk int,
	json string,
) models.Log {
	return models.Log{
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

// NewRandomnessRequestLog(t, r, emitter, blk) is a RandomnessRequest log for
// the randomness request log represented by r.
func NewRandomnessRequestLog(t *testing.T, r models.RandomnessRequestLog,
	emitter common.Address, blk int) models.Log {
	rawData, err := r.RawData()
	require.NoError(t, err)
	return models.Log{
		Address:     emitter,
		BlockNumber: uint64(blk),
		Data:        rawData,
		TxHash:      NewHash(),
		BlockHash:   NewHash(),
		Topics:      []common.Hash{models.RandomnessRequestLogTopic, r.JobID},
	}
}

func NewLink(t *testing.T, amount string) *assets.Link {
	link := assets.NewLink(0)
	link, ok := link.SetString(amount, 10)
	assert.True(t, ok)
	return link
}

func NewEth(t *testing.T, amount string) *assets.Eth {
	eth := assets.NewEth(0)
	eth, ok := eth.SetString(amount, 10)
	assert.True(t, ok)
	return eth
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
	case int: // Single case allows compiler to narrow x's type.
		return hexutil.Big(*big.NewInt(int64(x)))
	case uint32:
		return hexutil.Big(*big.NewInt(int64(x)))
	case uint64:
		return hexutil.Big(*big.NewInt(0).SetUint64(x))
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

func MustEVMUintHexFromBase10String(t *testing.T, strings ...string) string {
	var allBytes []byte
	for _, s := range strings {
		i, ok := big.NewInt(0).SetString(s, 10)
		require.True(t, ok)
		bs, err := utils.EVMWordBigInt(i)
		require.NoError(t, err)
		allBytes = append(allBytes, bs...)
	}
	return fmt.Sprintf("0x%0x", allBytes)
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

// NewJobRun returns a newly initialized job run for test purposes only
func NewJobRun(job models.JobSpec) models.JobRun {
	initiator := job.Initiators[0]
	now := time.Now()
	run := models.MakeJobRun(&job, now, &initiator, nil, &models.RunRequest{})
	return run
}

// NewJobRunPendingBridge returns a new job run in the pending bridge state
func NewJobRunPendingBridge(job models.JobSpec) models.JobRun {
	run := NewJobRun(job)
	run.SetStatus(models.RunStatusPendingBridge)
	run.TaskRuns[0].Status = models.RunStatusPendingBridge
	return run
}

// CreateJobRunWithStatus returns a new job run with the specified status that has been persisted
func CreateJobRunWithStatus(t testing.TB, store *strpkg.Store, job models.JobSpec, status models.RunStatus) models.JobRun {
	run := NewJobRun(job)
	run.SetStatus(status)
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

func NewRunInput(value models.JSON) models.RunInput {
	jobRunID := models.NewID()
	taskRunID := models.NewID()
	return *models.NewRunInput(jobRunID, *taskRunID, value, models.RunStatusUnstarted)
}

func NewRunInputWithString(t testing.TB, value string) models.RunInput {
	jobRunID := models.NewID()
	taskRunID := models.NewID()
	data := JSONFromString(t, value)
	return *models.NewRunInput(jobRunID, *taskRunID, data, models.RunStatusUnstarted)
}

func NewRunInputWithResult(value interface{}) models.RunInput {
	jobRunID := models.NewID()
	taskRunID := models.NewID()
	return *models.NewRunInputWithResult(jobRunID, *taskRunID, value, models.RunStatusUnstarted)
}

func NewRunInputWithResultAndJobRunID(value interface{}, jobRunID *models.ID) models.RunInput {
	taskRunID := models.NewID()
	return *models.NewRunInputWithResult(jobRunID, *taskRunID, value, models.RunStatusUnstarted)
}

func NewPollingDeviationChecker(t *testing.T, s *strpkg.Store) *fluxmonitor.PollingDeviationChecker {
	fluxAggregator := new(mocks.FluxAggregator)
	runManager := new(mocks.RunManager)
	fetcher := new(mocks.Fetcher)
	initr := models.Initiator{
		JobSpecID: models.NewID(),
		InitiatorParams: models.InitiatorParams{
			PollTimer: models.PollTimerConfig{
				Period: models.MustMakeDuration(time.Second),
			},
		},
	}
	checker, err := fluxmonitor.NewPollingDeviationChecker(s, fluxAggregator, initr, nil, runManager, fetcher, func() {})
	require.NoError(t, err)
	return checker
}

func NewEthHeader(height interface{}) *types.Header {
	var h int64
	switch v := height.(type) {
	case int64:
		h = v
	case int:
		h = int64(v)
	case uint64:
		h = int64(v)
	default:
		panic(fmt.Sprintf("invalid type: %t", height))
	}

	return &types.Header{
		Number:     big.NewInt(h),
		ParentHash: NewHash(),
		Time:       uint64(time.Now().Unix()),
	}
}

func MustInsertTaskRun(t *testing.T, store *strpkg.Store) models.ID {
	taskRunID := models.NewID()

	job := NewJobWithWebInitiator()
	require.NoError(t, store.CreateJob(&job))
	jobRun := NewJobRun(job)
	jobRun.TaskRuns = []models.TaskRun{models.TaskRun{ID: taskRunID, Status: models.RunStatusUnstarted, TaskSpecID: job.Tasks[0].ID}}
	require.NoError(t, store.CreateJobRun(&jobRun))

	return *taskRunID
}

func MustInsertKey(t *testing.T, store *strpkg.Store, address common.Address) models.Key {
	a, err := models.NewEIP55Address(address.Hex())
	require.NoError(t, err)
	key := models.Key{
		Address: a,
		JSON:    JSONFromString(t, "{}"),
	}
	require.NoError(t, store.DB.Save(&key).Error)
	return key
}

func NewEthTx(t *testing.T, store *strpkg.Store) models.EthTx {
	return models.EthTx{
		FromAddress:    GetDefaultFromAddress(t, store),
		ToAddress:      NewAddress(),
		EncodedPayload: []byte{1, 2, 3},
		Value:          assets.NewEthValue(142),
		GasLimit:       uint64(1000000000),
	}
}

func MustInsertUnconfirmedEthTxWithBroadcastAttempt(t *testing.T, store *strpkg.Store, nonce int64) models.EthTx {
	timeNow := time.Now()
	etx := NewEthTx(t, store)

	etx.BroadcastAt = &timeNow
	n := nonce
	etx.Nonce = &n
	etx.State = models.EthTxUnconfirmed
	require.NoError(t, store.DB.Save(&etx).Error)
	attempt := NewEthTxAttempt(t, etx.ID)

	tx := types.NewTransaction(uint64(nonce), NewAddress(), big.NewInt(142), 242, big.NewInt(342), []byte{1, 2, 3})
	rlp := new(bytes.Buffer)
	require.NoError(t, tx.EncodeRLP(rlp))
	attempt.SignedRawTx = rlp.Bytes()

	attempt.State = models.EthTxAttemptBroadcast
	require.NoError(t, store.DB.Save(&attempt).Error)
	etx, err := store.FindEthTxWithAttempts(etx.ID)
	require.NoError(t, err)
	return etx
}

func MustInsertConfirmedEthTxWithAttempt(t *testing.T, store *strpkg.Store, nonce int64, broadcastBeforeBlockNum int64) models.EthTx {
	timeNow := time.Now()
	etx := NewEthTx(t, store)

	etx.BroadcastAt = &timeNow
	etx.Nonce = &nonce
	etx.State = models.EthTxConfirmed
	require.NoError(t, store.DB.Save(&etx).Error)
	attempt := NewEthTxAttempt(t, etx.ID)
	attempt.BroadcastBeforeBlockNum = &broadcastBeforeBlockNum
	attempt.State = models.EthTxAttemptBroadcast
	require.NoError(t, store.DB.Save(&attempt).Error)
	etx.EthTxAttempts = append(etx.EthTxAttempts, attempt)
	return etx
}

func GetDefaultFromAddress(t *testing.T, store *strpkg.Store) common.Address {
	keys, err := store.Keys()
	require.NoError(t, err)
	require.Len(t, keys, 1)
	key := keys[0]
	return key.Address.Address()
}

func NewEthTxAttempt(t *testing.T, etxID int64) models.EthTxAttempt {
	gasPrice := utils.NewBig(big.NewInt(1))
	return models.EthTxAttempt{
		EthTxID:  etxID,
		GasPrice: *gasPrice,
		// Just a random signed raw tx that decodes correctly
		// Ignore all actual values
		SignedRawTx: hexutil.MustDecode("0xf889808504a817c8008307a12094000000000000000000000000000000000000000080a400000000000000000000000000000000000000000000000000000000000000000000000025a0838fe165906e2547b9a052c099df08ec891813fea4fcdb3c555362285eb399c5a070db99322490eb8a0f2270be6eca6e3aedbc49ff57ef939cf2774f12d08aa85e"),
		Hash:        NewHash(),
	}
}

func MustInsertBroadcastEthTxAttempt(t *testing.T, etxID int64, store *strpkg.Store, gasPrice int64) models.EthTxAttempt {
	attempt := NewEthTxAttempt(t, etxID)
	attempt.State = models.EthTxAttemptBroadcast
	attempt.GasPrice = *utils.NewBig(big.NewInt(gasPrice))
	require.NoError(t, store.DB.Create(&attempt).Error)
	return attempt
}

func MustInsertEthReceipt(t *testing.T, s *strpkg.Store, blockNumber int64, blockHash common.Hash, txHash common.Hash) models.EthReceipt {
	r := models.EthReceipt{
		BlockNumber:      blockNumber,
		BlockHash:        blockHash,
		TxHash:           txHash,
		TransactionIndex: uint(NewRandomInt64()),
		Receipt:          []byte(`{"foo":42}`),
	}
	require.NoError(t, s.DB.Save(&r).Error)
	return r
}

func MustInsertFatalErrorEthTx(t *testing.T, store *strpkg.Store) models.EthTx {
	etx := NewEthTx(t, store)
	errStr := "something exploded"
	etx.Error = &errStr
	etx.State = models.EthTxFatalError

	require.NoError(t, store.DB.Save(&etx).Error)
	return etx
}

func MustInsertRandomKey(t *testing.T, store *strpkg.Store) models.Key {
	k := models.Key{Address: models.EIP55Address(NewAddress().Hex()), JSON: JSONFromString(t, `{"key": "factory"}`)}
	require.NoError(t, store.UpsertKey(k))
	return k
}
