package evm_test

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	consensustypes "github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3/types"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
)

var (
	reportA = []byte{0x01, 0x02, 0x03}
	reportB = []byte{0xaa, 0xbb, 0xcc, 0xdd}

	// hex encoded 32 byte strings
	workflowID      = "15c631d295ef5e32deb99a10ee6804bc4af1385568f9b3363f6552ac6dbb2cef"
	donID           = "00010203"
	executionID     = "8d4e66421db647dd916d3ec28d56188c8d7dae5f808e03d03339ed2562f13bb0"
	workflowOwnerID = "0000000000000000000000000000000000000000000000000000000000000000"

	invalidID   = "not_valid"
	wrongLength = "8d4e66"
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
		workflowID +
			donID +
			executionID +
			workflowOwnerID +
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

func TestEVMEncoder_InvalidIDs(t *testing.T) {
	config := map[string]any{
		"abi": "mercury_reports bytes[]",
	}
	wrapped, err := values.NewMap(config)
	require.NoError(t, err)
	enc, err := evm.NewEVMEncoder(wrapped)
	require.NoError(t, err)

	// output of a DF2.0 aggregator + metadata fields appended by OCR
	// using an invalid ID
	input := map[string]any{
		"mercury_reports":                   []any{reportA, reportB},
		consensustypes.WorkflowIDFieldName:  invalidID,
		consensustypes.ExecutionIDFieldName: executionID,
	}
	wrapped, err = values.NewMap(input)
	require.NoError(t, err)
	_, err = enc.Encode(testutils.Context(t), *wrapped)
	assert.ErrorContains(t, err, "invalid byte")

	// using valid hex string of wrong length
	input = map[string]any{
		"mercury_reports":                   []any{reportA, reportB},
		consensustypes.WorkflowIDFieldName:  wrongLength,
		consensustypes.ExecutionIDFieldName: executionID,
	}
	wrapped, err = values.NewMap(input)
	require.NoError(t, err)
	_, err = enc.Encode(testutils.Context(t), *wrapped)
	assert.ErrorContains(t, err, "incorrect length for id")
}
