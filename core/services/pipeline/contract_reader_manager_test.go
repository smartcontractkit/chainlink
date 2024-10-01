package pipeline

import (
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/smartcontractkit/chainlink-ccip/mocks/pkg/contractreader"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/mocks"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestCheckForUnusedClients(t *testing.T) {
	var chStop services.StopChan
	lggr, logs := logger.TestLoggerObserved(t, zap.DebugLevel)
	networkID1 := types.NewRelayID("network1", "chain1")
	networkID2 := types.NewRelayID("network2", "chain2")
	relay1 := mocks.NewRelayer(t)
	relay2 := mocks.NewRelayer(t)

	relayers := map[types.RelayID]loop.Relayer{
		networkID1: relay1,
		networkID2: relay2,
	}
	contractReader1 := contractreader.NewMockContractReaderFacade(t)
	contractReader2 := contractreader.NewMockContractReaderFacade(t)
	contractReader3 := contractreader.NewMockContractReaderFacade(t)

	contractReader1.On("Bind", mock.Anything, mock.Anything).Return(nil)
	contractReader2.On("Bind", mock.Anything, mock.Anything).Return(nil)
	contractReader3.On("Bind", mock.Anything, mock.Anything).Return(nil)

	relay1.On("NewContractReader", mock.Anything, []byte("contractReader1")).Return(contractReader1, nil)
	relay1.On("NewContractReader", mock.Anything, []byte("contractReader2")).Return(contractReader2, nil)
	relay2.On("NewContractReader", mock.Anything, mock.Anything).Return(contractReader3, nil)

	manager, err := newContractReaderManager(relayers, chStop, lggr)
	require.NoError(t, err)
	manager.heartBeatCheckInterval = time.Millisecond * 500
	manager.hearthBeatTimeout = time.Second * 3

	//Create three contract readers
	_, _, err = manager.GetOrCreate(networkID1, "contract1", "contractAddress1", "method1", []byte("contractReader1"))
	require.NoError(t, err)
	_, _, err = manager.GetOrCreate(networkID1, "contract2", "contractAddress2", "method2", []byte("contractReader2"))
	require.NoError(t, err)
	_, _, err = manager.GetOrCreate(networkID2, "contract1", "contractAddress1", "method1", []byte{})
	require.NoError(t, err)

	//Make sure the manager has all three contract readers
	require.Len(t, manager.crs, 3)

	time.Sleep(time.Second * 2)
	//Call contract reader on relay1 to refresh the timeout
	_, _, err = manager.GetOrCreate(networkID1, "contract1", "contractAddress1", "method1", []byte("contractReader1"))
	require.NoError(t, err)

	//All contract readers should have timed out except contractReader1
	require.Eventually(t, func() bool { return len(manager.crs) == 1 }, time.Second*5, time.Millisecond*500)

	l := logs.FilterMessageSnippet("closing contractReader with ID \"network1_chain1_contractAddress2_method2\"").TakeAll()
	require.Len(t, l, 1)

	l = logs.FilterMessageSnippet("closing contractReader with ID \"network2_chain2_contractAddress1_method1\"").TakeAll()
	require.Len(t, l, 1)
}

func TestContractReaderManager(t *testing.T) {
	var chStop services.StopChan
	networkID1 := types.NewRelayID("network1", "chain1")
	networkID2 := types.NewRelayID("network2", "chain2")

	relay1 := mocks.NewRelayer(t)
	relay2 := mocks.NewRelayer(t)
	relayers := map[types.RelayID]loop.Relayer{
		networkID1: relay1,
		networkID2: relay2,
	}

	contractReader1 := contractreader.NewMockContractReaderFacade(t)
	contractReader2 := contractreader.NewMockContractReaderFacade(t)

	contractReader1.On("Bind", mock.Anything, mock.Anything).Return(nil)
	contractReader2.On("Bind", mock.Anything, mock.Anything).Return(nil)

	relay1.On("NewContractReader", mock.Anything, mock.Anything).Return(contractReader1, nil)
	relay2.On("NewContractReader", mock.Anything, mock.Anything).Return(contractReader2, nil)
	manager, err := newContractReaderManager(relayers, chStop, logger.TestLogger(t))
	require.NoError(t, err)

	//CSR created with no error
	_, id, err := manager.GetOrCreate(networkID1, "contract1", "contractAddress1", "method1", []byte{})
	require.NoError(t, err)

	//CSR already exists
	_, id2, err := manager.GetOrCreate(networkID1, "contract1", "contractAddress1", "method1", []byte{})
	require.Equal(t, id, id2)

	//CSR created with no error
	_, _, err = manager.GetOrCreate(networkID2, "contract1", "contractAddress1", "method1", []byte{})
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
