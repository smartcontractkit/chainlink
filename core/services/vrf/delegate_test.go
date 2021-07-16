package vrf

import (
	"bytes"
	"context"
	"encoding/json"
	"math/big"
	"testing"
	"time"

	gormpostgres "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/smartcontractkit/chainlink/core/logger"
	"gopkg.in/guregu/null.v4"

	"github.com/theodesp/go-heaps/pairing"

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

func waitForChannel(t *testing.T, c chan struct{}, timeout time.Duration, errMsg string) {
	select {
	case <-c:
	case <-time.After(timeout):
		t.Error(errMsg)
	}
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
		waitForChannel(t, listener.waitOnStop, time.Second, "did not clean up properly")
	})
	return vuni, listener, jb
}

func TestStartingCounts(t *testing.T) {
	db := pgtest.NewGormDB(t)
	counts := getStartingResponseCounts(db, logger.Default)
	assert.Equal(t, 0, len(counts))

	ks := keystore.New(db, utils.FastScryptParams)
	ks.Eth().Unlock("blah")
	k, err := ks.Eth().CreateNewKey()
	require.NoError(t, err)
	b := time.Now()
	n1, n2, n3, n4 := int64(0), int64(1), int64(2), int64(3)
	m1 := models.EthTxMetaV2{
		RequestID: utils.PadByteToHash(0x10),
	}
	md1, err := json.Marshal(&m1)
	require.NoError(t, err)
	m2 := models.EthTxMetaV2{
		RequestID: utils.PadByteToHash(0x11),
	}
	md2, err := json.Marshal(&m2)
	var txes = []bulletprooftxmanager.EthTx{
		{
			Nonce:       &n1,
			FromAddress: k.Address.Address(),
			Error:       null.String{},
			BroadcastAt: &b,
			CreatedAt:   b,
			State:       bulletprooftxmanager.EthTxConfirmed,
			Meta:        gormpostgres.Jsonb{},
		},
		{
			Nonce:       &n2,
			FromAddress: k.Address.Address(),
			Error:       null.String{},
			BroadcastAt: &b,
			CreatedAt:   b,
			State:       bulletprooftxmanager.EthTxConfirmed,
			Meta:        gormpostgres.Jsonb{RawMessage: md1},
		},
		{
			Nonce:       &n3,
			FromAddress: k.Address.Address(),
			Error:       null.String{},
			BroadcastAt: &b,
			CreatedAt:   b,
			State:       bulletprooftxmanager.EthTxConfirmed,
			Meta:        gormpostgres.Jsonb{RawMessage: md2},
		},
		{
			Nonce:       &n4,
			FromAddress: k.Address.Address(),
			Error:       null.String{},
			BroadcastAt: &b,
			CreatedAt:   b,
			State:       bulletprooftxmanager.EthTxConfirmed,
			Meta:        gormpostgres.Jsonb{RawMessage: md2},
		},
	}
	require.NoError(t, db.Create(&txes).Error)
	counts = getStartingResponseCounts(db, logger.Default)
	assert.Equal(t, 2, len(counts))
	assert.Equal(t, uint64(1), counts[utils.PadByteToHash(0x10)])
	assert.Equal(t, uint64(2), counts[utils.PadByteToHash(0x11)])
}

func TestConfirmedLogExtraction(t *testing.T) {
	lsn := listener{}
	lsn.reqs = []request{
		{
			confirmedAtBlock: 2,
			req: &solidity_vrf_coordinator_interface.VRFCoordinatorRandomnessRequest{
				RequestID: utils.PadByteToHash(0x02),
			},
		},
		{
			confirmedAtBlock: 1,
			req: &solidity_vrf_coordinator_interface.VRFCoordinatorRandomnessRequest{
				RequestID: utils.PadByteToHash(0x01),
			},
		},
		{
			confirmedAtBlock: 3,
			req: &solidity_vrf_coordinator_interface.VRFCoordinatorRandomnessRequest{
				RequestID: utils.PadByteToHash(0x03),
			},
		},
	}
	// None are confirmed
	lsn.latestHead = 0
	logs := lsn.extractConfirmedLogs()
	assert.Equal(t, 0, len(logs))     // None ready
	assert.Equal(t, 3, len(lsn.reqs)) // All pending
	lsn.latestHead = 2
	logs = lsn.extractConfirmedLogs()
	assert.Equal(t, 2, len(logs))     // 1 and 2 should be confirmed
	assert.Equal(t, 1, len(lsn.reqs)) // 3 is still pending
	assert.Equal(t, uint64(3), lsn.reqs[0].confirmedAtBlock)
	// Another block way in the future should clear it
	lsn.latestHead = 10
	logs = lsn.extractConfirmedLogs()
	assert.Equal(t, 1, len(logs))     // remaining log
	assert.Equal(t, 0, len(lsn.reqs)) // all processed
}

func TestResponsePruning(t *testing.T) {
	lsn := listener{}
	lsn.latestHead = 10000
	lsn.respCount = map[[32]byte]uint64{
		utils.PadByteToHash(0x00): 1,
		utils.PadByteToHash(0x01): 1,
	}
	lsn.blockNumberToReqID = pairing.New()
	lsn.blockNumberToReqID.Insert(fulfilledReq{
		blockNumber: 1,
		reqID:       utils.PadByteToHash(0x00),
	})
	lsn.blockNumberToReqID.Insert(fulfilledReq{
		blockNumber: 2,
		reqID:       utils.PadByteToHash(0x01),
	})
	lsn.pruneConfirmedRequestCounts()
	assert.Equal(t, 2, len(lsn.respCount))
	lsn.latestHead = 10001
	lsn.pruneConfirmedRequestCounts()
	assert.Equal(t, 1, len(lsn.respCount))
	lsn.latestHead = 10002
	lsn.pruneConfirmedRequestCounts()
	assert.Equal(t, 0, len(lsn.respCount))
}

func TestDelegate_ReorgAttackProtection(t *testing.T) {
	vuni, listener, jb := setup(t)

	// Same request has already been fulfilled twice
	reqID := utils.NewHash()
	var reqIDBytes [32]byte
	copy(reqIDBytes[:], reqID.Bytes())
	listener.respCount[reqIDBytes] = 2

	// Send in the same request again
	pk, err := secp256k1.NewPublicKeyFromHex(vuni.vrfkey.String())
	require.NoError(t, err)
	added := make(chan struct{})
	listener.reqAdded = func() {
		added <- struct{}{}
	}
	preSeed := common.BigToHash(big.NewInt(42)).Bytes()
	txHash := utils.NewHash()
	vuni.lb.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)
	vuni.lb.On("MarkConsumed", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
	}).Return(nil).Once()
	vuni.ec.On("CallContract", mock.Anything, mock.Anything, mock.Anything).Return(generateCallbackReturnValues(t), nil)
	listener.HandleLog(log.NewLogBroadcast(types.Log{
		// Data has all the NON-indexed parameters
		Data: bytes.Join([][]byte{pk.MustHash().Bytes(), // key hash
			preSeed,                  // preSeed
			utils.NewHash().Bytes(),  // sender
			utils.NewHash().Bytes(),  // fee
			reqID.Bytes()}, []byte{}, // requestID
		),
		// JobID is indexed, thats why it lives in the Topics.
		Topics:      []common.Hash{{}, jb.ExternalIDEncodeStringToTopic()}, // jobID
		BlockNumber: 10,
		TxHash:      txHash,
	}, nil))

	// Wait until the log is present
	waitForChannel(t, added, time.Second, "request not added to the queue")
	assert.Equal(t, 1, len(listener.reqs))
	// It should be confirmed at 10+6*(2^2)
	assert.Equal(t, uint64(34), listener.reqs[0].confirmedAtBlock)
}

func TestDelegate_ValidLog(t *testing.T) {
	vuni, listener, jb := setup(t)
	txHash := utils.NewHash()
	reqID1 := utils.NewHash()
	reqID2 := utils.NewHash()
	pk, err := secp256k1.NewPublicKeyFromHex(vuni.vrfkey.String())
	require.NoError(t, err)
	added := make(chan struct{})
	listener.reqAdded = func() {
		added <- struct{}{}
	}
	preSeed := common.BigToHash(big.NewInt(42)).Bytes()
	bh := utils.NewHash()
	var tt = []struct {
		reqID [32]byte
		log   types.Log
	}{
		{
			reqID: reqID1,
			log: types.Log{
				// Data has all the NON-indexed parameters
				Data: bytes.Join([][]byte{
					pk.MustHash().Bytes(),                    // key hash
					common.BigToHash(big.NewInt(42)).Bytes(), // seed
					utils.NewHash().Bytes(),                  // sender
					utils.NewHash().Bytes(),                  // fee
					reqID1.Bytes()},                          // requestID
					[]byte{}),
				// JobID is indexed, thats why it lives in the Topics.
				Topics:      []common.Hash{{}, jb.ExternalIDEncodeStringToTopic()}, // jobID STRING
				TxHash:      txHash,
				BlockNumber: 10,
				BlockHash:   bh,
			},
		},
		{

			reqID: reqID2,
			log: types.Log{
				Data: bytes.Join([][]byte{
					pk.MustHash().Bytes(),                    // key hash
					common.BigToHash(big.NewInt(42)).Bytes(), // seed
					utils.NewHash().Bytes(),                  // sender
					utils.NewHash().Bytes(),                  // fee
					reqID2.Bytes()},                          // requestID
					[]byte{}),
				Topics:      []common.Hash{{}, jb.ExternalIDEncodeBytesToTopic()}, // jobID BYTES
				TxHash:      txHash,
				BlockNumber: 10,
				BlockHash:   bh,
			},
		},
	}

	consumed := make(chan struct{})
	for i, tc := range tt {
		tc := tc
		vuni.lb.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)
		vuni.lb.On("MarkConsumed", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			consumed <- struct{}{}
		}).Return(nil).Once()
		// Expect a call to check if the req is already fulfilled.
		vuni.ec.On("CallContract", mock.Anything, mock.Anything, mock.Anything).Return(generateCallbackReturnValues(t), nil)

		// Ensure we queue up a valid eth transaction
		// Linked to  requestID
		vuni.txm.On("CreateEthTransaction", mock.AnythingOfType("*gorm.DB"), vuni.submitter, common.HexToAddress(jb.VRFSpec.CoordinatorAddress.String()), mock.Anything, uint64(500000), mock.MatchedBy(func(meta *models.EthTxMetaV2) bool {
			return meta.JobID > 0 && meta.RequestID == tc.reqID && meta.RequestTxHash == txHash
		}), bulletprooftxmanager.SendEveryStrategy{}).Once().Return(bulletprooftxmanager.EthTx{}, nil)

		listener.HandleLog(log.NewLogBroadcast(tc.log, nil))
		// Wait until the log is present
		waitForChannel(t, added, time.Second, "request not added to the queue")
		// Feed it a head which confirms it.
		listener.OnNewLongestChain(context.Background(), models.Head{Number: 16})
		waitForChannel(t, consumed, 2*time.Second, "did not mark consumed")

		// Ensure we created a successful run.
		runs, err := vuni.prm.GetAllRuns()
		require.NoError(t, err)
		require.Equal(t, i+1, len(runs))
		assert.False(t, runs[0].Errors.HasError())
		m, ok := runs[0].Meta.Val.(map[string]interface{})
		require.True(t, ok)
		_, ok = m["eth_tx_id"]
		assert.True(t, ok)
		assert.Len(t, runs[0].PipelineTaskRuns, 0)

		p, err := vuni.ks.VRF().GenerateProof(pk, utils.MustHash(string(bytes.Join([][]byte{preSeed, bh.Bytes()}, []byte{}))).Big())
		require.NoError(t, err)
		vuni.lb.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)
		vuni.lb.On("MarkConsumed", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			consumed <- struct{}{}
		}).Return(nil).Once()
		// If we send a completed log we should the respCount increase
		var reqIDBytes []byte
		copy(reqIDBytes[:], tc.reqID[:])
		listener.HandleLog(log.NewLogBroadcast(types.Log{
			// Data has all the NON-indexed parameters
			Data: bytes.Join([][]byte{reqIDBytes, // output
				p.Output.Bytes(),
			}, []byte{},
			),
			BlockNumber: 10,
			TxHash:      txHash,
		}, &solidity_vrf_coordinator_interface.VRFCoordinatorRandomnessRequestFulfilled{RequestId: tc.reqID}))
		waitForChannel(t, consumed, 2*time.Second, "fulfillment log not marked consumed")
		// Should record that we've responded to this request
		assert.Equal(t, uint64(1), listener.respCount[tc.reqID])
		vuni.Assert(t)
	}
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
	}, nil))
	waitForChannel(t, added, time.Second, "request not queued")
	// Feed it a head which confirms it.
	listener.OnNewLongestChain(context.Background(), models.Head{Number: 16})
	waitForChannel(t, done, time.Second, "log not consumed")

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
