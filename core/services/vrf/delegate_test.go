package vrf_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	eth_mocks "github.com/smartcontractkit/chainlink/core/services/eth/mocks"
	"github.com/smartcontractkit/chainlink/core/services/log"
	log_mocks "github.com/smartcontractkit/chainlink/core/services/log/mocks"
	pipeline_mocks "github.com/smartcontractkit/chainlink/core/services/pipeline/mocks"
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
	pr        *pipeline_mocks.Runner
	porm      *pipeline_mocks.ORM
	lb        *log_mocks.Broadcaster
	ec        *eth_mocks.Client
	vorm      vrf.ORM
	ks        *vrf.VRFKeyStore
	vrfkey    secp256k1.PublicKey
	submitter common.Address
}

func setup(t *testing.T, db *gorm.DB, s *store.Store) vrfUniverse {
	pr := new(pipeline_mocks.Runner)
	porm := new(pipeline_mocks.ORM)
	lb := new(log_mocks.Broadcaster)
	ec := new(eth_mocks.Client)
	vorm := vrf.NewORM(db)
	ks := vrf.NewVRFKeyStore(vrf.NewORM(db), utils.FastScryptParams)
	require.NoError(t, s.KeyStore.Unlock(cltest.Password))
	_, err := s.KeyStore.NewAccount()
	require.NoError(t, err)
	require.NoError(t, s.SyncDiskKeyStoreToDB())
	submitter, err := s.GetRoundRobinAddress(db)
	require.NoError(t, err)
	vrfkey, err := ks.CreateKey("blah")
	require.NoError(t, err)
	_, err = ks.Unlock("blah")
	require.NoError(t, err)
	return vrfUniverse{
		pr:        pr,
		porm:      porm,
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
	v.porm.AssertExpectations(t)
	v.pr.AssertExpectations(t)
	v.ec.AssertExpectations(t)
}

func TestDelegate(t *testing.T) {
	cfg, orm, cleanupDB := cltest.BootstrapThrowawayORM(t, "vrf_delegate", true)
	defer cleanupDB()
	store, cleanup := cltest.NewStoreWithConfig(t, cfg)
	defer cleanup()

	t.Run("creates a transaction on valid log", func(t *testing.T) {
		vuni := setup(t, orm.DB, store)
		vd := vrf.NewDelegate(orm.DB,
			vuni.vorm,
			store,
			vuni.ks,
			vuni.pr,
			vuni.porm,
			vuni.lb,
			vuni.ec,
			vrf.NewConfig(0, utils.FastScryptParams, 1000, 10))
		jb, err := vrf.ValidateVRFSpec(testspecs.GenerateVRFSpec(testspecs.VRFSpecParams{PublicKey: vuni.vrfkey.String()}))
		require.NoError(t, err)
		vl, err := vd.ServicesForSpec(jb)
		require.NoError(t, err)
		require.Len(t, vl, 1)

		listener := vl[0]
		done := make(chan struct{})
		unsubscribe := func() { done <- struct{}{} }

		var logListener log.Listener
		vuni.lb.On("Register", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			logListener = args.Get(0).(log.Listener)
		}).Return(unsubscribe)
		require.NoError(t, listener.Start())
		vuni.lb.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)
		vuni.pr.On("InsertFinishedRun", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(int64(1), nil)
		vuni.lb.On("MarkConsumed", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			done <- struct{}{}
		}).Return(nil)

		// Send a log with a valid key hash
		logListener.HandleLog(log.NewLogBroadcast(types.Log{
			Data: append(append(append(append(
				vuni.vrfkey.MustHash().Bytes(),
				common.BigToHash(big.NewInt(42)).Bytes()...),
				cltest.NewHash().Bytes()...),
				cltest.NewHash().Bytes()...),
				cltest.NewHash().Bytes()...),
			Topics:      []common.Hash{{}, common.BytesToHash([]byte("1234567890abcdef1234567890abcdef"))},
			Address:     common.Address{},
			BlockNumber: 0,
			TxHash:      common.Hash{},
			TxIndex:     0,
			BlockHash:   common.Hash{},
			Index:       0,
			Removed:     false,
		}))
		select {
		case <-time.After(1 * time.Second):
			t.Errorf("failed to consume log")
		case <-done:
			t.Log("woo done")
		}
		require.NoError(t, listener.Close())
		select {
		case <-time.After(1 * time.Second):
			t.Errorf("failed to unsubscribe")
		case <-done:
			t.Log("woo done")
		}
		vuni.lb.AssertExpectations(t)
	})

	t.Run("should not create an eth transction if invalid log", func(t *testing.T) {

	})
}
