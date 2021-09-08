package vrf

import (
	"bytes"
	"context"
	"encoding/json"
	"math/big"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/chains/evm"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/solidity_vrf_coordinator_interface"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	bptxmmocks "github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager/mocks"
	eth_mocks "github.com/smartcontractkit/chainlink/core/services/eth/mocks"
	"github.com/smartcontractkit/chainlink/core/services/headtracker"
	httypes "github.com/smartcontractkit/chainlink/core/services/headtracker/types"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/vrfkey"
	"github.com/smartcontractkit/chainlink/core/services/log"
	log_mocks "github.com/smartcontractkit/chainlink/core/services/log/mocks"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/services/signatures/secp256k1"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/testdata/testspecs"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/theodesp/go-heaps/pairing"
	"gopkg.in/guregu/null.v4"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type vrfUniverse struct {
	jrm       job.ORM
	pr        pipeline.Runner
	prm       pipeline.ORM
	lb        *log_mocks.Broadcaster
	ec        *eth_mocks.Client
	ks        keystore.Master
	vrfkey    vrfkey.KeyV2
	submitter common.Address
	txm       *bptxmmocks.TxManager
	hb        httypes.HeadBroadcaster
	cc        evm.ChainSet
	cid       big.Int
}

func buildVrfUni(t *testing.T, db *gorm.DB, cfg *configtest.TestGeneralConfig) vrfUniverse {
	// Mock all chain interactions
	lb := new(log_mocks.Broadcaster)
	lb.Test(t)
	lb.On("AddDependents", 1).Maybe()
	ec := new(eth_mocks.Client)
	ec.Test(t)
	ec.On("ChainID").Return(big.NewInt(0))
	hb := headtracker.NewHeadBroadcaster(logger.Default)

	// Don't mock db interactions
	eb := postgres.NewEventBroadcaster(cfg.DatabaseURL(), 0, 0)
	err := eb.Start()
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, eb.Close()) })
	prm := pipeline.NewORM(db)
	txm := new(bptxmmocks.TxManager)
	ks := keystore.New(db, utils.FastScryptParams)
	cc := evmtest.NewChainSet(t, evmtest.TestChainOpts{LogBroadcaster: lb, KeyStore: ks.Eth(), Client: ec, DB: db, GeneralConfig: cfg, TxManager: txm})
	jrm := job.NewORM(db, cc, prm, eb, &postgres.NullAdvisoryLocker{}, ks)
	pr := pipeline.NewRunner(prm, cfg, cc, ks.Eth(), ks.VRF())
	require.NoError(t, ks.Unlock("p4SsW0rD1!@#_"))
	_, err = ks.Eth().Create(big.NewInt(0))
	require.NoError(t, err)
	submitter, err := ks.Eth().GetRoundRobinAddress()
	require.NoError(t, err)
	vrfkey, err := ks.VRF().Create()
	require.NoError(t, err)

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
		cc:        cc,
		cid:       *ec.ChainID(),
	}
}

func (v vrfUniverse) Assert(t *testing.T) {
	v.lb.AssertExpectations(t)
	v.ec.AssertExpectations(t)
}

func generateCallbackReturnValues(t *testing.T, fulfilled bool) []byte {
	callback, err := abi.NewType("tuple", "", []abi.ArgumentMarshaling{
		{Name: "callback_contract", Type: "address"},
		{Name: "randomness_fee", Type: "int256"},
		{Name: "seed_and_block_num", Type: "bytes32"}})
	require.NoError(t, err)
	var args abi.Arguments = []abi.Argument{{Type: callback}}
	if fulfilled {
		// Empty callback
		b, err := args.Pack(solidity_vrf_coordinator_interface.Callbacks{
			RandomnessFee:   big.NewInt(10),
			SeedAndBlockNum: utils.EmptyHash,
		})
		require.NoError(t, err)
		return b
	}
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

func setup(t *testing.T) (vrfUniverse, *listenerV1, job.Job, func()) {
	db := pgtest.NewGormDB(t)
	c := configtest.NewTestGeneralConfig(t)
	vuni := buildVrfUni(t, db, c)

	vd := NewDelegate(
		db,
		vuni.ks,
		vuni.pr,
		vuni.prm,
		vuni.cc)
	vs := testspecs.GenerateVRFSpec(testspecs.VRFSpecParams{PublicKey: vuni.vrfkey.PublicKey.String()})
	jb, err := ValidatedVRFSpec(vs.Toml())
	require.NoError(t, err)
	jb, err = vuni.jrm.CreateJob(context.Background(), &jb, jb.Pipeline)
	require.NoError(t, err)
	vl, err := vd.ServicesForSpec(jb)
	require.NoError(t, err)
	require.Len(t, vl, 1)
	listener := vl[0].(*listenerV1)
	// Start the listenerV1
	go func() {
		listener.runLogListener([]func(){}, 6)
	}()
	go func() {
		listener.runHeadListener(func() {})
	}()
	return vuni, listener, jb, func() {
		listener.chStop <- struct{}{}
		waitForChannel(t, listener.waitOnStop, time.Second, "did not clean up properly")
		vuni.txm.AssertExpectations(t)
	}
}

func TestStartingCounts(t *testing.T) {
	db := pgtest.NewGormDB(t)
	counts := getStartingResponseCounts(db, logger.Default)
	assert.Equal(t, 0, len(counts))
	ks := keystore.New(db, utils.FastScryptParams)
	err := ks.Unlock("p4SsW0rD1!@#_")
	require.NoError(t, err)
	k, err := ks.Eth().Create(big.NewInt(0))
	require.NoError(t, err)
	b := time.Now()
	n1, n2, n3, n4 := int64(0), int64(1), int64(2), int64(3)
	m1 := bulletprooftxmanager.EthTxMeta{
		RequestID: utils.PadByteToHash(0x10),
	}
	md1, err := json.Marshal(&m1)
	require.NoError(t, err)
	m2 := bulletprooftxmanager.EthTxMeta{
		RequestID: utils.PadByteToHash(0x11),
	}
	md2, err := json.Marshal(&m2)
	require.NoError(t, err)
	var txes = []bulletprooftxmanager.EthTx{
		{
			Nonce:          &n1,
			FromAddress:    k.Address.Address(),
			Error:          null.String{},
			BroadcastAt:    &b,
			CreatedAt:      b,
			State:          bulletprooftxmanager.EthTxConfirmed,
			Meta:           datatypes.JSON{},
			EncodedPayload: []byte{},
		},
		{
			Nonce:          &n2,
			FromAddress:    k.Address.Address(),
			Error:          null.String{},
			BroadcastAt:    &b,
			CreatedAt:      b,
			State:          bulletprooftxmanager.EthTxConfirmed,
			Meta:           datatypes.JSON(md1),
			EncodedPayload: []byte{},
		},
		{
			Nonce:          &n3,
			FromAddress:    k.Address.Address(),
			Error:          null.String{},
			BroadcastAt:    &b,
			CreatedAt:      b,
			State:          bulletprooftxmanager.EthTxConfirmed,
			Meta:           datatypes.JSON(md2),
			EncodedPayload: []byte{},
		},
		{
			Nonce:          &n4,
			FromAddress:    k.Address.Address(),
			Error:          null.String{},
			BroadcastAt:    &b,
			CreatedAt:      b,
			State:          bulletprooftxmanager.EthTxConfirmed,
			Meta:           datatypes.JSON(md2),
			EncodedPayload: []byte{},
		},
	}
	require.NoError(t, db.Create(&txes).Error)
	counts = getStartingResponseCounts(db, logger.Default)
	assert.Equal(t, 2, len(counts))
	assert.Equal(t, uint64(1), counts[utils.PadByteToHash(0x10)])
	assert.Equal(t, uint64(2), counts[utils.PadByteToHash(0x11)])
}

func TestConfirmedLogExtraction(t *testing.T) {
	lsn := listenerV1{}
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
	lsn := listenerV1{}
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
	vuni, listener, jb, assertMocksCalled := setup(t)
	defer assertMocksCalled()

	// Same request has already been fulfilled twice
	reqID := utils.NewHash()
	var reqIDBytes [32]byte
	copy(reqIDBytes[:], reqID.Bytes())
	listener.respCount[reqIDBytes] = 2

	// Send in the same request again
	pk, err := secp256k1.NewPublicKeyFromHex(vuni.vrfkey.PublicKey.String())
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
	vuni.ec.On("CallContract", mock.Anything, mock.Anything, mock.Anything).Return(generateCallbackReturnValues(t, false), nil)
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
	}, vuni.cid, nil))

	// Wait until the log is present
	waitForChannel(t, added, time.Second, "request not added to the queue")
	assert.Equal(t, 1, len(listener.reqs))
	// It should be confirmed at 10+6*(2^2)
	assert.Equal(t, uint64(34), listener.reqs[0].confirmedAtBlock)
}

func TestDelegate_ValidLog(t *testing.T) {
	vuni, listener, jb, assertMocksCalled := setup(t)
	defer assertMocksCalled()
	txHash := utils.NewHash()
	reqID1 := utils.NewHash()
	reqID2 := utils.NewHash()
	keyID := vuni.vrfkey.PublicKey.String()
	pk, err := secp256k1.NewPublicKeyFromHex(keyID)
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

	runComplete := make(chan struct{})
	vuni.pr.OnRunFinished(func(run *pipeline.Run) {
		if run.State == pipeline.RunStatusCompleted {
			runComplete <- struct{}{}
		}
	})

	consumed := make(chan struct{})
	for i, tc := range tt {
		tc := tc
		vuni.lb.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)
		vuni.lb.On("MarkConsumed", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			consumed <- struct{}{}
		}).Return(nil).Once()
		// Expect a call to check if the req is already fulfilled.
		vuni.ec.On("CallContract", mock.Anything, mock.Anything, mock.Anything).Return(generateCallbackReturnValues(t, false), nil)

		// Ensure we queue up a valid eth transaction
		// Linked to  requestID
		vuni.txm.On("CreateEthTransaction", mock.AnythingOfType("*gorm.DB"),
			mock.MatchedBy(func(newTx bulletprooftxmanager.NewTx) bool {
				meta := newTx.Meta
				return newTx.FromAddress == vuni.submitter &&
					newTx.ToAddress == common.HexToAddress(jb.VRFSpec.CoordinatorAddress.String()) &&
					newTx.GasLimit == uint64(500000) &&
					(meta.JobID > 0 && meta.RequestID == tc.reqID && meta.RequestTxHash == txHash)
			}),
		).Once().Return(bulletprooftxmanager.EthTx{}, nil)

		listener.HandleLog(log.NewLogBroadcast(tc.log, vuni.cid, nil))
		// Wait until the log is present
		waitForChannel(t, added, time.Second, "request not added to the queue")
		// Feed it a head which confirms it.
		listener.OnNewLongestChain(context.Background(), models.Head{Number: 16})
		waitForChannel(t, consumed, 2*time.Second, "did not mark consumed")

		// Ensure we created a successful run.
		waitForChannel(t, runComplete, 2*time.Second, "pipeline not complete")
		runs, err := vuni.prm.GetAllRuns()
		require.NoError(t, err)
		require.Equal(t, i+1, len(runs))
		assert.False(t, runs[0].Errors.HasError())
		// Should have 4 tasks all completed
		assert.Len(t, runs[0].PipelineTaskRuns, 4)

		p, err := vuni.ks.VRF().GenerateProof(keyID, utils.MustHash(string(bytes.Join([][]byte{preSeed, bh.Bytes()}, []byte{}))).Big())
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
		}, vuni.cid, &solidity_vrf_coordinator_interface.VRFCoordinatorRandomnessRequestFulfilled{RequestId: tc.reqID}))
		waitForChannel(t, consumed, 2*time.Second, "fulfillment log not marked consumed")
		// Should record that we've responded to this request
		assert.Equal(t, uint64(1), listener.respCount[tc.reqID])
		vuni.Assert(t)
	}
}

func TestDelegate_InvalidLog(t *testing.T) {
	vuni, listener, jb, assertMocksCalled := setup(t)
	defer assertMocksCalled()
	vuni.lb.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)
	done := make(chan struct{})
	vuni.lb.On("MarkConsumed", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		done <- struct{}{}
	}).Return(nil).Once()
	// Expect a call to check if the req is already fulfilled.
	vuni.ec.On("CallContract", mock.Anything, mock.Anything, mock.Anything).Return(generateCallbackReturnValues(t, false), nil)

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
	}, vuni.cid, nil))
	waitForChannel(t, added, time.Second, "request not queued")
	// Feed it a head which confirms it.
	listener.OnNewLongestChain(context.Background(), models.Head{Number: 16})
	waitForChannel(t, done, time.Second, "log not consumed")

	// Should create a run that errors in the vrf task
	runs, err := vuni.prm.GetAllRuns()
	require.NoError(t, err)
	require.Equal(t, len(runs), 1)
	for _, tr := range runs[0].PipelineTaskRuns {
		if tr.Type == pipeline.TaskTypeVRF {
			assert.Contains(t, tr.Error.String, "invalid key hash")
		}
		// Log parsing task itself should succeed.
		if tr.Type != pipeline.TaskTypeETHABIDecodeLog {
			assert.Nil(t, tr.Output)
		}
	}

	// Ensure we have NOT queued up an eth transaction
	var ethTxes []bulletprooftxmanager.EthTx
	err = vuni.prm.DB().Find(&ethTxes).Error
	require.NoError(t, err)
	require.Len(t, ethTxes, 0)
}

func TestFulfilledCheck(t *testing.T) {
	vuni, listener, jb, assertMocksCalled := setup(t)
	defer assertMocksCalled()
	vuni.lb.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)
	done := make(chan struct{})
	vuni.lb.On("MarkConsumed", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		done <- struct{}{}
	}).Return(nil).Once()
	// Expect a call to check if the req is already fulfilled.
	// We return already fulfilled
	vuni.ec.On("CallContract", mock.Anything, mock.Anything, mock.Anything).Return(generateCallbackReturnValues(t, true), nil)

	added := make(chan struct{})
	listener.reqAdded = func() {
		added <- struct{}{}
	}
	// Send a invalid log (keyhash doesnt match)
	listener.HandleLog(log.NewLogBroadcast(
		types.Log{
			// Data has all the NON-indexed parameters
			Data: bytes.Join([][]byte{
				vuni.vrfkey.PublicKey.MustHash().Bytes(), // key hash
				common.BigToHash(big.NewInt(42)).Bytes(), // seed
				utils.NewHash().Bytes(),                  // sender
				utils.NewHash().Bytes(),                  // fee
				utils.NewHash().Bytes()},                 // requestID
				[]byte{}),
			// JobID is indexed, thats why it lives in the Topics.
			Topics: []common.Hash{{}, jb.ExternalIDEncodeStringToTopic()}, // jobID STRING
			//TxHash:      utils.NewHash().Bytes(),
			BlockNumber: 10,
			//BlockHash:   utils.NewHash().Bytes(),
		}, vuni.cid, nil))

	// Should queue the request, even though its already fulfilled
	waitForChannel(t, added, time.Second, "request not queued")
	listener.OnNewLongestChain(context.Background(), models.Head{Number: 16})
	waitForChannel(t, done, time.Second, "log not consumed")

	// Should consume the log with no run
	runs, err := vuni.prm.GetAllRuns()
	require.NoError(t, err)
	require.Equal(t, len(runs), 0)
}
