package mercury

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
	FeedID               string `json:"feedId"`
	FullReport           []byte `json:"fullReport"`
	BenchmarkPrice       int64  `json:"benchmarkPrice"`
	ObservationTimestamp int64  `json:"observationTimestamp"`
}

type Codec struct {
}

func (c Codec) Unwrap(raw values.Value) ([]FeedReport, error) {
	dest := []FeedReport{}
	err := raw.UnwrapTo(&dest)
	// TODO: validate reports
	return dest, err
}

func (c Codec) Wrap(reports []FeedReport) (values.Value, error) {
	// TODO: validate reports
	return values.Wrap(reports)
}

func NewCodec() Codec {
	return Codec{}
}
