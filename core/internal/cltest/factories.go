package cltest

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/json"
	"flag"
	"fmt"
	"math/big"
	mathrand "math/rand"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/lib/pq"
	"github.com/smartcontractkit/chainlink/core/adapters"

	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keeper"
	"github.com/smartcontractkit/chainlink/core/services/postgres"

	p2ppeer "github.com/libp2p/go-libp2p-core/peer"
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/flux_aggregator_wrapper"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/fluxmonitor"
	logmocks "github.com/smartcontractkit/chainlink/core/services/log/mocks"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	strpkg "github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
	"gorm.io/gorm"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	googleuuid "github.com/google/uuid"
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
func NewJobWithFluxMonitorInitiatorWithBridge(bridgeName string) models.JobSpec {
	j := NewJob()
	j.Initiators = []models.Initiator{{
		JobSpecID: j.ID,
		Type:      models.InitiatorFluxMonitor,
		InitiatorParams: models.InitiatorParams{
			Address:           NewAddress(),
			RequestData:       models.JSON{Result: gjson.Parse(`{"data":{"coin":"ETH","market":"USD"}}`)},
			Feeds:             models.JSON{Result: gjson.Parse(fmt.Sprintf("[{\"bridge\":\"%s\"}]", bridgeName))},
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

// NewHash return random Keccak256
func NewHash() common.Hash {
	return common.BytesToHash(randomBytes(32))
}

// NewAddress return a random new address
func NewAddress() common.Address {
	return common.BytesToAddress(randomBytes(20))
}

func NewEIP55Address() models.EIP55Address {
	a := NewAddress()
	e, err := models.NewEIP55Address(a.Hex())
	if err != nil {
		panic(err)
	}
	return e
}

func NewPeerID() p2ppeer.ID {
	id, err := p2ppeer.Decode("12D3KooWL3XJ9EMCyZvmmGXL2LMiVBtrVa2BuESsJiXkSj7333Jw")
	if err != nil {
		panic(err)
	}
	return id
}

func randomBytes(n int) []byte {
	b := make([]byte, n)
	rand.Read(b)
	return b
}

func Random32Byte() (b [32]byte) {
	copy(b[:], randomBytes(32))
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

// NewRunLog create types.Log for given jobid, address, block, and json
func NewRunLog(
	t *testing.T,
	jobID common.Hash,
	emitter common.Address,
	requester common.Address,
	blk int,
	json string,
) types.Log {
	return types.Log{
		Address:     emitter,
		BlockNumber: uint64(blk),
		Data:        StringToVersionedLogData20190207withoutIndexes(t, "internalID", requester, json),
		TxHash:      NewHash(),
		BlockHash:   NewHash(),
		Topics: []common.Hash{
			models.RunLogTopic20190207withoutIndexes,
			jobID,
		},
	}
}

// NewRandomnessRequestLog(t, r, emitter, blk) is a RandomnessRequest log for
// the randomness request log represented by r.
func NewRandomnessRequestLog(t *testing.T, r models.RandomnessRequestLog,
	emitter common.Address, blk int) types.Log {
	rawData, err := r.RawData()
	require.NoError(t, err)
	return types.Log{
		Address:     emitter,
		BlockNumber: uint64(blk),
		Data:        rawData,
		TxHash:      NewHash(),
		BlockHash:   NewHash(),
		Topics:      []common.Hash{models.RandomnessRequestLogTopic, r.JobID},
	}
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

func NewRunInputWithString(t testing.TB, value string) models.RunInput {
	taskRunID := uuid.NewV4()
	data := JSONFromString(t, value)
	jr := NewJobRun(NewJobWithRunLogInitiator())
	return *models.NewRunInput(jr, taskRunID, data, models.RunStatusUnstarted)
}

func NewRunInputWithResult(value interface{}) models.RunInput {
	jr := NewJobRun(NewJobWithRunLogInitiator())
	taskRunID := uuid.NewV4()
	return *models.NewRunInputWithResult(jr, taskRunID, value, models.RunStatusUnstarted)
}

func NewPollingDeviationChecker(t *testing.T, s *strpkg.Store) *fluxmonitor.PollingDeviationChecker {
	fluxAggregator := new(mocks.FluxAggregator)
	runManager := new(mocks.RunManager)
	fetcher := new(mocks.Fetcher)
	initr := models.Initiator{
		JobSpecID: models.NewJobID(),
		InitiatorParams: models.InitiatorParams{
			PollTimer: models.PollTimerConfig{
				Period: models.MustMakeDuration(time.Second),
			},
		},
	}
	lb := new(logmocks.Broadcaster)
	checker, err := fluxmonitor.NewPollingDeviationChecker(s, fluxAggregator, nil, lb, initr, nil, runManager, fetcher, big.NewInt(0), big.NewInt(100000000000))
	require.NoError(t, err)
	return checker
}

func MustInsertTaskRun(t *testing.T, store *strpkg.Store) (uuid.UUID, models.JobRun) {
	taskRunID := uuid.NewV4()

	job := NewJobWithWebInitiator()
	require.NoError(t, store.CreateJob(&job))
	jobRun := NewJobRun(job)
	jobRun.TaskRuns = []models.TaskRun{{ID: taskRunID, Status: models.RunStatusUnstarted, TaskSpecID: job.Tasks[0].ID}}
	require.NoError(t, store.CreateJobRun(&jobRun))

	return taskRunID, jobRun
}

func NewEthTx(t *testing.T, store *strpkg.Store, fromAddress common.Address) models.EthTx {
	return models.EthTx{
		FromAddress:    fromAddress,
		ToAddress:      NewAddress(),
		EncodedPayload: []byte{1, 2, 3},
		Value:          assets.NewEthValue(142),
		GasLimit:       uint64(1000000000),
		State:          models.EthTxUnstarted,
	}
}

func MustInsertUnconfirmedEthTx(t *testing.T, store *strpkg.Store, nonce int64, fromAddress common.Address, opts ...interface{}) models.EthTx {
	broadcastAt := time.Now()
	for _, opt := range opts {
		switch v := opt.(type) {
		case time.Time:
			broadcastAt = v
		}
	}
	etx := NewEthTx(t, store, fromAddress)

	etx.BroadcastAt = &broadcastAt
	n := nonce
	etx.Nonce = &n
	etx.State = models.EthTxUnconfirmed
	require.NoError(t, store.DB.Save(&etx).Error)
	return etx
}

func MustInsertUnconfirmedEthTxWithBroadcastAttempt(t *testing.T, store *strpkg.Store, nonce int64, fromAddress common.Address, opts ...interface{}) models.EthTx {
	etx := MustInsertUnconfirmedEthTx(t, store, nonce, fromAddress, opts...)
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

func MustInsertUnconfirmedEthTxWithInsufficientEthAttempt(t *testing.T, store *strpkg.Store, nonce int64, fromAddress common.Address) models.EthTx {
	timeNow := time.Now()
	etx := NewEthTx(t, store, fromAddress)

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

	attempt.State = models.EthTxAttemptInsufficientEth
	require.NoError(t, store.DB.Save(&attempt).Error)
	etx, err := store.FindEthTxWithAttempts(etx.ID)
	require.NoError(t, err)
	return etx
}

func MustInsertConfirmedEthTxWithAttempt(t *testing.T, store *strpkg.Store, nonce int64, broadcastBeforeBlockNum int64, fromAddress common.Address) models.EthTx {
	timeNow := time.Now()
	etx := NewEthTx(t, store, fromAddress)

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

func MustInsertInProgressEthTxWithAttempt(t *testing.T, store *strpkg.Store, nonce int64, fromAddress common.Address) models.EthTx {
	etx := NewEthTx(t, store, fromAddress)

	etx.BroadcastAt = nil
	etx.Nonce = &nonce
	etx.State = models.EthTxInProgress
	require.NoError(t, store.DB.Save(&etx).Error)
	attempt := NewEthTxAttempt(t, etx.ID)
	tx := types.NewTransaction(uint64(nonce), NewAddress(), big.NewInt(142), 242, big.NewInt(342), []byte{1, 2, 3})
	rlp := new(bytes.Buffer)
	require.NoError(t, tx.EncodeRLP(rlp))
	attempt.SignedRawTx = rlp.Bytes()
	attempt.State = models.EthTxAttemptInProgress
	require.NoError(t, store.DB.Save(&attempt).Error)
	etx, err := store.FindEthTxWithAttempts(etx.ID)
	require.NoError(t, err)
	return etx
}

func MustInsertUnstartedEthTx(t *testing.T, store *strpkg.Store, fromAddress common.Address) models.EthTx {
	etx := NewEthTx(t, store, fromAddress)
	etx.State = models.EthTxUnstarted
	require.NoError(t, store.DB.Save(&etx).Error)
	return etx
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
		State:       models.EthTxAttemptInProgress,
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

func MustInsertConfirmedEthTxWithReceipt(t *testing.T, s *strpkg.Store, fromAddress common.Address, nonce, blockNum int64) (etx models.EthTx) {
	etx = MustInsertConfirmedEthTxWithAttempt(t, s, nonce, blockNum, fromAddress)
	MustInsertEthReceipt(t, s, blockNum, NewHash(), etx.EthTxAttempts[0].Hash)
	return etx
}

func MustInsertFatalErrorEthTx(t *testing.T, store *strpkg.Store, fromAddress common.Address) models.EthTx {
	etx := NewEthTx(t, store, fromAddress)
	errStr := "something exploded"
	etx.Error = &errStr
	etx.State = models.EthTxFatalError

	require.NoError(t, store.DB.Save(&etx).Error)
	return etx
}

func MustAddRandomKeyToKeystore(t testing.TB, store *strpkg.Store, opts ...interface{}) (models.Key, common.Address) {
	t.Helper()

	k := MustGenerateRandomKey(t, opts...)
	err := store.KeyStore.Unlock(Password)
	require.NoError(t, err)
	MustAddKeyToKeystore(t, &k, store)
	return k, k.Address.Address()
}

func MustAddKeyToKeystore(t testing.TB, key *models.Key, store *strpkg.Store) {
	t.Helper()

	err := store.KeyStore.Unlock(Password)
	require.NoError(t, err)
	_, err = store.KeyStore.Import(key.JSON.Bytes(), Password)
	require.NoError(t, err)
	require.NoError(t, store.DB.Create(key).Error)
}

// MustInsertRandomKey inserts a randomly generated (not cryptographically
// secure) key for testing
// If using this with the keystore, it should be called before the keystore loads keys from the database
func MustInsertRandomKey(t testing.TB, db *gorm.DB, opts ...interface{}) models.Key {
	t.Helper()

	key := MustGenerateRandomKey(t, opts...)

	require.NoError(t, db.Create(&key).Error)
	return key
}

func MustGenerateRandomKey(t testing.TB, opts ...interface{}) models.Key {
	privateKeyECDSA, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
	require.NoError(t, err)
	//  < Geth 1.10 id type []byte
	//  >= Geth 1.10 id type [16]byte
	id := googleuuid.New()
	k := &keystore.Key{
		Id:         id,
		Address:    crypto.PubkeyToAddress(privateKeyECDSA.PublicKey),
		PrivateKey: privateKeyECDSA,
	}
	keyjsonbytes, err := keystore.EncryptKey(k, Password, utils.FastScryptParams.N, utils.FastScryptParams.P)
	require.NoError(t, err)
	keyjson, err := models.ParseJSON(keyjsonbytes)
	require.NoError(t, err)
	eip, err := models.EIP55AddressFromAddress(k.Address)
	require.NoError(t, err)

	var nextNonce int64
	var funding bool
	for _, opt := range opts {
		switch v := opt.(type) {
		case int:
			nextNonce = int64(v)
		case int64:
			nextNonce = v
		case bool:
			funding = v
		default:
			t.Fatalf("unrecognised option type: %T", v)
		}
	}

	key := models.Key{
		Address:   eip,
		JSON:      keyjson,
		NextNonce: nextNonce,
		IsFunding: funding,
	}
	return key
}

func MustInsertHead(t *testing.T, store *strpkg.Store, number int64) models.Head {
	h := models.NewHead(big.NewInt(number), NewHash(), NewHash(), 0)
	err := store.DB.Create(&h).Error
	require.NoError(t, err)
	return h
}

func MustInsertV2JobSpec(t *testing.T, store *strpkg.Store, transmitterAddress common.Address) job.Job {
	t.Helper()

	addr, err := models.NewEIP55Address(transmitterAddress.Hex())
	require.NoError(t, err)

	pipelineSpec := pipeline.Spec{}
	err = store.DB.Create(&pipelineSpec).Error
	require.NoError(t, err)

	oracleSpec := MustInsertOffchainreportingOracleSpec(t, store, addr)
	jb := job.Job{
		OffchainreportingOracleSpec:   &oracleSpec,
		OffchainreportingOracleSpecID: &oracleSpec.ID,
		Type:                          job.OffchainReporting,
		SchemaVersion:                 1,
		PipelineSpec:                  &pipelineSpec,
		PipelineSpecID:                pipelineSpec.ID,
	}

	err = store.DB.Create(&jb).Error
	require.NoError(t, err)
	return jb
}

func MustInsertOffchainreportingOracleSpec(t *testing.T, store *strpkg.Store, transmitterAddress models.EIP55Address) job.OffchainReportingOracleSpec {
	t.Helper()

	pid := models.PeerID(DefaultP2PPeerID)
	spec := job.OffchainReportingOracleSpec{
		ContractAddress:                        NewEIP55Address(),
		P2PPeerID:                              &pid,
		P2PBootstrapPeers:                      pq.StringArray{},
		IsBootstrapPeer:                        false,
		EncryptedOCRKeyBundleID:                &DefaultOCRKeyBundleIDSha256,
		TransmitterAddress:                     &transmitterAddress,
		ObservationTimeout:                     0,
		BlockchainTimeout:                      0,
		ContractConfigTrackerSubscribeInterval: 0,
		ContractConfigTrackerPollInterval:      0,
		ContractConfigConfirmations:            0,
	}
	require.NoError(t, store.DB.Create(&spec).Error)
	return spec
}

func MakeDirectRequestJobSpec(t *testing.T) *job.Job {
	t.Helper()
	drs := &job.DirectRequestSpec{}
	onChainJobSpecID := uuid.NewV4()
	copy(drs.OnChainJobSpecID[:], onChainJobSpecID[:])
	spec := &job.Job{
		Type:              job.DirectRequest,
		SchemaVersion:     1,
		DirectRequestSpec: drs,
		Pipeline:          *pipeline.NewTaskDAG(),
		PipelineSpec:      &pipeline.Spec{},
	}
	return spec
}

func MustInsertJobSpec(t *testing.T, s *strpkg.Store) models.JobSpec {
	j := NewJob()
	require.NoError(t, s.CreateJob(&j))
	return j
}

func MustInsertKeeperJob(t *testing.T, store *strpkg.Store, from models.EIP55Address, contract models.EIP55Address) job.Job {
	t.Helper()
	pipelineSpec := pipeline.Spec{}
	err := store.DB.Create(&pipelineSpec).Error
	require.NoError(t, err)
	keeperSpec := job.KeeperSpec{
		ContractAddress: contract,
		FromAddress:     from,
	}
	err = store.DB.Create(&keeperSpec).Error
	require.NoError(t, err)
	specDB := job.Job{
		KeeperSpec:     &keeperSpec,
		KeeperSpecID:   &keeperSpec.ID,
		Type:           job.Keeper,
		SchemaVersion:  1,
		PipelineSpec:   &pipelineSpec,
		PipelineSpecID: pipelineSpec.ID,
	}
	err = store.DB.Create(&specDB).Error
	require.NoError(t, err)
	return specDB
}

func MustInsertKeeperRegistry(t *testing.T, store *strpkg.Store) (keeper.Registry, job.Job) {
	key, _ := MustAddRandomKeyToKeystore(t, store)
	from := key.Address
	t.Helper()
	contractAddress := NewEIP55Address()
	job := MustInsertKeeperJob(t, store, from, contractAddress)
	registry := keeper.Registry{
		ContractAddress:   contractAddress,
		BlockCountPerTurn: 20,
		CheckGas:          10_000,
		FromAddress:       from,
		JobID:             job.ID,
		KeeperIndex:       0,
		NumKeepers:        1,
	}
	err := store.DB.Create(&registry).Error
	require.NoError(t, err)
	return registry, job
}

func MustInsertUpkeepForRegistry(t *testing.T, store *strpkg.Store, registry keeper.Registry) keeper.UpkeepRegistration {
	ctx, _ := postgres.DefaultQueryCtx()
	upkeepID, err := keeper.NewORM(store.DB).LowestUnsyncedID(ctx, registry)
	require.NoError(t, err)
	upkeep := keeper.UpkeepRegistration{
		UpkeepID:   upkeepID,
		ExecuteGas: int32(10_000),
		Registry:   registry,
		RegistryID: registry.ID,
		CheckData:  common.Hex2Bytes("ABC123"),
	}
	positioningConstant, err := keeper.CalcPositioningConstant(upkeepID, registry.ContractAddress)
	require.NoError(t, err)
	upkeep.PositioningConstant = positioningConstant
	err = store.DB.Create(&upkeep).Error
	require.NoError(t, err)
	return upkeep
}

func NewRoundStateForRoundID(store *strpkg.Store, roundID uint32, latestSubmission *big.Int) flux_aggregator_wrapper.OracleRoundState {
	return flux_aggregator_wrapper.OracleRoundState{
		RoundId:          roundID,
		EligibleToSubmit: true,
		LatestSubmission: latestSubmission,
		AvailableFunds:   store.Config.MinimumContractPayment().ToInt(),
		PaymentAmount:    store.Config.MinimumContractPayment().ToInt(),
	}
}

func MustInsertPipelineRun(t *testing.T, db *gorm.DB) pipeline.Run {
	run := pipeline.Run{
		Outputs:    pipeline.JSONSerializable{Null: true},
		Errors:     pipeline.RunErrors{},
		FinishedAt: nil,
	}
	require.NoError(t, db.Create(&run).Error)
	return run
}

func MustInsertUnfinishedPipelineTaskRun(t *testing.T, store *strpkg.Store, pipelineRunID int64) pipeline.TaskRun {
	/* #nosec G404 */
	p := pipeline.TaskRun{DotID: strconv.Itoa(mathrand.Int()), PipelineRunID: pipelineRunID}
	require.NoError(t, store.DB.Create(&p).Error)
	return p
}

func MustInsertSampleDirectRequestJob(t *testing.T, db *gorm.DB) job.Job {
	t.Helper()

	pspec := pipeline.Spec{DotDagSource: `
    // data source 1
    ds1          [type=bridge name=voter_turnout];
    ds1_parse    [type=jsonparse path="one,two"];
    ds1_multiply [type=multiply times=1.23];
`}

	require.NoError(t, db.Create(&pspec).Error)

	drspec := job.DirectRequestSpec{}
	require.NoError(t, db.Create(&drspec).Error)

	job := job.Job{Type: "directrequest", SchemaVersion: 1, DirectRequestSpecID: &drspec.ID, PipelineSpecID: pspec.ID}
	require.NoError(t, db.Create(&job).Error)

	return job
}

func RandomLog(t *testing.T) types.Log {
	t.Helper()

	topics := make([]common.Hash, 4)
	for i := range topics {
		topics[i] = NewHash()
	}

	return types.Log{
		Address:     NewAddress(),
		BlockHash:   NewHash(),
		BlockNumber: uint64(mathrand.Intn(9999999)),
		Index:       uint(mathrand.Intn(9999999)),
		Data:        MustRandomBytes(t, 512),
		Topics:      []common.Hash{NewHash(), NewHash(), NewHash(), NewHash()},
	}
}

func RawNewRoundLog(t *testing.T, contractAddr common.Address, blockHash common.Hash, blockNumber uint64, logIndex uint, removed bool) types.Log {
	t.Helper()
	topic := (flux_aggregator_wrapper.FluxAggregatorNewRound{}).Topic()
	topics := []common.Hash{topic, NewHash(), NewHash()}
	return RawNewRoundLogWithTopics(t, contractAddr, blockHash, blockNumber, logIndex, removed, topics)
}

func RawNewRoundLogWithTopics(t *testing.T, contractAddr common.Address, blockHash common.Hash, blockNumber uint64, logIndex uint, removed bool, topics []common.Hash) types.Log {
	t.Helper()
	return types.Log{
		Address:     contractAddr,
		BlockHash:   blockHash,
		BlockNumber: blockNumber,
		Index:       logIndex,
		Topics:      topics,
		Data:        []byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"),
		Removed:     removed,
	}
}
