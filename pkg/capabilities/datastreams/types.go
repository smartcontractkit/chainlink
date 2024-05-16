package datastreams

import (
	"encoding/hex"
	"errors"
	"strings"

	"github.com/smartcontractkit/chainlink-common/pkg/values"
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

type FeedReport struct {
	FeedID     string
	FullReport []byte
	Rs         [][]byte
	Ss         [][]byte
	Vs         []byte

	// Fields below are derived from FullReport
	// NOTE: BenchmarkPrice is a byte representation of big.Int. We can't use big.Int
	// directly due to Value serialization problems using mapstructure.
	BenchmarkPrice       []byte
	ObservationTimestamp int64
}

//go:generate mockery --quiet --name ReportCodec --output ./mocks/ --case=underscore
type ReportCodec interface {
	// validate each report and convert to a list of Mercury reports
	Unwrap(raw values.Value) ([]FeedReport, error)

	// validate each report and convert to Value
	Wrap(reports []FeedReport) (values.Value, error)
}
