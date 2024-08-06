package pipeline

import (
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/mocks"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestCheckForUnusedClients(t *testing.T) {
	var chStop services.StopChan
	lggr, logs := logger.TestLoggerObserved(t, zap.DebugLevel)
	rID1 := types.NewRelayID("network1", "chain1")
	rID2 := types.NewRelayID("network1", "chain2")
	rID3 := types.NewRelayID("network2", "chain1")
	relayer1 := mocks.NewRelayer(t)
	relayer2 := mocks.NewRelayer(t)
	relayer3 := mocks.NewRelayer(t)
	relayers := map[types.RelayID]loop.Relayer{
		rID1: relayer1,
		rID2: relayer2,
		rID3: relayer3,
	}

	relayer1.On("NewContractReader", mock.Anything, mock.Anything).Return(nil, nil)
	relayer2.On("NewContractReader", mock.Anything, mock.Anything).Return(nil, nil)
	relayer3.On("NewContractReader", mock.Anything, mock.Anything).Return(nil, nil)
	c, err := newContractReaderManager(relayers, chStop, lggr)
	require.NoError(t, err)
	c.heartBeatCheckInterval = time.Second
	c.hearthBeatTimeout = time.Second * 3

	relayer1.On("NewContractReader", mock.Anything, mock.Anything).Return(nil, nil)
	relayer2.On("NewContractReader", mock.Anything, mock.Anything).Return(nil, nil)

	//CSR created with no error
	_, err = c.Create(rID1, "contract1", "method1", []byte{})
	require.NoError(t, err)
	_, err = c.Create(rID2, "contract1", "method1", []byte{})
	require.NoError(t, err)
	_, err = c.Create(rID3, "contract1", "method1", []byte{})
	require.NoError(t, err)

	//There should be 3 CSR in the map
	require.Len(t, c.crs, 3)

	time.Sleep(time.Second * 2)
	//Call CSR on relayer3 to refresh the timeout
	_, err = c.Get(rID3, "contract1", "method1")
	require.NoError(t, err)

	time.Sleep(time.Second * 2)
	//All CSR should have timed out except relayer 3
	require.Len(t, c.crs, 1)

	l := logs.FilterMessageSnippet("closing").TakeAll()
	require.Len(t, l, 2)
	require.Contains(t, l[0].Entry.Message, "closing contractReader with ID \"network1_chain1_contract1_method1\"")
	require.Contains(t, l[1].Entry.Message, "closing contractReader with ID \"network1_chain2_contract1_method1\"")
}

func TestCSRM(t *testing.T) {
	var chStop services.StopChan
	rID1 := types.NewRelayID("network1", "chain1")
	rID2 := types.NewRelayID("network1", "chain2")

	relayer1 := mocks.NewRelayer(t)
	relayer2 := mocks.NewRelayer(t)
	relayers := map[types.RelayID]loop.Relayer{
		rID1: relayer1,
		rID2: relayer2,
	}

	relayer1.On("NewContractReader", mock.Anything, mock.Anything).Return(nil, nil)
	relayer2.On("NewContractReader", mock.Anything, mock.Anything).Return(nil, nil)
	csrm, err := newContractReaderManager(relayers, chStop, logger.TestLogger(t))
	require.NoError(t, err)

	//Manager not created for contract + method
	_, err = csrm.Get(rID1, "contract1", "method1")
	require.ErrorIs(t, err, ErrContractReaderNotFound)

	//CSR created with no error
	_, err = csrm.Create(rID1, "contract1", "method1", []byte{})
	require.NoError(t, err)

	//CSR already exists
	_, err = csrm.Create(rID1, "contract1", "method1", []byte{})
	require.ErrorContains(t, err, "contractReader already exists")

	//CSR created with no error
	_, err = csrm.Create(rID2, "contract1", "method1", []byte{})
	require.NoError(t, err)
}

func TestCSRMCreateId(t *testing.T) {
	//No chainID
	_, err := createID(types.NewRelayID("", ""), "", "")
	require.ErrorContains(t, err, "cannot create ID, chainID is empty")

	//No network
	_, err = createID(types.NewRelayID("", "chainID"), "", "")
	require.ErrorContains(t, err, "cannot create ID, network is empty")

	//No contractAddress
	_, err = createID(types.NewRelayID("network", "chainID"), "", "")
	require.ErrorContains(t, err, "cannot create ID, contractAddress is empty")

	//No methodName
	_, err = createID(types.NewRelayID("network", "chainID"), "contract", "")
	require.ErrorContains(t, err, "cannot create ID, methodName is empty")

	id, err := createID(types.NewRelayID("network", "chainID"), "contract", "method")
	require.NoError(t, err)

	require.Equal(t, "network_chainID_contract_method", id)
}
