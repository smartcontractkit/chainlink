package cltest

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/big"
	mathrand "math/rand"
	"net/url"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	ragep2ptypes "github.com/smartcontractkit/libocr/ragep2p/types"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink/v2/core/bridges"

	"github.com/smartcontractkit/chainlink-evm/pkg/heads"
	evmtypes "github.com/smartcontractkit/chainlink-evm/pkg/types"
	evmutils "github.com/smartcontractkit/chainlink-evm/pkg/utils"
	ubig "github.com/smartcontractkit/chainlink-evm/pkg/utils/big"

	"github.com/smartcontractkit/chainlink/v2/core/auth"
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
		btr.Name = bridges.MustParseBridgeName("test_bridge_" + rnd)
	}

	if opts.URL != "" {
		btr.URL = WebURL(t, opts.URL)
	} else {
		btr.URL = WebURL(t, "https://bridge.example.com/api?"+rnd)
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
func JSONFromString(t testing.TB, body string) models.JSON {
	return JSONFromBytes(t, []byte(body))
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

type RandomKey struct {
	Nonce    int64
	Disabled bool

	chainIDs []ubig.Big // nil: Fixture, set empty for none
}

func (r RandomKey) MustInsert(t testing.TB, keystore keystore.Eth) (ethkey.KeyV2, common.Address) {
	ctx := testutils.Context(t)
	chainIDs := r.chainIDs
	if chainIDs == nil {
		chainIDs = []ubig.Big{*ubig.New(&FixtureChainID)}
	}

	key := MustGenerateRandomKey(t)
	keystore.XXXTestingOnlyAdd(ctx, key)

	for _, cid := range chainIDs {
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

func MustGenerateRandomKey(t testing.TB) ethkey.KeyV2 {
	key, err := ethkey.NewV2()
	require.NoError(t, err)
	return key
}

func MustInsertHead(t *testing.T, ds sqlutil.DataSource, number int64) *evmtypes.Head {
	h := evmtypes.NewHead(big.NewInt(number), evmutils.NewHash(), evmutils.NewHash(), ubig.New(&FixtureChainID))
	horm := heads.NewORM(FixtureChainID, ds)

	err := horm.IdempotentInsertHead(testutils.Context(t), &h)
	require.NoError(t, err)
	return &h
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
