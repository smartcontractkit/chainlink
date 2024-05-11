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
	"github.com/google/uuid"
	ragep2ptypes "github.com/smartcontractkit/libocr/ragep2p/types"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"
	"gopkg.in/guregu/null.v4"

	"github.com/jmoiron/sqlx"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/jsonserializable"
	txmgrcommon "github.com/smartcontractkit/chainlink/v2/common/txmgr"
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/core/auth"
	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	evmutils "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/flux_aggregator_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keeper"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func NewEIP55Address() evmtypes.EIP55Address {
	a := testutils.NewAddress()
	e, err := evmtypes.NewEIP55Address(a.Hex())
	if err != nil {
		panic(err)
	}
	return e
}

func NewPeerID() (id ragep2ptypes.PeerID) {
	err := id.UnmarshalText([]byte("12D3KooWL3XJ9EMCyZvmmGXL2LMiVBtrVa2BuESsJiXkSj7333Jw"))
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
	rnd := uuid.New().String()

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
func MustCreateBridge(t testing.TB, ds sqlutil.DataSource, opts BridgeOpts) (bta *bridges.BridgeTypeAuthentication, bt *bridges.BridgeType) {
	bta, bt = NewBridgeType(t, opts)
	orm := bridges.NewORM(ds)
	err := orm.CreateBridgeType(testutils.Context(t), bt)
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

func NewEthTx(fromAddress common.Address) txmgr.Tx {
	return txmgr.Tx{
		FromAddress:    fromAddress,
		ToAddress:      testutils.NewAddress(),
		EncodedPayload: []byte{1, 2, 3},
		Value:          big.Int(assets.NewEthValue(142)),
		FeeLimit:       uint64(1000000000),
		State:          txmgrcommon.TxUnstarted,
	}
}

func NewLegacyTransaction(nonce uint64, to common.Address, value *big.Int, gasLimit uint32, gasPrice *big.Int, data []byte) *types.Transaction {
	tx := types.LegacyTx{
		Nonce:    nonce,
		To:       &to,
		Value:    value,
		Gas:      uint64(gasLimit),
		GasPrice: gasPrice,
		Data:     data,
	}
	return types.NewTx(&tx)
}

func MustInsertUnconfirmedEthTx(t *testing.T, txStore txmgr.TestEvmTxStore, nonce int64, fromAddress common.Address, opts ...interface{}) txmgr.Tx {
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
	etx := NewEthTx(fromAddress)

	etx.BroadcastAt = &broadcastAt
	etx.InitialBroadcastAt = &broadcastAt
	n := evmtypes.Nonce(nonce)
	etx.Sequence = &n
	etx.State = txmgrcommon.TxUnconfirmed
	etx.ChainID = chainID
	require.NoError(t, txStore.InsertTx(testutils.Context(t), &etx))
	return etx
}

func MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t *testing.T, txStore txmgr.TestEvmTxStore, nonce int64, fromAddress common.Address, opts ...interface{}) txmgr.Tx {
	etx := MustInsertUnconfirmedEthTx(t, txStore, nonce, fromAddress, opts...)
	attempt := NewLegacyEthTxAttempt(t, etx.ID)
	ctx := testutils.Context(t)

	tx := NewLegacyTransaction(uint64(nonce), testutils.NewAddress(), big.NewInt(142), 242, big.NewInt(342), []byte{1, 2, 3})
	rlp := new(bytes.Buffer)
	require.NoError(t, tx.EncodeRLP(rlp))
	attempt.SignedRawTx = rlp.Bytes()

	attempt.State = txmgrtypes.TxAttemptBroadcast
	require.NoError(t, txStore.InsertTxAttempt(ctx, &attempt))
	etx, err := txStore.FindTxWithAttempts(ctx, etx.ID)
	require.NoError(t, err)
	return etx
}

func MustInsertConfirmedEthTxWithLegacyAttempt(t *testing.T, txStore txmgr.TestEvmTxStore, nonce int64, broadcastBeforeBlockNum int64, fromAddress common.Address) txmgr.Tx {
	timeNow := time.Now()
	etx := NewEthTx(fromAddress)
	ctx := testutils.Context(t)

	etx.BroadcastAt = &timeNow
	etx.InitialBroadcastAt = &timeNow
	n := evmtypes.Nonce(nonce)
	etx.Sequence = &n
	etx.State = txmgrcommon.TxConfirmed
	etx.MinConfirmations.SetValid(6)
	require.NoError(t, txStore.InsertTx(ctx, &etx))
	attempt := NewLegacyEthTxAttempt(t, etx.ID)
	attempt.BroadcastBeforeBlockNum = &broadcastBeforeBlockNum
	attempt.State = txmgrtypes.TxAttemptBroadcast
	require.NoError(t, txStore.InsertTxAttempt(ctx, &attempt))
	etx.TxAttempts = append(etx.TxAttempts, attempt)
	return etx
}

func NewLegacyEthTxAttempt(t *testing.T, etxID int64) txmgr.TxAttempt {
	gasPrice := assets.NewWeiI(1)
	return txmgr.TxAttempt{
		ChainSpecificFeeLimit: 42,
		TxID:                  etxID,
		TxFee:                 gas.EvmFee{Legacy: gasPrice},
		// Just a random signed raw tx that decodes correctly
		// Ignore all actual values
		SignedRawTx: hexutil.MustDecode("0xf889808504a817c8008307a12094000000000000000000000000000000000000000080a400000000000000000000000000000000000000000000000000000000000000000000000025a0838fe165906e2547b9a052c099df08ec891813fea4fcdb3c555362285eb399c5a070db99322490eb8a0f2270be6eca6e3aedbc49ff57ef939cf2774f12d08aa85e"),
		Hash:        evmutils.NewHash(),
		State:       txmgrtypes.TxAttemptInProgress,
	}
}

func NewDynamicFeeEthTxAttempt(t *testing.T, etxID int64) txmgr.TxAttempt {
	gasTipCap := assets.NewWeiI(1)
	gasFeeCap := assets.NewWeiI(1)
	return txmgr.TxAttempt{
		TxType: 0x2,
		TxID:   etxID,
		TxFee: gas.EvmFee{
			DynamicTipCap: gasTipCap,
			DynamicFeeCap: gasFeeCap,
		},
		// Just a random signed raw tx that decodes correctly
		// Ignore all actual values
		SignedRawTx:           hexutil.MustDecode("0xf889808504a817c8008307a12094000000000000000000000000000000000000000080a400000000000000000000000000000000000000000000000000000000000000000000000025a0838fe165906e2547b9a052c099df08ec891813fea4fcdb3c555362285eb399c5a070db99322490eb8a0f2270be6eca6e3aedbc49ff57ef939cf2774f12d08aa85e"),
		Hash:                  evmutils.NewHash(),
		State:                 txmgrtypes.TxAttemptInProgress,
		ChainSpecificFeeLimit: 42,
	}
}

type RandomKey struct {
	Nonce    int64
	Disabled bool

	chainIDs []ubig.Big // nil: Fixture, set empty for none
}

func (r RandomKey) MustInsert(t testing.TB, keystore keystore.Eth) (ethkey.KeyV2, common.Address) {
	ctx := testutils.Context(t)
	if r.chainIDs == nil {
		r.chainIDs = []ubig.Big{*ubig.New(&FixtureChainID)}
	}

	key := MustGenerateRandomKey(t)
	keystore.XXXTestingOnlyAdd(ctx, key)

	for _, cid := range r.chainIDs {
		require.NoError(t, keystore.Add(ctx, key.Address, cid.ToInt()))
		require.NoError(t, keystore.Enable(ctx, key.Address, cid.ToInt()))
		if r.Disabled {
			require.NoError(t, keystore.Disable(ctx, key.Address, cid.ToInt()))
		}
	}

	return key, key.Address
}

func (r RandomKey) MustInsertWithState(t testing.TB, keystore keystore.Eth) (ethkey.State, common.Address) {
	ctx := testutils.Context(t)
	k, address := r.MustInsert(t, keystore)
	state, err := keystore.GetStateForKey(ctx, k)
	require.NoError(t, err)
	return state, address
}

// MustInsertRandomKey inserts a randomly generated (not cryptographically secure) key for testing.
// By default, it is enabled for the fixture chain. Pass chainIDs to override.
// Use MustInsertRandomKeyNoChains for a key associate with no chains.
func MustInsertRandomKey(t testing.TB, keystore keystore.Eth, chainIDs ...ubig.Big) (ethkey.KeyV2, common.Address) {
	r := RandomKey{}
	if len(chainIDs) > 0 {
		r.chainIDs = chainIDs
	}
	return r.MustInsert(t, keystore)
}

func MustInsertRandomKeyNoChains(t testing.TB, keystore keystore.Eth) (ethkey.KeyV2, common.Address) {
	return RandomKey{chainIDs: []ubig.Big{}}.MustInsert(t, keystore)
}

func MustInsertRandomKeyReturningState(t testing.TB, keystore keystore.Eth) (ethkey.State, common.Address) {
	return RandomKey{}.MustInsertWithState(t, keystore)
}

func MustGenerateRandomKey(t testing.TB) ethkey.KeyV2 {
	key, err := ethkey.NewV2()
	require.NoError(t, err)
	return key
}

func MustGenerateRandomKeyState(_ testing.TB) ethkey.State {
	return ethkey.State{Address: NewEIP55Address()}
}

func MustInsertHead(t *testing.T, ds sqlutil.DataSource, number int64) evmtypes.Head {
	h := evmtypes.NewHead(big.NewInt(number), evmutils.NewHash(), evmutils.NewHash(), 0, ubig.New(&FixtureChainID))
	horm := headtracker.NewORM(FixtureChainID, ds)

	err := horm.IdempotentInsertHead(testutils.Context(t), &h)
	require.NoError(t, err)
	return h
}

func MustInsertV2JobSpec(t *testing.T, db *sqlx.DB, transmitterAddress common.Address) job.Job {
	t.Helper()

	addr, err := evmtypes.NewEIP55Address(transmitterAddress.Hex())
	require.NoError(t, err)

	pipelineSpec := pipeline.Spec{}
	err = db.Get(&pipelineSpec, `INSERT INTO pipeline_specs (dot_dag_source,created_at) VALUES ('',NOW()) RETURNING *`)
	require.NoError(t, err)

	oracleSpec := MustInsertOffchainreportingOracleSpec(t, db, addr)
	jb := job.Job{
		OCROracleSpec:   &oracleSpec,
		OCROracleSpecID: &oracleSpec.ID,
		ExternalJobID:   uuid.New(),
		Type:            job.OffchainReporting,
		SchemaVersion:   1,
		PipelineSpec:    &pipelineSpec,
		PipelineSpecID:  pipelineSpec.ID,
	}

	jorm := job.NewORM(db, nil, nil, nil, logger.TestLogger(t))
	err = jorm.InsertJob(testutils.Context(t), &jb)
	require.NoError(t, err)
	return jb
}

func MustInsertOffchainreportingOracleSpec(t *testing.T, db *sqlx.DB, transmitterAddress evmtypes.EIP55Address) job.OCROracleSpec {
	t.Helper()

	ocrKeyID := models.MustSha256HashFromHex(DefaultOCRKeyBundleID)
	spec := job.OCROracleSpec{}
	require.NoError(t, db.Get(&spec, `INSERT INTO ocr_oracle_specs (created_at, updated_at, contract_address, p2pv2_bootstrappers, is_bootstrap_peer, encrypted_ocr_key_bundle_id, transmitter_address, observation_timeout, blockchain_timeout, contract_config_tracker_subscribe_interval, contract_config_tracker_poll_interval, contract_config_confirmations, database_timeout, observation_grace_period, contract_transmitter_transmit_timeout, evm_chain_id) VALUES (
NOW(),NOW(),$1,'{}',false,$2,$3,0,0,0,0,0,0,0,0,0
) RETURNING *`, NewEIP55Address(), &ocrKeyID, &transmitterAddress))
	return spec
}

func MakeDirectRequestJobSpec(t *testing.T) *job.Job {
	t.Helper()
	drs := &job.DirectRequestSpec{EVMChainID: (*ubig.Big)(testutils.FixtureChainID)}
	spec := &job.Job{
		Type:              job.DirectRequest,
		SchemaVersion:     1,
		ExternalJobID:     uuid.New(),
		DirectRequestSpec: drs,
		Pipeline:          pipeline.Pipeline{},
		PipelineSpec:      &pipeline.Spec{},
	}
	return spec
}

func MustInsertKeeperJob(t *testing.T, db *sqlx.DB, korm *keeper.ORM, from evmtypes.EIP55Address, contract evmtypes.EIP55Address) job.Job {
	t.Helper()
	ctx := testutils.Context(t)

	var keeperSpec job.KeeperSpec
	err := korm.DataSource().GetContext(ctx, &keeperSpec, `INSERT INTO keeper_specs (contract_address, from_address, created_at, updated_at,evm_chain_id) VALUES ($1, $2, NOW(), NOW(), $3) RETURNING *`, contract, from, testutils.SimulatedChainID.Int64())
	require.NoError(t, err)

	var pipelineSpec pipeline.Spec
	err = korm.DataSource().GetContext(ctx, &pipelineSpec, `INSERT INTO pipeline_specs (dot_dag_source,created_at) VALUES ('',NOW()) RETURNING *`)
	require.NoError(t, err)

	jb := job.Job{
		KeeperSpec:     &keeperSpec,
		KeeperSpecID:   &keeperSpec.ID,
		ExternalJobID:  uuid.New(),
		Type:           job.Keeper,
		SchemaVersion:  1,
		PipelineSpec:   &pipelineSpec,
		PipelineSpecID: pipelineSpec.ID,
	}

	cfg := configtest.NewTestGeneralConfig(t)
	tlg := logger.TestLogger(t)
	prm := pipeline.NewORM(db, tlg, cfg.JobPipeline().MaxSuccessfulRuns())
	btORM := bridges.NewORM(db)
	jrm := job.NewORM(db, prm, btORM, nil, tlg)
	err = jrm.InsertJob(testutils.Context(t), &jb)
	require.NoError(t, err)
	jb.PipelineSpec.JobID = jb.ID
	return jb
}

func MustInsertKeeperRegistry(t *testing.T, db *sqlx.DB, korm *keeper.ORM, ethKeyStore keystore.Eth, keeperIndex, numKeepers, blockCountPerTurn int32) (keeper.Registry, job.Job) {
	t.Helper()
	ctx := testutils.Context(t)
	key, _ := MustInsertRandomKey(t, ethKeyStore, *ubig.New(testutils.SimulatedChainID))
	from := key.EIP55Address
	contractAddress := NewEIP55Address()
	jb := MustInsertKeeperJob(t, db, korm, from, contractAddress)
	registry := keeper.Registry{
		ContractAddress:   contractAddress,
		BlockCountPerTurn: blockCountPerTurn,
		CheckGas:          150_000,
		FromAddress:       from,
		JobID:             jb.ID,
		KeeperIndex:       keeperIndex,
		NumKeepers:        numKeepers,
		KeeperIndexMap: map[evmtypes.EIP55Address]int32{
			from: keeperIndex,
		},
	}
	err := korm.UpsertRegistry(ctx, &registry)
	require.NoError(t, err)
	return registry, jb
}

func MustInsertUpkeepForRegistry(t *testing.T, db *sqlx.DB, registry keeper.Registry) keeper.UpkeepRegistration {
	ctx := testutils.Context(t)
	korm := keeper.NewORM(db, logger.TestLogger(t))
	upkeepID := ubig.NewI(int64(mathrand.Uint32()))
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
	err = korm.UpsertUpkeep(ctx, &upkeep)
	require.NoError(t, err)
	return upkeep
}

func MustInsertPipelineRun(t *testing.T, db *sqlx.DB) (run pipeline.Run) {
	require.NoError(t, db.Get(&run, `INSERT INTO pipeline_runs (state,pipeline_spec_id,pruning_key,created_at) VALUES ($1, 0, 0, NOW()) RETURNING *`, pipeline.RunStatusRunning))
	return run
}

func MustInsertPipelineRunWithStatus(t *testing.T, db *sqlx.DB, pipelineSpecID int32, status pipeline.RunStatus, jobID int32) (run pipeline.Run) {
	var finishedAt *time.Time
	var outputs jsonserializable.JSONSerializable
	var allErrors pipeline.RunErrors
	var fatalErrors pipeline.RunErrors
	now := time.Now()
	switch status {
	case pipeline.RunStatusCompleted:
		finishedAt = &now
		outputs = jsonserializable.JSONSerializable{
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
	require.NoError(t, db.Get(&run, `INSERT INTO pipeline_runs (state,pipeline_spec_id,pruning_key,finished_at,outputs,all_errors,fatal_errors,created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, NOW()) RETURNING *`, status, pipelineSpecID, jobID, finishedAt, outputs, allErrors, fatalErrors))
	return run
}

func MustInsertPipelineSpec(t *testing.T, db *sqlx.DB) (spec pipeline.Spec) {
	err := db.Get(&spec, `INSERT INTO pipeline_specs (dot_dag_source,created_at) VALUES ('',NOW()) RETURNING *`)
	require.NoError(t, err)
	return
}

func MustInsertUnfinishedPipelineTaskRun(t *testing.T, db *sqlx.DB, pipelineRunID int64) (tr pipeline.TaskRun) {
	/* #nosec G404 */
	require.NoError(t, db.Get(&tr, `INSERT INTO pipeline_task_runs (dot_id, pipeline_run_id, id, type, created_at) VALUES ($1,$2,$3, '', NOW()) RETURNING *`, strconv.Itoa(mathrand.Int()), pipelineRunID, uuid.New()))
	return tr
}

func RandomLog(t *testing.T) types.Log {
	t.Helper()

	topics := make([]common.Hash, 4)
	for i := range topics {
		topics[i] = evmutils.NewHash()
	}

	return types.Log{
		Address:     testutils.NewAddress(),
		BlockHash:   evmutils.NewHash(),
		BlockNumber: uint64(mathrand.Intn(9999999)),
		Index:       uint(mathrand.Intn(9999999)),
		Data:        MustRandomBytes(t, 512),
		Topics:      []common.Hash{evmutils.NewHash(), evmutils.NewHash(), evmutils.NewHash(), evmutils.NewHash()},
	}
}

func RawNewRoundLog(t *testing.T, contractAddr common.Address, blockHash common.Hash, blockNumber uint64, logIndex uint, removed bool) types.Log {
	t.Helper()
	topic := (flux_aggregator_wrapper.FluxAggregatorNewRound{}).Topic()
	topics := []common.Hash{topic, evmutils.NewHash(), evmutils.NewHash()}
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
	ctx := testutils.Context(t)
	var prefix string
	if opts.NamePrefix != "" {
		prefix = opts.NamePrefix
	} else {
		prefix = "ei"
	}
	ei.Name = fmt.Sprintf("%s-%s", prefix, uuid.New())
	ei.URL = opts.URL
	ei.OutgoingSecret = opts.OutgoingSecret
	ei.OutgoingToken = opts.OutgoingToken
	token := auth.NewToken()
	ei.AccessKey = token.AccessKey
	ei.Salt = utils.NewSecret(utils.DefaultSecretSize)
	hashedSecret, err := auth.HashedSecret(token, ei.Salt)
	require.NoError(t, err)
	ei.HashedSecret = hashedSecret
	err = orm.CreateExternalInitiator(ctx, &ei)
	require.NoError(t, err)
	return ei
}
