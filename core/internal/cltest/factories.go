package cltest

import (
	"bytes"
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
	p2ppeer "github.com/libp2p/go-libp2p-core/peer"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/auth"
	"github.com/smartcontractkit/chainlink/core/bridges"
	"github.com/smartcontractkit/chainlink/core/chains/evm/headtracker"
	"github.com/smartcontractkit/chainlink/core/chains/evm/txmgr"
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/flux_aggregator_wrapper"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	configtest "github.com/smartcontractkit/chainlink/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keeper"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func NewEIP55Address() ethkey.EIP55Address {
	a := testutils.NewAddress()
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

type BridgeOpts struct {
	Name string
	URL  string
}

// NewBridgeType create new bridge type given info slice
func NewBridgeType(t testing.TB, opts BridgeOpts) (*bridges.BridgeTypeAuthentication, *bridges.BridgeType) {
	btr := &bridges.BridgeTypeRequest{}

	// Must randomise default to avoid unique constraint conflicts with other parallel tests
	rnd := uuid.NewV4().String()

	if opts.Name != "" {
		btr.Name = bridges.MustParseBridgeName(opts.Name)
	} else {
		btr.Name = bridges.MustParseBridgeName(fmt.Sprintf("test_bridge_%s", rnd))
	}

	if opts.URL != "" {
		btr.URL = WebURL(t, opts.URL)
	} else {
		btr.URL = WebURL(t, fmt.Sprintf("https://bridge.example.com/api?%s", rnd))
	}

	bta, bt, err := bridges.NewBridgeType(btr)
	require.NoError(t, err)
	return bta, bt
}

// MustCreateBridge creates a bridge
// Be careful not to specify a name here unless you ABSOLUTELY need to
// This is because name is a unique index and identical names used across transactional tests will lock/deadlock
func MustCreateBridge(t testing.TB, db *sqlx.DB, opts BridgeOpts, cfg pg.QConfig) (bta *bridges.BridgeTypeAuthentication, bt *bridges.BridgeType) {
	bta, bt = NewBridgeType(t, opts)
	orm := bridges.NewORM(db, logger.TestLogger(t), cfg)
	err := orm.CreateBridgeType(bt)
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

func EmptyCLIContext() *cli.Context {
	set := flag.NewFlagSet("test", 0)
	return cli.NewContext(nil, set, nil)
}

func NewEthTx(t *testing.T, fromAddress common.Address) txmgr.EthTx {
	return txmgr.EthTx{
		FromAddress:    fromAddress,
		ToAddress:      testutils.NewAddress(),
		EncodedPayload: []byte{1, 2, 3},
		Value:          assets.NewEthValue(142),
		GasLimit:       uint32(1000000000),
		State:          txmgr.EthTxUnstarted,
	}
}

func MustInsertUnconfirmedEthTx(t *testing.T, borm txmgr.ORM, nonce int64, fromAddress common.Address, opts ...interface{}) txmgr.EthTx {
	broadcastAt := time.Now()
	chainID := &FixtureChainID
	for _, opt := range opts {
		switch v := opt.(type) {
		case time.Time:
			broadcastAt = v
		case *big.Int:
			chainID = v
		}
	}
	etx := NewEthTx(t, fromAddress)

	etx.BroadcastAt = &broadcastAt
	etx.InitialBroadcastAt = &broadcastAt
	n := nonce
	etx.Nonce = &n
	etx.State = txmgr.EthTxUnconfirmed
	etx.EVMChainID = *utils.NewBig(chainID)
	require.NoError(t, borm.InsertEthTx(&etx))
	return etx
}

func MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t *testing.T, borm txmgr.ORM, nonce int64, fromAddress common.Address, opts ...interface{}) txmgr.EthTx {
	etx := MustInsertUnconfirmedEthTx(t, borm, nonce, fromAddress, opts...)
	attempt := NewLegacyEthTxAttempt(t, etx.ID)

	tx := types.NewTransaction(uint64(nonce), testutils.NewAddress(), big.NewInt(142), 242, big.NewInt(342), []byte{1, 2, 3})
	rlp := new(bytes.Buffer)
	require.NoError(t, tx.EncodeRLP(rlp))
	attempt.SignedRawTx = rlp.Bytes()

	attempt.State = txmgr.EthTxAttemptBroadcast
	require.NoError(t, borm.InsertEthTxAttempt(&attempt))
	etx, err := borm.FindEthTxWithAttempts(etx.ID)
	require.NoError(t, err)
	return etx
}

func MustInsertUnconfirmedEthTxWithAttemptState(t *testing.T, borm txmgr.ORM, nonce int64, fromAddress common.Address, txAttemptState txmgr.EthTxAttemptState, opts ...interface{}) txmgr.EthTx {
	etx := MustInsertUnconfirmedEthTx(t, borm, nonce, fromAddress, opts...)
	attempt := NewLegacyEthTxAttempt(t, etx.ID)

	tx := types.NewTransaction(uint64(nonce), testutils.NewAddress(), big.NewInt(142), 242, big.NewInt(342), []byte{1, 2, 3})
	rlp := new(bytes.Buffer)
	require.NoError(t, tx.EncodeRLP(rlp))
	attempt.SignedRawTx = rlp.Bytes()

	attempt.State = txAttemptState
	require.NoError(t, borm.InsertEthTxAttempt(&attempt))
	etx, err := borm.FindEthTxWithAttempts(etx.ID)
	require.NoError(t, err)
	return etx
}

func MustInsertUnconfirmedEthTxWithBroadcastDynamicFeeAttempt(t *testing.T, borm txmgr.ORM, nonce int64, fromAddress common.Address, opts ...interface{}) txmgr.EthTx {
	etx := MustInsertUnconfirmedEthTx(t, borm, nonce, fromAddress, opts...)
	attempt := NewDynamicFeeEthTxAttempt(t, etx.ID)

	addr := testutils.NewAddress()
	dtx := types.DynamicFeeTx{
		ChainID:   big.NewInt(0),
		Nonce:     uint64(nonce),
		GasTipCap: big.NewInt(1),
		GasFeeCap: big.NewInt(1),
		Gas:       242,
		To:        &addr,
		Value:     big.NewInt(342),
		Data:      []byte{2, 3, 4},
	}
	tx := types.NewTx(&dtx)
	rlp := new(bytes.Buffer)
	require.NoError(t, tx.EncodeRLP(rlp))
	attempt.SignedRawTx = rlp.Bytes()

	attempt.State = txmgr.EthTxAttemptBroadcast
	require.NoError(t, borm.InsertEthTxAttempt(&attempt))
	etx, err := borm.FindEthTxWithAttempts(etx.ID)
	require.NoError(t, err)
	return etx
}

func MustInsertUnconfirmedEthTxWithInsufficientEthAttempt(t *testing.T, borm txmgr.ORM, nonce int64, fromAddress common.Address) txmgr.EthTx {
	timeNow := time.Now()
	etx := NewEthTx(t, fromAddress)

	etx.BroadcastAt = &timeNow
	etx.InitialBroadcastAt = &timeNow
	n := nonce
	etx.Nonce = &n
	etx.State = txmgr.EthTxUnconfirmed
	require.NoError(t, borm.InsertEthTx(&etx))
	attempt := NewLegacyEthTxAttempt(t, etx.ID)

	tx := types.NewTransaction(uint64(nonce), testutils.NewAddress(), big.NewInt(142), 242, big.NewInt(342), []byte{1, 2, 3})
	rlp := new(bytes.Buffer)
	require.NoError(t, tx.EncodeRLP(rlp))
	attempt.SignedRawTx = rlp.Bytes()

	attempt.State = txmgr.EthTxAttemptInsufficientEth
	require.NoError(t, borm.InsertEthTxAttempt(&attempt))
	etx, err := borm.FindEthTxWithAttempts(etx.ID)
	require.NoError(t, err)
	return etx
}

func MustInsertConfirmedMissingReceiptEthTxWithLegacyAttempt(
	t *testing.T, borm txmgr.ORM, nonce int64, broadcastBeforeBlockNum int64,
	broadcastAt time.Time, fromAddress common.Address) txmgr.EthTx {
	etx := NewEthTx(t, fromAddress)

	etx.BroadcastAt = &broadcastAt
	etx.InitialBroadcastAt = &broadcastAt
	etx.Nonce = &nonce
	etx.State = txmgr.EthTxConfirmedMissingReceipt
	require.NoError(t, borm.InsertEthTx(&etx))
	attempt := NewLegacyEthTxAttempt(t, etx.ID)
	attempt.BroadcastBeforeBlockNum = &broadcastBeforeBlockNum
	attempt.State = txmgr.EthTxAttemptBroadcast
	require.NoError(t, borm.InsertEthTxAttempt(&attempt))
	etx.EthTxAttempts = append(etx.EthTxAttempts, attempt)
	return etx
}

func MustInsertConfirmedEthTxWithLegacyAttempt(t *testing.T, borm txmgr.ORM, nonce int64, broadcastBeforeBlockNum int64, fromAddress common.Address) txmgr.EthTx {
	timeNow := time.Now()
	etx := NewEthTx(t, fromAddress)

	etx.BroadcastAt = &timeNow
	etx.InitialBroadcastAt = &timeNow
	etx.Nonce = &nonce
	etx.State = txmgr.EthTxConfirmed
	require.NoError(t, borm.InsertEthTx(&etx))
	attempt := NewLegacyEthTxAttempt(t, etx.ID)
	attempt.BroadcastBeforeBlockNum = &broadcastBeforeBlockNum
	attempt.State = txmgr.EthTxAttemptBroadcast
	require.NoError(t, borm.InsertEthTxAttempt(&attempt))
	etx.EthTxAttempts = append(etx.EthTxAttempts, attempt)
	return etx
}

func MustInsertInProgressEthTxWithAttempt(t *testing.T, borm txmgr.ORM, nonce int64, fromAddress common.Address) txmgr.EthTx {
	etx := NewEthTx(t, fromAddress)

	etx.Nonce = &nonce
	etx.State = txmgr.EthTxInProgress
	require.NoError(t, borm.InsertEthTx(&etx))
	attempt := NewLegacyEthTxAttempt(t, etx.ID)
	tx := types.NewTransaction(uint64(nonce), testutils.NewAddress(), big.NewInt(142), 242, big.NewInt(342), []byte{1, 2, 3})
	rlp := new(bytes.Buffer)
	require.NoError(t, tx.EncodeRLP(rlp))
	attempt.SignedRawTx = rlp.Bytes()
	attempt.State = txmgr.EthTxAttemptInProgress
	require.NoError(t, borm.InsertEthTxAttempt(&attempt))
	etx, err := borm.FindEthTxWithAttempts(etx.ID)
	require.NoError(t, err)
	return etx
}

func MustInsertUnstartedEthTx(t *testing.T, borm txmgr.ORM, fromAddress common.Address, opts ...interface{}) txmgr.EthTx {
	var subject uuid.NullUUID
	for _, opt := range opts {
		switch v := opt.(type) {
		case uuid.UUID:
			subject = uuid.NullUUID{UUID: v, Valid: true}
		}
	}
	etx := NewEthTx(t, fromAddress)
	etx.State = txmgr.EthTxUnstarted
	etx.Subject = subject
	require.NoError(t, borm.InsertEthTx(&etx))
	return etx
}

func NewLegacyEthTxAttempt(t *testing.T, etxID int64) txmgr.EthTxAttempt {
	gasPrice := assets.NewWeiI(1)
	return txmgr.EthTxAttempt{
		ChainSpecificGasLimit: 42,
		EthTxID:               etxID,
		GasPrice:              gasPrice,
		// Just a random signed raw tx that decodes correctly
		// Ignore all actual values
		SignedRawTx: hexutil.MustDecode("0xf889808504a817c8008307a12094000000000000000000000000000000000000000080a400000000000000000000000000000000000000000000000000000000000000000000000025a0838fe165906e2547b9a052c099df08ec891813fea4fcdb3c555362285eb399c5a070db99322490eb8a0f2270be6eca6e3aedbc49ff57ef939cf2774f12d08aa85e"),
		Hash:        utils.NewHash(),
		State:       txmgr.EthTxAttemptInProgress,
	}
}

func NewDynamicFeeEthTxAttempt(t *testing.T, etxID int64) txmgr.EthTxAttempt {
	gasTipCap := assets.NewWeiI(1)
	gasFeeCap := assets.NewWeiI(1)
	return txmgr.EthTxAttempt{
		TxType:    0x2,
		EthTxID:   etxID,
		GasTipCap: gasTipCap,
		GasFeeCap: gasFeeCap,
		// Just a random signed raw tx that decodes correctly
		// Ignore all actual values
		SignedRawTx:           hexutil.MustDecode("0xf889808504a817c8008307a12094000000000000000000000000000000000000000080a400000000000000000000000000000000000000000000000000000000000000000000000025a0838fe165906e2547b9a052c099df08ec891813fea4fcdb3c555362285eb399c5a070db99322490eb8a0f2270be6eca6e3aedbc49ff57ef939cf2774f12d08aa85e"),
		Hash:                  utils.NewHash(),
		State:                 txmgr.EthTxAttemptInProgress,
		ChainSpecificGasLimit: 42,
	}
}

func NewEthReceipt(t *testing.T, blockNumber int64, blockHash common.Hash, txHash common.Hash, status uint64) txmgr.EthReceipt {
	transactionIndex := uint(NewRandomPositiveInt64())

	receipt := evmtypes.Receipt{
		BlockNumber:      big.NewInt(blockNumber),
		BlockHash:        blockHash,
		TxHash:           txHash,
		TransactionIndex: transactionIndex,
		Status:           status,
	}

	r := txmgr.EthReceipt{
		BlockNumber:      blockNumber,
		BlockHash:        blockHash,
		TxHash:           txHash,
		TransactionIndex: transactionIndex,
		Receipt:          receipt,
	}
	return r
}

func MustInsertEthReceipt(t *testing.T, borm txmgr.ORM, blockNumber int64, blockHash common.Hash, txHash common.Hash) txmgr.EthReceipt {
	r := NewEthReceipt(t, blockNumber, blockHash, txHash, 0x1)
	require.NoError(t, borm.InsertEthReceipt(&r))
	return r
}

func MustInsertRevertedEthReceipt(t *testing.T, borm txmgr.ORM, blockNumber int64, blockHash common.Hash, txHash common.Hash) txmgr.EthReceipt {
	r := NewEthReceipt(t, blockNumber, blockHash, txHash, 0x0)
	require.NoError(t, borm.InsertEthReceipt(&r))
	return r
}

// Inserts into eth_receipts but does not update eth_txes or eth_tx_attempts
func MustInsertConfirmedEthTxWithReceipt(t *testing.T, borm txmgr.ORM, fromAddress common.Address, nonce, blockNum int64) (etx txmgr.EthTx) {
	etx = MustInsertConfirmedEthTxWithLegacyAttempt(t, borm, nonce, blockNum, fromAddress)
	MustInsertEthReceipt(t, borm, blockNum, utils.NewHash(), etx.EthTxAttempts[0].Hash)
	return etx
}

func MustInsertConfirmedEthTxBySaveFetchedReceipts(t *testing.T, borm txmgr.ORM, fromAddress common.Address, nonce int64, blockNum int64, chainID big.Int) (etx txmgr.EthTx) {
	etx = MustInsertConfirmedEthTxWithLegacyAttempt(t, borm, nonce, blockNum, fromAddress)
	receipt := evmtypes.Receipt{
		TxHash:           etx.EthTxAttempts[0].Hash,
		BlockHash:        utils.NewHash(),
		BlockNumber:      big.NewInt(nonce),
		TransactionIndex: uint(1),
	}
	borm.SaveFetchedReceipts([]evmtypes.Receipt{receipt}, chainID)
	return etx
}

func MustInsertFatalErrorEthTx(t *testing.T, borm txmgr.ORM, fromAddress common.Address) txmgr.EthTx {
	etx := NewEthTx(t, fromAddress)
	etx.Error = null.StringFrom("something exploded")
	etx.State = txmgr.EthTxFatalError

	require.NoError(t, borm.InsertEthTx(&etx))
	return etx
}

func MustAddRandomKeyToKeystore(t testing.TB, ethKeyStore keystore.Eth) (ethkey.KeyV2, common.Address) {
	t.Helper()

	k := MustGenerateRandomKey(t)
	MustAddKeyToKeystore(t, k, &FixtureChainID, ethKeyStore)

	return k, k.Address
}

func MustAddKeyToKeystore(t testing.TB, key ethkey.KeyV2, chainID *big.Int, ethKeyStore keystore.Eth) {
	t.Helper()
	ethKeyStore.XXXTestingOnlyAdd(key)
	err := ethKeyStore.Enable(key.Address, chainID)
	require.NoError(t, err)
}

// MustInsertRandomKey inserts a randomly generated (not cryptographically
// secure) key for testing
// By default it is enabled for the fixture chain
func MustInsertRandomKey(
	t testing.TB,
	keystore keystore.Eth,
	opts ...interface{},
) (ethkey.KeyV2, common.Address) {
	t.Helper()

	chainIDs := []utils.Big{*utils.NewBig(&FixtureChainID)}
	for _, opt := range opts {
		switch v := opt.(type) {
		case utils.Big:
			chainIDs[0] = v
		case []utils.Big:
			chainIDs = v
		}
	}

	key := MustGenerateRandomKey(t)
	keystore.XXXTestingOnlyAdd(key)

	for _, cid := range chainIDs {
		var nonce int64
		enabled := true
		for _, opt := range opts {
			switch v := opt.(type) {
			case int:
				nonce = int64(v)
			case int64:
				nonce = v
			case bool:
				enabled = v
			}
		}
		require.NoError(t, keystore.Enable(key.Address, cid.ToInt()))
		err := keystore.Reset(key.Address, cid.ToInt(), nonce)
		require.NoError(t, err)
		if !enabled {
			require.NoError(t, keystore.Disable(key.Address, cid.ToInt()))
		}
	}

	return key, key.Address
}

func MustInsertRandomKeyReturningState(t testing.TB,
	keystore keystore.Eth,
	opts ...interface{},
) (ethkey.State, common.Address) {
	k, address := MustInsertRandomKey(t, keystore, opts...)
	state := MustGetStateForKey(t, keystore, k)
	return state, address
}

func MustGenerateRandomKey(t testing.TB) ethkey.KeyV2 {
	key, err := ethkey.NewV2()
	require.NoError(t, err)
	return key
}

func MustGenerateRandomKeyState(t testing.TB) ethkey.State {
	return ethkey.State{Address: NewEIP55Address()}
}

func MustInsertHead(t *testing.T, db *sqlx.DB, cfg pg.QConfig, number int64) evmtypes.Head {
	h := evmtypes.NewHead(big.NewInt(number), utils.NewHash(), utils.NewHash(), 0, utils.NewBig(&FixtureChainID))
	horm := headtracker.NewORM(db, logger.TestLogger(t), cfg, FixtureChainID)

	err := horm.IdempotentInsertHead(testutils.Context(t), &h)
	require.NoError(t, err)
	return h
}

func MustInsertV2JobSpec(t *testing.T, db *sqlx.DB, transmitterAddress common.Address) job.Job {
	t.Helper()

	addr, err := ethkey.NewEIP55Address(transmitterAddress.Hex())
	require.NoError(t, err)

	pipelineSpec := pipeline.Spec{}
	err = db.Get(&pipelineSpec, `INSERT INTO pipeline_specs (dot_dag_source,created_at) VALUES ('',NOW()) RETURNING *`)
	require.NoError(t, err)

	oracleSpec := MustInsertOffchainreportingOracleSpec(t, db, addr)
	jb := job.Job{
		OCROracleSpec:   &oracleSpec,
		OCROracleSpecID: &oracleSpec.ID,
		ExternalJobID:   uuid.NewV4(),
		Type:            job.OffchainReporting,
		SchemaVersion:   1,
		PipelineSpec:    &pipelineSpec,
		PipelineSpecID:  pipelineSpec.ID,
	}

	jorm := job.NewORM(db, nil, nil, nil, nil, logger.TestLogger(t), configtest.NewTestGeneralConfig(t))
	err = jorm.InsertJob(&jb)
	require.NoError(t, err)
	return jb
}

func MustInsertOffchainreportingOracleSpec(t *testing.T, db *sqlx.DB, transmitterAddress ethkey.EIP55Address) job.OCROracleSpec {
	t.Helper()

	ocrKeyID := models.MustSha256HashFromHex(DefaultOCRKeyBundleID)
	spec := job.OCROracleSpec{}
	require.NoError(t, db.Get(&spec, `INSERT INTO ocr_oracle_specs (created_at, updated_at, contract_address, p2p_bootstrap_peers, is_bootstrap_peer, encrypted_ocr_key_bundle_id, transmitter_address, observation_timeout, blockchain_timeout, contract_config_tracker_subscribe_interval, contract_config_tracker_poll_interval, contract_config_confirmations, database_timeout, observation_grace_period, contract_transmitter_transmit_timeout) VALUES (
NOW(),NOW(),$1,'{}',false,$2,$3,0,0,0,0,0,0,0,0
) RETURNING *`, NewEIP55Address(), &ocrKeyID, &transmitterAddress))
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

func MustInsertKeeperJob(t *testing.T, db *sqlx.DB, korm keeper.ORM, from ethkey.EIP55Address, contract ethkey.EIP55Address) job.Job {
	t.Helper()

	var keeperSpec job.KeeperSpec
	err := korm.Q().Get(&keeperSpec, `INSERT INTO keeper_specs (contract_address, from_address, created_at, updated_at) VALUES ($1, $2, NOW(), NOW()) RETURNING *`, contract, from)
	require.NoError(t, err)

	var pipelineSpec pipeline.Spec
	err = korm.Q().Get(&pipelineSpec, `INSERT INTO pipeline_specs (dot_dag_source,created_at) VALUES ('',NOW()) RETURNING *`)
	require.NoError(t, err)

	jb := job.Job{
		KeeperSpec:     &keeperSpec,
		KeeperSpecID:   &keeperSpec.ID,
		ExternalJobID:  uuid.NewV4(),
		Type:           job.Keeper,
		SchemaVersion:  1,
		PipelineSpec:   &pipelineSpec,
		PipelineSpecID: pipelineSpec.ID,
	}

	cfg := configtest.NewTestGeneralConfig(t)
	tlg := logger.TestLogger(t)
	prm := pipeline.NewORM(db, tlg, cfg)
	btORM := bridges.NewORM(db, tlg, cfg)
	jrm := job.NewORM(db, nil, prm, btORM, nil, tlg, cfg)
	err = jrm.InsertJob(&jb)
	require.NoError(t, err)
	return jb
}

func MustInsertKeeperRegistry(t *testing.T, db *sqlx.DB, korm keeper.ORM, ethKeyStore keystore.Eth, keeperIndex, numKeepers, blockCountPerTurn int32) (keeper.Registry, job.Job) {
	key, _ := MustAddRandomKeyToKeystore(t, ethKeyStore)
	from := key.EIP55Address
	t.Helper()
	contractAddress := NewEIP55Address()
	job := MustInsertKeeperJob(t, db, korm, from, contractAddress)
	registry := keeper.Registry{
		ContractAddress:   contractAddress,
		BlockCountPerTurn: blockCountPerTurn,
		CheckGas:          150_000,
		FromAddress:       from,
		JobID:             job.ID,
		KeeperIndex:       keeperIndex,
		NumKeepers:        numKeepers,
		KeeperIndexMap: map[ethkey.EIP55Address]int32{
			from: keeperIndex,
		},
	}
	err := korm.UpsertRegistry(&registry)
	require.NoError(t, err)
	return registry, job
}

func MustInsertUpkeepForRegistry(t *testing.T, db *sqlx.DB, cfg keeper.Config, registry keeper.Registry) keeper.UpkeepRegistration {
	korm := keeper.NewORM(db, logger.TestLogger(t), cfg, txmgr.SendEveryStrategy{})
	mathrand.Seed(time.Now().UnixNano())
	upkeepID := utils.NewBigI(int64(mathrand.Uint32()))
	upkeep := keeper.UpkeepRegistration{
		UpkeepID:   upkeepID,
		ExecuteGas: uint32(150_000),
		Registry:   registry,
		RegistryID: registry.ID,
		CheckData:  common.Hex2Bytes("ABC123"),
	}
	positioningConstant, err := keeper.CalcPositioningConstant(upkeepID, registry.ContractAddress)
	require.NoError(t, err)
	upkeep.PositioningConstant = positioningConstant
	err = korm.UpsertUpkeep(&upkeep)
	require.NoError(t, err)
	return upkeep
}

func MustInsertPipelineRun(t *testing.T, db *sqlx.DB) (run pipeline.Run) {
	require.NoError(t, db.Get(&run, `INSERT INTO pipeline_runs (state,pipeline_spec_id,created_at) VALUES ($1, 0, NOW()) RETURNING *`, pipeline.RunStatusRunning))
	return run
}

func MustInsertPipelineRunWithStatus(t *testing.T, db *sqlx.DB, pipelineSpecID int32, status pipeline.RunStatus) (run pipeline.Run) {
	var finishedAt *time.Time
	var outputs pipeline.JSONSerializable
	var allErrors pipeline.RunErrors
	var fatalErrors pipeline.RunErrors
	now := time.Now()
	switch status {
	case pipeline.RunStatusCompleted:
		finishedAt = &now
		outputs = pipeline.JSONSerializable{
			Val:   "foo",
			Valid: true,
		}
	case pipeline.RunStatusErrored:
		finishedAt = &now
		allErrors = []null.String{null.StringFrom("oh no!")}
		fatalErrors = []null.String{null.StringFrom("oh no!")}
	case pipeline.RunStatusRunning, pipeline.RunStatusSuspended:
		// leave empty
	default:
		t.Fatalf("unknown status: %s", status)
	}
	require.NoError(t, db.Get(&run, `INSERT INTO pipeline_runs (state,pipeline_spec_id,finished_at,outputs,all_errors,fatal_errors,created_at) VALUES ($1, $2, $3, $4, $5, $6, NOW()) RETURNING *`, status, pipelineSpecID, finishedAt, outputs, allErrors, fatalErrors))
	return run
}

func MustInsertPipelineSpec(t *testing.T, db *sqlx.DB) (spec pipeline.Spec) {
	err := db.Get(&spec, `INSERT INTO pipeline_specs (dot_dag_source,created_at) VALUES ('',NOW()) RETURNING *`)
	require.NoError(t, err)
	return
}

func MustInsertUnfinishedPipelineTaskRun(t *testing.T, db *sqlx.DB, pipelineRunID int64) (tr pipeline.TaskRun) {
	/* #nosec G404 */
	require.NoError(t, db.Get(&tr, `INSERT INTO pipeline_task_runs (dot_id, pipeline_run_id, id, type, created_at) VALUES ($1,$2,$3, '', NOW()) RETURNING *`, strconv.Itoa(mathrand.Int()), pipelineRunID, uuid.NewV4()))
	return tr
}

func RandomLog(t *testing.T) types.Log {
	t.Helper()

	topics := make([]common.Hash, 4)
	for i := range topics {
		topics[i] = utils.NewHash()
	}

	return types.Log{
		Address:     testutils.NewAddress(),
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

func MustInsertExternalInitiator(t *testing.T, orm bridges.ORM) (ei bridges.ExternalInitiator) {
	return MustInsertExternalInitiatorWithOpts(t, orm, ExternalInitiatorOpts{})
}

type ExternalInitiatorOpts struct {
	NamePrefix     string
	URL            *models.WebURL
	OutgoingSecret string
	OutgoingToken  string
}

func MustInsertExternalInitiatorWithOpts(t *testing.T, orm bridges.ORM, opts ExternalInitiatorOpts) (ei bridges.ExternalInitiator) {
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
	ei.Salt = utils.NewSecret(utils.DefaultSecretSize)
	hashedSecret, err := auth.HashedSecret(token, ei.Salt)
	require.NoError(t, err)
	ei.HashedSecret = hashedSecret
	err = orm.CreateExternalInitiator(&ei)
	require.NoError(t, err)
	return ei
}
