package vrf_test

import (
	"bytes"
	"math/big"
	"testing"
	"time"

	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink/v2/core/assets"
	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm"
	evmclimocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker"
	httypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/log"
	log_mocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/log/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	txmmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr/mocks"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/solidity_vrf_coordinator_interface"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	configtest "github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/vrfkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/services/signatures/secp256k1"
	"github.com/smartcontractkit/chainlink/v2/core/services/srvctest"
	"github.com/smartcontractkit/chainlink/v2/core/services/vrf"
	vrf_mocks "github.com/smartcontractkit/chainlink/v2/core/services/vrf/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/testdata/testspecs"
	"github.com/smartcontractkit/chainlink/v2/core/utils"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type vrfUniverse struct {
	jrm       job.ORM
	pr        pipeline.Runner
	prm       pipeline.ORM
	lb        *log_mocks.Broadcaster
	ec        *evmclimocks.Client
	ks        keystore.Master
	vrfkey    vrfkey.KeyV2
	submitter common.Address
	txm       *txmmocks.MockEvmTxManager
	hb        httypes.HeadBroadcaster
	cc        evm.ChainSet
	cid       big.Int
}

func buildVrfUni(t *testing.T, db *sqlx.DB, cfg chainlink.GeneralConfig) vrfUniverse {
	// Mock all chain interactions
	lb := log_mocks.NewBroadcaster(t)
	lb.On("AddDependents", 1).Maybe()
	ec := evmclimocks.NewClient(t)
	ec.On("ConfiguredChainID").Return(testutils.FixtureChainID)
	lggr := logger.TestLogger(t)
	hb := headtracker.NewHeadBroadcaster(lggr)

	// Don't mock db interactions
	prm := pipeline.NewORM(db, lggr, cfg)
	btORM := bridges.NewORM(db, lggr, cfg)
	txm := txmmocks.NewMockEvmTxManager(t)
	ks := keystore.New(db, utils.FastScryptParams, lggr, cfg)
	cc := evmtest.NewChainSet(t, evmtest.TestChainOpts{LogBroadcaster: lb, KeyStore: ks.Eth(), Client: ec, DB: db, GeneralConfig: cfg, TxManager: txm})
	jrm := job.NewORM(db, cc, prm, btORM, ks, lggr, cfg)
	t.Cleanup(func() { jrm.Close() })
	pr := pipeline.NewRunner(prm, btORM, cfg, cc, ks.Eth(), ks.VRF(), lggr, nil, nil)
	require.NoError(t, ks.Unlock(testutils.Password))
	k, err := ks.Eth().Create(testutils.FixtureChainID)
	require.NoError(t, err)
	submitter := k.Address
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
		cid:       *ec.ConfiguredChainID(),
	}
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

func setup(t *testing.T) (vrfUniverse, *vrf.ListenerV1, job.Job) {
	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewTestGeneralConfig(t)
	vuni := buildVrfUni(t, db, cfg)

	mailMon := srvctest.Start(t, utils.NewMailboxMonitor(t.Name()))

	vd := vrf.NewDelegate(
		db,
		vuni.ks,
		vuni.pr,
		vuni.prm,
		vuni.cc,
		logger.TestLogger(t),
		cfg,
		mailMon)
	vs := testspecs.GenerateVRFSpec(testspecs.VRFSpecParams{PublicKey: vuni.vrfkey.PublicKey.String()})
	jb, err := vrf.ValidatedVRFSpec(vs.Toml())
	require.NoError(t, err)
	err = vuni.jrm.CreateJob(&jb)
	require.NoError(t, err)
	vl, err := vd.ServicesForSpec(jb)
	require.NoError(t, err)
	require.Len(t, vl, 1)
	listener := vl[0].(*vrf.ListenerV1)
	// Start the listenerV1
	go func() {
		listener.RunLogListener([]func(){}, 6)
	}()
	go func() {
		listener.RunHeadListener(func() {})
	}()
	t.Cleanup(func() { listener.Stop(t) })
	return vuni, listener, jb
}

func TestDelegate_ReorgAttackProtection(t *testing.T) {
	vuni, listener, jb := setup(t)

	// Same request has already been fulfilled twice
	reqID := utils.NewHash()
	var reqIDBytes [32]byte
	copy(reqIDBytes[:], reqID.Bytes())
	listener.SetRespCount(reqIDBytes, 2)

	// Send in the same request again
	pk, err := secp256k1.NewPublicKeyFromHex(vuni.vrfkey.PublicKey.String())
	require.NoError(t, err)
	added := make(chan struct{})
	listener.SetReqAdded(func() {
		added <- struct{}{}
	})
	preSeed := common.BigToHash(big.NewInt(42)).Bytes()
	txHash := utils.NewHash()
	vuni.lb.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil).Maybe()
	vuni.lb.On("MarkConsumed", mock.Anything, mock.Anything).Return(nil).Maybe()
	vuni.ec.On("CallContract", mock.Anything, mock.Anything, mock.Anything).Return(generateCallbackReturnValues(t, false), nil).Maybe()
	listener.HandleLog(log.NewLogBroadcast(types.Log{
		// Data has all the NON-indexed parameters
		Data: bytes.Join([][]byte{pk.MustHash().Bytes(), // key hash
			preSeed,                  // preSeed
			utils.NewHash().Bytes(),  // sender
			utils.NewHash().Bytes(),  // fee
			reqID.Bytes()}, []byte{}, // requestID
		),
		// JobID is indexed, thats why it lives in the Topics.
		Topics: []common.Hash{
			vrf.VRFRandomnessRequestLogTopic(),
			jb.ExternalIDEncodeStringToTopic(), // jobID
		},
		BlockNumber: 10,
		TxHash:      txHash,
	}, vuni.cid, nil))

	// Wait until the log is present
	waitForChannel(t, added, time.Second, "request not added to the queue")
	reqs := listener.ReqsConfirmedAt()
	if assert.Equal(t, 1, len(reqs)) {
		// It should be confirmed at 10+6*(2^2)
		assert.Equal(t, uint64(34), reqs[0])
	}
}

func TestDelegate_ValidLog(t *testing.T) {
	vuni, listener, jb := setup(t)
	txHash := utils.NewHash()
	reqID1 := utils.NewHash()
	reqID2 := utils.NewHash()
	keyID := vuni.vrfkey.PublicKey.String()
	pk, err := secp256k1.NewPublicKeyFromHex(keyID)
	require.NoError(t, err)
	added := make(chan struct{})
	listener.SetReqAdded(func() {
		added <- struct{}{}
	})
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
				Topics: []common.Hash{
					vrf.VRFRandomnessRequestLogTopic(),
					jb.ExternalIDEncodeStringToTopic(), // jobID STRING
				},
				TxHash:      txHash,
				BlockNumber: 10,
				BlockHash:   bh,
				Index:       1,
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
				Topics: []common.Hash{
					vrf.VRFRandomnessRequestLogTopic(),
					jb.ExternalIDEncodeBytesToTopic(), // jobID BYTES
				},
				TxHash:      txHash,
				BlockNumber: 10,
				BlockHash:   bh,
				Index:       2,
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
		// Linked to requestID
		vuni.txm.On("CreateEthTransaction",
			mock.MatchedBy(func(newTx txmgr.EvmNewTx) bool {
				meta := newTx.Meta
				return newTx.FromAddress == vuni.submitter &&
					newTx.ToAddress == common.HexToAddress(jb.VRFSpec.CoordinatorAddress.String()) &&
					newTx.FeeLimit == uint32(500000) &&
					meta.JobID != nil && meta.RequestID != nil && meta.RequestTxHash != nil &&
					(*meta.JobID > 0 && *meta.RequestID == tc.reqID && *meta.RequestTxHash == txHash)
			}),
		).Once().Return(txmgr.EvmTx{}, nil)

		listener.HandleLog(log.NewLogBroadcast(tc.log, vuni.cid, nil))
		// Wait until the log is present
		waitForChannel(t, added, time.Second, "request not added to the queue")
		// Feed it a head which confirms it.
		listener.OnNewLongestChain(testutils.Context(t), &evmtypes.Head{Number: 16})
		waitForChannel(t, consumed, 2*time.Second, "did not mark consumed")

		// Ensure we created a successful run.
		waitForChannel(t, runComplete, 2*time.Second, "pipeline not complete")
		runs, err := vuni.prm.GetAllRuns()
		require.NoError(t, err)
		require.Equal(t, i+1, len(runs))
		assert.False(t, runs[0].FatalErrors.HasError())
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
			Index:       uint(i),
		}, vuni.cid, &solidity_vrf_coordinator_interface.VRFCoordinatorRandomnessRequestFulfilled{RequestId: tc.reqID}))
		waitForChannel(t, consumed, 2*time.Second, "fulfillment log not marked consumed")
		// Should record that we've responded to this request
		assert.Equal(t, uint64(1), listener.RespCount(tc.reqID))
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
	vuni.ec.On("CallContract", mock.Anything, mock.Anything, mock.Anything).Return(generateCallbackReturnValues(t, false), nil)

	added := make(chan struct{})
	listener.SetReqAdded(func() {
		added <- struct{}{}
	})
	// Send an invalid log (keyhash doesnt match)
	listener.HandleLog(log.NewLogBroadcast(types.Log{
		// Data has all the NON-indexed parameters
		Data: append(append(append(append(
			utils.NewHash().Bytes(),                      // key hash
			common.BigToHash(big.NewInt(42)).Bytes()...), // seed
			utils.NewHash().Bytes()...), // sender
			utils.NewHash().Bytes()...), // fee
			utils.NewHash().Bytes()...), // requestID
		// JobID is indexed, that's why it lives in the Topics.
		Topics: []common.Hash{
			vrf.VRFRandomnessRequestLogTopic(),
			jb.ExternalIDEncodeBytesToTopic(), // jobID
		},
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
	listener.OnNewLongestChain(testutils.Context(t), &evmtypes.Head{Number: 16})
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
			assert.False(t, tr.Output.Valid)
		}
	}

	// Ensure we have NOT queued up an eth transaction
	var ethTxes []txmgr.DbEthTx
	err = vuni.prm.GetQ().Select(&ethTxes, `SELECT * FROM eth_txes;`)
	require.NoError(t, err)
	require.Len(t, ethTxes, 0)
}

func TestFulfilledCheck(t *testing.T) {
	vuni, listener, jb := setup(t)
	vuni.lb.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)
	done := make(chan struct{})
	vuni.lb.On("MarkConsumed", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		done <- struct{}{}
	}).Return(nil).Once()
	// Expect a call to check if the req is already fulfilled.
	// We return already fulfilled
	vuni.ec.On("CallContract", mock.Anything, mock.Anything, mock.Anything).Return(generateCallbackReturnValues(t, true), nil)

	added := make(chan struct{})
	listener.SetReqAdded(func() {
		added <- struct{}{}
	})
	// Send an invalid log (keyhash doesn't match)
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
			// JobID is indexed, that's why it lives in the Topics.
			Topics: []common.Hash{
				vrf.VRFRandomnessRequestLogTopic(),
				jb.ExternalIDEncodeBytesToTopic(), // jobID STRING
			},
			//TxHash:      utils.NewHash().Bytes(),
			BlockNumber: 10,
			//BlockHash:   utils.NewHash().Bytes(),
		}, vuni.cid, nil))

	// Should queue the request, even though its already fulfilled
	waitForChannel(t, added, time.Second, "request not queued")
	listener.OnNewLongestChain(testutils.Context(t), &evmtypes.Head{Number: 16})
	waitForChannel(t, done, time.Second, "log not consumed")

	// Should consume the log with no run
	runs, err := vuni.prm.GetAllRuns()
	require.NoError(t, err)
	require.Equal(t, len(runs), 0)
}

func Test_CheckFromAddressMaxGasPrices(t *testing.T) {
	t.Run("returns nil error if gasLanePrice not set in job spec", func(tt *testing.T) {
		spec := `
type            = "vrf"
schemaVersion   = 1
minIncomingConfirmations = 10
publicKey = "0x79BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F8179800"
coordinatorAddress = "0xB3b7874F13387D44a3398D298B075B7A3505D8d4"
requestTimeout = "168h" # 7 days
chunkSize = 25
backoffInitialDelay = "1m"
backoffMaxDelay = "2h"
observationSource = """
decode_log   [type=ethabidecodelog
              abi="RandomnessRequest(bytes32 keyHash,uint256 seed,bytes32 indexed jobID,address sender,uint256 fee,bytes32 requestID)"
              data="$(jobRun.logData)"
              topics="$(jobRun.logTopics)"]
vrf          [type=vrf
			  publicKey="$(jobSpec.publicKey)"
              requestBlockHash="$(jobRun.logBlockHash)"
              requestBlockNumber="$(jobRun.logBlockNumber)"
              topics="$(jobRun.logTopics)"]
encode_tx    [type=ethabiencode
              abi="fulfillRandomnessRequest(bytes proof)"
              data="{\\"proof\\": $(vrf)}"]
submit_tx  [type=ethtx to="%s"
			data="$(encode_tx)"
            txMeta="{\\"requestTxHash\\": $(jobRun.logTxHash),\\"requestID\\": $(decode_log.requestID),\\"jobID\\": $(jobSpec.databaseID)}"]
decode_log->vrf->encode_tx->submit_tx
"""
`
		jb, err := vrf.ValidatedVRFSpec(spec)
		require.NoError(tt, err)

		cfg := vrf_mocks.NewConfig(t)
		require.NoError(tt, vrf.CheckFromAddressMaxGasPrices(jb, cfg))
	})

	t.Run("returns nil error on valid gas lane <=> key specific gas price setting", func(tt *testing.T) {
		var fromAddresses []string
		for i := 0; i < 3; i++ {
			fromAddresses = append(fromAddresses, testutils.NewAddress().Hex())
		}

		cfg := vrf_mocks.NewConfig(t)
		for _, a := range fromAddresses {
			cfg.On("KeySpecificMaxGasPriceWei", common.HexToAddress(a)).Return(assets.GWei(100)).Once()
		}
		defer cfg.AssertExpectations(tt)

		jb, err := vrf.ValidatedVRFSpec(testspecs.GenerateVRFSpec(
			testspecs.VRFSpecParams{
				RequestedConfsDelay: 10,
				FromAddresses:       fromAddresses,
				ChunkSize:           25,
				BackoffInitialDelay: time.Minute,
				BackoffMaxDelay:     time.Hour,
				GasLanePrice:        assets.GWei(100),
			}).
			Toml())
		require.NoError(t, err)

		require.NoError(tt, vrf.CheckFromAddressMaxGasPrices(jb, cfg))
	})

	t.Run("returns error on invalid setting", func(tt *testing.T) {
		var fromAddresses []string
		for i := 0; i < 3; i++ {
			fromAddresses = append(fromAddresses, testutils.NewAddress().Hex())
		}

		cfg := vrf_mocks.NewConfig(t)
		cfg.On("KeySpecificMaxGasPriceWei", common.HexToAddress(fromAddresses[0])).Return(assets.GWei(100)).Once()
		cfg.On("KeySpecificMaxGasPriceWei", common.HexToAddress(fromAddresses[1])).Return(assets.GWei(100)).Once()
		// last from address has wrong key-specific max gas price
		cfg.On("KeySpecificMaxGasPriceWei", common.HexToAddress(fromAddresses[2])).Return(assets.GWei(50)).Once()
		defer cfg.AssertExpectations(tt)

		jb, err := vrf.ValidatedVRFSpec(testspecs.GenerateVRFSpec(
			testspecs.VRFSpecParams{
				RequestedConfsDelay: 10,
				FromAddresses:       fromAddresses,
				ChunkSize:           25,
				BackoffInitialDelay: time.Minute,
				BackoffMaxDelay:     time.Hour,
				GasLanePrice:        assets.GWei(100),
			}).
			Toml())
		require.NoError(t, err)

		require.Error(tt, vrf.CheckFromAddressMaxGasPrices(jb, cfg))
	})
}

func Test_CheckFromAddressesExist(t *testing.T) {
	t.Run("from addresses exist", func(t *testing.T) {
		db := pgtest.NewSqlxDB(t)
		cfg := configtest.NewTestGeneralConfig(t)
		lggr := logger.TestLogger(t)
		ks := keystore.New(db, utils.FastScryptParams, lggr, cfg)
		require.NoError(t, ks.Unlock(testutils.Password))

		var fromAddresses []string
		for i := 0; i < 3; i++ {
			k, err := ks.Eth().Create(big.NewInt(1337))
			assert.NoError(t, err)
			fromAddresses = append(fromAddresses, k.Address.Hex())
		}
		jb, err := vrf.ValidatedVRFSpec(testspecs.GenerateVRFSpec(
			testspecs.VRFSpecParams{
				RequestedConfsDelay: 10,
				FromAddresses:       fromAddresses,
				ChunkSize:           25,
				BackoffInitialDelay: time.Minute,
				BackoffMaxDelay:     time.Hour,
				GasLanePrice:        assets.GWei(100),
			}).
			Toml())
		assert.NoError(t, err)

		assert.NoError(t, vrf.CheckFromAddressesExist(jb, ks.Eth()))
	})

	t.Run("one of from addresses doesn't exist", func(t *testing.T) {
		db := pgtest.NewSqlxDB(t)
		cfg := configtest.NewTestGeneralConfig(t)
		lggr := logger.TestLogger(t)
		ks := keystore.New(db, utils.FastScryptParams, lggr, cfg)
		require.NoError(t, ks.Unlock(testutils.Password))

		var fromAddresses []string
		for i := 0; i < 3; i++ {
			k, err := ks.Eth().Create(big.NewInt(1337))
			assert.NoError(t, err)
			fromAddresses = append(fromAddresses, k.Address.Hex())
		}
		// add an address that isn't in the keystore
		fromAddresses = append(fromAddresses, testutils.NewAddress().Hex())
		jb, err := vrf.ValidatedVRFSpec(testspecs.GenerateVRFSpec(
			testspecs.VRFSpecParams{
				RequestedConfsDelay: 10,
				FromAddresses:       fromAddresses,
				ChunkSize:           25,
				BackoffInitialDelay: time.Minute,
				BackoffMaxDelay:     time.Hour,
				GasLanePrice:        assets.GWei(100),
			}).
			Toml())
		assert.NoError(t, err)

		assert.Error(t, vrf.CheckFromAddressesExist(jb, ks.Eth()))
	})
}

func Test_FromAddressMaxGasPricesAllEqual(t *testing.T) {
	t.Run("all max gas prices equal", func(tt *testing.T) {
		fromAddresses := []string{
			"0x498C2Dce1d3aEDE31A8c808c511C38a809e67684",
			"0x253b01b9CaAfbB9dC138d7D8c3ACBCDd47144b4B",
			"0xD94E6AD557277c6E3e163cefF90F52AB51A95143",
		}

		jb, err := vrf.ValidatedVRFSpec(testspecs.GenerateVRFSpec(testspecs.VRFSpecParams{
			RequestedConfsDelay: 10,
			FromAddresses:       fromAddresses,
			ChunkSize:           25,
			BackoffInitialDelay: time.Minute,
			BackoffMaxDelay:     time.Hour,
			GasLanePrice:        assets.GWei(100),
		}).Toml())
		require.NoError(tt, err)

		cfg := vrf_mocks.NewConfig(t)
		for _, a := range fromAddresses {
			cfg.On("KeySpecificMaxGasPriceWei", common.HexToAddress(a)).Return(assets.GWei(100))
		}
		defer cfg.AssertExpectations(tt)

		assert.True(tt, vrf.FromAddressMaxGasPricesAllEqual(jb, cfg))
	})

	t.Run("one max gas price not equal to others", func(tt *testing.T) {
		fromAddresses := []string{
			"0x498C2Dce1d3aEDE31A8c808c511C38a809e67684",
			"0x253b01b9CaAfbB9dC138d7D8c3ACBCDd47144b4B",
			"0xD94E6AD557277c6E3e163cefF90F52AB51A95143",
			"0x86E7c45Bf013Bf1Df3C22c14d5fd6fc3051AC569",
		}

		jb, err := vrf.ValidatedVRFSpec(testspecs.GenerateVRFSpec(testspecs.VRFSpecParams{
			RequestedConfsDelay: 10,
			FromAddresses:       fromAddresses,
			ChunkSize:           25,
			BackoffInitialDelay: time.Minute,
			BackoffMaxDelay:     time.Hour,
			GasLanePrice:        assets.GWei(100),
		}).Toml())
		require.NoError(tt, err)

		cfg := vrf_mocks.NewConfig(t)
		for _, a := range fromAddresses[:3] {
			cfg.On("KeySpecificMaxGasPriceWei", common.HexToAddress(a)).Return(assets.GWei(100))
		}
		cfg.On("KeySpecificMaxGasPriceWei", common.HexToAddress(fromAddresses[len(fromAddresses)-1])).
			Return(assets.GWei(200)) // doesn't match the rest
		defer cfg.AssertExpectations(tt)

		assert.False(tt, vrf.FromAddressMaxGasPricesAllEqual(jb, cfg))
	})
}
