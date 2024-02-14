package evm_test

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"

	consensustypes "github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3/types"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
)

var (
	reportA     = []byte{0x01, 0x02, 0x03}
	reportB     = []byte{0xaa, 0xbb, 0xcc, 0xdd}
	workflowID  = "my_id"
	executionID = "my_execution_id"
)

func TestEVMEncoder(t *testing.T) {
	config := map[string]any{
		"abi": "mercury_reports bytes[]",
	}
	wrapped, err := values.NewMap(config)
	require.NoError(t, err)
	enc, err := evm.NewEVMEncoder(wrapped)
	require.NoError(t, err)

	// output of a DF2.0 aggregator + metadata fields appended by OCR
	input := map[string]any{
		"mercury_reports":                   []any{reportA, reportB},
		consensustypes.WorkflowIDFieldName:  workflowID,
		consensustypes.ExecutionIDFieldName: executionID,
	}
	wrapped, err = values.NewMap(input)
	require.NoError(t, err)
	encoded, err := enc.Encode(testutils.Context(t), *wrapped)
	require.NoError(t, err)

	expected :=
		// start of the outer tuple ((user_fields), workflow_id, workflow_execution_id)
		"6d795f6964000000000000000000000000000000000000000000000000000000" + // workflow ID
			"6d795f657865637574696f6e5f69640000000000000000000000000000000000" + // execution ID
			// start of the inner tuple (user_fields)
			"0000000000000000000000000000000000000000000000000000000000000020" + // offset of mercury_reports array
			"0000000000000000000000000000000000000000000000000000000000000002" + // length of mercury_reports array
			"0000000000000000000000000000000000000000000000000000000000000040" + // offset of reportA
			"0000000000000000000000000000000000000000000000000000000000000080" + // offset of reportB
			"0000000000000000000000000000000000000000000000000000000000000003" + // length of reportA
			"0102030000000000000000000000000000000000000000000000000000000000" + // reportA
			"0000000000000000000000000000000000000000000000000000000000000004" + // length of reportB
			"aabbccdd00000000000000000000000000000000000000000000000000000000" // reportB
	// end of the inner tuple (user_fields)

	require.Equal(t, expected, hex.EncodeToString(encoded))
}
