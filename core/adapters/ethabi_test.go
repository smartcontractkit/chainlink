package adapters

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestEVMTranscodeJSONWithFormat(t *testing.T) {
	reallyLongHexString := "0x" + strings.Repeat("0123456789abcdef", 100)
	tests := []struct {
		name   string
		format string
		input  string
		output string
	}{
		{
			"result is string",
			FormatBytes,
			`{"result": "hello world"}`,
			"0x" +
				"000000000000000000000000000000000000000000000000000000000000000b" +
				"68656c6c6f20776f726c64000000000000000000000000000000000000000000",
		},
		{
			"result is number",
			FormatUint256,
			`{"result": 31223}`,
			"0x" +
				"0000000000000000000000000000000000000000000000000000000000000020" +
				"00000000000000000000000000000000000000000000000000000000000079f7",
		},
		{
			"result is negative number",
			FormatInt256,
			`{"result": -123481273.1}`,
			"0x" +
				"0000000000000000000000000000000000000000000000000000000000000020" +
				"fffffffffffffffffffffffffffffffffffffffffffffffffffffffff8a3d347",
		},
		{
			"result is true",
			FormatBool,
			`{"result": true}`,
			"0x" +
				"0000000000000000000000000000000000000000000000000000000000000020" +
				"0000000000000000000000000000000000000000000000000000000000000001",
		},
		{
			"result is preformatted",
			FormatPreformattedHexArguments,
			fmt.Sprintf(`{"result": "%s"}`, reallyLongHexString),
			reallyLongHexString,
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			input := gjson.GetBytes([]byte(test.input), "result")
			out, err := EVMTranscodeJSONWithFormat(input, test.format)
			require.NoError(t, err)
			assert.Equal(t, test.output, hexutil.Encode(out))
		})
	}
}

func TestEVMTranscodeJSONWithFormat_UnsupportedEncoding(t *testing.T) {
	_, err := EVMTranscodeJSONWithFormat(gjson.Result{}, "burgh")
	assert.Error(t, err)
}
