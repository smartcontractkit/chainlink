package mercury

import (
	"encoding/hex"
	"errors"
	"strings"

	"github.com/shopspring/decimal"

	"github.com/smartcontractkit/chainlink-common/pkg/values"
)

// hex-encoded 32-byte value, prefixed with "0x", all lowercase
type FeedID string

const FeedIDBytesLen = 32

var ErrInvalidFeedID = errors.New("invalid feed ID")

func (id FeedID) String() string {
	return string(id)
}

func (id FeedID) Validate() error {
	if len(id) != 2*FeedIDBytesLen+2 {
		return ErrInvalidFeedID
	}
	if !strings.HasPrefix(string(id), "0x") {
		return ErrInvalidFeedID
	}
	if strings.ToLower(string(id)) != string(id) {
		return ErrInvalidFeedID
	}
	_, err := hex.DecodeString(string(id)[2:])
	return err
}

type ReportSet struct {
	// feedID -> report
	Reports map[FeedID]Report
}

type Report struct {
	Info       ReportInfo // minimal data extracted from the report for convenience
	FullReport []byte     // full report, acceptable by the verifier contract
}

type ReportInfo struct {
	Timestamp uint32
	Price     float64
}

// TODO implement an actual codec
type Codec struct {
}

func (m Codec) Unwrap(raw values.Value) (ReportSet, error) {
	return ReportSet{
		Reports: map[FeedID]Report{
			FeedID("012345678901234567890123456789012345678901234567890123456789000000"): {
				Info: ReportInfo{
					Timestamp: 42,
					Price:     100.00,
				},
			},
		},
	}, nil
}

func (m Codec) Wrap(reportSet ReportSet) (values.Value, error) {
	return values.NewMap(
		map[string]any{
			"0123456789": map[string]any{
				"timestamp": 42,
				"price":     decimal.NewFromFloat(100.00),
			},
		},
	)
}

func NewCodec() Codec {
	return Codec{}
}
