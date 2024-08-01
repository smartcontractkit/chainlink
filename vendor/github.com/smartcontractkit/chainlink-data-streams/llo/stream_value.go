package llo

import (
	"encoding"
	"errors"
	"fmt"
	"regexp"

	"google.golang.org/protobuf/proto"

	"github.com/shopspring/decimal"

	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
)

type StreamValue interface {
	// Binary marshaler/unmarshaler used for protobufs
	// Unmarshal should NOT panic on nil receiver, but instead return ErrNilStreamValue
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
	// TextMarshaler needed for JSON serialization and logging
	// Unmarshal should NOT panic on nil receiver, but instead return ErrNilStreamValue
	encoding.TextMarshaler
	encoding.TextUnmarshaler
	// Type is needed for proto serialization so we know how to unserialize it
	Type() LLOStreamValue_Type
}

var (
	ErrNilStreamValue = errors.New("nil stream value")
)

func UnmarshalProtoStreamValue(enc *LLOStreamValue) (sv StreamValue, err error) {
	if enc == nil {
		return nil, ErrNilStreamValue
	}
	switch enc.Type {
	case LLOStreamValue_Quote:
		sv = new(Quote)
	case LLOStreamValue_Decimal:
		sv = new(Decimal)
	default:
		return nil, fmt.Errorf("cannot unmarshal protobuf stream value; unknown StreamValueType %d", enc.Type)
	}
	if err := sv.UnmarshalBinary(enc.Value); err != nil {
		return nil, err
	}
	return sv, nil
}

func Decode(value StreamValue, data []byte) error {
	return value.UnmarshalBinary(data)
}

// Values for a set of streams, e.g. "eth-usd", "link-usd", "eur-chf" etc
// StreamIDs are uint32
type StreamValues map[llotypes.StreamID]StreamValue
type StreamAggregates map[llotypes.StreamID]map[llotypes.Aggregator]StreamValue

// Quote implements StreamValue for a {Bid, Benchmark, Ask} tuple

type Quote struct {
	Bid       decimal.Decimal
	Benchmark decimal.Decimal
	Ask       decimal.Decimal
}

var _ StreamValue = (*Quote)(nil)

func (v *Quote) MarshalBinary() (b []byte, err error) {
	if v == nil {
		return nil, ErrNilStreamValue
	}
	q := LLOStreamValueQuote{}
	q.Bid, err = v.Bid.MarshalBinary()
	if err != nil {
		return nil, err
	}
	q.Benchmark, err = v.Benchmark.MarshalBinary()
	if err != nil {
		return nil, err
	}
	q.Ask, err = v.Ask.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return proto.Marshal(&q)
}

func (v *Quote) UnmarshalBinary(data []byte) error {
	q := new(LLOStreamValueQuote)
	if err := proto.Unmarshal(data, q); err != nil {
		return err
	}
	if err := (&v.Bid).UnmarshalBinary(q.Bid); err != nil {
		return err
	}
	if err := (&v.Benchmark).UnmarshalBinary(q.Benchmark); err != nil {
		return err
	}
	return (&v.Ask).UnmarshalBinary(q.Ask)
}

func (v *Quote) MarshalText() ([]byte, error) {
	if v == nil {
		return nil, ErrNilStreamValue
	}
	return []byte(fmt.Sprintf("Q{Bid: %s, Benchmark: %s, Ask: %s}", v.Bid.String(), v.Benchmark.String(), v.Ask.String())), nil
}

var quoteRegex = regexp.MustCompile(`Q\{Bid: ([0-9.]+), Benchmark: ([0-9.]+), Ask: ([0-9.]+)\}`)

func (v *Quote) UnmarshalText(data []byte) error {
	if v == nil {
		return ErrNilStreamValue
	}

	matches := quoteRegex.FindStringSubmatch(string(data))
	if len(matches) != 4 {
		return fmt.Errorf("unexpected input for quote, expected format Q{Bid: <bid>, Benchmark: <benchmark>, Ask: <ask>}, got %s", string(data))
	}

	bid := matches[1]
	benchmark := matches[2]
	ask := matches[3]
	if err := v.Bid.UnmarshalText([]byte(bid)); err != nil {
		return err
	}
	if err := v.Benchmark.UnmarshalText([]byte(benchmark)); err != nil {
		return err
	}
	return v.Ask.UnmarshalText([]byte(ask))
}

func (v *Quote) Type() LLOStreamValue_Type {
	return LLOStreamValue_Quote
}

func (v *Quote) IsValid() bool {
	return v.Bid.Cmp(v.Benchmark) <= 0 && v.Benchmark.Cmp(v.Ask) <= 0
}

// Decimal implements StreamValue for a simple decimal value
// Use this also for integers

type Decimal decimal.Decimal

var _ StreamValue = (*Decimal)(nil)

func ToDecimal(d decimal.Decimal) *Decimal {
	return (*Decimal)(&d)
}

func (v *Decimal) Decimal() decimal.Decimal {
	return decimal.Decimal(*v)
}

func (v *Decimal) MarshalBinary() ([]byte, error) {
	if v == nil {
		return nil, ErrNilStreamValue
	}
	return decimal.Decimal(*v).MarshalBinary()
}

func (v *Decimal) UnmarshalBinary(data []byte) error {
	return (*decimal.Decimal)(v).UnmarshalBinary(data)
}

func (v *Decimal) String() string {
	return decimal.Decimal(*v).String()
}

func (v *Decimal) MarshalText() ([]byte, error) {
	if v == nil {
		return nil, ErrNilStreamValue
	}
	return []byte(v.String()), nil
}

func (v *Decimal) UnmarshalText(data []byte) error {
	if v == nil {
		return ErrNilStreamValue
	}
	return (*decimal.Decimal)(v).UnmarshalText(data)
}

func (v *Decimal) Type() LLOStreamValue_Type {
	return LLOStreamValue_Decimal
}
