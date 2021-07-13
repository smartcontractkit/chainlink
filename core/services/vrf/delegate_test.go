package vrf

import (
	"bytes"
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/services/headtracker"
	httypes "github.com/smartcontractkit/chainlink/core/services/headtracker/types"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/store/models"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/solidity_vrf_coordinator_interface"

	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/stretchr/testify/assert"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	bptxmmocks "github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager/mocks"
	eth_mocks "github.com/smartcontractkit/chainlink/core/services/eth/mocks"
	"github.com/smartcontractkit/chainlink/core/services/log"
	log_mocks "github.com/smartcontractkit/chainlink/core/services/log/mocks"
	"github.com/smartcontractkit/chainlink/core/services/signatures/secp256k1"
	"github.com/smartcontractkit/chainlink/core/testdata/testspecs"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type vrfUniverse struct {
	jrm       job.ORM
	pr        pipeline.Runner
	prm       pipeline.ORM
	lb        *log_mocks.Broadcaster
	ec        *eth_mocks.Client
	ks        *keystore.Master
	vrfkey    secp256k1.PublicKey
	submitter common.Address
	txm       *bptxmmocks.TxManager
	hb        httypes.HeadBroadcaster
}

func buildVrfUni(t *testing.T, db *gorm.DB, cfg *orm.Config) vrfUniverse {
	// Mock all chain interactions
	lb := new(log_mocks.Broadcaster)
	ec := new(eth_mocks.Client)
	hb := headtracker.NewHeadBroadcaster()

	// Don't mock db interactions
	eb := postgres.NewEventBroadcaster(cfg.DatabaseURL(), 0, 0)
	err := eb.Start()
	require.NoError(t, err)
	t.Cleanup(func() { eb.Close() })
	prm := pipeline.NewORM(db)
	jrm := job.NewORM(db, cfg, prm, eb, &postgres.NullAdvisoryLocker{})
	pr := pipeline.NewRunner(prm, cfg, ec, nil)
	ks := keystore.New(db, utils.FastScryptParams)
	require.NoError(t, ks.Eth().Unlock("blah"))
	_, err = ks.Eth().CreateNewKey()
	require.NoError(t, err)
	submitter, err := ks.Eth().GetRoundRobinAddress()
	require.NoError(t, err)
	_, err = ks.VRF().Unlock("blah")
	require.NoError(t, err)
	vrfkey, err := ks.VRF().CreateKey()
	require.NoError(t, err)
	txm := new(bptxmmocks.TxManager)
	t.Cleanup(func() { txm.AssertExpectations(t) })

	return vrfUniverse{
		jrm:       jrm,
		pr:        pr,
		prm:       prm,
		lb:        lb,
		ec:        ec,
		ks:        ks,
		vrfkey:    vrfkey,
		submitter: submitter,
		txm:       txm,
		hb:        hb,
	}
}

func (v vrfUniverse) Assert(t *testing.T) {
	v.lb.AssertExpectations(t)
	v.ec.AssertExpectations(t)
}

func generateCallbackReturnValues(t *testing.T) []byte {
	callback, err := abi.NewType("tuple", "", []abi.ArgumentMarshaling{
		{Name: "callback_contract", Type: "address"},
		{Name: "randomness_fee", Type: "int256"},
		{Name: "seed_and_block_num", Type: "bytes32"}})
	require.NoError(t, err)
	var args abi.Arguments = []abi.Argument{{Type: callback}}
	b, err := args.Pack(solidity_vrf_coordinator_interface.Callbacks{
		RandomnessFee:   big.NewInt(10),
		SeedAndBlockNum: utils.NewHash(),
	})
	require.NoError(t, err)
	return b
}

func setup(t *testing.T) (vrfUniverse, *listener, job.Job) {
	db := pgtest.NewGormDB(t)
	c := orm.NewConfig()
	vuni := buildVrfUni(t, db, c)

	vd := NewDelegate(
		db,
		vuni.txm,
		vuni.ks,
		vuni.pr,
		vuni.prm,
		vuni.lb,
		vuni.hb,
		vuni.ec,
		c)
	vs := testspecs.GenerateVRFSpec(testspecs.VRFSpecParams{PublicKey: vuni.vrfkey.String()})
	t.Log(vs)
	jb, err := ValidatedVRFSpec(vs.Toml())
	require.NoError(t, err)
	require.NoError(t, vuni.jrm.CreateJob(context.Background(), &jb, pipeline.Pipeline{}))
	vl, err := vd.ServicesForSpec(jb)
	require.NoError(t, err)
	require.Len(t, vl, 1)
	listener := vl[0].(*listener)
	// Start the listener
	go func() {
		listener.runLogListener([]func(){}, 6)
	}()
	go func() {
		listener.runHeadListener(func() {})
	}()
	t.Cleanup(func() {
		listener.chStop <- struct{}{}
		select {
		case <-listener.waitOnStop:
		case <-time.After(1 * time.Second):
			t.Error("did not clean up properly")
		}
	})
	return vuni, listener, jb
}

func TestDelegate_ReorgAttackProtection(t *testing.T) {
	//vuni, listener, jb := setup(t)
	//TODO
}

func TestDelegate_ValidLog(t *testing.T) {
	vuni, listener, jb := setup(t)
	txHash := utils.NewHash()
	reqID := utils.NewHash()

	vuni.lb.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)
	done := make(chan struct{})
	vuni.lb.On("MarkConsumed", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		done <- struct{}{}
	}).Return(nil).Once()
	// Expect a call to check if the req is already fulfilled.
	vuni.ec.On("CallContract", mock.Anything, mock.Anything, mock.Anything).Return(generateCallbackReturnValues(t), nil)

	// Ensure we queue up a valid eth transaction
	// Linked to  requestID
	vuni.txm.On("CreateEthTransaction", mock.AnythingOfType("*gorm.DB"), vuni.submitter, common.HexToAddress(jb.VRFSpec.CoordinatorAddress.String()), mock.Anything, uint64(500000), mock.MatchedBy(func(meta *models.EthTxMetaV2) bool {
		return meta.JobID > 0 && meta.RequestID == reqID && meta.RequestTxHash == txHash
	}), bulletprooftxmanager.SendEveryStrategy{}).Once().Return(bulletprooftxmanager.EthTx{}, nil)

	// Send a valid log
	pk, err := secp256k1.NewPublicKeyFromHex(vuni.vrfkey.String())
	require.NoError(t, err)
	added := make(chan struct{})
	listener.reqAdded = func() {
		added <- struct{}{}
	}
	listener.HandleLog(log.NewLogBroadcast(types.Log{
		// Data has all the NON-indexed parameters
		Data: bytes.Join([][]byte{pk.MustHash().Bytes(), // key hash
			common.BigToHash(big.NewInt(42)).Bytes(), // seed
			utils.NewHash().Bytes(),                  // sender
			utils.NewHash().Bytes(),                  // fee
			reqID.Bytes()}, []byte{},                 // requestID
		),
		// JobID is indexed, thats why it lives in the Topics.
		Topics:      []common.Hash{{}, jb.ExternalIDEncodeStringToTopic()}, // jobID
		Address:     common.Address{},
		BlockNumber: 10,
		TxHash:      txHash,
		TxIndex:     0,
		BlockHash:   common.Hash{},
		Index:       0,
		Removed:     false,
	}))

	// Wait until the log is present
	select {
	case <-added:
	case <-time.After(1 * time.Second):
		t.FailNow()
	}

	// Feed it a head which confirms it.
	listener.OnNewLongestChain(context.Background(), models.Head{Number: 16})
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.FailNow()
	}

	// Ensure we created a successful run.
	runs, err := vuni.prm.GetAllRuns()
	require.NoError(t, err)
	require.Equal(t, 1, len(runs))
	assert.False(t, runs[0].Errors.HasError())
	m, ok := runs[0].Meta.Val.(map[string]interface{})
	require.True(t, ok)
	_, ok = m["eth_tx_id"]
	assert.True(t, ok)
	assert.Len(t, runs[0].PipelineTaskRuns, 0)

	vuni.Assert(t)
}

func TestDelegate_InvalidLog(t *testing.T) {
	vuni, listener, jb := setup(t)
	vuni.lb.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)
	done := make(chan struct{})
	vuni.lb.On("MarkConsumed", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		done <- struct{}{}
	}).Return(nil).Once()
	// Expect a call to check if the req is already fulfilled.
	vuni.ec.On("CallContract", mock.Anything, mock.Anything, mock.Anything).Return(generateCallbackReturnValues(t), nil)

	added := make(chan struct{})
	listener.reqAdded = func() {
		added <- struct{}{}
	}
	// Send a invalid log (keyhash doesnt match)
	listener.HandleLog(log.NewLogBroadcast(types.Log{
		// Data has all the NON-indexed parameters
		Data: append(append(append(append(
			utils.NewHash().Bytes(),                      // key hash
			common.BigToHash(big.NewInt(42)).Bytes()...), // seed
			utils.NewHash().Bytes()...), // sender
			utils.NewHash().Bytes()...), // fee
			utils.NewHash().Bytes()...), // requestID
		// JobID is indexed, thats why it lives in the Topics.
		Topics:      []common.Hash{{}, jb.ExternalIDEncodeStringToTopic()}, // jobID
		Address:     common.Address{},
		BlockNumber: 10,
		TxHash:      common.Hash{},
		TxIndex:     0,
		BlockHash:   common.Hash{},
		Index:       0,
		Removed:     false,
	}))
	// Wait until the log is present
	select {
	case <-added:
	case <-time.After(1 * time.Second):
		t.FailNow()
	}
	// Feed it a head which confirms it.
	listener.OnNewLongestChain(context.Background(), models.Head{Number: 16})
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.FailNow()
	}

	// Ensure we have not created a run.
	runs, err := vuni.prm.GetAllRuns()
	require.NoError(t, err)
	require.Equal(t, len(runs), 0)

	// Ensure we have NOT queued up an eth transaction
	var ethTxes []bulletprooftxmanager.EthTx
	err = vuni.prm.DB().Find(&ethTxes).Error
	require.NoError(t, err)
	require.Len(t, ethTxes, 0)
}
