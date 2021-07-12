package pipeline_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	ethmocks "github.com/smartcontractkit/chainlink/core/services/eth/mocks"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
)

func TestETHCallTask(t *testing.T) {
	tests := []struct {
		name                  string
		contract              string
		data                  string
		vars                  pipeline.Vars
		inputs                []pipeline.Result
		setupClientMock       func(ethClient *ethmocks.Client)
		expected              interface{}
		expectedErrorCause    error
		expectedErrorContains string
	}{
		{
			"happy",
			"0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbeeF",
			"$(foo)",
			pipeline.NewVarsFrom(map[string]interface{}{
				"foo": []byte("foo bar"),
			}),
			nil,
			func(ethClient *ethmocks.Client) {
				contractAddr := common.HexToAddress("0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbeeF")
				ethClient.
					On("CallContract", mock.Anything, ethereum.CallMsg{To: &contractAddr, Data: []byte("foo bar")}, (*big.Int)(nil)).
					Return([]byte("baz quux"), nil)
			},
			[]byte("baz quux"), nil, "",
		},
		{
			"bad contract address",
			"0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbee",
			"$(foo)",
			pipeline.NewVarsFrom(map[string]interface{}{
				"foo": []byte("foo bar"),
			}),
			nil,
			func(ethClient *ethmocks.Client) {},
			nil, pipeline.ErrBadInput, "contract",
		},
		{
			"missing data var",
			"0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbeeF",
			"$(foo)",
			pipeline.NewVarsFrom(map[string]interface{}{
				"zork": []byte("foo bar"),
			}),
			nil,
			func(ethClient *ethmocks.Client) {},
			nil, pipeline.ErrKeypathNotFound, "data",
		},
		{
			"no data",
			"0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbeeF",
			"$(foo)",
			pipeline.NewVarsFrom(map[string]interface{}{
				"foo": []byte(nil),
			}),
			nil,
			func(ethClient *ethmocks.Client) {},
			nil, pipeline.ErrBadInput, "data",
		},
		{
			"errored input",
			"0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbeeF",
			"$(foo)",
			pipeline.NewVarsFrom(map[string]interface{}{
				"foo": []byte("foo bar"),
			}),
			[]pipeline.Result{{Error: errors.New("uh oh")}},
			func(ethClient *ethmocks.Client) {},
			nil, pipeline.ErrTooManyErrors, "task inputs",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			task := pipeline.ETHCallTask{
				BaseTask: pipeline.NewBaseTask(0, "ethcall", nil, nil, 0),
				Contract: test.contract,
				Data:     test.data,
			}

			ethClient := new(ethmocks.Client)
			test.setupClientMock(ethClient)
			task.HelperSetDependencies(ethClient)

			result := task.Run(context.Background(), test.vars, test.inputs)

			if test.expectedErrorCause != nil {
				require.Equal(t, test.expectedErrorCause, errors.Cause(result.Error))
				require.Nil(t, result.Value)
				if test.expectedErrorContains != "" {
					require.Contains(t, result.Error.Error(), test.expectedErrorContains)
				}
			} else {
				require.NoError(t, result.Error)
				require.Equal(t, test.expected, result.Value)
			}
		})
	}
}
