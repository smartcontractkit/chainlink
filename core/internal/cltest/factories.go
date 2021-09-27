package cltest

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"flag"
	"fmt"
	"math/big"
	mathrand "math/rand"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/lib/pq"
	p2ppeer "github.com/libp2p/go-libp2p-core/peer"
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/auth"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/flux_aggregator_wrapper"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keeper"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	strpkg "github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/sjson"
	"github.com/urfave/cli"
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
)

// NewAddress return a random new address
func NewAddress() common.Address {
	return common.BytesToAddress(randomBytes(20))
}

func NewEIP55Address() ethkey.EIP55Address {
	a := NewAddress()
	e, err := ethkey.NewEIP55Address(a.Hex())
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
	_, _ = rand.Read(b) // Assignment for errcheck. Only used in tests so we can ignore.
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

func MustJSONMarshal(t *testing.T, val interface{}) string {
	t.Helper()
	bs, err := json.Marshal(val)
	require.NoError(t, err)
	return string(bs)
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

var (
	// RunLogTopic20190207withoutIndexes was the new RunRequest filter topic as of 2019-01-28,
	// after renaming Solidity variables, moving data version, and removing the cast of requestId to uint256
	RunLogTopic20190207withoutIndexes = utils.MustHash("OracleRequest(bytes32,address,bytes32,uint256,address,bytes4,uint256,uint256,bytes)")
)

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
		TxHash:      utils.NewHash(),
		BlockHash:   utils.NewHash(),
		Topics: []common.Hash{
			RunLogTopic20190207withoutIndexes,
			jobID,
		},
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

	payment := hexutil.MustDecode(configtest.MinimumContractPayment.ToHash().Hex())
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

func EmptyCLIContext() *cli.Context {
	set := flag.NewFlagSet("test", 0)
	return cli.NewContext(nil, set, nil)
}

func NewEthTx(t *testing.T, fromAddress common.Address) bulletprooftxmanager.EthTx {
	return bulletprooftxmanager.EthTx{
		FromAddress:    fromAddress,
		ToAddress:      NewAddress(),
		EncodedPayload: []byte{1, 2, 3},
		Value:          assets.NewEthValue(142),
		GasLimit:       uint64(1000000000),
		State:          bulletprooftxmanager.EthTxUnstarted,
	}
}

func MustInsertUnconfirmedEthTx(t *testing.T, db *gorm.DB, nonce int64, fromAddress common.Address, opts ...interface{}) bulletprooftxmanager.EthTx {
	broadcastAt := time.Now()
	for _, opt := range opts {
		switch v := opt.(type) {
		case time.Time:
			broadcastAt = v
		}
	}
	etx := NewEthTx(t, fromAddress)

	etx.BroadcastAt = &broadcastAt
	n := nonce
	etx.Nonce = &n
	etx.State = bulletprooftxmanager.EthTxUnconfirmed
	require.NoError(t, db.Save(&etx).Error)
	return etx
}

func MustInsertUnconfirmedEthTxWithBroadcastAttempt(t *testing.T, db *gorm.DB, nonce int64, fromAddress common.Address, opts ...interface{}) bulletprooftxmanager.EthTx {
	etx := MustInsertUnconfirmedEthTx(t, db, nonce, fromAddress, opts...)
	attempt := NewEthTxAttempt(t, etx.ID)

	tx := types.NewTransaction(uint64(nonce), NewAddress(), big.NewInt(142), 242, big.NewInt(342), []byte{1, 2, 3})
	rlp := new(bytes.Buffer)
	require.NoError(t, tx.EncodeRLP(rlp))
	attempt.SignedRawTx = rlp.Bytes()

	attempt.State = bulletprooftxmanager.EthTxAttemptBroadcast
	require.NoError(t, db.Save(&attempt).Error)
	etx, err := FindEthTxWithAttempts(db, etx.ID)
	require.NoError(t, err)
	return etx
}

// FindEthTxWithAttempts finds the EthTx with its attempts and receipts preloaded
func FindEthTxWithAttempts(db *gorm.DB, etxID int64) (bulletprooftxmanager.EthTx, error) {
	etx := bulletprooftxmanager.EthTx{}
	err := db.Preload("EthTxAttempts", func(db *gorm.DB) *gorm.DB {
		return db.Order("gas_price asc, id asc")
	}).Preload("EthTxAttempts.EthReceipts").First(&etx, "id = ?", &etxID).Error
	return etx, err
}

func MustInsertUnconfirmedEthTxWithInsufficientEthAttempt(t *testing.T, db *gorm.DB, nonce int64, fromAddress common.Address) bulletprooftxmanager.EthTx {
	timeNow := time.Now()
	etx := NewEthTx(t, fromAddress)

	etx.BroadcastAt = &timeNow
	n := nonce
	etx.Nonce = &n
	etx.State = bulletprooftxmanager.EthTxUnconfirmed
	require.NoError(t, db.Save(&etx).Error)
	attempt := NewEthTxAttempt(t, etx.ID)

	tx := types.NewTransaction(uint64(nonce), NewAddress(), big.NewInt(142), 242, big.NewInt(342), []byte{1, 2, 3})
	rlp := new(bytes.Buffer)
	require.NoError(t, tx.EncodeRLP(rlp))
	attempt.SignedRawTx = rlp.Bytes()

	attempt.State = bulletprooftxmanager.EthTxAttemptInsufficientEth
	require.NoError(t, db.Save(&attempt).Error)
	etx, err := FindEthTxWithAttempts(db, etx.ID)
	require.NoError(t, err)
	return etx
}

func MustInsertConfirmedEthTxWithAttempt(t *testing.T, db *gorm.DB, nonce int64, broadcastBeforeBlockNum int64, fromAddress common.Address) bulletprooftxmanager.EthTx {
	timeNow := time.Now()
	etx := NewEthTx(t, fromAddress)

	etx.BroadcastAt = &timeNow
	etx.Nonce = &nonce
	etx.State = bulletprooftxmanager.EthTxConfirmed
	require.NoError(t, db.Save(&etx).Error)
	attempt := NewEthTxAttempt(t, etx.ID)
	attempt.BroadcastBeforeBlockNum = &broadcastBeforeBlockNum
	attempt.State = bulletprooftxmanager.EthTxAttemptBroadcast
	require.NoError(t, db.Save(&attempt).Error)
	etx.EthTxAttempts = append(etx.EthTxAttempts, attempt)
	return etx
}

func MustInsertInProgressEthTxWithAttempt(t *testing.T, db *gorm.DB, nonce int64, fromAddress common.Address) bulletprooftxmanager.EthTx {
	etx := NewEthTx(t, fromAddress)

	etx.BroadcastAt = nil
	etx.Nonce = &nonce
	etx.State = bulletprooftxmanager.EthTxInProgress
	require.NoError(t, db.Save(&etx).Error)
	attempt := NewEthTxAttempt(t, etx.ID)
	tx := types.NewTransaction(uint64(nonce), NewAddress(), big.NewInt(142), 242, big.NewInt(342), []byte{1, 2, 3})
	rlp := new(bytes.Buffer)
	require.NoError(t, tx.EncodeRLP(rlp))
	attempt.SignedRawTx = rlp.Bytes()
	attempt.State = bulletprooftxmanager.EthTxAttemptInProgress
	require.NoError(t, db.Save(&attempt).Error)
	etx, err := FindEthTxWithAttempts(db, etx.ID)
	require.NoError(t, err)
	return etx
}

func MustInsertUnstartedEthTx(t *testing.T, db *gorm.DB, fromAddress common.Address, opts ...interface{}) bulletprooftxmanager.EthTx {
	var subject uuid.NullUUID
	for _, opt := range opts {
		switch v := opt.(type) {
		case uuid.UUID:
			subject = uuid.NullUUID{UUID: v, Valid: true}
		}
	}
	etx := NewEthTx(t, fromAddress)
	etx.State = bulletprooftxmanager.EthTxUnstarted
	etx.Subject = subject
	require.NoError(t, db.Save(&etx).Error)
	return etx
}

func NewEthTxAttempt(t *testing.T, etxID int64) bulletprooftxmanager.EthTxAttempt {
	gasPrice := utils.NewBig(big.NewInt(1))
	return bulletprooftxmanager.EthTxAttempt{
		EthTxID:  etxID,
		GasPrice: *gasPrice,
		// Just a random signed raw tx that decodes correctly
		// Ignore all actual values
		SignedRawTx: hexutil.MustDecode("0xf889808504a817c8008307a12094000000000000000000000000000000000000000080a400000000000000000000000000000000000000000000000000000000000000000000000025a0838fe165906e2547b9a052c099df08ec891813fea4fcdb3c555362285eb399c5a070db99322490eb8a0f2270be6eca6e3aedbc49ff57ef939cf2774f12d08aa85e"),
		Hash:        utils.NewHash(),
		State:       bulletprooftxmanager.EthTxAttemptInProgress,
	}
}

func MustInsertBroadcastEthTxAttempt(t *testing.T, etxID int64, db *gorm.DB, gasPrice int64) bulletprooftxmanager.EthTxAttempt {
	attempt := NewEthTxAttempt(t, etxID)
	attempt.State = bulletprooftxmanager.EthTxAttemptBroadcast
	attempt.GasPrice = *utils.NewBig(big.NewInt(gasPrice))
	require.NoError(t, db.Create(&attempt).Error)
	return attempt
}

func MustInsertEthReceipt(t *testing.T, db *gorm.DB, blockNumber int64, blockHash common.Hash, txHash common.Hash) bulletprooftxmanager.EthReceipt {
	transactionIndex := uint(NewRandomInt64())

	receipt := bulletprooftxmanager.Receipt{
		BlockNumber:      big.NewInt(blockNumber),
		BlockHash:        blockHash,
		TxHash:           txHash,
		TransactionIndex: transactionIndex,
	}

	data, err := json.Marshal(receipt)
	require.NoError(t, err)

	r := bulletprooftxmanager.EthReceipt{
		BlockNumber:      blockNumber,
		BlockHash:        blockHash,
		TxHash:           txHash,
		TransactionIndex: transactionIndex,
		Receipt:          data,
	}
	require.NoError(t, db.Save(&r).Error)
	return r
}

func MustInsertConfirmedEthTxWithReceipt(t *testing.T, db *gorm.DB, fromAddress common.Address, nonce, blockNum int64) (etx bulletprooftxmanager.EthTx) {
	etx = MustInsertConfirmedEthTxWithAttempt(t, db, nonce, blockNum, fromAddress)
	MustInsertEthReceipt(t, db, blockNum, utils.NewHash(), etx.EthTxAttempts[0].Hash)
	return etx
}

func MustInsertFatalErrorEthTx(t *testing.T, db *gorm.DB, fromAddress common.Address) bulletprooftxmanager.EthTx {
	etx := NewEthTx(t, fromAddress)
	etx.Error = null.StringFrom("something exploded")
	etx.State = bulletprooftxmanager.EthTxFatalError

	require.NoError(t, db.Save(&etx).Error)
	return etx
}

func MustAddRandomKeyToKeystore(t testing.TB, ethKeyStore keystore.Eth) (ethkey.KeyV2, common.Address) {
	t.Helper()

	k := MustGenerateRandomKey(t)
	MustAddKeyToKeystore(t, k, ethKeyStore)

	return k, k.Address.Address()
}

func MustAddKeyToKeystore(t testing.TB, key ethkey.KeyV2, ethKeyStore keystore.Eth) {
	t.Helper()
	err := ethKeyStore.Add(key)
	require.NoError(t, err)
}

// MustInsertRandomKey inserts a randomly generated (not cryptographically
// secure) key for testing
// If using this with the keystore, it should be called before the keystore loads keys from the database
func MustInsertRandomKey(
	t testing.TB,
	keystore keystore.Eth,
	opts ...interface{},
) (ethkey.KeyV2, common.Address) {
	t.Helper()
	key := MustGenerateRandomKey(t)
	require.NoError(t, keystore.Add(key))
	state, err := keystore.GetState(key.ID())
	require.NoError(t, err)

	for _, opt := range opts {
		switch v := opt.(type) {
		case int:
			state.NextNonce = int64(v)
		case int64:
			state.NextNonce = v
		case bool:
			state.IsFunding = v
		default:
			t.Fatalf("unrecognised option type: %T", v)
		}
	}
	err = keystore.SetState(state)
	require.NoError(t, err)

	return key, key.Address.Address()
}

func MustGenerateRandomKey(t testing.TB) ethkey.KeyV2 {
	key, err := ethkey.NewV2()
	require.NoError(t, err)
	return key
}

func MustInsertHead(t *testing.T, store *strpkg.Store, number int64) models.Head {
	h := models.NewHead(big.NewInt(number), utils.NewHash(), utils.NewHash(), 0)
	err := store.DB.Create(&h).Error
	require.NoError(t, err)
	return h
}

func MustInsertV2JobSpec(t *testing.T, store *strpkg.Store, transmitterAddress common.Address) job.Job {
	t.Helper()

	addr, err := ethkey.NewEIP55Address(transmitterAddress.Hex())
	require.NoError(t, err)

	pipelineSpec := pipeline.Spec{}
	err = store.DB.Create(&pipelineSpec).Error
	require.NoError(t, err)

	oracleSpec := MustInsertOffchainreportingOracleSpec(t, store.DB, addr)
	jb := job.Job{
		OffchainreportingOracleSpec:   &oracleSpec,
		OffchainreportingOracleSpecID: &oracleSpec.ID,
		ExternalJobID:                 uuid.NewV4(),
		Type:                          job.OffchainReporting,
		SchemaVersion:                 1,
		PipelineSpec:                  &pipelineSpec,
		PipelineSpecID:                pipelineSpec.ID,
	}

	err = store.DB.Create(&jb).Error
	require.NoError(t, err)
	return jb
}

func MustInsertOffchainreportingOracleSpec(t *testing.T, db *gorm.DB, transmitterAddress ethkey.EIP55Address) job.OffchainReportingOracleSpec {
	t.Helper()

	pid := p2pkey.PeerID(DefaultP2PPeerID)
	ocrKeyID := models.MustSha256HashFromHex(DefaultOCRKeyBundleID)
	spec := job.OffchainReportingOracleSpec{
		ContractAddress:                        NewEIP55Address(),
		P2PPeerID:                              &pid,
		P2PBootstrapPeers:                      pq.StringArray{},
		IsBootstrapPeer:                        false,
		EncryptedOCRKeyBundleID:                &ocrKeyID,
		TransmitterAddress:                     &transmitterAddress,
		ObservationTimeout:                     0,
		BlockchainTimeout:                      0,
		ContractConfigTrackerSubscribeInterval: 0,
		ContractConfigTrackerPollInterval:      0,
		ContractConfigConfirmations:            0,
	}
	require.NoError(t, db.Create(&spec).Error)
	return spec
}

func MakeDirectRequestJobSpec(t *testing.T) *job.Job {
	t.Helper()
	drs := &job.DirectRequestSpec{}
	spec := &job.Job{
		Type:              job.DirectRequest,
		SchemaVersion:     1,
		ExternalJobID:     uuid.NewV4(),
		DirectRequestSpec: drs,
		Pipeline:          pipeline.Pipeline{},
		PipelineSpec:      &pipeline.Spec{},
	}
	return spec
}

func MustInsertKeeperJob(t *testing.T, store *strpkg.Store, from ethkey.EIP55Address, contract ethkey.EIP55Address) job.Job {
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
		ExternalJobID:  uuid.NewV4(),
		Type:           job.Keeper,
		SchemaVersion:  1,
		PipelineSpec:   &pipelineSpec,
		PipelineSpecID: pipelineSpec.ID,
	}
	err = store.DB.Create(&specDB).Error
	require.NoError(t, err)
	return specDB
}

func MustInsertKeeperRegistry(t *testing.T, store *strpkg.Store, ethKeyStore keystore.Eth) (keeper.Registry, job.Job) {
	key, _ := MustAddRandomKeyToKeystore(t, ethKeyStore)
	from := key.Address
	t.Helper()
	contractAddress := NewEIP55Address()
	job := MustInsertKeeperJob(t, store, from, contractAddress)
	registry := keeper.Registry{
		ContractAddress:   contractAddress,
		BlockCountPerTurn: 20,
		CheckGas:          150_000,
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
	upkeepID, err := keeper.NewORM(store.DB, nil, store.Config, bulletprooftxmanager.SendEveryStrategy{}).LowestUnsyncedID(ctx, registry.ID)
	require.NoError(t, err)
	upkeep := keeper.UpkeepRegistration{
		UpkeepID:   upkeepID,
		ExecuteGas: uint64(150_000),
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
		AvailableFunds:   configtest.MinimumContractPayment.ToInt(),
		PaymentAmount:    configtest.MinimumContractPayment.ToInt(),
	}
}

func MustInsertPipelineRun(t *testing.T, db *gorm.DB) pipeline.Run {
	run := pipeline.Run{
		State:      pipeline.RunStatusRunning,
		Outputs:    pipeline.JSONSerializable{Null: true},
		Errors:     pipeline.RunErrors{},
		FinishedAt: null.Time{},
	}
	require.NoError(t, db.Create(&run).Error)
	return run
}

func MustInsertUnfinishedPipelineTaskRun(t *testing.T, store *strpkg.Store, pipelineRunID int64) pipeline.TaskRun {
	/* #nosec G404 */
	p := pipeline.TaskRun{DotID: strconv.Itoa(mathrand.Int()), PipelineRunID: pipelineRunID, ID: uuid.NewV4()}
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
		topics[i] = utils.NewHash()
	}

	return types.Log{
		Address:     NewAddress(),
		BlockHash:   utils.NewHash(),
		BlockNumber: uint64(mathrand.Intn(9999999)),
		Index:       uint(mathrand.Intn(9999999)),
		Data:        MustRandomBytes(t, 512),
		Topics:      []common.Hash{utils.NewHash(), utils.NewHash(), utils.NewHash(), utils.NewHash()},
	}
}

func RawNewRoundLog(t *testing.T, contractAddr common.Address, blockHash common.Hash, blockNumber uint64, logIndex uint, removed bool) types.Log {
	t.Helper()
	topic := (flux_aggregator_wrapper.FluxAggregatorNewRound{}).Topic()
	topics := []common.Hash{topic, utils.NewHash(), utils.NewHash()}
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

func MustInsertExternalInitiator(t *testing.T, db *gorm.DB) (ei models.ExternalInitiator) {
	return MustInsertExternalInitiatorWithOpts(t, db, ExternalInitiatorOpts{})
}

type ExternalInitiatorOpts struct {
	NamePrefix     string
	URL            *models.WebURL
	OutgoingSecret string
	OutgoingToken  string
}

func MustInsertExternalInitiatorWithOpts(t *testing.T, db *gorm.DB, opts ExternalInitiatorOpts) (ei models.ExternalInitiator) {
	var prefix string
	if opts.NamePrefix != "" {
		prefix = opts.NamePrefix
	} else {
		prefix = "ei"
	}
	ei.Name = fmt.Sprintf("%s-%s", prefix, uuid.NewV4())
	ei.URL = opts.URL
	ei.OutgoingSecret = opts.OutgoingSecret
	ei.OutgoingToken = opts.OutgoingToken
	token := auth.NewToken()
	ei.AccessKey = token.AccessKey
	err := db.Create(&ei).Error
	require.NoError(t, err)
	return ei
}
