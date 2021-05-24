package vrf_test

import (
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	eth_mocks "github.com/smartcontractkit/chainlink/core/services/eth/mocks"
	"github.com/smartcontractkit/chainlink/core/services/log"
	log_mocks "github.com/smartcontractkit/chainlink/core/services/log/mocks"
	pipeline_mocks "github.com/smartcontractkit/chainlink/core/services/pipeline/mocks"
	"github.com/smartcontractkit/chainlink/core/services/vrf"
	"github.com/smartcontractkit/chainlink/core/services/vrf/mocks"
	"github.com/smartcontractkit/chainlink/core/testdata/testspecs"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDelegate(t *testing.T) {
	_, orm, cleanupDB := cltest.BootstrapThrowawayORM(t, "vrf_delegate", true)
	defer cleanupDB()
	gethks := new(mocks.GethKeyStore)
	pr := new(pipeline_mocks.Runner)
	porm := new(pipeline_mocks.ORM)
	lb := new(log_mocks.Broadcaster)
	ec := new(eth_mocks.Client)
	vd := vrf.NewDelegate(orm.DB, gethks, pr, porm, lb, ec, vrf.NewConfig(0, utils.FastScryptParams, 1000, 10))
	jb, err := vrf.ValidateVRFSpec(testspecs.GenerateVRFSpec(testspecs.VRFSpecParams{}))
	require.NoError(t, err)
	vl, err := vd.ServicesForSpec(jb)
	require.NoError(t, err)
	require.Len(t, vl, 1)
	listener := vl[0]
	// Expect it to register
	unsubscribed := false
	unsubscribe := func() { unsubscribed = true }
	var logListener log.Listener
	lb.On("Register", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		logListener = args.Get(0).(log.Listener)
	}).Return(unsubscribe)
	require.NoError(t, listener.Start())

	lb.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)
	t.Log(logListener)
	//solidity_vrf_coordinator_interface.VRFCoordinator{}.ParseRandomnessRequest()
	// Should enqueue the log and get picked up
	//
	//data := cltest.MustGenericEncode([]string{"uint256"}, big.NewInt(0))
	//logListener.HandleLog(log.NewLogBroadcast(types.Log{
	//	Address:     common.Address{},
	//	Topics:      nil,
	//	Data:        data,
	//	BlockNumber: 0,
	//	TxHash:      common.Hash{},
	//	TxIndex:     0,
	//	BlockHash:   common.Hash{},
	//	Index:       0,
	//	Removed:     false,
	//}))
	//time.Sleep(1 * time.Second)

	require.NoError(t, listener.Close())
	lb.AssertExpectations(t)
	assert.True(t, unsubscribed)
}
