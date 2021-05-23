package pipeline_test

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/services/pipeline"
)

func TestETHABIDecodeTask(t *testing.T) {
	task := pipeline.ETHABIDecodeTask{
		BaseTask: pipeline.NewBaseTask("decode", nil, 0, 0),
		ABI:      "uint256 u256, bool boolean, int256 i256, string s",
	}

	expected := map[string]interface{}{
		"u256":    big.NewInt(123),
		"boolean": true,
		"i256":    big.NewInt(-321),
		"s":       "foo bar baz",
	}

	inputs := []pipeline.Result{
		{Value: "0x000000000000000000000000000000000000000000000000000000000000007b0000000000000000000000000000000000000000000000000000000000000001fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffebf0000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000000000000b666f6f206261722062617a000000000000000000000000000000000000000000"},
	}

	result := task.Run(context.Background(), pipeline.NewVars(), pipeline.JSONSerializable{}, inputs)
	bs, _ := json.MarshalIndent(result.Value, "", "    ")
	fmt.Println(string(bs))
	fmt.Printf("%+v\n", result.Error)

	require.Equal(t, expected, result.Value)
}
