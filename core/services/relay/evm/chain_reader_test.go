package evm_test

import (
	"context"
	"log"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/codec"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/testutils"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	evmtypes "github.com/smartcontractkit/chainlink-integrations/evm/types"

	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

func TestChainReaderPrimitive(t *testing.T) {
	t.Parallel()

	tester := &simpleTester{}
	wrapped := testutils.WrapContractReaderTesterForLoop(tester)
	wrapped.Setup(t)

	svc := wrapped.GetContractReader(t)
	binding := commontypes.BoundContract{Address: "0x45", Name: "Contract"}

	require.NoError(t, svc.Bind(t.Context(), []commontypes.BoundContract{binding}))

	var value values.Value
	require.NoError(t, svc.GetLatestValue(t.Context(), binding.ReadIdentifier("GetValue"), primitives.Finalized, nil, &value))

	t.Fail()
}

type mockedClient struct{}

func (_m *mockedClient) BatchCallContext(_ context.Context, _ []rpc.BatchElem) error {
	return nil
}

func (_m *mockedClient) CallContract(_ context.Context, msg ethereum.CallMsg, _ *big.Int) ([]byte, error) {
	log.Println(msg.Data)

	tp, _ := abi.NewType("uint256", "", nil)
	arg := abi.Argument{Type: tp}

	return abi.Arguments{arg}.Pack(big.NewInt(42))
}

func (_m *mockedClient) CodeAt(_ context.Context, _ common.Address, _ *big.Int) ([]byte, error) {
	return []byte{0, 1, 2}, nil
}

const contractABI = "[{\"inputs\":[],\"name\":\"GetValue\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"type\":\"function\",\"name\":\"GetPrimitiveValue\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"pure\"}]"

type simpleTester struct{}

func (s *simpleTester) Setup(t *testing.T) {}

func (s *simpleTester) Name() string { return "" }

func (s *simpleTester) GetAccountBytes(i int) []byte { return []byte{} }

func (s *simpleTester) GetAccountString(i int) string { return "" }

func (s *simpleTester) IsDisabled(testID string) bool { return false }

func (s *simpleTester) DisableTests(testIDs []string) {}

func (s *simpleTester) GetContractReader(t *testing.T) commontypes.ContractReader {
	config := types.ChainReaderConfig{
		Contracts: map[string]types.ChainContractReader{
			"Contract": {
				ContractABI: contractABI,
				Configs: map[string]*types.ChainReaderDefinition{
					"GetValue": {
						ChainSpecificName:   "GetValue",
						OutputModifications: codec.ModifiersConfig{},
					},
					"GetPrimitiveValue": {
						ChainSpecificName: "GetPrimitiveValue",
					},
				},
			},
		},
	}

	svc, err := evm.NewChainReaderService(t.Context(), logger.Nop(), nil, new(simpleHeadTracker), new(mockedClient), config)

	require.NoError(t, err)

	return svc
}

func (s *simpleTester) GetContractWriter(t *testing.T) commontypes.ContractWriter { return nil }

func (s *simpleTester) GetBindings(t *testing.T) []commontypes.BoundContract { return nil }

func (s *simpleTester) DirtyContracts() {}

func (s *simpleTester) MaxWaitTimeForEvents() time.Duration { return time.Second }

func (s *simpleTester) GenerateBlocksTillConfidenceLevel(t *testing.T, contractName, readName string, confidenceLevel primitives.ConfidenceLevel) {
}

type simpleHeadTracker struct {
}

func (h *simpleHeadTracker) Close() error { return nil }

func (h *simpleHeadTracker) HealthReport() map[string]error { return nil }

func (h *simpleHeadTracker) Name() string { return "" }

func (h *simpleHeadTracker) Ready() error { return nil }

func (h *simpleHeadTracker) Start(context.Context) error { return nil }

func (h *simpleHeadTracker) LatestAndFinalizedBlock(ctx context.Context) (latest, finalized *evmtypes.Head, err error) {
	return &evmtypes.Head{}, &evmtypes.Head{}, nil
}
