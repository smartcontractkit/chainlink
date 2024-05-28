package vrf_test

import (
	"bytes"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/mailbox/mailboxtest"

	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	evmclimocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker"
	httypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/log"
	log_mocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/log/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	evmutils "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/solidity_vrf_coordinator_interface"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/vrfkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	evmrelay "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
	"github.com/smartcontractkit/chainlink/v2/core/services/signatures/secp256k1"
	"github.com/smartcontractkit/chainlink/v2/core/services/vrf"
	vrf_mocks "github.com/smartcontractkit/chainlink/v2/core/services/vrf/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/vrf/solidity_cross_tests"
	v1 "github.com/smartcontractkit/chainlink/v2/core/services/vrf/v1"
	"github.com/smartcontractkit/chainlink/v2/core/services/vrf/vrfcommon"
	"github.com/smartcontractkit/chainlink/v2/core/testdata/testspecs"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type vrfUniverse struct {
	jrm          job.ORM
	pr           pipeline.Runner
	prm          pipeline.ORM
	lb           *log_mocks.Broadcaster
	ec           *evmclimocks.Client
	ks           keystore.Master
	vrfkey       vrfkey.KeyV2
	submitter    common.Address
	txm          *txmgr.TxManager
	hb           httypes.HeadBroadcaster
	legacyChains legacyevm.LegacyChainContainer
	cid          big.Int
}

func buildVrfUni(t *testing.T, db *sqlx.DB, cfg chainlink.GeneralConfig) vrfUniverse {
	ctx := testutils.Context(t)
	// Mock all chain interactions
	lb := log_mocks.NewBroadcaster(t)
	lb.On("AddDependents", 1).Maybe()
	lb.On("Register", mock.Anything, mock.Anything).Return(func() {}).Maybe()
	ec := evmclimocks.NewClient(t)
	ec.On("ConfiguredChainID").Return(testutils.FixtureChainID)
	ec.On("LatestBlockHeight", mock.Anything).Return(big.NewInt(51), nil).Maybe()
	lggr := logger.TestLogger(t)
	hb := headtracker.NewHeadBroadcaster(lggr)

	// Don't mock db interactions
	prm := pipeline.NewORM(db, lggr, cfg.JobPipeline().MaxSuccessfulRuns())
	btORM := bridges.NewORM(db)
	ks := keystore.NewInMemory(db, utils.FastScryptParams, lggr)
	_, dbConfig, evmConfig := txmgr.MakeTestConfigs(t)
	txm, err := txmgr.NewTxm(db, evmConfig, evmConfig.GasEstimator(), evmConfig.Transactions(), nil, dbConfig, dbConfig.Listener(), ec, logger.TestLogger(t), nil, ks.Eth(), nil)
	orm := headtracker.NewORM(*testutils.FixtureChainID, db)
	require.NoError(t, orm.IdempotentInsertHead(testutils.Context(t), cltest.Head(51)))
	jrm := job.NewORM(db, prm, btORM, ks, lggr)
	t.Cleanup(func() { assert.NoError(t, jrm.Close()) })
	relayExtenders := evmtest.NewChainRelayExtenders(t, evmtest.TestChainOpts{LogBroadcaster: lb, KeyStore: ks.Eth(), Client: ec, DB: db, GeneralConfig: cfg, TxManager: txm})
	legacyChains := evmrelay.NewLegacyChainsFromRelayerExtenders(relayExtenders)
	pr := pipeline.NewRunner(prm, btORM, cfg.JobPipeline(), cfg.WebServer(), legacyChains, ks.Eth(), ks.VRF(), lggr, nil, nil)
	require.NoError(t, ks.Unlock(ctx, testutils.Password))
	k, err2 := ks.Eth().Create(testutils.Context(t), testutils.FixtureChainID)
	require.NoError(t, err2)
	submitter := k.Address
	require.NoError(t, err)
	vrfkey, err3 := ks.VRF().Create(ctx)
	require.NoError(t, err3)

	return vrfUniverse{
		jrm:          jrm,
		pr:           pr,
		prm:          prm,
		lb:           lb,
		ec:           ec,
		ks:           ks,
		vrfkey:       vrfkey,
		submitter:    submitter,
		txm:          &txm,
		hb:           hb,
		legacyChains: legacyChains,
		cid:          *ec.ConfiguredChainID(),
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
		b, err2 := args.Pack(solidity_vrf_coordinator_interface.Callbacks{
			RandomnessFee:   big.NewInt(10),
			SeedAndBlockNum: evmutils.EmptyHash,
		})
		require.NoError(t, err2)
		return b
	}
	b, err := args.Pack(solidity_vrf_coordinator_interface.Callbacks{
		RandomnessFee:   big.NewInt(10),
		SeedAndBlockNum: evmutils.NewHash(),
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

func setup(t *testing.T) (vrfUniverse, *v1.Listener, job.Job) {
	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewTestGeneralConfig(t)
	vuni := buildVrfUni(t, db, cfg)

	mailMon := servicetest.Run(t, mailboxtest.NewMonitor(t))

	vd := vrf.NewDelegate(
		db,
		vuni.ks,
		vuni.pr,
		vuni.prm,
		vuni.legacyChains,
		logger.TestLogger(t),
		mailMon)
	vs := testspecs.GenerateVRFSpec(testspecs.VRFSpecParams{PublicKey: vuni.vrfkey.PublicKey.String(), EVMChainID: testutils.FixtureChainID.String()})
	jb, err := vrfcommon.ValidatedVRFSpec(vs.Toml())
	require.NoError(t, err)
	ctx := testutils.Context(t)
	err = vuni.jrm.CreateJob(ctx, &jb)
	require.NoError(t, err)
	vl, err := vd.ServicesForSpec(testutils.Context(t), jb)
	require.NoError(t, err)
	require.Len(t, vl, 1)
	listener := vl[0].(*v1.Listener)
	// Start the listenerV1
	go func() {
		listener.RunLogListener([]func(){}, 6)
	}()
	go func() {
		listener.RunHeadListener(func() {})
	}()
	servicetest.Run(t, listener)
	return vuni, listener, jb
}

func TestDelegate_ReorgAttackProtection(t *testing.T) {
	vuni, listener, jb := setup(t)

	// Same request has already been fulfilled twice
	reqID := evmutils.NewHash()
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
	txHash := evmutils.NewHash()
	vuni.lb.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil).Maybe()
	vuni.lb.On("MarkConsumed", mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()
	vuni.ec.On("CallContract", mock.Anything, mock.Anything, mock.Anything).Return(generateCallbackReturnValues(t, false), nil).Maybe()
	ctx := testutils.Context(t)
	listener.HandleLog(ctx, log.NewLogBroadcast(types.Log{
		// Data has all the NON-indexed parameters
		Data: bytes.Join([][]byte{pk.MustHash().Bytes(), // key hash
			preSeed,                    // preSeed
			evmutils.NewHash().Bytes(), // sender
			evmutils.NewHash().Bytes(), // fee
			reqID.Bytes()}, []byte{},   // requestID
		),
		// JobID is indexed, thats why it lives in the Topics.
		Topics: []common.Hash{
			solidity_cross_tests.VRFRandomnessRequestLogTopic(),
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
	txHash := evmutils.NewHash()
	reqID1 := evmutils.NewHash()
	reqID2 := evmutils.NewHash()
	keyID := vuni.vrfkey.PublicKey.String()
	pk, err := secp256k1.NewPublicKeyFromHex(keyID)
	require.NoError(t, err)
	added := make(chan struct{})
	listener.SetReqAdded(func() {
		added <- struct{}{}
	})
	preSeed := common.BigToHash(big.NewInt(42)).Bytes()
	bh := evmutils.NewHash()
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
					evmutils.NewHash().Bytes(),               // sender
					evmutils.NewHash().Bytes(),               // fee
					reqID1.Bytes()},                          // requestID
					[]byte{}),
				// JobID is indexed, thats why it lives in the Topics.
				Topics: []common.Hash{
					solidity_cross_tests.VRFRandomnessRequestLogTopic(),
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
					evmutils.NewHash().Bytes(),               // sender
					evmutils.NewHash().Bytes(),               // fee
					reqID2.Bytes()},                          // requestID
					[]byte{}),
				Topics: []common.Hash{
					solidity_cross_tests.VRFRandomnessRequestLogTopic(),
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
		ctx := testutils.Context(t)
		vuni.lb.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)
		vuni.lb.On("MarkConsumed", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			consumed <- struct{}{}
		}).Return(nil).Once()
		// Expect a call to check if the req is already fulfilled.
		vuni.ec.On("CallContract", mock.Anything, mock.Anything, mock.Anything).Return(generateCallbackReturnValues(t, false), nil)

		listener.HandleLog(ctx, log.NewLogBroadcast(tc.log, vuni.cid, nil))
		// Wait until the log is present
		waitForChannel(t, added, time.Second, "request not added to the queue")
		// Feed it a head which confirms it.
		listener.OnNewLongestChain(testutils.Context(t), &evmtypes.Head{Number: 16})
		waitForChannel(t, consumed, 2*time.Second, "did not mark consumed")

		// Ensure we created a successful run.
		waitForChannel(t, runComplete, 2*time.Second, "pipeline not complete")
		runs, err := vuni.prm.GetAllRuns(ctx)
		require.NoError(t, err)
		require.Equal(t, i+1, len(runs))
		assert.False(t, runs[0].FatalErrors.HasError())
		// Should have 4 tasks all completed
		assert.Len(t, runs[0].PipelineTaskRuns, 4)

		p, err := vuni.ks.VRF().GenerateProof(keyID, evmutils.MustHash(string(bytes.Join([][]byte{preSeed, bh.Bytes()}, []byte{}))).Big())
		require.NoError(t, err)
		vuni.lb.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)
		vuni.lb.On("MarkConsumed", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			consumed <- struct{}{}
		}).Return(nil).Once()
		// If we send a completed log we should the respCount increase
		var reqIDBytes []byte
		copy(reqIDBytes[:], tc.reqID[:])
		listener.HandleLog(ctx, log.NewLogBroadcast(types.Log{
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
	vuni.lb.On("MarkConsumed", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		done <- struct{}{}
	}).Return(nil).Once()
	// Expect a call to check if the req is already fulfilled.
	vuni.ec.On("CallContract", mock.Anything, mock.Anything, mock.Anything).Return(generateCallbackReturnValues(t, false), nil)

	added := make(chan struct{})
	listener.SetReqAdded(func() {
		added <- struct{}{}
	})
	// Send an invalid log (keyhash doesnt match)
	ctx := testutils.Context(t)
	listener.HandleLog(ctx, log.NewLogBroadcast(types.Log{
		// Data has all the NON-indexed parameters
		Data: append(append(append(append(
			evmutils.NewHash().Bytes(),                   // key hash
			common.BigToHash(big.NewInt(42)).Bytes()...), // seed
			evmutils.NewHash().Bytes()...), // sender
			evmutils.NewHash().Bytes()...), // fee
			evmutils.NewHash().Bytes()...), // requestID
		// JobID is indexed, that's why it lives in the Topics.
		Topics: []common.Hash{
			solidity_cross_tests.VRFRandomnessRequestLogTopic(),
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
	runs, err := vuni.prm.GetAllRuns(ctx)
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

	db := pgtest.NewSqlxDB(t)
	txStore := txmgr.NewTxStore(db, logger.TestLogger(t))

	txes, err := txStore.GetAllTxes(testutils.Context(t))
	require.NoError(t, err)
	require.Len(t, txes, 0)
}

func TestFulfilledCheck(t *testing.T) {
	vuni, listener, jb := setup(t)
	vuni.lb.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)
	done := make(chan struct{})
	vuni.lb.On("MarkConsumed", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
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
	ctx := testutils.Context(t)
	listener.HandleLog(ctx, log.NewLogBroadcast(
		types.Log{
			// Data has all the NON-indexed parameters
			Data: bytes.Join([][]byte{
				vuni.vrfkey.PublicKey.MustHash().Bytes(), // key hash
				common.BigToHash(big.NewInt(42)).Bytes(), // seed
				evmutils.NewHash().Bytes(),               // sender
				evmutils.NewHash().Bytes(),               // fee
				evmutils.NewHash().Bytes()},              // requestID
				[]byte{}),
			// JobID is indexed, that's why it lives in the Topics.
			Topics: []common.Hash{
				solidity_cross_tests.VRFRandomnessRequestLogTopic(),
				jb.ExternalIDEncodeBytesToTopic(), // jobID STRING
			},
			//TxHash:      evmutils.NewHash().Bytes(),
			BlockNumber: 10,
			//BlockHash:   evmutils.NewHash().Bytes(),
		}, vuni.cid, nil))

	// Should queue the request, even though its already fulfilled
	waitForChannel(t, added, time.Second, "request not queued")
	listener.OnNewLongestChain(testutils.Context(t), &evmtypes.Head{Number: 16})
	waitForChannel(t, done, time.Second, "log not consumed")

	// Should consume the log with no run
	runs, err := vuni.prm.GetAllRuns(ctx)
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
		jb, err := vrfcommon.ValidatedVRFSpec(spec)
		require.NoError(tt, err)

		cfg := vrf_mocks.NewFeeConfig(t)
		require.NoError(tt, vrf.CheckFromAddressMaxGasPrices(jb, cfg.PriceMaxKey))
	})

	t.Run("returns nil error on valid gas lane <=> key specific gas price setting", func(tt *testing.T) {
		var fromAddresses []string
		for i := 0; i < 3; i++ {
			fromAddresses = append(fromAddresses, testutils.NewAddress().Hex())
		}

		cfg := vrf_mocks.NewFeeConfig(t)
		for _, a := range fromAddresses {
			cfg.On("PriceMaxKey", common.HexToAddress(a)).Return(assets.GWei(100)).Once()
		}
		defer cfg.AssertExpectations(tt)

		jb, err := vrfcommon.ValidatedVRFSpec(testspecs.GenerateVRFSpec(
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

		require.NoError(tt, vrf.CheckFromAddressMaxGasPrices(jb, cfg.PriceMaxKey))
	})

	t.Run("returns error on invalid setting", func(tt *testing.T) {
		var fromAddresses []string
		for i := 0; i < 3; i++ {
			fromAddresses = append(fromAddresses, testutils.NewAddress().Hex())
		}

		cfg := vrf_mocks.NewFeeConfig(t)
		cfg.On("PriceMaxKey", common.HexToAddress(fromAddresses[0])).Return(assets.GWei(100)).Once()
		cfg.On("PriceMaxKey", common.HexToAddress(fromAddresses[1])).Return(assets.GWei(100)).Once()
		// last from address has wrong key-specific max gas price
		cfg.On("PriceMaxKey", common.HexToAddress(fromAddresses[2])).Return(assets.GWei(50)).Once()
		defer cfg.AssertExpectations(tt)

		jb, err := vrfcommon.ValidatedVRFSpec(testspecs.GenerateVRFSpec(
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

		require.Error(tt, vrf.CheckFromAddressMaxGasPrices(jb, cfg.PriceMaxKey))
	})
}

func Test_CheckFromAddressesExist(t *testing.T) {
	t.Run("from addresses exist", func(t *testing.T) {
		ctx := testutils.Context(t)
		db := pgtest.NewSqlxDB(t)
		lggr := logger.TestLogger(t)
		ks := keystore.NewInMemory(db, utils.FastScryptParams, lggr)
		require.NoError(t, ks.Unlock(ctx, testutils.Password))

		var fromAddresses []string
		for i := 0; i < 3; i++ {
			k, err := ks.Eth().Create(testutils.Context(t), big.NewInt(1337))
			assert.NoError(t, err)
			fromAddresses = append(fromAddresses, k.Address.Hex())
		}
		jb, err := vrfcommon.ValidatedVRFSpec(testspecs.GenerateVRFSpec(
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

		assert.NoError(t, vrf.CheckFromAddressesExist(testutils.Context(t), jb, ks.Eth()))
	})

	t.Run("one of from addresses doesn't exist", func(t *testing.T) {
		ctx := testutils.Context(t)
		db := pgtest.NewSqlxDB(t)
		lggr := logger.TestLogger(t)
		ks := keystore.NewInMemory(db, utils.FastScryptParams, lggr)
		require.NoError(t, ks.Unlock(ctx, testutils.Password))

		var fromAddresses []string
		for i := 0; i < 3; i++ {
			k, err := ks.Eth().Create(testutils.Context(t), big.NewInt(1337))
			assert.NoError(t, err)
			fromAddresses = append(fromAddresses, k.Address.Hex())
		}
		// add an address that isn't in the keystore
		fromAddresses = append(fromAddresses, testutils.NewAddress().Hex())
		jb, err := vrfcommon.ValidatedVRFSpec(testspecs.GenerateVRFSpec(
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

		assert.Error(t, vrf.CheckFromAddressesExist(testutils.Context(t), jb, ks.Eth()))
	})
}

func Test_FromAddressMaxGasPricesAllEqual(t *testing.T) {
	t.Run("all max gas prices equal", func(tt *testing.T) {
		fromAddresses := []string{
			"0x498C2Dce1d3aEDE31A8c808c511C38a809e67684",
			"0x253b01b9CaAfbB9dC138d7D8c3ACBCDd47144b4B",
			"0xD94E6AD557277c6E3e163cefF90F52AB51A95143",
		}

		jb, err := vrfcommon.ValidatedVRFSpec(testspecs.GenerateVRFSpec(testspecs.VRFSpecParams{
			RequestedConfsDelay: 10,
			FromAddresses:       fromAddresses,
			ChunkSize:           25,
			BackoffInitialDelay: time.Minute,
			BackoffMaxDelay:     time.Hour,
			GasLanePrice:        assets.GWei(100),
		}).Toml())
		require.NoError(tt, err)

		cfg := vrf_mocks.NewFeeConfig(t)
		for _, a := range fromAddresses {
			cfg.On("PriceMaxKey", common.HexToAddress(a)).Return(assets.GWei(100))
		}
		defer cfg.AssertExpectations(tt)

		assert.True(tt, vrf.FromAddressMaxGasPricesAllEqual(jb, cfg.PriceMaxKey))
	})

	t.Run("one max gas price not equal to others", func(tt *testing.T) {
		fromAddresses := []string{
			"0x498C2Dce1d3aEDE31A8c808c511C38a809e67684",
			"0x253b01b9CaAfbB9dC138d7D8c3ACBCDd47144b4B",
			"0xD94E6AD557277c6E3e163cefF90F52AB51A95143",
			"0x86E7c45Bf013Bf1Df3C22c14d5fd6fc3051AC569",
		}

		jb, err := vrfcommon.ValidatedVRFSpec(testspecs.GenerateVRFSpec(testspecs.VRFSpecParams{
			RequestedConfsDelay: 10,
			FromAddresses:       fromAddresses,
			ChunkSize:           25,
			BackoffInitialDelay: time.Minute,
			BackoffMaxDelay:     time.Hour,
			GasLanePrice:        assets.GWei(100),
		}).Toml())
		require.NoError(tt, err)

		cfg := vrf_mocks.NewFeeConfig(t)
		for _, a := range fromAddresses[:3] {
			cfg.On("PriceMaxKey", common.HexToAddress(a)).Return(assets.GWei(100))
		}
		cfg.On("PriceMaxKey", common.HexToAddress(fromAddresses[len(fromAddresses)-1])).
			Return(assets.GWei(200)) // doesn't match the rest
		defer cfg.AssertExpectations(tt)

		assert.False(tt, vrf.FromAddressMaxGasPricesAllEqual(jb, cfg.PriceMaxKey))
	})
}

func Test_VRFV2PlusServiceFailsWhenVRFOwnerProvided(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewTestGeneralConfig(t)
	vuni := buildVrfUni(t, db, cfg)

	mailMon := servicetest.Run(t, mailboxtest.NewMonitor(t))

	vd := vrf.NewDelegate(
		db,
		vuni.ks,
		vuni.pr,
		vuni.prm,
		vuni.legacyChains,
		logger.TestLogger(t),
		mailMon)
	chain, err := vuni.legacyChains.Get(testutils.FixtureChainID.String())
	require.NoError(t, err)
	vs := testspecs.GenerateVRFSpec(testspecs.VRFSpecParams{
		VRFVersion:    vrfcommon.V2Plus,
		PublicKey:     vuni.vrfkey.PublicKey.String(),
		FromAddresses: []string{vuni.submitter.Hex()},
		GasLanePrice:  chain.Config().EVM().GasEstimator().PriceMax(),
	})
	toml := "vrfOwnerAddress=\"0xF62fEFb54a0af9D32CDF0Db21C52710844c7eddb\"\n" + vs.Toml()
	jb, err := vrfcommon.ValidatedVRFSpec(toml)
	require.NoError(t, err)
	ctx := testutils.Context(t)
	err = vuni.jrm.CreateJob(ctx, &jb)
	require.NoError(t, err)
	_, err = vd.ServicesForSpec(testutils.Context(t), jb)
	require.Error(t, err)
	require.Equal(t, "VRF Owner is not supported for VRF V2 Plus", err.Error())
}
