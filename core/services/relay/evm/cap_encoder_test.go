package evm_test

import (
	"encoding/hex"
	"math/big"
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

	workflowID       = "15c631d295ef5e32deb99a10ee6804bc4af1385568f9b3363f6552ac6dbb2cef"
	workflowName     = "aabbccddeeaabbccddee"
	donID            = "00010203"
	executionID      = "8d4e66421db647dd916d3ec28d56188c8d7dae5f808e03d03339ed2562f13bb0"
	workflowOwnerID  = "0000000000000000000000000000000000000000"
	reportID         = "9988"
	timestampInt     = uint32(1234567890)
	timestampHex     = "499602d2"
	configVersionInt = uint32(1)
	configVersionHex = "00000001"

	invalidID   = "not_valid"
	wrongLength = "8d4e66"
)

func TestEVMEncoder_SingleField(t *testing.T) {
	config := map[string]any{
		"abi": "bytes[] Full_reports",
	}
	wrapped, err := values.NewMap(config)
	require.NoError(t, err)
	enc, err := evm.NewEVMEncoder(wrapped)
	require.NoError(t, err)

	// output of a DF2.0 aggregator + metadata fields appended by OCR
	input := map[string]any{
		"Full_reports":                   []any{reportA, reportB},
		consensustypes.MetadataFieldName: getMetadata(workflowID),
	}
	wrapped, err = values.NewMap(input)
	require.NoError(t, err)
	encoded, err := enc.Encode(testutils.Context(t), *wrapped)
	require.NoError(t, err)

	expected :=
		// start of the outer tuple
		getHexMetadata() +
			// start of the inner tuple (user_fields)
			"0000000000000000000000000000000000000000000000000000000000000020" + // offset of Full_reports array
			"0000000000000000000000000000000000000000000000000000000000000002" + // length of Full_reports array
			"0000000000000000000000000000000000000000000000000000000000000040" + // offset of reportA
			"0000000000000000000000000000000000000000000000000000000000000080" + // offset of reportB
			"0000000000000000000000000000000000000000000000000000000000000003" + // length of reportA
			"0102030000000000000000000000000000000000000000000000000000000000" + // reportA
			"0000000000000000000000000000000000000000000000000000000000000004" + // length of reportB
			"aabbccdd00000000000000000000000000000000000000000000000000000000" // reportB

	require.Equal(t, expected, hex.EncodeToString(encoded))
}

func TestEVMEncoder_TwoFields(t *testing.T) {
	config := map[string]any{
		"abi": "uint256[] Prices, uint32[] Timestamps",
	}
	wrapped, err := values.NewMap(config)
	require.NoError(t, err)
	enc, err := evm.NewEVMEncoder(wrapped)
	require.NoError(t, err)

	// output of a DF2.0 aggregator + metadata fields appended by OCR
	input := map[string]any{
		"Prices":                         []any{big.NewInt(234), big.NewInt(456)},
		"Timestamps":                     []any{int64(111), int64(222)},
		consensustypes.MetadataFieldName: getMetadata(workflowID),
	}
	wrapped, err = values.NewMap(input)
	require.NoError(t, err)
	encoded, err := enc.Encode(testutils.Context(t), *wrapped)
	require.NoError(t, err)

	expected :=
		// start of the outer tuple
		getHexMetadata() +
			// start of the inner tuple (user_fields)
			"0000000000000000000000000000000000000000000000000000000000000040" + // offset of Prices array
			"00000000000000000000000000000000000000000000000000000000000000a0" + // offset of Timestamps array
			"0000000000000000000000000000000000000000000000000000000000000002" + // length of Prices array
			"00000000000000000000000000000000000000000000000000000000000000ea" + // Prices[0]
			"00000000000000000000000000000000000000000000000000000000000001c8" + // Prices[1]
			"0000000000000000000000000000000000000000000000000000000000000002" + // length of Timestamps array
			"000000000000000000000000000000000000000000000000000000000000006f" + // Timestamps[0]
			"00000000000000000000000000000000000000000000000000000000000000de" // Timestamps[1]

	require.Equal(t, expected, hex.EncodeToString(encoded))
}

func TestEVMEncoder_Tuple(t *testing.T) {
	config := map[string]any{
		"abi": "(uint256[] Prices, uint32[] Timestamps) Elem",
	}
	wrapped, err := values.NewMap(config)
	require.NoError(t, err)
	enc, err := evm.NewEVMEncoder(wrapped)
	require.NoError(t, err)

	// output of a DF2.0 aggregator + metadata fields appended by OCR
	input := map[string]any{
		"Elem": map[string]any{
			"Prices":     []any{big.NewInt(234), big.NewInt(456)},
			"Timestamps": []any{int64(111), int64(222)},
		},
		consensustypes.MetadataFieldName: getMetadata(workflowID),
	}
	wrapped, err = values.NewMap(input)
	require.NoError(t, err)
	encoded, err := enc.Encode(testutils.Context(t), *wrapped)
	require.NoError(t, err)

	expected :=
		// start of the outer tuple
		getHexMetadata() +
			// start of the inner tuple (user_fields)
			"0000000000000000000000000000000000000000000000000000000000000020" + // offset of Elem tuple
			"0000000000000000000000000000000000000000000000000000000000000040" + // offset of Prices array
			"00000000000000000000000000000000000000000000000000000000000000a0" + // offset of Timestamps array
			"0000000000000000000000000000000000000000000000000000000000000002" + // length of Prices array
			"00000000000000000000000000000000000000000000000000000000000000ea" + // Prices[0] = 234
			"00000000000000000000000000000000000000000000000000000000000001c8" + // Prices[1] = 456
			"0000000000000000000000000000000000000000000000000000000000000002" + // length of Timestamps array
			"000000000000000000000000000000000000000000000000000000000000006f" + // Timestamps[0] = 111
			"00000000000000000000000000000000000000000000000000000000000000de" // Timestamps[1] = 222

	require.Equal(t, expected, hex.EncodeToString(encoded))
}

func TestEVMEncoder_ListOfTuples(t *testing.T) {
	config := map[string]any{
		"abi": "(uint256 Price, uint32 Timestamp)[] Elems",
	}
	wrapped, err := values.NewMap(config)
	require.NoError(t, err)
	enc, err := evm.NewEVMEncoder(wrapped)
	require.NoError(t, err)

	// output of a DF2.0 aggregator + metadata fields appended by OCR
	input := map[string]any{
		"Elems": []any{
			map[string]any{
				"Price":     big.NewInt(234),
				"Timestamp": int64(111),
			},
			map[string]any{
				"Price":     big.NewInt(456),
				"Timestamp": int64(222),
			},
		},
		consensustypes.MetadataFieldName: getMetadata(workflowID),
	}
	wrapped, err = values.NewMap(input)
	require.NoError(t, err)
	encoded, err := enc.Encode(testutils.Context(t), *wrapped)
	require.NoError(t, err)

	expected :=
		// start of the outer tuple
		getHexMetadata() +
			// start of the inner tuple (user_fields)
			"0000000000000000000000000000000000000000000000000000000000000020" + // offset of Elem list
			"0000000000000000000000000000000000000000000000000000000000000002" + // length of Elem list
			"00000000000000000000000000000000000000000000000000000000000000ea" + // Elem[0].Price = 234
			"000000000000000000000000000000000000000000000000000000000000006f" + // Elem[0].Timestamp = 111
			"00000000000000000000000000000000000000000000000000000000000001c8" + // Elem[1].Price = 456
			"00000000000000000000000000000000000000000000000000000000000000de" // Elem[1].Timestamp = 222

	require.Equal(t, expected, hex.EncodeToString(encoded))
}

func TestEVMEncoder_InvalidIDs(t *testing.T) {
	config := map[string]any{
		"abi": "bytes[] Full_reports",
	}
	wrapped, err := values.NewMap(config)
	require.NoError(t, err)
	enc, err := evm.NewEVMEncoder(wrapped)
	require.NoError(t, err)

	// output of a DF2.0 aggregator + metadata fields appended by OCR
	// using an invalid ID
	input := map[string]any{
		"Full_reports":                   []any{reportA, reportB},
		consensustypes.MetadataFieldName: getMetadata(invalidID),
	}
	wrapped, err = values.NewMap(input)
	require.NoError(t, err)
	_, err = enc.Encode(testutils.Context(t), *wrapped)
	assert.ErrorContains(t, err, "invalid byte")

	// using valid hex string of wrong length
	input = map[string]any{
		"Full_reports":                   []any{reportA, reportB},
		consensustypes.MetadataFieldName: getMetadata(wrongLength),
	}
	wrapped, err = values.NewMap(input)
	require.NoError(t, err)
	_, err = enc.Encode(testutils.Context(t), *wrapped)
	assert.ErrorContains(t, err, "incorrect length for id")
}

func getHexMetadata() string {
	return "01" + executionID + timestampHex + donID + configVersionHex + workflowID + workflowName + workflowOwnerID + reportID
}

func getMetadata(cid string) consensustypes.Metadata {
	return consensustypes.Metadata{
		Version:          1,
		ExecutionID:      executionID,
		Timestamp:        timestampInt,
		DONID:            donID,
		DONConfigVersion: configVersionInt,
		WorkflowID:       cid,
		WorkflowName:     workflowName,
		WorkflowOwner:    workflowOwnerID,
		ReportID:         reportID,
	}
}
