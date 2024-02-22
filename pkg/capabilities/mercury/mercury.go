package mercury

import (
	"encoding/hex"
	"errors"
	"strings"
	"time"

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
	now := uint32(time.Now().Unix())
	return ReportSet{
		Reports: map[FeedID]Report{
			FeedID("0x1111111111111111111100000000000000000000000000000000000000000000"): {
				Info: ReportInfo{
					Timestamp: now,
					Price:     100.00,
				},
			},
			FeedID("0x2222222222222222222200000000000000000000000000000000000000000000"): {
				Info: ReportInfo{
					Timestamp: now,
					Price:     100.00,
				},
			},
			FeedID("0x3333333333333333333300000000000000000000000000000000000000000000"): {
				Info: ReportInfo{
					Timestamp: now,
					Price:     100.00,
				},
			},
		},
	}, nil
}

func (m Codec) Wrap(reportSet ReportSet) (values.Value, error) {
	return values.NewMap(
		map[string]any{
			"0x1111111111111111111100000000000000000000000000000000000000000000": map[string]any{
				"timestamp": 42,
				"price":     decimal.NewFromFloat(100.00),
			},
		},
	)
}

func NewCodec() Codec {
	return Codec{}
}
