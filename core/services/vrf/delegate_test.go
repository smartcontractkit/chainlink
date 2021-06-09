package vrf_test

import (
	"context"
	"encoding/json"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/solidity_vrf_coordinator_interface"

	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/stretchr/testify/assert"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	eth_mocks "github.com/smartcontractkit/chainlink/core/services/eth/mocks"
	"github.com/smartcontractkit/chainlink/core/services/log"
	log_mocks "github.com/smartcontractkit/chainlink/core/services/log/mocks"
	"github.com/smartcontractkit/chainlink/core/services/signatures/secp256k1"
	"github.com/smartcontractkit/chainlink/core/services/vrf"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/testdata/testspecs"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type vrfUniverse struct {
	jpv2      cltest.JobPipelineV2TestHelper
	lb        *log_mocks.Broadcaster
	ec        *eth_mocks.Client
	vorm      vrf.ORM
	ks        *vrf.VRFKeyStore
	vrfkey    secp256k1.PublicKey
	submitter common.Address
}

func setup(t *testing.T, db *gorm.DB, cfg *cltest.TestConfig, s store.KeyStoreInterface) vrfUniverse {
	// Mock all chain interactions
	lb := new(log_mocks.Broadcaster)
	ec := new(eth_mocks.Client)

	// Don't mock db interactions
	jpv2 := cltest.NewJobPipelineV2(t, cfg, db)
	vorm := vrf.NewORM(db)
	ks := vrf.NewVRFKeyStore(vorm, utils.FastScryptParams)
	require.NoError(t, s.Unlock(cltest.Password))
	_, err := s.CreateNewKey()
	require.NoError(t, err)
	submitter, err := s.GetRoundRobinAddress()
	require.NoError(t, err)
	vrfkey, err := ks.CreateKey("blah")
	require.NoError(t, err)
	_, err = ks.Unlock("blah")
	require.NoError(t, err)

	return vrfUniverse{
		jpv2:      jpv2,
		lb:        lb,
		ec:        ec,
		vorm:      vorm,
		ks:        ks,
		vrfkey:    vrfkey,
		submitter: submitter,
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
		SeedAndBlockNum: cltest.NewHash(),
	})
	require.NoError(t, err)
	return b
}

func TestDelegate(t *testing.T) {
	cfg, cfgcleanup := cltest.NewConfig(t)
	t.Cleanup(cfgcleanup)
	store, cleanup := cltest.NewStoreWithConfig(t, cfg)
	t.Cleanup(cleanup)
	vuni := setup(t, store.DB, cfg, store.KeyStore)

	vd := vrf.NewDelegate(store.DB,
		vuni.vorm,
		store.KeyStore,
		vuni.ks,
		vuni.jpv2.Pr,
		vuni.jpv2.Prm,
		vuni.lb,
		vuni.ec,
		cfg)
	vs := testspecs.GenerateVRFSpec(testspecs.VRFSpecParams{PublicKey: vuni.vrfkey.String()})
	t.Log(vs)
	jb, err := vrf.ValidatedVRFSpec(vs.Toml())
	require.NoError(t, err)
	require.NoError(t, vuni.jpv2.Jrm.CreateJob(context.Background(), &jb, *pipeline.NewTaskDAG()))
	vl, err := vd.ServicesForSpec(jb)
	require.NoError(t, err)
	require.Len(t, vl, 1)

	listener := vl[0]
	unsubscribeAwaiter := cltest.NewAwaiter()
	unsubscribe := func() { unsubscribeAwaiter.ItHappened() }

	var logListener log.Listener
	vuni.lb.On("Register", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		logListener = args.Get(0).(log.Listener)
	}).Return(unsubscribe)
	require.NoError(t, listener.Start())

	t.Run("valid log", func(t *testing.T) {
		vuni.lb.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)
		a := cltest.NewAwaiter()
		vuni.lb.On("MarkConsumed", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			a.ItHappened()
		}).Return(nil).Once()
		// Expect a call to check if the req is already fulfilled.
		vuni.ec.On("CallContract", mock.Anything, mock.Anything, mock.Anything).Return(generateCallbackReturnValues(t), nil)

		// Send a valid log
		pk, err := secp256k1.NewPublicKeyFromHex(vs.PublicKey)
		require.NoError(t, err)
		reqID := cltest.NewHash()
		logListener.HandleLog(log.NewLogBroadcast(types.Log{
			// Data has all the NON-indexed parameters
			Data: append(append(append(append(
				pk.MustHash().Bytes(),                        // key hash
				common.BigToHash(big.NewInt(42)).Bytes()...), // seed
				cltest.NewHash().Bytes()...), // sender
				cltest.NewHash().Bytes()...), // fee
				reqID.Bytes()...), // requestID
			// JobID is indexed, thats why it lives in the Topics.
			Topics:      []common.Hash{{}, jb.ExternalIDToTopicHash()}, // jobID
			Address:     common.Address{},
			BlockNumber: 0,
			TxHash:      common.Hash{},
			TxIndex:     0,
			BlockHash:   common.Hash{},
			Index:       0,
			Removed:     false,
		}))
		a.AwaitOrFail(t)

		// Ensure we created a successful run.
		runs, err := vuni.jpv2.Prm.GetAllRuns()
		require.NoError(t, err)
		require.Equal(t, 1, len(runs))
		assert.False(t, runs[0].Errors.HasError())
		m, ok := runs[0].Meta.Val.(map[string]interface{})
		require.True(t, ok)
		_, ok = m["eth_tx_id"]
		assert.True(t, ok)
		assert.Len(t, runs[0].PipelineTaskRuns, 0)

		// Ensure we have queued up a valid eth transaction
		// Linked to  requestID
		var ethTxes []models.EthTx
		err = store.DB.Find(&ethTxes).Error
		require.NoError(t, err)
		require.Equal(t, 1, len(ethTxes))
		assert.Equal(t, vs.CoordinatorAddress, ethTxes[0].ToAddress.String())
		var em models.EthTxMetaV2
		err = json.Unmarshal(ethTxes[0].Meta.RawMessage, &em)
		require.NoError(t, err)
		assert.Equal(t, reqID, em.RequestID)
		require.NoError(t, store.DB.Exec(`TRUNCATE eth_txes,pipeline_runs CASCADE`).Error)
	})

	t.Run("invalid log", func(t *testing.T) {
		vuni.lb.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)
		a := cltest.NewAwaiter()
		vuni.lb.On("MarkConsumed", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			a.ItHappened()
		}).Return(nil).Once()
		// Expect a call to check if the req is already fulfilled.
		vuni.ec.On("CallContract", mock.Anything, mock.Anything, mock.Anything).Return(generateCallbackReturnValues(t), nil)

		// Send a invalid log (keyhash doesnt match)
		logListener.HandleLog(log.NewLogBroadcast(types.Log{
			// Data has all the NON-indexed parameters
			Data: append(append(append(append(
				cltest.NewHash().Bytes(),                     // key hash
				common.BigToHash(big.NewInt(42)).Bytes()...), // seed
				cltest.NewHash().Bytes()...), // sender
				cltest.NewHash().Bytes()...), // fee
				cltest.NewHash().Bytes()...), // requestID
			// JobID is indexed, thats why it lives in the Topics.
			Topics:      []common.Hash{{}, jb.ExternalIDToTopicHash()}, // jobID
			Address:     common.Address{},
			BlockNumber: 0,
			TxHash:      common.Hash{},
			TxIndex:     0,
			BlockHash:   common.Hash{},
			Index:       0,
			Removed:     false,
		}))
		a.AwaitOrFail(t)

		// Ensure we have not created a run.
		runs, err := vuni.jpv2.Prm.GetAllRuns()
		require.NoError(t, err)
		require.Equal(t, len(runs), 0)

		// Ensure we have NOT queued up an eth transaction
		var ethTxes []models.EthTx
		err = store.DB.Find(&ethTxes).Error
		require.NoError(t, err)
		require.Len(t, ethTxes, 0)
	})

	require.NoError(t, listener.Close())
	unsubscribeAwaiter.AwaitOrFail(t, 1*time.Second)
	vuni.Assert(t)
}
