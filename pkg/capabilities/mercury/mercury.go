package mercury

import (
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/shopspring/decimal"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
)

var (
	feedOne   = FeedID("0x1111111111111111111100000000000000000000000000000000000000000000")
	feedTwo   = FeedID("0x2222222222222222222200000000000000000000000000000000000000000000")
	feedThree = FeedID("0x3333333333333333333300000000000000000000000000000000000000000000")
)

// hex-encoded 32-byte value, prefixed with "0x", all lowercase
type FeedID string

const FeedIDBytesLen = 32

var ErrInvalidFeedID = errors.New("invalid feed ID")

func (id FeedID) String() string {
	return string(id)
}

// Bytes() converts the FeedID string into a [32]byte
// value.
// Note: this function panics if the underlying
// string isn't of the right length. For production (i.e.)
// non-test uses, please create the FeedID via the NewFeedID
// constructor, which will validate the string.
func (id FeedID) Bytes() [FeedIDBytesLen]byte {
	b, _ := hex.DecodeString(string(id)[2:])
	return [FeedIDBytesLen]byte(b)
}

func (id FeedID) validate() error {
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

func NewFeedID(s string) (FeedID, error) {
	id := FeedID(s)
	return id, id.validate()
}

func FeedIDFromBytes(b [FeedIDBytesLen]byte) FeedID {
	return FeedID("0x" + hex.EncodeToString(b[:]))
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

type FeedReport struct {
	FeedID               string `json:"feedId"`
	FullReport           []byte `json:"fullReport"`
	BenchmarkPrice       int64  `json:"benchmarkPrice"`
	ObservationTimestamp int64  `json:"observationTimestamp"`
}

// TODO implement an actual codec
type Codec struct {
}

func (m Codec) Unwrap(raw values.Value) (ReportSet, error) {
	now := uint32(time.Now().Unix())
	return ReportSet{
		Reports: map[FeedID]Report{
			feedOne: {
				Info: ReportInfo{
					Timestamp: now,
					Price:     100.00,
				},
			},
			feedTwo: {
				Info: ReportInfo{
					Timestamp: now,
					Price:     100.00,
				},
			},
			feedThree: {
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
			feedOne.String(): map[string]any{
				"timestamp": 42,
				"price":     decimal.NewFromFloat(100.00),
			},
		},
	)
}

func (m Codec) WrapMercuryTriggerEvent(event capabilities.TriggerEvent) (values.Value, error) {
	return values.Wrap(event)
}

func (m Codec) UnwrapMercuryTriggerEvent(raw values.Value) (capabilities.TriggerEvent, error) {
	mercuryTriggerEvent := capabilities.TriggerEvent{}
	val, err := raw.Unwrap()
	if err != nil {
		return mercuryTriggerEvent, err
	}
	event := val.(map[string]any)
	mercuryTriggerEvent.TriggerType = event["TriggerType"].(string)
	mercuryTriggerEvent.ID = event["ID"].(string)
	mercuryTriggerEvent.Timestamp = event["Timestamp"].(string)
	mercuryTriggerEvent.BatchedPayload = make(map[string]any)
	for id, report := range event["BatchedPayload"].(map[string]any) {
		reportMap := report.(map[string]any)
		var mercuryReport FeedReport
		err = mapstructure.Decode(reportMap, &mercuryReport)
		if err != nil {
			return mercuryTriggerEvent, err
		}

		mercuryTriggerEvent.BatchedPayload[id] = mercuryReport
	}
	return mercuryTriggerEvent, nil
}

func NewCodec() Codec {
	return Codec{}
}
