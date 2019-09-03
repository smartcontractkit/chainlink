package adapters_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/core/adapters"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/stretchr/testify/assert"
)

func TestCondition_Perform(t *testing.T) {
	tests := []struct {
		name       string
		input      interface{}
		adapter    adapters.Condition
		wantResult bool
	}{
		{
			"equals string",
			"inputVal",
			adapters.Condition{
				Operator: "eq",
				Value:    "inputVal",
			},
			true,
		},
		{
			"equals integer as string",
			"1",
			adapters.Condition{
				Operator: "eq",
				Value:    "1",
			},
			true,
		},
		{
			"equals integer as integer",
			1,
			adapters.Condition{
				Operator: "eq",
				Value:    "1",
			},
			true,
		},
		{
			"equals integer as float",
			1.00,
			adapters.Condition{
				Operator: "eq",
				Value:    "1",
			},
			true,
		},
		{
			"equals float as string",
			"1.11",
			adapters.Condition{
				Operator: "eq",
				Value:    "1.11",
			},
			true,
		},
		{
			"equals float as float",
			1.11,
			adapters.Condition{
				Operator: "eq",
				Value:    "1.11",
			},
			true,
		},
		{
			"equals string true",
			"true",
			adapters.Condition{
				Operator: "eq",
				Value:    "true",
			},
			true,
		},
		{
			"equals bool true",
			true,
			adapters.Condition{
				Operator: "eq",
				Value:    "true",
			},
			true,
		},
		{
			"equals string false",
			"false",
			adapters.Condition{
				Operator: "eq",
				Value:    "false",
			},
			true,
		},
		{
			"equals bool false",
			false,
			adapters.Condition{
				Operator: "eq",
				Value:    "false",
			},
			true,
		},
		{
			"greater than integer as string",
			"2",
			adapters.Condition{
				Operator: "gt",
				Value:    "1",
			},
			true,
		},
		{
			"greater than integer as integer",
			2,
			adapters.Condition{
				Operator: "gt",
				Value:    "1",
			},
			true,
		},
		{
			"greater than integer as float",
			2.12,
			adapters.Condition{
				Operator: "gt",
				Value:    "1",
			},
			true,
		},
		{
			"greater than or equals to integer as string 1",
			"2",
			adapters.Condition{
				Operator: "gte",
				Value:    "2",
			},
			true,
		},
		{
			"greater than or equals to integer as integer 1",
			2,
			adapters.Condition{
				Operator: "gte",
				Value:    "2",
			},
			true,
		},
		{
			"greater than or equals to integer as float",
			2.0,
			adapters.Condition{
				Operator: "gte",
				Value:    "2",
			},
			true,
		},
		{
			"greater than or equals to integer as string 2",
			"2",
			adapters.Condition{
				Operator: "gte",
				Value:    "1",
			},
			true,
		},
		{
			"greater than or equals to integer as integer 2",
			2,
			adapters.Condition{
				Operator: "gte",
				Value:    "1",
			},
			true,
		},
		{
			"greater than or equals to integer as float",
			2.0,
			adapters.Condition{
				Operator: "gte",
				Value:    "1",
			},
			true,
		},
		{
			"greater than float as string",
			"2.12",
			adapters.Condition{
				Operator: "gt",
				Value:    "1.11",
			},
			true,
		},
		{
			"greater than float as float",
			2.12,
			adapters.Condition{
				Operator: "gt",
				Value:    "1.11",
			},
			true,
		},
		{
			"greater than float as integer",
			2,
			adapters.Condition{
				Operator: "gt",
				Value:    "1.11",
			},
			true,
		},
		{
			"greater than or equals to float as string 1",
			"2.12",
			adapters.Condition{
				Operator: "gte",
				Value:    "1.11",
			},
			true,
		},
		{
			"greater than or equals to float as float 1",
			2.12,
			adapters.Condition{
				Operator: "gte",
				Value:    "1.11",
			},
			true,
		},
		{
			"greater than or equals to float as string 2",
			"2.12",
			adapters.Condition{
				Operator: "gte",
				Value:    "2.12",
			},
			true,
		},
		{
			"greater than or equals to float as float 2",
			2.12,
			adapters.Condition{
				Operator: "gte",
				Value:    "2.12",
			},
			true,
		},
		{
			"greater than or equals to float as integer",
			2,
			adapters.Condition{
				Operator: "gte",
				Value:    "2.12",
			},
			true,
		},
		{
			"less than integer as string",
			"1",
			adapters.Condition{
				Operator: "lt",
				Value:    "2",
			},
			true,
		},
		{
			"less than integer as integer",
			1,
			adapters.Condition{
				Operator: "lt",
				Value:    "2",
			},
			true,
		},
		{
			"less than integer as float",
			1.0,
			adapters.Condition{
				Operator: "lt",
				Value:    "2",
			},
			true,
		},
		{
			"less than or equals to integer as string 1",
			"1",
			adapters.Condition{
				Operator: "lte",
				Value:    "1",
			},
			true,
		},
		{
			"less than or equals to integer as integer 1",
			1,
			adapters.Condition{
				Operator: "lte",
				Value:    "1",
			},
			true,
		},
		{
			"less than or equals to integer as string 2",
			"1",
			adapters.Condition{
				Operator: "lte",
				Value:    "2",
			},
			true,
		},
		{
			"less than or equals to integer as integer 2",
			1,
			adapters.Condition{
				Operator: "lte",
				Value:    "2",
			},
			true,
		},
		{
			"less than or equals to integer as float",
			1.0,
			adapters.Condition{
				Operator: "lte",
				Value:    "2",
			},
			true,
		},
		{
			"less than float as string",
			"1.11",
			adapters.Condition{
				Operator: "lt",
				Value:    "2.12",
			},
			true,
		},
		{
			"less than float as float",
			1.11,
			adapters.Condition{
				Operator: "lt",
				Value:    "2.12",
			},
			true,
		},
		{
			"less than float as integer",
			1,
			adapters.Condition{
				Operator: "lt",
				Value:    "2.12",
			},
			true,
		},
		{
			"less than or equals to float as string 1",
			"1.11",
			adapters.Condition{
				Operator: "lte",
				Value:    "2.12",
			},
			true,
		},
		{
			"less than or equals to float as float 1",
			1.11,
			adapters.Condition{
				Operator: "lte",
				Value:    "2.12",
			},
			true,
		},
		{
			"less than or equals to float as string 2",
			"1.11",
			adapters.Condition{
				Operator: "lte",
				Value:    "1.11",
			},
			true,
		},
		{
			"less than or equals to float as float 2",
			1.11,
			adapters.Condition{
				Operator: "lte",
				Value:    "1.11",
			},
			true,
		},
		{
			"less than or equals to float as integer",
			1,
			adapters.Condition{
				Operator: "lte",
				Value:    "1.11",
			},
			true,
		},
		{
			"not equals string",
			"inputVal",
			adapters.Condition{
				Operator: "eq",
				Value:    "inputVal2",
			},
			false,
		},
		{
			"not equals integer as string",
			"1",
			adapters.Condition{
				Operator: "eq",
				Value:    "2",
			},
			false,
		},
		{
			"not equals integer as integer",
			1,
			adapters.Condition{
				Operator: "eq",
				Value:    "2",
			},
			false,
		},
		{
			"not equals integer as float",
			1.11,
			adapters.Condition{
				Operator: "eq",
				Value:    "2",
			},
			false,
		},
		{
			"not equals float as float",
			1.11,
			adapters.Condition{
				Operator: "eq",
				Value:    "2.12",
			},
			false,
		},
		{
			"not equals bool as string",
			"false",
			adapters.Condition{
				Operator: "eq",
				Value:    "true",
			},
			false,
		},
		{
			"not equals bool as bool",
			false,
			adapters.Condition{
				Operator: "eq",
				Value:    "true",
			},
			false,
		},
		{
			"not greater than integer as string",
			"1",
			adapters.Condition{
				Operator: "gt",
				Value:    "2",
			},
			false,
		},
		{
			"not greater than integer as integer",
			1,
			adapters.Condition{
				Operator: "gt",
				Value:    "2",
			},
			false,
		},
		{
			"not greater than integer as float",
			1.11,
			adapters.Condition{
				Operator: "gt",
				Value:    "2",
			},
			false,
		},
		{
			"not greater than float as float",
			1.11,
			adapters.Condition{
				Operator: "gt",
				Value:    "2.12",
			},
			false,
		},
		{
			"not greater than float as integer",
			1,
			adapters.Condition{
				Operator: "gt",
				Value:    "2.12",
			},
			false,
		},
		{
			"not greater than or equals to integer as string",
			"1",
			adapters.Condition{
				Operator: "gte",
				Value:    "2",
			},
			false,
		},
		{
			"not greater than or equals to integer as integer",
			1,
			adapters.Condition{
				Operator: "gte",
				Value:    "2",
			},
			false,
		},
		{
			"not greater than or equals to integer as float",
			1.11,
			adapters.Condition{
				Operator: "gte",
				Value:    "2",
			},
			false,
		},
		{
			"not greater than or equals to float as float",
			1.11,
			adapters.Condition{
				Operator: "gte",
				Value:    "2.12",
			},
			false,
		},
		{
			"not greater than or equals to float as integer",
			1,
			adapters.Condition{
				Operator: "gte",
				Value:    "2.12",
			},
			false,
		},
		{
			"not less than integer as string",
			"2",
			adapters.Condition{
				Operator: "lt",
				Value:    "1",
			},
			false,
		},
		{
			"not less than integer as integer",
			2,
			adapters.Condition{
				Operator: "lt",
				Value:    "1",
			},
			false,
		},
		{
			"not less than integer as float",
			2.12,
			adapters.Condition{
				Operator: "lt",
				Value:    "1",
			},
			false,
		},
		{
			"not less than float as float",
			2.12,
			adapters.Condition{
				Operator: "lt",
				Value:    "1.11",
			},
			false,
		},
		{
			"not less than float as integer",
			2,
			adapters.Condition{
				Operator: "lt",
				Value:    "1.11",
			},
			false,
		},
		{
			"not less than or equals to integer as string",
			"2",
			adapters.Condition{
				Operator: "lte",
				Value:    "1",
			},
			false,
		},
		{
			"not less than or equals to integer as integer",
			2,
			adapters.Condition{
				Operator: "lte",
				Value:    "1",
			},
			false,
		},
		{
			"not less than or equals to integer as float",
			2.12,
			adapters.Condition{
				Operator: "lte",
				Value:    "1",
			},
			false,
		},
		{
			"not less than or equals to float as float",
			2.12,
			adapters.Condition{
				Operator: "lte",
				Value:    "1.11",
			},
			false,
		},
		{
			"not less than or equals to float as integer",
			2,
			adapters.Condition{
				Operator: "lte",
				Value:    "1.11",
			},
			false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			input := cltest.RunResultWithResult(test.input)
			adapter := test.adapter
			result := adapter.Perform(input, nil)
			val := result.Result()
			assert.NoError(t, result.GetError())
			assert.Equal(t, test.wantResult, val.Bool())
		})
	}
}

func TestConditionError_Perform(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		adapter  adapters.Condition
		expected error
	}{
		{
			"greater than not number in result",
			"a",
			adapters.Condition{
				Operator: "gt",
				Value:    "2",
			},
			adapters.ErrResultNotNumber,
		},
		{
			"greater than not number in desired",
			"2",
			adapters.Condition{
				Operator: "gt",
				Value:    "a",
			},
			adapters.ErrValueNotNumber,
		},
		{
			"less than not number in result",
			"a",
			adapters.Condition{
				Operator: "lt",
				Value:    "2",
			},
			adapters.ErrResultNotNumber,
		},
		{
			"less than not number in desired",
			"2",
			adapters.Condition{
				Operator: "lt",
				Value:    "a",
			},
			adapters.ErrValueNotNumber,
		},
		{
			"missing operator",
			"2",
			adapters.Condition{
				Operator: "",
				Value:    "3",
			},
			adapters.ErrOperatorNotSpecified,
		},
		{
			"missing desired",
			"2",
			adapters.Condition{
				Operator: "eq",
				Value:    "",
			},
			adapters.ErrValueNotSpecified,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			input := cltest.RunResultWithResult(test.input)
			adapter := test.adapter
			result := adapter.Perform(input, nil)
			_, err := result.ResultString()
			assert.NoError(t, err)
			assert.Equal(t, test.expected, result.GetError())
		})
	}
}
