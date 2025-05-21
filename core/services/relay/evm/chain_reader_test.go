package evm_test

import (
	"context"
	"fmt"
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/codec"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/testutils"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	evmtypes "github.com/smartcontractkit/chainlink-evm/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

func TestChainReaderSizedBigIntTypes(t *testing.T) {
	t.Parallel()

	tests := []string{}

	// 8, 16, 32, and 64 bits have their own type in go that is used by abi.
	for i := 24; i <= 256; i += 8 {
		if i == 32 || i == 64 {
			continue
		}

		tp := fmt.Sprintf("int%d", i)
		tests = append(tests, tp, "u"+tp)
	}

	for _, test := range tests {
		t.Run(test, func(t *testing.T) {
			t.Parallel()

			tester := &simpleTester{returnVal: big.NewInt(42), internalType: test}
			wrapped := testutils.WrapContractReaderTesterForLoop(tester)
			wrapped.Setup(t)

			svc := wrapped.GetContractReader(t)
			binding := commontypes.BoundContract{Address: "0x21", Name: "Contract"}

			require.NoError(t, svc.Bind(t.Context(), []commontypes.BoundContract{binding}))

			var value values.Value
			require.NoError(t, svc.GetLatestValue(t.Context(), binding.ReadIdentifier("GetValue"), primitives.Finalized, nil, &value))

			out := new(big.Int)
			require.NoError(t, value.UnwrapTo(out))

			assert.Equal(t, int64(42), out.Int64())
		})
	}
}

func TestChainReaderPrimitiveTypes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		abiType  string
		expected any
	}{
		{"int8", int8(42)},
		{"int16", int16(42)},
		{"int32", int32(42)},
		{"int64", int64(42)},
		{"string", "42"},
	}

	for _, test := range tests {
		t.Run(test.abiType, func(t *testing.T) {
			t.Parallel()

			tester := &simpleTester{returnVal: test.expected, internalType: test.abiType}
			wrapped := testutils.WrapContractReaderTesterForLoop(tester)
			wrapped.Setup(t)

			svc := wrapped.GetContractReader(t)
			binding := commontypes.BoundContract{Address: "0x21", Name: "Contract"}

			require.NoError(t, svc.Bind(t.Context(), []commontypes.BoundContract{binding}))

			var value values.Value
			require.NoError(t, svc.GetLatestValue(t.Context(), binding.ReadIdentifier("GetValue"), primitives.Finalized, nil, &value))

			out := reflect.New(reflect.TypeOf(test.expected)).Interface()
			require.NoError(t, value.UnwrapTo(out))

			assert.Equal(t, test.expected, reflect.Indirect(reflect.ValueOf(out)).Interface())
		})
	}
}

type mockedClient struct {
	value        any
	internalType abi.Type
}

func newMockedClient(t *testing.T, value any, internalType string) *mockedClient {
	t.Helper()

	internal, err := abi.NewType(internalType, "", nil)

	require.NoError(t, err)

	return &mockedClient{
		value:        value,
		internalType: internal,
	}
}

func (_m *mockedClient) BatchCallContext(_ context.Context, _ []rpc.BatchElem) error { return nil }

func (_m *mockedClient) CallContract(_ context.Context, msg ethereum.CallMsg, _ *big.Int) ([]byte, error) {
	return abi.Arguments{abi.Argument{Type: _m.internalType}}.Pack(_m.value)
}

func (_m *mockedClient) CodeAt(_ context.Context, _ common.Address, _ *big.Int) ([]byte, error) {
	return []byte{0, 1, 2}, nil
}

const contractABI = `[{"inputs":[],"name":"GetValue","outputs":[{"internalType":"%s","name":"","type":"%s"}],"stateMutability":"pure","type":"function"}]`

type simpleTester struct {
	returnVal    any
	internalType string
}

func (s *simpleTester) Setup(t *testing.T) {}

func (s *simpleTester) Name() string { return "" }

func (s *simpleTester) GetAccountBytes(i int) []byte { return []byte{} }

func (s *simpleTester) GetAccountString(i int) string { return "" }

func (s *simpleTester) IsDisabled(testID string) bool { return false }

func (s *simpleTester) DisableTests(testIDs []string) {}

func (s *simpleTester) GetContractReader(t *testing.T) commontypes.ContractReader {
	t.Helper()

	config := types.ChainReaderConfig{
		Contracts: map[string]types.ChainContractReader{
			"Contract": {
				ContractABI: fmt.Sprintf(contractABI, s.internalType, s.internalType),
				Configs: map[string]*types.ChainReaderDefinition{
					"GetValue": {
						ChainSpecificName:   "GetValue",
						OutputModifications: codec.ModifiersConfig{},
					},
				},
			},
		},
	}

	client := newMockedClient(t, s.returnVal, s.internalType)
	svc, err := evm.NewChainReaderService(t.Context(), logger.Nop(), nil, new(simpleHeadTracker), client, config)

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
