package adapters_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/core/adapters"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"

	"github.com/stretchr/testify/assert"
)

func TestCompare_Perform(t *testing.T) {
	tests := []struct {
		name       string
		input      interface{}
		adapter    adapters.Compare
		wantResult bool
	}{
		{
			"equals string",
			"inputVal",
			adapters.Compare{
				Operator: "eq",
				Value:    "inputVal",
			},
			true,
		},
		{
			"equals integer as string",
			"1",
			adapters.Compare{
				Operator: "eq",
				Value:    "1",
			},
			true,
		},
		{
			"equals integer as integer",
			1,
			adapters.Compare{
				Operator: "eq",
				Value:    "1",
			},
			true,
		},
		{
			"equals integer as float",
			1.00,
			adapters.Compare{
				Operator: "eq",
				Value:    "1",
			},
			true,
		},
		{
			"equals float as string",
			"1.11",
			adapters.Compare{
				Operator: "eq",
				Value:    "1.11",
			},
			true,
		},
		{
			"equals float as float",
			1.11,
			adapters.Compare{
				Operator: "eq",
				Value:    "1.11",
			},
			true,
		},
		{
			"equals string true",
			"true",
			adapters.Compare{
				Operator: "eq",
				Value:    "true",
			},
			true,
		},
		{
			"equals bool true",
			true,
			adapters.Compare{
				Operator: "eq",
				Value:    "true",
			},
			true,
		},
		{
			"equals string false",
			"false",
			adapters.Compare{
				Operator: "eq",
				Value:    "false",
			},
			true,
		},
		{
			"equals bool false",
			false,
			adapters.Compare{
				Operator: "eq",
				Value:    "false",
			},
			true,
		},
		{
			"notequals string",
			"somethingElse",
			adapters.Compare{
				Operator: "neq",
				Value:    "inputVal",
			},
			true,
		},
		{
			"notequals integer as string",
			"2",
			adapters.Compare{
				Operator: "neq",
				Value:    "1",
			},
			true,
		},
		{
			"notequals integer as integer",
			2,
			adapters.Compare{
				Operator: "neq",
				Value:    "1",
			},
			true,
		},
		{
			"notequals integer as float",
			2.00,
			adapters.Compare{
				Operator: "neq",
				Value:    "1",
			},
			true,
		},
		{
			"notequals float as string",
			"2.12",
			adapters.Compare{
				Operator: "neq",
				Value:    "1.11",
			},
			true,
		},
		{
			"notequals float as float",
			2.12,
			adapters.Compare{
				Operator: "neq",
				Value:    "1.11",
			},
			true,
		},
		{
			"notequals string true",
			"true",
			adapters.Compare{
				Operator: "neq",
				Value:    "false",
			},
			true,
		},
		{
			"notequals bool true",
			true,
			adapters.Compare{
				Operator: "neq",
				Value:    "false",
			},
			true,
		},
		{
			"notequals string false",
			"false",
			adapters.Compare{
				Operator: "neq",
				Value:    "true",
			},
			true,
		},
		{
			"notequals bool false",
			false,
			adapters.Compare{
				Operator: "neq",
				Value:    "true",
			},
			true,
		},
		{
			"greater than integer as string",
			"2",
			adapters.Compare{
				Operator: "gt",
				Value:    "1",
			},
			true,
		},
		{
			"greater than integer as integer",
			2,
			adapters.Compare{
				Operator: "gt",
				Value:    "1",
			},
			true,
		},
		{
			"greater than integer as float",
			2.12,
			adapters.Compare{
				Operator: "gt",
				Value:    "1",
			},
			true,
		},
		{
			"greater than or equals to integer as string 1",
			"2",
			adapters.Compare{
				Operator: "gte",
				Value:    "2",
			},
			true,
		},
		{
			"greater than or equals to integer as integer 1",
			2,
			adapters.Compare{
				Operator: "gte",
				Value:    "2",
			},
			true,
		},
		{
			"greater than or equals to integer as float",
			2.0,
			adapters.Compare{
				Operator: "gte",
				Value:    "2",
			},
			true,
		},
		{
			"greater than or equals to integer as string 2",
			"2",
			adapters.Compare{
				Operator: "gte",
				Value:    "1",
			},
			true,
		},
		{
			"greater than or equals to integer as integer 2",
			2,
			adapters.Compare{
				Operator: "gte",
				Value:    "1",
			},
			true,
		},
		{
			"greater than or equals to integer as float",
			2.0,
			adapters.Compare{
				Operator: "gte",
				Value:    "1",
			},
			true,
		},
		{
			"greater than float as string",
			"2.12",
			adapters.Compare{
				Operator: "gt",
				Value:    "1.11",
			},
			true,
		},
		{
			"greater than float as float",
			2.12,
			adapters.Compare{
				Operator: "gt",
				Value:    "1.11",
			},
			true,
		},
		{
			"greater than float as integer",
			2,
			adapters.Compare{
				Operator: "gt",
				Value:    "1.11",
			},
			true,
		},
		{
			"greater than or equals to float as string 1",
			"2.12",
			adapters.Compare{
				Operator: "gte",
				Value:    "1.11",
			},
			true,
		},
		{
			"greater than or equals to float as float 1",
			2.12,
			adapters.Compare{
				Operator: "gte",
				Value:    "1.11",
			},
			true,
		},
		{
			"greater than or equals to float as string 2",
			"2.12",
			adapters.Compare{
				Operator: "gte",
				Value:    "2.12",
			},
			true,
		},
		{
			"greater than or equals to float as float 2",
			2.12,
			adapters.Compare{
				Operator: "gte",
				Value:    "2.12",
			},
			true,
		},
		{
			"greater than or equals to float as integer",
			2,
			adapters.Compare{
				Operator: "gte",
				Value:    "2.12",
			},
			false,
		},
		{
			"less than integer as string",
			"1",
			adapters.Compare{
				Operator: "lt",
				Value:    "2",
			},
			true,
		},
		{
			"less than integer as integer",
			1,
			adapters.Compare{
				Operator: "lt",
				Value:    "2",
			},
			true,
		},
		{
			"less than integer as float",
			1.0,
			adapters.Compare{
				Operator: "lt",
				Value:    "2",
			},
			true,
		},
		{
			"less than or equals to integer as string 1",
			"1",
			adapters.Compare{
				Operator: "lte",
				Value:    "1",
			},
			true,
		},
		{
			"less than or equals to integer as integer 1",
			1,
			adapters.Compare{
				Operator: "lte",
				Value:    "1",
			},
			true,
		},
		{
			"less than or equals to integer as string 2",
			"1",
			adapters.Compare{
				Operator: "lte",
				Value:    "2",
			},
			true,
		},
		{
			"less than or equals to integer as integer 2",
			1,
			adapters.Compare{
				Operator: "lte",
				Value:    "2",
			},
			true,
		},
		{
			"less than or equals to integer as float",
			1.0,
			adapters.Compare{
				Operator: "lte",
				Value:    "2",
			},
			true,
		},
		{
			"less than float as string",
			"1.11",
			adapters.Compare{
				Operator: "lt",
				Value:    "2.12",
			},
			true,
		},
		{
			"less than float as float",
			1.11,
			adapters.Compare{
				Operator: "lt",
				Value:    "2.12",
			},
			true,
		},
		{
			"less than float as integer",
			1,
			adapters.Compare{
				Operator: "lt",
				Value:    "2.12",
			},
			true,
		},
		{
			"less than or equals to float as string 1",
			"1.11",
			adapters.Compare{
				Operator: "lte",
				Value:    "2.12",
			},
			true,
		},
		{
			"less than or equals to float as float 1",
			1.11,
			adapters.Compare{
				Operator: "lte",
				Value:    "2.12",
			},
			true,
		},
		{
			"less than or equals to float as string 2",
			"1.11",
			adapters.Compare{
				Operator: "lte",
				Value:    "1.11",
			},
			true,
		},
		{
			"less than or equals to float as float 2",
			1.11,
			adapters.Compare{
				Operator: "lte",
				Value:    "1.11",
			},
			true,
		},
		{
			"less than or equals to float as integer",
			1,
			adapters.Compare{
				Operator: "lte",
				Value:    "1.11",
			},
			true,
		},
		{
			"not equals string",
			"inputVal",
			adapters.Compare{
				Operator: "eq",
				Value:    "inputVal2",
			},
			false,
		},
		{
			"not equals integer as string",
			"1",
			adapters.Compare{
				Operator: "eq",
				Value:    "2",
			},
			false,
		},
		{
			"not equals integer as integer",
			1,
			adapters.Compare{
				Operator: "eq",
				Value:    "2",
			},
			false,
		},
		{
			"not equals integer as float",
			1.11,
			adapters.Compare{
				Operator: "eq",
				Value:    "2",
			},
			false,
		},
		{
			"not equals float as float",
			1.11,
			adapters.Compare{
				Operator: "eq",
				Value:    "2.12",
			},
			false,
		},
		{
			"not equals bool as string",
			"false",
			adapters.Compare{
				Operator: "eq",
				Value:    "true",
			},
			false,
		},
		{
			"not equals bool as bool",
			false,
			adapters.Compare{
				Operator: "eq",
				Value:    "true",
			},
			false,
		},
		{
			"not notequals string",
			"inputVal",
			adapters.Compare{
				Operator: "neq",
				Value:    "inputVal",
			},
			false,
		},
		{
			"not notequals integer as string",
			"2",
			adapters.Compare{
				Operator: "neq",
				Value:    "2",
			},
			false,
		},
		{
			"not notequals integer as integer",
			2,
			adapters.Compare{
				Operator: "neq",
				Value:    "2",
			},
			false,
		},
		{
			"not notequals integer as float",
			1.00,
			adapters.Compare{
				Operator: "neq",
				Value:    "1",
			},
			false,
		},
		{
			"not notequals float as float",
			1.11,
			adapters.Compare{
				Operator: "neq",
				Value:    "1.11",
			},
			false,
		},
		{
			"not notequals bool as string",
			"false",
			adapters.Compare{
				Operator: "neq",
				Value:    "false",
			},
			false,
		},
		{
			"not notequals bool as bool",
			false,
			adapters.Compare{
				Operator: "neq",
				Value:    "false",
			},
			false,
		},
		{
			"not greater than integer as string",
			"1",
			adapters.Compare{
				Operator: "gt",
				Value:    "2",
			},
			false,
		},
		{
			"not greater than integer as integer",
			1,
			adapters.Compare{
				Operator: "gt",
				Value:    "2",
			},
			false,
		},
		{
			"not greater than integer as float",
			1.11,
			adapters.Compare{
				Operator: "gt",
				Value:    "2",
			},
			false,
		},
		{
			"not greater than float as float",
			1.11,
			adapters.Compare{
				Operator: "gt",
				Value:    "2.12",
			},
			false,
		},
		{
			"not greater than float as integer",
			1,
			adapters.Compare{
				Operator: "gt",
				Value:    "2.12",
			},
			false,
		},
		{
			"not greater than or equals to integer as string",
			"1",
			adapters.Compare{
				Operator: "gte",
				Value:    "2",
			},
			false,
		},
		{
			"not greater than or equals to integer as integer",
			1,
			adapters.Compare{
				Operator: "gte",
				Value:    "2",
			},
			false,
		},
		{
			"not greater than or equals to integer as float",
			1.11,
			adapters.Compare{
				Operator: "gte",
				Value:    "2",
			},
			false,
		},
		{
			"not greater than or equals to float as float",
			1.11,
			adapters.Compare{
				Operator: "gte",
				Value:    "2.12",
			},
			false,
		},
		{
			"not greater than or equals to float as integer",
			1,
			adapters.Compare{
				Operator: "gte",
				Value:    "2.12",
			},
			false,
		},
		{
			"not less than integer as string",
			"2",
			adapters.Compare{
				Operator: "lt",
				Value:    "1",
			},
			false,
		},
		{
			"not less than integer as integer",
			2,
			adapters.Compare{
				Operator: "lt",
				Value:    "1",
			},
			false,
		},
		{
			"not less than integer as float",
			2.12,
			adapters.Compare{
				Operator: "lt",
				Value:    "1",
			},
			false,
		},
		{
			"not less than float as float",
			2.12,
			adapters.Compare{
				Operator: "lt",
				Value:    "1.11",
			},
			false,
		},
		{
			"not less than float as integer",
			2,
			adapters.Compare{
				Operator: "lt",
				Value:    "1.11",
			},
			false,
		},
		{
			"not less than or equals to integer as string",
			"2",
			adapters.Compare{
				Operator: "lte",
				Value:    "1",
			},
			false,
		},
		{
			"not less than or equals to integer as integer",
			2,
			adapters.Compare{
				Operator: "lte",
				Value:    "1",
			},
			false,
		},
		{
			"not less than or equals to integer as float",
			2.12,
			adapters.Compare{
				Operator: "lte",
				Value:    "1",
			},
			false,
		},
		{
			"not less than or equals to float as float",
			2.12,
			adapters.Compare{
				Operator: "lte",
				Value:    "1.11",
			},
			false,
		},
		{
			"not less than or equals to float as integer",
			2,
			adapters.Compare{
				Operator: "lte",
				Value:    "1.11",
			},
			false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			input := cltest.NewRunInputWithResult(test.input)
			adapter := test.adapter
			result := adapter.Perform(input, nil, nil)
			val := result.Result()
			assert.NoError(t, result.Error())
			assert.Equal(t, test.wantResult, val.Bool())
		})
	}
}

func TestCompareError_Perform(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		adapter  adapters.Compare
		expected error
	}{
		{
			"greater than not number in result",
			"a",
			adapters.Compare{
				Operator: "gt",
				Value:    "2",
			},
			adapters.ErrResultNotNumber,
		},
		{
			"greater than not number in desired",
			"2",
			adapters.Compare{
				Operator: "gt",
				Value:    "a",
			},
			adapters.ErrValueNotNumber,
		},
		{
			"less than not number in result",
			"a",
			adapters.Compare{
				Operator: "lt",
				Value:    "2",
			},
			adapters.ErrResultNotNumber,
		},
		{
			"less than not number in desired",
			"2",
			adapters.Compare{
				Operator: "lt",
				Value:    "a",
			},
			adapters.ErrValueNotNumber,
		},
		{
			"missing operator",
			"2",
			adapters.Compare{
				Operator: "",
				Value:    "3",
			},
			adapters.ErrOperatorNotSpecified,
		},
		{
			"missing desired",
			"2",
			adapters.Compare{
				Operator: "eq",
				Value:    "",
			},
			adapters.ErrValueNotSpecified,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			input := cltest.NewRunInputWithResult(test.input)
			adapter := test.adapter
			result := adapter.Perform(input, nil, nil)
			assert.Equal(t, test.expected, result.Error())
		})
	}
}
