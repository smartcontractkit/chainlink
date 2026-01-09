package ocr3_1

import (
	"encoding/hex"
	"strings"
	"testing"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/require"
)

func Test_hexStringTo32ByteArray_success(t *testing.T) {
	// create a valid 32-byte hex (64 hex chars)
	h := strings.Repeat("aa", 32)
	got, err := HexStringTo32ByteArray(h)
	require.NoError(t, err)
	require.NotNil(t, got)

	expectedBytes, err := hex.DecodeString(h)
	require.NoError(t, err)
	var expectedArr [32]byte
	copy(expectedArr[:], expectedBytes)
	require.Equal(t, &expectedArr, got)
}

func Test_hexStringTo32ByteArray_invalidHex(t *testing.T) {
	_, err := HexStringTo32ByteArray("zz")
	require.Error(t, err)
}

func Test_hexStringTo32ByteArray_wrongLength(t *testing.T) {
	_, err := HexStringTo32ByteArray("aa") // only 1 byte
	require.Error(t, err)
}

func Test_verifyAndExtractOCR3_1Fields_allNil(t *testing.T) {
	prevConfig, prevHist, err := VerifyAndExtractOCR3_1Fields("", 0, "")
	require.NoError(t, err)
	require.Nil(t, prevConfig)
	require.Nil(t, prevHist)
}

func Test_verifyAndExtractOCR3_1Fields_allSet_success(t *testing.T) {
	h1 := strings.Repeat("01", 32)
	h2 := strings.Repeat("ff", 32)
	cd, hd, err := VerifyAndExtractOCR3_1Fields(h1, 5, h2)
	require.NoError(t, err)
	require.NotNil(t, cd)
	require.NotNil(t, hd)

	// compare bytes
	b1, err := hex.DecodeString(h1)
	require.NoError(t, err)
	var expectedCD types.ConfigDigest
	copy(expectedCD[:], b1)
	require.Equal(t, &expectedCD, cd)

	b2, err := hex.DecodeString(h2)
	require.NoError(t, err)
	var expectedHD types.HistoryDigest
	copy(expectedHD[:], b2)
	require.Equal(t, &expectedHD, hd)
}

func Test_verifyAndExtractOCR3_1Fields_onlySeqSet_error(t *testing.T) {
	_, _, err := VerifyAndExtractOCR3_1Fields("", 1, "")
	require.Error(t, err)
	require.Contains(t, err.Error(), "must all be set or all be nil")
}

func Test_verifyAndExtractOCR3_1Fields_invalidPrevConfigHex_error(t *testing.T) {
	// seq non-zero and both strings set so validation proceeds to parse hex
	invalid := "zzzz"
	valid := strings.Repeat("aa", 32)
	_, _, err := VerifyAndExtractOCR3_1Fields(invalid, 1, valid)
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid OracleConfig.PrevConfigDigest")
}

func Test_verifyAndExtractOCR3_1Fields_invalidPrevHistoryHex_error(t *testing.T) {
	valid := strings.Repeat("aa", 32)
	invalid := "zzzz"
	_, _, err := VerifyAndExtractOCR3_1Fields(valid, 1, invalid)
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid OracleConfig.PrevHistoryDigest")
}
