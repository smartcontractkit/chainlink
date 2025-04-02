package operations

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/Masterminds/semver/v3"
	chainsel "github.com/smartcontractkit/chain-selectors"
	mcmslib "github.com/smartcontractkit/mcms"
	"github.com/smartcontractkit/mcms/types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func Test_IsSerializable(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		v    any
		want bool
	}{
		{
			name: "should serialise nil",
			v:    nil,
			want: true,
		},
		{
			name: "should serialise pointer to nil",
			v:    (*int)(nil),
			want: true,
		},
		{
			name: "should serialise string",
			v:    "test",
			want: true,
		},
		{
			name: "should serialise int",
			v:    42,
			want: true,
		},
		{
			name: "should serialise slice",
			v:    []string{"a", "b", "c"},
			want: true,
		},
		{
			name: "should serialise map with string keys",
			v:    map[string]int{"a": 1, "b": 2},
			want: true,
		},
		{
			name: "should serialise map with non-string keys",
			v:    map[int]string{1: "a", 2: "b"},
			want: true,
		},
		{
			name: "should serialise serializableStruct",
			v:    serializableStruct{Name: "test", Age: 30, Numbers: []int{1, 2, 3}},
			want: true,
		},
		{
			name: "should serialise pointer to serializableStruct",
			v:    &serializableStruct{Name: "test", Age: 30},
			want: true,
		},
		{
			name: "should serialise nil pointer",
			v:    (*serializableStruct)(nil),
			want: true,
		},
		{
			name: "should serialise pointer customMarshaler - implements MarshalJSON & UnmarshalJSON",
			v:    &customMarshaler{data: "test"},
			want: true,
		},
		{
			name: "should fail to serialise value customMarshaler - does not implements MarshalJSON & UnmarshalJSON",
			v:    customMarshaler{data: "test"},
			want: false,
		},
		{
			name: "should serialise nested struct",
			v: struct {
				Inner serializableStruct
				Data  string
			}{
				Inner: serializableStruct{Name: "inner", Age: 25},
				Data:  "outer",
			},
			want: true,
		},
		{
			name: "should fail to serialise nested private field",
			v: struct {
				Inner customPrivateStruct
				Data  string
			}{
				Inner: customPrivateStruct{
					Public:  "public",
					private: "private",
				},
				Data: "outer",
			},
			want: false,
		},
		{
			name: "should fail to serialise function",
			v: customFunctionStruct{Action: func() bool {
				return true
			}},
			want: false,
		},
		{
			name: "should serialise function with json tag '-'",
			v: customFunctionStructWithTag{Action: func() bool {
				return true
			}},
			want: true,
		},
		{
			name: "should serialise success report",
			v: NewReport(Definition{
				ID:          "test",
				Version:     semver.MustParse("1.0.0"),
				Description: "test description",
			}, "test", "test", nil, "child-id"),
			want: true,
		},
		{
			name: "should serialise error report",
			v: NewReport(Definition{
				ID:          "test",
				Version:     semver.MustParse("1.0.0"),
				Description: "test description",
			}, "test", "test", errors.New("error report"), "child-id"),
			want: true,
		},
		{
			name: "should serialise TimelockProposal",
			v:    createMCMSTimelockProposal(t),
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			lggr := logger.TestLogger(t)
			if got := IsSerializable(lggr, tt.v); got != tt.want {
				t.Errorf("IsSerializable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func createMCMSTimelockProposal(t *testing.T) *mcmslib.TimelockProposal {
	futureTime := time.Now().Add(time.Hour * 72).Unix()
	builder := mcmslib.NewTimelockProposalBuilder()
	builder.
		SetVersion("v1").
		SetAction(types.TimelockActionSchedule).
		// #nosec G115
		SetValidUntil(uint32(futureTime)).
		SetDescription("mcms description").
		SetDelay(types.NewDuration(1 * time.Hour)).
		SetOverridePreviousRoot(true).
		SetChainMetadata(map[types.ChainSelector]types.ChainMetadata{
			types.ChainSelector(chainsel.ETHEREUM_MAINNET.Selector): {
				StartingOpCount:  1,
				MCMAddress:       "abc",
				AdditionalFields: nil,
			},
		}).
		SetTimelockAddresses(map[types.ChainSelector]string{
			types.ChainSelector(chainsel.ETHEREUM_MAINNET.Selector): "0x123",
		}).
		SetOperations([]types.BatchOperation{
			{
				ChainSelector: types.ChainSelector(chainsel.ETHEREUM_MAINNET.Selector),
				Transactions: []types.Transaction{
					{
						OperationMetadata: types.OperationMetadata{
							ContractType: "test",
						},
						To:               "0x123",
						Data:             []byte{1, 2, 3},
						AdditionalFields: json.RawMessage(`{"test": "test"}`),
					},
				},
			},
		})

	proposal, err := builder.Build()
	require.NoError(t, err)
	return proposal
}

type customFunctionStruct struct {
	Action func() bool
}

type customFunctionStructWithTag struct {
	Action func() bool `json:"-"`
}

type customPrivateStruct struct {
	Public  string
	private string
}

type serializableStruct struct {
	Name    string
	Age     int
	Numbers []int
}

// Implement json.Marshaler and json.Unmarshaler via pointer receiver
type customMarshaler struct {
	data string
}

func (c *customMarshaler) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.data)
}

func (c *customMarshaler) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &c.data)
}
