package llo

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
)

// Chose uint32 to represent StreamID and ChannelID for the following reasons:
// - 4 billion is more than enough to cover our conceivable channel/stream requirements
// - It is the most compatible, supported everywhere, can be serialized into JSON and used in Javascript without problems
// - It is the smallest reasonable data type that balances between a large set of possible identifiers and not using too much space
// - If randomly chosen, low chance of off-by-one ids being valid
// - Is not specific to any chain, e.g. [32]byte is not fully supported on starknet etc
// - Avoids any possible encoding/copypasta issues e.g. UUIDs which can convert to [32]byte in multiple different ways
//
// It is recommended not to use 0 as a stream ID to avoid confusion with
// uninitialized values
type StreamID = uint32

type LifeCycleStage string

// ReportFormat represents different formats for different targets e.g. EVM,
// Solana, JSON, kalechain etc
type ReportFormat uint32

const (
	_ ReportFormat = 0 // reserved

	// NOTE: Only add something here if you actually need it, because it has to
	// be supported forever and can't be changed
	ReportFormatEVMPremiumLegacy ReportFormat = 1
	ReportFormatJSON             ReportFormat = 2

	_ ReportFormat = math.MaxUint32 // reserved
)

var ReportFormats = []ReportFormat{
	ReportFormatEVMPremiumLegacy,
	ReportFormatJSON,
}

func (rf ReportFormat) String() string {
	switch rf {
	case ReportFormatEVMPremiumLegacy:
		return "evm_premium_legacy"
	case ReportFormatJSON:
		return "json"
	default:
		return fmt.Sprintf("unknown(%d)", rf)
	}
}

func ReportFormatFromString(s string) (ReportFormat, error) {
	switch s {
	case "evm_premium_legacy":
		return ReportFormatEVMPremiumLegacy, nil
	case "json":
		return ReportFormatJSON, nil
	default:
		return 0, fmt.Errorf("unknown report format: %q", s)
	}
}

func (rf ReportFormat) MarshalText() ([]byte, error) {
	return []byte(rf.String()), nil
}

func (rf *ReportFormat) UnmarshalText(text []byte) error {
	val, err := ReportFormatFromString(string(text))
	if err != nil {
		return err
	}
	*rf = val
	return nil
}

func (rf ReportFormat) MarshalJSON() ([]byte, error) {
	return json.Marshal(rf.String())
}

func (rf *ReportFormat) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		val, err := ReportFormatFromString(str)
		if err != nil {
			return err
		}
		*rf = val
		return nil
	}

	var num uint32
	if err := json.Unmarshal(data, &num); err == nil {
		*rf = ReportFormat(num)
		return nil
	}

	return fmt.Errorf("invalid JSON value for ReportFormat, expected number or string: %s", data)
}

type Aggregator uint32

const (
	_ Aggregator = 0 // reserved

	// NOTE: Only add something here if you actually need it, because it has to
	// be supported forever and can't be changed
	AggregatorMedian = 1
	AggregatorMode   = 2
	// AggregatorQuote is a special aggregator that is used to aggregate
	// a list of Bid/Mid/Ask price tuples
	AggregatorQuote = 3

	_ Aggregator = math.MaxUint32 // reserved
)

func (a Aggregator) String() string {
	switch a {
	case AggregatorMedian:
		return "median"
	case AggregatorMode:
		return "mode"
	case AggregatorQuote:
		return "quote"
	default:
		return fmt.Sprintf("unknown(%d)", a)
	}
}

func (a Aggregator) MarshalText() ([]byte, error) {
	return []byte(a.String()), nil
}

func AggregatorFromString(s string) (Aggregator, error) {
	switch s {
	case "median":
		return AggregatorMedian, nil
	case "mode":
		return AggregatorMode, nil
	case "quote":
		return AggregatorQuote, nil
	default:
		return 0, fmt.Errorf("unknown aggregator: %q", s)
	}
}

func (a *Aggregator) UnmarshalText(text []byte) error {
	val, err := AggregatorFromString(string(text))
	if err != nil {
		return err
	}
	*a = val
	return nil
}

func (a Aggregator) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.String())
}

func (a *Aggregator) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		val, err := AggregatorFromString(str)
		if err != nil {
			return err
		}
		*a = val
		return nil
	}

	var num uint32
	if err := json.Unmarshal(data, &num); err == nil {
		*a = Aggregator(num)
		return nil
	}

	return fmt.Errorf("invalid JSON value for Aggregator, expected number or string, got: '%s'", data)
}

type ReportInfo struct {
	LifeCycleStage LifeCycleStage
	ReportFormat   ReportFormat
}

type Transmitter ocr3types.ContractTransmitter[ReportInfo]

type Stream struct {
	// ID is the ID of the stream to be observed
	StreamID StreamID `json:"streamId"`
	// Aggregator is the method used by consensus protocol to aggregate
	// multiple stream observations e.g. "median", "mode" or other more exotic
	// methods
	Aggregator Aggregator `json:"aggregator"`
}

type ChannelDefinition struct {
	// ReportFormat controls the output format of the report. It might be a
	// different chain target, different report format etc. Custom logic can be
	// implemented for each report format.
	ReportFormat ReportFormat `json:"reportFormat"`
	// Streams is the list of streams to be observed and aggregated
	// by the protocol.
	Streams []Stream `json:"streams"`
	// Opts contains configuration data for use in report generation
	// for this channel, e.g. feed ID, expiry window, USD base fee etc
	//
	// How this is encoded is up to the implementer but should be consistent
	// for a given ReportFormat
	//
	// May be nil
	Opts ChannelOpts `json:"opts"`
}

func (a ChannelDefinition) Equals(b ChannelDefinition) bool {
	if a.ReportFormat != b.ReportFormat {
		return false
	}
	if len(a.Streams) != len(b.Streams) {
		return false
	}
	for i, strm := range a.Streams {
		if strm != b.Streams[i] {
			return false
		}
	}
	return bytes.Equal(a.Opts, b.Opts)
}

type ChannelOpts []byte

// UnmarshalJSON allows taking an actual JSON object as ChannelOpts and passes it through
// in canonicalized form.
func (m *ChannelOpts) UnmarshalJSON(data []byte) (err error) {
	if string(data) == "null" {
		return nil
	}
	*m, err = formatJSON(data)
	return err
}

// MarshalJSON passes raw json directly through in a canonicalized form.
func (m ChannelOpts) MarshalJSON() ([]byte, error) {
	if m == nil {
		return []byte("null"), nil
	}
	return formatJSON(m)
}

// formatJSON ensures that the JSON string is in a deterministic shape
func formatJSON(input []byte) ([]byte, error) {
	var obj map[string]interface{}

	// Unmarshal the JSON string into a map
	if err := json.Unmarshal(input, &obj); err != nil {
		return nil, err
	}

	// Marshal the map back to a JSON string with sorted keys
	formatted, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	return formatted, nil
}

type ChannelDefinitions map[ChannelID]ChannelDefinition

func (c *ChannelDefinitions) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan Data: value is not []byte")
	}
	if bytes == nil {
		*c = nil
		return nil
	}
	if len(bytes) == 0 {
		*c = ChannelDefinitions{}
		return nil
	}

	return json.Unmarshal(bytes, c)
}

func (c ChannelDefinitions) Value() (driver.Value, error) {
	return json.Marshal(c)
}

// It is recommended not to use 0 as a channel ID to avoid confusion with
// uninitialized values
type ChannelID = uint32

type ChannelDefinitionCache interface {
	Definitions() ChannelDefinitions
	services.Service
}
