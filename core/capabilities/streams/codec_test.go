package streams_test

import (
	"crypto/rand"
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	ocrTypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/chains/evmutil"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/datastreams"
	v3 "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v3"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/streams"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/v3/reportcodec"
)

func TestCodec_WrapUnwrap(t *testing.T) {
	codec := streams.NewCodec(logger.TestLogger(t))

	id1, id1Str := newFeedID(t)
	id2, id2Str := newFeedID(t)
	price1, price2 := big.NewInt(1), big.NewInt(2)
	timestamp1, timestamp2 := int64(100), int64(200)
	report1, report2 := newReport(t, id1, price1, timestamp1), newReport(t, id2, price2, timestamp2)
	reportCtx := ocrTypes.ReportContext{}
	rawCtx := rawReportContext(reportCtx)

	keyBundle1, err := ocr2key.New(chaintype.EVM)
	require.NoError(t, err)
	keyBundle2, err := ocr2key.New(chaintype.EVM)
	require.NoError(t, err)

	signatureK1R1, err := keyBundle1.Sign(reportCtx, report1)
	require.NoError(t, err)
	signatureK1R2, err := keyBundle1.Sign(reportCtx, report2)
	require.NoError(t, err)
	signatureK2R1, err := keyBundle2.Sign(reportCtx, report1)
	require.NoError(t, err)
	signatureK2R2, err := keyBundle2.Sign(reportCtx, report2)
	require.NoError(t, err)

	allowedSigners := [][]byte{keyBundle1.PublicKey(), keyBundle2.PublicKey()} // bad name - see comment on evmKeyring.PublicKey

	wrapped, err := codec.Wrap([]datastreams.FeedReport{
		{
			FeedID:        id1Str,
			FullReport:    report1,
			ReportContext: rawCtx,
			Signatures:    [][]byte{signatureK1R1, signatureK2R1},
		},
		{
			FeedID:        id2Str,
			FullReport:    report2,
			ReportContext: rawCtx,
			Signatures:    [][]byte{signatureK1R2, signatureK2R2},
		},
	})
	require.NoError(t, err)

	// wrong type
	_, err = codec.Unwrap(values.NewBool(true))
	require.Error(t, err)

	// correct reports but wrong signatures
	unwrapped, err := codec.Unwrap(wrapped)
	require.NoError(t, err)
	require.Equal(t, 2, len(unwrapped))
	require.Equal(t, price1.Bytes(), unwrapped[0].BenchmarkPrice)
	require.Equal(t, price2.Bytes(), unwrapped[1].BenchmarkPrice)
	require.Equal(t, timestamp1, unwrapped[0].ObservationTimestamp)
	require.Equal(t, timestamp2, unwrapped[1].ObservationTimestamp)
	for _, report := range unwrapped {
		require.Error(t, codec.Validate(report, nil, 1))
	}

	// valid signatures
	for _, report := range unwrapped {
		require.NoError(t, codec.Validate(report, allowedSigners, 2))
	}

	// invalid FeedID
	wrappedInvalid, err := codec.Wrap([]datastreams.FeedReport{
		{
			FeedID:        id2Str, // ID #2 doesn't match what's in report #1
			FullReport:    report1,
			ReportContext: rawCtx,
			Signatures:    [][]byte{signatureK1R1, signatureK2R1},
		},
	})
	require.NoError(t, err)
	_, err = codec.Unwrap(wrappedInvalid)
	require.Error(t, err)
	require.Contains(t, err.Error(), "feed ID mismatch")
}

func newFeedID(t *testing.T) ([32]byte, string) {
	buf := [32]byte{}
	_, err := rand.Read(buf[:])
	require.NoError(t, err)
	return buf, "0x" + hex.EncodeToString(buf[:])
}

func newReport(t *testing.T, feedID [32]byte, price *big.Int, timestamp int64) []byte {
	v3Codec := reportcodec.NewReportCodec(feedID, logger.TestLogger(t))
	raw, err := v3Codec.BuildReport(v3.ReportFields{
		BenchmarkPrice:     price,
		Timestamp:          uint32(timestamp),
		ValidFromTimestamp: uint32(timestamp),
		Bid:                price,
		Ask:                price,
		LinkFee:            price,
		NativeFee:          price,
		ExpiresAt:          uint32(timestamp + 1000000),
	})
	require.NoError(t, err)
	return raw
}

func rawReportContext(reportCtx ocrTypes.ReportContext) []byte {
	rc := evmutil.RawReportContext(reportCtx)
	flat := []byte{}
	for _, r := range rc {
		flat = append(flat, r[:]...)
	}
	return flat
}
